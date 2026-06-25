package db

/**
This is an implementation of the document vector model.

Reference the documentation at docs/software/database/vector_database.md

The schema is:

```sql
CREATE TABLE document_vectors (
    id TEXT PRIMARY KEY,
    document_id TEXT NOT NULL,
    chunk_number INTEGER NOT NULL,
    embedding BLOB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```
*/

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

// DocumentVector represents a record in the "document_vectors" table.
type DocumentVector struct {
	ID          uuid.UUID
	DocumentID  uuid.UUID
	ChunkNumber int
	Embedding   Float32Vector
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CreateDocumentVector inserts a new document_vector into the database.
func CreateDocumentVector(db *sql.DB, documentVector *DocumentVector) error {
	if documentVector.ID == uuid.Nil {
		documentVector.ID = uuid.New()
	}
	if documentVector.CreatedAt.IsZero() {
		documentVector.CreatedAt = time.Now()
	}
	if documentVector.UpdatedAt.IsZero() {
		documentVector.UpdatedAt = time.Now()
	}

	query := `INSERT INTO document_vectors (id, document_id, chunk_number, embedding, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, documentVector.ID.String(), documentVector.DocumentID.String(), documentVector.ChunkNumber, documentVector.Embedding, documentVector.CreatedAt, documentVector.UpdatedAt)
	return err
}

// GetDocumentVectorByID retrieves a document_vector from the database by its ID.
func GetDocumentVectorByID(db *sql.DB, id uuid.UUID) (*DocumentVector, error) {
	var documentVector DocumentVector
	query := `SELECT id, document_id, chunk_number, embedding, created_at, updated_at FROM document_vectors WHERE id = ?`
	err := db.QueryRow(query, id.String()).Scan(&documentVector.ID, &documentVector.DocumentID, &documentVector.ChunkNumber, &documentVector.Embedding, &documentVector.CreatedAt, &documentVector.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &documentVector, nil
}

// QueryDocumentVectorByEmbedding queries the database for documents by the embedding vector.
// It uses cosine similarity to find the most similar documents.
//
// @param embedding: The embedding vector to query by.
// @param limit: The maximum number of documents to return.
// @param offset: The number of documents to skip.
//
// @return []DocumentVector: The documents that are most similar to the embedding vector.
// @return error: The error that occurred during the query.
func QueryDocumentVectorByEmbedding(db *sql.DB, embedding Float32Vector, limit int, offset int) ([]DocumentVector, error) {
	var documentVectors []DocumentVector

	// Use raw SQL query to use the cosine similarity function found at https://docs.turso.tech/guides/vector-search#cosine-distance
	query := `
		SELECT 
			id,
			document_id,
			chunk_number,
			embedding,
			created_at,
			updated_at
		FROM 
			document_vectors
		WHERE 
			vector_distance_cos(embedding, ?)
		ORDER BY
			vector_distance_cos(embedding, ?)
		LIMIT ? OFFSET ?
	`

	rows, err := db.Query(query, embedding, embedding, limit, offset)
	if err != nil {
		return nil, errors.New("failed to query document_vectors")
	}
	defer rows.Close()

	for rows.Next() {
		var dv DocumentVector
		if err := rows.Scan(&dv.ID, &dv.DocumentID, &dv.ChunkNumber, &dv.Embedding, &dv.CreatedAt, &dv.UpdatedAt); err != nil {
			return nil, err
		}
		documentVectors = append(documentVectors, dv)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return documentVectors, nil
}

// ListDocumentVectors retrieves all document_vectors from the database.
func ListDocumentVectors(db *sql.DB) ([]DocumentVector, error) {
	query := `SELECT id, document_id, chunk_number, embedding, created_at, updated_at FROM document_vectors`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documentVectors []DocumentVector
	for rows.Next() {
		var dv DocumentVector
		if err := rows.Scan(&dv.ID, &dv.DocumentID, &dv.ChunkNumber, &dv.Embedding, &dv.CreatedAt, &dv.UpdatedAt); err != nil {
			return nil, err
		}
		documentVectors = append(documentVectors, dv)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return documentVectors, nil
}

// UpdateDocumentVector updates a document_vector in the database.
func UpdateDocumentVector(db *sql.DB, documentVector *DocumentVector) error {
	documentVector.UpdatedAt = time.Now()
	query := `UPDATE document_vectors SET document_id = ?, chunk_number = ?, embedding = ?, updated_at = ? WHERE id = ?`
	_, err := db.Exec(query, documentVector.DocumentID.String(), documentVector.ChunkNumber, documentVector.Embedding, documentVector.UpdatedAt, documentVector.ID.String())
	return err
}

// DeleteDocumentVector deletes a document_vector from the database.
func DeleteDocumentVector(db *sql.DB, documentVector *DocumentVector) error {
	query := `DELETE FROM document_vectors WHERE id = ?`
	_, err := db.Exec(query, documentVector.ID.String())
	return err
}

