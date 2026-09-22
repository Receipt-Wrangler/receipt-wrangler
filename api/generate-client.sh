#!/bin/bash

# Function to display usage
show_usage() {
    echo "Usage: $0 <platform> <output-dir>"
    echo ""
    echo "Arguments:"
    echo "  platform     Either 'desktop', 'mobile', or 'mcp'"
    echo "  output-dir   Directory path for generated code output"
    echo ""
    echo "Example:"
    echo "  $0 desktop /home/user/project/src/open-api"
    exit 1
}

# Check if correct number of arguments is provided
if [ $# -ne 2 ]; then
    echo "Error: Exactly 2 arguments are required"
    show_usage
fi

platform=$1
output_dir=$2

# Validate platform argument
if [ "$platform" != "desktop" ] && [ "$platform" != "mobile" ] && [ "$platform" != "mcp" ]; then
    echo "Error: Platform must be either 'desktop', 'mobile', or 'mcp'"
    show_usage
fi

# Set generator based on platform
if [ "$platform" = "desktop" ]; then
    generator="typescript-angular"
elif [ "$platform" = "mcp" ]; then
    generator="typescript"
else
    generator="dart-dio"
fi

# Execute the OpenAPI generator command
npx @openapitools/openapi-generator-cli generate \
    -i swagger.yml \
    -g "$generator" \
    -o "$output_dir"

# Check if command executed successfully
if [ $? -ne 0 ]; then
    echo "Error: API code generation failed"
    exit 1
fi

# The dart-dio generator emits four things the app cannot use as generated - two that do not
# compile and two `fallback: true` annotations that keep an unrecognized enum value from
# failing the whole payload. Re-applying them here makes them part of generation instead of a
# step someone has to remember; the script fails loudly if one no longer applies.
if [ "$platform" = "mobile" ]; then
    script_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
    if ! "$script_dir/patches/apply-dart-dio-patches.sh" "$output_dir"; then
        echo "Error: dart-dio patches could not be applied"
        exit 1
    fi
fi

echo "API code successfully generated in: $output_dir"
