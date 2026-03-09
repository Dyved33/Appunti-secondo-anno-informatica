### 26-02-26 lz.1
**Compilatori e Interpreti**
La traduzione del codice sorgente in un linguaggio comprensibile dalla macchina può avvenire in modalità differenti (è anche possibile la traduzione da un linguaggio di programmazione a un altro).
* **Compilatore:** Esegue la traduzione in blocco *prima* dell'esecuzione. Prende in input il **codice sorgente** e, tramite il traduttore, restituisce il **codice oggetto** (oppure delle notifiche di errore).
* **Interprete:** Esegue la traduzione *istruzione per istruzione*, simultaneamente all'esecuzione del programma. Prende in input sia il **codice sorgente** sia l'**input** dell'utente, producendo direttamente l'**output**.

Esistono compilatori Ibridi (Es. Java) che combinano i due approcci:
1. Traduttore analizza il codice sorgente e lo traduce in un **programma intermedio** (ad esempio, il *Java bytecode*).
2. Successivamente, una **macchina virtuale** (ad esempio, la JVM) riceve gli input dell'utente e *interpreta* il codice intermedio per produrre l'output finale.

**Le Fasi della Compilazione: Dall'Analisi alla Sintesi**
La compilazione si divide concettualmente in due macro-fasi principali:
* **Fase Analitica (Analisi):** Riguarda la comprensione del codice in ingresso. Ha lo scopo di costruire l'**Albero Sintattico Astratto (AST)** partendo dal codice sorgente.
* **Fase Sintetica (Sintesi):** Riguarda la generazione dell'output. Si occupa di produrre il codice oggetto finale a partire dall'AST generato.

**Le 8 Fasi del Sistema di Compilazione**
Il passaggio dal codice sorgente al codice macchina attraversa stadi sequenziali, che comunicano tra loro avvalendosi di strutture condivise come la **tabella dei simboli** e il modulo di **gestione degli errori**.

=>*FASE DI ANALISI*
1. **Gestione Input (Preprocessore / Lettura):** Prende in ingresso il codice sorgente puro.
   > [!NOTE] Nota del Prof
   > Il codice sorgente è tipicamente un file di testo che deve essere letto e processato sequenzialmente dal sistema.

2. **Analisi Lessicale (Scanner):** Il file in ingresso è una sequenza di caratteri. Lo scanner raggruppa questi caratteri in sequenze elementari chiamate **lessemi**. A ogni lessema viene associata una classe lessicale (**token**) e un **attributo** (che viene registrato nella tabella dei simboli).
   * *Esempio sull'espressione `position = initial + rate * 60`:*
     * `position` $\rightarrow$ Token: *id*, Attributo: *1*.
     * `=` $\rightarrow$ Token: *assign*.
     * `initial` $\rightarrow$ Token: *id*, Attributo: *2*.
     * `60` $\rightarrow$ Token: *number*, Attributo: *4*.
   > [!EXAMPLE] La scelta dei token
   > L'assegnazione dei token e degli attributi segue regole specifiche che cambiano da linguaggio a linguaggio. La scelta dei token è a discrezione del progettista del linguaggio (ad esempio, nel Fortran, la sintassi e gli spazi richiedono regole di scansione particolari).

3. **Analisi Sintattica (Parser):** Stabilisce le regole grammaticali. Il parser prende la sequenza di token prodotta dallo scanner e la organizza in un **Albero Sintattico Astratto (AST)**.
   * I *nodi interni* dell'albero rappresentano le operazioni.
   * I *nodi figli* rappresentano gli argomenti di tali operazioni.
   > [!EXAMPLE] Struttura dell'AST
   > Confrontando una funzione scritta in 3 linguaggi diversi (Pascal, ML, Java) con 3 parametri reali, una variabile locale e un'operazione, la struttura ad albero logica rimane la stessa. Il nodo principale (radice) viene identificato come `FuncDef`.

4. **Analisi Semantica Statica:** Utilizza l'albero sintattico astratto e la tabella dei simboli per verificare che il programma sorgente sia semanticamente coerente con la definizione del linguaggio.
   * **Type checking:** Verifica che ogni operatore abbia tutti gli operandi del tipo corretto. Se necessario, può eseguire la **coercizione** (conversione forzata del tipo).
   * *Esempio in C:* È permesso moltiplicare un intero per un reale; il compilatore converte silenziosamente l'intero in un reale.
   * *Esempio in Pascal:* Il costrutto `WHILE` richiede obbligatoriamente un valore `BOOLEAN`. L'espressione `WHILE X DO X := X - 1 ;` farà segnalare un errore all'analizzatore semantico.

=>*FASE DI SINTESI*
4. **Generazione del Codice Intermedio:** Dall'AST e dalla tabella dei simboli si ottiene un codice intermedio basato su istruzioni elementari, facili da tradurre in codice macchina. È **indipendente dall'architettura** (non specifica gli indirizzi fisici in memoria né i registri).
   * *Esempio pratico:* L'espressione precedente viene scomposta in operazioni atomiche: `t1 = inttofloat(60)`, `t2 = id3 * t1`, `t3 = id2 + t2`, `id1 = t3`.
6. **Ottimizzazione:** L'ottimizzatore elabora il codice intermedio per ridurre il tempo o lo spazio necessario all'esecuzione.
   > [!NOTE] Taglio delle Istruzioni
   > L'ottimizzatore snellisce il codice eliminando passaggi inutili. Nel nostro esempio, riduce le operazioni accorpando i calcoli: `t2 = id3 * 60.0` e `id1 = id2 + t2`.
