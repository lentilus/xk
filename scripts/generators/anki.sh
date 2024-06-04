#!/bin/bash

if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    echo script was run directly   
    exit
fi

COLLECTION="current"
DECK_NAME="default"
TEMPDIR="/tmp/converter"

# a nice helper function to make building the request json
# less of a pain. :)
query() {
    action=$1
    params=$2
    echo "{
              \"action\": \"$action\",
              \"version\": 6,
              \"params\": { $params }
          }"
}

# helper function for anki connect posts
connect_request() {
    curl -s localhost:8765 -X POST -d "$1"
}

# Anki does not allow us to manually assign the id of new flashcards
# therefor we need to query the collection using the texflash id
# to retrieve the actual id assigned by anki.
get_anki_id() {
    # we could probably just use the notes id in the first place
    deck=$1
    tex_id=$2
    id_request="$(query "findCards" "\"query\": \"deck:$deck id:$tex_id\"" )"
    response="$(connect_request "$id_request")"

    card_id="$(echo "$response" | jq ".result[0]")"
    response="$(connect_request "$(query "cardsToNotes" "\"cards\": [ $card_id ]")")"
    note_id="$(echo "$response" | jq ".result[0]")"
    echo "$note_id"
}

get_hash() {
    # the json from texflash parse
    src_json=$1
    echo "$src_json" | shasum | head -c 10
}

# We check if the current hash matches the already saved hash.
# If so, there are no changes and we move on to the next card
detect_changes() {
    id=$1 # anki_id
    json=$2
    # new_hash="$(echo "$json" | shasum | head -c 10)"
    new_hash="$(get_hash "$json" )"

    # get last hash from anki
    hash_request="$(query "cardsInfo" "\"cards\":[$id]")"
    old_hash="$(connect_request "$hash_request" | jq -r ".result[0] .fields .hash .value")" 

    # if the new hash is equals the old one, there are no changes.
    [[ "$new_hash" = "$old_hash" ]] && return 1
    return 0
}

# store a file in out anki db
store_file() {
    path=$1
    filename=$2
    json="$(query "storeMediaFile"  "\"filename\":\"$filename\",\"path\":\"$path\"" )"
    connect_request "$json" &>/dev/null
}

update_card() {
    anki_id="$1"
    hash="$2"
    front_file="$3"
    back_file="$4"
    
    store_file "$front_file" "front_$hash.svg"
    store_file "$back_file" "back_$hash.svg"

    front="\"Front\": \"<img src=front_$hash.svg>\""
    back="\"Back\": \"<img src=back_$hash.svg>\""
    hash="\"hash\": \"$hash\""
    # id="\"id\": \"updated id\""
    update="$(query "updateNoteFields" "\"note\": { \"id\": $anki_id, \"fields\": {$front, $back, $hash}}")"
    # echo trying update:
    connect_request "$update" &>/dev/null
}

new_card() {
    texflash_id="$1" 
    hash="$2"
    front_file="$3"
    back_file="$4"

    store_file "$front_file" "front_$hash.svg"
    store_file "$back_file" "back_$hash.svg"

    fields="\"Front\": \"<img src=front_$hash.svg>\",
            \"Back\": \"<img src=back_$hash.svg>\",
            \"id\": \"$texflash_id\",
            \"hash\": \"$hash\""

    new_note="\"deckName\": \"$DECK_NAME\", \"modelName\": \"Basic\", \"fields\": {$fields} "
    json="$(query "addNote" "\"note\": { $new_note }" | jq)"
    connect_request "$json"
}

# MAIN part of the script
echo extracting flashcards

# laumch anki to make anki_connect available
flatpak run net.ankiweb.Anki &>/dev/null &
sleep 3

# iterate over all zettels that compiled successfully in health-check
for zettel in "${HEALTHY_LIST[@]}"; do
    name="$(basename "$zettel")"
    tex_source="$(cat "$zettel/zettel.tex")"
    json="$(flashtexparse "$tex_source")"
    num_cards="$(jq "length"<<< "$json")"

    # for each zettel we iterate over all flashcards that were parsed
    # there may be multiple flashcards in one zettel
    for i in $(seq 0 $(( "$num_cards" -1 ))); do
        # get corresponding id in anki
        id="$(jq -r ".[$i] .id"<<< "$json")"
        
        # flashcard invalid
        [ "$id" = 'null' ] && break
        anki_id="$(get_anki_id "$COLLECTION" "$id")"
        hash="$(get_hash "$json")"

        if [ "$anki_id" = 'null' ]; then
            # if there is no card with this id, we need to create it.
            echo "$name: creating new"
            srcdir="new_$(uuidgen)"
            mode="new"
        else
            # if the code for the flashcard has not changed,
            # we move on to the next one
            if ! detect_changes "$anki_id" "$json"; then
                echo "$name: nothing to do"
                break
            fi
            
            # but if it did, we update the flashcard
            echo "$name: updating"
            srcdir="$hash"
            mode="update"
        fi

        tmpdir="$TEMPDIR/$srcdir"
        zettel_dir="$zettel"
        mkdir -p "$tmpdir"

        # save the source code
        jq -r ".[$i] .front"<<< "$json" > "$zettel_dir/front.tex"
        jq -r ".[$i] .back"<<< "$json" > "$zettel_dir/back.tex"
        
        # compile both front and back
        latexmk -pdf -f -norc -lualatex -interaction=batchmode -outdir="$tmpdir" -cd "$zettel_dir/front.tex" &>/dev/null
        latexmk -pdf -f -norc -lualatex -interaction=batchmode -outdir="$tmpdir" -cd "$zettel_dir/back.tex" &>/dev/null

        rm -rf "$zettel_dir/front.tex" "$zettel_dir/back.tex"

        # convert to svg
        pdf2svg "$tmpdir/front.pdf" "$tmpdir/front.svg" || break # no pdf produced
        pdf2svg "$tmpdir/back.pdf" "$tmpdir/back.svg" || break
        
        # crop to content
        inkscape --actions "select-all;fit-canvas-to-selection" --export-overwrite "$tmpdir/front.svg"
        inkscape --actions "select-all;fit-canvas-to-selection" --export-overwrite "$tmpdir/back.svg"
        
        if [ "$mode" = "update" ]; then
            update_card "$anki_id" "$hash" "$tmpdir/front.svg" "$tmpdir/back.svg"
        elif [ "$mode" = "new" ]; then
            new_card "$id" "$hash" "$tmpdir/front.svg" "$tmpdir/back.svg"
        fi
    done
done

# clean up cached stuff
rm -rf "${TEMPDIR:?}/*"

# quit anki when we are done
flatpak kill net.ankiweb.Anki
