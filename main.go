package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func main() {
	// load config

	// check if template found

	go func() {
		time.Sleep(2 * time.Second)
		cName, cParams := GetOpen()
		cParams = append(cParams, "http://127.0.0.1:8080")
		cPath, err := exec.LookPath(cName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		exec.Command(cPath, cParams...).Run()
	}()

	s := http.Server{
		Addr: "127.0.0.1:8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				fmt.Println(r.Form)
				fmt.Println(r.FormValue("input1"))
			}
			http.ServeFile(w, r, "template.html")
		}),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	if err := s.ListenAndServe(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
