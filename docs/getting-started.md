# Getting started

ciso-gophish is a single binary that runs **two HTTP servers**:

- the **admin server** (default `:3333`, HTTPS) — the console and REST API you use to build and run campaigns;
- the **phishing server** (default `:80`, HTTP) — what targets actually reach: it serves landing pages and records opens, clicks and submissions.

Keep these two surfaces separate in any real deployment: the phishing server is public to targets, while the admin server should be restricted (e.g. bound to localhost, behind a VPN, or IP-allowlisted).

## Run with the container image (recommended)

The tool is published as a multi-arch (amd64 + arm64) image on the GitHub Container Registry:

```bash
podman pull ghcr.io/cisotoolbox/ciso-gophish:latest   # or a pinned tag, e.g. :v0.12.2

podman run -d --name ciso-gophish \
  -p 3333:3333 \
  -p 8080:8080 \
  -e PHISH_LISTEN_URL=0.0.0.0:8080 \
  -v ciso-gophish-data:/opt/gophish \
  ghcr.io/cisotoolbox/ciso-gophish:latest
```

`docker` works too (it is a shim to Podman on the reference environment). The container is configured through environment variables applied at startup by `docker/run.sh`; see [configuration](configuration.md) for the full list. Database migrations run automatically on first boot.

## Build from source

Requires **Go 1.26+** and **Node.js 22** (for the front-end assets).

```bash
git clone https://github.com/CISOToolbox/ciso-gophish.git
cd ciso-gophish

# 1. Front-end assets (esbuild)
corepack enable
yarn install --frozen-lockfile
yarn build            # runs `node build.mjs`

# 2. Backend binary
go build             # produces the ./gophish binary

./gophish            # reads ./config.json
```

## First login

On first start the server creates an `admin` account and prints a one-time password to the log:

```
time="..." level=info msg="Please login with the username admin and the password 4304d5255378177d"
```

Open the admin console at <https://localhost:3333> (accept the self-signed certificate), log in, and set a new password when prompted. From there, follow the [typical engagement](README.md#a-typical-engagement) flow.

## Where data lives

All state is stored in a single SQLite database (`gophish.db` by default; path configurable, MySQL also supported). In the container, persist it with a volume mounted at the working directory (`/opt/gophish`) — see the `-v` flag above. See [configuration](configuration.md) for database options.
