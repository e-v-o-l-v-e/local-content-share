FROM alpine:latest

# Install curl for downloading the binary
RUN apk add --no-cache curl

# Download the latest release binary
RUN curl -sL $(curl -s https://api.github.com/repos/tanq16/local-content-share/releases/latest | grep "browser_download_url.*local-content-share-linux-amd64\"" | cut -d '"' -f 4) -o local-content-share && \
    chmod +x local-content-share

# Create app directory
WORKDIR /app

# Create data directory for persistent storage
RUN mkdir data

# Expose the port the app runs on
EXPOSE 8080

# Run the binary
CMD ["./local-content-share"]
