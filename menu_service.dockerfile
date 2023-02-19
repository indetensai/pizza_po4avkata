FROM golang:latest
WORKDIR /usr/src/app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY menu_service/*.go ./
COPY menu_service/.env ./
COPY service/ ./service/
RUN go build -o ./menu
EXPOSE 50001
CMD [ "./menu" ]