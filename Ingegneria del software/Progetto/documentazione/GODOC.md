# Documentazione Tecnica Godoc 

Per avere una documentazione precisa di tutti i metodi descritti nel programma, abbiamo fatto uso dello strumento proprietario di golanf pkgsite. Una volta installato con il comando `go install golang.org/x/pkgsite/cmd/pkgsite@latest`, sarà necessario eseguire `$(go env GOPATH)/bin/pkgsite -http=localhost:8081` (la porta indicata è diversa da quella usata dai client per comunicare con il server perché quest'ultima è occupata). A questo punto sarà sufficiente scrivere nella barra di ricerca del proprio browser `http://localhost:8081` e da questo verrà aperta una pagina web con la possibilità di ricercare una delle sotto-cartelle del progetto. Una volta selezionata la sottocartella desiderata verrà mostrata una descrizione dettagliata e indicizzata delle sezioni di codici commentate. 

Alternativamente si può consultare questa analisi del progetto, la quale si concentra sulle rotte e sulle api usate:

## Indice
- [Documentazione Tecnica Godoc](#documentazione-tecnica-godoc)
  - [Indice](#indice)
  - [1. Introduzione e Riferimenti Architetturali](#1-introduzione-e-riferimenti-architetturali)
  - [2. Package `router` (`server/router`)](#2-package-router-serverrouter)
    - [Panoramica e Design](#panoramica-e-design)
    - [Inizializzazione e Dependency Injection](#inizializzazione-e-dependency-injection)
      - [Descrizione](#descrizione)
      - [Parametri](#parametri)
      - [Workflow di Inizializzazione:](#workflow-di-inizializzazione)
    - [Specifica Completa degli Endpoint REST](#specifica-completa-degli-endpoint-rest)
      - [1. Autenticazione (Pubblica)](#1-autenticazione-pubblica)
      - [2. Catalogo Prodotti \& Giacenze (`/api/v1/products`)](#2-catalogo-prodotti--giacenze-apiv1products)
      - [3. Ordini, DDT, Fatture e Resi (`/api/v1/orders`)](#3-ordini-ddt-fatture-e-resi-apiv1orders)
      - [4. Categorie Merceologiche (`/api/v1/categories`)](#4-categorie-merceologiche-apiv1categories)
      - [5. Gestione Utenti (`/api/v1/users`)](#5-gestione-utenti-apiv1users)
  - [3. Package `api` (`client/api`)](#3-package-api-clientapi)
    - [Panoramica e Pattern Client](#panoramica-e-pattern-client)
    - [Tipi, Costanti e Strutture DTO](#tipi-costanti-e-strutture-dto)
      - [Costanti di Stato](#costanti-di-stato)
      - [Modelli DTO Esportati](#modelli-dto-esportati)
    - [Client API ed Autenticazione](#client-api-ed-autenticazione)
      - [`type APIClient`](#type-apiclient)
      - [Costruttore e Metodi di Sessione](#costruttore-e-metodi-di-sessione)
    - [Metodi di Dominio Esportati](#metodi-di-dominio-esportati)
      - [1. Catalogo \& Prodotti (`products.go`)](#1-catalogo--prodotti-productsgo)
      - [2. Ordini, DDT, Fatture e Resi (`orders.go`)](#2-ordini-ddt-fatture-e-resi-ordersgo)
      - [3. Gestione Utenti (`users.go`)](#3-gestione-utenti-usersgo)
  - [4. Integrazione con la Documentazione di Progetto (`documentazione/`)](#4-integrazione-con-la-documentazione-di-progetto-documentazione)
    - [Decision Log Mapping](#decision-log-mapping)
    - [Garanzie Transazionali e di Sicurezza](#garanzie-transazionali-e-di-sicurezza)

---

## 1. Introduzione e Riferimenti Architetturali

Il sistema **GiaLoSai** adotta un'architettura **Layered** con riferimenti all'architettura **client-server**. Adotta una separazione netta tra il backend RESTful scritto in Go e l'interfaccia utente TUI (Text-based User Interface) basata sul framework Bubble Tea.

Come sancito nel *Design Document* e nel *Decision Log*:
- **DC-001 (Framework TUI):** Il frontend opera in terminale con architettura Model-View-Update (MVU / Elm architecture).
- **DC-002 (Router HTTP):** Il backend usa `chi/v5` per il routing, preservando piena compatibilità con `net/http` e la gestione a catena dei middleware.
- **DC-003 & DC-004 (MySQL & Repository Pattern):** I dati risiedono su MySQL 8 (InnoDB con vincoli ACID); ogni package applicativo espone un'interfaccia di repository e service separata dai controller di trasporto HTTP.

---

## 2. Package `router` (`server/router`)

```go
package router // import "github.com/tommasogarofani/ingsw_25-26_group-11/server/router"
```

### Panoramica e Design
Il package `router` funge da punto centrale di instradamento e composizione del server. Il suo compito principale è definire l'albero di routing HTTP, montare i middleware di autenticazione e iniettare le dipendenze nei rispettivi controller.

### Inizializzazione e Dependency Injection

```go
func NewRouter(db *sql.DB) chi.Router
```

#### Descrizione
`NewRouter` istanzia tutti i repository, i service di dominio, i controller HTTP e l'helper transazionale, configurando e restituendo il router `chi.Router`.

#### Parametri
- `db *sql.DB`: Connessione al pool del database MySQL 8.

#### Workflow di Inizializzazione:
1. **Helper Transazionale:** Inizializzazione di `domain.NewTransactionHelper(db)`, propagato a tutti i service per garantire transazioni atomiche cross-repository tramite `context.Context`.
2. **Layer di Dominio:**
   - `User` & `Auth`: `user.NewUserRepository`, `user.NewUserService`, `auth.NewSessionRepository`, `auth.NewAuthService`, `auth.NewAuthController`.
   - `Product`: `product.NewProductRepository`, `product.NewProductService`, `product.NewProductController`.
   - `Order`: `order.NewOrderRepository`, `order.NewOrderService`, `order.NewOrderController`.
   - `Category`: `category.NewCategoryRepository`, `category.NewCategoryService`, `category.NewCategoryController`.
3. **Definizione Rotte:**
   - Rotta pubblica per `/api/v1/login`.
   - Gruppo protetto via middleware JWT (`middleware.NewAuthMiddleware(sessionRepo)`).

---

### Specifica Completa degli Endpoint REST

#### 1. Autenticazione (Pubblica)
- **`POST /api/v1/login`**
  - **Controller:** `authController.AuthController_Login`
  - **Body Richiesta:** `auth.LoginRequestDTO` (`{"username": "...", "password": "..."}`)
  - **Risposta (200 OK):** `auth.LoginResponseDTO` (`{"token": "JWT...", "role": "..."}`)
  - **Errori:** `400 Bad Request` (payload malformato), `401 Unauthorized` (credenziali errate o utente bloccato).

---

#### 2. Catalogo Prodotti & Giacenze (`/api/v1/products`)
*Tutti gli endpoint richiedono header `Authorization: Bearer <token>`.*

- **`GET /api/v1/products`**
  - **Descrizione:** Recupera l'elenco di tutti i prodotti in catalogo.
  - **Risposta (200 OK):** JSON array di `product.ProductOutputDTO`.
  - **Errori:** `500 Internal Server Error`.

- **`GET /api/v1/products/inventory-report`**
  - **Descrizione:** Genera il report aggregato delle giacenze, inclusi prodotti sottoscorta e valorizzazione totale del magazzino.
  - **Risposta (200 OK):** JSON `domain.InventoryReport`.
  - **Errori:** `500 Internal Server Error`.

- **`POST /api/v1/products`**
  - **Descrizione:** Registra un nuovo prodotto anagrafico con codice a barre EAN-13 validato.
  - **Body Richiesta:** `product.ProductInputDTO`.
  - **Risposta (200 OK):** `product.ProductOutputDTO` creato.
  - **Errori:** `400 Bad Request`, `500 Internal Server Error`.

- **`PUT /api/v1/products/{id}/price`**
  - **Descrizione:** Aggiorna il prezzo di vendita per il prodotto indicato.
  - **URL Param:** `id` (intero).
  - **Body Richiesta:** `product.PriceUpdateInputDTO` (`{"prezzo_vendita": float64}`).
  - **Risposta (204 No Content)**.
  - **Errori:** `400 Bad Request`, `500 Internal Server Error`.

- **`DELETE /api/v1/products/{id}`**
  - **Descrizione:** Cancella un prodotto dal catalogo.
  - **URL Param:** `id` (intero).
  - **Risposta (204 No Content)**.
  - **Errori:** `400 Bad Request`, `500 Internal Server Error`.

- **`GET /api/v1/products/barcode/{barcode}`**
  - **Descrizione:** Ricerca rapida di un articolo tramite codice a barre (utilizzato dalla funzione Scanner).
  - **URL Param:** `barcode` (stringa formato EAN-13).
  - **Risposta (200 OK):** `product.ProductOutputDTO`.
  - **Errori:** `404 Not Found`.

---

#### 3. Ordini, DDT, Fatture e Resi (`/api/v1/orders`)
*Tutti gli endpoint richiedono header `Authorization: Bearer <token>`.*

- **`GET /api/v1/orders` / `GET /api/v1/orders/uscita`**
  - **Descrizione:** Restituisce la lista degli ordini cliente di uscita.
  - **Risposta (200 OK):** JSON array di `order.OrderOutputDTO`.

- **`POST /api/v1/orders/crea-uscita`**
  - **Descrizione:** Registra un nuovo ordine cliente di uscita merci, verificando disponibilità e decrementando le giacenze in transazione. Inietta il `CreatedByID` dal contesto JWT.
  - **Body Richiesta:** `order.OrdineUscitaDTO`.
  - **Risposta (201 Created)**.
  - **Errori:** `400 Bad Request`, `500 Internal Server Error` (es. giacenza insufficiente).

- **`GET /api/v1/orders/entrata`**
  - **Descrizione:** Restituisce l'elenco degli ordini fornitore di entrata merci.
  - **Risposta (200 OK):** JSON array di `order.OrderOutputDTO`.

- **`POST /api/v1/orders/crea-entrata`**
  - **Descrizione:** Registra un arrivo merci da fornitore incrementando atomicamente le giacenze a magazzino.
  - **Body Richiesta:** `order.OrdineEntrataDTO`.
  - **Risposta (201 Created)**.

- **`PUT /api/v1/orders/update-ordine-uscita/{ordine_id}`**
  - **Descrizione:** Aggiorna lo stato o i dettagli di un ordine in uscita.
  - **Risposta (204 No Content)**.

- **`DELETE /api/v1/orders/delete-uscita/{ordine_id}` & `/delete-entrata/{ordine_id}`**
  - **Descrizione:** Cancellazione ordini specificati per ID.
  - **Risposta (204 No Content)**.

- **`POST /api/v1/orders/{order_id}/create-ddt`**
  - **Descrizione:** Genera il Documento di Trasporto (DDT) ed evade l'ordine di uscita collegato.
  - **Body Richiesta:** `order.DDTEntrataDTO`.
  - **Risposta (200 OK):** `order.DDTOutputDTO`.

- **`GET /api/v1/orders/invoices`**
  - **Descrizione:** Elenco di tutte le fatture emesse.
  - **Risposta (200 OK):** JSON array di `order.FatturaOutputDTO`.

- **`POST /api/v1/orders/create-fattura/{order_id}`**
  - **Descrizione:** Emette la fattura fiscale a partire da un ordine di uscita evaso.
  - **Risposta (200 OK):** `order.FatturaOutputDTO`.

- **`DELETE /api/v1/orders/delete-fattura/{fattura_id}`**
  - **Descrizione:** Rimuove la fattura indicata.
  - **Risposta (204 No Content)**.

- **`GET /api/v1/orders/credit-notes`**
  - **Descrizione:** Elenco delle note di credito registrate.
  - **Risposta (200 OK):** JSON array di `order.FatturaOutputDTO`.

- **`POST /api/v1/orders/create-credit-notes/{order_id}`**
  - **Descrizione:** Emette nota di credito per storno merci o rettifica.
  - **Risposta (200 OK):** `order.FatturaOutputDTO`.

- **`GET /api/v1/orders/resi`**
  - **Descrizione:** Elenco delle richieste di reso merce.
  - **Risposta (200 OK):** JSON array di `order.ResoOutputDTO`.

- **`POST /api/v1/orders/crea-reso`**
  - **Descrizione:** Inserisce una nuova segnalazione di reso.
  - **Body Richiesta:** `order.ResoDTO`.
  - **Risposta (201 Created)**.

- **`PUT /api/v1/orders/update-resi`**
  - **Descrizione:** Modifica lo stato di lavorazione del reso (`da_verificare`, `accettato`, `rifiutato`), con reintegro giacenze in caso di accettazione (il controller mappa `approvato` sull'accettazione).
  - **Risposta (204 No Content)**.

- **`DELETE /api/v1/orders/delete-reso/{reso_id}`**
  - **Descrizione:** Cancella un record di reso.
  - **Risposta (204 No Content)**.

---

#### 4. Categorie Merceologiche (`/api/v1/categories`)
*Tutti gli endpoint richiedono header `Authorization: Bearer <token>`.*

- **`GET /api/v1/categories`**: Recupera tutte le categorie.
- **`POST /api/v1/categories`**: Crea una nuova categoria (`category.CategoryInputDTO`).
- **`DELETE /api/v1/categories/{id}`**: Elimina una categoria dato l'ID.

---

#### 5. Gestione Utenti (`/api/v1/users`)
*Tutti gli endpoint richiedono header `Authorization: Bearer <token>`.*

- **`POST /api/v1/users`**: Crea un nuovo utente cifrando la password con bcrypt.
- **`GET /api/v1/users/{username}`**: Restituisce le informazioni del profilo utente.
- **`PUT /api/v1/users/{username}`**: Modifica credenziali, ruolo, sblocco account o azzeramento tentativi falliti.
- **`DELETE /api/v1/users/{username}`**: Elimina l'utente specificato.

---

## 3. Package `api` (`client/api`)

```go
package api // import "github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
```

### Panoramica e Pattern Client
Il package `api` implementa l'SDK client HTTP REST per l'interfaccia terminale. Funge da intermediario tra i modelli TUI di Bubble Tea e gli endpoint del server, occupandosi di:
1. Gestire la sessione utente (conservazione token JWT e ruolo in memoria).
2. Iniettare automaticamente l'header `Authorization: Bearer <token>`.
3. Mappare i DTO server nelle strutture dati di frontend.

---

### Tipi, Costanti e Strutture DTO

#### Costanti di Stato
```go
type StatoOrdine string

const (
    StatoConfermato StatoOrdine = "confermato" // Ordine registrato, non evaso
    StatoEvaso       StatoOrdine = "evaso"       // Ordine evaso con DDT associato
    StatoAnnullato   StatoOrdine = "annullato"   // Ordine annullato
)

type StatoReso string

const (
    StatoDaVerificare StatoReso = "da_verificare" // Reso in attesa di perizia
    StatoAccettato    StatoReso = "accettato"    // Reso approvato con ripristino scorte
    StatoRifiutato    StatoReso = "rifiutato"   // Reso respinto
)
```

#### Modelli DTO Esportati
- **`Prodotto`**: Rappresenta l'articolo mostrato nelle tabelle catalogo e giacenze.
- **`Ordine` & `OrdineEntrata`**: Strutture ordine con righe aggregate, cliente/fornitore e totale.
- **`Fattura`**: Modello per fatture di vendita e note di credito.
- **`Reso`**: Modello per tracciamento pratiche di reso cliente/magazzino.
- **`User`**: Dati anagrafici e stato di blocco dell'operatore.
- **`GiacenzaReportItem`**: Elemento aggregato per il report scorte e sottoscorta.

---

### Client API ed Autenticazione

#### `type APIClient`
```go
type APIClient struct {
    indirizzoServer string
    token           string
    role            string
    httpClient      *http.Client
}
```

#### Costruttore e Metodi di Sessione
```go
func New(url string) *APIClient
```
Inizializza un'istanza di `APIClient` configurata sull'URL base fornito.

```go
func (c *APIClient) Token() string
```
Restituisce il token JWT memorizzato nella sessione corrente.

```go
func (c *APIClient) Role() string
```
Restituisce la stringa del ruolo autenticato (`magazziniere`, `contabile`, `admin`).

```go
func (c *APIClient) Login(username, password string) error
```
Esegue la chiamata HTTP POST a `/api/v1/login`, decodifica la risposta ed estrae token e ruolo memorizzandoli nel client.

---

### Metodi di Dominio Esportati

#### 1. Catalogo & Prodotti (`products.go`)
```go
// Recupera tutti i prodotti associandoli al nome della rispettiva categoria.
func (c *APIClient) GetCatalogo() ([]Prodotto, []string, error)

// Aggiorna il prezzo di vendita per il prodotto indicato.
func (c *APIClient) AggiornaPrezzo(prodottoID int, nuovoPrezzo float64) error

// Recupera un prodotto per codice interno.
func (c *APIClient) GetProductByCode(code string) (Prodotto, error)

// Recupera un prodotto tramite lettura barcode EAN-13.
func (c *APIClient) GetProductByBarcode(barcode string) (Prodotto, error)

// Crea una nuova voce nel catalogo prodotti.
func (c *APIClient) CreaProdotto(input ProdottoInput) (Prodotto, error)

// Rimuove un articolo dal catalogo.
func (c *APIClient) DeleteProdotto(prodottoID int) error

// Genera i dati aggregati del report giacenze e sottoscorta.
func (c *APIClient) GetReportGiacenze() ([]GiacenzaReportItem, error)
```

#### 2. Ordini, DDT, Fatture e Resi (`orders.go`)
```go
// Recupera tutti gli ordini di uscita clienti.
func (c *APIClient) GetOrdini() ([]Ordine, error)

// Registra un nuovo ordine cliente di uscita.
func (c *APIClient) CreaOrdineUscita(input OrdineUscitaInput) error

// Recupera gli ordini di entrata da fornitore.
func (c *APIClient) GetOrdiniEntrata() ([]OrdineEntrata, error)

// Registra un nuovo carico merci fornitore.
func (c *APIClient) CreaOrdineEntrata(input OrdineEntrataInput) error

// Evade un ordine di uscita emettendo il Documento di Trasporto (DDT).
func (c *APIClient) EvadiOrdine(ordineID int, ddt DDTInput) error

// Recupera fatture ordinarie e note di credito.
func (c *APIClient) GetFatture() ([]Fattura, error)

// Emette la fattura per un ordine evaso.
func (c *APIClient) CreaFattura(ordineID int) error

// Emette una nota di credito per rettifica/reso.
func (c *APIClient) CreaNotaCredito(ordineID int) error

// Recupera l'elenco dei resi.
func (c *APIClient) GetResi() ([]Reso, error)

// Registra una richiesta di reso.
func (c *APIClient) CreaReso(input ResoInput) error

// Aggiorna lo stato di una pratica di reso.
func (c *APIClient) AggiornaStatoReso(resoID int, stato StatoReso) error
```

#### 3. Gestione Utenti (`users.go`)
```go
// Recupera le informazioni di un utente per username.
func (c *APIClient) GetUsersByUsername(username string) (User, error)

// Registra un nuovo utente nel sistema.
func (c *APIClient) CreateUser(input UserCreateInput) (User, error)

// Aggiorna username, password, ruolo o stato di un utente.
func (c *APIClient) UpdateUser(username string, input UserUpdateInput) (User, error)

// Blocca o sblocca un utente azzerando i tentativi falliti.
func (c *APIClient) ToggleUserBlock(username string, bloccato bool) (User, error)

// Cancella un account utente.
func (c *APIClient) DeleteUser(username string) error
```

---

## 4. Integrazione con la Documentazione di Progetto (`documentazione/`)

### Decision Log Mapping

| Decisione | Ambito | Impatto su `server/router` e `client/api` |
|---|---|---|
| **DC-001** | Framework TUI | `client/api` incapsula tutte le chiamate REST necessarie ai modelli di pagina Bubble Tea. |
| **DC-002** | Router `chi` | `server/router` definisce le 32 rotte modulari con handler standard `net/http` e URL matching. |
| **DC-003** | MySQL 8 | La persistenza di backend è isolata dai repository; il client comunica solo tramite payload JSON. |
| **DC-004** | Layered + Repository | Disaccoppiamento tra controller (HTTP), service (business logic) e repository (SQL). |
| **DC-005** | Autenticazione JWT | Whitelist sessioni hashato SHA-256 nel DB (`server/auth`); Bearer token conservato in `APIClient`. |

### Garanzie Transazionali e di Sicurezza
1. **Transazioni Atomiche:** La creazione di ordini (entrata/uscita) e la lavorazione dei resi utilizzano `domain.TransactionHelper`. In caso di errore (ad es. `ErrGiacenzaInsufficiente`), il rollback è automatico e trasparente.
2. **Protezione da Brute Force:** Al raggiungimento di 3 tentativi falliti di login (`domain.MaxTentativiLogin`), l'account viene bloccato (`Bloccato = true`), con l'eccezione dell'account `admin`. L'amministratore può sbloccarlo e ricordare le credenziali tramite `ToggleUserBlock` (`PUT /api/v1/users/{username}`).
3. **Validazione Barcode:** La registrazione dei prodotti convalida la conformità del checksum EAN-13 (`domain.IsValidEAN13`).
