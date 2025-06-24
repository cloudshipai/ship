### Cloudship "Ship CLI" — Product Requirements Document (PRD)

*version 0.9 · 24 Jun 2025*

---

## 1 · Why we're building this

Non-technical founders want Cloudship's **"answers-first"** agents without wrestling with Terraform or secret rotation.
Power users want the **same value in CI** pipelines they already run.
A single CLI gives both groups:

* **One-liner artifact pushes** (tfplan, SBOM, CSV, …).
* **One-command investigations** powered by Steampipe + Dagger, no local installs.
* (Optional) an **MCP host** so any laptop LLM can call Cloudship tools offline.

---

## 2 · Personas & success criteria

| Persona                      | "Job to be done"                                                    | Success signal                                                                                |
| ---------------------------- | ------------------------------------------------------------------- | --------------------------------------------------------------------------------------------- |
| **Founding PM** (Heroku, CF) | Wants cost & perf answers via chat; installs CLI only if team asks. | Auth in <2 min; runs `ship investigate --heroku` once; sees PR suggestions next business day. |
| **SRE**                      | Adds `ship push plan.json` to GitHub Action; expects PR comments.   | First merge shows cost diff comment; FinOps trend in Cloudship UI updates.                    |
| **Platform engineer**        | Automates Steampipe sweeps via MCP & LLM for exploratory audits.    | Local `ship mcp start`; Copilot chooses "steampipe.query" tool and returns JSON.              |

---

## 3 · Top-level flows

1. **Auth & config**

   ```bash
   ship auth --token sk_prod_...   # writes ~/.ship/config.yaml
   ```

   *Looks first at `SHIP_TOKEN` env var, else config file.*

2. **Push artifact** (CI / manual)

   ```bash
   terraform show -json plan.out > plan.json
   ship push plan.json --kind tfplan --env prod
   ```

   *Max upload 100 MB; Worker rejects larger.*

3. **Investigate** (automatic reconnaissance)

   ```bash
   ship investigate --aws  --env prod   # or --cloudflare / --heroku / --gcp
   ```

   *Steps:*

   1. GET `/goals?env=prod` → returns panels & goal IDs.
   2. Local LLM (or built-in prompt templates) maps goal → SQL list.
   3. Boot **Dagger pipeline** → Steampipe container → run SQL.
   4. Artifacts bundle pushed; progress bar shown.

4. **MCP host (optional)**

   ```bash
   ship mcp start --port 9190 [--aws]
   ```

   *Registers tools:*

   * `goal.list` – lists env goals & panels.
   * `goal.run` – runs investigation for given goal.
   * `steampipe.query` – raw SQL pass-through.

---

## 4 · Detailed functional requirements

### 4.1 Auth & config

| Req-ID | Description                                                                  |
| ------ | ---------------------------------------------------------------------------- |
| A-1    | `ship auth` writes `~/.ship/config.yaml` (`token`, `org_id`, `default_env`). |
| A-2    | All sub-commands read `SHIP_TOKEN` env var first, fall back to config file.  |
| A-3    | `ship auth --logout` removes token and cached data.                          |

### 4.2 `ship push`

| Req-ID | Description                                                      |
| ------ | ---------------------------------------------------------------- |
| P-1    | Accepts file path or `-` (stdin).                               |
| P-2    | Requires `--kind` unless auto-detected by file header.          |
| P-3    | Optional flags `--env`, `--tag <key=value>`.                    |
| P-4    | Reject files >100 MB; print error.                              |
| P-5    | On success print artifact SHA & URL to Cloudship UI.            |

### 4.3 `ship investigate`

| Req-ID | Description                                                                                                                          |
| ------ | ------------------------------------------------------------------------------------------------------------------------------------ |
| I-1    | Pulls goal list: `GET /v1/goals?env=<env>&panels=enabled`.                                                                          |
| I-2    | If user passes `--ai`, pipe goals to local LLM via `mcp-go`; else use built-in rule templates.                                      |
| I-3    | Spin Dagger engine (bundled CLI); mount host creds; run composite container with: **steampipe-MCP host + orchestrator**.            |
| I-4    | Auto-install needed Steampipe plugins based on flags.                                                                               |
| I-5    | Bundle each query result JSON in `/out/`; gzip; `ship push` as `kind=steampipe`.                                                    |
| I-6    | Display progress stepper; exit 0 even if some queries error; errors reported in summary table.                                      |

