package virtuous

import (
	"sort"
)

type clientSpec struct {
	Services []clientService
}

type clientService struct {
	Name    string
	Methods []clientMethod
}

type clientMethod struct {
	Name       string
	HTTPMethod string
	Path       string
	PathParams []string
	HasBody    bool
	HasAuth    bool
	Auth       GuardSpec
}

func buildClientSpec(routes []Route) clientSpec {
	serviceMap := make(map[string]*clientService)
	for _, route := range routes {
		if route.Handler == nil {
			continue
		}
		service := route.Meta.Service
		methodName := camelizeDown(route.Meta.Method)
		if service == "" || methodName == "" {
			continue
		}
		cs, ok := serviceMap[service]
		if !ok {
			cs = &clientService{Name: service}
			serviceMap[service] = cs
		}
		hasBody := route.Handler.RequestType() != nil
		method := clientMethod{
			Name:       methodName,
			HTTPMethod: route.Method,
			Path:       route.Path,
			PathParams: route.PathParams,
			HasBody:    hasBody,
		}
		if len(route.Guards) > 0 {
			method.HasAuth = true
			method.Auth = route.Guards[0]
		}
		cs.Methods = append(cs.Methods, method)
	}

	services := make([]clientService, 0, len(serviceMap))
	for _, svc := range serviceMap {
		sort.Slice(svc.Methods, func(i, j int) bool {
			return svc.Methods[i].Name < svc.Methods[j].Name
		})
		services = append(services, *svc)
	}
	sort.Slice(services, func(i, j int) bool {
		return services[i].Name < services[j].Name
	})

	return clientSpec{Services: services}
}
