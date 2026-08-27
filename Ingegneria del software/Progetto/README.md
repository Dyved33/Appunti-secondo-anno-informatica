# GiaLoSai

GiaLoSai è un progetto universitario di Ingegneria del Software per la gestione di magazzino e informazioni amministrative collegate. L'implementazione è client-server: un client testuale comunica via HTTP/JSON con un backend REST che gestisce regole applicative e persistenza MySQL.

Questa pagina è il riferimento principale per architettura, struttura e avvio del software. La relazione tecnica, i diagrammi UML e i documenti di progetto restano in [documentazione/](documentazione/).

## Architettura

```text
Operatore
   |
   v
Client TUI (Bubble Tea, Lipgloss)
   | HTTP/JSON + Authorization: Bearer <JWT>
   v
Backend REST Go (chi)
   |-- middleware di autenticazione
   |-- controller -> service -> repository
   v
MySQL 8
```

Il client raccoglie l'input e mantiene lo stato delle viste; non accede direttamente al database. Il pacchetto `client/api` traduce le operazioni della TUI in chiamate REST e conserva in memoria il token di sessione. Il server è l'unico componente che legge o modifica MySQL.

Il backend è organizzato per area di dominio: `auth`, `user`, `product`, `category` e `order`. Ogni area separa controller, service e repository.

| Layer | Responsabilità |
| --- | --- |
| Router e middleware | Registra le rotte `/api/v1`, separa login pubblico e risorse protette, verifica l'header `Authorization`. |
| Controller | Decodifica e codifica DTO JSON, delegando la logica ai service. |
| Service | Applica le regole di business e coordina operazioni atomiche. |
| Repository | Incapsula query SQL parametrizzate con `database/sql` e il driver MySQL. |
| Domain | Definisce entità, costanti, errori, validazione EAN-13 e infrastruttura transazionale. |

`router.NewRouter` compone repository, service e controller usando la medesima connessione `*sql.DB`. I service dipendono da interfacce, quindi possono essere sottoposti a test unitari con mock.

### Autenticazione, autorizzazione e transazioni

`POST /api/v1/login` verifica le credenziali e restituisce un JWT con identificativo, username e ruolo. Le password sono gestite con bcrypt. Il server salva l'hash SHA-256 del token e la sua scadenza nella tabella `Sessioni`.

Per ogni rotta protetta il middleware controlla firma e scadenza del JWT, poi la presenza della relativa sessione non scaduta nel database; infine inserisce l'identificativo utente nel contesto della richiesta. Il servizio utenti blocca l'account dopo tre tentativi falliti; l'account `admin` non viene mai bloccato. Il ruolo è presente nel token, ma il router corrente non applica ancora autorizzazioni diverse per ruolo sulle singole rotte.

`domain.TransactionHelper` inserisce una transazione nel `context.Context`; i repository usano la transazione quando presente e la connessione altrimenti. Un errore provoca il rollback dell'intera operazione. Il meccanismo è usato per prodotti, ordini e resi. Un ordine di uscita verifica tutte le quantità disponibili prima di ridurre le giacenze; un ordine di entrata registra le righe e incrementa le quantità nello stesso confine transazionale.

### Modello applicativo

Lo schema è in [docker/mysql/init.sql](docker/mysql/init.sql). Le entità principali sono:

- `Utenti` e `Sessioni`: identità, ruoli e validità delle sessioni;
- `Categorie`, `Prodotti` e `Giacenze`: catalogo e quantità disponibili;
- `OrdiniEntrata` e `OrdiniUscita`, con le relative righe;
- `DDT`, `Fatture` e `Resi`, collegati agli ordini di uscita.

Lo schema usa InnoDB e `utf8mb4`. Le chiavi esterne mantengono le relazioni: ogni prodotto ha una sola giacenza, le righe sono eliminate insieme al rispettivo ordine e un DDT è associato a un solo ordine di uscita.

## Flussi di dominio

