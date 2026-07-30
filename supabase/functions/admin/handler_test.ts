import { assertEquals } from "jsr:@std/assert@1";
import { createAdminHandler, type AdminClient, type RpcResult, type ServiceClient } from "./handler.ts";

const request = (body: unknown) => new Request("https://functions.test/admin", {
  method: "POST", headers: { Authorization: "Bearer token", "Content-Type": "application/json" }, body: JSON.stringify(body),
});

function setup({ user = { id: "admin" }, rpcResult = { data: { invite_code: "invite-1" }, error: null }, emailError = null }: {
  user?: unknown; rpcResult?: RpcResult; emailError?: { message: string } | null;
} = {}) {
  const calls: { rpc: [string, Record<string, unknown> | undefined][]; email: [string, { redirectTo: string }][] } = { rpc: [], email: [] };
  const db: AdminClient = {
    auth: { getUser: () => Promise.resolve({ data: { user } }) },
    rpc: (name, args) => { calls.rpc.push([name, args]); return Promise.resolve(structuredClone(rpcResult)); },
  };
  const service: ServiceClient = { auth: { admin: { inviteUserByEmail: (email, options) => { calls.email.push([email, options]); return Promise.resolve({ error: emailError }); } } } };
  return { calls, handler: createAdminHandler(() => db, () => service, "https://app.example") };
}

Deno.test("rejects admin operations without an authenticated user", async () => {
  const { handler } = setup({ user: null });
  const response = await handler(request({ action: "list_candidates" }));
  assertEquals(response.status, 401);
  assertEquals(await response.json(), { error: "unauthorized" });
});

Deno.test("creates an invite and emails its exact acceptance URL", async () => {
  const { handler, calls } = setup();
  const response = await handler(request({ action: "create_invite", email: "candidate@example.com", expiry_hours: 24 }));
  assertEquals(response.status, 200);
  assertEquals(calls.rpc, [["admin_create_company_invite", { invitee_email: "candidate@example.com", expiry_hours: 24 }]]);
  assertEquals(calls.email, [["candidate@example.com", { redirectTo: "https://app.example/app/admin.html?invite=invite-1" }]]);
  assertEquals(await response.json(), { data: { invite_code: "invite-1", email_sent: true, email_error: null } });
});

Deno.test("returns an email-delivery failure without losing the invite", async () => {
  const { handler } = setup({ emailError: { message: "mail unavailable" } });
  const response = await handler(request({ action: "create_invite", email: "candidate@example.com" }));
  assertEquals(response.status, 200);
  assertEquals(await response.json(), { data: { invite_code: "invite-1", email_sent: false, email_error: "mail unavailable" } });
});

Deno.test("routes candidate and team roster actions to their protected RPCs", async () => {
  const { handler, calls } = setup({ rpcResult: { data: [], error: null } });
  await handler(request({ action: "list_candidates" }));
  await handler(request({ action: "list_team_members", team_id: 42 }));
  await handler(request({ action: "add_member", team_id: 42, user_id: "candidate", role: "lead" }));
  assertEquals(calls.rpc, [
    ["admin_get_company_candidates", undefined],
    ["admin_get_team_members", { target_team: 42 }],
    ["admin_add_team_member", { target_team: 42, target_user: "candidate", member_role: "lead" }],
  ]);
});
