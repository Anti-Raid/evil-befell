package api

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/anti-raid/evil-befall/pkg/state"
	"github.com/go-viper/mapstructure/v2"
)

type ApiRequestFuncWithOnlyResp[RespType any] func(ctx context.Context, state *state.State) (*RespType, error)
type ApiRequestFuncWithOnlyReq[ReqType any] func(ctx context.Context, state *state.State, data ReqType) error
type ApiRequestFuncWithReqAndResp[ReqType any, RespType any] func(ctx context.Context, state *state.State, data ReqType) (*RespType, error)

type TestableRoute interface {
	ID() string
	PopulateWithArgs(args map[string]any) error
	Exec(ctx context.Context, state *state.State) (any, error)
}

func IsTestableRoute(r TestableRoute) {}

// TestableRouteWrapper is a wrapper to make TestableRoute inside functions
type TestableRouteWrapper[Data any] struct {
	FuncID               func(self *TestableRouteWrapper[Data]) string
	FuncPopulateWithArgs func(self *TestableRouteWrapper[Data], args map[string]any) error
	FuncExec             func(self *TestableRouteWrapper[Data], ctx context.Context, state *state.State) (any, error)
	Data                 Data
}

func (r *TestableRouteWrapper[Data]) ID() string {
	return r.FuncID(r)
}

func (r *TestableRouteWrapper[Data]) PopulateWithArgs(args map[string]any) error {
	return r.FuncPopulateWithArgs(r, args)
}

func (r *TestableRouteWrapper[Data]) Exec(ctx context.Context, state *state.State) (any, error) {
	return r.FuncExec(r, ctx, state)
}

func CreateTestableRouteWithOnlyResp[RespType any](id string, fn ApiRequestFuncWithOnlyResp[RespType]) TestableRoute {
	trw := &TestableRouteWrapper[struct{}]{}

	// Implement methods on trw
	trw.FuncID = func(self *TestableRouteWrapper[struct{}]) string {
		return id
	}

	trw.FuncPopulateWithArgs = func(self *TestableRouteWrapper[struct{}], args map[string]any) error {
		return nil
	}

	trw.FuncExec = func(self *TestableRouteWrapper[struct{}], ctx context.Context, state *state.State) (any, error) {
		return fn(ctx, state)
	}

	return trw
}

func CreateTestableRouteWithOnlyReq[ReqType any](id string, fn ApiRequestFuncWithOnlyReq[ReqType]) TestableRoute {
	trw := &TestableRouteWrapper[ReqType]{}

	// Implement methods on trw
	trw.FuncID = func(self *TestableRouteWrapper[ReqType]) string {
		return id
	}

	trw.FuncPopulateWithArgs = func(self *TestableRouteWrapper[ReqType], args map[string]any) error {
		// Use mapstructure to create a ReqType from args
		var reqData ReqType

		err := mapstructure.Decode(args, &reqData)

		if err != nil {
			return err
		}

		self.Data = reqData

		return nil
	}

	trw.FuncExec = func(self *TestableRouteWrapper[ReqType], ctx context.Context, state *state.State) (any, error) {
		err := fn(ctx, state, self.Data)

		if err != nil {
			return nil, err
		}

		return map[string]string{}, nil
	}

	return trw
}

func CreateTestableRouteWithReqAndResp[ReqType any, RespType any](id string, fn ApiRequestFuncWithReqAndResp[ReqType, RespType]) TestableRoute {
	trw := &TestableRouteWrapper[ReqType]{}

	// Implement methods on trw
	trw.FuncID = func(self *TestableRouteWrapper[ReqType]) string {
		return id
	}

	trw.FuncPopulateWithArgs = func(self *TestableRouteWrapper[ReqType], args map[string]any) error {
		// Use mapstructure to create a ReqType from args
		var reqData ReqType

		err := mapstructure.Decode(args, &reqData)

		if err != nil {
			return err
		}

		self.Data = reqData

		return nil
	}

	trw.FuncExec = func(self *TestableRouteWrapper[ReqType], ctx context.Context, state *state.State) (any, error) {
		return fn(ctx, state, self.Data)
	}

	return trw
}

func init() {
	// Asset that TestableRouteWrapper implements TestableRoute
	IsTestableRoute(&TestableRouteWrapper[struct{}]{})
}

// API test registry
var testableRoutes = map[string]TestableRoute{}

func RegisterTestableRoute(r TestableRoute) {
	testableRoutes[r.ID()] = r
}

func CreateAndRegisterTestableRouteWithReqAndResp[ReqType any, RespType any](id string, fn ApiRequestFuncWithReqAndResp[ReqType, RespType]) {
	RegisterTestableRoute(CreateTestableRouteWithReqAndResp(id, fn))
}

func CreateAndRegisterTestableRouteWithOnlyReq[ReqType any](id string, fn ApiRequestFuncWithOnlyReq[ReqType]) {
	RegisterTestableRoute(CreateTestableRouteWithOnlyReq(id, fn))
}

func CreateAndRegisterTestableRouteWithOnlyResp[RespType any](id string, fn ApiRequestFuncWithOnlyResp[RespType]) {
	RegisterTestableRoute(CreateTestableRouteWithOnlyResp(id, fn))
}

// Other utilities
func StructToQueryParamsList(s any) map[string]any {
	// Reflect to get fields
	refType := reflect.TypeOf(s)

	var cols = map[string]any{}

	for _, f := range reflect.VisibleFields(refType) {
		jsonTag := f.Tag.Get("json")
		reflectOpts := f.Tag.Get("reflect")

		if !strings.HasPrefix(jsonTag, "query:") || reflectOpts == "ignore" {
			continue
		}

		// Get the value of the field
		val := reflect.ValueOf(s).FieldByName(f.Name).Interface()

		// Add to cols
		cols[strings.Split(jsonTag, "query:")[0]] = val
	}

	return cols
}

func QueryParamsListToString(qp map[string]any) string {
	var parts []string

	for k, v := range qp {
		parts = append(parts, k+"="+fmt.Sprint(v))
	}

	if len(parts) == 0 {
		return ""
	}

	return "?" + strings.Join(parts, "&")
}

func StructToQueryParamsString(s any) string {
	return QueryParamsListToString(StructToQueryParamsList(s))
}
