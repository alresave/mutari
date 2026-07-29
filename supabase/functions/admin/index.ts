import { createClient } from "jsr:@supabase/supabase-js@2";

const cors = { "Access-Control-Allow-Origin": "*", "Access-Control-Allow-Headers": "authorization, x-client-info, apikey, content-type" };
Deno.serve(async (request) => {
  if (request.method === "OPTIONS") return new Response("ok", { headers: cors });
  const db = createClient(Deno.env.get("SUPABASE_URL")!, Deno.env.get("SUPABASE_ANON_KEY")!, { global: { headers: { Authorization: request.headers.get("Authorization") ?? "" } } });
  const { data: { user } } = await db.auth.getUser();
  if (!user) return Response.json({ error: "unauthorized" }, { status: 401, headers: cors });
  const body = await request.json().catch(() => ({}));
  let result;
  if (body.action === "create_company") result = await db.rpc("admin_create_company", { company_name: body.name, company_website: body.website ?? "" });
  else if (body.action === "create_team") result = await db.rpc("admin_create_team", { team_name: body.name, team_description: body.description ?? "" });
  else if (body.action === "add_member") result = await db.rpc("admin_add_team_member", { target_team: body.team_id, target_user: body.user_id, member_role: body.role ?? "member" });
  else if (body.action === "create_invite") result = await db.rpc("admin_create_company_invite", { invitee_email: body.email ?? null, expiry_hours: body.expiry_hours ?? 168 });
  else if (body.action === "list_candidates") result = await db.rpc("admin_get_company_candidates");
  else if (body.action === "list_team_members") result = await db.rpc("admin_get_team_members", { target_team: body.team_id });
  else if (body.action === "accept_invite") result = await db.rpc("accept_company_invite", { invite_code: body.invite_code });
  else return Response.json({ error: "unknown action" }, { status: 400, headers: cors });
  return Response.json(result.error ? { error: result.error.message } : { data: result.data }, { status: result.error ? 400 : 200, headers: cors });
});
