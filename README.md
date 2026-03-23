# len

`len` is a meta programming language.

It is intended for writing programs that generate, transform, or reason about other programs and code structures.

📖 **[Documentation →  ewiger.github.io/len](https://ewiger.github.io/len/)**

## Hello World

```len
type String
rel Hello(output: String)

fn hello() -> output: String
    ensures Hello(output)
    quasi using style ProceduralAlgorithm:
        let greeting := "Hello, world!"
        return greeting
```

Validate it with the CLI:

```bash
go run ./cmd/len-cli validate examples/basic/helloworld/hello.l1
```

## Language Layers

Each layer has a distinct purpose and design focus:

| Layer | Name | Purpose | Typical contents |
|------|------|---------|------------------|
| **L0** | Natural spec | Capture human intent before formalization. | Goals, examples, rationale, edge cases, sketches, informal pseudo-code |
| **L1** | Structural core | Define the canonical semantic model. | Types, relations, contracts, laws, function obligations, examples |
| **L2** | Generation contract | Connect the abstract model to concrete output targets. | Package layout, naming rules, AST shapes, runtime choices, error strategy, performance requirements |

- **Level 0** or `len.l0` remains in natural language and informal pseudo-code. It is intended for high-level design, brainstorming, communication, and early exploration.
- **Level 1** or `len.l1` defines the canonical structural core of len: types, relations, formulas, and syntax-level logical constructs. It is intentionally lightweight and declarative.
- **Level 2** or `len.l2` builds on this core by introducing interpretation, evaluation, and target-aware generation guidance.

Taken together, the layers support a specification-first workflow:

1. Write **L0** to clarify intent, examples, and edge cases.
2. Write **L1** to formalize the model, relations, and laws.
3. Write **L2** to describe how that model maps into executable systems or generated code.
4. Use **L0 + L1 + L2** together when driving code generation with LLMs or other tooling.
5. Validate the generated result against tests and, where possible, against the L1 obligations.

## the Big Picture

The top down loop is:

1. Write L0 design notes and examples to clarify intent and edge cases.
2. Write L1 declarations to capture the core model, relations, and laws.
3. Write L2 generation contracts to connect the abstract model to concrete code structures.
4. Use LLMs to generate code from the L2 contracts, guided by the L1 model and L0 intent.
5. Validate the generated code against tests and, where possible, against the L1 obligations.

However it is also possible to move in both directions and experiment: 

- *interpret* to impose semantical and relational interpretation on the structures, types, relations
- *verbalize*  (the opposite of interpret) to enrich previous or next layer
- *interpolate* to fill in gaps or connect layers
- *extrapolate* to generalize and extend the model or generation contracts
- *verify* to check logical consistency and correctness
- *iterate* to refine and improve
- *code generate* to produce executable code
- *test* to validate functionality
- *document* to clarify and communicate
- *evolve* to adapt to new requirements or insights through many experimental pipelines, selection and refinement rounds

This changes the workflow or even the whole software lifecycle from a code-centric one to a specification-first one, where the specification is the primary artifact and the code is generated from it.

## Status

> This project is **early-stage and experimental**.

| Area | Status | Notes |
|------|--------|-------|
| Project | ![Experimental](https://img.shields.io/badge/status-experimental-orange) | Expect breaking changes to syntax, semantics, and tooling. |
| L0 | ![Exploratory](https://img.shields.io/badge/L0-exploratory-lightgrey) | Natural-language design space for intent, rationale, and examples. |
| L1 | ![Most stable](https://img.shields.io/badge/L1-most%20stable-brightgreen) | The most complete and stable layer today. |
| L2 | ![In progress](https://img.shields.io/badge/L2-in%20progress-blue) | Actively being designed and implemented; not yet usable. |
| CLI | ![MVP](https://img.shields.io/badge/cli-validator%20mvp-teal) | Go implementation validates `.l1` files structurally. |
| Docs | ![GitHub Pages](https://img.shields.io/badge/docs-GitHub%20Pages-6f42c1) | Static site generated with MkDocs and deployed with GitHub Actions. |

- Expect breaking changes to syntax, semantics, and tooling.
- Level 0 is natural language and is not expected to have a formal syntax or semantics, but the design may evolve as we explore how to best support informal design and communication.
- Level 1 is the most complete and stable layer.
- Level 2 is actively being designed and implemented, but it is not yet in a usable state.
- The CLI has an MVP Go implementation that can structurally validate `.l1` files, but it does not yet execute code, interpret formulas, or perform proof checking.

# Documentation

Language is still being designed and implemented, but given that this is a meta-programming language, the documentation it can be productive for code generation already now. 

Documentation is available in three forms:

- Published website: [ewiger.github.io/len](https://ewiger.github.io/len/)
- Docs source in this repository: [doc/lang/index.md](doc/lang/index.md)
- Generated static site output: [doc/lang-html](doc/lang-html)

The website is generated from `doc/lang/**` with MkDocs and deployed as static GitHub Pages via GitHub Actions.

For a quick start, check out the [concepts guide](doc/lang/concepts.md) and the [Hello World tutorial](doc/lang/tutorials/01-hello-world.md).


## Developer

The repository now includes the first Go implementation for the accepted CLI parser and validator milestone.

Current implementation:

- `len-cli validate` for parsing and validating `.l1` files
- hand-written lexer and parser for the current MVP declaration surface
- semantic validation for duplicate declarations, unresolved imports, binder scope, unresolved names, and function or relation arity
- quasi surface validation through YAML style profiles

Build the CLI:

```bash
make build
```

This writes the CLI binary to `bin/len`.

Run all unit tests:

```bash
make test
```

Build the language docs into `doc/lang-html`:

```bash
make docs
```

Run a live preview server for `doc/lang/**`:

```bash
make docs-serve
```

If `mkdocs` is not installed yet:

```bash
python3 -m pip install -r requirements-docs.txt
```

To publish the docs with GitHub Pages, enable **Settings > Pages > Build and deployment > Source: GitHub Actions**.
The workflow in `.github/workflows/docs-pages.yml` will build the site from `doc/lang/**` and deploy it on pushes to `main` or `master`, and it can also be run manually.

Run the validator on the hello world example:

```bash
go run ./cmd/len-cli validate examples/basic/helloworld/hello.l1
```

Current scope limits:

- `lang/l1/**` is intentionally not the acceptance target yet
- quasi bodies are captured as raw indented lines and surface-validated, not parsed into a dedicated internal statement AST
- the tool validates structure only; it does not execute code, interpret formulas, or perform proof checking
