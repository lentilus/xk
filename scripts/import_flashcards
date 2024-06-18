#!/bin/bash

if ! command -v xk &> /dev/null
then
    echo "xettelkasten-core could not be found."
    echo "the asset generation depends on it."
    echo "please run the installer first."
    exit 1
fi

if ! command -v texgoanki &> /dev/null
then
    echo "texgoanki could not be found. Exiting"
    exit 1
fi

kasten="$(xk path)"
deck="$(basename "$kasten")"

for z in $(xk ls); do
    echo "exporting flashcards from $z"
    path="$(xk path -z "$z")"
    texgoanki "$deck" "$(cat "$path/zettel.tex")" "$z/flashcard.tex" "$kasten" "standalone_preamble.tex" "$z"
done
