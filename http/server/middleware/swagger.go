package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpswagger "github.com/swaggo/http-swagger"
)

// Info is used to load environment variable with the value we want to modify at
// run time.
type Info struct {
	// Host need to be ingress value without schema. Eg.: example.com:8080.
	Host string `env:"SWAGGER_HOST,required"`
}

// Handler refToSwaggerHost point to SwaggerInfo.Host generated value and return
// an http.Handler to be mounted.
//
// First generate swagger files:
//
//	make swagger
//
// Mount desired route invoking this method with SwaggerInfo from generated files:
//
//	r.Mount("/swagger", dep.Components.Swagger.Handler(&docs.SwaggerInfo.Host))
func (s Info) Handler(refToSwaggerHost *string) http.Handler {
	r := chi.NewRouter()

	*refToSwaggerHost = s.Host
	r.Get("/*", httpswagger.Handler(
		httpswagger.URL("docs/doc.json"),
	))

	return r
}
