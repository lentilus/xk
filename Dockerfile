FROM mfisherman/texlive-full

COPY . /xk
WORKDIR /xk

RUN /xk/install
RUN echo "ZETTEL_DATA=/data">/xk/config

ENTRYPOINT ["/usr/local/bin/xk"]
