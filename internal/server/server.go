// server.go
package server

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/gorilla/mux"
	"github.com/sarabrajsingh/restful-openapi/config"
	"github.com/sarabrajsingh/restful-openapi/internal/global_errors"
	"github.com/sarabrajsingh/restful-openapi/internal/handlers"
	"github.com/sarabrajsingh/restful-openapi/internal/logging"
	"github.com/sarabrajsingh/restful-openapi/internal/utils"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

type serverImpl struct {
	config        *config.Config
	logger        logging.Logger
	errorStore    global_errors.ErrorStore
	bodyReader    func(io.Reader) ([]byte, error)
	specification *openapi3.T
}

func NewServer(config *config.Config, logger logging.Logger, errorStore global_errors.ErrorStore, bodyReader func(io.Reader) ([]byte, error)) Server {
	return &serverImpl{
		config:     config,
		logger:     logger,
		errorStore: errorStore,
		bodyReader: bodyReader,
	}
}

func (s *serverImpl) GetSpecification() *openapi3.T {
	return s.specification
}

func (s *serverImpl) NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	// Load and validate OpenAPI spec
	spec, err := openapi3.NewLoader().LoadFromFile(s.config.OpenAPI3YamlFileLocation)
	if err != nil {
		s.logger.Fatalf("Failed to load OpenAPI spec: %v", err)
	}

	s.specification = spec

	err = spec.Validate(openapi3.NewLoader().Context)
	if err != nil {
		s.logger.Fatalf("Failed to validate OpenAPI spec: %v", err)
	}

	// Create OpenAPI router
	oapiRouter, err := gorillamux.NewRouter(spec)
	if err != nil {
		s.logger.Fatalf("Failed to create OpenAPI router: %v", err)
	}

	s.logger.Printf("Validating Contract")

	// Define routes with /api/v1/ prefix
	routes := []Route{
		{
			Name:        "Index",
			Method:      "GET",
			Pattern:     "/",
			HandlerFunc: Index,
		},
		{
			Name:        "ErrorsDelete",
			Method:      strings.ToUpper("DELETE"),
			Pattern:     "/errors",
			HandlerFunc: handlers.DeleteErrors(s.logger, s.errorStore.DeleteErrors),
		},
		{
			Name:        "ErrorsGet",
			Method:      strings.ToUpper("GET"),
			Pattern:     "/errors",
			HandlerFunc: handlers.GetErrors(s.logger, s.errorStore.GetErrors),
		},
		{
			Name:        "TempPost",
			Method:      strings.ToUpper("POST"),
			Pattern:     "/temp",
			HandlerFunc: handlers.TempPost(s.logger, s.errorStore.AddError, utils.DefaultBodyReader),
		},
	}

	// Register routes with middleware
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		// logging middleware for handlers
		handler = LoggerMiddleware(s.logger, handler, route.Name)
		// openapi3 validaton middleware for each handler request
		handler = OpenAPIMiddleware(oapiRouter, handler)

		router.
			Methods(route.Method).
			Path("/api/v1" + route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	s.logger.Printf("Successfully validated the contract")

	// Serve Swagger UI
	fs := http.FileServer(http.Dir(s.config.SwaggerUIFolder))
	router.PathPrefix("/").Handler(http.StripPrefix("/", fs))

	return router
}

// LoggerMiddleware logs the HTTP request details
func LoggerMiddleware(logger logging.Logger, inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			r.RemoteAddr,
		)
		inner.ServeHTTP(w, r)
	})
}

func OpenAPIMiddleware(router routers.Router, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route, pathParams, err := router.FindRoute(r)
		if err != nil {
			fmt.Printf("Request: %+v\n", r)
			response := fmt.Sprintf("OpenAPI Middleware: Error finding route: %v\n", err)
			utils.WriteErrorResponse(w, response, http.StatusBadRequest)
			return
		}

		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    r,
			PathParams: pathParams,
			Route:      route,
		}

		if err := openapi3filter.ValidateRequest(r.Context(), requestValidationInput); err != nil {
			response := fmt.Sprintf("OpenAPI Middleware: Request validation failed: %v\n", err)
			utils.WriteErrorResponse(w, response, http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Home/Landing page for the API here
func Index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/index.html", http.StatusFound)
}
