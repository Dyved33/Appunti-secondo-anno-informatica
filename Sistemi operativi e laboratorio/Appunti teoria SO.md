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
* Il sistema operativo alterna rapidamente l'esecuzione di ciascun programma utente assegnando un **quanto di elaborazione** (un breve periodo di tempo).
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
- CPU = costituita dai registri (es. Program counter che contiene l'indirizzo dell'istruzione successiva da eseguire o l'instruction register che contiene l'istruzione da eseguire,...)
- Memoria principale = vengono caricati i programmi che sono costituiti da istruzioni (quelli che la CPU può eseguire) e dati (locazioni di memoria in cui mettiamo le nostre informazioni). 
- Moduli I/O = relativi a tastiera, monitor e dischi
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
  Quindi il procedimento è: estrazione istruzione -> esecuzione istruzione -> check delle interruzioni. Se c'è l'interruzione se ne guarda il tipo si compie un salto all'indirizzo al quale inizia il programma per la gestione delle interruzioni (*interrupt handler*) e successivamente si salta ad un'altra istruzione. Una volta che il gestore delle interruzioni ha finito di gestire il problema legato ad un'interruzione, questa viene eseguita. 
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
> ![[Screenshot From 2026-03-05 12-31-03.png|300]]

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

**Generazione sistema operativo:** Scritto il codice sorgente del sistema $\rightarrow$ va compilato e configurato per il sistema su sui verrà eseguito $\rightarrow$ poi si istalla il sistema $\rightarrow$ avvia il computer e il nuovo sistema operativo. La configurazione può realizzarsi modificando il codice prima della compilazione in base al sistema su cui lo vado a installare, oppure tramite la selezione di tabelle e moduli precompilati da collegare al momento dell'installazione, oppure tramite un sistema guidato da tabelle di dati che guidano il sistema in base all'hardware.<u> Per far partire il sistema l'esecuzione inizia in un predefinito indirizzo della ROM (parte di memoria non modificabile nemmeno a computer spento = *read only memory*) che contiene il *bootstrap loader*. Questo contiene le istruzioni necessarie a caricare il *bootstrap program* che a sua volta carica il kernel e avvia l'esecuzione del sistema</u>
###### Processi
**Definizione:**<u> un processo (o job) è un programma (lista di istruzioni con riferimenti a delle locazioni di memoria) in esecuzione</u>. Gli elementi del processo sono:
- Program counter (indirizzo istruzione da eseguire) e contenuto dei registri (immagine del processore)
- Stack che contiene variabili locali, indirizzi di return o parametri delle funzioni
- Sezione testo (programma)
- Sezione dati che contiene i valori delle variabili locali
- Heap per l'allocazione dinamica

| Programma                        | Processo                                                   |
| -------------------------------- | ---------------------------------------------------------- |
| Entità passiva / file eseguibile | Entità che esegue le istruzioni e usa le risorse associate |
> [!tip] Ricorda
>  - Un programma può essere associato a più processi
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
    Supponiamo che venga creato il processo P 0 che è in stato di running. Ad un certo punto si verifica un'interruzione e tutte le informazioni del processo vengono salvate nel PCB 0. Queste saranno ricaricate dal sistema di gestione delle interruzioni dopo che il processore avrà terminato un altro processo P 1 (ricaricato da PCB 1 quando magari P 1 era stato interrotto precedentemente). A questo punto può accadere una nuova interruzione, si salva lo stato di P 1 dentro a PCB 1 e si ricarica lo stato contenuto in PCB 0 (si riprende esecuzione di P 0). Questo ci permette di usare il sistema in maniera efficiente purché le operazioni di salvataggio e ricarica degli stati del processore siano eseguite in maniera efficiente.
  </div>
</div>

