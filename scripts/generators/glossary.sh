#!/bin/bash

if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    echo script was run directly   
    exit
fi

dir=$(pwd)
cd "$ASSET_WORKDIR/../resources/glossary/" || exit
latexmk -pdf -f -norc -pdflatex="lualatex -interaction=batchmode" -outdir="$ASSET_WORKDIR/../assets/glossary" "glossary.tex"

cd $dir
