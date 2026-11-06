---
name: coderabbit
description: Address all CodeRabbit review comments on the current PR
---

## 1. Preflight

- `gh auth status`. If it isn't authenticated, check where you're running:
  - **In a Claude Code cloud session** (claude.ai/code), `gh` is unauthenticated by design. It only
    works for search there, and GitHub access goes through the GitHub MCP tools (`mcp__github__*`,
    loaded with ToolSearch if deferred). Use the MCP column of the table below for every GitHub
    call in this skill.
  - **Anywhere else**, stop and tell me. Do not retry blindly.
- `gh pr view --json number,url,headRefName,baseRefName` to identify the PR.

### GitHub MCP fallback

| Need | `gh` | GitHub MCP |
|---|---|---|
| Identify the PR | `gh pr view --json …` | `list_pull_requests` with `head: "<owner>:<branch>"` |
| Review threads, with resolved/outdated state | the GraphQL query in step 2 | `pull_request_read`, `method: get_review_comments` (page with `after`) |
| Review bodies | `…/pulls/{n}/reviews` | `pull_request_read`, `method: get_reviews` |
| PR conversation comments | `gh pr view {n} --comments` | `pull_request_read`, `method: get_comments` |
| Reply in a thread | `gh api …/pulls/{n}/comments -F in_reply_to=…` | `add_reply_to_pull_request_comment` |
| PR-level comment | `gh pr comment` | `add_issue_comment` |

`add_reply_to_pull_request_comment` takes the comment's **numeric** id. That's the number after
`#discussion_r` in the comment's `html_url`, not the thread's `PRRT_…` node id.

## 2. Fetch EVERY finding (two sources, not one)

Inline comments and review-body findings live at different endpoints. Fetch both:

- **Inline (line-anchored) comments:**
  `gh api --paginate repos/{owner}/{repo}/pulls/{n}/comments`
- **Review bodies.** This is where CodeRabbit's "outside diff range" items, nitpick
  blocks, and summary findings live:
  `gh api --paginate repos/{owner}/{repo}/pulls/{n}/reviews`

Filter to CodeRabbit. Its login is `coderabbitai[bot]` over REST, but `coderabbitai` over GraphQL
and the MCP tools. Skip already-resolved threads. Extract the individual findings out of review
bodies: one collapsed section ("Outside diff range comments", "Nitpick comments") often contains
several distinct items, each with its own `file:line` header. Each is a separate finding.

**Always run this GraphQL query too.** The REST comments payload has no resolved flag, so this is
the only read that tells you which threads to skip. It's also the authoritative cross-check: if
either REST call 404s under a token that works elsewhere, don't trust the empty result. If
`totalCount` is over 100, page with `after`.

```bash
gh api graphql -f query='
{ repository(owner:"{owner}", name:"{repo}") {
    pullRequest(number:{n}) {
      reviewThreads(first:100){ totalCount nodes { isResolved isOutdated path
        comments(first:10){ nodes { databaseId author{login} body } } } }
      reviews(first:50){ totalCount nodes { author{login} state body } } } } }'
```

Build a numbered finding list before touching anything, and state the count.

## 2a. Zero findings? Check for an outage BEFORE re-pinging

A "Review failed — An error occurred during the review process" comment, or a
review that simply never arrives, is an outage until proven otherwise. Check both
status pages **before** commenting `@coderabbitai review`:

- **CodeRabbit**: https://statusgator.com/services/coderabbit
- **GitHub**: `curl -s https://www.githubstatus.com/api/v2/components.json` —
  CodeRabbit reads the PR through GitHub's API and is driven by its webhooks, so a
  `major_outage` on **Pull Requests**, **Issues** or **Webhooks** takes CodeRabbit
  down with it even when CodeRabbit's own page is green.

**One nudge maximum.** If it fails again, stop and report. Do not retry into an
outage. A comment posted while Webhooks are degraded can be accepted by the API
and never delivered to the bot, so it *looks* like a nudge was sent when it
wasn't. Before nudging at all, confirm the previous nudge actually landed (compare
the comment count) so you don't post a duplicate.

Two tells worth reading:

- **Fast vs slow failure.** Under a minute means it never fetched the diff;
  ~10 minutes means it got the diff and broke during analysis.
- **Failure timestamp vs the branch's commits.** A failure that predates later
  pushes means those commits were never reviewed at all, rather than partially
  covered. Say so explicitly in the report.

Then stop and tell me. Do not invent work for steps 3–7 when the finding list is
empty; there is nothing to verify, fix, or reply to.

## 3. Verify in parallel with subagents — never one-by-one

**Generated code is rejected, always.** `desktop/src/open-api/` and `mobile/api/` are
openapi-generator output from `api/swagger.yml` and are never edited by hand. Mark every finding in
those paths INVALID yourself, without dispatching an agent. If one exposes a real problem in the
spec, say so in its reason column. The fix would be an `api/swagger.yml` change plus a regen, and
whether that's in scope is my call.

Dispatch subagents (Agent tool, `Explore` or `general-purpose`) to verify the rest
**concurrently**, in a single message with multiple tool calls. Batch
roughly 3–5 related findings per agent (group by file or subsystem so each agent
loads context once), and cap concurrency around 6 agents.

Each agent gets: the finding text, its `file:line`, and this instruction —

