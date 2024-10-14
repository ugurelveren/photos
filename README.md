
# Image Gallery Project

This project consists of two main parts:

1. A web-based image gallery that displays images from the `images` folder, with zoom and navigation functionality.
2. A CLI tool built in Go that scans the `images` folder, generates a JSON file with metadata for each image, and saves it in the `data` folder.

## Features

### Web-Based Image Gallery:
- **Responsive Thumbnails**: Displays small clickable thumbnails of the images.
- **Zoom-In Modal**: Clicking on an image opens a zoomed-in view of the image.
- **Navigation**: Allows users to navigate between images in the modal (Next/Previous buttons).
- **Image Metadata**: Each image includes a caption for additional context.

### Go CLI Tool:
- **Recursive Scanning**: Recursively scans the `images` folder and its subfolders for supported image formats.
- **Supported Image Formats**: Includes `.jpg`, `.jpeg`, `.png`, `.gif`, `.bmp`, `.tiff`, and `.webp`.
- **JSON Generation**: Creates an `images.json` file containing the metadata for each image, including:
  - **Title**: Automatically generated from the file name (without the extension).
  - **URL**: The relative path to the image from the `images` folder.
  - **Alt Text**: Generated automatically from the title (can be customized later if needed).

## Project Structure

```
/project-root
  ├── /src              # The source code for the Go CLI tool
  ├── /images           # Folder that contains the image files and subfolders
  ├── /data             # The folder where the generated JSON file is saved
  └── README.md         # This file!
```

## How to Use

### Web-Based Image Gallery:
1. **Clone the repository**:
   ```bash
   git clone git@github.com:your-username/image-gallery.git
   ```

2. **Navigate into the project directory**:
   ```bash
   cd image-gallery
   ```

3. **Open the `index.html` file** in a browser to view the image gallery.

4. **To add new images**:
   - Update the `data/images.json` file with new image URLs and alt text.
   - Add the corresponding image files to the `images/` folder.

### Go CLI Tool:
1. **Prerequisites**:
   - **Go**: You need to have Go installed on your machine. You can download it from [here](https://golang.org/dl/).
   - **Images Folder**: Place your images inside an `images` folder in the project’s parent directory (relative to the `src` folder).

2. **Running the CLI**:
   Navigate to the `src` folder in your terminal:

   ```bash
   cd src
   ```

   Run the CLI tool:

   ```bash
   go run main.go
   ```

   The tool will:
   - Look for an `images` folder in the parent directory.
   - Scan the folder and its subdirectories for image files.
   - Generate a `data/images.json` file containing the image metadata.

### Example of the Generated JSON:
```json
{
  "images": [
    {
      "title": "golden",
      "url": "golden.jpeg",
      "alt": "golden"
    },
    {
      "title": "sunset",
      "url": "subfolder/sunset.jpeg",
      "alt": "sunset"
    }
  ]
}
```

## Customization

- **Alt Text**: The alt attribute in the JSON file defaults to the image title, but you can manually edit the JSON to add more descriptive alt text if needed.
- **Additional Features**: Feel free to extend the CLI tool to handle custom metadata, filtering, or additional image processing!

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.


## License

- **Source Code**: This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
- **Images**: All images in the `/images/` folder are licensed under the **Creative Commons Attribution-NonCommercial 4.0 International Public License**. See the [LICENSE](IMAge-LICENSE). You can learn more about this license [here](https://creativecommons.org/licenses/by-nc/4.0/).
