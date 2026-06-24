### Cenni Storici e Famiglie di Sistemi Operativi
*Unix*
* **Unix** è una "famiglia" di sistemi operativi **multitasking** e **multiuser** che deriva dall'originale AT&T Unix.
* Lo sviluppo di AT&T Unix è iniziato nel 1969 al Bell Labs research center, principalmente ad opera di **Ken Thompson**, **Dennis Ritchie** e altri collaboratori.
* Inizialmente si è assistito a una forte frammentazione dove "ognuno si fece il suo", dando vita a varianti come BSD (UC Berkeley), Xenix (Microsoft), SunOS/Solaris (Sun Microsystems), HP-UX (HP), e AIX (IBM).
* Nel corso del tempo, sono stati fatti sforzi per standardizzare le varie release e varianti attraverso iniziative come **POSIX**, **SUS**, e **The OpenGroup**.

*(GNU/)Linux*
* **Linux** è un discendente o derivato di Unix, definibile come **Unix-like**.
* È una "famiglia" **open-source** di sistemi operativi Unix-like.
* Si basa sul kernel Linux, rilasciato per la prima volta nel **1991** da **Linus Torvalds**.
* Tecnicamente, sarebbe più corretto parlare di sistema operativo **GNU/Linux**.
* È un sistema **POSIX compatibile**, seppur non ufficialmente POSIX-certified.
* Esistono varie **distribuzioni** (o "distro"), tra cui Debian, SUSE, RedHat, Fedora, Ubuntu e Mint.

*macOS*
* Rappresenta un gruppo di sistemi operativi proprietari di casa **Apple**.
* Sono sistemi Unix-like e, soprattutto, **Unix certified** (certificati da The Open Group).
* Hanno attraversato diverse transizioni di **processori**: inizialmente PowerPC, poi dal 2006 sono passati ad Intel, e dal 2020 utilizzano processori 64-bit ARM-based come l'Apple M1 (seguiti da M2, M3, ecc.).
* Ad Ottobre 2024, Apple ha presentato il chip **M4**.

*Android*
* Sistema basato su un **kernel Linux modificato**.
* Originariamente sviluppato dalla Open Handset Alliance, la sua versione più diffusa è oggi sviluppata principalmente da **Google**.
* Il resto del sistema, inclusi software vari e utilities, è scritto principalmente in **Java**.
* È specificamente pensato per sistemi **touchscreen mobile**, come telefoni e tablet.
### Concetti Base dei Sistemi Unix
**Struttura e Architettura:**
* Sistemi **multitasking e multiuser**.
* **Modularità**: Unix non è un sistema monolitico.
    * Il **kernel** fornisce i servizi fondamentali del sistema.
    * La **shell** (o in alternativa l'interfaccia grafica) permette all'utente di interagire con il sistema operativo.
    * Esistono tante semplici utility che svolgono compiti base; queste possono essere combinate insieme tramite **pipe line** per eseguire compiti anche molto complessi.
* È stato concepito appositamente per supportare gruppi di sviluppatori che lavorano insieme.
* **Portabilità**: Essendo stato quasi subito riscritto in **C**, risulta molto più portabile rispetto ai vecchi sistemi scritti in assembly.
* I device e i vari strumenti per l'IPC (Inter-Process Communication) sono trattati e mostrati dal sistema come se fossero dei semplici **file**.
* Il formato di memorizzazione dei file è volutamente semplice (**plain text**).

 **Autenticazione e Utenti:**
* Essendo multiutente, ogni utente possiede una **login**, una **password** e una **home directory** (uno spazio disco personale dove salvare i propri file).
* Esiste un utente privilegiato, chiamato **root**, che funge da amministratore del sistema.
* Ogni utente gestisce in autonomia i permessi di accesso sulle proprie cartelle e sui propri file.
* Gli utenti possono essere aggregati in **gruppi** per definire e semplificare privilegi comuni.

**Terminale, Console e BASH**
**Definizioni:**
* **Terminal**: Un ambiente di input e output testuale. Una finestra di terminale (o emulatore di terminale) è una finestra di solo testo che emula una console all'interno di una GUI (Graphical User Interface).
* **Console**: Un terminale fisico; un pannello della strumentazione contenente i controlli del computer. Una console è quindi un tipo specifico di terminale.
* **Shell**: È l'interfaccia vera e propria tra il sistema operativo e l'utente. Permette di eseguire programmi e offre funzionalità simili a un linguaggio di programmazione per automatizzare operazioni. Può essere a riga di comando (Bash, csh, ksh) o grafica (GNOME, KDE).

**BASH (Bourne Again Shell):**
* È un interprete interattivo di comandi per sistemi Unix-like.
* È stata creata nel **1989** da **Brian Fox** per il progetto GNU, nata come alternativa ed estensione della preesistente sh (Bourne shell).
* È la shell di default in moltissime distribuzioni Linux ed è stata la default shell di macOS prima del 2019.
* È un normale programma che gira nello "spazio utente".
* Espone il **prompt** e attende che l'utente inserisca una linea di comando; successivamente la interpreta, la esegue e torna in attesa.
* Per terminare la sessione BASH si usa il comando `exit` oppure la scorciatoia da tastiera **EOF (Ctrl+d)**.

> [!IMPORTANT] Il Prompt di BASH
> Il prompt è la stringa di caratteri con cui la shell segnala di essere in attesa di un comando. 
> Esempio tipico: `fabio@fabiopc:/usr/lib/binfmt.d$`
> * **fabio**: username.
> * **fabiopc**: nome host (da non confondere col nome di dominio).
> * **/usr/lib/binfmt.d**: path corrente della directory di default in quel momento.
> * **$**: livello di privilegio dell'utente (il carattere `#` indica l'utente root, mentre `$` indica gli altri utenti).

**Perché usare una shell testuale?** Anche se esistono le interfacce grafiche, la shell testuale si usa perché:
* Aiuta a capire come funziona il sistema "sotto".
* Offre il massimo controllo sui comandi e sulle loro opzioni.
* Permette di configurare le distribuzioni in modo preciso.
* I comandi Unix sono stati nativamente progettati per le interfacce a riga di comando.
* Una volta imparati, risultano più veloci e flessibili.
* Sono estremamente veloci e facilmente accessibili da remoto.
* Funzionano sempre, anche in caso di problemi all'interfaccia grafica o con hardware obsoleto e poco performante.

**Il Filesystem:** Il filesystem è una struttura ad albero al cui interno sono memorizzati tutti i dati del sistema operativo.
* È organizzato in **File** e **Directory** (o cartelle).
* L'albero ha un'unica radice detta **root**, identificata dal path `/`. A differenza di Windows (che ha `C:`, `D:`, ecc.), in Unix tutti i dati convergono in quest'unico albero, indipendentemente dal numero di dischi o partizioni fisiche. 

**Navigazione e Pathnames:** Ogni processo nel sistema ha, in ogni singolo istante, una "directory corrente".
* **Pathname assoluto**: È il percorso che parte dalla radice root e porta all'oggetto. Comincia sempre con la barra `/` (es. `/etc/hosts` o `/home/fabio/Desktop`).
* **Pathname relativo**: È il percorso che porta all'oggetto partendo dalla directory corrente o da un punto specifico che non è la radice. Non inizia mai con `/` (es. `log/apache`).
* **Simboli speciali nei path**:
    * `.` (punto): indica la directory corrente.
    * `..` (due punti): indica la directory "padre", ovvero quella che contiene la directory corrente.

**Link e Device Files:** Nel filesystem possono essere presenti file speciali chiamati **link**, che fungono da puntatori/collegamenti ad altri oggetti (in modo analogo ai collegamenti di Windows). Compiere un'operazione su un link significa implicitamente compierla sull'oggetto a cui esso punta.
* Sono presenti anche **special files** o **device files**: espongono un dispositivo hardware all'interno del fylesystem come se fosse un normale file. Questo fornisce un'interfaccia verso l'hardware, permettendo alle system calls di eseguire operazioni di Input/Output in modo standardizzato.
### Primi Comandi BASH e variabili:
I comandi BASH hanno tipicamente questa sintassi: `nome_comando <opzioni> <argomenti>`. 
Le **opzioni** iniziano tipicamente con il carattere `-` e modificano il comportamento del comando; possono essere accorpate (es. `ls -ltra` equivale a `ls -l -t -r -a`) e a volte richiedono argomenti aggiuntivi. Ogni comando restituisce un valore di ritorno che informa sull'esito, dove solitamente `0` indica che è andato tutto a buon fine.

**Comandi Interni vs Esterni:** 
- **Comandi Esterni**: Sono veri e propri eseguibili registrati nei supporti di memorizzazione. La shell li cerca all'interno dei percorsi di filesystem specificati nella variabile di ambiente `$PATH`. Il comando `which` permette di scoprire dove è memorizzato l'eseguibile di un comando esterno.
* **Comandi Interni**: Sono riconosciuti ed eseguiti direttamente dalla shell stessa (es. i costrutti per lo scripting o il comando `cd`). È possibile vederne la lista utilizzando il comando interno `help`.

**Operazioni sul Filesystem**

| Comando        | Descrizione                                                                                                              |
| :------------- | :----------------------------------------------------------------------------------------------------------------------- |
| **`ls`**       | Elenca file e directory. Può prendere file/cartelle come argomenti specifici (es. `ls -l lettera foto`).                 |
| **`cd`**       | Cambia la directory corrente (è un comando interno). Esempio: `cd /etc`.                                                 |
| **`cp`**       | Copia file. Es: `cp pippo pluto`. Se pluto è una directory, pippo viene messo al suo interno.                            |
| **`mv`**       | Sposta file o cartelle. Se i file hanno percorsi riferiti alla stessa directory, il comando rinomina il file originario. |
| **`mkdir`**    | Crea nuove cartelle.                                                                                                     |
| **`rm`**       | Rimuove file.                                                                                                            |
| **`rmdir`**    | Rimuove cartelle, a condizione che siano **strettamente vuote** (altrimenti genera un errore).                           |
| **`ln`**       | Crea collegamenti a file o directory.                                                                                    |
| **`pwd`**      | Visualizza la directory corrente in cui ci si trova.                                                                     |
| **`df`**       | Elenca i filesystem attualmente montati.                                                                                 |
| **`cat`**      | Visualizza per intero il contenuto di un file di testo.                                                                  |
| **`realpath`** | Visualizza il path assoluto "risolto" (molto utile per calcolare la vera destinazione di link o riferimenti relativi).   |
| **`file`**     | Ritorna il tipo reale di un file (testo, binario, pdf) in maniera totalmente indipendente dall'estensione nominale.      |
| **`tail`**     | Visualizza le ultime "n" linee di un file.                                                                               |
| **`head`**     | Visualizza le prime "n" linee di un file.                                                                                |
| **echo**       | Stampa nel terminale cosa gli passo (sia stringa che variabile)                                                          |
**Opzioni più Comuni per `ls`:** `-l`: Formato lungo (long format), mostra informazioni aggiuntive sui file.
* `-a`: Elenca tutti i file, compresi quelli "nascosti" (quelli che iniziano col carattere `.`).
* `-t`: Ordina i file mostrati in base alla data di ultima modifica (dal più recente al meno recente).
* `-r`: Inverte il senso dell'ordinamento.
* `-h`: Visualizza le dimensioni dei file in un formato "human readable" e facilmente leggibile.

> [!EXAMPLE] Esempio combinato
> Il comando `ls -lrta` elenca tutti i file e le cartelle (compresi i nascosti), ordinandoli in base al tempo di ultima modifica (dal meno recente al più recente) formattandoli in modalità lunga.

**Opzioni più Comuni per CANCELLARE o COPIARE Cartelle:** `rm -r cartella`: Rimuove ricorsivamente una cartella e tutto il ramo di file/cartelle non vuote in essa contenute.
* `cp -ar dir1 dir2`: Copia ricorsivamente (`-r`) tutto l'albero sottostante `dir1` dentro `dir2`, preservando (`-a`) gli attributi originari dei file come la data di modifica.

 **Differenza nella creazione dei Link (`ln`):** 
 - **Hard Link** (`ln pippo link1`): Crea una vera e propria copia che occupa spazio fisico sul disco. Se il file originario "pippo" viene cancellato, il file "link1" e il suo contenuto rimangono inalterati.
 - **Symbolic Link** (`ln -s pippo link1`): Crea un collegamento simbolico (occupazione minimale su disco). "link1" contiene solo una stringa col percorso di "pippo". Se il file "pippo" originario viene cancellato, "link1" diventa semplicemente un collegamento non valido o rotto.

**Ottenere Aiuto e Documentazione:** Per imparare le funzioni e le opzioni dei comandi, ci sono diverse utility fondamentali:
* **`man <nome_comando>`**: Apre il manuale in linea completo del comando specificato.
* **`apropos <stringa>`**: Cerca e lista i comandi che contengono la stringa specificata all'interno della loro descrizione sintetica.
* **`whatis <comando>`**: Fornisce una descrizione estremamente sintetica del comando passato.
* **`locate <stringa>`**: Ricerca velocemente tutti i file aventi un path che contiene la stringa specificata, interrogando un database interno costruito tramite il tool `updatedb`.
* **`Switch di aiuto rapido`**: Generalmente ogni comando BASH, se invocato aggiungendo lo switch `-h` oppure `--help`, è in grado di stampare a terminale una rapida e sintetica descrizione delle proprie funzionalità.

