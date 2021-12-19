# **BUILD**

## **Table of contents** 
### 1. [**About**](#about)
### 2. [**Environment**](#environment)
### 3. [**Development**](#development)
>#### 3.1. [**Style Guide**](#style-guide)
>#### 3.2. [**Tests**](#tests)
>#### 3.3. [**Security**](#security)
### 4. [**Production**](#production)       

## **About**

The **BUILD.md** is a file to check the environment and build specifications of **horusec-operator** project.


## **Environment**

- [**Golang**](https://go.dev/dl/): ^1.17.X
- [**Kubectl**](https://kubernetes.io/docs/tasks/tools/#kubectl): ^1.20.X
- [**Helm**](https://helm.sh/docs/intro/install/#from-script): ^3.7.X
- [**Kind**](https://kind.sigs.k8s.io/docs/user/quick-start/#installation): ^0.11.1
- [**GNU Make**](https://www.gnu.org/software/make/): ^4.2.X

## **Development**

At the root of the project, run the following command to bring the cluster up along with its dependencies:

```bash
make up-sample
```

After running the above command, the cluster will be available. Then run the command below to apply the settings to the development environment:

```bash
make apply-sample
```

### **Style Guide**

For source code standardization, the project uses the [**golangci-lint**](https://golangci-lint.run) tool as a Go linter aggregator.
You can perform the lint check via the `make` command:

```bash
make lint
```

The project also has a dependency import pattern, and the commands below organize your code in the pattern defined by the Horusec team, run:

```bash
make fmt
```

Then, run the command:

```bash
make fix-imports
```

All project files must have the [**license header**](./copyright.txt). You can check if all files are in agreement, using the command:

```bash
make license
```

If it is necessary to add the license in any file, run the command below to insert it in all files that do not have it:

```bash
make license-fix
```

### **Tests**

Written with the [**Golang standard**](https://pkg.go.dev/testing), the unit tests were  package and some mock and assert snippets, we used the [**testify**](https://github.com/stretchr/testify). You can run the tests using the command below:

```bash
make tests
```

To check test coverage, run the command below:

```bash
make coverage
```

### **Security**

We use the latest version of the tool itself to maintain the security of our source code. Through the command below, you can perform this verification:

```bash
make security
```

## **Production**

The Horusec-Operator is one of the recommended ways to make [**Horusec-Platform**](https://github.com/ZupIT/horusec-platform) available in a production environment.

In our [**documentation**](https://docs.horusec.io/docs/web/installation/install-with-operator/overview/) you can check all the steps that must be followed.
