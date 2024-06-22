#!/bin/bash
# run with xk script checkhealth myzettel

[ "$SKIP_HEALTHCHECK" = 1 ] && exit 0

z=$1
out="/tmp/xk_checkhealth/$z"
mkdir -p "$out"

echo "$z"

p="$(xk path -z "$z")"
latexmk -pdf -g -norc -pdflatex="pdflatex -interaction=batchmode" -outdir="$out" -cd "$p/zettel.tex" &>/dev/null

if [ -f "$out/zettel.pdf" ]; then
    # make sure we dont affect next check
    rm -rf "$out/zettel.pdf"
    echo "$z : healthy"
    exit 0
else
    exit 1
fi
