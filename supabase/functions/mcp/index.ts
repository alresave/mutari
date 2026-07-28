import { createClient } from "jsr:@supabase/supabase-js@2";

const corsHeaders = {
  "Access-Control-Allow-Origin": "*",
  "Access-Control-Allow-Headers": "authorization, x-client-info, apikey, content-type",
};

Deno.serve(async (request) => {
  if (request.method === "OPTIONS") return new Response("ok", { headers: corsHeaders });

  const supabase = createClient(
    Deno.env.get("SUPABASE_URL")!,
    Deno.env.get("SUPABASE_ANON_KEY")!,
    { global: { headers: { Authorization: request.headers.get("Authorization") ?? "" } } },
  );
  const url = new URL(request.url);

  if (request.method === "GET" && url.pathname.endsWith("/questions")) {
    const { data, error } = await supabase.from("mcp_questions")
      .select("id, block_id, dimension_index, question_text, scale, mcp_blocks(name, description)")
      .order("dimension_index");
    return Response.json(error ? { error: error.message } : { questions: data }, { status: error ? 400 : 200, headers: corsHeaders });
  }

  const { data: { user } } = await supabase.auth.getUser();
  if (!user) return Response.json({ error: "unauthorized" }, { status: 401, headers: corsHeaders });

  if (request.method === "POST" && url.pathname.endsWith("/answers")) {
    const { answers } = await request.json();
    const { error } = await supabase.rpc("submit_mcp_answers", { answers });
    return Response.json(error ? { error: error.message } : { status: "answers stored" }, { status: error ? 400 : 200, headers: corsHeaders });
  }
  if (request.method === "GET" && url.pathname.endsWith("/matches")) {
    const { data, error } = await supabase.rpc("get_mcp_matches", { match_limit: 50 });
    return Response.json(error ? { error: error.message } : { matches: data }, { status: error ? 400 : 200, headers: corsHeaders });
  }
  return Response.json({ error: "not found" }, { status: 404, headers: corsHeaders });
});
