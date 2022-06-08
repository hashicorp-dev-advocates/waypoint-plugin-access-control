FROM golang:1.17-alpine3.15 as build

# Install the Protocol Buffers compiler and Go plugin
RUN apk add protobuf git make zip

COPY . $GOPATH/src/hashicorp/plugin
WORKDIR $GOPATH/src/hashicorp/plugin

RUN go get -d -v

RUN go get github.com/golang/protobuf/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc

# Build the plugin
RUN chmod +x ./print_arch
RUN make all

# Create the zipped binaries
RUN make zip

FROM scratch as export_stage

COPY --from=build /go/src/hashicorp/plugin/bin/*.zip .