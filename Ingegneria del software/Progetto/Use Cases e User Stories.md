# Use Cases e User Stories — Gestionale Magazzino

---

## UC1 — Effettuare Login

**Attore primario:** Utente (Magazziniere, Contabile, Amministratore)  
**Precondizioni:** L'applicazione è avviata. L'utente non è autenticato.  
**Postcondizioni:** L'utente accede all'applicazione con il proprio ruolo.

**Flusso principale:**
1. L'utente inserisce username e password.
2. Il sistema verifica le credenziali.
3. Il sistema apre la home page in base al ruolo dell'utente.

**Flusso alternativo — Credenziali errate:**
- 2a. Le credenziali non sono valide: il sistema mostra un messaggio di errore e incrementa il contatore di tentativi falliti.
- 2b. Al terzo tentativo fallito il sistema blocca l'account e mostra un avviso.

**Flusso alternativo — Recupero password:**
- L'utente seleziona "Recupera password", inserisce lo username, risponde alla domanda di sicurezza e imposta una nuova password.

---

### User Stories — UC1

#### US1.1 — Accesso con credenziali
> *Come utente, voglio accedere all'applicazione con username e password, in modo da poter utilizzare le funzionalità del mio ruolo.*

**Criteri di accettazione:**
1. Il sistema mostra una schermata di login con i campi username e password.
2. Se le credenziali sono corrette, l'utente accede alla home page corrispondente al suo ruolo.
3. Se le credenziali sono errate, il sistema mostra il messaggio "Credenziali non valide" senza indicare quale dei due campi è sbagliato.

#### US1.2 — Blocco account dopo tentativi falliti
> *Come sistema, voglio bloccare l'account dopo 3 tentativi di accesso falliti consecutivi, in modo da proteggere i dati da accessi non autorizzati.*

**Criteri di accettazione:**
1. Il sistema conta i tentativi falliti consecutivi per ciascun utente.
2. Al raggiungimento del terzo tentativo fallito, il sistema blocca l'account e mostra il messaggio "Account bloccato".
3. Un account bloccato non può essere sbloccato tramite la schermata di login ordinaria.

#### US1.3 — Recupero password
> *Come utente, voglio recuperare la password dimenticata rispondendo alla mia domanda di sicurezza, in modo da riottenere l'accesso al sistema autonomamente.*

**Criteri di accettazione:**
1. Dalla schermata di login è presente un link "Recupera password".
2. Il sistema richiede lo username e mostra la domanda di sicurezza associata.
3. Se la risposta è corretta, l'utente può impostare una nuova password.
4. Se la risposta è errata, il sistema mostra un messaggio di errore senza sbloccare l'accesso.
5. Dopo il cambio password, il contatore di tentativi falliti viene azzerato.

---

## UC2 — Visualizzare Ordini

**Attore primario:** Magazziniere  
**Precondizioni:** L'utente è autenticato come Magazziniere.  
**Postcondizioni:** Il Magazziniere visualizza l'elenco o il dettaglio degli ordini.

**Flusso principale:**
1. Il Magazziniere apre la sezione Ordini.
2. Il sistema mostra l'elenco degli ordini di entrata e uscita.
3. Il Magazziniere seleziona un ordine.
4. Il sistema mostra il dettaglio dell'ordine con le righe prodotto.

**Flusso alternativo — Nessun ordine presente:**
- 2a. Se non ci sono ordini, il sistema mostra il messaggio "Nessun ordine presente".

---

### User Stories — UC2

#### US2.1 — Visualizzazione elenco ordini
> *Come magazziniere, voglio vedere l'elenco di tutti gli ordini di entrata e uscita, in modo da avere una panoramica completa dei movimenti del magazzino.*

**Criteri di accettazione:**
1. L'elenco mostra per ogni ordine: numero, data, tipo (entrata/uscita) e stato.
2. Gli ordini sono visualizzati in ordine cronologico decrescente per impostazione predefinita.
3. Se non ci sono ordini, il sistema mostra un messaggio informativo.

#### US2.2 — Visualizzazione dettaglio ordine
> *Visualizzazione elenco ordini*

**Criteri di accettazione:**
1. Il dettaglio mostra: numero ordine, data, stato e l'elenco delle righe con prodotto e quantità.
2. Il Magazziniere può tornare all'elenco senza perdere la posizione di scorrimento.

---

## UC3 — Gestione Giacenze

