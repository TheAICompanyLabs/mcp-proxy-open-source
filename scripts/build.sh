#!/bin/bash
# scripts/build.sh

# Exit on any error
set -e

APP_NAME="mcp-proxy"
DIST_DIR="./dist"

# Clean previous builds
rm -rf ${DIST_DIR}
mkdir -p ${DIST_DIR}

echo "🔨 Building Universal Binaries for ${APP_NAME}..."

# Define the build matrix (OS/ARCH)
PLATFORMS=(
    "darwin/amd64"   # macOS Intel
    "darwin/arm64"   # macOS Apple Silicon (M1/M2/M3)
    "linux/amd64"    # Linux x86_64
    "linux/arm64"    # Linux ARM
    "windows/amd64"  # Windows x86_64
)

for PLATFORM in "${PLATFORMS[@]}"; do
    # Split the platform string
    GOOS=${PLATFORM%/*}
    GOARCH=${PLATFORM#*/}
    
    # Define output name
    OUTPUT_NAME="${APP_NAME}-${GOOS}-${GOARCH}"
    if [ "$GOOS" = "windows" ]; then
        OUTPUT_NAME="${OUTPUT_NAME}.exe"
    fi
    
    echo " -> Compiling ${GOOS}/${GOARCH}..."
    
    # Run the Go cross-compiler
    GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 go build -ldflags="-s -w" -o ${DIST_DIR}/${OUTPUT_NAME} .
    
    # Compress the binary for distribution (macOS/Linux)
    if [ "$GOOS" != "windows" ]; then
        tar -czvf ${DIST_DIR}/${OUTPUT_NAME}.tar.gz -C ${DIST_DIR} ${OUTPUT_NAME}
        rm ${DIST_DIR}/${OUTPUT_NAME} # Cleanup raw binary
    else
        zip -j -q ${DIST_DIR}/${OUTPUT_NAME}.zip ${DIST_DIR}/${OUTPUT_NAME}
        rm ${DIST_DIR}/${OUTPUT_NAME} # Cleanup raw exe
    fi
done

echo "✅ All builds complete! Check the ${DIST_DIR}/ directory."
