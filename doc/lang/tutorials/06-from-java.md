# Tutorial 6 — Coming from Java

This tutorial is for developers who know Java well and want to understand how len's concepts map to familiar territory. len is not an OOP language, but many of its ideas have direct analogues — just with a different emphasis.

---

## The Fundamental Shift

Java asks: *how is this thing structured and how does it behave at runtime?*

len asks: *what is true about this thing and what laws must hold?*

This shift from *behaviour* to *proposition* is the key mental adjustment.

---

## Classes → `type` + `struct`

**Java:**

```java
public class Product {
    private String name;
    private BigDecimal price;
    private String sku;

    public Product(String name, BigDecimal price, String sku) { ... }
    public String getName() { return name; }
    public BigDecimal getPrice() { return price; }
    public String getSku() { return sku; }
}
```

**len:**

```len
struct Product
    name  : String
    price : Money
    sku   : String
```

A `struct` is a named record with typed fields. There are no constructors, no access modifiers, and no getters. The structure is exposed directly.

If you need a richer type without fields, use `type`:

```len
type OrderId    # opaque — no structure exposed
type Timestamp  # opaque
```

---

## Interfaces → `contract`

**Java:**

```java
public interface Comparable<T> {
    int compareTo(T other);
}

public interface Serializable { }
```

**len:**

```len
contract Ord(T: Type)
    rel Leq(x: T, y: T)

    spec ord_reflexive
        given x: T
        must Leq(x, x)

    spec ord_antisymmetric
        given x, y: T
        must Leq(x, y) and Leq(y, x) implies x = y

    spec ord_total
        given x, y: T
        must Leq(x, y) or Leq(y, x)
```

The crucial differences:

| Java Interface | len Contract |
|----------------|--------------|
| Method signatures (dispatch targets) | `rel` declarations (propositions) |
| No laws unless you write javadoc | `spec` clauses are formal laws |
| `instanceof` check at runtime | Contract satisfaction is a logical statement |
| Vtable / dynamic dispatch | No dispatch mechanism |

In Java, you satisfy an interface by providing a method body. In len, you satisfy a contract by writing specs and functions that match the contract's obligations.

---

## Abstract Classes → `contract` + shared `spec`

**Java:**

```java
public abstract class Shape {
    public abstract double area();

    public boolean isLargerThan(Shape other) {
        return area() > other.area();
    }
}
```

**len:**

```len
contract Shape(T: Type)
    rel Area(self: T, a: Real)

    spec area_positive
        given self: T, a: Real
        must Area(self, a) implies Positive(a)

    spec larger_than_def
        given a, b: T, sa, sb: Real
        must Area(a, sa) and Area(b, sb) implies
            (LargerThan(a, b) iff sa > sb)
```

---

## Methods → `rel` + `fn`

**Java:**

```java
public List<Product> findByCategory(Category cat) { ... }
```

**len** — two steps:

```len
# Step 1: declare what the method means
rel FindByCategory(cat: Category, results: Seq)

# Step 2: implement it
fn find_by_category(cat: Category) -> results: Seq
    implements FindByCategory(cat, results)
    ensures forall p: Product . results contains p implies ProductCategory(p, cat)
    quasi using style ProceduralAlgorithm:
        let out := empty_seq()
        for p in all_products():
            if has_category(p, cat):
                append(out, p)
        return out
```

The `rel` is the specification; the `fn` is the implementation.

---

## Generics → Parameterised `contract` and `type`

**Java:**

```java
public interface Repository<T, ID> {
    Optional<T> findById(ID id);
    List<T> findAll();
    T save(T entity);
}
```

**len:**

```len
contract Repository(T: Type, Id: Type)
    rel FindById(id: Id, result: T)
    rel FindAll(results: Seq)
    rel Save(entity: T, saved: T)

    spec find_by_id_unique
        given id: Id, a, b: T
        must FindById(id, a) and FindById(id, b) implies a = b
```

