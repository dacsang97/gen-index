package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gen-index <directory>")
		os.Exit(1)
	}

	dir := os.Args[1]
	outputFile := filepath.Join(dir, "index.ts")

	exports := []string{}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if (filepath.Ext(path) == ".tsx" || filepath.Ext(path) == ".ts") && filepath.Base(path) != "index.tsx" {
			relPath, _ := filepath.Rel(dir, path)
			exportPath := "./" + strings.TrimSuffix(relPath, filepath.Ext(relPath))
			exports = append(exports, fmt.Sprintf("export * from '%s'", exportPath))
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		os.Exit(1)
	}

	content := strings.Join(exports, "\n")

	err = os.WriteFile(outputFile, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated %s\n", outputFile)
}
