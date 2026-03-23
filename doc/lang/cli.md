# CLI Reference

`len-cli` is the command-line tool for the `len` meta-programming language. The current release provides a parser and structural validator for `.l1` source files.

---

## Installation

### Build from Source

Prerequisites: Go 1.21 or later.

```bash
git clone <repo>
cd len
make build
```

The compiled binary is written to `bin/len`. Add it to your `PATH`:

```bash
export PATH="$PATH:$(pwd)/bin"
```

To run tests:

```bash
make test
```

---

## Commands

### `len-cli validate`

Parse and validate one or more `.l1` files or directories.

**Synopsis**

```
len-cli validate <path> [<path> ...]
```

**Arguments**

| Argument | Description |
|----------|-------------|
| `<path>` | A `.l1` file or a directory. Directories are walked recursively; only `.l1` files are processed. Multiple paths may be provided. |

**Examples**

Validate a single file:

```bash
len-cli validate examples/basic/helloworld/hello.l1
```

Validate an entire directory tree:

```bash
len-cli validate lang/l1
```

Validate multiple paths at once:

```bash
len-cli validate src/model.l1 src/contracts.l1
```

**Exit Codes**

| Code | Meaning |
|------|---------|
| `0` | All files parsed and validated without errors |
| `1` | One or more validation errors were found |
| `2` | Bad arguments or a filesystem error |

---

## Diagnostics

Diagnostics are written to **stderr**. Each line has the format:

```
<file>:<line>:<col>: <severity>: <message>
```

Example output:

```
src/model.l1:12:5: error: unresolved name 'Sorted'
src/model.l1:20:1: error: duplicate declaration 'sort_correct'
```

**Severity levels**

| Level | Meaning |
|-------|---------|
| `error` | The file is invalid. The validator exits with code `1`. |
| `warning` | The file is structurally valid but has a potential issue. |

---

## Validation Checks

`len-cli validate` performs the following checks:

| Check | Description |
|-------|-------------|
| **Parse** | The file is syntactically well-formed |
| **Duplicate declarations** | No two declarations share the same name in the same scope |
| **Unresolved imports** | Every `import` and `from … import` resolves to a file on disk |
| **Unresolved names** | Every identifier used in a formula refers to a declared `type`, `rel`, `const`, or `fn` |
| **Arity** | Every relation and function application uses the correct number of arguments |
| **Binder scope** | Variables bound by `forall`, `exists`, or `given` are used within their scope |
| **Quasi surface** | The body of a `quasi` block conforms to its declared style profile |

Validation does **not** check:

- Logical consistency or provability of specs
- Execution or evaluation of quasi bodies
- l2-level semantics

---

## Quasi Style Profiles

Style profiles are YAML files that describe what surface forms are allowed inside a `quasi` block.

The default style is `ProceduralAlgorithm`. It is looked up in:

```
doc/proposals/accepted/lip-0001-cli-parser-n-validation/
```

To use a named style in a quasi block:

```len
fn my_fn(x: T) -> y: T
    implements MyRel(x, y)
    quasi using style ProceduralAlgorithm:
        let result := x
        return result
```

---

## Shell Completion

Shell completion is not yet implemented. Track progress in the project TODO.

---

## Makefile Targets

| Target | Command | Description |
|--------|---------|-------------|
| `build` | `make build` | Compile `len-cli` to `bin/len` |
| `test` | `make test` | Run all Go unit tests |
| `clean` | `make clean` | Remove the `bin/` directory |
