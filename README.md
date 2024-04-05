# XeTtelkasten - A TeX based Zettelkasen

## Intro

Extensible manager for atomic notes in Tex-Format written in bash. 

Leverage Tmux and nvim for tight workflow integration

No extra nvim configuratio needed, as XeTtelkasen leverages nvim --remote to control an nvim instance of your choice, specified by NVIM_CMD

## Get started

```bash
git clone ...
cd xettelkasten
cp config.template config
```

Edit the config, most importantly specify the ZETTEL_DATA directory.

Note, that variables and ~ in the config are not expanded

```bash
./src/zettelkasten
    #    XETTELKASTEN CLI
    #    ------------------------------------
    #     Version: ALPHA 0.1
    #     Usage: zettelkasten [command]
    # 
    #    Commands:
    #    ------------------------------------
    #     init      initialize xettelkasten
    #     ref       manage references
    #     go        navigate
    #     zettel    create/remove zettel
    #     status    output open zettel
    #     *         Help

./src/zettelkasten init

./src/zettelkasten zettel
    #    zettel
    #    Usage:  zettel [command]
    #    Commands:
    #    new       create a new zettel
    #    del       delete an existing zettel
    #    *         Help

./src/zettelkasten zettel new -n "my_first_zettel"
```

You can also use the tool interactively using tmux run-shell

```bash
tmux run-shell -b "/path/to/src/zettelkasten zettel new"
tmux run-shell -b "/path/to/src/zettelkasten go fzf"
tmux run-shell -b "/path/to/src/zettelkasten ref add"
tmux run-shell -b "/path/to/src/zettelkasten go ref"
# ...
```
## Git integration

with the following syntax you can use your zettelkasten in combination with git.
```bash
zettelkasten git init
zettelkasten git foo
```

## Gitlab CI - Auto-Glossary

```bash
zettelkasten git remote add origin foo
zettelkasten publish
```
Note: publish is essentially like a combination of
```bash
zettelkasten git add -A
zettelkasten git commit -m "foo"
zettelkasten git push
```
If you commit a zettelkasten to gitlab (which is encuraged), a glossary of all zettel is created automatically. This is done using Gitlab CI with a predefined configuration that is part of the base directory.

To see an example for a zettelkasten on Gitlab along with an auto-generated glossary check out my "Zettelkasten of Mathematics"
- repo: https://gitlab.com/lentilus/zettelkasten-of-mathematics/ 
- most recent glossary: https://gitlab.com/lentilus/zettelkasten-of-mathematics/-/releases/permalink/latest/downloads/export/glossary.pdf

## Anki Export

> coming soon!
