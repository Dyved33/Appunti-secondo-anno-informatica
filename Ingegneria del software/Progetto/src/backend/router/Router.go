package router

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/tommasogarofani/ingsw_25-26_group-11/src/backend/auth"
	"github.com/tommasogarofani/ingsw_25-26_group-11/src/backend/category"
	"github.com/tommasogarofani/ingsw_25-26_group-11/src/backend/domain"
	"github.com/tommasogarofani/ingsw_25-26_group-11/src/backend/middleware"
	"github.com/tommasogarofani/ingsw_25-26_group-11/src/backend/order"
	"github.com/tommasogarofani/ingsw_25-26_group-11/src/backend/product"
	"github.com/tommasogarofani/ingsw_25-26_group-11/src/backend/user"
)

func NewRouter(db *sql.DB) chi.Router {
	r := chi.NewRouter()

	// --- Helper transazionale condiviso da tutti i service ---
	tx := domain.NewTransactionHelper(db) 

	// --- User ---
	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo, tx)
	userController := user.NewUserController(userService)

	authService := auth.NewAuthService(userService)
	authController := auth.NewAuthController(authService)

	// --- Product ---
	productRepo := product.NewProductRepository(db)
	productService := product.NewProductService(productRepo, tx)
	productController := product.NewProductController(productService)

	// --- Order ---
	orderRepo := order.NewOrderRepository(db)
	orderService := order.NewOrderService(orderRepo, productRepo, tx) 
	orderController := order.NewOrderController(orderService)

	// --- Category --- 
	categoryRepo := category.NewCategoryRepository(db)
	categoryService := category.NewCategoryService(categoryRepo, tx)
	categoryController := category.NewCategoryController(categoryService)

	_ = userController // usato solo per Login, non per rotte dirette 

	// --- Rotta pubblica: login ---
	r.Post("/api/v1/login", func(w http.ResponseWriter, req *http.Request) {
		var body auth.LoginRequestDTO
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			http.Error(w, "richiesta non valida", http.StatusBadRequest)
			return
		}
		result, err := authController.AuthController_Login(req.Context(), body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	// --- Rotte protette ---
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		r.Route("/api/v1/products", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, req *http.Request) {
				result, err := productController.ProductController_GetProducts(req.Context())
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(result)
			})

			r.Post("/", func(w http.ResponseWriter, req *http.Request) {
				var input product.ProductInputDTO
				if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
					http.Error(w, "richiesta non valida", http.StatusBadRequest)
					return
				}
				result, err := productController.ProductController_AddProduct(req.Context(), input)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(result)
			})
		})

		r.Route("/api/v1/orders", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, req *http.Request) {
				result, err := orderController.OrderController_GetAll(req.Context())
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(result)
			})
		})

		r.Route("/api/v1/categories", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, req *http.Request) {
				result, err := categoryController.CategoryController_GetAll(req.Context())
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(result)
			})
		})
        //Per ora non ho messo una rotta per gli utenti perché, a parte il login, non ci sono altre operazioni da fare sugli utenti
	})
	return r
}