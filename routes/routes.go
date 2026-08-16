package routes

import (
	"bytes"
	"encoding/json"
	"gamebiller/configs"
	"gamebiller/helpers"
	"gamebiller/models"
	"gamebiller/services/dashboard"
	"gamebiller/services/webapp"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

var JwtConfig echojwt.Config

func utils(a *echo.Group) {
	a.POST("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{"status": "UP", "time": timeNowStr()})
	})
	a.POST("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
}
func publicWebApp(a *echo.Group) {
	// Public
	a.POST("/register", webapp.Register)
	a.POST("/login", webapp.Login)
	a.POST("/reference-products", webapp.GetReferenceProduct)
	a.POST("/product-segments", webapp.GetProductSegment)
	a.POST("/product-segments/customer", webapp.GetProductSegmentByCustomer)
	a.POST("/popular-products", webapp.GetPopularProduct)
	a.POST("/payment-methods", webapp.GetPaymentMethod)
	a.POST("/callback/iak", webapp.IAKCallback)
	a.POST("/publicinquiry", webapp.InquiryUnSubscribe)
	a.POST("/publicpayment", webapp.PaymentUnSubscribe)
	a.POST("/publichistory", webapp.TransactionUnSubscribeHistory)
}
func privateWebApp(a *echo.Group) {
	// Private
	a.Use(echojwt.WithConfig(JwtConfig))
	a.POST("/inquiry", webapp.Inquiry)
	a.POST("/payment", webapp.Payment)
	a.POST("/history", webapp.TransactionHistory)
	a.POST("/update-password", webapp.UpdatePassword)
}

func dashboardRoutes(a *echo.Group) {
	a.Use(echojwt.WithConfig(JwtConfig))

	// Add custom middleware to enforce Admin roles
	a.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claimsVal := c.Get("user")
			if claimsVal == nil {
				return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-AUTH-403", nil))
			}
			token := claimsVal.(*jwt.Token)
			claims := token.Claims.(*models.JwtCustomClaims)
			if claims.RoleCode != "super_admin" && claims.RoleCode != "admin" && claims.RoleCode != "agent" && claims.RoleCode != "merchant" {
				return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-AUTH-403", nil))
			}
			return next(c)
		}
	})

	// Users CRUD
	a.POST("/users/list", dashboard.AdminGetUsers)
	a.POST("/users/create", dashboard.AdminCreateUser)
	a.POST("/users/update", dashboard.AdminUpdateUser)
	a.POST("/users/update-password", dashboard.AdminUpdatePassword)
	a.POST("/users/delete", dashboard.AdminDeleteUser)

	// Merchants CRUD
	a.POST("/merchants/list", dashboard.AdminGetMerchants)
	a.POST("/merchants/create", dashboard.AdminCreateMerchant)
	a.POST("/merchants/update", dashboard.AdminUpdateMerchant)
	a.POST("/merchants/delete", dashboard.AdminDeleteMerchant)

	// Products CRUD
	a.POST("/products/list", dashboard.AdminGetProducts)
	a.POST("/products/create", dashboard.AdminCreateProduct)
	a.POST("/products/update", dashboard.AdminUpdateProduct)
	a.POST("/products/delete", dashboard.AdminDeleteProduct)

	// Transactions CRUD
	a.POST("/transactions/list", dashboard.AdminGetTransactions)
	a.POST("/transactions/update", dashboard.AdminUpdateTransaction)

	// Roles CRUD
	a.POST("/roles/list", dashboard.AdminGetRoles)
	// a.POST("/roles/create", dashboard.AdminCreateRole)
	// a.POST("/roles/update", dashboard.AdminUpdateRole)
	// a.POST("/roles/delete", dashboard.AdminDeleteRole)

	// ModelHasRoles CRUD
	a.POST("/model-has-roles/list", dashboard.AdminGetModelHasRoles)
	a.POST("/model-has-roles/create", dashboard.AdminCreateModelHasRole)
	a.POST("/model-has-roles/update", dashboard.AdminUpdateModelHasRole)
	a.POST("/model-has-roles/delete", dashboard.AdminDeleteModelHasRole)

	// Providers CRUD
	a.POST("/providers/list", dashboard.AdminGetProviders)
	a.POST("/providers/create", dashboard.AdminCreateProvider)
	a.POST("/providers/update", dashboard.AdminUpdateProvider)
	a.POST("/providers/delete", dashboard.AdminDeleteProvider)

	// Product Categories CRUD
	a.POST("/product-categories/list", dashboard.AdminGetProductCategories)
	a.POST("/product-categories/create", dashboard.AdminCreateProductCategory)
	a.POST("/product-categories/update", dashboard.AdminUpdateProductCategory)
	a.POST("/product-categories/delete", dashboard.AdminDeleteProductCategory)

	// Product References CRUD
	a.POST("/product-references/list", dashboard.AdminGetProductReferences)
	a.POST("/product-references/create", dashboard.AdminCreateProductReference)
	a.POST("/product-references/update", dashboard.AdminUpdateProductReference)
	a.POST("/product-references/delete", dashboard.AdminDeleteProductReference)

	// Product Prefixes CRUD
	a.POST("/product-prefixes/list", dashboard.AdminGetProductPrefixes)
	a.POST("/product-prefixes/create", dashboard.AdminCreateProductPrefix)
	a.POST("/product-prefixes/update", dashboard.AdminUpdateProductPrefix)
	a.POST("/product-prefixes/delete", dashboard.AdminDeleteProductPrefix)

	// Product Providers CRUD
	a.POST("/product-providers/list", dashboard.AdminGetProductProviders)
	a.POST("/product-providers/create", dashboard.AdminCreateProductProvider)
	a.POST("/product-providers/update", dashboard.AdminUpdateProductProvider)
	a.POST("/product-providers/delete", dashboard.AdminDeleteProductProvider)

	// Product Segments CRUD
	a.POST("/product-segments/list", dashboard.AdminGetProductSegments)
	a.POST("/product-segments/create", dashboard.AdminCreateProductSegment)
	a.POST("/product-segments/update", dashboard.AdminUpdateProductSegment)
	a.POST("/product-segments/delete", dashboard.AdminDeleteProductSegment)

	// Product Masters CRUD
	a.POST("/product-masters/list", dashboard.AdminGetProductMasters)
	a.POST("/product-masters/create", dashboard.AdminCreateProductMaster)
	a.POST("/product-masters/update", dashboard.AdminUpdateProductMaster)
	a.POST("/product-masters/delete", dashboard.AdminDeleteProductMaster)

	// Agents CRUD
	a.POST("/agents/list", dashboard.AdminGetAgents)
	a.POST("/agents/create", dashboard.AdminCreateAgent)
	a.POST("/agents/update", dashboard.AdminUpdateAgent)
	a.POST("/agents/delete", dashboard.AdminDeleteAgent)

	// Segments CRUD
	a.POST("/segments/list", dashboard.AdminGetSegments)
	a.POST("/segments/create", dashboard.AdminCreateSegment)
	a.POST("/segments/update", dashboard.AdminUpdateSegment)
	a.POST("/segments/delete", dashboard.AdminDeleteSegment)

	// Payment Methods CRUD
	a.POST("/payment-methods/list", dashboard.AdminGetPaymentMethods)
	a.POST("/payment-methods/create", dashboard.AdminCreatePaymentMethod)
	a.POST("/payment-methods/update", dashboard.AdminUpdatePaymentMethod)
	a.POST("/payment-methods/delete", dashboard.AdminDeletePaymentMethod)

	// Payment Channels CRUD
	a.POST("/payment-channels/list", dashboard.AdminGetPaymentChannels)
	a.POST("/payment-channels/create", dashboard.AdminCreatePaymentChannel)
	a.POST("/payment-channels/update", dashboard.AdminUpdatePaymentChannel)
	a.POST("/payment-channels/delete", dashboard.AdminDeletePaymentChannel)

	// Payment Segments CRUD
	a.POST("/payment-segments/list", dashboard.AdminGetPaymentSegments)
	a.POST("/payment-segments/create", dashboard.AdminCreatePaymentSegment)
	a.POST("/payment-segments/update", dashboard.AdminUpdatePaymentSegment)
	a.POST("/payment-segments/delete", dashboard.AdminDeletePaymentSegment)

	// Saving Accounts CRUD
	a.POST("/saving-accounts/list", dashboard.AdminGetSavingAccounts)
	a.POST("/saving-accounts/detail", dashboard.AdminGetSavingAccountByID)
	a.POST("/saving-accounts/update", dashboard.AdminUpdateSavingAccount)

	// Saving Transactions CRUD
	a.POST("/saving-transactions/list", dashboard.AdminGetSavingTransactions)
}
func AppRoutes(e *echo.Echo) {
	// Print logs for every request, from header, request and response
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// --- 1. Log Request Headers ---
			reqHeaderBytes, _ := json.Marshal(c.Request().Header)
			log.Println("Request Endpoint :: ", c.Request().URL.Path)
			log.Println("Request Headers :: ", string(reqHeaderBytes))

			// --- 2. Safely Log Request Body ---
			var reqBodyBytes []byte
			if c.Request().Body != nil {
				// Read the body bytes
				reqBodyBytes, _ = io.ReadAll(c.Request().Body)
				// IMPORTANT: Restore the body so downstream handlers can read it too!
				c.Request().Body = io.NopCloser(bytes.NewBuffer(reqBodyBytes))
			}
			log.Println("Request Body :: ", string(reqBodyBytes))

			// --- 3. Set up Response Interceptor ---
			// We wrap the response writer to capture what goes out
			// --- 3. Set up Response Interceptor ---
			resBodyBuffer := new(bytes.Buffer)
			mw := io.MultiWriter(c.Response().Writer, resBodyBuffer)

			// Explicitly name the fields so the compiler maps them correctly
			writer := &responseBodyWriter{
				ResponseWriter: c.Response().Writer, // The original http.ResponseWriter
				Writer:         mw,                  // The multi-writer stream
			}
			c.Response().Writer = writer
			// Defer the response logging until after next(c) executes
			defer func() {
				resHeaderBytes, _ := json.Marshal(c.Response().Header())
				log.Println("Response Headers :: ", string(resHeaderBytes))
				log.Println("Response Body :: ", resBodyBuffer.String())
			}()
			return next(c)
		}
	})
	// Helper struct to intercept the response stream
	// 1. UTILS ROUTES
	utilsGroup := e.Group("/api/utils")

	utils(utilsGroup)

	// JWT Config
	JwtConfig = echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(models.JwtCustomClaims)
		},
		SigningKey: []byte(configs.APP_KEY),
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-AUTH-419", nil))
		},
	}

	// 2. WEBAPP ROUTES
	webappGroup := e.Group("/api/webapp")

	publicWebApp(webappGroup)

	// Authenticated (JWT)
	webappAuthGroup := e.Group("/api/webapp")
	privateWebApp(webappAuthGroup)
	// 3. DASHBOARD ROUTES
	dashboardGroup := e.Group("/api/dashboard")
	// Public Login
	dashboardGroup.POST("/login", dashboard.AdminLogin)
	// Authenticated Admin (JWT + Role validation)
	dashboardAuthGroup := e.Group("/api/dashboard")
	dashboardRoutes(dashboardAuthGroup)
}

func timeNowStr() string {
	return time.Now().Format(time.RFC3339)
}

type responseBodyWriter struct {
	http.ResponseWriter
	io.Writer
}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
