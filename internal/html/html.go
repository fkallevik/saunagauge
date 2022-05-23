package html

import (
	"embed"
	"io"
	"io/fs"
	"strings"
	"text/template"
)

//go:embed * static/*
var EmbedFS embed.FS

type HomeParams struct {
	Title       string
	Time        string
	Temperature float32
	Humidity    float32
}

func Home(w io.Writer, fsys fs.FS, p HomeParams) error {
	home := parse("home.html", fsys)
	return home.Execute(w, p)
}

var funcs = template.FuncMap{
	"uppercase": func(v string) string {
		return strings.ToUpper(v)
	},
}

func parse(file string, fsys fs.FS) *template.Template {
	return template.Must(
		template.New("layout.html").Funcs(funcs).ParseFS(fsys, "layout.html", file))

}
