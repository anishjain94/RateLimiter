package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"ratelimit/infra/middleware"
	"ratelimit/util"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := mux.NewRouter().StrictSlash(true)

	middleware.NewRateLimiter(util.CLEANUP_EXPIRY)

	addTestData()

	initializeMiddleware(router)

	InitializeRoutes(router)

	httpServer := &http.Server{
		Handler: router,
		Addr:    os.Getenv("DOMAIN") + ":" + os.Getenv("HOST"),
	}

	err = httpServer.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}

}

func InitializeRoutes(router *mux.Router) {
	router.HandleFunc("/health", handleHealth).Methods(http.MethodGet)
	router.HandleFunc("/info", handleInfo).Methods(http.MethodGet)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Health: up and running")
}

func handleInfo(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Info: up and running")
}

func initializeMiddleware(router *mux.Router) {

	router.Use(middleware.RateLimiterMiddleware)
}

func addTestData() {
	ip := "127.0.0.1"
	middleware.AddConfig("/health"+"_"+ip, 20, 1*time.Second)
	middleware.AddConfig("/info"+"_"+ip, 10, 1*time.Second)
}