Dopo il login, il client invia il JWT come Bearer token a ogni richiesta. Il catalogo combina prodotti e categorie; l'aggiornamento del prezzo di vendita modifica il prodotto in transazione.

Nei service, un ordine di entrata aumenta le giacenze; uno di uscita ne controlla la disponibilità, le riduce e registra l'ordine. La generazione DDT assegna una numerazione annuale e porta l'ordine nello stato `evaso`. I resi passano da `da_verificare` a `verificato`, quindi ad `approvato` o `rifiutato`; la verifica ripristina la giacenza e il rifiuto la riduce nuovamente. Per ordini evasi, i service calcolano imponibile, IVA al 22% e totale di fatture o note di credito.

I percorsi PDF sono registrati nei modelli, ma il codice non contiene un generatore di PDF o CSV.

## Struttura del repository

```text
.
├── client/                       # Modulo Go del client TUI
│   ├── api/                      # Client HTTP e modelli lato client
│   ├── ui/app/                   # Stato e ciclo di vita Bubble Tea
│   ├── ui/components/, styles/   # Componenti e stile Lipgloss
│   └── ui/pages/                 # Login, menu, magazzino e contabilità
├── server/                       # Modulo Go del backend
│   ├── auth/, user/              # Login, JWT, sessioni e utenti
│   ├── product/, category/       # Catalogo, giacenze e report inventario
│   ├── order/                    # Ordini, DDT, fatture, note di credito e resi
│   ├── domain/, middleware/      # Modello comune e protezione HTTP
│   ├── router/                   # Composizione dei componenti e rotte REST
│   ├── docker/mysql/             # Schema DDL e dati seed
│   └── cmd/gialosai/main.go      # Server HTTP sulla porta 8080
├── docker/                       # Dockerfile e Docker Compose
├── documentazione/               # Documenti LaTeX e diagrammi UML
├── .github/workflows/            # Notifiche degli eventi Git su Discord
└── Makefile                      # Compilazione, test e pulizia
```

`client` e `server` hanno `go.mod` distinti: i comandi Go vanno eseguiti nella rispettiva directory, salvo i target del `Makefile`.

## Tecnologie

- Go 1.25;
- Bubble Tea, Bubbles e Lipgloss per la TUI Model-View-Update;
- chi per il router HTTP;
- MySQL 8 e `go-sql-driver/mysql`;
- `golang-jwt/jwt/v5` e bcrypt;
- Docker e Docker Compose;
- `testify` e `go-sqlmock` nei test.

Le versioni esatte sono nei file `client/go.mod` e `server/go.mod`.

## API REST esposta

Tutte le rotte, tranne il login, richiedono autenticazione.

| Metodo e percorso | Operazione |
| --- | --- |
| `POST /api/v1/login` | Autentica e restituisce token e ruolo. |
| `GET, POST /api/v1/products/` | Elenca o crea prodotti. |
| `PUT /api/v1/products/{id}/price` | Aggiorna il prezzo di vendita. |
| `GET /api/v1/products/inventory-report` | Restituisce il report inventario. |
| `GET /api/v1/categories/` | Elenca categorie. |
| `GET /api/v1/orders/`, `/entrata` | Restituisce ordini di entrata. |
| `GET /api/v1/orders/resi` | Restituisce resi. |
| `GET /api/v1/orders/invoices`, `/credit-notes` | Restituisce documenti fiscali. |
| `PUT /api/v1/users/` | Crea un utente. |
| `GET, PUT, DELETE /api/v1/users/{username}` | Gestisce un utente. |

La tabella descrive il router, non tutti i metodi presenti nei controller e nei service: creazione ordini di uscita, DDT, fatture, ricerca barcode e recupero password sono implementati internamente ma non hanno una rotta nel router corrente.

## Avvio

### Prerequisiti

Sono necessari Go 1.25, Docker con Docker Compose e Git; `make` è opzionale.

### Server (macchina che ospita il backend)

