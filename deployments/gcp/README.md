# Deployment to GCP GKE

This document describes how to deploy to GCP.

## Requirement

* Kuberntes and GCP knowledge
* Google Cloud SDK: 327.0.0, bq 2.0.64, core 2021.02.05, gsutil 4.58
* AWS account, that can use EKS and ECR
* GCP account, that can use Translation API.
* Docker: version 20.10.2
* kubectl: GitVersion v1.19.3
* GCP account, that can use GKE.
* kustomize: v3.10.0 
* macOS 10.15.x

Tested in the above environment.

Ref. [Setting up a local shell](https://cloud.google.com/container-registry/docs/quickstart#local-shell)

## Design

TODO

## Step

**The description is written under the following settings, please change the values as your environment.**

* GCP region: asia-northeast1-c
* GCR region: asia.gcr.io
* GCP Project ID translator-123
* GCR repository name: translatorapp-api 
* GKE cluster name: translatorapp-cluster
* GCP virtual global IP name: api-example-com-vip

### 1. Push docker image

Configure the request to GCR.

Ref. [Configure authentication](https://cloud.google.com/container-registry/docs/quickstart)

```sh
$ gcloud auth configure-docker
```


@TODO

Build docker image.

```sh
$ ../../scripts/buildimage.sh 123456.dkr.ecr.ap-northeast-1.amazonaws.com/translatorapp-api
```

Push the image.

```sh
$ docker push 123456.dkr.ecr.ap-northeast-1.amazonaws.com/translatorapp-api
```

### 2. Deploy to GKE

Create cluster.

```sh
$ gcloud container clusters create translatorapp-cluster \
    --num-nodes=1 \
    --enable-ip-alias \
    --create-subnetwork="" \
    --network=default \
    --zone=asia-northeast1-c

...

WARNING: Starting in January 2021, clusters ...
...
Created [https://container.googleapis.com/v1/projects/translator-123/zones/asia-northeast1-c/clusters/translatorapp-cluster].
To inspect the contents of your cluster, go to: https://console.cloud.google.com/kubernetes/workload_/gcloud/asia-northeast1-c/translatorapp-cluster?project=translator-123
kubeconfig entry generated for translatorapp-cluster.
NAME                   LOCATION           MASTER_VERSION    MASTER_IP      MACHINE_TYPE  NODE_VERSION      NUM_NODES  STATUS
translatorapp-cluster  asia-northeast1-c  1.17.14-gke.1600  xxx.xxx.xxx.xxx  e2-medium     1.17.14-gke.1600  1          RUNNING
```

Ref. [VPC ネイティブ クラスタを作成する](https://cloud.google.com/kubernetes-engine/docs/how-to/standalone-neg#create_a_vpc-native_cluster)


Create global IP.

```
gcloud compute addresses create api-example-com-vip \
  --ip-version=IPV4 \
  --global
...
Created [https://www.googleapis.com/compute/v1/projects/translator-123/global/addresses/api-example-com-api].
```

Ref. [Attaching an external HTTP(S) load balancer to standalone NEGs](https://cloud.google.com/kubernetes-engine/docs/how-to/standalone-neg#attaching-ext-https-lb) 

Set kubectl context to GKE.

```
$ gcloud container clusters get-credentials translatorapp-cluster

# Check
$ kubectl config current-context
gke_translator-123_asia-northeast1-c_translatorapp-cluster
```

Add GCR image ARN in `manifest/api/api-patch.yaml`.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-deployment
spec:
  template:
    spec:
      containers:
        - name: api-pod
          # Add the GCR image name.
          image: asia.gcr.io/translator-123/translatorapp-api:latest
```

Add the domain name in `manifest/api/ingress-cert-patch.yaml`.

```yaml
apiVersion: networking.gke.io/v1
kind: ManagedCertificate
metadata:
  name: envoy-cert
spec:
  domains:
    # Add the domain.
    - api.example.com
```

Add the static ip name and the domain name in `manifest/api/ingress-patch.yaml`.

```yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: ingress
  annotations:
    # Add the static ip name.
    kubernetes.io/ingress.global-static-ip-name: api-example-com
spec:
  rules:
  # Add the domain.
  - host: api.example.com
    http:
      paths:
      - path: /*
        backend:
          serviceName: envoy-service
          servicePort: 9000```
```


Apply.

```sh
$ cd manifest/api
$ kustomize build . | kubectl apply -f -
```

### 5. Setup DNS

Add DNS record to route the domain to the global IP created Step 2. To see IP, run the following command:

```
$ gcloud compute addresses list
```

Ref. [Step 4: Update the DNS A and AAAA records to point to the load balancer's IP address](https://cloud.google.com/load-balancing/docs/ssl-certificates/google-managed-certs?hl=ja#update-dns)

After a while, the certificate will be approved and you will be able to access it with the domain you specified.

### 6. Clean up

Delete components.

```
$ kustomize build . | kubectl delete -f -
```

Delete GKE cluster.

```
% gcloud container clusters delete translatorapp-cluster
The following clusters will be deleted.
 - [translatorapp-cluster] in [asia-northeast1-c]

Do you want to continue (Y/n)?  y

Deleting cluster translatorapp-cluster...⠶    
```


---------
Old


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






@TODO: 未整理
https://cloud.google.com/kubernetes-engine/docs/how-to/load-balance-ingress?hl=ja
