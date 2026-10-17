package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {
	zipPath := run()
	if zipPath != "" {
		if explorerErr := exec.Command("explorer.exe", "/select,", zipPath).Start(); explorerErr != nil {
			fmt.Println("Error opening explorer:", explorerErr)
		}
	}

	fmt.Println()
	fmt.Println("Press Enter to close")
	fmt.Scanln()
}

func run() string {
	logsDir := filepath.Join(os.Getenv("LOCALAPPDATA"), "OnlineExpedition", "logs")

	logFiles, logsErr := filepath.Glob(filepath.Join(logsDir, "*.log"))
	if logsErr != nil {
		fmt.Println("Error getting logs:", logsErr.Error())
		return ""
	}

	if len(logFiles) == 0 {
		fmt.Println("No logs found")
		return ""
	}

	filestamp := time.Now().Format("2006-01-02_15-04-05")
	zipPath := filepath.Join(os.Getenv("USERPROFILE"), "Desktop", "OnlineExpedition-logs-"+filestamp+".zip")

	outputFile, outputFileErr := os.Create(zipPath)
	if outputFileErr != nil {
		fmt.Println("Error creating output file:", outputFileErr)
		return ""
	}

	defer func() {
		if closeErr := outputFile.Close(); closeErr != nil {
			fmt.Println("Error closing output file:", closeErr)
		}
	}()

	zipWriter := zip.NewWriter(outputFile)
	defer func() {
		if closeErr := zipWriter.Close(); closeErr != nil {
			fmt.Println("Error closing zip writer:", closeErr)
		}
	}()

	for _, logFile := range logFiles {
		log, openErr := os.Open(logFile)
		if openErr != nil {
			fmt.Println("Error opening log file:", openErr)
			continue
		}

		zipFile, createErr := zipWriter.Create(filepath.Base(logFile))
		if createErr != nil {
			fmt.Println("Error creating zip file:", createErr)
			continue
		}

		if _, copyErr := io.Copy(zipFile, log); copyErr != nil {
			fmt.Println("Error copying log file:", copyErr)
		}

		if closeErr := log.Close(); closeErr != nil {
			fmt.Println("Error closing log file:", closeErr)
		}
	}

	fmt.Println(len(logFiles), "logs zipped up")
	fmt.Println("Your logs are ready:", zipPath)
	fmt.Println("Share this file in your bug report")

	return zipPath
}
