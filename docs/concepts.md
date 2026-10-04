# Concepts

A phishing engagement is assembled from a few reusable building blocks, then run as a **campaign**. Each block is managed independently in the admin console and via the [REST API](api.md).

## The building blocks

| Object | Role |
|---|---|
| **Sending profile** | The SMTP server (and optional custom headers) used to send the emails. |
| **Email template** | The message itself: subject, HTML/text body, attachments, and the tracked link. |
| **Landing page** | The web page a target lands on after clicking the link. Can capture submitted form data. |
| **Educational page** *(fork feature)* | An awareness page shown to the target **after** they submit the landing-page form. |
| **Group** | A list of targets (recipients), each with email, name and position. |
| **Campaign** | Ties a template + landing page + sending profile + group(s) together, with a launch schedule and the public URL targets will reach. Optionally selects an educational page. |

## The workflow

```
   Campaign launches
        │
        ▼
   Email sent  ──────────────►  [Email Sent]
        │  target opens the message (tracking pixel)
        ▼
   [Email Opened]
        │  target clicks the tracked link  ({{.URL}})
        ▼
   Landing page served  ──────►  [Clicked Link]
        │  target submits the form
        ▼
   [Submitted Data]
        │
        ▼
   Educational page (if set)  ──  else the landing page is served again

   At any point the target may report the mail  ──►  [Email Reported]
```

Each step is recorded per target as an **event**, and aggregated into the campaign **funnel** (see [reporting](reporting.md)).

## The two servers

- The **admin server** hosts the console and API. You never expose it to targets.
- The **phishing server** hosts landing pages and the tracking endpoints. Targets reach it through the campaign **URL**. The tracked link in the template points here, carrying a per-recipient id (`rid`) so every open/click/submit is attributed to the right target.

## Result identifiers

When a campaign is created, each target gets a unique **result id** (`rid`). It is embedded in the tracked link and the tracking pixel, which is how opens, clicks and submissions are tied back to a specific person. Appending a `+` to an `rid` turns the link into a **transparency request** that returns campaign contact information instead of the landing page — useful for recipients who want to verify the engagement.
