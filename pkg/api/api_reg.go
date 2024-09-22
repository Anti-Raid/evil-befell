package api

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/anti-raid/evil-befall/pkg/state"
)

type ApiRequestFuncWithOnlyResp[RespType any] func(ctx context.Context, state *state.State) (*RespType, error)
type ApiRequestFuncWithOnlyReq[ReqType any] func(ctx context.Context, state *state.State, data ReqType) error
type ApiRequestFuncWithReqAndResp[ReqType any, RespType any] func(ctx context.Context, state *state.State, data ReqType) (*RespType, error)

type TestableRoute interface {
	ID() string
	PopulateWithArgs(args map[string]any) (TestableRoute, error)
	ReqType() any
	RespType() any
	Exec(ctx context.Context, state *state.State) (any, error)
}

func IsTestableRoute(r TestableRoute) {}

// TestableRouteWrapper is a wrapper to make TestableRoute inside functions
type TestableRouteWrapper[Data any] struct {
	FuncID               func(self *TestableRouteWrapper[Data]) string
	FuncPopulateWithArgs func(self *TestableRouteWrapper[Data], args map[string]any) (TestableRoute, error)
	FuncReqType          func(self *TestableRouteWrapper[Data]) any
	FuncRespType         func(self *TestableRouteWrapper[Data]) any
	FuncExec             func(self *TestableRouteWrapper[Data], ctx context.Context, state *state.State) (any, error)
	Data                 Data
}

func (r *TestableRouteWrapper[Data]) ID() string {
	return r.FuncID(r)
}

func (r *TestableRouteWrapper[Data]) PopulateWithArgs(args map[string]any) (TestableRoute, error) {
	return r.FuncPopulateWithArgs(r, args)
}

func (r *TestableRouteWrapper[Data]) ReqType() any {
	return r.FuncReqType(r)
}

func (r *TestableRouteWrapper[Data]) RespType() any {
	return r.FuncRespType(r)
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

	trw.FuncPopulateWithArgs = func(self *TestableRouteWrapper[struct{}], args map[string]any) (TestableRoute, error) {
		return self, nil
	}

	trw.FuncReqType = func(self *TestableRouteWrapper[struct{}]) any {
		return struct{}{}
	}

	trw.FuncRespType = func(self *TestableRouteWrapper[struct{}]) any {
		var respType RespType
		return respType
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

	trw.FuncPopulateWithArgs = func(self *TestableRouteWrapper[ReqType], args map[string]any) (TestableRoute, error) {
		// Use json Marshal/Unmarshal to create a ReqType from args
		var reqData ReqType

		reqBytes, err := json.Marshal(args)

		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(reqBytes, &reqData)

		if err != nil {
			return nil, err
		}

		// Copy self
		return &TestableRouteWrapper[ReqType]{
			FuncID:               self.FuncID,
			FuncPopulateWithArgs: self.FuncPopulateWithArgs,
			FuncReqType:          self.FuncReqType,
			FuncRespType:         self.FuncRespType,
			FuncExec:             self.FuncExec,
			Data:                 reqData,
		}, nil
	}

	trw.FuncReqType = func(self *TestableRouteWrapper[ReqType]) any {
		var reqType ReqType
		return reqType
	}

	trw.FuncRespType = func(self *TestableRouteWrapper[ReqType]) any {
		return struct{}{}
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

	trw.FuncPopulateWithArgs = func(self *TestableRouteWrapper[ReqType], args map[string]any) (TestableRoute, error) {
		// Use json Marshal/Unmarshal to create a ReqType from args
		var reqData ReqType

		reqBytes, err := json.Marshal(args)

		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(reqBytes, &reqData)

		if err != nil {
			return nil, err
		}

		// Copy self
		return &TestableRouteWrapper[ReqType]{
			FuncID:               self.FuncID,
			FuncPopulateWithArgs: self.FuncPopulateWithArgs,
			FuncReqType:          self.FuncReqType,
			FuncRespType:         self.FuncRespType,
			FuncExec:             self.FuncExec,
			Data:                 reqData,
		}, nil
	}

	trw.FuncReqType = func(self *TestableRouteWrapper[ReqType]) any {
		var reqType ReqType
		return reqType
	}

	trw.FuncRespType = func(self *TestableRouteWrapper[ReqType]) any {
		var respType RespType
		return respType
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

type TestableRouteCategory struct {
	Name   string
	Routes []TestableRoute
}

func NewTestableRouteCategory(name string, routes ...TestableRoute) TestableRouteCategory {
	return TestableRouteCategory{
		Name:   name,
		Routes: routes,
	}
}

// API test registry
var testableRoutesCategory = []TestableRouteCategory{}

// RegisterTestableRouteCategory registers a new TestableRouteCategory
func RegisterTestableRouteCategory(r TestableRouteCategory) {
	testableRoutesCategory = append(testableRoutesCategory, r)
}

// GetTestableRouteCategories returns all registered TestableRouteCategory
func GetTestableRouteCategories() []TestableRouteCategory {
	return testableRoutesCategory
}

// GetTestableRoutes returns all registered TestableRoute in a flat list
func GetTestableRoutes() []TestableRoute {
	var testableRoutes []TestableRoute

	for _, category := range testableRoutesCategory {
		testableRoutes = append(testableRoutes, category.Routes...)
	}

	return testableRoutes
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
