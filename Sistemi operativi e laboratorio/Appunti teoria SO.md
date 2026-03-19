###### Evoluzione Hardware e Multiprogrammazione
> [!NOTE] Contesto Storico
> Agli albori dell'informatica, quando i calcolatori funzionavano ancora a valvole ed erano estremamente costosi e lenti, la presenza di un operatore umano era indispensabile per risparmiare tempo e gestire manualmente l'esecuzione.

**Nuove Funzioni Hardware**
Con l'evoluzione dei sistemi di elaborazione, sono state introdotte nuove funzionalità a livello hardware per consentire al **Monitor** (il predecessore dell'odierno sistema operativo) di gestire meglio le risorse:
* **Protezione della memoria:** Funzionalità che impedisce ai programmi utente di accedere ad aree di memoria non autorizzate.
* **Timer:** Meccanismo che previene la monopolizzazione della CPU da parte di un singolo processo.
* **Istruzioni privilegiate:** Particolari istruzioni di controllo che possono essere eseguite solamente dal Monitor.
* **Interruzioni:** Eventi che permettono al sistema di sospendere l'esecuzione corrente per gestire operazioni asincrone.

> [!IMPORTANT] Approfondimento: Violazione di Memoria e Interruzioni
> In certe situazioni è strettamente necessario poter "violare" o accedere alla memoria in modo controllato durante l'esecuzione. 
> Le **interruzioni** permettono al sistema di compiere un salto verso un'istruzione specifica, che non fa parte delle normali istruzioni utente ma del linguaggio del Monitor. Ad esempio, se un processo richiede troppo tempo per l'esecuzione (*too much time to execute*), l'interruzione restituisce il controllo al sistema.

**Il Problema dell'Overhead**
L'utilizzo del Monitor ha però generato il problema dell'**Overhead** (sovraccarico di sistema):
* Una parte significativa della memoria centrale deve essere costantemente occupata dal Monitor stesso.
* Una frazione del tempo macchina (CPU) viene consumata dal Monitor per il suo stesso funzionamento, sottraendo tempo all'esecuzione effettiva dei programmi.

**Terza Generazione (1965–1980)**
La terza generazione di calcolatori è segnata dall'introduzione dei **circuiti integrati**.

> [!NOTE] Nota sui Circuiti Integrati
> I circuiti integrati sono costituiti da piastre di silicio che contengono migliaia di transistor. Questo ha permesso di costruire computer molto più potenti e decisamente meno dispendiosi in termini di energia.

**La Multiprogrammazione**
L'innovazione software chiave di questo periodo è la **multiprogrammazione**. È nata per risolvere l'enorme spreco di tempo della CPU durante le operazioni di Input/Output.

Ecco un esempio di quanto fosse inefficace l'esecuzione sequenziale:
* Lettura di un record: 15 micros 
* Esecuzione di 100 istruzioni: 1 micros 
* Scrittura di un record: 15 micros 
* **Totale:** 31 micros 
* **Utilizzo reale della CPU:** 1/31 = 3.2% 

 **Uniprogramming vs Multiprogramming**
* **Uniprogramming:** I processi vengono eseguiti uno alla volta, in modo strettamente sequenziale.
* **Multiprogramming:** Si caricano più job contemporaneamente in memoria e si assegna a turno la CPU ai vari processi.

> [!EXAMPLE] Dinamica di Assegnazione delle Risorse
> Immaginiamo 3 job che necessitano di risorse diverse. Nella multiprogrammazione, il sistema li carica tutti in memoria contemporaneamente, purché lo spazio sia sufficiente e periferiche come disco e stampante non debbano essere utilizzate nello stesso esatto momento. Se il primo job si blocca in attesa di un I/O e non usa la CPU, questa risorsa non rimane inattiva, ma viene assegnata al processo successivo pronto per l'esecuzione.

**I Problemi della Multiprogrammazione**
Questo approccio ha introdotto nuove sfide architetturali e di gestione:
1. **Gestione della memoria:** Nell'esempio precedente la memoria richiesta era sufficiente, ma nella realtà, se si attivano numerosi processi, questi potrebbero non entrare fisicamente nella memoria disponibile, richiedendo tecniche di gestione avanzate.
2. **Scheduling:** Determinare come, quando e a chi assegnare il processore. Quando la CPU passa dal primo processo agli altri a ruota, il sistema operativo deve "togliere" il processo corrente dal processore. Questo genera il problema di dover salvare i dati di stato del processo uscente, per poi poterli ripristinare in modo trasparente quando il processo otterrà nuovamente il turno di esecuzione.

**Il Time-Sharing**
Il **Time-sharing** (condivisione del tempo) è un paradigma sviluppato per supportare la gestione di più lavori interattivi:
* Permette a calcolatori potenti di supportare più utenti contemporaneamente sulla stessa macchina.
* Il processore viene condiviso in modo fluido tra tutti gli utenti attivi.
* L'accesso avviene simultaneamente attraverso l'uso di terminali.
* Il sistema operativo alterna raPIDamente l'esecuzione di ciascun programma utente assegnando un **quanto di elaborazione** (un breve periodo di tempo).
* Esempio storico notevole: il sistema **CTSS** (Compatible Time-Sharing System).

**Multiprogrammazione vs Time-Sharing**
Pur condividendo principi di esecuzione concorrente, hanno finalità diverse:

| Caratteristica | Multiprogrammazione | Time-Sharing |
| :--- | :--- | :--- |
| **Obiettivo Principale** | Massimizzare l'uso del processore (massimizzare i lavori completati) | Minimizzare il tempo di risposta percepito dall'utente |
| **Istruzioni al SO** | **JCL** (Job Control Language): istruzioni preparate in anticipo e fornite insieme al job | Comandi interattivi inseriti direttamente dall'utente al terminale |
**Quarta Generazione (1980 - 1990)**
* Nasce con l'introduzione della tecnologia **LSI** (Integrazione su larga scala).
* Vede la nascita e la diffusione di massa dei **Personal Computer**.
* È caratterizzata dallo sviluppo di un'**interfaccia amichevole** (user-friendly) per avvicinare l'informatica al grande pubblico.

 **Quinta Generazione (1990 - oggi)**
* Dominata dall'avvento di dispositivi mobili come gli **smartphone**.
* Impone nuovi vincoli fisici e progettuali:
  * Gestione dell'**energia limitata**.
  * Integrazione e gestione di **nuovi dispositivi**.

> [!NOTE] Nuovi Dispositivi
> I moderni sistemi operativi mobili devono tenere costantemente conto dell'energia utilizzata, dei forti limiti della **memoria** a disposizione e della presenza di innumerevoli componenti hardware aggiuntive, come i **sensori**, che richiedono driver e gestioni specifiche.
###### Funzionamento di un sistema di elaborazione
- *CPU:* costituita dai registri (es. Program counter che contiene l'indirizzo dell'istruzione successiva da eseguire o l'instruction register che contiene l'istruzione da eseguire,...)
- *Memoria principale:* vengono caricati i programmi che sono costituiti da istruzioni (quelli che la CPU può eseguire) e dati (locazioni di memoria in cui mettiamo le nostre informazioni). 
- *Moduli I/O:* relativi a tastiera, monitor e dischi

