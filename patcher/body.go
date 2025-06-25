package patcher

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"regexp"
)

var prefillRegex = regexp.MustCompile(`\|prefill:([^|]*)\|`)

func processMessage(message map[string]any) string {
	content, ok := message["content"].(string)
	if !ok {
		fmt.Fprintln(os.Stderr, "Message content is not a string")
		return ""
	}

	// Extract the prefill
	matches := prefillRegex.FindStringSubmatch(content)
	if len(matches) < 2 {
		return ""
	}

	prefill := strings.TrimSpace(matches[1])

	// Clean the message
	cleanContent := prefillRegex.ReplaceAllString(content, "")
	message["content"] = strings.TrimSpace(cleanContent)

	return prefill
}

func injectPrefill(payload map[string]any) {
	if messages, ok := payload["messages"].([]any); ok {
		var prefill string
		for i := range messages {
			if message, ok := messages[i].(map[string]any); ok {
				if newPrefill := processMessage(message); newPrefill != "" {
					prefill = newPrefill
				}
			} else {
				fmt.Fprintln(os.Stderr, "Error: unexpected message type")
			}
		}
		if prefill != "" {
			newMessage := map[string]any{
				"role":    "assistant",
				"content": prefill,
			}
			payload["messages"] = append(messages, newMessage)
		}
	} else {
		fmt.Fprintln(os.Stderr, "Error: expected an array of messages")
	}
}

func Body(data []byte) []byte {
	var payload map[string]any
	err := json.Unmarshal(data, &payload)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing JSON:", err)
		return data
	}
	injectPrefill(payload)
	modifiedData, err := json.Marshal(payload)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error marshaling JSON:", err)
		return data
	}
	return modifiedData
}
