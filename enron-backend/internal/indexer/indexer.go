package indexer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"time"

	"enron_corp_indexer.snowartz/internal/config"
	"enron_corp_indexer.snowartz/internal/email"
)

func CreateIndex(cfg *config.Config) error {
	url := fmt.Sprintf("%s/api/index", cfg.ZincURL)
	payload := map[string]interface{}{
		"name":         cfg.IndexName,
		"storage_type": "disk",
		"shard_num":    runtime.NumCPU(),
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"from": map[string]interface{}{
					"type":          "keyword",
					"index":         true,
					"store":         true,
					"highlightable": true,
				},
				"to": map[string]interface{}{
					"type":          "keyword",
					"index":         true,
					"store":         true,
					"highlightable": true,
				},
				"subject": map[string]interface{}{
					"type":          "text",
					"index":         true,
					"store":         true,
					"highlightable": true,
				},
				"body": map[string]interface{}{
					"type":          "text",
					"index":         true,
					"store":         true,
					"highlightable": true,
				},
				"date": map[string]interface{}{
					"type":     "date",
					"format":   "Mon, 02 Jan 2006 15:04:05 -0700",
					"locale":   "en",
					"index":    true,
					"sortable": true,
				},
			},
		},
		"settings": map[string]interface{}{},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error al serializar payload: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creando request: %v", err)
	}

	req.Header.Set("Authorization", cfg.BasicAuthToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error enviando request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("respuesta inesperada: %d, cuerpo: %s", resp.StatusCode, string(body))
	}

	return nil
}

func BulkIndexEmails(emails []email.Email, cfg *config.Config) error {
	url := fmt.Sprintf("%s/api/_bulkv2", cfg.ZincURL)
	payload := map[string]interface{}{
		"index":   cfg.IndexName,
		"records": emails,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error al serializar el lote: %v", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creando request: %v", err)
	}

	req.Header.Set("Authorization", cfg.BasicAuthToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error enviando request bulk: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("respuesta inesperada: %d, cuerpo: %s", resp.StatusCode, string(body))
	}
	return nil
}
