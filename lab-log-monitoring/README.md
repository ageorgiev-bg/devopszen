# Log Monitoring System - Loki
---
# What?
In this home lab, we will bring up a **Log Agreegation and Monitoring system** using **Grafana Loki**. This will allow us to bring up the whole stack on our local lab, Configure it and gain experience working with this tool. 
___
# Why?
Having a proper Log Aggregation and Monitoring is crucial to:
 - Extend observability
 - Better understanding of system under different load types
 - Detect and fix bugs proactively   

---
# How?
Among different Log aggregation tools avaulable, we will **demonstrate Grafana Loki**, Known for its efficiently, reliability and scaleability. We will go through this lab in following steps:
 - Bringing up a ***Kubernetes cluster***
 - Deploy ***ArgoCD*** as our deployment manager
 - Configure and deploy ***Loki***
 - Loki integration with ***Google Cloud Storage (GCS)*** as storage solution
 - Bringing up ***Kafka Cluster***
 - Configure ***Grafana Alloy*** to consume sent logs to Kafka and send them to ***Loki*** 

### Bringing up a Kubernetes cluster

We use Kubernetes as our container orchestration system. All components Including *Loki* and *Alloy* will be deployed on this Kubernetes cluster. 
**We use *Kind* for bringing up Kubernetes clusters**. Required Kind clusters configuration can be found under ***infrastructure*** directory. Following commands are being used to bring up the clusters:


```bash
cd infrastructure
kind create cluster --name lab-log-monitoring --config ./lab-log-monitoring.yaml 
```

We shall have the Kubernetes cluster with 2 worker nodes up and running.

Change the context to this cluster and then install ArgoCD via Helm

```bash
kubectx kind-lab-log-monitoring
helm repo add argo https://argoproj.github.io/argo-helm
helm repo update
kubectl create ns argocd
helm install argocd argo/argo-cd -n argocd
```
If not via ingress, You can always use `kubectl port-forward ...` to access to the UI. Consider following as an example:
```bash
kubectl port-forward service/argocd-server -n argocd --address 0.0.0.0 8080:443
``` 
NOTE: use `admin` as username. Password can be found via following command
```bash
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d
```

Following the steps, You will have fresh instance of ArgoCD up and running.

***NOTE: ARGOCD will be used to deploy all other required component for this Lab***

Next we need to add this git repository to ArgoCD attached repositories. We do that via argoCD UI. ***All applications that are supposed to be deployed by ArgoCD, including Strimzi operator, Kafka, Loki and Alloy will be placed under applications directory of this repository***

### Create a ArgoCD root application
As all the components are already places under ***applications*** Directory, you only need to create a root application and target this directory. ***Having *Recursive* feature selected, ArgoCD will iterate through this directory and deploy all components***