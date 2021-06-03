package gindoc

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz/openapi"
)

const ctxOpenAPIOperation = "_ctx_openapi_operation"

const CtxInputBindName = "_ctx_gindoc_input_bind_name"

// Primitive type helpers.
var (
	Integer  int32
	Long     int64
	Float    float32
	Double   float64
	String   string
	Byte     []byte
	Binary   []byte
	Boolean  bool
	DateTime time.Time
)

// GinDoc is an abstraction of a Gin engine that wraps the
// routes handlers with Tonic and generates an OpenAPI
// 3.0 specification from it.
type GinDoc struct {
	doc    *openapi3.T
	engine *gin.Engine
	*RouterGroup
}

// RouterGroup is an abstraction of a Gin router group.
type RouterGroup struct {
	group  *gin.RouterGroup
	engine *gin.Engine
	tags   openapi3.Tags

	Name        string
	Description string
}

// New creates a new GinDoc wrapper for
// a default Gin engine.
func New() *GinDoc {
	return NewFromEngine(gin.New())
}

// NewFromEngine creates a new GinDoc wrapper
// from an existing Gin engine.
func NewFromEngine(e *gin.Engine) *GinDoc {
	doc := &openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:   "API",
			Version: "1.0",
		},
	}
	return &GinDoc{
		engine: e,
		doc:    doc,
		RouterGroup: &RouterGroup{
			group: &e.RouterGroup,
			doc:   doc,
		},
	}
}

// ServeHTTP implements http.HandlerFunc for GinDoc.
func (g *GinDoc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.engine.ServeHTTP(w, r)
}

// Engine returns the underlying Gin engine.
func (g *GinDoc) Engine() *gin.Engine {
	return g.engine
}

func (g *GinDoc) DocumentInfo(info *openapi3.Info) {
	g.doc.Info = info
}

func (g *GinDoc) Document() *openapi3.T {
	return g.doc
}

func (g *GinDoc) OpenAPIHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, g.doc)
	}
}

// // Generator returns the underlying OpenAPI generator.
// func (g *GinDoc) Generate() *openapi3gen.Generator {
// 	return g.gen
// }

// // Errors returns the errors that may have occurred
// // during the spec generation.
// func (g *GinDoc) Errors() []error {
// 	return g.gen.Errors()
// }

// Group creates a new group of routes.
func (g *RouterGroup) Group(path string, tag *openapi3.Tag, handlers ...gin.HandlerFunc) *RouterGroup {
	// Create the tag in the specification
	// for this groups.
	if tag != nil {
		if len(g.tags) == 0 {
			g.tags = []*openapi3.Tag{}
		}
		g.tags = append([]*openapi3.Tag(g.tags), tag)
	}

	return &RouterGroup{
		tags:  g.tags,
		group: g.group.Group(path, handlers...),
	}
}

// Use adds middleware to the group.
func (g *RouterGroup) Use(handlers ...gin.HandlerFunc) {
	g.group.Use(handlers...)
}

type PathInfo struct {
	Summary        string
	Description    string
	RequestInfo    interface{}
	ResponseHeader interface{}
	ResponseBody   interface{}

	SecurityRequirements *openapi3.SecurityRequirements
	ExampleRequest       string
	ExampleResponse      string
}

type HandleInfo func(*openapi3.PathItem)

// Summary adds a summary to a path.
func Summary(summary string) func(*openapi3.Operation) {
	return func(o *openapi3.Operation) {
		o.Summary = summary
	}
}

// Description adds a description to a path.
func Description(desc string) func(*openapi3.Operation) {
	return func(o *openapi3.Operation) {
		o.Description = desc
	}
}

// Deprecated marks the operation as deprecated.
func Deprecated(deprecated bool) func(*openapi3.Operation) {
	return func(o *openapi3.Operation) {
		o.Deprecated = deprecated
	}
}


// ResponseWithExamples is a variant of Response that accept many examples.
func ResponseWithExamples(statusCode int, desc string, model interface{}, headers openapi3.Headers, examples map[string]interface{}) func(*openapi3.Operation) {
	return func(o *openapi3.Operation) {
		if len(o.Responses) == 0 {
			o.Responses = openapi3.NewResponses()
		}
		resp := openapi3.NewRes..;;;;;;;';]ponse()
		resp.
		resp.WithDescription(desc)
		openapi3.Response{
			ExtensionProps: openapi3.ExtensionProps{},
			Description:    desc,
			Headers:        headers,
			Content:        nil,
			Links:          nil,
		}
		o.AddResponse(statusCode)
		o.Responses = append([]*openapi3.Response(o.Responses), &openapi3.Responses{
			Code:        statusCode,
			Description: desc,
			Model:       model,
			Headers:     headers,
			Examples:    examples,
		})
	}
}

