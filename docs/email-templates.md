# Email templates

An **email template** is the message sent to targets: a subject and an HTML and/or plain-text body, with optional attachments. Templates use Go's templating syntax, so you can personalize each message and insert the tracked link.

## Fields

| Field | Description |
|---|---|
| **Name** | A label for the template. |
| **Envelope Sender** | Optional override of the SMTP envelope sender. |
| **Subject** | The email subject (supports template variables). |
| **Text** / **HTML** | The message body. Provide HTML for rich mail; a text part improves deliverability. |
| **Attachments** | Optional files attached to the message (also support template variables in their contents/filename). |
| **Add Tracking Image** | Inserts the open-tracking pixel automatically (you can also place it yourself with `{{.Tracker}}`). |

## Importing an existing email

You can import a raw email (e.g. an `.eml` or pasted message source) to seed a template, then edit it. This is handy to reproduce a real-looking message.

## Template variables

Insert these with `{{.Variable}}` anywhere in the subject, body or attachments:

| Variable | Expands to |
|---|---|
| `{{.FirstName}}` | Target's first name |
| `{{.LastName}}` | Target's last name |
| `{{.Email}}` | Target's email address |
| `{{.Position}}` | Target's position |
| `{{.From}}` | The sending profile's From address |
| `{{.URL}}` | The **tracked link** to the landing page (carries the recipient id). **Put this in every clickable link.** |
| `{{.TrackingURL}}` | The URL of the open-tracking endpoint |
| `{{.Tracker}}` | A ready-made invisible `<img>` open-tracking pixel |
| `{{.RId}}` | The recipient's unique result id |
| `{{.BaseURL}}` | The scheme+host of the phishing server (no path/query) |

### Example

```html
<p>Hi {{.FirstName}},</p>
<p>Please review the updated policy here:
   <a href="{{.URL}}">Security Policy</a>.</p>
<p>Thanks,<br>{{.From}}</p>
{{.Tracker}}
```

- `{{.URL}}` is what turns a click into a tracked **Clicked Link** event and sends the target to your landing page.
- `{{.Tracker}}` (or the *Add Tracking Image* option) records the **Email Opened** event when images load.

## API

See [REST API](api.md). Endpoints live under `/api/templates/`.
