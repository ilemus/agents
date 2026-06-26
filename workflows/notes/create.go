package notes

import (
	"fmt"
	"log"
	"strings"

	"github.com/ilemus/agents/db"
	"github.com/ilemus/agents/llm"
	"github.com/ilemus/agents/llm/ollama"

	"context"
)

/*
This file is used for implementing the notes create command functionality.

The notes command handles capturing the note as a block of text.
The purpose of this workflow is to:
1. Create tags based on the note's content
2. Store the note as a complete text block along with the tags generated
3. Semantically chunk the note
4. Store the semantic chunks into a vector store table
*/

const NOTE_MODEL = "gemma4:e4b"

const NOTE_TAG_PROMPT = `You are a tagging system for notes.

<rules>
- Generate up to 5 tags based on the note's content
- Tags should be comma separated
- Only output the tags, nothing else
</rules>

<output>
abc, def, ghi, jkl
</output>

<tag_examples>
<example>
input: "Store text blocks into a vector store database for easy retrieval for LLMs"
output: "vector store, embeddings, semantic search"
</example>
<example>
input: "The logs show that there is an error in processing complex JSON structures. The schema should be simplified."
output: "logs, error, JSON, schema, simplification"
</example>
</tag_examples>

<task>Generate comma-separated tags provided a note from the user.</task>`

func GenerateTags(ctx context.Context, note string) []string {
	client := ollama.NewClient("")

	// Create the request, which uses the NOTE_TAG_PROMPT as the system message
	req := llm.ChatRequest{
		Model: NOTE_MODEL,
		Messages: []llm.Message{
			{
				Role:    "system",
				Content: NOTE_TAG_PROMPT,
			},
			{
				Role:    "user",
				Content: note,
			},
		},
		Options: ollama.DefaultOptions,
	}

	resp, err := client.Chat(ctx, req)
	if err != nil {
		return []string{}
	}

	content := resp.Message.Content

	// Split the tags by comman and return a string array
	tags := strings.Split(content, ",")
	return tags
}

func Execute(note string) {
	ctx := context.Background()

	// Generate tags for the note
	tags := GenerateTags(ctx, note)

	// Connect to the database
	db.InitDB()
	defer db.DB.Close()

	// Store into the database
	noteObj := db.Note{
		Tags: tags,
		Note: note,
	}
	noteId, err := db.CreateNote(db.DB, &noteObj)
	if err != nil {
		log.Fatalf("Failed to create note: %v", err)
	}
	// TODO: Remove this after testing
	fmt.Println("Note created with ID: ", noteId)

}
