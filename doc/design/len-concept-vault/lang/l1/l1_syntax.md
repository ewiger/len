# `l1` Syntax

`l1` is the first real formal layer. But it is still designed for humans, so it keeps the syntax local, contract-oriented, and readable.

The stable core of `l1` is:

- declarations first
- specs as the main semantic unit
- formulas in lightweight FOL-style syntax
- optional quasi steps for proof hints or implementation steering

### `l1` Design Constraint

`l1` should remain relational and contract-first.
It should not collapse into ordinary pseudocode.

This means:

- `rel` usage in `contract` with `spec` remains primary
- `fn` is secondary and usually interpreted (remains open for interpretation)
- `quasi` is allowed for steering, not for replacing the contract language
- `contract` comes with documentation and reasoning aid except for rigorous `spec`


### `l1` Scope

`l1` should cover:

- types and signatures
- relations and helper functions
- contracts, policies, queries, and validation rules through `spec`
- local quantification
- witness-oriented statements
- partially constructive quasi blocks


### `l1` Comments and Docstrings

Basically they are ignored by the formal semantics, but they can be used for documentation, reasoning, and communication purposes. They can be attached to any element in `l1`, such as types, relations, functions, specs, contracts, policies, queries, and validation rules.


Furthermore, as we come from l0 natural language specifications, we can have a Empty keywords or reserved words in which can be droped without changing the semantics or formal specifications.

- `a` is empty
- `the` is empty
- `is` is empty
- `are` is empty
- `of` is empty
- `to` is empty
- `for` is empty ??? or is it a reserved word for quantification? TBC
- `each` is empty

However, these words are important syntactic sugar:

- `be` or `is` or `are` means `equals` 
- `contains` or `has` means `in`

Once normalized, these words can be dropped without changing the semantics, but they can be used for readability and for making the specifications more natural and intuitive for humans.

### `l1` Logical Formulas
Logical formulas in `l1` should be lightweight and FOL-style, with support for quantifiers, connectives, and basic predicates. They should be readable and intuitive for humans, while still being precise enough for formal reasoning.

- `and` is *conjunction*
- `or` is *disjunction*
- `not` is *negation*

therea also quantifiers:
- `forall` is *universal quantification*
- `exists` is *existential quantification*

plus set-theoretic predicates:

- `in` is *membership*
- `subset` is *subset*
- `union` is *union*
- `intersection` is *intersection*
- `difference` is *set difference*
- `empty` is *empty set*
- `superset` is *superset*
- `subseteq` is *subset or equal*
- `superseteq` is *superset or equal*


### `l1` Operators And Reserved Forms (Symbols)

These are part of the concrete syntax even if they are not word-keywords.

- `->` is `implies`
- `<->`  is `iff`
- `=` is `equals`
- `!=` is `not equals`

- `<` is reserved and stands for *less than*, but it can be declared as a symbol for other meanings in the language, so it is not reserved for that meaning and must be explicitly declared as a symbol in the language. The same applies to the other comparison operators.
- `<=` is *less than or equal*
- `>` is *greater than*
- `>=` is *greater than or equal*
- `:` is used for type annotations and for separating variable names from their types in quantifiers
- `:=` is `decl` 
- `.` is used for member access and for separating namespaces or for logical formula scoping
- `,` is used for separating arguments in relations and functions, and for separating conjuncts in formulas
- `(` and `)` are used for grouping and for function application
- `[` and `]` are used for array indexing and for denoting sets or lists

composition, application, and abstraction?

list is an ordered set, array is a fixed-size collection of elements, set is an unordered collection of unique elements, tuple is an ordered collection of fixed number of elements.


All of the above are syntactic sugar for the underlying logical meaning, but they are not reserved for that meaning and must be explicitly declared as symbols in the language.

> arrays? square brackets? parentheses? commas?

declared as symbols in the language, but they are reserved for their logical meaning.



### Recommended Core `l1` Keywords

These should be treated as the main `l1` keyword set.

1) namespace, module, import, export for structuring and modularity:
- `module`
- `import` and `import as` for importing other modules with optional aliasing
- `export` is optional, but if defined, then not listed elemetns are not exported by default, otherwise all elements are exported by default.

