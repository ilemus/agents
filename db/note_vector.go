package db

import (
	"database/sql"

	"github.com/google/uuid"
)

// NoteVector represents a record in the "note_vectors" table.
type NoteVector struct {
	ID        uuid.UUID
	ParentID  uuid.UUID
	Chunk     string
	Embedding Float32Vector
	Tags      JSONTags
}

// CreateNoteVector inserts a new note vector into the database.
func CreateNoteVector(db *sql.DB, noteVector *NoteVector) error {
	if noteVector.ID == uuid.Nil {
		noteVector.ID = uuid.New()
	}

	query := `INSERT INTO note_vectors (id, parent_id, chunk, embedding, tags) VALUES (?, ?, ?, ?, ?)`
	_, err := db.Exec(query, noteVector.ID.String(), noteVector.ParentID.String(), noteVector.Chunk, noteVector.Embedding, noteVector.Tags)
	return err
}

// GetNoteVectorByParentID retrieves a note vector from the database by its parent ID.
func GetNoteVectorByParentID(db *sql.DB, parentID uuid.UUID) (*NoteVector, error) {
	var noteVector NoteVector
	query := `SELECT id, parent_id, chunk, embedding, tags FROM note_vectors WHERE parent_id = ?`
	err := db.QueryRow(query, parentID.String()).Scan(&noteVector.ID, &noteVector.ParentID, &noteVector.Chunk, &noteVector.Embedding, &noteVector.Tags)
	if err != nil {
		return nil, err
	}
	return &noteVector, nil
}

// ListNoteVectors retrieves all note vectors from the database.
func ListNoteVectors(db *sql.DB) ([]NoteVector, error) {
	query := `SELECT id, parent_id, chunk, embedding, tags FROM note_vectors`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var noteVectors []NoteVector
	for rows.Next() {
		var nv NoteVector
		if err := rows.Scan(&nv.ID, &nv.ParentID, &nv.Chunk, &nv.Embedding, &nv.Tags); err != nil {
			return nil, err
		}
		noteVectors = append(noteVectors, nv)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return noteVectors, nil
}

