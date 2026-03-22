# Tutorial 2 — Types and Relations

In `len`, all modelling is built from two primitives: **types** and **relations**. This tutorial shows you how to think in those terms, introduce a domain, and wire it together with imports.

---

## Types

A `type` declaration introduces a named, opaque set of values.

```len
type Product
type Customer
type Order
type Money
```

Types carry no built-in structure. You cannot ask "does a Product have a name?" without also declaring a relation that captures that fact. This is intentional — it keeps the core language small and forces you to make relationships explicit.

---

## Relations

A `rel` declaration describes a property of, or relationship between, typed values.

```len
type Product
type String
type Money

rel HasName(p: Product, name: String)
rel HasPrice(p: Product, price: Money)
rel InStock(p: Product)
```

Relations replace what you might reach for methods or fields in an OOP language. `HasName(p, name)` is a proposition: it is either true or false for a given product and string. You do not define *how* it is true — only *that* it exists as a concept.

### Relations are not methods

In Java you might write:

```java
class Product {
    String getName() { ... }
}
```

In len you write:

```len
rel HasName(p: Product, name: String)
```

The difference: `HasName` makes no claim about uniqueness, about how the name is stored, or about calling conventions. It is a pure logical proposition.

---

## Constants

`const` declares a named value with a type:

```len
type Currency

const USD : Currency
const EUR : Currency
const GBP : Currency
```

Constants are useful for canonical sentinel values and named domain primitives.

---

## A Small Domain Model

Let's build a minimal e-commerce domain. Create `shop.l1`:

```len
# Types
type Product
type Customer
type Order
type OrderItem
type Money
type String
type Nat

# Product catalogue relations
rel ProductName(p: Product, name: String)
rel ProductPrice(p: Product, price: Money)
rel ProductSku(p: Product, sku: String)
rel InStock(p: Product)

# Customer relations
rel CustomerName(c: Customer, name: String)
rel CustomerEmail(c: Customer, email: String)

# Order relations
rel OrderCustomer(o: Order, c: Customer)
rel OrderItem(o: Order, item: OrderItem)
rel ItemProduct(item: OrderItem, p: Product)
rel ItemQuantity(item: OrderItem, qty: Nat)

# Money relations
rel MoneyAmount(m: Money, amount: Nat)
rel MoneyCurrency(m: Money, code: String)
```

This is a pure structural model. No algorithms, no storage decisions, no SQL yet — just the vocabulary of the domain.

---

## Imports

Larger projects split across files and use `import` to connect them.

Create `types.l1`:

```len
type Product
type Customer
type Order
type OrderItem
type Money
type String
type Nat
```

Create `catalogue.l1`:

```len
import shop.types

rel ProductName(p: Product, name: String)
rel ProductPrice(p: Product, price: Money)
rel ProductSku(p: Product, sku: String)
rel InStock(p: Product)
```

The module path `shop.types` maps to the file `shop/types.l1` on disk. Validate the project by pointing at the directory:

```bash
len-cli validate shop/
```

### Selective imports

If you only need a few names from a module:

```len
from shop.types import Product, Money
```

### Core library imports

The standard library ships under `core`:

```len
import core.math.set      # Set theory primitives
import core.math.nat      # Natural numbers
import core.math.logic    # Bool, Formula, Holds
```

---

## Naming Conventions

| Form | Used for |
|------|----------|
| `PascalCase` | Types, relations, constants |
| `snake_case` | Specs, functions, local variables |
| `dot.separated` | Module paths |

---

## What You Learned

- `type` names an opaque set of values — it carries no structure by itself
- `rel` declares a logical proposition about typed values — the workhorse of len
- `const` names a fixed value of a type
- `import` and `from … import` let you split a model across files
- Module paths are dot-separated and map directly to file paths

**Next:** [Writing Specifications →](03-writing-specifications.md)
