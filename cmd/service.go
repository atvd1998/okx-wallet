package cmd

import (
	"context"
	"net/http"
	"okx-wallet/config"
	"okx-wallet/internal/controller"
	"okx-wallet/internal/repository"
	"okx-wallet/internal/service"
	"okx-wallet/pkg/logger"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

// ServiceCommand returns the cobra command for running the service
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "API Command of service",
	Long:  "API Command of service",
	Run: func(_ *cobra.Command, _ []string) {
		invoke(start).Run()
	},
}

func registerRoute(e *echo.Echo, ctrl *controller.Controller) {
	g := e.Group("/api")
	v1 := g.Group("/v1")

	okxG := v1.Group("/okx")
	okxG.GET("/test", func(c echo.Context) error {
		ctrl.GetConnection()
		return c.JSON(http.StatusOK, "test")
	})
}

func start(
	lc fx.Lifecycle,
	conf *config.Config,
	ctrl *controller.Controller,
) {
	logger := logger.MustNamed("service")
	e := echo.New()
	registerRoute(e, ctrl)
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.Infow("Starting HTTP server", "addr", conf.App.HTTPAddr)
				if err := e.Start(conf.App.HTTPAddr); err != nil && err != http.ErrServerClosed {
					logger.Fatalw("Error starting HTTP server", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down HTTP server")
			return e.Shutdown(ctx)
		},
	})

}

func invoke(invokers ...any) *fx.App {
	conf := config.MustLoad()
	app := fx.New(
		fx.Provide(
			controller.NewController,
			fx.Annotate(service.NewOKXService, fx.As(new(controller.OKXService))),
			fx.Annotate(repository.NewOKXRepository, fx.As(new(service.OKXRepository))),
		),
		fx.Supply(conf),
		fx.Invoke(invokers...),
	)

	return app
}
