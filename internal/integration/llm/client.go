package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

// LLMClient - наш клиент к API
type LLMClient struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

func NewClient(baseURL, apiKey string) *LLMClient {
	return &LLMClient{
		httpClient: &http.Client{Timeout: 90 * time.Second},
		apiKey:     apiKey,
		baseURL:    baseURL,
	}
}

// --- 1. Chat Completion (gpt-4o-mini) ---

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream,omitempty"`
}
type ChatResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

func (c *LLMClient) GetChatCompletion(ctx context.Context, messages []ChatMessage) (string, error) {
	reqBody := ChatRequest{
		Model:    "groq/compound",
		Messages: messages,
	}

	respBody, err := c.doRequest(ctx, "/chat/completions", reqBody)
	if err != nil {
		return "", err
	}

	var resp ChatResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return "", fmt.Errorf("failed to parse LLM chat response: %w", err)
	}

	if len(resp.Choices) == 0 || resp.Choices[0].Message.Content == "" {
		return "", errors.New("no choice returned from LLM")
	}
	return resp.Choices[0].Message.Content, nil
}

// --- 2. Embeddings (text-embedding-3-small) ---

type EmbeddingRequest struct {
	Input string `json:"input"`
	Model string `json:"model"`
}
type EmbeddingData struct {
	Embedding []float32 `json:"embedding"`
}
type EmbeddingResponse struct {
	Data []EmbeddingData `json:"data"`
}

func (c *LLMClient) GetEmbedding(ctx context.Context, text string) ([]float32, error) {
	reqBody := EmbeddingRequest{Input: text, Model: "text-embedding-3-small"}
	respBody, err := c.doRequest(ctx, "/embeddings", reqBody)
	if err != nil {
		return nil, err
	}

	var resp EmbeddingResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse LLM embedding response: %w", err)
	}
	if len(resp.Data) == 0 {
		return nil, errors.New("no embedding data returned")
	}
	return resp.Data[0].Embedding, nil
}

// --- 3. Moderations ---

type ModerationRequest struct {
	Input string `json:"input"`
}
type ModerationResponse struct {
	Results []struct {
		Flagged bool `json:"flagged"`
	} `json:"results"`
}

func (c *LLMClient) CheckModeration(ctx context.Context, text string) (bool, error) {
	reqBody := ModerationRequest{Input: text}
	respBody, err := c.doRequest(ctx, "/moderations", reqBody)
	if err != nil {
		return false, err
	}

	var resp ModerationResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return false, fmt.Errorf("failed to parse moderation response: %w", err)
	}
	if len(resp.Results) == 0 {
		return false, errors.New("no moderation result")
	}
	return resp.Results[0].Flagged, nil
}

// --- 4. Speech-to-Text (whisper-1) ---

type WhisperResponse struct {
	Text string `json:"text"`
}

func (c *LLMClient) GetTranscription(ctx context.Context, audioData []byte, filename string) (string, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Добавляем файл
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}
	part.Write(audioData)

	// Добавляем модель
	writer.WriteField("model", "whisper-large-v3")
	writer.Close()

	url := c.baseURL + "/audio/transcriptions"
	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("whisper API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var whisperResp WhisperResponse
	if err := json.Unmarshal(respBody, &whisperResp); err != nil {
		return "", err
	}
	return whisperResp.Text, nil
}

// --- doRequest (Helper) ---
// ЭТОТ МЕТОД НЕ БЫЛ НАЙДЕН, ПОТОМУ ЧТО ПРЕДЫДУЩИЙ ФАЙЛ БЫЛ СЛОМАН
func (c *LLMClient) doRequest(ctx context.Context, path string, payload interface{}) ([]byte, error) {
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LLM API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}
	return respBody, nil
}
