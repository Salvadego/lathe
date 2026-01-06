# Lathe

Lathe is a simple, declarative C build tool.  
It focuses on fast iteration, minimal configuration, and editor/tooling integration.

Lathe is **not** a replacement for Make or CMake. It is intended for
small–to–medium C projects where explicitness and simplicity are preferred.

---

## Features

- Single-command builds (`lathe build`)
- Named targets and reusable aliases
- Debug and release build modes
- Incremental compilation
- Global and project-level configuration
- Generation of:
  - `compile_commands.json`
  - `compile_flags.txt`
  - `Makefile`
- Shell completions for targets and modes

---

## Installation

Build from source:

```sh
go build -o lathe
````

Place the binary somewhere in your `PATH`.
Or simply:

---

## Configuration

Lathe reads configuration from two locations:

1. **Global config**
   `~/.config/lathe/lathe.toml`

2. **Project config**
   `./lathe.toml`

Both files are merged. Project config overrides global config where applicable.

---

## Global Configuration Example

```toml
[alias.c99]
cflags = [
    "-std=c99",
    "-Wall",
    "-Wextra",
    "-Werror",
    "-Wpedantic",
    "-Wconversion",
    "-Wformat=2",
    "-Wno-unused-parameter",
    "-Wshadow",
    "-Wwrite-strings",
    "-Wstrict-prototypes",
    "-Wold-style-definition",
    "-Wredundant-decls",
    "-Wnested-externs",
    "-Wmissing-include-dirs",
]

[alias.raylib]
ldflags = [
    "-lraylib",
    "-lGL",
    "-lm",
    "-lpthread",
    "-ldl",
    "-lrt",
    "-lX11",
]

[alias.size]
cflags = ["-Os"]
```

Aliases are reusable flag groups that can be attached to targets.

---

## Project Configuration Example

```toml
[target.snake]
sources = ["main.c"]
cflags = ["-static"]
aliases = ["raylib", "c99", "size"]
```

A target defines:

* Source files (globs supported)
* Optional includes and defines
* A list of aliases to apply

---

## Build Modes

Lathe supports named build modes. By default:

* **debug**

  * `-ggdb`
  * `-O0`
  * Address + undefined behavior sanitizers enabled

* **release**

  * `-O3`

Modes can be extended or overridden in config:

```toml
[build.debug]
cflags = ["-g3"]
sanitize = true
```

---

## Commands

### `lathe build`

Build a target or a list of sources.

```sh
lathe build -t snake
lathe build main.c util.c
```

Options:

* `-m, --mode` — build mode (`debug` or `release`)
* `-t, --target` — target name
* `--workdir` — build output directory (default: `.lathe`)
* `-v, --verbose` — print compiler and linker commands

---

### `lathe run`

Build and run a target or sources.

```sh
lathe run -t snake
lathe run main.c -- arg1 arg2
```

Arguments after `--` are passed to the program.

---

### `lathe gen`

Generate build artifacts.

#### `compile-commands`

```sh
lathe gen compile-commands -t snake
```

Generates `compile_commands.json` in the project root.

#### `compile-flags`

```sh
lathe gen compile-flags -t snake
```

Generates `compile_flags.txt`.

#### `makefile`

```sh
lathe gen makefile -t snake
```

Generates a standalone `Makefile`.

---

## Output Layout

By default, build artifacts are placed under:

```
.lathe/
  debug/
    path/to/source.c.o
  release/
    path/to/source.c.o
```

Final binaries are written to:

```
.lathe/<target-name>
```

---

## Compiler Selection

The compiler defaults to `gcc`.

Override globally or per-project:

```toml
[compiler]
cc = "clang"
```

---

## Shell Completions

Lathe provides dynamic completion for:

* `--target`
* `--mode`

This works automatically when using a shell with completion enabled.

---

## Philosophy

* Explicit configuration
* Minimal magic
* Fast feedback
* Tooling-friendly

Lathe aims to stay small, predictable, and easy to reason about.


## Motivation

Lathe exists because I was tired of writing and maintaining Makefiles,
especially repeating the same compiler and linker flags across multiple
projects.

Many C libraries either:
- Require copying long lists of flags by hand, or
- Rely on external tools (`pkg-config`, `*-config`) that are not always available or consistent

This becomes particularly painful for libraries such as clients and others that
do not reliably expose their flags in a portable way.

Lathe 'tries' to solve this by:
- Centralizing common flags in reusable aliases
- Allowing those aliases to be shared globally across projects
- Keeping build logic simple, explicit, and transparent

The goal is not abstraction, but convenience without losing control.
