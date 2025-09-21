package router

import (
	"net/http"
	"time"

	"github.com/cmdb-lite/backend/internal/auth"
	"github.com/cmdb-lite/backend/internal/config"
	"github.com/cmdb-lite/backend/internal/handlers"
	"github.com/cmdb-lite/backend/internal/logging"
	"github.com/cmdb-lite/backend/internal/middleware"
	"github.com/cmdb-lite/backend/internal/metrics"
	"github.com/cmdb-lite/backend/internal/repositories"
	"github.com/cmdb-lite/backend/internal/tracing"
	"github.com/gorilla/mux"
)

// SetupRouter sets up the HTTP router with all endpoints
func SetupRouter(
	db *repositories.DB,
	jwtManager *auth.JWTManager,
	passwordManager *auth.PasswordManager,
	cfg *config.Config,
) *mux.Router {
	// Create a new router
	r := mux.NewRouter()

	// Initialize logger
	logger := logging.NewLogger("cmdb-lite", cfg.LogLevel)

	// Initialize tracing if enabled
	if cfg.TracingEnabled {
		tracingConfig := tracing.Config{
			Enabled:      true,
			ServiceName:  cfg.TracingService,
			Environment:  cfg.TracingEnv,
			JaegerURL:    cfg.TracingJaegerURL,
			ZipkinURL:    cfg.TracingZipkinURL,
			SamplingRate: cfg.TracingSamplingRate,
		}
		if err := tracing.InitializeDefaultTracerProvider(tracingConfig); err != nil {
			logger.WithError(err).Error("Failed to initialize tracing provider")
		}
	}

	// Create repositories
	userRepo := repositories.NewUserPostgresRepository(db.DB)
	refreshTokenRepo := repositories.NewRefreshTokenPostgresRepository(db.DB)
	ciRepo := repositories.NewCIPostgresRepository(db.DB)
	relRepo := repositories.NewRelationshipPostgresRepository(db.DB)
	auditRepo := repositories.NewAuditLogPostgresRepository(db.DB)

	// Create handlers
	authHandler := handlers.NewAuthHandler(userRepo, refreshTokenRepo, jwtManager, passwordManager)
	ciHandler := handlers.NewCIHandler(ciRepo, relRepo, auditRepo)
	relHandler := handlers.NewRelationshipHandler(relRepo, auditRepo)
	auditLogHandler := handlers.NewAuditLogHandler(auditRepo)
	metricsHandler := handlers.NewMetricsHandler()

	// Apply common middleware
	r.Use(middleware.CORS(cfg))
	r.Use(middleware.ObservabilityMiddleware)
	r.Use(middleware.DatabaseQueryMiddleware)
	r.Use(middleware.Recovery)
	
	// Add tracing middleware if enabled
	if cfg.TracingEnabled {
		r.Use(tracing.HTTPTracingMiddleware)
	}

	// Update system metrics
	metrics.DefaultMetrics.UpdateSystemInfo("1.0.0", "dev", cfg.TracingEnv)

	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok", "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`))
	}).Methods("GET")

	// Metrics endpoint (if enabled)
	if cfg.MetricsEnabled {
		r.Handle(cfg.MetricsPath, metricsHandler).Methods("GET")
	}

	// API version 1
	// All API endpoints are versioned using path-based versioning (e.g., /api/v1/)
	// This allows for future API versions without breaking existing client integrations
	apiV1 := r.PathPrefix("/api/v1").Subrouter()

	// Authentication endpoints (no authentication required)
	authRouter := apiV1.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/login", authHandler.Login).Methods("POST")
	authRouter.HandleFunc("/refresh", authHandler.RefreshToken).Methods("POST")
	authRouter.HandleFunc("/validate", middleware.AuthMiddleware(jwtManager)(http.HandlerFunc(authHandler.ValidateToken))).Methods("GET")
	authRouter.HandleFunc("/logout", middleware.AuthMiddleware(jwtManager)(http.HandlerFunc(authHandler.Logout))).Methods("POST")

	// CI endpoints (authentication required)
	ciRouter := apiV1.PathPrefix("/cis").Subrouter()
	ciRouter.Use(middleware.AuthMiddleware(jwtManager))
	
	// CI endpoints that require admin or viewer role
	ciAdminViewerRouter := ciRouter.NewRoute().Subrouter()
	ciAdminViewerRouter.Use(middleware.RBACMiddleware("admin", "viewer"))
	
	ciAdminViewerRouter.HandleFunc("", ciHandler.GetAllCIs).Methods("GET")
	ciAdminViewerRouter.HandleFunc("/{id}", ciHandler.GetCI).Methods("GET")
	ciAdminViewerRouter.HandleFunc("/{id}/graph", ciHandler.GetCIGraph).Methods("GET")
	
	// CI endpoints that require admin role
	ciAdminRouter := ciRouter.NewRoute().Subrouter()
	ciAdminRouter.Use(middleware.RBACMiddleware("admin"))
	
	ciAdminRouter.HandleFunc("", ciHandler.CreateCI).Methods("POST")
	ciAdminRouter.HandleFunc("/{id}", ciHandler.UpdateCI).Methods("PUT")
	ciAdminRouter.HandleFunc("/{id}", ciHandler.DeleteCI).Methods("DELETE")

	// Relationship endpoints (authentication required)
	relRouter := apiV1.PathPrefix("/relationships").Subrouter()
	relRouter.Use(middleware.AuthMiddleware(jwtManager))
	
	// Relationship endpoints that require admin or viewer role
	relAdminViewerRouter := relRouter.NewRoute().Subrouter()
	relAdminViewerRouter.Use(middleware.RBACMiddleware("admin", "viewer"))
	
	relAdminViewerRouter.HandleFunc("", relHandler.GetAllRelationships).Methods("GET")
	relAdminViewerRouter.HandleFunc("/{id}", relHandler.GetRelationship).Methods("GET")
	
	// Relationship endpoints that require admin role
	relAdminRouter := relRouter.NewRoute().Subrouter()
	relAdminRouter.Use(middleware.RBACMiddleware("admin"))
	
	relAdminRouter.HandleFunc("", relHandler.CreateRelationship).Methods("POST")
	relAdminRouter.HandleFunc("/{id}", relHandler.UpdateRelationship).Methods("PUT")
	relAdminRouter.HandleFunc("/{id}", relHandler.DeleteRelationship).Methods("DELETE")

	// Audit log endpoints (authentication required)
	auditRouter := apiV1.PathPrefix("/audit-logs").Subrouter()
	auditRouter.Use(middleware.AuthMiddleware(jwtManager))
	
	// Audit log endpoints that require admin or viewer role
	auditAdminViewerRouter := auditRouter.NewRoute().Subrouter()
	auditAdminViewerRouter.Use(middleware.RBACMiddleware("admin", "viewer"))
	
	auditAdminViewerRouter.HandleFunc("", auditLogHandler.GetAllAuditLogs).Methods("GET")
	auditAdminViewerRouter.HandleFunc("/{id}", auditLogHandler.GetAuditLog).Methods("GET")
	auditAdminViewerRouter.HandleFunc("/entity-type/{entity_type}", auditLogHandler.GetAuditLogsByEntityType).Methods("GET")
	auditAdminViewerRouter.HandleFunc("/entity-id/{entity_id}", auditLogHandler.GetAuditLogsByEntityID).Methods("GET")
	auditAdminViewerRouter.HandleFunc("/changed-by/{changed_by}", auditLogHandler.GetAuditLogsByChangedBy).Methods("GET")
	
	// Audit log endpoints that require admin role
	auditAdminRouter := auditRouter.NewRoute().Subrouter()
	auditAdminRouter.Use(middleware.RBACMiddleware("admin"))
	
	auditAdminRouter.HandleFunc("/{id}", auditLogHandler.DeleteAuditLog).Methods("DELETE")

	// User endpoints (authentication required)
	userRouter := apiV1.PathPrefix("/users").Subrouter()
	userRouter.Use(middleware.AuthMiddleware(jwtManager))
	
	// User endpoints that require admin role
	userAdminRouter := userRouter.NewRoute().Subrouter()
	userAdminRouter.Use(middleware.RBACMiddleware("admin"))
	
	// TODO: Add user management endpoints

	return r
}