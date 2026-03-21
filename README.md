# Focusbar

Focusbar is a small ADHD-friendly menubar todo app built in Go.

It is designed for lightweight focus sessions:

- add a task quickly
- keep one active task at a time
- pause and resume without losing context
- mark tasks done when your brain is ready
- keep the menu simple and low-friction

The app lives in the macOS menu bar and uses `systray` for the UI.

## Why this app

This project is aimed at people who do better with:

- one clear task instead of a large task list
- fast capture with minimal clicks
- visible focus state
- simple momentum over heavy productivity systems

It is not trying to be a full project manager. It is more like a tiny focus companion.

## Features

- menubar-based task flow
- quick add task prompt
- one active working task at a time
- pause, resume, complete, and delete actions
- native macOS notifications from the bundled app
- persistent task storage in `~/Library/Application Support/Focusbar/tasks.json`
- automatic cleanup for old completed tasks
- live timer in the menu bar while working on a task
- short tray-safe titles for multi-display menu bars

## Project structure

```text
focusbar/
├── cmd/
│   └── app/
│       └── main.go
├── internal/
│   ├── app/
│   ├── logx/
│   ├── notifier/
│   ├── reminder/
│   ├── task/
│   ├── timer/
│   ├── tray/
│   └── ui/
├── macos/
├── scripts/
├── go.mod
└── README.md
```

### Package overview

- `cmd/app`: application entrypoint
- `internal/app`: startup wiring and app lifecycle
- `internal/logx`: file logging for debugging app and notification flow
- `internal/notifier`: native macOS notification bridge
- `internal/task`: task model, persistence, cleanup, and task management
- `internal/timer`: focus timer handling
- `internal/tray`: short menu bar title formatting
- `internal/ui`: systray rendering, menu state, and add-task prompt
- `internal/reminder`: reminder evaluation and escalation logic
- `macos`: app bundle metadata
- `scripts`: build and run helpers for the macOS app bundle

## Requirements

- macOS
- Go 1.25+

This app is macOS-specific. It uses:

- `osascript` for the quick add-task prompt
- native macOS notifications through an in-process Objective-C bridge

## Build And Run

```bash
./scripts/run-app-bundle.sh
```

That script:

- builds the app bundle
- ad-hoc signs it
- launches `Focusbar.app`

If you only want to build the bundle:

```bash
./scripts/build-app-bundle.sh
```

The generated app bundle lives at:

```text
dist/Focusbar.app
```

## Install

Install Focusbar into your Applications folder:

```bash
./scripts/install-app.sh
```

Create a release zip for sharing or GitHub Releases:

```bash
./scripts/package-release.sh
```

## Website

A GitHub Pages landing page lives in:

```text
docs/
```

Once GitHub Pages is enabled for the repository, the Pages workflow will publish that folder automatically on pushes to `main`.

## Development Build

If you want to compile the code without launching the app:

```bash
go build ./...
```

## Logs

Focusbar writes logs to:

```text
~/Library/Logs/Focusbar/focusbar.log
```

Tail the log during testing:

```bash
tail -f ~/Library/Logs/Focusbar/focusbar.log
```

## How it works

1. The app starts in the menu bar.
2. Tasks are loaded from the app support directory.
3. If a task was already active, the timer resumes.
4. You can add a new task from the menu.
5. Starting one task pauses any previously active task.
6. Completed tasks are cleaned up automatically after they get old enough.
7. Reminders only trigger when there are pending tasks.

## Data storage

Tasks are stored locally in:

```text
~/Library/Application Support/Focusbar/tasks.json
```

Each task includes:

- id
- title
- state
- created timestamp

## Current task states

- `todo`
- `working`
- `paused`
- `done`

## Future ideas

- recurring reminders
- better daily review flow
- lightweight task history
- SQLite-backed history and sessions
- configurable cleanup windows
- keyboard shortcuts

## License

Add a license here if you want to open source the project.
