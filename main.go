package main

import (
	"bytes"
	"fmt"
	"io"
	"llm_proxy/parsers"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

func main() {
	target, _ := url.Parse("https://openrouter.ai")

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.Host = target.Host
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/chat/completions" {
			// Payload
			body, _ := io.ReadAll(r.Body)
			modifiedPayload := parsers.ParsePayload(body)
			fmt.Printf("=== PAYLOAD ===\n%s\n\n", string(modifiedPayload))

			r.Body = io.NopCloser(bytes.NewReader(modifiedPayload))
			r.ContentLength = int64(len(modifiedPayload))
			r.Header.Set("Content-Length", strconv.Itoa(len(modifiedPayload)))

			// Response wrapper
			proxy.ModifyResponse = func(resp *http.Response) error {
				resp.Body = &LoggingReader{resp.Body}
				return nil
			}
		}

		proxy.ServeHTTP(w, r)
	})

	log.Println("Reverse proxy running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type LoggingReader struct {
	io.ReadCloser
}

func (lr *LoggingReader) Read(p []byte) (n int, err error) {
	n, err = lr.ReadCloser.Read(p)
	if n > 0 {
		fmt.Printf("=== CHUNK ===\n%s\n", string(p[:n]))
	}
	return n, err
}
