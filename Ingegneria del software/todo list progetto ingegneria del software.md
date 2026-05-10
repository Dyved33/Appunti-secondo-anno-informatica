# Contesto 
- Scelta obbligatoria: ogni studente deve definire il gruppo (1, 2 o 3 persone) e il contesto del progetto. 
- Repository del progetto: creare un repository privato con prefisso obbligatorio **ingsw_25-26_group-x**, dove <u>x è il numero progressivo assegnato dal docente</u>. 
- Ogni gruppo sceglie un contesto applicativo realistico e circoscritto (e.g.ospedale/clinica, assicurazione, banca, e-commerce, gestione universitaria/segreteria studenti)
- Il giorno dell’esame è richiesta la consegna tramite export completo in formato ZIP (o archivio equivalente); nome consigliato: **ingsw_25-26_group-x.zip** (esempio). 
- Non pubblicare il repository. 
- <u>Regola operativa: scegliere un dominio che consenta di definire chiaramente i casi d’uso richiesti dallo scope</u>.  

# Gruppi e scope proporzionale 
Sono ammessi gruppi da 1, 2 o 3 persone. Lo scope del progetto deve essere proporzionale alla dimensione del gruppo.

| Gruppo    | Ampiezza attesa | Indicazione minima                                                              |
| --------- | --------------- | ------------------------------------------------------------------------------- |
| 3 persone | Estesa          | Scope più articolato, con integrazione tra più componenti e maggiore profondità |
# Requisiti minimi comuni (tutti i gruppi) 
- **Repository**: creare un repository GitHub dedicato al progetto, usarlo per tutta la consegna e mantenerlo ordinato. 
- **Struttura**: organizzare il progetto in cartelle con responsabilità chiare (es. docs/, src/, test/, uml/). 
- **Requisiti**: niente Word; consegnare i requisiti in LATEX (es. file requisiti.tex) con funzionali, non funzionali e acceptance criteria. 
- **Modellazione** UML: usare PlantUML, esportare i diagrammi in PDF e includerli nella relazione in LATEX. 
- **Implementazione**: linguaggio e stack liberi (Java, C++, Python, web stack, script Bash, . . . ), purché coerenti e funzionanti. 
- **Qualità del design**: applicare in modo ragionato i principi SOLID, con attenzione particolare a SRP, OCP e DIP. 
- *Test*: eseguire una suite di test utile e significativa rispetto ai casi d’uso implementati. 
- D**ocumentazione codice**: <u>README obbligatorio; codice commentato; JavaDoc in Java o documentazione equivalente per altri linguaggi</u>

# Vincoli comuni di consegna 
- **Versionamento**: creare un tag finale di release (es. v 1.0.0) e un changelog minimo. 
- **Tracciabilità**: collegare casi d’uso, implementazione e test (anche con una tabella sintetica). 
- **Riproducibilità**: dichiarare prerequisiti ambiente e comando unico di build/test. 
- **Contributi**: indicare in modo esplicito il contributo dei membri del gruppo. 
- **Igiene del repository**: .gitignore corretto, struttura pulita, nessun file inutile o artefatto temporaneo. 
- **Pacchetto finale**: specificare chiaramente cosa costituisce la consegna (repo + eventuale relazione PDF).  

# Dettaglio scope: gruppo da 3 persone 
- Implementare 6 casi d’uso principali end-to-end. 
- Suddividere il lavoro per aree funzionali o componenti, con confini tecnici espliciti. 
- Includere test unitari e 2 test di integrazione. 
- Verificare coerenza tra contributi diversi (stile codice, naming, struttura moduli/cartelle).  

# Cosa verrà valutato 
- **Sul progetto**: 
	- qualità della soluzione proposta,
	- coerenza con i principi del corso, 
	- chiarezza del lavoro svolto e delle scelte. 
- **Nella discussione**: 
	- motivazioni delle decisioni progettuali, 
	- comprensione di test, architettura e workflow, 
	- capacità di argomentare trade-off e limiti.   

