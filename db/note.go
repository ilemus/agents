package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// JSONTags is a custom type representing an array of tags stored as a JSON array in SQLite/Turso.
type JSONTags []string

// Value implements the driver.Valuer interface.
func (j JSONTags) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface.
func (j *JSONTags) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("failed to scan JSONTags: unsupported type %T", value)
		}
		bytes = []byte(str)
	}
	return json.Unmarshal(bytes, j)
}

// Note represents a record in the "notes" table.
type Note struct {
	ID        uuid.UUID `gorm:"type:text;primaryKey"`
	Note      string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
	Tags      JSONTags  `gorm:"type:json"`
}

// TableName overrides GORM's default naming convention to ensure the table is explicitly named "notes".
func (Note) TableName() string {
	return "notes"
}

// BeforeCreate is a GORM hook that runs before inserting a note.
// It assigns a new Google UUID and timestamp if they are not already set.
func (n *Note) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	if n.CreatedAt.IsZero() {
		n.CreatedAt = time.Now()
	}
	return nil
}

// CreateNote inserts a new note into the database.
func CreateNote(db *gorm.DB, note *Note) error {
	return db.Create(note).Error
}

// GetNoteByID retrieves a note from the database by its UUID.
func GetNoteByID(db *gorm.DB, id uuid.UUID) (*Note, error) {
	var note Note
	if err := db.First(&note, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &note, nil
}

// ListNotes retrieves all notes from the database.
func ListNotes(db *gorm.DB) ([]Note, error) {
	var notes []Note
	if err := db.Order("created_at DESC").Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}