Il sistema ha bisogno di conoscere la situazione dei vari processi e per questo i PCB sono inseriti in varie code:
- Coda di job che contiene tutti i processi 
- Coda di ready che contiene tutti i processi pronti a partire (quelli tra i quali il processore deve decidere a chi assegnare le risorse)
- Coda dei dispositivi (una per dispositivo) che contiene i processi in attesa di interagire con quel dispositivo (il dispositivo genera una nuova interruzione)
###### Schedulatori e fasi del processo
Il sistema operativo deve decidere a quali processi assegnare il processore e per fare ciò ci sono diversi schedulatori:
- **Schedulatore a breve termine (CPU scheduler):** Prende i processi dalla coda ready e gli assegna un processore. Siccome deve essere eseguito con grande frequenza, deve essere molto veloce (i processi running vengono anche messi in coda ready proprio perché lo scheduler assegna la CPU a un processo, poi la toglie a quello per assegnarla ad un altro, così da far avanzare più processi contemporaneamente).
- **Schedulatore a lungo termine:** Seleziona tra i processi nuovi quali fa passare in coda ready, ma non è opportuno che i processi finiscano tutti in questa coda. Per questo esiste il *job scheduler*: Seleziona i processi da inserire in coda ready e per fare ciò controlla il grado di multiprogrammazione, cioè i programmi eseguiti contemporaneamente. Non viene eseguito con frequenza, per questo può applicare procedure più complesse. Ci sono ad esempio processi CPU-bound, ovvero processi che non rilasciano mai a meno che non forzati la CPU a loro assegnata 

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
Il sistema operativo fornisce meccanismi per creazione/terminazione dei processi. Generalmente un processo genitore chiede al kernel di creare un processo figlio che a sua volta può creare altri processi (*albero dei processi*). Ogni processo è identificato da un numero detto *identificatore di processo (pid)* $\rightarrow$ su linux è il primo numero che vedo quando scrivo `ps` nel terminale. Quando un processo termina, viene eliminato dall'albero.

**Assegnazione delle risorse ad un processo figlio**
Avviene in base al tipo di sistema operativo:
- Un processo figlio riceve le risorse dal processo padre 
 - Risorse ripartite = il padre non può usarle finché il figlio non ha finito di usarle
 - Risorse condivise 
- Il processo figlio ottiene le risorse dal sistema operativo

Genitore e figlio possono essere eseguiti concorrentemente mentre il genitore attende la terminazione del figlio. Per quanto riguarda la memoria il processo figlio può avere un duplicato di memoria del genitore e quindi avere lo stesso programma, oppure può averne uno nuovo.

> [!info] La funzione `fork` è una syscall che crea una copia esatta del processo e restituisce al padre il pid del figlio:
> - **0 (Zero):** Indica che il processo corrente è il **processo figlio** appena creato.
> - **Maggiore di 0 (Positivo):** Indica che il processo corrente è il **processo padre** e il valore corrisponde al PID del figlio.
> - **-1 (Negativo):** Indica che la `fork()` è fallita e nessun processo figlio è stato creato.

**Terminazione del processo**
Il processo termina con la syscall *exit ()*, così:
- Restituisce un valore di stato al processo genitore, via *wait ()*
- Le risorse sono deallocate dal sistema operativo
Il genitore può terminare l'esecuzione del figlio con la syscall *abort ()* e questo per varie motivazioni:
- Eccesso di risorse allocate
- Compito non più richiesto
La terminazione di un processo figlio fa si che vengano terminati anche i figli del processo terminato (terminazione a cascata). Dolo dopo il wait () del genitore viene deallocato il pid: 
- il pid però non scompare subito e il processo rimane in stato di zombie e scompare solo dopo che il genitore ha eseguito il wait ()
- il genitore termini senza che vengano terminati i suoi processi figli, in tal caso quei processi si dicono orfani.
###### **Threads**
I processi sono caratterizzati da:
- **Risorse allocate:** spazio di indirizzamento virtuale (immagine in memoria), la traccia di esecuzione, RAM, dispositivi I/O, file
- **Scheduling:** traccia di esecuzione (alternata con altri processi), stati e priorità
Ma le due caratteristiche potrebbero essere trattate in maniera indipendente:
- **Processo (o task):** possiede le risorse
- **Thread (o processo leggero):** unità che viene assegnata al processore. 
Si parla di multithreading la capacità di un SO di permettere più tracce di esecuzione concorrenti in un singolo processo (quindi ad un processo o più processi vengono assegnati più di un singolo thread ciascuno)

