package artifactrepositories

import (
	"context"
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"
	apierr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"

	wfv1 "github.com/argoproj/argo-workflows/v4/pkg/apis/workflow/v1alpha1"
	"github.com/argoproj/argo-workflows/v4/pkg/client/clientset/versioned"
	"github.com/argoproj/argo-workflows/v4/util/deprecation"
	errorsutil "github.com/argoproj/argo-workflows/v4/util/errors"
	"github.com/argoproj/argo-workflows/v4/util/logging"
	"github.com/argoproj/argo-workflows/v4/util/retry"
	waitutil "github.com/argoproj/argo-workflows/v4/util/wait"
)

// DefaultName is the name of the WorkflowArtifactRepository used when a workflow has no artifactRepositoryRef.
const DefaultName = "default"

type Interface interface {
	// Resolve Figures out the correct repository to for a workflow.
	Resolve(ctx context.Context, ref *wfv1.ArtifactRepositoryRef, workflowNamespace string) (*wfv1.ArtifactRepositoryRefStatus, error)
	// Get returns the referenced repository. May return nil (if no default artifact repository is configured).
	Get(ctx context.Context, ref *wfv1.ArtifactRepositoryRefStatus) (*wfv1.ArtifactRepository, error)
}

func New(kubernetesInterface kubernetes.Interface, wfClientset versioned.Interface, namespace string, defaultArtifactRepository *wfv1.ArtifactRepository) Interface {
	return &artifactRepositories{kubernetesInterface, wfClientset, namespace, defaultArtifactRepository}
}

type artifactRepositories struct {
	kubernetesInterface       kubernetes.Interface
	wfClientset               versioned.Interface
	namespace                 string
	defaultArtifactRepository *wfv1.ArtifactRepository
}

func (s *artifactRepositories) Resolve(ctx context.Context, ref *wfv1.ArtifactRepositoryRef, workflowNamespace string) (*wfv1.ArtifactRepositoryRefStatus, error) {
	var refs []*wfv1.ArtifactRepositoryRefStatus
	switch {
	case ref != nil && ref.Name != "":
		refs = []*wfv1.ArtifactRepositoryRefStatus{
			{Namespace: workflowNamespace, ArtifactRepositoryRef: wfv1.ArtifactRepositoryRef{Name: ref.Name}},
			{Namespace: s.namespace, ArtifactRepositoryRef: wfv1.ArtifactRepositoryRef{Name: ref.Name}},
		}
	case ref != nil:
		refs = []*wfv1.ArtifactRepositoryRefStatus{
			{Namespace: workflowNamespace, ArtifactRepositoryRef: wfv1.ArtifactRepositoryRef{ConfigMap: ref.ConfigMap, Key: ref.Key}},
			{Namespace: s.namespace, ArtifactRepositoryRef: wfv1.ArtifactRepositoryRef{ConfigMap: ref.ConfigMap, Key: ref.Key}},
		}
	default:
		refs = []*wfv1.ArtifactRepositoryRefStatus{
			{Namespace: workflowNamespace, ArtifactRepositoryRef: wfv1.ArtifactRepositoryRef{Name: DefaultName}},
			{Namespace: workflowNamespace},
			{Namespace: s.namespace, ArtifactRepositoryRef: wfv1.ArtifactRepositoryRef{Name: DefaultName}},
			{Default: true},
		}
	}
	logger := logging.RequireLoggerFromContext(ctx)
	for _, r := range refs {
		resolvedRef, err := s.get(ctx, r)
		// Forbidden: the controller or server may not yet have RBAC for WorkflowArtifactRepositories.
		if err != nil && (apierr.IsNotFound(err) || (r.Name != "" && apierr.IsForbidden(err)) || strings.Contains(err.Error(), "config map missing key")) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf(`error getting artifact repository ref "%v": %w`, r, err)
		}
		logger.WithField("artifactRepositoryRef", r).Info(ctx, "resolved artifact repository")
		switch {
		case resolvedRef.Default && resolvedRef.ArtifactRepository.Get() != nil:
			logger.Warn(ctx, "artifactRepository in the workflow-controller-configmap is deprecated, use a WorkflowArtifactRepository named default instead")
			deprecation.Record(ctx, deprecation.ControllerConfigArtifactRepository)
		case !resolvedRef.Default && resolvedRef.Name == "":
			logger.WithField("artifactRepositoryRef", r).Warn(ctx, "artifact repository config maps are deprecated, use a WorkflowArtifactRepository instead")
			deprecation.Record(ctx, deprecation.ArtifactRepositoryConfigMap)
		}
		return resolvedRef, nil
	}
	return nil, fmt.Errorf(`failed to find any artifact repository for artifact repository ref "%v"`, ref)
}

