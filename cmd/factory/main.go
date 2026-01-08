package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version information (set during build)
	version   = "dev"
	buildTime = "unknown"
	commitSHA = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "factory",
	Short: "Code-Factory - Spec-Driven Software Development with AI",
	Long: `Code-Factory transforms natural language specifications into production-ready code.

Built with love by developers, for developers.

Features:
  • Beautiful Terminal UI with Charm.sh
  • AI-powered spec generation and code implementation
  • Works with local Ollama or cloud LLMs (OpenAI, Claude)
  • Git-native workflow with optional GitHub integration
  • Four core modes: INTAKE, REVIEW, CHANGE_ORDER, RESCUE

Get started:
  factory init        # One-click onboarding
  factory intake      # Create specifications
  factory review      # Analyze code against specs
  factory change-order # Implement changes
  factory rescue      # Debug and fix issues`,
	Version: version,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Code-Factory for the current project",
	Long: `Initialize Code-Factory with interactive onboarding.

This will:
  • Configure your LLM provider (Ollama or BYOK)
  • Optionally connect to GitHub
  • Set up project structure
  • Create initial templates`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🏭 Initializing Code-Factory...")
		fmt.Println("\nThis feature is coming soon!")
		fmt.Println("\nFor now, you can explore the specifications in:")
		fmt.Println("  • contracts/factory_bootstrap_spec.md")
		fmt.Println("  • contracts/system_architecture.md")
		fmt.Println("  • contracts/mode_specifications.md")
		fmt.Println("  • docs/ARCHITECTURE.md")
	},
}

var intakeCmd = &cobra.Command{
	Use:   "intake",
	Short: "Create specifications from requirements",
	Long: `Transform natural language requirements into structured specifications.

Interactive mode guides you through requirement gathering, clarifying questions,
and AI-powered spec generation.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("📝 INTAKE Mode")
		fmt.Println("\nThis feature is under development.")
		fmt.Println("See contracts/mode_specifications.md for detailed specs.")
	},
}

var reviewCmd = &cobra.Command{
	Use:   "review [path]",
	Short: "Analyze code against specifications",
	Long: `Review existing codebase against specifications.

Generates comprehensive reports with:
  • Spec compliance analysis
  • Security vulnerabilities
  • Performance issues
  • Best practice violations
  • Actionable recommendations`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔍 REVIEW Mode")
		fmt.Println("\nThis feature is under development.")
		fmt.Println("See contracts/mode_specifications.md for detailed specs.")
	},
}

var changeOrderCmd = &cobra.Command{
	Use:   "change-order",
	Short: "Implement changes from specifications",
	Long: `Generate code changes based on specifications or change requests.

Workflow:
  • Analyze impact of requested change
  • Generate implementation plan
  • Create code with AI assistance
  • Review diffs and approve
  • Create feature branch and optional PR`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("⚙️  CHANGE_ORDER Mode")
		fmt.Println("\nThis feature is under development.")
		fmt.Println("See contracts/mode_specifications.md for detailed specs.")
	},
}

var rescueCmd = &cobra.Command{
	Use:   "rescue",
	Short: "Debug and fix issues",
	Long: `AI-powered debugging and problem solving.

Helps diagnose:
  • Test failures
  • Runtime errors
  • Performance issues
  • Build problems
  
Provides solutions and can apply fixes automatically.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚑 RESCUE Mode")
		fmt.Println("\nThis feature is under development.")
		fmt.Println("See contracts/mode_specifications.md for detailed specs.")
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Code-Factory configuration",
	Long: `View and modify Code-Factory configuration.

Subcommands:
  config           Show current configuration
  config llm       Configure LLM provider
  config github    Configure GitHub integration
  config reset     Reset to defaults`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("⚙️  Configuration")
		fmt.Println("\nThis feature is under development.")
		fmt.Println("\nConfiguration file location:")
		fmt.Println("  ~/.config/factory/config.yaml")
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Code-Factory %s\n", version)
		fmt.Printf("Built: %s\n", buildTime)
		fmt.Printf("Commit: %s\n", commitSHA)
	},
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(intakeCmd)
	rootCmd.AddCommand(reviewCmd)
	rootCmd.AddCommand(changeOrderCmd)
	rootCmd.AddCommand(rescueCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(versionCmd)

	// Global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().String("config", "", "Config file (default: ~/.config/factory/config.yaml)")

	// Init flags
	initCmd.Flags().Bool("no-questions", false, "Skip interactive questions")
	initCmd.Flags().String("llm", "", "LLM provider (ollama, openai, anthropic)")

	// Intake flags
	intakeCmd.Flags().String("prompt", "", "Initial requirement prompt")
	intakeCmd.Flags().String("file", "", "Load requirements from file")
	intakeCmd.Flags().String("output", "", "Output file path")
	intakeCmd.Flags().Bool("voice", false, "Use voice input")

	// Review flags
	reviewCmd.Flags().String("spec", "", "Specific spec file to review against")
	reviewCmd.Flags().StringSlice("focus", []string{}, "Focus areas (security, performance, etc.)")
	reviewCmd.Flags().String("format", "markdown", "Output format (markdown, json, html)")
	reviewCmd.Flags().String("output", "", "Output file path")

	// Change order flags
	changeOrderCmd.Flags().String("description", "", "Change description")
	changeOrderCmd.Flags().String("spec", "", "Spec file to implement")
	changeOrderCmd.Flags().Bool("dry-run", false, "Show plan only, don't apply changes")
	changeOrderCmd.Flags().Bool("auto-approve", false, "Auto-approve changes (dangerous!)")
	changeOrderCmd.Flags().Bool("no-pr", false, "Don't create pull request")

	// Rescue flags
	rescueCmd.Flags().String("problem", "", "Problem description")
	rescueCmd.Flags().String("error", "", "Error message or log")
	rescueCmd.Flags().Bool("auto-fix", false, "Automatically apply fix if confident")
	rescueCmd.Flags().Bool("explain-only", false, "Explain problem without fixing")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
