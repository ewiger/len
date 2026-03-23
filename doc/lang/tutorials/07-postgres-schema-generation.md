# Tutorial 7 — Generating a PostgreSQL Schema for an Online Shop

This tutorial is the most advanced in the series. You will:

1. Model a complete online shop domain in `len.l1`
2. Understand how `fn` + `quasi` express a code generator
3. Walk through a generator that emits `CREATE TABLE` SQL for PostgreSQL
4. See how the `rel` + `spec` contract ensures the generated schema is semantically consistent

This is a showcase of len as a **meta-programming language**: the generator itself is written in len, and it operates on len types and relations as its input.

---

## Part 1 — The Domain Model

Create the following files under `shop/`:

### `shop/types.l1`

```len
# Primitive value types
type String
type Nat
type Bool

# Domain types — opaque identifiers
type ProductId
type CustomerId
type OrderId
type OrderItemId
type CategoryId
type AddressId
```

### `shop/money.l1`

```len
import shop.types

type Money
type CurrencyCode

const USD : CurrencyCode
const EUR : CurrencyCode
const GBP : CurrencyCode

rel MoneyAmount(m: Money, amount: Nat)
rel MoneyCurrency(m: Money, code: CurrencyCode)

spec money_amount_unique
    given m: Money, a, b: Nat
    must MoneyAmount(m, a) and MoneyAmount(m, b) implies a = b

spec money_currency_unique
    given m: Money, a, b: CurrencyCode
    must MoneyCurrency(m, a) and MoneyCurrency(m, b) implies a = b
```

### `shop/catalogue.l1`

```len
import shop.types
import shop.money

struct Product
    id       : ProductId
    name     : String
    sku      : String
    price    : Money
    category : CategoryId

struct Category
    id   : CategoryId
    name : String
    slug : String

# A product belongs to exactly one category
spec product_has_category
    must forall p: Product .
        exists c: Category . ProductCategory(p, c)

# SKUs are unique across the catalogue
spec sku_unique
    must forall a, b: Product, s: String .
        ProductSku(a, s) and ProductSku(b, s) implies a = b

rel ProductCategory(p: Product, c: Category)
rel ProductSku(p: Product, sku: String)
rel InStock(p: Product)
rel StockQuantity(p: Product, qty: Nat)
```

### `shop/customers.l1`

```len
import shop.types

struct Customer
    id    : CustomerId
    name  : String
    email : String

struct Address
    id         : AddressId
    customer   : CustomerId
    street     : String
    city       : String
    country    : String
    postalCode : String

# Every customer has a unique email
spec customer_email_unique
    must forall a, b: Customer, e: String .
        CustomerEmail(a, e) and CustomerEmail(b, e) implies a = b

rel CustomerEmail(c: Customer, email: String)
rel CustomerAddress(c: Customer, addr: Address)
```

### `shop/orders.l1`

```len
import shop.types
import shop.money
import shop.catalogue
import shop.customers

type OrderStatus

const StatusPending   : OrderStatus
const StatusPaid      : OrderStatus
const StatusShipped   : OrderStatus
const StatusCancelled : OrderStatus

struct Order
    id       : OrderId
    customer : CustomerId
    status   : OrderStatus

struct OrderItem
    id       : OrderItemId
    order    : OrderId
    product  : ProductId
    quantity : Nat
    unitPrice : Money

rel OrderStatus(o: Order, s: OrderStatus)
rel OrderContains(o: Order, item: OrderItem)
rel OrderTotal(o: Order, total: Money)

spec order_has_customer
    must forall o: Order .
        exists c: Customer . OrderBelongsTo(o, c)

spec shipped_implies_paid
    given o: Order
    must OrderStatus(o, StatusShipped) implies
        exists prev: OrderStatus . OrderStatus(o, StatusPaid)

spec order_total_covers_items
    given o: Order, total: Money
    must OrderTotal(o, total) iff
        SumOfItemPrices(o, total)

rel OrderBelongsTo(o: Order, c: Customer)
rel SumOfItemPrices(o: Order, total: Money)
```

---

## Part 2 — The SQL Schema Generator

The generator is itself a len program. It takes the domain model as input — represented as relations over a meta-type `SqlTable` — and produces `CREATE TABLE` statements.

### `codegen/sql_types.l1`

```len
# Meta-types for SQL generation
type SqlTable
type SqlColumn
type SqlType
type SqlConstraint
type SqlScript

# SQL primitive type constants
const SqlText      : SqlType
const SqlInteger   : SqlType
const SqlNumeric   : SqlType
const SqlBoolean   : SqlType
const SqlTimestamp : SqlType
const SqlUuid      : SqlType

# Table structure relations
rel TableName(t: SqlTable, name: String)
rel TableColumn(t: SqlTable, col: SqlColumn)
rel ColumnName(col: SqlColumn, name: String)
rel ColumnType(col: SqlColumn, ty: SqlType)
rel ColumnNullable(col: SqlColumn, nullable: Bool)
rel ColumnDefault(col: SqlColumn, expr: String)
rel PrimaryKey(t: SqlTable, col: SqlColumn)
rel ForeignKey(col: SqlColumn, ref_table: SqlTable, ref_col: SqlColumn)
rel UniqueConstraint(t: SqlTable, col: SqlColumn)
rel NotNullConstraint(t: SqlTable, col: SqlColumn)
```

### `codegen/schema_gen.l1`

