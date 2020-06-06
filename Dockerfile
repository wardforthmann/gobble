FROM golang:alpine AS build

RUN apk add git
WORKDIR /go/src
RUN git clone https://github.com/wardforthmann/gobble.git
WORKDIR /go/src/gobble
RUN go get ./
RUN go build -o gobble

FROM alpine
COPY --from=build /go/src/gobble/gobble /go/src/gobble/gobble
WORKDIR /go/src/gobble
EXPOSE 80
ENV GIN_MODE=release
ENTRYPOINT [ "./gobble" ] 