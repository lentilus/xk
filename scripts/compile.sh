latexmk -f -norc -pdflatex="lualatex -interaction=batchmode" -auxdir=./out -outdir=./out -cd "../resources/glossary/" -pdf glossary.tex
