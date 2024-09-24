package apiexec_exec

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/anti-raid/evil-befall/pkg/api"
	"github.com/anti-raid/evil-befall/pkg/state"
	"github.com/anti-raid/shellcli/shell"
	"github.com/anti-raid/spintrack/structstring"
	"github.com/anti-raid/spintrack/strutils"
	"github.com/go-andiamo/splitter"
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
		{"__debug", "Print debug information", "bool"},
		{"__spew.req", "Spew the request", "bool"},
		{"__spew.resp", "Spew the response", "bool"},
		{"__file", "Write the response to a file", "string"},
		{"__file.mode", "File mode (json, spew)", "string"},
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

	// Print the request
	if spewReq, ok := args["__spew.req"]; ok && spewReq == "true" {
		fmt.Println(structstring.SpewStruct(route))
	}

	// Send the request
	resp, err := route.Exec(context.TODO(), state)

	if err != nil {
		return fmt.Errorf("failed to execute route: %w", err)
	}

	fmt.Println("Route Resp Recv:")

	// Print the response
	if spewResp, ok := args["__spew.resp"]; ok && spewResp == "true" {
		fmt.Println(structstring.SpewStruct(resp))
		return nil
	}

	// If __file is set, write to file
	if file, ok := args["__file"]; ok && file != "" {
		f, err := os.Create(file)

		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", file, err)
		}

		defer f.Close()

		mode, ok := args["__file.mode"]

		if !ok {
			mode = "json"
		}

		switch mode {
		case "json":
			enc := json.NewEncoder(f)

			err = enc.Encode(resp)

			if err != nil {
				return fmt.Errorf("failed to encode response to file: %w", err)
			}

			fmt.Println("JSON response written to file:", file)
		case "spew":
			_, err = f.WriteString(structstring.SpewStruct(resp))

			if err != nil {
				return fmt.Errorf("failed to write response to file: %w", err)
			}

			fmt.Println("Spew response written to file:", file)
		default:
			return fmt.Errorf("unsupported mode %s", mode)
		}

		return nil
	}

	// Otherwise, convert to JSON
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
// # Note that array support is pretty lacking and so using raw JSON is recommended for arrays
//
// Arrays: KEY::[]{SEP}[VAL1SEPVAL2SEPVAL3]
func setValue(key, typ, v string, setMap map[string]any) error {
	// Array support
	if strings.HasPrefix(typ, "[]") {
		// Handle array types
		typ = strings.TrimPrefix(typ, "[]")
		// All chars between {} are the separator
		if typ[0] != '{' {
			return fmt.Errorf("invalid array type %s", typ)
		}

		var sep string

		// Keep going from typ[1] until we hit a }
		var gotSep bool
		for i := 1; i < len(typ); i++ {
			if typ[i] == '}' {
				sep = typ[1:i]
				typ = typ[i+1:]
				gotSep = true
				break
			}
		}

		if !gotSep {
			return fmt.Errorf("invalid array type %s", typ)
		}

		if len(sep) > 1 {
			return fmt.Errorf("only single character separators are supported")
		}

		var err error
		parseArraySplitter, err := splitter.NewSplitter(rune(sep[0]), splitter.DoubleQuotesBackSlashEscaped, splitter.SingleQuotesBackSlashEscaped)

		if err != nil {
			panic("error initializing array tokenizer: " + err.Error())
		}

		vals, err := parseArraySplitter.Split(v)

		if err != nil {
			return fmt.Errorf("failed to split array %s: %w", v, err)
		}

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

	// Base64 URL support
	if strings.HasPrefix(typ, "[b64url]") {
		decoder := base64.NewDecoder(base64.URLEncoding, strings.NewReader(v))

		decoded, err := io.ReadAll(decoder)

		if err != nil {
			return fmt.Errorf("failed to decode base64url: %w", err)
		}

		return setValue(key, strings.TrimPrefix(typ, "[b64url]"), string(decoded), setMap)
	}

	val, err := parseValueImpl(key, typ, v)

	if err != nil {
		return err
	}

	setMap[key] = val
	return nil
}

func parseValueImpl(key, typ, v string) (any, error) {
	switch strings.ToLower(typ) {
	// Unsigned int types
	case "uint":
		uintVal, err := strconv.ParseUint(strings.TrimSpace(v), 10, strconv.IntSize)

		if err != nil {
			return nil, fmt.Errorf("failed to convert %s=%s to uint: %w", key, v, err)
		}

		return uintVal, nil
	case "uint8":
		uintVal, err := strconv.ParseUint(strings.TrimSpace(v), 10, 8)

		if err != nil {
			return nil, fmt.Errorf("failed to convert %s=%s to uint8: %w", key, v, err)
		}

		return uint8(uintVal), nil
	case "uint16":
		uintVal, err := strconv.ParseUint(strings.TrimSpace(v), 10, 16)

		if err != nil {
			return nil, fmt.Errorf("failed to convert %s=%s to uint16: %w", key, v, err)
		}

		return uint16(uintVal), nil
	case "uint32":
		uintVal, err := strconv.ParseUint(strings.TrimSpace(v), 10, 32)

		if err != nil {
			return nil, fmt.Errorf("failed to convert %s=%s to uint32: %w", key, v, err)
		}

		return uint32(uintVal), nil
	case "uint64":
		uintVal, err := strconv.ParseUint(strings.TrimSpace(v), 10, 64)

		if err != nil {
			return nil, fmt.Errorf("failed to convert %s=%s to uint64: %w", key, v, err)
		}

		return uintVal, nil
	case "uintptr":
		uintVal, err := strconv.ParseUint(strings.TrimSpace(v), 10, strconv.IntSize)

		if err != nil {
			return nil, fmt.Errorf("failed to convert %s=%s to uintptr: %w", key, v, err)
		}

		return uintptr(uintVal), nil
	case "byte":
		uintVal, err := strconv.ParseUint(strings.TrimSpace(v), 10, 8)

		if err != nil {
			return nil, fmt.Errorf("failed to convert %s=%s to byte: %w", key, v, err)
		}

		return byte(uintVal), nil
	// Signed int types
	case "int":
		intVal, err := strconv.ParseInt(strings.TrimSpace(v), 10, strconv.IntSize)

		if err != nil {
			return nil, fmt.Errorf("failed to convert %s=%s to int: %w", key, v, err)
		}

		return intVal, nil
	case "int8":
		intVal, err := strconv.ParseInt(strings.TrimSpace(v), 10, 8)

		if err != nil {
			return nil, fmt.Errorf("failed to convert %s=%s to int8: %w", key, v, err)
		}

		return int8(intVal), nil
	case "int16":
		intVal, err := strconv.ParseInt(strings.TrimSpace(v), 10, 16)

		if err != nil {
			return nil, fmt.Errorf("failed to convert %s=%s to int16: %w", key, v, err)
		}

		return int16(intVal), nil
	case "int32":
		intVal, err := strconv.ParseInt(strings.TrimSpace(v), 10, 32)

		if err != nil {
			return nil, fmt.Errorf("failed to convert %s=%s to int32: %w", key, v, err)
		}

		return int32(intVal), nil
	case "int64":
		intVal, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)

		if err != nil {
			return nil, fmt.Errorf("failed to convert %s=%s to int64: %w", key, v, err)
		}

		return intVal, nil
	// Floating point types
	case "float32":
		floatVal, err := strconv.ParseFloat(strings.TrimSpace(v), 32)

		if err != nil {
			return nil, fmt.Errorf("failed to convert %s=%s to float32: %w", key, v, err)
		}

		return float32(floatVal), nil
	case "float64":
		floatVal, err := strconv.ParseFloat(strings.TrimSpace(v), 64)

		if err != nil {
			return nil, fmt.Errorf("failed to convert %s=%s to float64: %w", key, v, err)
		}

		return floatVal, nil
	// Other types
	case "bool", "boolean":
		boolVal, err := strconv.ParseBool(strings.TrimSpace(v))

		if err != nil {
			return nil, fmt.Errorf("failed to convert %s=%s to bool: %w", key, v, err)
		}

		return boolVal, nil
	case "json":
		var patch any

		err := json.Unmarshal([]byte(v), &patch)

		if err != nil {
			return nil, fmt.Errorf("failed to parse %s=%s as JSON: %w", key, v, err)
		}

		return patch, nil
	default:
		return v, nil // Just set it as a string/default type
	}
}