func (s *artifactRepositories) Get(ctx context.Context, ref *wfv1.ArtifactRepositoryRefStatus) (*wfv1.ArtifactRepository, error) {
	ref, err := s.get(ctx, ref)
	if err != nil {
		return nil, err
	}
	return ref.ArtifactRepository, nil
}

func (s *artifactRepositories) get(ctx context.Context, ref *wfv1.ArtifactRepositoryRefStatus) (*wfv1.ArtifactRepositoryRefStatus, error) {
	if ref.ArtifactRepository != nil {
		return ref, nil
	}
	if ref.Default {
		return &wfv1.ArtifactRepositoryRefStatus{
			ArtifactRepositoryRef: ref.ArtifactRepositoryRef,
			Namespace:             ref.Namespace,
			Default:               true,
			ArtifactRepository:    s.defaultArtifactRepository,
		}, nil
	}
	if ref.Name != "" {
		var repo *wfv1.WorkflowArtifactRepository
		err := waitutil.Backoff(retry.DefaultRetry(ctx), func() (bool, error) {
			var err error
			repo, err = s.wfClientset.ArgoprojV1alpha1().WorkflowArtifactRepositories(ref.Namespace).Get(ctx, ref.Name, metav1.GetOptions{})
			return !errorsutil.IsTransientErrQuiet(ctx, err), err
		})
		if err != nil {
			return nil, err
		}
		return &wfv1.ArtifactRepositoryRefStatus{
			Namespace:             ref.Namespace,
			ArtifactRepositoryRef: wfv1.ArtifactRepositoryRef{Name: ref.Name},
			ArtifactRepository:    &repo.Spec,
		}, nil
	}
	var cm *v1.ConfigMap
	namespace := ref.Namespace
	configMap := ref.GetConfigMapOr("artifact-repositories")
	err := waitutil.Backoff(retry.DefaultRetry(ctx), func() (bool, error) {
		var err error
		cm, err = s.kubernetesInterface.CoreV1().ConfigMaps(namespace).Get(ctx, configMap, metav1.GetOptions{})
		return !errorsutil.IsTransientErrQuiet(ctx, err), err
	})
	if err != nil {
		return nil, err
	}
	key := ref.GetKeyOr(cm.Annotations["workflows.argoproj.io/default-artifact-repository"])
	value, ok := cm.Data[key]
	if !ok {
		return nil, fmt.Errorf(`config map missing key "%s" for artifact repository ref "%v"`, key, ref)
	}
	repo := &wfv1.ArtifactRepository{}
	if err := yaml.Unmarshal([]byte(value), repo); err != nil {
		return nil, fmt.Errorf(`failed to unmarshall config map key %q for artifact repository ref "%v": %w`, key, ref, err)
	}
	// we need the fully filled out ref so we can store it in the workflow status and it will never change
	// (even if the config map default annotation is changed)
	// this means users can change the default
	return &wfv1.ArtifactRepositoryRefStatus{
		Namespace:             namespace,
		ArtifactRepositoryRef: wfv1.ArtifactRepositoryRef{ConfigMap: configMap, Key: key},
		ArtifactRepository:    repo,
	}, nil
}
