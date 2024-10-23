package main

import (
	"crypto"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
)

// Helper function: add a new flashcard to Anki
func addNewCard(flashcard Flashcard, frontSVG, backSVG string) {
	// Define fields for the new flashcard
	fields := map[string]string{
		"front": frontSVG,
		"back":  backSVG,
		"id":    flashcard.ID,
		"hash":  flashcard.Hash,
		"fixme": "",
	}

	// Add the card to the deck
	_, err := AddCard(&connect, deck, modelName, fields)
	if err != nil {
		log.Fatalf("Failed to add new flashcard: %v", err)
	}
	log.Printf("Added new flashcard with ID: %s", flashcard.ID)
}

// Tex2Base64 compiles LaTeX content into a base64-encoded SVG (or PDF) string using temporary files.
func Tex2Base64(texPath string) (string, error) {
	log.Printf("Compiling %s", texPath)

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "tex2base64-")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up the temporary directory after use

	if err := os.Chmod(tempDir, 0766); err != nil {
		return "", err
	}

	// Compile LaTeX to PDF using latexmk
	latexmkCmd := exec.Command(
		"latexmk",
		"-f",
		"-pdf",
		"-cd",
		"-outdir="+tempDir,
		"-jobname=pdfout", // strip .pdf extension
		texPath,
	)

	log.Printf("Writing output to %s", tempDir)

	compiled := filepath.Join(tempDir, "pdfout.pdf")
	cropped := filepath.Join(tempDir, "cropped.pdf")
	vector := filepath.Join(tempDir, "svgout.svg")

	log.Printf("Trying to write to %s", compiled)

	err = latexmkCmd.Run()
	if err != nil {
		// We are in force mode so latexmk is still trying to create the pdf
		log.Printf("error running latexmk: %v", err)
	}

	// Check if the PDF was successfully created
	if _, err := os.Stat(compiled); os.IsNotExist(err) {
		return "", fmt.Errorf("PDF not generated: %v", err)
	}

	// Convert the PDF to SVG using pdf2svg
	pdfcropCmd := exec.Command("pdfcrop", compiled, cropped)
	err = pdfcropCmd.Run()
	if err != nil {
		return "", fmt.Errorf("error running pdfcrop: %v", err)
	}

	// Convert the PDF to SVG using pdf2svg
	pdf2svgCmd := exec.Command("pdf2svg", cropped, vector)
	err = pdf2svgCmd.Run()
	if err != nil {
		return "", fmt.Errorf("error running pdf2svg: %v", err)
	}

	// Check if the SVG was successfully created
	if _, err := os.Stat(vector); os.IsNotExist(err) {
		return "", fmt.Errorf("SVG not generated: %v", err)
	}

	// Read the SVG file and encode it as base64
	svgBytes, err := os.ReadFile(vector)
	if err != nil {
		return "", fmt.Errorf("error reading SVG file: %v", err)
	}

	base64SVG := base64.StdEncoding.EncodeToString(svgBytes)

	// Return the Base64-encoded SVG content
	return base64SVG, nil
}

// Helper function: update an existing flashcard in Anki
func updateCard(cardID int, flashcard Flashcard, frontSVG string, backSVG string) {
	// Define fields to be updated
	fields := map[string]string{
		"front": frontSVG,
		"back":  backSVG,
		"hash":  flashcard.Hash,
	}

	// Update the card
	noteID, err := Card2Note(&connect, cardID)
	if err != nil {
		log.Fatalf("Failed to convert card to note: %v", err)
	}

	_, err = UpdateNoteFields(&connect, noteID, fields)
	if err != nil {
		log.Fatalf("Failed to update flashcard: %v", err)
	}
	log.Printf("Updated flashcard with ID: %s", flashcard.ID)
}

func FindFixme() ([]int, error) {
	var res GenericResponse[[]int]

	params := map[string]any{
		"query": fmt.Sprintf("deck:%s fixme:_*", deck),
	}

	err := connect.Request("findCards", params, &res)
	if err != nil {
		return []int{}, err
	}

	if err := checkAPIError(res.Error); err != nil {
		return []int{}, err
	}

	return res.Result, nil
}

