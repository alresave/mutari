# Mutari

Psychological compatibility engine for teams and relationships.

## Features

- 64 dimension compatibility model
- pgvector similarity search
- team fit analysis
- relationship stability simulation
- JWT authentication
- OpenAPI / Swagger documentation

## Run locally

make setup
make dev

## API Docs

http://localhost:8080/swagger/index.html

## Demo

http://localhost:8080/app/dashboard.html

## Supabase migration

The Supabase implementation lives in `supabase/` and uses Supabase Auth rather than the
legacy custom JWT flow. Apply it to a linked project with:

```sh
supabase db push --include-seed
supabase functions deploy mcp
supabase functions deploy teams
supabase functions deploy admin
```

The `mcp` Edge Function exposes `GET /questions`, `POST /answers`, and `GET /matches`.
The `teams` function exposes the team analytics routes, and `admin` covers company and
team administration. Authenticated requests use the Supabase access token; authorization
is enforced again inside the database RPCs. The Go service is now legacy-only and is not
used by the deployed web clients.

Admins can create a company invite through the `admin` function (`action: "create_invite"`).
The invited account accepts it with `action: "accept_invite"`; after completing the MCP
questionnaire it becomes available as a candidate. Admin roster actions are
`list_candidates`, `list_team_members`, and `add_member`.

Set the project's URL and **anon** key in `web/app/supabase-config.js`. Both interfaces
load that shared configuration. Those keys are designed to be public; never expose a
Supabase service-role key in a browser file.

## Vercel

The root `vercel.json` deploys `web/` as a static site; it does not build or deploy the
legacy Go API. Import the repository in Vercel, leave the framework preset as **Other**,
and deploy. The root URL redirects to `/app/index.html`; the dating experience is at
`/dating/index.html`. `supabase-config.js` is deliberately sent with `no-store` so a
configuration update is available immediately after deployment.
