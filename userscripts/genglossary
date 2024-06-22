#!/bin/bash

userscripts="$(dirname "$0")"
latexmk --shell-escape -pdf -g -f -norc -pdflatex="lualatex -interaction=batchmode" -outdir="/tmp/glossary" -cd "$userscripts/glossary/glossary.tex"
