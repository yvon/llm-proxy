package patcher

import (
	"fmt"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"os"
	"regexp"
	"strings"
)

type Payload struct {
	Model    string         `json:"model,omitempty"`
	Messages []Message      `json:"messages"`
	Unkown   jsontext.Value `json:",unknown"`
}

type Message struct {
	Role         string         `json:"role"`
	Content      string         `json:"content"`
	CacheControl CacheControl   `json:"cache_control,omitzero"`
	Unkown       jsontext.Value `json:",unknown"`
}

type CacheControl struct {
	Type string `json:"type"`
}

func findTag(payload *Payload, tag string) []string {
	var last []string = nil

	pattern := fmt.Sprintf(`\|%s(?::\s*([^|]*))?\|`, regexp.QuoteMeta(tag))
	regexp := regexp.MustCompile(pattern)

	for i := range payload.Messages {
		content := payload.Messages[i].Content
		submatch := regexp.FindStringSubmatch(content)

		if submatch != nil {
			last = submatch
		}

		clean := regexp.ReplaceAllString(content, "")
		payload.Messages[i].Content = strings.TrimSpace(clean)
	}

	return last
}

func getLastUserMessage(messages []Message) *Message {
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == "user" {
			return &messages[i]
		}
	}
	return nil
}

func addCache(payload *Payload) {
	var match = findTag(payload, "cache")

	if match != nil {
		lastUserMessage := getLastUserMessage(payload.Messages)

		if lastUserMessage != nil {
			lastUserMessage.CacheControl = CacheControl{
				Type: "ephemeral",
			}
		}
	}
}

func injectPrefill(payload *Payload) {
	var match = findTag(payload, "prefill")

	if len(match) > 1 {
		newMessage := Message{
			Role:    "assistant",
			Content: match[1],
		}
		payload.Messages = append(payload.Messages, newMessage)
	}
}

func Body(body []byte) []byte {
	var payload Payload

	err := json.Unmarshal(body, &payload)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing JSON:", err)
		return body
	}

	injectPrefill(&payload)
	addCache(&payload)

	modified, err := json.Marshal(payload)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error marshaling JSON:", err)
		return body
	}

	return modified
}
