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
| **`df`**       | Elenca i file system attualmente montati.                                                                                |
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
| -       | File "regolare"    |
| d       | Cartella           |
| b       | Device a blocchi   |
| c       | Device a caratteri |
| l       | Link simbolico     |
| p       | Pipe               |
| s       | Socket             |
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

| Nome                | Cosa fa                                                                        |
| ------------------- | ------------------------------------------------------------------------------ |
| open                | Apre un file in lettura e/o scrittura o crea un nuovo file                     |
| close               | Chiude un file precedentemente aperto                                          |
| read                | Legge da un file                                                               |
| write               | Scrive su un file                                                              |
| lseek               | Sposta il puntatore di lettura/scrittura ad una determinata posizione nel file |
| unlink              | Rimuove una entry di un file nel filesystem                                    |
| stat (fstat, Istat) | Fornisce informazioni sugli attributi di un file                               |
| dup, dup2, dup3     | Duplicano un file descriptor                                                   |
**System calls per le directory:**

| Nome  | Cosa fa                           |
| ----- | --------------------------------- |
| mkdir | Crea una directory                |
| rmdir | Cancella una directory (se vuota) |
| chdir | Cambia la directory corrente      |
**Open:** ` #include <unistd.h> int open(const char *pathname, int flags, mode_t mode)`
- Apre il file specificato da *pathname*
- *flags* deve includere uno tra `O_RDONLY, O_WRONLY o O_RDWR`
- *mode* specifica i permessi da applicare al file in caso di creazione (permessi che vengono stabiliti anche in base a `UMASK`). La creazione può essere richiesta specificando `O_CREAT` tra i flags, altrimenti mode viene ignorato e può essere omesso

| Altri flags (separati da \|) | Funzionamento                                                                              |
| ---------------------------- | ------------------------------------------------------------------------------------------ |
| O_CREAT                      | Se il file no esiste lo crea                                                               |
| O_APPEND                     | Il file è aperto in accodo per la scrittura                                                |
| O_TRUNC                      | Se il file esiste il contenuto viene azzerato                                              |
| O_EXCL                       | Se usato insieme a O_CREAT e il file esiste, open () fallisce ed errno è settata ad EEXIST |
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

| Comando      | Contenuto                                              |
| ------------ | ------------------------------------------------------ |
| Man syscalls | Linux system calls                                     |
| Man 7 libc   | Overview of standard C libraries on Linux              |
| Strace -c ls | Summary dell'uso delle syscalls                        |
| Man 2 open   | Mostra la pagina del manuale per la system call *open* |
| Man 3 stdio  | Funzioni I/O                                           |

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

| Syscall (no bufferizzazione nello spazio utante)         | StdIO library (bufferizzazione nello spazio utente)                                                         |
| -------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------- |
| Open                                                     | Fopen                                                                                                       |
| Close                                                    | Fclose                                                                                                      |
| Read/write                                               | getchar/putchar, getc/putc, fgetc/fputc, fread/fwrite, gets/puts, fgets/fputs, scanf/printf, fscanf/fprintf |
| Lseek                                                    | Fseek                                                                                                       |
| Richiedono cambio contesto user-kernel ad ogni carattere | Tipicamente più efficaci e più veloci                                                                       |
###### Bash: standard stream e redirezioni
**command-line completion:** *bash* ha la funzione di completamento delle linee di comando premendo il tasto *TAB*. Premendo una volta il tasto *TAB* ottengo:
1. se c'è una sola possibilità il comando viene completato 
2. altrimenti il comando viene eventualmente completato fin dove è possibile farlo in modo univoco e viene emesso un segnale sonoro

Nel secondo caso, premendo una seconda volta, viene mostrata la lista delle scelte possibili

**Gestione della history:** *bash* tiene nota degli ultimi *n* comandi dati alla shell. *n* dipende dalle variabili *HISTFILESIZE* e *HISTSIZE* (*man bash*)

| Comando (interno) | Descrizione                                                                                                                                |
| ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------ |
| history           | visualizza gli ultimi n comandi eseguiti (il numero di comandi visualizzato e memorizzato dipende dalle variabili HISTSIZE e HISTFILESIZE) |
| !a                | esegue l’ultimo comando nella history che inizia per ‘a’                                                                                   |
| !68               | esegue il comando numero 68 della history                                                                                                  |
| !!                | ripete l’ultimo comando eseguito                                                                                                           |
**Bash: Metacaratteri:** si tratta di caratteri che separano le parole. Un metacarattere può essere: spazio, tab, newline, `|`, `&`,`;`,`(`, `)`, `<`, `>` (da *Bash Reference Manual*). Questi caratteri hanno un significato per la shell, come anche i caratteri `*`, `?`, `Ctrl + d (EOF)`, `Ctrl + c`, `'`,…
> [!IMPORTANT]
> Se si desidera usare questi caratteri come testo normale, devono essere racchiusi tra **apici** (quoting) per evitare che la Bash li interpreti.

