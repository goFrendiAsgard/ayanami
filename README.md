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

# How

* Developer create functions. The functions can be written in any language, even binary
* Developer define flows (how the functions are connected to each others)
* Ayanami compose flows and functions into several microservices that can talk to each other using nats messaging.

# Project Structure

TODO: adjust this with the new one
TODO: flows: exposing all endpoints for each service except analytics, compose things using weather, populations, and analytics

```
examples/greeter/
├── components/              # your functions
│   ├── user-function.go
│   ├── user-function.py
│   └── user-function.sh
├── dist/                    # generated microservices
├── flows.yml                # flow definitions
└── templates/               # your templates
    ├── service-template-1/
    ├── service-template-2/
    ├── trigger-template-1/
    └── trigger-template-2/
```

# Flow Definition

TODO: datatype

```yml
# filename: flows.yml

serviceTemplate:
    template1: service-template-1-location
    template2: service-template-2-location

triggerNames:
    - trigger1
    - trigger2

triggerConfigs:

    # e.g: http
    triggerGroup1:
        template: trigger-template1-location
        configs: {}

    # e.g: scheduler
    triggerGroup2:
        template: trigger-template1-location
        configs: {}

components:
    fn-1:
        path: user-function.go
        functionName: Fn
    fn-2:
        path: user-function.sh
    fn-3:
        path: user-function.py
        functionNamae: fn
flows:
    flow-name-1:
        - trigger1 |> fn-1 |> trigger2
        - trigger1 |> fn-3 |> trigger3
        - trigger3 |> fn-4
```
