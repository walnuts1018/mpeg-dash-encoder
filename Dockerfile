FROM golang AS builder
ENV ROOT=/build
RUN mkdir ${ROOT}
WORKDIR ${ROOT}

RUN --mount=type=cache,target=/go/pkg/mod/,sharing=locked \
    --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod/,sharing=locked \
    --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    CGO_ENABLED=0 GOOS=linux go build -o main /$ROOT && chmod +x ./main

FROM debian:bookworm-slim
WORKDIR /app

RUN echo "Types: deb\n\
URIs: http://deb.debian.org/debian\n\
Suites: bookworm bookworm-updates\n\
Components: main non-free non-free-firmware\n\
Signed-By: /usr/share/keyrings/debian-archive-keyring.gpg\n\
\n\
Types: deb\n\
URIs: http://deb.debian.org/debian-security\n\
Suites: bookworm-security\n\
Components: main\n\
Signed-By: /usr/share/keyrings/debian-archive-keyring.gpg\n" > /etc/apt/sources.list.d/debian.sources

RUN cat /etc/apt/sources.list.d/debian.sources

RUN --mount=type=cache,target=/var/lib/apt,sharing=locked \
    --mount=type=cache,target=/var/cache/apt,sharing=locked \
    apt-get -y update && apt-get install -y ffmpeg libva-dev libmfx-dev intel-media-va-driver-non-free vainfo

COPY --from=builder /build/main ./
EXPOSE 8080

CMD ["./main"]
