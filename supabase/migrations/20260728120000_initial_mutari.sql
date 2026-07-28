create extension if not exists vector with schema extensions;

create table public.companies (
  id bigint generated always as identity primary key,
  name text not null,
  website text not null default '',
  created_at timestamptz not null default now()
);

create table public.profiles (
  id uuid primary key references auth.users(id) on delete cascade,
  name text not null default '',
  photo_url text not null default '',
  company_id bigint references public.companies(id) on delete set null,
  created_at timestamptz not null default now()
);

create table public.teams (
  id bigint generated always as identity primary key,
  name text not null,
  description text,
  company_id bigint not null references public.companies(id) on delete cascade
);

create table public.team_members (
  team_id bigint not null references public.teams(id) on delete cascade,
  user_id uuid not null references public.profiles(id) on delete cascade,
  role text not null default 'member',
  primary key (team_id, user_id)
);

create table public.mcp_blocks (
  id smallint primary key,
  name text not null unique,
  description text not null
);

create table public.mcp_questions (
  id bigint generated always as identity primary key,
  block_id smallint not null references public.mcp_blocks(id) on delete cascade,
  dimension_index smallint not null check (dimension_index between 0 and 63),
  question_text text not null,
  scale text not null default '0-1',
  unique (dimension_index)
);

create table public.user_mcp_answers (
  user_id uuid not null references public.profiles(id) on delete cascade,
  question_id bigint not null references public.mcp_questions(id) on delete cascade,
  value double precision not null check (value between 0 and 1),
  updated_at timestamptz not null default now(),
  primary key (user_id, question_id)
);

create table public.user_mcp_profiles (
  user_id uuid primary key references public.profiles(id) on delete cascade,
  effective_vector extensions.vector(64) not null,
  confidence_score double precision not null default 1 check (confidence_score between 0 and 1),
  updated_at timestamptz not null default now()
);
create index user_mcp_profiles_vector_idx on public.user_mcp_profiles using hnsw (effective_vector extensions.vector_cosine_ops);

create or replace function public.handle_new_user()
returns trigger language plpgsql security definer set search_path = public as $$
begin
  insert into public.profiles (id, name)
  values (new.id, coalesce(new.raw_user_meta_data ->> 'name', split_part(new.email, '@', 1)));
  return new;
end;
$$;
create trigger on_auth_user_created after insert on auth.users for each row execute procedure public.handle_new_user();

create or replace function public.rebuild_mcp_profile(target_user uuid default auth.uid())
returns void language sql security definer set search_path = public, extensions as $$
  with dimensions as (
    select q.dimension_index, coalesce(avg(a.value), 0.5) - 0.5 as value
    from public.mcp_questions q
    left join public.user_mcp_answers a on a.question_id = q.id and a.user_id = target_user
    group by q.dimension_index
  ), normalized as (
    select dimension_index, value / nullif(sqrt(sum(value * value) over ()), 0) as value
    from dimensions
  ), vector_value as (
    select ('[' || string_agg(coalesce(value, 0)::text, ',' order by dimension_index) || ']')::extensions.vector as value
    from normalized
  )
  insert into public.user_mcp_profiles (user_id, effective_vector, updated_at)
  select target_user, value, now() from vector_value
  on conflict (user_id) do update set effective_vector = excluded.effective_vector, updated_at = now();
$$;

create or replace function public.submit_mcp_answers(answers jsonb)
returns void language plpgsql security definer set search_path = public as $$
begin
  if auth.uid() is null then raise exception 'authentication required'; end if;
  insert into public.user_mcp_answers (user_id, question_id, value, updated_at)
  select auth.uid(), item.question_id, item.value, now()
  from jsonb_to_recordset(answers) as item(question_id bigint, value double precision)
  on conflict (user_id, question_id) do update set value = excluded.value, updated_at = now();
  perform public.rebuild_mcp_profile(auth.uid());
end;
$$;

create or replace function public.get_mcp_matches(match_limit integer default 50)
returns table(user_id uuid, name text, photo_url text, score integer)
language sql stable security definer set search_path = public, extensions as $$
  select candidate.user_id, profile.name, profile.photo_url,
    round((1 - (candidate.effective_vector <=> source.effective_vector)) * 100)::integer as score
  from public.user_mcp_profiles source
  join public.user_mcp_profiles candidate on candidate.user_id <> source.user_id
  join public.profiles profile on profile.id = candidate.user_id
  where source.user_id = auth.uid()
    and profile.company_id is not distinct from (select company_id from public.profiles where id = auth.uid())
  order by candidate.effective_vector <=> source.effective_vector
  limit greatest(1, least(match_limit, 100));
$$;

alter table public.profiles enable row level security;
alter table public.mcp_questions enable row level security;
alter table public.user_mcp_answers enable row level security;
alter table public.user_mcp_profiles enable row level security;
alter table public.teams enable row level security;
alter table public.team_members enable row level security;

create policy "read own profile" on public.profiles for select to authenticated using (id = auth.uid());
create policy "update own profile" on public.profiles for update to authenticated using (id = auth.uid()) with check (id = auth.uid());
create policy "questions are readable" on public.mcp_questions for select to anon, authenticated using (true);
create policy "answers belong to caller" on public.user_mcp_answers for all to authenticated using (user_id = auth.uid()) with check (user_id = auth.uid());
create policy "profile vectors belong to caller" on public.user_mcp_profiles for select to authenticated using (user_id = auth.uid());

grant execute on function public.submit_mcp_answers(jsonb) to authenticated;
grant execute on function public.get_mcp_matches(integer) to authenticated;
