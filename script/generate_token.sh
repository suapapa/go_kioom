#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Default configuration
LENGTH=32
FORMAT="hex"

# Print usage instructions
usage() {
    echo "Usage: $0 [-l length] [-f format]"
    echo "Options:"
    echo "  -l LENGTH  Length of random bytes to generate (default: 32)"
    echo "  -f FORMAT  Format of the token: 'hex' or 'base64' (default: hex)"
    echo "  -h         Show this help message"
    exit 1
}

# Parse command line options
while getopts "l:f:h" opt; do
    case "${opt}" in
        l)
            LENGTH=${OPTARG}
            if ! [[ "${LENGTH}" =~ ^[0-9]+$ ]]; then
                echo "Error: Length must be a positive integer." >&2
                exit 1
            fi
            ;;
        f)
            FORMAT=${OPTARG}
            if [[ "${FORMAT}" != "hex" && "${FORMAT}" != "base64" ]]; then
                echo "Error: Format must be 'hex' or 'base64'." >&2
                exit 1
            fi
            ;;
        h|*)
            usage
            ;;
    esac
done

# Generate token using openssl if available (preferred)
if command -v openssl >/dev/null 2>&1; then
    if [ "${FORMAT}" = "hex" ]; then
        openssl rand -hex "${LENGTH}"
    else
        openssl rand -base64 "${LENGTH}"
    fi
else
    # Fallback to /dev/urandom if openssl is not installed
    if [ "${FORMAT}" = "hex" ]; then
        head -c "${LENGTH}" /dev/urandom | od -An -tx1 | tr -d ' \n'
        echo ""
    else
        # macOS base64 has slightly different flags, but standard base64 works for piping
        head -c "${LENGTH}" /dev/urandom | base64 | tr -d '\n'
        echo ""
    fi
fi
