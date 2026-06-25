package db

/**
This is an implementation of the document table.

Reference the documentation at docs/software/database/vector_database.md

The schema is:

```sql
CREATE TABLE documents (
    id TEXT PRIMARY KEY,
    file_path TEXT NOT NULL,
    title TEXT,
    summary TEXT,
    tags TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```
*/

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Document represents a record in the "documents" table.
type Document struct {
	ID        uuid.UUID
	FilePath  string
	Title     string
	Summary   string
	Tags      JSONTags
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateDocument inserts a new document into the database.
func CreateDocument(db *sql.DB, document *Document) error {
	if document.ID == uuid.Nil {
		document.ID = uuid.New()
	}
	if document.CreatedAt.IsZero() {
		document.CreatedAt = time.Now()
	}
	if document.UpdatedAt.IsZero() {
		document.UpdatedAt = time.Now()
	}

	query := `INSERT INTO documents (id, file_path, title, summary, tags, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, document.ID.String(), document.FilePath, document.Title, document.Summary, document.Tags, document.CreatedAt, document.UpdatedAt)
	return err
}

// GetDocumentByFilePath retrieves a document from the database by its file path.
func GetDocumentByFilePath(db *sql.DB, filePath string) (*Document, error) {
	var document Document
	query := `SELECT id, file_path, title, summary, tags, created_at, updated_at FROM documents WHERE file_path = ?`
	err := db.QueryRow(query, filePath).Scan(&document.ID, &document.FilePath, &document.Title, &document.Summary, &document.Tags, &document.CreatedAt, &document.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &document, nil
}

// GetDocumentByID retrieves a document from the database by its ID.
func GetDocumentByID(db *sql.DB, id uuid.UUID) (*Document, error) {
	var document Document
	query := `SELECT id, file_path, title, summary, tags, created_at, updated_at FROM documents WHERE id = ?`
	err := db.QueryRow(query, id.String()).Scan(&document.ID, &document.FilePath, &document.Title, &document.Summary, &document.Tags, &document.CreatedAt, &document.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &document, nil
}

// ListDocuments retrieves all documents from the database.
func ListDocuments(db *sql.DB) ([]Document, error) {
	query := `SELECT id, file_path, title, summary, tags, created_at, updated_at FROM documents`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []Document
	for rows.Next() {
		var document Document
		if err := rows.Scan(&document.ID, &document.FilePath, &document.Title, &document.Summary, &document.Tags, &document.CreatedAt, &document.UpdatedAt); err != nil {
			return nil, err
		}
		documents = append(documents, document)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return documents, nil
}

// UpdateDocument updates a document in the database.
func UpdateDocument(db *sql.DB, document *Document) error {
	document.UpdatedAt = time.Now()
	query := `UPDATE documents SET file_path = ?, title = ?, summary = ?, tags = ?, updated_at = ? WHERE id = ?`
	_, err := db.Exec(query, document.FilePath, document.Title, document.Summary, document.Tags, document.UpdatedAt, document.ID.String())
	return err
}

// DeleteDocument deletes a document from the database.
func DeleteDocument(db *sql.DB, document *Document) error {
	query := `DELETE FROM documents WHERE id = ?`
	_, err := db.Exec(query, document.ID.String())
	return err
}

