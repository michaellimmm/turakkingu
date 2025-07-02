package api

import (
	"html/template"
	"net/http"
)

func render404(w http.ResponseWriter) {
	tmpl, err := template.ParseFiles("./web/404.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	_ = tmpl.Execute(w, nil)
}