2) declaration keywords for declaring the main concepts of the domain:
- `decl` 
- `symbol` helps with syntactic sugar for declaring symbols without a fixed signature, which can be useful for sketching out ideas or for domains that are not well understood yet. It can also be used for declaring a symbol of discourse for quantification.
- `domain` allows for open-ended declarations of symbols and relations without a fixed signature, which can be useful for sketching out ideas or for domains that are not well understood yet. It can also be used for declaring a domain of discourse for quantification.

- `type`
- `rel`
- `fn`

-  `contract` is a synonym for `open spec`, but is a form of documentation. It is more open ended suggestive of the intended meaning in `l1` which is to write a contract that can be later refined into an implementation or proof. It can still contain verbal or natural language parts as doc strings, but the main point is that it is a contract that can be refined into an implementation or proof.

    - `goal` (optional)
    - `context` (optional)
    - `given` (optional)
    - `when` (optional)
    - `then` (optional)
    - `because` (optional)
    - `example` (optional)
    - `counterexample` (optional)
    - `ambiguity` (optional)
    - `question` (optional)
    - `note` (optional)

    - `spec`
        - `requires`
        - `ensures <contraint_name>: <logical formula>`

            is followed by name of the constrint and a logical formula that states the constraint. The logical formula can be a simple formula or a complex formula with quantifiers and connectives.

            there can be multiple `ensures` clauses for a given `spec`, and they can be named for reference in proofs or implementations, but they are not required to be named. Importantly, the `ensures` clauses are joined into a big logical AND clause and must all hold together, so they are not independent constraints but rather parts of a single contract that must be satisfied as a whole. This allows for more modular and compositional reasoning about the contract, as well as for better error reporting and debugging when the contract is violated.

early logical formulas coem with the following keywords:
- `forall`
- `exists`
- `and`
- `or`
- `not`


- `open` is synonym for abstract, but it is more suggestive of the intended meaning in `l1` which is to open a contract or policy for later filling in of details, proof hints, or implementation steering.

- `quasi`
    - `impl` or `proof` for proof hints or implementation steering, but not for replacing the contract language. It can be used for sketching out ideas or for domains that are not well understood yet. It can also be used for declaring a proof sketch or an implementation sketch for a given contract.

    can use 
        - `choose`
        - `let`
        - `show`
        - `assume`
        - `witness`
        - `case`
        - `induction`
        - `then`

- `choose` is *witness selection* for existential quantification, but it is not a replacement for `exists` and should be used for sketching out ideas or for domains that are not well understood yet. It can also be used for declaring a witness for a given existential quantification.
- `let` is for defining local variables or for introducing intermediate steps in a proof or implementation, but it is not a replacement for quantifiers and should be used for sketching out ideas or for domains that are not well understood yet. It can also be used for declaring a local variable or an intermediate step for a given contract or policy.
- `show` is for stating a goal or a subgoal in a proof or implementation, but it is not a replacement for `ensures` and should be used for sketching out ideas or for domains that are not well understood yet. It can also be used for declaring a goal or a subgoal for a given contract or policy.


union, intersection, set difference, subset, membership?

true and false as keywords or symbols?


### Example `l1`

```len
type Expr
type Op
type Int

rel Const(Expr, Int)
rel Binary(Expr, Op, Expr, Expr)
rel AddOp(Op)
rel DivOp(Op)
rel Eval(Expr, Int)

fn add(a: Int, b: Int) -> Int

spec eval_const(e: Expr) -> v: Int
	requires Const(e, v)
	ensures Eval(e, v)

spec eval_add(e: Expr, op: Op, l: Expr, r: Expr) -> v: Int
	requires Binary(e, op, l, r)
	requires AddOp(op)
	requires exists a: Int. Eval(l, a)
	requires exists b: Int. Eval(r, b)
	ensures exists a: Int. exists b: Int.
		Eval(l, a) and Eval(r, b) and v = add(a, b)
	ensures Eval(e, v)

spec no_div_by_zero(e: Expr, op: Op, l: Expr, r: Expr) -> ok: Int
	requires Binary(e, op, l, r)
	requires DivOp(op)
	ensures exists b: Int. Eval(r, b) and b != 0
	open DivisionResultPolicy
```
## See Also

- [[program]] for the big picture of how `l1` fits into the overall language design and its relationship with other layers.