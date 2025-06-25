package main

import (
	"fmt"
	"io"
	"llmproxy/patcher"
	"os"
)

func main() {
	body, err := io.ReadAll(os.Stdin)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading body:", err)
		return
	}

	patched := patcher.Body(body)
	fmt.Println(string(patched))
}
