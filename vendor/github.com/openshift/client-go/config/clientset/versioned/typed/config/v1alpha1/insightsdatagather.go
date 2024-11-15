// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"

	v1alpha1 "github.com/openshift/api/config/v1alpha1"
	configv1alpha1 "github.com/openshift/client-go/config/applyconfigurations/config/v1alpha1"
	scheme "github.com/openshift/client-go/config/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	gentype "k8s.io/client-go/gentype"
)

// InsightsDataGathersGetter has a method to return a InsightsDataGatherInterface.
// A group's client should implement this interface.
type InsightsDataGathersGetter interface {
	InsightsDataGathers() InsightsDataGatherInterface
}

// InsightsDataGatherInterface has methods to work with InsightsDataGather resources.
type InsightsDataGatherInterface interface {
	Create(ctx context.Context, insightsDataGather *v1alpha1.InsightsDataGather, opts v1.CreateOptions) (*v1alpha1.InsightsDataGather, error)
	Update(ctx context.Context, insightsDataGather *v1alpha1.InsightsDataGather, opts v1.UpdateOptions) (*v1alpha1.InsightsDataGather, error)
	// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
	UpdateStatus(ctx context.Context, insightsDataGather *v1alpha1.InsightsDataGather, opts v1.UpdateOptions) (*v1alpha1.InsightsDataGather, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.InsightsDataGather, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.InsightsDataGatherList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.InsightsDataGather, err error)
	Apply(ctx context.Context, insightsDataGather *configv1alpha1.InsightsDataGatherApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.InsightsDataGather, err error)
	// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
	ApplyStatus(ctx context.Context, insightsDataGather *configv1alpha1.InsightsDataGatherApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.InsightsDataGather, err error)
	InsightsDataGatherExpansion
}

// insightsDataGathers implements InsightsDataGatherInterface
type insightsDataGathers struct {
	*gentype.ClientWithListAndApply[*v1alpha1.InsightsDataGather, *v1alpha1.InsightsDataGatherList, *configv1alpha1.InsightsDataGatherApplyConfiguration]
}

// newInsightsDataGathers returns a InsightsDataGathers
func newInsightsDataGathers(c *ConfigV1alpha1Client) *insightsDataGathers {
	return &insightsDataGathers{
		gentype.NewClientWithListAndApply[*v1alpha1.InsightsDataGather, *v1alpha1.InsightsDataGatherList, *configv1alpha1.InsightsDataGatherApplyConfiguration](
			"insightsdatagathers",
			c.RESTClient(),
			scheme.ParameterCodec,
			"",
			func() *v1alpha1.InsightsDataGather { return &v1alpha1.InsightsDataGather{} },
			func() *v1alpha1.InsightsDataGatherList { return &v1alpha1.InsightsDataGatherList{} }),
	}
}