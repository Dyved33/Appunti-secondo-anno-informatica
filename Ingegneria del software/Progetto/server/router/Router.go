package router

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/tommasogarofani/ingsw_25-26_group-11/server/auth"
	"github.com/tommasogarofani/ingsw_25-26_group-11/server/category"
	"github.com/tommasogarofani/ingsw_25-26_group-11/server/domain"
	"github.com/tommasogarofani/ingsw_25-26_group-11/server/middleware"
	"github.com/tommasogarofani/ingsw_25-26_group-11/server/order"
	"github.com/tommasogarofani/ingsw_25-26_group-11/server/product"
	"github.com/tommasogarofani/ingsw_25-26_group-11/server/user"
)

func NewRouter(db *sql.DB) chi.Router {
	r := chi.NewRouter()

	// --- Helper transazionale condiviso da tutti i service ---
	tx := domain.NewTransactionHelper(db)

	// --- User & Auth ---
	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo, tx)
	userController := user.NewUserController(userService)

	sessionRepo := auth.NewSessionRepository(db)
	authService := auth.NewAuthService(userService, sessionRepo)
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
		r.Use(middleware.NewAuthMiddleware(sessionRepo))

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

			r.Get("/inventory-report", func(w http.ResponseWriter, req *http.Request) {
				result, err := productController.ProductController_GetInventoryReport(req.Context())
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

			r.Put("/{id}/price", func(w http.ResponseWriter, req *http.Request) {
				id, err := strconv.Atoi(chi.URLParam(req, "id"))
				if err != nil {
					http.Error(w, "id prodotto non valido", http.StatusBadRequest)
					return
				}

				var input product.PriceUpdateInputDTO
				if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
					http.Error(w, "richiesta non valida", http.StatusBadRequest)
					return
				}

				if err := productController.ProductController_UpdateSellPrice(req.Context(), id, input.PrezzoVend); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusNoContent)
			})

			r.Delete("/{id}", func(w http.ResponseWriter, req *http.Request) {
				id, err := strconv.Atoi(chi.URLParam(req, "id"))
				if err != nil {
					http.Error(w, "id prodotto non valido", http.StatusBadRequest)
					return
				}
				if err := productController.ProductController_RemoveProduct(req.Context(), id); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusNoContent)
			})

			r.Get("/barcode/{barcode}", func(w http.ResponseWriter, req *http.Request) {
				barcode := chi.URLParam(req, "barcode")
				result, err := productController.ProductController_GetProductByBarcode(req.Context(), barcode)
				if err != nil {
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(result)
			})

		})

		r.Route("/api/v1/orders", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, req *http.Request) {
				result, err := orderController.OrderController_GetOrdiniUscita(req.Context())
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(result)
			})

			r.Get("/uscita", func(w http.ResponseWriter, req *http.Request) {
				result, err := orderController.OrderController_GetOrdiniUscita(req.Context())
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(result)
			})

			r.Post("/crea-uscita", func(w http.ResponseWriter, req *http.Request) {
				var input order.OrdineUscitaDTO
				if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
					http.Error(w, "richiesta non valida", http.StatusBadRequest)
					return
				}
				input.CreatedByID, _ = req.Context().Value(middleware.UserIDKey).(int)

				if err := orderController.OrderController_RegisterExitOrder(req.Context(), input); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusCreated)
			})

			r.Post("/crea-entrata", func(w http.ResponseWriter, req *http.Request) {
				var input order.OrdineEntrataDTO
				if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
					http.Error(w, "richiesta non valida", http.StatusBadRequest)
					return
				}
				input.CreatedByID, _ = req.Context().Value(middleware.UserIDKey).(int)

				if err := orderController.OrderController_RegisterEntryOrder(req.Context(), input); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusCreated)
			})

			r.Get("/entrata", func(w http.ResponseWriter, req *http.Request) {
				result, err := orderController.OrderController_GetOrdiniEntrata(req.Context())
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(result)
			})

			r.Delete("/delete-entrata/{ordine_id}", func(w http.ResponseWriter, req *http.Request) {
				var id = chi.URLParam(req, "ordine_id")
				ordineID, err := strconv.Atoi(id)
				if err != nil {
					http.Error(w, "id ordine non valido", http.StatusBadRequest)
					return
				}

				err = orderController.OrderController_DeleteOrdineEntrata(req.Context(), ordineID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusNoContent)
			})

			r.Delete("/delete-uscita/{ordine_id}", func(w http.ResponseWriter, req *http.Request) {
				var id = chi.URLParam(req, "ordine_id")
				ordineID, err := strconv.Atoi(id)
				if err != nil {
					http.Error(w, "id ordine non valido", http.StatusBadRequest)
					return
				}

				err = orderController.OrderController_DeleteOrdineUscita(req.Context(), ordineID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusNoContent)
			})

			r.Put("/update-ordine-uscita/{ordine_id}", func(w http.ResponseWriter, req *http.Request) {
				id := chi.URLParam(req, "ordine_id")

				ordineID, err := strconv.Atoi(id)
				if err != nil {
					http.Error(w, "id ordine non valido", http.StatusBadRequest)
					return
				}

				var input order.OrdineUscitaDTO
				if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
					http.Error(w, "richiesta non valida", http.StatusBadRequest)
					return
				}

				err = orderController.OrderController_UpdateOrdineUscita(
					req.Context(),
					ordineID,
					input,
				)

				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusNoContent)
			})

			r.Delete("/delete-reso/{reso_id}", func(w http.ResponseWriter, req *http.Request) {
				var id = chi.URLParam(req, "reso_id")
				resoID, err := strconv.Atoi(id)
				if err != nil {
					http.Error(w, "id reso non valido", http.StatusBadRequest)
					return
				}

				err = orderController.OrderController_DeleteReso(req.Context(), resoID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusNoContent)
			})

			r.Get("/resi", func(w http.ResponseWriter, req *http.Request) {
				result, err := orderController.OrderController_GetReturns(req.Context())
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(result)
			})

			r.Post("/crea-reso", func(w http.ResponseWriter, req *http.Request) {
				var input order.ResoDTO
				if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
					http.Error(w, "richiesta non valida", http.StatusBadRequest)
					return
				}
				input.CreatedByID, _ = req.Context().Value(middleware.UserIDKey).(int)

				if err := orderController.OrderController_CreateReturn(req.Context(), input); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusCreated)
			})

			r.Put("/update-resi", func(w http.ResponseWriter, req *http.Request) {
				var input order.ResoDTO
				if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
					http.Error(w, "richiesta non valida", http.StatusBadRequest)
					return
				}

				err := orderController.OrderController_UpdateStatoReso(req.Context(), input.ID, input.Stato)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusNoContent)
			})

			r.Get("/invoices", func(w http.ResponseWriter, req *http.Request) {
				result, err := orderController.OrderController_GetInvoices(req.Context())
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(result)
			})

			r.Post("/create-credit-notes/{order_id}", func(w http.ResponseWriter, req *http.Request) {
				id, err := strconv.Atoi(chi.URLParam(req, "order_id"))
				if err != nil {
					http.Error(w, "id ordine non valido", http.StatusBadRequest)
					return
				}
				result, err := orderController.OrderController_GeneraNotaCredito(req.Context(), id)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(result)
			})

			r.Get("/credit-notes", func(w http.ResponseWriter, req *http.Request) {
				result, err := orderController.OrderController_GetCreditNotes(req.Context())
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(result)
			})

			r.Post("/{order_id}/create-ddt", func(w http.ResponseWriter, req *http.Request) {
				var input order.DDTEntrataDTO
				var id = chi.URLParam(req, "order_id")
				orderID, err := strconv.Atoi(id)
				if err != nil {
					http.Error(w, "id ordine non valido", http.StatusBadRequest)
					return
				}
				if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
					http.Error(w, "richiesta non valida", http.StatusBadRequest)
					return
				}

				result, err := orderController.OrderController_GeneraDDT(req.Context(), orderID, input)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(result)
			})

			r.Post("/create-fattura/{order_id}", func(w http.ResponseWriter, req *http.Request) {
				id, err := strconv.Atoi(chi.URLParam(req, "order_id"))
				if err != nil {
					http.Error(w, "id ordine non valido", http.StatusBadRequest)
					return
				}

				result, err := orderController.OrderController_GeneraFattura(req.Context(), id)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(result)
			})

			r.Delete("/delete-fattura/{fattura_id}", func(w http.ResponseWriter, req *http.Request) {
				id, err := strconv.Atoi(chi.URLParam(req, "fattura_id"))
				if err != nil {
					http.Error(w, "id fattura non valido", http.StatusBadRequest)
					return
				}

				if err := orderController.OrderController_DeleteFattura(req.Context(), id); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusNoContent)
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

			r.Post("/", func(w http.ResponseWriter, req *http.Request) {
				var input category.CategoryInputDTO
				if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
					http.Error(w, "richiesta non valida", http.StatusBadRequest)
					return
				}
				result, err := categoryController.CategoryController_Create(req.Context(), input)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(result)
			})

			r.Delete("/{id}", func(w http.ResponseWriter, req *http.Request) {
				id, err := strconv.Atoi(chi.URLParam(req, "id"))
				if err != nil {
					http.Error(w, "id categoria non valido", http.StatusBadRequest)
					return
				}
				if err := categoryController.CategoryController_Delete(req.Context(), id); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusNoContent)
			})
		})

		r.Route("/api/v1/users", func(r chi.Router) {
			r.Post("/", func(w http.ResponseWriter, req *http.Request) {
				var input struct {
					Username string `json:"username"`
					Password string `json:"password"`
					Ruolo    string `json:"ruolo"`
				}
				if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				passHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				utente := domain.Utente{
					Username:         input.Username,
					PasswordHash:     string(passHash),
					Ruolo:            domain.Ruolo(input.Ruolo),
					TentativiFalliti: 0,
					Bloccato:         false,
					CreatedAt:        time.Now(),
				}

				if err := userRepo.UserRepository_Create(req.Context(), &utente); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(utente)
			})

			r.Get("/{username}", func(w http.ResponseWriter, req *http.Request) {
				username := chi.URLParam(req, "username")
				result, err := userRepo.UserRepository_GetByUsername(req.Context(), username)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(result)
			})

			r.Put("/{username}", func(w http.ResponseWriter, req *http.Request) {
				username := chi.URLParam(req, "username")
				existing, err := userRepo.UserRepository_GetByUsername(req.Context(), username)
				if err != nil {
					http.Error(w, "utente non trovato", http.StatusNotFound)
					return
				}

				var payload struct {
					Username         *string `json:"username"`
					Ruolo            *string `json:"ruolo"`
					Bloccato         *bool   `json:"bloccato"`
					TentativiFalliti *int    `json:"tentativi_falliti"`
					Password         *string `json:"password"`
				}
				if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				if payload.Username != nil && *payload.Username != "" {
					existing.Username = *payload.Username
				}
				if payload.Ruolo != nil && *payload.Ruolo != "" {
					existing.Ruolo = domain.Ruolo(*payload.Ruolo)
				}
				if payload.Bloccato != nil {
					existing.Bloccato = *payload.Bloccato
					if !existing.Bloccato {
						existing.TentativiFalliti = 0
					}
				}
				if payload.TentativiFalliti != nil {
					existing.TentativiFalliti = *payload.TentativiFalliti
				}
				if payload.Password != nil && *payload.Password != "" {
					hash, err := bcrypt.GenerateFromPassword([]byte(*payload.Password), bcrypt.DefaultCost)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					existing.PasswordHash = string(hash)
				}

				err = userRepo.UserRepository_Update(req.Context(), existing)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(existing)
			})

			r.Delete("/{username}", func(w http.ResponseWriter, req *http.Request) {
				username := chi.URLParam(req, "username")
				utente, err := userRepo.UserRepository_GetByUsername(req.Context(), username)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				err = userRepo.UserRepository_Delete(req.Context(), utente)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusNoContent)
			})
		})
	})
	return r
}
