package log

import (
	"html/template"
	"net/http"
	"os/exec"
)

type handler struct {
	templates *template.Template
}

func NewHandler(templates *template.Template) *handler {
	return &handler{
		templates: templates,
	}
}

func (h *handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("docker", "compose", "-f", "/home/dima/Study/Go/avito_informer/deploy/compose/core/docker-compose.yml", "logs", "collector", "http", "notification")
	out, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, "Error getting logs: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"Logs": string(out),
	}

	err = h.templates.ExecuteTemplate(w, "log", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
