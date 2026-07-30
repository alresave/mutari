create or replace function public.get_teams()
returns table(id bigint, name text, description text, company_id bigint)
language sql stable security definer set search_path = public as $$
  select t.id, t.name, t.description, t.company_id
  from public.teams t join public.profiles p on p.company_id = t.company_id
  where p.id = auth.uid() order by t.name;
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
  ), company as (select company_id from public.teams where id = target_team)
  select up.user_id, p.name, p.photo_url,
    round((1 - (up.effective_vector <=> tv.value)) * 100)::integer
  from public.user_mcp_profiles up join public.profiles p on p.id = up.user_id
  cross join team_vector tv
  where p.company_id = (select company_id from company)
    and not exists (select 1 from public.team_members tm where tm.team_id = target_team and tm.user_id = up.user_id)
  order by up.effective_vector <=> tv.value limit greatest(1, least(match_limit, 100));
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
  from values_by_dimension v join public.mcp_blocks b on b.id = v.block;
$$;

create or replace function public.get_team_gaps(target_team bigint)
returns jsonb language sql stable security definer set search_path = public, extensions as $$
  with culture as (select public.get_team_culture(target_team)->'block_scores' as scores),
  gaps as (select key from culture, jsonb_each_text(scores) where value::numeric < 20)
  select jsonb_build_object('gaps', coalesce(jsonb_agg(key), '[]'::jsonb), 'recommendation', 'Prioritize candidates who strengthen the lowest-scoring dimensions.') from gaps;
$$;

grant execute on function public.get_teams() to authenticated;
grant execute on function public.get_team_matches(bigint, integer) to authenticated;
grant execute on function public.get_team_culture(bigint) to authenticated;
grant execute on function public.get_team_gaps(bigint) to authenticated;
