# Lathe

Lathe is a simple, declarative C build tool.
It focuses on fast iteration, minimal configuration, and editor/tooling integration.

Lathe is **not** a replacement for Make or CMake. It is intended for
small–to–medium C projects where explicitness and simplicity are preferred.

---

## Features

* Single-command builds (`lathe build`)
* Named targets and reusable aliases
* Debug and release build modes
* Incremental compilation (flag-aware)
* Global and project-level configuration
* Generation of:
  * `compile_commands.json`
  * `compile_flags.txt`
  * `Makefile`
* Shell completions for targets and modes
* Strict configuration validation (unknown keys are errors)

---

## Installation

Build from source:

```sh
go build -o lathe
```

Place the binary somewhere in your `PATH`.

Or install directly:

```sh
go install github.com/Salvadego/lathe@latest
```

---

## Configuration

Lathe reads configuration from two locations:

1. **Global config**
   `~/.config/lathe/lathe.toml`

2. **Project config**
   `./lathe.toml`

Both files are merged. Project configuration overrides global configuration.

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

[alias.sanitize]
ldflags = [
    "-fsanitize=address,undefined",
    "-fno-omit-frame-pointer",
]
```

Aliases are reusable groups of compiler and linker flags that can be shared
across projects.

---

## Project Configuration Example

```toml
[target.snake]
sources = ["main.c"]
aliases = ["raylib", "c99"]
```

A target defines:

* Source files (globs supported)
* Optional includes and defines
* A list of aliases to apply

Compiler and linker flags are supplied via aliases and build modes.

---

## Build Modes

Lathe supports named build modes. By default:

### `debug`

* `-ggdb`
* `-O0`

### `release`

* `-O3`

Modes can be extended or overridden:

```toml
[build.debug]
aliases = ["sanitize"]
```

Build-mode aliases are applied before target aliases.

---

## Incremental Builds

Incremental compilation is **signature-based**.

A source file is recompiled if **any** of the following change:

* Compiler
* Build mode
* Source path
* Compiler flags
* Defines
* Include paths

This prevents stale object files when flags (e.g. sanitizers) are added or removed.

Enable incremental builds with:

```sh
lathe build --incremental
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
* `--incremental` — enable incremental compilation
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

This works automatically in shells with completion enabled.

---

## Philosophy

* Explicit configuration
* Minimal magic
* Fast feedback
* Tooling-friendly

Lathe aims to stay small, predictable, and easy to reason about.

---

## Motivation

Lathe exists because maintaining Makefiles becomes tedious for small projects,
especially when repeatedly copying large sets of compiler and linker flags.

Many C libraries either:

* Require manually copying long flag lists, or
* Depend on external tools (`pkg-config`, `*-config`) that are not always available or consistent

This is especially painful for libraries that do not reliably expose portable
build metadata.

Lathe tries to solve this by:
* Centralizing common flags in reusable aliases
* Allowing aliases to be shared globally across projects
* Keeping build logic explicit and transparent

The goal is convenience without abstraction or loss of control.
