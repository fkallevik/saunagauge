package server

import (
	"html/template"
	"net/http"
)

type htmlResponse struct {
	html template.HTML
}

func respond(w http.ResponseWriter, resp htmlResponse) error {
	w.Header().Set("Content-Type", "text/html")

	if _, err := w.Write([]byte(resp.html)); err != nil {
		return err
	}

	return nil
}
