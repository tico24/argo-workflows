package artifactrepositories

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apierr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kubefake "k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"

	wfv1 "github.com/argoproj/argo-workflows/v4/pkg/apis/workflow/v1alpha1"
	wffake "github.com/argoproj/argo-workflows/v4/pkg/client/clientset/versioned/fake"
	"github.com/argoproj/argo-workflows/v4/util/logging"
)

func TestArtifactRepositories(t *testing.T) {
	defaultArtifactRepository := &wfv1.ArtifactRepository{
		S3: &wfv1.S3ArtifactRepository{KeyFormat: "foo"},
	}
	defaultArtifactRepositoryRefStatus := &wfv1.ArtifactRepositoryRefStatus{
		Default:            true,
		ArtifactRepository: defaultArtifactRepository,
	}
	k := kubefake.NewClientset()
	wf := wffake.NewSimpleClientset()
	i := New(k, wf, "my-ctrl-ns", defaultArtifactRepository)
	createRepo := func(t *testing.T, namespace, name, keyFormat string) {
		t.Helper()
		ctx := logging.TestContext(t.Context())
		_, err := wf.ArgoprojV1alpha1().WorkflowArtifactRepositories(namespace).Create(ctx, &wfv1.WorkflowArtifactRepository{
			ObjectMeta: metav1.ObjectMeta{Name: name},
			Spec:       wfv1.ArtifactRepository{S3: &wfv1.S3ArtifactRepository{KeyFormat: keyFormat}},
		}, metav1.CreateOptions{})
		require.NoError(t, err)
		t.Cleanup(func() {
			require.NoError(t, wf.ArgoprojV1alpha1().WorkflowArtifactRepositories(namespace).Delete(ctx, name, metav1.DeleteOptions{}))
		})
	}
	createConfigMap := func(t *testing.T, namespace string) {
		t.Helper()
		ctx := logging.TestContext(t.Context())
		_, err := k.CoreV1().ConfigMaps(namespace).Create(ctx, &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "artifact-repositories",
				Annotations: map[string]string{"workflows.argoproj.io/default-artifact-repository": "default-v1"},
			},
			Data: map[string]string{"default-v1": "s3:\n  keyFormat: configmap\n"},
		}, metav1.CreateOptions{})
		require.NoError(t, err)
		t.Cleanup(func() {
			require.NoError(t, k.CoreV1().ConfigMaps(namespace).Delete(ctx, "artifact-repositories", metav1.DeleteOptions{}))
		})
	}
	t.Run("Name.WorkflowNamespace", func(t *testing.T) {
		ctx := logging.TestContext(t.Context())
		createRepo(t, "my-wf-ns", "my-repo", "wf-ns")
		createRepo(t, "my-ctrl-ns", "my-repo", "ctrl-ns")

		ref, err := i.Resolve(ctx, &wfv1.ArtifactRepositoryRef{Name: "my-repo", ConfigMap: "ignored"}, "my-wf-ns")
		require.NoError(t, err)
		assert.Equal(t, &wfv1.ArtifactRepositoryRefStatus{
			Namespace:             "my-wf-ns",
			ArtifactRepositoryRef: wfv1.ArtifactRepositoryRef{Name: "my-repo"},
			ArtifactRepository:    &wfv1.ArtifactRepository{S3: &wfv1.S3ArtifactRepository{KeyFormat: "wf-ns"}},
		}, ref)

		repo, err := i.Get(ctx, &wfv1.ArtifactRepositoryRefStatus{Namespace: "my-wf-ns", ArtifactRepositoryRef: wfv1.ArtifactRepositoryRef{Name: "my-repo"}})
		require.NoError(t, err)
		assert.Equal(t, "wf-ns", repo.S3.KeyFormat)
	})
	t.Run("Name.ControllerNamespace", func(t *testing.T) {
		ctx := logging.TestContext(t.Context())
		createRepo(t, "my-ctrl-ns", "my-repo", "ctrl-ns")

		ref, err := i.Resolve(ctx, &wfv1.ArtifactRepositoryRef{Name: "my-repo"}, "my-wf-ns")
		require.NoError(t, err)
		assert.Equal(t, "my-ctrl-ns", ref.Namespace)
		assert.Equal(t, "ctrl-ns", ref.ArtifactRepository.S3.KeyFormat)
	})
	t.Run("Name.NotFound", func(t *testing.T) {
		ctx := logging.TestContext(t.Context())
		_, err := i.Resolve(ctx, &wfv1.ArtifactRepositoryRef{Name: "my-repo"}, "my-wf-ns")
		require.Error(t, err)
	})
	t.Run("DefaultRepo.WorkflowNamespaceBeatsConfigMap", func(t *testing.T) {
		ctx := logging.TestContext(t.Context())
		createRepo(t, "my-wf-ns", "default", "wf-ns")
		createConfigMap(t, "my-wf-ns")

		ref, err := i.Resolve(ctx, nil, "my-wf-ns")
		require.NoError(t, err)
		assert.Equal(t, "my-wf-ns", ref.Namespace)
		assert.Equal(t, "default", ref.Name)
		assert.Equal(t, "wf-ns", ref.ArtifactRepository.S3.KeyFormat)
	})
	t.Run("DefaultRepo.ConfigMapBeatsControllerNamespace", func(t *testing.T) {
		ctx := logging.TestContext(t.Context())
		createRepo(t, "my-ctrl-ns", "default", "ctrl-ns")
		createConfigMap(t, "my-wf-ns")

		ref, err := i.Resolve(ctx, nil, "my-wf-ns")
		require.NoError(t, err)
		assert.Equal(t, "artifact-repositories", ref.ConfigMap)
		assert.Equal(t, "configmap", ref.ArtifactRepository.S3.KeyFormat)
	})
	t.Run("DefaultRepo.ControllerNamespaceBeatsControllerConfig", func(t *testing.T) {
		ctx := logging.TestContext(t.Context())
		createRepo(t, "my-ctrl-ns", "default", "ctrl-ns")

		ref, err := i.Resolve(ctx, nil, "my-wf-ns")
		require.NoError(t, err)
		assert.Equal(t, "my-ctrl-ns", ref.Namespace)
		assert.False(t, ref.Default)
		assert.Equal(t, "ctrl-ns", ref.ArtifactRepository.S3.KeyFormat)
	})
	t.Run("Explicit.WorkflowNamespace", func(t *testing.T) {
		ctx := logging.TestContext(t.Context())
		_, err := k.CoreV1().ConfigMaps("my-wf-ns").Create(ctx, &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: "artifact-repositories"},
			Data: map[string]string{"my-key": `
s3:
  keyFormat: bar
`},
		}, metav1.CreateOptions{})
		require.NoError(t, err)

		ref, err := i.Resolve(ctx, &wfv1.ArtifactRepositoryRef{Key: "my-key"}, "my-wf-ns")
		require.NoError(t, err)
		assert.Equal(t, "my-wf-ns", ref.Namespace)
		assert.Equal(t, "artifact-repositories", ref.ConfigMap)
		assert.Equal(t, "my-key", ref.Key)
		assert.False(t, ref.Default)
		assert.NotNil(t, ref.ArtifactRepository)

		repo, err := i.Get(ctx, ref)
		require.NoError(t, err)
		assert.Equal(t, &wfv1.ArtifactRepository{S3: &wfv1.S3ArtifactRepository{KeyFormat: "bar"}}, repo)

		err = k.CoreV1().ConfigMaps("my-wf-ns").Delete(ctx, "artifact-repositories", metav1.DeleteOptions{})
		require.NoError(t, err)
	})
	t.Run("Explicit.ControllerNamespace", func(t *testing.T) {
		ctx := logging.TestContext(t.Context())
		_, err := k.CoreV1().ConfigMaps("my-ctrl-ns").Create(ctx, &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: "artifact-repositories"},
			Data: map[string]string{"my-key": `
s3:
  keyFormat: bar
`},
		}, metav1.CreateOptions{})
		require.NoError(t, err)

		ref, err := i.Resolve(ctx, &wfv1.ArtifactRepositoryRef{Key: "my-key"}, "my-wf-ns")
		require.NoError(t, err)
		assert.Equal(t, "my-ctrl-ns", ref.Namespace)
		assert.Equal(t, "artifact-repositories", ref.ConfigMap)
		assert.Equal(t, "my-key", ref.Key)
		assert.False(t, ref.Default)
		assert.NotNil(t, ref.ArtifactRepository)

		repo, err := i.Get(ctx, ref)
		require.NoError(t, err)
		assert.Equal(t, &wfv1.ArtifactRepository{S3: &wfv1.S3ArtifactRepository{KeyFormat: "bar"}}, repo)

		err = k.CoreV1().ConfigMaps("my-ctrl-ns").Delete(ctx, "artifact-repositories", metav1.DeleteOptions{})
		require.NoError(t, err)
	})
	t.Run("Explicit.ConfigMapNotFound", func(t *testing.T) {
		ctx := logging.TestContext(t.Context())
		_, err := i.Resolve(ctx, &wfv1.ArtifactRepositoryRef{}, "my-wf-ns")
		require.Error(t, err)
	})
	t.Run("Explicit.ConfigMapMissingKey", func(t *testing.T) {
		ctx := logging.TestContext(t.Context())
		_, err := k.CoreV1().ConfigMaps("my-ns").Create(ctx, &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: "artifact-repositories"},
		}, metav1.CreateOptions{})
		require.NoError(t, err)

		_, err = i.Resolve(ctx, &wfv1.ArtifactRepositoryRef{}, "my-wf-ns")
		require.Error(t, err)

		err = k.CoreV1().ConfigMaps("my-ns").Delete(ctx, "artifact-repositories", metav1.DeleteOptions{})
		require.NoError(t, err)
	})
	t.Run("WorkflowNamespaceDefault", func(t *testing.T) {
		ctx := logging.TestContext(t.Context())
		_, err := k.CoreV1().ConfigMaps("my-wf-ns").Create(ctx, &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "artifact-repositories",
				Annotations: map[string]string{"workflows.argoproj.io/default-artifact-repository": "default-v1"},
			},
			Data: map[string]string{"default-v1": `
s3:
  keyFormat: bar
`},
		}, metav1.CreateOptions{})
		require.NoError(t, err)

		ref, err := i.Resolve(ctx, nil, "my-wf-ns")
		require.NoError(t, err)
		assert.Equal(t, "my-wf-ns", ref.Namespace)
		assert.Equal(t, "artifact-repositories", ref.ConfigMap)
		assert.Equal(t, "default-v1", ref.Key)
		assert.False(t, ref.Default)
		assert.NotNil(t, ref.ArtifactRepository)

		repo, err := i.Get(ctx, ref)
		require.NoError(t, err)
		assert.Equal(t, &wfv1.ArtifactRepository{S3: &wfv1.S3ArtifactRepository{KeyFormat: "bar"}}, repo)

		err = k.CoreV1().ConfigMaps("my-wf-ns").Delete(ctx, "artifact-repositories", metav1.DeleteOptions{})
		require.NoError(t, err)
	})
	t.Run("DefaultWithNamespace", func(t *testing.T) {
		ctx := logging.TestContext(t.Context())
		_, err := k.CoreV1().ConfigMaps("my-wf-ns").Create(ctx, &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name: "artifact-repositories",
			},
		}, metav1.CreateOptions{})
		require.NoError(t, err)

		ref, err := i.Resolve(ctx, nil, "my-wf-ns")
		require.NoError(t, err)
		assert.Equal(t, defaultArtifactRepositoryRefStatus, ref)

		err = k.CoreV1().ConfigMaps("my-wf-ns").Delete(ctx, "artifact-repositories", metav1.DeleteOptions{})
		require.NoError(t, err)
	})
	t.Run("Default", func(t *testing.T) {
		ctx := logging.TestContext(t.Context())
		ref, err := i.Resolve(ctx, nil, "my-wf-ns")
		require.NoError(t, err)
		assert.Equal(t, defaultArtifactRepositoryRefStatus, ref)
	})
}

func TestArtifactRepositoriesForbidden(t *testing.T) {
	ctx := logging.TestContext(t.Context())
	wf := wffake.NewSimpleClientset()
	wf.PrependReactor("get", "workflowartifactrepositories", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, apierr.NewForbidden(schema.GroupResource{Resource: "workflowartifactrepositories"}, "default", nil)
	})
	defaultArtifactRepository := &wfv1.ArtifactRepository{S3: &wfv1.S3ArtifactRepository{KeyFormat: "foo"}}
	i := New(kubefake.NewClientset(), wf, "my-ctrl-ns", defaultArtifactRepository)

	ref, err := i.Resolve(ctx, nil, "my-wf-ns")
	require.NoError(t, err)
	assert.True(t, ref.Default)

	_, err = i.Get(ctx, &wfv1.ArtifactRepositoryRefStatus{Namespace: "my-wf-ns", ArtifactRepositoryRef: wfv1.ArtifactRepositoryRef{Name: "default"}})
	require.True(t, apierr.IsForbidden(err))
}
