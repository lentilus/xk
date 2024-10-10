package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"xk/src/userscripts-go/pkg/api"
	"xk/src/userscripts-go/pkg/logging"
	"xk/src/userscripts-go/pkg/treesitter"

	sitter "github.com/smacker/go-tree-sitter"
)

type FlashCard struct {
	id    string
	front string
	back  string
}

// EnvToFlashcard extracts flashcards from the form
// \begin{flashcard}[<id>]{<question>} <content> \end{flashcard}
func EnvToFlashcard(env treesitter.GenericEnvironment, source []byte) (FlashCard, error) {
	// get id
	idNode := env.ArgumentNodes[0]
	if idNode.Type() != "brack_group" {
		return FlashCard{"", "", ""}, fmt.Errorf("flashcard is malformatted")
	}
	id := string(source[idNode.StartByte()+1 : idNode.EndByte()-1])

	// get front
	frontNode := env.EnvironmentNode.Child(0).NextSibling()
	if frontNode.Type() != "curly_group" {
		return FlashCard{"", "", ""}, fmt.Errorf("flashcard is malformatted")
	}
	front := string(source[frontNode.StartByte()+1 : frontNode.EndByte()-1])

	// get back
	back := env.EnvironmentNode.Content(source)

	return FlashCard{id, front, back}, nil
}

// CheckAndRemoveObsoleteFiles removes files with IDs not matching the extracted ones
func CheckAndRemoveObsoleteFiles(validIDs map[string]bool, zettelDir string) error {
	// Get all files in the zettel directory matching the pattern
	files, err := filepath.Glob(filepath.Join(zettelDir, "card_*_*.tex"))
	if err != nil {
		return err
	}

	for _, file := range files {
		parts := strings.Split(file, "_")
		if len(parts) < 3 {
			continue
		}

		id := parts[1]
		if !validIDs[id] {
			log.Println("Removing obsolete file:", file)
			err := os.Remove(file)
			if err != nil {
				log.Printf("Error removing file %s: %v", file, err)
				return err
			}
		}
	}

	return nil
}

// CompareAndUpdateFile compares the current file content with the new content and updates if necessary
func CompareAndUpdateFile(filename, newContent string) error {
	if _, err := os.Stat(filename); err == nil {
		existingContent, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}

		if string(existingContent) == newContent {
			log.Println("No changes in:", filename)
			return nil
		}
	}

	log.Println("Updating file:", filename)
	return ioutil.WriteFile(filename, []byte(newContent), 0644)
}

// SaveToFile saves the LaTeX content to a .tex file, comparing it with the existing content
func SaveToFile(filename, preamble, content string) error {
	texContent := preamble + "\\begin{document}\n" + content + "\n\\end{document}"
	return CompareAndUpdateFile(filename, texContent)
}

func main() {
	// Add command-line flag for Zettel name
	zettelName := flag.String("z", "", "Name of the Zettel to extract flashcards from")
	flag.Parse()

	// Validate that the Zettel name was provided
	if *zettelName == "" {
		logging.PanicWithLog("You must provide a Zettel name using the -z option.")
	}

	// Fetch the Zettel path using the xk API
	zettelPaths, err := api.Xk("path", map[string]string{"z": *zettelName})
	if err != nil {
		logging.PanicWithLog("Error fetching Zettel path: %v", err)
	}

	if len(zettelPaths) == 0 {
		logging.PanicWithLog("No paths returned for the provided Zettel name.")
	}

	// Use the first path and determine the directory of zettel.tex
	zettelPath := zettelPaths[0]
	texFilePath := filepath.Join(zettelPath, "zettel.tex")
	logFilePath := filepath.Join(zettelPath, "flashcards.log")

	// Get the directory of zettel.tex
	zettelDir := filepath.Dir(texFilePath)

	// Set up logging to file
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
		logging.PanicWithLog("Error reading zettel.tex: %v", err)
	}

	// Initialize parser for LaTeX
	parser := sitter.NewParser()
	defer parser.Close()
	lang := sitter.NewLanguage(treesitter.Language())
	parser.SetLanguage(lang)

	// Parse LaTeX content
	tree := parser.Parse(nil, source)
	defer tree.Close()

	// Find document environment
	rootNode := tree.RootNode()
	documentEnv := treesitter.FindGenericEnvironment(rootNode, source, "document")
	if len(documentEnv) == 0 {
		logging.PanicWithLog("No document environment found in zettel.tex")
	}

	preamble := string(source[:documentEnv[0].EnvironmentNode.StartByte()])

	// Find all flashcard environments
	cardEnvs := treesitter.FindGenericEnvironment(rootNode, source, "flashcard")
	var flashcards []FlashCard
	validIDs := make(map[string]bool)

	for _, env := range cardEnvs {
		card, err := EnvToFlashcard(env, source)
		if err != nil {
			log.Printf("Error parsing flashcard: %v", err)
			continue
		}
		flashcards = append(flashcards, card)
		validIDs[card.id] = true
	}

	// Remove obsolete files in the zettel directory
	err = CheckAndRemoveObsoleteFiles(validIDs, zettelDir)
	if err != nil {
		logging.PanicWithLog("Error checking obsolete files: %v", err)
	}

	// Save front and back of flashcards to .tex files in the zettel directory
	for _, card := range flashcards {
		frontFile := filepath.Join(zettelDir, fmt.Sprintf("card_%s_front.tex", card.id))
		backFile := filepath.Join(zettelDir, fmt.Sprintf("card_%s_back.tex", card.id))

		if err := SaveToFile(frontFile, preamble, card.front); err != nil {
			log.Printf("Error saving front of card %s: %v", card.id, err)
			continue
		}

		if err := SaveToFile(backFile, preamble, card.back); err != nil {
			log.Printf("Error saving back of card %s: %v", card.id, err)
			continue
		}
	}

	log.Println("Flashcards processed successfully.")
}
