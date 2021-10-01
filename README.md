<p></p>
<p></p>
<p align="center" margin="20 0"><a href="https://horusec.io/"><img src="https://raw.githubusercontent.com/ZupIT/horusec-devkit/main/assets/horusec_logo.png" alt="logo_header" width="65%" style="max-width:100%;"/></a></p>
<p></p>
<p></p>

<p align="center">
    <a href="https://github.com/ZupIT/horusec-engine/pulse" alt="activity">
        <img src="https://img.shields.io/github/commit-activity/m/ZupIT/horusec-operator?label=activity"/></a>
    <a href="https://github.com/ZupIT/horusec-operator/graphs/contributors" alt="contributors">
        <img src="https://img.shields.io/github/contributors/ZupIT/horusec-operator?label=contributors"/></a>
    <a href="https://github.com/ZupIT/horusec-operator/actions/workflows/lint.yml" alt="lint">
        <img src="https://img.shields.io/github/workflow/status/ZupIT/horusec-operator/Lint?label=lint"/></a>
    <a href="https://github.com/ZupIT/horusec-operator/actions/workflows/tests.yml" alt="tests">
        <img src="https://img.shields.io/github/workflow/status/ZupIT/horusec-operator/Test?label=test"/></a>
    <a href="https://github.com/ZupIT/horusec-operator/actions/workflows/security.yml" alt="security">
        <img src="https://img.shields.io/github/workflow/status/ZupIT/horusec-operator/Security?label=security"/></a>
    <a href="https://github.com/ZupIT/horusec-operator/actions/workflows/coverage.yml" alt="coverage">
        <img src="https://img.shields.io/github/workflow/status/ZupIT/horusec-operator/Coverage?label=coverage"/></a>
    <a href="https://opensource.org/licenses/Apache-2.0" alt="license">
        <img src="https://img.shields.io/badge/license-Apache%202-blue"/></a>
</p>

