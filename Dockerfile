FROM andersfylling/disgord:v0.10 as builder
MAINTAINER https://github.com/andersfylling
WORKDIR /build
COPY . /build
RUN export GO111MODULE=on
RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o discordbot .

FROM gcr.io/distroless/base
WORKDIR /bot
COPY --from=builder /build/discordbot .
CMD ["/bot/discordbot"]