7. **Generazione del Codice Oggetto:** È la parte più complessa e dipende strettamente dall'architettura della piattaforma hardware. Il compilatore genera codice assembly o file oggetto ricollocabili. I passaggi chiave sono:
   * Fissare le locazioni di memoria esatte per i dati.
   * Generare il codice di basso livello per accedere a tali dati.
   * Selezionare i registri fisici della CPU per i calcoli intermedi.
   > [!IMPORTANT] Il ruolo dell'Assemblatore e del Linker
   > A rigore, l'Assemblatore (se il compilatore emette assembly) e il Linker sono strumenti esterni che intervengono al termine della compilazione. Il Linker si occuperà di collegare il codice oggetto generato in questa fase con eventuali librerie esterne, producendo il file eseguibile finale.
8. **Ottimizzazione Peep-hole:** Un'ulteriore e ultimissima fase di ottimizzazione legata alle specificità e ai "trucchi" dell'architettura hardware di destinazione.

**L'Architettura Software del Compilatore**
L'organizzazione interna del compilatore definisce come i moduli comunichino effettivamente tra loro durante il processo di traduzione.

* **Il Nucleo Centrale (Parser & Semantica):** la parte centrale di tutto il processo è l'**analizzatore sintattico**. In questa architettura, esso realizza contestualmente anche l'**analisi semantica statica**.
   > [!IMPORTANT] Avanzamento Parallelo
   > Dato che le due funzioni sono accorpate nel nucleo, l'analisi sintattica e l'analisi semantica avanzano parallelamente.
* **Gestione Input e Lettura:** Il sistema inizia con il modulo **gestore input**, che riceve il codice sorgente e può generare un "listato sorgente". Per formare i token, l'analizzatore lessicale deve leggere il codice un carattere alla volta invocando la funzione `readChar()` direttamente sul gestore input.
* **Dialogo Scanner-Parser:** L'analizzatore sintattico coordina il lavoro chiedendo i dati all'analizzatore lessicale tramite la chiamata alla funzione `prossimoToken()`. Questa funzione, implementata nello scanner, restituisce il `token` richiesto per continuare a costruire l'albero.
* **Gestione degli Errori:** Il **notificatore degli errori** riceve comunicazioni dirette e incrociate dal gestore input, dall'analizzatore lessicale e dal blocco sintattico/semantico ogni volta che viene rilevata un'anomalia. Alla fine, produce in output la "lista degli errori".
* **Generazione del Codice Finale:** Completata l'analisi senza errori bloccanti, è sempre l'analizzatore sintattico ad attivare la sintesi. Tramite chiamate a funzioni opportune, innesca il **generatore del codice**, che a sua volta emetterà il **codice oggetto** definitivo.

---
### 26-02-27 lz.2
**Il problema dell'infinito e il Teorema di Cantor**
Per descrivere qualsiasi cosa, che sia la Divina Commedia o il codice genetico umano, usiamo delle parole. 
In informatica teorica, un **linguaggio è definito come un insieme di parole**.
* **Linguaggi Finiti:** Se un linguaggio è finito, la sua descrizione è banale poiché è possibile elencarne tutti gli elementi uno per uno. Un esempio è $L=\{a,ab,abb,abbb\}$.
* **Linguaggi Infiniti:** Il problema sorge quando cerchiamo di descrivere linguaggi con infinite parole.

Il **Teorema di Cantor** ci pone un limite matematico insuperabile: non esiste una funzione iniettiva dall'insieme delle parti $\mathcal{P}(\Sigma^{*})$ all'insieme delle parole $\Sigma^{*}$. Di conseguenza, non è possibile associare a *ciascun* linguaggio possibile una parola univoca (la sua "descrizione") che lo caratterizzi in modo perfetto. In altre parole, non possiamo inventare una "parola descrittiva" univoca per tutti i possibili linguaggi infiniti esistenti.

> [!NOTE] L'approccio Informatico
> Nonostante il Teorema di Cantor, in informatica non abbiamo bisogno di descrivere *tutti* i linguaggi infiniti. È comunque possibile fornire **descrizioni finite** per *alcuni* linguaggi infiniti che ci sono particolarmente utili. Ad esempio, possiamo usare una regola descrittiva semplice come "tutte le parole sull'alfabeto $\{a,b\}$ di lunghezza dispari", oppure la notazione formale $\{a^{n}b^{n}|n\ge0\}$.

Un costrutto fondamentale per generare infinità è la **chiusura di Kleene** di un linguaggio $L$ (indicata con $L^{*}$). Essa è costituita dalla parola vuota, dalle parole appartenenti a $L$ e da tutte le parole che si ottengono concatenando due o più parole di $L$. Se $L$ è un linguaggio finito, il linguaggio derivato $L^{*}$ è ben definito e, in generale, infinito.

**Alfabeti, Parole e le loro Proprietà (Definizioni Formali)**
Per poter costruire delle regole rigorose, dobbiamo prima definire matematicamente gli elementi base:

* **Alfabeto ($\Sigma$):** È un insieme finito e non vuoto di simboli. I suoi elementi sono chiamati **lettere**.
    * *Esempi di alfabeti:* $\Sigma_{0}=\{a,b\}$, $\Sigma_{1}=\{0,1\}$, $\Sigma_{2}=\{a,b,c\}$ e l'alfabeto esadecimale $\Sigma_{3}=\{0,1,2,3,4,5,6,7,8,9,A,B,C,D,E,F\}$.
* **Parola:** Ogni sequenza finita di lettere appartenenti a un alfabeto $\Sigma$ è detta parola sull'alfabeto $\Sigma$. 
    * *Esempi:* "a", "abb", "ababbabb" sono parole su $\Sigma_{0}$; "2020" e "CD078B" sono parole su $\Sigma_{3}$.