```bash
git clone https://github.com/tommasogarofani/ingsw_25-26_group-11.git
cd ingsw_25-26_group-11
docker compose -f docker/docker-compose.yml up --build
```

Compose avvia MySQL sulla porta host `3307` e il backend sulla porta `8080`. Entrambe le porte sono esposte su tutte le interfacce di rete (accessibili dall'esterno).

Per fermare l'ambiente:

```bash
docker compose -f docker/docker-compose.yml down -v
```

Il volume `mysql_data` conserva i dati. Gli script in `docker/mysql/` vengono eseguiti solo alla prima creazione del volume, in ordine alfabetico: `init.sql` (DDL) e `prova.sql` (dati dimostrativi: utenti di prova con password bcrypt, prodotti, ordini, fatture). Per ripartire da un database vuoto e rieseguirli: `docker compose -f docker/docker-compose.yml down -v`.

### Client locale (stessa macchina del server)

```bash
cd client
go run .
```

Il client usa `http://localhost:8080` di default.

### Client da un'altra macchina (rete locale)

Sulla macchina server, trovare l'IP locale:

```bash
hostname -I
```

Sulla macchina client:

1. Installare Go 1.25 e Git.
2. Clonare il repository e avviare il client:

```bash
git clone https://github.com/tommasogarofani/ingsw_25-26_group-11.git
cd ingsw_25-26_group-11/client
GIALOSAI_SERVER=http://<IP_SERVER>:8080 go run .
```

Sostituire `<IP_SERVER>` con l'indirizzo IP del server (es. `192.168.1.10`). La variabile `GIALOSAI_SERVER` è opzionale: se non impostata, il client usa `http://localhost:8080`.

### Avvio server sull'host (senza Compose)

Per avviare il server sull'host, mantenendo il database Compose attivo:

```bash
DB_DSN='root:rootpassword@tcp(127.0.0.1:3307)/gialosai?parseTime=true' make run
```

Senza `DB_DSN`, il server tenta di usare l'host `db`, disponibile normalmente solo nella rete Compose.

Lo schema di riferimento (`docker/mysql/init.sql`) crea soltanto tabelle, senza utenti o dati dimostrativi: non esistono quindi credenziali predefinite avviando il server sull'host. Poiché la creazione utenti è una rotta protetta, il primo utente deve essere predisposto nel database con password bcrypt, oppure deve essere aggiunto un meccanismo di bootstrap al codice. Nell'ambiente Compose sono invece disponibili utenti di prova tramite `docker/mysql/prova.sql`.

## Test

```bash
make test
```

Il target esegue `go test -v ./...` nei due moduli. La suite include test di service, repository SQL mockati, controller, middleware, router, client API e modelli TUI. In ambienti che impediscono l'apertura di socket locali, i test basati su `httptest.NewServer` possono fallire per il vincolo dell'ambiente.

## Stato dell'integrazione e limiti residui

La documentazione è stata allineata al codice. Le funzionalità principali sono integrate (v2.0.0). I limiti residui sono:

- il ruolo è presente nel token JWT, ma il router non applica ancora autorizzazioni diverse per ruolo sulle singole rotte;
- i percorsi dei documenti PDF e il volume `output` esistono, ma la generazione materiale dei file non è implementata;
- `GET /api/v1/orders/` restituisce ordini di entrata, non di uscita;
- il client chiede il report a `/api/v1/product/inventory-report` e i resi a `/api/v1/resi`, percorsi diversi da quelli esposti dal router.

## Contributi e autori

Il flusso previsto richiede branch dedicati, pull request revisionate e aggiornamento di test e documentazione quando necessario; non si lavora direttamente su `main`.

Progetto del Gruppo 11 per il corso di Ingegneria del Software, anno accademico 2025/2026.

- Tommaso Garofani — [@tommasogarofani](https://github.com/tommasogarofani)
- David Cardinali — [@Dyved33](https://github.com/Dyved33)
