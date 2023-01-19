# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.195.0/containers/go/.devcontainer/base.Dockerfile
# [Choice] Go version (use -bullseye variants on local arm64/Apple Silicon): 1, 1.16, 1.17, 1-bullseye, 1.16-bullseye, 1.17-bullseye, 1-buster, 1.16-buster, 1.17-buster
ARG VARIANT=1.19-bullseye
FROM mcr.microsoft.com/devcontainers/go:${VARIANT}

ARG DOCKER_USER_NAME=vscode
# ADD --chown=${DOCKER_USER_NAME}:${DOCKER_USER_NAME} . /home/${DOCKER_USER_NAME}/



RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends git \
    apt-transport-https \
    build-essential \
    ca-certificates \
    curl \
    gnupg \
    jq \
    protobuf-compiler \
    libprotobuf-dev \
# [Optional] Uncomment the next line to use go get to install anything else you need
    && go install github.com/posener/goreadme/cmd/goreadme@latest \
    && go install github.com/smartystreets/goconvey@latest \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install github.com/golang/protobuf/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest \
    && go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@latest \
    && go install github.com/infobloxopen/protoc-gen-gorm@latest \
    && go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest \
    && go install github.com/rakyll/statik@latest \
    && chown -R ${DOCKER_USER_NAME}:${DOCKER_USER_NAME} /go