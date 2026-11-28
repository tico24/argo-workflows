package convert

import (
	"fmt"
	"maps"
	"slices"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	wf "github.com/argoproj/argo-workflows/v4/pkg/apis/workflow"
	wfv1 "github.com/argoproj/argo-workflows/v4/pkg/apis/workflow/v1alpha1"
	"github.com/argoproj/argo-workflows/v4/workflow/artifactrepositories"
)

const defaultArtifactRepositoryAnnotation = "workflows.argoproj.io/default-artifact-repository"

// ArtifactRepositoryConfigMap moves deprecated artifact repositories out of a config map into WorkflowArtifactRepositories.
// The artifactRepository key of a controller config map becomes a WorkflowArtifactRepository named default.
// An artifact repository config map becomes one WorkflowArtifactRepository per key,
// and its default key is also written as one named default.
// Any other config map is returned unchanged.
func ArtifactRepositoryConfigMap(cm *corev1.ConfigMap) ([]any, error) {
	if data, ok := cm.Data["artifactRepository"]; ok {
		repo, err := newWorkflowArtifactRepository(cm.Namespace, artifactrepositories.DefaultName, data)
		if err != nil {
			return nil, err
		}
		delete(cm.Data, "artifactRepository")
		return []any{cm, repo}, nil
	}
	defaultKey, ok := cm.Annotations[defaultArtifactRepositoryAnnotation]
	if !ok && cm.Name != "artifact-repositories" {
		return []any{cm}, nil
	}
	var out []any
	for _, key := range slices.Sorted(maps.Keys(cm.Data)) {
		names := []string{key}
		if key == defaultKey {
			names = append(names, artifactrepositories.DefaultName)
		}
		for _, name := range names {
			repo, err := newWorkflowArtifactRepository(cm.Namespace, name, cm.Data[key])
			if err != nil {
				return nil, err
			}
			out = append(out, repo)
		}
	}
	return out, nil
}

func newWorkflowArtifactRepository(namespace, name, data string) (*wfv1.WorkflowArtifactRepository, error) {
	repo := &wfv1.WorkflowArtifactRepository{
		TypeMeta:   metav1.TypeMeta{APIVersion: wfv1.SchemeGroupVersion.String(), Kind: wf.WorkflowArtifactRepositoryKind},
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
	}
	if err := yaml.UnmarshalStrict([]byte(data), &repo.Spec); err != nil {
		return nil, fmt.Errorf("failed to parse artifact repository %q: %w", name, err)
	}
	return repo, nil
}
