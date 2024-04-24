#!/bin/bash

echo exporting anki cards

deck="$(basename "$ZETTELKASTEN_DATA")"
echo $deck

for zettel in "${HEALTHY_NAMES[@]}"; do
    path="$( xettelkasten path -z "$zettel")/zettel.tex"
    json="$(texanki-parse "$(cat "$path")")"
    
    num_cards="$(jq ".[] .source"<<< "$json")"

    counter=0
    for _ in $num_cards; do
        ques="$(jq -r ".[$counter] .question"<<< "$json")"  
        ans="$(jq -r ".[$counter] .source"<<< "$json")"  

        if [[ -z $ques ]]; then
            ques="${zettel/_/ }"
        fi

        prepreamble='\documentclass[12pt]{article}
\usepackage[paperwidth=5in, paperheight=100in]{geometry}
\pagestyle{empty}'

        begdoc='\begin{document}'
        enddoc='\end{document}'
        
        preamble="$(cat "$ZETTEL_DATA/standalone_preamble.tex")"
        # front="[latex]"$'\n'"$prepreamble"$'\n'"$preamble"$'\n'"$begdoc"$'\n'"Hello World"$'\n'"$enddoc"$'\n'"[/latex]"
        # front="[latex]"$'\n'"$prepreamble"$'\n'"$begdoc"$'\n'"Hello World"$'\n'"$enddoc"$'\n'"[/latex]"
        # back="[latex] $prepreamble $preamble $begdoc $ans $enddoc [/latex]"
        #
        #
        #
        front='[latex]
        \documentclass[12pt]{article}
        \usepackage{geometry}
        \pagestyle{empty}

        % basics
        \usepackage{inputenc}
        \usepackage{babel}
        \usepackage{amsmath,amssymb,amsfonts,amsthm}
        \usepackage{faktor}
        \usepackage{thmtools}
        \usepackage{mathtools} % prettier math
        \usepackage{mathdots}
        \usepackage{enumitem}
        \usepackage{comment}

        % frames
        \usepackage{mdframed}
        \usepackage{framed}
        \usepackage{tcolorbox}
        \definecolor{shadecolor}{rgb}{0.9,0.9,0.9}

        % graphics
        \usepackage{import}
        \usepackage{xifthen}
        \usepackage{pdfpages}

        % amsthm config
        \declaretheoremstyle[notebraces={[}{]},headpunct=\newline,]{custom}
        \theoremstyle{custom}
        \newtheorem*{theorem}{Theorem}
        \newtheorem*{lemma}{Lemma}
        \newtheorem*{corollary}{Corollary}

        \theoremstyle{custom}
        \newtheorem*{axiom}{Axiom}
        \newtheorem*{definition}{Definition}
        \newtheorem*{example}{Example}
        \newtheorem*{remark}{Remark}

        % symbol shortcuts
        \newcommand{\zz}{\mathrm{Z\kern-.4em\raise-0.5ex\hbox{Z}}}
        \newcommand{\tx}[1]{\text{ #1 }}
        \newcommand{\from}{\colon}
        \newcommand{\La}{\mathcal{L}}
        \newcommand{\N}{\mathbb{N}}
        \newcommand{\K}{\mathbb{K}}
        \newcommand{\R}{\mathbb{R}}
        \newcommand{\Q}{\mathbb{Q}}
        \newcommand{\C}{\mathbb{C}}

        % math operators
        \DeclareMathOperator{\Mat}{Mat}
        \DeclareMathOperator{\sgn}{sgn}
        \DeclareMathOperator{\Eig}{Eig}
        \DeclareMathOperator{\Image}{Im}
        \DeclareMathOperator{\Hom}{Hom}
        \DeclareMathOperator{\End}{End}
        \DeclareMathOperator{\GL}{GL}
        \DeclareMathOperator{\HP}{HP}
        \DeclareMathOperator{\rand}{rand}
        \DeclareMathOperator{\ord}{ord}

        % misc
        \setlength\parindent{0pt}

        % quantifiers
        \let\oldforall\forall
        \let\oldexists\exists
        \renewcommand{\forall}{\ \oldforall}
        \renewcommand{\exists}{\ \oldexists}

        \makeatletter
        % nice way to write sets
        \newcommand{\set}[1]{\@ifnextchar\bgroup {\left\{#1\filteredset} {\{#1\}} }
        \newcommand{\filteredset}[1]{\ \setseperator \ #1 \right\}}
        \newcommand{\setseperator}{\middle|}
        \makeatother
        \newcommand{\norm}[1]{\left|\left|#1\right|\right|}

            \newenvironment{flashcard}{}{}
        \newenvironment{question}{\paragraph{Question}}{\vspace{5pt}}
        \begin{document}
        Hello World
        \end{document}
        [/latex]'
        back="$front"

        echo "$front"
        echo "$back"

        apy --base-path "/home/lentilus/.var/app/net.ankiweb.Anki/data/Anki2/" add-single -m zettelkasten -d testing "$front" "$back"
        #
        counter=$(( "$counter" + 1 ))
    done
done

