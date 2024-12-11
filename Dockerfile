FROM golang:alpine

ADD . .

RUN go build  -o main ./

CMD ["./main"]
EXPOSE 8000