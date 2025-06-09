package parsers

import (
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
)

func processMessage(message map[string]any) string {
	content, ok := message["content"].(string)

	if !ok {
		println("Message content is not a string")
		return ""
	}

	t, err := template.New("content").Delims("[[", "]]").Parse(content)
	if err != nil {
		println("Error parsing template:", err)
		return ""
	}

	var contentWriter strings.Builder
	err = t.Execute(&contentWriter, nil)
	if err != nil {
		println("Error executing main template:", err)
		return ""
	}

	message["content"] = contentWriter.String()

	var prefillWriter strings.Builder
	err = t.ExecuteTemplate(&prefillWriter, "prefill", nil)
	if err != nil {
		println("Error executing prefill template:", err)
		return ""
	}

	return prefillWriter.String()
}

func ParsePayload(data []byte) []byte {
	var err error

	// Parse the JSON data into a struct
	var payload map[string]any
	err = json.Unmarshal(data, &payload)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return data
	}

	if messages, ok := payload["messages"].([]any); ok {
		var prefill string

		for i := range messages {
			if message, ok := messages[i].(map[string]any); ok {
				prefill = processMessage(message)
			}
		}

		if prefill != "" {
			newMessage := map[string]any{
				"role":    "assistant",
				"content": prefill,
			}

			payload["messages"] = append(messages, newMessage)
		}
	}

	modifiedData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return data
	}

	return modifiedData
}
