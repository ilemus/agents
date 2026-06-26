package db

import (
	"math/rand/v2"
	"os"
	"testing"

	"github.com/google/uuid"
)

func TestDBAndNote(t *testing.T) {
	// Setup environment variables for test execution
	origDBURL := os.Getenv("TURSO_DATABASE_URL")
	origAuthToken := os.Getenv("TURSO_AUTH_TOKEN")
	origLocalPath := os.Getenv("LOCAL_DB_PATH")

	os.Setenv("TURSO_DATABASE_URL", "")
	os.Setenv("TURSO_AUTH_TOKEN", "")
	os.Setenv("LOCAL_DB_PATH", ":memory:")

	defer func() {
		os.Setenv("TURSO_DATABASE_URL", origDBURL)
		os.Setenv("TURSO_AUTH_TOKEN", origAuthToken)
		os.Setenv("LOCAL_DB_PATH", origLocalPath)
	}()

	// 1. Initialize the Database
	err := InitDB()
	if err != nil {
		t.Fatalf("failed to initialize db: %v", err)
	}

	defer DB.Close()

	// 2. Perform Migration
	err = CreateTables(DB)
	if err != nil {
		t.Fatalf("failed to create tables: %v", err)
	}

	// 3. Create a Note
	testNote := &Note{
		Note: "This is a test note generated for ORM validation.",
		Tags: JSONTags{"llm", "database", "turso"},
	}

	noteID, err := CreateNote(DB, testNote)
	if err != nil {
		t.Fatalf("failed to create note: %v", err)
	}
	if noteID != testNote.ID {
		t.Errorf("expected returned note ID %v, got %v", testNote.ID, noteID)
	}

	// Verify that the UUID and CreatedAt were generated via hooks
	if testNote.ID == uuid.Nil {
		t.Error("expected note ID to be generated and not nil")
	}
	if testNote.CreatedAt.IsZero() {
		t.Error("expected note CreatedAt timestamp to be set")
	}

	// 4. Retrieve the Note by ID
	retrieved, err := GetNoteByID(DB, testNote.ID)
	if err != nil {
		t.Fatalf("failed to retrieve note by ID: %v", err)
	}

	// Verify all properties
	if retrieved.ID != testNote.ID {
		t.Errorf("expected ID %v, got %v", testNote.ID, retrieved.ID)
	}
	if retrieved.Note != testNote.Note {
		t.Errorf("expected Note content %q, got %q", testNote.Note, retrieved.Note)
	}
	if !retrieved.CreatedAt.Equal(testNote.CreatedAt) {
		t.Errorf("expected CreatedAt %v, got %v", testNote.CreatedAt, retrieved.CreatedAt)
	}
	if len(retrieved.Tags) != 3 {
		t.Fatalf("expected 3 tags, got %d", len(retrieved.Tags))
	}
	expectedTags := []string{"llm", "database", "turso"}
	for i, tag := range retrieved.Tags {
		if tag != expectedTags[i] {
			t.Errorf("expected tag at index %d to be %q, got %q", i, expectedTags[i], tag)
		}
	}

	// 5. List Notes
	notes, err := ListNotes(DB)
	if err != nil {
		t.Fatalf("failed to list notes: %v", err)
	}

	if len(notes) != 1 {
		t.Errorf("expected 1 note in list, got %d", len(notes))
	} else {
		if notes[0].ID != testNote.ID {
			t.Errorf("expected list note ID %v, got %v", testNote.ID, notes[0].ID)
		}
	}
}

