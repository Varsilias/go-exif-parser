package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/varsilias/exif-parser/parser"
)

func main() {
	image := flag.String("image", "", "The path to image to be processed")
	output := flag.String("output", "output.json", "The path to the file where the result of the parsing will be stored")

	// Register short flag
	flag.StringVar(image, "i", "", "Short for --image")
	flag.StringVar(output, "o", "output.json", "Short for --output")

	flag.Parse()

	_, err := os.Stat(*image)
	if err != nil {
		fmt.Println("Error reading file")
		os.Exit(1)
	}

	payload, err := parser.ParseImageFile(*image)
	if err != nil {
		fmt.Printf("Error parsing image file: %v", err)
		os.Exit(1)
	}

	content, err := json.MarshalIndent(payload, "", " ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
		os.Exit(1)
	}
	err = os.WriteFile(*output, content, 0644)

	if err != nil {
		log.Fatalf("Failed to write file: %v", err)
		os.Exit(1)
	}

	fmt.Println(string(content))

}
