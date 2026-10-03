# Campaigns

A **campaign** runs an engagement: it ties together an [email template](email-templates.md), a [landing page](landing-pages.md), a [sending profile](sending-profiles.md) and one or more [groups](groups.md), then sends the mail on a schedule and tracks what each recipient does.

## Creating a campaign

| Field | Description |
|---|---|
| **Name** | A label for the campaign. |
| **Email Template** | The message to send. |
| **Landing Page** | The page shown after a click. |
| **Educational Page (Optional)** | *(fork feature)* An [educational page](educational-pages.md) shown after the target submits the form. Takes precedence over the landing page's redirect URL. Leave empty to keep classic behavior. |
| **URL** | The public base URL targets will reach — i.e. where your **phishing server** is reachable from the targets' network (e.g. `https://links.example.com`). This is what `{{.URL}}` is built from, so it must be reachable by targets. |
| **Launch Date** | When sending starts. |
| **Send Emails By** (optional) | If set, emails are spread evenly between the launch date and this date, instead of all at once. |
| **Sending Profile** | The SMTP profile used to send. |
| **Groups** | The recipient lists to target (de-duplicated by email). |

> The **URL** must point at the phishing server as seen by your targets, not at `localhost`. Getting this wrong is the most common reason links don't track.

## Lifecycle

1. **Queued / In progress** — on launch, a result and a scheduled mail are created for each unique recipient.
2. **Sending** — emails go out (immediately, or staggered if *Send Emails By* is set).
3. **Tracking** — opens, clicks and submissions are recorded per recipient as they happen (see [reporting](reporting.md)).
4. **Completed** — mark a campaign complete to stop further processing. Completed campaigns no longer record new events.

You can **launch**, **complete/stop**, **export**, and **delete** campaigns from the console.

## Previewing

Use a [test email](sending-profiles.md#sending-a-test-email) to preview how the template renders and to confirm the tracked link works before launching to real groups.

## API

See [REST API](api.md). Endpoints live under `/api/campaigns/`.
