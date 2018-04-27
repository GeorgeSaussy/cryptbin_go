FROM golang:latest
RUN mkdir /app

RUN go get github.com/lib/pq
RUN go get github.com/mattn/go-sqlite3

ADD . /app/
WORKDIR /app
RUN go build -o cryptbin .
CMD ["/app/cryptbin"]
