package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var buf bytes.Buffer
var finalLen int
var prevBufLen int

func Handler(cfg *config) http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method)
		if r.Method == http.MethodPost {
			cmd, err := command(*cfg, r)
			buf.Reset()
			finalLen = 0
			cmd.Stderr = &buf
			cmd.Stdout = &buf
			if err != nil {
				fmt.Println("Cmd creation failed:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			go func() {
				err = cmd.Run()
				if err != nil {
					fmt.Println("Running command failed", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				finalLen = buf.Len()
			}()
			return
		}
		http.ServeFile(w, r, "template.html")
	})

	m.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		for {
			if buf.Len() > prevBufLen {
				str := buf.String()
				if prevBufLen != 0 {
					str = str[prevBufLen:]
				}
				itemarr := strings.Split(str, "\n")
				for _, it := range itemarr {
					if it != "" {
						fmt.Fprintf(w, "data: %s\n\n", it)
					}
				}
				w.(http.Flusher).Flush()
				prevBufLen = buf.Len()
			}

			if finalLen != 0 && finalLen == buf.Len() {
				fmt.Fprintf(w, "data:%s\n\n", "finished")
				w.(http.Flusher).Flush()
				prevBufLen = 0
				finalLen = 0
				buf.Reset()
			}
			time.Sleep(1 * time.Second)
		}
	})
	return m
}
