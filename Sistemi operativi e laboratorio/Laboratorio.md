### 26-02-27 
**Cenni Storici e Famiglie di Sistemi Operativi**
=>*Unix*
* **Unix** è una "famiglia" di sistemi operativi **multitasking** e **multiuser** che deriva dall'originale AT&T Unix.
* Lo sviluppo di AT&T Unix è iniziato nel 1969 al Bell Labs research center, principalmente ad opera di **Ken Thompson**, **Dennis Ritchie** e altri collaboratori.
* Inizialmente si è assistito a una forte frammentazione dove "ognuno si fece il suo", dando vita a varianti come BSD (UC Berkeley), Xenix (Microsoft), SunOS/Solaris (Sun Microsystems), HP-UX (HP), e AIX (IBM).
* Nel corso del tempo, sono stati fatti sforzi per standardizzare le varie release e varianti attraverso iniziative come **POSIX**, **SUS**, e **The OpenGroup**.

=>*(GNU/)Linux*
* **Linux** è un discendente o derivato di Unix, definibile come **Unix-like**.
* È una "famiglia" **open-source** di sistemi operativi Unix-like.
* Si basa sul kernel Linux, rilasciato per la prima volta nel **1991** da **Linus Torvalds**.
* Tecnicamente, sarebbe più corretto parlare di sistema operativo **GNU/Linux**.
* È un sistema **POSIX compatibile**, seppur non ufficialmente POSIX-certified.
* Esistono varie **distribuzioni** (o "distro"), tra cui Debian, SUSE, RedHat, Fedora, Ubuntu e Mint.

=>*macOS*
* Rappresenta un gruppo di sistemi operativi proprietari di casa **Apple**.
* Sono sistemi Unix-like e, soprattutto, **Unix certified** (certificati da The Open Group).
* Hanno attraversato diverse transizioni di **processori**: inizialmente PowerPC, poi dal 2006 sono passati ad Intel, e dal 2020 utilizzano processori 64-bit ARM-based come l'Apple M1 (seguiti da M2, M3, ecc.).
* Ad Ottobre 2024, Apple ha presentato il chip **M4**.

=>*Android*
* Sistema basato su un **kernel Linux modificato**.
* Originariamente sviluppato dalla Open Handset Alliance, la sua versione più diffusa è oggi sviluppata principalmente da **Google**.
* Il resto del sistema, inclusi software vari e utilities, è scritto principalmente in **Java**.
* È specificamente pensato per sistemi **touchscreen mobile**, come telefoni e tablet.

**Concetti Base dei Sistemi Unix**

