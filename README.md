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

# **Horusec-Operator**
Horusec-operator performs management between horus web services and its Kubernetes cluster. It was created based on the community's idea and it can have a simpler way to install the services in an environment using Kubernetes.

This is the Kubernetes operator that enhances the Horusec-Platform installation system in a unified way. 
 
## **Requirements**
To use horusec-operator you need to configure some secrets and dependencies, see them below:
* [**Kubectl**](https://kubernetes.io/docs/tasks/tools/#kubectl) and connection with your cluster
* Connection with a database:
    * You can upload a pod from a PostgreSQL database [**you can see in the development environment example**](#development-environment), or you can create secrets of connection with your database.
    * Create two databases for horusec-platform and horusec-analytic. 
* Connection with a message broker:
    * You can upload a pod from a RabbitMQ message broker or you can create secrets of connection with your message broker.
* Other secrets necessary:
    * The secrets you need to configure may vary depending on how you use horusec. [**Check out the configuration options**](https://horusec.io/docs/web/installation/install-with-horusec-operator#resources).

## **Installing Operator**
After configuring your machine according to the requirements, install horusec-operator on your cluster, see an example below:

1. Run the command: 

```bash
kubectl apply -k "https://github.com/ZupIT/horusec-operator/config/default?ref=v2.2.1"
```
2. See if the resource was installed: 
```bash
kubectl api-resources | grep horus
```
3. You may see an output like this:
```text
$ kubectl api-resources | grep horus                                                           
horusecplatforms                  horus        install.horusec.io             true         HorusecPlatform
```

## **Usage**

After installing, you need to send the changes you want to Kubernetes. In this example we are using an [**example yaml file**](./config/samples/install_v2alpha1_horusecplatform.yaml), if you send an empty yaml file like the example below, it will take the [**default horusec settings**](./api/v2alpha1/horusec_platform_defaults.json): 

```yaml
apiVersion: install.horusec.io/v2alpha1
kind: HorusecPlatform
metadata:
  name: horusecplatform-sample
spec: {}
```

And now you apply your changes: 

```bash
kubectl apply -f "https://raw.githubusercontent.com/ZupIT/horusec-operator/main/config/samples/install_v2alpha1_horusecplatform.yaml"
```

You can see all horusec web services upload in your cluster, like this example:
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

## **Development Environment**
This is a development environment example on how to use horusec-operator.

You will need to install: 
- [**Helm**](https://helm.sh/docs/intro/install/#from-script) 
- [**Kind**](https://kind.sigs.k8s.io/docs/user/quick-start/#installation) 

Now, you can run the commands and see horusec-operator up all horusec web services. See the steps below: 

**Step 1.** Clone this project:

```bash
git clone https://github.com/ZupIT/horusec-operator.git && cd horusec-operator
```

**Step 2.** Up kubernetes cluster with all dependencies and wait finish: 

```bash
make up-sample
```

If you see this message:

```text
Creating horusec_analytic_db...
If you don't see a command prompt, try pressing enter.
psql: could not connect to server: Connection refused
        Is the server running on host "postgresql" (10.96.182.42) and accepting
        TCP/IP connections on port 5432?
pod "postgresql-client" deleted
pod default/postgresql-client terminated (Error)
```
Don't worry, this is normal because the script is trying to create new database, but the pod of the postgresql is not ready, it will run again until create new database.

**Step 3.** After the script finishes, install Horusec-Operator:
```bash
kubectl apply -k "https://github.com/ZupIT/horusec-operator/config/default?ref=v2.2.1"
```

**Step 4.** Check if the resource was installed: 

```bash
kubectl api-resources | grep horus
```
You can see an output like this:
```text
$ kubectl api-resources | grep horus                                                           
horusecplatforms                  horus        install.horusec.io             true         HorusecPlatform
```

And you can see the pod manager by the resource below: 

```text
$ kubectl get pods -n horusec-operator-system
NAME                                                   READY   STATUS              RESTARTS   AGE
horusec-operator-controller-manager-7b9696d4c4-t7w2q   2/2     Running             0          2m10s
```

**Step 5.** And now, pass the [**yaml with your configuration**](https://horusec.io/docs/web/installation/install-with-horusec-operator#resources) to upload in your Kubernetes cluster. See the example:

```bash
kubectl apply -f ./config/samples/install_v2alpha1_horusecplatform.yaml
```

You can see all horusec web services uploaded in your cluster, like this:
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

## **Documentation**

For more information about Horusec, please check out the [**documentation**](https://horusec.io/docs/).


## **Contributing**

If you want to contribute to this repository, access our [**Contributing Guide**](https://github.com/ZupIT/charlescd/blob/main/CONTRIBUTING.md). 
And if you want to know more about Horusec, check out some of our other projects:


- [**Admin**](https://github.com/ZupIT/horusec-admin)
- [**Charts**](https://github.com/ZupIT/charlescd/tree/main/circle-matcher)
- [**Devkit**](https://github.com/ZupIT/horusec-devkit)
- [**Jenkins**](https://github.com/ZupIT/horusec-jenkins-sharedlib)
- [**Platform**](https://github.com/ZupIT/horusec-platform)
- [**VSCode plugin**](https://github.com/ZupIT/horusec-vscode-plugin)
- [**Kotlin**](https://github.com/ZupIT/horusec-tree-sitter-kotlin)
- [**Vulnerabilities**](https://github.com/ZupIT/horusec-examples-vulnerabilities)

## **Community**
Feel free to reach out to us at:

- [**GitHub Issues**](https://github.com/ZupIT/horusec-devkit/issues)
- [**Zup Open Source Forum**](https://forum.zup.com.br)


This project exists thanks to all the contributors. You rock! ❤️🚀
