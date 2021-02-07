# Deployment to GCP

This section describes how to deploy to GCP.

## Requirement

* GCP account, that can use GKE.
* Docker
* kubectl

## Design


## Deployment

### 1. Build container image

Build API app container image.

```shell
$ cd deployments/docker-image/api
$ docker build -t <your name for GCR> .
```

After building, push image to Container Registry.

### 2. Set secrets, image name and environement variables

Set credentials in `apisecrets.yaml`, environment variables in `apisecrets.yaml` and container image name in `apiserver.yaml`.

See [Kubernetes components] for details.

### 3. Create GKE cluster.

Create cluster in GCP console or Cloud SDK.

```Shell
# Create cluster by Cloud SDK.
$ gcloud container clusters create cluster-name --num-nodes=1
```

### 4. Deploy to GKE.

Deploy each component to the cluster.

```shell
$ kubectl apply -f ./namespace.yaml

$ kubectl apply -f ./apisecrets.yaml
$ kubectl apply -f ./apiservice.yaml

$ kubectl apply -f ./envoyservice.yaml
```

## Kubernetes components

### apisecrets.yaml

The secrets for the app uses AWS and Google API. Set each value to your environment.

See [Environment variables]() for details.

**GOOGLE_API_KEY**

Set credential JSON as string to GOOGLE_API_KEY. The below command output value to set.

```shell
$ cat <GOOGLE_APPLICATION_CREDENTIALS key.json> | tr -d '\n' | base64
```

Ref.
https://kubernetes.io/docs/concepts/configuration/secret/

### apiservice.yaml

Service and Deployment of the app, which serves API with gRPC on port 50051.

Set container image name in the deployment containers section.

**Environment variables**

set DEBUG_LEVEL.

See [Environment variables]() for details.

### envoyconfigmap.yaml

This config map defines config.yaml in which envoy service.

### envoyservcie.yaml

Service and Deployment of envoy, which dispatch requests as the load balancer. Port 80 is connected to port 50051 on the gRPC server and port 9901 is connected envoy admin view.

### namespace.yaml

Defines the app's namespace.


## Ref.


[grpc-web/envoy.yaml at master · grpc/grpc-web](https://github.com/grpc/grpc-web/blob/master/net/grpc/gateway/examples/echo/envoy.yaml)

[Deploying a containerized web application](https://cloud.google.com/kubernetes-engine/docs/tutorials/hello-app)

[Envoy で HTTPS 接続をする設定を学べる「Securing traffic with HTTPS and SSL/TLS」を試した](https://kakakakakku.hatenablog.com/entry/2019/12/06/143207)

[HTTP(S) 負荷分散用 GKE Ingress](https://cloud.google.com/kubernetes-engine/docs/concepts/ingress)
[Ingress 機能の構成](https://cloud.google.com/kubernetes-engine/docs/how-to/ingress-features)
[ヘルスチェックの概要
](https://cloud.google.com/load-balancing/docs/health-check-concepts)
[Envoy プロキシを使用して GKE 上で gRPC サービスの負荷分散を行う](https://cloud.google.com/solutions/exposing-grpc-services-on-gke-using-envoy-proxy)
[gRPC & HTTP servers on GKE Ingress failing healthcheck for gRPC backend
](https://stackoverflow.com/questions/56277949/grpc-http-servers-on-gke-ingress-failing-healthcheck-for-grpc-backend)
