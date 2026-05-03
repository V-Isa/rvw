# Roadmap

This roadmap keeps `rvw` focused on a small, useful CLI while leaving room for a more capable terminal review layer.

## MVP: Captured Output Review

- Run `rvw run <command> [args...]` and `rvw run -- <command> [args...]`.
- Run commands with `rvw run --timeout=<duration> -- <command> [args...]`.
- Capture command output after completion.
- Bound captured output with a default 25 MiB cap and transcript truncation marker.
- Review output in a split TUI.
- Navigate line by line with `j`/`k` and arrow keys.
- Add, edit, or clear one comment per output line.
- Export Markdown review to clipboard.
- Keep non-zero child command exits reviewable while showing the child exit code.
- Support macOS first.

## Near-Term

- Add generated-log demo examples to the README.
- Add a screenshot, terminal recording, or demo GIF.
- Add page up/down and jump-to-top/bottom navigation.
- Add search within captured output.
- Add release builds to GitHub Actions.

## Portability

- Add build-tagged clipboard providers:
  - macOS: `pbcopy`
  - Linux: `wl-copy`, `xclip`, or `xsel`
  - Windows: PowerShell `Set-Clipboard`
- Add a non-PTY fallback runner for platforms where PTY support is limited.
- Verify TUI behavior on Linux terminals and Windows Terminal.
- Publish prebuilt binaries for common platforms.
- Add GoReleaser for repeatable tagged releases.
- Add a Homebrew tap or formula after the module path and release process are stable.
- Add shell completions if the command surface grows.

## Product Expansion

- Streaming captured-output mode where the TUI opens immediately and output is appended while a non-interactive command runs.
- Live PTY passthrough mode for interactive commands.
- Optional file export in addition to clipboard export.
- Comment persistence and reloadable review sessions.
- Multiple comments per line if user workflows require it.
- Inline markers or annotations in the output pane.
- Richer Markdown templates for PRs, Slack, Jira, and issue trackers.
- Add man page or generated CLI reference documentation.
- Add a short architecture note for contributors.

## Out of Scope For Now

- Daemon mode.
- Mouse support.
- AI-tool-specific plugin behavior.
- External API calls.
- Telemetry.
- Docker runtime support.
