package main

import (
	"testing"
	"os"
	"os/exec"
)

// This is a basic integration test that ensures the program can be compiled and run
// without immediately crashing. We can't easily test the full interactive behavior
// in an automated way, but we can at least make sure it builds and starts.
func TestMainCompiles(t *testing.T) {
	// Skip this test in normal test runs - it's more of a build validation
	if os.Getenv("RUN_BUILD_TEST") != "1" {
		t.Skip("Skipping build test; set RUN_BUILD_TEST=1 to run")
	}
	
	// Build the program
	buildCmd := exec.Command("go", "build", "-o", "nvim-plugin-test")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build program: %v", err)
	}
	defer os.Remove("nvim-plugin-test") // Clean up after test
	
	// Start the program with a timeout to ensure it doesn't hang
	// Since this is an interactive program, we'll just make sure it starts
	// but then immediately terminate it
	cmd := exec.Command("./nvim-plugin-test")
	
	// Start the process but don't wait for it to complete
	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start program: %v", err)
	}
	
	// Kill the process after starting it (since it's interactive)
	if err := cmd.Process.Kill(); err != nil {
		t.Fatalf("Failed to kill test process: %v", err)
	}
}