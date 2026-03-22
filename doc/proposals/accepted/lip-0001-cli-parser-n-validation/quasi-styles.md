# quasi-styles


```len.l1

fn hello() -> output: String
    quasi using style ProceduralAlgorithm:
        let greeting := "Hello, "
        let name := "world!"
        set output := greeting + name

```

## Validation of QUASI-BLOCK using many Quasi Styles

Here we describe how different style of pseudo-code that parser can encounter inside the `quasi` block can be syntactically validate without executing or going too deep.

We will call such process `surface validation` of `quasi` block. The main goal of this process is to make sure that the content of `quasi` block is syntactically correct and follows the expected structure, without necessarily understanding the semantics of the code or even parsing the specific grammar.

What is checked:

- keywords - each style has a set of expected keywords, and we can check if the content of the `quasi` block contains those keywords in the right places. For example, for `ProceduralAlgorithm` style, we expect to see keywords like `let`, `if`, `else`, `while`, etc. We can check if these keywords are used correctly.
- expression structure - we have a list of regular expressions. At least one of them should match the line of the quasi block.

Again, quasi block can be recognized as indented block with `QUASI_BLOCK` token, and then we can apply the surface validation on each line of the block. If any line fails the validation, we can report a syntax error.

In the accepted MVP direction, the host grammar does not parse those inner lines into a dedicated quasi AST. It only captures the raw block and the optional style name from the quasi header. Style-specific surface validation happens later, after parse.

## ProceduralAlgorithm Style

The current sorting examples under [examples/quasi/sorting/bubble_sort.l1](/Users/yy/code/len/len-feat-cli-l1-validation/examples/quasi/sorting/bubble_sort.l1), [examples/quasi/sorting/insertion_sort.l1](/Users/yy/code/len/len-feat-cli-l1-validation/examples/quasi/sorting/insertion_sort.l1), [examples/quasi/sorting/merge_sorting.l1](/Users/yy/code/len/len-feat-cli-l1-validation/examples/quasi/sorting/merge_sorting.l1), and [examples/quasi/sorting/quick_sort.l1](/Users/yy/code/len/len-feat-cli-l1-validation/examples/quasi/sorting/quick_sort.l1) are a good corpus for a first explicit style profile.

The sibling YAML file [procedural-algorithm.quasi-style.yaml](/Users/yy/code/len/len-feat-cli-l1-validation/doc/proposals/accepted/lip-0001-cli-parser-n-validation/procedural-algorithm.quasi-style.yaml) captures the observed surface contract for `quasi using style ProceduralAlgorithm:`.

That profile intentionally stays at the level of surface validation.

- it fixes the allowed leading keywords
- it lists the accepted line schemas as regular expressions
- it distinguishes simple statements from block-opening statements
- it encodes continuation rules for `else` and `else if`
- it forbids free-form narrative lines for this style

### Why This YAML Shape

The sorting corpus uses a stable set of statement families:

- `let name := expr`
- `set target := expr`
- `append expr to target`
- `return expr`
- `if formula:`
- `else if formula:`
- `else:`
- `while formula:`
- `for name in expr:`

That is enough structure to validate the examples without fully interpreting the algorithm.

The YAML separates three concerns:

1. lexical expectations for line starts via `keywords`
2. shape validation via regex-based `rules`
3. structural constraints via `validation` and block metadata such as `opensBlock`, `attachesTo`, and `mustAlignWithParent`

### Suggested Go Validation Algorithm

For a Go CLI, the cleanest design is a two-phase validation path after the host parser has already captured the raw quasi block.

1. Parse the enclosing `fn` normally and capture each quasi block as raw lines with indentation metadata.
2. Read the style name from the clause header. For `quasi:` use the configured default profile. For `quasi using style ProceduralAlgorithm:` load the matching profile.
3. Unmarshal the YAML profile with `gopkg.in/yaml.v3` into a typed struct such as `QuasiStyleProfile`.
4. Normalize each raw quasi line into:
    - original text
    - trimmed text
    - indentation depth in spaces
    - source position
5. Reject lines that violate the indentation mode or width from the profile.
6. For each non-blank line, detect the leading keyword and reject it early if it is not allowed by the style.
7. Match the full trimmed line against the ordered rule list for that keyword.
8. Build a lightweight block stack while scanning lines top to bottom:
    - push `if`, `while`, and `for` rules that open a block
    - require the next nested line to be indented one level deeper
    - when indentation decreases, pop completed blocks
    - when `else` or `else if` appears, ensure it attaches to the most recent compatible `if` family block at the same indentation level
9. Optionally validate captured slots with host helpers:
    - identifiers with a simple identifier checker
    - formulas and expressions either with permissive regex acceptance in MVP or by delegating to the host expression/formula parser later
10. Emit diagnostics with file, line, column, rule id, and style name so users can see whether a failure is lexical, structural, or slot-related.

### Practical Go Data Model

The YAML can map directly to structs of this shape:

```go
type QuasiStyleProfile struct {
     Version    int                `yaml:"version"`
     Style      StyleMeta          `yaml:"style"`
     Layout     LayoutRules        `yaml:"layout"`
     Keywords   KeywordGroups      `yaml:"keywords"`
     Slots      map[string]SlotDef `yaml:"slots"`
     Rules      []RuleDef          `yaml:"rules"`
     Validation ValidationRules    `yaml:"validation"`
}
```

Then compile all regexes once when the profile is loaded, not on every line.

### Validation Scope

This profile is deliberately narrower than a full quasi language.

- it is intended to accept the current sorting examples
- it does not try to prove semantic correctness
- it does not require runtime evaluation
- it only guarantees that a block looks like a valid ProceduralAlgorithm block on the surface

That makes it appropriate for `len-cli validate` in the current Go MVP.

