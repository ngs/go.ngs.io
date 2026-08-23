---
title: slack-bridge-mcp-server
import_path: go.ngs.io/slack-bridge-mcp-server
repo_url: https://github.com/ngs/slack-bridge-mcp-server
description: MCP server that bridges a resident Claude CLI session with Slack via Socket Mode - chat with your local agent from your phone
version: v0.3.0
documentation_url: https://pkg.go.dev/go.ngs.io/slack-bridge-mcp-server
license: MIT
author: Atsushi Nagase
created_at: 2026-08-23T03:43:57Z
updated_at: 2026-08-23T20:25:46Z
---

# slack-bridge-mcp-server

Chat with your local Claude CLI session from your phone, through a private
Slack channel.

This is an [MCP](https://modelcontextprotocol.io) server that bridges a
resident Claude CLI session to your Slack over Socket Mode. The session calls a
blocking `slack_wait` tool in a loop; the bridge holds the WebSocket, catches up
on anything missed while you were away, and posts the replies back.

One private channel is its home, where everything you say reaches the session.
In any other channel you have added the app to, mentioning it opens a
conversation thread — and inside that thread you talk to it without mentioning
it again. Nobody else's messages are relayed, anywhere.

Messages live in Slack, so the bridge only has to remember its position in the
channel. The laptop can sleep, the session can restart, the network can drop —
when the next `slack_wait` runs, everything sent in the meantime comes back as a
backlog.

## Architecture

```
   phone / desktop Slack
            │
            │  private channel
            ▼
     ┌──────────────┐
     │    Slack     │   messages persist here
     └──────┬───────┘
            │  Socket Mode WebSocket (live)
            │  conversations.history  (catch-up)
            ▼
 ┌───────────────────────────────┐
 │  slack-bridge-mcp-server      │
 │  connection · cursor · filter │
 └───────────────┬───────────────┘
                 │  MCP over stdio
                 ▼
        ┌──────────────────┐
        │   Claude CLI     │   parent process
        │  resident session│
        └──────────────────┘
```

The bridge is a child process of the CLI with the same lifetime as the session.
No daemon, no HTTP listener, no launchd job. It does not connect to Slack until
the first tool call that needs Slack, so it can sit in every project's
`.mcp.json` without opening a socket in sessions that never use it.

The full rationale is in [docs/design.md](docs/design.md).

## Setup

**[docs/setup.md](docs/setup.md) is the full walkthrough**, from creating the
Slack app to the first message: scopes, tokens, the channel, installing the
binary, `.mcp.json`, and what to check when something does not work. There is a
[Slack app manifest](docs/slack-app-manifest.yaml) that skips most of the app
configuration.

The short version, if you have done this kind of thing before:

```sh
brew install ngs/tap/slack-bridge-mcp-server   # or: go install go.ngs.io/slack-bridge-mcp-server@latest
```

Create a Slack app with Socket Mode and interactivity on, an app-level token
with `connections:write`, the bot scopes `chat:write`, `groups:history`,
`reactions:write` and `users:read`, and the `message.groups` bot event. Install
it, invite it to a private channel, and set the four variables below.

## Configuration

All configuration is environment variables:

| Variable | Required | Description |
|---|---|---|
| `SLACK_BOT_TOKEN` | yes | Bot user OAuth token, `xoxb-…` |
| `SLACK_APP_TOKEN` | yes | App-level token with `connections:write`, `xapp-…` |
| `SLACK_BRIDGE_CHANNEL` | yes | Channel ID to bridge, `C…` |
| `SLACK_BRIDGE_OWNER` | yes | Your Slack user ID, `U…`. Only this user's messages are relayed. |
| `SLACK_BRIDGE_STATE_DIR` | no | Override the state directory |
| `SLACK_BRIDGE_AUTO_ACK` | no | `off` disables the automatic receipt reaction. Enabled by default. |
| `SLACK_BRIDGE_AUTO_ACK_EMOJI` | no | Emoji for the receipt reaction, without colons. Default `eyes`. |
| `SLACK_BRIDGE_INDICATOR` | no | `off` disables the processing indicator. Enabled by default. |
| `SLACK_BRIDGE_INDICATOR_GRACE` | no | Seconds before the indicator appears. Default 10, clamped to 3–120. |
| `SLACK_BRIDGE_INDICATOR_INTERVAL` | no | Seconds between indicator updates. Default 10, clamped to 5–60. |

An unusable value for either number falls back to the default with a note on
stderr; these settings can never keep the bridge from starting.

If any required variable is missing, the process still starts and serves MCP —
`slack_status` will tell you exactly which ones are unset. The other tools fail
with the same message.

State lives in `~/.config/slack-bridge/` (honouring `XDG_CONFIG_HOME`):

- `state.json` — the cursors, created `0600`: `channels` holds the position in
  the home channel, `threads` the conversations open elsewhere and how far each
  has been read, and `mention_cursor` how far the search for missed mentions has
  looked. A file written by an older version loads unchanged.
- `bridge.lock` — an exclusive lock taken by the first tool call that connects
  to Slack, whichever that is. A second concurrent bridge fails immediately
  rather than splitting your messages between two listeners.

## Wiring it into Claude Code

`.mcp.json`:

```json
{
  "mcpServers": {
    "slack-bridge": {
      "type": "stdio",
      "command": "slack-bridge-mcp-server",
      "env": {
        "SLACK_BOT_TOKEN": "xoxb-...",
        "SLACK_APP_TOKEN": "xapp-...",
        "SLACK_BRIDGE_CHANNEL": "C0123456789",
        "SLACK_BRIDGE_OWNER": "U0123456789"
      }
    }
  }
}
```

`slack_status` answers `"connected": false` until the first tool call that
needs Slack — `slack_wait`, `slack_post`, `slack_ack`, `slack_ask`,
`slack_history`, or `slack_progress` when it has to start an indicator,
whichever comes first — which is the lazy connect working as
intended rather than a problem. [docs/setup.md](docs/setup.md) covers the rest of the first run.

## Tools

| Tool | Arguments | Returns |
|---|---|---|
| `slack_wait` | `timeout_seconds` (optional, default 300, clamped to 5–1500) | `{"messages": [{"ts", "thread_ts"?, "user", "text", "channel", "files"?}…], "timed_out": false}`, oldest first. On timeout, `{"messages": [], "timed_out": true}`. |
| `slack_post` | `text` (required), `thread_ts`, `channel` (optional) | `{"ts", "channel"}` — where the message landed |
| `slack_ack` | `ts` (required), `emoji` (optional, default `eyes`), `channel` (optional) | Confirmation. Receipt is marked automatically, so this is for a deliberate signal beyond it. |
| `slack_ask` | `question` (required), `options` (required, 2–10), `timeout_seconds`, `thread_ts` and `channel` (optional) | `{"choice_index", "choice_label", "ts", "timed_out": false}`. On timeout, `{"choice_index": -1, "timed_out": true}`. |
| `slack_history` | `limit` (optional, default 50, clamped to 1–200), `oldest`, `latest` (exclusive bounds), `thread_ts`, `channel` (all optional) | `{"messages": [{"ts", "user"?, "user_name", "text", "thread_ts"?, "bot", "reply_count"?, "files"?}…], "has_more"}`, oldest first, every author. A limit keeps the newest end of the window. |
| `slack_progress` | `text` (required), `thread_ts` and `channel` (optional — where the status belongs; they start the indicator, or move a running one) | `{"ok": true, "ts"}` — the indicator message the label went on. `ts` is left out while the indicator has yet to post, and `ok` is false when the label had nowhere to go: the indicator is turned off, or the turn ended before it could be applied. |
| `slack_status` | — | `{connected, channel, owner, last_ts, pending_backlog_count, config_error?, state_file}` |

Every tool that speaks takes an optional `channel`, and leaving it out means the
home channel — so a session that never leaves it never has to think about the
argument. Messages come back with the channel they were sent in; pass it back
with `thread_ts` and the reply lands in the conversation it answers.

`slack_progress` is the one exception, and only while an indicator is running:
there the omitted `channel` means the channel that indicator is already in
rather than the home one, because the argument is asking whether to move it.
With nothing running it starts one, and an omitted `channel` means home as
everywhere else.

`slack_post` messages and `slack_ask` questions go out as Block Kit
[`markdown` blocks](https://docs.slack.dev/reference/block-kit/blocks/markdown-block),
so they take standard Markdown — the dialect a model writes by habit — and let
Slack render it: `**bold**`, `# headings`, tables and fenced code all arrive
rendered instead of showing their markup. The text is sent as the message's
`text` as well, which is the fallback Slack puts in push notifications.

A markdown block holds 12,000 characters per message. A longer post goes as the
message body instead — whole, but with only Slack's own mrkdwn applied to it,
so Markdown headings and tables show their markup. Splitting it across blocks
would not help, since the budget is for the whole message rather than for each
block. The `slack_progress` label goes the same way: it shares one line with a
stopwatch, where there is nothing to render.

`slack_wait` caps at 1500 seconds because Claude Code aborts a stdio MCP tool
call after 30 minutes with no response bytes; 25 minutes keeps a margin. A
`timed_out: true` result is not an error — just call it again.

## Talking to the agent in other channels

The home channel is yours, and everything you say in it goes to the session. The
rest of your Slack is not, so the rule there is different: **mention the app and
it starts listening in that thread.**

> **you** in #project-atlas: `@slack-bridge` can you check why the nightly build
> is red?
>
> ⤷ the agent replies in the thread under it, and from then on the two of you
> talk in that thread with no further mentions

That is the whole of it. A mention on a channel message opens the thread under
that message; a mention inside a thread opens that thread. Every message you
send in an open thread reaches the session, carrying its channel so the reply
comes back to the right place. Threads stay open — there is no expiry in this
version — and they survive a restart, so a conversation you had yesterday is
still a conversation today.

Everything else in those channels is ignored: your own messages out on the
surface, your messages in threads nobody opened, and **everybody else's
messages everywhere**, mentions included. A colleague cannot address your agent,
and a channel the app has been added to does not become an inbox. When you want
the agent to read what other people said, ask it to — that is `slack_history`,
below.

Catch-up outside the home channel is best effort, which is the one place this
promises less than the home channel does. On reconnect the bridge re-reads every
open thread from where it left off, and looks for new mentions in the newest
hundred messages of up to twenty channels it belongs to. A mention further back
than that, or in the twenty-first channel, is not found — the bridge says so on
stderr, and mentioning it again is all it takes. The home channel's guarantees
are unchanged.

## Reading the channel

`slack_wait` relays only your messages, which is the right rule for a relay and
the wrong one for "summarise what we decided up there". `slack_history` is the
other mode: ask the agent to read the channel and it gets everything, including
the colleagues, the bots and the incoming webhooks, with display names resolved
and a `reply_count` pointing at any thread worth opening (`thread_ts` reads
that thread).

It is strictly a read. It does not consume anything `slack_wait` would have
delivered, does not move the cursor, does not react, and does not disturb the
indicator — calling it changes nothing except what the model knows.

Names come from `users.info`, which needs the `users:read` scope. Without it
the tool still works and shows raw user IDs, so an app installed before this
existed keeps working until it suits you to reinstall it with the scope.

## Attachments

Send a screenshot and it reaches the session like anything else you say. The
message arrives with a `files` array beside its text, and a caption is optional:
an upload with nothing typed alongside it is delivered on the strength of the
file alone.

```json
{
  "ts": "1787504922.733289",
  "channel": "C0BRIDGE",
  "thread_ts": "1787504116.958859",
  "text": "why is this red?",
  "files": [
    {
      "name": "image.png",
      "mimetype": "image/png",
      "size": 88731,
      "url_private": "https://files.slack.com/files-pri/T0TEAM-F0FILEID/image.png",
      "permalink": "https://example.slack.com/files/U0OWNER/F0FILEID/image.png"
    }
  ]
}
```

The bridge transports metadata, not bytes: it does not download anything, and
nothing is written to disk. `url_private` is where the file actually is, and it
is not a public link — fetching it takes the bot token in an `Authorization`
header and the `files:read` scope:

```sh
curl -H "Authorization: Bearer $SLACK_BOT_TOKEN" -o image.png \
  "https://files.slack.com/files-pri/T0TEAM-F0FILEID/image.png"
```

`files:read` is optional, and only the download needs it. Slack returns the file
metadata on the message whether or not the app has the scope, so delivery, the
names, the sizes and the links all work without it; what changes is that the
`curl` above answers with a Slack login page instead of the file. Add the scope
under **OAuth & Permissions** and reinstall the app if you want the session to
be able to open what you send it. `permalink` opens the file in Slack and needs
no token at all, which makes it the thing to quote back at you.

`slack_history` reports the same `files` on any message that has them, so an
attachment somebody else posted is visible when you ask the agent to read the
channel.

## Receipt and progress

Two signals, and they mean different things:

> **👀** — received: the message reached the session
>
> **⏳ Working… (1m 29s)** — still busy with it

The 👀 goes on every message the moment `slack_wait` hands it over, added by the
server itself rather than by the model, so the owner sees delivery at second
zero instead of whenever the model gets around to it. It is best effort and
happens off to the side: if Slack refuses the reaction, the message is still
delivered and the failure only reaches stderr.

`SLACK_BRIDGE_AUTO_ACK=off` turns it off, and `SLACK_BRIDGE_AUTO_ACK_EMOJI`
picks a different emoji (bare name, no colons). `slack_ack` stays for the
deliberate signals — done, rejected, picked up by hand — with whatever emoji
the moment calls for.

## Asking the owner a question

`slack_ask` is the bridge's version of stopping to ask. The agent posts a
question with one button per answer and blocks until you tap one:

> Ship the fix now, or hold it until the release?
>
> `[ Ship now ]` `[ Hold ]` `[ Ask me later ]`

Tapping rewrites the message with your choice and takes the buttons away, so a
question is answered exactly once; a question nobody answers is marked expired
when the timeout runs out and returns `timed_out: true`. Only the owner's clicks
count, and only one question can be outstanding at a time — a second `slack_ask`
while one is pending is refused rather than queued.

The elapsed-time indicator below stops while the question is up, since the agent
is not the one working, and starts again on your answer.

The app needs **Interactivity** turned on for the clicks to arrive; the
[manifest](docs/slack-app-manifest.yaml) sets it, and for an app installed
before this existed it is one toggle and no reinstall. See
[docs/setup.md](docs/setup.md#1-create-the-slack-app).

## Processing indicator

Once `slack_wait` hands the agent a message, the bridge keeps the channel
posted on how long the answer is taking:

> ⏳ Working… (1m 29s)

It is posted only if the agent is still busy after the grace period, so a quick
reply leaves no trace. From then on the same message is updated in place every
interval, and it is deleted as soon as the agent replies with `slack_post` or
goes back to `slack_wait`. Neither the automatic receipt nor an explicit
`slack_ack` disturbs it — "seen, still working" is exactly when the elapsed time
is worth showing.

It appears wherever the conversation is. A message sent inside a thread gets its
indicator in that thread; one sent on the channel surface gets it out there. A
`slack_ask` asked in a thread starts the next one in the same thread, and when a
batch of messages arrives at once the newest decides, since that is where the
owner spoke last.

The whole feature is best effort: if Slack refuses any of these calls, the
failure is logged to stderr and the tools carry on unaffected. Set
`SLACK_BRIDGE_INDICATOR=off` to turn it off.

## Saying what you are waiting on

A stopwatch says the agent is busy, not what with. When it starts something long
— a CI run, a release pipeline, a build — one `slack_progress` call puts the
answer next to the clock:

> ⏳ Working… (4m 10s) — release chain: waiting for CI

The agent says it once and the server does the rest: the label rides every
update from then on and goes away with the indicator, so there is no second
message to keep alive or clean up. Calling it again replaces the label, and the
next turn starts with a bare stopwatch again.

It also brings the indicator forward. The grace period exists because most
answers arrive in seconds; an agent calling `slack_progress` has just said this
one will not, so the message is posted straight away instead of waiting the
grace period out. If no indicator is running at all — long work started after a
`slack_wait` timed out, say — this starts one, in `thread_ts` if you pass it,
and it retires on the next reply or wait like any other.

**Where the status goes.** The indicator starts wherever you last spoke, which
is a guess about what the agent is working on — and a wrong one as soon as two
topics are in flight, when the label for one lands under the other. So the
`channel` and `thread_ts` on `slack_progress` say where the status belongs, and
naming a different conversation moves the indicator into it: anything standing
in the wrong place is taken down and the indicator appears in the right one,
still counting from when the work started rather than from the move. Inside the
grace period there is nothing to take down yet, and the move is invisible — the
first message the owner sees is the one in the right conversation. Only what
you name changes — a `thread_ts` with no `channel` moves within the channel the
indicator is already in — with one exception: a thread cannot come along to a
different channel, since it identifies a message only within its own, so a move
across channels that names no thread goes to the new channel's surface. A call
that names nothing labels the indicator where it is, which is what a session
talking in one place only ever needs.

With `SLACK_BRIDGE_INDICATOR=off` there is nowhere to put a label, so the call
does nothing and answers `{"ok": false}`.

## Running a resident session

Start a session and give it a loop like this:

```
You are bridged to my Slack via the slack-bridge MCP server. Run this loop and
do not stop:

1. Call slack_wait.
2. If it returns timed_out, go back to step 1.
3. For each message: do what it asks, then reply with slack_post, passing back
   the message's channel and, if it has one, its thread_ts — that is what puts
   the reply in the conversation it answers rather than in my home channel.
   Receipt is already marked for you, so reach for slack_ack only to say
   something an emoji says well — done, rejected, picked up by hand.
4. If you need a decision from me before you can go on, call slack_ask with the
   question and the answers to choose from, in the same channel and thread, and
   act on what I tap.
5. If something is going to take a while — CI, a release, a long build — call
   slack_progress once with what you are waiting on, so I can see it from the
   channel.
6. Go back to step 1.

Keep replies short — I am reading them on a phone.
```

Then message the channel from anywhere — or mention the app in any other
channel it has been added to.

### Example skill

For a session you start this way often, the loop is better kept as a Claude
Code skill than pasted in each time.
[examples/attend/SKILL.md](examples/attend/SKILL.md) is a generic one: the same
loop, plus when to reach for each of the other tools and how to treat what
comes back. Copy it to `.claude/skills/attend/SKILL.md` in a project, or to
`~/.claude/skills/attend/SKILL.md` for every project, and start the session
with `/attend`.

## Manual smoke test

With the binary built, this drives the MCP handshake by hand and should list
the seven tools. It needs no Slack credentials: listing the tools calls none of
them, and it is the calls that connect.

```sh
{ printf '%s\n' \
  '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"smoke","version":"0"}}}' \
  '{"jsonrpc":"2.0","method":"notifications/initialized"}' \
  '{"jsonrpc":"2.0","id":2,"method":"tools/list"}'; sleep 1; } \
| ./slack-bridge-mcp-server 2>/dev/null
```

The `sleep` matters: without it the shell closes stdin as soon as the last
request is written, and the server shuts down on EOF before it has flushed its
replies.

Once credentials are set, calling `slack_status` the same way reports the
channel, the owner and the cursor without opening a socket.

## Security

The bridge relays only messages written by `SLACK_BRIDGE_OWNER`. That is the
whole of the authentication, and it holds in every channel: a colleague
mentioning the app, or replying inside a conversation you opened, reaches
nobody. Bot messages are dropped, which is what stops the agent's own replies
from being read back as new instructions, and so are edits, deletions and join
notices.

Where a message is sent decides whether it is relayed, not whether it is
trusted. `SLACK_BRIDGE_CHANNEL` relays everything you say; anywhere else you
open a conversation by mentioning the app, and only that thread is relayed. The
effect is that adding the app to a channel does not put that channel's traffic
into the agent's context — only your own half of a conversation you started
there.

**Messages arriving over the bridge are external input reaching an agent with
local tool access.** The owner filter authenticates them as coming from your
Slack account, which is exactly as strong as that account and the phone it is
signed in on. Treat access to the channel as equivalent to terminal access on
the machine running the session, and set that session's tool permissions
accordingly.

Attachments are relayed as metadata and nothing more. The bridge never fetches
a file, so a link that arrives is a link the session decided to follow, with
the same care any other external content deserves.

Tokens are read from the environment and are never logged or written to the
state file.

## Development

```sh
make build        # build the binary
make test         # go test ./...
make test-coverage
make fmt
make lint         # golangci-lint
```

stdout carries the JSON-RPC stream, so every diagnostic goes to stderr. A guard
test in `stdout_guard_test.go` parses the source tree and fails if any non-test
file writes to stdout or points `log.SetOutput` anywhere but `os.Stderr`.

Releases: `make set-version VERSION=vX.Y.Z`, then push the tag. CI refuses to
publish if `server/version.go` and the tag disagree.

## License

MIT. See [LICENSE](LICENSE).