* **L'insieme Universo ($\Sigma^{*}$):** L'insieme di *tutte* le parole generabili sull'alfabeto $\Sigma$ è denotato con $\Sigma^{*}$.
* **Parola Vuota:** Esiste anche la sequenza composta da zero lettere, che si denota con $\epsilon$ (oppure $\Lambda$) e si dice **parola vuota**.
* **Lunghezza ($|u|$):** Data una parola $u=a_{1}a_{2}\cdot\cdot\cdot a_{k}$ (con $k\ge0$), l'intero $k$ rappresenta la sua lunghezza e si denota con $|u|$ oppure $l(u)$.
    * *Esempi:* $|abb|=3$, $|\epsilon|=0$.

**Operazioni sulle Parole**
* **Concatenazione:** Date due parole $u$ e $v$, la loro concatenazione $uv$ è semplicemente l'unione in sequenza delle loro lettere. Ad esempio, concatenando "abb" e "aaab" si ottiene "abbaaab".
    * *Proprietà Binaria:* È un'operazione binaria totale su $\Sigma^{*}$.
    * *Associatività:* Per ogni $u, v, w \in \Sigma^{*}$, vale $(uv)w=u(vw)$.
    * *Elemento Neutro:* La parola vuota $\epsilon$ fa da elemento neutro: $u\epsilon=\epsilon u=u$.
    * *Cancellatività:* Gode della cancellatività a destra e a sinistra (se $uw=vw$ oppure $wu=wv$, allora $u=v$).
* **Potenza n-esima ($w^{n}$):** Si ottiene concatenando $n$ copie della parola $w$ (es. $w^{n}=ww\cdot\cdot\cdot w$). Si definisce $w^{0}=\epsilon$. Se $u=abb$, allora $u^{2}=abbabb$.
* **Fattori, Prefissi e Suffissi:** Una parola $v$ è un **fattore** di $w$ se è contenuta al suo interno, ovvero se $w=xvy$ per opportune parole $x, y$.
    * **Prefisso:** Se $x=\epsilon$, allora $v$ si trova all'inizio ed è detto prefisso.
    * **Suffisso:** Se $y=\epsilon$, allora $v$ si trova alla fine ed è detto suffisso.
    * **Fattore Proprio:** Diremo che $v$ è un fattore proprio se $v\ne w$.

**Linguaggio Formale (Definizione Rigorosa)**
Ogni sottoinsieme di $\Sigma^{*}$ si dice **linguaggio formale** (o semplicemente linguaggio) sull'alfabeto $\Sigma$. Possono esserci linguaggi limitati come $L_{1}=\{a,ab,abb\}$, linguaggi infiniti come $L_{2}=\{ab^{n}a|n\ge0\}$, l'insieme vuoto $L_{3}=\emptyset$ o l'intero universo $L_{4}=\Sigma^{*}$.

**L'Intuizione: Il Gioco delle Sostituzioni**
Per capire operativamente come le "regole finite" possano generare "linguaggi infiniti", analizziamo un esempio pratico basato sul linguaggio naturale.

Vogliamo generare liste di nomi come "Aldo, Bianca e Carlo". Le regole logiche che governano queste frasi sono:
1.  I nomi sono separati da virgola, tranne gli ultimi due che sono separati dalla congiunzione 'e'.
2.  Sono ammesse ripetizioni di nomi (es. "Aldo, Bianca, Aldo...").
3.  Non c'è alcun limite alla lunghezza della lista.

Per istruire una macchina, dobbiamo usare dei **segnaposto** (come i termini "nome" o "frase") ed eseguire sostituzioni meccaniche secondo regole predefinite. 

> [!IMPORTANT] Il Problema delle Regole Globali vs Locali
> Un primo tentativo di creare queste regole umane includerebbe una direttiva del tipo: "prima di terminare, se sono presenti virgole, l'ultima va sostituita con la congiunzione e".
> Questa regola, tuttavia, è anomala rispetto alle altre: impone di analizzare l'intera frase simultaneamente per trovare "l'ultima virgola". Un calcolatore preferisce invece **regole di sostituzione rigorose e puramente locali**, dove un simbolo viene sostituito ciecamente a prescindere dal contesto globale.

**La Soluzione Formale: Ricorsione e Segnaposto**
La soluzione "informatica" parte sempre da un segnaposto iniziale (es. `frase`) e procede sostituendolo finché non rimangono solo parole e punteggiatura reali. 

1. **Regole di Base:** Il segnaposto $\langle nome \rangle$ può essere sostituito da Aldo, Bianca o Carlo.
2. **Inizio:** $\langle frase \rangle$ può diventare un singolo $\langle nome \rangle$ oppure il blocco $\langle lista Nomi \rangle \langle fineLista \rangle$.
3.  **Il Motore dell'Infinito (Ricorsione):** $\langle lista Nomi \rangle$ può essere sostituito da $\langle nome \rangle$ oppure da $\langle nome \rangle, \langle lista Nomi \rangle$. 

> [!EXAMPLE] Come si genera l'infinito
    > Contenendo se stessa nella parte destra dell'assegnazione, può essere applicata all'infinito per aggiungere quante virgole e quanti nomi desideriamo, creando una lunghezza illimitata a partire da una singola regola.

4.  **Chiusura:** Il blocco $\langle nome \rangle \langle fine Lista \rangle$ si sostituisce rigorosamente con $e \langle nome \rangle$. Questo trucco garantisce che, alla fine del processo meccanico, l'ultimo elemento sia sempre preceduto dalla 'e' e non dalla virgola.

Ci si arresta solo quando i segnaposto non compaiono più.
**Il concetto chiave è che un linguaggio di complessità e lunghezza infinita può essere generato meccanicamente e senza ambiguità tramite un set ristretto di regole di sostituzione locali**.

**Grammatiche a Struttura di Frase (Formalizzazione)**
Le regole di sostituzione viste nell'esempio vengono astratte nel concetto matematico di **Grammatica a Struttura di Frase**, formalmente definita come una quadrupla $G=\langle V,\Sigma,P,S\rangle$:

