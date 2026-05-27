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

drift=0
locale_suffix="$(printf '%s%s' T W)"
upstream_locale="zh-$locale_suffix"
traditional_en="$(printf '\124\162\141\144\151\164\151\157\156\141\154\040\103\150\151\156\145\163\145')"
traditional_zh="$(printf '\347\271\201\351\253\224')"

normalized_source() {
  local rel="$1"
  local source_file="$2"

  if [ "$rel" = "README.md" ]; then
    sed \
      -e "s#README\\.${upstream_locale}\\.md#README.zh.md#g" \
      -e "s#CLAUDE\\.md\\.fragment\\.${upstream_locale}\\.md#CLAUDE.md.fragment.zh.md#g" \
      -e "s#${upstream_locale}#zh#g" \
      -e "s#${traditional_en}#Simplified Chinese#g" \
      -e "s#${traditional_zh}中文#简体中文#g" \
      -e "s#${traditional_zh}中文版#简体中文版#g" \
      "$source_file"
  else
    cat "$source_file"
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
	if [ ! -f "$target_file" ]; then
    echo "missing embedded schema file: $rel"
    drift=1
		continue
	fi

	tmp="$(mktemp)"
	normalized_source "$rel" "$source_file" > "$tmp"
	if ! cmp -s "$tmp" "$target_file"; then
		echo "schema drift: $rel"
		drift=1
	fi
	rm -f "$tmp"
done < <(find "$SOURCE_DIR" -type f -print0 | sort -z)

if [ "$drift" -ne 0 ]; then
  echo "Run: make sync-schemas OPENSCHEMA_DIR=$SOURCE_ROOT"
  exit 1
fi

echo "Schema bundle in sync: $TARGET_DIR"
