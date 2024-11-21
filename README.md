# Local Content Share

A simple web application for sharing content within your local network. The app allows you to:

- Share and render text content (available for viewing by any machine in the local network)
- Upload and share files across your local network
- Access content through a clean, modern interface with dark mode support
- Auto-render plain text content for easy viewing

## Features

- Modern, responsive UI with dark mode support
- File upload functionality
- Text content sharing
- Clean link/content management with delete capability
- Cross-platform support (binaries available for Windows, Linux, and macOS, both AMD64 and ARM64)

## Quick Start

### Using Docker

```bash
docker run --name local-content-share -p 8080:8080 -v $(pwd)/data:/app/data tanq16/local-content-share:latest
```

The application will be available at `http://localhost:8080`

> [!WARNING]
> The docker build uses releases from the repo and only builds an x86-64 image.

### Using Binary

1. Download the appropriate binary for your system from the [latest release](https://github.com/tanq16/local-content-share/releases/latest)
2. Make the binary executable (Linux/macOS):
   ```bash
   chmod +x local-content-share-*
   ```
3. Run the binary:
   ```bash
   ./local-content-share-*
   ```

The application will be available at `http://localhost:8080`

## Building from Source

Requirements:
- Go 1.23 or later

```bash
git clone https://github.com/tanq16/local-content-share.git
cd local-content-share
go build
./local-content-share
```

## Usage

1. Access the web interface at `http://localhost:8080`
2. To share text content:
   - Type or paste your text in the text area
   - Click the send button
3. To share files:
   - Click the upload button
   - Select your file
4. To view/download content:
   - Text content: Click the eye icon
   - Files: Click the download icon
5. To delete content:
   - Click the trash icon next to the item

## Directory Structure

The application creates a `data` directory to store all uploaded files and text content. Make sure the application has write permissions for the directory where it runs.