**Attore primario:** Magazziniere  
**Precondizioni:** L'utente è autenticato come Magazziniere. I prodotti esistono a catalogo.  
**Postcondizioni:** Le giacenze sono aggiornate, o è stato generato un report inventario.

**Flusso principale — Aggiungere merci:**
1. Il Magazziniere crea un nuovo ordine di entrata.
2. Aggiunge le righe prodotto con le relative quantità.
3. Conferma l'ordine.
4. Il sistema aggiorna le giacenze per ogni riga e segna l'ordine come ricevuto.

**Flusso principale — Visualizzare / Stampare:**
1. Il Magazziniere apre la sezione Giacenze.
2. Il sistema mostra l'elenco dei prodotti con le quantità disponibili.
3. Il Magazziniere può richiedere la stampa di un report inventario in PDF.

**Flusso alternativo — Nessun prodotto in catalogo:**
- 2a. Se non esiste alcun prodotto, non è possibile aggiungere righe all'ordine.

---

### User Stories — UC3

#### US3.1 — Registrazione ordine di entrata
> *Come magazziniere, voglio registrare un ordine di entrata con i prodotti ricevuti, in modo da aggiornare automaticamente le giacenze.*

**Criteri di accettazione:**
1. Il sistema permette di aggiungere una o più righe, ciascuna con prodotto e quantità.
2. Alla conferma dell'ordine, le giacenze di ogni prodotto coinvolto vengono incrementate della quantità indicata.
3. L'ordine confermato compare nell'elenco degli ordini con stato RICEVUTO.
4. Non è possibile confermare un ordine privo di righe.

#### US3.2 — Visualizzazione giacenze
> *Come magazziniere, voglio visualizzare le giacenze attuali di tutti i prodotti, in modo da conoscere la disponibilità reale in magazzino.*

**Criteri di accettazione:**
1. L'elenco mostra per ogni prodotto: codice, nome e quantità disponibile.
2. I prodotti con giacenza zero sono visibili e distinguibili dagli altri.

#### US3.3 — Stampa inventario
> *Come magazziniere, voglio stampare un documento inventario, in modo da avere un registro fisico dello stato del magazzino.*

**Criteri di accettazione:**
1. Il sistema genera un documento PDF con l'elenco completo dei prodotti e le relative giacenze.
2. Il documento riporta la data di generazione.
3. Il PDF può essere aperto per la stampa direttamente dall'applicazione.

---

## UC4 — Generare DDT

**Attore primario:** Magazziniere  
**Precondizioni:** L'utente è autenticato come Magazziniere. Esiste almeno un ordine di uscita con stato NON_EVASO.  
**Postcondizioni:** Il DDT è generato e l'ordine di uscita passa allo stato EVASO.

**Flusso principale:**
1. Il Magazziniere apre l'elenco degli ordini di uscita non evasi.
2. Seleziona un ordine e ne visualizza il dettaglio.
3. Richiede la generazione del DDT.
4. Il sistema genera il documento DDT e imposta l'ordine come EVASO.

**Flusso alternativo — Nessun ordine non evaso:**
- 1a. Se non esistono ordini non evasi, il sistema mostra il messaggio "Nessun ordine da evadere".

---

### User Stories — UC4

#### US4.1 — Generazione DDT
> *Come magazziniere, voglio generare un DDT per un ordine di uscita, in modo da documentare la spedizione della merce.*

**Criteri di accettazione:**
1. Il DDT è generabile solo per ordini con stato NON_EVASO.
2. Il documento DDT riporta: numero DDT, data, elenco prodotti e quantità dell'ordine associato.
3. Il DDT può essere aperto per la stampa direttamente dall'applicazione.

#### US4.2 — Cambio stato automatico a EVASO
> *Come magazziniere, voglio che alla generazione del DDT l'ordine venga automaticamente segnato come evaso, in modo da non dover aggiornare lo stato manualmente.*

**Criteri di accettazione:**
1. Immediatamente dopo la generazione del DDT, lo stato dell'ordine diventa EVASO.
2. L'ordine evaso non compare più nell'elenco degli ordini non evasi.
3. Il DDT generato rimane consultabile dal dettaglio dell'ordine.

---

## UC5 — Gestione Resi

**Attori primari:** Contabile (registra il reso), Magazziniere (verifica il reso)  
**Precondizioni:** Esiste almeno un ordine di uscita con stato EVASO.  
**Postcondizioni:** Il reso è verificato e le giacenze sono aggiornate.

