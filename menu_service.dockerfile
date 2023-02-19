FROM golang:alpine
WORKDIR /usr/src/app
RUN apk update && apk add protoc \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2 \
    && export PATH="$PATH:$(go env GOPATH)/bin" 
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY proto/ ./proto/
RUN protoc --proto_path=./proto --go-grpc_out=. service.proto \
    && protoc --proto_path=./proto --go_out=. service.proto
COPY menu_service/*.go ./
COPY menu_service/.env ./
RUN go build -o ./menu
EXPOSE 50001
CMD [ "./menu" ]