Description: Configure artifact repositories with the WorkflowArtifactRepository kind
Authors: [Tim Collins](https://github.com/tico24)
Component: General
Issues: 00000

A WorkflowArtifactRepository holds one artifact repository, using the same fields as `artifactRepository` in the `workflow-controller-configmap`.
Refer to one from a workflow with `artifactRepositoryRef.name`, which looks in the workflow's namespace and then the controller's namespace.
A workflow without `artifactRepositoryRef` uses the WorkflowArtifactRepository named `default`, in its own namespace and then the controller's namespace.
Artifact repository config maps, `artifactRepositoryRef.configMap`, `artifactRepositoryRef.key` and `artifactRepository` in the `workflow-controller-configmap` are deprecated and still work.
`argo convert` converts artifact repository config maps and `artifactRepositoryRef` keys to WorkflowArtifactRepositories.
The workflow controller and Argo Server need `get` on `workflowartifactrepositories`, which the install manifests grant.
