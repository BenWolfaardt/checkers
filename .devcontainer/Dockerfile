FROM --platform=linux ubuntu:22.04
ARG BUILDARCH

# Change your versions here
ENV GO_VERSION=1.18.3
ENV IGNITE_VERSION=0.22.1
ENV NODE_VERSION=18.x

ENV LOCAL=/usr/local
ENV GOROOT=$LOCAL/go
ENV HOME=/root
ENV GOPATH=$HOME/go
ENV PATH=$GOROOT/bin:$GOPATH/bin:$PATH

RUN mkdir -p $GOPATH/bin

ENV PACKAGES curl gcc jq openssh-client
RUN apt-get update
RUN apt-get install -y $PACKAGES

# Install Go
RUN curl -L https://go.dev/dl/go${GO_VERSION}.linux-$BUILDARCH.tar.gz | tar -C $LOCAL -xzf -

RUN go install github.com/go-delve/delve/cmd/dlv@v1.20.2 \
    && go install golang.org/x/tools/gopls@v0.14.2 \
    && go install github.com/haya14busa/goplay/cmd/goplay@v1.0.0 \
    && go install github.com/josharian/impl@latest \
    && go install github.com/fatih/gomodifytags@latest \
    && go install honnef.co/go/tools/cmd/staticcheck@v0.3.3

# Install ignite
RUN curl -L https://get.ignite.com/cli@v${IGNITE_VERSION}! | bash

# Install Node
RUN curl -fsSL https://deb.nodesource.com/setup_${NODE_VERSION} | bash -
RUN apt-get install -y nodejs

RUN apt-get clean

EXPOSE 1317 3000 4500 5000 26657

WORKDIR /checkers

COPY go.mod /checkers/go.mod
RUN go mod download
RUN rm /checkers/go.mod
