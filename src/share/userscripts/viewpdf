#!/bin/bash

zettel="$(xk ls | rofi -i -async-pre-read 1 -dmenu)"

if [[ -z $zettel ]]; then
    exit
fi

path="$(xk path -z "$zettel")"

zathura "$path/zettel.pdf" || exit
