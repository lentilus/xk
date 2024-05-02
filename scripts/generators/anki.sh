#!/bin/bash
#
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    echo script was run directly   
    exit
fi

store_json() { path="$1"
    filename="$2"
    json="{\"action\":\"storeMediaFile\",\"params\":{\"filename\":\"$filename\",\"path\":\"$path\"}}"
    echo "$json"
}


import_card() {
    front_file="$1"
    back_file="$2"

    front_id="flashtex$(sha1sum "$front_file" |  awk '{print $1}').svg"
    back_id="flashtex$(sha1sum "$back_file" |  awk '{print $1}').svg"

    save_front="$(store_json "$front_file" "$front_id")"
    save_back="$(store_json "$back_file" "$back_id")"

    addcard=" { \"action\": \"addNote\",
                \"params\": {
                    \"note\": {
                        \"deckName\": \"Default\",
                        \"modelName\": \"Basic\",
                        \"fields\": {
                            \"Front\": \"<img src=$front_id>\",
                            \"Back\": \"<img src=$back_id>\" },
                            \"tags\": []
                        }
                    }
                }"

    requestjson=" { \"action\": \"multi\",
                    \"params\": {
                        \"actions\": [
                            $save_front,
                            $save_back,
                            $addcard
                        ]
                    },
                    \"version\": 6}"

    response="$(curl localhost:8765 -X POST -d "$requestjson")"
    error="$(jq ".result[2].error" <<< "$response")"

    if [[ $error == "null" ]]; then
        echo "$response", error "$error"
        echo added new card
        return
    fi

    findjson="
    {
        \"action\": \"findNotes\",
        \"version\": 6,
        \"params\": {
            \"query\": \"$front_id or $back_id\"
        }
    }"
    findresponse="$(curl localhost:8765 -X POST -d "$findjson")"
    echo find response: "$findresponse"
    id="$(jq ".result[0]" <<< "$findresponse")"
    echo id is $id

    updatejson="
    {
        \"action\": \"updateNote\",
        \"version\": 6,
        \"params\": {
            \"note\": {
                \"id\": $id,
                \"fields\": {
                    \"Front\": \"<img src=$front_id>\",
                    \"Back\": \"<img src=$back_id>\"
                }
            }
        }
    }"
    echo "echo updating request: $updatejson"
    curl localhost:8765 -X POST -d "$updatejson"
}

echo exporting anki cards...

# counter=0

for zettel in "${HEALTHY_LIST[@]}"; do
    tex_file="$zettel/zettel.tex"
    tex_source="$(cat "$tex_file")"
    name="$(basename "$zettel")"
    name="${name//_/ }"
    echo "$name"
    json="$(flashtexparse "$tex_source" "$name")"
    num_cards="$(jq "length"<<< "$json")"

    # [[ $counter -gt 10 ]] && break
    # counter=$(( "$counter" + 1 ))

    for i in $(seq 0 $(( "$num_cards" -1 ))); do


        temp_dir="/tmp/converter_$(uuidgen)"
        mkdir -p "$temp_dir"
        dir="$(dirname "$tex_file")"

        jq -r ".[$i] .question"<<< "$json" > "$dir/front.tex"
        jq -r ".[$i] .source"<<< "$json" > "$dir/back.tex"

        latexmk -pdf -f -norc -lualatex -interaction=batchmode -outdir="$temp_dir" -cd "$dir/front.tex" &>/dev/null
        latexmk -pdf -f -norc -lualatex -interaction=batchmode -outdir="$temp_dir" -cd "$dir/back.tex" &>/dev/null

        rm -rf "$dir/front.tex"
        rm -rf "$dir/back.tex"

        pdf2svg "$temp_dir/front.pdf" "$temp_dir/front.svg"
        pdf2svg "$temp_dir/back.pdf" "$temp_dir/back.svg"


        inkscape --actions "select-all;fit-canvas-to-selection" --export-overwrite "$temp_dir/front.svg"
        inkscape --actions "select-all;fit-canvas-to-selection" --export-overwrite "$temp_dir/back.svg"

        prepared_cards+=("$temp_dir")
    done
done

flatpak  run net.ankiweb.Anki &>/dev/null &
sleep 3

for tmp in "${prepared_cards[@]}"; do
    import_card "$tmp/front.svg" "$tmp/back.svg"
    echo ""
    rm -rf "$tmp"
done

flatpak kill net.ankiweb.Anki
