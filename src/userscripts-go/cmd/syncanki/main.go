package main

import (
	"fmt"
	"log"
	"os"
	"xk/src/userscripts-go/pkg/api"
)

// the Anki-Connect API
var url, _ = os.LookupEnv("ANKI_CONNECT_URL")
var deck, _ = os.LookupEnv("ANKI_DECK_NAME")
var modelName, _ = os.LookupEnv("ANKI_MODEL_NAME")
var connect = AnkiConnect{Url: url}

// Flashcard structure, representing front, back, id, hash
type Flashcard struct {
	Front string
	Back  string
	ID    string
	Hash  string
}

func Tex2Anki(flashcard Flashcard) (string, string, error) {
	frontFilenameAnki := fmt.Sprintf("%s_front.svg", flashcard.ID)
	backFilenameAnki := fmt.Sprintf("%s_back.svg", flashcard.ID)

	frontHtml := fmt.Sprintf("<img src=%s>", frontFilenameAnki)
	backHtml := fmt.Sprintf("<img src=%s>", backFilenameAnki)

	// Compile LaTex
	frontSVG, err := Tex2Base64(flashcard.Front)
	if err != nil {
		return "", "", err
	}
	backSVG, err := Tex2Base64(flashcard.Back)
	if err != nil {
		return "", "", err
	}

	// Store SVG in Anki
	err = StoreMedia(&connect, frontSVG, frontFilenameAnki)
	if err != nil {
		return "", "", err
	}
	err = StoreMedia(&connect, backSVG, backFilenameAnki)
	if err != nil {
		return "", "", err
	}

	return frontHtml, backHtml, nil
}

// processZettel retrieves the Zettel path and handles the retrieval and comparison of multiple flashcards
func processZettel(zettel string) error {
	log.Printf("Processing zettel %s", zettel)

	// Retrieve the Zettel's path using api.Xk("path")
	zettelPaths, err := api.Xk("path", map[string]string{"z": zettel})
	if err != nil {
		log.Fatalf("Unable to retrieve path for zettel '%s': %v", zettel, err)
	}

	// Assuming api.Xk returns a list of strings and we take the first element
	if len(zettelPaths) == 0 {
		log.Fatalf("No path found for zettel '%s'", zettel)
	}
	zettelPath := zettelPaths[0]

	// Continue with processing flashcards in the zettel path
	flashcards, err := findFlashcards(zettelPath)
	if err != nil {
		return err
	}

	for _, flashcard := range flashcards {
		// Search for the card in Anki by its ID
		log.Printf("---%s---", flashcard.ID)
		ankiID, err := FindCard(&connect, deck, flashcard.ID)

		if ankiID != -1 || err != nil {
			log.Printf("Found Flashard %s in anki: %d", flashcard.ID, ankiID)
			noteID, err := Card2Note(&connect, ankiID)
			if err != nil {
				log.Printf("Failed to get note of card %d", ankiID)
				continue
			}
			cardHash, err := GetCardField(&connect, noteID, "hash")
			if err != nil {
				log.Printf("Failed to get hash from note %d. Updating card.", noteID)
				cardHash = ""
			}

			log.Printf("Hash in Anki: %s, current hash: %s", cardHash, flashcard.Hash)
			if flashcard.Hash == cardHash {
				log.Printf("Flashcard %s unchanged, skipping update", flashcard.ID)
				continue
			}
			front, back, err := Tex2Anki(flashcard)
			if err != nil {
				log.Println(err)
				continue
			}
			updateCard(ankiID, flashcard, front, back)
			continue
		}

		log.Println("Flashcard does not exist. Creating it.")
		front, back, err := Tex2Anki(flashcard)
		if err != nil {
			log.Println(err)
			continue
		}
		addNewCard(flashcard, front, back)
	}
	return nil
}

// Main function
func main() {
	// find cards to fix
	pathResp, err := api.Xk("path", map[string]string{})
	if err != nil {
		log.Println("Unable to retrieve zettel kasten path.")
		os.Exit(1)
	}

	// get the note ids of flashcards to fix
	cardsToFix, err := FindFixme()
	if err != nil {
		log.Println("Unable to retrieve flashcards to fix.")
		os.Exit(1)
	}

	for _, card := range cardsToFix {
		// get the flashcard id
		defer RemoveCard(&connect, card)

		cardID, err := GetCardField(&connect, card, "id")
		if err != nil {
			log.Println("Unable to get cards id. Skipping")
			continue
		}
		cardIDstring, ok := cardID.(string)
		if !ok {
			log.Println("Invalid ID type. Skipping")
			continue
		}

		fixme, err := GetCardField(&connect, card, "fixme")
		if err != nil {
			log.Println("Unable to get cards fixme. Skipping")
			continue
		}
		fixmeString, ok := fixme.(string)
		if !ok {
			log.Println("fixme of invalid type. Skipping")
			continue
		}

		// Find the Zettel the card originated from
		originZettel, err := Card2Zettel(pathResp[0], cardIDstring)
		if err != nil {
			log.Println("Unable to find origin zettel. Skipping")
			continue
		}

		InsertFixme(originZettel, cardIDstring, fixmeString)
	}

	zettels, err := api.Xk("ls", map[string]string{})
	if err != nil {
		log.Println("Unable to retrieve zettels.")
		os.Exit(1)
	}

	if len(zettels) == 0 {
		log.Println("No zettels found.")
		os.Exit(1)
	}

	// if the deck does not exist, create it
	decks, err := GetDecks(&connect)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(decks)

	// If deck does not exist, create it
	if !deckExists(decks, deck) {
		_, err := CreateDeck(&connect, deck)
		if err != nil {
			log.Fatalf("Failed to create deck: %v", err)
		}
	}

	// If model does not exist, create it
	fields := []string{"front", "back", "id", "hash", "fixme"}
	template := []map[string]string{
		{
			"Front": "{{front}}",
			"Back":  "{{back}}",
		},
	}

	CreateModel(
		&connect,
		modelName,
		fields,
		template,
	)

	// Process each zettel
	for _, z := range zettels {
		// fmt.Println(z)
		processZettel(z)
	}
}
