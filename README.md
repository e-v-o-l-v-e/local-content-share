# Local Content Share

<div align="center">
  <img src="assets/logo.png" alt="Local Content Share Logo" width="200"><br>

  <a href="https://github.com/tanq16/local-content-share/actions/workflows/binary-build.yml"><img alt="Build Workflow" src="https://github.com/tanq16/local-content-share/actions/workflows/binary-build.yml/badge.svg"></a>&nbsp;<a href="https://github.com/tanq16/local-content-share/actions/workflows/docker-publish.yml"><img alt="Container Workflow" src="https://github.com/tanq16/local-content-share/actions/workflows/docker-publish.yml/badge.svg"></a><br>
  <a href="https://github.com/Tanq16/local-content-share/releases"><img alt="GitHub Release" src="https://img.shields.io/github/v/release/tanq16/local-content-share"></a>&nbsp;<a href="https://hub.docker.com/r/tanq16/local-content-share"><img alt="Docker Pulls" src="https://img.shields.io/docker/pulls/tanq16/local-content-share"></a><br>
</div>

A simple web application for sharing content (files and text) within your local network across any device. The app can be launched via its binary or as a container. The primary features are:

- Make arbitrary text content available for viewing on any machine/smartphone in the local network
- Upload files and make them available across your local network
- Access content through a clean and modern interface with dark mode support that maintains its good looks even on mobile aspect ratio
- Pure HTTP API and no use of websockets, which means no external communications (like for TURN) needed
- It can also be installed as a PWA (so it shows as an icon in a mobile home screens)

> [!NOTE]
> This application is meant to be deployed within your homelab only. There is no authentication mechanism here.

## Screenshots

| | Desktop View | Mobile View |
| --- | --- | --- |
| Light | <img src="assets/desktop-light.png" alt="Desktop Light Mode"> | <img src="assets/mobile-light.png" alt="Mobile Light Mode"> |
| Dark | <img src="assets/desktop-dark.png" alt="Desktop Dark Mode"> | <img src="assets/mobile-dark.png" alt="Mobile Dark Mode"> |

## Usage

### Using Docker

```bash
mkdir $HOME/.localcontentshare
```
```bash
docker run --name local-content-share -p 8080:8080 \
-v $HOME/.localcontentshare:/app/data tanq16/local-content-share:latest
```

The application will be available at `http://localhost:8080`

You could also use the following compose file with container managers like Portainer and Dockge:

```yaml
services:
  contentshare:
    image: tanq16/local-content-share:main
    container_name: local-content-share
    volumes:
      - /home/tanq/lcshare:/app/data
    ports:
      - 5000:8080
```

> [!WARNING]
> The public image built via GitHub actions only builds an x86-64 image. If you need an ARM variant, just run `docker build -t lcshare:local .` after cloning the repo. Then for the image use `lcshare:local` instead of the tag mentioned above.

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

## Building Binary from Source

Requirements: `Go 1.23` or later

```bash
git clone https://github.com/tanq16/local-content-share.git
cd local-content-share
go build .
./local-content-share
```

## Interface Tips

- To share text content:
   - Type or paste your text in the text area
   - Click the send button (like the telegram arrow)
   - It will set a timestamp-based file name
- To rename files (text items only):
   - Click the pencil icon and provide the new name
   - It will automatically append 4 random digits if your input isn't unique
- To share files:
   - Click the upload button
   - Select your file
   - It will automatically append 4 random digits if filename isn't unique
- To view/download content:
   - Text content: click the eye icon (shows raw text)
   - Files: click the download icon
- To delete content:
   - Click the trash icon next to the item box

## Directory Structure

The application creates a `data` directory to store all uploaded files and text content. Make sure the application has write permissions for the directory where it runs.
