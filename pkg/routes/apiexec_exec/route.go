package apiexec_exec

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/anti-raid/evil-befall/pkg/api"
	"github.com/anti-raid/evil-befall/pkg/state"
	"github.com/anti-raid/spintrack/structstring"
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

type ApiExecExecRoute struct {
}

func (r *ApiExecExecRoute) Command() string {
	return "apiexec.exec"
}

func (r *ApiExecExecRoute) Description() string {
	return "Execute/Make a request to an API endpoint that can be tested"
}

func (r *ApiExecExecRoute) Arguments() [][3]string {
	return [][3]string{
		{"route", "Show detailed information about a specific route", "string"},
		{"debug", "Print debug information", "bool"},
	}
}

func (r *ApiExecExecRoute) Setup(state *state.State) error {
	return nil
}

func (r *ApiExecExecRoute) Destroy(state *state.State) error {
	return nil
}

func (r *ApiExecExecRoute) Render(state *state.State, args map[string]string) error {
	if debug, ok := args["__debug"]; ok && debug == "true" {
		for k, v := range args {
			fmt.Println(k, v, []byte(v))
		}
	}

	show, ok := args["route"]

	if !ok {
		return fmt.Errorf("no route specified")
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

	mkMap := make(map[string]any)
	for k, v := range args {
		if k == "route" || strings.HasPrefix(k, "__") {
			continue
		}

		// Handle types
		kSplit := strings.Split(k, "::")

		if len(kSplit) != 2 {
			kSplit = append(kSplit, "string")
		}

		setKey := kSplit[0]
		keyTyp := kSplit[1]

		err := setValue(setKey, keyTyp, v, mkMap)

		if err != nil {
			return err
		}
	}

	fmt.Println("Route ID:", route.ID())
	fmt.Println("Route Req Send:")
	fmt.Println(structstring.SpewStruct(mkMap))

	// Create the reqtype
	route, err := route.PopulateWithArgs(mkMap)

	if err != nil {
		return fmt.Errorf("failed to populate route with args: %w", err)
	}

	fmt.Println(structstring.SpewStruct(route))

	// Send the request
	resp, err := route.Exec(context.TODO(), state)

	if err != nil {
		return fmt.Errorf("failed to execute route: %w", err)
	}

	fmt.Println("Route Resp Recv:")

	// Print the response
	if spewResp, ok := args["__spew"]; ok && spewResp == "true" {
		fmt.Println(structstring.SpewStruct(resp))
		return nil
	}

	// Convert to JSON
	respJSON, err := json.Marshal(resp)

	if err != nil {
		return fmt.Errorf("failed to convert response to JSON: %w", err)
	}

	fmt.Println(string(respJSON))

	return nil
}

// Format for KV's are as follows:
//
// KEY::TYPE=VALUE for normal values
// For inputting raw JSON, typ is JSON and value is a json value
//
// Note that array support is pretty lacking and so using raw JSON is recommended for arrays
func setValue(key, typ, v string, setMap map[string]any) error {
	if strings.HasPrefix("[]", typ) {
		// Handle array types
		typ = strings.TrimPrefix(typ, "[]")
		vals := strings.Split(v, ",")

		var arr []any

		for _, val := range vals {
			val, err := parseValueImpl(key, typ, val)

			if err != nil {
				return err
			}

			arr = append(arr, val)
		}

		setMap[key] = arr
		return nil
	}

	if strings.ToLower(typ) == "json" {
		var patch any

		err := json.Unmarshal([]byte(v), &patch)

		if err != nil {
			return fmt.Errorf("failed to parse %s=%s as JSON: %w", key, v, err)
		}

		setMap[key] = patch
		return nil
	}

	val, err := parseValueImpl(key, typ, v)

	if err != nil {
		return err
	}

	setMap[key] = val
	return nil
}

func parseValueImpl(key, typ, v string) (any, error) {
	switch typ {
	// Unsigned int types
	case "uint":
		uintVal, err := strconv.ParseUint(v, 10, strconv.IntSize)

		if err != nil {
			return nil, fmt.Errorf("failed to convert %s=%s to uint: %w", key, v, err)
		}

		return uintVal, nil
	case "uint8":
		uintVal, err := strconv.ParseUint(v, 10, 8)

		if err != nil {
			return nil, fmt.Errorf("failed to convert %s=%s to uint8: %w", key, v, err)
		}

		return uint8(uintVal), nil
	case "uint16":
		uintVal, err := strconv.ParseUint(v, 10, 16)

		if err != nil {
			return nil, fmt.Errorf("failed to convert %s=%s to uint16: %w", key, v, err)
		}

		return uint16(uintVal), nil
	case "uint32":
		uintVal, err := strconv.ParseUint(v, 10, 32)

		if err != nil {
			return nil, fmt.Errorf("failed to convert %s=%s to uint32: %w", key, v, err)
		}

		return uint32(uintVal), nil
	case "uint64":
		uintVal, err := strconv.ParseUint(v, 10, 64)

		if err != nil {
			return nil, fmt.Errorf("failed to convert %s=%s to uint64: %w", key, v, err)
		}

		return uintVal, nil
	case "uintptr":
		uintVal, err := strconv.ParseUint(v, 10, strconv.IntSize)

		if err != nil {
			return nil, fmt.Errorf("failed to convert %s=%s to uintptr: %w", key, v, err)
		}

		return uintptr(uintVal), nil
	case "byte":
		uintVal, err := strconv.ParseUint(v, 10, 8)

		if err != nil {
			return nil, fmt.Errorf("failed to convert %s=%s to byte: %w", key, v, err)
		}

		return byte(uintVal), nil
	// Signed int types
	case "int":
		intVal, err := strconv.ParseInt(v, 10, strconv.IntSize)

		if err != nil {
			return nil, fmt.Errorf("failed to convert %s=%s to int: %w", key, v, err)
		}

		return intVal, nil
	case "int8":
		intVal, err := strconv.ParseInt(v, 10, 8)

		if err != nil {
			return nil, fmt.Errorf("failed to convert %s=%s to int8: %w", key, v, err)
		}

		return int8(intVal), nil
	case "int16":
		intVal, err := strconv.ParseInt(v, 10, 16)

		if err != nil {
			return nil, fmt.Errorf("failed to convert %s=%s to int16: %w", key, v, err)
		}

		return int16(intVal), nil
	case "int32":
		intVal, err := strconv.ParseInt(v, 10, 32)

		if err != nil {
			return nil, fmt.Errorf("failed to convert %s=%s to int32: %w", key, v, err)
		}

		return int32(intVal), nil
	case "int64":
		intVal, err := strconv.ParseInt(v, 10, 64)

		if err != nil {
			return nil, fmt.Errorf("failed to convert %s=%s to int64: %w", key, v, err)
		}

		return intVal, nil
	default:
		return v, nil // Just set it as a string/default type
	}
}