=>*Struttura e Architettura*
* Sistemi **multitasking e multiuser**.
* **Modularità**: Unix non è un sistema monolitico.
    * Il **kernel** fornisce i servizi fondamentali del sistema.
    * La **shell** (o in alternativa l'interfaccia grafica) permette all'utente di interagire con il sistema operativo.
    * Esistono tante semplici utility che svolgono compiti base; queste possono essere combinate insieme tramite **pipe line** per eseguire compiti anche molto complessi.
* È stato concepito appositamente per supportare gruppi di sviluppatori che lavorano insieme.
* **Portabilità**: Essendo stato quasi subito riscritto in **C**, risulta molto più portabile rispetto ai vecchi sistemi scritti in assembly.
* I device e i vari strumenti per l'IPC (Inter-Process Communication) sono trattati e mostrati dal sistema come se fossero dei semplici **file**.
* Il formato di memorizzazione dei file è volutamente semplice (**plain text**).

 =>*Autenticazione e Utenti*
* Essendo multiutente, ogni utente possiede una **login**, una **password** e una **home directory** (uno spazio disco personale dove salvare i propri file).
* Esiste un utente privilegiato, chiamato **root**, che funge da amministratore del sistema.
* Ogni utente gestisce in autonomia i permessi di accesso sulle proprie cartelle e sui propri file.
* Gli utenti possono essere aggregati in **gruppi** per definire e semplificare privilegi comuni.

**Terminale, Console e BASH**

=>*Definizioni*
* **Terminal**: Un ambiente di input e output testuale. Una finestra di terminale (o emulatore di terminale) è una finestra di solo testo che emula una console all'interno di una GUI (Graphical User Interface).
* **Console**: Un terminale fisico; un pannello della strumentazione contenente i controlli del computer. Una console è quindi un tipo specifico di terminale.
* **Shell**: È l'interfaccia vera e propria tra il sistema operativo e l'utente. Permette di eseguire programmi e offre funzionalità simili a un linguaggio di programmazione per automatizzare operazioni. Può essere a riga di comando (Bash, csh, ksh) o grafica (GNOME, KDE).

=>*BASH (Bourne Again Shell)*
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

=>*Perché usare una shell testuale?*
Anche se esistono le interfacce grafiche, la shell testuale si usa perché:
* Aiuta a capire come funziona il sistema "sotto".
* Offre il massimo controllo sui comandi e sulle loro opzioni.
* Permette di configurare le distribuzioni in modo preciso.
* I comandi Unix sono stati nativamente progettati per le interfacce a riga di comando.
* Una volta imparati, risultano più veloci e flessibili.
* Sono estremamente veloci e facilmente accessibili da remoto.
* Funzionano sempre, anche in caso di problemi all'interfaccia grafica o con hardware obsoleto e poco performante.

**Il Filesystem**
Il filesystem è una struttura ad albero al cui interno sono memorizzati tutti i dati del sistema operativo.
* È organizzato in **File** e **Directory** (o cartelle).
* L'albero ha un'unica radice detta **root**, identificata dal path `/`. A differenza di Windows (che ha `C:`, `D:`, ecc.), in Unix tutti i dati convergono in quest'unico albero, indipendentemente dal numero di dischi o partizioni fisiche. 

**Navigazione e Pathnames**
Ogni processo nel sistema ha, in ogni singolo istante, una "directory corrente".
* **Pathname assoluto**: È il percorso che parte dalla radice root e porta all'oggetto. Comincia sempre con la barra `/` (es. `/etc/hosts` o `/home/fabio/Desktop`).
* **Pathname relativo**: È il percorso che porta all'oggetto partendo dalla directory corrente o da un punto specifico che non è la radice. Non inizia mai con `/` (es. `log/apache`).
* **Simboli speciali nei path**:
    * `.` (punto): indica la directory corrente.
    * `..` (due punti): indica la directory "padre", ovvero quella che contiene la directory corrente.

**Link e Device Files**
* Nel filesystem possono essere presenti file speciali chiamati **link**, che fungono da puntatori/collegamenti ad altri oggetti (in modo analogo ai collegamenti di Windows). Compiere un'operazione su un link significa implicitamente compierla sull'oggetto a cui esso punta.
* Sono presenti anche **special files** o **device files**: espongono un dispositivo hardware all'interno del filesystem come se fosse un normale file. Questo fornisce un'interfacversale verso l'hardware, permettendo alle system calls di eseguire operazioni dicia uni Input/Output in modo standardizzato.

**Primi Comandi BASH**
I comandi BASH hanno tipicamente questa sintassi: `nome_comando <opzioni> <argomenti>`. 
Le **opzioni** iniziano tipicamente con il carattere `-` e modificano il comportamento del comando; possono essere accorpate (es. `ls -ltra` equivale a `ls -l -t -r -a`) e a volte richiedono argomenti aggiuntivi. Ogni comando restituisce un valore di ritorno che informa sull'esito, dove solitamente `0` indica che è andato tutto a buon fine.

**Comandi Interni vs Esterni**
* **Comandi Esterni**: Sono veri e propri eseguibili registrati nei supporti di memorizzazione. La shell li cerca all'interno dei percorsi di filesystem specificati nella variabile di ambiente `$PATH`. Il comando `which` permette di scoprire dove è memorizzato l'eseguibile di un comando esterno.
* **Comandi Interni**: Sono riconosciuti ed eseguiti direttamente dalla shell stessa (es. i costrutti per lo scripting o il comando `cd`). È possibile vederne la lista utilizzando il comando interno `help`.

***Operazioni sul Filesystem***

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
**Opzioni più Comuni per `ls`**
* `-l`: Formato lungo (long format), mostra informazioni aggiuntive sui file.
* `-a`: Elenca tutti i file, compresi quelli "nascosti" (quelli che iniziano col carattere `.`).
* `-t`: Ordina i file mostrati in base alla data di ultima modifica (dal più recente al meno recente).
* `-r`: Inverte il senso dell'ordinamento.
* `-h`: Visualizza le dimensioni dei file in un formato "human readable" e facilmente leggibile.
	
> [!EXAMPLE] Esempio combinato
> Il comando `ls -lrta` elenca tutti i file e le cartelle (compresi i nascosti), ordinandoli in base al tempo di ultima modifica (dal meno recente al più recente) formattandoli in modalità lunga.

**Opzioni più Comuni per CANCELLARE o COPIARE Cartelle**
* `rm -r cartella`: Rimuove ricorsivamente una cartella e tutto il ramo di file/cartelle non vuote in essa contenute.
* `cp -ar dir1 dir2`: Copia ricorsivamente (`-r`) tutto l'albero sottostante `dir1` dentro `dir2`, preservando (`-a`) gli attributi originari dei file come la data di modifica.

 **Differenza nella creazione dei Link (`ln`)**
* **Hard Link** (`ln pippo link1`): Crea una vera e propria copia che occupa spazio fisico sul disco. Se il file originario "pippo" viene cancellato, il file "link1" e il suo contenuto rimangono inalterati.
* **Symbolic Link** (`ln -s pippo link1`): Crea un collegamento simbolico (occupazione minimale su disco). "link1" contiene solo una stringa col percorso di "pippo". Se il file "pippo" originario viene cancellato, "link1" diventa semplicemente un collegamento non valido o rotto.

**Ottenere Aiuto e Documentazione**
Per imparare le funzioni e le opzioni dei comandi, ci sono diverse utility fondamentali:
* **`man <nome_comando>`**: Apre il manuale in linea completo del comando specificato.
* **`apropos <stringa>`**: Cerca e lista i comandi che contengono la stringa specificata all'interno della loro descrizione sintetica.
* **`whatis <comando>`**: Fornisce una descrizione estremamente sintetica del comando passato.
* **`locate <stringa>`**: Ricerca velocemente tutti i file aventi un path che contiene la stringa specificata, interrogando un database interno costruito tramite il tool `updatedb`.
* **`Switch di aiuto rapido`**: Generalmente ogni comando BASH, se invocato aggiungendo lo switch `-h` oppure `--help`, è in grado di stampare a terminale una rapida e sintetica descrizione delle proprie funzionalità.
---
### 26-03-03 
***Variabile $PATH:***
Se voglio vedere dove è memorizzato l'eseguibile che corrisponde ad un comando (esterno) uso il comando **which**. Uso la variabile d'ambiente $**PATH** per dire al terminale dove cercare i programmi (sono una serie di cartelle). Quindi, quando si desidera eseguire un programma, è sufficiente digitare il nome del file eseguibile e il terminale cercherà nelle cartelle specificate in PATH se è presente un file di programma eseguibile denominato in quel modo.

> [!info] Con il comando help variables vedo nome e significato di alcune delle shell variables

| Comando | Descrizione                                          | Tipologia |
| ------- | ---------------------------------------------------- | --------- |
| Set     | Visualizza, imposta o annulla le opzioni della shell | Interno   |
| Env     | Esegue un programma in un ambiente modificato        | Esterno   |
> [!info] 
> - Sudo è il comando che mi permette di richiedere i permessi da utente root.
> - I software vengono distribuiti in Linux sotto forma di pacchetti. 
> - cd senza argomenti mi riporta alla home directory

***Attributi del file/cartella (bash)*** 
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

---
### 26-03-06
Le system calls sono l'interfaccia per il programmatore ai servizi del sistema operativo (chiamata diretta al sistema). Vi si accede generalmente attraverso un API/linguaggi che la supportino. Ce ne sono vari tipi (accesso ai file, gestione di processi, IPC, semafori,...) e richiedono un "cambio di contesto": da *kernel mode* a *user mode* e viceversa.
- Kernel mode = accesso senza limitazioni alle risorse HW/SW del sistema (memoria, dispositivi I/O, registri, codice del SO,...). Un errore in kernel mode può provocare un arresto del sistema
- User mode = accesso limitato alle risorse di sistema (es. Memoria). Un errore in user mode può provocare il blocco del solo processo utente

Quando si passa da un mode ad un altro il SO cambia il contesto di esecuzione del processore. Ciò implica un certo consumo di risorse di sistema (tempo di CPU, salvataggio e ripristino del processo e del contesto precedente,...)

**Funzioni wrappers: Glibc**
> [!info] Glibc = "GNU C library" è la libreria C standard che di solito si usa in Linux (contiene le funzioni standard richiamabili da C)

Le system calls sono invocate non direttamente ma tramite funzioni *wrappers* della glibc:
- Le *wrappers* hanno di solito lo stesso nome della corrispondente *system call* sottostante
- Ove possibile è  preferibile chiamare le funzioni *wrappers* piuttosto che direttamente la *system call*. Questo perché spesso si limitano a fare poco altro in più rispetto alle *system call*, ma in altri casi compiono anche altre operazioni come controllo degli argomenti inseriti dal programmatore,... 

***Accesso ai files:***
Per l'accesso ai  file ci sono due possibilità:
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
| dup, dup 2, dup 3   | Duplicano un file descriptor                                                   |
**System calls per le directory:**

| Nome  | Cosa fa                           |
| ----- | --------------------------------- |
| mkdir | Crea una directory                |
| rmdir | Cancella una directory (se vuota) |
| chdir | Cambia la directory corrente      |
**Open**
` #include int open(const char *pathname, int flags, mode_t mode)`
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
- Un intero positivo che è il *file descriptor* da usare per le operazioni successive
- -1 in caso di errore che viene settato da errno

> [!info] 
> - di default il puntatore (file offset) nel file è posizionato all'inizio (a meno che non si specifichi O_APPEND)
> - i file aperti andrebbero chiusi prima di terminare il programma con `#include int close(int fd)`, dove fd è il file descriptor ritornato da open () 
> - l'apertura di un file crea una nuova entry nella tabella di sistema dei file aperti. Quando un file descriptor viene duplicato (con dup 2 o simili), questo si riferisce alla stessa *open file description* del descrittore del file originale $\implies$  i due file descriptors condividono lo stesso *file offset* e le stesse *file status flags*
> - su Linux si può usare l'operazione kcmp (2) KCMP_FILE per testare se due file descriptors si riferiscono alla stessa open file description
 
> [!tip] 
> - Errno = è una libreria di C (errno.h) che mi dice precisamente quale sia l'errore in corso 
> - con `strace -c ls` posso vedere le system call attive e alcuni dati che le riguardano: percentuale del tempo totale con cui il kernel è stato dentro alla syscall, tempo totale speso in quel tipo di syscall, media per ogni call, numero di volte che la syscall è stata eseguita, quante syscall hanno ritornato un errore, nome della syscall
> - **buffering** = in generale il buffer è una certa quantità di memoria che deve essere riempita prima che accada una qualche chiamata o azione specifica (ad esempio se scrivo una stringa questa potrebbe non essere visualizzata subito perché la memoria buffer non è ancora piena). Esistono 3 tipi di buffer: interattivo (buffer si svuota quando appare `\n`), di pipe o di file (buffer si svuota solo quando la memoria predisposta è piena) e l'assenza di buffer (ogni carattere fa si che la memoria buffer si svuoti subito)

``` C
#include <fcntl.h> //serve per le costanti che usa open() come O_CREAT, O_RDWR
#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <string.h>
#include <unistd.h>  //syscalls come write() e close()

int main()
{
	int iRet = EXIT_SUCCESS, fd; //iRet è il valore che ritorna il programma e fd è il file descriptor
	char sBuf[] = "Ciao!\n"; //buffer che verrà scritto nel file
	
	fd = open("pippo.txt", O_RDWR | O_CREAT | O_TRUNC, S_IRUSR | S_IWUSR); //chiamata al kernel con la syscall open()
	//S_IRUSR = read permission for user
	//S_IWUSR = write permission for users
	if (fd >= 0)
	{
		if (write(fd, sBuf, strlen(sBuf)) == strlen(sBuf)) //chiamata al kernel con la syscall write()
		//write prende come argomenti il file descriptor, il memory buffer e numero di bytes
		{
			printf("Ho scritto nel file pippo.txt, fd = %d\n", fd);
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

> [!info] fd =
> - 0 = stdin
> - 1 = stdout
> - 2 = stderr
> - 3 = riguarda l'apertura del file

| Comando      | Contenuto                                              |
| ------------ | ------------------------------------------------------ |
| Man syscalls | Linux system calls                                     |
| Man 7 libc   | Overview of standard C libraries on Linux              |
| Strace -c ls | Summary dell'uso delle syscalls                        |
| Man 2 open   | Mostra la pagina del manuale per la system call *open* |

---
### 26-03-10
Attraverso `#include <unisrd.h>` posso usare anche syscall (wrapped) come write, read,...

**write**
`ssize_t write(int fd, const void *buf, size_t count);`
Scrive i primi count byte contenuti in buf nel file aperto identificato da fd (file descriptor ritornato da open(2)). Ritorna il numero di byte scritti (che può anche essere inferiore a count se, per esempio, il device di scrittura si è riempito) o -1 in caso di errore (errno settata)

**read**
`ssize_t read(int fd, void *buf, size_t count);`
Legge count byte dal file identificato da fd (file descriptor ritornato da open(2)) e li memorizza in buf. Ritorna il numero di byte letti (che può anche essere inferiore a count se c'erano a disposizione nel device di lettura meno byte) o -1 in caso di errore (errno settata). Il valore di ritorno 0 significa EOF (end of file).

> [!info] le operazioni di lettura/scrittura sono eseguite a partire dalla posizione corrente (puntatore, cursore o file offset) nel file

**lseek**
`off_t lseek(int fd, off_t offset, int whence);`
riposiziona il puntatore nel file. A seconda dei seguenti valori di whence il puntatore è posizionato: 
- SEEK_SET: a offset byte dall'inizio 
- SEEK_CUR: a offset byte dalla posizione corrente 
- SEEK_END: a offset byte dopo la fine del file 

Ritorna la nuova posizione (distanza in byte) da inizio file o -1 in caso di errore (setta errno)

**dup, dup 2**: 
`int dup(int oldfd);`
`int dup2(int oldfd, int newfd);` 
- dup() torna una copia di oldfd usando il primo file descriptor libero 
- dup2() torna una copia di oldfd ma usa newfd come file descriptor copia. 
N.B.: Se c'era già un file aperto con descrittore pari a newfd viene chiuso 
- dopo la duplicazione i due file descriptor possono essere usati indifferentemente, ma si riferiscono alla stessa entry nella tabella di sistema che elenca i file aperti e quindi, per esempio, condividono la posizione corrente nel file (se viene modificata tramite uno dei due la modifica vale anche per l'altro) 
- Ritornano il nuovo file descriptor copia; o -1 in caso di errore (errno settata)