**Variabile $PATH:** Se voglio vedere dove è memorizzato l'eseguibile che corrisponde ad un comando (esterno) uso il comando **which**. Uso la variabile d'ambiente $**PATH** per dire al terminale dove cercare i programmi (è un elenco di directory in cui il sistema operativo cerca i programmi eseguibili quando digiti un comando nel terminale. Essa permette di avviare applicazioni da qualsiasi posizione senza dover digitare l'intero percorso del file). Quindi, quando si desidera eseguire un programma, è sufficiente digitare il nome del file eseguibile e il terminale cercherà nelle cartelle specificate in PATH se è presente un file di programma eseguibile denominato in quel modo.

> [!info] Con il comando help variables vedo nome e significato di alcune delle shell variables

| Comando | Descrizione                                          | Tipologia |
| ------- | ---------------------------------------------------- | --------- |
| Set     | Visualizza, imposta o annulla le opzioni della shell | Interno   |
| Env     | Esegue un programma in un ambiente modificato        | Esterno   |
> [!info] 
> - Sudo è il comando che mi permette di richiedere i permessi da utente root.
> - I software vengono distribuiti in Linux sotto forma di pacchetti. 
> - cd senza argomenti mi riporta alla home directory
### Attributi del file/cartella (bash)
Quando eseguo il comando *ls -l*, nella prima colonna, compare un carattere che mi identifica il tipo di file:

| Simbolo | Tipologia          |
| ------- | ------------------ |
| `-`     | File "regolare"    |
| `d`     | Cartella           |
| `b`     | Device a blocchi   |
| `c`     | Device a caratteri |
| `l`     | Link simbolico     |
| `p`     | Pipe               |
| `s`     | Socket             |
> [!example] Esempio: `-rw-r--r--` (trattino iniziale = file regolare), `drwxr-xr-x` (d = directory)

Su Linux ogni file ha un proprietario (utente) e un gruppo proprietario. Quest'ultimo è una raccolta di utenti che condividono un insieme di permessi sul file (un utente appartenente ad un gruppo proprietario del file riceve i premessi per il "group" quando si accede al file stesso). Questo è utile per condividere cartelle/progetti tra più persone senza dare permessi a tutti gli utenti del sistema. Se non si è ne il proprietario ne si appartiene al gruppo proprietario si fa parte degli "altri". Questa distinzione tra ruoli mi serve per distinguere i permessi che ognuno può avere.

Ogni file/cartella ha tre set di permessi, nell'ordine **user (u)**, **group (g)**, **others (o)**. Per ciascuno ci sono tre bit:
- `r` = lettura (read) — per file: leggere il contenuto; per cartelle: vedere i nomi dei file.
- `w` = scrittura (write) — per file: modificare; per cartelle: creare/eliminare file.
- `x` = esecuzione (execute) — per file: eseguire; per cartelle: entrare (`cd`) / attraversare.

> [!info] Questa differenza file vs cartella è importante: ad esempio per entrare in una cartella serve `x`, mentre per elencarne il contenuto serve `r`

Per cambiare proprietario e permessi:

| Comando | Descrizione                           |
| ------- | ------------------------------------- |
| Chown   | Modifica utente e gruppi proprietari  |
| Chmod   | Cambia i permessi di file e directory |
**Parametri Chmod**:
 - *-R* = esegue il comando ricorsivamente
 - *mode* = nuova maschera dei permessi. Questa non è una parola chiave ma serve a capire come impostare i permessi e può essere scritto in 2 modi diversi
   1. *mode simbolico* = chmod (chi) (operatore) (permessi file). Si riferiscono a: **chi può cambiare i permessi** (u = user, g = group, o = others, a all), gli **operatori** (+ aggiunge, - toglie, = imposta esattamente) e ai **permessi** (r = read, w = write, x = execute).

> [!example] 
> - chmod u+x file.sh => Aggiungi il permesso di esecuzione al proprietario
> - chmod g-w file.txt => Togli la scrittura al gruppo
> - chmod o=r file.txt => Gli "others" avranno solo lettura

  2. *mode numerico (ottale)* = qui mode è un numero di 3 cifre. Ogni cifra rappresenta user, group e others e ogni cifra è la somma dei permessi (r = 4, w = 2,   x = 1). Ogni cifra viene calcolata nel seguente modo: si considera un 1 se il permesso è garantito, 0 altrimenti. Si scrive la sequenza di 1 e 0 rispettivamente per i permessi di lettura, scrittura ed esecuzione. L’equivalente ottale del numero binario ottenuto è la cifra che identifica il set di permessi. Le tre cifre rappresentano i permessi associati al proprietario, quelli associati agli utenti appartenenti al gruppo proprietario del file e quelli per tutti gli altri utenti.

> [!example] rw- r-- ---: (lettura scrittura al proprietario, lettura al gruppo, niente agli altri utenti): 110 100 000 = 640

> [!info]
> - chmod +w miofile.c aggiunge il permesso di scrittura per tutti gli utenti
> - chmod 664 miofile.c assegna la maschera rw-rw-r--

 - *pathname* = oggetto a cui applicare la nuova maschera dei permessi

**Parametri Chown:**
- -R = esegue il comando ricorsivamente (anche sottocartelle e file in esse contenuti)
- owner = nuovo utente proprietario 
- group = nuovo gruppo proprietario 
- pathname = oggetti interessati dalla modifica

> [!info] esiste anche il comando chgrp che cambia solo il gruppo: chgrp [-R] group pathname...

> [!example] 
> - cambia solo owner: chown utente file
> - cambia solo gruppo: chown :gruppo file
> - cambia entrambi: chown utente:gruppo file
### System call, wrappers e accesso ai file
<u>Le system calls sono l'interfaccia per il programmatore ai servizi del sistema operativo (chiamata diretta al sistema)</u>. Vi si accede generalmente attraverso un API/linguaggi che la supportino. Ce ne sono vari tipi (accesso ai file, gestione di processi, IPC, semafori,...) e richiedono un "cambio di contesto": da *kernel mode* a *user mode* e viceversa.
- **Kernel mode** = accesso senza limitazioni alle risorse HW/SW del sistema (memoria, dispositivi I/O, registri, codice del SO,...). Un errore in kernel mode può provocare un arresto del sistema
- **User mode** = accesso limitato alle risorse di sistema (es. Memoria). Un errore in user mode può provocare il blocco del solo processo utente

Quando si passa da un mode ad un altro il SO cambia il contesto di esecuzione del processore. Ciò implica un certo consumo di risorse di sistema (tempo di CPU, salvataggio e ripristino del processo e del contesto precedente,...)

**Funzioni wrappers: Glibc**
> [!info] Glibc = "GNU C library" è la libreria C standard che di solito si usa in Linux (contiene le funzioni standard richiamabili da C)

Le system calls sono invocate non direttamente ma tramite funzioni *wrappers* della glibc:
- Le *wrappers* hanno di solito lo stesso nome della corrispondente *system call* sottostante
- Ove possibile è  preferibile chiamare le funzioni *wrappers* piuttosto che direttamente la *system call*. Questo perché spesso si limitano a fare poco altro in più rispetto alle *system call*, ma in altri casi compiono anche altre operazioni come controllo degli argomenti inseriti dal programmatore,... 

**Accesso ai files:** Per l'accesso ai  file ci sono due possibilità:
- System calls come open, write,... (o meglio i corrispondenti wrappers)
- Standard I/O library come fopen, fwrite,... 

**System calls per i file:** 

| Nome                  | Cosa fa                                                                        |
| --------------------- | ------------------------------------------------------------------------------ |
| `open`                | Apre un file in lettura e/o scrittura o crea un nuovo file                     |
| `close`               | Chiude un file precedentemente aperto                                          |
| `read`                | Legge da un file                                                               |
| `write`               | Scrive su un file                                                              |
| `lseek`               | Sposta il puntatore di lettura/scrittura ad una determinata posizione nel file |
| `unlink`              | Rimuove una entry di un file nel filesystem                                    |
| `stat (fstat, Istat)` | Fornisce informazioni sugli attributi di un file                               |
| `dup, dup2, dup3`     | Duplicano un file descriptor                                                   |
**System calls per le directory:**

| Nome    | Cosa fa                           |
| ------- | --------------------------------- |
| `mkdir` | Crea una directory                |
| `rmdir` | Cancella una directory (se vuota) |
| `chdir` | Cambia la directory corrente      |
**Open:** ` #include <unistd.h> int open(const char *pathname, int flags, mode_t mode)`
- Apre il file specificato da *pathname*
- *flags* deve includere uno tra `O_RDONLY, O_WRONLY o O_RDWR`
- *mode* specifica i permessi da applicare al file in caso di creazione (permessi che vengono stabiliti anche in base a `UMASK`). La creazione può essere richiesta specificando `O_CREAT` tra i flags, altrimenti mode viene ignorato e può essere omesso

| Altri flags (separati da \|) | Funzionamento                                                                              |
| ---------------------------- | ------------------------------------------------------------------------------------------ |
| `O_CREAT`                    | Se il file no esiste lo crea                                                               |
| `O_APPEND`                   | Il file è aperto in accodo per la scrittura                                                |
| `O_TRUNC`                    | Se il file esiste il contenuto viene azzerato                                              |
| `O_EXCL`                     | Se usato insieme a O_CREAT e il file esiste, open () fallisce ed errno è settata ad EEXIST |
Questa system call ha due valori di ritorno possibili:
 1. Un intero positivo che è il *file descriptor* da usare per le operazioni successive
 2.  -1 in caso di errore che viene settato da errno

> [!info] 
> - il file offset è la posizione corrente dentro al file, alla quale il prossimo read o write avverrà (è simile ad un cursore dentro al file)
> - di default il puntatore (file offset) nel file è posizionato all'inizio (a meno che non si specifichi O_APPEND)
> - i file aperti andrebbero chiusi prima di terminare il programma con `#include int close(int fd)`, dove fd è il file descriptor ritornato da `open()` 
> - l'apertura di un file crea una nuova entry nella tabella di sistema dei file aperti. Quando un file descriptor viene duplicato (con `dup2` o simili), questo si riferisce alla stessa *open file description* del descrittore del file originale $\implies$  i due file descriptors condividono lo stesso *file offset* e le stesse *file status flags*
> - su Linux si può usare l'operazione kcmp (2) KCMP_FILE per testare se due file descriptors si riferiscono alla stessa open file description
 
> [!tip] 
> - Errno = è una libreria di C (errno.h) che mi dice precisamente quale sia l'errore in corso 
> - con `strace -c ls` posso vedere le system call attive e alcuni dati che le riguardano: percentuale del tempo totale con cui il kernel è stato dentro alla syscall, tempo totale speso in quel tipo di syscall, media per ogni call, numero di volte che la syscall è stata eseguita, quante syscall hanno ritornato un errore, nome della syscall
> - **buffering** = in generale il buffer è una certa quantità di memoria che deve essere riempita prima che accada una qualche chiamata o azione specifica (ad esempio se scrivo una stringa questa potrebbe non essere visualizzata subito perché la memoria buffer non è ancora piena). Esistono 3 tipi di buffer: interattivo (buffer si svuota quando appare `\n`), di pipe o di file (buffer si svuota solo quando la memoria predisposta è piena) e l'assenza di buffer (ogni carattere fa si che la memoria buffer si svuoti subito)

> [!info] fd = è un intero piccolo che il  SO usa per identificare un file aperto dentro ad un processo (non mi riferisco più a file. txt, ma al valore assunto da fd)
> - 0 = stdin
> - 1 = stdout
> - 2 = stderr
> - 3 = primo file aperto dal programma

| Comando        | Contenuto                                              |
| -------------- | ------------------------------------------------------ |
| `man syscalls` | Linux system calls                                     |
| `man 7 libc`   | Overview of standard C libraries on Linux              |
| `strace -c ls` | Summary dell'uso delle syscalls                        |
| `man 2 open`   | Mostra la pagina del manuale per la system call *open* |
| `man 3 stdio`  | Funzioni I/O                                           |

**write:** `ssize_t write(int fd, const void *buf, size_t count);` <u>Scrive i primi #count (numero di byte da scrivere) byte contenuti in buf (puntatore alla memoria che contiene i dati da scrivere) nel file aperto identificato da fd </u>(file descriptor ritornato da open(2)). Ritorna il numero di byte scritti (che può anche essere inferiore a count se, per esempio, il device di scrittura si è riempito) o -1 in caso di errore (se errno settata). 

**read:** `ssize_t read(int fd, void *buf, size_t count);` <u>Legge #count (numero massimo di byte da leggere) byte dal file identificato da fd (file descriptor ritornato da open(2)) e li memorizza in buf (buffer dove mette i dati letti)</u>. Ritorna il numero di byte letti (che può anche essere inferiore a count se c'erano a disposizione nel device di lettura meno byte) o -1 in caso di errore (errno settata). Il valore di ritorno 0 significa EOF (end of file).

> [!info] le operazioni di lettura/scrittura sono eseguite a partire dalla posizione corrente (puntatore, cursore o file offset) nel file

**lseek:** `off_t lseek(int fd, off_t offset, int whence);` Riposiziona il puntatore nel file. A seconda dei seguenti valori di whence il puntatore è posizionato: 
- SEEK_SET: a offset byte dall'inizio 
- SEEK_CUR: a offset byte dalla posizione corrente 
- SEEK_END: a offset byte dopo la fine del file 

Ritorna la nuova posizione (distanza in byte) da inizio file o -1 in caso di errore (setta errno)

**dup, dup 2**: `int dup(int oldfd);` `int dup2(int oldfd, int newfd);` 
- dup() torna una copia di oldfd usando il primo file descriptor libero 
- dup2() torna una copia di oldfd ma usa newfd come file descriptor copia. 

N.B.: Se c'era già un file aperto con descrittore pari a newfd viene chiuso 
- dopo la duplicazione i due file descriptor possono essere usati indifferentemente, ma si riferiscono alla stessa entry nella tabella di sistema che elenca i file aperti e quindi, per esempio, condividono la posizione corrente nel file (se viene modificata tramite uno dei due la modifica vale anche per l'altro) 
- Ritornano il nuovo file descriptor copia, o -1 in caso di errore (errno settata)

```C
#include <fcntl.h> //serve per le costanti che usa open() come O_CREAT, O_RDWR
#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <string.h>
#include <unistd.h>  //syscalls come write() e close()

int main(int argc, char *argv[]) //sono gli unici due parametri che può prendere la funzione main: argc è //il numero di argomenti scritti nella linea di comando del terminale, argv è il puntatore ad un array di //puntatori a ciascun elemento scritto nella linea di comando
{
	int iRet = EXIT_SUCCESS, fd, iReadRet; //iRet è il valore che ritorna il programma e fd è il file 			//descriptor
	char sBuf[] = "Ciao!\n"; //buffer che verrà scritto nel file
	char sReadBuf[128];
	
	if(argc < 2) // se non viene passato il nome del file (argv[1] è il nome del file ed è il secondo 			//parametro passato)
	{
	 printf("Dammi il nome del file");
	 return EXIT_FAILURE;
	}
	
	fd = open(argv[1], O_RDWR | O_CREAT | O_TRUNC, S_IRUSR | S_IWUSR); //chiamata al kernel con la syscall 		//open()
	//le flag che compongono il secondo parametro mi dicono che viene aperto il file in scrittura e 			//lettura, viene creato il file se non esiste e viene svuotato il file se esiste
	//S_IRUSR = read permission for user
	//S_IWUSR = write permission for users
	if (fd >= 0)
	{
		if (write(fd, sBuf, strlen(sBuf)) == strlen(sBuf)) //chiamata al kernel con la syscall write()
		//write prende come argomenti il file descriptor, l'indirizzo in memoria del buffer e il numero di 			//bytes da scrivere
		{
			printf("Ho scritto nel file %s, fd = %d\n", argv[1], fd);
			if (lseek(fd, 0, SEEK_SET) != -1) //sposto il file offset all'inizio del file (perché dopo la 				//lettura sono alla fine dello stesso) e mi sposto di zero posizioni
			{
				if ((iReadRet = read(fd, sReadBuf, sizeof(sReadBuf) - 1)) != -1) //sReadBuf è l'indirizzo 					//al primo byte del buffer
				{
				 char sWriteBuf[] = "Questa proviene da write()\n";
				 int iCount = strlen(sWriteBuf);
				 sReadBuf[iReadRet] = '\0'; //terminazione della stringa (read non la aggiunge di default)
				 dup2(fd,1); // duplica fd su STDOUT_FILENO, ora scrivere su STDOUT_FILENO scrive su fd
				 //siccome newfd = 1, ora stdout (nello specifico printf) non stampa più a terminale ma 					 //direttamente nel file
				 printf("Ho riletto per %d byte: %s", iReadRet, sReadBuf);//printf scrive sul descrittore 1
				 /* write(STDOUT_FILENO, sWriteBuf, iCount); */
				 }
				 else
				 {
				 perror("Errore in read()");
				 iRet = EXIT_FAILURE;
				 }
			}
			else
			{
			perror("Errore in lseek()");
			iRet = EXIT_FAILURE;
			}
			close(fd); //chiiamata al kernel con la syscall close() che rilascia il fd e flusha il kernel buffer
		}
		else //write error handling
		{
			printf("errno = %d\n", errno);
			perror("Errore in write()"); // stampa un messaggio leggibile
			iRet = EXIT_FAILURE;
		}
	}
	else //open error handling
	{
		printf("errno = %d\n", errno); 
		perror("Errore in open()"); /* stampa la stringa a video: descrizione errore*/
		iRet = EXIT_FAILURE;
	}
	return iRet; //valori che possono essere ritornati: 0 per il successo e -1 per il fallimento
}
//se eseguo il programma con strace ./nomeProgramma vedo tutto il percorso delle varie syscalls
```

### Standard I/O library
È un insieme di funzioni per accesso ai file, agli *streams* standard, etc che astraggono le operazioni di input-output sotto forma di streams. Le operazioni sono bufferizzate in spazio utente e ciò riduce il numero di chiamate alle sottostanti syscall e i relativi cambi di contesto ottimizzando le prestazioni. Forniscono un'interfaccia verso il SO di più alto livello rispetto alle syscall

| Syscall (no bufferizzazione nello spazio utante)         | StdIO library (bufferizzazione nello spazio utente)                                                           |
| -------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------- |
| `open`                                                   | `fopen`                                                                                                       |
| `close`                                                  | `fclose`                                                                                                      |
| `read/write`                                             | `getchar/putchar, getc/putc, fgetc/fputc, fread/fwrite, gets/puts, fgets/fputs, scanf/printf, fscanf/fprintf` |
| `lseek`                                                  | `fseek`                                                                                                       |
| Richiedono cambio contesto user-kernel ad ogni carattere | Tipicamente più efficaci e più veloci                                                                         |
### Bash: standard stream e redirezioni
**command-line completion:** *bash* ha la funzione di completamento delle linee di comando premendo il tasto `TAB`. Premendo una volta il tasto *TAB* ottengo:
1. se c'è una sola possibilità il comando viene completato 
2. altrimenti il comando viene eventualmente completato fin dove è possibile farlo in modo univoco e viene emesso un segnale sonoro

Nel secondo caso, premendo una seconda volta, viene mostrata la lista delle scelte possibili

**Gestione della history:** *bash* tiene nota degli ultimi *n* comandi dati alla shell. *n* dipende dalle variabili `HISTFILESIZE` e `HISTSIZE` (*man bash*)

| Comando (interno) | Descrizione                                                                                                                                |
| ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------ |
| `history`         | visualizza gli ultimi n comandi eseguiti (il numero di comandi visualizzato e memorizzato dipende dalle variabili HISTSIZE e HISTFILESIZE) |
| `!a`              | esegue l’ultimo comando nella history che inizia per ‘a’                                                                                   |
| `!68`             | esegue il comando numero 68 della history                                                                                                  |
| `!!`              | ripete l’ultimo comando eseguito                                                                                                           |
**Bash: Metacaratteri:** si tratta di caratteri che separano le parole. Un metacarattere può essere: spazio, tab, newline, `|`, `&`,`;`,`(`, `)`, `<`, `>` (da *Bash Reference Manual*). Questi caratteri hanno un significato per la shell, come anche i caratteri `*`, `?`, `Ctrl + d (EOF)`, `Ctrl + c`, `'`,…

> [!info] i caratteri `*` e `?` sono **metacaratteri** (wildcard) utilizzati per l'espansione dei nomi di file (globbing) e per il pattern matching nelle stringhe
> - `?` rappresenta esattamente un singolo carattere qualsiasi, per cui il pattern `lettera?` corrisponde solo a nomi come `lettera1` o `lettera2`, ma non a `lettera10` o `lettera` 
> - `*` rappresenta qualsiasi sequenza di caratteri (inclusa una sequenza vuota), quindi `*.txt` corrisponde a tutti i file che terminano con `.txt`, mentre `*` corrisponde a tutti i file nella directory corrente

> [!tip]
> Se si desidera usare questi caratteri come testo normale, devono essere racchiusi tra **apici** (quoting) per evitare che la Bash li interpreti.

**Bash: standard streams:** In un sistema Unix-like, ogni processo all'avvio ha tre flussi di dati (file descriptor) predefiniti aperti. Questi flussi permettono al comando di interagire con l'ambiente esterno (tastiera, terminale, file) tramite anche la redirezione e il comando pipe (`|`).
* **0 (stdin):** Standard Input. Tipicamente collegato alla **tastiera**, è utilizzato come sorgente da cui prendere i dati da elaborare.
* **1 (stdout):** Standard Output. Tipicamente collegato al **terminale** (monitor), è utilizzato per mostrare i risultati della computazione.
* **2 (stderr):** Standard Error. Tipicamente collegato al **terminale**, è utilizzato per comunicare informazioni addizionali sullo stato della computazione come i messaggi di errore.
* **3+ (stdfile):** Descrittori di file aggiuntivi definiti dall'utente o dal sistema.

> [!NOTE] Nota del Prof
> Se esegui un comando come `sort` senza alcun parametro, la Bash "apre lo 0": invece di leggere da un file, resta in attesa di input direttamente dalla tastiera finché non viene inviato un segnale di EOF (Ctrl+D).

**Bash: redirezione dell’input:** <u>La redirezione dell'input in Bash permette di inviare il contenuto di un file (o altro input) a un comando</u>, sostituendo la tastiera standard (stdin, descrittore file 0) con l'operatore `<`: `comando [n] < filename` (<u>la shell apre il file, lo collega a stdin e il comando legge da stdin</u>). 
- Il descrittore *n* è associato a *filename* che viene aperto in lettura
- Se *n* è omesso, lo *standard input* (*n* = 0) è associato al file, cioè il comando legge l’input dal file

> [!example] 
> - `sort < /etc/hosts`: Collega il descrittore 0 (stdin) al file specificato. Il comando `sort` leggerà il contenuto di `/etc/hosts` come se fosse digitato dall'utente. 
> - `cat /percorso-file` . Se non viene specificato un file, `cat` apre il descrittore 0 (la tastiera nel terminale). 

> [!IMPORTANT] Differenza Tecnica
> C'è una sottile differenza tra l'uso di un comando con `<` e senza:
> * **Senza `<` :** Il comando gestisce l'apertura del file internamente (spesso usando il descrittore **3**).
> * **Con `<` :** È la Shell (Bash) ad aprire il file e collegarlo allo **0** prima di lanciare il comando.

> [!info] Posso avere anche le modalità seguenti:
> - **Here Document (`<<`):** permette di passare blocchi di testo come input
> - **Here String (`<<<`):** passa una singola stringa come input
>> [!example] `grep "cerca" <<< "$variabile"`

**Bash: redirezione dell'Output:** I simboli `>` e `>>` permettono di deviare/salvare l'output di un comando (che normalmente andrebbe a video) verso un file (sovrascrivendolo). `comando [n] > filename` (la shell apre il file in questione in scrittura e collega stdout a quel file)
- Il descrittore *n* è associato a *filename* che viene aperto in scrittura
- Se *n* è omesso lo *standard output* (*n* = 1) è associato al file, cioè il comando scrive sul file l’output
> [!example] `ls -l > lista.txt`: Salva la lista dei file in `lista.txt`, sovrascrivendolo.

**BASH: redirezione di stderr:** `comando 2> filename` (default `n=1`), reindirizza solo lo standard error in un file
> [!example] 
> - `cat /etc/hosts > file.out`: Il contenuto di `/etc/hosts` viene scritto in `file.out` . Se il file esiste, viene sovrascritto.
> - **Redirezione Errori:** `ls cartella_inesistente 2> errori.log` Invia solo i messaggi d'errore al file `errori.log` .

**BASH: redirezione dell'output in append:** `comando >> file` viene aggiunto l’output in coda al file senza sovrascriverlo
> [!example] `echo "nuova riga" >> scan.txt`: Il testo viene aggiunto in coda al file senza cancellare il contenuto precedente.

**BASH: reindirizzare/duplicare uno stream:** È possibile duplicare o spostare i file descriptor per operazioni avanzate con gli operatori della redirezione (<u>si tratta di far puntare due flussi allo stesso oggetto tipo un file</u>).
- `[n]<&x`: duplica il file descriptor in input x, nel senso che n fa riferimento allo stesso device/file cui si riferisce x (il quale è un fd). Se n è omesso viene assunto 0 (stdin). x deve essere un descrittore aperto in input altrimenti viene restituito un errore
- `[n]>&x`: duplica il file descriptor di output x. Se n è omesso viene assunto 1 (stdout). x deve essere un descrittore aperto in output altrimenti viene restituito un errore (semplicemente viene redirezionato l'output al file/device indicato da x)

> [!tip]
> Ogni flusso, a prescindere dalla redirezione, mantiene la sua natura (es.stdin legge o stdout scrive,...), ciò che cambia è da dove legge o dove scrive. Quando si scrive `0<&3` ("stdin punta ciò che punta 3") significa che stdin usa lo stesso oggetto (file aperto) che usa 3

> [!example] `3< /percorso-file comando 0<&3`: 
>  Viene aperto il file sul descrittore 3, poi il comando (es. `cat` ) viene istruito a leggere dallo 0 che è stato duplicato dal 3 (simile alla system call `dup2` ).

> [!NOTE] Nota del Prof
> La funzione di sistema `read` torna sempre il numero di byte effettivamente letti. Per analizzare queste operazioni a basso livello, si può usare il comando `strace` , che mostra tutte le **system call** eseguite da un comando.

**BASH: redirezione simultanea di stdout e stderr:** 
- Esistono due possibili sintassi: `comando &> filename` oppure  `comando >& filename`. La prima delle due è preferibile e semanticamente equivale a`comando > filename 2>&1`. 
- Analogamente si realizza l’append:`comando &>> filename`, che semanticamente equivale a `comando >> filename 2>&1`

**Bash:Pipeline:** `cmd1 | cmd2 | ... | cmdN`. Posso scrivere ed eseguire vari comandi separati solo dal carattere `|` (pipe). La pipe collega direttamente lo Standard Output di un comando allo Standard Input di quello successivo. Al posto del carattere `|` posso anche utilizzare solo `|&` con cui stdout e stderr del comando precedente vengono entrambi connessi all’input del comando successivo. Ogni comando è eseguito in un processo differente (*subshell*)

> [!example] `sort /etc/passwd | less`: Il file `/etc/passwd` viene ordinato da `sort` , l'output prodotto non viene stampato a video ma passato a `less` , che permette di consultarlo pagina per pagina.

| Comando                                     | Flusso                                                                                                                                 | Comportamento                                                                                                                                                                                                                                                                                                                          |
| ------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `cat`                                       | Tastiera $\rightarrow$ cat $\rightarrow$ tastiera                                                                                      | scrivo e lui stampa fino al `Ctrl + d`. `cat` usa il flusso stdin e scrive su stdout                                                                                                                                                                                                                                                   |
| `cat /etc/hosts`                            | /etc/hosts $\rightarrow$ cat $\rightarrow$ stdout (terminale)                                                                          | `cat` apre `/etc/hosts` e ne stampa il contenuto su stdout (legge il file passato direttamente come argomento)                                                                                                                                                                                                                         |
| `cat < /etc/hosts`                          | /etc/hosts $\rightarrow$ stdin $\rightarrow$ cat $\rightarrow$ terminale                                                               | la shell apre il file, lo collega a stdin (fd 0) e `cat` lo legge da stdin                                                                                                                                                                                                                                                             |
| `cat 3< /etc/hosts`                         | /etc/hosts $\rightarrow$ fd 3 (inutilizzato) $\rightarrow$ stdin (terminale) $\rightarrow$ cat $\rightarrow$ terminale                 | La shell con `3< /etc/hosts` apre il file in lettura e lo collega a fd 3, poi si esegue cat. Siccome `cat` usa SOLO fd 0 e fd 1, "ignora" fd 3 e stdin resta la tastiera (quindi il comando si comporta come se ci fosse solo `cat`)                                                                                                   |
| `cat 3< /etc/hosts 0<&3`                    | /etc/hosts $\rightarrow$ fd 3 $\rightarrow$ stdin $\rightarrow$ cat $\rightarrow$ terminale                                            | `3< /etc/hosts` apre il file su fd 3 e `0<&3` <u>duplica</u> fd 3 su stdin. Ora stdin punta al file e `cat` legge da li                                                                                                                                                                                                                |
| `cat 4> /tmp/my.txt`                        | stdin (tastiera) $\rightarrow$ cat $\rightarrow$ terminale  <br>fd 4 $\rightarrow$ file (inutilizzato)                                 | la shell apre il file in scrittura su fd 4 (`cat` non usa fd 4). `cat` legge da stdin e stampa su stdout. Il file viene creato (vuoto) ma non scritto                                                                                                                                                                                  |
| `cat 4> /tmp/my.txt /etc/hosts`             | /etc/hosts $\rightarrow$cat $\rightarrow$ terminale  <br>fd 4 $\rightarrow$ file vuoto                                                 | `4> /tmp/my.txt` apre fd 4 (inutile per cat) e `cat /etc/hosts` legge il file e stampa su stdout (quindi il file.txt resta vuoto e il contenuto dell'argomento viene stampato a terminale)                                                                                                                                             |
| `cat 4> /tmp/my.txt /etc/hosts 1>&4`        | /etc/hosts $\rightarrow$ cat $\rightarrow$ stdout (1) $\rightarrow$ fd 4 $\rightarrow$ file                                            | `4> /tmp/my.txt` apre fd 4 e `1>&4` duplica il flusso stdout in fd 4 (quindi l'output di `cat /etc/hosts` va nel file). In sostanza non avrò output a schermo ma `/tmp/my.txt` conterrà `/etc/hosts`                                                                                                                                   |
| `cat 3< /etc/hosts 4>/tmp/my.txt 1>&4 0<&3` | /etc/hosts $\rightarrow$ fd 3 $\rightarrow$ stdin $\rightarrow$ cat $\rightarrow$ stdout $\rightarrow$ fd 4 $\rightarrow$ /tmp/my. Txt | `3< /etc/hosts` apre il file in lettura su fd 3, `4> /tmp/my.txt` apre il file in scrittura su fd 4, `0<&3` duplica stdin su fd 3 e `1>&4` duplica stdout su fd 4. In sostanza il comando legge (stdin) da `/etc/hosts` e scrive (stdout) su `/tmp/my.txt` (è come fare una copia del file). Equivale a `cat /etc/hosts > /tmp/my.txt` |
> [!info] Il comando `cat` legge sempre e solo da stdin (fd 0) e scrive sempre e solo su stdout (fd 1)

> [!tip] Nota bene:
> - un comando come `<` fa si che la shell apra il file usando il pathname, ottenga così il pathname e lo colleghi a stdin (fd 0). Quindi si usa un file esterno che viene aperto nel momento dell'esecuzione del comando
> - un comando come `<&` implica che la shell NON debba aprire un file, ma usare un file descriptor già aperto (lo duplica). Pertanto si lavora solo con connessioni già esistenti

> [!warning] 
> - `<` apre il file in lettura
> - `>` apre il file in scrittura
> - `<&` duplica (solitamente) stdin (se c'è 0 prima) sul file descriptor passato
> - `>&` duplica (solitamente) stdout (se c'è 1 prima) sul file descriptor passato

**Bash- comandi informativi:**

| Comando    | Comportamento                                            |
| ---------- | -------------------------------------------------------- |
| `id`       | Visualizza gli ID di utente e gruppo                     |
| `hostname` | visualizza/imposta l'hostname della macchina             |
| `uname`    | visualizza informazioni sul sistema                      |
| `which`    | visualizza il path dell'eseguibile che esegue un comando |
| `date`     | visualizza/imposta data e ora correnti                   |
| `who, w`   | visualizzano chi è loggato nel sistema                   |
| `tty`      | visualizza il nome del terminale                         |
| `df`       | visualizza l’occupazione dello spazio disco              |
**Bash- trattamento testi:**

| Comando        | Comportamento                                                                       |
| -------------- | ----------------------------------------------------------------------------------- |
| `grep (-rvec)` | visualizza linee che soddisfano criteri (patterns) di ricerca                       |
| `awk`          | implementa un linguaggio di programmazione per il «pattern scanning and processing» |
| `tr`           | sostituisce, comprime e/o cancella caratteri                                        |
| `uniq`         | visualizza o sostituisce linee ripetute                                             |
| `cut`          | seleziona parti di linee di un file (es.: a campi delimitati)                       |
| `sed(-ei)`     | stream editor per filtrare e modificare un testo                                    |
| `wc(-l)`       | conta linee, parole e caratteri                                                     |
| `more`         | visualizza il contenuto di un file una pagina alla volta                            |
| `less`         | permette di navigare in un file una pagina alla volta (ed anche altro...)           |
| `diff`         | confronta il contenuto di due file per linee                                        |
**Bash- gestione processi:**

| Comando       | Comportamento                                              |
| ------------- | ---------------------------------------------------------- |
| `ps (axuj)`   | visualizza i processi attivi                               |
| `top`         | visualizza in modo interattivo processi e risorse occupate |
| `nice`        | esegue un comando modificandone la priorità                |
| `renice`      | modifica la priorità di un processo                        |
| `kill`        | invia un signal ad un processo                             |
| `strace (-p)` | monitora system calls e signal di un processo              |
**Bash- altri comandi:**

| Comando     | Comportamento                                                                                                                                                         |
| ----------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `find`      | cerca file che corrispondono a determinati criteri                                                                                                                    |
| `uptime`    | risponde a: "da quanto la macchina è accesa?"                                                                                                                         |
| `last`      | riporta le ultime login                                                                                                                                               |
| `lastlog`   | riporta le ultime login di tutti o di un singolo utente                                                                                                               |
| `xargs`     | da man: "build and execute command lines from standard input"                                                                                                         |
| `su(-,-c)`  | permette di "agire" come altro utente                                                                                                                                 |
| `echo(-en)` | visualizza una linea di testo (interno)                                                                                                                               |
| `help`      | fornisce la lista dei comandi built-in della shell (interno)                                                                                                          |
| `tee`       | "scrive" il suo stdin su stdout e uno o più file                                                                                                                      |
| `basename`  | visualizza in output il solo nome file senza le precedenti directory da un pathname. Se richiesto rimuove il suffisso. Es.: basename /home/user 1 visualizza "user 1" |
| `dirname`   | visualizza le directory che precedono il nome del file in un pathname Es.: dirname /home/user 1 visualizza "/home"                                                    |
| `sleep n`   | attende n secondi senza fare nulla                                                                                                                                    |
**Comando ps:** è un'abbreviazione per *precess status*, visualizza informazioni sui processi attivi nel sistema. Ci sono 3 tipi di opzioni (opzioni di diverso tipo possono essere mescolate tra di loro ma ci possono essere dei conflitti):
1. *Unix style:* precedute da un trattino
2. *BSD style:* non precedute da un trattino
3. *GNU long options:* precedute da 2 trattini

| Comando                     | Comportamento                                                                                                            |
| --------------------------- | ------------------------------------------------------------------------------------------------------------------------ |
| `ps -e/-ef/-eF/-ely`        | Per vedere ogni processo nel sistema usando la sintassi standard                                                         |
| `ps ax/axu`                 | Per vedere ogni processo nel sistema usando la sintassi BSD                                                              |
| `ps -ejH/axjf`              | Per stampare un albero dei processi                                                                                      |
| `ps -eLf/axms`              | Per avere informazioni sui thread                                                                                        |
| `ps`                        | Lista i processi dell'utente che lo esegue ed associati allo stesso terminale cui è associata la shell da cui è invocato |
| `ps ax`                     | Lista tutti i processi                                                                                                   |
| `ps axf`                    | Visualizzazione a foresta (gerarchica)                                                                                   |
| `ps –U user1`               | Visualizza i processi dell'utente *user1*                                                                                |
| `ps ax –L`                  | Lista tutti i processi e i thread                                                                                        |
| `ps –C name1,name2,…,nameN` | Lista i processi il cui nome è specificato tra quelli che seguono -C                                                     |

> [!info] Useremo anche:
> - `ps ax | grep –e <string1> –e <string2>` per selezionare delle righe dall'output di ps, e precisamente quelle che contengono una delle stringhe string1 o string2 o entrambe 
> - watch `–n3 'ps ax | grep –e <string1> –e <string2>'` per eseguire lo stesso comando sopra descritto ogni 3 secondi fin quando il comando watch non viene interrotto (per esempio premendo Ctrl+c)
### System calls per i processi
**fork():** `include <unistd.h> pid_t fork();`

È una syscall che non vuole argomenti. Questo metodo ritorna -1 difficilmente (il *child* non viene creato e errno viene settata), altrimenti:
- Crea un nuovo processo duplicando il chiamante, tuttavia alcuni attributi non sono duplicati: PID, memory locks, segnali pendenti, timer, etc,
- Il nuovo processo  è detto *child*, il processo chiamante è detto *parent*
- Lo spazio di indirizzamento del nuovo processo è un duplicato di quello del *parent* (sono separati, ognuno ha il proprio)

> [!info] 
> - fork() nel *child* torna 0, mentre nel *parent* torna il valore del PID del *child* (valore > 0). In questo modo si può differenziare il codice tra i due
> - con `fprintf()` posso decidere su quale descrittore andare a stampare 

I due processi (child e parent) sono eseguiti concorrentemente. Il child eredita una copia della *tabella dei descrittori* aperti dal padre ma ogni copia di un file descriptor si riferisce alla stessa voce nella tabella dei file aperti (di sistema) cui si riferisce il descrittore originale. Originale e copia condividono:
- I flag di stato del file
- Il puntatore alla locazione corrente di ogni file
- Etc. (vedi *man open (2) e man fork (2))

```C
#include <sys/wait.h>
#include <stdio.h>
#include <unistd.h>

int main()
{
	pid_t pid;	
	PID= fork();/* fork a child process */
	if (PID< 0) { /* error occurred */
		fprintf(stderr, "Fork Failed");
		return 1;
	}
	else if (PID== 0) { /* child process */
		printf("Child: my PID= %d, parent PID= %d\n", getpid(), getppid());//execlp("/bin/ls","ls",NULL);
		//sleep(10); //mette il processo in attesa passiva (no consumo risorse) per 10 secondi
		_exit(0); //termina il processo figlio senza eseguire le operazioni di pulizia del processo (come la chiusura dei file //aperti)
	}
		else { /* parent process */
		/* parent will wait for the child to complete */
		wait(NULL);
		printf("Parent: my PID= %d, parent PID= %d. PID del child %d\n",getpid(),getppid(),PID);
		//sleep(10);
	}
	sleep(10);
	return 0;
}
```

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260605185108.png" width="300">
</div>

**PID e PPID:** `include <unistd.h> pid_t getpid(void); pid_t getppid(void)`

- getpid() torna il process ID
- getppid() torna il parent process ID 

Queste due funzioni hanno sempre successo

> [!info] Esiste anche `getpgrp()` che ritorna il *process group ID*

**wait () e waitPID():** `#include <sys/wait.h> pid_t wait(int *wstatus); pid_t waitpid(pid_t pid, int *wstatus, int options)`

- `wait()` sospende il chiamante finché un *child* termina
- `waitpid()` sospende il chiamante finché un *child* cambia stato (terminato; stoppato da un *signal*; ha ripreso dopo un *signal*), dipende da *options*

Permettono di rilasciare le risorse assegnate al *child*. Se un *child* termina e una *wait* non viene eseguita il *child* rimane nello stato di *zombie* (*man wait(2)*). Questi metodi ritornano il PID del processo il cui stato è cambiato (`waitpid()` può tornare 0 in certe circostanze) oppure -1 in caso di errore. Tra i valori che possiamo dare a PID ci sono i seguenti:
- -1: attende per un cambiamento di stato di uno qualsiasi dei child
- >0: attende per un cambiamento di stato del child il cui process ID è uguale a PID
- Etc

`waitpid(-1, &wstatus, 0)` è equivalente a `wait(&wstatus)`. Se `&wstatus` non è NULL riceve informazioni sul *child*. Queste informazioni possono essere evinte con delle macro:
- WIFEXITED(wstatus) torna true se il child è terminato normalmente (exit() o return) 
- WEXITSTATUS(wstatus) ritorna il codice di uscita del child 
- Etc. (man 2 wait)

> [!info] Il comando `watch` in bash esegue ripetutamente un comando specifico a intervalli regolari (di default ogni 2 secondi), mostrando l'output nel terminale. È fondamentale per monitorare in tempo reale le modifiche all'output di un comando senza doverlo rieseguire manualmente, aggiornando la visualizzazione. Con `-n` imposto l’intervallo di tempo e con `-d` evidenzia le differenze tra l’esecuzione precedente e l’attuale

**System calls per i processi**

| Nome                | Cosa fa                                          |
| ------------------- | ------------------------------------------------ |
| `fork`              | Crea un nuovo processo                           |
| `wait`              | Attende finché uno dei child termina             |
| `waitPID`           | Attende finché un determinato child cambia stato |
| `waitid`            | Maggiore controllo su quali child attendere      |
| `_exit(2), exit(3)` | Terminano il chiamante                           |
| `execve(2)`         | Esegue un programma che rimpiazza il chiamante   |
| `exec*(3)`          | Wrappers per execve(2) (sono 6)                  |
### File Bash
**Backtick:** in BASH le stringhe contenute tra caratteri backtick $(`)$ sono sostituite dall'output del comando che rappresentano  

> [!example] 
> user1@host1:~$ pwd 
> 
> /home/user1 
> 
> user1@host1:~$ ls -ld `pwd`
>  
> drwxr-xr-x 93 user1 group1 12288 mar 17 11:52 /home/user1 
> 
> user1@host1:~$

> [!info] La sintassi `$(cmd)` è equivalente (l'output del comando `cmd` si sostituisce a `$(cmd)`)

**Esecuzione in background:** Esegue un comando senza interazione con l’utente. La shell viene subito «liberata», cioè non devo attendere la fine del comando per impartirne altri. Alla fine del comando si inserisce il carattere `&` (il processo in background continuerà comunque a mandare l’output al terminale). I processi in background sono detti «jobs». Eseguendo un comando in background la shell mostra il PID del processo ed il codice del job. Per visualizzare i job in esecuzione uso il comando `jobs`. Per far tornare nella modalità interattiva (foreground) un job uso il comando `fg [%x]`, dove `%x` è facoltativo e `x` è il codice del job (1,2,...). Per terminare un job uso il comando `kill -9 <pid>` , dove `<pid>` è il PID del processo. Per mandare un processo da foreground in background: 
- Stopparlo con Ctrl+z 
- Mandarlo in background col comando bg

**Variabili:** array di coppie nella forma nome=valore:
- nome: stringa alfanumerica che comincia per lettera 
- valore: stringa di caratteri 

Si dichiarano/assegnano senza spazi intorno all'uguale (per convenzione il nome in maiuscolo) 

> [!example] 
> - VAR1=2
> - VAR1=`pwd`
> - VAR1=$VAR2

Si referenziano con $NOME e si può utilizzare anche la notazione ${NOME} che è utile se può esserci ambiguità

> [!example]  
> user1@host1 :~$ VAR 1="NUMBER " 
> 
> user1@host1 :~$ VAR 1 ONE="scratch" 
> 
> user1@host1 :~$ echo $VAR1ONE 
> 
> scratch
> 
> user1@host1 :~$ echo ${VAR1}ONE 
> 
> NUMBER ONE 
> 
> user1@host1 :~$

> [!info] 
> - nelle costanti stringa delimitate da DOPPI apici $NOME viene interpretato 
> - nelle costanti stringa delimitate da SINGOLI apici $NOME non viene interpretato 
>> [!example] 
>> - `VAR1="33"` echo "$VAR 1" # visualizza '33' 
>> - `echo '$VAR1'` # visualizza '$VAR 1'

**Ambiente (environment):** variabili impostate inizialmente dalla BASH quando viene invocata in base ai file di configurazione, "ereditate" da ogni comando eseguito dalla BASH.
- `NOME=VALORE` per definire una nuova variabile $NOME 
- `export NOME` per renderla disponibile ai processi child 
- `unset NOME` per cancellarla 
- `set` o `env` per visualizzare l'environment. set mostra tutte le variabili e funzioni definite nella shell. set permette anche di settare gli attributi (opzioni) che regolano il funzionamento della BASH

**Variabili d'ambiente predefinite:** 
- `PATH`: contiene una lista di directory, separate da ':', in cui la BASH cerca i comandi da eseguire (se non viene specificato un path per essi) 
- `HOSTNAME`: nome host del sistema 
- `USER`: nome dell'utente 
- `PWD`: directory corrente

```C
// Prima riga = quale interprete deve interpretare i comandi: shabang
#!/bin/bash

RET=0
LABEL="${1}" #non mettere spazi tra le parole, altrimenti da errore 

if test $# -ge 1
then

gcc -o $LABEL $LABEL.c
RET=$? #ciò che gcc ritorna alla shell lo ho in $?

if test $RET -eq 0 #-eq significa equal e si usa per i valori numerici
then

./$LABEL

echo "Pid dello script = $$"
echo "Pid della bash = $BASHPID"

fi //fine if
LABEL="fork" #non mettere spazi tra le parole, altrimenti da errore

gcc -o $LABEL $LABEL.c

RET=$? #ciò che gcc ritorna alla shell ce lo ho in $?

if test $RET -eq 0 #-eq significa equal e si usa per i valori numerici

then

./$LABEL

echo "PIDdello script = $$"

echo "PIDdella bash = $BASHPID"

fi #fine if

exit $RET
else

echo "Usage: `basename $0` nome_file" >&2
RET=1

fi
  
exit $RET
```

> [!attention] Nei file bash (e nella shell in generale), non bisogna mettere spazi attorno al segno di uguale quando si definisce una variabile perché **bash usa gli spazi per separare i comandi dai loro argomenti** (nome variabile verrebbe interpretato come nome di un comando da eseguire e `=` e il `valore` della variabile verrebbero interpretati come argomenti di questo comando)

**Script e il suo contenuto:** uno script è un file di testo che contiene comandi per la shell. Se eseguito la shell esegue i comandi contenuti. Viene usato per automatizzare operazioni senza dover sempre ripetere la scrittura dei comandi che le eseguono

> [!example]
> bash:~$ cat sc1.sh
> 
> cd /etc
> 
> echo "Ora siamo in /etc"
> 
> bash:~$

Per eseguire uno script posso:
- `bash sc1.sh` 
- oppure renderlo eseguibile aggiungendo il permesso di esecuzione x: `chmod u+x sc1.sh` e richiamandolo con `./sc1.sh` (se nella directory corrente) 
- lo script è eseguito in una subshell 
- `source sc1.sh` (o . sc1.sh) che fa un'altra cosa: "legge" lo script e ne esegue i comandi all'interno della shell corrente

La *shabang line* è la prima riga dello script ,indica l'interprete dei comandi contenuti: `#!/bin/bash`. 

Il carattere `#` è usato in generale per inserire commenti nel codice dello script

Quando si parla di variabili, come in qualsiasi linguaggio di programmazione, si tratta di coppie NOME - VALORE. Non sono tipizzate e non ammettono dichiarazione, si definiscono assegnando un valore:es. VAR 1=“Waiting…”, VAR 2=10. Una volta definite si referenziano anteponendo un $ prima del nome: es. echo `$VAR1`, `VAR2=$VAR1`

Possiamo passare ad uno script degli argomenti, per lo script sono variabili nella forma `$n` dove n è un intero da 0 in poi. `$0` contiene il nome dello script e gli altri argomenti (se passati) sono da `$1` a `$n` dove n è il loro numero. Quelli oltre il nono devono essere referenziati nella forma `${n}` (per es.: $10 è sostituito col primo + il carattere "0" e non il decimo) 

> [!info]  
> - il numero di argomenti passati è nella variabile `$#`
> - il comando interno `shift` rinomina i parametri posizionali (argomenti) decrementando l'intero che li identifica. Cioè `$2` diventa `$1`, `$3` diventa `$2` e così via (utile nel parsing degli argomenti)

Le variabili della shell non sono visibili ai processi figli (anche script) o alle sottoshell, a meno che non siano state esportate (con comando `export`) e siano quindi state inserite nell'environment (variabili d'ambiente).

Esistono alcune variabili speciali:
- `$0`: nome dello script 
- `$*`: tutti gli argomenti passati allo script in una sola parola ("$*" equivale a "$1 $2 $3 …") 
- `$@`: come il precedente ma ogni argomento in una parola distinta ("$@" equivale a "$1" "$2" "$3" …) 
- `$#`: numero di argomenti passati
- `$!`: PID dell'ultimo comando eseguito in background 
- `$_`: ultimo argomento passato al comando precedente 
- `$?`: valore di ritorno del comando precedente 
- `$$`: process identifier (PID) della shell, tranne che per una subshell tra `()` per cui torna il PID del chiamante. Si può utilizzare la variabile BASHPID

La shell riconosce certi caratteri come "speciali" e li "interpreta", cioè attribuisce loro un significato diverso dal carattere in sé. Possiamo prevenire questo comportamento col quoting: `\`. La shell ignora il significato "speciale" del carattere preceduto da `\`. Per certi comandi `\` ottiene l'effetto contrario, cioè abilita un significato speciale per il carattere che lo segue. Per es. per avere comportamenti simili al C (inserire un newline in una stringa, etc.)

| Sequenze di escape | Uso per echo e sed                   |
| ------------------ | ------------------------------------ |
| `\n`               | Newline                              |
| `\r`               | Return                               |
| `\t`               | Tab                                  |
| `\v`               | Vertical tab                         |
| `\b`               | Backspace                            |
| `\a`               | Alert (beep or flash)                |
| `\0xx`             | Traduce all'equivalente ottale ACSII |
Inoltre il quoting con singolo apice preserva il significato letterale dei caratteri, mentre quello con doppi apici preserva il significato letterale dei caratteri, eccezion fatta per '$', ''' e '\'. Sostituisce quindi le variabili e i comandi racchiusi da apici obliqui.

> [!info] echo ls stampa ls (cioè ciò che c’è dopo echo)
> 'echo ls' (Alt Gr + ') stampa invece il contenuto di ls

**Creazione processi:execve(2):** `#include int execve(const char *pathname, char *const argv[], char *const envp[]);`
- **Pathname:** path del programma da eseguire. Può essere un eseguibile o uno script (che contenga la shabang)
- **Argv:** array di puntatori a stringa che contiene gli argomenti che verranno passati al nuovo programma. Il primo elemento è il "filename" del programma da eseguire (`argv[0]`). L’ultimo elemento deve essere un puntatore nullo
- **Envp:** array di puntatori a stringhe, nella forma nome=valore, che verranno passare come enviroment al nuovo programma. L’ultimo elemento deve essere un puntatore nullo

**Creazione processi: exec(3):** `#include extern char **environ;`

`int execl(const char *pathname, const char *arg, ... /* (char *) NULL */); `

`int execlp(const char *file, const char *arg, ... /* (char *) NULL */); `

`int execle(const char *pathname, const char *arg, ... /*, (char *) NULL, char *const envp[] */); `

`int execv(const char *pathname, char *const argv[]); `

`int execvp(const char *file, char *const argv[]); 
`
`int execvpe(const char *file, char *const argv[], char *const envp[]);
`
`#include <unistd.h>`

- Rimpiazzano il processo chiamante con un nuovo processo 
- Chiamano la sottostante `execve(2) `
- Il primo argomento è il nome del file da eseguire

Se non c'è errore non ritornano nulla, -1 in caso di errore (errno settata)

Dopo il prefisso "exec" posso mettere:
- l: (execl(), execlp(), execle()) gli argomenti da passare al programma chiamato sono specificati come lista
- v: (execv(), execvp(), execvpe()) gli argomenti da passare al programma chiamato sono specificati in un array
- e: (execle(), execvpe()) envp specifica l'enviroment. Gli altri wrappers prendono l'enviroment dalla variabile esterna environ
- p: (execlp(),execvp(), execvpe()) se il primmo argomento del wrapper non è un path assoluto il programma da eseguire è cercato in base alla variabile di ambiente PATH

**Execl(3):** `#include int execl(const char *pathname, const char *arg, ... /* (char *) NULL */);`
- Pathname è il programma da eseguire
- Const char *arg* ed i susseguenti argomenti (in numero variabile) equivalgono ad arg0,...,argN e rappresentano gli argomenti del programma da eseguire:
	- Il primo dei suddetti argomenti (arg0) per convenzione è il nome del programma
	- L'ultimo (argN) dovrebbe essere NULL e castato: (char*)NULL

Se non c'è errore non ritornano nulla, -1 in caso di errore (errno settata)

**_ exit(2):** `#include void _exit(int status);`
- Termina il chiamante immediatamente
- I file descriptors sono chiusi
- Il parent ID di eventuali child diventa init (o il più vicino "subreaper process"...)

Non c'è ritorno, ma il processo parent può ottenere il byte meno significativo di status (status & 0xFF) con wait(2)

**exit(3):** `#include void exit(int status);`
- Termina il chiamante
- `exit(3)` fa tutto ciò che fa *_ exit(2)*, ma inoltre:
	- Chiama eventuali funzioni registrate con *atexit(3)* e *on_exit(3)*
	- Esegue il flush dei buffer di I/O relativi alle funzioni di libreria C per l'I/O (stdio(3)streams) e chiude i relativi stream
- Lo standard C definisce EXIT_SUCCESS e EXIT_FAILURE (per portabilità) come valori di uscita

Non c'è ritorno, ma il processo parent può ottenere il byte meno significativo di status (status & 0xFF) con wait(2)

```C
#include <sys/wait.h>
#include <stdio.h>
#include <unistd.h>
#include <sys/stat.h>
#include <fcntl.h>
int main()
{
	pid_t pid;
	/* char *args[] = {"ls", "-l", NULL}; */ // Array di argomenti per cat
	/* char *args[] = {"cat", NULL} */
	char *args[] = {"go.sh", NULL};
	char *envs[] = {"MY_VAR=addio",NULL}; // Array di variabili d'ambiente (vuoto)
	/* fork a child process */
	pid = fork();
	if (pid < 0) { /* error occurred */
		fprintf(stderr, "Fork Failed");
		return 1;
	}
	else if (pid == 0) { /* child process */
		/* execlp("/bin/ls","ls",NULL); */
		int fd = open("etc/passwd", O_RDONLY);
		printf("Child. my pid = %d\n", getpid());
		int x = execve("/usr/bin/ls", args, envs);
		if (x) fprintf(stderr, "Errore:""Program failed to execute\n""|[%d]\n");
		else fprintf(stdout, "Program executed successfully\n");
		perror("Errore in execve");
	}
	else { /* parent process */
	/* parent will wait for the child to complete */
	wait(NULL);
	printf("Child Complete");
	}
	return 0;
}
```

### Pipe
**IPC:** vari strumenti per la comunicazione tra processi:
- Signal
- Socket
- Pipe
- Semafori
- Shared memory
- Memory-mapped file
- Message queue
- Etc

Le pipe sono di solito unidirezionali (in certi sistemi sono bidirezionali come in Windows), sono uno dei primi meccanismi di IPC e sono disponibili praticamente su tutti i sistemi. Connettono due o più processi tramite un canale di comunicazione. Hanno un'estremità per la lettura (*read end*) ed una per la scrittura (*write end*). Sono caratterizzate da comunicazione non strutturata (*byte stream*), cioè <u>le pipe trasmettono un flusso continuo di byte grezzi senza alcuna informazione di controllo, formattazione o metadati</u> associati al contenuto dei dati stessi

Ne esistono 2 tipi:
1. Senza nome (*unnamed*): sono crate con la system call `pipe()`, i processi comunicanti condividono file descriptor per la lettura e la scrittura (`fork(2)`). Solo tra processi imparentati (padre-figlio)
2. Con nome (*named o FIFO*): sono un file speciale (type `p` in `ls -l`) visibile nel filesystem (ma non occupano spazio né implicano operazioni in esso).

> [!info] Una volta create (ed aperte) si usano allo stesso modo
> "The only difference between pipes and FIFOs is the manner in which they are created and opened. Once these tasks have been accomplished, I/O on pipes and FIFOs has exactly the same semantics" (man 7 pipe)

**(unnamed) Pipe:** solo processi relati/discendenti (`fork(2)`). I file descriptor per le due estremità (di lettura e di scrittura) sono duplicati da fork(2) e quindi disponibili anche nel child. Non ha nome nel file system ed esiste solo temporaneamente (sparisce quando i processi terminano)

**Named Pipe o FIFO:** sono visibili nel filesystem ma il "nome" nel filesystem è solo un riferimento. Non occupano spazio disco né implicano operazioni in esso (i dati transitano nel kernel)

> [!info] Da man 7 pipe: 
> "Although FIFOs have a pathname in the filesystem, I/O on FIFOs does not involve operations on the underlying device (if there is one)"

Possono essere trattate come qualsiasi file (SC per l'I/O) e sono visibili quindi nel sistema come gli altri file:
- Qualsiasi processo può fare la open (se ne ha il permesso)
- Possiamo importare permessi per gli utenti
- Possiamo rimuoverle con `rm`

| Comando man per le pipe | Comportamento                                                                                                     |
| ----------------------- | ----------------------------------------------------------------------------------------------------------------- |
| `man 7 pipe`            | Per una panoramica (overview) sulle pipe                                                                          |
| `man 7 fifo`            | Per una panoramica sulle named pipe                                                                               |
| `man 2 pipe`            | Per il manuale della system calls pipe(2) che crea una unnamed pipe (solo tra processi "imparentati" con fork(2)) |
| `man 3 mkfifo`          | Per il manuale della funzione di libreria C che crea una named pipe (visibile nel filesystem)                    |
| `man 1 mkfifo`          | Per il manuale sul comando unix che crea una named pipe nella shell (visibile nel filesystem)                     |
**Creazione di una pipe unnamed:** `#include int pipefd[2]; if (pipe(pipefd) == -1) { perror("pipe"); exit(EXIT_FAILURE); }`
Crea la pipe e apre le due estremità:
- `pipefd[0]` file descriptor per la read-end della pipe 
- `pipefd[1]` file descriptor per la write-end della pipe 
- `Ritorno : 0` se tutto ok. -1 per errore (errno settata)

Creata la pipe è possibile leggere e scrivere dati, per es.: con `read()` e `write()`. Una sequenza tipica di utilizzo è:
- il parent crea la pipe (pipe()) 
- il parent crea un child (fork()). I file descriptor ritornati dalla pipe() sono duplicati nel child 
- lo scrivente chiude il fd di lettura 
- il lettore chiude il fd di scrittura 
- i processi comunicano con read() e write() 
- al termine entrambi chiudono il fd utilizzato

> [!info] Differenza tra `pipe()` e `|`
> - `pipe()`: è una system call che crea una pipe anonima e restituisce 2 file descriptors: int fd[2], pipe(fd)
> - `|`: è la sintassi della bash che usa internamente la funzione `pipe()` (chiama prima `pipe()`, poi `fork()` per creare i processi e infine collega la pipe ai comandi che poi verranno eseguiti)

**Creazione di una named pipe:** `int mkfifo(const char *pathname, mode_t mode)`
- crea una named pipe di nome pathname 
- visibile nel filesystem ma non occupa spazio 
- mode specifica i permessi, in accordo con umask (mode & ~umask) 
- Persistenti, permangono anche dopo il reboot 
- utilizzabile da qualsiasi processo (se ha il permesso…)

**Apertura di una named pipe:** `int open(const char *pathname, int flags);` 
- una volta creata, la named pipe deve essere aperta 
- nel caso in cui il lato opposto non sia ancora aperto `open` è bloccante (a meno che sia settato il flag `O_NONBLOCK`) 
- in Linux è possibile aprire una named pipe in lettura e scrittura. Ciò permette di raggiungere certi scopi ma va usato con cautela. 

> [!info] Da man 7 fifo: 
> "Under Linux, opening a FIFO for read and write will succeed both in blocking and nonblocking mode. POSIX leaves this behavior undefined"

**Capacità delle pipe:** (numero di byte) la pipe hanno una capacità limitata che varia seconda dell'implementazione. In Linux di default attualmente è 16 "system page size" (che dipende dal sistema), ma è configurabile (da 2.6.35) con `fcntl(2)`

**Bloccaggio delle pipe:** le operazioni di apertura/lettura/scrittura possono essere bloccanti o non bloccanti, cioè: 
- ritornare al chiamante solo dopo che si sia verificato un certo evento 
- oppure ritornare subito 

Questo fatto dipende anche da: 
- la condizione della pipe, per esempio se vi sono stati scritti dati, se dall'altro lato della pipe qualcuno legge/scrive, se è piena o c'è ancora spazio, etc. 
- sia stato settato o meno per la pipe il flag `O_NONBLOCK` (vedi anche open(2)) 

Non sono possibili operazioni di seek (no lseek(2)). Di default il flag `O_NONBLOCK` non è settato. Le operazioni allora sono bloccanti nei seguenti casi: 
- `open(2)`: si blocca se dall'altro lato nessun processo ha aperto la pipe nella modalità di accesso opposta (read/write) 
- `read(2)`: si blocca se la pipe è vuota (non ci sono dati) 
- `write(2)`: si blocca se la pipe è piena

**Pipe- scrittura:** `ret = write(fdw, buffer, size)` (vedi man 7 pipe (paragrafo PIPE_BUF) per maggiori dettagli)
- I dati sono scritti nella pipe sequenzialmente
- Se la pipe è piena ulteriori `write(2)` si bloccheranno (o falliranno se `O_NONBLOCK` è settato) 
- POSIX stabilisce che le `write(2)` di ampiezza minore di `PIPE_BUF` sono garantite essere atomiche (ma possono essere bloccanti se la pipe è piena) 
- Stabilisce anche che `PIPE_BUF` deve essere almeno 512 byte
- se tutti i descrittori di lettura sono stati chiusi (se non c'è più un fd dal quale leggo le informazioni) una write provoca l’invio di un segnale `SIGPIPE` al processo: 
	- se `SIGPIPE` è gestito da un handler o ignorato: write ritorna -1 con errno settato a `EPIPE` 
	- se `SIGPIPE` non è gestito: la gestione di default provoca la terminazione del processo 

**Pipe- lettura** `ret = read(fdr, buffer, size)` 
- i dati sono letti in ordine di arrivo 
- i dati letti sono tolti dalla pipe 
- `O_NONBLOCK`:
	- se non settato (generalmente il default) : 
		- se la pipe è vuota: 
			- se i descrittori di scrittura non sono ancora stati tutti chiusi la read si blocca in attesa 
			- se durante tale blocco in attesa tutti i descrittori in scrittura vengono chiusi ritorna 0 (end_of_file) 
			- se invocata quando i descrittori di scrittura sono stati tutti chiusi ritorna 0 (end_of_file) 
		- altrimenti legge al più size byte; se al momento dell’invocazione la pipe contiene meno di size byte, li legge tutti. Ritorna il numero dei byte letti (ret) 
	- se settato: 
		- se la pipe è vuota la read ritorna subito -1; `errno EAGAIN`

**Pipe- close():** `ret = close(fd)` 
- quando l’ultimo descrittore di scrittura è chiuso genera EOF per i lettori (facendo ritornare 0 ad eventuali read in attesa) 
- se viene effettuata una write su una pipe in cui tutti i descrittori di lettura sono stati chiusi il processo riceve `SIGPIPE`

**Named Pipe da shell:** 
- comando `mkfifo(1)` 
- crea un file speciale di tipo `p` 
- ci possiamo scrivere e leggere da shell con i comandi usuali per i file 
- anche qui le operazioni possono essere bloccanti, in accordo con quanto detto sopra 

> [!example]
> - mkfifo mypipe 
> - echo "Ciao" > mypipe 
> - cat mypipe

```C
unn_pipe.c
#include <sys/wait.h>
#include <stdio.h>
#include <unistd.h>
#include <stdlib.h> 

#define READ_END 0
#define WRITE_END 1  

int main()
{
	pid_t pid;
	int pipefd[2], iRet = EXIT_SUCCESS;
	char sBuf[128];
	/* fork a child process */
	if (pipe(pipefd) != -1)
	{
		pid = fork();
		if (pid == 0)
		{
			close(pipefd[READ_END]);
			printf("Child: my pid is %d\n", getpid());
			sprintf(sBuf, "Hello from child process\n");
			if (write(pipefd[WRITE_END], sBuf, sizeof(sBuf)) != -1)
			{
				printf("Child: wrote to pipe\n");
			}
			else
			{
				fprintf(stderr, "Child: Failed to write to pipe\n");
				iRet = EXIT_FAILURE;
			}
			close(pipefd[WRITE_END]);
			return iRet;
		}
		else
		{ /* parent process */
		/* parent will wait for the child to complete */
			close(pipefd[WRITE_END]);
			wait(NULL);
			int iReadCount;
			if ((iReadCount = read(pipefd[READ_END],sBuf,sizeof(sBuf) - 1)) >= 0)
			{
				sBuf[iReadCount] = '\0';
				printf("Parent: read from pipe: \"%s\"\n", sBuf);
			}
			else{
				if (iReadCount == 0)
				{
					fprintf(stderr, "Parent: No data to read from pipe\n");
					iRet = EXIT_FAILURE;
				}
				else
				{
					perror("Parent: Failed to read from pipe\n");
					iRet = EXIT_FAILURE;
				}
			}
			close(pipefd[READ_END]);
		}
	}
	else /* parent success */
	{
		fprintf(stderr, "Pipe Failed\n");
		iRet = EXIT_FAILURE;
	}
	return iRet;
}
```

```C
wfifo.c
#include <sys/wait.h>
#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <fcntl.h>
  
#define MY_FIFO "myfifo"
  
int main()
{
	int iRet = EXIT_SUCCESS, fd;
	char sBuf[128];
	  
	if (access(MY_FIFO, F_OK) != 0)
	{
		if (mkfifo(MY_FIFO, S_IRUSR | S_IWUSR) != 0)
		{
			perror("Error in mkfifo\n");
			return EXIT_FAILURE;
		}
	}
	  
	if (fd = open(MY_FIFO, O_WRONLY) != -1)
	{
		printf("Open in write mode\n");
		sleep(5);
		sprintf(sBuf, "I wrote in the pipe\n");
		if (write(fd, sBuf, sizeof(sBuf)) != -1)
		{
			printf("Wrote in the pipe\n");
		}
		else
		{
			perror("Error in write");
			iRet = EXIT_FAILURE;
		}
		close(fd);
	}
	else
	{
		perror("Error in open");
		iRet = EXIT_FAILURE;
	}
	  
	printf("Write exiting...\n");
	return iRet;
}
```

```C
rfifo.c
#include <sys/wait.h>
#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <fcntl.h>
  
#define MY_FIFO "myfifo"
  
int main()
{
	int iRet  EXIT_SUCCESS, fd;
	char sBuf[128];
	
	if (fd = open(MY_FIFO, O_RDONLY) != -1)
	{
		int iReadCount;
		if ((iReadCount = read(fd, sBuf, sizeof(sBuf)) - 1) != -1)
			{
			printf("Open in read mode\n");
			if (iReadCount == 0)
			{
				printf("Every descriptor closed\n");
			}
			else
			{
				sBuf[iReadCount] = '\0';
				printf("Read from the pipe: \"%s\"\n", sBuf);
			}
		}
		else
		{
			perror("Error in read");
			iRet = EXIT_FAILURE;
		}
		close(fd);
	}
	else
	{
		perror("Error in open");
		iRet = EXIT_FAILURE;
	}
	  
	printf("Read exiting...\n");
	return iRet;
}
```

> [!info]
> - `|` (pipe):  passa lo stdout del comando precedente allo stdin del comando successivo, i comandi vengono eseguiti in maniera concorrente
> - `&&` (AND logico): collega i comandi così che il secondo in successione venga eseguito subito dopo il primo. Se il primo comando fallisce, allora il resto viene saltato
> - `&` (background):  va posto dopo un comando e esegue quest‘ultimo in modo asincrono, consentendo alla shell di continuare immediatamente senza attendere il completamento del processo (ATTENZIONE! NON legge eventuali input da terminale come password) 

### Bash scripting
**`if`:** 
> [!example] 
> ```
> if grep -q Bash file 
> then 
> 	echo "occurrence(s) found" 
> else 
> 	echo "Bash string not found!" 
> fi 
> ```
> verifica se il comando grep va a buon fine (cioè se ritorna 0, non true…)

`&&, ||, !` (and, or, not) si possono usare per combinare gli exit status nelle condizioni 

> [!example]
> ```
> if grep str 1 testo.txt && grep str 2 testo.txt 
> then 
> 	echo "entrambe presenti" 
> fi
> ```

**`test`:** Con test possiamo confrontare stringhe ed interi, testare alcune proprietà dei file (esistenza, permessi, etc.). Possiamo combinare logicamente le condizioni (-o, -a)

> [!example] 
> ```
> if test -d dir1 
> then 
> 	echo "la directory esiste" 
> else 
> 	echo "la directory non esiste" 
> fi
> ```

> [!info] Confronti tra interi (tratto da man test): 
> - INTEGER1 -eq INTEGER2: INTEGER1 is equal to INTEGER2 
> - INTEGER1 -ge INTEGER2: INTEGER1 is greater than or equal to INTEGER2
> - INTEGER1 -gt INTEGER2: INTEGER1 is greater than INTEGER2 
> - INTEGER1 -le INTEGER2: INTEGER1 is less than or equal to INTEGER2
> - INTEGER1 -lt INTEGER2: INTEGER1 is less than INTEGER2
> - INTEGER1 -ne INTEGER2: INTEGER1 is not equal to INTEGER2

> [!info] Opzioni per i file 
> - `-a file` o `-e file:` Vero se il file esiste 
> - `-d file:` Vero se il file esiste ed è una directory 
> - `-f file:` Vero se il file esiste ed è un file regolare 
> - `-h file:` Vero se il file esiste ed è un link simbolico 
> - `-r file`: Vero se il file esiste ed è leggibile 
> - `-w file:` Vero se il file esiste ed è scrivibile 
> - `-x file`: Vero se il file esiste ed è eseguibile

**`if [...]:`**  
> [!example] 
> ```bash
> if [ -f file 1 ] 
> then 
> 	echo "file 1 esiste" 
> else 
> 	echo "file 1 non esiste" 
> fi 
> ```

`[]` è sinonimo di `test`, è builtin 

**`if [[...]]:`** 
> [!example] 
> ```bash
> if [[ -f file 1 ]]
> then 
> 	echo "file 1 esiste" 
> else 
> 	echo "file 1 non esiste" 
> fi 
> ```

`[[ ]]` è simile a `[ ]`, possiamo considerarlo un test esteso e con sintassi più familiare al programmatore 

**Costrutti condizionali:** `if [[ $a -lt $b ]]` considerato dalla Bash come un unico elemento con un valore di ritorno. Alcune espansioni non sono attuate e le variabili e i comandi sono sostituiti. Alcuni errori logici sono evitati. Per es. posso usare `&&`, `||`, `<` e `>` 

**`if ((...)):`** 
> [!example] 
> ```bash
> if (( var 1 > var 2+10 )) 
> then # ^ ^ Note: Not $var 1, $var2… 
> 	echo "$var 1 is greater than $var2 + 10" 
> else 
> 	echo "$var 1 is less/equal than/to $var2 + 10" 
> fi
> ```

Si usa per valutare espressioni aritmetiche 

**Costrutto case:** analogo al switch/case del C: 
> [!example]
> ```bash
>  case "$VAR 1" in 
> 	 "a") echo "VAR 1 = a";; 
> 	 "b") echo "VAR 1 = a";; 
> 	 *) echo "Nulla di tutto questo";; 
>  esac
> ``` 

in generale la struttura è:

```bash
case <expr> in
	<pattern 1>)
		<comandi>;;
	<pattern 2>)
		<comandi>;;
		...
	<pattern n>)
		<comandi>;;
	*) <comandi>;; # default, i.e. anything else 
esac
```

**Iterazioni:** 
- `while test-commands; do <commands>; done`>: Permane nel ciclo finché il valore di ritorno dei test-commands è 0. `while` torna il valore di ritorno dell'ultimo comando eseguito o 0 se non sono stati eseguiti comandi
- `until test-commands; do <commands>; done`: Permane nel ciclo finché il valore di ritorno dei test-commands è diverso 0. `until` torna il valore di ritorno dell'ultimo comando eseguito o 0 se non sono stati eseguiti comandi 
- `for i [in words … ]; do <commands>; done`: Per ogni stringa presente nell’espansione di words esegue commands. Ad ogni ciclo `$i` contiene la stringa corrente. `for` torna il valore di ritorno dell'ultimo comando eseguito o 0 se non sono stati eseguiti comandi. Se la clausola in non è specificata `for` considera le stringhe derivanti dall'espansione di `$@` 
- `for (( expr1; expr2; expr3 )); do commands; done`: Sintassi C-like. 

> [!example] 
> ```bash
> for ((i=1;i<=10;i++)); 
> do 
> 	ls file$i; 
> done 
> ```

- `break e continue`: analoghi ad altri linguaggi 

`for`: itera l'esecuzione di un blocco di istruzioni. La variabile di loop assume ad ogni iterazione un valore diverso in base agli argomenti passati. Non specifichiamo quante iterazioni fare (come in C) ma una lista di valori che la variabile di loop deve assumere 

> [!example]
> ```bash
> for <i> [ in <lista> ]; do
> 	<comando 1>
> 	<comando 2>
> 	...
> 	<comando n>
> done
> ``` 

Se la clausola in non è specificata si assumono come la lista degli argomenti dello script (`$@`) 

- la shell espande 
- `for` esegue un ciclo per ogni elemento nell'espansione (il separatore tra gli elementi è il primo carattere della variabile IFS, di default spazio) 
- ad ogni iterazione `<i>` contiene l'elemento corrente 
- `for` ha come valore di ritorno quello dell’ultimo comando eseguito (o 0 se nessun comando è stato eseguito)

> [!example]
> ```bash
> for root@dbserver :~# ls -w1 
> 1.sh 
> 1.txt 
> dead.letter 
> upgrade-stretch1.script 
> upgrade-stretch1.time 
> root@dbserver :~# cat 1.sh 
> #!/bin/bash 
> for i in * ; do 
> 	echo "La variabile \$i contiene: ${i}" 
> done 
> root@dbserver :~# ./1.sh 
> La variabile $i contiene: 1.sh 
> La variabile $i contiene: 1.txt 
> La variabile $i contiene: dead.letter 
> La variabile $i contiene: upgrade-stretch 1.script 
> La variabile $i contiene: upgrade-stretch 1.time 
> root@dbserver :~# 
>```
 
> [!example] 
> ```bash
> for ((i=0; i<=10; i++)); 
> do 
> 	echo $i 
> done
> ``` 
> stampa i numeri da 1 a 10 
 
> [!example] for da linea di comando 
> ```bash
> root@dbserver :~# for((i=0; i<=10; i++)); do echo $i; done 
> 0
> 1 
> 2 
> 3 
> 4 
> 5 
> 6 
> 7 
> 8 
> 9 
> 10 
> root@dbserver:~# 
> ``` 

`while`: itera fin quando una condizione è vera: 

> [!example] 
> ```bash
> root@dbserver :~# cat 2.sh 
> #!/bin/bash 
> COUNT=0 
> while test $COUNT -lt 5; do 
> 	echo $COUNT 
> 	COUNT=$((COUNT+1)) 
> done 
> root@dbserver:~# ./2.sh 
> 0 
> 1 
> 2 
> 3 
> 4 
> root@dbserver :~#
> ```

> [!example] itera all'infinito, interrompibile con Ctrl+c: 
> ```bash
> while test 1; do 
> 	clear 
> 	df -h 
> 	sleep 5 
> done
> ``` 

**`read` (comando interno):** serve per acquisire input, per esempio, da tastiera: 
> [!example] 
> ```bash
> user1@host1:~$ read -p "Inserisci un saluto: " var1 
> Inserisci un saluto: Ciao 
> user1@host1:~$ echo $var1 
> Ciao 
> user1@host1:~$ set |grep var 1 
> var 1=Ciao 
> user1@host1:~$ 
> ```

L'opzione `-p` visualizza un prompt prima di attendere per l'input. In generale per avere l'help su di un comando interno: `help comando`. `read` ha l'opzione `-u` che permette di redirigere l'input da un file descriptor invece che da stdin. 

> [!example]
> ```bash
> user1@host1:~$ exec 3< testo.txt 
> user1@host1:~$ while read -u 3 linea ; do echo $linea; done 
> fabio 
> ciccio 
> pippo 
> user1@host1:~$
> ``` 

`exec`, in generale, esegue il comando passato come argomento creando un processo che sostituisce la shell corrente. Senza questo argomento qualsiasi redirezione è operata nella shell corrente. In questo caso testo.txt è aperto in lettura con fd = 3. `read -u 3` fa si che `read` prenda l'input da testo.txt invece che da stdin 

> [!info] In altre parole `exec` e `read -u` possono essere utilizzati in una shell, e quindi in uno script, per aprire in lettura (o scrittura) un file e per elaborarne (o scriverne) il contenuto

**Array:** Per definirne uno la sintassi è:

```bash
area[11]=23 
area[13]=37 # indici anche non adiacenti 
area[51]=UFOs 
area2=( zero one two three four ) 
declare -a array1 
```

Gli elementi si referenziano tra parentesi graffe: 

```bash
echo ${area[11]} 
echo ${area2[0]} # visualizza "zero"; 
# indici zero-based se non specificati 
```

> [!example] 
> ```bash
> fabio@fabiopc:~$ area[11] =23 
> fabio@fabiopc:~$ area[13]=37 
> fabio@fabiopc:~$ area[51]=UFOs 
> fabio@fabiopc:~$ area 2=( zero one two three ) 
> fabio@fabiopc:~$ set|grep area 
> area=([1]="0" [11]="23" [13]="37" [51]="UFOs") 
> area 2=([0]="zero" [1]="one" [2]="two" [3]="three) 
> fabio@fabiopc:~$ echo ${area 2[0]} ${area[51]} zero UFOs fabio@fabiopc:~$ 
> fabio@fabiopc:~$ echo ${area2[@]} # tutti gli elementi 
> zero one two three four 
> fabio@fabiopc:~$ echo ${area2[*]} # tutti gli elementi 
> zero one two three four 
> fabio@fabiopc:~$ 
> ```

Per determinare il numero di elementi di un array si usa la sintassi `${ #arrayname [@]}:`, mentre si usa `unset nomearray` per cancellare l'array

> [!example]
> ```bash
> fabio@fabiopc:~$ area 2=( zero one two three four ) 
> fabio@fabiopc:~$ echo ${ #area2 [@]} 
> 5 
> fabio@fabiopc:~$
> ```

**Funzioni:** Una funzione è un insieme di comandi che può essere richiamato più volte ed è identificato dal nome della funzione stessa, richiamabile come un comando interno. Sono eseguite nella shell corrente (non subshell) e lo scopo è di evitare di riscrivere lo stesso codice più volte e rendere gli script più comprensibili. Rispetto ad altri linguaggi il supporto di Bash per le funzioni è più limitato:
- Cancellabili con `unset` 
- Per vedere le funzioni definite: 
	- `bash:~$ declare -f`: visualizza il nome di tutte le funzioni ed il loro codice
	- `bash:~$ declare -F`: visualizza solo i nomi di tutte le funzioni 
	- `bash:~$ type -all name_function`: fornisce nome e codice della funzione di nome name_function  
	
Per definire una funzione ho due differenti formati:
- `function_name () { commands }` oppure, a linea singola, `function_name () { commands; }` 
- `function function_name { commands }` oppure, a linea singola, `function function_name { commands; }` 
- non solo negli script ma anche a linea di comando 

> [!info] Note da https://linuxize.com/post/bash-functions/
> - The commands between the curly braces (`{}`) are called the body of the function. The curly braces must be separated from the body by spaces or newlines. 
> - Defining a function doesn’t execute it. To invoke a bash function, simply use the function name. Commands between the curly braces are executed whenever the function is called in the shell script. 
> - The function definition must be placed before any calls to the function. 
> - When using single line “compacted” functions, a semicolon `;` must follow the last command in the function

> [!example] 
> ```bash
> #!/bin/bash 
> hello_world () { 
> 	echo 'hello, world' 
> } 
> 
> hello_world 
> ```

In Bash di default tutte le variabili sono definite come globali, anche se dichiarate all'interno di una funzione. È possibile accedere alle variabili globali da qualsiasi punto dello script, indipendentemente dall'ambito. Le variabili locali possono essere dichiarate all'interno del corpo della funzione con la parola chiave `local` e possono essere utilizzate solo all'interno di quella funzione. È possibile avere variabili locali con lo stesso nome in funzioni diverse 

> [!example] 
> ```bash
> #!/bin/bash 
> var 1='A' 
> var 2='B' 
> 
> my_function () { 
> 	local var 1='C' 
> 	var 2='D' 
> 	echo "Inside function: var 1: $var 1, var 2: $var 2" 
> } 
> 
> echo "Before executing function: var 1: $var 1, var 2: $var 2" 
> 
> my_function 
> 
> echo "After executing function: var 1: $var 1, var 2: $var 2"
> 
> Before executing function: var 1: A, var 2: B 
> Inside function: var 1: C, var 2: D 
> After executing function: var 1: A, var 2: D
> ```

Per quanto riguarda gli argomenti vale quanto detto per i parametri posizionali passati ad uno script: 
- sono variabili nella forma `$n` dove n è un intero da 0 in poi 
- `$0` contiene il nome della funzione 
- gli altri argomenti (se passati) sono da `$1` a `$n` dove n è il loro numero 
- quelli oltre il nono devono essere referenziati nella forma `${n}` (per es.: `$10` è sostituito col primo + il carattere "0" e non il decimo) 

> [!example]
> ```bash
> fabio@fabiopc:~$ cat sc11.sh 
> #!/bin/bash 
> 
> greeting () { 
> 	echo "Grettings from $0:" 
> 	while test $# -gt 0; do 
> 		echo -e "Hello $1" 
> 		shift 
> 	done 
> } 
> 
> greeting Fabio Paolo Mario 
> fabio@fabiopc:~$ ./sc11.sh 
> Grettings from ./sc11.sh: 
> Hello Fabio 
> Hello Paolo 
> Hello Mario 
> fabio@fabiopc:~$ 
> ```

Di default il valore restituito da una funzione è lo stato dell'ultima istruzione eseguita, 0 per successo e un numero decimale diverso da zero, nell'intervallo 1 - 255, in caso di errore. Le funzioni Bash non consentono di assegnare il valore restituito ad una variabile a piacimento; viene invece assegnato alla variabile `$?`. Possiamo specificare un valore "arbitrario" da restituire utilizzando la parola chiave `return`. L'istruzione return termina la funzione: 

> [!example]
> ```bash
> #!/bin/bash 
> my_function () { 
> 	echo "some result" 
> 	return 55 } 
> my_function 
> echo $? 
> 
> some result 
> 55 
> ```

Per assegnare il valore di ritorno ad una variabile arbitraria dobbiamo utilizzare altri metodi:
- Assegnare un valore ad una variabile globale: 

> [!example] 
> ```bash
> #!/bin/bash 
> my_function () { 
> 	func_result="some result" 
> } 
> my_function 
> echo $func_result 
> 
> some result
> ```

- inviare il ritorno a stdout: 

> [!example]
> ```bash
> #!/bin/bash 
> my_function () { 
> 	local func_result="some result" 
> 	echo "$func_result" 
> } 
> func_result="$(my_function)" # oppure `my_function` 
> echo $func_result 
> 
> some result
> ```
### Threads
**Posix threads:** sono parti di uno stesso programma eseguite contemporaneamente e schedulate indipendentemente l'una dall'altra. In altre parole sono "più tracce di esecuzione concorrenti in un singolo processo". Un singolo processo può "contenere" più thread. Questi condividono la memoria globale (data and heap) ma ognuno ha il proprio stack. Non è un'implementazione, ma uno standard. Sono implementati in molti SO, ma con delle differenze (vengono definiti una serie di interfacce come funzioni, file header per la programmazione di threads). Sono conosciuti come POSIX Threads o Pthreads 

> [!info] `man 7 pthreads`

Nel corso del tempo due implementazioni nella GNU C library:
- *LinuxThreads:* Implementazione originale. Non più supportata (a partire da glibc 2.4) 
- NPTL (Native POSIX Threads Library) 
	- implementazione "più moderna" 
	- più rispondente allo standard POSIX.1 
	- prestazioni migliori se si crea un gran numero di threads 
	- disponibile a partire da glibc 2.3.2 e richiede caratteristiche presenti a partire da Linux 2.6 

Per determinare la versione di glibc: `ldd --version`. Per conoscere quale implementazione dei threads: `getconf GNU_LIBPTHREAD_VERSION` 

POSIX.1 richiede che i seguenti attributi siano condivisi tra i threads di uno stesso processo: 
- process ID 
- parent process ID 
- process group ID and session ID 
- controlling terminal 
- user and group IDs 
- open file descriptors 
- record locks (see fcntl(2)) 
- signal dispositions 
- file mode creation mask (umask(2)) 
- current directory (chdir(2)) and root directory (chroot(2)) 
- interval timers (setitimer(2)) and POSIX timers (timer_create(2)) 
- nice value (setpriority(2)) 
- resource limits (setrlimit(2)) 
- measurements of the consumption of CPU time (times(2)) and resources (getrusage(2)) 

Attributi NON condivisi (per-thread):
- thread ID (the pthread_t data type) 
- signal mask (pthread_sigmask(3)) 
- the errno variable 
- alternate signal stack (sigaltstack(2)) 
- real-time scheduling policy and priority (sched(7))

**Posix Threads: ritorno delle funzioni:** La maggior parte delle funzioni pthreads tornano 0 se non ci sono errori ed un codice di errore altrimenti. Questo codice di errore ha lo stesso significato dei valori convenzionali di errno, ma le funzioni Pthreads NON settano errno. POSIX.1-2001 stabilisce che EINTR non può essere ritornato 

**Threads IDs:** Ogni thread ha un unico thread ID* utilizzabile nelle funzioni pthreads. Questo viene tornato al chiamante da pthread_create(3) o può essere ottenuto con pthread_self(3). Gli ID sono unici nel processo ma non nel sistema. Il SO riusa gli ID di thread terminati, ma se un'applicazione usa un ID di un thread terminato il risultato non è definito

**Thread Safety:** 
- *POSIX Safety Concepts:* definiscono se una funzione può essere chiamata in determinati contesti senza che si verifichino comportamenti indesiderati 
- *MT-Safe (MT sta per Multi-threading) o Thread-Safe:* una funzione lo è se può essere chiamata contemporaneamente da più thread senza modificare il comportamento. MT-Safe non significa che l'esecuzione sia atomica, né che usi alcun meccanismo di sincronizzazione nell'accesso alla memoria. È possible che chiamare più funzioni MT-Safe in sequenza possa sortire effetti diversi da un'esecuzione atomica delle stesse chiamate concorrenti da altri thread possono interferire 

> [!info] vedi `man 7 attributes` per chiarimenti su MT-Safe

**Compilazione:** 
- includere `pthread.h `
- se richiesto compilare col flag `-pthread` 

**Creazione di threads:** `#include <pthread.h> int pthread_create(pthread_t *thread, const pthread_attr_t *attr, void *(*start_routine) (void *), void *arg);` 
- avvia un nuovo thread nel processo chiamante 
- imposta in thread il thread ID del thread creato che può essere utilizzato dalle altre funzioni pthreads
- il thread avvia l'esecuzione chiamando la funzione `start_routine` 
- `arg` viene passato come unico argomento a `start_routine` 
- un esempio d'uso in `man 3 pthread_create` 
- `attr` punta ad una struttura di tipo `pthread_attr_t` che può essere usata per settare alcuni attributi del nuovo thread, per esempio per crearlo `detached` 
- `attr` può essere inizializzata con `pthread_attr_init(3).`
- Se è `NULL` sono usati gli attributi di default 
- `pthread_create(3)` torna 0 per successo; un codice di errore altrimenti 

**Terminazione di un thread:** un thread può terminare: 
- Se il thread chiama `pthread_exit(3)` specificando un codice d'uscita, che sarà disponibile per un altro thread del processo che abbia chiamato `pthread_join(3)` 
- quando `start_routine()` (vedi `pthread_create()`) termina con return. Equivalente a chiamare `pthread_exit(3)` 
- è cancellato con `pthread_cancel(3)` 
- un thread chiama `exit(3)` o quello principale chiama return. L'intero processo termina (e tutti i suoi threads)

`pthread_exit(void *retval);` termina il thread chiamante. Se il thread è joinable, il valore ritornato in `retval` è reso disponibile ad un altro thread del processo che chiami `pthread_join(3)`. Il valore cui `retval` punta non dovrebbe essere allocato nello stack del thread chiamante, perché dopo la terminazione del thread il contenuto dello stack è "undefined"

> [!example] mt 1.c 
> ```C
> #include <pthread.h>
> #include <stdio.h>
> #include <unistd.h>
> 
> int iCondivisa = 13; 
> 
> void thread1(int *arg) { 
> 	int iLocale = 101; 
> 	printf("Sono il primo thread. Argomento = %d, var. condivisa = %d, locale = %d\n", *arg, iCondivisa, iLocale); 
> 	pthread_exit(NULL); 
> } 
> 
> void thread2(char *arg) { 
> 	int iLocale = 202; 
> 	printf("Sono il secondo thread. Argomento = \"%s\", var. condivisa = %d, locale = %d\n", arg, iCondivisa, iLocale);
> 	pthread_exit(NULL); 
> } 
> 
> void main() { 
> 	pthread_t th1, th2; 
> 	int arg1 = 1; 
> 	char arg2[] = "Ciao"; 
> 	pthread_create(&th1, NULL, (void *)thread1, (void *)&arg1);
> 	pthread_create(&th2, NULL, (void *)thread2, (void *)arg2);
> 	sleep(2); 
> 	return; 
> }
> ``` 
> - variabili globali condivise tra threads 
> - comunicazione tra threads realizzabile con le risorse (variabili) condivise ma non c'è nessuna sincronizzazione 
> - stack (var. locali) specifico per ogni thread 
> - nessun vincolo alla schedulazione dei threads. POSIX definisce funzioni che consentono di modificare schedulazione e priorità 

**`join e detach`:** un thread può essere joinable o detached. Di default i threads sono creati joinable. Se un thread è joinable un altro thread può richiamare `pthread_join(3)` per attendere che termini e per ottenere il suo exit status 

Le risorse allocate per un thread joinable sono rilasciate solo se un altro thread chiama `pthread_join(3)` per attenderlo. Se nessun `join` viene fatto il thread terminato è in stato zombie e le risorse sono rilasciate solo all'uscita dal processo. Se i threads zombie si accumulano oltre misura non sarà più possible creare altri thread. Ogni thread nel processo può attendere con `pthread_join(3)` qualsiasi altro thread (del processo stesso). Non c'è per i threads un'operazione analoga a quella vista per i processi child (`fork(2)`): `waitpid(-1, &status, 0)` cioè una chiamata che esegue join su qualsiasi thread termini 

**Thread detached:** E' possibile creare un thread detached modificandone gli attributi in fase di creazione con `pthread_attr_setdetachstate(3)` oppure richiamando `pthread_detach(3)` per un thread già creato. Le risorse allocate per un thread detached sono subito liberate alla sua terminazione. Non è possibile eseguire il `join` ad un thread detached per ottenere il suo exit status, né è possibile renderlo joinable. Creare thread detached può essere utile quando non ci curiamo dell'exit status dei thread; per esempio nello sviluppo di certi daemon 

**Conclusioni su join e detach:** una tra `pthread_join()` e `pthread_detach()*` dovrebbe essere chiamata per ogni thread creato affinché le risorse vengano liberate, anche se, all'uscita del processo, tutte le risorse dei thread sono liberate comunque 

**Joining threads:** `int pthread_join(pthread_t thread, void **retval);` Attende che il thread specificato da thread termini (se è joinable). Se retval non è NULL riceve il valore tornato da `pthread_exit(3)` oppure `PTHREAD_CANCELED` se il thread è stato cancellato. Se più thread cercano simultaneamente di fare il join allo stesso thread il risultato non è definito 

**Detaching threads:** `int pthread_detach(pthread_t thread);` rende il thread specificato da thread detached 

> [!example] mt 1.c con join 
> ```C
> #include <pthread.h> 
> #include <stdio.h>
> #include <unistd.h>
> 
> int iCondivisa = 13; 
> void thread 1(int *arg) { 
> 	int iLocale = 101; 
> 	printf("Sono il primo thread. Parametro = %d, var. condivisa = %d, locale = %d\n", *arg, iCondivisa, iLocale); 
> 	pthread_exit((void *)1000);
> } 
> 
> void thread 2(char *arg) { 
> 	 int iLocale = 202; 
> 	 printf("Sono il secondo thread. Parametro = \"%s\", var. condivisa = %d, locale = %d\n", arg, iCondivisa, iLocale); 
> 	 pthread_exit((void *)2000); 
> } 
> 
> void main() { 
> 	pthread_t th 1, th 2; 
> 	int arg 1 = 1; 
> 	char arg 2[] = "Ciao"; 
> 	long l 1, l 2; 
> 	pthread_create(&th 1, NULL, (void *)thread 1, (void *)&arg 1); 
> 	pthread_create(&th 2, NULL, (void *)thread 2, (void *)arg 2); 
> 	pthread_join(th 1, (void **)&l 1); 
> 	pthread_join(th 2, (void **)&l 2); 
> 	printf("\nRitorno thread 1 = %ld\nRitorno thread 2 = %ld\n", l 1, l 2); 
> 	return; 
> }
> ```

> [!info] In sostanza:
> - con `join` il thread chiamante attende che il thread interpellato termini così da usarne le risorse (un solo thread può fare `join` con successo). Se non faccio mai `join` accumulo zombie e ad un certo punto non posso più creare thread
> - quando il thread termina, con `detached`, le risorse vengono liberate subito e questo sparisce immediatamente. Non posso farci il `join` ne ottenere il suo valore di ritorno ne tornare indietro
### Race condition
Riguarda il problema della concorrenza. 

> [!quote] Dalle slide del Prof. Carpi:  
> "Si verifica quando più processi o threads leggono e scrivono dati in modo che il risultato finale dipende dall’ordine di esecuzione delle istruzioni dei processi"

**Sezione critica:** porzione di codice eseguita da più threads/processi per la quale il risultato delle esecuzioni concorrenti dipende dalla sequenza con cui i threads/processi accedono a questa sezione di codice (e il risultato può non essere quello desiderato...). In questo caso si dice che la sezione critica contiene una race condition 

```C
// Esempio tratto dall’IA: 
#include <stdio.h> 
#include <pthread.h>

int counter = 0; // variabile condivisa 

void *worker(void *arg) { 
	for (int i = 0; i < 100000; i++) { 
		counter++; // <-- NON atomico 
	} return NULL; 
} 

int main(void) { 
	pthread_t t1, t2; 
	pthread_create(&t1, NULL, worker, NULL); 
	pthread_create(&t2, NULL, worker, NULL); 
	
	pthread_join(t1, NULL); 
	pthread_join(t2, NULL); 
	
	printf("counter = %d\n", counter); 
	return 0; 
} 
```

Ci aspetteremmo che counter assumesse a fine esecuzione il valore 200.000: può non essere così. `++` NON è un’operazione atomica: `counter++` è una sezione critica. Dobbiamo quindi garantire che i thread abbiano accesso mutuamente esclusivo alle sezioni critiche. Si hanno più soluzioni, tra cui: 
- *pthread_mutex:* parte della libreria pthreads 
- supporto per le operazioni atomiche in C: Atomic operations library, a partire da C11
- *semafori* 

> [!info] Ricorda:
> Dire che qualcosa è atomico significa che viene eseguito senza essere interrotto in una volta sola (nessun altro thread può vedere uno stato a "metà"). Una variabile atomica è tale se letture e scritture avvengono in modo indivisibile (e due threads non possono incastrarsi mentre la modificano). Un'operazione è atomica, invece, se avviene tutta o non avviene per nulla (anche qui senza interferenze da parte di altri thread). 
>> [!example]
>> `x = x + 1;` anche se sembra un'unica operazione implica di leggere x, sommargli 1 e poi scrivere x. Tra questi passaggi può succedere che un altro thread intervenga
 
**Pthreads: Mutex:** 

> [!quote] Da `man pthread_mutex_init(3)*`: 
"A mutex is a MUTual EXclusion device, and is useful for protecting shared data structures from concurrent modifications, and implementing critical sections and monitors " 

In pthreads abbiamo a disposizione un "MUTual EXclusion device", le variabili mutex. Una mutex può avere due possibili stati: locked ("owned" da un thread) e unlocked (nessun thread è "owner"). Solo un thread alla volta può essere " owner ". Se un altro thread tenta di acquisire ownership su una mutex già owned da altro thread si blocca finché il thread attualmente owner non la rilascia 

**Fasi per una mutex:**
1. si definisce una variabile `pthread_mutex_t mutex` per la singola sezione critica e la si inizializza (con pthread_mutex_init(3), PTHREAD_MUTEX_INITIALIZER, etc. Vedi man pthread_mutex_init(3)) 
2. si acquisisce l'accesso esclusivo alla sezione critica con `pthread_mutex_lock(&mutex)` 
3. si elabora la sezione critica 
4. si rilascia l'accesso esclusivo con `pthread_mutex_unlock(&mutex)` 
5. alla fine, quando la variabile mutex non serve più si può chiamare `pthread_mutex_destroy` che libera le risorse allocate per la mutex (deve essere unlocked) 

```C
#include <stdio.h>
#include <pthread.h>

pthread_mutex_t my_mutex = PTHREAD_MUTEX_INITIALIZER; // my_mutex per proteggere l'accesso a counter
int counter = 0; // variabile condivisa
  
void *worker(void *arg)
{
	for (int i = 0; i < 100000; i++) {
		pthread_mutex_lock(&my_mutex); // blocca il my_mutex prima di accedere a counter
		counter++; // <-- NON atomico
		pthread_mutex_unlock(&my_mutex); // sblocca il my_mutex dopo aver modificato counter
	}
	return NULL;
}
  
int main(void)
{
	pthread_t t1, t2;
	  
	pthread_create(&t1, NULL, worker, NULL);
	pthread_create(&t2, NULL, worker, NULL);
	  
	pthread_join(t1, NULL);
	pthread_join(t2, NULL);
	  
	printf("counter = %d\n", counter);
	pthread_mutex_destroy(&my_mutex); // distrugge il my_mutex
	
	return 0;
}
```
 
**C11 Atomic support:** un'alternativa in certi casi per gestire la concorrenza è stata introdotta in C11. 
- Si tratta di uno standard del passato per il linguaggio C 
- Tra le altre cose introduce il supporto per il *multithreading* (da non confondere con le funzioni di libreria pthreads) 

Introduce anche il supporto per variabili ed operazioni atomiche: `stdatomic.h`, qualificatore `_Atomic` e varie primitive per le operazioni. Il qualificatore `_Atomic` permette di dichiarare variabili alle quali è garantito un accesso atomico. Per esempio, un'operazione di incremento, come `a++`, che implica un accesso alla variabile (lettura del valore) ed un incremento, è garantita essere atomica (occorre includere `stdatomic.h`). In `stdatomic.h` sono definite delle typedef che creano "nuovi" tipi di dato analoghe alle seguenti: 

| Typedef name      | Full type name         |
| ----------------- | ---------------------- |
| atomic_char       | _Atomic char           |
| atomic_schar      | _Atomic signed char    |
| atomic_uchar      | _Atomic unsigned har   |
| atomic_short      | _Atomic short          |
| atomic_ushort     | _Atomic unsigned short |
| atomic_int_Atomic | _Atomic int            |
> [!info]
>  C11 introduce anche varie funzioni per compiere più operazioni su variabili in maniera atomica. Ad esempio C11 introduce funzioni analoghe a: `atomic_fetch_add( volatile A* obj, M arg )` che atomicamente aggiunge arg al valore di obj

```C
#include <stdio.h>
#include <pthread.h>
#include <stdatomic.h>

_Atomic int counter = 0; // variabile condivisa

void *worker(void *arg)
{
	for (int i = 0; i < 100000; i++) {
		counter++; 
	}
	return NULL;
}
  
int main(void)
{
	pthread_t t1, t2;
	  
	pthread_create(&t1, NULL, worker, NULL);
	pthread_create(&t2, NULL, worker, NULL);
	  
	pthread_join(t1, NULL);
	pthread_join(t2, NULL);
	  
	printf("counter = %d\n", counter);
	return 0;
}
```

### Semafori (Posix):
Servono per la sincronizzazione dei processi o dei threads per eseguire coerentemente una serie di azioni, per esempio per accedere ordinatamente a sezioni critiche di un programma o a risorse condivise. Si usano anche per la sincronizzazione dell’accesso ai dati per garantirne integrità e coerenza

Storicamente ci sono 2 possibili implementazioni: POSIX o SYSTEM V. POSIX è più semplice ma meno diffusa (soprattutto tra sistemi più vecchi). SYSTEM V è una API più tradizionale

I semafori (in POSIX, `man 7 sem_overview`) sono rappresentabili come variabili intere, ma mai minori di 0.

> [!warning] Attenzione!
> Non confondere questa informazione "i semafori non possono essere mai minori di 0", con la teoria (per la quale i semafori non binari possono assumere valori negativi)

Ci sono 2 operazioni possibili:
1. `POST`: incrementa di 1
2. `WAIT`: decrementa di 1

Non potendo essere negativi se una `WAIT` trova la variabile uguale a 0 si blocca fin quando il valore non diviene positivo. Solo allora il decremento può essere effettuato e la WAIT ritorna. 

I semafori POSIX possono essere:
- *Named:* 
	- sono identificati da un nome nella forma "/somename". 
	- `sem_open(3)` apre o crea (se non esiste) un semaforo *named*. 
	- Due processi/thread possono operare sullo stesso semaforo passando lo stesso nome a `sem_open(3)`.
	- Dopo averlo aperto i processi possono operare con `sem_post(3)` e `sem_wait(3)`
	- Quando un semaforo non serve più il processo può chiudere il suo accesso con `sem_close(3)`.
	- Quando tutti i processi hanno terminato di usarlo, il semaforo può essere rimosso con `sem_unlink(3)` 

> [!info] Ricorda:
> Sono persistenti nel kernel, cioè se non rimossi permangono fino al riavvio della macchina

**FUNZIONAMENTO DI SEM_OPEN:**
```C
#include <fcntl.h> /* For O_* constants */ 
#include <sys/stat.h> /* For mode constants */ 
#include <semaphore.h>

sem_t *sem_open(const char *name, int oflag); 
sem_t *sem_open(const char *name, int oflag, mode_t mode, unsigned int value);
```

- Apre un semaforo *named* (o lo crea se si specifica il flag `O_CREAT`)
- Gli `oflag` sono analoghi a `open(2)`
- `mode` specifica i permessi come in `open(2)`
- `value` è il valore iniziale per il semaforo
- Se `O_CREAT` viene specificato allora `mode` e `value` devono essere forniti (ma se il semaforo già esiste sono ignorati)
- L'utente proprietario del semaforo eventualmente creato è lo stesso utente proprietario del processo chiamante 

Ritorna l’indirizzo del semaforo. In caso di errore torna `SEM_FAILED` e setta `errno`

**Sequenza di utilizzo:**
```C
#include <fcntl.h> /* For O_* constants */ 
#include <sys/stat.h> /* For mode constants */ 
#include <semaphore.h>

sem_t *sem; 
... 
sem = sem_open("/mysem", O_CREAT, S_IRWXU, 0); 
... 
sem_post(sem); o sem_wait(sem); (o sem_timedwait()) 
... 
sem_close(sem); 
... 
sem_unlink("/mysem");
```

In linux sono "visibili" ed accessibili da filesystem, generalmente in `/dev/shm`. Hanno nomi nella forma *sem.somename* e possiamo assegnare permessi (come ogni altro oggetto nel filesystem). 

- Unnamed: 
	- non hanno nome e sono allocati in un'area di memoria condivisa tra i processi/thread che li usano:
	- *Thread-shared:* l'area di memoria può essere una variabile globale del processo (ovviamente per thread di uno stesso processo)
	- *Process-shared:* allocati un una shred memory "mappata" nello spazio di indirizzamento dei processi/threads (con un puntatore) 
		- SYSTEM V: creata con shmget(2) 
		- POSIX: creata con shm_open(3)
	- Vengono inizializzati con `sem_init(3)`. 
	- Dopo l'inizializzazione si può operare con `sem_post(3), sem_wait(3)/sem_timedwait(3)` 
	- Quando non più necessari, prima di deallocare un'eventuale area di memoria allocata per essi, dovrebbero essere distrutti con sem_destroy(3), altrimenti in certe implementazioni certe risorse potrebbero rimanere occupate (leak)

**Semafori POSIX Unnamed:** `#include int sem_init(sem_t *sem, int pshared, unsigned int value);` Inizializza un semaforo unnamed all'indirizzo specificato da `sem`. `pshared` può essere uguale a: 
- 0: semaforo condiviso tra threads di uno stesso processo (può essere allocato in una variabile globale) 
- non 0: semaforo condiviso tra processi. Allocato in una shared memory (vedi `shm_open(3), mmap(2), and shmget(2)`). No variabile globale tra processi discendenti (fork()) 

`value` è il valore iniziale per il semaforo. Avremo un comportamento "indefinito" se inizializziamo un semaforo già inizializzato. Ritorna: 0 per successo. -1 per errore (errno settata) 

**Sequenza di utilizzo:** 
```C
... 
#include <semaphore.h>
... 
sem_t sem;
... 
sem_init(&sem, 0, 0); 
... 
sem_post(&sem); o sem_wait(&sem); (o sem_timedwait()) 
... 
sem_destroy(&sem);
```

 **Semafori POSIX, decremento (operazione di wait):** `#include <semaphore.h> int sem_wait(sem_t *sem);`
 - decrementa il valore di un semaforo (sia esso unnamed o named) 
 - se il semaforo ha il valore zero si blocca in attesa che il valore diventi positivo. Altrimenti decrementa e ritorna 
 - la funzione torna anche se, pur essendo il semaforo a 0, un signal handler la interrompe (cioè il processo riceve un signal gestito). In questo caso torna -1 ed `errno` è settata a `EINTR`

`int sem_timedwait(sem_t *sem, const struct timespec *abs_timeout);`
- come `sem_wait()` ma, nel caso il semaforo sia a 0 e la funzione si blocchi, l'attesa non si protrae oltre il tempo specificato da `abs_timeout` 
- se `sem_timedwait()` torna perché il timeout è spirato prima che il semaforo diventasse positivo allora ritorna -1 ed `errno` viene settata a `ETIMEDOUT` 
- `abs_timeout` punta ad una struttura che specifica il numero di secondi dal 01/01/1970 trascorsi i quali il timeout spira e la chiamata comunque ritorna  

Ritorna: 0 per successo. -1 per errore, errno settata 

**Utilizzo di sem_timedwait(3):** `#include <semaphore.h> int sem_timedwait(sem_t *sem, const struct timespec *abs_timeout);`
- `struct timespec` è definita come segue (`man 3 timespec`): 

```
struct timespec { 
	time_t tv_sec; /* Seconds */ 
	long tv_nsec; /* Nanoseconds [0 .. 999999999]*/ 
};
```
 
> [!example] Esempio di utilizzo di `set_timedwait(3)` con `clock_gettime(2)`, in cui è specificato il numero di secondi dal 01/01/1970 trascorsi i quali la chiamata deve comunque ritornare. : 
> ```C
> struct timespec ts; sem_t sem; 
> ......... 
> clock_gettime(CLOCK_REALTIME, &ts); /* imposta in ts il numero di secondi trascorsi dal 1-1-70 al momento della chiamata */ 
> ts.tv_sec += 8; // imposta un timeout di 8 secondi sem_timedwait(&sem, &ts); /* decrementa il semaforo ma non attende più di 8 secondi */ 
> ```

**Incremento (operazione di post):** `#include <semaphore.h> int sem_post(sem_t *sem);` 
- incrementa il valore di un semaforo (sia esso unnamed o named) 
- di conseguenza il valore del semaforo diventa positivo ed un eventuale processo in attesa su sem_wait() viene «risvegliato» 

Ritorna: 0 per successo, -1 per errore (errno viene settata ed il valore del semaforo rimane invariato)

Esempio con semaforo named:
```C
#include <stdio.h>
#include <pthread.h>
#include <semaphore.h>
#include <fcntl.h>
#include <sys/stat.h>

#define MAX		1000
#define VERDE	1
#define ROSSO	0

sem_t *sem_pari, *sem_dispari;

void pari(void *arg)
{
	int i;
    for (i = 2; i <= MAX; i += 2) {
    	sem_wait(sem_pari);
    	printf("%d\n", i);
    	sem_post(sem_dispari);
    }
    pthread_exit(NULL);
}

void dispari(void *arg)
{
	int i;
    for (i = 1; i <= MAX; i += 2) {
    	sem_wait(sem_dispari);
    	printf("%d\n", i);
    	sem_post(sem_pari);
    }
    pthread_exit(NULL);
}

int main(void)
{
    pthread_t t1, t2;

	/*********************************************************************/
	// La visualizzazione deve partire da 1 e non da 2
    sem_pari    = sem_open("/pari",    O_CREAT|O_RDWR, S_IRUSR|S_IWUSR, ROSSO);
    sem_dispari = sem_open("/dispari", O_CREAT|O_RDWR, S_IRUSR|S_IWUSR, VERDE);
	/*********************************************************************/

    pthread_create(&t1, NULL, (void *)pari, NULL);
    pthread_create(&t2, NULL, (void *)dispari, NULL);
    pthread_join(t1, NULL);
    pthread_join(t2, NULL);

	/*********************/
    sem_close(sem_pari);
    sem_close(sem_dispari);
    sem_unlink("/pari");
    sem_unlink("/dispari");
	/*********************/

    return 0;
}

```

Esempio con semaforo unnamed:
```C
#include <stdio.h>
#include <pthread.h>
#include <semaphore.h>
#include <fcntl.h>
#include <sys/stat.h>

#define MAX		1000
#define VERDE	1
#define ROSSO	0

sem_t sem_pari, sem_dispari;

void pari(void *arg)
{
	int i;
    for (i = 2; i <= MAX; i += 2) {
    	sem_wait(&sem_pari);
    	printf("%d\n", i);
    	sem_post(&sem_dispari);
    }
    pthread_exit(NULL);
}

void dispari(void *arg)
{
	int i;
    for (i = 1; i <= MAX; i += 2) {
    	sem_wait(&sem_dispari);
    	printf("%d\n", i);
    	sem_post(&sem_pari);
    }
    pthread_exit(NULL);
}

int main(void)
{
    pthread_t t1, t2;

	/*********************************************************************/
    sem_init(&sem_pari,    0, ROSSO); // La visualizzazione deve partire
    sem_init(&sem_dispari, 0, VERDE); // da 1 e non da 2
	/*********************************************************************/

    pthread_create(&t1, NULL, (void *)pari, NULL);
    pthread_create(&t2, NULL, (void *)dispari, NULL);
    pthread_join(t1, NULL);
    pthread_join(t2, NULL);

	/*********************/
    sem_destroy(&sem_pari);
    sem_destroy(&sem_dispari);
	/*********************/

    return 0;
}
```

### Signal
Possono essere considerati come IPC (inteprocess communication) alcuni strumenti per la comunicazione tra processi: 
- signals 
- socket 
- pipe 
- shared memory 
- memory-mapped file 
- etc. 

Tipi di signal:
- POSIX Reliable Signals (nel seguito Standard Signals) 
- Linux supporta anche i POSIX Real-Time Signals 
- definiti negli standard POSIX.1-1990 e POSIX.1- 2001 

I signal sono una forma di IPC limitata, nel senso che non contengono ulteriore informazione (il dato scambiato è il signal stesso), Real-Time a parte. Sono simili ad "interrupt software" o a notifiche asincrone, che vengono inviati ad un processo (od a un thread) al verificarsi di un evento. Gli standard hanno un significato intrinseco. In altre parole <u>un signal è una notifica asincrona che il sistema (o un altro processo) invia a un processo/thread per dirgli: “è successo qualcosa”</u>. Non trasporta dati (quasi mai), è solo un **evento**.

> [!quote]
> «un meccanismo asincrono di notifica di eventi tra kernel e processi, che induce una variazione di flusso o di flusso di controllo nel processo destinatario» 

**Cosa accade:** il flusso di esecuzione del processo viene interrotto dal sistema operativo per eseguire la gestione del signal: 
- gestione definita dal processo (se lo ha fatto) 
- gestione di default (altrimenti) 

La gestione può essere varia: il processo termina, viene sospeso, viene eseguita una funzione registrata dal processo all'uopo, il signal viene ignorato, etc. 

> [!info]
> Signals: gestioni di default (da `man 7 signal`) 

**Signal dispositions:** In ogni processo ogni signal ha una **disposition corrente** che determina come il processo stesso si comporta se il signal viene ricevuto (cioè “cosa faccio se arriva questo segnale?“). Se non diversamente specificato, tale disposition è quella di default ed è una delle seguenti: 
- `Term:` Default action is to terminate the process. 
- `Ign`: Default action is to ignore the signal. 
- `Core:` Default action is to terminate the process and dump core (see `core(5)`). 
- `Stop`: Default action is to stop the process. 
- `Cont`: Default action is to continue the process if it is currently stopped. 

La disposition è quindi un attributo per-process ed è condiviso da tutti i suoi threads. È così impostato di default ma è modificabile (con `sigaction()`)

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260424145500.png" width="300">
  <img src="Pasted image 20260424145515.png" width="300">
</div>

**Cambiare la gestione (disposition) corrente:** si fa con `sigaction(2)` o `signal(2)`. La seconda è meno portabile, meglio usare l'altra. Ho tre opzioni: 
- ignorare il signal = `SIG_IGN`
- lasciare/ripristinare la gestione di default del kernel = `SIG_DFL` 
- specificare una funzione handler che gestisca il segnale = `void handler(int sig) {// codice eseguito quando arriva il segnale}`

Un handler è: una funzione chiamata in modo asincrono, in un momento imprevedibile:
- può interrompere qualsiasi punto del programma
- può creare inconsistenze
- deve essere brevissimo

> [!info]
> `SIGKILL` e `SIGSTOP` non possono essere gestiti, ignorati o bloccati

Un handler è una funzione definita dal programmatore che il sistema invoca in maniera asincrona (*callback*) alla ricezione del signal (il processo/programmatore non sa se e quando sarà invocata). Le disposition sono duplicate da `fork(2)` nel child. Dopo `execve(2)` le disposition che prevedono un handler sono resettate al default, quelle che prevedono di ignorare rimangono inalterate

**Chi li invia:** I signal possono essere inviati:
- dal sistema operativo, per comunicare al processo/thread che si è verificato un certo evento (es.: SIGPIPE) 
- dall'utente da terminale premendo certi tasti: 
	- `Ctrl+C`: `SIGINT` 
	- `Ctrl+Z`: `SIGTSTP` 
	- `Ctrl+\`: `SIGQUIT` 

```bash
fabio@fabiopc :~$ stty -a 
speed 38400 baud; rows 46; columns 173; line = 0; intr = ^C; quit = ^\; erase = ^?; kill = ^U; eof = ^D; eol = <undef>; eol2 = <undef>; swtch = <undef>; start = ^Q; stop = ^S; susp = ^Z; rprnt = ^R; werase = ^W; lnext = ^V; ……… 
```

- con il comando kill dalla shell 
- da un processo/thread all'altro con apposite funzioni
- I signals inviati dall'utente dalla tastiera del terminale (`Ctrl+C`, `Ctrl+\`, etc.) sono ricevuti da tutti i processi facenti parte del *process group*

> [!info] Ricorda:
> Un process group:
> - Raggruppa una collezione di processi 
> - tra le altre cose usato per controllare la distribuzione di un signal; quando è diretto al process group viene ricevuto da tutti i processi nel gruppo 
> - di default dopo `fork(2)` il child appartiene allo stesso process group del parent 
> - dopo il richiamo di una delle funzioni exec la nuova immagine del processo appartiene allo stesso process group del chiamante 
> - di default tutti i processi di una singola pipeline appartengono allo stesso process group. Ciò per gestire l’invio di segnali da tastiera, per gestire la sospensione (`SIGTSTP` con Ctrl+Z) e la gestione del background/foreground (comandi interni di Bash bg e fg), etc.

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260424150410.png" width="300">
</div>

**Chi li riceve:** *process-directed* o *thread-directed* 
- process-directed: 
	- inviati all'intero processo 
	- per esempio generati dal kernel per motivi diversi da un'eccezione hardware; o perché inviato con `kill()` o `sigqueue()` 
	- ricevuti da uno dei thread che non li bloccano. Se sono ricevuti da più di uno è il kernel che sceglie a quale thread destinarli 
- thread-directed: 
	- inviati al thread 
	- per esempio perché generati da una specifica istruzione del linguaggio macchina che ha provocato un'eccezione hardware (per es.: segmentation fault, etc.); o inviati con `pthread_kill(), tgkill()`

**Cambiare la disposition: `sigaction(2)`:** `#include <signal.h> int sigaction(int signum, const struct sigaction *act, struct sigaction *oldact);` Usata per variare la gestione di un signal da parte di un processo 
- *signum:* specifica il signal di cui variare la gestione 
- *act:* struct che definisce la nuova gestione 
- *oldact:* struct che riceve le informazioni relative alla gestione precedente 

Ritorna 0 (successo) o -1 (errore; errno settata) 

```C
struct sigaction { 
	void (*sa_handler)(int); 
	void (*sa_sigaction)(int, siginfo_t *, void *); 
	sigset_t sa_mask; 
	int sa_flags; 
	void (*sa_restorer)(void); 
}; 
```

**`sa_handler`:** indica la nuova gestione 
- `SIG_IGN`: ignorato 
- `SIG_DFL`: funzione di default 
- puntatore a funzione dell'handler. Tale funzione riceve come unico argomento il numero del signal ricevuto 

`SIGKILL` e `SIGSTOP` non possono essere gestiti, ignorati o bloccati 

> [!info] N.B.
> Non assegnare contemporaneamente sa_handler e sa_sigaction (in certe architetture la struct è oggetto di union). In questa sede non parliamo di sa_sigaction. Per maggiori informazioni vedere `man sigaction(2)` 

**Registrare Signal Handlers: `sigaction()`:** 

```C
struct sigaction { 
	void (*sa_handler)(int); 
	void (*sa_sigaction)(int, siginfo_t *, void *); 
	sigset_t sa_mask; 
	int sa_flags; 
	void (*sa_restorer)(void); 
};
``` 

**`sa_mask`:** segnali da mascherare durante l'esecuzione dell'handler (in aggiunta a quello che ne ha provocato la chiamata che viene automaticamente bloccato) 

**Quando è eseguito un handler:** 
> [!quote] da `man 7 signal`: 
> "Whenever there is a transition from kernel-mode to user-mode execution (e.g., on return from a system call or scheduling of a thread onto the CPU), the kernel checks whether there is a pending 1 unblocked 2 signal for which the process has established a signal handler. If there is such a pending signal, the following steps occur:"

**Passi per l'esecuzione dell'handler:** 
- il kernel esegue operazioni preparatorie per l'esecuzione dell'handler: 
	- rimuove il signal dalla lista dei pendenti 
	- salva lo stato del processo/thread da interrompere (i registri, etc.) 
- l'handler viene eseguito in user-mode 
- quando l'handler ritorna, lo stato del processo/thread viene ripristinato 
- il processo/thread che ha ricevuto il signal continua la sua esecuzione 

> [!info] Nota: 
> Dal punto di vista del kernel non c'è differenza tra eseguire l'handler o il codice del processo/threads. Il kernel non tiene traccia del fatto che il codice oggetto di esecuzione sia o meno in un handler, e le informazioni sullo stato del processo sono salvate nei registri e nello stack dello «user-space» (non nelle aree di memoria del kernel). Per maggiori dettagli vedi `man 7 signal` alla sezione "Execution of signal handlers "

> [!example] 
> ![[Screenshot From 2026-04-24 15-18-56.png|200]]
> 
> ![[Pasted image 20260424160112.png|200]]
> 
> nell'esempio 2 sopra descritto la open(2) si interrompe e torna -1 (errno viene settata a EINTR)

**Interruzione di SC bloccate:** se un signal handler è invocato mentre una syscall è bloccata in determinate operazioni allora può verificarsi: 
- la funzione è automaticamente riavviata 
- la funzione fallisce e torna `EINTR` 

Ciò dipende da: 
- se in fase di creazione dell'handler (sigaction) era stato o meno specificato il flag `SA_RESTART` 
- l'interfaccia oggetto della chiamata

Tale comportamento può variare tra sistemi Unix diversi. Alcuni dettagli per Linux: esempi di syscall che se bloccate vengono restartate se `SA_RESTART` è attivo: 
- `read(2), write(2), ioctl(2)`, etc. su "slow devices", cioè devices in cui una chiamata I/O può bloccarsi indefinitamente, come pipe, terminale, socket 
- `open(2)`, per esempio su una FIFO 
- `wait(2)`, `waitpid(2)` 
- `sem_wait(3)`, `sem_timedwait(3)` 

Considerazioni per gli handlers:
- se più segnali sono pendenti l'ordine in cui sono inviati non è specificato 
- gli standard signal non si accumulano. Se più istanze dello stesso segnale arrivano mentre esso è bloccato solo un'istanza verrà inviata 
- durante l'esecuzione dell'handler possono arrivare altri segnali. Questo può generare inconsistenze
- come visto prima, quando (e se) l'handler ritorna il programma riprende da dove aveva lasciato

Siamo sicuri che tutto sia consistente?
- l'handler dovrebbe essere breve e fare il meno possibile 
- non tutte le funzioni di libreria possono essere chiamate nell'handler senza che possa accadere qualcosa di imprevisto 
- POSIX.1 definisce un elenco di funzioni che DEBBONO essere async-signal-safe (vedi `man 7 signal-safety`), cioè che possono essere invocate «safely» in un handler 
- i.e: `printf()` (chiamata nell'esempio di prima) non è in elenco. Non c'è garanzia che non si verifichino comportamenti non previsti

**Signal mask:** 
- un signal può essere bloccato, cioè non sarà ricevuto fino a che non verrà sbloccato 
- la signal mask definisce quali signal sono bloccati e quali no 
- una signal mask per ogni thread (per-thread) 
- `SIGKILL` e `SIGSTOP` non possono essere gestiti, ignorati o bloccati 
- ereditata dai child dopo `fork(2)`; preservata dopo `execve(2)`

**Segnali pendenti:** 
- nel tempo tra la generazione e l'invio al processo il signal è detto pendente 
- può essere pendente perché bloccato (da signal mask)
- un thread può ottenere la lista dei segnali pendenti con `sigpending()` 
- la lista dei segnali pendenti è resettata (vuota) dopo `fork(2)` ed è preservata dopo `execve(2)`

**Funzioni per bloccare signals:** 
- `sigprocmask()`: per processi single-threaded 
- `pthread_sigmask()`: per processi multithreaded

> [!info] Ricorda:
> - Il tuo programma sta eseguendo normalmente
> - Arriva un evento (Ctrl+C, errore, kill, ecc.)
> - Il kernel **interrompe temporaneamente** il tuo flusso
> - Esegue una gestione del segnale
> - Poi il programma **riprende da dove era**
> 
> È molto simile a un **interrupt**, ma a livello software. Nello specifico, quando arriva un signal: 
> 1. Il kernel salva lo stato del processo (registri, stack)
> 2. Decide **come gestire il segnale**
> 3. Esegue:
>     - handler (se definito)
>     - oppure azione di default
> 4. Finito l’handler → ripristina lo stato
> 5. Il programma riprende

### Socket
Astrazione software per la comunicazione attraverso API standard e condivisa: 
- networking socket layer user interface 
- "uniform interface between the user process and the network protocol stacks in the kernel" (man 7 socket) 

<u>Comunicazione IPC anche attraverso una rete: scambio di pacchetti di dati ed altre operazioni</u>. Esistono diversi communication domain, che stabiliscono la "protocol family" (`man 7 socket`), e diversi tipi di socket, che stabiliscono la modalità di comunicazione nello scambio di dati. Sono stati introdotti nel 1983 in Unix BSD 4.2. Le system calls per i socket sono più complesse rispetto, per esempio, alle pipe; le operazioni realizzate da `mkfifo(3)` ed `open(2)` per le named pipe sono divise tra più system call 

**Socket domains (`man 2 socket`):**  

| Name              | Purpose                        | Man page     |
| ----------------- | ------------------------------ | ------------ |
| AF_UNIX, AF_LOCAL | Local comunication             | `unix(7)`    |
| AF_INET           | IPv4 internet protocols        | `ip(7)`      |
| AF_INET 6         | IPv6 internet protocols        | `ipv6(7)`    |
| AF_IPX            | IPX-novell protocols           |              |
| AF_NETLINK        | Kernel user interface dev      | `netlink(7)` |
| AF_X25            | ITU-T X.25/ISO-8208 protocol   | `x25(7)`     |
| AF_AX25           | Amateur radio AX.25 protocol   |              |
| AF_APPLETALK      | Apple Talk                     | `ddp(7)`     |
| AF_PACKET         | Low level packet interface     | `packet(7)`  |
| AF_ALG            | Interface to kernel crypto API |              |
**Socket type (`man 2 socket`):** 
- **SOCK_STREAM:** Provides sequenced, reliable, twoway, connection-based byte streams. An out-of-band data transmission mechanism may be supported. 
- **SOCK_DGRAM:** Supports datagrams (connectionless, unreliable messages of a fixed maximum length). 
- **SOCK_SEQPACKET:** Provides a sequenced, reliable, two-way connectionbased data transmission path for datagrams of fixed maximum length; a consumer is required to read an entire packet with each input system call. 
- **SOCK_RAW:** Provides raw network protocol access. 
- **SOCK_RDM:** Provides a reliable datagram layer that does not guarantee ordering. 
- **SOCK_PACKET:** Obsolete and should not be used in new programs; see packet(7)

> [!info] N.B.: 
> Some socket types may not be implemented by all protocol families

**Socket: passi per il collegamento:** 

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260505135541.png" width="300">
</div> 

- Il server chiama `socket(2)`, che crea un "endpoint" per la comunicazione e ritorna un file descriptor (*fd_s*) per le successive call 
- Il server chiama `bind(2)` per associare il socket con un "indirizzo/nome" (vari tipi a seconda del domain). Storicamente tale operazione è anche detta "assegnare un nome al socket (naming)" ("My_Sock", un indirizzo IP e una porta su cui stare in ascolto, etc.) 
- Il server chiama `listen(2)` per specificare che sul socket creato sono accettate connessioni in ingresso 
- Con `accept(2)` il server rimane in attesa di connessioni da parte del/i client 
- Anche il client crea un socket (endpoint) con `socket(2)` e riceve un file descriptor (fd_c) 
- Con `connect(2)` il client si collega col socket del server usando l'indirizzo/nome associato a fd_s con `bind(2)` 
- A connessione avvenuta il server crea un nuovo socket (*fd_new*) 
- *fd_new* è ritornato al server dalla `accept(2)` ed è utilizzato per la comunicazione col client. Il server può chiamare nuovamente `accept(2)` per rimettersi in attesa di altre connessioni 

**Creare un socket (sia lato server che client):** 
```C
#include <sys/socket.h>
int socket(int domain, int type, int protocol); 

domain: uno dei domini (AF_UNIX, AF_INET, etc.) 
type: in generale uno dei tipi di socket (SOCK_STREAM, SOCK_DGRAM, etc.) 
protocol:in molti casi ce n'è uno per ogni coppia family/type; in questo caso possiamo passare il valore 0 

Ritorna il file descriptor per il nuovo socket; -1 per errore e setta errno 
```

> [!info] N.B.
> Né POSIX.1 né Linux richiedono di includere anche sys/types.h ma in altre implementazioni potrebbe essere richiesto 

**Chiamate lato server successive a `socket(2)`** 
- `bind(2)` 
- `listen(2) `
- `accept(2)` 

**`bind(2)`: assegnare un nome ad un socket(1):** 

```C
#include <sys/socket.h> 
int bind(int sockfd, const struct sockaddr *addr, socklen_t addrlen); 

sockfd: file descriptor del socket 
addr : indirizzo/nome del socket 
addrlen: lunghezza dell'indirizzo/nome 

Ritorna 0 per successo; -1 per errore e setta errno 
```

Le regole per il name binding variano a seconda del domain. Consultare le seguenti pagine del man per i dettagli: 
- For `AF_INET, see ip(7)` 
- For `AF_INET 6, see ipv 6(7)`
- For `AF_UNIX, see unix(7)` 
- For `AF_APPLETALK, see ddp(7)` 
- For `AF_PACKET, see packet(7)` 
- For `AF_X 25, see x 25(7)` 
- For `AF_NETLINK, see netlink(7)`

La struttura `sockaddr` è definita come segue: `struct sockaddr { sa_family_t sa_family; char sa_data[14];}`. `struct sockaddr` è un tipo di dato "generico" che ha il solo ruolo di evitare una warning del compilatore

Come secondo argomento occorre passare il puntatore alla struttura specifica per la famiglia (domain) di socket (`man 2 socket`). Per tale puntatore si fa il cast a `struct sockaddr *`, e ciò evita la warning.  

> [!quote] Da `man 7 socket`:
> The purpose of this type is purely to allow casting of domain-specific socket address types to a "generic" type, so as to avoid compiler warnings about type mismatches in calls to the sockets API.

**`bind(2)`: assegnare un nome ad un socket `AF_UNIX`:** Un UNIX domain socket address è rappresentato nella seguente struttura: 
```
struct sockaddr_un 
{ 
sa_family_t sun_family; /* AF_UNIX */ 
char sun_path[108]; /* pathname */ 
}; 
```

`sun_family` contiene sempre `AF_UNIX`. `sun_path` è il pathname del socket 

**`bind(2)`: assegnare un nome ad un socket `AF_INET`:**

```C
#include <sys/socket.h>
#include <netinet/in.h>
#include <netinet/ip.h>/* superset of previous two */ 

tcp_socket = socket(AF_INET, SOCK_STREAM, 0); 
udp_socket = socket(AF_INET, SOCK_DGRAM, 0); 

int bind(int sockfd, const struct sockaddr *addr, socklen_t addrlen); 

//addr: combinazione di indirizzo IP e porta di connessione: 

struct sockaddr_in 
{ 
sa_family_t sin_family; /* address family: AF_INET */ 
in_port_t sin_port; /* port in network byte order */ 
struct in_addr sin_addr; /* internet address */ 
}; 

/* Internet address. */ 
struct in_addr 
{ 
uint 32_t s_addr; /* address in network byte order */ 
}; 
```

**Endianness: host e network byte order:** 
- *endianness:* ordine dei byte nella word 
	- **big endian:** il byte più significativo è all'indirizzo di memoria minore 
	- **little endian:** il byte più significativo è all'indirizzo di memoria maggiore 
- *host byte order:* come ordina la macchina. Little-endian 
- *network byte order:* usato dai protocolli TCP/IP. Corrisponde a big-endian 

**host e network byte order:** 
```C
#include <arpa/inet.h>
uint 32_t htonl(uint 32_t hostlong); 
uint 16_t htons(uint 16_t hostshort); 
uint 32_t ntohl(uint 32_t netlong); 
uint 16_t ntohs(uint 16_t netshort); 
```

Convertono da host byte order (little) a network byte order (big) e viceversa. Sono utilizzate, per esempio, per specificare un numero di porta (`man 3 byteorder`) 

**IPv 4 "numbers-and-dots" e network byte order:** 

```C
#include <arpa/inet.h>
#include <sys/socket.h>
#include <netinet/in.h>

int inet_aton(const char *cp, struct in_addr *inp);// 127.0.0.1 => net b. order 
in_addr_t inet_addr(const char *cp); // 127.0.0.1 => net b. order 
in_addr_t inet_network(const char *cp); // 127.0.0.1 => host b. order 
char *inet_ntoa(struct in_addr in); // net b. order => 127.0.0.1 
struct in_addr inet_makeaddr(in_addr_t net, in_addr_t host); 
in_addr_t inet_lnaof(struct in_addr in); 
in_addr_t inet_netof(struct in_addr in); 
```

Sono funzioni che convertono nei diversi formati indirizzi IP: da IPv 4 "numbers-and-dots" (es.: 127.0.0.1) a intero e viceversa (`man 3 inet_aton`)

**`listen(2)`: restare in ascolto:** 

```C
#include <sys/types.h>
#include <sys/socket.h>

int listen(int sockfd, int backlog); 
``` 

- *sockfd* è il file descriptor tornato da `socket(2)`
- "marca" il socket come "passive", cioè un socket che riceverà connessioni in ingresso tramite `accept(2)` 
- *backlog* definisce la lunghezza massima fino alla quale può crescere la coda di connessioni in sospeso 
- se il *backlog* specificato come argomento è maggiore di /proc/sys/net/core/somaxconn esso viene ridotto a tale valore 

**`accept(2)`: accettare connessioni:** 

```C
#include <sys/types.h>
#include <sys/socket.h>

int accept(int sockfd, struct sockaddr *addr, socklen_t *addrlen);
``` 

- usata con connection-based socket types (`SOCK_STREAM`, `SOCK_SEQPACKET`). 
- esegue le seguenti operazioni: 
	- estrae la prima delle connessioni in coda per il socket *sockfd* in "listening" (cioè tornato da `socket(2)` e che stato oggetto di `bind(2)` e `listen(2)`)
	- crea un nuovo socket e ne ritorna il file descriptor 
	- il file descriptor ritornato può essere utilizzato lato server per la comunicazione col client che aveva richiesto la connessione 
- il processo può chiamare di nuovo `accept(2)` per accettare nuove connessioni 
- il socket referenziato da *sockfd* non è modificato 

 ```C
 #include <sys/types.h>
 #include <sys/socket.h>
  
 int accept(int sockfd, struct sockaddr *addr, socklen_t *addrlen);
 ``` 
 
- `addr` punta ad una struttura `sockaddr`. Riceve l'indirizzo/naming del socket dall'altro lato della comunicazione. Se NULL non vengono memorizzati dati 
- se non ci sono connessioni entranti in coda e il socket NON è configurato in modalità "non bloccante" `accept()` si blocca finché non riceve connessioni. Se il socket è in modalità "non bloccante" e non ci sono connessioni in coda `accept()` fallisce ed errno è settata a `EAGAIN` or `EWOULDBLOCK`. 

**`connect(2)`: connessione lato client:** 

```C
#include <sys/types.h>
#include <sys/socket.h>

int connect(int sockfd, const struct sockaddr *addr, socklen_t addrlen);
``` 

- connette il socket referenziato da *sockfd* con l'indirizzo specificato da `addr` 
- `addr` è lo stesso indirizzo/nome usato nella `bind(2)` dal lato server 
- al ritorno *sockfd* è pronto per scambiare dati (se non ci sono stati errori ovviamente) 

**Socket: inviare e ricevere dati:** una volta stabilita la connessione ci sono diverse funzioni per inviare e ricevere dati attraverso i sockets 
- Invio: `send(2), sendto(2) e sendmsg(2)` 
- Ricezione: `recv(2), recvfrom(2) e recvmsg(2)` 
- in aggiunta possiamo utilizzare le system calls viste per i file: `write(2), read(2)`, etc. 
- non si possono fare operazioni di seek che non siano alla posizione 0 

```C
#include <sys/socket.h>

ssize_t send(int sockfd, const void buf[.len], size_t len, int flags); 
```

Utilizzabile solo se il socket è nello stato di connesso. La sola differenza tra `send(2)` e `write(2)` è la presenza dell'argomento flags. Se flags è 0 `send(2)` è equivalente a `write(2)` (per ulteriori informazioni si rimanda al man delle specifiche funzioni, in particolare per i flags di `send(2)`)

```C
#include <sys/socket.h>
ssize_t recv(int sockfd, const void buf[.len], size_t len, int flags); 
```

La sola differenza tra `recv(2)` e `read(2)` è la presenza dell'argomento flags. Se flags è 0, `recv(2)` è "generalmente" * equivalente a `read(2)` (per ulteriori informazioni si rimanda al man delle specifiche funzioni, in particolare per i flags di `recv(2)`) 

> [!info]
> Differenza di comportamento segnalata nelle "Notes" di `man 2 recv`: "If a zero-length datagram is pending, read(2) and recv() with a flags argument of zero provide different behavior. In this circumstance, read(2) has no effect (the datagram remains pending), while recv() consumes the pending datagram." 

### Shared memory
Esistono 2 possibili implementazioni: POSIX e SYSTEM V:
- **POSIX:** interfaccia più semplice e meglio congegnata; ma meno diffusa, specialmente in sistemi vecchi (man 7 shm_overview) 
- **SYSTEM V:** API più datata ma presente anche in sistemi molto vecchi (POSIX potrebbe non esserci) (man 7 sysvipc) 

Si tratta di processi che condividono un'area di memoria e, conseguentemente, attraverso essa possono comunicare/sincronizzarsi. Tipicamente i processi devono sincronizzare l'accesso alla memoria condivisa, per esempio con dei semafori (linkare con `–lrt` (librt) se necessario). Hanno persistenza nel kernel, le aree non deallocate permangono sino al prossimo shutdown 

In Linux sono visibili nel filesystem virtuale `/dev/shm`. Possiamo impostare permessi 

**Interfacce dell'API POSIX (da `man 7 shm_overview)`:**
- `shm_open(3)`: Create and open a new object, or open an existing object. This is analogous to `open(2)`. The call returns a file descriptor for use by the other interfaces listed below. 
- `ftruncate(2)`: Set the size of the shared memory object. (A newly created shared memory object has a length of zero.) 
- `mmap(2)`: Map the shared memory object into the virtual address space of the calling process. 
- `munmap(2)`: Unmap the shared memory object from the virtual address space of the calling process. 
- `fstat(2)`: Obtain a stat structure that describes the shared memory object. Among the information returned by this call are the object's `size (st_size)`, permissions (`st_mode), owner (st_uid)`, and group `(st_gid)`. 
- `close(2)`: Close the file descriptor allocated by `shm_open(3)` when it is no longer needed. 
- `shm_unlink(3)`: Remove a shared memory object name. 
- `fchown(2)`: To change the ownership of a shared memory object. 
- `fchmod(2)`: To change the permissions of a shared memory object. 

**`shm_open`:**

```C
#include <sys/mman.h>
#include <sys/stat.h> /* For mode constants */ 
#include <fcntl.h> /* For O_* constants */ 

int shm_open(const char *name, int oflag, mode_t mode); 
```

- analoga ad `open(2)` 
- apre una shared memory (o la crea se non esiste e si specifica il flag `O_CREAT`) 
- torna un file descriptor (o -1 per errore settando errno) 
- `name` identifica la shared memory. Per portabilità dovrebbe essere nella forma "/somename" 
- `oflag` devono comprendere uno tra `O_RDONLY` e `O_RDWR`. Tra i possibili oflag vi sono `O_CREAT` ed `O_TRUNC` 
- `mode` specifica i permessi, come in `open(2)`

Un'area di nuova creazione ha dimensione iniziale 0 byte, la dimensione può essere variata con `ftruncate(2)`. I bytes di una nuova area sono inizializzati a 0 (cod. ASCII) 

**`shm_unlink`:**

```C
#include <sys/mman.h>
#include <sys/stat.h> /* For mode constants */ 
#include <fcntl.h> /* For O_* constants */ 

int shm_unlink(const char *name);
``` 

- analoga ad `unlink(2)` 
- rimuove un'area di memoria creata precedentemente con `shm_open(2)`. 
- Più precisamente rimuove "il nome" dall'elenco degli oggetti "shared memory" e, quando tutti i processi hanno fatto l'unmapping dell'area, dealloca e distrugge il suo contenuto 
- in altre parole, si può richiamare `shm_unlink()` anche quando ancora vi siano processi che usano l'area 

**`ftruncate`:**

```C
#include <unistd.h>

int ftruncate(int fd, off_t length);
```

 - in generale tronca un file alla dimensione specificata da length 
 - nel nostro caso stabilisce la dimensione della shared memory (considerate che `shm_open(2)` torna un file descriptor…) 
 - torna 0 per successo, -1 per errore (errno settata) 
 
```C
#include <sys/mman.h>

void *mmap(void *addr, size_t length, int prot, int flags, int fd, off_t offset);
``` 

- in generale mappa file o device nello spazio di indirizzamento virtuale del chiamante. Nel nostro caso mappa (anche) shared memory (cioè collega la shared memory allo spazio di indirizzi del processo)
- torna un puntatore all'area mappata 
- `addr`: indirizzo iniziale nello spazio di indirizzamento per la mappatura. Se NULL, sceglie il kernel dove mappare (più portabile) 
- `length`: lunghezza dell'area mappata 
- `offset`: l'area mappata è inizializzata col contenuto della shared memory, referenziata da fd, a partire da offset per lunghezza `length` 

> [!info] Ricorda:
> `offset` si riferisce alla posizione nel/nella file/shared memory mappato/a. Deve essere un multiplo di "page size", che viene ritornato (cioè può essere ottenuto) da `sysconf(_SC_PAGESIZE)`

- `prot`: indica i permessi sull'area (non deve essere in conflitto con il modo di apertura della shared memory). Può essere: 
	- `PROT_NONE`: l'area non può essere acceduta oppure una combinazione di: 
	- `PROT_EXEC`: l'area ha il permesso di esecuzione 
	- `PROT_READ`: l'area può essere letta 
	- `PROT_WRITE`: l'area può essere scritta 

- `flags`: determinano, tra l'altro, se le modifiche all'area mappata sono o meno visibili agli altri processi che la usano (e nel caso di file, se le modifiche sono trasferite al file sottostante): 
	- `MAP_SHARED`: le modifiche sono visibili agli altri processi 
	- `MAP_PRIVATE`: modifiche non visibili agli altri processi 
	- `MAP_SHARED_VALIDATE`: è un'estensione di Linux (i due precedenti sono definiti da POSIX). Come `MAP_SHARED` ma, mentre con `MAP_SHARED` eventuali altri flags "sconosciuti" sono ignorati, con `MAP_SHARED_VALIDATE` il kernel controlla che gli altri flags specificati siano validi 
	- ce ne sono vari altri specifici di Linux (no POSIX) che sono descritti nella sezione "The flags argument" di `man 2 mmap`
	
- `fd`: il file descriptor tornato da `shm_open(2)` 
- dopo che `mmap(2)` è ritornata il file descriptor `fd` può essere subito chiuso senza invalidare il mapping 
- ritorna un puntatore all'area mappata. In caso di errore torna `MAP_FAILED` (cioè `(void *) -1`) e setta errno 
- il mapping è preservato dopo `fork(2)`
