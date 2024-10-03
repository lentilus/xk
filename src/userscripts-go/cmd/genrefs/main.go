package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"xk/src/userscripts-go/pkg/api"
	"xk/src/userscripts-go/pkg/treesitter"

	sitter "github.com/smacker/go-tree-sitter"
)

const logMaxLines = 100 // Set the maximum number of log lines to keep

// Function to truncate the log file to the last N lines
func truncateLogFile(logFilePath string, maxLines int) {
	file, err := os.Open(logFilePath)
	if err != nil {
		log.Fatalf("Error opening log file for truncation: %v", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading log file: %v", err)
	}

	// Keep only the last `maxLines` lines
	if len(lines) > maxLines {
		lines = lines[len(lines)-maxLines:]
	}

	// Write the truncated lines back to the log file
	err = ioutil.WriteFile(logFilePath, []byte(strings.Join(lines, "\n")+"\n"), 0644)
	if err != nil {
		log.Fatalf("Error writing to log file: %v", err)
	}
}

func panicWithLog(msg string, args ...interface{}) {
	log.Printf(msg, args...)
	panic(fmt.Sprintf(msg, args...))
}

func main() {
	// Add a command-line flag for the Zettel name
	zettelName := flag.String("z", "", "Name of the Zettel to extract references from")
	flag.Parse()

	// Check if the Zettel name is provided
	if *zettelName == "" {
		panicWithLog("You must provide a Zettel name using the -z option.")
	}

	// Use xk API to get the path to the Zettel
	zettelPaths, err := api.Xk("path", map[string]string{"z": *zettelName})
	if err != nil {
		panicWithLog("Error fetching Zettel path: %v", err)
	}

	// Check if we have any paths; if not, throw an error
	if len(zettelPaths) == 0 {
		panicWithLog("Error: No paths returned for the given Zettel name.")
	}

	// Use only the first path
	zettelPath := zettelPaths[0]

	// Define paths for zettel.tex and references.log
	texFilePath := filepath.Join(zettelPath, "zettel.tex")
	referencesFilePath := filepath.Join(zettelPath, "references")
	logFilePath := filepath.Join(zettelPath, "references.log")

	// Open or create the log file
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panicWithLog("Error opening log file: %v", err)
	}
	defer logFile.Close()

	// Set log output to the log file
	log.SetOutput(logFile)

	// Truncate the log file to keep only the last N lines
	truncateLogFile(logFilePath, logMaxLines)

	// Read the content of zettel.tex
	source, err := ioutil.ReadFile(texFilePath)
	if err != nil {
		panicWithLog("Error reading zettel.tex file: %v", err)
	}

	// Initialize the parser for LaTeX
	parser := sitter.NewParser()
	defer parser.Close()
	lang := sitter.NewLanguage(treesitter.Language())
	parser.SetLanguage(lang)

	// Parse the source code (LaTeX content)
	tree := parser.Parse(nil, source)
	defer tree.Close()

	// Get the root node of the parsed tree
	rootNode := tree.RootNode()

	// Extract commands \zref and \zinc
	refCommands := append(
		treesitter.FindGenericCommand(rootNode, source, "zref"),
		treesitter.FindGenericCommand(rootNode, source, "zinc")...,
	)

	// Store references in a map to avoid duplicates
	refs := map[string]bool{}
	for _, c := range refCommands {
		start := c.ArgumentNode.StartByte() + 1
		end := c.ArgumentNode.EndByte() - 1
		if !(start < end) {
			continue
		}
		arg := string(source[start:end])

		// Validate the reference by checking if it can be found with `xk path`
		_, err := api.Xk("path", map[string]string{"z": arg})
		if err != nil {
			// Log the error but continue with the next reference
			log.Printf("Invalid reference %s: %v", arg, err)
			continue
		}
		refs[arg] = true
	}

	// Create a slice from the map keys and sort them alphabetically
	var sortedRefs []string
	for ref := range refs {
		sortedRefs = append(sortedRefs, ref)
	}
	sort.Strings(sortedRefs) // Alphabetically sort the references

	// Create a temporary file to store the references
	tmpFile, err := ioutil.TempFile("", "references-*.tmp")
	if err != nil {
		panicWithLog("Error creating temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Ensure the temp file is removed after use

	// Write the sorted references into the temporary file, one per line
	writer := bufio.NewWriter(tmpFile)
	for _, ref := range sortedRefs {
		_, err := writer.WriteString(ref + "\n")
		if err != nil {
			panicWithLog("Error writing to temporary file: %v", err)
		}
	}
	writer.Flush()
	tmpFile.Close()

	// Check if the references file exists
	if _, err := os.Stat(referencesFilePath); os.IsNotExist(err) {
		// If it doesn't exist, create it
		log.Println("References file does not exist, creating it.")
		if _, err := os.Create(referencesFilePath); err != nil {
			panicWithLog("Error creating references file: %v", err)
		}
	}

	// Show a diff between the old references file and the temporary file
	diffCmd := exec.Command("diff", "-u", referencesFilePath, tmpFile.Name())
	diffOutput, err := diffCmd.CombinedOutput()
	if err != nil &&
		err.Error() != "exit status 1" { // exit status 1 means diff found differences, not an actual error
		panicWithLog("Error running diff: %v", err)
	}
	if len(diffOutput) > 0 {
		log.Println("Changes in references:")
		log.Println(string(diffOutput))
	} else {
		log.Println("No changes in references.")
	}

	// Overwrite the references file with the temporary file contents
	tempContents, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		panicWithLog("Error reading temporary file: %v", err)
	}
	err = ioutil.WriteFile(referencesFilePath, tempContents, 0644)
	if err != nil {
		panicWithLog("Error writing to references file: %v", err)
	}

	log.Println("References file updated successfully.")
}
