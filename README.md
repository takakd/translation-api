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

### Create AWS and GCP account

This API uses AWS IAM and GCP service account to use each translation API.

For instructions on how to create an AWS account, see [Creating an IAM user in your AWS account
](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_users_create.html).

For a GCP account, see [Creating and managing service accounts](https://cloud.google.com/iam/docs/creating-managing-service-accounts).

### Prepare a TLS certificate

Create a server key file and TLS certificate if an API used TLS.

e.g., Self-signed certificate

```
$ cd manifest/api/secret
$ openssl genrsa -aes256 -passout pass:gsahdg -out server.pass.key 4096
$ openssl rsa -passin pass:gsahdg -in server.pass.key -out server.key
$ rm server.pass.key
$ openssl req -new -key server.key -out server.csr
...
$ openssl x509 -req -sha256 -days 365 -in server.csr -signkey server.key -out server.crt
```

Ref. [Generate private key and certificate signing request](https://devcenter.heroku.com/articles/ssl-certificate-self)

### Set environment variables

Need Several environment variables.

For details, see [.env.example](deployments/local/.env.example).

### Deployment

Several examples are here.

* [AWS EKS](deployments/eks/README.md)
* [GCP GKE](deployments/gke/README.md)

## Development

### Tech stacks

- Golang
- gRPC
- Kubernetes

### Requirements

- Golang: 1.14.4 darwin/amd64
- Docker: 20.10.2
- AWS IAM credentials, which can use Amazon Translate.
- Google service account, which can use Google Translation API.

We tested in the above environment.

### Setup

1. Install Golang by following [Download and install](https://golang.org/doc/install).
2. Run `go mod vendor` to get modules.

### Helper command

#### make

**Build**

```
$ make build
```

**go test**

Run test with details: "-v" and "-cover"

```
$ make test
```

**Format sources**

Run "go fmt", "goimports", and "go lint".

```
$ make fmt
```

#### On local 

```sh
$ ./scripts/local.sh
Usage:
  ./scripts/local.sh Command

Example.
  ./scripts/local.sh build

Command:
  run           Run envoy and gRPC server on local.
  run:go        Run gRPC server on local.
  run:envoy     Run envoy on local.
  down:envoy    Stop envoy on local.
  grpc          Generate gRPC codes.
```

### Structure

- Directory structure refers to [golang-standards/project-layout](https://github.com/golang-standards/project-layout).
- API services a translation service with gRPC.

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
|   |-- docker-image            <-- Working directory for a container image
|   |-- eks                     <-- EKS deployment examples
|   |-- gcp                     <-- GKE deployment examples
|   `-- local                   <-- Example running on local
|
|-- internal
|   |-- app                     <-- This api directory
|   |   |-- controller          <-- Controller layer
|   |   |   `-- translator      <-- gRPC handler
|   |   |       `-- ...
|   |   |-- driver              <-- Driver layer
|   |   |   |-- aws             <-- Codes related to handling AWS translate service
|   |   |   |   `-- ...
|   |   |   |-- config          <-- Concrete implementation of Config methods
|   |   |   |   `-- ...
|   |   |   |
|   |   |   |-- google          <-- Codes related to handling Google Translation API.
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
`-- scripts             <-- Scripts for this app
    |-- buildimage.sh   <-- For building container image
    |-- local.sh        <-- For local running
    |-- makefile.sh     <-- For Makefile
    `-- mock.pl         <-- To generate go mock file in the same directory
```

## Get in touch

- [Dev.to](https://dev.to/takakd)
- [Twitter](https://twitter.com/takakdkd)

## Contributing

Issues and reviews are welcome. Don't hesitate to create issues and PR.

## License

- Copyright 2020 © takakd.



