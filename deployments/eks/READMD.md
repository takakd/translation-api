# Deployment to EKS

This document describes how to deploy to EKS.

## Requirement

* AWS CLI: 2.1.22
* Docker: version 20.10.2
* kubectl: GitVersion v1.19.3
* AWS account, that can use EKS and ECR
* macOS 10.15.x
* Kuberntes and AWS knowledge

## Design

TODO

## Step

**The description is written under the following settings, please change the values as your environment.**

* AWS region: ap-northeast-1
* AWS account: 123456
* ECR repository name: translatorapp-api 
* Working directory: `deployments/eks`

### 1. Create ECR repository

```sh
$ aws ecr create-repository --repository-name translatorapp-api --region ap-northeast-1
...
{
    "repository": {
        "repositoryArn": "arn:aws:ecr:ap-northeast-1:123456:repository/translatorapp-api",
        ...
```
       
### 2. Push docker image

Logged in ECR.

```sh
$ aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin 123456.dkr.ecr.ap-northeast-1.amazonaws.com
```

@TODO

Build docker image.

```sh
$ docker build -t 123456.dkr.ecr.ap-northeast-1.amazonaws.com/translatorapp-api .
```

Push image.

```sh
$ docker push 123456.dkr.ecr.ap-northeast-1.amazonaws.com/translatorapp-api   
```

### 3. Deploy to EKS

Create cluster.

```sh
$ eksctl create cluster -f cluster.yaml
```

Create IAM Policy for the ingress controller and note `Arn`.

```sh
$ curl -O https://raw.githubusercontent.com/kubernetes-sigs/aws-alb-ingress-controller/v1.1.8/docs/examples/iam-policy.json

$ aws iam create-policy --policy-name EKSALBIngressControllerPolicy --policy-document file://iam-policy.json

"Policy": {
        "PolicyName": "EKSALBIngressControllerPolicy",
        "PolicyId": "AN.....",
        "Arn": "arn:aws:iam::123456:policy/EKSALBIngressControllerPolicy",
        ...
    }
$ rm iam-policy.json
```

Set kubectl config of EKS.

```sh
$ aws eks --region ap-northeast-1 update-kubeconfig --name translatorapp-cluster
```

Create service account component and attach IAM Policy.

```sh
$ kubectl apply -f manifest/rbac-role.yaml
```

Create service account.

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

Create ALB Ingress Controller.

*Set `VPCID` in manifest/alb-ingress-controller.yaml before applying.*

```sh
$ kubectl apply -f manifest/alb-ingress-controller.yaml
```

Create namespace.

```sh
$ kubectl apply -f manifest/namespace.yaml
```

Create each service.

Note: Rename example files to `apisecrets-env.yaml` and `apisecrets-file.yaml` in `manifest/secrets`, then set secret values to them before applying. See each example file for secret value details.

@TODO: generate secret and set
@TODO: ACM

```sh
# gRPC service.
$ kubectl apply -f manifest/apiconfig-env.yaml
$ kubectl apply -f manifest/secrets/apisecrets-env.yaml
$ kubectl apply -f manifest/secrets/apisecrets-file.yaml
$ kubectl apply -f manifest/apiservice.yaml

# Envoy for gRPC web proxy.

# No TSL
$ kubectl apply -f manifest/notls/envoyconfig-file.yaml
# Use TSL
# Set enabled certificates to manifest.
# $ kubectl apply -f manifest/envoyconfig-file.yaml

$ kubectl apply -f manifest/envoyservice.yaml
```

Create Ingress ALB.

```sh
# No TSL
$ kubectl apply -f manifest/notls/alb-ingress.yaml
# Use TSL
# $ kubectl apply -f manifest/tls/alb-ingress.yaml
```

Done.  
The endpoint is showed by the below command.

```sh
$ kubectl get ingress -n translatorapp
```

## Ref. 
**[Envoy documentation](https://www.envoyproxy.io/docs/envoy/latest/)**

- [TLS](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/security/ssl)
- [Securing Envoy](https://www.envoyproxy.io/docs/envoy/latest/start/quick-start/securing)
- [extensions.transport_sockets.tls.v3.UpstreamTlsContext](https://www.envoyproxy.io/docs/envoy/latest/api-v3/extensions/transport_sockets/tls/v3/tls.proto#extensions-transport-sockets-tls-v3-upstreamtlscontext)
- [How do I configure SNI for listeners?](https://www.envoyproxy.io/docs/envoy/latest/faq/configuration/sni#how-do-i-configure-sni-for-listeners)

**Others**
- [Ingress annotations](https://kubernetes-sigs.github.io/aws-load-balancer-controller/guide/ingress/annotations/#ingress-annotations)
- [How to configure HTTPS backends in envoy](https://farcaller.medium.com/how-to-configure-https-backends-in-envoy-b446727b2eb3)
- [Simple SSL Termination with Envoy](https://timburks.me/2019/12/06/simple-ssl-termination-with-envoy)
- [New – Application Load Balancer Support for End-to-End HTTP/2 and gRPC](https://aws.amazon.com/jp/blogs/aws/new-application-load-balancer-support-for-end-to-end-http-2-and-grpc/)
- [HTTP/2 Adventure in the Go World](https://posener.github.io/http2/)
- [Introduction to HTTP/2](https://developers.google.com/web/fundamentals/performance/http2)
- [Envoy Proxy Configuration](https://docs.build.security/docs/envoy)
- [ALBでgRPCを使う際にターゲット側もTLSしてみた](https://dev.classmethod.jp/articles/alb-e2e-tls/)
