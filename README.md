# ChartLab

## About

ChartLab is a lightweight Kubernetes-native application which enables the use of private GitLab projects as Helm chart repositories. This is done by converting Helm's HTTP requests into GitLab API calls. Simply install ChartLab on your Kubernetes cluster and point Helm at it when adding a new repository.

## Installation

```sh
kubectl apply -f https://raw.githubusercontent.com/aiden-deloryn/chart-lab/main/k8s.yaml
```

## Usage

### Add a new private GitLab repository to Helm
Requirements:
- A [GitLab personal access token](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html).
- The GitLab Project ID for the repository you wish to add. You can find this on the project page at https://gitlab.com below the project name.
- Your project must contain a [Helm index file](https://helm.sh/docs/helm/helm_repo_index/).

```sh
# Get the node port of the ChartLab service
kubectl get service chartlab-service -n chartlab

# Add your private repository to Helm
helm repo add <repo-name> https://<node-ip>:<node-port>/<gitlab-project-id> --username '<username>' --password '<gitlab-personal-access-token>' --insecure-skip-tls-verify
```