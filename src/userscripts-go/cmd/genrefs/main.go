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
	"xk/src/userscripts-go/pkg/api"
	"xk/src/userscripts-go/pkg/logging"
	"xk/src/userscripts-go/pkg/treesitter"

	sitter "github.com/smacker/go-tree-sitter"
)

// Array to hold the LaTeX commands to extract as references
var latexCommands = []string{"cite"}

func main() {
	// Add a command-line flag for the Zettel name
	zettelName := flag.String("z", "", "Name of the Zettel to extract references from")
	refQuery, set := os.LookupEnv("TS_QUERY_REF")
	if !set || refQuery == "" {
		logging.PanicWithLog("Treesitter query for references not configured.")
	}

	flag.Parse()

	// Check if the Zettel name is provided
	if *zettelName == "" {
		logging.PanicWithLog("You must provide a Zettel name using the -z option.")
	}

	// Use xk API to get the path to the Zettel
	zettelPaths, err := api.Xk("path", map[string]string{"z": *zettelName})
	if err != nil {
		logging.PanicWithLog("Error fetching Zettel path: %v", err)
	}

	// Check if we have any paths; if not, throw an error
	if len(zettelPaths) == 0 {
		logging.PanicWithLog("Error: No paths returned for the given Zettel name.")
	}

	// Use only the first path
	zettelPath := zettelPaths[0]

	// Define paths for zettel.tex and references.log
	texFilePath := filepath.Join(zettelPath, "zettel.tex")
	referencesFilePath := filepath.Join(zettelPath, "references")
	logFilePath := filepath.Join(zettelPath, "references.log")

	// Open or create the log file and set log output
	logFile, err := logging.SetLogOutput(logFilePath)
	if err != nil {
		logging.PanicWithLog("Error setting log output: %v", err)
	}
	defer logFile.Close()

	// Truncate the log file to keep only the last N lines
	err = logging.TruncateLogFile(logFilePath, logging.LogMaxLines)
	if err != nil {
		logging.PanicWithLog("Error truncating log file: %v", err)
	}

	// Read the content of zettel.tex
	source, err := ioutil.ReadFile(texFilePath)
	if err != nil {
		logging.PanicWithLog("Error reading zettel.tex file: %v", err)
	}

	// Initialize the parser for LaTeX
	parser := sitter.NewParser()
	defer parser.Close()
	lang := sitter.NewLanguage(treesitter.Language())
	parser.SetLanguage(lang)

	// Parse the source code (LaTeX content)
	tree := parser.Parse(nil, source)
	defer tree.Close()

	// Query the tree
	query, err := sitter.NewQuery([]byte(refQuery), lang)
	if err != nil {
		panic(err)
	}
	cursor := sitter.NewQueryCursor()
	cursor.Exec(query, tree.RootNode())

	refs := map[string]bool{}
	for {
		m, ok := cursor.NextMatch()
		if !ok {
			break
		}
		// Apply predicates filtering
		m = cursor.FilterPredicates(m, source)
		for _, c := range m.Captures {
			ref := c.Node.Content(source)

			// strip brackets if necessary
			sref := ref[1 : len(ref)-1]
			if ref == fmt.Sprintf("{%s}", sref) || ref == fmt.Sprintf("[%s]", sref) {
				ref = sref
			}

			// validate zettels existence
			_, err := api.Xk("path", map[string]string{"z": ref})
			if err != nil {
				// Log the error but continue with the next reference
				log.Printf("Invalid reference %s: %v", ref, err)
				continue
			}
			refs[ref] = true
		}
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
		logging.PanicWithLog("Error creating temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Ensure the temp file is removed after use

	// Write the sorted references into the temporary file, one per line
	writer := bufio.NewWriter(tmpFile)
	for _, ref := range sortedRefs {
		_, err := writer.WriteString(ref + "\n")
		if err != nil {
			logging.PanicWithLog("Error writing to temporary file: %v", err)
		}
	}
	writer.Flush()
	tmpFile.Close()

	// Check if the references file exists
	if _, err := os.Stat(referencesFilePath); os.IsNotExist(err) {
		// If it doesn't exist, create it
		log.Println("References file does not exist, creating it.")
		if _, err := os.Create(referencesFilePath); err != nil {
			logging.PanicWithLog("Error creating references file: %v", err)
		}
	}

	// Show a diff between the old references file and the temporary file
	diffCmd := exec.Command("diff", "-u", referencesFilePath, tmpFile.Name())
	diffOutput, err := diffCmd.CombinedOutput()
	if err != nil &&
		err.Error() != "exit status 1" { // exit status 1 means diff found differences, not an actual error
		logging.PanicWithLog("Error running diff: %v", err)
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
		logging.PanicWithLog("Error reading temporary file: %v", err)
	}
	err = ioutil.WriteFile(referencesFilePath, tempContents, 0644)
	if err != nil {
		logging.PanicWithLog("Error writing to references file: %v", err)
	}

	log.Println("References file updated successfully.")
}
