# ctenterd

**ctenterd** is a lightweight container shell agent — a minimal interactive shell designed for use inside containers. It provides a REPL (Read-Eval-Print Loop) with a curated set of built-in commands and falls back to executing binaries found on `PATH` for anything else.

## Features

- Interactive REPL with a clean prompt
- Built-in shell commands that work without external dependencies
- POSIX-style shell syntax parsing via [`mvdan.cc/sh`](https://github.com/mvdan/sh)
- External binary execution with helpful error diagnostics
- Supports single-command mode (`ctenterd <command>`)
- Fully static binary support for distroless / scratch containers

## Built-in Commands

| Command    | Description                                          |
|------------|------------------------------------------------------|
| `cat`      | Print file contents                                  |
| `cd`       | Change the shell working directory                   |
| `cp`       | Copy files or directories                            |
| `echo`     | Display line of text                                 |
| `env`      | Show environment variables                           |
| `kill`     | Terminate processes by PID                           |
| `ls`       | List directory contents                              |
| `mkdir`    | Create directories                                   |
| `mv`       | Move/rename files or directories                     |
| `nslookup` | Resolve DNS records (A/AAAA, CNAME, NS, TXT)         |
| `ping`     | Send ICMP echo requests (requires `CAP_NET_RAW`)     |
| `ps`       | List running processes                               |
| `pwd`      | Print working directory                              |
| `rmdir`    | Remove empty directories                             |
| `touch`    | Create empty files or update timestamps              |
| `whoami`   | Show current user (UID/GID)                          |

## Usage

### Interactive shell

```sh
ctenterd
```

### Run a single command

```sh
ctenterd ls /etc
ctenterd "echo hello world"
```

### Version

```sh
ctenterd --version
ctenterd -V
```

## Building

### Prerequisites

- [Go](https://go.dev/) 1.21+
- GNU Make

### Dynamic build (default)

Produces `bin/ctenterd` (or `bin/ctenterd.exe` on Windows) linked against the host's libc.

```sh
make build
```

### Static build

Produces `bin/static/ctenterd` — a fully self-contained binary with no external library dependencies, suitable for deployment into distroless or scratch container images.

```sh
make static
```

### Build both

```sh
make all
```

### Clean

```sh
make clean
```

## Project Structure

```
ctenterd/
├── main.go            # Entry point — REPL and exec handler
├── builtin/           # Built-in command implementations + registry
├── internal/          # Internal helpers (fs, process utilities)
├── pkg/color/         # Terminal colour helpers
├── bin/               # Build output (created by make)
│   └── static/        # Static build output
└── Makefile
```

## License

See [LICENSE](LICENSE).