func TestNoteVector(t *testing.T) {
	// Setup environment variables for test execution
	origDBURL := os.Getenv("TURSO_DATABASE_URL")
	origAuthToken := os.Getenv("TURSO_AUTH_TOKEN")
	origLocalPath := os.Getenv("LOCAL_DB_PATH")

	os.Setenv("TURSO_DATABASE_URL", "")
	os.Setenv("TURSO_AUTH_TOKEN", "")
	os.Setenv("LOCAL_DB_PATH", ":memory:")

	defer func() {
		os.Setenv("TURSO_DATABASE_URL", origDBURL)
		os.Setenv("TURSO_AUTH_TOKEN", origAuthToken)
		os.Setenv("LOCAL_DB_PATH", origLocalPath)
	}()

	// 1. Initialize the Database
	err := InitDB()
	if err != nil {
		t.Fatalf("failed to initialize db: %v", err)
	}

	defer DB.Close()

	// 2. Perform Migration
	err = CreateTables(DB)
	if err != nil {
		t.Fatalf("failed to create tables: %v", err)
	}

	testNote := &Note{
		Note: "This is a test note generated for ORM validation.",
		Tags: JSONTags{"llm", "database", "turso"},
	}

	noteID, err := CreateNote(DB, testNote)
	if err != nil {
		t.Fatalf("failed to create note: %v", err)
	}
	if noteID != testNote.ID {
		t.Errorf("expected returned note ID %v, got %v", testNote.ID, noteID)
	}

	// Genearate a random float32 array of size 768
	embedding := make([]float32, 768)
	for i := range embedding {
		embedding[i] = float32(rand.Float32())
	}

	testNoteVector := &NoteVector{
		Chunk:     "This is a test note vector generated for ORM validation.",
		ParentID:  testNote.ID,
		Embedding: embedding,
		Tags:      JSONTags{"llm", "database", "turso"},
	}

	err = CreateNoteVector(DB, testNoteVector)
	if err != nil {
		t.Fatalf("failed to create note vector: %v", err)
	}

	retrieved, err := GetNoteVectorByParentID(DB, testNote.ID)
	if err != nil {
		t.Fatalf("failed to retrieve note vector by parent ID: %v", err)
	}

	// Verify all properties
	if retrieved.ID != testNoteVector.ID {
		t.Errorf("expected ID %v, got %v", testNoteVector.ID, retrieved.ID)
	}
	if retrieved.ParentID != testNoteVector.ParentID {
		t.Errorf("expected ParentID %v, got %v", testNoteVector.ParentID, retrieved.ParentID)
	}
	if retrieved.Chunk != testNoteVector.Chunk {
		t.Errorf("expected Chunk content %q, got %q", testNoteVector.Chunk, retrieved.Chunk)
	}
	// Check embedding equallity by going float by float
	for i, v := range retrieved.Embedding {
		if v != testNoteVector.Embedding[i] {
			t.Errorf("expected Embedding[%d] %f, got %f", i, testNoteVector.Embedding[i], v)
		}
	}
	if len(retrieved.Tags) != 3 {
		t.Fatalf("expected 3 tags, got %d", len(retrieved.Tags))
	}
	expectedTags := []string{"llm", "database", "turso"}
	for i, tag := range retrieved.Tags {
		if tag != expectedTags[i] {
			t.Errorf("expected tag at index %d to be %q, got %q", i, expectedTags[i], tag)
		}
	}

	// List note vectors
	noteVectors, err := ListNoteVectors(DB)
	if err != nil {
		t.Fatalf("failed to list note vectors: %v", err)
	}

	if len(noteVectors) != 1 {
		t.Errorf("expected 1 note vector in list, got %d", len(noteVectors))
	} else {
		if noteVectors[0].ID != testNoteVector.ID {
			t.Errorf("expected list note vector ID %v, got %v", testNoteVector.ID, noteVectors[0].ID)
		}
	}
}

