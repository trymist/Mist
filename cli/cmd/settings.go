package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/corecollectives/mist/models"
)

func HandleSettingsCommand(args []string) {
	if len(args) == 0 {
		printSettingsUsage()
		os.Exit(1)
	}

	subcommand := args[0]

	switch subcommand {
	case "get":
		getSettings(args[1:])
	case "set":
		setSetting(args[1:])
	case "list":
		listSettings(args[1:])
	case "help", "-h", "--help":
		printSettingsUsage()
	default:
		fmt.Printf("Unknown settings subcommand: %s\n\n", subcommand)
		printSettingsUsage()
		os.Exit(1)
	}
}

func printSettingsUsage() {
	fmt.Println("System Settings Management Commands")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  mist-cli settings <subcommand> [options]")
	fmt.Println()
	fmt.Println("Available Subcommands:")
	fmt.Println("  get    Get a specific setting value")
	fmt.Println("  set    Set a specific setting value")
	fmt.Println("  list   List all system settings")
	fmt.Println("  help   Show this help message")
	fmt.Println()
	fmt.Println("Available Settings Keys:")
	fmt.Println("  wildcard_domain          - Wildcard domain for auto-generated app domains")
	fmt.Println("  mist_app_name            - Subdomain name for Mist dashboard")
	fmt.Println("  production_mode          - Enable production mode (true/false)")
	fmt.Println("  secure_cookies           - Enable secure cookies for HTTPS (true/false)")
	fmt.Println("  auto_cleanup_containers  - Auto cleanup stopped containers (true/false)")
	fmt.Println("  auto_cleanup_images      - Auto cleanup dangling images (true/false)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  mist-cli settings list")
	fmt.Println("  mist-cli settings get --key wildcard_domain")
	fmt.Println("  mist-cli settings set --key wildcard_domain --value example.com")
	fmt.Println("  mist-cli settings set --key production_mode --value true")
}

func getSettings(args []string) {
	fs := flag.NewFlagSet("get", flag.ExitOnError)
	key := fs.String("key", "", "Setting key (required)")
	fs.Parse(args)

	if *key == "" {
		fmt.Println("Error: --key is required")
		fmt.Println()
		printSettingsUsage()
		os.Exit(1)
	}

	// Initialize database
	if err := initDB(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	settings, err := models.GetSystemSettings()
	if err != nil {
		fmt.Printf("Error fetching settings: %v\n", err)
		os.Exit(1)
	}

	var value interface{}
	found := true

	switch *key {
	case "wildcard_domain":
		if settings.WildcardDomain != nil {
			value = *settings.WildcardDomain
		} else {
			value = ""
		}
	case "mist_app_name":
		value = settings.MistAppName
	case "production_mode":
		value = settings.ProductionMode
	case "secure_cookies":
		value = settings.SecureCookies
	case "auto_cleanup_containers":
		value = settings.AutoCleanupContainers
	case "auto_cleanup_images":
		value = settings.AutoCleanupImages
	default:
		found = false
	}

	if !found {
		fmt.Printf("Error: Unknown setting key '%s'\n", *key)
		fmt.Println()
		printSettingsUsage()
		os.Exit(1)
	}

	fmt.Printf("%s: %v\n", *key, value)
}

func setSetting(args []string) {
	fs := flag.NewFlagSet("set", flag.ExitOnError)
	key := fs.String("key", "", "Setting key (required)")
	value := fs.String("value", "", "Setting value (required)")
	fs.Parse(args)

	if *key == "" || *value == "" {
		fmt.Println("Error: --key and --value are required")
		fmt.Println()
		printSettingsUsage()
		os.Exit(1)
	}

	// Initialize database
	if err := initDB(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	settings, err := models.GetSystemSettings()
	if err != nil {
		fmt.Printf("Error fetching settings: %v\n", err)
		os.Exit(1)
	}

	switch *key {
	case "wildcard_domain":
		if *value == "" {
			settings.WildcardDomain = nil
		} else {
			settings.WildcardDomain = value
		}
	case "mist_app_name":
		settings.MistAppName = *value
	case "production_mode":
		settings.ProductionMode = (*value == "true" || *value == "1")
	case "secure_cookies":
		settings.SecureCookies = (*value == "true" || *value == "1")
	case "auto_cleanup_containers":
		settings.AutoCleanupContainers = (*value == "true" || *value == "1")
	case "auto_cleanup_images":
		settings.AutoCleanupImages = (*value == "true" || *value == "1")
	default:
		fmt.Printf("Error: Unknown setting key '%s'\n", *key)
		fmt.Println()
		printSettingsUsage()
		os.Exit(1)
	}

	if err := settings.UpdateSystemSettings(); err != nil {
		fmt.Printf("Error updating settings: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Setting '%s' updated to '%s'\n", *key, *value)
}

func listSettings(args []string) {
	// Initialize database
	if err := initDB(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	settings, err := models.GetSystemSettings()
	if err != nil {
		fmt.Printf("Error fetching settings: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("System Settings:")
	fmt.Println("----------------------------------------------")
	fmt.Printf("%-30s %s\n", "Setting", "Value")
	fmt.Println("----------------------------------------------")

	wildcardDomain := ""
	if settings.WildcardDomain != nil {
		wildcardDomain = *settings.WildcardDomain
	}

	fmt.Printf("%-30s %s\n", "wildcard_domain", wildcardDomain)
	fmt.Printf("%-30s %s\n", "mist_app_name", settings.MistAppName)
	fmt.Printf("%-30s %v\n", "production_mode", settings.ProductionMode)
	fmt.Printf("%-30s %v\n", "secure_cookies", settings.SecureCookies)
	fmt.Printf("%-30s %v\n", "auto_cleanup_containers", settings.AutoCleanupContainers)
	fmt.Printf("%-30s %v\n", "auto_cleanup_images", settings.AutoCleanupImages)
	fmt.Println("----------------------------------------------")
}