**Bash: standard streams:** In un sistema Unix-like, ogni processo all'avvio ha tre flussi di dati (file descriptor) predefiniti aperti. Questi flussi permettono al comando di interagire con l'ambiente esterno (tastiera, terminale, file) tramite anche la redirezione e il comando pipe (`|`).
* **0 (stdin):** Standard Input. Tipicamente collegato alla **tastiera**, è utilizzato come sorgente da cui prendere i dati da elaborare.
* **1 (stdout):** Standard Output. Tipicamente collegato al **terminale** (monitor), è utilizzato per mostrare i risultati della computazione.
* **2 (stderr):** Standard Error. Tipicamente collegato al **terminale**, è utilizzato per comunicare informazioni addizionali sullo stato della computazione come i messaggi di errore.
* **3+ (stdfile):** Descrittori di file aggiuntivi definiti dall'utente o dal sistema.

> [!NOTE] Nota del Prof
> Se esegui un comando come `sort` senza alcun parametro, la Bash "apre lo 0": invece di leggere da un file, resta in attesa di input direttamente dalla tastiera finché non viene inviato un segnale di EOF (Ctrl+D).

**Bash: redirezione dell’input:** La redirezione dell'input in Bash permette di inviare il contenuto di un file (o altro input) a un comando, sostituendo la tastiera standard (stdin, descrittore file 0) con l'operatore `<`: `comando [n] < filename`. 
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

**Bash: redirezione dell'Output:** I simboli `>` e `>>` permettono di deviare/salvare l'output di un comando (che normalmente andrebbe a video) verso un file (sovrascrivendolo). `comando [n] > filename`
- Il descrittore *n* è associato a *filename* che viene aperto in scrittura
- Se *n* è omesso lo *standard output* (*n* = 1) è associato al file, cioè il comando scrive sul file l’output
> [!example] `ls -l > lista.txt`: Salva la lista dei file in `lista.txt`, sovrascrivendolo.

**BASH: redirezione di stderr:** `comando [n]> filename` (default `n=1`), reindirizza solo lo standard error in un file
> [!example] 
> - `cat /etc/hosts > file.out`: Il contenuto di `/etc/hosts` viene scritto in `file.out` . Se il file esiste, viene sovrascritto.
> - **Redirezione Errori:** `ls cartella_inesistente 2> errori.log` Invia solo i messaggi d'errore al file `errori.log` .

**BASH: redirezione dell'output in append:** `comando >> file` viene aggiunto l’output in coda al file senza sovrascriverlo
> [!example] `echo "nuova riga" >> scan.txt`: Il testo viene aggiunto in coda al file senza cancellare il contenuto precedente.

**BASH: reindirizzare/duplicare uno stream:** È possibile duplicare o spostare i descrittori di file per operazioni avanzate con gli operatori della redirezione.
- `[n]<&x`: duplica il file descriptor in input x, nel senso che n fa riferimento allo stesso device/file cui si riferisce x. Se n è omesso viene assunto 0 (stdin). x deve essere un descrittore aperto in input altrimenti viene restituito un errore
- `[n]>&x`: duplica il file descriptor di output x. Se n è omesso viene assunto 1 (stdout). x deve essere un descrittore aperto in output altrimenti viene restituito un errore
> [!example] `3< /percorso-file comando 0<&3`:  Viene aperto il file sul descrittore 3, poi il comando (es. `cat` ) viene istruito a leggere dallo 0 che è stato duplicato dal 3 (simile alla system call `dup2` ).

> [!NOTE] Nota del Prof
> La funzione di sistema `read` torna sempre il numero di byte effettivamente letti. Per analizzare queste operazioni a basso livello, si può usare il comando `strace` , che mostra tutte le **system call** eseguite da un comando.

**BASH: redirezione simultanea di stdout e stderr:** 
- Esistono due possibili sintassi: `comando &> filename` oppure  `comando >& filename`. La prima delle due è preferibile e semanticamente equivale a`comando > filename 2>&1`. 
- Analogamente si realizza l’append:`comando &>> filename`, che semanticamente equivale a `comando >> filename 2>&1`

