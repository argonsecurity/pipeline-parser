#!/usr/bin/env bash

set -e

function print_usage() {
    echo "Usage: make tag LEVEL=<patch|minor|major>"
}

function validate_input() {
    if [ -z "${LEVEL}" ]; then
        echo "Parameter LEVEL is not allowed to be empty"
        print_usage
        exit 1
    fi

    case "${LEVEL}" in
        patch|minor|major) ;;
        *)
            echo "LEVEL value is not valid"
            print_usage
            exit 2
        ;;
    esac
}

# Usage: incr_semver <VERSION> <patch|minor|major>
# Example: incr_semver 1.2.3 patch
# Output: 1.2.4
function incr_semver() { 
    IFS='.' read -ra ver <<< "$1"
    [[ "${#ver[@]}" -ne 3 ]] && echo "Invalid semver string" && return 1
    [[ "$#" -eq 1 ]] && level='patch' || level=$2

    patch=${ver[2]}
    minor=${ver[1]}
    major=${ver[0]}

    case $level in
        patch)
            patch=$((patch+1))
        ;;
        minor)
            patch=0
            minor=$((minor+1))
        ;;
        major)
            patch=0
            minor=0
            major=$((major+1))
        ;;
        *)
            return 2
    esac
    echo "${major}.${minor}.${patch}"
}

validate_input

git checkout main
git fetch --tags --all
git pull

LATEST_TAG=$(git describe --tags --abbrev=0)
TAG=$(incr_semver ${LATEST_TAG} ${LEVEL})
echo "Using tag: ${TAG}"

git tag -a ${TAG} -m ${TAG}
git push --tag