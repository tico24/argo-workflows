

# IoArgoprojWorkflowV1alpha1ArtifactRepositoryRef

ArtifactRepositoryRef is a reference to an artifact repository.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**configMap** | **String** | The name of the config map. Defaults to \&quot;artifact-repositories\&quot;. The config map is deprecated in v4.2, use name instead. |  [optional]
**key** | **String** | The config map key. Defaults to the value of the \&quot;workflows.argoproj.io/default-artifact-repository\&quot; annotation. The config map is deprecated in v4.2, use name instead. |  [optional]
**name** | **String** | The name of a WorkflowArtifactRepository in the workflow&#39;s namespace or the controller&#39;s namespace. If set, configMap and key are ignored. |  [optional]



