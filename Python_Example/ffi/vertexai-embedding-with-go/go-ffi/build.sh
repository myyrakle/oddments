#!/bin/bash

SCRIPT_DIR="$(dirname "${BASH_SOURCE[0]}")"

case "$(uname -s)" in
    Linux*)   OUTPUT="libvertexai.so" ;;
    Darwin*)  OUTPUT="libvertexai.dylib" ;;
    MINGW*|MSYS*|CYGWIN*) OUTPUT="libvertexai.dll" ;;
    *) echo "Unsupported platform"; exit 1 ;;
esac

(cd "$SCRIPT_DIR" && CGO_ENABLED=1 go build -buildmode=c-shared -ldflags="-s -w" -trimpath -o "$OUTPUT" lib.go)
