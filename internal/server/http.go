package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	fiberSwagger "github.com/swaggo/fiber-swagger"

	_ "app/internal/docs"
)

func StartHttp(opts StartHttpOpts) error {
	if err := opts.Init(); err != nil {
		return fmt.Errorf("failed to initialise http server: %s", err)
	}

	serverHeader := opts.ServerName
	if opts.Version != "" {
		serverHeader += "-" + opts.Version
	}

	app := fiber.New(fiber.Config{
		AppName:               opts.ServerName,
		DisableStartupMessage: true,
		ReadTimeout:           10 * time.Second,
		WriteTimeout:          10 * time.Second,
		IdleTimeout:           15 * time.Second,
		EnablePrintRoutes:     true,
		ServerHeader:          opts.ServerName,
	})
	app.Use(requestid.New(requestid.Config{
		Generator: func() string {
			return uuid.New().String()
		},
	}))
	app.Use(logger.New(logger.Config{
		Format:     "${time} ${locals:requestid} ${status} ${method} ${url} ${ip}:${port} ${ua} ${latency} ${bytesReceived} ${bytesSent}\n${resBody}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "UTC",
	}))

	if !opts.DisableLivenessProbe {
		if opts.LivenessProbes == nil {
			opts.LivenessProbes = []Healthcheck{}
		}
		app.Get("/healthz", func(ctx *fiber.Ctx) error {
			errors := []string{}
			for _, probe := range opts.LivenessProbes {
				if err := probe(); err != nil {
					errors = append(errors, err.Error())
				}
			}
			if len(errors) > 0 {
				return ctx.Status(http.StatusBadRequest).JSON("not ok, see logs for more")
			}
			return ctx.JSON("ok")
		})
	}

	if !opts.DisableReadinessProbe {
		if opts.ReadinessProbes == nil {
			opts.ReadinessProbes = []Healthcheck{}
		}
		app.Get("/readyz", func(ctx *fiber.Ctx) error {
			return ctx.JSON("ok")
		})
	}

	if !opts.DisableMetrics {
		app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
	}

	if !opts.DisableSwagger {
		app.Get("/swagger/*", fiberSwagger.WrapHandler)
	}

	if opts.Router != nil {
		opts.Router(app)
	}

	bindAddr := fmt.Sprintf("%s:%v", opts.BindInterface, opts.BindPort)
	if err := app.Listen(bindAddr); err != nil {
		return fmt.Errorf("failed to start server: %s", err)
	}
	return nil
}
