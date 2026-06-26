package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Note represents a record in the "notes" table.
type Note struct {
	ID        uuid.UUID
	Note      string
	CreatedAt time.Time
	Tags      JSONTags
}

// CreateNote inserts a new note into the database.
func CreateNote(db *sql.DB, note *Note) (uuid.UUID, error) {
	if note.ID == uuid.Nil {
		note.ID = uuid.New()
	}
	if note.CreatedAt.IsZero() {
		note.CreatedAt = time.Now()
	}

	query := `INSERT INTO notes (id, note, created_at, tags) VALUES (?, ?, ?, ?)`
	_, err := db.Exec(query, note.ID.String(), note.Note, note.CreatedAt, note.Tags)
	if err != nil {
		return uuid.Nil, err
	}
	return note.ID, nil
}

// GetNoteByID retrieves a note from the database by its UUID.
func GetNoteByID(db *sql.DB, id uuid.UUID) (*Note, error) {
	var note Note
	query := `SELECT id, note, created_at, tags FROM notes WHERE id = ?`
	err := db.QueryRow(query, id.String()).Scan(&note.ID, &note.Note, &note.CreatedAt, &note.Tags)
	if err != nil {
		return nil, err
	}
	return &note, nil
}

// ListNotes retrieves all notes from the database.
func ListNotes(db *sql.DB) ([]Note, error) {
	query := `SELECT id, note, created_at, tags FROM notes ORDER BY created_at DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		if err := rows.Scan(&note.ID, &note.Note, &note.CreatedAt, &note.Tags); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return notes, nil
}

