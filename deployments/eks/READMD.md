# Deployment to AWS EKS

This document describes how to deploy to EKS.

## Requirement

* Kubernetes and AWS knowledge
* AWS account, which can use EKS, ECR and Amazon Translate.
* GCP account, which can use Translation API.
* AWS CLI: 2.1.22
* Docker: 20.10.2
* kubectl: GitVersion v1.19.3
* kustomize: v3.10.0 
* macOS: 10.15.x

We tested in the above environment.

## Design

TODO

## Step

**This description is written under the following settings. Change the values as your environment.**

* AWS region: ap-northeast-1
* AWS account: 123456
* ECR repository name: translatorapp-api 
* Domain is api.example.com

### 1. Push docker image to ECR

Create an ECR repository.

```sh
$ aws ecr create-repository --repository-name translatorapp-api --region ap-northeast-1
...
{
    "repository": {
        "repositoryArn": "arn:aws:ecr:ap-northeast-1:123456:repository/translatorapp-api",
        ...
```

Logged in ECR.

```sh
$ aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin 123456.dkr.ecr.ap-northeast-1.amazonaws.com
```

Build docker image.

```sh
$ <this repository root>/scripts/buildimage.sh 123456.dkr.ecr.ap-northeast-1.amazonaws.com/translatorapp-api
```

Push the image.

```sh
$ docker push 123456.dkr.ecr.ap-northeast-1.amazonaws.com/translatorapp-api   
```

### 2. Create ACM

Create a certificate by ACM, which an ALB uses.

For instruction on how to create, see the following resources: [Issuing and Managing Certificates](https://docs.aws.amazon.com/acm/latest/userguide/gs.html)

### 3. Prepare secrets.

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

### 4. Deploy to EKS

Create an EKS cluster.

```sh
$ eksctl create cluster -f cluster.yaml
```

Create IAM Policy for the ingress controller and note `Arn`.


```sh
$ curl -O https://raw.githubusercontent.com/kubernetes-sigs/aws-alb-ingress-controller/v1.1.5/docs/examples/iam-policy.json

$ aws iam create-policy \
    --policy-name EKSALBIngressControllerPolicy \
    --policy-document file://iam-policy.json

"Policy": {
        "PolicyName": "EKSALBIngressControllerPolicy",
        "PolicyId": "AN.....",
        "Arn": "arn:aws:iam::123456:policy/EKSALBIngressControllerPolicy",
        ...
    }
```

Ref. https://github.com/kubernetes-sigs/aws-load-balancer-controller/issues/1171#issuecomment-593725742

Set kubectl config context.

```sh
$ aws eks --region ap-northeast-1 update-kubeconfig --name translatorapp-cluster
```

Create a service account.

*Set policy Arn to --attach-policy-arn option.*

```sh
$ eksctl create iamserviceaccount \
    --name alb-ingress-controller \
    --namespace kube-system \
    --cluster translatorapp-cluster \
    --attach-policy-arn arn:aws:iam::123456:policy/EKSALBIngressControllerPolicy \
    --override-existing-serviceaccounts \
    --approve
```

Add AWS Region and Cluster VPC ID in `manifest/eks-kube-system/alb-ingress-patch.yaml`

```yaml
spec:
  template:
    spec:
      containers:
        - name: alb-ingress-controller
          args:
          - --ingress-class=alb
          - --cluster-name=translatorapp-cluster
          - --aws-region=<Cluster Region>   #<-- Add the region where the cluster is.
          - --aws-vpc-id=<Cluster VPC ID>   #<-- Add the VPC ID 
```

Apply to the cluster.

```sh
$ cd manifest/eks-kube-system
$ kustomize build . | kubectl apply -f -
```

Add ECR image ARN in `manifest/api/api-patch.yaml`.

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
          # Add the ECR image ARN.
          image: 123456.dkr.ecr.ap-northeast-1.amazonaws.com/translatorapp-api:latest

```

Add the certificate ARN and the domain name in `manifest/api/ingress-patch.yaml`.

```yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: ingress
  annotations:
    # Add the certificate ARN.
    alb.ingress.kubernetes.io/certificate-arn: arn:aws:acm:ap-northeast-1:123456:certificate/....
spec:
  rules:
    # Add the domain.
    - host: api.example.com
      http:
        paths:
          - path: /*
            backend:
              serviceName: envoy-service
              servicePort: 9000

```

Apply to the cluster.

```sh
$ cd manifest/api
$ kustomize build . | kubectl apply -f -
```

### 5. Setup DNS

Add DNS record to route the domain to ALB CNAME.

To see ALB CNAME, run the following command:

```
$ kubectl get ingress -n translatorapp
```

*The command output should have the load balancer's fully qualified domain name (FQDN).*

### 6. Check

If some issues are, run the following command.:

```sh
$ kubectl logs your-alb-ingress-controller -n kube-system
```

### 7. Clean up

Delete all components.

```
$ cd manifest/api
$ kustomize build . | kubectl delete -f -

$ cd ../eks-kube-system
$ kustomize build . | kubectl delete -f -
```

Delete EKS cluster.

```
$ eksctl delete cluster -f cluster.yaml
```

**Sometimes some resources are still active, so check AWS console page to see whether all resources are deleted.**

## Ref. 
- AWS
    - [Ingress annotations](https://kubernetes-sigs.github.io/aws-load-balancer-controller/guide/ingress/annotations/#ingress-annotations)
    - [New – Application Load Balancer Support for End-to-End HTTP/2 and gRPC](https://aws.amazon.com/jp/blogs/aws/new-application-load-balancer-support-for-end-to-end-http-2-and-grpc/)
    - [ALBでgRPCを使う際にターゲット側もTLSしてみた](https://dev.classmethod.jp/articles/alb-e2e-tls/)
- Others
    - [Envoy documentation](https://www.envoyproxy.io/docs/envoy/latest/)
    - [grpc-web/envoy.yaml at master · grpc/grpc-web](https://github.com/grpc/grpc-web/blob/master/net/grpc/gateway/examples/echo/envoy.yaml)
    - [Envoy で HTTPS 接続をする設定を学べる「Securing traffic with HTTPS and SSL/TLS」を試した](https://kakakakakku.hatenablog.com/entry/2019/12/06/143207)
    - [How to configure HTTPS backends in envoy](https://farcaller.medium.com/how-to-configure-https-backends-in-envoy-b446727b2eb3)
    - [Simple SSL Termination with Envoy](https://timburks.me/2019/12/06/simple-ssl-termination-with-envoy)
    - [HTTP/2 Adventure in the Go World](https://posener.github.io/http2/)
    - [Introduction to HTTP/2](https://developers.google.com/web/fundamentals/performance/http2)
    - [Envoy Proxy Configuration](https://docs.build.security/docs/envoy)    
