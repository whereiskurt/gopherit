FROM golang

ARG releaseVersion="v1.0.0"
ENV VERSION=$releaseVersion

ARG hash="0xABCD1234"
ENV HASH=$hash

ARG goos="linux"
ENV GOOS=$goos

## NOTE: GOFLAGS won't be needed in go1.12
ARG goflags="-mod=vendor"
ENV GOFLAGS=$goflags

RUN mkdir /gophercli

ADD . /gophercli/

WORKDIR /gophercli

RUN go test -v ./...

RUN GOOS=$GOOS go build \
    --ldflags \
    "-X 00-newapp-template/internal/app/cmd.ReleaseVersion=$VERSION \
     -X 00-newapp-template/internal/app/cmd.GitHash=$HASH" \
    -o ./gophercli \
    cmd/gophercli.go