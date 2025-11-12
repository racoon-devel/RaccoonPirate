FROM golang:1.25.4-bookworm as builder
WORKDIR /src/service
COPY . .
RUN apt-get update && apt-get install fuse libfuse-dev -y && make

FROM debian:bookworm
RUN apt-get update && apt-get install fuse ca-certificates -y
RUN mkdir /app
WORKDIR /app
COPY --from=builder /src/service/.build/raccoon-pirate .
COPY --from=builder /src/service/configs/raccoon-pirate.yml .
COPY --from=builder /src/service/fuse.conf /etc/fuse.conf
CMD ["./raccoon-pirate", "-config", "raccoon-pirate.yml"]