**Flusso principale:**
1. Il Contabile registra un nuovo reso collegandolo a un ordine di uscita evaso.
2. Il sistema crea il reso con stato DA_VERIFICARE.
3. Il Magazziniere apre l'elenco dei resi e seleziona quello da verificare.
4. Il Magazziniere visualizza il dettaglio e conferma la verifica fisica della merce.
5. Il sistema aggiorna le giacenze e imposta il reso come VERIFICATO.

**Flusso alternativo — Nessun ordine evaso disponibile:**
- 1a. Se non esistono ordini evasi, il Contabile non può aprire un reso.

---

### User Stories — UC5

#### US5.1 — Registrazione reso da parte del Contabile
> *Come contabile, voglio registrare un reso collegandolo a un ordine di uscita evaso, in modo da tracciare la restituzione della merce da parte del cliente.*

**Criteri di accettazione:**
1. Il sistema permette di selezionare solo ordini con stato EVASO come origine del reso.
2. Il reso viene creato con stato DA_VERIFICARE e compare nell'elenco dei resi.
3. Non è possibile creare un reso senza collegarlo a un ordine evaso.

#### US5.2 — Visualizzazione elenco e dettaglio resi
> *Come magazziniere, voglio vedere l'elenco dei resi e il dettaglio di ciascuno, in modo da sapere quali restituzioni devo verificare fisicamente.*

**Criteri di accettazione:**
1. L'elenco mostra per ogni reso: numero, data e stato.
2. Il dettaglio mostra i prodotti e le quantità restituite.
3. I resi con stato DA_VERIFICARE sono visivamente distinguibili dagli altri.

#### US5.3 — Verifica reso e aggiornamento giacenze
> *Come magazziniere, voglio confermare la verifica fisica di un reso, in modo da reintegrare le quantità restituite nelle giacenze del magazzino.*

**Criteri di accettazione:**
1. Solo i resi con stato DA_VERIFICARE possono essere verificati.
2. Dopo la conferma, le giacenze di ogni prodotto del reso vengono incrementate della quantità restituita.
3. Lo stato del reso diventa VERIFICATO.

---

## UC6 — Contabile: Ordini e Fattura

**Attore primario:** Contabile  
**Precondizioni:** L'utente è autenticato come Contabile.  
**Postcondizioni:** Il Contabile ha visualizzato gli ordini e/o generato una fattura.

**Flusso principale:**
1. Il Contabile apre la sezione Ordini di uscita.
2. Il sistema mostra l'elenco completo degli ordini di uscita.
3. Il Contabile può filtrare gli ordini per stato (Evaso / Non Evaso).
4. Il Contabile seleziona un ordine e ne visualizza il dettaglio.
5. Il Contabile richiede la generazione di una fattura per un ordine evaso.
6. Il sistema genera il documento fattura.

**Flusso alternativo — Ordine non evaso:**
- 5a. Se l'ordine selezionato non è evaso, il sistema blocca la generazione e mostra un messaggio di errore.

---

### User Stories — UC6

#### US6.1 — Visualizzazione ordini di uscita
> *Come contabile, voglio visualizzare l'elenco degli ordini di uscita con il loro stato, in modo da monitorare i movimenti in uscita dal magazzino.*

**Criteri di accettazione:**
1. L'elenco mostra per ogni ordine: numero, data e stato (Evaso / Non Evaso).
2. Il Contabile può aprire il dettaglio di qualsiasi ordine.

#### US6.2 — Filtro ordini per stato
> *Come contabile, voglio filtrare gli ordini di uscita per stato, in modo da separare facilmente quelli evasi da quelli ancora da evadere.*

**Criteri di accettazione:**
1. Sono disponibili almeno due opzioni di filtro: "Evaso" e "Non Evaso".
2. L'applicazione del filtro aggiorna immediatamente l'elenco visualizzato.
3. È possibile rimuovere il filtro per tornare alla vista completa.

#### US6.3 — Generazione fattura
> *Come contabile, voglio generare una fattura a partire da un ordine di uscita evaso, in modo da produrre la documentazione contabile della vendita.*

**Criteri di accettazione:**
1. La fattura è generabile solo per ordini con stato EVASO.
2. Il documento fattura riporta: numero fattura, data, elenco prodotti con quantità e prezzi unitari, totale.
3. Se l'ordine non è evaso, il sistema mostra il messaggio "Impossibile generare fattura: ordine non ancora evaso".
4. La fattura generata è consultabile in seguito dalla sezione Fatture.

---

