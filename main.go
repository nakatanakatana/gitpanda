package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

var port int

func init() {
	flag.IntVar(&port, "port", 8000, "HTTP port")
	flag.Parse()
}

func main() {
	fmt.Printf("gitpanda started: port=%d\n", port)
	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text")

	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("It works"))

	case http.MethodPost:
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()

		s := NewSlackWebhook(
			os.Getenv("SLACK_TOKEN"),
			&GitLabURLParserParams{
				APIEndpoint:  os.Getenv("GITLAB_API_ENDPOINT"),
				BaseURL:      os.Getenv("GITLAB_BASE_URL"),
				PrivateToken: os.Getenv("GITLAB_PRIVATE_TOKEN"),
			},
		)
		response, err := s.Request(
			body,
			true,
		)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write([]byte(response))
	}
}
