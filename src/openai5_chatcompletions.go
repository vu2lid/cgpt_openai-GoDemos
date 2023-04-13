// This program prints out the next chat completion tokens generated
// based on the context for a specified set of messages - gpt-3.5-turbo - model.
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

	// Set up the request body as a JSON string
	body := strings.NewReader(`{
        "model": "gpt-3.5-turbo",
        "messages": [
            {"role": "system", "content": "You are a helpful assistant."},
            {"role": "user", "content": "These are the voyages of the?"},
            {"role": "assistant", "content": "These are the voyages of the starship Enterprise"},
            {"role": "user", "content": "What comes next?"}
        ]
    }`)

	fmt.Printf("Request sent was: %s\n\n", body)

	// Set up the HTTP request with the API key in the headers
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", body)
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
