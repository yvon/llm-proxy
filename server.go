package main

import (
	"bytes"
	"fmt"
	"io"
	"llmproxy/patcher"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

func serve() {
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
			patched := patcher.Body(body)
			fmt.Printf("=== PAYLOAD ===\n%s\n\n", string(patched))

			r.Body = io.NopCloser(bytes.NewReader(patched))
			r.ContentLength = int64(len(patched))
			r.Header.Set("Content-Length", strconv.Itoa(len(patched)))
		}

		proxy.ServeHTTP(w, r)
	})

	log.Println("Reverse proxy running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
