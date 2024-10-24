package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/eddie023/tenantx/pkg/store"
)

type Handler struct {
	db *store.DB
}

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		slog.Error("startup", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	ctx = store.SetTenantID(ctx, 0)

	dbConnectionURI := os.Getenv("DB_CONNECTION_URI")
	if dbConnectionURI == "" {
		return fmt.Errorf("DB_CONNECTION_URI environment variable must be set")
	}

	db, err := store.NewDB(ctx, dbConnectionURI)
	if err != nil {
		return fmt.Errorf("failed db connection: %w", err)
	}

	h := Handler{
		db: db,
	}

	http.HandleFunc("/product", h.getProducts)

	slog.Info("server listening on", "port", "8848")
	if err := http.ListenAndServe(":8848", nil); err != nil {
		return fmt.Errorf("server failed: %w", err)
	}

	return nil
}

func (h *Handler) getProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// NOTE: in Production you will get the tenantID ideally from authentication middleware
	// for the demonstration purpose, I am getting it from query param.
	tenantId := r.URL.Query().Get("tenantId")
	if tenantId == "" {
		ctx = store.SetTenantID(ctx, 0)
	} else {
		id, err := strconv.Atoi(tenantId)
		if err != nil {
			slog.Error("parsing tenantId", "invalid integer", tenantId, "err", err)
			http.Error(w, "invalid tenantId: tenantId must be a valid integer", http.StatusBadGateway)
			return
		}

		ctx = store.SetTenantID(ctx, id)
	}

	products, err := h.db.GetProducts(ctx)
	if err != nil {
		slog.Error("getting products failed", "err", err)
		http.Error(w, "db failed", http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(products)
	if err != nil {
		http.Error(w, "unable to marshal", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}
