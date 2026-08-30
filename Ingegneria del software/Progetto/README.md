# GiaLoSai

GiaLoSai è un software gestionale client-server sviluppato in Go per la gestione di magazzino, ordini e documentazione amministrativo-contabile. Il sistema adotta un'architettura disaccoppiata: un client a riga di comando basato su un'interfaccia TUI (Terminal User Interface) reattiva comunica tramite API REST JSON con un backend a tre livelli (Controller-Service-Repository) collegato a un database relazionale MySQL 8.

**Obiettivi e funzionalità principali:**
- **Gestione Catalogo e Giacenze:** Monitoraggio in tempo reale di prodotti, categorie merci e quantità disponibili.
- **Flusso Documentale e Contabilità:** Tracciamento degli ordini di entrata/uscita e generazione di DDT, fatture, note di credito e resi.
- **Controllo Accessi (RBAC):** Autenticazione sicura tramite token JWT con permessi differenziati in base al ruolo dell'utente.
- **Interfaccia TUI Reattiva:** Navigazione da terminale sviluppata secondo il pattern Model-View-Update (framework Bubble Tea e Lipgloss).

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


## Quickstart

Requisiti minimi: **Go 1.25**, **Docker** e **Docker Compose**.

1. **Avviare Backend e Database MySQL via Docker Compose:**
   ```bash
   git clone https://github.com/tommasogarofani/ingsw_25-26_group-11.git
   cd ingsw_25-26_group-11
   docker compose -f docker/docker-compose.yml up --build
   ```
   Compose avvia MySQL sulla porta host `3307` e il backend sulla porta `8080`. Entrambe le porte sono esposte su tutte le interfacce di rete (accessibili dall'esterno).

2. **Avviare il Client TUI:**
   ``` bash
   cd client/cmd/gialosai
   go run .
   ```
Il backend risponderà su `http://localhost:8080`.

3. Primo Accesso (Credenziali di Default):
Alla schermata di login della TUI, inserire le credenziali di prova precaricate dallo script prova.sql:
- **Username:** admin
- **Password:** admin

Per arrestare l'ambiente e rimuovere i volumi temporanei:
```bash
docker compose -f docker/docker-compose.yml down -v
```

## Usage
### 1. **Autenticazione e Navigazione Operativa**

All'avvio della TUI, inserire le credenziali di accesso. Una volta autenticato, il sistema mostrerà un menu contestuale basato sul ruolo assegnato (es. Amministratore o Operatore), permettendo di consultare il catalogo, registrare movimenti di magazzino o generare documenti fiscali senza interagire direttamente con il database.

### 2. **Connessione Client da una Macchina Remota (LAN)**

Per collegare un client a un server avviato su un'altra macchina della stessa rete locale:
```bash
cd client/cmd/gialosai
GIALOSAI_SERVER=http://<IP_SERVER>:8080 go run .
```
(Sostituire <IP_SERVER> con l'indirizzo IP locale della macchina che ospita il backend).

## Configurazione

| Parametro / File | Valore Default / Posizione | Descrizione |
|---|---|---|
| GIALOSAI_SERVER | http://localhost:8080 | Variabile d'ambiente letta dal client per impostare l'URL del backend REST. |
| Porta Backend | 8080 | Porta HTTP del server, configurabile in docker/docker-compose.yml e server/cmd/gialosai/main.go. |
| Porta MySQL | 3307 | Porta host del database MySQL (docker/docker-compose.yml). |
| Schema DDL | docker/mysql/init.sql | Definizione di tabelle, vincoli di integrità e chiavi esterne. |
| Seed Dati | docker/mysql/prova.sql | Dati dimostrativi precaricati (utenti di prova con password bcrypt, prodotti e ordini). |

**Nota di Sicurezza:** Dopo il primo accesso con l'account amministratore predefinito, si raccomanda di accedere alla sezione Gestione Utenti per modificare le credenziali di default.

## Test
Per eseguire l'intera suite di test unitari e di integrazione (inclusi i mock SQL e i test dei componenti TUI):
```bash
# Test del modulo client
cd client && go test ./... -v
# Test del modulo server
cd ../server && go test ./... -v
```

## Contributing

Per garantire coerenza e qualità nel codice, il flusso di lavoro prevede le seguenti regole:

- **Branching:** Non effettuare commit direttamente sul ramo `main`. Creare branch dedicati con nomi descrittivi (`feature/nome-feature` o `fix/nome-bug`).
- **Pull Request:** Ogni modifica deve essere sottoposta tramite Pull Request e approvata tramite revisione prima del merge.
- **Stile e Verifiche:** Eseguire `gofmt` prima di effettuare il commit e assicurarsi che tutti i test passino senza errori.
- **Documentazione:** Modifiche all'architettura o alle API REST devono essere riflesse nei documenti presenti nella cartella `documentazione/`.

### Strumenti per lo Sviluppo Locale (Makefile)

Per velocizzare le operazioni di build, testing e pulizia durante lo sviluppo senza dover eseguire manualmente i singoli comandi nei rispettivi moduli, è disponibile un **`Makefile`** nella radice del progetto:

```bash
make test     # Esegue i test di unità per client e server
make compile  # Compila il backend e genera l'eseguibile in bin/gialosai
make clean    # Pulisce la cartella bin/ e la cache di go clean
```

**Autori — Progetto Ingegneria del Software (Gruppo 11, A.A. 2025/2026):**

- Tommaso Garofani — @tommasogarofani
- David Cardinali — @Dyved33
