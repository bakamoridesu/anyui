package runner

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
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

func command(cfg config, r *http.Request) (*exec.Cmd, error) {
	cParams := []string{}
	cName := cfg.Command
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return nil, err
	}
	for _, flag := range cfg.Flags {
		value := ""
		if r.FormValue(flag.Param) != "" {
			value = r.FormValue(flag.Param)
		} else if flag.Default != "" {
			value = flag.Default
		}
		if value != "" {
			cParams = append(cParams, flag.Name+"="+value)
		}
	}

	for _, arg := range cfg.Args {
		value := ""
		if r.FormValue(arg.Param) != "" {
			value = r.FormValue(arg.Param)
		} else {
			value = arg.Default
		}
		cParams = append(cParams, value)
	}
	return exec.Command(cPath, cParams...), nil
}

type Runner struct {
	template string
	config   string
}

func NewRunner(template, config string) (*Runner, error) {
	if _, err := os.Stat(template); errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	if _, err := os.Stat(config); errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	return &Runner{template, config}, nil
}

func (r Runner) Run() error {
	cfg := &config{}
	cfgFile, err := os.Open(r.config)
	if err != nil {
		return err
	}
	defer cfgFile.Close()

	if err := json.NewDecoder(cfgFile).Decode(cfg); err != nil {
		return err
	}

	// check if template found

	cName, cParams := GetOpen()
	cParams = append(cParams, "http://127.0.0.1:8080")
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}

	go func() {
		exec.Command(cPath, cParams...).Run()
	}()

	s := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: Handler(cfg, r.template),
		// ReadTimeout:  10 * time.Second,
		// WriteTimeout: 10 * time.Second,
	}
	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
