FROM golang:1.22.2
WORKDIR /build
COPY / .

ENTRYPOINT ["go", "run", "cmd/app.go"]
CMD $1
