# Deployment to GKE

This document describes how to deploy to GCP GKE.

## Requirement

* Kuberntes and GCP knowledge
* AWS account, which can use Amazon Translate.
* GCP account, which can use GKE, GCR, Networking Service, and Translation API.
* Google Cloud SDK: 327.0.0, bq 2.0.64, core 2021.02.05, gsutil 4.58
* Docker: 20.10.2
* kubectl: GitVersion v1.19.3
* kustomize: v3.10.0 
* macOS 10.15.x

We tested in the above environment.

## Design

![Design](design.jpg?raw=true)

## Step

**This description is written under the following settings. Change the values as your environment.**

* GCP region: asia-northeast1-c
* GCR region: asia.gcr.io
* GCP ProjectID: translator-123
* GCR repository name: translatorapp-api 
* GKE cluster name: translatorapp-cluster
* GCP virtual global IP name: api-example-com-vip

### 1. Push docker image to GCR

Configure the request to GCR.

```sh
$ gcloud auth configure-docker
```

Ref. [Configure authentication](https://cloud.google.com/container-registry/docs/quickstart)

Build docker image.

```sh
$ <this repository root>/scripts/buildimage.sh asia.gcr.io/translator-123/translatorapp-api
```

Push the image.

```sh
$ docker push asia.gcr.io/translator-123/translatorapp-api
```

### 2. Prepare secrets.

First, create a GCP service account to use Google Translation API, then save its certificate JSON as `manifest/api/secret/google.key.json`. 

Ref. [Creating and managing service account keys
](https://cloud.google.com/iam/docs/creating-managing-service-account-keys)

Generate a self-certified certification as `server.key` and `server.crt` in `manifest/api/secret`.

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


Add some secret values in `/manifest/api/secret/.env`.

```
AWS_ACCESS_KEY_ID=AK...         #<-- Add IAM Access Key ID
AWS_SECRET_ACCESS_KEY=XU...     #<-- Add IAM Secret Access Key
AWS_REGION=ap-northeast-1       #<-- Add IAM Region
GOOGLE_PROJECT_ID=...           #<-- Add GCP ProjectID where the GCP service account is.
```

### 3. Deploy to GKE

Create a GKE cluster.

```sh
$ gcloud container clusters create translatorapp-cluster \
    --num-nodes=1 \
    --enable-ip-alias \
    --create-subnetwork="" \
    --network=default \
    --zone=asia-northeast1-c
    --machine-type=e2-small

...

WARNING: Starting in January 2021, clusters ...
...
Created [https://container.googleapis.com/v1/projects/translator-123/zones/asia-northeast1-c/clusters/translatorapp-cluster].
To inspect the contents of your cluster, go to: https://console.cloud.google.com/kubernetes/workload_/gcloud/asia-northeast1-c/translatorapp-cluster?project=translator-123
kubeconfig entry generated for translatorapp-cluster.
NAME                   LOCATION           MASTER_VERSION    MASTER_IP      MACHINE_TYPE  NODE_VERSION      NUM_NODES  STATUS
translatorapp-cluster  asia-northeast1-c  1.17.14-gke.1600  xxx.xxx.xxx.xxx  e2-medium     1.17.14-gke.1600  1          RUNNING
```

Ref. [Create a VPC-native cluster](https://cloud.google.com/kubernetes-engine/docs/how-to/standalone-neg#create_a_vpc-native_cluster)


Create a global virtual IP.

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

Add GCR image name in `manifest/api/api-patch.yaml`.

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
          # Add a GCR image name.
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
    # Add a domain.
    - api.example.com
```

Add the global virtual IP name and the domain name in `manifest/api/ingress-patch.yaml`.

```yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: ingress
  annotations:
    # Add a static IP name.
    kubernetes.io/ingress.global-static-ip-name: api-example-com
spec:
  rules:
  # Add a domain.
  - host: api.example.com
    http:
      paths:
      - path: /*
        backend:
          serviceName: envoy-service
          servicePort: 9000```
```

Apply to the cluster.

```sh
$ cd manifest/api
$ kustomize build . | kubectl apply -f -
```

### 5. Setup DNS

Add DNS record to route the domain to the global virtual IP created Step 2. To see it, run the following command:

```
$ gcloud compute addresses list
```

Ref. [Step 4: Update the DNS A and AAAA records to point to the load balancer's IP address](https://cloud.google.com/load-balancing/docs/ssl-certificates/google-managed-certs?hl=ja#update-dns)

After a while, the certificate will be approved, and you will be able to access it with the domain you specified.

### 6. Clean up

Delete all components.

```
$ cd manifest/api
$ kustomize build . | kubectl delete -f -
```

Delete GKE cluster.

```
$ gcloud container clusters delete translatorapp-cluster
The following clusters will be deleted.
 - [translatorapp-cluster] in [asia-northeast1-c]

Do you want to continue (Y/n)?  y

Deleting cluster translatorapp-cluster...⠶    
```

**Sometimes, some resources are still active, so check GCP console page to see whether all resources are deleted.**

## Ref.
- GCP
    - [Deploying a containerized web application](https://cloud.google.com/kubernetes-engine/docs/tutorials/hello-app)
    - [GKE Ingress for HTTP(S) Load Balancing](https://cloud.google.com/kubernetes-engine/docs/concepts/ingress)
    - [Configuring Ingress features](https://cloud.google.com/kubernetes-engine/docs/how-to/ingress-features)
    - [Health checks overview](https://cloud.google.com/load-balancing/docs/health-check-concepts)
    - [Using Envoy Proxy to load-balance gRPC services on GKE](https://cloud.google.com/solutions/exposing-grpc-services-on-gke-using-envoy-proxy)
    - [gRPC & HTTP servers on GKE Ingress failing healthcheck for gRPC backend
    ](https://stackoverflow.com/questions/56277949/grpc-http-servers-on-gke-ingress-failing-healthcheck-for-grpc-backend)
- Others
    - [Envoy documentation](https://www.envoyproxy.io/docs/envoy/latest/)
    - [grpc-web/envoy.yaml at master · grpc/grpc-web](https://github.com/grpc/grpc-web/blob/master/net/grpc/gateway/examples/echo/envoy.yaml)
    - [How to configure HTTPS backends in envoy](https://farcaller.medium.com/how-to-configure-https-backends-in-envoy-b446727b2eb3)
    - [Simple SSL Termination with Envoy](https://timburks.me/2019/12/06/simple-ssl-termination-with-envoy)
    - [HTTP/2 Adventure in the Go World](https://posener.github.io/http2/)
    - [Introduction to HTTP/2](https://developers.google.com/web/fundamentals/performance/http2)
    - [Envoy Proxy Configuration](https://docs.build.security/docs/envoy)
    - [Envoy で HTTPS 接続をする設定を学べる「Securing traffic with HTTPS and SSL/TLS」を試した](https://kakakakakku.hatenablog.com/entry/2019/12/06/143207)
    
## License

&copy; 2021 takakd
