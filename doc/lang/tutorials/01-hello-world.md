# Tutorial 1 — Hello World

This tutorial walks you through installing `len-cli`, writing your first `.l1` file, and running the validator.

---

## Prerequisites

- Go 1.21 or later (`go version` to check)
- Git

---

## Step 1 — Build the CLI

Clone the repository and build:

```bash
git clone <repo-url>
cd len
make build
```

The binary is written to `bin/len`. Add it to your path for convenience:

```bash
export PATH="$PATH:$(pwd)/bin"
```

Verify the build:

```bash
len-cli validate --help
# or
len-cli
# usage: len-cli validate <path-or-directory> [more paths]
```

---

## Step 2 — Write Your First File

Create a file called `hello.l1`:

```len
type String

rel Hello(output: String)

fn hello() -> output: String
    ensures Hello(output)
    quasi using style ProceduralAlgorithm:
        let greeting := "Hello, world!"
        return greeting
```

Let's break this down line by line.

### `type String`

Declares a type named `String`. In len, types are opaque names — you introduce their structure through relations.

### `rel Hello(output: String)`

Declares a relation named `Hello` that holds for a value of type `String`. This is the specification: we are saying there exists a notion of "Hello output".

### `fn hello() -> output: String`

Declares a function named `hello` that returns a value of type `String` named `output`.

### `ensures Hello(output)`

A postcondition: after `hello` runs, the relation `Hello(output)` must hold for the returned value.

### `quasi using style ProceduralAlgorithm:`

Opens the required pseudocode body. This is not executed — it is a structured sketch that the validator checks for surface well-formedness.

### `let greeting := "Hello, world!"`

Inside the quasi block: introduces a local variable `greeting`.

### `return greeting`

Returns the greeting.

---

## Step 3 — Validate

Run the validator:

```bash
len-cli validate hello.l1
```

If the file is valid, the command exits silently with code `0`.

Try introducing an error — for example, reference an undefined relation:

```len
fn hello() -> output: String
    ensures Missing(output)
    quasi:
        return "hi"
```

```bash
len-cli validate hello.l1
# hello.l1:2:13: error: unresolved name 'Missing'
```

The diagnostic tells you the file, line, and column of the problem.

---

## Step 4 — Validate a Directory

When a project grows into multiple files, you can validate the whole tree at once:

```bash
len-cli validate src/
```

All `.l1` files under `src/` are discovered and validated together, so cross-file imports resolve correctly.

---

## What You Learned

- `type` declares a named type
- `rel` declares a named relation — the primary modelling tool
- `fn` declares a function with a signature, contract clauses, and a required `quasi` body
- `ensures` states a postcondition
- `quasi using style ProceduralAlgorithm:` opens a procedural pseudocode sketch
- `len-cli validate` checks syntax, names, arity, and quasi structure

**Next:** [Types and Relations →](02-types-and-relations.md)
