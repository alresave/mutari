create table public.company_invites (
  id uuid primary key default gen_random_uuid(),
  company_id bigint not null references public.companies(id) on delete cascade,
  email text,
  created_by uuid not null references public.profiles(id) on delete cascade,
  expires_at timestamptz not null default now() + interval '7 days',
  accepted_at timestamptz,
  accepted_by uuid references public.profiles(id) on delete set null,
  created_at timestamptz not null default now(),
  check (email is null or email = lower(email))
);

alter table public.company_invites enable row level security;

create or replace function public.protect_profile_privileges()
returns trigger
language plpgsql
security definer
set search_path = public
as $$
begin
  if (new.role is distinct from old.role or new.company_id is distinct from old.company_id)
     and not public.is_admin()
     and current_setting('app.allow_company_join', true) is distinct from 'true' then
    raise exception 'role and company are managed by an administrator';
  end if;
  return new;
end;
$$;

create or replace function public.admin_create_company_invite(invitee_email text default null, expiry_hours integer default 168)
returns jsonb
language plpgsql
security definer
set search_path = public
as $$
declare
  invite public.company_invites;
  current_company bigint;
begin
  if not public.is_admin() then raise exception 'admin role required'; end if;
  select company_id into current_company from public.profiles where id = auth.uid();
  if current_company is null then raise exception 'create or join a company first'; end if;
  if expiry_hours not between 1 and 720 then raise exception 'expiry_hours must be between 1 and 720'; end if;

  insert into public.company_invites(company_id, email, created_by, expires_at)
  values (current_company, nullif(lower(trim(invitee_email)), ''), auth.uid(), now() + make_interval(hours => expiry_hours))
  returning * into invite;

  return jsonb_build_object('invite_code', invite.id, 'expires_at', invite.expires_at, 'email', invite.email);
end;
$$;

create or replace function public.accept_company_invite(invite_code uuid)
returns jsonb
language plpgsql
security definer
set search_path = public
as $$
declare
  invite public.company_invites;
  caller_email text;
begin
  if auth.uid() is null then raise exception 'authentication required'; end if;
  select * into invite from public.company_invites where id = invite_code for update;
  if not found then raise exception 'invite not found'; end if;
  if invite.accepted_at is not null then raise exception 'invite already accepted'; end if;
  if invite.expires_at <= now() then raise exception 'invite expired'; end if;

  caller_email := lower(coalesce(auth.jwt() ->> 'email', ''));
  if invite.email is not null and invite.email <> caller_email then
    raise exception 'invite was issued for a different email address';
  end if;

  perform set_config('app.allow_company_join', 'true', true);
  update public.profiles set company_id = invite.company_id where id = auth.uid();
  update public.company_invites set accepted_at = now(), accepted_by = auth.uid() where id = invite.id;
  return jsonb_build_object('company_id', invite.company_id, 'status', 'joined');
end;
$$;

create or replace function public.admin_get_company_candidates()
returns table(user_id uuid, name text, photo_url text, profile_complete boolean, team_count bigint)
language sql
stable
security definer
set search_path = public
as $$
  with candidates as (
    select p.id as user_id, p.name, p.photo_url,
      exists(select 1 from public.user_mcp_profiles up where up.user_id = p.id) as profile_complete,
      count(tm.team_id) as team_count
    from public.profiles p
    left join public.team_members tm on tm.user_id = p.id
    where public.is_admin()
      and p.company_id = (select company_id from public.profiles where id = auth.uid())
    group by p.id, p.name, p.photo_url
  )
  select user_id, name, photo_url, profile_complete, team_count
  from candidates
  order by profile_complete desc, name;
$$;

create or replace function public.admin_get_team_members(target_team bigint)
returns table(user_id uuid, name text, photo_url text, member_role text, profile_complete boolean)
language sql
stable
security definer
set search_path = public
as $$
  select p.id, p.name, p.photo_url, tm.role,
    exists(select 1 from public.user_mcp_profiles up where up.user_id = p.id)
  from public.team_members tm
  join public.profiles p on p.id = tm.user_id
  where public.is_admin()
    and tm.team_id = target_team
    and exists (
      select 1 from public.teams t
      join public.profiles caller on caller.id = auth.uid()
      where t.id = target_team and t.company_id = caller.company_id
    )
  order by p.name;
$$;

grant execute on function public.admin_create_company_invite(text, integer) to authenticated;
grant execute on function public.accept_company_invite(uuid) to authenticated;
grant execute on function public.admin_get_company_candidates() to authenticated;
grant execute on function public.admin_get_team_members(bigint) to authenticated;
