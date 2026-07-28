// Set these public values in supabase-config.js. Never place a service-role key here.
window.Mutari = (() => {
  const config = window.MUTARI_SUPABASE_CONFIG || {};
  let clientPromise;

  function client() {
    if (!config.url || !config.anonKey) {
      return Promise.reject(new Error("Supabase is not configured. Set web/app/supabase-config.js."));
    }
    clientPromise ||= import("https://esm.sh/@supabase/supabase-js@2")
      .then(({ createClient }) => createClient(config.url, config.anonKey));
    return clientPromise;
  }

  async function session() {
    const { data } = await (await client()).auth.getSession();
    return data.session;
  }

  async function mcp(path, init = {}) {
    const current = await session();
    const headers = new Headers(init.headers);
    headers.set("apikey", config.anonKey);
    if (current) headers.set("Authorization", `Bearer ${current.access_token}`);
    return fetch(`${config.url}/functions/v1/mcp${path}`, { ...init, headers });
  }

  async function teams(path, init = {}) {
    const current = await session();
    if (!current) throw new Error("Authentication required");
    const headers = new Headers(init.headers);
    headers.set("apikey", config.anonKey);
    headers.set("Authorization", `Bearer ${current.access_token}`);
    return fetch(`${config.url}/functions/v1/teams${path}`, { ...init, headers });
  }

  return { client, session, mcp, teams };
})();
