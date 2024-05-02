#!/bin/bash

if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    echo script was run directly   
    exit
fi

dir=$(pwd)

# glossary.lua needs it
# export HEALTHY_LIST

# echo healthy is "${HEALTHY_LIST[@]}"
healthy_string="${HEALTHY_LIST[*]}"

export healthy_string

cd "$ASSET_WORKDIR/../resources/glossary/" || exit
latexmk -pdf -g -f -norc -pdflatex="lualatex -interaction=batchmode" -outdir="$ASSET_WORKDIR/../assets/glossary" "glossary.tex"

cd "$dir"
