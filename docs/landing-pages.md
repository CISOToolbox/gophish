# Landing pages

A **landing page** is the web page a target reaches after clicking the tracked link. It can display content and capture submitted form data. What the target sees after submitting is controlled per campaign by the optional [educational page](educational-pages.md).

## Fields

| Field | Description |
|---|---|
| **Name** | A label for the page. |
| **HTML** | The page content (edited in the WYSIWYG/HTML editor). |
| **Capture Submitted Data** | Record the form fields the target submits. See the warning below. |
| **Capture Passwords** | Also record password fields. **Off by default and strongly discouraged** — see below. |

## Capturing data — read this

- With **Capture Submitted Data** on, submitted field values are stored in the database.
- With **Capture Passwords** on, password values are stored **in cleartext**.

For awareness training you rarely need the actual values — knowing that a target *submitted* is enough. **Recommended policy:** leave *Capture Passwords* off, and prefer capturing nothing (or non-sensitive fields only). Only enable value capture when it is explicitly authorized and justified, and protect the database accordingly.

## Importing a site

The **Import Site** button clones an existing page by URL to use as a starting point. Options:

- **Download & inline CSS** — fetch each linked stylesheet and embed it in the page, so the styling stays local and independent of the original site.
- **Embed images & fonts** — also download images (and CSS-referenced resources) and embed them as `data:` URIs, producing a fully self-contained page. This makes a larger page but removes all external dependencies.

> Single-page apps that build their form in JavaScript may not clone cleanly, since the importer captures the served HTML, not the JS-rendered DOM. Review and adjust imported pages before use.

## Forms and post-submit behavior

- Any `<form>` on the page is rewritten so submissions post back to the phishing server (recording a **Submitted Data** event).
- After submission, if the campaign selects an [educational page](educational-pages.md) it is shown; otherwise the landing page is served again. (The legacy per-page "Redirect to" URL has been removed — use an educational page instead.)

## API

See [REST API](api.md). Endpoints live under `/api/pages/`; site import is `POST /api/import/site`.
