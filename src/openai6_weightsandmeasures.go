// This program tries to scale and convert weights and measures from one
// unit to another with correct multiplication factors. It uses a predefied
// set of conversion factors defined as CSV data.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Replace "OPENAI_API_KEY" with your OpenAI API key
	// apiKey := "OPENAI_API_KEY"
	// OR define as an environment variable OPENAI_API_KEY
	apiKey := os.Getenv("OPENAI_API_KEY")

	// Items to convert
	scaleFactor := 1.0
	quantity := 2.0
	fromUnit := "kilogram"
	toUnit := "pound"
	// Create input prompt for OpenAI
	prompt := generatePrompt(generateQuestion(scaleFactor, quantity, fromUnit, toUnit))

	data := map[string]interface{}{
		"model":       "text-davinci-003",
		"prompt":      prompt,
		"temperature": 0.5,
		"max_tokens":  50,
		"n":           1,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error Marshalling JSON data:", err)
		return
	}

	fmt.Println("Request sent was: %s", string(jsonData))

	body := strings.NewReader(string(jsonData))

	// Set up the HTTP request with the API key in the headers
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", body)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send the HTTP request and read the response
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading HTTP response:", err)
		return
	}

	// Pretty print the JSON response
	var jsonResponse interface{}
	err = json.Unmarshal(respBody, &jsonResponse)
	if err != nil {
		fmt.Println("Error unmarshaling JSON response:", err)
		return
	}
	prettyJson, err := json.MarshalIndent(jsonResponse, "", "    ")
	if err != nil {
		fmt.Println("Error pretty printing JSON response:", err)
		return
	}

	// Print the generated text
	fmt.Println(string(prettyJson))
}

// The CSV data table given below used for training comes from:
// https://en.wikipedia.org/wiki/Imperial_units
func generatePrompt(question string) string {
	return `Read the CSV data given below and answer the question: ` + question + `.

Unit,Imperial ounces,Imperial pints,Millilitres,Cubic inches,US ounces,US pints
fluid ounce (fl oz),1 ,1/20 ,28.4130625 ,1.7339 ,0.96076 ,0.060047
gill (gi),5 ,1/4 ,142.0653125 ,8.6694 ,4.8038 ,0.30024
pint (pt),20 ,1 ,568.26125 ,34.677 ,19.215 ,1.2009
quart (qt),40 ,2 ,1136.5225 ,69.355 ,38.43 ,2.4019
gallon (gal),160 ,8 ,4546.09 ,277.42 ,153.72 ,9.6076

`
}

func generateQuestion(scaleFactor float64, quantity float64, fromUnit string, toUnit string) string {
	return fmt.Sprintf("Scale multiply with a scale factor of %.2f and convert %.2f %s to %s", scaleFactor, quantity, fromUnit, toUnit)
}
