FROM golang:1.4.2

COPY . /go/src/github.com/hiroosak/golgtm
RUN go get github.com/hiroosak/golgtm/cmd/golgtm

EXPOSE 8000
CMD ["golgtm", "-http=:8000"]
