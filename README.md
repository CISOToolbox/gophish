![gophish logo](https://raw.github.com/gophish/gophish/master/static/images/gophish_purple.png)

CISO Toolbox — gophish
======================

[![CI](https://github.com/CISOToolbox/gophish/actions/workflows/ci.yml/badge.svg)](https://github.com/CISOToolbox/gophish/actions/workflows/ci.yml)

The phishing-simulation tool of the [CISO Toolbox](https://cisotoolbox.org) project — an authorized phishing toolkit for businesses, penetration testers and security-awareness training. It lets you set up and run phishing engagements end to end, then measure who opened, clicked, submitted data, or reported the message.

> **Authorized use only.** This tool is intended for phishing simulations against users and systems you are explicitly authorized to test, as part of security-awareness training or sanctioned penetration testing.

---

## About this fork

CISO Toolbox is built as a set of modular, self-hostable security-governance tools. For phishing simulation, rather than reimplement a mature capability from scratch, we chose to **reuse and improve an existing, proven tool**: this project is a **fork of [gophish/gophish](https://github.com/gophish/gophish)**, kept close to upstream and extended where it helps our users.

Unlike the other CISO Toolbox modules, **this tool does not integrate into the unified CISO Toolbox suite** (no Pilot SSO, no shared UI): it remains a **standalone** tool, deliberately, and is developed and released on its own cadence in this repository.

What this fork adds on top of upstream gophish:

- **Educational Pages** — a managed entity (full CRUD) selectable per campaign, shown to recipients after they submit the landing-page form.
- **Landing-page import** that can fetch & inline CSS and optionally embed images/fonts as data URIs for self-contained pages.
- Dependency and toolchain modernization (Go modules, GORM v2, esbuild/Biome front-end build) and CSRF handling via the standard-library `CrossOriginProtection`.

## Maintenance & modernization

Upstream gophish has seen little active maintenance in recent years (no new release since **v0.12.1**). As a result the codebase had accumulated significant technical debt and aging — in places vulnerable — dependencies. Bringing the fork up to date was a prerequisite before building new features on it. What was addressed:

- **Dependencies & Go runtime** — updated the Go modules and the Go toolchain (now **Go 1.26**). Known vulnerabilities were cleared: `govulncheck` reports no reachable vulnerabilities and `osv-scanner` reports no npm advisories.
- **ORM migration (GORM v1 → v2)** — the old GORM v1 is unmaintained. The migration fixed the behavioral differences it introduced: statement/`WHERE` reuse across finisher calls, `ErrRecordNotFound` handling (`Find` vs `First`), struct fields mistaken for associations, `Count` into `*int64`, and association upserts.
- **Front-end toolchain** — replaced the legacy webpack/gulp build (which only ran with Node's `--openssl-legacy-provider` workaround) with **esbuild + Biome**, and regenerated the lockfile.
- **CSRF handling** — the previous `gorilla/csrf` dependency was unmaintained and broke cleartext-HTTP admin access; it was replaced by the Go standard library's `net/http` `CrossOriginProtection`.
- **Regressions caught during the refactors** — campaign creation and campaign fetching (GORM v2 fallout), and mail deliverability (Message-Id derived from the sender domain, custom-header handling).
- **CI / supply-chain hardening** — scoped workflow permissions, pinned base image, and build provenance/SBOM in the release pipeline.

The result is a maintained baseline (gophish feature parity, **v0.12.2**) on which this fork's own features — such as [educational pages](docs/educational-pages.md) — are built.

---

## Quickstart

Run the published container image (multi-arch amd64 + arm64):

```bash
podman run -d --name ciso-gophish \
  -p 3333:3333 \        # admin console
  -p 8080:8080 \        # phishing server (landing pages + tracking)
  -e PHISH_LISTEN_URL=0.0.0.0:8080 \
  -v ciso-gophish-data:/opt/gophish \
  ghcr.io/cisotoolbox/ciso-gophish:latest
```

Then open the admin console at <https://localhost:3333> and log in with the username and password printed in the container logs. See **[docs/getting-started.md](docs/getting-started.md)** for the full walkthrough (including building from source).

## Documentation

All documentation lives in this repository under [`docs/`](docs/) — it is self-contained and does not rely on any external site.

| Guide | What it covers |
|---|---|
| [Getting started](docs/getting-started.md) | Run via container, build from source, first login |
| [Concepts](docs/concepts.md) | Object model and the phishing workflow |
| [Configuration](docs/configuration.md) | `config.json`, environment variables, database |
| [Sending profiles](docs/sending-profiles.md) | SMTP configuration and test emails |
| [Email templates](docs/email-templates.md) | Message content, template variables, tracking |
| [Landing pages](docs/landing-pages.md) | Capture pages, site import, redirects |
| [Educational pages](docs/educational-pages.md) | Post-submit awareness pages (fork feature) |
| [Groups & targets](docs/groups.md) | Recipient lists and CSV import |
| [Campaigns](docs/campaigns.md) | Launching, scheduling, tracking |
| [Reporting](docs/reporting.md) | Funnel, events, CSV export |
| [REST API](docs/api.md) | Authentication and endpoints |

## License

This project keeps upstream gophish's MIT license. The CISO Toolbox additions are released under the same terms. See [`LICENSE`](LICENSE).

```
Gophish - Open-Source Phishing Framework

The MIT License (MIT)

Copyright (c) 2013 - 2020 Jordan Wright

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software ("Gophish Community Edition") and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
```
