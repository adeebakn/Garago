package apis

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"garago/internal/config"
	"garago/internal/database"
	"garago/internal/database/postgres"
	mailservice "garago/internal/platform/mailer"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const shutdownTimeout = 10 * time.Second

// Run initializes dependencies, serves requests, and drains in-flight work on shutdown.
func Run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	queries, pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer pool.Close()
	databasePool = pool

	if err := initializeAuth(ctx, cfg); err != nil {
		return fmt.Errorf("initialize authentication: %w", err)
	}
	defer redisClient.Close()
	passwordMailer, err = mailservice.New(cfg)
	if err != nil {
		return fmt.Errorf("initialize email: %w", err)
	}
	passwordResetFrontendURL = cfg.FrontendURL

	server := &http.Server{
		Addr:              cfg.ServerAddr,
		Handler:           newRouter(queries, cfg.FrontendURL),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() { serverErrors <- server.ListenAndServe() }()

	select {
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("serve HTTP: %w", err)
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown HTTP server: %w", err)
		}
		return nil
	}
}

func newRouter(queries *postgres.Queries, frontendURL string) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{frontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	api := router.Group("/api")
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Garago API is running"})
	})
	api.GET("/organization", getOrganizationSettings)
	api.POST("/auth/signup", func(c *gin.Context) { signup(c, queries) })
	api.POST("/auth/login", func(c *gin.Context) { login(c, queries) })
	api.POST("/auth/forgot-password", func(c *gin.Context) { forgotPassword(c, queries) })
	api.POST("/auth/reset-password", func(c *gin.Context) { resetPassword(c, queries) })

	protected := api.Group("")
	protected.Use(authenticationRequired())
	protected.GET("/auth/me", func(c *gin.Context) { currentUser(c, queries) })
	protected.POST("/auth/logout", logout)
	protected.PUT("/organization", updateOrganizationSettings)

	protected.POST("/customers", func(c *gin.Context) { createCustomer(c, queries) })
	protected.GET("/customers", func(c *gin.Context) { getCustomers(c, queries) })
	protected.GET("/customers/:id", func(c *gin.Context) { getCustomerByID(c, queries) })
	protected.PUT("/customers/:id", func(c *gin.Context) { updateCustomerByID(c, queries) })
	protected.DELETE("/customers/:id", func(c *gin.Context) { deleteCustomerByID(c, queries) })

	protected.POST("/vehicles", func(c *gin.Context) { createVehicle(c, queries) })
	protected.GET("/vehicles", func(c *gin.Context) { getVehicles(c, queries) })
	protected.GET("/vehicles/:id", func(c *gin.Context) { getVehicleByID(c, queries) })
	protected.GET("/customers/:id/vehicles", func(c *gin.Context) { getVehiclesByCustomerID(c, queries) })
	protected.PUT("/vehicles/:id", func(c *gin.Context) { updateVehicleByID(c, queries) })
	protected.DELETE("/vehicles/:id", func(c *gin.Context) { deleteVehicleByID(c, queries) })

	protected.POST("/services", func(c *gin.Context) { createService(c, queries) })
	protected.GET("/services", func(c *gin.Context) { getServices(c, queries) })
	protected.GET("/services/:id", func(c *gin.Context) { getServiceByID(c, queries) })
	protected.PUT("/services/:id", func(c *gin.Context) { updateServiceByID(c, queries) })
	protected.DELETE("/services/:id", func(c *gin.Context) { deleteServiceByID(c, queries) })

	protected.POST("/quotations", func(c *gin.Context) { createQuotation(c, queries) })
	protected.GET("/quotations", func(c *gin.Context) { getQuotations(c, queries) })
	protected.GET("/quotations/:id", func(c *gin.Context) { getQuotationByID(c, queries) })
	protected.PUT("/quotations/:id", func(c *gin.Context) { updateQuotationByID(c, queries) })
	protected.PUT("/quotations/:id/status", func(c *gin.Context) { updateQuotationStatus(c, queries) })
	protected.DELETE("/quotations/:id", func(c *gin.Context) { deleteQuotationByID(c, queries) })

	protected.POST("/work-orders", func(c *gin.Context) { createWorkOrder(c, queries) })
	protected.GET("/work-orders", func(c *gin.Context) { getWorkOrders(c, queries) })
	protected.GET("/work-orders/:id", func(c *gin.Context) { getWorkOrderByID(c, queries) })
	protected.PUT("/work-orders/:id/status", func(c *gin.Context) { updateWorkOrderStatus(c, queries) })
	protected.PUT("/work-orders/:id", func(c *gin.Context) { updatePendingWorkOrder(c, queries) })
	protected.DELETE("/work-orders/:id", func(c *gin.Context) { deleteWorkOrderByID(c, queries) })

	protected.POST("/invoices", func(c *gin.Context) { createInvoice(c, queries) })
	protected.GET("/invoices", func(c *gin.Context) { getInvoices(c, queries) })
	protected.GET("/invoices/:id", func(c *gin.Context) { getInvoiceByID(c, queries) })
	protected.DELETE("/invoices/:id", func(c *gin.Context) { deleteInvoiceByID(c, queries) })
	protected.POST("/payments", func(c *gin.Context) { createPayment(c, queries) })
	protected.GET("/invoices/:id/payments", func(c *gin.Context) { getPaymentsByInvoice(c, queries) })
	protected.GET("/payments", func(c *gin.Context) { getPayments(c, queries) })
	protected.GET("/payments/:id", func(c *gin.Context) { getPaymentByID(c, queries) })
	protected.DELETE("/payments/:id", func(c *gin.Context) { deletePaymentByID(c, queries) })

	protected.GET("/work-orders/status/:status", func(c *gin.Context) { getWorkOrdersByStatus(c, queries) })
	protected.GET("/quotations/status/:status", func(c *gin.Context) { getQuotationsByStatus(c, queries) })
	protected.GET("/invoices/status/:status", func(c *gin.Context) { getInvoicesByStatus(c, queries) })
	protected.GET("/dashboard", func(c *gin.Context) { getDashboard(c, queries) })

	return router
}