## UC7 — Contabile: Prezzi e Valore Magazzino

**Attore primario:** Contabile  
**Precondizioni:** L'utente è autenticato come Contabile. I prodotti esistono a catalogo.  
**Postcondizioni:** Il prezzo di un prodotto è aggiornato, o il Contabile ha visualizzato il valore totale del magazzino.

**Flusso principale — Aggiornare prezzo:**
1. Il Contabile apre la lista prodotti.
2. Seleziona un prodotto e inserisce il nuovo prezzo unitario.
3. Il sistema salva il nuovo prezzo.

**Flusso principale — Valore magazzino:**
1. Il Contabile apre la sezione Valore Magazzino.
2. Il sistema calcola e mostra la quantità totale dei prodotti e il loro valore complessivo.

---

### User Stories — UC7

#### US7.1 — Aggiornamento prezzo unitario prodotto
> *Come contabile, voglio aggiornare il prezzo unitario di un prodotto, in modo da mantenere i valori del magazzino sempre allineati ai prezzi correnti.*

**Criteri di accettazione:**
1. Il Contabile può modificare il prezzo unitario di qualsiasi prodotto presente a catalogo.
2. Il prezzo può essere inserito al pezzo o al chilo, in base all'unità di misura del prodotto.
3. Il sistema accetta solo valori numerici positivi come prezzo.
4. Dopo il salvataggio, il nuovo prezzo è immediatamente riflesso nel calcolo del valore magazzino.

#### US7.2 — Visualizzazione valore totale magazzino
> *Come contabile, voglio vedere il valore complessivo del magazzino, in modo da avere una stima economica immediata delle scorte disponibili.*

**Criteri di accettazione:**
1. Il sistema mostra la quantità totale dei prodotti in magazzino.
2. Il sistema mostra il valore totale calcolato come somma di (quantità × prezzo unitario) per ogni prodotto.
3. Il riepilogo si aggiorna automaticamente ad ogni modifica di giacenza o di prezzo.

---

## UC8 — Amministratore: Gestione Completa

**Attore primario:** Amministratore  
**Precondizioni:** L'utente è autenticato come Amministratore.  
**Postcondizioni:** L'Amministratore ha eseguito le operazioni richieste su merci, fatture o resi.

**Flusso principale:**
1. L'Amministratore sceglie la sezione da gestire (Merci, Fatture o Resi).
2. Il sistema carica e mostra i dati con il valore associato dove applicabile.
3. L'Amministratore seleziona un elemento ed esegue l'azione disponibile.
4. Il sistema aggiorna i dati e mostra una conferma.

**Flusso alternativo — Azione non consentita:**
- 3a. Se l'elemento non può essere modificato o eliminato (es. prodotto con giacenze attive), il sistema mostra un messaggio di avviso.

---

### User Stories — UC8

#### US8.1 — Gestione merci
> *Come amministratore, voglio poter modificare ed eliminare i prodotti a catalogo, in modo da mantenere il catalogo accurato e aggiornato.*

**Criteri di accettazione:**
1. L'elenco merci mostra per ogni prodotto: codice, nome, quantità disponibile e valore totale (quantità × prezzo).
2. L'Amministratore può modificare tutti i campi di un prodotto.
3. L'Amministratore può eliminare un prodotto; il sistema richiede conferma prima di procedere.

#### US8.2 — Gestione fatture
> *Come amministratore, voglio poter annullare o eliminare le fatture, in modo da correggere eventuali errori nella documentazione contabile.*

**Criteri di accettazione:**
1. L'elenco mostra tutte le fatture con numero, data e stato.
2. L'Amministratore può annullare una fattura impostandone lo stato ad ANNULLATA senza eliminarla dal sistema.
3. L'Amministratore può eliminare definitivamente una fattura; il sistema richiede conferma prima di procedere.
4. Una fattura annullata rimane visibile nell'elenco ma è visivamente distinguibile dalle attive.

#### US8.3 — Gestione resi
> *Come amministratore, voglio poter approvare o rifiutare i resi verificati, in modo da avere il controllo finale sul processo di reso.*

**Criteri di accettazione:**
1. L'elenco mostra tutti i resi con numero, data e stato corrente.
2. L'Amministratore può approvare un reso: il sistema imposta lo stato ad APPROVATO.
3. L'Amministratore può rifiutare un reso: il sistema imposta lo stato a RIFIUTATO.
4. Solo i resi con stato VERIFICATO possono essere approvati o rifiutati.
