package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/ssdajoker/Code-Factory/internal/tui"
)

const version = "1.0.0"

var rootCmd = &cobra.Command{
	Use:   "factory",
	Short: "Spec-Driven Software Factory",
	Long: `Factory - A Spec-Driven Software Factory

Factory helps you capture project vision, generate specifications,
and maintain alignment between code and contracts.`,
	Run: func(cmd *cobra.Command, args []string) {
		// No subcommand: launch TUI
		if err := tui.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Factory v%s\n", version)
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Factory in current project",
	Run: func(cmd *cobra.Command, args []string) {
		quick, _ := cmd.Flags().GetBool("quick")
		fmt.Println("Running factory init...")
		if quick {
			fmt.Println("Using quick start mode")
		}
		fmt.Println("(Implementation coming soon)")
	},
}

var intakeCmd = &cobra.Command{
	Use:   "intake",
	Short: "Start INTAKE mode (capture vision)",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		fmt.Println("Starting INTAKE mode...")
		if name != "" {
			fmt.Printf("Project: %s\n", name)
		}
		if err := tui.RunIntake(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Start REVIEW mode (check code against specs)",
	Run: func(cmd *cobra.Command, args []string) {
		spec, _ := cmd.Flags().GetString("spec")
		fmt.Println("Starting REVIEW mode...")
		if spec != "" {
			fmt.Printf("Spec: %s\n", spec)
		}
		fmt.Println("(Implementation coming soon)")
	},
}

var rescueCmd = &cobra.Command{
	Use:   "rescue",
	Short: "Start RESCUE mode (reverse-engineer codebase)",
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("path")
		fmt.Println("Starting RESCUE mode...")
		fmt.Printf("Path: %s\n", path)
		fmt.Println("(Implementation coming soon)")
	},
}

var changeOrderCmd = &cobra.Command{
	Use:   "change-order",
	Short: "Start CHANGE_ORDER mode (track drift)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting CHANGE_ORDER mode...")
		fmt.Println("(Implementation coming soon)")
	},
}

var githubCmd = &cobra.Command{
	Use:   "github",
	Short: "GitHub integration",
	Run: func(cmd *cobra.Command, args []string) {
		login, _ := cmd.Flags().GetBool("login")
		status, _ := cmd.Flags().GetBool("status")

		if login {
			fmt.Println("Starting GitHub OAuth flow...")
			fmt.Println("(Implementation coming soon)")
			return
		}
		if status {
			fmt.Println("GitHub connection status: Not connected")
			return
		}
		fmt.Println("GitHub integration commands:")
		fmt.Println("  factory github --login   Authenticate with GitHub")
		fmt.Println("  factory github --status  Show connection status")
	},
}

var llmCmd = &cobra.Command{
	Use:   "llm",
	Short: "LLM provider management",
	Run: func(cmd *cobra.Command, args []string) {
		status, _ := cmd.Flags().GetBool("status")
		setup, _ := cmd.Flags().GetBool("setup")

		if status || (!status && !setup) {
			fmt.Println("LLM Status:")
			fmt.Println("  Checking Ollama... (localhost:11434)")
			fmt.Println("  Checking API keys...")
			fmt.Println("  (Full implementation coming soon)")
			return
		}
		if setup {
			fmt.Println("LLM Setup:")
			fmt.Println("  1. Install Ollama: https://ollama.ai")
			fmt.Println("  2. Or set OPENAI_API_KEY / ANTHROPIC_API_KEY")
		}
	},
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(intakeCmd)
	rootCmd.AddCommand(reviewCmd)
	rootCmd.AddCommand(rescueCmd)
	rootCmd.AddCommand(changeOrderCmd)
	rootCmd.AddCommand(githubCmd)
	rootCmd.AddCommand(llmCmd)

	// Add flags
	initCmd.Flags().Bool("quick", false, "Quick start with defaults")
	intakeCmd.Flags().String("name", "", "Project name")
	reviewCmd.Flags().String("spec", "", "Specification file to review against")
	rescueCmd.Flags().String("path", ".", "Path to codebase")
	githubCmd.Flags().Bool("login", false, "Authenticate with GitHub")
	githubCmd.Flags().Bool("status", false, "Show GitHub connection status")
	llmCmd.Flags().Bool("status", false, "Show LLM status")
	llmCmd.Flags().Bool("setup", false, "Setup LLM provider")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
