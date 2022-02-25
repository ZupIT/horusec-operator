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

## **Table of contents**
### 1. [**About**](#about)
### 2. [**Usage**](#usage)
>#### 2.1. [**Requirements**](#requirements)
>#### 2.2. [**Installation**](#installation)
>#### 2.3. [**Configuration**](#configuration)
### 3. [**Development Environment**](#development-environment)
### 4. [**Documentation**](#documentation)
### 5. [**Issues**](#issues)
### 6. [**Contributing**](#contributing)
### 7. [**License**](#license)
### 8. [**Community**](#community)

## **About**
Horusec-operator performs management between horus web services and its Kubernetes cluster. It was created based on the community's idea and it can have a simpler way to install the services in an environment using Kubernetes.

This is the Kubernetes operator that enhances the Horusec-Platform installation system in a unified way. 
 
## **Usage**
### **Requirements**
See below the requirements to install and configure Horusec-Operator:
* [**Kubectl**](https://kubernetes.io/docs/tasks/tools/#kubectl) and connection with your cluster
* Connection with a database:
    * You can upload a pod from a PostgreSQL database. [**Check out the development environment example**](#development-environment), or you can create secrets of connection with your database.
    * Create two databases for Horusec-Platform and Horusec-Analytic. 
* Connection with a message broker:
    * You can upload a pod from a RabbitMQ message broker or you can create secrets of connection with your message broker.
* Other necessary secrets:
    * The secrets you need to configure may vary depending on how you use Horusec. [**Check out the configuration options**](https://docs.horusec.io/docs/web/installation/install-with-operator/yaml-definition/).

### **Installation**
Install Horusec-Operator on your cluster, see an example below:

1. Run the command: 

```bash
kubectl apply -k "https://github.com/ZupIT/horusec-operator/config/default?ref=v2.2.3"
```
2. Check if the resource was installed: 
```bash
kubectl api-resources | grep horus
```
3. You may see an output like this:
```text
$ kubectl api-resources | grep horus                                                           
horusecplatforms                  horus        install.horusec.io             true         HorusecPlatform
```

### **Configuration**

After installing, you need to send the changes you want to Kubenernetes. 

- In this example we are using a [**YAML file**](./config/samples/install_v2alpha1_horusecplatform.yaml).If you send an empty YAML file, it will take the [**default Horusec settings**](./api/v2alpha1/horusec_platform_defaults.json): 

```yaml
apiVersion: install.horusec.io/v2alpha1
kind: HorusecPlatform
metadata:
  name: horusecplatform-sample
spec: {}
```

- Apply your changes: 

```bash
kubectl apply -f "https://raw.githubusercontent.com/ZupIT/horusec-operator/main/config/samples/install_v2alpha1_horusecplatform.yaml"
```

- You can see all Horusec web services upload in your cluster, like the example below:
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
This is a development environment example on how to use Horusec-Operator.

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
kubectl apply -k "https://github.com/ZupIT/horusec-operator/config/default?ref=v2.2.3"
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

**Step 5.** And now, pass the [**yaml with your configuration**](https://docs.horusec.io/docs/web/installation/install-with-operator/yaml-definition/) to upload in your Kubernetes cluster. See the example:

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

For more information about Horusec, please check out the [**documentation**](https://docs.horusec.io/docs/).

## **Issues**

To open or track an issue for this project, in order to better coordinate your discussions, we recommend that you use the [**Issues tab**](https://github.com/ZupIT/horusec/issues) in the main [**Horusec**](https://github.com/ZupIT/horusec) repository.

## **Contributing**

If you want to contribute to this repository, access our [**Contributing Guide**](https://github.com/ZupIT/horusec-operator/blob/main/CONTRIBUTING.md). 

### **Developer Certificate of Origin - DCO**

 This is a security layer for the project and for the developers. It is mandatory.
 
 Follow one of these two methods to add DCO to your commits:
 
**1. Command line**
 Follow the steps: 
 **Step 1:** Configure your local git environment adding the same name and e-mail configured at your GitHub account. It helps to sign commits manually during reviews and suggestions.

 ```
git config --global user.name “Name”
git config --global user.email “email@domain.com.br”
```
**Step 2:** Add the Signed-off-by line with the `'-s'` flag in the git commit command:

```
$ git commit -s -m "This is my commit message"
```

**2. GitHub website**
You can also manually sign your commits during GitHub reviews and suggestions, follow the steps below: 

**Step 1:** When the commit changes box opens, manually type or paste your signature in the comment box, see the example:

```
Signed-off-by: Name < e-mail address >
```

For this method, your name and e-mail must be the same registered on your GitHub account.

## **License**
[**Apache License 2.0**](https://github.com/ZupIT/horusec-operator/blob/main/LICENSE).

## **Community**
Do you have any question about Horusec? Let's chat in our [**forum**](https://forum.zup.com.br/).


This project exists thanks to all the contributors. You rock! ❤️🚀

