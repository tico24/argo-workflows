# Artifact Repository Ref

You can reduce duplication in your templates by configuring repositories that can be accessed by any workflow. This can also remove sensitive information from your templates.

Create a WorkflowArtifactRepository in either (a) your workflow's namespace or (b) the controller's namespace:

```yaml
apiVersion: argoproj.io/v1alpha1
kind: WorkflowArtifactRepository
metadata:
  name: my-repo
spec:
  s3:
    bucket: my-bucket
    endpoint: minio:9000
    insecure: true
    accessKeySecret:
      name: my-minio-cred
      key: accesskey
    secretKeySecret:
      name: my-minio-cred
      key: secretkey
```

The `spec` takes the same fields as the `artifactRepository` in the `workflow-controller-configmap`.
WorkflowArtifactRepositories need v4.2 or later.

Refer to it from a workflow by name:

```yaml
spec:
  artifactRepositoryRef:
    name: my-repo
```

The controller looks for `my-repo` in the workflow's namespace first, then in the controller's namespace.
If `name` is set, `configMap` and `key` are ignored.

This feature gives maximum benefit when used with [key-only artifacts](key-only-artifacts.md).

## Default Artifact Repository

A workflow without `artifactRepositoryRef` uses the first of these that exists:

1. A WorkflowArtifactRepository named `default` in the workflow's namespace.
1. The deprecated `artifact-repositories` config map in the workflow's namespace, if it has the `workflows.argoproj.io/default-artifact-repository` annotation.
1. A WorkflowArtifactRepository named `default` in the controller's namespace.
1. The deprecated `artifactRepository` in the `workflow-controller-configmap`.

The controller stores the resolved repository in the workflow's `status.artifactRepositoryRef`, so changing or deleting a WorkflowArtifactRepository doesn't affect workflows that have already started.

## RBAC

The workflow controller and the Argo Server need `get` on `workflowartifactrepositories`.
The install manifests grant it.
If they don't have it, the controller skips WorkflowArtifactRepositories and uses the deprecated config maps instead.

## Artifact Repository Config Maps (deprecated)

Artifact repository config maps are deprecated in v4.2, use WorkflowArtifactRepositories instead.
See [deprecations](deprecations.md#artifact-repository-configmap).

Create a suitable config map in either (a) your workflow's namespace or (b) the controller's namespace:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  # If you want to use this config map by default, name it "artifact-repositories". Otherwise, you can provide a reference to a
  # different config map in `artifactRepositoryRef.configMap`.
  name: my-artifact-repository
  annotations:
    # If you want to use a specific key, put that key into this annotation.
    workflows.argoproj.io/default-artifact-repository: default-v1-s3-artifact-repository
data:
  default-v1-s3-artifact-repository: |
    s3:
      bucket: my-bucket
      endpoint: minio:9000
      insecure: true
      accessKeySecret:
        name: my-minio-cred
        key: accesskey
      secretKeySecret:
        name: my-minio-cred
        key: secretkey
  v2-s3-artifact-repository: |
    s3:
      ...
```

You can override the artifact repository for a workflow as follows:

```yaml
spec:
  artifactRepositoryRef:
    configMap: my-artifact-repository # default is "artifact-repositories"
    key: v2-s3-artifact-repository # default can be set by the `workflows.argoproj.io/default-artifact-repository` annotation in config map.
```

[Reference](fields.md#artifactrepositoryref).
