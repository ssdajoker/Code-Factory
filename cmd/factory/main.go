package main

import (
	"fmt"
	"os"
)

const version = "1.0.0"

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v":
			fmt.Printf("Factory v%s\n", version)
			return
		case "--help", "-h":
			printHelp()
			return
		case "init":
			fmt.Println("Running factory init...")
			fmt.Println("(Implementation coming soon)")
			return
		case "intake":
			fmt.Println("Starting INTAKE mode...")
			fmt.Println("(Implementation coming soon)")
			return
		case "review":
			fmt.Println("Starting REVIEW mode...")
			fmt.Println("(Implementation coming soon)")
			return
		case "change-order":
			fmt.Println("Starting CHANGE_ORDER mode...")
			fmt.Println("(Implementation coming soon)")
			return
		case "rescue":
			fmt.Println("Starting RESCUE mode...")
			fmt.Println("(Implementation coming soon)")
			return
		}
	}

	// Default: start TUI
	fmt.Println("Starting Factory TUI...")
	fmt.Println("(Implementation coming soon)")
}

func printHelp() {
	fmt.Printf(`Factory v%s - Spec-Driven Software Factory

USAGE:
    factory [COMMAND]

COMMANDS:
    init            Initialize Factory in current project
    intake          Start INTAKE mode (capture vision)
    review          Start REVIEW mode (check code against specs)
    change-order    Start CHANGE_ORDER mode (track drift)
    rescue          Start RESCUE mode (reverse-engineer codebase)
    
    --version, -v   Show version
    --help, -h      Show this help

EXAMPLES:
    factory init                 # Initialize Factory
    factory                      # Start TUI
    factory intake               # Create new specification
    factory review               # Review code against specs

DOCUMENTATION:
    https://github.com/ssdajoker/Code-Factory

`, version)
}
