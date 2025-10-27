#!/usr/bin/env bash
set -euo pipefail

# create-github-release.sh
# Create a GitHub release with all template zip files
# Usage: create-github-release.sh <version>

if [[ $# -ne 1 ]]; then
  echo "Usage: $0 <version>" >&2
  exit 1
fi

VERSION="$1"

# Remove 'v' prefix from version for release title
VERSION_NO_V=${VERSION#v}

gh release create "$VERSION" \
  .genreleases/taskkit-template-copilot-sh-"$VERSION".zip \
  .genreleases/taskkit-template-copilot-ps-"$VERSION".zip \
  .genreleases/taskkit-template-claude-sh-"$VERSION".zip \
  .genreleases/taskkit-template-claude-ps-"$VERSION".zip \
  .genreleases/taskkit-template-gemini-sh-"$VERSION".zip \
  .genreleases/taskkit-template-gemini-ps-"$VERSION".zip \
  .genreleases/taskkit-template-cursor-agent-sh-"$VERSION".zip \
  .genreleases/taskkit-template-cursor-agent-ps-"$VERSION".zip \
  .genreleases/taskkit-template-opencode-sh-"$VERSION".zip \
  .genreleases/taskkit-template-opencode-ps-"$VERSION".zip \
  .genreleases/taskkit-template-qwen-sh-"$VERSION".zip \
  .genreleases/taskkit-template-qwen-ps-"$VERSION".zip \
  .genreleases/taskkit-template-codex-sh-"$VERSION".zip \
  .genreleases/taskkit-template-codex-ps-"$VERSION".zip \
  .genreleases/taskkit-template-codebuddy-sh-"$VERSION".zip \
  .genreleases/taskkit-template-codebuddy-ps-"$VERSION".zip \
  --title "Task Kit Templates - $VERSION_NO_V" \
  --notes-file release_notes.md
