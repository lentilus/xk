FROM debian
RUN apt-get update && apt-get install -y git texlive-full latexmk

COPY . ./xettelkasten

RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
