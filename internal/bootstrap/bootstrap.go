package bootstrap

import (
	"context"
	"fmt"
	"log"

	"ecommerce-go/internal/config"
	"ecommerce-go/internal/db"
	"ecommerce-go/internal/module/admin"
	"ecommerce-go/internal/module/order"
	"ecommerce-go/internal/module/product"
	"ecommerce-go/internal/module/user"
	"ecommerce-go/internal/router"
)

// Run wires dependencies and starts the HTTP server.
func Run() {
	cfg := config.Load()

	dbConn, err := db.New(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}
	if err := dbConn.InitSchema(context.Background()); err != nil {
		log.Fatal(err)
	}

	r := router.NewRouter(cfg)
	api := r.Group("/api/v1")

	// Module wiring (DI): controller -> service -> repository (MySQL from .env).
	userRepo := user.NewRepository(dbConn.SQL)
	userSvc := user.NewService(userRepo)
	userCtrl := user.NewController(userSvc)
	user.RegisterUserRoutes(api, userCtrl)

	orderRepo := order.NewRepository(dbConn.SQL)
	orderSvc := order.NewService(orderRepo)
	orderCtrl := order.NewController(orderSvc)
	order.RegisterOrderRoutes(api, orderCtrl)

	productRepo := product.NewRepository(dbConn.SQL)
	productSvc := product.NewService(productRepo)
	productCtrl := product.NewController(productSvc)
	product.RegisterProductRoutes(api, productCtrl)

	adminRepo := admin.NewRepository(dbConn.SQL)
	adminSvc := admin.NewService(adminRepo)
	adminCtrl := admin.NewController(adminSvc)
	admin.RegisterAdminRoutes(api, adminCtrl)

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("HTTP server listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
