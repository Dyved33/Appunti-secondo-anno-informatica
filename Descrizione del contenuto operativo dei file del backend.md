### Contenuto del folder domain: 
Contiene i file di configurazione con: 
1. Costanti relative agli attori, allo stato dell'ordine, al tipo di fattura, ad alcune costanti numeriche (es. Numero massimo di tentativi per il login), ai prefissi necessari per ordini, fatture,…, ai path. Vi sono, inoltre, le funzioni per:
	- assicurarsi che un ruolo sia valido (cioè effettivamente presente tra quelli disponibili) 
	- Produrre le labels relative a ciascun ruolo
2. Errori relativi a comportamenti anomali (errori contestuali). Al momento si tratta solo di helper e di funzioni che restituiscono errore con messaggio annesso. Vi è anche la gestione dell'errore relativo ad una giacenza insufficiente
3. Struttura dei modelli che definiscono le proprietà e le caratteristiche degli attori e degli oggetti materiali e immateriali trattati (es.fattura, reso,…). Per ciascuno vengono richiamate le colonne del database che riguardano e caratterizzano l'oggetto 
4. Struttura del contenuto di un inventario. Si tratta di un paio di strutture che vanno a descrivere o l'inventario di uno specifico prodotto o l'inventario di tutto il magazzino (mi servono per restituire tali informazioni al "frontend" tramite json)
5. Verifica che una determinata stringa sia un barcode (così vengono rifiutati i codici solo apparentemente validi)

Da qui in poi un elemento fondamentale sarà il context. Si tratta di una libreria che ho importato per utilizzarla per mandare (ed eliminare) le richieste verso il database che (in ogni caso) si troverà nel server
### Contenuto del folder repository:
L'elemento centrale è il file db.go. Al suo interno vi sono: un interfaccia con tutte le funzioni necessarie a gestire il context, due funzioni che permettono alle funzioni che accettano un QueryExecutor di funzionare sia con transazioni che senza, senza dover cambiare il codice che le utilizza. Per transazioni si intende un blocco di operazioni che devono essere eseguite tutte insieme o nessuna, garantendo l'integrità dei dati. Ad esempio, se si sta creando un ordine e si devono inserire sia l'ordine che le righe dell'ordine, si può utilizzare una transazione per assicurarsi che entrambi gli inserimenti avvengano correttamente o che nessuno dei due venga inserito in caso di errore. Infine 

### Contenuto del folder handler:


### Contenuto del folder service:

