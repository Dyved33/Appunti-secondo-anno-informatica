# Architectural Design Document - GiaLoSai

## 1. Repository Structure
La disposizione dei file segue una struttura modulare per separare gli artefatti di design, il codice sorgente e la suite di test. 

```text
ingsw_25-26_group-11/
├── documentazione/         # Files e materiali relativi alla documentazione
│   ├── build/              # Output della compilazione LaTeX (esclusa dal versionamento)  
│   │   ├── *.aux
│   │   ├── *.fdb_latexmk
│   │   ├── *.fls
│   │   ├── *.log
│   │   ├── *.out
│   │   ├── *.synctex.gz
│   │   ├── *.toc
│   │   └── *.pdf  
│   ├── img/                # Immagini usate nella documentazione
│   │   └── *.png       
│   ├── *.md   
│   └── *.tex          
├── src/                    # Codice sorgente dell'applicazione
│   ├── backend/          
│   │   ├── DTO/
│   │   ├── handler/
│   │   ├── model/
│   │   ├── repository/
│   │   └── service/ 
│   ├── controller/
│   ├── frontend/
│   │   ├── app/
│   │   ├── config/
│   │   └── routes/
│   ├── model/
│   ├── view/
│   ├── database.db
│   └── main.go
├── test/                   # Suite di test automatizzati
├── uml/                    # Modelli PlantUML
│   └── *.puml
├── .github/                # Impostazioni GitHub
├── .gitignore              # Esclusione file temporanei e artefatti di build
├── README.md               
└── Makefile                # automazione routine di compilazione e test
```
### Caratteri: 

`├`, `└`, `─`. `│`  

## 2. Backend Structure
La struttura del backend è organizzata in più moduli tra loro comunicanti, ciascuno con responsabilità ben definite secondo un’architettura a livelli.

L’accesso ai dati viene gestito tramite la componente Repository, incaricata dell’interazione con il database e delle operazioni di persistenza. La logica applicativa e le regole di business sono invece demandate alla componente Service, che coordina i flussi operativi dell’applicazione.

La comunicazione tra frontend e backend avviene attraverso appositi Controller, responsabili della gestione delle richieste HTTP e della restituzione delle relative risposte. Per il trasferimento dei dati verranno utilizzati i DTO (Data Transfer Object), serializzati in formato JSON, così da garantire uno scambio di informazioni strutturato e indipendente dai modelli interni dell’applicazione.

A questo livello viene inoltre gestita l’esposizione delle API tramite specifici Handler, mentre la componente Model definisce la struttura delle entità e dei dati utilizzati all’interno del sistema.

Backend/
   ├── handler/
   ├── service/
   ├── repository/
   ├── model/
   └── dto/
   └── controller/
