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
  where up.user_id = candidate;
$$;

create or replace function public.get_team_hiring_profile(target_team bigint)
returns jsonb language sql stable security definer set search_path = public as $$
  with gap_data as (select public.get_team_gaps(target_team) as value),
  gaps as (select coalesce(array_agg(item), array[]::text[]) as items from gap_data, jsonb_array_elements_text(value->'gaps') item)
  select jsonb_build_object(
    'suggested_role', case when cardinality(items) >= 3 then 'Versatile team builder' else 'Specialist contributor' end,
    'recommended_traits', to_jsonb(coalesce(items, array['Collaboration Depth', 'Execution Drive']::text[])),
    'recommendation', 'Prioritize candidates who complement the team''s lowest-scoring dimensions.'
  ) from gaps;
$$;

grant execute on function public.get_team_fit(bigint, uuid) to authenticated;
grant execute on function public.get_team_hiring_profile(bigint) to authenticated;
