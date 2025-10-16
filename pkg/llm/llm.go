package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/nvidia/k8s-launch-kit/pkg/config"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func SelectPrompt(promptPath string, config config.ClusterConfig, llmApiKey string, llmApiUrl string, llmVendor string) (map[string]string, error) {
	options := []openai.Option{}

	if llmVendor == "openai-azure" {
		options = append(options, openai.WithAPIType(openai.APITypeAzure))
		options = append(options, openai.WithToken(llmApiKey))
		options = append(options, openai.WithBaseURL(llmApiUrl))
		options = append(options, openai.WithModel("model-router"))
		options = append(options, openai.WithEmbeddingModel("text-embedding-3-small"))
		options = append(options, openai.WithAPIVersion("2025-02-01-preview"))
	} else {
		return nil, fmt.Errorf("unsupported LLM vendor: %s", llmVendor)
	}

	llm, err := openai.New(options...)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile("system-prompt")
	if err != nil {
		return nil, err
	}

	prompt := string(data)

	configJson, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	prompt = fmt.Sprintf("%s\n%s\nUSER:", prompt, string(configJson))

	data, err = os.ReadFile(promptPath)
	if err != nil {
		return nil, err
	}
	prompt = fmt.Sprintf("%s\n%s", prompt, string(data))

	log.Log.V(1).Info("User prompt", "prompt", string(data))

	response, err := llms.GenerateFromSinglePrompt(context.Background(), llm, prompt, llms.WithTemperature(0.5))
	if err != nil {
		return nil, err
	}

	log.Log.V(1).Info("LLM Response", "response", response)

	jsonResponse := make(map[string]string)
	err = json.Unmarshal([]byte(response), &jsonResponse)
	if err != nil {
		return nil, err
	}

	return jsonResponse, nil
}
