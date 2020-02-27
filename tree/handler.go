package tree

import (
	"fmt"
	"net/http"
)

// NewHandler defines dependencies for our http handler and then returns the handler
func NewHandler(urlQueryKey string, view ViewFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var tree string

		input, inputIsPresent := req.URL.Query()[urlQueryKey]
		if inputIsPresent && len(input) > 0 && input[0] != "" {
			tree = input[0]
		} else {
			w.WriteHeader(http.StatusNotFound)
		}

		w.Header().Set("Content-Type", "text/html")

		// sanitization is performed by html/template automatically
		if err := view(w, tree); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error while writing the tree: %s\n", err)
		}
	})
}
