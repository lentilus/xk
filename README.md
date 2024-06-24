# xk - a LaTeX centric zettelkasten
> an extensible light weight LaTeX centric zettelkasten

## The Idea

Generally speaking - a [zettelkasten (slipbox)](https://en.wikipedia.org/wiki/Zettelkasten) is strategy for organizing scrientific or literary work. Information is condenced into atomic concepts. Each zettel is the manifestation of a concepts. Zettels with related concepts should link to another.

> This strategy is very well suited for digial work. Each zettel is a file. A link to a zettel is a hyperref to the zettels file.
For this reason there are tons of apps (Obsidian, Zettlr, Logseq, etc) implementing this paradigm. The issue for me is, that most of them are markdown not LaTeX based, which is a dealbreaker for me. I want to use the full power of LaTeX in my notes. MathJax etc dont cut it for me :D.

xk is light weight, extensible, and could be used with any plain text file format. (It just happes to be LaTeX for me - hence the name XeTtelkasten : TeX-Zettelkasen)


## Demo (with xettelkasten.nvim)

https://github.com/lentilus/xk/assets/170900031/8bedf9c5-04ba-4ffa-b534-0004cb29456f


This video shows me creating two zettels and referencing one to anoter.

## Installation

Make sure you have the basic gnu utils such as xargs etc available (you probably have).
Just clone the repo and run the installation script.

```bash
git clone https://github.com/lentilus/xk.git
cd xettelkasten-core
./install
```
> The install script just creates a symlink at `$HOME/.local/bin/xk`.

Note: You may need additional dependencies for running the scripts in `./scripts`. These are not part of the core functionality, but provide additional features.


## Usage
> make sure you have done the installation steps
To get started, edit the configuration file found at `~/.config/xettelkasten/config`.
Most importantly specify a directory where your zettelkasten should live.
To get going run
```bash
xk init
```
now you have the following commands at your disposal:

Basics
```bash
xk git init               # runs my git command on the zettelkasten
xk insert -z "foo"        # inserts a zettel with the name foo
xk ls                     # list all zettels
xk mv -z "foo" -n "bar"   # rename zettel foo to bar (updates references too)
xk path -z "bar"          # get path of zettel "bar"
xk rm -z "bar"            # remove "bar"
```

References
```bash
xk ref insert -z "foo" -r "bar" # let foo reference bar
xk ref ls -z "foo"              # list references from foo
xk ref rm -z "foo" -r "bar"     # remove reference to bar
```

Tags
```bash
xk tag insert -z "foo" -r "bar" # add tag bar to "foo"
xk tag ls -z "foo"              # list tags of "foo"
xk tag rm -z "foo" -r "bar"     # remove tag "bar" from "foo"
```

If you are a neovim user I recommend the plugin `xettelkasten.nvim`, coming to Github soon but currently hosetet at gitlab.com/lentilus/xettelkasten.nvim.git.

## Docker
The Dockerfile builds an image based on a texlive-full (alpine) image to enable in container compilation.

You can build the docker container yourself and use it to run xk commands on your local zettelkasten.

```bash
docker build . --progress=plain -t lentilus/xkfoo
docker run -v /my/zettel/kasten:/data lentilus/xk ls
```
> Note that the image will not contain all dependencies for all commands / userscripts

The image will soon be available on dockerhub.

## Userscripts
> coming soon

## TODOS and ideas for Contribution

[x] make xk run in alpine contianers
[x] rewrite user-scripts
[ ] complete readme
[ ] cronjob for asset generation
[ ] write tests
[ ] man entry
[ ] help menu
