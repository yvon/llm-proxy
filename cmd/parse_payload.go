package main

import (
	"os"
	"fmt"
	"io"
	"llm_proxy/parsers"
)

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println("Failed to read STDIN:", err)
		return
	}

	fmt.Println(string(parsers.ParsePayload(data)))
}
