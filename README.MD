# Gen Index

Gen-Index is a simple command-line tool written in Go to automatically generate an index.ts file for TypeScript or React projects.

## Description

This tool scans a specified directory and creates an index.ts file containing export statements for all .ts and .tsx files in that directory (excluding the index.tsx file itself).

## Usage

1. Compile the program:

   ```
   go build
   ```

2. Run the program:

   ```
   ./gen-index <directory_path>
   ```

   Replace `<directory_path>` with the path to the directory where you want to create the index.ts file.

## Features

- Automatically generates an index.ts file in the specified directory
- Supports both .ts and .tsx files
- Ignores index.tsx if present
- Creates relative export statements for each file
