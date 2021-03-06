package handlers

import (
	"database/sql"
	"github.com/vardius/go-api-boilerplate/internal/errors"
	"github.com/vardius/go-api-boilerplate/internal/http/response"
	"net/http"

	grpc_utils "github.com/vardius/go-api-boilerplate/internal/grpc"
	"google.golang.org/grpc"
)

// BuildLivenessHandler provides liveness handler
func BuildLivenessHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}

	return http.HandlerFunc(fn)
}

// BuildReadinessHandler provides readiness handler
func BuildReadinessHandler(db *sql.DB, connMap map[string]*grpc.ClientConn) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if err := db.PingContext(r.Context()); err != nil {
			response.RespondJSONError(r.Context(), w, errors.Wrap(err, errors.INTERNAL, "Database is not responding"))
			return
		}

		for name, conn := range connMap {
			if !grpc_utils.IsConnectionServing(name, conn) {
				response.RespondJSONError(r.Context(), w, errors.Newf(errors.INTERNAL, "gRPC connection %s is not serving", name))
				return
			}
		}

		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}

	return http.HandlerFunc(fn)
}