Nei sistemi moderni ci sono più processori (cores) su un singolo chip (socket) . Ognuno di questi può svolgere contemporaneamente agli altri delle istruzioni/programmi. Meno si accede alla memoria principale e più è veloce la computazione => per questo esiste la memoria cache (unità di memoria piccola così da essere veloce e poco costosa) dalla quale la CPU prende un piccolo numero di dati e informazioni. I dati contenuti nella cache sono quelli usati di recente (molto probabilmente li riuseremo). Poi di solito c'è un processore grafico (GPU). Questo è in grado di applicare la stessa istruzione a tutti i dati di un array (in un tempo ridotto). Usati principalmente per aspetti grafici (es. Pixel). La tecnica usata è la *Single-instruction multiple-data* (SIMD) e si può usare anche per applicazioni che non sono grafiche. Esiste poi un processore di segnale digitale (DSP) che codifica/decodifica segnali audio e video e garantisce un supporto crittografico. Ad oggi tutto ciò si può trovare su un solo chip (Soc) => funzionale per i cellulari

**Ciclo di estrazione ed esecuzione:**
L'istruzione viene estratta, viene eseguita e, se il processo è terminato, si torna in dietro e si ricomincia dall'inizio. Tipi di istruzione riguardano:
- Memoria
- I/O
- Elaborazione dati (es. Operazioni logiche o aritmetiche)
- Controllo (es. Goto)

> [!tip]
> RICORDA:
> - PC = memorizza l'indirizzo di memoria della prossima istruzione da prelevare. Durante la fase di fetch il valore contenuto nel PC viene inviato alla memoria e, dopo il prelievo, viene automaticamente incrementato per puntare all'istruzione successiva nella sequenza. Nel caso di un salto il PC viene caricato con un nuovo indirizzo specifico
> - IR = contiene l'istruzione che la CPU sta elaborando. L'istruzione prelevata dalla RAM viene copiata nell'IR. Dopodiché il suo contenuto viene letto per decodificare l'operazione e i dati da utilizzare

**Interruzioni:** Se la macchina funzionasse sempre in questo modo non avrei modo di intervenire. Per questo esistono le interruzioni che permettono agli altri moduli di interrompere il ciclo di fetch-execute. In questo modo si evita di lasciare la CPU inattiva o in attesa di I/O e si può sospendere l'esecuzione di un programma. Le classi di interruzione sono:
- Interruzioni di programma = generate da errori nell'esecuzione di un programma (es. Overflow, eccezioni di indirizzo = un programma non può andare a modificare le informazioni di altri programmi e deve restare nello spazio di memoria a lui assegnato, ecc ...)
- Timer = permettono al sistema operativo di eseguire delle funzioni a intervalli regolari
- I/O = generate da un controllore I/O per segnalare il completamento di un'operazione o un errore
- Errori HW = generati da problemi hardware

<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Screenshot From 2026-03-04 17-58-44.png" style="width: 100%; border-radius: 8px;">
  </div>
  <div style="flex: 1.5;">
  Quindi il procedimento è: estrazione istruzione -> esecuzione istruzione -> check delle interruzioni. Se c'è l'interruzione se ne guarda il tipo si compie un salto all'indirizzo al quale inizia il programma per la gestione delle interruzioni (<i>interrupt handler</i>) e successivamente si salta ad un'altra istruzione. Una volta che il gestore delle interruzioni ha finito di gestire il problema legato ad un'interruzione, questa viene eseguita. 
  </div>
</div>

Se avviene un'interruzione durante la gestione di un'interruzione ci sono almeno 2 possibilità:
1. Disabilito le interruzioni durante la gestione di un'interruzione (può funzionare dato che le interruzioni sono brevi)
2. Definisco uno schema di priorità: in questo caso il sistema permette a un'interruzione di interromperne un'altra in base ad un sistema di priorità. Ciò significa che se l'interruzione in arrivo ha una priorità **più alta** di quella attualmente in esecuzione, il processore mette in pausa la prima, viene salvato il PC e viene eseguita la seconda (quella più urgente). Alla fine della seconda istruzione, il processore recupera i dati relativi al PC e riprende la prima delle due istruzioni che ha causato un'interruzione.
###### Sviluppo di un sistema operativo
**Aspetti principali nello sviluppo dei sistemi operativi:**
- Processi
- Gestione della memoria
- Protezione e sicurezza dell'informazione
- Schedulazione e gestione delle risorse
- Struttura del sistema

**Approcci e elementi progettuali dei sistemi operativi:**
- Architettura micro-kernel
- Multithread
- Multiprocessing simmetrico
- Sistemi operativi distribuiti
- Architettura orientata agli oggetti

Il sistema operativo può realizzare una macchina virtuale. Cioè il sistema operativo crea una macchina che abbia una sua memoria e un suo sistema operativo. Esiste un virtual machine monitor (VMM)  che viene eseguito sul sistema operativo e gestisce le comunicazioni tra i singoli sistemi operativi virtuali e il processore

