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
```

The `mcp` Edge Function exposes `GET /questions`, `POST /answers`, and `GET /matches`
below its function URL. The caller must send its Supabase access token for answers and
matches. Keep the Go API running during rollout, then point the web clients to the
Supabase project and remove it only after data and UI migration are verified.

Set the project's URL and **anon** key in `web/app/supabase-config.js`. Both interfaces
load that shared configuration. Those keys are designed to be public; never expose a
Supabase service-role key in a browser file.

## Vercel

The root `vercel.json` deploys `web/` as a static site; it does not build or deploy the
legacy Go API. Import the repository in Vercel, leave the framework preset as **Other**,
and deploy. The root URL redirects to `/app/index.html`; the dating experience is at
`/dating/index.html`. `supabase-config.js` is deliberately sent with `no-store` so a
configuration update is available immediately after deployment.
