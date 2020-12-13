package webservicetest

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/go-openapi/spec"
)

type RouteCheck struct {
	Validators []RoutePredicate
}

func (c RouteCheck) Check(route restful.Route) []error {
	var results []error

	for _, predicate := range c.Validators {
		if !predicate.Matches(route) {
			results = append(results, RouteCheckError{
				Validator: predicate,
			})
		}
	}

	return results

}

type RoutePredicate struct {
	Description string
	Matches func(route restful.Route) bool
}

type RouteCheckError struct {
	Validator RoutePredicate
}

func (c RouteCheckError) Error() string {
	return fmt.Sprintf("Failed route validator: %s", c.Validator.Description)
}

func RouteHasReturnCode(code int) RoutePredicate {
	return RoutePredicate{
		Description: fmt.Sprintf("route.ResponseErrors contains %d", code),
		Matches: func(route restful.Route) bool {
			_, ok := route.ResponseErrors[code]
			return ok
		},
	}
}

func RouteHasReturnCodes(codes ...int) []RoutePredicate {
	var result []RoutePredicate
	for _, code := range codes {
		result = append(result, RouteHasReturnCode(code))
	}
	return result
}

func RouteHasDefaultReturnCode(code int) RoutePredicate {
	return RoutePredicate{
		Description: fmt.Sprintf("route.Metadata[DefaultReturnCode] == %d", code),
		Matches: func(route restful.Route) bool {
			defaultReturnCode, ok := route.Metadata["DefaultReturnCode"]
			if !ok {
				return false
			}
			return code == defaultReturnCode.(int)
		},
	}
}

func RouteHasConsumes(consumes string) RoutePredicate {
	return RoutePredicate{
		Description: fmt.Sprintf("route.Consumes contains %q", consumes),
		Matches: func(route restful.Route) bool {
			return types.StringStack(route.Consumes).Contains(consumes)
		},
	}
}

func RouteHasProduces(produces string) RoutePredicate {
	return RoutePredicate{
		Description: fmt.Sprintf("route.Produces contains %q", produces),
		Matches: func(route restful.Route) bool {
			return types.StringStack(route.Produces).Contains(produces)
		},
	}
}

func RouteHasPermissions(permissions ...string) []RoutePredicate {
	var result []RoutePredicate
	for _, permission := range permissions {
		result = append(result, RouteHasPermission(permission))
	}
	return result

}

func RouteHasPermission(permission string) RoutePredicate {
	return RoutePredicate{
		Description: fmt.Sprintf("route.Metadata[Permissions] contains %q", permission),
		Matches: func(route restful.Route) bool {
			permissionsIface, ok := route.Metadata["Permissions"]
			if !ok {
				return false
			}
			permissions := permissionsIface.([]string)
			return types.StringStack(permissions).Contains(permission)
		},
	}
}

func RouteHasTag(tag string) RoutePredicate {
	return RoutePredicate{
		Description: fmt.Sprintf("route.Metadata[TagDefinition].Name == %q", tag),
		Matches: func(route restful.Route) bool {
			tagIface, ok := route.Metadata["TagDefinition"]
			if !ok {
				return false
			}
			tagProps := tagIface.( spec.TagProps)
			return tagProps.Name == tag
		},
	}
}

func RouteHasWriteSample(payload interface{}) RoutePredicate {
	return RoutePredicate{
		Description: fmt.Sprintf("route.WriteSample == %v", payload),
		Matches: func(route restful.Route) bool {
			return route.WriteSample == payload
		},
	}
}

func RouteHasAnyWriteSample() RoutePredicate {
	return RoutePredicate{
		Description: fmt.Sprintf("route.WriteSample not nil"),
		Matches: func(route restful.Route) bool {
			return route.WriteSample != nil
		},
	}
}

func RouteHasParameter(kind int, name string) RoutePredicate {
	return RoutePredicate{
		Description: fmt.Sprintf("route.%s(%q) exists", kind, name),
		Matches: func(route restful.Route) bool {
			for _, parameter := range route.ParameterDocs {
				if parameter.Data().Name == name && parameter.Data().Kind == kind {
					return true
				}
			}
			return false
		},
	}
}
