# Educational pages

> **Fork feature** — educational pages are an addition in ciso-gophish, not present in upstream gophish.

An **educational page** is an awareness page shown to a target **after they submit the landing-page form**. It turns the moment a user falls for a simulation into a teaching moment: explain what happened, what the warning signs were, and what to do next time.

Educational pages are a first-class managed entity with full create / read / update / delete, exactly like landing pages — but they never capture any data. You build a library of them and pick one per campaign.

## Fields

| Field | Description |
|---|---|
| **Name** | A label for the page. |
| **HTML** | The page content (edited in the WYSIWYG/HTML editor). You can also use **Import Site** to start from an existing page. |

Educational pages support the same [template variables](email-templates.md#template-variables) as other content, so you can personalize the message (e.g. `Hi {{.FirstName}}`).

## How it is shown

When a campaign has an educational page selected, the phishing server serves that page **in response to the landing-page form submission** — so a campaign can steer users to a chosen awareness message after they submit.

If no educational page is selected, the landing page is simply served again after submission.

## Using one in a campaign

In the **New Campaign** dialog, the **Educational Page (Optional)** selector sits just below **Landing Page**. Leave it empty to simply re-serve the landing page after submission, or pick a page to show it instead.

## API

See [REST API](api.md). Endpoints live under `/api/educational_pages/` (mirrors `/api/pages/`). A campaign references one via its educational-page id.
