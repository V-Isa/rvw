# Security Policy

## Reporting A Vulnerability

Please report security issues privately before opening a public issue.

Preferred contact: open a private security advisory on GitHub after the repository is public.

If private advisories are unavailable, contact the maintainer through the GitHub profile associated with this repository and include:

- Affected version or commit.
- Steps to reproduce.
- Impact and expected behavior.
- Any relevant logs or terminal output, with secrets removed.

## Scope

`rvw` is a local CLI tool that executes commands provided by the user. Reports are most useful when they involve behavior outside that expected model, such as unintended shell invocation, unsafe file handling, clipboard leaks, terminal escape handling issues, or dependency vulnerabilities.

Please do not include secrets, tokens, private command output, or personal data in reports.

## Supported Versions

Before the first stable release, only the latest commit on the default branch is supported.
