# Data Streaming Using Kafka and Kafka Connect
---
# What?
This effort is to demonstrate how we can create a data streaming system using Kafka and Kafka connect and practice on how to deploy, config and create monitoring for this system. 
___
# Why?
There are many cases in which we need to replicate data stored in one source into another system. A good example is a scenario that you need to have databases data in a data warehouse, running heavy queries and generating reports.  

---
# How?
Following steps will be taken care of:
 - Bringing up a Kubernetes cluster 
 - Deploy and configure a Kafka cluster
 - Deploy and configure a Kafka Connect cluster
 - Create a source connector 
 - Connect the connector to Postgresql database and establish a data pipeline.

### Bringing up a Kubernetes cluster

We use Kubernetes as our container orchestration system. Both Kafka and Kafka connect clusters will be managed inside Kubernetes. For the purpose of tutorial, We use *Kind* for bringing up Kubernetes clusters. Required Kind clusters configuration can be found under **infrastructure** directory. Following commands are being used to bring up the clusters:


```
cd infrastructure
kind create cluster --name lab-kafka-connect --config ./lab-kafka-connect.yaml
```
**NOTE: Do not forget to replace your interface IP address to config files**

We shall have the Kubernetes cluster with 2 worker nodes up and running.

Change the context to this cluster and then install ArgoCD via Helm

```
kubectx lab-kafka-connect
helm repo add argo https://argoproj.github.io/argo-helm
helm repo update
kubectl create ns argocd
helm install argocd argo/argo-cd -n argocd
```
If not via ingress, You can always use `kubectl port-forward ...` to access to the UI. Consider following as an example:
```
kubectl port-forward service/argocd-server -n argocd --address 0.0.0.0 8080:443
``` 
NOTE: use `admin` as username. Password can be found via following command
```
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d
```

Following the steps, You will have fresh instance of ArgoCD up and running.

Next we need to add this git repository to ArgoCD attached repositories. We do that via argoCD UI. ***All applications that are supposed to be deployed by ArgoCD, including Strimzi operator, Kafka and Kafka connect will be placed under applications directory of this repository***

Having the repository added to the ArgoCD repositories, components should be deployed in the following order:
 - **Strimzi-operator**: Kubernetes operator which we use to deploy Kafka and Kafka connect cluster. To have it deployed on the Kubernetes cluster, Create an ArgoCD application from ArgoCD ui. The target should be this repository and path should be ***lab-kafka-connect/applications/strimzi-operator***. Set the type of the application to ***directory*** 
 - **Kafka**: Kafka cluster that we use. As we have Operator installed in the first step, We just apply our **Custom Resource (CR)** to have kafka cluster up and running.Operator will act upon the applied CR and create the Kafka cluster in Kubernetes cluster.  To have it deployed on the Kubernetes cluster, Create another ArgoCD application from ArgoCD UI. The target should be this repository and path should be ***lab-kafka-connect/applications/kafka***. Set the type of the application to ***directory***.
 ***NOTE: This will also create the rest of the required components including kafka-ui, kafka-connect and kafka connector*** 
 - Kafka-UI: Simple interface that ease us working with Kafka
 - Kafka connect: Kafka connect cluster to manage connectors 
 - Kafka Connector: The source connector responsible to connect to the Postgresql service and stream Data changes to Kafka 
 

If all the configs set correctly, you will see the Zookeeper, Kafka, Kafka-UI, Kafka-connect and the connector being deployed in `kafka` namespace.
