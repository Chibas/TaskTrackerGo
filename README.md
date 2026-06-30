# TaskTrackerGo
https://roadmap.sh/projects/task-tracker
A simple command-line task tracker built with Go. It lets you create, update, delete, and change the status of tasks, with data stored in a JSON file.

## Features

- Add new tasks
- Update task descriptions
- Delete tasks
- Mark tasks as done or in progress
- List tasks, optionally filtered by status

## Requirements

- Go 1.26 or newer

## Getting Started

Run the application from the project root:

```bash
go run ./src <command> [arguments]
```

## Available Commands

### Add a task

```bash
go run ./src add "Write documentation"
```

### Update a task

```bash
go run ./src update 1 "Write better documentation"
```

### Delete a task

```bash
go run ./src delete 1
```

### Mark a task as in progress

```bash
go run ./src mark-in-progress 1
```

### Mark a task as done

```bash
go run ./src mark-done 1
```

### List tasks

```bash
go run ./src list
```

### List tasks by status

Supported statuses:

- `todo`
- `in-progress`
- `done`

```bash
go run ./src list done
```

## Data Storage

Tasks are stored in a JSON file named `tasks.json` in the project root.

## Project Structure

```text
src/
  main.go
  mapper/
  storage/
  task/
```

## Build

To build the binary:

```bash
go build -o tasktracker ./src
```
