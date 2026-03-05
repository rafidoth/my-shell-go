# My own written shell in Go

## Features I've Implemented

- **REPL loop** - Interactive prompt that reads, evaluates, and loops
- **Raw terminal input** - Character-by-character reading with raw mode for real-time key handling
- **Builtin commands** - `echo`, `pwd`, `cd`, `type`, `exit`
- **External command execution** - Runs executables found in `PATH`
- **Command parsing** - Supports single quotes, double quotes, and backslash escaping
- **Tab autocompletion** - Completes builtins and PATH executables from partial input
- **Multi-match completion** - Shows all candidates on tab when multiple matches exist, with longest common prefix fill
- **Bell character** - Rings terminal bell on ambiguous or no match
- **Stdout/Stderr redirection** - Supports `>`, `1>`, `2>`, `>>`, `1>>`, `2>>` to redirect output to files
- **Backspace handling** - Proper terminal backspace with cursor cleanup
- **Ctrl+C handling** - Graceful interrupt signal handling