**Servizi del sistema operativo:**
- Esecuzione programmi
- Gestione di operazioni di I/O
- Filesystem
- Comunicazione con altri sistemi
- Allocazione delle risorse
- Accounting (verifica chi usa il sistema e per quanto tempo)
- Sicurezza e protezione
- Riconoscimento degli errori (nell'esecuzione dei programmi)

Queste funzioni sono utilizzate attraverso le chiamate di sistema = interruzioni richieste dall'utente stesso. Si accede al sistema tramite delle interfacce:
- GUI
- Batch (file di istruzioni che vengono eseguite una dietro l'altra)
- Linea di comando

**Servizi per l'utente:**
- Interfaccia utente: devo essere in grado di interagire con il sistema operativo (o da interfaccia grafica o da linea di comando o dai file di batch)
- Esecuzione di programmi: caricamento, esecuzione e terminazione di programmi
- Operazioni I/O
- Organizzare il filesystem: creare/eliminare file, leggere/scrivere in un file, cercare informazioni e gestire i permessi (di lettura, scrittura o esecuzione per utente, gruppi e altri)
- Comunicazione fra processi: ogni processo ha il suo spazio e non può invadere lo spazio di altri processi (di solito va bene ma a volte voglio che i processi interagiscano tra di loro). A volte ci può essere uno spazio di memoria condivisa tra i vari processi o in alternativa c'è uno scambio di messaggi (un processo manda un messaggio ad un altro processo attraverso il sistema operativo)
- Rilevazione errori: hardware, di I/O o dei programmi utente

**Funzioni che servono ad assicurare l'efficienza del sistema:**
- Allocazione delle risorse: alcuni programmi necessitano di risorse di processore/i, memoria, disco o altri dispositivi I/O. L'ordine con cui assegno le risorse determina se il sistema è efficiente o meno
- Accounting: soprattutto nei processi condivisi è bene sapere chi usa il sistema e quanto lo usa 
- Protezione e sicurezza: utilizzo di tutte le risorse va riservato a chi ne ha diritto (protezione) e controllo degli accessi al sistema (sicurezza)

**Interprete dei comandi (CLI) vs interfaccia grafica (GUI):**

| CLI                                                                                                                                                                                                                              | GUI                                                                                                                                                                                                                                       |
| -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Esegue in sequenza i comandi specificati dall'utente. Viene implementato nel nucleo o programma di sistema (in questo caso ci sono molte versioni (*shells*). I comandi sono implementati nell'interprete o in programmi esterni | È un'interfaccia amichevole (*user friendly*) in cui le icone rappresentano file, directories, programmi e funzioni. Inoltre ci sono menù che permettono di selezionare i comandi e la possibilità di avere touch-screen e comandi vocali |
| Per compiti ripetititvi                                                                                                                                                                                                          | Più intuitiva ma non accede a tutte le funzioni                                                                                                                                                                                           |
###### System call
Le system call sono un'interfaccia per il programmatore ai servizi del SO. Vi si accede attraverso l'API che nasconde i dettagli al programmatore (gestita da libreria di supporto run-time). I parametri sono passati attraverso:
- Registri
- In un blocco o tabella di memoria dell'indirizzo in un registro
- Nello stack da cui sono prelevati (pop) dal sistema operativo

**Categorie di chiamate di sistema:**
1. **Controllo dei processi:
 - Creazione e arresto
 - Caricamento e esecuzione
 - Terminazione normale e anormale
 - Esame e impostazione degli attributi tramite funzioni apposite
 - Funzioni che mettono in pausa per un tempo indicato
 - Attesa segnalazione evento 
 - Assegnazione e rilascio di memoria
2. **Gestione file:**
 - Creazione e cancellazione
 - Apertura e chiusura
 - Lettura, scrittura, posizionamento (posso leggere un carattere a scelta e non per forza uno dopo l'altro)
 - Esame e impostazione degli attributi
3. **Gestione dei dispositivi I/O:**
 - Richiesta (di assegnazione) e rilascio
 - Lettura, scrittura e posizionamento
 - Esame e impostazione degli attributi
 - Inserimento ed esclusione logica
4. **Gestione delle informazioni:**
 - Esame e impostazione dell'ora e della data
 - Esame e impostazione dei dati di sistema
 - Esame e impostazione degli attributi dei processi, file e dispositivi
5. **Comunicazione:**
 - Creazione e chiusura di connessione
 - Invio e ricezione di messaggi
 - Creazione e accesso a regioni di memoria condivisa (interazione con condivisione)
 - Informazioni sullo stato di un trasferimento
 - Inserimento ed esclusione dei dispositivi remoti
6. **Protezione:**
 - Controllo degli accessi alle risorse
 - Visualizzazione e impostazione dei permessi
 - Autorizzazione e divieto di accesso all'utente

Un sistema operativo è fatto anche di programmi detti *utilities* che servono a realizzare l'ambiente utile per lo sviluppo e l'esecuzione dei programmi. Ce ne sono diverse categorie:
- Gestione files
- Informazione di stato
- Modifica di file (text editor)
- Supporto alla programmazione (compilatore, debugger, linker,...)
- Caricamento ed esecuzione dei programmi
- Comunicazioni
- Servizi in background
- Programmi applicativi
###### Progettazione di un sistema operativo
Non è semplice dato che devo tenere conto di che hardware si dispone e di che tipo di richiesta abbiamo da parte dell'utilizzatore. Quindi le specifiche da rispettare sono:
- Decidere hardware
- Decidere che sistema usare
- Obiettivi dell'utilizzatore
- Obiettivi del sistema (es obiettivo dell'utilizzatore = far terminare i programmi nel minor tempo possibile, obiettivo del sistema = non lasciare il processore inutilizzato)
- Separazione delle politiche (se ho un server servo il primo client e non l'ultimo che ha fatto una richiesta) dai meccanismi (es come si accede al disco)

I moderni sistemi operativi sono scritti in linguaggi come C e C++ (non più Assembler) perché ciò mi garantisce:
- Facilità di sviluppo e manutenzione
- Compattezza
- Portabilità

Inoltre sono linguaggi più efficienti di linguaggi interpretati come Python (computer non sono fatti per mandare avanti il sistema operativo, ma per far partire programmi di sistema). Alcune sezioni critiche che richiedono altissima efficienza sono scritte in Assembler.
###### Struttura del sistema operativo
**Sistemi monolitici (Non sono divisi in moduli):** MS-DOS o UNIX tradizionale sono sistemi monolitici complessi e scritti tramite molte funzioni. Il kernel di UNIX tradizionale ha un'interfaccia verso l'hardware per la gestione delle chiamate di sistema (questo all'interno del kernel che gestisce segnali I/O, file system, driver, dispositivi di archiviazione e scheduling)

**Sistemi stratificati:** Ogni layer (strato) è l'implementazione di un oggetto astratto fatto di dati e operazioni per manipolarli. Ogni strato nasconde i dettagli dell'implementazione agli strati superiori.
- *Vantaggi* = semplicità di scrittura, aggiornamento, debugging e portabilità
- *Svantaggi* = difficoltà nel definire gli strati e overhead

**Microkernel:** Sono sistemi monolitici che facciano il minimo necessario. Il resto è trasferito a livello-utente (server). Le richieste di servizi avvengano tramite scambio di messaggi tra client e server.
- *Vantaggi* = semplicità di aggiornamento e debugging, portabilità e affidabilità
- *Svantaggi* = overhead (il kernel affida più compiti ai programmi e per questo c'è bisogno di un maggiore tempo di esecuzione)

> [!example]
> ![[Screenshot From 2026-03-05 12-31-03.png|200]]

Non c'è differenza dal punto di vista del sistema tra questi programmi, ma ogni volta che un programma applicativo deve accedere a un file o un dispositivo invierà una richiesta tramite messaggio che sarà svolta un un altro programma che è il gestore del dispositivo e non dal kernel, dato che questo contiene solo le funzioni fondamentali (altre funzionalità sono in moduli, ogni modulo nasconde i dettagli di implementazione agli altri moduli)

<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Screenshot From 2026-03-05 12-34-21.png" style="width: 100%; border-radius: 8px;">
  </div>
  <div style="flex: 1.5;">
  Solaris ha questo tipo di architettura: il kernel contiene tutte le funzioni fondamentali e una serie di moduli da collegare quando si vuole per le loro specifiche funzioni
  </div>
</div>

**Debugging:** Per il fatto che i SO sono sistemi complessi, ci possono essere degli errori e per questo nella loro progettazione e gestione c'è il problema del debugging. Gli errori, però, si possono presentare in situazioni particolari (difficili da individuare), considerando che il SO deve fare il suo compito nel minor tempo possibile. Quindi valutare le prestazioni di un SO (e per monitorare il comportamento del sistema) richiede una lunga sperimentazione tramite anche la generazione di file di log che registrano il comportamento del sistema e *core dump*

**Generazione sistema operativo:** Scritto il codice sorgente del sistema $\rightarrow$ va compilato e configurato per il sistema su sui verrà eseguito $\rightarrow$ poi si istalla il sistema $\rightarrow$ avvia il computer e il nuovo sistema operativo. La configurazione può realizzarsi modificando il codice prima della compilazione in base al sistema su cui lo vado a installare, oppure tramite la selezione di tabelle e moduli precompilati da collegare al momento dell'installazione, oppure tramite un sistema guidato da tabelle di dati che guidano il sistema in base all'hardware.<u> Per far partire il sistema l'esecuzione inizia in un predefinito indirizzo della ROM (parte di memoria non modificabile nemmeno a computer spento = <i>read only memory</i>) che contiene il <i>bootstrap loader</i>. Questo contiene le istruzioni necessarie a caricare il <i>bootstrap program</i> che a sua volta carica il kernel e avvia l'esecuzione del sistema</u>
###### Processi 
**Definizione:**<u> un processo (o job) è un programma (lista di istruzioni con riferimenti a delle locazioni di memoria) in esecuzione</u>. Gli elementi del processo sono:
- Program counter (indirizzo istruzione da eseguire) e contenuto dei registri (immagine del processore)
- Stack che contiene variabili locali, indirizzi di return o parametri delle funzioni
- Sezione testo (programma)
- Sezione dati che contiene i valori delle variabili locali
- Heap per l'allocazione dinamica

| Programma                                                                                                         | Processo                                                       |
| ----------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------- |
| Entità passiva su disco / file eseguibile                                                                         | Entità che esegue le istruzioni e usa le risorse associate     |
| Sono istruzioni che non cambiano finché non avviene una modifica e non ha memoria allocata né stato di esecuzione | Possiede risorse proprie ed evolve attraverso stati successivi |
> [!tip] Ricorda
>  - Un programma può essere associato a più processi attraverso il multitasking = situazione in cui il SO permette di avviare lo stesso codice eseguibile più volte, assegnando a ciascuna istanza uno spazio in memoria, un PID e delle risorse
>  - nella memoria il processo occupa un certo spazio occupato dal codice, i dati, l'heap e lo stack dove vengono eseguite le procedure (la heap si allarga verso l'alto e lo stack verso il basso e bisogna fare attenzione affinché non collidano)

**Stati del processo:**
- *Nuovo* = in fase di creazione
- *In esecuzione* = il processore sta eseguendo le istruzioni
- *In attesa* = il processo attende un evento
- *Pronto* = il processo attende la disponibilità di un processore
- *Terminato* = il processo ha terminato l'esecuzione

<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Screenshot From 2026-03-05 13-17-14.png" style="width: 100%; border-radius: 8px;">
  </div>
  <div style="flex: 1.5;">
  Quando il processo viene creato può trovarsi in stato di attesa di un processore assegnato ad esso che esegua le istruzioni e lo termini. Può capitare che l'esecuzione del processo sia temporaneamente interrotta => non è necessario che le risorse del processore siano ancora assegnate a questo processo => si salva l'immagine del processore e il processo rimane in attesa dell'evento che lo libera dallo stato di attesa. Quando l'evento avviene saranno riassegnate al processo delle risorse con l'immagine salvata prima e si completerà l'esecuzione/terminerà il processo
  </div>
</div>

**PCB (descrittore di processo):** A ogni processo è associata una struttura dati detta descrittore di processo o PCB (*process control block*) contenente tutte le informazioni relative al processo:
- Stato del processo
- Program counter
- Registri CPU (immagine processore)
- Parametri per lo scheduling della CPU 
- Informazioni per la gestione della memoria
- Informazioni di accounting per sistemi multiutente
- Informazioni su I/O

Un processo può avere multiple sequenze di istruzioni (*threads*) che possono essere eseguite in parallelo (anche su core differenti) per la realizzazione di task differenti. In tal caso il PCB include informazioni su ciascun thread.

**Salvataggio del PCB:**
<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Screenshot From 2026-03-05 13-28-52.png" style="width: 100%; border-radius: 8px;">
  </div>
  <div style="flex: 1.5;">
    Supponiamo che venga creato il processo P0 che è in stato di running. Ad un certo punto si verifica un'interruzione e tutte le informazioni del processo vengono salvate nel PCB0. Queste saranno ricaricate dal sistema di gestione delle interruzioni dopo che il processore avrà terminato un altro processo P1 (ricaricato da PCB1 quando magari P1 era stato interrotto precedentemente). A questo punto può accadere una nuova interruzione, si salva lo stato di P1 dentro a PCB1 e si ricarica lo stato contenuto in PCB0 (si riprende esecuzione di P0). Questo ci permette di usare il sistema in maniera efficiente purché le operazioni di salvataggio e ricarica degli stati del processore siano eseguite in maniera efficiente.
  </div>
</div>

Il sistema ha bisogno di conoscere la situazione dei vari processi e per questo i PCB sono inseriti in varie code:
- Coda di job che contiene tutti i processi 
- Coda di ready che contiene tutti i processi pronti a partire (quelli tra i quali il processore deve decidere a chi assegnare le risorse)
- Coda dei dispositivi (una per dispositivo) che contiene i processi in attesa di interagire con quel dispositivo (il dispositivo genera una nuova interruzione)
###### Schedulatori e fasi del processo
Il sistema operativo deve decidere a quali processi assegnare il processore e per fare ciò ci sono diversi schedulatori:
- **Schedulatore a breve termine (CPU scheduler):** Prende i processi dalla coda ready e gli assegna un processore. Siccome deve essere eseguito con grande frequenza, deve essere molto veloce (i processi running vengono anche messi in coda ready proprio perché lo scheduler assegna la CPU a un processo, poi la toglie a quello per assegnarla ad un altro, così da far avanzare più processi contemporaneamente).
- **Schedulatore a lungo termine:** Seleziona tra i processi nuovi quali fa passare in coda ready, ma non è opportuno che i processi finiscano tutti in questa coda. Per questo esiste il *job scheduler*. Seleziona i processi da inserire in coda ready e per fare ciò controlla il grado di multiprogrammazione, cioè i programmi eseguiti contemporaneamente. Non viene eseguito con frequenza, per questo può applicare procedure più complesse. Ci sono ad esempio processi CPU-bound, ovvero processi che non rilasciano mai la CPU a loro assegnata a meno che non siano forzati 

> [!example]
> - I/O bound processi che interagiscono con input e output
> - In UNIX se entro un certo periodo di tempo il processo non rilascia il processore viene costretto a rilasciarlo e viene rimesso in coda (timebound)
> È bene avere entrambi i tipi di processi in quanto ci potrebbero essere o processi che non possono avanzare o processi che richiedono troppe risorse se ci fosse una prevalenza dei due

- **Schedulatore a medio termine:** Determina con che velocità avanzano i processi nel sistema. Si occupa di assegnare o revocare la memoria ai processi (*swapping*). Quando i processi non hanno bisogno di memoria, gli viene revocata la RAM e la memoria del processo viene copiata su disco rigido. La memoria liberata viene occupata da un nuovo processo. 

La schedulazione realizza delle politiche per l'assegnazione della CPU. Generalmente distingue tra:
- Processi *I/O bound*
- *Processi CPU bound*

Oltre a questo c'è un problema di overhead (risorse accessorie sono ridondanti rispetto a quelle necessarie) nel cambio dei contesto, il quale varia in funzione di:
- Velocità della memoria
- Numero registri del processore
- Istruzioni speciali (riservate al kernel)

**Creazione del processo**
Il sistema operativo fornisce meccanismi per creazione/terminazione dei processi. Generalmente un processo genitore chiede al kernel di creare un processo figlio che a sua volta può creare altri processi (*albero dei processi*). Ogni processo è identificato da un numero detto *identificatore di processo (PID)* $\rightarrow$ su linux è il primo numero che vedo quando scrivo `ps` nel terminale. Quando un processo termina, viene eliminato dall'albero.

**Assegnazione delle risorse ad un processo figlio**
Avviene in base al tipo di sistema operativo:
- Un processo figlio riceve le risorse dal processo padre 
 - Risorse ripartite = il padre non può usarle finché il figlio non ha finito di usarle
 - Risorse condivise 
- Il processo figlio ottiene le risorse dal sistema operativo

Genitore e figlio possono essere eseguiti concorrentemente mentre il genitore attende la terminazione del figlio. Per quanto riguarda la memoria il processo figlio può avere un duplicato di memoria del genitore e quindi avere lo stesso programma, oppure può averne uno nuovo.

> [!info] La funzione `fork` è una syscall che crea una copia esatta del processo e restituisce al padre il PID del figlio:
> - **0 (Zero):** Indica che il processo corrente è il **processo figlio** appena creato.
> - **Maggiore di 0 (Positivo):** Indica che il processo corrente è il **processo padre** e ritorna il valore corrisponde al PID del figlio.
> - **-1 (Negativo):** Indica che la `fork()` è fallita e nessun processo figlio è stato creato.

**Terminazione del processo**
Il processo termina con la syscall *exit()*, così:
- Restituisce un valore di stato al processo genitore, via *wait()*
- Le risorse sono deallocate dal sistema operativo

Il genitore può terminare l'esecuzione del figlio con la syscall *abort()* e questo per varie motivazioni:
- Eccesso di risorse allocate
- Compito non più richiesto

La terminazione di un processo figlio fa si che vengano terminati anche i figli del processo terminato (terminazione a cascata). Solo dopo il wait() del genitore viene deallocato il PID: 
- Il PID però non scompare subito e il processo rimane in stato di zombie e scompare solo dopo che il genitore ha eseguito il wait()
- Il genitore termina senza che vengano terminati i suoi processi figli, in tal caso quei processi si dicono orfani.
###### **Threads**
I processi sono caratterizzati da:
- **Risorse allocate:** spazio di indirizzamento virtuale (immagine in memoria), la traccia di esecuzione, RAM, dispositivi I/O, file
- **Scheduling:** traccia di esecuzione (alternata con altri processi), stati e priorità

Ma le due caratteristiche potrebbero essere trattate in maniera indipendente:
- **Processo (o task):** possiede le risorse
- **Thread (o processo leggero):** unità di lavoro all’interno di un processo. 

Si parla di multithreading la capacità di un SO di permettere più tracce di esecuzione concorrenti in un singolo processo (quindi ad un processo o più processi vengono assegnati più di un singolo thread ciascuno)

| Single-threading                                         | Multi-threading                                                |
| -------------------------------------------------------- | -------------------------------------------------------------- |
| Un processo, un thread (MD-DOS)                          | Un processo, più thread (Java run-time)                        |
| Più processi, un thread per ciascuno (UNIX tradizionale) | Più processi, ciascuno con più thread (Windows, Solaris, UNIX) |
**Definizione:** in un ambiente multi-thread, il processo è l'unità di allocazione delle risorse e un'unità di protezione (nessun altro può accedere alle risorse assegnate al processo). Associati al processo ci sono:
- Spazio di indirizzamento virtuale (contiene immagine del processo = è l'insieme completo e strutturato dei dati che definiscono lo stato di un programma in esecuzione, includendo il codice, i dati in RAM, lo stack e i registri)
- Accesso protetto al processore, ad altri processi (IPC), files, dispositivi I/O
- Uno o più thread

Associati al thread ci sono:
- Stato = pronto, in attesa, in esecuzione
- Contesto del thread (immagine del processo), viene salvata quando non è in esecuzione
- Stack di esecuzione
- Spazio di memoria statico per le variabili locali
- Accesso a memoria e risorse del processo, condivise con gli altri thread del processo

Questo genera un problema: due thread che accedono agli stessi dati/stessa memoria (magari uno vuole incrementare e l'altro decrementare la stessa variabile)

<div style="display: flex; justify-content: center;">
  <img src="Screenshot From 2026-03-12 19-20-44.png" width="200">
</div>

I benefici del multithreading sono:
- La creazione di un thread è più rapida della creazione di un processo
- La terminazione di un thread è più rapida della terminazione di un processo (no questione di protezione della memoria)
- Lo switch tra due thread del medesimo processo è più rapida dello switch tra due processi (sempre per via della memoria)
- Efficienza della comunicazione fra programmi (eseguiti da thread dello stesso processo)

In un sistema a utente singolo:
- Esecuzione in foreground e background
- Elaborazione asincrona
- Velocità di esecuzione
- Struttura modulare dei programmi

**Osservazione:** Si nota che la schedulazione è gestita principalmente nel thread, ma non la sospensione e la terminazione del processo.

Gli stati del thread sono gli stessi dei processi (in esecuzione, pronto e in attesa), così come le operazioni (creazione, block/sospensione, unblock/riattivazione, terminazione).
<div style="display: flex; justify-content: center;">
  <img src="Screenshot From 2026-03-12 19-24-20.png" width="200">
</div>

**Sincronizzazione dei thread:** I thread di un processo condividono spazio di indirizzamento e altre risorse. Ogni modifica simultanea di un dato da parte di due thread può corrompere il dato stesso. Per questi motivi è necessario sincronizzare l'accesso dei thread alle risorse

**Thread a livello utente:** La gestione dei thread è fatta dall'applicazione e il kernel non sa dell'esistenza dei thread.
- *Vantaggi*: lo scambio di thread non richiede privilegi della modalità kernel, la schedulazione può essere specifica per l'applicazione e portabilità.
- *Svantaggi:* le chiamate di sistema sono bloccanti per l'intero processo e non ci sono vantaggi dai multiprocessori
- *Contromisure:* programmare le applicazioni in processi multipli e *jacketing* (trasforma una chiamata di sistema bloccante in una non bloccante)

**Thread a livello kernel:** La gestione dei thread è fatta dal kernel e l'applicazione usa un'API per la gestione dei thread. Il kernel gestisce i PCB dei processi e dei thread e inoltre la schedulazione del processore sulla base dei thread.
- *Vantaggi:* possibilità di eseguire thread differenti su diversi processori, il blocco di un thread non blocca il processo e anche le routine del kernel possono essere gestite in multithreading
- *Svantaggi:* i cambi di modalità (utente/kernel) producono overhead

**Approcci misti:** Creazione dei thread, schedulazione e sincronizzazione nello spazio utente
- Thread modalità utente mappati su thread modalità kernel (in numero minore o uguale)
- Permette esecuzione parallela dei thread e il blocco di un thread non blocca il processo (come solaris)

| Modello | Descrizione                                                                                                                                              |
| ------- | -------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1:1     | Ogni thread di esecuzione è un processo unico con il suo spazio e le sue risorse                                                                         |
| M: 1     | Un processo definisce uno spazio di indirizzo e il possesso di risorse dinamiche. Molti threads possono essere creati ed eseguiti in quel processo       |
| 1: M     | Un thread può migrare dall'ambiente di un processo all'ambiente di un altro processo. Ciò permette al thread di muoversi facilmente tra sistemi distinti |
| M: N     | Combina gli attributi di 1: M e M: 1                                                                                                                       |

###### Legge di Amdahl
$\text{speedup} = \frac{\text{tempo di esecuzione su un processore}}{\text{tempo di esecuzione su N processori}} = \frac{1}{(1-f)+ \frac{f}{n}}$
La legge di Amdahl stabilisce che il miglioramento massimo delle prestazioni di un sistema, grazie alla parallelizzazione (uso di più processori), è limitato dalla parte del software che non può essere parallelizzata. In breve: l'algoritmo, non il numero di processori, determina il limite massimo di velocità.

Il vantaggio lo hanno:
- Applicazioni multithread native (pochi processi e molti thread)
- Applicazioni multiprocesso (molti processi single-thread)
- Applicazioni Java (JVM è un processo multithread e assicura schedulazione e gestione della memoria per le applicazioni Java)
- Multiple istanze di un'applicazione

**Threads in UNIX:**
- Il kernel non distingue tra processi e thread
- Per ogni task (processo o thread) c'è un *task control block* (TCB) con fork()
- Fork() crea una nuova task, con una copia del task control block
- Clone() crea una nuova task, ma nel suo TCB ci sono puntatori ai campi di quello del genitore (in base ai parametri passati a clone()
###### Mutua esclusione
I processi possono interferire tra di loro in maniera positiva o negativa, negativa ad esempio quando sottrae risorse ad altri processi, in maniera positiva quando due o più processi interagiscono e ci permettono di realizzare un’attività che sfrutta l’uso di diversi processi in un’unica applicazione. La gestione dei processi può avvenire tramite: 
 - **Multiprogrammazione**: molti processi con un singolo processore, i processi avanzano attendendo che il processore si liberi; 
 - **Multiprocessing**: molti processi gestiti da un multiprocessore; 
 - **Processi distribuiti**: molti processi, gestiti da molti processori distribuiti (**cluster**), con memorie distribuite, rende le cose più complicate da gestire, ma sicuramente più veloci. 
 
 I differenti processi hanno bisogno di risorse che tipicamente vengono condivise tra i vari processi, si crea quindi una *concorrenza*:
 - Nelle applicazioni multiple, o nelle applicazioni strutturate (con più processi o threads che devono accedere a delle risorse)
 - La modifica dei dati disponibili da parte di un processo influenza l’andamento di un altro processo
 - Nella struttura del sistema operativo.

**Principi della concorrenza**
- Nell'uniprocessore multiprogrammato, le esecuzioni dei processi si alternano, interrompendo certi processi e riprendendone altri
- Nel multiprocessore, le esecuzioni dei processi si alternano e si sovrappongono

Questi processi detti *concorrenti* potrebbero dover accedere ad una risorsa e non si può prevedere quale processo specifico entrerà per primo nella risorsa, poiché la velocità relativa di esecuzione dei processi non è prevedibile. Dipende da:
 - Attività degli altri processi
 - Gestione delle interruzioni
 - Politiche di schedulazione

Vengono generate le seguenti difficoltà:
- Condivisione di risorse globali
- Allocazione ottimale delle risorse
- Debugging (può comparire un errore specifico dovuto a una serie di eventi sfortunati che si verifica raramente e sarà difficile da ri-individuare)

> [!info] L’uniprocessore prende il processo P1, esegue chin=getchar(), sospende l’esecuzione di P1 passando a P2, rifà getchar() e cambia chin, stampa chout, poi si ritorna al processo P1, dove chout = chin, ma chin è stato cambiato da P2, perciò verrà stampato nuovamente il secondo carattere. Se un solo processo alla volta potesse eseguire la procedura P2 va in stato di block perchè la risorsa (in questo caso la procedura) non è disponibile poiché usata da un altro processo. La computazione va avanti finché viene ripreso P1, viene conclusa la procedura e si ritorna a P2 con la risorsa ora disponibile, i due caratteri ora sono stampati in ordine. Anche nel multiprocessore ci sarebbe questo problema.
> ![[Screenshot From 2026-03-18 19-47-13.png|200]]
> ![[Screenshot From 2026-03-18 19-47-24.png|200]]
> ![[Screenshot From 2026-03-18 19-48-41.png|200]] > 

È necessario pertanto proteggere le variabili condivise e si deve controllare il codice che accede alla variabile.

**Race condition:** si verifica quando più processi o threads leggono e scrivono dati in modo che il risultato finale dipenda dall'ordine di esecuzione delle istruzioni dei processi

> [!example] I processi P1 e P2 assegnano rispettivamente 1 e 2 alla variabile condivisa a. Il valore finale sarà quello del ‘secondo arrivato’. I processi P1 e P2 condividono le variabili b e c con valori iniziali b = 1, c = 2. P1 esegue b = b + c, P2 esegue c = b + c. Se arriva primo P1, avremo b = 3; c = 5. Nel caso opposto, b = 4; c = 3

**Compiti del SO:** La concorrenza quindi può portare a problemi, è compito del sistema operativo fare in modo che la concorrenza si svolga nella maniera corretta:
- Tenere traccia dei processi (PCB)
- Allocare e deallocare le risorse per i processi attivi (tempo di elaborazione, memoria, files, dispositivi I/O)
- Proteggere dati e le risorse di ogni processo da interferenze involontarie di altri processi
- Il risultato di ogni processo deve essere indipendente dalla susa velocità relativa a quella dei processi concorrenti (salvo l'interazione voluta tra questi processi)

| Grado di consapevolezza                                                                                                     | Relazione                          | Influenza che un processo ha su un altro                                                                                                                                 | Potenziali problemi di controllo                                                                    |
| --------------------------------------------------------------------------------------------------------------------------- | ---------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | --------------------------------------------------------------------------------------------------- |
| Processi non consapevoli degli altri                                                                                        | Competizione                       | <ul><li>I risultati di un processo indipendente dall’azione degli altri</li><li>I tempi del processo potrebbero subire variazioni.</li></ul>                             | <ul><li>Mutua esclusione</li><li>Deadlock (risorsa rinnovabile)</li><li>Starvation</li></ul>        |
| Processi consapevoli indirettamente degli altri (cioè oggetti condivisi)                                                    | Cooperazione tramite condivisione  | <ul><li>I risultati di un processo potrebbero dipendere dale informazioni ottenute da altri</li><li>I tempi del processo potrebbero subire variazioni.</li></ul>         | <ul><li>Mutua esclusione</li><li>Deadlock (risorsa rinnovabile)</li><li>Coerenza dei dati</li></ul> |
| Processi consapevoli direttamente degli altri (hanno a disposizione<br>strumenti di comunicazione<br>primitivi<br>per loro) | Cooperazione tramite comunicazione | <br><ul><li>I risultati di un processo potrebbero dipendere dale informazioni ottenute da altri</li><li>I tempi del processo potrebbero subire variazioni.</li></ul><br> | <ul><li>Deadlock (risorsa rinnovabile)</li><li>Starvation</li></ul>                                 |
**Competizione per le risorse:** Cosa cambia quindi quando abbiamo processi indipendenti e processi cooperanti? La relazione più semplice è la competizione. I processi non sanno nulla (unaware) degli altri processi in corso e il risultato di un processo deve quindi essere indipendente dall’azione degli altri. Il tempo di elaborazione, però, potrebbe essere influenzato poiché una risorsa potrebbe non essere disponibile (come prima quando si blocca l’accesso a una procedura già usata) e si deve attendere che si liberi. I problemi di controllo possibili possono essere:
- Mutua esclusione (risorsa critica e sezione critica): cioè si impedisce l’accesso alla stessa risorsa da parte di più di un processo
- Stallo (*deadlock*): cioè tanti processi bloccati perché richiedono le stesse risorse già usate da altri processi e quindi al momento inaccessibili
- Starvation: un processo ha bisogno di una risorsa, ma gli altri processi la prendono prima di lui e esso non riesce mai ad accedere alla risorsa.  

```Mutua_esclusione
/* PROCESS 1 */ 
void P1 { 
	while (true) { 
		/* preceding code */; 
		entercritical (Ra); 
		/*critical section */; 
		exitcritical (Ra);
		/* following code */;
	}
}

/* PROCESS 2 */ 
void P2 { 
	while (true) { 
		/* preceding code */; 
		entercritical (Ra);  
		/*critical section */; 
		exitcritical (Ra); 
		/* following code */; 
	} 
}
```

**Collaborazione tramite la condivisione di risorse:** in questo caso i processi interagiscono senza conoscersi, ma usano e modificano i dati condivisi. Il risultato di un processo può dipendere dalle azioni di un altro e anche qui i tempi di elaborazione ne sono influenzati. I problemi sono sempre gli stessi, la mutua esclusione, poiché il meccanismo di controllo deve assicurare l’integrità dei dati (tranne per le operazioni di lettura), la deadlock e la starvation, ma anche un problema di data coherence. I processi modificano dei dati condivisi e se uno dei processi sbaglia a modificare i dati essi diventano incoerenti e portano a errori anche gli altri processi.

> [!example] 
> P1: a = a + 1, b = b + 1
> P2: b = 2 * b, a = 2 * b
> Posso ottenere situazioni come: a = a + 1, b = 2 * b, b = b + 1, a = 2 * b $\implies$ due risultati diversi

**Collaborazione per condivisione (condivisione per scambio di messaggi):** in questo caso i processi collaborano per un obiettivo comune e la comunicazione permette di sincronizzarsi e coordinarsi. Ci sono delle primitive che consentono a un processo di mandare informazioni a un altro processo o di riceverle fornite dal kernel o dal linguaggio di programmazione. In questo caso il risultato è sicuramente influenzato da ciò che il processo riceve, così come il tempo di elaborazione che è limitato a quando il processo riceve l’informazione necessaria. In questo caso non c’è problema di mutua esclusione poiché i processi non condividono dati, si hanno però sempre problemi di deadlock e starvation.

**Requisiti della mutua esclusione:** Quindi per la mutua esclusione dobbiamo far sì che: 
- Un solo processo alla volta è ammesso alla sezione critica per ciascuna risorsa condivisa
- Un processo che si ferma fuori dalla sua sezione critica non deve interferire con gli altri processi
- Si deve evitare deadlock e stravation
- Se nessun processo è nella  sezione critica deve essere concesso a ogni processo di entrare senza ritardo
- Nessuna ipotesi sulla velocità relativa dei processi e sul numero dei processori
- Un processo resa nella sezione critica per un tempo finito

**Mutua esclusione con supporto harware:** Per assicurare che un processore non possa uscire da una sezione critica di un processo su uniprocessore si possono disabilitare le interruzioni tramite delle primitive definite dal kernel, ciò però rallenta il sistema ed è inapplicabile per multiprocessori:
```
while(true){
	//disabilita interruzioni
	//sezione critica
	//abilita interruzioni
	//resto del programma
}
```

È anche possibile utilizzare una serie di istruzioni speciali:
- Un solo processore alla volta può accedere a una locazione di memoria 
- Speciali istruzioni macchina eseguono due operazioni su un dato in maniera atomica

Ci sono degli hardware con delle istruzioni come *compare and swap* in un’unica operazione:
```
in compare_and_swap(int *word, int testval, int newval){
	int oldval;
	oldval = *word;
	if (oldval == testval) *word = newval;
	return oldval
}

/* program mutualexclusion */ 
const int n = /* number of processes */; 
int bolt; 
void P( int i) { 
	while (true) { 
		while (compare_and_swap(bolt, 0, 1) == 1) 
			/* do nothing */;
		/* critical section */; 
		bolt = 0; 
		/* remainder */; 
	}
} 

void main() { 
	bolt = 0; 
	parbegin (P(1), P(2), ... ,P(n)); 
}
```
Un altro tipo di istruzione speciale è la *exchange*, una procedura che scambia il valore di 2 variabili.
```
/* program mutualexclusion */ 
int const n = /* number of processes */; 
int bolt; 
void P( int i) { 
	int keyi = 1; 
	while (true) { 
		do exchange (&keyi, &bolt) 
		while (keyi != 0); 
			/* critical section */; 
			bolt = 0; 
			/* remainder */; 
	}
} 

void main() { 
	bolt = 0; 
	parbegin (P(1), P(2), ..., P(n)); 
}
```

*Vantaggi:* 
- Applicabile a qualsiasi numero di processi su uniprocessore o multiprocessore 
- Semplicità 
- Può supportare multiple sezioni critiche

*Svantaggi:*
- Attesa attiva (*busy-waiting*) 
- Possibilità di stallo 
- Possibilità di starvation
###### Algoritmo di Dekker
L’obiettivo è quello di avere un protocollo che permetta a due processi di assicurare la mutua esclusione:
- Senza supporto hardware (salvo impossibilità di accesso simultaneo al medesimo indirizzo di memoria)
- Senza stallo
- Senza starvation

1. Primo tentativo: Apparentemente questo primo approccio funziona, c’è però un problema, supponiamo che il processo 0 esegua la sua sezione critica e che poi setti turn = 1, ma poi process 0 va in attesa, magari si vuole eseguire la sezione critica solo una volta, ma così le esecuzioni dei due processi si alternano.
<div style="display: flex; justify-content: center;">
  <img src="Screenshot From 2026-03-18 22-22-37.png" width="200">
</div>
2. Secondo tentativo: Si usano due flag che dovrebbero segnalare che la sezione critica sta venendo eseguita da uno dei due processi, questo non funziona poiché se i due processi arrivano quasi contemporaneamente con entrambe le flag false, entrano entrambi nella sezione critica poiché pensano che essa sia libera. 
<div style="display: flex; justify-content: center;">
  <img src="Screenshot From 2026-03-18 23-28-57.png" width="200">
</div>
3. Terzo tentativo: Scambiamo l’ordine delle operazioni di scrittura e lettura, prima occupo, poi aspetto che l’altra sia liberata, ma anche qui se i due processi partissero quasi insieme le due variabili sarebbero inizializzate entrambe a true, così entrambi i loop while sarebbero infiniti e i due processi non avanzeranno mai, è un esempio palese di stallo.
<div style="display: flex; justify-content: center;">
  <img src="Screenshot From 2026-03-18 23-29-45.png" width="200">
</div>
4. Quarto tentativo: Se un processo rileva che un altro processo ha flag true blocca se stesso, poi fa un delay che permette all’altro processo di avanzare e di entrare nella sezione critica per poi impostare flag a false. 
<div style="display: flex; justify-content: center;">
  <img src="Screenshot From 2026-03-18 23-30-23.png" width="200">
</div>

La soluzione corretta:
```Algoritmo_di_Dekker
boolean flag [2]; 
int turn; 
void P0() { 
	while (true) {
		flag [0] = true; 
		while (flag [1]) { 
			if (turn == 1) { 
				flag [0] = false; 
				while (turn == 1) /* do nothing */; 
				flag [0] = true; 
			} 
		} 
		/* critical section */; 
		turn = 1; 
		flag [0] = false; 
		/* remainder */; 
	} 
}
```
Con questa soluzione però paghiamo l’attesa attiva che è un costo non indifferente. L’algoritmo di Dekker si può scrivere in maniera più elegante con l’algoritmo di Peterson: 
```Algoritmo_di_Peterson
boolean flag [2]; 
int turn; 
void P0() { 
	while (true) { 
		flag [0] = true; 
		turn = 1; 
		while (flag [1] && turn == 1) /* do nothing */;
		/* critical section */; 
		flag [0] = false; 
		/* remainder */; 
	} 
}
```

###### Semafori
Una variabile con valore intero, una coda e tre operazioni:
- Inizializzazione a un valore non negativo 
- *semWait* decrementa il calore del semaforo. Se il valore diventa negativo, il processo è sospeso e inserito nella coda
- *semSignal* incrementa il valore del semaforo. Se il valore non diventa positivo, uno dei processi della coda è riattivato

**Osservazione:** il processo non ha modo di sapere in anticipo:
- Se con l’esecuzione di *semWait* sarà sospeso
- Se c’è un porcesso che attende l’esecuzione di *semSignal*
- Quale fra i processi in attesa sarà riattivato da *semSignal*

```Semaforo_generico
struct semaphore { 
	int count; 
	queueType queue; 
};

void semWait(semaphore s) { 
	s.count--; 
	if (s.count < 0) { 
		/* place this process in s.queue */; 
		/* block this process */; 
	}
}

void semSignal(semaphore s) { 
	s.count++; 
	if (s.count<= 0) { 
		/* remove a process P from s.queue */; 
		/* place process P on ready list */; 
	}
}
```

**Semaforo binario:** una variabile con valore binario (o 0 o 1), una coda e tre operazioni:
- Inizializzazione a 0 o 1
- *semWaitB* verifica il valore del semaforo. Se il valore è 1 allora è portato a 0, nel caso contrario, il processo è sospeso e inserito nella coda
- *semSignalB* verifica se ci sono processi in attesa nella coda. In caso affermativo, uno dei processi della coda è riattivato, nel caso opposto, il valore del semaforo è portato a 1

**Gestione della coda:** 
- Deve evitare la starvation
- Generalmente FIFO (*semaforo forte*)

```Semaforo_binario
struct binary_semaphore { 
	enum {zero, one} value; 
	queueType queue; 
}; 

void semWaitB(binary_semaphore s) { 
	if (s.value == one) 
		s.value = zero; 
	else { 
		/* place this process in s.queue */; 
		/* block this process */; 
	}
}

void semSignalB(semaphore s) { 
	if (s.queue is empty()) 
		s.value = one; 
	else { 
		/* remove a process P from s.queue */;
		/* place process P on ready list */; 
	}
}
```

> [!example] ![[Pasted image 20260319122744.png|200]]

```Mutua_esclusione
/* program mutualexclusion */ 
const int n = /* number of processes */; 
semaphore s = 1; 

void P(int i) {
	 while (true) { 
	 	semWait(s); 
	 	/* critical section */; 
	 	semSignal(s); 
	 	/* remainder */; 
	 } 
} 

void main() { 
	parbegin (P(1), P(2), ..., P(n)); 
}
```

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319123611.png" width="200">
</div>

**Implementazione:** *semWait* e *semSignal* devono essere primitive atomiche
- hardware o firmware
- con istruzioni speciali
- con l’algoritmo di Peterson

**Messaggi:** nell’interazione fra processi sono necessarie:
1. *Sincronizzazione:* per la realizzare la mutua esclusione
2. *Comunicazione:* per lo scambio di informazioni
3. *Scambio di messaggi:* permette di realizzare entrambe
	- Naturale nei sistemi distribuiti
	- Utile anche nei sistemi con memoria condivisa
4. *Due primitive:* 
	- *Send (destination, message):* un processo invia informazioni in forma di messaggio a un altro processo designato come destinatario
	- *Receive (source, message):* un processo riceve informazioni indicando il mittente e l’indirizzo in cui salvare il messaggio

**Caratteristiche dei sistemi di messaggistica:** 
1. *Sincronizzazione:* 
	- Send: può essere o bloccante (finché il destinatario non esegue la receive, il processo che ha eseguito la send si blocca) o non bloccante
	- Receive: può essere bloccante o non bloccante. In alternativa viene fatto un test di presenza messaggi (se ci sono viene lanciata una receive ?????)

> [!info] ci sono situazioni in cui sia la send che la receive sono bloccanti

2. *Indirizzamento:*
	- Diretto: sia per la send che per la receive (esplicito o implicito). 
	- Indiretto: statico o dinamico
3. *Formato:*
	- Contenuto
	- Lunghezza (fissa o variabile)
4. *Organizzazione delle code:* se ci sono più messaggi, vanno organizzati
	- FIFO (cioè in code)
	- Priorità

**Sincronizzazione:** il destinatario non può ricevere il messaggio prima che questo sia stato inviato. Quando si esegue la receive:
- Se è presente un messaggio, il processo lo riceve e continua l’esecuzione
- In caso contrario ci sono due possibilità: o si sospende finché il messaggio non è inviato oppure si continua l’esecuzione senza più ricevere il messaggio

| Tipo di send e receive                 | Comportamento                                                                                                                                                                                                                                                                     |
| -------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Send e receive bloccanti               | <ul><li>Mittente e destinatario bloccati fino al completamento dell’operazione</li><li>Rendez-vous</li><li>Permette la sincronizzazione dei processi</li></ul>                                                                                                                    |
| Send non bloccante e receive bloccante | <ul><li>Il mittente procede ma il destinatario è bloccato fino alla ricezione del messaggio</li><li>Combinazione più comune</li><li>Il mittente invia più messaggi a tanti destinatari senza perdita di tempo (es. Processo server che fornisce servizi a più processi)</li></ul> |
| Send e receive NON bloccanti           | <ul><li>Nessuno dei due processi attende</li></ul>                                                                                                                                                                                                                                |
**Indirizzamento:** 
1. *Indirizzamento diretto:* 
	- La primitiva send include l’identificatore del processo destinatario
	- La primitiva receive può avere sia l’indirizzamento esplicito (destinatario designato esplicitamente, è utile per processi concorrenti cooperanti) che quello implicito (sender è un parametro che restituisce l’identificatore del mittente)
2. *Indirizzamento indiretto:*
	- Messaggi inviati a una struttura dati condivisa consistente in code che contengono temporaneamente i messaggi (mailbox)
	- Il mittente invia un messaggio a una mailbox
	- Il destinatario preleva un messaggio dalla mailbox
	- Flessibilità

> [!info] Indirizzamento indiretto 
> ![[Pasted image 20260319132049.png|200]]

**Formato del messaggio:** 
<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319132443.png" width="200">
</div>

**Gestione delle code:** 
- FIFO
- Priorità (per tipo di messaggio o designata dal mittente)
- Il destinatario può ispezionare la coda e scegliere il messaggio

```Mutua_esclusione
/* program mutualexclusion */ 
const int n = /* number of process */; 
void P(int i) { 
	message msg; 
	while (true) { 
		receive (box, msg); 
		/* critical section */; 
		send (box, msg); 
		/* remainder */; 
	} 
} 

void main() { 
	create mailbox (box); 
	send (box, null); 
	parbegin (P(1), P(2),..., P(n)); 
}
```

**Problema produttore/conduttore:** Condizioni:
- Uno o più produttori inseriscono dati in buffer
- Un consumatore estrae i dati dal buffer, uno alla volta
- Un solo produttore o consumatore alla volta accede al buffer

Il problema è che:
- Assicurare la mutua esclusione
- Evitare che un produttore aggiunga dati a un buffer pieno
- Evitare cue il consumatore cerchi di estrarre dati da un buffer vuoto

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319133107.png" width="200">
</div>
