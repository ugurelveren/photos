package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
)

// Image struct to represent the image metadata
type Image struct {
	Title     string `json:"title"`
	URL       string `json:"url"`
	Alt       string `json:"alt"`
	Thumbnail string `json:"thumbnail"`
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
			thumbPath := handleCreateThumbnail(fullPath, baseDir)
			// Create the path relative to the base project directory
			relativePath, _ := filepath.Rel(baseDir, fullPath)
			relativeThumbPath, _ := filepath.Rel(baseDir, thumbPath)

			relativePath = strings.Replace(relativePath, "\\", "/", -1) // Handle Windows paths

			// Add the image to the image map
			imageMap.Images = append(imageMap.Images, Image{
				Title:     title,
				URL:       relativePath,
				Thumbnail: relativeThumbPath,
				Alt:       title, // You can add more specific alt text here if needed
			})
		}
	}

	return nil
}

// ensureDataFolder ensures the "data" folder exists in the parent directory
func ensureFolderExists(parentDir string, folderName string) error {
	dataFolder := filepath.Join(parentDir, folderName)
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
	if err := ensureFolderExists(parentDir, "data"); err != nil {
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

func handleCreateThumbnail(file string, outputPath string) string {
	img, format, err := openImage(file)
	if err != nil {
		fmt.Println("Error opening image:", err)
	}
	fileName := filepath.Base(file)
	outputPath = filepath.Join(outputPath, "assets", "thumbnail", fileName)

	resizedImg := resizeImage(img)

	thumbnail := cropCenter(resizedImg, 415, 415)
	if err := saveImage(thumbnail, outputPath, format); err != nil {
		fmt.Println("Error saving thumbnail:", err)
	}

	return outputPath
}

// Resize the image to 1/5th of its original size
func resizeImage(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx() / 5
	height := bounds.Dy() / 5

	resizedImg := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(resizedImg, resizedImg.Bounds(), img, bounds, draw.Over, nil)

	return resizedImg
}

// Open the source image and decode it
func openImage(filePath string) (image.Image, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return nil, "", err
	}

	return img, format, nil
}

// Save the new thumbnail image
func saveImage(img image.Image, filePath string, format string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		return jpeg.Encode(file, img, &jpeg.Options{Quality: 80})
	case "png":
		return png.Encode(file, img)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// Generate the thumbnail
func cropCenter(img image.Image, cropWidth, cropHeight int) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Calculate the crop's top-left corner coordinates
	startX := (width - cropWidth) / 2
	startY := (height - cropHeight) / 2

	// Define the crop rectangle
	cropRect := image.Rect(startX, startY, startX+cropWidth, startY+cropHeight)

	// Crop the image
	croppedImg := img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(cropRect)

	return croppedImg
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

	if err := ensureFolderExists(projectDir, "assets/thumbnail"); err != nil {
		log.Fatalf("Error finding thumbnail folder: %v", err)
	}

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
