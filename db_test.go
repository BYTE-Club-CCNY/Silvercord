package main

import (
	"os/exec"
	"strings"
	"testing"
	"fmt"
)

func TestPipeline(t *testing.T) {
	urls := []string{
		"https://www.ratemyprofessors.com/professor/2946510",
		"https://www.ratemyprofessors.com/professor/432142",
		"https://www.ratemyprofessors.com/professor/2581017",
		"https://www.ratemyprofessors.com/professor/422536",
		"https://www.ratemyprofessors.com/professor/2818033",
		"https://www.ratemyprofessors.com/professor/354797",
	}
	expectedOutput := "ChromaDB Store Successful!"

	for _, url := range urls {
		t.Run("Testing URL: "+url, func(t *testing.T) {
			cmd := exec.Command("python", "db_store.py", url)

			output, err := cmd.CombinedOutput() // <- CombinedOutput() function RUNS the cmd, and then gives us the stdout & the stderr
			fmt.Println(string(output))
			if err != nil {
				t.Fatalf("Failed to run pipeline for URL %s: %v", url, err)
			}

			outputStr := strings.TrimSpace(string(output))
			if outputStr != expectedOutput {
				t.Errorf("Unexpected output for URL %s:\nGot: %s\nWant: %s", url, outputStr, expectedOutput)
			}
		})
	}
}