Type parameters in contracts work like Java generics — they are placeholders for any type that satisfies the contract.

---

## Annotations → `spec`

**Java:**

```java
@NotNull
@Size(min = 1, max = 100)
private String name;

@Positive
private BigDecimal price;
```

**len:**

```len
spec product_name_not_empty
    given p: Product, name: String
    must ProductName(p, name) implies not EmptyString(name)

spec product_name_length
    given p: Product, name: String, n: Nat
    must ProductName(p, name) and Length(name, n) implies n <= 100

spec product_price_positive
    given p: Product, m: Money
    must ProductPrice(p, m) implies Positive(m)
```

Where Java annotations are runtime-checked metadata, len specs are formal logical laws.

---

## Enums → `type` + `const` + `spec`

**Java:**

```java
public enum OrderStatus {
    PENDING, PAID, SHIPPED, CANCELLED
}
```

**len:**

```len
type OrderStatus

const Pending   : OrderStatus
const Paid      : OrderStatus
const Shipped   : OrderStatus
const Cancelled : OrderStatus

spec order_status_distinct
    must not (Pending = Paid)
      and not (Pending = Shipped)
      and not (Pending = Cancelled)
      and not (Paid = Shipped)
      and not (Paid = Cancelled)
      and not (Shipped = Cancelled)

spec order_status_exhaustive
    must forall s: OrderStatus .
        s = Pending or s = Paid or s = Shipped or s = Cancelled
```

---

## Exception / Validation Rules → `requires` + `ensures`

**Java:**

```java
public Order createOrder(Customer customer, List<OrderItem> items) {
    if (customer == null) throw new IllegalArgumentException("customer required");
    if (items.isEmpty()) throw new IllegalArgumentException("at least one item required");
    ...
}
```

**len:**

```len
rel CreateOrder(customer: Customer, items: Seq, order: Order)

fn create_order(customer: Customer, items: Seq) -> order: Order
    implements CreateOrder(customer, items, order)
    requires NonEmpty(items)
    ensures OrderCustomer(order, customer)
    ensures forall item: OrderItem . item in items implies OrderContains(order, item)
    quasi using style ProceduralAlgorithm:
        let o := new_order(customer)
        for item in items:
            append(o.items, item)
        return o
```

`requires` and `ensures` are formal contract clauses — stronger than runtime checks because they must be logically consistent with the model's specs.

---

## Packages / Namespaces → Modules

**Java:**

```java
package com.example.shop.catalogue;
import com.example.shop.types.*;
```

**len:**

```len
# In shop/catalogue.l1
import shop.types
```

Module paths map directly to file paths: `shop.types` → `shop/types.l1`.

---

## Quick Reference

| Java | len |
|------|-----|
| `class Product { String name; }` | `struct Product \n    name: String` |
| `interface Comparable<T>` | `contract Ord(T: Type)` |
| `abstract class Shape` | `contract Shape(T: Type)` |
| Method signature | `rel` declaration |
| Method body | `fn` with `quasi` block |
| Generic type parameter | Type parameter in `contract` |
| `@NotNull` annotation | `spec` with `not Null(...)` |
| `enum Status { A, B }` | `type` + `const` + exhaustiveness `spec` |
| `throws` / precondition check | `requires` clause on `fn` |
| Postcondition / return invariant | `ensures` clause on `fn` |
| `package` / `import` | `import module.path` |
| Javadoc comment | Triple-quoted docstring `"""..."""` |

---

## What You Learned

- Java classes map to `struct` (with fields) or `type` (opaque)
- Java interfaces map to `contract` — but contracts carry formal laws via `spec`
- Java methods map to `rel` (the specification) + `fn` (the implementation)
- Java generics map to parameterised `contract` and `type`
- Java annotations map to `spec` declarations
- Java enums map to `type` + `const` + exhaustiveness specs
- Java `throws` / validation maps to `requires` clauses

**Next:** [PostgreSQL Schema Generation →](07-postgres-schema-generation.md)
