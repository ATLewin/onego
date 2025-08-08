# This is the official library to use for accessing OneLLM with golang / go 

## Usage:

1. Run the command: 
```zsh
go get github.com/OneLLM-Dev/onego
```

2. [Get a OneLLM API Key](https://onellm.dev) 

3. Start using OneLLM!
```go
package main

import (
	"fmt"
	"log"
	"os"
	"github.com/OneLLM-Dev/onego"
	"github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()
    if err != nil {
    	log.Fatalf("Error loading .env file: %v", err)
    }
	apikey := os.Getenv("API_KEY")

	message := onego.Message {Role: "user", Content: "Hi there!"}
	messages := []onego.Message{message}
	api := onego.NewAPIInput("https://api.openai.com/v1/chat/completions", onego.ModelGpt4_1, messages, 200)
	output, _ := api.Send(apikey)
	fmt.Println(output.Output.Content)
}
```

