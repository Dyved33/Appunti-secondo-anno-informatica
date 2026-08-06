# Architectural Design Document - GiaLoSai

## 1. Repository Structure
La disposizione dei file segue una struttura modulare per separare gli artefatti di design, il codice sorgente, la base dati locale e la documentazione del progetto.

```text
ingsw_25-26_group-11/
├── .github/                # Impostazioni di GitHub e workflow CI/CD
├── cmd/                    # Entrypoint dell'applicazione
│   └── gialosai/
│       └── main.go         # Inizializzatore del database e avvio del ciclo Bubble Tea
├── documentazione/         # File e materiali relativi alla documentazione (LaTeX e UML)
│   ├── build/              # Output della compilazione LaTeX (escluso dal versionamento)
│   ├── img/                # Immagini usate nella documentazione e grafici UML esportati
│   ├── uml/                # Modelli e diagrammi PlantUML (*.puml)
│   ├── Decision-log.tex    # Registro delle decisioni architetturali (es. TUI vs Web-App)
│   ├── Design-Document.tex # Documento di alto livello per il design dell'applicazione
│   ├── architettura.md     # Questo documento
│   └── documentazione-progetto-ingegneria-del-software.tex # Tesi di progetto principale
├── src/                    # Codice sorgente dell'applicazione (TUI & Backend)
│   ├── backend/            # Strato logico e di persistenza (Backend)
│   │   ├── domain/         # Modello del dominio (entità, costanti ed errori comuni)
│   │   ├── repository/     # DAO per l'interazione con il database SQLite locale
│   │   └── service/        # Logica di business e casi d'uso (es. AuthService, MagazzinoService)
│   │   // Nota: i pacchetti DTO, handler e model (in backend) sono legacy/in via di dismissione
│   ├── ui/                 # Interfaccia utente TUI basata su Bubble Tea (Frontend)
│   │   └── root.go         # Gestione degli stati dell'applicazione e della navigazione
│   └── database.db         # Database locale SQLite (persistenza locale offline)
├── testdata/               # Dati e suite di test automatizzati
├── .gitignore              # Esclusione di file temporanei e build binary
├── Makefile                # Automazione delle routine di compilazione e test
├── README.md               # Istruzioni di avvio e linee guida di sviluppo
├── go.mod                  # Modulo Go principale (gialosai)
└── go.sum                  # Somme di controllo per le dipendenze di terze parti
```

---

## 2. Dettaglio dei Livelli

In linea con la decisione **DC-001** (adozione di una TUI locale invece di una web-app client-server), l'applicazione è strutturata come un monolite locale organizzato in livelli (Layered Architecture):

### A. Presentation Layer (Frontend — `src/ui`)
- Sviluppato in Go utilizzando il framework **Bubble Tea** (architettura ELM-like con `Init`, `Update`, `View`) e **Lipgloss** per la formattazione grafica del terminale.
- Gestisce l'interazione diretta con l'utente (input da tastiera e lettore di codici a barre in emulazione tastiera).
- I modelli della UI (`LoginModel`, `MagazzinoModel`) non contengono logica di business: interrogano direttamente i servizi locali dello strato sottostante per convalidare le operazioni.

### B. Service Layer (Business Logic — `src/backend/service`)
- Rappresenta il cuore della logica applicativa.
- Riceve le chiamate direttamente dal Presentation Layer (in-process, senza necessità di chiamate di rete o serializzazione JSON/HTTP).
- Implementa le regole di business (es. verifica dei tentativi di login, blocco account, controllo delle soglie minime di giacenza prima di confermare uno scarico, gestione delle transazioni con eventuale rollback).

### C. Data Access Layer (Persistence — `src/backend/repository`)
- Implementa il pattern **DAO (Data Access Object)**.
- Incapsula tutte le query SQL eseguite sulla base dati **SQLite** embedded (`database.db`), isolando il resto del sistema dai dettagli fisici del database.
- Restituisce allo strato dei servizi oggetti di dominio puri definiti nel pacchetto `domain`.

### D. Domain Layer (Shared Types — `src/backend/domain`)
- Contiene le definizioni di tutte le entità (`Utente`, `Prodotto`, `Giacenza`, `DDT`, ecc.), le costanti globali del sistema (soglie, ruoli RBAC) e gli errori sentinella.
- Essendo privo di dipendenze verso altri strati del progetto, viene importato liberamente sia dalla UI che da servizi e repository per garantire la coerenza dei tipi di dato scambiati.
