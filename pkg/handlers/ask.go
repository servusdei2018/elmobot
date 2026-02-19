package handlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/servusdei2018/elmobot/pkg/nim"
)

const (
	defaultModel       = "meta/llama-4-maverick-17b-128e-instruct"
	defaultTemperature = 1.0
	maxMessageLength   = 2000
	updateInterval     = 500 * time.Millisecond
	streamTimeout      = 5 * time.Minute
)

func Ask(s *discordgo.Session, i *discordgo.InteractionCreate) {
	apiKey := os.Getenv("NVIDIA_API_KEY")
	if apiKey == "" {
		respondWithError(s, i, "NVIDIA_API_KEY environment variable not set")
		return
	}

	client, err := nim.NewClient(nim.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		respondWithError(s, i, fmt.Sprintf("Failed to initialize NIM client: %v", err))
		return
	}

	var question string
	if len(i.ApplicationCommandData().Options) > 0 {
		question = i.ApplicationCommandData().Options[0].StringValue()
	}
	if question == "" {
		respondWithError(s, i, "Question cannot be empty")
		return
	}

	req := &nim.CompletionRequest{
		Model: defaultModel,
		Messages: []nim.Message{
			{
				Role:    "user",
				Content: question,
			},
		},
		MaxTokens:   2048,
		Temperature: defaultTemperature,
		Stream:      true,
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "🤔 Thinking...",
		},
	})
	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), streamTimeout)
	defer cancel()

	eventChan, errChan, err := client.CreateCompletionStream(ctx, req)
	if err != nil {
		respondWithErrorEdit(s, i, fmt.Sprintf("Failed to start stream: %v", err))
		return
	}

	fullResponse := strings.Builder{}
	lastUpdate := time.Now()
	ticker := time.NewTicker(updateInterval)
	defer ticker.Stop()

	for {
		select {
		case event, ok := <-eventChan:
			if !ok {
				goto streamDone
			}
			if len(event.Choices) > 0 && event.Choices[0].Delta.Content != "" {
				fullResponse.WriteString(event.Choices[0].Delta.Content)
			}

		case err := <-errChan:
			if err != nil && err != io.EOF {
				respondWithErrorEdit(s, i, fmt.Sprintf("Stream error: %v", err))
				return
			}
			goto streamDone

		case <-ticker.C:
			if time.Since(lastUpdate) >= updateInterval && fullResponse.Len() > 0 {
				content := fullResponse.String()
				if len(content) > maxMessageLength {
					content = content[:maxMessageLength-3] + "..."
				}
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &content,
				})
				lastUpdate = time.Now()
			}

		case <-ctx.Done():
			respondWithErrorEdit(s, i, "Request timeout")
			return
		}
	}

streamDone:
	content := fullResponse.String()
	if content == "" {
		content = "(no response)"
	}

	if len(content) > maxMessageLength {
		content = content[:maxMessageLength-3] + "..."
	}

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
}

func respondWithError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("❌ Error: %s", message),
		},
	})
}

func respondWithErrorEdit(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	content := fmt.Sprintf("❌ Error: %s", message)
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
}
