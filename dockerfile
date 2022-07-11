FROM golang:1.16-alpine

WORKDIR /app

COPY ./ ./

RUN go build main.go

EXPOSE 8005

CMD [ "./main" ]