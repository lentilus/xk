# XeTtelkasten : Backend for LaTeX based Zettelkasten

## The Idea

Generally speaking - a [Zettelkasten (slipbox)](https://en.wikipedia.org/wiki/Zettelkasten) is strategy for organizing scrientific or literary work. Information is condenced into atomic concepts. Each zettel is the manifestation of a concepts. Zettels with related concepts should link to another.

This strategy is very well suited for digial work. Each zettel is a file. A link to a zettel is a hyperref to the zettels file.
For this reason there are tons of apps (Obsidian, Zettlr, Logseq, etc) implementing this paradigm. The issue for me is, that most of them are markdown or plain text based, which is a dealbreaker for me, because I want to use the full power of LaTeX (not just some speced-down markdown-version).

The solution I came up with is very light weight, extensible, and could be used with any plain text file format. (It just happes to be LaTeX for me - hence XeTtelkasten : TeX-Zettelkasen)

## Installation

### dependencies
- bash and gnu utils such as xargs etc (probably installed already)

for asset generation (totally optional)
- latexmk
- anki
- inkscape

### locally
Just clone the repo and add src/xettelkasten to your path.

### containerized
Pull the docker from dockerhub. It expects your zettelkasten to be mounted to its /root/zettelkasten.
Use compose or do directly with docker.

## The Core Utility Commands

## My Workflow
I edit and browse my zettelkasten almost exclusively through nvim. For that I use the plugin that I wrote specifically for the xettelkasten. I use the xettelkasten to learn mathematics. I fill it with content during the lecture or if I have ideas in between. Generally I found the following quite sensible to determine what to put and what not to put on a zettel:

- A defintion deserves a zettel
- different defintions for the same thing may be on the same zettel
- Almost never should there be definitions for different things on the same zettel.
- lemmas for theorems go on the same zettel as the theorem
- Theorems with Names deserve their own zettel
- Sketch longer well known proofs instead of polluting a zettel with technical details


### auto generate Assets

- Anki flashcards
- A glossary pdf

## Neovim Plugin

## Gitlab
