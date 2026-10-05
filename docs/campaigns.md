# Campaigns

A **campaign** runs an engagement: it ties together an [email template](email-templates.md), a [landing page](landing-pages.md), a [sending profile](sending-profiles.md) and one or more [groups](groups.md), then sends the mail on a schedule and tracks what each recipient does.

## Creating a campaign

| Field | Description |
|---|---|
| **Name** | A label for the campaign. |
| **Email Template** | The message to send. |
| **Landing Page** | The page shown after a click. |
| **Educational Page (Optional)** | *(fork feature)* An [educational page](educational-pages.md) shown after the target submits the form. Leave empty to simply re-serve the landing page. |
| **URL** | The public base URL targets will reach — i.e. where your **phishing server** is reachable from the targets' network (e.g. `https://links.example.com`). This is what `{{.URL}}` is built from, so it must be reachable by targets. |
| **Launch Date** | When sending starts. |
| **Send Emails By** (optional) | If set, emails are spread evenly between the launch date and this date, instead of all at once. |
| **Ignore security scanner interactions** | *(fork feature)* When enabled, opens/clicks that look like a mail security scanner or sandbox are logged but **not counted** as recipient activity. See below. |
| **Sending Profile** | The SMTP profile used to send. |
| **Groups** | The recipient lists to target (de-duplicated by email). |

## Ignoring security-scanner interactions

Mail security products (e.g. Microsoft Defender **Safe Links**, Proofpoint, Mimecast) automatically **detonate** links and pre-fetch tracking pixels, which would otherwise inflate a campaign's *opened* and *clicked* numbers with activity that is not a real recipient.

When **Ignore security scanner interactions** is enabled on a campaign, the phishing server reclassifies an interaction as a scanner when **both** of these hold:

1. it occurs within a short window (~2 minutes) of the email being sent to that recipient (sandboxes detonate near-instantly), **and**
2. its `User-Agent` looks like a scanner — either a known scanner/automation string (headless browsers; vendors such as Microsoft, Proofpoint, Mimecast, Barracuda, …), **or an outdated browser version**. The latter is key for Microsoft Defender Safe Links, which detonates with a clean but *pinned old* Chrome build (e.g. Chrome 109) while real recipients run current, auto-updating browsers.

So a genuine recipient clicking from an up-to-date browser is still counted even if it lands inside the window, while a clean-UA-but-outdated sandbox detonation is caught.

The window, the minimum browser versions (per engine), and the list of risky User-Agent substrings are all editable by an admin under Settings → **Scanner Detection** (tab) (no config file or restart needed). Set a minimum version to `0` to disable the outdated-browser check for that engine.

Matched interactions are recorded on the timeline as **`Clicked Link (scanner)`** / **`Email Opened (scanner)`** for transparency, but they do **not** advance the recipient's status, so they are excluded from the funnel. A genuine open/click from the same recipient afterwards is counted normally. Form submissions (POST) are always counted. The option is **off by default**.

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
