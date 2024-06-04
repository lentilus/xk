# xk : XeTtelkasten - a LaTeX centric zettelkasten
> an extensible light weight LaTeX centric zettelkasten

## The Idea

Generally speaking - a [zettelkasten (slipbox)](https://en.wikipedia.org/wiki/Zettelkasten) is strategy for organizing scrientific or literary work. Information is condenced into atomic concepts. Each zettel is the manifestation of a concepts. Zettels with related concepts should link to another.

> This strategy is very well suited for digial work. Each zettel is a file. A link to a zettel is a hyperref to the zettels file.
For this reason there are tons of apps (Obsidian, Zettlr, Logseq, etc) implementing this paradigm. The issue for me is, that most of them are markdown not LaTeX based, which is a dealbreaker for me. I want to use the full power of LaTeX in my notes. MathJax etc dont cut it for me :D.

xk is light weight, extensible, and could be used with any plain text file format. (It just happes to be LaTeX for me - hence the name XeTtelkasten : TeX-Zettelkasen)

## Installation

Make sure you have the basic gnu utils such as xargs etc available (you probably have).
Just clone the repo and run the installation script.

```bash
git clone https://github.com/lentilus/xettelkasten-core.git
cd xettelkasten-core
./install
```
> The install script just creates a symlink at `$HOME/.local/bin/xk`.

Note: You may need additional dependencies for running the scripts in `./scripts`. These are not part of the core functionality, but provide additional features.

### Docker
Pull the docker from dockerhub. It expects your zettelkasten to be mounted to its /root/zettelkasten.
Use compose or do directly with docker.

## Usage
> cli usage coming soon

If you are a neovim user I recommend the plugin `xettelkasten.nvim`, coming to Github soon but currently hosetet at gitlab.com/lentilus/xettelkasten.nvim.git.

## additional scripts
> coming soon

## TODOS and ideas for Contribution

- man entry
- help menu
- make xk run in alpine contianers
- complete readme
- write tests
- complete asset scripts
- cronjob for asset generation
