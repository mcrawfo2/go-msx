package restops

import "net/http"

type Middleware func(next http.Handler) http.Handler

type Middlewares []Middleware

func (m Middlewares) Compose(final http.Handler) http.Handler {
	here := final
	for _, mw := range m {
		here = mw(here)
	}
	return here
}
