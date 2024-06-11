package manager

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"semay.com/configs"
	"semay.com/models/controlers"

	"github.com/spf13/cobra"
)

var (
	runCmd = &cobra.Command{
		Use:   "dev",
		Short: "Run Development server ",
		Long:  `Run Gofr development server`,
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

func run() {

	app := echo.New()
	configs.NewEnvFile("./configs")

	//  prometheus metrics middleware
	app.Use(echoprometheus.NewMiddleware("echo_blue"))

	// Rate Limiting to throttle overload
	app.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(1000)))

	// Recover incase of panic attacks
	app.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
		LogLevel:  log.ERROR,
	}))

	setupRoutes(app)
	// starting on provided port
	go func(app *echo.Echo) {
		//  Http serving port
		HTTP_PORT := configs.AppConfig.Get("HTTP_PORT")
		app.Logger.Fatal(app.Start("0.0.0.0:" + HTTP_PORT))
		// log.Fatal(app.ListenTLS(":" + port_1, "server.pem", "server-key.pem"))
	}(app)

	c := make(chan os.Signal, 1)   // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt) // When an interrupt or termination signal is sent, notify the channel

	<-c // This blocks the main thread until an interrupt is received
	fmt.Println("Gracefully shutting down...")

	fmt.Println("Running cleanup tasks...")
	// Your cleanup tasks go here
	fmt.Println("Blue was successful shutdown.")

}

func init() {
	goBlueCmd.AddCommand(runCmd)

}

func setupRoutes(app *echo.Echo) {
	gapp := app.Group("/admin")
	gapp.GET("/roles", controlers.GetRoles).Name = "get_all_roles"
	gapp.GET("/role/:role_id", controlers.GetRoleByID).Name = "get_one_roles"
	gapp.POST("/role", controlers.PostRole).Name = "post_role"
	gapp.PATCH("/role/:role_id", controlers.PatchRole).Name = "patch_role"
	gapp.DELETE("/role/:role_id", controlers.PostRole).Name = "post_role"

}
