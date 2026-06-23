package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "agents",
	Short: "agents is a project planner, tracker, and LLM benchmarking tool",
	Long: `A CLI tool designed to:
1. Help plan, document, and track project work using a "reverse agent" approach.
2. Benchmark LLMs on various agentic tasks (tool calling, long context, reasoning, etc.).`,
	Version: "0.1.0",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// You can define global flags here if needed.
}
