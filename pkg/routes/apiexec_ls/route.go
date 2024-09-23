package apiexec_ls

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/anti-raid/evil-befall/pkg/api"
	"github.com/anti-raid/evil-befall/pkg/state"
	"github.com/anti-raid/spintrack/structstring"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var ssCfg = structstring.NewDefaultConvertStructToStringConfig()

func init() {
	ssCfg.StructRecurseOverride = func(t reflect.Type) (*string, bool) {
		switch t.PkgPath() {
		case "github.com/wk8/go-ordered-map/v2":
			omapVal := "orderedmap.OrderedMap"
			return &omapVal, true
		}

		return structstring.BaseStructRecurseOverride(t)
	}
}

type ApiExecLsRoute struct {
}

func (r *ApiExecLsRoute) Command() string {
	return "apiexec.ls"
}

func (r *ApiExecLsRoute) Description() string {
	return "Lists all available API endpoints that can be tested"
}

func (r *ApiExecLsRoute) Arguments() [][3]string {
	return [][3]string{
		{"route", "Show detailed information about a specific route", "string"},
	}
}

func (r *ApiExecLsRoute) Setup(state *state.State) error {
	return nil
}

func (r *ApiExecLsRoute) Destroy(state *state.State) error {
	return nil
}

func (r *ApiExecLsRoute) Render(state *state.State, args map[string]string) error {
	show, ok := args["route"]

	if !ok {
		// Print all route ID's only
		for _, cat := range api.GetTestableRouteCategories() {
			fmt.Println(cases.Title(language.English).String(cat.Name))

			// Print 2x = for each character in the category name
			var eqs = ""

			for i := 0; i < len(cat.Name); i++ {
				eqs += "=="
			}

			fmt.Println(eqs)

			for _, route := range cat.Routes {
				fmt.Println(route.ID())
			}

			fmt.Println()
		}

		return nil
	}

	// Print detailed information about a specific route
	var route api.TestableRoute

	for _, r := range api.GetTestableRoutes() {
		if r.ID() == show {
			route = r
			break
		}
	}

	if route == nil {
		return fmt.Errorf("route %s not found", show)
	}

	fmt.Println("Route ID:", route.ID())
	fmt.Println("Route ReqType:")
	fmt.Println(structstring.ConvertStructToString(route.ReqType(), ssCfg))
	fmt.Println("Route RespType:")
	fmt.Println(structstring.ConvertStructToString(route.RespType(), ssCfg))

	return nil
}

func (r *ApiExecLsRoute) Completion(state *state.State, line string, args map[string]string) ([]string, error) {
	route, ok := args["route"]
	if !ok || route == "" {
		routes := make([]string, 0, len(api.GetTestableRoutes()))

		for _, route := range api.GetTestableRoutes() {
			routes = append(routes, "apiexec.ls "+route.ID())
		}
		return routes, nil
	}

	routeStr := strings.TrimSpace(strings.ToLower(route))

	var completions = []string{}

	for _, route := range api.GetTestableRoutes() {
		if strings.HasPrefix(strings.ToLower(route.ID()), routeStr) {
			completions = append(completions, "apiexec.ls "+route.ID())
		}
	}

	return completions, nil
}
