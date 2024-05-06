FROM debian

RUN apt-get update && apt-get install -y git texlive-full latexmk python3-pip
COPY . /xettelkasten

# configuration file
RUN echo 'ZETTEL_DATA="/root/zettelkasten"' >/xettelkasten/config &&\
    echo 'LOG_FILE="/tmp/xettelkasten.log"'>>/xettelkasten/config &&\
    echo 'GEN_ANKI=0'                      >>/xettelkasten/config &&\
    echo 'GEN_GLOSSARY=1'                  >>/xettelkasten/config

# source scripts
ENV PATH "/xettelkasten/src:${PATH}"
ENV PATH "/xettelkasten/scripts:${PATH}"

RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
