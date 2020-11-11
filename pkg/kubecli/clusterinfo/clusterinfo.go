package clusterinfo

import (
	"fmt"
	"github.com/gookit/color"
	corev1 "k8s.io/api/core/v1"
	utilnet "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/scheme"
	"strconv"
	"strings"
)

func ClusterInfo(factory util.Factory) (string, error) {
	client, err := factory.ToRESTConfig()
	if err != nil {
		return "", err
	}
	builder := factory.NewBuilder()
	b := builder.
		WithScheme(scheme.Scheme, scheme.Scheme.PrioritizedVersionsAllGroups()...).
		NamespaceParam("kube-system").DefaultNamespace().
		LabelSelectorParam("kubernetes.io/cluster-service=true").
		ResourceTypeOrNameArgs(false, []string{"services"}...).
		Latest()

	infoArr := make([]string, 0)

	if err := b.Do().Visit(func(r *resource.Info, err error) error {
		if err != nil {
			return err
		}
		infoArr = append(infoArr, fmt.Sprintf("%s %s", color.Green.Sprint("Kubernetes control plane"), client.Host))
		services := r.Object.(*corev1.ServiceList).Items
		for _, service := range services {
			var link string
			if len(service.Status.LoadBalancer.Ingress) > 0 {
				ingress := service.Status.LoadBalancer.Ingress[0]
				ip := ingress.IP
				if ip == "" {
					ip = ingress.Hostname
				}
				for _, port := range service.Spec.Ports {
					link += "http://" + ip + ":" + strconv.Itoa(int(port.Port)) + " "
				}
			} else {
				name := service.ObjectMeta.Name

				if len(service.Spec.Ports) > 0 {
					port := service.Spec.Ports[0]

					// guess if the scheme is https
					scheme := ""
					if port.Name == "https" || port.Port == 443 {
						scheme = "https"
					}

					// format is <scheme>:<service-name>:<service-port-name>
					name = utilnet.JoinSchemeNamePort(scheme, service.ObjectMeta.Name, port.Name)
				}

				if len(client.GroupVersion.Group) == 0 {
					link = client.Host + "/api/" + client.GroupVersion.Version + "/namespaces/" + service.ObjectMeta.Namespace + "/services/" + name + "/proxy"
				} else {
					link = client.Host + "/api/" + client.GroupVersion.Group + "/" + client.GroupVersion.Version + "/namespaces/" + service.ObjectMeta.Namespace + "/services/" + name + "/proxy"

				}
			}
			name := service.ObjectMeta.Labels["kubernetes.io/name"]
			if len(name) == 0 {
				name = service.ObjectMeta.Name
			}

			infoArr = append(infoArr, fmt.Sprintf("%s %s", color.Green.Sprint(name), link))
		}
		return nil
	}); err != nil {
		return "", nil
	}
	return strings.Join(infoArr, "\n"), nil
}