* **$V$ (Vocabolario Totale):** Un alfabeto finito che comprende tutti i simboli.
* **$\Sigma$ (Alfabeto dei Simboli Terminali):** Sottoinsieme di $V$ ($\Sigma\subseteq V$), comprende i simboli definitivi (es. "Aldo", la virgola, la "e").
* **$P$ (Insieme delle Produzioni):** Un insieme finito di regole di riscrittura della forma $\alpha\rightarrow\beta$.
    * La parte sinistra $\alpha \in V^{*} \setminus \Sigma^{*}$  = deve contenere almeno una variabile.
    * La parte destra $\beta \in V^{*}$ =  può essere una sequenza qualsiasi di variabili e terminali.
* **$S$ (Simbolo Iniziale o Assioma):** È il punto di partenza del processo generativo ($S\in N$).
* **$N$ (Variabili):** Le lettere dell'insieme $N=V\backslash\Sigma$ si dicono variabili (i nostri "segnaposto").

**Convenzioni Notazionali**
* Variabili ($N$): Si indicano con lettere maiuscole come A, B, C, X, Y, Z.
* Simboli Terminali ($\Sigma$): Si indicano con lettere minuscole come a, b, c.
* Parole sull'alfabeto $V$: Si indicano con lettere greche come $\alpha, \beta, \gamma$.
* Parole sui terminali $\Sigma$: Si indicano con $u, v, w$.

**Il Linguaggio Generato**
Il processo con cui una grammatica crea parole si basa sulla **derivazione:**
* **Conseguenza Diretta ($\alpha\Rightarrow\beta$):** $\beta$ è una conseguenza diretta di $\alpha$ se posso applicare una singola produzione $\gamma\rightarrow\gamma^{\prime}$ (appartenente a $P$) a un frammento di $\alpha$. Formalmente: esistono $\gamma_{1}, \gamma_{2} \in V^{*}$ tali che $\alpha=\gamma_{1}\gamma\gamma_{2}$ e $\beta=\gamma_{1}\gamma^{\prime}\gamma_{2}$.
* **Derivazione ($\alpha\Rightarrow^{*}\beta$):** $\beta$ si deriva da $\alpha$ nella grammatica $G$ se esiste una catena di conseguenze dirette di lunghezza $n \ge 0$ passaggi ($\alpha=\alpha_{0}\Rightarrow\alpha_{1}\Rightarrow\cdot\cdot\cdot\Rightarrow\alpha_{n}=\beta$).
* **Forme Sentenziali:** Tutte le stringhe che si possono derivare a partire dal simbolo iniziale $S$.
* **Linguaggio Generato da $G$:** È l'insieme di tutte e sole le forme sentenziali che sono **prive di variabili** (cioè composte interamente da simboli terminali appartenenti a $\Sigma$).

 ***LEGENDA***

| **Concetto**                  | **Simbolo / Notazione**              | **Descrizione**                                                   |
| ----------------------------- | ------------------------------------ | ----------------------------------------------------------------- |
| **Alfabeto**                  | $\Sigma$                             | Insieme finito di simboli (terminali).                            |
| **Vocabolario Totale**        | $V$                                  | Unione di terminali e variabili ($\Sigma \cup N$).                |
| **Variabili (Non-terminali)** | $N$ (o $V_{N}$)                      | Segnaposto indicati con lettere maiuscole ($A, B, S$).            |
| **Parola Vuota**              | $\epsilon$ (o $\Lambda$)             | Sequenza di lunghezza zero.                                       |
| **Chiusura di Kleene**        | $\Sigma^*$                           | Insieme di tutte le parole possibili su un alfabeto.              |
| **Produzione**                | $\gamma \rightarrow \gamma^{\prime}$ | Regola di riscrittura da una stringa a un'altra.                  |
| **Assioma**                   | $S$                                  | Il simbolo iniziale da cui parte ogni derivazione.                |
| **Conseguenza Diretta**       | $\alpha \Rightarrow \beta$           | Trasformazione di $\alpha$ in $\beta$ tramite **una** produzione. |
| **Conseguenza Derivazione**   | $\alpha \Rightarrow^* \beta$         | Trasformazione tramite **zero o più** passaggi.                   |
| **Lunghezza**                 | $                                    | w                                                                 |
| **Linguaggio generato**       | $L(G)$                               | Insieme delle parole terminali derivabili da $S$.                 |

**Il nostro esempio: Formalizzazione della Grammatica**
Applicando le definizioni formali alla lista di nomi generata nei passaggi precedenti, definiamo i parametri esatti della nostra grammatica:
 1. Le Variabili ($N$)
     L'insieme delle variabili (i nostri segnaposto) è composto da:
      * $\langle frase \rangle$
      * $\langle listaNomi \rangle$
      * $\langle fineLista \rangle$
      * $\langle nome \rangle$

 2. I Simboli Terminali ($\Sigma$)
     L'alfabeto dei simboli terminali (ovvero gli elementi che compariranno nella stringa finale prodotta) comprende:
      * Le lettere delle parole **Aldo**, **Bianca**, **Carlo**
      * La congiunzione **e**
      * La **virgola** (`,`)
      * Lo **spazio**

 3. L'Insieme delle Produzioni ($P$)
     Le regole di sostituzione formali della grammatica sono le seguenti:

* $\langle nome \rangle \rightarrow Aldo$
* $\langle nome \rangle \rightarrow Bianca$
* $\langle nome \rangle \rightarrow Carlo$
* $\langle frase \rangle \rightarrow \langle nome \rangle$
* $\langle frase \rangle \rightarrow \langle listaNomi \rangle \langle fineLista \rangle$
* $\langle listaNomi \rangle \rightarrow \langle nome \rangle$
* $\langle listaNomi \rangle \rightarrow \langle nome \rangle, \langle listaNomi \rangle$
* $\langle nome \rangle \langle fineLista \rangle \rightarrow e \langle nome \rangle$
---
### 26-03-05 lz. 3
Le grammatiche mi permettono di definire dei linguaggi (anche infiniti). Questa è una definizione operativa: la descrizione mi permette di derivare tutte le parole del linguaggio. 
 
