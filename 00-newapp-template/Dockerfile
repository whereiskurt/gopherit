FROM golang

RUN mkdir /gopherit

ADD . /gopherit/
WORKDIR /gopherit

## NOTE:GOFLAGS won't be needed in go1.12
RUN GOFLAGS="-mod=vendor" go test -v ./...