> Verify this claim against the actual current code. Read the file and its
> callers. **Do not edit anything.** Return for each finding: VALID (the defect
> is real — quote the exact lines that prove it), INVALID (with the concrete
> reason: stale line reference, already handled elsewhere, misreads the
> framework, conflicts with an existing convention at `<file>`), or NEEDS-HUMAN
> (correct but the fix is a judgment call about product behavior or scope).

Never fix an unverified claim. CodeRabbit flags stale and incorrect items
routinely, and its suggested patches often ignore project conventions.
Cross-reference existing patterns in the codebase rather than applying blindly.

### Always present the table, then wait

Before applying anything, show me the full table. Do this every time, with every finding,
whatever the mix of verdicts:

| # | File:line | Finding | Verdict | Planned action | Reason |
|---|---|---|---|---|---|

The verdict is VALID / INVALID / NEEDS-HUMAN. The planned action is fix / push back / ask me. Then
stop and wait for my go-ahead. The table is a planning artifact for me. Never post it to the PR as
a bulk comment.

## 4. Apply fixes

- Fix every VALID finding I approved, and resolve NEEDS-HUMAN items the way I decided them.
- Never hand-edit `desktop/src/open-api/` or `mobile/api/`. If an approved fix changes
  `api/swagger.yml`, regenerate **both** clients in the same commit (root `CLAUDE.md` →
  "API Client Regeneration").
- Add or update a **unit test** for every behavioral fix. It should fail before
  the fix and pass after.
- If a fix touches user-facing flow that lacks e2e coverage, note it for step 5;
  don't write e2e speculatively.

## 5. Verify

Run the suite and the production build for **each component the fixes touched**, and actually
execute them. `--list`, a typecheck, or a compile pass is **not** a test run.

| Component | Unit suite | Production build |
|---|---|---|
| `api/` | `go test -count=1 ./...` | `go build -v ./...` (what CI runs) |
| `desktop/` | `npm run test:ci` | `npm run build` (`ng build`; `production` is the default configuration) |
| `mobile/` | `flutter analyze` and `flutter test` | `flutter build apk --release` (release signs with the debug key, so no keystore is needed). Without an Android SDK, use `flutter build linux --release`. |
| `mobile/api/` (regen only) | `dart analyze` in `mobile/api`, which must report 0 errors | — |

In the Claude Code cloud sandbox:

- **`api/`** needs ImageMagick 7 built from source, plus
  `CHROMIUM_BINARY_PATH=/opt/pw-browsers/chromium` for the tests. Follow `api/CLAUDE.md` →
  "Running in the Claude Code Web/Cloud Sandbox".
- **`desktop/`** needs `npm install` first. Then `git restore desktop/package-lock.json`, since the
  sandbox's npm rewrites it.
- **`mobile/`** needs Flutter installed first, since it isn't there by default. Follow
  `mobile/CLAUDE.md` → "Flutter SDK Setup (Claude Code Environment)", then run the suite and build
  as above. There's no Android SDK, so the build is `flutter build linux --release`.

### Servers: you start them

When a check needs the running app (an e2e run, or seeing a UI fix in the browser), start the
servers yourself, in the background. Don't ask me to. Check the port first, and reuse a server
that's already up rather than starting a second one.

- **API (`:8081`)**: `cd api && go run main.go`. In the cloud sandbox, follow the `api/CLAUDE.md`
  playbook above. It's ready when `curl -s -o /dev/null -w '%{http_code}' localhost:8081/api/featureConfig`
  prints `200`.
- **Desktop (`:4200`)**: `cd desktop && npm start`. For a `localhost` base URL, Playwright's
  `webServer` starts this itself, so start it by hand only for a manual check. In the cloud sandbox,
  follow `desktop/CLAUDE.md` → "Running in the Claude Code Web/Cloud Sandbox".
- **Mobile e2e** runners start nothing themselves and need the API up.

### E2E: ask first

E2E is user-initiated in this repo. If the changes warrant e2e, name the specs and ask. Only run
them after I say go. The setup is in `desktop/CLAUDE.md` → "E2E Testing" and `mobile/CLAUDE.md` →
"E2E Testing".

## 6. Commit, push, reply

Commit and push first, so replies can reference the resolving commit. Then reply
to **every** finding you evaluated: applied, partially applied, or skipped. An
unanswered CodeRabbit comment looks like it was missed.

- **Inline comments**: reply in the thread so the disposition sits next to the
  finding:
  `gh api repos/{owner}/{repo}/pulls/{n}/comments -f body='...' -F in_reply_to={comment_id}`
  (MCP: `add_reply_to_pull_request_comment`)
- **Outside-diff-range / review-body findings**: no thread exists, so post a
  regular PR comment that quotes the location:
  `gh pr comment {n} --body '...'` (MCP: `add_issue_comment`)

Reply content:
- **Applied**: acknowledge, and cite the resolving commit SHA.
- **Skipped**: give the concrete reason ("matches the existing convention in
  `<file>`", "line reference is stale, the code now reads X", "out of scope for
  this PR", "generated from `api/swagger.yml`, never hand-edited"). Never skip silently.

## 7. Report

A table: finding → file:line → VALID/INVALID/NEEDS-HUMAN → fixed or pushed-back
→ test added. Plus the suite/build results and anything left for me to run.