func (r *ApiExecExecRoute) Completion(state *state.State, line string, args map[string]string) ([]string, error) {
	route, ok := args["route"]

	// Case #1: No route
	if !ok || route == "" {
		routes := make([]string, 0, len(api.GetTestableRoutes()))

		for _, route := range api.GetTestableRoutes() {
			routes = append(routes, "apiexec.exec "+route.ID())
		}
		return routes, nil
	}

	routeStr := strings.TrimSpace(strings.ToLower(route))

	var completions = []string{}
	var completionRoutes = []api.TestableRoute{}

	for _, route := range api.GetTestableRoutes() {
		if strings.HasPrefix(strings.ToLower(route.ID()), routeStr) {
			completions = append(completions, "apiexec.exec "+route.ID())

			if strings.ToLower(route.ID()) == routeStr {
				completionRoutes = append(completionRoutes, route)
				break
			}
		}
	}

	if len(completionRoutes) == 1 {
		// Case #2: Only one completion means we have gotten to a full non-partial route
		// Move on to stage 2 completions
		return r.stage2Completion(line, args, completionRoutes[0])
	}

	return completions, nil
}

// Stage 2 completion occurs when the user has fully typed out a route, at this point, we return the whole line combined with request options in reqType
func (r *ApiExecExecRoute) stage2Completion(line string, args map[string]string, route api.TestableRoute) (c []string, err error) {
	// Case 1: In the middle of typing out an argument
	argsStr := strings.Replace(line, "apiexec.exec "+route.ID(), "", 1)

	// Check if the user is at an '=' sign. This means that we should not provide completions at all as they want to type out a value
	lastArg := shell.UtilFindLastArgInArgStr(argsStr)

	if strings.HasSuffix(lastArg, "=") {
		return
	}

	reqType := route.ReqType()

	if reqType == nil {
		return
	}

	structFields := structstring.StructFields(reqType, structstring.StructFieldsConfig{
		FieldFilter: func(f reflect.StructField) (*string, bool) {
			typeOverride := ""

			typ := f.Type
			flag := false
			for !flag {
				switch typ.Kind() {
				case reflect.Ptr:
					typ = typ.Elem()
					continue
				case reflect.Bool:
					typeOverride = "bool"
					flag = true
				case reflect.Struct, reflect.Interface, reflect.Map, reflect.Slice, reflect.Array:
					typeOverride = "json"
					flag = true
				default:
					flag = true
				}
			}

			jsonTag := f.Tag.Get("json")

			candidateName := jsonTag

			if typeOverride != "" {
				candidateName += "::" + typeOverride
			}

			return &candidateName, jsonTag != "" && jsonTag != "-"
		},
	})

	// Look for an untyped arg, args are in format a=b
	untypedArg := shell.UtilFindUntypedArgInArgStr(argsStr)

	if untypedArg != "" {
		// Find all fields that start with the untyped arg for completion
		for _, candidate := range structFields {
			if !strings.HasPrefix(candidate, untypedArg) {
				continue
			}

			c = append(c, strings.TrimSpace(strutils.ReplaceFromBack(line, untypedArg, "", 1))+" "+candidate+"=")
		}
		return
	}

	for _, candidate := range structFields {
		if _, ok := args[candidate]; ok {
			continue // Skip if already set
		}

		c = append(c, strings.TrimSpace(line)+" "+candidate+"=")
	}

	return
}
