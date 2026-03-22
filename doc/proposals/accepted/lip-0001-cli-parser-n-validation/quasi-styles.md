# quasi-styles


```len.l1

fn hello() -> output: String
    quasi using style ProceduralAlgorithm:
        let greeting := "Hello, "
        let name := "world!"
        output := greeting + name

```

## Validation of QUASI-BLOCK using many Quasi Styles

Here we describe how different style of pseudo-code that parser can encounter inside the `quasi` block can be syntactically validate without executing or going too deep.

We will call such process `surface validation` of `quasi` block. The main goal of this process is to make sure that the content of `quasi` block is syntactically correct and follows the expected structure, without necessarily understanding the semantics of the code or even parsing the specific grammar.

What is checked:

- keywords - each style has a set of expected keywords, and we can check if the content of the `quasi` block contains those keywords in the right places. For example, for `ProceduralAlgorithm` style, we expect to see keywords like `let`, `if`, `else`, `while`, etc. We can check if these keywords are used correctly.
- expression structure - we have a list of regular expressions. At least one of them should match the line of the quasi block.

Again, quasi block can be recognized as indented block with `QUASI_BLOCK` token, and then we can apply the surface validation on each line of the block. If any line fails the validation, we can report a syntax error.

## ProceduralAlgorithm Style

