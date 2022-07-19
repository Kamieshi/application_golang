FROM golang:1.18-alpine

WORKDIR /app

COPY ./ ./

RUN go build main.go

EXPOSE 8005
EXPOSE 5300

CMD [ "./main" ]

