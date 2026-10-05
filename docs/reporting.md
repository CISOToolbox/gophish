# Reporting

Every campaign tracks what each recipient does, from delivery to reporting. Results are visible per campaign in the console and via the [REST API](api.md).

## The funnel

Each recipient advances through a funnel. The campaign summary aggregates the counts:

| Stage | Event | Meaning |
|---|---|---|
| **Email Sent** | `Email Sent` | The message was handed to the SMTP server. |
| **Email Opened** | `Email Opened` | The tracking pixel loaded (implies sent). |
| **Clicked Link** | `Clicked Link` | The target opened the landing page via `{{.URL}}` (implies opened). |
| **Submitted Data** | `Submitted Data` | The target submitted the landing-page form (implies clicked). |
| **Email Reported** | `Email Reported` | The target reported the message as phishing. |

Counts are **cumulative up the funnel**: a submission also counts as a click and an open, so the numbers are internally consistent (every "submitted" is also "clicked", etc.). An **Error** count tracks recipients whose send failed.

On the campaign results page each funnel stage is shown as a **donut** displaying the count and, below it, the percentage of recipients. The donuts are **clickable**: clicking one filters the results table to the recipients **at that status** (for *Email Reported*, those who reported); click it again, or use **Clear filter**, to remove the filter.

## Events timeline

Beyond the aggregate counts, each result has a timeline of individual **events** (with timestamps and details such as browser/user-agent for opens and clicks). This lets you see exactly when and how a given target interacted.

## Per-recipient status

Each result carries a status reflecting how far that recipient got: scheduled → sending → sent → opened → clicked → submitted, plus reported and error states.

## Security-scanner interactions

If a campaign has **[Ignore security scanner interactions](campaigns.md#ignoring-security-scanner-interactions)** enabled, detonations by mail sandboxes (e.g. Safe Links) appear on the timeline as `Clicked Link (scanner)` / `Email Opened (scanner)` but are **not** counted in the funnel above. This keeps the open/click rates reflective of real recipients.

## Declaring a report manually

Recipients often report a simulated phish **out of band** — to the helpdesk, over Slack/Teams, or by phone — rather than through an automated button. An admin can record these: on the campaign results table, click the report icon in the **Reported** column for a recipient and specify the **channel** (e.g. Email, Slack, Phone) and the **time** it was reported. This flags the result as reported (counted in the funnel) and adds a timeline event — placed at the **specified time**, so it appears in its correct chronological position in the recipient's timeline — noting it was a manual report and through which channel.

## Exporting

Campaign results can be **exported to CSV** from the console for offline analysis or to feed an awareness dashboard.

## Reading results via API

- `GET /api/campaigns/{id}/summary` — the aggregate funnel stats.
- `GET /api/campaigns/{id}/results` — per-recipient results and events.
- `GET /api/campaigns/summary` — a roll-up across all campaigns.

See [REST API](api.md).
