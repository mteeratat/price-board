package main

import (
	"bytes"
	"context"
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v5"

	"github.com/mteeratat/price-board/app/internal/env"
	"github.com/mteeratat/price-board/app/internal/store"
)

//go:embed index.html
var files embed.FS

// Parsed once at startup: a broken template crashes at boot, not on first request.
var page = template.Must(template.ParseFS(files, "index.html"))

func main() {
	// No default: a missing DB secret must fail loudly.
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	port := env.Get("PORT", "8080")

	// Cancelled on Ctrl+C or SIGTERM (what Kubernetes sends before killing a pod).
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	st, err := store.New(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer st.Close()

	e := echo.New()

	// Liveness: process is up. Never touches the DB, so a DB outage doesn't restart pods.
	e.GET("/healthz", func(c *echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	// Readiness: DB reachable. 503 takes the pod out of traffic without restarting it.
	e.GET("/readyz", func(c *echo.Context) error {
		if err := st.Ping(c.Request().Context()); err != nil {
			return c.String(http.StatusServiceUnavailable, "db unreachable")
		}
		return c.String(http.StatusOK, "ok")
	})

	e.GET("/prices", func(c *echo.Context) error {
		prices, err := st.LatestPrices(c.Request().Context())
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, "internal error")
		}
		return c.JSON(http.StatusOK, prices)
	})

	e.GET("/", func(c *echo.Context) error {
		prices, err := st.LatestPrices(c.Request().Context())
		if err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, "internal error")
		}
		// Render to a buffer first, so a template error never sends half a page.
		var buf bytes.Buffer
		if err := page.Execute(&buf, prices); err != nil {
			log.Print(err)
			return c.String(http.StatusInternalServerError, "internal error")
		}
		return c.HTMLBlob(http.StatusOK, buf.Bytes())
	})

	sc := echo.StartConfig{Address: ":" + port}
	if err := sc.Start(ctx, e); err != nil {
		log.Fatal(err)
	}
}
