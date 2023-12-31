package builder

import "html/template"

type InputData struct {
	Id    string
	Top   string
	Left  string
	Width string
}

type FormData struct {
	Inputs template.HTML
}

const Input = `
<input type="text" id="{{ .Id}}" name="{{ .Id}}" style="top:{{ .Top}};left:{{ .Left}};width:{{ .Width}}"/>`

const Form = `
<form>
	{{.Inputs}}
</form>
`
