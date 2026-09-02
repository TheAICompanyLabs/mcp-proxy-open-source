#!/bin/bash
# scripts/install.sh
set -e

install_mcp_proxy() {
    GITHUB_REPO="your-github-username/mcp-proxy"
    APP_NAME="mcp-proxy"

    echo "📥 Installing Universal MCP Proxy..."

    OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
    ARCH="$(uname -m)"

    if [ "$ARCH" = "x86_64" ]; then
        ARCH="amd64"
    elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
        ARCH="arm64"
    else
        echo "❌ Unsupported architecture: ${ARCH}"
        exit 1
    fi

    # Fetch the latest release version
    VERSION=$(curl -sL "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

    if [ -z "$VERSION" ]; then
        echo "❌ Failed to fetch latest version. Ensure your GitHub repository is public."
        exit 1
    fi

    FILE_NAME="${APP_NAME}-${OS}-${ARCH}.tar.gz"
    DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/download/${VERSION}/${FILE_NAME}"

    echo " -> Downloading ${VERSION} for ${OS}/${ARCH}..."

    curl -sSfL -o ${FILE_NAME} ${DOWNLOAD_URL}
    tar -xzf ${FILE_NAME}

    echo " -> Installing to /usr/local/bin (may require sudo password)..."
    sudo mv ${APP_NAME}-${OS}-${ARCH} /usr/local/bin/${APP_NAME}

    rm ${FILE_NAME}
    echo "✅ Installation complete! Run '${APP_NAME}' to get started."
}

# Execute the function only after the entire script has downloaded
install_mcp_proxy
