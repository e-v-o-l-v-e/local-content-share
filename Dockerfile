FROM ubuntu:jammy

# Install curl for downloading the binary
RUN apt update && apt install curl -y

# Create app directory
WORKDIR /app

# Download the latest release binary
RUN curl -sL $(curl -s https://api.github.com/repos/tanq16/local-content-share/releases/latest | grep "browser_download_url.*local-content-share-linux-amd64\"" | cut -d '"' -f 4) -o local-content-share && \
    chmod +x local-content-share

# Create data directory for persistent storage
RUN mkdir -p data

# Expose the port the app runs on
EXPOSE 8080

# Run the binary
CMD ["/app/local-content-share"]