| Single-threading                                         | Multi-threading                                                |
| -------------------------------------------------------- | -------------------------------------------------------------- |
| Un processo, un thread (MD-DOS)                          | Un processo, più thread (Java run-time)                        |
| Più processi, un thread per ciascuno (UNIX tradizionale) | Più processi, ciascuno con più thread (Windows, Solaris, UNIX) |
**Definizione:** in un ambiente multi-thread, il processo è l'unità di allocazione delle risorse e un'unità di protezione (nessun altro può accedere alle risorse assegnate al processo). Associati al processo ci sono:
- Spazio di indirizzamento virtuale (contiene immagine del processo)
- Accesso protetto al processore, ad altri processi (IPC), files, dispositivi I/O
- Uno o più thread
Associati al thread ci sono:
- Stato = pronto, in attesa, in esecuzione
- Contesto del thread (immagine del processore). Viene salvata quando non è in esecuzione
- Stack di esecuzione
- Spazio di memoria statico per le variabili locali
- Accesso a memoria e risorse del processo, condivise con gli altri thread del processo
Questo genera un problema: due thread che accedono agli stessi dati/stessa memoria (magari uno vuole incrementare e l'altro decrementare la stessa variabile)
<div style="display: flex; justify-content: center;">
  <img src="Screenshot From 2026-03-12 19-20-44.png" width="300">
</div>
I benefici del multithreading sono:
- La creazione di un thread è più rapida della creazione di un processo
- La teminazione di un thread è più rapida della terminazione di un processo (no questione di protezione della memoria)
- Lo switch tra due thread del medesimo processo è più rapida dello switch tra due processi (sempre per via della memoria)
- Efficianza della comunicazione fra programmi (eseguiti da thread dello stesso processo)
In un sistema a utente singolo:
- Esecuzione in foreground e background
- Elaborazione asincrona
- Velocità di esecuzione
- Struttura modulare dei programmi

**Osservazione:** Si nota che la schedulazione è gestita principalmente nel thread, ma non la sospensione e la terminazione del processo.

Gli stati del thread sono gli stessi dei processi (in esecuzione, pronto e in attesa), così come le operazioni (creazione, block/sospensione, unblock/riattivazione, terminazione).
<div style="display: flex; justify-content: center;">
  <img src="Screenshot From 2026-03-12 19-24-20.png" width="300">
</div>

**Sincronizzazione dei thread:** I thread di un processo condividono spazio di indirizzamento e altre risorse. Ogni modifica simultanea di un dato da parte di due thread può corrompere il dato stesso. Per questi motivi è necessario sincronizzare l'accesso dei thread alle risorse

**Thread a livello utente:** La gestione dei thread è fatta dall'applicazione e il kernel non sa dell'esistenza dei thread.
- *Vantaggi*: lo scambio di thread non richiede privilegi della modalità kernel, la schedulazione può essere specifica per l'applicazione e portabilità.
- *Svantaggi:* le chiamate di sistema sono bloccanti per l'intero processo e non ci sono vantaggi dai multiprocessori
- *Contromisure:* programmare le applicazioni in processi multipli e *jacketing* (trasforma una chiamata di sistema bloccante in una non bloccante)

**Thread a livello kernel:** La gestione dei thread è fatta dal kernel e l'applicazione usa un'API per la gestione dei thread. Il kernel gestisce i PCB dei processi e dei thread e inoltre la schedulazione del processore sulla base dei thread.
- *Vantaggi:* possibilità di eseguire thread differenti su diversi processori, il blocco di un thread non blocca il processo e anche le routine del kernel possono essere gestite in multithreading
- *Svantaggi:* i cambi di modalità (utente/kernel) producono overhead

**Approcci misti:** - Creazione dei thread, schedulazione e sincronizzazione nello spazio utente
- Thread modalità utente mappati su thread modalità kernel (in numero minore o uguale)
- Permette esecuzione parallela dei thread e il blocco di un thread non blocca il processo (come solaris)

| Modello | Descrizione                                                                                                                                              |
| ------- | -------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1:1     | Ogni thread di esecuzione è un processo unico con il suo spazio e le sue risorse                                                                         |
| M: 1     | Un processo definisce uno spazio di indirizzo e il possesso di risorse dinamiche. Molti threads possono essere creati ed eseguiti in quel processo       |
| 1: M     | Un thread può migrare dall'ambiente di un processo all'ambiente di un altro processo. Ciò permette al thread di muoversi facilmente tra sistemi distinti |
| M: N     | Combina gli attributi di 1: M e M: 1                                                                                                                       |

**Multicore e multithreading:** **Legge di Amdahl:** $\text{speedup} = \frac{\text{tempo di esecuzione su un processore}}{\text{tempo di esecuzione su N processori}} = \frac{1}{(1-f)+ \frac{f}{n}}$
Il vantaggio lo hanno:
- Applicazioni multithread native (pochi processi e molti thread)
- Applicazioni multiprocesso (molti processi single-thread)
- Applicazioni Java (JVM è un processo multithread e assicura schedulazione e gestione della memoria per le aplllicazioni Java)
- Multiple istanze di un'applicazione

**Threads in UNIX:**
- Il kernel non distingue tra processi e thread
- Per ogni task (processo o thread) c'è un *task control block* (TCB) con fork ()
- Fork () crea una nuova task, con una copia del task control block
- Clone () crea una nuova task, ma nel suo TCB ci sono puntatori ai campi di quello del genitore (in base ai parametri passati a clone ()

**Mutua esclusione:** I processi possono interferire tra di loro in maniera positiva o negativa, negativa ad esempio quando sottrae risorse ad altri processi, in maniera positiva quando due o più processi interagiscono e ci permettono di realizzare un’attività che sfrutta l’uso di diversi processi in un’unica applicazione. La gestione dei processi può avvenire tramite: 
 - Multiprogrammazione: molti processi con un singolo processore, i processi avanzano attendendo che il processore si liberi; 
 - Multiprocessing: molti processi gestiti da un multiprocessore; 
 - Processi distribuiti: molti processi, gestiti da molti processori distribuiti (*cluster*), con memorie distribuite, rende le cose più complicate da gestire, ma sicuramente più veloci. 
 I differenti processi hanno bisogno di risorse che tipicamente vengono condivise tra i vari processi, si crea quindi una *concorrenza*:

~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
 - Nelle applicazioni multiple, o nelle applicazioni strutturate (con più processi o threads che devono accedere a delle risorse)
 - La modifica dei dati disponibili da parte di un processo influenza l’andamento di un altro processo
 - Nella struttura del sistema operativo.

**Principi della concorrenza**
- Nell'uniprocessore multiprogrammato, le esecuzioni dei processi si alternano
- Nel multiprocessore, le esecuzioni dei processi si alternano e si sovrappongono
- In entrambi i casi, la velocità relativa di esecuzione dei processi non è prevedibile. Dipende da:
 - Attività degli altri processi
 - Gestione delle interruzioni
 - Politiche di schedulazione
Vengono generate le seguenti difficoltà:
- Condivisione di risorse globali
- Allocazione ottimale delle risorse
- Debugging 

È necessario proteggere le variabili condivise e si deve controllare il codice che accede alla variabile.
**Race condition:** si verifica quando più processi o threads leggono e scrivono dati in modo che il risultato finale dipende dall'ordine di esecuzione delle istruzioni dei processi

Compiti del SO:
- Tenere traccia dei processi (PCB)
- Allocare e deallocare le risorse per i processi attivi (tempo di elaborazione, memoria, files, dispositivi I/O)
- Proteggere dati e le risorse di ogni processo da interferenze involontarie di altri processi
- Il risultato di ogni processo deve essere indipendente dalla susa velocità relativa a quella dei processi concorrenti (salvo l'interazione voluta tra questi processi)

Processi concorrenti sono in conflitto per l'uso delle risorse. I problemi sono:
- Mutua esclusione (risorsa critica e sezione critica)
- Stallo (*deadlock*)
- Starvation
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~