# Implementation Plan: Adding Note Vector Model and Database Support

This plan outlines the changes required to introduce a `note_vectors` database table in Go using GORM, aligned with Turso's native vector search features.

## User Review Required

> [!NOTE]
> - We define a custom type `Embedding []float32` in Go.
> - `Embedding` implements `driver.Valuer` to automatically serialize the float32 array to standard little-endian binary bytes for storage in the `F32_BLOB(768)` column.
> - `Embedding` implements `sql.Scanner` to automatically deserialize the binary float32 array back into `[]float32`. It also includes a JSON fallback to handle text/JSON representations if encountered.
> - This eliminates the need for separate standalone conversion functions, making model interaction seamless.

## Proposed Changes

---

### Database Component

#### [MODIFY] [note.go](file:///home/lemusi/vscode/agents/db/note.go)
- Define `type Embedding []float32` custom type.
- Implement `Value()` and `Scan()` on `Embedding`:
  - `Value()` serializes `[]float32` to binary (`[]byte`) representation.
  - `Scan()` deserializes binary (or JSON arrays) back to `[]float32`.
- Add a new `NoteVector` struct mapping to the `note_vectors` table:
  - `ID`: primary key (`uuid.UUID`)
  - `ParentID`: foreign key (`uuid.UUID`), maps to `parent_id` and references `notes(id)`
  - `Embedding`: custom `Embedding` type mapped to `F32_BLOB(768)`
  - `Tags`: custom `JSONTags` array
- Add `BeforeCreate` hook for `NoteVector` to generate UUIDs automatically.
- Add helpers:
  - `CreateNoteVector(db *gorm.DB, noteVector *NoteVector) error`
  - `GetNoteVectorByParentID(db *gorm.DB, parentID uuid.UUID) (*NoteVector, error)`

#### [MODIFY] [db_test.go](file:///home/lemusi/vscode/agents/db/db_test.go)
- Update migration step to auto-migrate both `Note` and `NoteVector`.
- Add test coverage for `NoteVector` CRUD operations, verifying that `[]float32` values are correctly saved, retrieved, and matched.

## Verification Plan

### Automated Tests
Run unit tests for the database package:
```bash
go test -v ./db/...
```
This will verify:
- Database schema migration for `note_vectors` with `parent_id` foreign key.
- Creating a `NoteVector` with float32 array, retrieving it, and asserting the float32 values are exactly identical.
