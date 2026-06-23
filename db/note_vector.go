package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Float32Vector []float32

// Value converts the slice into a JSON-marshaled string/bytes for the database
func (v Float32Vector) Value() (driver.Value, error) {
	if v == nil {
		return nil, nil
	}
	// Marshal the slice to a JSON byte array
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal float32 vector: %w", err)
	}
	// Return it as a string (or []byte, which drivers handle cleanly)
	return string(bytes), nil
}

// Scan converts the database value (JSON string/bytes) back into a Go slice
func (v *Float32Vector) Scan(value interface{}) error {
	if value == nil {
		*v = nil
		return nil
	}

	var bytes []byte
	switch data := value.(type) {
	case []byte:
		bytes = data
	case string:
		bytes = []byte(data)
	default:
		return fmt.Errorf("unsupported data type for Float32Vector: %T", value)
	}

	return json.Unmarshal(bytes, v)
}

// Note represents a record in the "notes" table.
type NoteVector struct {
	ID        uuid.UUID     `gorm:"type:text;primaryKey"`
	ParentID  uuid.UUID     `gorm:"type:text;not null"`
	Chunk     string        `gorm:"type:text;not null"`
	Embedding Float32Vector `gorm:"type:F32_BLOB(768)"`
	Tags      JSONTags      `gorm:"type:json"`
}

// TableName overrides GORM's default naming convention to ensure the table is explicitly named "notes".
func (NoteVector) TableName() string {
	return "note_vectors"
}

// BeforeCreate is a GORM hook that runs before inserting a note.
// It assigns a new Google UUID and timestamp if they are not already set.
func (n *NoteVector) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}

// CreateNoteVector inserts a new note vector into the database.
func CreateNoteVector(db *gorm.DB, noteVector *NoteVector) error {
	return db.Create(noteVector).Error
}

// GetNoteVectorByParentID retrieves a note vector from the database by its parent ID.
func GetNoteVectorByParentID(db *gorm.DB, parentID uuid.UUID) (*NoteVector, error) {
	var noteVector NoteVector
	if err := db.Where("parent_id = ?", parentID).First(&noteVector).Error; err != nil {
		return nil, err
	}
	return &noteVector, nil
}

// ListNoteVectors retrieves all note vectors from the database.
func ListNoteVectors(db *gorm.DB) ([]NoteVector, error) {
	var noteVectors []NoteVector
	if err := db.Find(&noteVectors).Error; err != nil {
		return nil, err
	}
	return noteVectors, nil
}
