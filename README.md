# ChartLab

- [About](#about)
- [Installation](#installation)
- [Usage](#usage)

## About

ChartLab is a lightweight Kubernetes-native application which enables the use of private GitLab projects as Helm chart repositories. This is done by converting Helm's HTTP requests into GitLab API calls. Simply install ChartLab on your Kubernetes cluster and point Helm at it when adding a new repository.

## Installation

### Create a TLS secret (optional)

If you want to use HTTPS/TLS to create a secure connection between Helm and ChartLab, create a new [Kubernetes TLS Secret](https://kubernetes.io/docs/concepts/configuration/secret/#tls-secrets) called `chartlab-tls` in the namespace `chartlab`.

```sh
# Create the chartlab namespace
kubectl apply -f https://raw.githubusercontent.com/aiden-deloryn/chart-lab/v1.0.0/config/namespace.yaml

# Create a TLS secret
kubectl create secret tls chartlab-tls --cert=<path-to-cert-file> --key=<path-to-key-file> --namespace chartlab
```

### Install ChartLab

Install chartlab using the command below.

```sh
kubectl apply -f https://raw.githubusercontent.com/aiden-deloryn/chart-lab/v1.0.0/config/deployment.yaml
```

## Usage

### Add a new private GitLab repository to Helm
Requirements:
- A [GitLab personal access token](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html).
- The GitLab Project ID for the repository you wish to add. You can find this on the project page at https://gitlab.com below the project name.
- Your project must contain a [Helm index file](https://helm.sh/docs/helm/helm_repo_index/).

Get the node port of the ChartLab service. This is the port you will use when adding the Helm repository. Note that there are two exposed ports, one for HTTP and the other for HTTPS.

```sh
kubectl get service chartlab-service -n chartlab
```

#### Adding a repository with HTTP

Add your private repository to Helm using `helm repo add`.

```sh
helm repo add <repo-name> http://<node-ip>:<node-port>/<gitlab-project-id> --username '<username>' --password '<gitlab-personal-access-token>'
```

#### Adding a repository with HTTPS/TLS

**Note:** If you are using a self-signed certificate, you will need to use the `--insecure-skip-tls-verify` for commands such as `helm repo add` and `helm install`.

Add your private repository to Helm using `helm repo add`.

```sh
helm repo add <repo-name> https://<node-ip>:<node-port>/<gitlab-project-id> --username '<username>' --password '<gitlab-personal-access-token>'
```