***Grammatiche equivalenti:***
**DEF** = Due grammatiche si dicono equivalenti se generano lo stesso linguaggio.  
Capire se una parola appartiene a un linguaggio (ricognizione) è molto importante perché ci permette di capire se il nostro codice è sintatticamente corretto (equivale a verificare se il nostro codice come sequenza di token, è una delle parole appartenenti al linguaggio e alla grammatica) 

A questo punto si pone il problema di generare un linguaggio e questo a sua volta genera 2 problemi:
1. **Problema di ricognizione:**
Prende quindi in input una grammatica e una parola sull'alfabeto determinato dalla mia grammatica, e avrà come output SI se la parola è generata dalla grammatica oppure NO.
 - *Input:* una grammatica (linguaggio) $G = (V, \Sigma, P,S)$ e una parola *w* $\in$ $\Sigma^*$
 - *Output:* SI se w $\in$ $L(G)$, NO altrimenti
2. **Problema di parsing:** 
Prende come input una grammatica, una parola del linguaggio generato, e vuole trovare una derivazione della parola, cioè una sequenza di forme sentenziali ognuna conseguenza diretta della precedente, in maniera tale che partendo dal simbolo iniziale si ottenga la nostra parola.
 - *Input:* una grammatica (linguaggio) $G = (V, \Sigma, P,S)$ e una parola *w* $\in$ $L(G)$
 - *Output:* una derivazione *S*  => $\alpha_1$ => $\alpha_2$ => ... => $\alpha_n$ => *w* di *w* in *G*

