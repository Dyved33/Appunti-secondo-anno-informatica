###### Cenni Storici e Famiglie di Sistemi Operativi
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
###### Concetti Base dei Sistemi Unix
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
###### Primi Comandi BASH e variabili:
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
| **`df`**       | Elenca i filesystem attualmente montati.                                                                                |
| **`cat`**      | Visualizza per intero il contenuto di un file testo.                                                                     |
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

 **Differenza nella creazione dei Link (`ln`):** **Hard Link** (`ln pippo link1`): Crea una vera e propria copia che occupa spazio fisico sul disco. Se il file originario "pippo" viene cancellato, il file "link1" e il suo contenuto rimangono inalterati.
* **Symbolic Link** (`ln -s pippo link1`): Crea un collegamento simbolico (occupazione minimale su disco). "link1" contiene solo una stringa col percorso di "pippo". Se il file "pippo" originario viene cancellato, "link1" diventa semplicemente un collegamento non valido o rotto.

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
###### Attributi del file/cartella (bash)
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
   1. *mode simbolico* = chmod (chi) (operatore) (permessi) file. Si riferiscono a: **chi può cambiare i permessi** (u user, g group, o others, a all), gli **operatori** (+ aggiunge, - toglie, = imposta esattamente) e ai **permessi** (r read, w write, x execute).

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
###### System call, wrappers e accesso ai file
Le system calls sono l'interfaccia per il programmatore ai servizi del sistema operativo (chiamata diretta al sistema). Vi si accede generalmente attraverso un API/linguaggi che la supportino. Ce ne sono vari tipi (accesso ai file, gestione di processi, IPC, semafori,...) e richiedono un "cambio di contesto": da *kernel mode* a *user mode* e viceversa.
- Kernel mode = accesso senza limitazioni alle risorse HW/SW del sistema (memoria, dispositivi I/O, registri, codice del SO,...). Un errore in kernel mode può provocare un arresto del sistema
- User mode = accesso limitato alle risorse di sistema (es. Memoria). Un errore in user mode può provocare il blocco del solo processo utente

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
| 'rmdir' | Cancella una directory (se vuota) |
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
> - i file aperti andrebbero chiusi prima di terminare il programma con `#include int close(int fd)`, dove fd è il file descriptor ritornato da open () 
> - l'apertura di un file crea una nuova entry nella tabella di sistema dei file aperti. Quando un file descriptor viene duplicato (con dup2 o simili), questo si riferisce alla stessa *open file description* del descrittore del file originale $\implies$  i due file descriptors condividono lo stesso *file offset* e le stesse *file status flags*
> - su Linux si può usare l'operazione kcmp (2) KCMP_FILE per testare se due file descriptors si riferiscono alla stessa open file description
 
> [!tip] 
> - Errno = è una libreria di C (errno.h) che mi dice precisamente quale sia l'errore in corso 
> - con `strace -c ls` posso vedere le system call attive e alcuni dati che le riguardano: percentuale del tempo totale con cui il kernel è stato dentro alla syscall, tempo totale speso in quel tipo di syscall, media per ogni call, numero di volte che la syscall è stata eseguita, quante syscall hanno ritornato un errore, nome della syscall
> - **buffering** = in generale il buffer è una certa quantità di memoria che deve essere riempita prima che accada una qualche chiamata o azione specifica (ad esempio se scrivo una stringa questa potrebbe non essere visualizzata subito perché la memoria buffer non è ancora piena). Esistono 3 tipi di buffer: interattivo (buffer si svuota quando appare `\n`), di pipe o di file (buffer si svuota solo quando la memoria predisposta è piena) e l'assenza di buffer (ogni carattere fa si che la memoria buffer si svuoti subito)

> [!info] fd = è un intero piccolo che il  SO usa per identificare un file aperto dentro ad un processo (non mi riferisco più a file. txt, ma al valore assunto da fd. )
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

**write:** `ssize_t write(int fd, const void *buf, size_t count);` <u>Scrive i primi count (numero di byte da scrivere) byte contenuti in buf (puntatore alla memoria che contiene i dati da scrivere) nel file aperto identificato da fd </u>(file descriptor ritornato da open(2)). Ritorna il numero di byte scritti (che può anche essere inferiore a count se, per esempio, il device di scrittura si è riempito) o -1 in caso di errore (se errno settata). 

