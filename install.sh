#!/usr/bin/env bash
#
# TODO(everyone): Keep this script simple and easily auditable.
#
#   /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/takumakei/androidpublisher-get-version-code/refs/heads/main/install.sh)"
#
set -euo pipefail

       OWNER=takumakei
       REPOS=androidpublisher-get-version-code
ASSET_PREFIX=androidpublisher-get-version-code-
  LOCAL_NAME=androidpublisher-get-version-code

main() {
  local tag="${1:-latest}" suffix url dir

  case "$(uname -sm)" in
    "Darwin x86_64") suffix="darwin-amd64" ;;
    "Darwin arm64" ) suffix="darwin-arm64" ;;
    "Linux x86_64" ) suffix="linux-amd64"  ;;
    "Linux aarch64") suffix="linux-arm64"  ;;
    *)
      echo "Error: the pre-built binary for $(uname -sm) are not available." 1>&2
      exit 1
      ;;
  esac

  [[ "$tag" != "latest" ]] && tag="tags/$tag"

  url="$(
    curl -fsSL "https://api.github.com/repos/$OWNER/$REPOS/releases/$tag" \
      | jq -r '.assets[]|select(.name=="'"$ASSET_PREFIX$suffix"'")|.browser_download_url'
  )"
  if [[ -z "$url" ]]; then
    echo "Error: not found $tag"
    exit 1
  fi

  dir="$(install_dir)"
  curl --compressed -fL -o "$dir/$LOCAL_NAME" "$url"
  chmod +x "$dir/$LOCAL_NAME"
  echo "$dir/$LOCAL_NAME"
}

install_dir() {
  local dirs=( "$HOME/.local/bin" "$HOME/bin" )
  for i in "${dirs[@]}"; do
    if [[ -d "$i" ]]; then
      echo "$i"
      return
    fi
  done
  mkdir -p "${dirs[0]}"
  echo "${dirs[0]}"
}

main "$@"
