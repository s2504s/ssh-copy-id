package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// Define the flag for specifying the SSH key path with -i
	keyPath := flag.String("i", filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa.pub"), "Path to the public SSH key")

	// Parse the command-line flags
	flag.Parse()

	// Check if the remote user@host is provided
	if len(flag.Args()) < 1 {
		log.Fatal("Usage: ssh-copy-id -i /path/to/id_rsa.pub user@host")
	}

	// Get the remote server from the arguments
	remote := flag.Args()[0]

	// Expand ~ to the user home directory if needed
	if (*keyPath)[:2] == "~/" {
		*keyPath = filepath.Join(os.Getenv("HOME"), (*keyPath)[2:])
	}

	// Check if public key exists
	if _, err := os.Stat(*keyPath); os.IsNotExist(err) {
		log.Fatalf("Public key not found at %s. Provide a valid path.", *keyPath)
	}

	// Read public key
	pubKey, err := os.ReadFile(*keyPath)
	if err != nil {
		log.Fatalf("Failed to read public key: %v", err)
	}

	// Define the SSH command to copy the key
	cmd := exec.Command("ssh", remote, "mkdir -p ~/.ssh && chmod 700 ~/.ssh && echo '"+string(pubKey)+"' >> ~/.ssh/authorized_keys && chmod 600 ~/.ssh/authorized_keys")

	// Start the command
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Copying public key to %s...\n", remote)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Failed to copy key: %v", err)
	}

	fmt.Println("Public key copied successfully!")
}