// Helper function: check if the deck exists in the list of decks
func deckExists(decks []string, deck string) bool {
	for _, d := range decks {
		if d == deck {
			return true
		}
	}
	return false
}

func Hash(objs ...interface{}) []byte {
	digester := crypto.MD5.New()
	for _, ob := range objs {
		fmt.Fprint(digester, reflect.TypeOf(ob))
		fmt.Fprint(digester, ob)
	}
	return digester.Sum(nil)
}

// Card2Zettel finds a file in the given directory that matches the pattern <kastenPath>/*/card_<id>_front.tex
func Card2Zettel(kastenPath string, cardID string) (string, error) {
	// Define the pattern for the file search
	pattern := filepath.Join(kastenPath, "*", fmt.Sprintf("card_%s_front.tex", cardID))

	// Use filepath.Glob to find the matching files
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", err // Return an error if there is an issue with the search
	}

	// If no files match the pattern, return an error
	if len(matches) == 0 {
		return "", fmt.Errorf("no file found matching card ID %s in %s", cardID, kastenPath)
	}

	// Return the first matching file path
	return filepath.Dir(matches[0]), nil
}

func InsertFixme(zettelPath string, cardID string, fix string) error {
	if _, err := os.Stat(zettelPath); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", zettelPath)
	}
	fileName := fmt.Sprintf("fix_%s", cardID)
	filePath := filepath.Join(zettelPath, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(fix)
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}
	return nil
}

// findFlashcards is a helper function to find flashcards in a given Zettel path
func findFlashcards(zettelPath string) ([]Flashcard, error) {
	// Slice to store flashcards
	var flashcards []Flashcard

	// Regular expressions for card front and back file names (ID can be alphanumeric)
	frontPattern := regexp.MustCompile(`card_([a-zA-Z0-9]+)_front\.tex`)
	backPattern := regexp.MustCompile(`card_([a-zA-Z0-9]+)_back\.tex`)

	// Maps to hold matched files (keyed by card ID)
	frontFiles := make(map[string]string)
	backFiles := make(map[string]string)

	// List all files in the zettelPath directory
	files, err := os.ReadDir(zettelPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read directory: %v", err)
	}

	// Iterate over files in the directory and match patterns
	for _, file := range files {
		fileName := file.Name()

		// Match the front card pattern
		if matches := frontPattern.FindStringSubmatch(fileName); matches != nil {
			cardID := matches[1]
			frontFiles[cardID] = filepath.Join(zettelPath, fileName)
		}

		// Match the back card pattern
		if matches := backPattern.FindStringSubmatch(fileName); matches != nil {
			cardID := matches[1]
			backFiles[cardID] = filepath.Join(zettelPath, fileName)
		}
	}

	// Now check if both front and back files exist for every ID
	// and the card has no fixme file.
	for id, frontPath := range frontFiles {
		// Check if card has a fixme file
		fixmePath := fmt.Sprintf("%s/fix_%s", zettelPath, id)
		_, err := os.Stat(fixmePath)
		if !os.IsNotExist(err) {
			log.Println("Card must be fixed. Not syncing to anki")
			continue
		}

		// Check that both front and back files exist
		backPath, exists := backFiles[id]
		if !exists {
			log.Printf("Missing back file for card ID %s. Skipping card.\n", id)
			continue
		}

		// Read content from the front and back files
		frontContent, err := os.ReadFile(frontPath)
		if err != nil {
			log.Printf("Error reading front file for card ID %s: %v. Skipping card.\n", id, err)
			continue
		}

		backContent, err := os.ReadFile(backPath)
		if err != nil {
			log.Printf("Error reading back file for card ID %s: %v. Skipping card.\n", id, err)
			continue
		}

		hash := hex.EncodeToString(Hash(frontContent, backContent))

		// Create a new Flashcard instance and append it to the flashcards slice
		flashcards = append(flashcards, Flashcard{
			ID:    id,
			Front: string(frontPath),
			Back:  string(backPath),
			Hash:  string(hash),
		})
	}

	// Return the list of flashcards
	return flashcards, nil
}