// GET is a shortcut to register a new handler with the GET method.
func (g *RouterGroup) GET(path string, item *openapi3.PathItem, handlers ...gin.HandlerFunc) *RouterGroup {
	return g.Handle(path, "GET", item, handlers...)
}

// GET is a shortcut to register a new handler with the GET method.
func (g *RouterGroup) GET(path string, infos []OperationOption, handlers ...gin.HandlerFunc) *RouterGroup {
	return g.Handle(path, "GET", infos, handlers...)
}

// POST is a shortcut to register a new handler with the POST method.
func (g *RouterGroup) POST(path string, infos []OperationOption, handlers ...gin.HandlerFunc) *RouterGroup {
	return g.Handle(path, "POST", infos, handlers...)
}

// PUT is a shortcut to register a new handler with the PUT method.
func (g *RouterGroup) PUT(path string, infos []OperationOption, handlers ...gin.HandlerFunc) *RouterGroup {
	return g.Handle(path, "PUT", infos, handlers...)
}

// PATCH is a shortcut to register a new handler with the PATCH method.
func (g *RouterGroup) PATCH(path string, infos []OperationOption, handlers ...gin.HandlerFunc) *RouterGroup {
	return g.Handle(path, "PATCH", infos, handlers...)
}

// DELETE is a shortcut to register a new handler with the DELETE method.
func (g *RouterGroup) DELETE(path string, infos []OperationOption, handlers ...gin.HandlerFunc) *RouterGroup {
	return g.Handle(path, "DELETE", infos, handlers...)
}

// OPTIONS is a shortcut to register a new handler with the OPTIONS method.
func (g *RouterGroup) OPTIONS(path string, infos []OperationOption, handlers ...gin.HandlerFunc) *RouterGroup {
	return g.Handle(path, "OPTIONS", infos, handlers...)
}

// HEAD is a shortcut to register a new handler with the HEAD method.
func (g *RouterGroup) HEAD(path string, infos []OperationOption, handlers ...gin.HandlerFunc) *RouterGroup {
	return g.Handle(path, "HEAD", infos, handlers...)
}

// TRACE is a shortcut to register a new handler with the TRACE method.
func (g *RouterGroup) TRACE(path string, infos []OperationOption, handlers ...gin.HandlerFunc) *RouterGroup {
	return g.Handle(path, "TRACE", infos, handlers...)
}

// Handle registers a new request handler that is wrapped
// with Tonic and documented in the OpenAPI specification.
func (g *RouterGroup) Handle(path, method string, infos []OperationOption, handlers ...gin.HandlerFunc) *RouterGroup {
	oi := &openapi.OperationInfo{}
	for _, info := range infos {
		info(oi)
	}
	type wrap struct {
		h gin.HandlerFunc
		r *tonic.Route
	}
	var wrapped []wrap

	// Find the handlers wrapped with Tonic.
	for _, h := range handlers {
		r, err := tonic.GetRouteByHandler(h)
		if err == nil {
			wrapped = append(wrapped, wrap{h: h, r: r})
		}
	}
	// Check that no more that one tonic-wrapped handler
	// is registered for this operation.
	if len(wrapped) > 1 {
		panic(fmt.Sprintf("multiple tonic-wrapped handler used for operation %s %s", method, path))
	}
	// If we have a tonic-wrapped handler, generate the
	// specification of this operation.
	if len(wrapped) == 1 {
		hfunc := wrapped[0].r

		// Set an operation ID if none is provided.
		if oi.ID == "" {
			oi.ID = hfunc.HandlerName()
		}
		oi.StatusCode = hfunc.GetDefaultStatusCode()

		// Set an input type if provided.
		it := hfunc.InputType()
		if oi.InputModel != nil {
			it = reflect.TypeOf(oi.InputModel)
		}

		// Consolidate path for OpenAPI spec.
		operationPath := joinPaths(g.group.BasePath(), path)

		// Add operation to the OpenAPI spec.
		operation, err := g.gen.AddOperation(operationPath, method, g.Name, it, hfunc.OutputType(), oi)
		if err != nil {
			panic(fmt.Sprintf(
				"error while generating OpenAPI spec on operation %s %s: %s",
				method, path, err,
			))
		}
		// If an operation was generated for the handler,
		// wrap the Tonic-wrapped handled with a closure
		// to inject it into the Gin context.
		if operation != nil {
			for i, h := range handlers {
				if funcEqual(h, wrapped[0].h) {
					orig := h // copy the original func
					handlers[i] = func(c *gin.Context) {
						c.Set(ctxOpenAPIOperation, operation)
						orig(c)
					}
				}
			}
		}
	}
	// Register the handlers with Gin underlying group.
	g.group.Handle(method, path, handlers...)

	return g
}

