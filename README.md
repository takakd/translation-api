<h1 align="center">Translation gRPC API</h1>

<p align="center">Translate text with <a href="https://aws.amazon.com/jp/translate/" alt="Amazon Translate">Amazon Translate</a> and <a href="https://cloud.google.com/translate/" alt="Google Translation">Google Translation</a>.</p>

<p align="center">
<a target="_blank" rel="noopener noreferrer" href="https://camo.githubusercontent.com/a568b3692dcc72af17d4abfed1b2c81d47f05dcaaefb021c9f9d3d6a856d3e6e/68747470733a2f2f696d672e736869656c64732e696f2f62616467652f4c6963656e73652d4d49542d696e666f726d6174696f6e616c3f7374796c653d666c6174"><img src="https://camo.githubusercontent.com/a568b3692dcc72af17d4abfed1b2c81d47f05dcaaefb021c9f9d3d6a856d3e6e/68747470733a2f2f696d672e736869656c64732e696f2f62616467652f4c6963656e73652d4d49542d696e666f726d6174696f6e616c3f7374796c653d666c6174" alt="License-MIT" data-canonical-src="https://img.shields.io/badge/License-MIT-informational?style=flat" style="max-width:100%;"></a>
</p>

<br>

## Table of Contents

- [Features](#features)
- [Usage](#usage)
- [Development](#development)
- [License](#license)

## Features

- Translate text with Google Translation API and Amazon Translate API.

## Setup

### Set environment variables

Set each environment variable to each Lambda function the following.

Key | Value | e.g.
---- | ---- | ---
DEBUG_LEVEL | Log level. Set `DEBUG` or `INFO` | DEBUG
GRPC_PORT | API port number |  50051
AWS_ACCESS_KEY_ID | AWS AccessKeyID to use Amazon Translate | AKI...
AWS_SECRET_ACCESS_KEY | AWS SecretAccessKey to use Amazon Translate | 4pfWR38...
AWS_REGION | AWS region to use Amazon Translate | ap-northeast-1
GOOGLE_PROJECT_ID | GOOGLE projectID to use Google Translation API | translator-123456
GOOGLE_API_KEY | GOOGLE service account credential JSON use Google Translation API |  {  "type": "service_account",  "project_id": "example",  "private_key_id": "0000000000" ... }

See also [.env.example](cmd/api/.env.example).

## Usage

Run gRPC server.

```
$ ./cmd/api/api
```

### Use .env file

Set `ENV_FILE` to use .env file.

```
$ ENV_FILE=/some/where/.env ./cmd/api/api
```

## Example

### Run on Kubernets service.

- [AWS EKS]()
- [Google GKE]()

### Translation application with this API.

- []


## Development

### Tech stacks

- Golang
- gRPC

### Setup

1. Install Golang by following [Download and install](https://golang.org/doc/install).
2. Run `go mod vendor` to get modules.

#### Requirements

- go version go1.14.4 darwin/amd64
- AWS IAM credentials that can use Amazon Translate.
- Google service account that can use Google Translation API.

### Command

**Testing**

Run test With details: "-v" and "-cover"

```
$ make test
```

**Formatting codes**

Run "go fmt", "goimports", and "go lint".

```
$ make fmt
```

**Run**

Run on local, use `cmd/api/.env.local` if it exists.

```
$ make run
```

### Structure

- Directory structure refers to [golang-standards/project-layout](https://github.com/golang-standards/project-layout).
- Serve gRPC API with envoy where the client apps requests.

#### Design

![Design](docs/design.jpg?raw=true)

#### Sources

```sh
.
|-- Makefile            <-- Defines make command targets
|-- README.layout.md    <-- golang-standards/project-layout README
|-- README.md           <-- This instruction file
|-- cmd
|   `-- api
|       |-- .env.local      <-- Environment variables on local
|       |-- .env.example    <-- Environment variables example
|       |-- api             <-- This API binary
|       `-- api.go          <-- main func
|
|-- deployments
|   `-- envoy                   <-- Envoy config directory
|       |-- docker-compose.yml  <-- docker-compose config for local
|       `-- envoy.yaml          <-- Envoy config
|
|-- docs
|   `-- logo.svg
|-- go.mod      <-- go module list
|-- go.sum      <-- go module hash list
|
|-- internal
|   |-- app                     <-- This app directory
|   |   |-- controller          <-- Controller layer
|   |   |   `-- translator      <-- gRPC handler
|   |   |       `-- ...
|   |   |-- driver              <-- Driver layer
|   |   |   |-- aws             <-- Codes related to handle AWS translate service
|   |   |   |   `-- ...
|   |   |   |-- config          <-- Concrete implementation of Config methods
|   |   |   |   `-- ...
|   |   |   |
|   |   |   |-- google          <-- Codes related to handle Google Translation API.
|   |   |   |   `-- ...
|   |   |   |-- grpcserver      <-- gRPC server
|   |   |   |   `-- ...
|   |   |   `-- log             <-- Concrete implementation of Logger methods
|   |   |       `-- ...
|   |   |-- grpc                <-- Auto generated gRPC codes
|   |   |   `-- translator
|   |   |       |-- translator.proto    <-- gRPC service definition
|   |   |       `-- ...
|   |   |-- initializer.go      <-- App initializer func
|   |   `-- util                <-- Codes shared throughout the app
|   |       |-- config          <-- Config
|   |       |   |-- config.go
|   |       |   `-- ...
|   |       |-- di              <-- DI
|   |       |   |-- container   <-- Concrete implementation of DI methods
|   |       |   |   |-- dev
|   |       |   |   `-- ...
|   |       |   |-- di.go
|   |       |   `-- ...
|   |       `-- log             <-- Logging
|   |           |-- log.go
|   |           `-- ...
|   `-- pkg             <-- Codes shared, which are not dependent on the app
|       `-- util        <-- Helper functions
|           |-- file.go
|           |-- http.go
|           |-- time.go
|           `-- type.go
|
`-- scripts         <-- Scripts for this app
    |-- local.sh    <-- Script used by Makefile
    `-- mock.pl     <-- Script to generate go mock file in the same directory
```

## Get in touch

- [Dev.to](https://dev.to/takakd)
- [Twitter](https://twitter.com/takakdkd)

## Contributing

Issues and reviews are welcome. Don't hesitate to create issues and PR.

## License

- Copyright 2020 © takakd.



