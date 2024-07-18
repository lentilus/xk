# This docker file is not very robust
# as it pull in external stuff, that could
# change in the future - just be aware of that.

FROM mfisherman/texlive-full

# install xk
COPY . /xk
WORKDIR /xk

# some dependencies
RUN apk add git curl

RUN /xk/install && \
    echo "ZETTEL_DATA=/data">/xk/config

# dependencies for flashcard script
RUN apk add py3-pip && \
    pip install --break-system-packages git+https://github.com/lentilus/texflash.git

COPY --from=golang:1.21-alpine /usr/local/go/ /usr/local/go/
ENV PATH="/usr/local/go/bin:${PATH}"

RUN git clone https://github.com/lentilus/texgoanki.git /texgoanki
RUN cd /texgoanki && go build -o /bin/texgoanki main.go

RUN apk add inkscape
RUN apk add pdf2svg --repository=http://dl-cdn.alpinelinux.org/alpine/edge/testing/

# xk as entry
ENTRYPOINT ["/usr/local/bin/xk"]
