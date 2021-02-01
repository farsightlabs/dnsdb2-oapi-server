FROM golang:alpine AS build
WORKDIR /src

COPY go.mod /src
COPY go.sum /src
RUN go mod download

COPY . /src
RUN CGO_ENABLED=0 go build .

FROM scratch AS bin
COPY --from=build /src/dnsdb2-oapi-proxy .

EXPOSE 8086
ENTRYPOINT ["/dnsdb2-oapi-proxy"]
