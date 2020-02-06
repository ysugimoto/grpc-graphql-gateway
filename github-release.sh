#!/bin/sh

### Github create release script
### This script requires following UNIX commands:
### - jq
### - file
### - curl
### - basename
### Those commands might not be installed on CI environment due to tiny OS package,
### So probably you need to install manually.
### For example, on Ubuntu: sudo apt-get install jq file curl

## Configuration for your project

### If you run this script as inline code like Jenkins, define as following:
### But we strongly recommend that specify access token at external environement variable.
#GITHUB_TOKEN="your github access token"

### Directory contains build artifacts if you have.
### If you don't have any artifacts, please keep it empty.
ASSETS_DIR="./dist"

### Determine  project repository
#REPOSITORY=":owner/:repo"

### If you are using via some CI service, you can use following server specific variable.

### In Circle CI:
REPOSITORY="${CIRCLE_PROJECT_USERNAME}/${CIRCLE_PROJECT_REPONAME}"

### In Travis CI:
#REPOSITORY="${TRAVIS_REPO_SLUG}"

### Determine release tag name
#TAG=

### If you are using via some CI service, you can use following server specific variable.

### In Circle CI:
TAG="${CIRCLE_TAG}"

### In Travis CI:
#TAG="${TRAVIS_TAG}"


####### You don't need to modify following area #######

ACCEPT_HEADER="Accept: application/vnd.github.jean-grey-preview+json"
TOKEN_HEADER="Authorization: token ${GITHUB_TOKEN}"
ENDPOINT="https://api.github.com/repos/${REPOSITORY}/releases"

echo "Creatting new release as version ${TAG}..."
REPLY=$(curl -H "${ACCEPT_HEADER}" -H "${TOKEN_HEADER}" -d "{\"tag_name\": \"${TAG}\", \"name\": \"${TAG}\"}" "${ENDPOINT}")

# Check error
RELEASE_ID=$(echo "${REPLY}" | jq .id)
if [ "${RELEASE_ID}" = "null" ]; then
  echo "Failed to create release. Please check your configuration. Github replies:"
  echo "${REPLY}"
  exit 1
fi

echo "Github release created as ID: ${RELEASE_ID}"
RELEASE_URL="https://uploads.github.com/repos/${REPOSITORY}/releases/${RELEASE_ID}/assets"

# If assets is empty, skip it.
if [ "${ASSETS_DIR}" = "" ]; then
  echo "No upload assets, finished."
  exit
fi

# Uploads artifacts
for FILE in ${ASSETS_DIR}/*; do
  MIME=$(file -b --mime-type "${FILE}")
  echo "Uploading assets ${FILE} as ${MIME}..."
  NAME=$(basename "${FILE}")
  curl -v \
    -H "${ACCEPT_HEADER}" \
    -H "${TOKEN_HEADER}" \
    -H "Content-Type: ${MIME}" \
    --data-binary "@${FILE}" \
    "${RELEASE_URL}?name=${NAME}"
done

echo "Finished."
