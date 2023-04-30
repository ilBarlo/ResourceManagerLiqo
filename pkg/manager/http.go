package manager

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SetupRouterAndServeHTTP setups and starts the http server. If everything goes well, it never returns.
func SetupRouterAndServeHTTP(ctx context.Context, cl client.Client) error {
	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		klog.Fatal("environment variable named ADDRS not set")
	}

	router := routes(ctx, cl)
	server := &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: 5 * time.Second,
		Handler:           cors.Default().Handler(router),
	}
	return server.ListenAndServe()
}

func routes(ctx context.Context, cl client.Client) http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	// Routes to serve
	router.HandleFunc("/api/foreign_clusters", foreignClusters(ctx, cl))
	router.HandleFunc("/api/available_resources", availableResources(ctx, cl))

	return router
}

func foreignClusters(ctx context.Context, cl client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clusters, err := getForeignClusters(ctx, cl)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(ErrorResponse{
				Message: "An error occurred during retrieving foreign clusters and metrics",
				Status:  http.StatusInternalServerError,
			})
			if err != nil {
				klog.Errorf("error encoding error response: %s", err)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(*clusters)
		if err != nil {
			klog.Errorf("error encoding foreign clusters response: %s", err)
		}
	}
}

func availableResources(ctx context.Context, cl client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clustersMap, err := aggregateAvailableResources(ctx, cl)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(ErrorResponse{
				Message: "An error occurred during retrieving maps",
				Status:  http.StatusInternalServerError,
			})
			if err != nil {
				klog.Errorf("error encoding error response: %s", err)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(*clustersMap)
		if err != nil {
			klog.Errorf("error encoding foreign clusters response: %s", err)
		}
	}
}