// OpenAPI returns a Gin HandlerFunc that serves
// the marshalled OpenAPI specification of the API.
func (f *Fizz) OpenAPI(info *openapi.Info, ct string) gin.HandlerFunc {
	f.gen.SetInfo(info)

	ct = strings.ToLower(ct)
	if ct == "" {
		ct = "json"
	}
	switch ct {
	case "json":
		return func(c *gin.Context) {
			c.JSON(200, f.gen.API())
		}
	case "yaml":
		return func(c *gin.Context) {
			c.YAML(200, f.gen.API())
		}
	}
	panic("invalid content type, use JSON or YAML")
}

// OperationOption represents an option-pattern function
// used to add informations to an operation.
type OperationOption func(*openapi.OperationInfo)

// StatusDescription sets the default status description of the operation.
func StatusDescription(desc string) func(*openapi.OperationInfo) {
	return func(o *openapi.OperationInfo) {
		o.StatusDescription = desc
	}
}

// Summaryf adds a summary to an operation according
// to a format specifier.
func Summaryf(format string, a ...interface{}) func(*openapi.OperationInfo) {
	return func(o *openapi.OperationInfo) {
		o.Summary = fmt.Sprintf(format, a...)
	}
}

// Descriptionf adds a description to an operation
// according to a format specifier.
func Descriptionf(format string, a ...interface{}) func(*openapi.OperationInfo) {
	return func(o *openapi.OperationInfo) {
		o.Description = fmt.Sprintf(format, a...)
	}
}

// ID overrides the operation ID.
func ID(id string) func(*openapi.OperationInfo) {
	return func(o *openapi.OperationInfo) {
		o.ID = id
	}
}



// Response adds an additional response to the operation.
func Response(statusCode, desc string, model interface{}, headers []*openapi.ResponseHeader, example interface{}) func(*openapi.OperationInfo) {
	return func(o *openapi.OperationInfo) {
		o.Responses = append(o.Responses, &openapi.OperationResponse{
			Code:        statusCode,
			Description: desc,
			Model:       model,
			Headers:     headers,
			Example:     example,
		})
	}
}

// ResponseWithExamples is a variant of Response that accept many examples.
func ResponseWithExamples(statusCode, desc string, model interface{}, headers []*openapi.ResponseHeader, examples map[string]interface{}) func(*openapi.OperationInfo) {
	return func(o *openapi.OperationInfo) {
		o.Responses = append(o.Responses, &openapi.OperationResponse{
			Code:        statusCode,
			Description: desc,
			Model:       model,
			Headers:     headers,
			Examples:    examples,
		})
	}
}

// Header adds a header to the operation.
func Header(name, desc string, model interface{}) func(*openapi.OperationInfo) {
	return func(o *openapi.OperationInfo) {
		o.Headers = append(o.Headers, &openapi.ResponseHeader{
			Name:        name,
			Description: desc,
			Model:       model,
		})
	}
}

// InputModel overrides the binding model of the operation.
func InputModel(model interface{}) func(*openapi.OperationInfo) {
	return func(o *openapi.OperationInfo) {
		o.InputModel = model
	}
}

// XCodeSample adds a code sample to the operation.
func XCodeSample(cs *openapi.XCodeSample) func(*openapi.OperationInfo) {
	return func(o *openapi.OperationInfo) {
		o.XCodeSamples = append(o.XCodeSamples, cs)
	}
}

// OperationFromContext returns the OpenAPI operation from
// the givent Gin context or an error if none is found.
func OperationFromContext(c *gin.Context) (*openapi.Operation, error) {
	if v, ok := c.Get(ctxOpenAPIOperation); ok {
		if op, ok := v.(*openapi.Operation); ok {
			return op, nil
		}
		return nil, errors.New("invalid type: not an operation")
	}
	return nil, errors.New("operation not found")
}

func joinPaths(abs, rel string) string {
	if rel == "" {
		return abs
	}
	final := path.Join(abs, rel)
	as := lastChar(rel) == '/' && lastChar(final) != '/'
	if as {
		return final + "/"
	}
	return final
}

func lastChar(str string) uint8 {
	if str == "" {
		panic("empty string")
	}
	return str[len(str)-1]
}

func funcEqual(f1, f2 interface{}) bool {
	v1 := reflect.ValueOf(f1)
	v2 := reflect.ValueOf(f2)

	if v1.Kind() == reflect.Func && v2.Kind() == reflect.Func { // prevent panic on call to Pointer()
		return runtime.FuncForPC(v1.Pointer()).Entry() == runtime.FuncForPC(v2.Pointer()).Entry()
	}
	return false
}
