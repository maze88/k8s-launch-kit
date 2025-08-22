package llm

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

func Init() {
	// Initialize LLM
	llm, err := openai.New(
		openai.WithAPIType(openai.APITypeAzure),
		openai.WithToken("eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjMxMGZlNjA0LWY2YmUtNDEyYy05ZWE4LWZlZjI3ZmQ0NzRlMCIsInNlY3JldCI6IlUwWkZyZ3k0dis1bGlJQWx2VWZweXBxM1NmYmZPb3lmSzVlNGY4b2pMUEU9In0.n4H3Wbl8H15TGlTEd9jil5J1mFxjRRCMXM3JnXg3rc8"),
		openai.WithBaseURL("https://llm-proxy.perflab.nvidia.com"),
		openai.WithModel("model-router"),
		openai.WithEmbeddingModel("text-embedding-3-small"),
	)
	if err != nil {
		log.Fatal(err)
	}

	data, err := os.ReadFile("prompt")
	if err != nil {
		log.Fatal(err)
	}
	prompt := string(data)

	response, err := llms.GenerateFromSinglePrompt(context.Background(), llm, prompt, llms.WithTemperature(0.7))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response)
}
