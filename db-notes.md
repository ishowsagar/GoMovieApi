# Database Calls Notes (Go `database/sql`)

A practical guide for choosing the right DB method in APIs.

## Core Methods

### 1. `QueryContext(ctx, query, args...)`

Use this when your query returns multiple rows.

Best for:

- list endpoints
- search endpoints
- reports

Pattern:

```go
rows, err := db.QueryContext(ctx, query, args...)
if err != nil {
    return nil, err
}
defer rows.Close()

var items []Item
for rows.Next() {
    var it Item
    if err := rows.Scan(&it.ID, &it.Name); err != nil {
        return nil, err
    }
    items = append(items, it)
}

if err := rows.Err(); err != nil {
    return nil, err
}

return items, nil
```

### 2. `QueryRowContext(ctx, query, args...)`

Use this when your query should return exactly one row.

Best for:

- get by id
- get by unique field (email, slug, etc.)

Pattern:

```go
var item Item
err := db.QueryRowContext(ctx, query, id).Scan(
    &item.ID,
    &item.Name,
    &item.CreatedAt,
)
if err != nil {
    return nil, err
}
return &item, nil
```

### 3. `ExecContext(ctx, query, args...)`

Use this when you do not need row data back.

Best for:

- UPDATE
- DELETE
- CREATE/ALTER
- INSERT without returning columns

Pattern:

```go
result, err := db.ExecContext(ctx, query, args...)
if err != nil {
    return err
}

affected, err := result.RowsAffected()
if err != nil {
    return err
}
if affected == 0 {
    return fmt.Errorf("no rows affected")
}
```

### 4. `PrepareContext(ctx, query)`

Use when the same statement runs many times in a loop.

Pattern:

```go
stmt, err := db.PrepareContext(ctx, query)
if err != nil {
    return err
}
defer stmt.Close()

_, err = stmt.ExecContext(ctx, args...)
return err
```

## Rows Object Methods (after `QueryContext`)

- `rows.Next()` -> move to next row
- `rows.Scan(...)` -> read current row into pointers
- `rows.Close()` -> release resources (always `defer`)
- `rows.Err()` -> final iteration error check

## Transactions

Use when multiple steps must succeed/fail together.

Pattern:

```go
tx, err := db.BeginTx(ctx, nil)
if err != nil {
    return err
}
defer tx.Rollback()

if _, err := tx.ExecContext(ctx, q1, a1); err != nil {
    return err
}
if _, err := tx.ExecContext(ctx, q2, a2); err != nil {
    return err
}

return tx.Commit()
```

## Quick Decision Rule

1. Need many rows -> `QueryContext`
2. Need one row -> `QueryRowContext`
3. Need no returned row data -> `ExecContext`
4. Need atomic multi-step flow -> transaction (`BeginTx`)

## Your Movie API Mapping

- `GetAllMovies` -> `QueryContext` (many rows)
- `GetMovieByID` -> `QueryRowContext` (single row)
- `CreateMovie`:
  - if you need inserted id/timestamps back -> `QueryRowContext` with `RETURNING ...`
  - if you only need success/failure -> `ExecContext`

Example create with return:

```sql
INSERT INTO movies(name, genre, description, ratings, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, name, genre, description, ratings, created_at, updated_at;
```

## Common Mistakes To Avoid

- Using `QueryContext` for a single-row lookup
- Calling `Scan` without pointer destinations
- Forgetting `defer rows.Close()`
- Forgetting to check `rows.Err()` after loop
- Not returning after writing an HTTP error response
- Passing non-UUID id when DB column type is UUID

## URL Param Reminder (chi)

- Define route as `/movies/{id}`
- Read id with `chi.URLParam(r, "id")`
- Route param name and `URLParam` key must match exactly

Example:

```go
r.Get("/movies/{id}", controller.GetMovieByID)

func (c Controller) GetMovieByID(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    if id == "" {
        http.Error(w, "id is required", http.StatusBadRequest)
        return
    }
}
```
