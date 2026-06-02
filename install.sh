#!/bin/sh
# Install the binpath CLI on Linux or macOS.
#
#   curl -fsSL https://raw.githubusercontent.com/loosestack-dev/binpathutil/main/install.sh | sh
#
# Install a specific release tag (default: latest):
#   ... | sh -s -- v1.0.0            # or:  BINPATH_VERSION=v1.0.0 ... | sh
#
# Override the install directory (default: ~/.local/bin):
#   BINPATH_INSTALL_DIR=/usr/local/bin ... | sh
set -eu

REPO="loosestack-dev/binpathutil"
BIN="binpath"
INSTALL_DIR="${BINPATH_INSTALL_DIR:-$HOME/.local/bin}"
VERSION="${1:-${BINPATH_VERSION:-latest}}"

err() { echo "install.sh: $*" >&2; exit 1; }

# --- detect platform -------------------------------------------------------
os=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$os" in
	linux | darwin) ;;
	*) err "unsupported OS '$os' (only linux and macOS are supported)" ;;
esac

arch=$(uname -m)
case "$arch" in
	x86_64 | amd64) arch="amd64" ;;
	arm64 | aarch64) arch="arm64" ;;
	*) err "unsupported architecture '$arch' (only amd64 and arm64 are supported)" ;;
esac

asset="${BIN}_${os}_${arch}.tar.gz"
if [ "$VERSION" = "latest" ]; then
	base="https://github.com/$REPO/releases/latest/download"
else
	base="https://github.com/$REPO/releases/download/$VERSION"
fi

# --- download helper (curl or wget) ----------------------------------------
download() {
	# download <url> <dest>
	curl -fsSL "$1" -o "$2"
}

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

echo "Downloading $base/$asset"
download "$base/$asset" "$tmp/$asset" || err "download failed (is $VERSION published for ${os}/${arch}?)"

# --- verify checksum (best-effort) -----------------------------------------
if download "$base/checksums.txt" "$tmp/checksums.txt" 2>/dev/null; then
	expected=$(grep " $asset\$" "$tmp/checksums.txt" | awk '{print $1}')
	if [ -n "$expected" ]; then
		if command -v sha256sum >/dev/null 2>&1; then
			actual=$(sha256sum "$tmp/$asset" | awk '{print $1}')
		elif command -v shasum >/dev/null 2>&1; then
			actual=$(shasum -a 256 "$tmp/$asset" | awk '{print $1}')
		fi
		if [ -n "${actual:-}" ] && [ "$actual" != "$expected" ]; then
			err "checksum mismatch for $asset (expected $expected, got $actual)"
		fi
	fi
fi

# --- extract and install ---------------------------------------------------
tar -xzf "$tmp/$asset" -C "$tmp"
mkdir -p "$INSTALL_DIR"
install -m 0755 "$tmp/$BIN" "$INSTALL_DIR/$BIN"
echo "Installed $BIN to $INSTALL_DIR/$BIN"

# --- PATH guidance ---------------------------------------------------------
case ":$PATH:" in
	*":$INSTALL_DIR:"*) ;;
	*) echo "Note: $INSTALL_DIR is not on your PATH. Add this to your shell profile:"
		echo "  export PATH=\"$INSTALL_DIR:\$PATH\"" ;;
esac
