#!/usr/bin/env bash
set -euo pipefail

output_dir=${1:-./publish}
input_dir=$(mktemp -d)
trap 'rm -rf "$input_dir"' EXIT

mkdir -p "$output_dir"

download() {
  curl -fsSL --connect-timeout 5 --retry 5 "$1" -o "$2"
}

write_domain_rules() {
  local name=$1
  local source=$2
  local value rule count=0
  local yaml_file="$output_dir/$name.txt"
  local rules_file="$input_dir/$name.rules"

  printf 'payload:\n' > "$yaml_file"
  : > "$rules_file"

  while IFS= read -r value || [[ -n "$value" ]]; do
    value=${value%$'\r'}
    value=${value%%#*}
    value=${value#"${value%%[![:space:]]*}"}
    value=${value%"${value##*[![:space:]]}"}
    value=${value%%[[:space:]]*}
    [[ -z "$value" ]] && continue

    case "$value" in
      full:*)
        rule=${value#full:}
        ;;
      domain:*)
        rule="+.${value#domain:}"
        ;;
      keyword:*|regexp:*|include:*)
        continue
        ;;
      +.*)
        rule=$value
        ;;
      *)
        rule="+.$value"
        ;;
    esac

    [[ -z "$rule" ]] && continue
    printf '%s\n' "$rule" >> "$rules_file"
    printf "  - '%s'\n" "$rule" >> "$yaml_file"
    count=$((count + 1))
  done < "$source"

  if [[ "$count" -eq 0 ]]; then
    echo "::error::$name has no supported domain rules."
    return 1
  fi
}

write_ipcidr_rules() {
  local name=$1
  local source=$2
  local value count=0
  local yaml_file="$output_dir/$name.txt"
  local rules_file="$input_dir/$name.rules"

  printf 'payload:\n' > "$yaml_file"
  : > "$rules_file"

  while IFS= read -r value || [[ -n "$value" ]]; do
    value=${value%$'\r'}
    value=${value%%#*}
    value=${value#"${value%%[![:space:]]*}"}
    value=${value%"${value##*[![:space:]]}"}
    [[ -z "$value" || "$value" != */* ]] && continue

    printf '%s\n' "$value" >> "$rules_file"
    printf "  - '%s'\n" "$value" >> "$yaml_file"
    count=$((count + 1))
  done < "$source"

  if [[ "$count" -eq 0 ]]; then
    echo "::error::$name has no supported CIDR rules."
    return 1
  fi
}

for name in icloud private tld-!cn; do
  cp "$output_dir/$name.txt" "$input_dir/$name.source"
done

download "https://raw.githubusercontent.com/Loyalsoldier/v2ray-rules-dat/release/direct-list.txt" "$input_dir/direct.source"
download "https://raw.githubusercontent.com/Loyalsoldier/v2ray-rules-dat/release/proxy-list.txt" "$input_dir/proxy.source"
download "https://raw.githubusercontent.com/Loyalsoldier/v2ray-rules-dat/release/reject-list.txt" "$input_dir/reject.source"
download "https://raw.githubusercontent.com/Loyalsoldier/v2ray-rules-dat/release/gfw.txt" "$input_dir/gfw.source"
download "https://raw.githubusercontent.com/Loyalsoldier/v2ray-rules-dat/release/greatfire.txt" "$input_dir/greatfire.source"
download "https://raw.githubusercontent.com/felixonmars/dnsmasq-china-list/master/apple.china.conf" "$input_dir/apple.conf"
download "https://raw.githubusercontent.com/felixonmars/dnsmasq-china-list/master/google.china.conf" "$input_dir/google.conf"
download "https://raw.githubusercontent.com/Loyalsoldier/geoip/release/text/cn.txt" "$input_dir/cncidr.source"
download "https://raw.githubusercontent.com/Loyalsoldier/geoip/release/text/private.txt" "$input_dir/lancidr.source"
download "https://raw.githubusercontent.com/Loyalsoldier/geoip/release/text/telegram.txt" "$input_dir/telegramcidr.source"
download "https://raw.githubusercontent.com/Loyalsoldier/clash-rules/release/applications.txt" "$output_dir/applications.txt"

awk -F/ '/^server=\// { print $2 }' "$input_dir/apple.conf" > "$input_dir/apple.source"
awk -F/ '/^server=\// { print $2 }' "$input_dir/google.conf" > "$input_dir/google.source"

for name in apple direct gfw google greatfire icloud private proxy reject tld-not-cn; do
  source_name=$name
  [[ "$name" == tld-not-cn ]] && source_name=tld-!cn
  write_domain_rules "$name" "$input_dir/$source_name.source"
done

for name in cncidr lancidr telegramcidr; do
  write_ipcidr_rules "$name" "$input_dir/$name.source"
done

go run ./ --mrsrulespath="$input_dir" --outputpath="$output_dir"
