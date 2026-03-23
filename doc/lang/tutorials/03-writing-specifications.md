# Tutorial 3 — Writing Specifications

A `spec` is a logical statement that constrains how relations behave. Specs are the semantic core of a len model — they express the laws, definitions, and invariants that must hold.

This tutorial shows you how to write specs for a domain model and use the full range of logical connectives.

---

## The `spec` Form

```len
spec <name>
    given <var>: <Type> [, ...]
    must <formula>
```

- `given` introduces universally-bound variables — they range over all values of their type.
- `must` states the formula that must hold.

A spec with no `given` clause states an unconditional axiom:

```len
spec false_is_not_true
    must not (True = False)
```

---

## Naming Conventions for Specs

Prefix the spec name to signal its role:

| Prefix | Role |
|--------|------|
| (no prefix) | General law or constraint |
| `axiom_*` | Fundamental assumption |
| `definition_*` | Introduces a new notion in terms of existing ones |
| `theorem_*` | Derivable from other specs |
| `lemma_*` | Auxiliary statement for a proof |
| `corollary_*` | Follows directly from a theorem |
| `postulate_*` | Local primitive assumption |

---

## Logical Connectives

### `and` — Conjunction

Both sides must hold:

```len
spec valid_order
    given o: Order
    must HasCustomer(o) and HasAtLeastOneItem(o)
```

### `or` — Disjunction

At least one side must hold:

```len
spec payment_method
    given o: Order
    must PaidByCreditCard(o) or PaidByBankTransfer(o)
```

### `not` — Negation

The formula must not hold:

```len
spec no_empty_order
    given o: Order
    must not EmptyOrder(o)
```

### `implies` — Implication

If the left holds, the right must hold too:

```len
spec shipped_implies_paid
    given o: Order
    must Shipped(o) implies Paid(o)
```

### `iff` — Biconditional

Both directions of implication — a definition-style law:

```len
spec definition_of_active
    given c: Customer
    must Active(c) iff
        HasEmail(c) and not Banned(c)
```

---

## Quantifiers

### `forall` — Universal

The formula must hold for every value:

```len
spec price_is_positive
    must forall p: Product, m: Money .
        HasPrice(p, m) implies Positive(m)
```

The `.` separates the binders from the formula body.

Multiple variables of the same type can be declared together:

```len
spec order_uniqueness
    must forall a, b: Order .
        SameId(a, b) implies a = b
```

### `exists` — Existential

At least one value must satisfy the formula:

```len
spec catalogue_not_empty
    must exists p: Product . InStock(p)
```

Combine `forall` and `exists` freely:

```len
spec every_order_has_a_customer
    must forall o: Order .
        exists c: Customer . OrderCustomer(o, c)
```

---

## Equality

Use `=` to state that two things are identical:

```len
spec customer_email_unique
    must forall c1, c2: Customer, e: Email .
        CustomerEmail(c1, e) and CustomerEmail(c2, e) implies c1 = c2
```

---

## Syntax Sugar

If you import `core.math.logic.sugar` you can use shorthand operators:

```len
import core.math.logic.sugar

spec de_morgan_1
    given A, B: Formula
    must !(A /\ B) <=> (!A \/ !B)
```

| Sugar | Keyword | Meaning |
|-------|---------|---------|
| `A => B` | `implies` | Implication |
| `A <=> B` | `iff` | Biconditional |
| `A /\ B` | `and` | Conjunction |
| `A \/ B` | `or` | Disjunction |
| `!A` | `not` | Negation |

---

## Building a Spec Suite for the Shop Model

Continuing from the shop model in [Tutorial 2](02-types-and-relations.md), add a `specs.l1` file:

```len
import shop.types
import shop.catalogue
import shop.orders

# Every product has exactly one name
spec product_name_unique
    must forall p: Product, a, b: String .
        ProductName(p, a) and ProductName(p, b) implies a = b

# Every product has exactly one price
spec product_price_unique
    must forall p: Product, a, b: Money .
        ProductPrice(p, a) and ProductPrice(p, b) implies a = b

# Every order must have a customer
spec order_has_customer
    must forall o: Order .
        exists c: Customer . OrderCustomer(o, c)

# An order may only be shipped if it has been paid
spec shipped_implies_paid
    given o: Order
    must Shipped(o) implies Paid(o)

# An order may only be paid if it has at least one item
spec paid_implies_has_items
    given o: Order
    must Paid(o) implies
        exists item: OrderItem . OrderContains(o, item)

# A cancelled order is neither paid nor shipped
spec cancelled_is_terminal
    given o: Order
    must Cancelled(o) implies
        not Paid(o) and not Shipped(o)
```

Run validation:

```bash
len-cli validate shop/
```

---

## What You Learned

- `spec` declares a logical law using `given` (binders) and `must` (the formula)
- `and`, `or`, `not`, `implies`, `iff` are the core connectives
- `forall` and `exists` express statements about all or some values
- `=` states identity
- Naming conventions (`definition_*`, `axiom_*`, etc.) communicate the role of each spec
- Sugar operators (`=>`, `/\`, etc.) are available through `core.math.logic.sugar`

**Next:** [Functions and Quasi →](04-functions-and-quasi.md)
