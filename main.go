package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type ChatResponse struct {
	Model     string  `json:"model"`
	Message   Message `json:"message"`
	CreatedAt string  `json:"created_at"`
	Done      bool    `json:"done"`
}

type FunctionDefinition struct {
	// Type      string `json:"type"`
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
	FollowUp  string `json:"follow_up,omitempty"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Format   string    `json:"format"`
	Stream   bool      `json:"stream"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

var functionDefinitions = []FunctionDefinition{
	{
		// Type:      "function",
		Name:      "get_weather",
		Arguments: "San Francisco",
	},
	{

		// Type:      "function",
		Name:      "get_weather",
		Arguments: "Toronto",
	},
	{
		// Type:      "function",
		Name:      "get_weather",
		Arguments: "",
		FollowUp:  "In which city you are located?",
	},
}

var second_messages = []Message{
	{
		Role:    "system",
		Content: "You are a classification bot, and your task is to classify user input and fit into one of two functions: get_weather and not_applicable. You will have to identify the arguments can be passed into the function which is a location. If no location is provided, ask user to specify location. Respond using JSON",
	},
	{
		Role:    "user",
		Content: "What's the weather like in San Francisco?"},
	{
		Role:    "assistant",
		Content: "",
	},
	{
		Role:    "user",
		Content: "I'm visiting my parents in Toronto"},
	{
		Role:    "assistant",
		Content: "",
	},
	{
		Role:    "user",
		Content: "Is it going to be sunny next week?"},
	{
		Role:    "assistant",
		Content: "",
	},
	{
		Role:    "user",
		Content: "is it going to be raining next week?",
	},
}

var messages = []Message{
	{
		Role:    "system",
		Content: "You are a classification bot, and your task is to classify user input and fit into one of two functions: get_weather and not_applicable. You will have to identify the arguments can be passed into the function which is a location. If no location is provided, ask user to specify location. Respond using JSON",
	},
	{
		Role:    "user",
		Content: "What's the weather like in San Francisco?"},
	{
		Role:    "assistant",
		Content: "",
	},
	{
		Role:    "user",
		Content: "I'm visiting my parents in Toronto"},
	{
		Role:    "assistant",
		Content: "",
	},
	{
		Role:    "user",
		Content: "Is it going to be sunny next week?"},
	{
		Role:    "assistant",
		Content: "",
	},
	{
		Role:    "user",
		Content: "is it going to be raining in Montreal?",
	},
}

func main() {
	url := "http://localhost:11434/api/chat"
	method := "POST"

	chatRequest := ChatRequest{
		Model:    "mistral",
		Format:   "json",
		Stream:   false,
		Messages: messages,
	}

	second_chatRequest := ChatRequest{
		Model:    "mistral",
		Format:   "json",
		Stream:   false,
		Messages: second_messages,
	}

	for _, f := range functionDefinitions {
		jsons, err := json.Marshal(f)
		if err != nil {
			fmt.Println("Error line 145:", err)
			return
		}

		for i, message := range chatRequest.Messages {
			if message.Content == "" {
				chatRequest.Messages[i].Content = string(jsons)
				break
			}

		}
	}

	for _, f := range functionDefinitions {
		jsons, err := json.Marshal(f)
		if err != nil {
			fmt.Println("Error line 145:", err)
			return
		}

		for i, message := range second_chatRequest.Messages {
			if message.Content == "" {
				second_chatRequest.Messages[i].Content = string(jsons)
				break
			}

		}
	}

	payload2, err := json.Marshal(chatRequest)
	if err != nil {
		fmt.Errorf("Error at line 153 %s", err)
	}

	payload, err := json.Marshal(second_chatRequest)
	if err != nil {
		fmt.Errorf("Error at line 153 %s", err)
	}

	client := &http.Client{}
	req2, err := http.NewRequest(method, url, bytes.NewBuffer(payload2))
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	req2.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Type", "application/json")

	requests := []*http.Request{req, req2}

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, req := range requests {

		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		var payload ChatResponse
		if err := json.Unmarshal(body, &payload); err != nil {
			fmt.Println("Error line 166")
			panic(err)
		}

		// Unmarshal the content of Payload.Message.Content into FunctionDefinition struct
		var function FunctionDefinition
		if err := json.Unmarshal([]byte(payload.Message.Content), &function); err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Print the function data
		fmt.Println("Function:", function.Name)
		fmt.Println("Arguments:", function.Arguments)
		fmt.Println("FollowUp:", function.FollowUp)

		if function.Name == "get_weather" && function.Arguments != "" {
			slog.Info("Got value", "temperature", get_weather())
		}
	}

}

// get_weather() is a function return an int indicating current temperature in celsius
func get_weather() int {
	slog.Info("Invoking get_weather()...")
	return 20
}
