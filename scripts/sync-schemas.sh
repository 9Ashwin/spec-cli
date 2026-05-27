#!/usr/bin/env bash
set -euo pipefail

SOURCE_ROOT="${1:-${OPENSCHEMA_DIR:-/Users/mervyn/workspaces/github/openspec-schemas}}"
SOURCE_DIR="${SOURCE_ROOT%/}/superpowers-bridge"
TARGET_DIR="${SPEC_CLI_SCHEMA_EMBED_DIR:-embed/assets/schemas/superpowers-bridge}"

if [ ! -d "$SOURCE_DIR" ]; then
  echo "schema source not found: $SOURCE_DIR" >&2
  echo "Set OPENSCHEMA_DIR or pass the openspec-schemas repo path." >&2
  exit 2
fi

if [ ! -d "$TARGET_DIR" ]; then
	echo "embedded schema target not found: $TARGET_DIR" >&2
	exit 2
fi

locale_suffix="$(printf '%s%s' T W)"
upstream_locale="zh-$locale_suffix"
traditional_en="$(printf '\124\162\141\144\151\164\151\157\156\141\154\040\103\150\151\156\145\163\145')"
traditional_zh="$(printf '\347\271\201\351\253\224')"

sync_file() {
  local rel="$1"
  local source_file="$2"
  local target_file="$3"

  if [ "$rel" = "README.md" ]; then
    sed \
      -e "s#README\\.${upstream_locale}\\.md#README.zh.md#g" \
      -e "s#CLAUDE\\.md\\.fragment\\.${upstream_locale}\\.md#CLAUDE.md.fragment.zh.md#g" \
      -e "s#${upstream_locale}#zh#g" \
      -e "s#${traditional_en}#Simplified Chinese#g" \
      -e "s#${traditional_zh}中文#简体中文#g" \
      -e "s#${traditional_zh}中文版#简体中文版#g" \
      "$source_file" > "$target_file"
  else
    cp "$source_file" "$target_file"
  fi
}

while IFS= read -r -d '' source_file; do
	rel="${source_file#$SOURCE_DIR/}"
	case "$rel" in
		"README.${upstream_locale}.md"|"templates/adopters/CLAUDE.md.fragment.${upstream_locale}.md")
			continue
			;;
	esac

	target_file="$TARGET_DIR/$rel"
	mkdir -p "$(dirname "$target_file")"
	sync_file "$rel" "$source_file" "$target_file"
	echo "synced $rel"
done < <(find "$SOURCE_DIR" -type f -print0 | sort -z)
