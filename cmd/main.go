package main

import (
	extractor "MicrosoftFormsExtractor/pkg"
	"fmt"
)

func main() {
	result := extractor.Extract("https://example.com/form")
	fmt.Println(result)
}
