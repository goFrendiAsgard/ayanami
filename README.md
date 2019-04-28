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

```
.
├── composistions
│   └── composition.go   # Flows & Triggers definition (should be in golang).
│
├── dists                # Generated packages.
│
├── services             # Service definitions.
│   │                    #   Service can be written in any language
│   │                    #   as long as there is a template for it.
│   │                    #   But functions in a particular service
│   │                    #   should be written in the same language
│   ├── analytics
│   │   └── lib.js
│   ├── cities
│   │   └── lib.js
│   ├── Makefile
│   ├── populations
│   │   └── lib.js
│   ├── README.md
│   ├── test.js
│   └── weathers
│       └── lib.js
│
└── templates           # templates
    ├── composition/
    └── nodejs/
```

# Terminologies

* Composition: The functionality definition of your program
    - Flow: Composition of functions, usually triggered by a trigger (e.g: when there is a HTTP request to `/order` end point, the system should execute several functions from different services and return a response).
    - Trigger: Event that trigger flows (e.g: Scheduler, HTTP request, etc)
    - Service: Collection of functions. Usually from the same domain
        - Function: The atomic part of your business logic. Functions from different services should be independent from each other.
* Template: The template we use to generate package
* Package: The final source code of your program, ready for deployment
