package api

import (
	"net/http"
	"wasa-photo/service/api/auth"
	"wasa-photo/service/api/errors"
	"wasa-photo/service/api/structures"

	"github.com/julienschmidt/httprouter"
)

// httpRouterHandler is the signature for functions that accepts a reqcontext.RequestContext in addition to those
// required by the httprouter package.
type httpRouterHandler func(http.ResponseWriter, *http.Request, httprouter.Params, structures.User)

// wrap parses the request and adds a reqcontext.RequestContext instance related to the request.
func (rt *_router) wrap(fn httpRouterHandler) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		res, user := auth.CheckAuth(rt.db, r)
		if !res {
			errors.WriteResponse(rt.baseLogger, w, "Authentication failed", http.StatusUnauthorized, "Unauthorized access")
			return
		}

		fn(w, r, ps, user)
	}
}
