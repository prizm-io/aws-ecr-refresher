# AWS ECR Credentials Refresher

This is a simple application that runs in a container and refreshes the AWS ECR credentials every 10 hours.
This is useful when you are running a container in a Kubernetes cluster and
you need to refresh the credentials to pull images from AWS ECR.

## Prerequisites

1. An IAM AWS account that has access to pulling images from ECR.
2. A Kubernetes cluster running somewhere.
3. A kubernetes service account with the necessary permissions to
create/delete secrets in the namespace where the application is running.

## Environment variables

```sh
AWS_ACCESS_KEY_ID=your-access-key-id
AWS_SECRET_ACCESS_KEY=your-secret-access-key
AWS_REGION=your-region

K8S_NAMESPACE=your-k8s-namespace
K8S_SECRET_NAME=your-k8s-secret-name

DOCKER_SERVER=your-docker-server # formatted as <aws-account-id>.dkr.ecr.<region>.amazonaws.com
DOCKER_EMAIL=your-docker-email
```

## Example Service Account
```yaml
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: sa-health-check
  namespace: default
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: role-full-access-to-secrets
rules:
  - apiGroups: [""]
    resources: ["secrets"]
    verbs: ["delete", "create"]
---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: health-check-role-binding
  namespace: default
subjects:
  - kind: ServiceAccount
    name: sa-health-check
    namespace: default
    apiGroup: ""
roleRef:
  kind: Role
  name: role-full-access-to-secrets
  apiGroup: ""
---
```
