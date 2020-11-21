package kubecli

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (cli *KubeCLI) GetResourceGroupVersionKind(resourceOrKindArg string) schema.GroupVersionKind {
	fullySpecifiedGVR, groupResource := schema.ParseResourceArg(resourceOrKindArg)
	gvk := schema.GroupVersionKind{}

	restMapper, err := cli.factory.ToRESTMapper()
	if err != nil {
		return gvk
	}

	if fullySpecifiedGVR != nil {
		gvk, _ = restMapper.KindFor(*fullySpecifiedGVR)
	}
	if gvk.Empty() {
		gvk, _ = restMapper.KindFor(groupResource.WithVersion(""))
	}
	return gvk
}