**read:** `ssize_t read(int fd, void *buf, size_t count);` <u>Legge count (numero massimo di byte da leggere) byte dal file identificato da fd (file descriptor ritornato da open(2)) e li memorizza in buf (buffer dove mette i dati letti)</u>. Ritorna il numero di byte letti (che può anche essere inferiore a count se c'erano a disposizione nel device di lettura meno byte) o -1 in caso di errore (errno settata). Il valore di ritorno 0 significa EOF (end of file).

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

###### Standard I/O library
È un insieme di funzioni per accesso ai file, agli *streams* standard, etc che astraggono le operazioni di input-output sotto forma di streams. Le operazioni sono bufferizzate in spazio utente e ciò riduce il numero di chiamate alle sottostanti syscall e i relativi cambi di contesto ottimizzando le prestazioni. Forniscono un'interfaccia verso il SO di più alto livello rispetto alle syscall

| Syscall (no bufferizzazione nello spazio utante)         | StdIO library (bufferizzazione nello spazio utente)                                                           |
| -------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------- |
| `open`                                                   | `fopen`                                                                                                       |
| `close`                                                  | `fclose`                                                                                                      |
| `read/write`                                             | `getchar/putchar, getc/putc, fgetc/fputc, fread/fwrite, gets/puts, fgets/fputs, scanf/printf, fscanf/fprintf` |
| `lseek`                                                  | `fseek`                                                                                                       |
| Richiedono cambio contesto user-kernel ad ogni carattere | Tipicamente più efficaci e più veloci                                                                         |
###### Bash: standard stream e redirezioni
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
> [!info] Il comando `cat` legge sempre e solo da stdin (fd 0)

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
###### System calls per i processi
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
**PID e PPID:** `include <unistd.h> pid_t getpid(void); pid_t getppid(void)`

- getpid() torna il process ID
- getppid() torna il parent process ID 

Queste due funzioni hanno sempre successo

> [!info] Esiste anche `getpgrp()` che ritorna il *process group ID*

**wait () e waitPID():** `#include <sys/wait.h> pid_t wait(int *wstatus); pid_t waitpid(pid_t pid, int *wstatus, int options)`

- *wait()* sospende il chiamante finché un *child* termina
- *waitpid()* sospende il chiamante finché un *child* cambia stato (terminato; stoppato da un *signal*; ha ripreso dopo un *signal*), dipende da *options*

Permettono di rilasciare le risorse assegnate al *child*. Se un *child* termina e una *wait* non viene eseguita il *child* rimane nello stato di *zombie* (*man wait(2)*). Questi metodi ritornano il PID del processo il cui stato è cambiato (*waitpid()* può tornare 0 in certe circostanze) oppure -1 in caso di errore. Se un/il *child* è terminato/ha già cambiato stato tornano immediatamente. Tra i valori che possiamo dare a PID ci sono i seguenti:
- -1: attende per un cambiamento di stato di uno qualsiasi dei child
- >0: attende per un cambiamento di stato del child il cui process ID è uguale a PID
- Etc

*waitpid(-1, &wstatus, 0)* è equivalente a *wait(&wstatus)*. Se *&wstatus* non è NULL riceve informazioni sul *child*. Queste informazioni possono essere evinte con delle macro:
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
###### File Bash
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
- PATH: contiene una lista di directory, separate da ':', in cui la BASH cerca i comandi da eseguire (se non viene specificato un path per essi) 
- HOSTNAME: nome host del sistema 
- USER: nome dell'utente 
- PWD: directory corrente

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

**Scripting:** 
- Possiamo passare ad uno script degli argomenti

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
- Chiamano la sottostante execve(2) 
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
- Termina il chiamante "immediately"
- I file descriptors sono chiusi
- Il parent ID di eventuali child diventa init (o il più vicino "subreaper process"...)

Non c'è ritorno, ma il processo parent può ottenere il byte meno significativo di status (status & 0xFF) con wait(2)

**_ exit(2) vs exit(3):** `#include void exit(int status);`
- Termina il chiamante
- *Exit(3)* fa tutto ciò che fa *_ exit(2)*, ma inoltre:
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

###### Pipe
**IPC:** vari strumenti per la comunicazione tra processi:
- Signal
- Socket
- Pipe
- Semafori
- Shared memory
- Memory-mapped file
- Message queue
- Etc

Le pipe sono di solito unidirezionali (in certi sistemi sono bidirezionali come in Windows), sono uno dei primi meccanismi di IPC e sono disponibili praticamente su tutti i sistemi. Connettono due o più processi tramite un canale di comunicazione. Hanno un'estremità per la lettura (*read end*) ed una per la scrittura (*write end*). Sono caratterizzate da comunicazione non strutturata (*byte stream*), cioè **le pipe trasmettono un flusso continuo di byte grezzi senza alcuna informazione di controllo, formattazione o metadati** associati al contenuto dei dati stessi

Ne esistono 2 tipi:
1. Senza nome (*unnamed*): sono crate con la system call `pipe()`, i processi comunicanti condividono file descriptor per la lettura e la scrittura (`fork(2)`). Solo tra processi imparentati (padre-figlio)
2. Con nome (*named o FIFO*): sono un file speciale (type `p` in `ls -l`) visibile nel filesystem (ma non occupano né implicano operazioni in esso).

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
- `pipefd[0]` file descriptor per la read end della pipe 
- `pipefd[1]` file descriptor per la write end della pipe 
- `Ritorno : 0` se tutto ok. -1 per errore (errno settata)

Creata la pipe è possibile leggere e scrivere dati, per es.: con `read()` e `write()`. Una sequenza tipica di utilizzo è:
- il parent crea la pipe (pipe()) 
- il parent crea un child (fork()). I file descriptor ritornati dalla pipe() sono duplicati nel child 
- lo scrivente chiude il fd di lettura* 
- il lettore chiude il fd di scrittura* 
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
- Se la pipe è piena ulteriori write(2) si bloccheranno (o falliranno se `O_NONBLOCK` è settato) 
- POSIX stabilisce che le `write(2)` di ampiezza minore di `PIPE_BUF` sono garantite essere atomiche (ma possono essere bloccanti se la pipe è piena) 
- Stabilisce anche che `PIPE_BUF` deve essere almeno 512 byte
- se tutti i descrittori di lettura sono stati chiusi una write provoca l’invio di un segnale `SIGPIPE` al processo: 
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