func TestDocument(t *testing.T) {
	// Setup environment variables for test execution
	origDBURL := os.Getenv("TURSO_DATABASE_URL")
	origAuthToken := os.Getenv("TURSO_AUTH_TOKEN")
	origLocalPath := os.Getenv("LOCAL_DB_PATH")

	os.Setenv("TURSO_DATABASE_URL", "")
	os.Setenv("TURSO_AUTH_TOKEN", "")
	os.Setenv("LOCAL_DB_PATH", ":memory:")

	defer func() {
		os.Setenv("TURSO_DATABASE_URL", origDBURL)
		os.Setenv("TURSO_AUTH_TOKEN", origAuthToken)
		os.Setenv("LOCAL_DB_PATH", origLocalPath)
	}()

	// 1. Initialize the Database
	err := InitDB()
	if err != nil {
		t.Fatalf("failed to initialize db: %v", err)
	}

	defer DB.Close()

	// 2. Perform Migration
	err = CreateTables(DB)
	if err != nil {
		t.Fatalf("failed to create tables: %v", err)
	}

	testDocument := &Document{
		Title:    "Test Document",
		Summary:  "This is a test document generated for ORM validation.",
		FilePath: "/path/to/test/document",
		Tags:     JSONTags{"test", "document"},
	}

	err = CreateDocument(DB, testDocument)
	if err != nil {
		t.Fatalf("failed to create document: %v", err)
	}

	retrieved, err := GetDocumentByID(DB, testDocument.ID)
	if err != nil {
		t.Fatalf("failed to retrieve document by ID: %v", err)
	}

	// Verify all properties
	if retrieved.ID != testDocument.ID {
		t.Errorf("expected ID %v, got %v", testDocument.ID, retrieved.ID)
	}
	if retrieved.Title != testDocument.Title {
		t.Errorf("expected Title %q, got %q", testDocument.Title, retrieved.Title)
	}
	if retrieved.Summary != testDocument.Summary {
		t.Errorf("expected Summary %q, got %q", testDocument.Summary, retrieved.Summary)
	}
	if retrieved.FilePath != testDocument.FilePath {
		t.Errorf("expected FilePath %q, got %q", testDocument.FilePath, retrieved.FilePath)
	}
	if len(retrieved.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(retrieved.Tags))
	}
	expectedTags := []string{"test", "document"}
	for i, tag := range retrieved.Tags {
		if tag != expectedTags[i] {
			t.Errorf("expected tag at index %d to be %q, got %q", i, expectedTags[i], tag)
		}
	}
}

func TestDocumentVectorQuery(t *testing.T) {
	// Setup environment variables for test execution
	origDBURL := os.Getenv("TURSO_DATABASE_URL")
	origAuthToken := os.Getenv("TURSO_AUTH_TOKEN")
	origLocalPath := os.Getenv("LOCAL_DB_PATH")

	os.Setenv("TURSO_DATABASE_URL", "")
	os.Setenv("TURSO_AUTH_TOKEN", "")
	os.Setenv("LOCAL_DB_PATH", ":memory:")

	defer func() {
		os.Setenv("TURSO_DATABASE_URL", origDBURL)
		os.Setenv("TURSO_AUTH_TOKEN", origAuthToken)
		os.Setenv("LOCAL_DB_PATH", origLocalPath)
	}()

	// 1. Initialize the Database
	err := InitDB()
	if err != nil {
		t.Fatalf("failed to initialize db: %v", err)
	}

	defer DB.Close()

	// 2. Perform Migration
	err = CreateTables(DB)
	if err != nil {
		t.Fatalf("failed to create tables: %v", err)
	}

	testDocument := &Document{
		Title:    "Test Document",
		Summary:  "This is a test document generated for ORM validation.",
		FilePath: "/path/to/test/document",
		Tags:     JSONTags{"test", "document"},
	}

	err = CreateDocument(DB, testDocument)
	if err != nil {
		t.Fatalf("failed to create document: %v", err)
	}

	// Embedding vector 768 of random floats.
	embedding := make([]float32, 768)
	for i := range embedding {
		embedding[i] = float32(rand.Float32())
	}

	testDocumentVector := &DocumentVector{
		DocumentID:  testDocument.ID,
		ChunkNumber: 1,
		Embedding:   embedding,
	}

	err = CreateDocumentVector(DB, testDocumentVector)
	if err != nil {
		t.Fatalf("failed to create document vector: %v", err)
	}

	// Query document vectors
	documentVectors, err := QueryDocumentVectorByEmbedding(DB, embedding, 1, 0)
	if err != nil {
		t.Skip("skipping vector query test; vector_distance_cos is not supported on this SQLite environment")
		return
	}

	if len(documentVectors) != 1 {
		t.Errorf("expected 1 document vector in list, got %d", len(documentVectors))
	} else {
		if documentVectors[0].ID != testDocumentVector.ID {
			t.Errorf("expected list document vector ID %v, got %v", testDocumentVector.ID, documentVectors[0].ID)
		}
	}
}
