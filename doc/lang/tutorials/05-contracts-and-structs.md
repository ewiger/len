# Tutorial 5 — Contracts and Structs

This tutorial introduces two surface-level organisational forms: `contract` and `struct`. Both are sugar over the relational core — they make models more readable without adding new semantic primitives.

---

## Contracts

A `contract` groups related `rel`, `fn`, and `spec` declarations under a single named proposition.

```len
contract Eq(T: Type)
    rel Equal(x: T, y: T)

    spec eq_reflexive
        given x: T
        must Equal(x, x)

    spec eq_symmetric
        given x, y: T
        must Equal(x, y) implies Equal(y, x)

    spec eq_transitive
        given x, y, z: T
        must Equal(x, y) and Equal(y, z) implies Equal(x, z)
```

Reading this: `Eq(T)` is a contract proposition. Any type `T` that satisfies `Eq` must come equipped with an `Equal` relation that obeys the three laws above.

### Contracts vs. Java interfaces

| Java Interface | len Contract |
|----------------|--------------|
| `interface Comparable<T>` | `contract Ord(T: Type)` |
| Method signature (dispatch) | `rel` declaration (logical proposition) |
| `@Override` implementation | `fn` linked via `implements` |
| No laws on the interface | `spec` clauses state the laws |
| Runtime dispatch | No dispatch — satisfaction is a logical statement |

The key difference: a Java interface is a runtime mechanism. A len contract is a logical proposition. `Ord(T)` does not create a vtable — it states that `T` satisfies a set of relational laws.

### A richer contract example

```len
contract Ord(T: Type)
    rel Leq(x: T, y: T)

    spec ord_reflexive
        given x: T
        must Leq(x, x)

    spec ord_antisymmetric
        given x, y: T
        must Leq(x, y) and Leq(y, x) implies x = y

    spec ord_transitive
        given x, y, z: T
        must Leq(x, y) and Leq(y, z) implies Leq(x, z)

    spec ord_total
        given x, y: T
        must Leq(x, y) or Leq(y, x)
```

---

## Structs

A `struct` declares a composite record type — sugar over `type` and field relations.

```len
struct Product
    name  : String
    price : Money
    sku   : String
```

This expands semantically to:

```len
type Product
rel ProductField_name(self: Product, value: String)
rel ProductField_price(self: Product, value: Money)
rel ProductField_sku(self: Product, value: String)
```

(The exact expansion form is implementation-defined; treat `struct` as a readable shorthand.)

### Adding specs to a struct

You can write ordinary specs that reference the struct's field relations:

```len
struct Product
    name  : String
    price : Money
    sku   : String

spec product_name_required
    must forall p: Product .
        exists name: String . ProductField_name(p, name)

spec product_price_positive
    must forall p: Product, m: Money .
        ProductField_price(p, m) implies Positive(m)
```

---

## Relation Refinement

A relation can refine another, stating that it is a semantically more specific version:

```len
rel Sort(input: Seq, output: Seq)

rel StableSort(input: Seq, output: Seq) refines Sort
rel BubbleSort(input: Seq, output: Seq) refines Sort
rel MergeSort(input: Seq, output: Seq) refines Sort
```

`BubbleSort` inherits all obligations of `Sort` and may add its own specs. Refinement is a logical relationship, not inheritance in the OOP sense.

---

## Combining Contracts, Structs, and Specs

A realistic domain model for an online shop:

```len
import core.math.nat

# --- Shared types ---

type String
type Nat
type Money
type Email

# --- Contracts ---

contract Named(T: Type)
    rel Name(self: T, value: String)
    spec name_unique
        given self: T, a, b: String
        must Name(self, a) and Name(self, b) implies a = b

contract Priced(T: Type)
    rel Price(self: T, amount: Money)
    spec price_unique
        given self: T, a, b: Money
        must Price(self, a) and Price(self, b) implies a = b
    spec price_positive
        given self: T, m: Money
        must Price(self, m) implies Positive(m)

# --- Structs ---

struct Product
    name  : String
    price : Money
    sku   : String

struct Customer
    name  : String
    email : Email

struct OrderItem
    product  : Product
    quantity : Nat

struct Order
    customer : Customer
    items    : OrderItem

# --- Cross-cutting specs ---

spec order_must_have_items
    must forall o: Order .
        exists item: OrderItem . OrderContains(o, item)

spec order_total_is_sum_of_items
    given o: Order, total: Money
    must OrderTotal(o, total) iff
        SumOfItemPrices(o, total)
```

---

## What You Learned

- `contract` groups `rel`, `fn`, and `spec` declarations under a named proposition
- Contracts are interface-like in appearance but carry no runtime dispatch
- `struct` is record sugar over `type` + field relations
- `refines` states semantic refinement between relations
- `satisfies` and `derives` are the bridges that link a type to a contract

**Next:** [Coming from Java →](06-from-java.md)
