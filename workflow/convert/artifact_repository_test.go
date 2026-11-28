package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	wfv1 "github.com/argoproj/argo-workflows/v4/pkg/apis/workflow/v1alpha1"
)

func TestArtifactRepositoryConfigMap(t *testing.T) {
	t.Run("ControllerConfigMap", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: "workflow-controller-configmap", Namespace: "argo"},
			Data:       map[string]string{"artifactRepository": "archiveLogs: true\ns3:\n  bucket: my-bucket\n", "parallelism": "10"},
		}
		out, err := ArtifactRepositoryConfigMap(cm)
		require.NoError(t, err)
		require.Len(t, out, 2)
		assert.Equal(t, map[string]string{"parallelism": "10"}, out[0].(*corev1.ConfigMap).Data)
		repo := out[1].(*wfv1.WorkflowArtifactRepository)
		assert.Equal(t, "WorkflowArtifactRepository", repo.Kind)
		assert.Equal(t, "argoproj.io/v1alpha1", repo.APIVersion)
		assert.Equal(t, "default", repo.Name)
		assert.Equal(t, "argo", repo.Namespace)
		assert.True(t, *repo.Spec.ArchiveLogs)
		assert.Equal(t, "my-bucket", repo.Spec.S3.Bucket)
	})
	t.Run("ArtifactRepositoriesConfigMap", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "my-repos",
				Namespace:   "my-ns",
				Annotations: map[string]string{"workflows.argoproj.io/default-artifact-repository": "b"},
			},
			Data: map[string]string{"b": "s3:\n  bucket: b\n", "a": "gcs:\n  bucket: a\n"},
		}
		out, err := ArtifactRepositoryConfigMap(cm)
		require.NoError(t, err)
		var names []string
		for _, o := range out {
			repo := o.(*wfv1.WorkflowArtifactRepository)
			assert.Equal(t, "my-ns", repo.Namespace)
			names = append(names, repo.Name)
		}
		assert.Equal(t, []string{"a", "b", "default"}, names)
		assert.Equal(t, "a", out[0].(*wfv1.WorkflowArtifactRepository).Spec.GCS.Bucket)
		assert.Equal(t, "b", out[2].(*wfv1.WorkflowArtifactRepository).Spec.S3.Bucket)
	})
	t.Run("ArtifactRepositoriesConfigMapWithoutAnnotation", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: "artifact-repositories"},
			Data:       map[string]string{"a": "s3:\n  bucket: a\n"},
		}
		out, err := ArtifactRepositoryConfigMap(cm)
		require.NoError(t, err)
		require.Len(t, out, 1)
		assert.Equal(t, "a", out[0].(*wfv1.WorkflowArtifactRepository).Name)
	})
	t.Run("OtherConfigMap", func(t *testing.T) {
		cm := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "other"}, Data: map[string]string{"a": "b"}}
		out, err := ArtifactRepositoryConfigMap(cm)
		require.NoError(t, err)
		assert.Equal(t, []any{cm}, out)
	})
	t.Run("Invalid", func(t *testing.T) {
		cm := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "artifact-repositories"}, Data: map[string]string{"a": "nope: true"}}
		_, err := ArtifactRepositoryConfigMap(cm)
		require.Error(t, err)
	})
}

func TestLegacyWorkflowSpec_ToCurrent_ArtifactRepositoryRef(t *testing.T) {
	for name, tt := range map[string]struct {
		in, want *wfv1.ArtifactRepositoryRef
	}{
		"Nil":       {nil, nil},
		"Key":       {&wfv1.ArtifactRepositoryRef{ConfigMap: "my-cm", Key: "my-key"}, &wfv1.ArtifactRepositoryRef{Name: "my-key"}},
		"Name":      {&wfv1.ArtifactRepositoryRef{Name: "my-repo"}, &wfv1.ArtifactRepositoryRef{Name: "my-repo"}},
		"ConfigMap": {&wfv1.ArtifactRepositoryRef{ConfigMap: "my-cm"}, &wfv1.ArtifactRepositoryRef{ConfigMap: "my-cm"}},
	} {
		t.Run(name, func(t *testing.T) {
			lws := &LegacyWorkflowSpec{WorkflowSpec: wfv1.WorkflowSpec{ArtifactRepositoryRef: tt.in}}
			assert.Equal(t, tt.want, lws.ToCurrent().ArtifactRepositoryRef)
		})
	}
}
