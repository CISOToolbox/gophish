# Sending profiles

A **sending profile** describes the SMTP server used to send a campaign's emails, plus any custom headers.

## Fields

| Field | Description |
|---|---|
| **Name** | A label for the profile. |
| **Interface Type** | SMTP (the mail transport). |
| **From** | The envelope/From address recipients see, e.g. `IT Support <it@example.com>`. |
| **Host** | SMTP server `host:port`, e.g. `smtp.example.com:587`. |
| **Username** / **Password** | Credentials for authenticated SMTP (leave empty for an open relay you control). |
| **Ignore Certificate Errors** | Skip TLS verification of the SMTP server. Use only against servers you trust. |
| **Email Headers** | Optional custom headers added to every message (name/value pairs). |

## Custom headers

Headers are useful to influence deliverability or to tag simulation mail. For example, some senders add an `X-Mailer` header or a custom marker so downstream filters can treat simulation mail differently. Add them as name/value pairs on the profile.

> **Message-Id:** outgoing messages use a Message-Id derived from the sender domain, which helps legitimate delivery. Even so, outbound spam filtering at your provider is driven mostly by **content and reputation** — a header alone will not bypass a provider that classifies the message body as spam. If simulation mail is rejected, review the message content and the sending domain's reputation (SPF/DKIM/DMARC).

## Sending a test email

From a sending profile (or the campaign dialog) you can send a **test email** to a single address to verify SMTP connectivity and preview how the template renders for a given recipient. Test emails use the same [template variables](email-templates.md#template-variables) as a real campaign.

## API

See [REST API](api.md). Endpoints live under `/api/smtp/`.
