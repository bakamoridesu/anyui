package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var buffers = make(map[int]*Buffer)

type Buffer struct {
	buf      *bytes.Buffer
	finalLen int
}

func registerBuffer(id int) {
	var newBuf bytes.Buffer
	buffers[id] = &Buffer{&newBuf, 0}
}

func unregisterBuffer(id int) {
	delete(buffers, id)
}

func writeBufToCh(ctx context.Context, ch chan<- string, bufId int) {
	var prevLen int
	buffer := buffers[bufId]
	buf := buffer.buf
	ticker := time.NewTicker(time.Second)

outerloop:
	for {
		select {
		case <-ticker.C:
			if buf.Len() > prevLen {
				str := buf.String()
				if prevLen != 0 {
					str = str[prevLen:]
				}
				itemarr := strings.Split(str, "\n")
				for _, it := range itemarr {
					if it != "" {
						ch <- fmt.Sprintf("data:%s\n\n", it)
					}
				}
				prevLen = buf.Len()
			}
			if buffers[bufId].finalLen != 0 && buffers[bufId].finalLen == buf.Len() {
				time.Sleep(1 * time.Second)
				ch <- fmt.Sprintf("event:done\ndata:%s\n\n", "finished")
				prevLen = 0
				buffer.finalLen = 0
				buf.Reset()
			}
		case <-ctx.Done():
			break outerloop
		}
	}

	ticker.Stop()
	close(ch)
}

func Handler(cfg *config) http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			id, err := strconv.Atoi(r.FormValue("id"))
			if err != nil {
				http.Error(w, "incorrect id", http.StatusInternalServerError)
				return
			}
			buffer, ok := buffers[id]
			if !ok {
				http.Error(w, "fatal error", http.StatusInternalServerError)
				return
			}
			buf := buffer.buf
			buf.Reset()
			buffer.finalLen = 0
			cmd, err := command(*cfg, r)
			if err != nil {
				fmt.Println("Cmd creation failed:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			cmd.Stderr = buf
			cmd.Stdout = buf
			go func() {
				err = cmd.Run()
				if err != nil {
					fmt.Println("Running command failed", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				buffer.finalLen = buf.Len()
			}()
			return
		}
		http.ServeFile(w, r, "template.html")
	})

	m.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		id, err := strconv.Atoi(q.Get("id"))
		if err != nil {
			http.Error(w, "incorrect id", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		registerBuffer(id)
		outCh := make(chan string)
		go writeBufToCh(r.Context(), outCh, id)

		for msg := range outCh {
			fmt.Fprint(w, msg)
			w.(http.Flusher).Flush()
		}

		unregisterBuffer(id)
	})
	return m
}
