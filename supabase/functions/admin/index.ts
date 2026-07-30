import { createClient } from "jsr:@supabase/supabase-js@2";
import { createAdminHandler } from "./handler.ts";

Deno.serve(createAdminHandler(
  (authorization) => createClient(Deno.env.get("SUPABASE_URL")!, Deno.env.get("SUPABASE_ANON_KEY")!, { global: { headers: { Authorization: authorization } } }),
  () => createClient(Deno.env.get("SUPABASE_URL")!, Deno.env.get("SUPABASE_SERVICE_ROLE_KEY")!),
  Deno.env.get("APP_URL") ?? "https://mutari-iota.vercel.app",
));
