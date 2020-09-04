FROM golang:1.13-alpine3.11

RUN apk add --update --no-cache \
    bash \
    bash-completion \
    binutils \
    ca-certificates \
    coreutils \
    curl \
    findutils \
    g++ \
    git \
    grep \
    less \
    make \
    openssl \
    shellcheck \
    util-linux

WORKDIR /calc

deps:
    RUN go get golang.org/x/tools/cmd/goimports
    RUN go get golang.org/x/lint/golint
    RUN go get github.com/gordonklaus/ineffassign
    COPY go.mod go.sum ./
    RUN go mod download
    SAVE ARTIFACT go.mod AS LOCAL go.mod
    SAVE ARTIFACT go.sum AS LOCAL go.sum
    SAVE IMAGE

code:
    FROM +deps
    COPY --dir cmd ./
    COPY ./parser+parser/*.go ./parser/
    SAVE IMAGE

lint:
    FROM +code
    RUN output="$(ineffassign . | grep -v '/calc/parser/.*\.go')" ; \
        if [ -n "$output" ]; then \
            echo "$output" ; \
            exit 1 ; \
        fi
    RUN output="$(goimports -d . 2>&1)" ; \
        if [ -n "$output" ]; then \
            echo "$output" ; \
            exit 1 ; \
        fi
    RUN golint -set_exit_status ./...
    RUN output="$(go vet ./... 2>&1)" ; \
        if [ -n "$output" ]; then \
            echo "$output" ; \
            exit 1 ; \
        fi

unit-test:
    FROM +code
    RUN go test ./...

calc:
    FROM +code
    ARG GOOS=linux
    ARG GOARCH=amd64
    RUN test -n "$GOOS" && test -n "$GOARCH"
    ARG GOCACHE=/go-cache
    RUN mkdir -p build
    RUN --mount=type=cache,target=$GOCACHE \
        go build \
            -o build/calc \
            cmd/calc/*.go
    SAVE ARTIFACT ./build/tags
    SAVE ARTIFACT ./build/ldflags
    SAVE ARTIFACT build/calc AS LOCAL "build/$GOOS/$GOARCH/calc"