Non esiste un algoritmo che risolva questi problemi in caso generale (c'è una perdita di efficienza). Occorre quindi restringersi a classi particolari di grammatiche (restringere la forma delle produzioni). Se questa forma, però, è troppo stretta non posso generare linguaggi complessi, per questo occorre bilanciare:
- Efficienza algoritmi
- Espressività delle grammatiche

La strategia è quindi quella di classificare (**gerarchia di Chomsky**) le nostre grammatiche in base alla forma delle produzioni: produzioni con una forma particolarmente semplice avranno algoritmi di ricognizione e parsing => semplici e veloci ma non in grado di generare linguaggi espressivi. Altri tipi di grammatiche con produzioni di forma più generale avranno la capacità di produrre linguaggi molto espressivi, ma avranno algoritmi di ricognizione e di parsing molto inefficienti o inesistenti.

***Gerarchia di Chomsky:*** divide le grammatiche e i linguaggi generati da esse in 4 tipi (dai più complessi ai più semplici)
**Linguaggi di tipo 0:**
**Def:** le grammatiche a struttura di frase di dicono anche *grammatiche di tipo 0*. I linguaggi generati da grammatiche di tipo 0 si dicono *linguaggi di tipo 0* oppure *linguaggi ricorsivamente numerabili*. Fanno parte di questo livello le grammatiche a struttura di frase. Queste grammatiche hanno la massima espressività possibile: se io ho una procedura effettiva infinita che mi genera una lista di parole, quella lista di parole è un linguaggio di tipo 0.

> [!example] 
> ![[Screenshot From 2026-03-05 10-15-42.png]]

**Linguaggi di tipo 1 (grammatiche sensibili al contesto / contestuali):**
**Def:** una grammatica a struttura di frase si dice sensibile al contesto se le produzioni hanno la forma: $\alpha_1X\alpha_2 \to \alpha_1\beta\alpha_2$, con $X \in N$, $\alpha_1, \alpha_2, \beta \in V^*$ , $\beta \neq \epsilon$.  Ciò significa che le produzioni di queste grammatiche hanno il lato sx che contiene una variabile, mentre il lato dx si ottiene sostituendo quella variabile con una parola qualsiasi purché non vuota. Quindi mentre prima potevamo sostituire una parola con una parola, qui possiamo sostituire solo una lettera all'interno di una parola con una parola, che per di più non deve essere vuota. I linguaggi generati da grammatiche di tipo *q* si dicono *linguaggi di tipo 1*  o anche *sensibili al contesto*.  A differenza dei linguaggi di tipo 0, per questo tipo di linguaggi esistono algoritmi che risolvono i problemi di parsing e ricognizione, ma sono molto costosi.

> [!example]
> ![[Screenshot From 2026-03-05 10-22-40.png]]
> VN / VU è possibile $\iff$ V è seguita da N / U (per questo dipende da contesto)

**Grammatiche monotone (sempre di tipo 1):**
**Def:** una grammatica si dice *monotona* se tutte le produzioni hanno la forma: $\alpha \to \beta$ con $|\alpha| \leq |\beta|$ . Cioè una grammatica in cui tutte le produzioni hanno il lato dx che ha lunghezza maggiore o uguale a quella del lato sx. Quindi non c’è più l’obbligo di sostituire singole variabili, ma c’è la sola regola che il lato dx non può essere più corto di quello sx. Inoltre: 
- Le grammatiche sensibili al contesto sono monotone, perché in esse si sostituisce (all'interno del contesto) una variabile con qualcosa che non è vuoto, quindi quello che otterremo sarà qualcosa di più lungo o al massimo della stessa lunghezza di quello che già abbiamo.
- Non tutte le grammatiche monotone sono sensibili al contesto
- Ogni grammatica monotona è equivalente a una grammatica sensibile al contesto (che genera lo stesso linguaggio)

> [!example]
>  ![[Screenshot From 2026-03-05 10-26-52.png]]

***Sintesi:***
1. I linguaggi di tipo 1 costituiscono una sottoclasse propria dei linguaggi di tipo 0
2. Esistono algoritmi che risolvono il problema di ricognizione e il problema di parsing per grammatiche di tipo 1 (in generale sono molto costosi)
3. Esistono linguaggi di tipo 0 per cui un tale algoritmo non esiste

**Linguaggi di tipo 2 (Grammatiche non contestuali):**
Una grammatica a struttura di frase si dice *non contestuale* o di *tipo 2* se le produzioni hanno la forma: $X \to \beta$, con $X \in N$, $\beta \in V^*$. Dunque una grammatica a struttura di frase si dice non contestuale se le produzioni hanno a sx una singola variabile. I linguaggi generati da grammatiche di tipo 2 si dicono *linguaggi di tipo 2* o anche *non contestuali*. Sono non contestuali perché la variabile può essere sostituita indipendentemente dal contesto in cui appare.

> [!example]
> ![[Screenshot From 2026-03-05 10-33-04.png]]

I linguaggi di tipo 2 costituiscono una sottoclasse propria dei linguaggi di tipo 1. Per esempio $L = \{ a^{n}b^{n}c^{n} | n > 0\}$ è un linguaggio di tipo 1 ma non è di tipo 2.
Gli algoritmi di ricognizione e parsing per le grammatiche di tipo 2 si comportano abbastanza bene, alcuni sono moderatamente efficienti . La maggior parte dei linguaggi di programmazione sono generati da una grammatica non contestuale.

**Linguaggi di tipo 3 Grammatiche regolari):**
**Def:** una grammatica a struttura di frase si dice *regolare* o di *tipo 3* se le produzioni hanno la forma $X \to aY$ oppure $X \to a$, con $X,Y \in N$, $a \in \Sigma$. Dunque Le produzioni hanno una forma ancora più semplice: a sx sempre una variabile, come nelle grammatiche di tipo 2, ma a dx ho un terminale da solo o un terminale seguito da una variabile. I linguaggi generati da grammatiche di tipo 3 si dicono *linguaggi di tipo 3* o anche *regolari*. I linguaggi di tipo 3 costituiscono una sottoclasse dei linguaggi di tipo 2 perché le produzioni a sx hanno una singola variabile. Quindi un linguaggio generato da una grammatica di tipo 3, è anche un linguaggio generato da una grammatica di tipo 2. Anche in questo caso, però, ci sono linguaggi di tipo 2 che non sono linguaggi di tipo 3.

> [!example]
> ![[Pasted image 20260305103801.png]]

> [!info] Gli algoritmi di ricognizione per i linguaggi di tipo 3 esistono e sono estremamente efficienti, sono molto importanti per l’analisi lessicale ovvero l’analisi di token che costituiscono il nostro codice. Questo si ottiene attraverso la Teoria degli Automi a Stati Finiti.

***Automi a stati finiti deterministici:***
Un automa è un dispositivo che può raggiungere un certo numero finito di configurazioni interne. Il nastro di input è diviso in celle su cui si può scrivere l'input e a ogni passo il dispositivo:
- Legge il contenuto di una cella
- Passa in un nuovo stato (dipende esclusivamente dallo stato precedente e dal contenuto della cella letta)
- Passa a esaminare la cella successiva
Terminata la lettura del nastro, in base allo stato raggiunto, accetta o rifiuta la parola letta. Le parole accettate costituiscono un linguaggio di tipo 3 che, a sua volta, è sempre accettato da un automa a stati finiti. Questo algoritmo è lineare e in tempo reale, perché ad ogni cella letta so se tutto ciò che è venuto prima è stato accettato o meno.

**Def:** un *automa a stati finiti deterministico* è una quintupla $A = (Q, \Sigma, \delta, q_0, F)$ ove:
- $Q$ è un insieme finito, detto insieme degli stati = tutte le configurazioni che il dispositivo può assumere
- $\Sigma$ è un alfabeto detto alfabeto di input
- $\delta$: $Q \times \Sigma \to Q$ è la funzione di transizione  = legge con cui i nostri stati cambiano in base ad uno stato precedente e alla lettera scritta sul nastro
- $q_0 \in Q$ è lo stato iniziale
- $F \subseteq Q$ è l'insieme degli stati finiti = se al termine della computazione, il dispositivo si trova in uno stato che è compreso nell'insieme degli stati finali, la parola viene accettata

> [!example] ![[Screenshot From 2026-03-05 21-44-54.png]]

**Funzione di transizione estesa:**
La funzione $\delta$ , affinché non lavori solo con le coppie stato-lettera ma anche con le coppie stato-parola, si estende a $\hat{\delta}$: $Q \times \Sigma^* \to Q$ ponendo:
- $\hat{\delta}(q, \epsilon) = q$, $\forall q \in Q$ 
- $\hat{\delta}(q,va) = \delta(\hat{\delta}(q,v),a)$, $\forall q \in Q$, $v \in \Sigma^*$, $a \in \Sigma$
Per ogni $q \in Q, u, v \in \Sigma^*$,  si ha $\hat{\delta}(q,uv) = \hat{\delta}(\hat{\delta}(q,u),v)$  

> [!example] ![[Screenshot From 2026-03-05 21-47-18.png]]

**Linguaggio accettato:**
**Def:** una parola $w \in \Sigma^*$ è accettata da $A$ se $\hat{\delta}(q_0, w) \in F$. L'insieme delle parole accettate da $A$ si dice linguaggio riconosciuto (o accettato) da $A$ e si denota come $L(A) = \{w \in \Sigma^* | \hat{\delta}(q_0,w) \in F\}$ 

**Grafo di un automa:**
A ogni automa $A = (Q, \Sigma, \delta, q_0, F)$ associamo un grafo diretto con frecce etichettate:
- I vertici sono gli stati
- Le frecce sono le triple ($q, a, \delta (q, a)$) con $q \in Q , a \in \Sigma$ 
- Lo stato iniziale è denotato da una freccia entrante
- Gli stati finali sono identificati dal bordo doppio

> [!example] ![[Screenshot From 2026-03-05 21-53-02.png]]

- Per ogni stato $q$ e ogni lettera $a$ c'è esattamente una freccia con etichetta $a$ che esce da $q$
- C'è un cammino da uno stato $p$ a uno stato $q$ con etichetta $w$ $\iff$ $q = \hat{\delta}(p,w)$
- Le parole accettate dall'automa sono esattamente le etichette dei cammini da $q_0$ a uno stato finale

> [!example] ![[Screenshot From 2026-03-05 21-54-11.png]]

---
### 26-03-06 lz. 4
> [!example] 
> - ![[Screenshot From 2026-03-06 10-06-15.png]]
> - ![[Screenshot From 2026-03-06 10-06-27.png]]
> - ![[Screenshot From 2026-03-06 10-06-39.png]]

***Automi a stati finiti non deterministici:***
Abbiamo detto che un automa è rappresentato da un grafo e che per ogni stato $q$ e ogni lettera $a$ esiste esattamente un arco uscente da $q$ con etichetta $a$. Le parole accettate da $A$ sono le etichette dei cammini dallo stato iniziale a uno stato finale.
Se da questa definizione togliessi l'aspetto degli archi otterrei: una maggiore libertà nella progettazione ma avrei più difficoltà a verificare se una parola è accettata o meno. La soluzione a ciò è un modello non-deterministico che mi permette di fare una conversione automatica al modello deterministico.  Il modello non-deterministico:
- A ogni passo il dispositivo esegue una fra le transizioni possibili (purché esista)
- L'input è accettato se almeno una delle possibili computazioni ha successo (termina in uno stato finale)
- L'input è rifiutato se tutte le possibili computazioni falliscono (terminano in uno stato non finale o prima che il nastro sia esaurito)

**Sintesi** = Possiamo definire Automi a Stati finiti non Deterministici, dei meccanismi che non hanno la restrizione dei modelli deterministici: un meccanismo che può assumere un numero predefinito di configurazioni interne, che ha un nastro di input, ma che ad ogni passo non esegue più una sola transizione ben determinata ma dispone di un ventaglio di possibilità e può quindi eseguire varie transizioni (o nessuna). L’input viene accettato se tra le tante computazioni possibili c’è n’è almeno una che termina in uno stato finale. Quindi una parola viene rifiutata se tutte le computazioni terminano in uno stato non finale o terminano prima che il nastro sia esaurito. 

**Def:** un automa a stati finiti non deterministico è una quintupla $A = \langle Q, \Sigma, \delta, q_0, F \rangle$ dove:
- Q è l'insieme degli stati
- $\Sigma$ è l'alfabeto di input
- $\delta: Q \times \Sigma \to \wp(Q)$ è la funzione di transizione
- $q_0 \in Q$ è lo stato iniziale
- $F \subseteq Q$ è l'insieme degli stati finali

In sostanza vuol dire che nel grafo, nel quale non c’è più il limite al numero di frecce che escono da uno stato, c’è un cammino che ha come etichetta $w$, che parte da uno stato iniziale e arriva in uno stato finale.
Una parola $w = a_1, a_2,...,a_n, (a_1,...a_n \in \Sigma, n \geq 0)$ è *accettata* da $A$ se esistono stati $q_1, q_2,...,q_n \in Q$ tali che $q_i \in \delta(q_{i-1}, a_i), 1 \leq i \leq n, q_n \in F$. L'insieme delle parole accettate da $A$ si dice *linguaggio accettato (o riconosciuto)* da $A$ e si denota con $L(A)$.

A ogni automa non-deterministico $A = \langle Q, \Sigma, \delta, q_0, F \rangle$ associamo un grafo diretto con frecce etichettate:
- I vertici sono gli stati
- Le frecce sono le triple $(q,a,p)$ con $q \in Q$, $a \in \Sigma$, $p \in \delta(q,a)$
- Lo stato iniziale è denotato da una freccia entrante
- Gli stati finali sono identificati dal doppio bordo

Ci si convince facilmente a questo punto che:
- Una parola *w* può essere etichetta di molti cammini uscenti dallo stato iniziale, o anche nessuno
- La parola *w* è accettata $\iff$ c'è un cammino dallo stato iniziale a uno stato finale con etichetta *w*
- Questo non esclude però che vi possano essere altri cammini con etichetta *w* che partono dallo stato iniziale e terminano in uno stato non finale

> [!example] 
> ![[Screenshot From 2026-03-06 11-03-22.png]]
> ![[Screenshot From 2026-03-06 11-03-35.png]]
  >![[Pasted image 20260306110549.png]]
  > 
> Questa computazione non è accettata perché non finisce in uno stato finale, ma questo non significa che la parola non sia accettata perché possono esserci altre computazioni (come la verde) che terminano in uno stato finale

Apparentemente questo tipo di progettazione non ci da nessun vantaggio, perché le computazioni da fare diventerebbero troppe. Supponiamo infatti di avere 2 cammini possibili per ogni scelta, e dovessimo leggere una parola di lunghezza 10, questo vorrebbe dire che dovremmo fare 210 computazioni. Proviamo però a ragionare come segue: se mi trovo in q0 e il mio input è a, posso andare o in q0 o in q1, ecc… Se mi trovo in q1 e l’input è b, andrò in q2,… Quindi ad ogni passo, in base alla colonna che avevo al passo precedente, e alla lettera di input, produco una colonna nuova. Questo procedimento è praticamente uguale a quello di un automa a stati deterministico, perché in base alla configurazione del mio dispositivo e alla lettera di input, determina la configurazione successiva. In questo caso la differenza è che abbiamo delle colonnine di stati.

> [!info] ![[Pasted image 20260306110859.png]]

**Teorema**: Sia $A$ un automa a stati finiti non deterministico, esiste effettivamente un equivalente automa a stati finiti deterministico $A’$ tale che $L (A) = L (A’)$, cioè esiste l’automa deterministico e l’algoritmo per costruire questo automa a partire dall'automa non deterministico. Quindi, prima andremo a disegnare l’automa a stati non deterministico (prima immagine con più percorsi in verde e rosso), cosa che risulta molto più facile rispetto al crearne uno deterministico, e poi sapremo che esiste l’algoritmo in grado di trasformare quel grafo nel grafo di un automa a stati deterministico (seconda immagine con colonne). 
**Dimostrazione:** Gli stati dell’automa deterministico registreranno l’insieme gli stati raggiungibili dalle computazioni dell’automa non deterministico sul medesimo input $\rightarrow$ Accetto se fra questi ce n’è uno finale

**Costruzione:**
Sia $A = \langle Q, \Sigma, \delta, q_0, F \rangle$ l'automa non deterministico. Definisco l'automa deterministico $A' = \langle Q', \Sigma, \delta', s_0, F' \rangle$ come segue: 
- L'insieme degli stati è l'insieme $Q' = \wp(Q)$ costituito dai sottoinsiemi di $Q$ (uno stato dell' $A$ deterministico è in realtà un'insieme di stati dell' $A$ non deterministico)
- Lo stato iniziale è $s_0 = \{q_0\}$ (singleton costituito dallo stato iniziale dell' $A$ non deterministico)
- Gli stati finali sono tutti i sottoinsiemi di $Q$ che contengono almeno un elemento di $F$ (quindi almeno uno stato finale dell' $A$ di partenza), cioè: $F' = \{s \in p(Q) | s \cap F \neq \emptyset\}$ 
- La funzione di transizione $\delta'(s,a) = \displaystyle\bigcup_{q \in s} \delta(q,a), s \in \wp(Q), a \in \Sigma$ (prende in input uno stato dell' $A$, cioè un insieme degli stati dell'automa non-deterministico, e una lettera $a$)
**Dimostrazione:** per ogni $w \in \Sigma^*$ risulta che $\hat{\delta'}(s_0,w) = \{q \in Q | \text{nel grafo di A c'è un cammino da}\ q_0 \ \text{a} \ q \ \text{con etichetta} \ w\}$ (dimostrazione per induzione sulla lunghezza di $w$). Pertanto sono equivalenti: 
- $w\ \text{è accettata da} \ A$
- $\hat{\delta'}(s_0,w) \ \text{contiene uno stato di} \ F$
- $\hat{\delta'}(s_0,w) \in F'$
- $w \ \text{è accettata da} \ A'$
$\implies L(A) = L(A')$ 

Come abbiamo accennato precedentemente, però, sarebbero presenti troppi stati da controllare (se $Q=n$, allora $Q’=2^n$). Per fortuna non tutti questi stati sono necessari: possiamo definire uno stato necessario/accessibile come un cammino nel grafo dell’automa che, dallo stato iniziale, mi porta a $q$. Tutti gli altri stati si dicono inaccessibili.

**Def:** sia $A = \langle Q, \Sigma, \delta, q_0, F \rangle$ un automa a stati finiti deterministico. Uno stato $q \in Q$ si dice *accessibile* se $q = \hat{\delta}(q_0,w)$ per qualche $w \in \Sigma^*$

Se io elimino tutti questi stati inaccessibili e restringo la funzione di transizione (delta cappello) solo agli stati accessibili, ottengo un automa con un numero di stati più piccolo, equivalente a quello di partenza, che accetta lo stesso linguaggio di quello di partenza.

> [!example] ![[Screenshot From 2026-03-06 11-26-56.png]]
> Nell’esempio vediamo come lo stato q1 è stato eliminato perché, iniziando da q0 (stato di partenza), non c’è nessun cammino che mi porta in q1.

Attraverso un Algoritmo di Determinizzazione andremo a progettare un automa che prevede solo gli stati accessibili, e di conseguenza avrà un numero di stati molto inferiore di 2 n.

**Algoritmo di determinizzazione:**
1. Prima di tutto inserisco nella lista degli stati, lo stato iniziale dell’automa deterministico, che come abbiamo visto prima corrisponde all'insieme costituito dal solo stato iniziale dell’automa non deterministico di partenza: $s_0=\{q_0\}$ 
2. Per ogni stato $r$ della lista (che per ora contiene il solo elemento $q0$, singleton) e ogni lettera $a \in \Sigma$ 
3. Calcolo il valore della funzione di transizione sulla coppia stato-lettera: $s= \delta'(r, a)$; Una volta calcolato ho due possibilità: o lo stato è già nella lista degli stati, oppure no.   
4. Se $s$ non è nella lista degli stati allora lo appendo alla lista degli stati. 
5. Se $s$ contiene uno stato finale di $A$, aggiungo $s$ alla lista degli stati finali. 

> [!example] ![[Screenshot From 2026-03-06 11-32-20.png]]
> ![[Screenshot From 2026-03-06 11-32-32.png]]
> ![[Screenshot From 2026-03-06 11-33-08.png]]
> Come possiamo notare, l’automa deterministico ha più stati rispetto a quelli che avevamo trovato noi manualmente nel primo esempio abac (pag.11), successivamente vedremo infatti un altro algoritmo in grado di trovare partendo da un automa non deterministico, un automa deterministico con il minimo numero di stati possibile (automa minimo).

> [!info] Caso difficile:
> ![[Pasted image 20260308174311.png]]

---
