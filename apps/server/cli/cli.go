package cli

import (
	"flag"
	"fmt"
	"os"
)

const version = "1.0.0"

func Run(args []string) {
	flag.Usage = printUsage

	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	switch args[0] {
	case "user":
		HandleUserCommand(args[1:])
	case "settings":
		HandleSettingsCommand(args[1:])
	case "version":
		fmt.Printf("Mist CLI version %s\n", version)
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n\n", args[0])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Mist CLI - Command-line tool for managing Mist")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  mist-cli <command> [options]")
	fmt.Println()
	fmt.Println("Available Commands:")
	fmt.Println("  user        Manage users")
	fmt.Println("  settings    Manage system settings")
	fmt.Println("  version     Show CLI version")
	fmt.Println("  help        Show this help message")
	fmt.Println()
	fmt.Println("Use 'mist-cli <command> --help' for more information about a command")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  mist-cli user change-password --username admin")
	fmt.Println("  mist-cli settings get --key wildcard_domain")
	fmt.Println("  mist-cli settings set --key wildcard_domain --value example.com")
}