# Horusec-Operator
The main function of horusec-operator is to perform management between horusec web services and its kubernetes cluster.
Its creation came from an idea of the community where it can have a simpler way to install the services in an environment using kubernetes.
See all horusec operator details in [our documentation](https://horusec.io/docs/web/installation/install-with-horusec-operator/)

## Requirements
To use horusec-operator you need to configure some secrets and dependencies of horusec, they are:
* [Kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl) and connection with your cluster
* Connection with a database
    * You can upload a pod from a PostgreSQL database [as shown in our example](#development-environment), or you can only create secrets of connection with your database.
    * Create two databases for the horusec-platform and horusec-analytic 
* Connection with a message broker
    * You can upload a pod from a RabbitMQ message broker [as shown in our example](#development-environment), or you can only create secrets of connection with your message broker.
* Others secrets necessary
    * The secrets you need to configure may vary depending on how you use horusec. [See possible configuration options](https://horusec.io/docs/web/installation/install-with-horusec-operator#resources).

## Installing
After configuring your database connection, connecting to your broker and creating the secrets you need to install horusec-operator on your cluster, see an example below:
```bash
kubectl apply -k "https://github.com/ZupIT/horusec-operator/config/default?ref=v2.2.0"
```
See the resource if was installed with success!
```bash
kubectl api-resources | grep horus
```
you can see an output like this:
```text
$ kubectl api-resources | grep horus                                                           
horusecplatforms                  horus        install.horusec.io             true         HorusecPlatform
```

## Usage
And now just send the changes you want to kubenernetes. In this example we are using an [example yaml file](./config/samples/install_v2alpha1_horusecplatform.yaml), if you happen to send an empty yaml file like for example: 
```yaml
apiVersion: install.horusec.io/v2alpha1
kind: HorusecPlatform
metadata:
  name: horusecplatform-sample
spec: {}
```
It will take the [default horusec settings](./api/v2alpha1/horusec_platform_defaults.json)

And now you apply your changes
```bash
kubectl apply -f "https://raw.githubusercontent.com/ZupIT/horusec-operator/main/config/samples/install_v2alpha1_horusecplatform.yaml"
```

and you can see all horusec web services upload in your cluster, like this example:
```text
$ kubectl get pods
NAME                                                    READY   STATUS      RESTARTS   AGE
analytic-6f6bffb5d6-f8pl9                               1/1     Running     0          74s
api-5cc5b7545-km925                                     1/1     Running     0          73s
auth-8fbc876d9-62r6d                                    1/1     Running     0          73s
core-6bf7f9c9fc-fdv5c                                   1/1     Running     0          73s
horusecplatform-sample-analytic-migration-wwdzc-r9th2   0/1     Completed   0          74s
horusecplatform-sample-analytic-v1-2-v2-8zchl-445mz     0/1     Completed   2          74s
horusecplatform-sample-api-v1-2-v2-5lndp-w2rbd          0/1     Completed   3          74s
horusecplatform-sample-platform-migration-8g5ml-zmntl   0/1     Completed   0          74s
manager-c959f4f67-fz7r4                                 1/1     Running     0          74s
postgresql-postgresql-0                                 1/1     Running     0          7m54s
rabbitmq-0                                              1/1     Running     0          7m54s
vulnerability-7d789fd655-tpjp8                          1/1     Running     0          74s
webhook-7b5c45c859-cq4nf                                1/1     Running     0          73s
```

## Development Environment
This only an dev example how usage horusec-operator.
For usage this example is necessary installing [helm](https://helm.sh/docs/intro/install/#from-script) and [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation) in your local machine
After of you install you can run follow commands and see horusec-operator up all horusec web services.

Clone this project
```bash
git clone https://github.com/ZupIT/horusec-operator.git && cd horusec-operator
```

Up kubernetes cluster with all dependencies and wait finish! 
```bash
make up-sample
```

If you see this message
```text
Creating horusec_analytic_db...
If you don't see a command prompt, try pressing enter.
psql: could not connect to server: Connection refused
        Is the server running on host "postgresql" (10.96.182.42) and accepting
        TCP/IP connections on port 5432?
pod "postgresql-client" deleted
pod default/postgresql-client terminated (Error)
```
Don't worry this is normal because the script is trying create new database, but the pod of the postgresql is not ready, it will run again until create new database.

After script finish. Install Horusec-Operator
```bash
kubectl apply -k "https://github.com/ZupIT/horusec-operator/config/default?ref=v2.2.0"
```

See the resource if was installed with sucess!
```bash
kubectl api-resources | grep horus
```
you can see an output like this:
```text
$ kubectl api-resources | grep horus                                                           
horusecplatforms                  horus        install.horusec.io             true         HorusecPlatform
```

And you can see the pod manager by this resource 
```text
$ kubectl get pods -n horusec-operator-system
NAME                                                   READY   STATUS              RESTARTS   AGE
horusec-operator-controller-manager-7b9696d4c4-t7w2q   2/2     Running             0          2m10s
```

And now, you can pass [yaml with your configuration](https://horusec.io/docs/web/installation/install-with-horusec-operator#resources) to upload in your kubernetes cluster. See this example
```bash
kubectl apply -f ./config/samples/install_v2alpha1_horusecplatform.yaml
```

and you can see all horusec web services upload in your cluster, like this example:
```text
$ kubectl get pods
NAME                                                    READY   STATUS      RESTARTS   AGE
analytic-6f6bffb5d6-f8pl9                               1/1     Running     0          74s
api-5cc5b7545-km925                                     1/1     Running     0          73s
auth-8fbc876d9-62r6d                                    1/1     Running     0          73s
core-6bf7f9c9fc-fdv5c                                   1/1     Running     0          73s
horusecplatform-sample-analytic-migration-wwdzc-r9th2   0/1     Completed   0          74s
horusecplatform-sample-analytic-v1-2-v2-8zchl-445mz     0/1     Completed   2          74s
horusecplatform-sample-api-v1-2-v2-5lndp-w2rbd          0/1     Completed   3          74s
horusecplatform-sample-platform-migration-8g5ml-zmntl   0/1     Completed   0          74s
manager-c959f4f67-fz7r4                                 1/1     Running     0          74s
postgresql-postgresql-0                                 1/1     Running     0          7m54s
rabbitmq-0                                              1/1     Running     0          7m54s
vulnerability-7d789fd655-tpjp8                          1/1     Running     0          74s
webhook-7b5c45c859-cq4nf                                1/1     Running     0          73s
```

## Contributing Guide

Read our [contributing guide](CONTRIBUTING.md) to learn about our development process, how to propose bugfixes and improvements, and how to build and test your changes to horusec.

## Communication

We have a few channels for contact, feel free to reach out to us at:

- [GitHub Issues](https://github.com/ZupIT/horusec-operator/issues)
- [Zup Open Source Forum](https://forum.zup.com.br)

## Contributing with others projects

Feel free to use, recommend improvements, or contribute to new implementations.

If this is our first repository that you visit, or would like to know more about Horusec,
check out some of our other projects.

- [Horusec CLI](https://github.com/ZupIT/horusec)
- [Horusec Platform](https://github.com/ZupIT/horusec-platform)
- [Horusec DevKit](https://github.com/ZupIT/horusec-devkit)
- [Horusec Engine](https://github.com/ZupIT/horusec-engine)
- [Horusec Admin](https://github.com/ZupIT/horusec-admin)
- [Horusec VsCode](https://github.com/ZupIT/horusec-vscode-plugin)

This project exists thanks to all the contributors. You rock! ❤️🚀
