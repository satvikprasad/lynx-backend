FROM golang:1.21

WORKDIR /usr/src/app

copy . . 

RUN go mod download && go mod verify

RUN go build -o /usr/local/bin/lynx-backend

CMD ["lynx-backend"]
