# Ayanami

A FaaS framework for your own infrastructure.

The name is inspired from Evangelion-Unit-00's pilot: Ayanami Rei. The name `Rei` itself has the same pronunciation as in Heraclitus's philosophy, `Panta Rhei` (lit: everything flows). We believe that the developer should focus more on data flows and transformations rather than managing infrastructures.

# Why

* FaaS is good since it let you focus on the code instead of infrastructure
* Any FaaS providers are prone to vendor lock-in
* Having your own infrastructure (e.g: kubernetes) while developing/deploying in FaaS manner is a good solution
* At some point, developers need to run the entire infrastructure in their local machine. In this case, installing kubernetes/minikube could be overkill

# Goal

Providing an environment with minimum dependencies in order to:

* Build & deploy FaaS
* Make kubernetes-ready artifacts
* Run the entire infrastructure locally without kubernetes

# Dependencies

* golang 1.2
* nats.io

# How

* Developer create functions. The functions can be written in any language, even binary
* Developer define flows (how the functions are connected to each others)
* Ayanami compose flows and functions into several microservices that can talk to each other using nats messaging.

# Project Structure

```
Projects
```
