package middleware

import (
	"context"
	"github.com/go-chi/chi"
	"net/http"
	"strings"
)

// Contexts extract all of the params related to their route
type contextMapKey string

func (c contextMapKey) String() string {
	return "pkg.server.context" + string(c)
}

var (
	ctxMapKey = contextMapKey("ctxMap")
)

// ContextMap extract from request and type asserts it (helper function.)
func ContextMap(r *http.Request) map[string]string {
	return (r.Context().Value(ctxMapKey)).(map[string]string)
}

// GopherID extracts from request context
func GopherID(r *http.Request) string {
	return ContextMap(r)["GopherID"]
}

// ThingID extracts from request context
func ThingID(r *http.Request) string {
	return ContextMap(r)["ThingID"]
}

// GopherName extracts from request context
func GopherName(r *http.Request) string {
	return ContextMap(r)["GopherName"]
}

// GopherDescription extracts from request context
func GopherDescription(r *http.Request) string {
	return ContextMap(r)["GopherDescription"]
}

// IsAuthenticated looks at request context AccessKey and SecretKey, if these values are present then the user is Authenticated.
func IsAuthenticated(r *http.Request) (authed bool) {
	ctxMap := ContextMap(r)
	if ctxMap["AccessKey"] != "" && ctxMap["SecretKey"] != "" {
		authed = true
	}
	return authed
}

// InitialCtx runs for every route, sets the response to JSON for all responses and unpacks AccessKey&SecretKey
func InitialCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		ctxMap := make(map[string]string)

		xKeys := strings.Split(r.Header.Get("X-ApiKeys"), ";")
		for x := range xKeys {
			keys := strings.Split(xKeys[x], "=")
			switch {
			case strings.ToLower(keys[0]) == "accesskey":
				ctxMap["AccessKey"] = keys[1]

			case strings.ToLower(keys[0]) == "secretkey":
				ctxMap["SecretKey"] = keys[1]
			}
		}
		ctx := context.WithValue(r.Context(), ctxMapKey, ctxMap)
		w = NewPrettyPrint(w)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GopherCtx enforces non-authenicated read-only (GET/HEAD) requests and sets:
func GopherCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(r) {
			switch strings.ToUpper(r.Method) {
			case "GET", "HEAD":
				// ALLOW!
			default:
				// DENY ALL OTHERS!
				http.Error(w, http.StatusText(403), 403)
				return
			}
		}

		ctxMap := r.Context().Value(ctxMapKey).(map[string]string)
		ctxMap["GopherID"] = chi.URLParam(r, "GopherID")
		ctxMap["GopherName"] = r.FormValue("GopherName")
		ctxMap["GopherDescription"] = r.FormValue("GopherDescription")
		ctx := context.WithValue(r.Context(), ctxMapKey, ctxMap)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ThingCtx requires IsAuthenticated() for ALL HTTP methods
func ThingCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(r) {
			// PS: Never /actually/ do this.. create a proper SecurityCtx and evaluates the uri etc. :-)
			http.Error(w, http.StatusText(403), 403)
			return
		}

		ctxMap := r.Context().Value(ctxMapKey).(map[string]string)
		ctxMap["ThingID"] = chi.URLParam(r, "ThingID")
		ctxMap["ThingName"] = r.FormValue("ThingName")
		ctxMap["ThingDescription"] = r.FormValue("ThingDescription")
		ctx := context.WithValue(r.Context(), ctxMapKey, ctxMap)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
