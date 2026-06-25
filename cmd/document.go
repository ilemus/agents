package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var documentCmd = &cobra.Command{
	Use:   "document",
	Short: "Manage documents",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var documentCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a document, potentially from existing notes or documents",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Implement document create logic
		// Executes prompt with reference text and search tool
		// Creates file, registers in table, updates vector table
		fmt.Println("document create called")
	},
}

var documentImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Ingests an existing document file and updates tables",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Implement document import logic
		fmt.Println("document import called")
	},
}

var documentUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates an existing registered document",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Implement document update logic
		// Carry forward tags, re-tag based on content
		fmt.Println("document update called")
	},
}

var documentQueryCmd = &cobra.Command{
	Use:   "query",
	Short: "Queries the documents vector store",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Implement document query logic
		fmt.Println("document query called")
	},
}

func init() {
	rootCmd.AddCommand(documentCmd)

	// document create flags
	documentCmd.AddCommand(documentCreateCmd)
	documentCreateCmd.Flags().StringSlice("note-ids", []string{}, "Comma-separated note IDs to reference")
	documentCreateCmd.Flags().StringSlice("doc-ids", []string{}, "Comma-separated document IDs to reference")
	documentCreateCmd.Flags().String("prompt", "", "Prompt string or prompt ID")
	documentCreateCmd.Flags().String("template", "", "File path to a document template")
	documentCreateCmd.Flags().StringP("output", "o", "", "Output file path")

	// document import flags
	documentCmd.AddCommand(documentImportCmd)
	documentImportCmd.Flags().StringP("file", "f", "", "Path to the file to import")

	// document update flags
	documentCmd.AddCommand(documentUpdateCmd)
	documentUpdateCmd.Flags().StringSlice("note-ids", []string{}, "Comma-separated note IDs to reference")
	documentUpdateCmd.Flags().StringSlice("doc-ids", []string{}, "Comma-separated document IDs to reference")
	documentUpdateCmd.Flags().String("prompt", "", "Prompt string or prompt ID")
	documentUpdateCmd.Flags().String("file", "", "File path of the document being updated")
	documentUpdateCmd.Flags().String("doc-id", "", "ID of the document being updated")
	documentUpdateCmd.Flags().Bool("create", false, "Create document entry if not registered")

	// document query flags
	documentCmd.AddCommand(documentQueryCmd)
	documentQueryCmd.Flags().IntP("limit", "l", 10, "Limit the number of results returned")
	documentQueryCmd.Flags().BoolP("rank", "r", false, "Enables reranking")
}