# Step necessari
### Fase 1: Comprensione del Problema e Requisiti
Per prima cosa bisogna capire cosa deve fare il gestionale.
- **Il Backlog e le User Stories:** Creare un elenco ordinato delle funzionalità. Scrivere le necessità sotto forma di User Stories ("Come [attore], voglio [azione] per [obiettivo]").
- **Criteri di Accettazione:** Per ogni User Story, mettere per iscritto i criteri di accettazione. Questa è la prima "documentazione": servirà per sapere quando una funzionalità è davvero finita e per scrivere i test in seguito.
- **Mappa Concettuale (Domain Model):** Disegnare in modo informale le entità del vostro dominio (es. Cliente, Ordine, Prodotto) e le loro relazioni. Aiuterà tutto il team a usare lo stesso vocabolario.
### Fase 2: Setup Iniziale e Modellazione UML
Ora che sapete cosa fare, iniziate a impostare il progetto e a formalizzare le interazioni
- **Inizializzazione del README.md:** Creare il file alla radice del repository Git. In questa fase deve solo rispondere alla prima domanda: _Che cos'è il progetto?_ (Overview generale e obiettivi).
- **Use Case Diagram:** Documentate visivamente chi sono gli attori (es. Utente, Amministratore) e i casi d'uso principali.
- **Aggiornamento del Modello di Dominio:** Affinate le entità aggiungendo attributi essenziali e molteplicità (es. 1 a molti). Niente dettagli tecnici o database per ora.
### Fase 3: Scelte Architetturali
È il momento di decidere come organizzare il codice per rispettare principi come la separazione delle responsabilità e la modularità.
- **Il Design Document (1-3 pagine):** Scrivere un documento in `/docs/design/` che spieghi l'architettura scelta (es. stile Layered con CLI/UI, Service, Domain, Repository). Inseriteci un UML minimale (Class Diagram) che mostri i confini tra i layer.
- **Decision Log (ADR - Architecture Decision Record):** Per ogni scelta importante scrivere un breve ADR in `/docs/decisions/`. Segnare: ID, Data, Contesto, Decisione, Alternative scartate e Conseguenze.
### Fase 4: Sviluppo Quotidiano (Codice e Test)
Questa è la fase iterativa. La documentazione qui viene scritta a fianco del codice.
- **Javadoc mirata:** Non commentate ogni singola riga di codice. Aggiornate la Javadoc **solo per le API pubbliche e i Service**. Documentare il _contratto_ (cosa fa il metodo, precondizioni, eccezioni lanciate), non l'implementazione interna. 
- **I Test come Documentazione:** Gli Unit Test e Integration Test sono documentazione eseguibile. Bisogna dare ai metodi di test nomi lunghi e descrittivi (es. `iscrizioneFallisceSeUtenteGiaIscritto()`) in modo che chiunque capisca la regola di business validata.  
- **Messaggi di Commit (Git):** Trattare i commit come documentazione storica. Bisogna evitare messaggi come "fix" o "aggiornamento". Invece va scritto cosa è cambiato e _perché_ è stato fatto.
- **Aggiornamento del README:** Man mano che il progetto prende forma, aggiungete le sezioni _Quickstart_ (comandi reali per far partire il gestionale), _Usage_ e _Test_ (come lanciare la suite JUnit).
### Fase 5: Modifiche, Bugfix e Refactoring
Il gestionale è in corso d'opera, i requisiti cambiano o emergono dei code smells (es. Large Class o Long Method).
- **Definition of Done (DoD) della Documentazione:** Non fare mai una Pull Request o un merge su `main` se la documentazione è disallineata.
- **Regola pratica degli aggiornamenti:**
    - Cambia il modo di avviare l'app? $\rightarrow$ Aggiornate il **README**.
    - Cambiano i confini tra i livelli dell'architettura? $\rightarrow$ Aggiornate il **Design Document**.    
    - Prendete una nuova decisione tecnica non banale? $\rightarrow$ Aggiungete un **ADR**.     
    - Cambiano i parametri o i ritorni di un metodo pubblico? $\rightarrow$ Aggiornate la **Javadoc**.
### Fase 6: Rilasci (Manutenzione ed Evoluzione)
Quando il gestionale è pronto per la consegna o per presentare una prima versione funzionante.
- **Tag di Git e Semantic Versioning:** Usate i tag di Git (`MAJOR.MINOR.PATCH`) per segnare le versioni stabili (es. `v1.0.0`).  
- **Changelog / Release Notes:** Quando viene rilasciata una versione, bisogna creare un riepilogo testuale delle issue chiuse, dei bug risolti e delle nuove feature aggiunte. Questo dice a chi legge cosa aspettarsi da quell'aggiornamento.