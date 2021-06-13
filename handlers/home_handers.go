package handlers

import (
	// ThirdParty Packages
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func HomeHandler(db *gorm.DB) (fn AppHandler) {
	fn = func(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
		vars := mux.Vars(req)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Category: %v\n", vars["category"])

		return nil
	}

	return
}
