package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Image struct to represent the image metadata
type Image struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Alt   string `json:"alt"`
}

// ImageMap struct to map image files
type ImageMap struct {
	Images []Image `json:"images"`
}

// isImage checks if a file is an image based on its extension
func isImage(fileName string) bool {
	extensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp"}
	for _, ext := range extensions {
		if strings.HasSuffix(strings.ToLower(fileName), ext) {
			return true
		}
	}
	return false
}

// getImageTitle generates an image title from the file name
func getImageTitle(fileName string) string {
	name := strings.TrimSuffix(fileName, filepath.Ext(fileName)) // Remove extension
	return name
}

// findImages recursively searches the folder and subfolders for image files
func findImages(dir string, imageMap *ImageMap, baseDir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		fullPath := filepath.Join(dir, file.Name())
		if file.IsDir() {
			// Recursively search subdirectories
			if err := findImages(fullPath, imageMap, baseDir); err != nil {
				return err
			}
		} else if isImage(file.Name()) {
			// Create the image metadata
			title := getImageTitle(file.Name())

			// Create the path relative to the base project directory
			relativePath, err := filepath.Rel(baseDir, fullPath)
			if err != nil {
				return err
			}
			relativePath = strings.Replace(relativePath, "\\", "/", -1) // Handle Windows paths

			// Add the image to the image map
			imageMap.Images = append(imageMap.Images, Image{
				Title: title,
				URL:   relativePath,
				Alt:   title, // You can add more specific alt text here if needed
			})
		}
	}

	return nil
}

// ensureDataFolder ensures the "data" folder exists in the parent directory
func ensureDataFolder(parentDir string) error {
	dataFolder := filepath.Join(parentDir, "data")
	if _, err := os.Stat(dataFolder); os.IsNotExist(err) {
		err := os.Mkdir(dataFolder, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

// writeJSONToFile writes the JSON data to a file in the parent "data" folder
func writeJSONToFile(imageMap *ImageMap, parentDir string) error {
	// Ensure the data folder exists in the parent directory
	if err := ensureDataFolder(parentDir); err != nil {
		return err
	}

	// Convert the image map to JSON
	jsonData, err := json.MarshalIndent(imageMap, "", "  ")
	if err != nil {
		return err
	}

	// Write JSON to data/images.json in the parent directory
	jsonFilePath := filepath.Join(parentDir, "data", "images.json")
	err = os.WriteFile(jsonFilePath, jsonData, 0644)
	if err != nil {
		return err
	}

	fmt.Println("JSON file successfully created at", jsonFilePath)
	return nil
}

// getImagesFolderPath returns the path to the parent "images" folder
func getImagesFolderPath() string {
	// Check the parent directory for the "images" folder
	parentDir, err := filepath.Abs("..")
	if err != nil {
		log.Fatalf("Error finding parent directory: %v", err)
	}

	imagesFolder := filepath.Join(parentDir, "images")
	if _, err := os.Stat(imagesFolder); os.IsNotExist(err) {
		log.Fatalf("Images folder not found in the parent directory.")
	}

	return imagesFolder
}

func main() {
	// Get the images folder path
	dir := getImagesFolderPath()

	// Get the parent directory (the project folder)
	projectDir, err := filepath.Abs("..")
	if err != nil {
		log.Fatalf("Error finding project directory: %v", err)
	}

	// Create an ImageMap to store the images
	imageMap := &ImageMap{}

	// Find images in the "images" folder, using projectDir as the base
	err = findImages(dir, imageMap, projectDir)
	if err != nil {
		log.Fatalf("Error finding images: %v", err)
	}

	// Write the image map to a JSON file in the parent directory's "data" folder
	err = writeJSONToFile(imageMap, projectDir)
	if err != nil {
		log.Fatalf("Error writing JSON file: %v", err)
	}
}
