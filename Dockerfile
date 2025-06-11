# 開発環境
FROM golang:1.24

RUN curl -sSf https://temporal.download/cli.sh | sh
ENV TEMPORAL_HOME="/root/.temporalio/bin"
ENV PATH="$PATH:$TEMPORAL_HOME"

WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download
