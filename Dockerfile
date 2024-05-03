FROM debian

RUN apt-get update && apt-get install -y git texlive-full latexmk python3-pip

# anki stuff
# we need --break-system-packages in order to use normal pip
# RUN pip install git+https://gitlab.com/lentilus/texflash.git --break-system-packages

COPY . /xettelkasten

# configuration
RUN echo 'ZETTEL_DATA="/root/.local/zettelkasten"'>/xettelkasten/config &&\
    echo 'LOG_FILE="/tmp/xettelkasten.log"'>>/xettelkasten/config

# add to path
ENV PATH "$PATH:/xettelkasten/src/xettelkasten"
ENV PATH "$PATH:/xettelkasten/srcipts/generate_assets"

RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
