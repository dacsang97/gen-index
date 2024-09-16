package main

import (
	"bufio"
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

	existingExports := make(map[string]bool)
	validExports := []string{}

	if _, err := os.Stat(outputFile); err == nil {
		file, err := os.Open(outputFile)
		if err != nil {
			fmt.Printf("Error opening index.ts file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "export * from") {
				exportPath := strings.Trim(strings.TrimPrefix(line, "export * from"), "'\"")
				fullPath := filepath.Join(dir, exportPath)
				if _, err := os.Stat(fullPath + ".ts"); err == nil {
					existingExports[line] = true
					validExports = append(validExports, line)
				} else if _, err := os.Stat(fullPath + ".tsx"); err == nil {
					existingExports[line] = true
					validExports = append(validExports, line)
				}
				// Removed the original condition to check for both .ts and .tsx
			}
		}
	}

	newExports := []string{}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip if file or directory is inaccessible
		}

		if info.IsDir() {
			return nil
		}

		if (filepath.Ext(path) == ".tsx" || filepath.Ext(path) == ".ts") && filepath.Base(path) != "index.ts" {
			relPath, _ := filepath.Rel(dir, path)
			exportPath := "./" + strings.TrimSuffix(relPath, filepath.Ext(relPath))
			exportLine := fmt.Sprintf("export * from '%s'", exportPath)
			if !existingExports[exportLine] {
				newExports = append(newExports, exportLine)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		os.Exit(1)
	}

	allExports := append(validExports, newExports...)

	if len(allExports) > 0 {
		file, err := os.Create(outputFile)
		if err != nil {
			fmt.Printf("Error opening file for writing: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		content := strings.Join(allExports, "\n") + "\n"
		if _, err := file.WriteString(content); err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully updated %s\n", outputFile)
	} else {
		fmt.Printf("No valid exports for %s\n", outputFile)
	}
}
