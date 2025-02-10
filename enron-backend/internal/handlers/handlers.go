// internal/handlers/handlers.go
package handlers

import (
	"encoding/json"
	"net/http"

	"enron_corp_indexer.snowartz/internal/config"
	"enron_corp_indexer.snowartz/internal/indexer"
	"enron_corp_indexer.snowartz/internal/processor"
	"enron_corp_indexer.snowartz/internal/search"

	"enron_corp_indexer.snowartz/internal/email"
	"github.com/go-chi/chi/v5"
)

type APIHandler struct {
	Config *config.Config
}

func (h *APIHandler) RegisterRoutes(r chi.Router) {
	r.Post("/index", h.StartIndexing)
	r.Post("/search", h.SearchEmails)
}

func (h *APIHandler) StartIndexing(w http.ResponseWriter, r *http.Request) {
	err := indexer.CreateIndex(h.Config)
	if err != nil {
		http.Error(w, "Error creating index: "+err.Error(), http.StatusInternalServerError)
		return
	}
	go processor.ProcessDirectory(h.Config.InputDir, h.Config)
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Indexing process started"))
}

func (h *APIHandler) SearchEmails(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Query   string            `json:"query"`
		Filters map[string]string `json:"filters"`
		From    int               `json:"from"`
		Size    int               `json:"size"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	total, emails, err := search.SearchEmails(req.Query, req.Filters, req.From, req.Size, h.Config)
	if err != nil {
		http.Error(w, "Error searching emails: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Total  int           `json:"total"`
		Emails []email.Email `json:"emails"`
	}{
		Total:  total,
		Emails: emails,
	}
	json.NewEncoder(w).Encode(response)
}
