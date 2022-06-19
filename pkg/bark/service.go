package bark

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Service is the base type for Bark services.
type Service struct {
	Logger interface {
		Printf(format string, v ...interface{})
	}
	Name string
}

// Logf provides a standard logging mechanism for services.
func (service *Service) Logf(r *http.Request, format string, v ...interface{}) {
	if service.Logger != nil {
		if service.Name == "" {
			service.Name = os.Args[0]
		}

		requestID := middleware.GetReqID(r.Context())
		if requestID == "" {
			requestID = "???"
		}

		var paramFields string
		params := chi.RouteContext(r.Context()).URLParams
		for i := 0; i < len(params.Keys); i++ {
			paramFields += fmt.Sprintf(" %s=%s", params.Keys[i], params.Values[i])
		}

		defaultFields := fmt.Sprintf("serviceName=%s requestID=%s method=%s endpoint=%s%s",
			service.Name, requestID, r.Method, r.URL.EscapedPath(), paramFields)

		if format == "" {
			service.Logger.Printf(defaultFields)
		} else {
			fullFormat := defaultFields + " " + format
			service.Logger.Printf(fullFormat, v...)
		}
	}
}
