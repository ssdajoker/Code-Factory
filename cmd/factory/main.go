package main

import (
        "flag"
        "fmt"
        "os"

        "github.com/ssdajoker/Code-Factory/internal/tui"
)

const version = "1.0.0"

func main() {
        // Global flags
        showVersion := flag.Bool("version", false, "Show version")
        showHelp := flag.Bool("help", false, "Show help")
        flag.BoolVar(showVersion, "v", false, "Show version (shorthand)")
        flag.BoolVar(showHelp, "h", false, "Show help (shorthand)")

        // Custom usage
        flag.Usage = printHelp

        // Parse global flags only if no subcommand
        if len(os.Args) < 2 {
                // No args: launch TUI
                launchTUI()
                return
        }

        // Check for subcommand
        cmd := os.Args[1]

        // Handle global flags first
        if cmd == "--version" || cmd == "-v" {
                fmt.Printf("Factory v%s\n", version)
                return
        }
        if cmd == "--help" || cmd == "-h" || cmd == "help" {
                printHelp()
                return
        }

        // Subcommand dispatch
        switch cmd {
        case "init":
                cmdInit(os.Args[2:])
        case "intake":
                cmdIntake(os.Args[2:])
        case "review":
                cmdReview(os.Args[2:])
        case "rescue":
                cmdRescue(os.Args[2:])
        case "change-order":
                cmdChangeOrder(os.Args[2:])
        case "github":
                cmdGitHub(os.Args[2:])
        case "llm":
                cmdLLM(os.Args[2:])
        case "version":
                fmt.Printf("Factory v%s\n", version)
        default:
                fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", cmd)
                printHelp()
                os.Exit(1)
        }
}

func launchTUI() {
        if err := tui.Run(); err != nil {
                fmt.Fprintf(os.Stderr, "Error: %v\n", err)
                os.Exit(1)
        }
}

func cmdInit(args []string) {
        fs := flag.NewFlagSet("init", flag.ExitOnError)
        quick := fs.Bool("quick", false, "Quick start with defaults")
        fs.Parse(args)

        fmt.Println("Running factory init...")
        if *quick {
                fmt.Println("Using quick start mode")
        }
        fmt.Println("(Implementation coming soon)")
}

func cmdIntake(args []string) {
        fs := flag.NewFlagSet("intake", flag.ExitOnError)
        name := fs.String("name", "", "Project name")
        fs.Parse(args)

        fmt.Println("Starting INTAKE mode...")
        if *name != "" {
                fmt.Printf("Project: %s\n", *name)
        }
        
        // Launch intake TUI
        if err := tui.RunIntake(); err != nil {
                fmt.Fprintf(os.Stderr, "Error: %v\n", err)
                os.Exit(1)
        }
}

func cmdReview(args []string) {
        fs := flag.NewFlagSet("review", flag.ExitOnError)
        spec := fs.String("spec", "", "Specification file to review against")
        fs.Parse(args)

        fmt.Println("Starting REVIEW mode...")
        if *spec != "" {
                fmt.Printf("Spec: %s\n", *spec)
        }
        fmt.Println("(Implementation coming soon)")
}

func cmdRescue(args []string) {
        fs := flag.NewFlagSet("rescue", flag.ExitOnError)
        path := fs.String("path", ".", "Path to codebase")
        fs.Parse(args)

        fmt.Println("Starting RESCUE mode...")
        fmt.Printf("Path: %s\n", *path)
        fmt.Println("(Implementation coming soon)")
}

func cmdChangeOrder(args []string) {
        fs := flag.NewFlagSet("change-order", flag.ExitOnError)
        fs.Parse(args)

        fmt.Println("Starting CHANGE_ORDER mode...")
        fmt.Println("(Implementation coming soon)")
}

func cmdGitHub(args []string) {
        fs := flag.NewFlagSet("github", flag.ExitOnError)
        login := fs.Bool("login", false, "Authenticate with GitHub")
        status := fs.Bool("status", false, "Show GitHub connection status")
        fs.Parse(args)

        if *login {
                fmt.Println("Starting GitHub OAuth flow...")
                fmt.Println("(Implementation coming soon)")
                return
        }
        if *status {
                fmt.Println("GitHub connection status: Not connected")
                return
        }
        fmt.Println("GitHub integration commands:")
        fmt.Println("  factory github --login   Authenticate with GitHub")
        fmt.Println("  factory github --status  Show connection status")
}

func cmdLLM(args []string) {
        fs := flag.NewFlagSet("llm", flag.ExitOnError)
        status := fs.Bool("status", false, "Show LLM status")
        setup := fs.Bool("setup", false, "Setup LLM provider")
        fs.Parse(args)

        if *status || len(args) == 0 {
                fmt.Println("LLM Status:")
                fmt.Println("  Checking Ollama... (localhost:11434)")
                fmt.Println("  Checking API keys...")
                fmt.Println("  (Full implementation coming soon)")
                return
        }
        if *setup {
                fmt.Println("LLM Setup:")
                fmt.Println("  1. Install Ollama: https://ollama.ai")
                fmt.Println("  2. Or set OPENAI_API_KEY / ANTHROPIC_API_KEY")
                return
        }
}

func printHelp() {
        fmt.Printf(`Factory v%s - Spec-Driven Software Factory

USAGE:
    factory [COMMAND] [FLAGS]

COMMANDS:
    init            Initialize Factory in current project
                    --quick    Quick start with defaults
    
    intake          Start INTAKE mode (capture vision)
                    --name     Project name
    
    review          Start REVIEW mode (check code against specs)
                    --spec     Specification file path
    
    change-order    Start CHANGE_ORDER mode (track drift)
    
    rescue          Start RESCUE mode (reverse-engineer codebase)
                    --path     Path to codebase (default: .)
    
    github          GitHub integration
                    --login    Authenticate with GitHub
                    --status   Show connection status
    
    llm             LLM provider management
                    --status   Show LLM status
                    --setup    Setup LLM provider
    
    version         Show version
    help            Show this help

GLOBAL FLAGS:
    --version, -v   Show version
    --help, -h      Show this help

EXAMPLES:
    factory                      # Start TUI
    factory init                 # Initialize Factory
    factory init --quick         # Quick initialization
    factory intake --name myapp  # Create new specification
    factory review --spec spec.md # Review code against spec
    factory github --login       # Connect to GitHub
    factory llm --status         # Check LLM availability

DOCUMENTATION:
    https://github.com/ssdajoker/Code-Factory

`, version)
}
