FROM debian:bookworm-slim as builder

# Install required packages to compile
RUN apt-get update && \
        apt-get install -y gcc g++ ca-certificates git make wget tar xz-utils && \
        rm -rf /var/lib/apt/lists/*

# Install go 1.24.0
RUN wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz && \
    rm go1.24.0.linux-amd64.tar.gz

# Install upx
RUN wget https://github.com/upx/upx/releases/download/v4.0.1/upx-4.0.1-amd64_linux.tar.xz && \
    tar -xf upx-4.0.1-amd64_linux.tar.xz && \
    mv upx-4.0.1-amd64_linux/upx /usr/local/bin/upx && \
    rm -rf upx-4.0.1-amd64_linux*

ENV PATH=$PATH:/usr/local/go/bin

COPY . .

# Download deps and compile, then compress binary with upx
RUN go mod download && go build -ldflags="-s -w" -v -o vidsh . && upx -9 --lzma ./vidsh

FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update && \
    apt-get install --no-install-recommends -y ca-certificates ffmpeg && \
    rm -rf /var/lib/apt/lists/* /usr/share/doc/* /usr/share/man/* /usr/share/info/*

# Copy the UPX-compressed Go binary from the builder
COPY --from=builder /vidsh /app/

RUN useradd -m vidsh-runner
USER vidsh-runner

EXPOSE 8888

ENTRYPOINT ["./vidsh"]