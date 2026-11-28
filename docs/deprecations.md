# Deprecations

Sometimes a feature of Argo Workflows is deprecated.
You should stop using a deprecated feature as it may be removed in a future minor or major release of Argo Workflows.

To determine if you are using a deprecated feature the [`deprecated_feature`](metrics.md#deprecated_feature) metric can help.
This metric will go up for each use of a deprecated feature by the workflow controller.
This means it may go up once or many times for a single event.
If the number is going up the feature is still in use by your system.
If the metric is not present or no longer increasing are no longer using the monitored deprecated features.

## Deprecated features

These features are deprecated in v4.2 and still work.

### `artifact repository configmap`

Artifact repository config maps, and the `configMap` and `key` fields of `artifactRepositoryRef`, are replaced by WorkflowArtifactRepositories.
Create a WorkflowArtifactRepository for each config map key and refer to it with `artifactRepositoryRef.name`.
Name the one you want as the default `default`, in place of the `workflows.argoproj.io/default-artifact-repository` annotation.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: artifact-repositories
  annotations:
    workflows.argoproj.io/default-artifact-repository: my-key
data:
  my-key: |
    s3:
      bucket: my-bucket
```

is replaced with

```yaml
apiVersion: argoproj.io/v1alpha1
kind: WorkflowArtifactRepository
metadata:
  name: default
spec:
  s3:
    bucket: my-bucket
```

and `artifactRepositoryRef` with `key: my-key` is replaced with `name: default`.
See [Artifact Repository Ref](artifact-repository-ref.md).

[`argo convert`](cli/argo_convert.md) converts artifact repository config maps and `artifactRepositoryRef` keys for you:

```bash
kubectl get configmap artifact-repositories -n my-ns -o yaml | argo convert - | kubectl apply -f -
argo convert my-workflow.yaml
```

It writes one WorkflowArtifactRepository per config map key, named after the key, and writes the default key a second time as `default`.
It changes `artifactRepositoryRef` with a `key` to `name` set to that key.
It converts config maps named `artifact-repositories` and config maps with the `workflows.argoproj.io/default-artifact-repository` annotation.

### `controller configmap artifactRepository`

The `artifactRepository` in the `workflow-controller-configmap` is replaced by a WorkflowArtifactRepository named `default` in the controller's namespace.
Move the value of `artifactRepository` into the `spec` of that WorkflowArtifactRepository.
[`argo convert`](cli/argo_convert.md) does this for you, writing the config map without `artifactRepository` and the WorkflowArtifactRepository:

```bash
kubectl get configmap workflow-controller-configmap -n argo -o yaml | argo convert - | kubectl apply -f -
```

## Removed features

These features are removed in this version of workflows.
The [`convert`](cli/argo_convert.md) command in the CLI can update your Workflows, WorkflowTemplates, ClusterWorkflowTemplate and CronWorkflows to the modern syntax.

### `cronworkflow schedule`

The spec field `schedule` which takes a single value is replaced by `schedules` which takes a list.
To update this replace the `schedule` with `schedules` as in the following example

```yaml
spec:
  schedule: "30 1 * * *"
```

is replaced with

```yaml
spec:
  schedules:
    - "30 1 * * *"
```

### `synchronization mutex`

The synchronization field `mutex` which takes a single value is replaced by `mutexes` which takes a list.
To update this replace `mutex` with `mutexes` as in the following example

```yaml
synchronization:
  mutex:
    name: foobar
```

is replaced with

```yaml
synchronization:
  mutexes:
    - name: foobar
```

### `synchronization semaphore`

The synchronization field `semaphore` which takes a single value is replaced by `semaphores` which takes a list.
To update this replace `semaphore` with `semaphores` as in the following example

```yaml
synchronization:
  semaphore:
    configMapKeyRef:
      name: my-config
      key: workflow
```

is replaced with

```yaml
synchronization:
  semaphores:
    - configMapKeyRef:
        name: my-config
        key: workflow
```

### `workflow podpriority`

The Workflow spec field `podPriority` which takes a numeric value is deprecated and `podPriorityClassName` should be used instead.
To update this you will need a [PriorityClass](https://kubernetes.io/docs/concepts/scheduling-eviction/pod-priority-preemption/#priorityclass) in your cluster and refer to that using `podPriorityClassName`.

This field cannot be updated by the CLI `convert` command.
