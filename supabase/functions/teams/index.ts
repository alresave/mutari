import { createClient } from "jsr:@supabase/supabase-js@2";

const cors = { "Access-Control-Allow-Origin": "*", "Access-Control-Allow-Headers": "authorization, x-client-info, apikey, content-type" };
const response = (body: unknown, status = 200) => Response.json(body, { status, headers: cors });

Deno.serve(async (request) => {
  if (request.method === "OPTIONS") return new Response("ok", { headers: cors });
  const client = createClient(Deno.env.get("SUPABASE_URL")!, Deno.env.get("SUPABASE_ANON_KEY")!, { global: { headers: { Authorization: request.headers.get("Authorization") ?? "" } } });
  const { data: { user } } = await client.auth.getUser();
  if (!user) return response({ error: "unauthorized" }, 401);
  const path = new URL(request.url).pathname.replace(/^.*\/teams/, "") || "/";
  const parts = path.split("/").filter(Boolean);
  let result;
  if (!parts.length) result = await client.rpc("get_teams");
  else {
    const teamId = Number(parts[0]);
    if (!Number.isInteger(teamId)) return response({ error: "invalid team id" }, 400);
    if (parts[1] === "matches") result = await client.rpc("get_team_matches", { target_team: teamId, match_limit: 50 });
    else if (parts[1] === "culture") result = await client.rpc("get_team_culture", { target_team: teamId });
    else if (parts[1] === "gaps") result = await client.rpc("get_team_gaps", { target_team: teamId });
    else if (parts[1] === "hiring-profile") result = await client.rpc("get_team_hiring_profile", { target_team: teamId });
    else if (parts[1] === "fit" && parts[2]) result = await client.rpc("get_team_fit", { target_team: teamId, candidate: parts[2] });
    else return response({ error: "not found" }, 404);
  }
  return result.error ? response({ error: result.error.message }, 400) : response(parts.length ? result.data : { teams: result.data });
});
