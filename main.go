package main

import (
	"flag"
	"fmt"
	"os"

	"bakamori.com/anyui/builder"
	"bakamori.com/anyui/runner"
)

func runTemplate(template, config string) error {

	r, err := runner.NewRunner(template, config)
	if err != nil {
		return err
	}

	if err := r.Run(); err != nil {
		return err
	}

	return nil
}

func main() {
	template := flag.String("template", "template.html", "Template file")
	config := flag.String("config", "config.json", "Config file")
	run := flag.Bool("run", false, "Run UI with template")
	flag.Parse()

	if *run {
		if err := runTemplate(*template, *config); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		builder.Run()
	}
}
