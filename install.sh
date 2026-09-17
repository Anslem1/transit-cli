#!/bin/sh
set -e

REPO="Anslem1/transit-cli"
BIN_NAME="transit"

# Detect OS
OS="$(uname -s)"
case "${OS}" in
  Darwin*) OS="Darwin" ;;
  Linux*)  OS="Linux" ;;
  *)
    echo "Error: Unsupported operating system ${OS}"
    exit 1
    ;;
esac

# Detect Architecture
ARCH="$(uname -m)"
case "${ARCH}" in
  x86_64*|amd64*) ARCH="x86_64" ;;
  arm64*|aarch64*) ARCH="arm64" ;;
  *)
    echo "Error: Unsupported architecture ${ARCH}"
    exit 1
    ;;
esac

# Find latest release tag if not specified
if [ -z "${VERSION}" ]; then
  VERSION=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
fi

if [ -z "${VERSION}" ]; then
  echo "Error: Could not determine latest release version."
  exit 1
fi

# Clean version (strip leading v for tarball name)
CLEAN_VERSION="${VERSION#v}"
TARBALL="transit-cli_${CLEAN_VERSION}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${VERSION}/${TARBALL}"

echo "==> Downloading Transit CLI ${VERSION} (${OS}/${ARCH})..."
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "${TMP_DIR}"' EXIT

curl -sSL "${URL}" -o "${TMP_DIR}/${TARBALL}"

echo "==> Extracting..."
tar -xzf "${TMP_DIR}/${TARBALL}" -C "${TMP_DIR}"

# Determine installation directory
INSTALL_DIR="/usr/local/bin"
if [ ! -w "${INSTALL_DIR}" ]; then
  INSTALL_DIR="${HOME}/.local/bin"
  mkdir -p "${INSTALL_DIR}"
fi

echo "==> Installing to ${INSTALL_DIR}/${BIN_NAME}..."
if [ -w "${INSTALL_DIR}" ]; then
  mv "${TMP_DIR}/transit-cli" "${INSTALL_DIR}/${BIN_NAME}"
  chmod +x "${INSTALL_DIR}/${BIN_NAME}"
else
  sudo mv "${TMP_DIR}/transit-cli" "${INSTALL_DIR}/${BIN_NAME}"
  sudo chmod +x "${INSTALL_DIR}/${BIN_NAME}"
fi

echo "==> Transit CLI successfully installed!"
"${INSTALL_DIR}/${BIN_NAME}" version
