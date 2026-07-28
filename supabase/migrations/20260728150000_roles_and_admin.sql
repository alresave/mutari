alter table public.profiles add column if not exists role text not null default 'member'
  check (role in ('admin', 'manager', 'member'));

-- Promotes the existing account when present; the trigger covers a future signup.
update public.profiles set role = 'admin'
where id in (select id from auth.users where lower(email) = 'alresave@gmail.com');

create or replace function public.handle_new_user()
returns trigger language plpgsql security definer set search_path = public as $$
begin
  insert into public.profiles (id, name, role)
  values (
    new.id,
    coalesce(new.raw_user_meta_data ->> 'name', split_part(new.email, '@', 1)),
    case when lower(new.email) = 'alresave@gmail.com' then 'admin' else 'member' end
  );
  return new;
end;
$$;

create or replace function public.is_admin()
returns boolean language sql stable security definer set search_path = public as $$
  select exists(select 1 from public.profiles where id = auth.uid() and role = 'admin');
$$;

create or replace function public.admin_create_company(company_name text, company_website text default '')
returns bigint language plpgsql security definer set search_path = public as $$
declare new_id bigint;
begin
  if not public.is_admin() then raise exception 'admin role required'; end if;
  insert into public.companies(name, website) values (company_name, company_website) returning id into new_id;
  update public.profiles set company_id = new_id where id = auth.uid();
  return new_id;
end;
$$;

create or replace function public.admin_create_team(team_name text, team_description text default '')
returns bigint language plpgsql security definer set search_path = public as $$
declare new_id bigint; current_company bigint;
begin
  if not public.is_admin() then raise exception 'admin role required'; end if;
  select company_id into current_company from public.profiles where id = auth.uid();
  if current_company is null then raise exception 'create or join a company first'; end if;
  insert into public.teams(name, description, company_id) values (team_name, team_description, current_company) returning id into new_id;
  insert into public.team_members(team_id, user_id, role) values (new_id, auth.uid(), 'manager');
  return new_id;
end;
$$;

create or replace function public.admin_add_team_member(target_team bigint, target_user uuid, member_role text default 'member')
returns void language plpgsql security definer set search_path = public as $$
begin
  if not public.is_admin() then raise exception 'admin role required'; end if;
  if not exists (select 1 from public.teams t join public.profiles p on p.company_id=t.company_id where t.id=target_team and p.id=target_user) then
    raise exception 'user and team must belong to the same company';
  end if;
  insert into public.team_members(team_id,user_id,role) values(target_team,target_user,member_role)
  on conflict(team_id,user_id) do update set role=excluded.role;
end;
$$;

grant execute on function public.admin_create_company(text, text) to authenticated;
grant execute on function public.admin_create_team(text, text) to authenticated;
grant execute on function public.admin_add_team_member(bigint, uuid, text) to authenticated;
