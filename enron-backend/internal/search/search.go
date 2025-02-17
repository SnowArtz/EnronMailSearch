package search

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"enron_corp_indexer.snowartz/internal/config"
	"enron_corp_indexer.snowartz/internal/email"
)

// SearchRequest representa la estructura de la consulta a ZincSearch
type SearchRequest struct {
	Query struct {
		Bool struct {
			Must []struct {
				Bool *struct {
					Should             []map[string]map[string]string `json:"should,omitempty"`
					MinimumShouldMatch int                            `json:"minimum_should_match,omitempty"`
				} `json:"bool,omitempty"`
				QueryString *struct {
					Query string `json:"query,omitempty"`
				} `json:"query_string,omitempty"`
			} `json:"must"`
		} `json:"bool"`
	} `json:"query"`
	Size      int `json:"size"`
	From      int `json:"from"`
	Highlight struct {
		Fields map[string]interface{} `json:"fields"`
	} `json:"highlight"`
}

// SearchResponse representa la estructura de la respuesta de ZincSearch
type SearchResponse struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []struct {
			Id        string      `json:"_id"`
			Source    email.Email `json:"_source"`
			Highlight struct {
				Body    []string `json:"body"`
				Subject []string `json:"subject"`
				From    []string `json:"from"`
				To      []string `json:"to"`
			} `json:"highlight"`
		} `json:"hits"`
	} `json:"hits"`
}

// SearchEmails realiza una búsqueda en ZincSearch
func SearchEmails(query string, filters map[string]string, from int, size int, cfg *config.Config) (int, []email.Email, error) {
	searchRequest := SearchRequest{
		Size: size,
		From: from,
	}

	// Configurar highlight
	searchRequest.Highlight.Fields = map[string]interface{}{
		"body":    struct{}{},
		"subject": struct{}{},
		"from":    struct{}{},
		"to":      struct{}{},
	}

	// Construcción de la consulta principal con match_phrase
	boolQuery := struct {
		Should             []map[string]map[string]string `json:"should,omitempty"`
		MinimumShouldMatch int                            `json:"minimum_should_match,omitempty"`
	}{
		Should: []map[string]map[string]string{
			{"match_phrase": {"body": query}},
			{"match_phrase": {"subject": query}},
			{"match_phrase": {"from": query}},
		},
		MinimumShouldMatch: 1,
	}

	searchRequest.Query.Bool.Must = append(searchRequest.Query.Bool.Must, struct {
		Bool *struct {
			Should             []map[string]map[string]string `json:"should,omitempty"`
			MinimumShouldMatch int                            `json:"minimum_should_match,omitempty"`
		} `json:"bool,omitempty"`
		QueryString *struct {
			Query string `json:"query,omitempty"`
		} `json:"query_string,omitempty"`
	}{
		Bool: &boolQuery,
	})

	// Agregar filtros adicionales si existen
	for field, value := range filters {
		searchRequest.Query.Bool.Must = append(searchRequest.Query.Bool.Must, struct {
			Bool *struct {
				Should             []map[string]map[string]string `json:"should,omitempty"`
				MinimumShouldMatch int                            `json:"minimum_should_match,omitempty"`
			} `json:"bool,omitempty"`
			QueryString *struct {
				Query string `json:"query,omitempty"`
			} `json:"query_string,omitempty"`
		}{
			QueryString: &struct {
				Query string `json:"query,omitempty"`
			}{
				Query: fmt.Sprintf("%s:%s", field, value),
			},
		})
	}

	jsonData, err := json.Marshal(searchRequest)

	if err != nil {
		return 0, nil, fmt.Errorf("error serializando la consulta: %v", err)
	}

	// Construcción del endpoint
	url := fmt.Sprintf("%s/es/%s/_search", cfg.ZincURL, cfg.IndexName)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, nil, fmt.Errorf("error creando la petición de búsqueda: %v", err)
	}

	// Configuración de encabezados
	req.Header.Set("Authorization", cfg.BasicAuthToken)
	req.Header.Set("Content-Type", "application/json")

	// Realizar la petición
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("error enviando la petición de búsqueda: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, nil, fmt.Errorf("respuesta inesperada: %d, cuerpo: %s", resp.StatusCode, string(body))
	}

	// Decodificar la respuesta
	var result SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, nil, fmt.Errorf("error decodificando la respuesta: %v", err)
	}

	// Extraer los emails de la respuesta
	var emails []email.Email

	for _, hit := range result.Hits.Hits {
		hit.Source.Id = hit.Id
		hit.Source.Highlight = hit.Highlight
		emails = append(emails, hit.Source)
	}

	return result.Hits.Total.Value, emails, nil
}
