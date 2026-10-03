# REST API

Everything you can do in the admin console is available over a JSON REST API served by the **admin server** under `/api/`.

## Authentication

Every request must carry your **API key**, either as a bearer token (preferred) or a query parameter:

```bash
# Bearer token
curl -k -H "Authorization: Bearer $API_KEY" https://localhost:3333/api/campaigns/

# Query parameter
curl -k "https://localhost:3333/api/campaigns/?api_key=$API_KEY"
```

Find your API key under **Settings** in the console; you can rotate it there (or via `POST /api/reset`). The key identifies a user, and permissions are enforced per user (some endpoints require admin).

> Requests to `/api/*` are authenticated by API key only — the console's session/CSRF protection does not apply to the token-authenticated API. Keep the admin server private and treat the API key as a secret.

## Response shape

Reads return the requested object(s) as JSON. Errors return:

```json
{ "success": false, "message": "…", "data": null }
```

## Endpoints

### Campaigns
| Method | Path | Purpose |
|---|---|---|
| GET / POST | `/api/campaigns/` | List / create campaigns |
| GET | `/api/campaigns/summary` | Roll-up stats across campaigns |
| GET / DELETE | `/api/campaigns/{id}` | Fetch / delete a campaign |
| GET | `/api/campaigns/{id}/results` | Per-recipient results and events |
| GET | `/api/campaigns/{id}/summary` | Aggregate funnel for one campaign |
| GET | `/api/campaigns/{id}/complete` | Mark a campaign complete |

### Groups
| Method | Path | Purpose |
|---|---|---|
| GET / POST | `/api/groups/` | List / create groups |
| GET | `/api/groups/summary` | Group summaries (with target counts) |
| GET / PUT / DELETE | `/api/groups/{id}` | Fetch / update / delete a group |

### Templates
| Method | Path | Purpose |
|---|---|---|
| GET / POST | `/api/templates/` | List / create templates |
| GET / PUT / DELETE | `/api/templates/{id}` | Fetch / update / delete a template |

### Landing pages
| Method | Path | Purpose |
|---|---|---|
| GET / POST | `/api/pages/` | List / create landing pages |
| GET / PUT / DELETE | `/api/pages/{id}` | Fetch / update / delete a landing page |

### Educational pages *(fork feature)*
| Method | Path | Purpose |
|---|---|---|
| GET / POST | `/api/educational_pages/` | List / create educational pages |
| GET / PUT / DELETE | `/api/educational_pages/{id}` | Fetch / update / delete an educational page |

### Sending profiles
| Method | Path | Purpose |
|---|---|---|
| GET / POST | `/api/smtp/` | List / create sending profiles |
| GET / PUT / DELETE | `/api/smtp/{id}` | Fetch / update / delete a sending profile |

### Import helpers
| Method | Path | Purpose |
|---|---|---|
| POST | `/api/import/group` | Import targets from CSV |
| POST | `/api/import/email` | Parse a raw email into a template |
| POST | `/api/import/site` | Clone a site into a landing page |

### Utilities & admin
| Method | Path | Purpose |
|---|---|---|
| POST | `/api/util/send_test_email` | Send a test email |
| POST | `/api/reset` | Rotate your API key |
| GET / POST | `/api/users/` | List / create users *(admin)* |
| GET / PUT / DELETE | `/api/users/{id}` | Manage a user *(admin)* |
| — | `/api/webhooks/…` | Manage webhooks *(admin)* |

## Example: create an educational page

```bash
curl -k -X POST "https://localhost:3333/api/educational_pages/" \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"name":"Awareness","html":"<html><body>You took part in a phishing simulation…</body></html>"}'
```
