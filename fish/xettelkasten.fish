# add xettelkasten to PATH
set -l bin "$(dirname (status -f))/../src"
set -l scripts "$(dirname (status -f))/../scripts"
set PATH $PATH $bin
set PATH $PATH $scripts
