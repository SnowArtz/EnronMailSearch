// internal/processor/processor.go
package processor

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"enron_corp_indexer.snowartz/internal/config"
	"enron_corp_indexer.snowartz/internal/email"
	"enron_corp_indexer.snowartz/internal/indexer"
)

func ProcessFolders(folderPaths []string, cfg *config.Config) {
	var combinedEmails []email.Email
	for _, folder := range folderPaths {
		emails, err := email.ReadEmails(folder)
		if err != nil {
			log.Printf("Error al leer emails en %s: %v", folder, err)
			continue
		}
		combinedEmails = append(combinedEmails, emails...)
	}
	if len(combinedEmails) == 0 {
		log.Printf("No se encontraron emails en las carpetas: %v", folderPaths)
		return
	}
	err := indexer.BulkIndexEmails(combinedEmails, cfg)
	if err != nil {
		log.Printf("Error indexando lote para carpetas %v: %v", folderPaths, err)
	} else {
		fmt.Printf("Se indexaron correctamente %d emails\n", len(combinedEmails))
	}
}

func ProcessDirectory(inputDir string, cfg *config.Config) {
	entries, err := os.ReadDir(inputDir)
	if err != nil {
		log.Fatalf("Error leyendo directorio %s: %v", inputDir, err)
	}
	var folderPaths []string
	for _, entry := range entries {
		if entry.IsDir() {
			folderPaths = append(folderPaths, filepath.Join(inputDir, entry.Name()))
		}
	}

	var groups [][]string
	for i := 0; i < len(folderPaths); i += cfg.FoldersPerRoutine {
		end := i + cfg.FoldersPerRoutine
		if end > len(folderPaths) {
			end = len(folderPaths)
		}
		groups = append(groups, folderPaths[i:end])
	}
	start_time := time.Now()
	var wg sync.WaitGroup
	for _, group := range groups {
		wg.Add(1)
		go func(folders []string) {
			defer wg.Done()
			ProcessFolders(folders, cfg)
		}(group)
	}
	wg.Wait()
	elapsedTime := time.Since(start_time)
	fmt.Printf("Tiempo transcurrido: %s\n", elapsedTime)
}
