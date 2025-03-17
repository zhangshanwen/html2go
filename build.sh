#!/bin/bash
# HTML2Go Build Script
# This script generates component mappings from UI files

set -e  # Exit on error

# Define directories
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
UI_DIR="/Users/zhangshanwen/Desktop/qor5/x/ui"
VUETIFYX_DIR="${UI_DIR}/vuetifyx"
VUETIFY_DIR="${UI_DIR}/vuetify"
OUTPUT_DIR="${SCRIPT_DIR}/embed/data"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Make sure output directory exists
mkdir -p "$OUTPUT_DIR"

# Function to generate component mappings
generate_component_map() {
    local input_dir=$1
    local output_file=$2
    local component_name=$3

    log_info "Generating component mappings for ${component_name}..."
    
    # Run the generator
    if go run ${SCRIPT_DIR}/gen/main.go -dir="${input_dir}" -output="${output_file}"; then
        log_success "Generated ${component_name} mappings: ${output_file}"
    else
        log_error "Failed to generate ${component_name} mappings"
        exit 1
    fi
}

# Main function
main() {
    # Generate VuetifyX components
    if [ -d "$VUETIFYX_DIR" ]; then
        generate_component_map "$VUETIFYX_DIR" "${OUTPUT_DIR}/vuetifyx.json" "VuetifyX"
    else
        log_warning "VuetifyX directory not found: $VUETIFYX_DIR"
    fi

    # Generate Vuetify components
    if [ -d "$VUETIFY_DIR" ]; then
        generate_component_map "$VUETIFY_DIR" "${OUTPUT_DIR}/vuetify.json" "Vuetify"
    else
        log_warning "Vuetify directory not found: $VUETIFY_DIR"
    fi

    log_success "All component mappings generated successfully"
}

# Execute the main function
main
