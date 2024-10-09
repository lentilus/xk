package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// generic type for request body parameters
type P any

// GenericResponse encapsulates a standard API response structure with a result and error
type GenericResponse[T any] struct {
	Result T   `json:"result"`
	Error  any `json:"error"`
}

type Body struct {
	Action  string `json:"action"`
	Version int    `json:"version"`
	Params  P      `json:"params,omitempty"`
}

// API encapsulates the actual API calls to allow mocking during testing
type API interface {
	Request(string, P, any) error
}

// AnkiConnect is the concrete implementation of the API interface
type AnkiConnect struct {
	Url string
}

// Make a request to the AnkiConnect REST-API
func (a *AnkiConnect) Request(action string, params P, response any) error {
	body := Body{action, 6, params}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, a.Url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close() // Ensure the response body is closed

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code: %d", res.StatusCode)
	}

	if err := json.NewDecoder(res.Body).Decode(response); err != nil {
		return err
	}

	return nil
}

// Helper function to check for errors in API response
func checkAPIError(errField any) error {
	if errField != nil {
		msg, ok := errField.(string)
		if ok {
			return errors.New(msg)
		}
		return errors.New("failed to decode error message from response")
	}
	return nil
}

// returns a list of Anki deck names (and an error)
func GetDecks(api API) ([]string, error) {
	var res GenericResponse[[]string]
	if err := api.Request("deckNames", nil, &res); err != nil {
		return nil, err
	}

	if err := checkAPIError(res.Error); err != nil {
		return nil, err
	}

	return res.Result, nil
}

// creates a new deck and returns its ID (and an error)
func CreateDeck(api API, name string) (int, error) {
	var res GenericResponse[int]
	err := api.Request("createDeck", map[string]any{"deck": name}, &res)
	if err != nil {
		return -1, err
	}

	if err := checkAPIError(res.Error); err != nil {
		return -1, err
	}

	return res.Result, nil
}

// stores file with base64 encoded content
func StoreMediaFile(api API, name string, data string) (string, error) {
	var res GenericResponse[string]
	params := map[string]any{"filename": name, "data": data}
	err := api.Request("storeMediaFile", params, &res)
	if err != nil {
		return "", err
	}

	if err := checkAPIError(res.Error); err != nil {
		return "", err
	}

	return res.Result, nil
}

// retrieves available models
func GetModels(api API) ([]string, error) {
	var res GenericResponse[[]string]
	err := api.Request("modelNames", nil, &res)
	if err != nil {
		return nil, err
	}

	if err := checkAPIError(res.Error); err != nil {
		return nil, err
	}

	return res.Result, nil
}

// creates a new model with specified fields and templates
func CreateModel(
	api API,
	name string,
	fields []string,
	templates []map[string]string,
) (any, error) {
	var res GenericResponse[any]
	params := map[string]any{
		"modelName":     name,
		"inOrderFields": fields,
		"isCloze":       false,
		"cardTemplates": templates,
	}

	err := api.Request("createModel", params, &res)
	if err != nil {
		return nil, err
	}

	if err := checkAPIError(res.Error); err != nil {
		return nil, err
	}

	return res.Result, nil
}

// finds a card based on deck name and card id
func FindCard(api API, deck string, id string) (int, error) {
	var res GenericResponse[[]int]

	params := map[string]any{
		"query": fmt.Sprintf("deck:%s id:%s", deck, id),
	}

	err := api.Request("findCards", params, &res)
	if err != nil {
		return -1, err
	}

	if err := checkAPIError(res.Error); err != nil {
		return -1, err
	}

	// We expect exactly one result
	if len(res.Result) == 0 {
		return -1, nil
	}

	if len(res.Result) != 1 {
		return -1, errors.New("We expected 0 or 1 matches.")
	}

	// return Card2Note(api, res.Result[0])
	return res.Result[0], nil
}

// converts a card to its corresponding note
func Card2Note(api API, id int) (int, error) {
	var res GenericResponse[[]int]

	params := map[string]any{
		"cards": []int{id},
	}

	err := api.Request("cardsToNotes", params, &res)
	if err != nil {
		return -1, err
	}

	if err := checkAPIError(res.Error); err != nil {
		return -1, err
	}

	if len(res.Result) != 1 {
		return -1, errors.New("expected exactly one match")
	}

	return res.Result[0], nil
}

// retrieves a specific field value from a card by note ID
func GetCardField(
	api API,
	noteid int,
	field string,
) (any, error) {
	var res GenericResponse[[]map[string]any]

	params := map[string]any{
		"notes": []int{noteid},
	}

	err := api.Request("notesInfo", params, &res)
	if err != nil {
		return nil, err
	}

	if err := checkAPIError(res.Error); err != nil {
		return nil, err
	}

	if len(res.Result) != 1 {
		return nil, errors.New("expected exactly one result")
	}

	fields, ok := res.Result[0]["fields"].(map[string]any)
	if !ok {
		return nil, errors.New("field not present")
	}

	v, ok := fields[field].(map[string]any)
	if !ok {
		return nil, errors.New("field in invalid format")
	}

	return v["value"], nil
}

// UpdateNoteFields updates specific fields of a note identified by its note ID.
func UpdateNoteFields(
	api API,
	noteID int,
	fields map[string]string, // Generalized to update any fields dynamically
) (any, error) {

	// Safety checks
	if noteID <= 0 {
		return nil, errors.New("invalid note ID")
	}
	if len(fields) == 0 {
		return nil, errors.New("no fields provided for update")
	}

	// Prepare the request payload
	params := map[string]any{
		"note": map[string]any{
			"id":     noteID,
			"fields": fields, // Dynamic fields input
		},
	}

	// Response struct using the generic response abstraction
	var res GenericResponse[any]
	err := api.Request("updateNoteFields", params, &res)
	if err != nil {
		return nil, err
	}

	// Check if the API returned an error
	if err := checkAPIError(res.Error); err != nil {
		return nil, err
	}

	// Return the result of the update operation
	return res.Result, nil
}

// AddCard adds a card to the specified deck with dynamic fields.
func AddCard(
	api API,
	deck string,
	model string,
	fields map[string]string, // Generalized to allow any fields
) (any, error) {
	// Safety checks
	if deck == "" {
		return nil, errors.New("deck name cannot be empty")
	}
	if model == "" {
		return nil, errors.New("model name cannot be empty")
	}
	if len(fields) == 0 {
		return nil, errors.New("no fields provided for the card")
	}

	// Prepare the request parameters
	params := map[string]any{
		"note": map[string]any{
			"deckName":  deck,
			"modelName": model,
			"fields":    fields, // Dynamic fields input
		},
	}

	// Use the GenericResponse struct to capture the response
	var res GenericResponse[any]

	// Make the API request
	err := api.Request("addNote", params, &res)
	if err != nil {
		return nil, err
	}

	// Check for any API error
	if err := checkAPIError(res.Error); err != nil {
		return nil, err
	}

	// Return the result of adding the card
	return res.Result, nil
}

func RemoveCard(
	api API,
	id int,
) error {
	params := map[string]any{
		"notes": []int{id},
	}
	var res GenericResponse[any]
	err := api.Request("deleteNotes", params, &res)
	if err != nil {
		return err
	}
	return nil
}

func StoreMedia(
	api API,
	data string,
	filename string,
) error {
	params := map[string]any{"filename": filename, "data": data}
	var res GenericResponse[string]

	err := api.Request("storeMediaFile", params, &res)
	if err != nil {
		return err
	}
	return nil
}
