package builder

import (
	"bufio"
	"bytes"
	"html/template"
	"io"
	"net/http"
	"os"
)

const address = "127.0.0.1:8081"

var fileName = "t.html"

func CreateTemplate() error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	// TODO: make template in go file
	start, err := os.Open("html_start.html")
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(f)
	io.Copy(writer, start)

	return nil
}

func Run() error {
	tmpl()
	// CreateTemplate()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "builder.html")
	})

	// go open.Open(address)

	http.ListenAndServe(address, nil)
	return nil
}

func tmpl() {

	t, _ := template.New("webpage").Parse(Form)

	inputTempl, _ := template.New("input").Parse(Input)

	input_data := InputData{"param1", "40px", "100px", "200px"}
	input1_data := InputData{"param2", "100px", "100px", "300px"}

	var buf bytes.Buffer
	inputTempl.Execute(&buf, input_data)
	inputTempl.Execute(&buf, input1_data)
	data := FormData{template.HTML(buf.String())}

	_ = t.Execute(os.Stdout, data)

}
