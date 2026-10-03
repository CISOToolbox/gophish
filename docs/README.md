# ciso-gophish documentation

Self-contained documentation for **ciso-gophish**, the phishing-simulation tool of the [CISO Toolbox](https://cisotoolbox.org) project (a fork of [gophish](https://github.com/gophish/gophish)).

This documentation lives entirely in the repository and does not depend on any external website.

## Contents

1. [Getting started](getting-started.md) — run via container, build from source, first login
2. [Concepts](concepts.md) — the object model and the phishing workflow
3. [Configuration](configuration.md) — `config.json`, environment variables, database
4. [Sending profiles](sending-profiles.md) — SMTP configuration and test emails
5. [Email templates](email-templates.md) — message content, template variables, tracking
6. [Landing pages](landing-pages.md) — capture pages, site import, redirects
7. [Educational pages](educational-pages.md) — post-submit awareness pages (fork feature)
8. [Groups & targets](groups.md) — recipient lists and CSV import
9. [Campaigns](campaigns.md) — launching, scheduling, tracking
10. [Reporting](reporting.md) — funnel, events, CSV export
11. [REST API](api.md) — authentication and endpoints

## A typical engagement

1. Create a **[sending profile](sending-profiles.md)** (the SMTP server used to send mail).
2. Write an **[email template](email-templates.md)** (the message, with a tracked link).
3. Build a **[landing page](landing-pages.md)** (what the target sees after clicking).
4. Optionally create an **[educational page](educational-pages.md)** (shown after the target submits the form).
5. Import your **[groups of targets](groups.md)**.
6. Launch a **[campaign](campaigns.md)** tying these together.
7. Watch results in **[reporting](reporting.md)**.

> **Authorized use only** — run simulations exclusively against users and systems you are authorized to test.
