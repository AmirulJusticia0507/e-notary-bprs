// Package handler adalah entry point serverless untuk deploy backend
// di Vercel (project e-notary-bprs-api dengan Root Directory = backend).
// Vercel memetakan file ini menjadi function yang menangani semua route,
// lalu meneruskannya ke Gin router yang sama dengan server lokal.
//
// Env yang wajib diisi di dashboard Vercel (Settings → Environment Variables):
//   DATABASE_URL  (postgres hosted, mis. Neon/Supabase — localhost tidak bisa dijangkau Vercel)
//   JWT_SECRET, JWT_EXPIRY,
//   CBS_*, BPN_*, PEGADAIAN_*, EMETERAI_*, ESIGN_* (opsional, kosongkan bila belum ada)
package handler

import (
	"log"
	"net/http"
	"sync"

	"github.com/e-notary-bprs/backend/internal/config"
	"github.com/e-notary-bprs/backend/internal/delivery"
	"github.com/e-notary-bprs/backend/pkg/db"
	"github.com/gin-gonic/gin"
)

var (
	router   *gin.Engine
	initOnce sync.Once
	initErr  error
)

func initRouter() {
	initOnce.Do(func() {
		gin.SetMode(gin.ReleaseMode)
		cfg, err := config.Load()
		if err != nil {
			initErr = err
			return
		}
		database, err := db.Open(cfg.Database)
		if err != nil {
			initErr = err
			return
		}
		router = delivery.NewRouter(cfg, database)
	})
}

// Handler adalah entry point Vercel serverless function.
func Handler(w http.ResponseWriter, r *http.Request) {
	initRouter()
	if initErr != nil {
		log.Printf("init error: %v", initErr)
		http.Error(w, `{"status":"error","error":"service unavailable"}`, http.StatusServiceUnavailable)
		return
	}
	router.ServeHTTP(w, r)
}
