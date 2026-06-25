package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var notesCmd = &cobra.Command{
	Use:   "notes",
	Short: "Manage notes",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var notesImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Imports text blobs or files into the notes vector store",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Implement notes import logic
		// When imported through a file, the content is semantically chunked and stored
		fmt.Println("notes import called")
	},
}

var notesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new note using an external text editor (e.g., vi)",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Implement the actual note creation logic (e.g., saving to DB/vector store)
		// For now, we open the editor, capture content, and stub the import.
		content, err := openEditorInTempFile()
		if err != nil {
			fmt.Printf("Error opening editor: %v\n", err)
			return
		}

		if content == "" {
			fmt.Println("Note content was empty. Aborting creation.")
			return
		}

		// This is example code
		fmt.Println("Note created with content:")
		fmt.Println(content)
		fmt.Println("\nTODO: Automatically importing this note content into the notes vector store...")
	},
}

var notesQueryCmd = &cobra.Command{
	Use:   "query",
	Short: "Queries the notes vector store",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Implement notes query logic
		// Attempts general answer, then vector embedding query
		// Results returned in table format
		fmt.Println("notes query called")
	},
}

func init() {
	rootCmd.AddCommand(notesCmd)

	// notes create command
	notesCmd.AddCommand(notesCreateCmd)

	// notes import flags
	notesCmd.AddCommand(notesImportCmd)
	notesImportCmd.Flags().BoolP("import", "i", false, "Trigger import")
	notesImportCmd.Flags().StringP("file", "f", "", "Path to the file being imported")

	// notes query flags
	notesCmd.AddCommand(notesQueryCmd)
	notesQueryCmd.Flags().IntP("limit", "l", 10, "Limit the number of results returned")
	notesQueryCmd.Flags().BoolP("rank", "r", false, "Enables reranking")
}

// openEditorInTempFile creates a temporary file, opens the system default editor (or vi/nano),
// captures the edited content, and returns it.
func openEditorInTempFile() (string, error) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "agent-note-*.md")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Determine editor to use
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi" // Default to vi
	}

	// Prepare command
	cmd := exec.Command(editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start editor and wait for it to finish
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to run editor %s: %w", editor, err)
	}

	// Read content from the temp file
	contentBytes, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to read temporary file: %w", err)
	}

	return strings.TrimSpace(string(contentBytes)), nil
}
