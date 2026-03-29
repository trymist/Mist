package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/corecollectives/mist/models"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

func HandleUserCommand(args []string) {
	if len(args) == 0 {
		printUserUsage()
		os.Exit(1)
	}

	subcommand := args[0]

	switch subcommand {
	case "change-password":
		changePassword(args[1:])
	case "list":
		listUsers(args[1:])
	case "help", "-h", "--help":
		printUserUsage()
	default:
		fmt.Printf("Unknown user subcommand: %s\n\n", subcommand)
		printUserUsage()
		os.Exit(1)
	}
}

func printUserUsage() {
	fmt.Println("User Management Commands")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  mist-cli user <subcommand> [options]")
	fmt.Println()
	fmt.Println("Available Subcommands:")
	fmt.Println("  change-password   Change a user's password")
	fmt.Println("  list              List all users")
	fmt.Println("  help              Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  mist-cli user change-password --username admin")
	fmt.Println("  mist-cli user change-password --username admin --password newpass123")
	fmt.Println("  mist-cli user list")
}

func changePassword(args []string) {
	fs := flag.NewFlagSet("change-password", flag.ExitOnError)
	username := fs.String("username", "", "Username (required)")
	password := fs.String("password", "", "New password (if not provided, will prompt)")
	fs.Parse(args)

	if *username == "" {
		fmt.Println("Error: --username is required")
		fmt.Println()
		printUserUsage()
		os.Exit(1)
	}

	// Initialize database
	if err := initDB(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Check if user exists
	user, err := models.GetUserByUsername(*username)
	if err != nil {
		fmt.Printf("Error: User '%s' not found\n", *username)
		os.Exit(1)
	}

	var newPassword string

	if *password == "" {
		fmt.Print("Enter new password: ")
		bytePwd, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Printf("\nError reading password: %v\n", err)
			os.Exit(1)
		}
		fmt.Println()

		fmt.Print("Confirm new password: ")
		byteConfirm, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Printf("\nError reading password: %v\n", err)
			os.Exit(1)
		}
		fmt.Println()

		if string(bytePwd) != string(byteConfirm) {
			fmt.Println("Error: Passwords do not match")
			os.Exit(1)
		}

		newPassword = string(bytePwd)
	} else {
		newPassword = *password
	}

	// WARNING: password validation isn't implemented yet, will be implemented in future releases
	// if len(newPassword) < 8 {
	// 	fmt.Println("Error: Password must be at least 8 characters long")
	// 	os.Exit(1)
	// }

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Error hashing password: %v\n", err)
		os.Exit(1)
	}

	if err := models.UpdateUserPassword(user.ID, string(hashedPassword)); err != nil {
		fmt.Printf("Error updating password: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Password changed successfully for user '%s'\n", *username)
}

func listUsers(args []string) {
	if err := initDB(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	users, err := models.GetAllUsers()
	if err != nil {
		fmt.Printf("Error fetching users: %v\n", err)
		os.Exit(1)
	}

	if len(users) == 0 {
		fmt.Println("No users found")
		return
	}

	fmt.Println("Users:")
	fmt.Println("----------------------------------------------")
	fmt.Printf("%-5s %-20s %-30s %-10s\n", "ID", "Username", "Email", "Role")
	fmt.Println("----------------------------------------------")

	for _, user := range users {
		email := user.Email
		if email == "" {
			email = "N/A"
		}
		fmt.Printf("%-5d %-20s %-30s %-10s\n", user.ID, user.Username, email, user.Role)
	}
	fmt.Println("----------------------------------------------")
	fmt.Printf("Total: %d users\n", len(users))
}

func promptConfirm(message string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s (y/N): ", message)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}
