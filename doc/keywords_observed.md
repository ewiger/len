# Observed Keywords In `lang/**`

This document is a snapshot of the words and symbols currently observed in `lang/**/*.l1`.
It is descriptive, not normative.

The goal is to separate:

- words already treated as reserved keywords
- words that behave like syntax in practice
- symbols and operators used by the current surface language
- identifiers that appear in source but should not be promoted to keywords

## Declared Reserved Keywords

These are explicitly declared in [lang/core/syntax/syntax.l1](/Users/yy/code/len/len-feat-init-core/lang/core/syntax/syntax.l1):

- `rel`
- `type`
- `syntax`
- `implies`
- `iff`
- `import`
- `keyword`
- `trait`
- `impl`
- `def`
- `spec`
- `given`
- `must`

## Observed Logic And Binder Words

These words appear inside formulas or spec bodies:

- `forall`
- `exists`
- `and`
- `or`
- `not`
- `implies`
- `iff`

## Observed Surface Syntax Words

These are not declared via `keyword` today, but they are used as part of surface notation:

- `in`
- `where`
- `from`
- `subsetof`

Notes:

- `in` is used as an infix membership form in [lang/core/math/set/set.l1](/Users/yy/code/len/len-feat-init-core/lang/core/math/set/set.l1).
- `where` is used in syntax declarations in [lang/core/math/set/set.l1](/Users/yy/code/len/len-feat-init-core/lang/core/math/set/set.l1) and [lang/core/math/set/nat.l1](/Users/yy/code/len/len-feat-init-core/lang/core/math/set/nat.l1).
- `from` appears in typed import syntax in [lang/core/math/logic/logic.l1](/Users/yy/code/len/len-feat-init-core/lang/core/math/logic/logic.l1).
- `subsetof` currently behaves more like user-facing surface syntax than a reserved core word.

## Observed Operators And Symbolic Tokens

These symbolic forms are present in `lang/**` source today:

- `:`
- `.`
- `=`
- `=>`
- `|s|`

There are also structural punctuation characters used in declarations and applications, including parentheses and commas.

## Observed But Not Keywords

These appear in source, but should not be treated as language keywords:

- type names such as `Set`, `Nat`, `Expr`, `Object`, `fn`
- relation and definition names such as `Member`, `Subset`, `Successor`, `BijectiveFun`
- spec labels such as `empty_def`, `pairing`, `zero_exists`
- module path segments such as `core`, `math`, `set`
- local binders such as `a`, `b`, `x`, `y`, `n`, `s`, `f`
- surface aliases such as `succ`

## Current Practical Dump

If the goal is a compact working keyword inventory for the current `l1` corpus, the strongest observed set is:

- `type`
- `rel`
- `syntax`
- `implies`
- `import`
- `keyword`
- `trait`
- `impl`
- `def`
- `spec`
- `given`
- `must`
- `forall`
- `exists`
- `and`
- `or`
- `not`
- `iff`
- `where`
- `from`
- `in`

## Gaps Worth Resolving

The current sources expose a few design questions:

- Should `forall`, `exists`, `and`, `or`, and `not` be explicitly declared as reserved keywords in the core syntax file?
- Should `where`, `from`, and `in` be treated as reserved words or only as grammar-level forms?
- Should `=>` replace `implies` in `l1`, or should both remain with distinct roles?
- Should `subsetof` stay user-definable syntax sugar rather than become a reserved word?

For now, this file records observed usage only.