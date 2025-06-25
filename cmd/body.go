package main

import (
	"io"
	"os"
	"fmt"
	"llmproxy/patcher"
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
