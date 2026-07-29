-- Team analytics are exposed through SECURITY DEFINER RPCs. Keep the
-- authorization check in the database so it cannot be bypassed by calling an
-- RPC directly instead of going through an Edge Function.
alter table public.companies enable row level security;
alter table public.team_members enable row level security;
alter table public.mcp_blocks enable row level security;

create policy "blocks are readable" on public.mcp_blocks
  for select to anon, authenticated using (true);

-- A user may edit their public details, but may never promote themselves or
-- attach themselves to another company through the browser API.
create or replace function public.protect_profile_privileges()
returns trigger
language plpgsql
security definer
set search_path = public
as $$
begin
  if (new.role is distinct from old.role or new.company_id is distinct from old.company_id)
     and not public.is_admin() then
    raise exception 'role and company are managed by an administrator';
  end if;
  return new;
end;
$$;

drop trigger if exists protect_profile_privileges on public.profiles;
create trigger protect_profile_privileges
before update on public.profiles
for each row execute procedure public.protect_profile_privileges();

create or replace function public.can_access_team(target_team bigint)
returns boolean
language sql
stable
security definer
set search_path = public
as $$
  select exists (
    select 1
    from public.teams t
    join public.profiles caller on caller.id = auth.uid()
    where t.id = target_team
      and t.company_id = caller.company_id
  );
$$;

create or replace function public.get_team_matches(target_team bigint, match_limit integer default 50)
returns table(user_id uuid, name text, photo text, score integer)
language sql stable security definer set search_path = public, extensions as $$
  with dimension_values as (
    select idx, avg(value::double precision) as value
    from public.team_members tm join public.user_mcp_profiles up on up.user_id = tm.user_id,
         unnest(string_to_array(trim(both '[]' from up.effective_vector::text), ',')::double precision[]) with ordinality x(value, idx)
    where tm.team_id = target_team group by idx
  ), team_vector as (
    select ('[' || string_agg(value::text, ',' order by idx) || ']')::extensions.vector as value from dimension_values
  )
  select up.user_id, p.name, p.photo_url,
    round((1 - (up.effective_vector <=> tv.value)) * 100)::integer
  from public.user_mcp_profiles up
  join public.profiles p on p.id = up.user_id
  cross join team_vector tv
  where public.can_access_team(target_team)
    and p.company_id = (select company_id from public.teams where id = target_team)
    and not exists (select 1 from public.team_members tm where tm.team_id = target_team and tm.user_id = up.user_id)
  order by up.effective_vector <=> tv.value limit greatest(1, least(match_limit, 100));
$$;

create or replace function public.get_team_fit(target_team bigint, candidate uuid)
returns jsonb language sql stable security definer set search_path = public, extensions as $$
  with dimension_values as (
    select idx, avg(value::double precision) as value
    from public.team_members tm join public.user_mcp_profiles up on up.user_id = tm.user_id,
         unnest(string_to_array(trim(both '[]' from up.effective_vector::text), ',')::double precision[]) with ordinality x(value, idx)
    where tm.team_id = target_team group by idx
  ), team_vector as (
    select ('[' || string_agg(value::text, ',' order by idx) || ']')::extensions.vector as value from dimension_values
  )
  select jsonb_build_object(
    'user_id', p.id, 'name', p.name, 'photo', p.photo_url,
    'score', round((1 - (up.effective_vector <=> tv.value)) * 100)::integer,
    'suggested_role', case when 1 - (up.effective_vector <=> tv.value) >= .8 then 'Strong fit' else 'Complementary contributor' end
  )
  from public.user_mcp_profiles up join public.profiles p on p.id = up.user_id cross join team_vector tv
  where public.can_access_team(target_team)
    and up.user_id = candidate
    and p.company_id = (select company_id from public.teams where id = target_team);
$$;

create or replace function public.get_team_culture(target_team bigint)
returns jsonb language sql stable security definer set search_path = public, extensions as $$
  with values_by_dimension as (
    select floor((idx - 1) / 8)::integer + 1 as block, avg(value::double precision) as score
    from public.team_members tm join public.user_mcp_profiles p on p.user_id = tm.user_id,
         unnest(string_to_array(trim(both '[]' from p.effective_vector::text), ',')::double precision[]) with ordinality x(value, idx)
    where tm.team_id = target_team group by block
  )
  select jsonb_build_object('block_scores', coalesce(jsonb_object_agg(b.name, round(v.score * 100)), '{}'::jsonb))
  from values_by_dimension v join public.mcp_blocks b on b.id = v.block
  where public.can_access_team(target_team);
$$;

create or replace function public.get_team_gaps(target_team bigint)
returns jsonb language sql stable security definer set search_path = public, extensions as $$
  with culture as (select public.get_team_culture(target_team)->'block_scores' as scores),
  gaps as (select key from culture, jsonb_each_text(scores) where value::numeric < 20)
  select jsonb_build_object('gaps', coalesce(jsonb_agg(key), '[]'::jsonb), 'recommendation', 'Prioritize candidates who strengthen the lowest-scoring dimensions.')
  from gaps
  where public.can_access_team(target_team);
$$;

create or replace function public.get_team_hiring_profile(target_team bigint)
returns jsonb language sql stable security definer set search_path = public as $$
  with gap_data as (select public.get_team_gaps(target_team) as value),
  gaps as (select coalesce(array_agg(item), array[]::text[]) as items from gap_data, jsonb_array_elements_text(value->'gaps') item)
  select jsonb_build_object(
    'suggested_role', case when cardinality(items) >= 3 then 'Versatile team builder' else 'Specialist contributor' end,
    'recommended_traits', to_jsonb(coalesce(items, array['Collaboration Depth', 'Execution Drive']::text[])),
    'recommendation', 'Prioritize candidates who complement the team''s lowest-scoring dimensions.'
  ) from gaps where public.can_access_team(target_team);
$$;

revoke execute on function public.can_access_team(bigint) from public;
grant execute on function public.can_access_team(bigint) to authenticated;
