package web

import (
	"bytes"
	"html/template"
	"net/http"
)

func render(t *template.Template, w http.ResponseWriter, name string, data any) {
	var buf bytes.Buffer

	if err := t.ExecuteTemplate(&buf, name, data); err != nil {
		http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = buf.WriteTo(w)
}
