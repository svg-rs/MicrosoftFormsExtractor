package main

import (
	extractor "MicrosoftFormsExtractor/pkg"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load("info.env"); err != nil {
		log.Println("Error loading .env file")
		time.Sleep(10 * time.Second)
	}

	// Read the URL and AUTH from environment variables
	url := os.Getenv("URL")
	auth := os.Getenv("AUTH")

	if url == "" || auth == "" {
		log.Println("URL or AUTH environment variable is missing")
		time.Sleep(10 * time.Second)
	}

	// Call the Extract function with URL and AUTH
	fmt.Println("URL:", url)
	fmt.Println("AUTH:", auth)
	extractor.Extract(url, auth)
}