**Bash:Pipeline:** `cmd1 | cmd2 | ... | cmdN`. Posso scrivere ed eseguire vari comandi separati solo dal carattere `|` (pipe). La pipe collega direttamente lo Standard Output di un comando allo Standard Input di quello successivo. Al posto del carattere `|` posso anche utilizzare solo `|&` con cui stdout e stderr del comando precedente vengono entrambi connessi all’input del comando successivo. Ogni comando è eseguito in un processo differente (*subshell*)
> [!example] `sort /etc/passwd | less`: Il file `/etc/passwd` viene ordinato da `sort` , l'output prodotto non viene stampato a video ma passato a `less` , che permette di consultarlo pagina per pagina.

###### System calls per i processi
**fork():** `include <unistd.h> pid_t fork();`

È una syscall che non vuole argomenti. Questo metodo ritorna -1 difficilmente (il *child* non viene creato e errno viene settata).
- Crea un nuovo processo duplicando il chiamante, tuttavia alcuni attributi non sono duplicati: PID, memory locks, segnali pendenti, timer, etc,
- Il nuovo processo  è detto *child*, il processo chiamante è detto *parent*
- Lo spazio di indirizzamento del nuovo processo è un duplicato di quello del *parent* (sono separati, ognuno ha il proprio)

> [!info] 
> - fork() nel *child* torna 0, mentre nel *parent* torna il valore del PIDdel *child* (valore > 0). In questo modo si può differenziare il codice tra i due
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
/* fork a child process */
PID= fork();
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
printf("Parent: my PID= %d, parent PID= %d. PIDdel child %d\n",getpid(),getppid(),PID);
//sleep(10);
}
sleep(10);
return 0;
}
```
**PIDe PPID:** `include <unistd.h> pid_t getpid(void); pid_t getppid(void)`

- getpid() torna il process ID
- getppid() torna il parent process ID 

Queste due funzioni hanno sempre successo

> [!info] esiste anche `getpgrp()` che ritorna il *process group ID*

**wait () e waitPID():** `#include <sys/wait.h> pid_t wait(int *wstatus); pid_t waitpid(pid_t pid, int *wstatus, int options)`

- *wait()* sospende il chiamante finché un *child* termina
- *waitpid()* sospende il chiamante finché un *child* cambia stato (terminato; stoppato da un *signal*; ha ripreso dopo un *signal*), dipende da *options*

Permettono di rilasciare le risorse assegnate al *child*. Se un *child* termina e una *wait* non viene eseguita il *child* rimane nello stato di *zombie* (*man wait(2)*). Questi metodi ritornano il PIDdel processo il cui stato è cambiato (*waitpid()* può tornare 0 in certe circostanze) oppure -1 in caso di errore. Se un/il *child* è terminato/ha già cambiato stato tornano immediatamente. Tra i valori che possiamo dare a PIDci sono i seguenti:
- -1: attende per un cambiamento di stato di uno qualsiasi dei child
- >0: attende per un cambiamento di stato del child il cui process ID è uguale a pid
- Etc

*waitpid(-1, &wstatus, 0)* è equivalente a *wait(&wstatus)*. Se *&wstatus* non è NULL riceve informazioni sul *child*. Queste informazioni possono essere evinte con delle macro:
- WIFEXITED(wstatus) torna true se il child è terminato normalmente (exit() o return) 
- WEXITSTATUS(wstatus) ritorna il codice di uscita del child 
- Etc. (man 2 wait)

> [!info] Il comando `watch` in bash esegue ripetutamente un comando specifico a intervalli regolari (di default ogni 2 secondi), mostrando l'output nel terminale. È fondamentale per monitorare in tempo reale le modifiche all'output di un comando senza doverlo rieseguire manualmente, aggiornando la visualizzazione. Con `-n` imposto l’intervallo di tempo e con `-d` evidenzia le differenze tra l’esecuzione precedente e l’attuale

**System calls per i processi**

| Nome              | Cosa fa                                          |
| ----------------- | ------------------------------------------------ |
| Fork              | Crea un nuovo processo                           |
| Wait              | Attende finché uno dei child termina             |
| WaitPID          | Attende finché un determinato child cambia stato |
| Waitid            | Maggiore controllo su quali child attendere      |
| _exit(2), exit(3) | Terminano il chiamante                           |
| Execvve(2)        | Esegue un programma che rimpiazza il chiamante   |
| Exec*(3)          | Wrappers per execve(2) (sono 6)                  |
###### File Bash
```C
// Prima riga = quale interprete deve interpretare i comandi: shabang
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
```

> [!attention] Nei file bash (e nella shell in generale), non bisogna mettere spazi attorno al segno di uguale quando si definisce una variabile perché **bash usa gli spazi per separare i comandi dai loro argomenti** (nome variabile verrebbe interpretato come nome di un comando da eseguire e `=` e il `valore` della variabile verrebbero interpretati come argomenti di questo comando)

