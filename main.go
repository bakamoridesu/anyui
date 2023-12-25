package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type arg struct {
	Param   string `json:"param"`
	Default string `json:"default"`
}

type flag struct {
	Name    string `json:"name"`
	Param   string `json:"param"`
	Default string `json:"default"`
}

type config struct {
	Command string `json:"command"`
	Args    []arg  `json:"args"`
	Flags   []flag `json:"flags"`
}

func (c config) String() string {
	args := "args:\n"
	for _, arg := range c.Args {
		args += fmt.Sprintf("\tparam:%s, defalt:%s", arg.Param, arg.Default)
	}
	flags := "flags:\n"
	for _, flag := range c.Flags {
		flags += fmt.Sprintf("\tname:%s, param:%s, defalt:%s", flag.Name, flag.Param, flag.Default)
	}

	return fmt.Sprintf("command:%s\n%s\n%s", c.Command, args, flags)
}

func run() error {
	cfg := &config{}
	cfgFile, err := os.Open("config.json")
	if err != nil {
		return err
	}
	defer cfgFile.Close()

	if err := json.NewDecoder(cfgFile).Decode(cfg); err != nil {
		return err
	}

	fmt.Println(cfg)
	// check if template found

	cName, cParams := GetOpen()
	cParams = append(cParams, "http://127.0.0.1:8080")
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}

	go func() {
		time.Sleep(2 * time.Second)
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
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
