FROM golang:latest

WORKDIR /go/github.com/sede-res/renamosy
COPY . .
# RUN go build

# CMD ["./renamosy"]
CMD ["go", "run", "main.go"]
