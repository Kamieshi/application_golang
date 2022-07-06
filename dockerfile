FROM golang:1.16-alpine

WORKDIR /app

COPY ./ ./

RUN go build Server.go

EXPOSE 8005

CMD [ "./Server" ]