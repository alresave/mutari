export type RpcResult = { data: unknown; error: { message: string } | null };

export interface AdminClient {
  auth: { getUser(): Promise<{ data: { user: unknown } }> };
  rpc(name: string, args?: Record<string, unknown>): PromiseLike<RpcResult>;
}

export interface ServiceClient {
  auth: {
    admin: {
      inviteUserByEmail(email: string, options: { redirectTo: string }): Promise<{ error: { message: string } | null }>;
    };
  };
}

const cors = { "Access-Control-Allow-Origin": "*", "Access-Control-Allow-Headers": "authorization, x-client-info, apikey, content-type" };

export function createAdminHandler(
  createClient: (authorization: string) => AdminClient,
  createServiceClient: () => ServiceClient,
  appUrl = "https://mutari-iota.vercel.app",
) {
  return async (request: Request) => {
    if (request.method === "OPTIONS") return new Response("ok", { headers: cors });
    const db = createClient(request.headers.get("Authorization") ?? "");
    const { data: { user } } = await db.auth.getUser();
    if (!user) return Response.json({ error: "unauthorized" }, { status: 401, headers: cors });

    const body = await request.json().catch(() => ({}));
    let result: RpcResult;
    if (body.action === "create_company") result = await db.rpc("admin_create_company", { company_name: body.name, company_website: body.website ?? "" });
    else if (body.action === "create_team") result = await db.rpc("admin_create_team", { team_name: body.name, team_description: body.description ?? "" });
    else if (body.action === "add_member") result = await db.rpc("admin_add_team_member", { target_team: body.team_id, target_user: body.user_id, member_role: body.role ?? "member" });
    else if (body.action === "create_invite") {
      result = await db.rpc("admin_create_company_invite", { invitee_email: body.email ?? null, expiry_hours: body.expiry_hours ?? 168 });
      if (!result.error && body.email) {
        const invite = result.data as { invite_code: string };
        const { error } = await createServiceClient().auth.admin.inviteUserByEmail(body.email, {
          redirectTo: `${appUrl}/app/admin.html?invite=${invite.invite_code}`,
        });
        result.data = { ...invite, email_sent: !error, email_error: error?.message ?? null };
      }
    } else if (body.action === "list_candidates") result = await db.rpc("admin_get_company_candidates");
    else if (body.action === "list_team_members") result = await db.rpc("admin_get_team_members", { target_team: body.team_id });
    else if (body.action === "accept_invite") result = await db.rpc("accept_company_invite", { invite_code: body.invite_code });
    else return Response.json({ error: "unknown action" }, { status: 400, headers: cors });

    return Response.json(result.error ? { error: result.error.message } : { data: result.data }, { status: result.error ? 400 : 200, headers: cors });
  };
}