### 4.4 `ship mcp start`

| Req-ID | Description                                                                              |
| ------ | ---------------------------------------------------------------------------------------- |
| M-1    | Exposes JSON-RPC over TCP (default 9190) & stdio.                                       |
| M-2    | Registers tools from orchestrator (`goal.list`, `goal.run`, `steampipe.query`).         |
| M-3    | Uses same Dagger engine & container; keeps it hot until CTRL-C.                         |
| M-4    | Prints markdown block showing tool schema for quick paste into LLM settings.            |

---

## 5 · Non-functional requirements

| Category                 | Target                                                                                           |
| ------------------------ | ------------------------------------------------------------------------------------------------ |
| **Install UX**           | Homebrew, Scoop, curl bash; single static ≈50 MB binary                                         |
| **Runtime dependencies** | Docker/Podman present **or** Daytona sandbox fallback message                                   |
| **Performance**          | First `investigate` cold pull ≤ 150 s; subsequent ≤ 30 s                                        |
| **Security**             | No cloud secrets stored in ~/.ship; only org token.                                             |
| **Artifact integrity**   | SHA-256 computed client-side; server rejects dup SHA unless `--force`.                          |
| **Telemetry opt-in**     | CLI prompts "Send anonymous metrics? (y/N)" on first run.                                       |

---

## 6 · Hi-level technical design

```
ship CLI
 ├─ embeds dagger CLI (10 MB)
 ├─ embeds module index + custom images manifest
 └─ ~/.ship/
        ├─ config.yaml
        └─ cache/

investigate →
   dagger engine start (docker)
      ↳ composite image
          - steampipe-mcp (AGPL)
          - orchestrator (MIT)
          - plugins (aws, cloudflare, heroku, github)
   orchestrator parses goals → SQL list
   for q in list: call `steampipe query -o=json`
   bundle out/* → ship push
```

---

## 7 · Open questions

1. **LLM choice for local AI mapping** – use system Python+OpenAI, or embed `ollama`?
2. **Plugin licensing** – embed only Apache-2 plugins; download others on-demand?
3. **Windows Docker Desktop** cold-start risk? Consider WSL-fallback docs.
4. **Multi-env config** – should `ship auth --env staging` store separate tokens or reuse?

---

## 8 · MVP timeline (8 weeks)

| Week | Milestone                                                                       |
| ---- | ------------------------------------------------------------------------------- |
| 1-2  | Auth & config scaffolding; `/push` command GA; Worker validation.               |
| 3-4  | Bundle Dagger; simple `investigate --aws` runs static SQL list; push artifacts. |
| 5-6  | Goal fetch API; AI mapping using OpenAI (optional flag).                        |
| 7    | MCP host (`ship mcp start`) exposing `goal.list/run`; internal dog-food.        |
| 8    | Beta doc site, Homebrew formula, user onboarding emails.                        |

---

### "Cherry on top": LLM self-serve flow

* User installs `ship mcp start`, pastes the tool schema into ChatGPT *or* local LM.
* As they type "How can I harden my Cloudflare firewall?", their LM calls `goal.list`, picks `tighten-security`, then `goal.run` with `--cloudflare`.
* CLI container executes, pushes artifact, Cloudship FinOps agent writes summary back; LM surfaces it.

---

## 9 · Definition of done

* Push & Investigate commands run on macOS, Linux, Windows (WSL) with Docker present.
* CLI installs in under 30 s; first `ship push tfplan.json` uploads and appears in Cloudship UI.
* At least **3 default providers** supported: `--aws`, `--heroku`, `--cloudflare`.
* `ship mcp start` works with GPT-4o (tool mode) and Ollama (stdio mode).
* Documentation covers CI snippet, local investigations, and non-tech chat path.