FROM golang

## NOTE:GOFLAGS won't be needed in go1.12
ARG goflags="-mod=vendor"
ENV GOFLAGS=$goflags

RUN mkdir /gophercli

ADD . /gophercli/
WORKDIR /gophercli

RUN go test -v ./...
RUN GOOS=windows go build cmd/gophercli.go
RUN GOOS=linux go build cmd/gophercli.go
RUN GOOS=darwin go build -o gophercli.darwin cmd/gophercli.go
