FROM golang:1.19

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY . .
RUN go mod download && go mod verify

RUN go build -v -o /usr/local/bin/app ./...
COPY templates/index.html .
COPY static/* /static/
EXPOSE 5000
CMD ["app"]
