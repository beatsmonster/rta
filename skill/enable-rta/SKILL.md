---
name: enable-rta
description: Resume current Claude Code session in tmux for remote access via rta
disable-model-invocation: true
allowed-tools: Bash(tmux *), Bash(mkdir *), Bash(cat *)
---

# Enable Remote tmux Access

Move this Claude Code session into a tmux session so you can reconnect via `rta`.

## Steps

1. Use `${CLAUDE_SESSION_ID}` as the session ID.

2. Generate a short session name (<32 chars):
   ```
   claude-<first-8-chars-of-session-id>
   ```

3. Generate a one-line recap description summarizing what this session is working on.

4. Create the recap JSON:
   ```bash
   mkdir -p ~/.rta/recaps
   cat > ~/.rta/recaps/<name>.json <<'RECAP'
   {"name": "<name>", "recap": "<description>", "project_path": "<cwd>", "session_id": "<session-id>", "updated_at": "<ISO8601>"}
   RECAP
   ```

5. Create a tmux session with auto-start and retry backoff:
   ```bash
   tmux new-session -d -s "<name>" "for i in 1 2 3 4 5; do claude --resume <session-id> && break; sleep \$((2**i)); done"
   ```

6. Set a tmux `client-attached` hook to update the recap on re-attach:
   ```bash
   tmux set-hook -t "<name>" client-attached 'run-shell "~/.rta/bin/update-recap <name>"'
   ```

7. Tell the user:
   > Session '<name>' is ready. Type /exit to leave this session, then run `rta` to connect.
