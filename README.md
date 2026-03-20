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
- task persistence in `tasks.json`
- automatic cleanup for old completed tasks
- live timer in the menu bar while working on a task

## Project structure

```text
focusbar/
├── cmd/
│   └── app/
│       └── main.go
├── internal/
│   ├── app/
│   ├── reminder/
│   ├── task/
│   ├── timer/
│   └── ui/
├── tasks.json
├── go.mod
└── README.md
```

### Package overview

- `cmd/app`: application entrypoint
- `internal/app`: startup wiring and app lifecycle
- `internal/task`: task model, persistence, cleanup, and task management
- `internal/timer`: focus timer handling
- `internal/ui`: systray rendering, menu state, and add-task prompt
- `internal/reminder`: notification-related logic

## Requirements

- macOS
- Go 1.25+

This app uses `osascript` for prompts and notifications, so it is currently macOS-specific.

## Run locally

```bash
go run ./cmd/app
```

## Build

```bash
go build ./...
```

If you want a specific binary output:

```bash
go build -o focusbar ./cmd/app
```

## How it works

1. The app starts in the menu bar.
2. Tasks are loaded from `tasks.json`.
3. If a task was already active, the timer resumes.
4. You can add a new task from the menu.
5. Starting one task pauses any previously active task.
6. Completed tasks are cleaned up automatically after they get old enough.

## Data storage

Tasks are stored locally in:

```text
tasks.json
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
- configurable cleanup windows
- keyboard shortcuts

## License

Add a license here if you want to open source the project.
