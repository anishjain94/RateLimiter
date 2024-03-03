package main

import (
	"encoding/json"
	"net/http"
	"ratelimit/infra/middleware"
	"ratelimit/util"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter().StrictSlash(true)

	middleware.NewRateLimiter(util.CLEANUP_EXPIRY)
	middleware.AddConfig("v1/health_127.0.0.1", 30, 1*time.Second)
	initializeMiddleware(router)

	InitializeRoutes(router)

	httpServer := &http.Server{
		Handler: router,
		Addr:    "0.0.0.0:8000",
	}

	err := httpServer.ListenAndServe()

	if err != nil {
		panic(err.Error())
	}

}

func InitializeRoutes(router *mux.Router) {
	router.HandleFunc("/health", handleHealth).Methods(http.MethodGet)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("up and running")
}

func initializeMiddleware(router *mux.Router) {

	router.Use(middleware.RateLimiterMiddleware)
	// router.Use(middleware.GetThrottlingMiddleWare(1*time.Second, 10)) //NOTE: Alternate Approach
}