```len
import codegen.sql_types
import shop.types

# The top-level relation: a set of tables forms a schema
rel Schema(tables: Seq)

# The generator relation: given domain structs, produce a schema
rel GenerateSchema(domain_structs: Seq, schema: SqlScript)

fn generate_schema(domain_structs: Seq) -> schema: SqlScript
    implements GenerateSchema(domain_structs, schema)
    ensures WellFormed(schema)
    quasi using style ProceduralAlgorithm:
        let lines := empty_seq()

        for table in domain_structs:
            set lines := append_all(lines, emit_table(table))

        return join_lines(lines)

# Per-table emission
rel EmitTable(table: SqlTable, sql: Seq)

fn emit_table(table: SqlTable) -> sql: Seq
    implements EmitTable(table, sql)
    quasi using style ProceduralAlgorithm:
        let name := table_name(table)
        let cols := columns(table)
        let lines := empty_seq()

        append(lines, "CREATE TABLE " + name + " (")

        for col in cols:
            let col_def := emit_column(col)
            if not last(col, cols):
                set col_def := col_def + ","
            append(lines, "    " + col_def)

        append(lines, ");")
        append(lines, "")

        return lines

# Per-column emission
rel EmitColumn(col: SqlColumn, definition: String)

fn emit_column(col: SqlColumn) -> definition: String
    implements EmitColumn(col, definition)
    quasi using style ProceduralAlgorithm:
        let name := column_name(col)
        let ty   := sql_type_name(column_type(col))
        let def  := name + " " + ty

        if not column_nullable(col):
            set def := def + " NOT NULL"

        if has_default(col):
            set def := def + " DEFAULT " + column_default(col)

        return def
```

---

## Part 3 — The Generated Output

Given the shop domain, the generator produces SQL like this:

```sql
CREATE TABLE categories (
    id   UUID        NOT NULL,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    CONSTRAINT categories_pkey PRIMARY KEY (id),
    CONSTRAINT categories_slug_unique UNIQUE (slug)
);

CREATE TABLE products (
    id          UUID           NOT NULL,
    name        VARCHAR(255)   NOT NULL,
    sku         VARCHAR(100)   NOT NULL,
    price_amount NUMERIC(12,2) NOT NULL,
    price_currency CHAR(3)     NOT NULL DEFAULT 'USD',
    category_id UUID           NOT NULL,
    CONSTRAINT products_pkey         PRIMARY KEY (id),
    CONSTRAINT products_sku_unique   UNIQUE (sku),
    CONSTRAINT products_category_fk  FOREIGN KEY (category_id)
                                     REFERENCES categories (id)
);

CREATE TABLE customers (
    id    UUID         NOT NULL,
    name  VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    CONSTRAINT customers_pkey         PRIMARY KEY (id),
    CONSTRAINT customers_email_unique UNIQUE (email)
);

CREATE TABLE addresses (
    id          UUID         NOT NULL,
    customer_id UUID         NOT NULL,
    street      VARCHAR(255) NOT NULL,
    city        VARCHAR(100) NOT NULL,
    country     CHAR(2)      NOT NULL,
    postal_code VARCHAR(20)  NOT NULL,
    CONSTRAINT addresses_pkey        PRIMARY KEY (id),
    CONSTRAINT addresses_customer_fk FOREIGN KEY (customer_id)
                                     REFERENCES customers (id)
);

CREATE TABLE orders (
    id          UUID      NOT NULL,
    customer_id UUID      NOT NULL,
    status      VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at  TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT orders_pkey        PRIMARY KEY (id),
    CONSTRAINT orders_customer_fk FOREIGN KEY (customer_id)
                                  REFERENCES customers (id)
);

CREATE TABLE order_items (
    id              UUID           NOT NULL,
    order_id        UUID           NOT NULL,
    product_id      UUID           NOT NULL,
    quantity        INTEGER        NOT NULL,
    unit_price_amount   NUMERIC(12,2) NOT NULL,
    unit_price_currency CHAR(3)       NOT NULL,
    CONSTRAINT order_items_pkey       PRIMARY KEY (id),
    CONSTRAINT order_items_order_fk   FOREIGN KEY (order_id)
                                      REFERENCES orders (id),
    CONSTRAINT order_items_product_fk FOREIGN KEY (product_id)
                                      REFERENCES products (id)
);
```

---

## Part 4 — How the Model Drives the Schema

The generated SQL is not invented ad hoc — it is derived directly from the len model:

| len model element | Generated SQL |
|-------------------|---------------|
| `struct Product` | `CREATE TABLE products (...)` |
| Field `name : String` | `name VARCHAR(255) NOT NULL` |
| Field `id : ProductId` | `id UUID NOT NULL` + `PRIMARY KEY` |
| `spec sku_unique` | `CONSTRAINT products_sku_unique UNIQUE (sku)` |
| `spec product_has_category` | `category_id` column + `FOREIGN KEY` to `categories` |
| `rel OrderBelongsTo(o: Order, c: Customer)` | `customer_id` + `FOREIGN KEY` to `customers` |
| `const StatusPending : OrderStatus` | `DEFAULT 'pending'` on the `status` column |

The domain model is the single source of truth. The SQL is a derivative artefact.

---

## Part 5 — Validating the Generator

Run the validator on the entire project:

```bash
len-cli validate shop/ codegen/
```

The validator checks:

- All imports resolve (`shop.types`, `codegen.sql_types`, etc.)
- All names used in specs and quasi blocks are declared
- All `fn` declarations have a `quasi` block
- Arity of every relation application is correct

---

## What You Learned

- A len domain model is a complete, self-consistent specification of the domain
- A code generator is itself a len program: a set of `rel` + `spec` + `fn` declarations
- The generator's `fn` declarations use `quasi` blocks to sketch the emission logic
- The generated SQL is derived directly from the model; constraints come from `spec` laws
- `len-cli validate` catches structural errors in both the domain model and the generator before any SQL is emitted

**You have reached the end of the tutorial series.** Return to the [index](../index.md) for reference material.
