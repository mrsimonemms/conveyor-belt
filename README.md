# Conveyor Belt

<!-- toc -->

* [Why?](#why)
* [Purpose](#purpose)
* [Getting Started](#getting-started)
  * [Including Dynamic Data](#including-dynamic-data)
* [Kubernetes](#kubernetes)
* [Project Status](#project-status)
* [Contributing](#contributing)
  * [Open in Gitpod](#open-in-gitpod)
  * [Open in devbox](#open-in-devbox)

<!-- Regenerate with "pre-commit run -a markdown-toc" -->

<!-- tocstop -->

Define your own pipeline of commands

## Why?

I was looking at creating an [OpenFaaS](https://openfaas.com) function which
made more sense to run as a pipeline. I did a bit of Googling, but wasn't able
to find anything that would _EASILY_ do what I wanted, so I spent a couple of
days doing this.

## Purpose

It is explicitly designed to chain together OpenFaaS functions. That's not to
say it can only be used with OpenFaaS, but that was in my mind when I was
building it.

## Getting Started

Define a config file, using
[examples/basic/config.yaml](./examples/basic//config.yaml) as a guide.

```yaml
# The apiVersion and kind are required
apiVersion: conveyor-belt.simonemms.com/v1alpha1
kind: Pipeline

# metadata.name is the name of the pipeline
metadata:
  name: basic

# The spec defines the pipeline
spec:
  port: 3000 # Port that this runs on

  # stages defines the names of the stages and the order that they're run in
  stages:
    - stage1
    - stage2

  # jobs define the jobs in the pipeline - this is an array
  jobs:
    - name: item1 # Name of the job - this can be anything
      stage: stage1 # Name of the stage - this must be in spec.stages
      timeout: 10s # Timeout, as a go duration - set to 0 to never timeout. Defaults to 30s
      # action - this is what happens
      action:
        # http - this is currently the only supported action
        http:
          method: POST # method - the HTTP method (eg, GET, POST, PUT, DELETE, PATCH etc)
          url: https://eosv8e8x84ccn8d.m.pipedream.net?stage=stage1&name=item1 # url - the URL to call
          # data is simple key/value pairs of HTTP body data to send. This is optional and can include dynamic data (see below)
          data:
            hello: world
            oi oi: true
            number: 2

  # This is a special job that is triggered if any of the jobs error
  error:
    # The action is as-above in jobs
    action:
      http:
        method: POST
        url: https://eosv8e8x84ccn8d.m.pipedream.net?stage=error&name=errorHandler

  # triggers defines how a pipeline can be started - this is an array of objects
  triggers:
    - type: webhook # webhook - this will receive on POST:/webhook/basic (basic is the name set in metadata.name)
```

### Including Dynamic Data

The action can include dynamic data received from a previous stage/job. In order
to guarantee that the data is present, this should refer to a previous stage as
jobs can be run in any order.

You can call it by specifying "{{ .Response.\<stage-name\>.\<item-name\>.Body.\<key\> }}"
in your `data` key/value pairs. A fully worked example would be
`name: "{{ .Response.stage1.item1.Body.name }}"`.

## Kubernetes

At some point, I may publish this to my [Helm registry](https://helm.simonemms.com).
Until then, use the following template:

```yaml
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: conveyor-belt-config
  labels:
    app: conveyor-belt
data:
  config.yaml: | # Insert your config file as-above
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: conveyor-belt
  labels:
    app: conveyor-belt
spec:
  replicas: 1
  selector:
    matchLabels:
      app: conveyor-belt
  template:
    metadata:
      labels:
        app: conveyor-belt
    spec:
      containers:
        - name: conveyor-belt
          image: ghcr.io/mrsimonemms/conveyor-belt:0.0.2
          args:
            - run
            - --config=/config/config.yaml
          ports:
            - containerPort: 3000
          volumeMounts:
            - name: config
              mountPath: /config
      volumes:
        - name: config
          configMap:
            name: conveyor-belt-config
---
apiVersion: v1
kind: Service
metadata:
  name: conveyor-belt
  labels:
    app: conveyor-belt
spec:
  selector:
    app: conveyor-belt
  ports:
    - protocol: TCP
      port: 3000
      targetPort: 3000
```

## Project Status

Currently, it only supports the `webhook` trigger and the only supported action is
`http`. This may change in future.

This is the result of a couple of days of work. If it proves useful, I may do
additional triggers/actions.

At the moment, all dynamic data is put through as a string. In future, I will
figure out a way of getting this to set a different data type (probably using
a ` | ToBool` style template function).

## Contributing

Please open an issue and propose a change before raising a PR.

### Open in Gitpod

* [Open in Gitpod](https://gitpod.io/from-referrer/)

### Open in devbox

* `curl -fsSL https://get.jetpack.io/devbox | bash`
* `devbox shell`
