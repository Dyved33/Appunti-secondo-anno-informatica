### **Compilatori e Interpreti**
La traduzione del codice sorgente in un linguaggio comprensibile dalla macchina può avvenire in modalità differenti (è anche possibile la traduzione da un linguaggio di programmazione a un altro):
* **Compilatore:** Esegue la traduzione in blocco *prima* dell'esecuzione. Prende in input il **codice sorgente** e, tramite il traduttore, restituisce il **codice oggetto** (oppure delle notifiche di errore).
* **Interprete:** Esegue la traduzione *istruzione per istruzione*, simultaneamente all'esecuzione del programma. Prende in input sia il **codice sorgente** sia l'**input** dell'utente, producendo direttamente l'**output**.

Esistono compilatori Ibridi (Es. Java) che combinano i due approcci:
1. Traduttore analizza il codice sorgente e lo traduce in un **programma intermedio** (ad esempio, il *Java bytecode*).
2. Successivamente, una **macchina virtuale** (ad esempio, la JVM) riceve gli input dell'utente e *interpreta* il codice intermedio per produrre l'output finale.
### **Le Fasi della Compilazione: Dall'Analisi alla Sintesi**
La compilazione si divide concettualmente in due macro-fasi principali:
* **Fase Analitica (Analisi):** Riguarda la comprensione del codice in ingresso. Ha lo scopo di costruire l'**Albero Sintattico Astratto (AST)** partendo dal codice sorgente.
* **Fase Sintetica (Sintesi):** Riguarda la generazione dell'output. Si occupa di produrre il codice oggetto finale a partire dall'AST generato.

**Le 8 Fasi del Sistema di Compilazione**
Il passaggio dal codice sorgente al codice macchina attraversa stadi sequenziali, che comunicano tra loro avvalendosi di strutture condivise come la **tabella dei simboli** e il modulo di **gestione degli errori**.

**Fase di analisi:**
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

3. **Analisi Semantica Statica:** Utilizza l'albero sintattico astratto e la tabella dei simboli per verificare che il programma sorgente sia semanticamente coerente con la definizione del linguaggio.
   * **Type checking:** Verifica che ogni operatore abbia tutti gli operandi del tipo corretto. Se necessario, può eseguire la **coercizione** (conversione forzata del tipo).
   * *Esempio in C:* È permesso moltiplicare un intero per un reale; il compilatore converte silenziosamente l'intero in un reale.
   * *Esempio in Pascal:* Il costrutto `WHILE` richiede obbligatoriamente un valore `BOOLEAN`. L'espressione `WHILE X DO X := X - 1 ;` farà segnalare un errore all'analizzatore semantico.

**Fase di sintesi:**
5. **Generazione del Codice Intermedio:** Dall'AST e dalla tabella dei simboli si ottiene un codice intermedio basato su istruzioni elementari, facili da tradurre in codice macchina. È **indipendente dall'architettura** (non specifica gli indirizzi fisici in memoria né i registri).
   * *Esempio pratico:* L'espressione precedente viene scomposta in operazioni atomiche: `t1 = inttofloat(60)`, `t2 = id3 * t1`, `t3 = id2 + t2`, `id1 = t3`.
5. **Ottimizzazione:** L'ottimizzatore elabora il codice intermedio per ridurre il tempo o lo spazio necessario all'esecuzione.
   > [!NOTE] Taglio delle Istruzioni
   > L'ottimizzatore snellisce il codice eliminando passaggi inutili. Nel nostro esempio, riduce le operazioni accorpando i calcoli: `t2 = id3 * 60.0` e `id1 = id2 + t2`.
6. **Generazione del Codice Oggetto:** È la parte più complessa e dipende strettamente dall'architettura della piattaforma hardware. Il compilatore genera codice assembly o file oggetto ricollocabili. I passaggi chiave sono:
   * Fissare le locazioni di memoria esatte per i dati.
   * Generare il codice di basso livello per accedere a tali dati.
   * Selezionare i registri fisici della CPU per i calcoli intermedi.
   > [!IMPORTANT] Il ruolo dell'Assemblatore e del Linker
   > A rigore, l'Assemblatore (se il compilatore emette assembly) e il Linker sono strumenti esterni che intervengono al termine della compilazione. Il Linker si occuperà di collegare il codice oggetto generato in questa fase con eventuali librerie esterne, producendo il file eseguibile finale.
5. **Ottimizzazione Peep-hole:** Un'ulteriore e ultimissima fase di ottimizzazione legata alle specificità e ai "trucchi" dell'architettura hardware di destinazione.
### **L'Architettura Software del Compilatore**
L'organizzazione interna del compilatore definisce come i moduli comunichino effettivamente tra loro durante il processo di traduzione.
* **Il Nucleo Centrale (Parser & Semantica):** la parte centrale di tutto il processo è l'**analizzatore sintattico**. In questa architettura, esso realizza contestualmente anche l'**analisi semantica statica**.
   > [!IMPORTANT] Avanzamento Parallelo
   > Dato che le due funzioni sono accorpate nel nucleo, l'analisi sintattica e l'analisi semantica avanzano parallelamente.
* **Gestione Input e Lettura:** Il sistema inizia con il modulo **gestore input**, che riceve il codice sorgente e può generare un "listato sorgente". Per formare i token, l'analizzatore lessicale deve leggere il codice un carattere alla volta invocando la funzione `readChar()` direttamente sul gestore input.
* **Dialogo Scanner-Parser:** L'analizzatore sintattico coordina il lavoro chiedendo i dati all'analizzatore lessicale tramite la chiamata alla funzione `prossimoToken()`. Questa funzione, implementata nello scanner, restituisce il `token` richiesto per continuare a costruire l'albero.
* **Gestione degli Errori:** Il **notificatore degli errori** riceve comunicazioni dirette e incrociate dal gestore input, dall'analizzatore lessicale e dal blocco sintattico/semantico ogni volta che viene rilevata un'anomalia. Alla fine, produce in output la "lista degli errori".
* **Generazione del Codice Finale:** Completata l'analisi senza errori bloccanti, è sempre l'analizzatore sintattico ad attivare la sintesi. Tramite chiamate a funzioni opportune, innesca il **generatore del codice**, che a sua volta emetterà il **codice oggetto** definitivo.
### Il problema dell'infinito e il Teorema di Cantor
Per descrivere qualsiasi cosa, che sia la Divina Commedia o il codice genetico umano, usiamo delle parole. 
In informatica teorica, un **linguaggio è definito come un insieme di parole**.
* **Linguaggi Finiti:** Se un linguaggio è finito, la sua descrizione è banale poiché è possibile elencarne tutti gli elementi uno per uno. Un esempio è $L=\{a,ab,abb,abbb\}$.
* **Linguaggi Infiniti:** Il problema sorge quando cerchiamo di descrivere linguaggi con infinite parole.

Il **Teorema di Cantor** ci pone un limite matematico insuperabile: non esiste una funzione iniettiva dall'insieme delle parti $\mathcal{P}(\Sigma^{*})$ all'insieme delle parole $\Sigma^{*}$. Di conseguenza, non è possibile associare a *ciascun* linguaggio possibile una parola univoca (la sua "descrizione") che lo caratterizzi in modo perfetto. In altre parole, non possiamo inventare una "parola descrittiva" univoca per tutti i possibili linguaggi infiniti esistenti (linguaggio infinito non può essere descritto da una descrizione finita).

> [!NOTE] L'approccio Informatico
> Nonostante il Teorema di Cantor, in informatica non abbiamo bisogno di descrivere *tutti* i linguaggi infiniti. È comunque possibile fornire **descrizioni finite** per *alcuni* linguaggi infiniti che ci sono particolarmente utili. Ad esempio, possiamo usare una regola descrittiva semplice come "tutte le parole sull'alfabeto $\{a,b\}$ di lunghezza dispari", oppure la notazione formale $\{a^{n}b^{n}|n\ge0\}$.

Un costrutto fondamentale per generare infinità è la **chiusura di Kleene** di un linguaggio $L$ (indicata con $L^{*}$). Essa è costituita dalla parola vuota, dalle parole appartenenti a $L$ e da tutte le parole che si ottengono concatenando due o più parole di $L$. Se $L$ è un linguaggio finito, il linguaggio derivato $L^{*}$ è ben definito e, in generale, infinito.
### Alfabeti , Parole e le loro Proprietà (Definizioni Formali):
Per poter costruire delle regole rigorose, dobbiamo prima definire matematicamente gli elementi base:
* **Alfabeto ($\Sigma$):** È un insieme finito e non vuoto di simboli. I suoi elementi sono chiamati **lettere**.
    * *Esempi di alfabeti:* $\Sigma_{0}=\{a,b\}$, $\Sigma_{1}=\{0,1\}$, $\Sigma_{2}=\{a,b,c\}$ e l'alfabeto esadecimale $\Sigma_{3}=\{0,1,2,3,4,5,6,7,8,9,A,B,C,D,E,F\}$.
* **Parola:** Ogni sequenza finita di lettere appartenenti a un alfabeto $\Sigma$ è detta parola sull'alfabeto $\Sigma$. 
    * *Esempi:* "a", "abb", "ababbabb" sono parole su $\Sigma_{0}$; "2020" e "CD078B" sono parole su $\Sigma_{3}$.
* **L'insieme Universo ($\Sigma^{*}$):** L'insieme di *tutte* le parole generabili sull'alfabeto $\Sigma$ è denotato con $\Sigma^{*}$.
* **Parola Vuota:** Esiste anche la sequenza composta da zero lettere, che si denota con $\epsilon$ (oppure $\Lambda$) e si dice **parola vuota**.
* **Lunghezza ($|u|$):** Data una parola $u=a_{1}a_{2}\cdot\cdot\cdot a_{k}$ (con $k\ge0$), l'intero $k$ rappresenta la sua lunghezza e si denota con $|u|$ oppure $l(u)$.
    * *Esempi:* $|abb|=3$, $|\epsilon|=0$.

**Operazioni sulle Parole:**
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
### Linguaggi e grammatiche
**Linguaggio Formale:** Ogni sottoinsieme di $\Sigma^{*}$ si dice **linguaggio formale** (o semplicemente linguaggio) sull'alfabeto $\Sigma$. Possono esserci linguaggi limitati come $L_{1}=\{a,ab,abb\}$, linguaggi infiniti come $L_{2}=\{ab^{n}a|n\ge0\}$, l'insieme vuoto $L_{3}=\emptyset$ o l'intero universo $L_{4}=\Sigma^{*}$.

**L'Intuizione: Il Gioco delle Sostituzioni:** Per capire operativamente come le "regole finite" possano generare "linguaggi infiniti", analizziamo un esempio pratico basato sul linguaggio naturale. Vogliamo generare liste di nomi come "Aldo, Bianca e Carlo". Le regole logiche che governano queste frasi sono:
1.  I nomi sono separati da virgola, tranne gli ultimi due che sono separati dalla congiunzione 'e'.
2.  Sono ammesse ripetizioni di nomi (es. "Aldo, Bianca, Aldo...").
3.  Non c'è alcun limite alla lunghezza della lista.

Per istruire una macchina, dobbiamo usare dei **segnaposto** (come i termini "nome" o "frase") ed eseguire sostituzioni meccaniche secondo regole predefinite. 

> [!IMPORTANT] Il Problema delle Regole Globali vs Locali
> Un primo tentativo di creare queste regole umane includerebbe una direttiva del tipo: "prima di terminare, se sono presenti virgole, l'ultima va sostituita con la congiunzione e".
> Questa regola, tuttavia, è anomala rispetto alle altre: impone di analizzare l'intera frase simultaneamente per trovare "l'ultima virgola". Un calcolatore preferisce invece **regole di sostituzione rigorose e puramente locali**, dove un simbolo viene sostituito ciecamente a prescindere dal contesto globale.

**La Soluzione Formale: Ricorsione e Segnaposto:** La soluzione "informatica" parte sempre da un segnaposto iniziale (es. `frase`) e procede sostituendolo finché non rimangono solo parole e punteggiatura reali. 
1. **Regole di Base:** Il segnaposto $\langle nome \rangle$ può essere sostituito da Aldo, Bianca o Carlo.
2. **Inizio:** $\langle frase \rangle$ può diventare un singolo $\langle nome \rangle$ oppure il blocco $\langle lista Nomi \rangle \langle fineLista \rangle$.
3. **Il Motore dell'Infinito (Ricorsione):** $\langle lista Nomi \rangle$ può essere sostituito da $\langle nome \rangle$ oppure da $\langle nome \rangle, \langle lista Nomi \rangle$. 

> [!EXAMPLE] Come si genera l'infinito
    > Contenendo se stessa nella parte destra dell'assegnazione, può essere applicata all'infinito per aggiungere quante virgole e quanti nomi desideriamo, creando una lunghezza illimitata a partire da una singola regola.

4. **Chiusura:** Il blocco $\langle nome \rangle \langle fine Lista \rangle$ si sostituisce rigorosamente con $e \langle nome \rangle$. Questo trucco garantisce che, alla fine del processo meccanico, l'ultimo elemento sia sempre preceduto dalla 'e' e non dalla virgola. 

Ci si arresta solo quando i segnaposto non compaiono più. <u>Il concetto chiave è che un linguaggio di complessità e lunghezza infinita può essere generato meccanicamente e senza ambiguità tramite un set ristretto di regole di sostituzione locali</u>.

**Grammatiche a Struttura di Frase (Formalizzazione):** Le regole di sostituzione viste nell'esempio vengono astratte nel concetto matematico di **Grammatica a Struttura di Frase**, formalmente definita come una quadrupla $G=\langle V,\Sigma,P,S\rangle$:
* **$V$ (Vocabolario Totale):** Un alfabeto finito che comprende tutti i simboli.
* **$\Sigma$ (Alfabeto dei Simboli Terminali):** Sottoinsieme di $V$ ($\Sigma\subseteq V$), comprende i simboli definitivi, cioè quelli usati (es. "Aldo", la virgola, la "e").
* **$P$ (Insieme delle Produzioni):** Un insieme finito di regole di riscrittura della forma $\alpha\rightarrow\beta$.
    * La parte sinistra $\alpha \in V^{*} \setminus \Sigma^{*}$  = deve contenere almeno una variabile.
    * La parte destra $\beta \in V^{*}$ =  può essere una sequenza qualsiasi di variabili e terminali.
* **$S$ (Simbolo Iniziale o Assioma):** È il punto di partenza del processo generativo ($S\in N$).
* **$N$ (Variabili):** Le lettere dell'insieme $N=V\backslash\Sigma$ si dicono variabili (i nostri "segnaposto").

**Convenzioni Notazionali:** 
- Variabili ($N$): Si indicano con lettere maiuscole come A, B, C, X, Y, Z.
* Simboli Terminali ($\Sigma$): Si indicano con lettere minuscole come a, b, c.
* Parole sull'alfabeto $V$: Si indicano con lettere greche come $\alpha, \beta, \gamma$.
* Parole sui terminali $\Sigma$: Si indicano con $u, v, w$.

**Il Linguaggio Generato:** Il processo con cui una grammatica crea parole si basa sulla **derivazione:**
* **Conseguenza Diretta ($\alpha\Rightarrow\beta$):** $\beta$ è una conseguenza diretta di $\alpha$ se posso applicare una singola produzione $\gamma\rightarrow\gamma^{\prime}$ (appartenente a $P$) a un frammento di $\alpha$. Formalmente: esistono $\gamma_{1}, \gamma_{2} \in V^{*}$ tali che $\alpha=\gamma_{1}\gamma\gamma_{2}$ e $\beta=\gamma_{1}\gamma^{\prime}\gamma_{2}$.
* **Derivazione ($\alpha\Rightarrow^{*}\beta$):** $\beta$ si deriva da $\alpha$ nella grammatica $G$ se esiste una catena di conseguenze dirette di lunghezza $n \ge 0$ passaggi ($\alpha=\alpha_{0}\Rightarrow\alpha_{1}\Rightarrow\cdot\cdot\cdot\Rightarrow\alpha_{n}=\beta$).
* **Forme Sentenziali:** Tutte le stringhe che si possono derivare a partire dal simbolo iniziale $S$.
* **Linguaggio Generato da $G$:** È l'insieme di tutte e sole le forme sentenziali che sono **prive di variabili** (cioè composte interamente da simboli terminali appartenenti a $\Sigma$).

 **LEGENDA**

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
| **Lunghezza**                 | $\|u\|$                              | Lunghezza finita di una parola                                    |
| **Linguaggio generato**       | $L(G)$                               | Insieme delle parole terminali derivabili da $S$.                 |

**Il nostro esempio: Formalizzazione della Grammatica:** Applicando le definizioni formali alla lista di nomi generata nei passaggi precedenti, definiamo i parametri esatti della nostra grammatica:
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

Le grammatiche mi permettono di definire dei linguaggi (anche infiniti). Questa è una definizione operativa: la descrizione mi permette di derivare tutte le parole del linguaggio. 
 
**Grammatiche equivalenti:** 
*Definizione:* Due grammatiche si dicono equivalenti se generano lo stesso linguaggio.  
Capire se una parola appartiene a un linguaggio (ricognizione) è molto importante perché ci permette di capire se il nostro codice è sintatticamente corretto (equivale a verificare se il nostro codice come sequenza di token, è una delle parole appartenenti al linguaggio e alla grammatica). A questo punto si pone la di generare un linguaggio e questo a sua volta genera 2 problemi:
1. **Problema di ricognizione (la parola appartiene al linguaggio?):** Prende quindi in input una grammatica e una parola sull'alfabeto determinato dalla mia grammatica, e avrà come output SI se la parola è generata dalla grammatica oppure NO.
 - *Input:* una grammatica (linguaggio) $G = (V, \Sigma, P,S)$ e una parola *w* $\in$ $\Sigma^*$
 - *Output:* SI se w $\in$ $L(G)$, NO altrimenti.
2. **Problema di parsing (come è stata generata la parola?):** Prende come input una grammatica, una parola del linguaggio generato, e vuole trovare una derivazione della parola, cioè una sequenza di forme sentenziali ognuna conseguenza diretta della precedente, in maniera tale che partendo dal simbolo iniziale si ottenga la nostra parola.
 
 - *Input:* una grammatica (linguaggio) $G = (V, \Sigma, P,S)$ e una parola *w* $\in$ $L(G)$
 - *Output:* una derivazione *S*  => $\alpha_1$ => $\alpha_2$ => ... => $\alpha_n$ => *w* di *w* in *G*

Non esiste un algoritmo che risolva questi problemi in caso generale (c'è una perdita di efficienza). Occorre quindi restringersi a classi particolari di grammatiche (restringere la forma delle produzioni). Se questa forma, però, è troppo stretta non posso generare linguaggi complessi, per questo occorre bilanciare:
- Efficienza algoritmi
- Espressività delle grammatiche

La strategia è quindi quella di classificare (**gerarchia di Chomsky**) le nostre grammatiche in base alla forma delle produzioni: produzioni con una forma particolarmente semplice avranno algoritmi di ricognizione e parsing => semplici e veloci ma non in grado di generare linguaggi espressivi. Altri tipi di grammatiche con produzioni di forma più generale avranno la capacità di produrre linguaggi molto espressivi, ma avranno algoritmi di ricognizione e di parsing molto inefficienti o inesistenti.
### Gerarchia di Chomsky e tipi di linguaggio
Divide le grammatiche e i linguaggi generati da esse in 4 tipi (dai più complessi ai più semplici)

**Linguaggi di tipo 0:** 

*Definizione:* le grammatiche a struttura di frase (in base alle loro regole possono scrivere intere frasi e non solo singoli caratteri) si dicono anche *grammatiche di tipo 0*. I linguaggi generati da grammatiche di tipo 0 si dicono *linguaggi di tipo 0* oppure *linguaggi ricorsivamente numerabili*. Fanno parte di questo livello le grammatiche a struttura di frase. Queste grammatiche hanno la massima espressività possibile: <u>se io ho una procedura effettiva infinita che mi genera una lista di parole, quella lista di parole è un linguaggio di tipo 0.</u>

> [!example] 
> ![[Screenshot From 2026-03-05 10-15-42.png|300]]

**Linguaggi di tipo 1 (grammatiche sensibili al contesto / contestuali):**

*Definizione:* una grammatica a struttura di frase si dice sensibile al contesto se le produzioni hanno la forma: $\alpha_1X\alpha_2 \to \alpha_1\beta\alpha_2$, con $X \in N$, $\alpha_1, \alpha_2, \beta \in V^*$ , $\beta \neq \epsilon$.  Ciò significa che le produzioni di queste grammatiche hanno il lato sinistra che contiene una variabile, mentre il lato destra si ottiene sostituendo quella variabile con una parola qualsiasi purché non vuota. Quindi mentre prima potevamo sostituire una parola con una parola, qui possiamo sostituire solo una lettera all'interno di una parola con una parola, che per di più non deve essere vuota (posso sostituire una variabile solo se è contenuta in un  contesto). I linguaggi generati da grammatiche di tipo *q* si dicono *linguaggi di tipo 1*  o anche *sensibili al contesto*.  A differenza dei linguaggi di tipo 0, per questo tipo di linguaggi esistono algoritmi che risolvono i problemi di parsing e ricognizione, ma sono molto costosi.

<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Screenshot From 2026-03-05 10-22-40.png" style="width: 100%; border-radius: 8px;">
  </div>
  <div style="flex: 1.5;">
  VN / VU è possibile se e solo so V è seguita da N / U (per questo dipende da contesto)
  </div>
</div>

> [!example]
> ![[Screenshot From 2026-03-05 10-22-40.png|300]]

**Grammatiche monotone (sempre di tipo 1):**

*Definizione:* una grammatica si dice *monotona* se tutte le produzioni hanno la forma: $\alpha \to \beta$ con $|\alpha| \leq |\beta|$. Cioè una grammatica in cui tutte le produzioni hanno il lato destra che ha lunghezza maggiore o uguale a quella del lato sinistra. Quindi non c'è più l'obbligo di sostituire singole variabili, ma c'è la sola regola che il lato destra non può essere più corto di quello sinistra. Inoltre: 
- Le grammatiche sensibili al contesto sono monotone, perché in esse si sostituisce (all'interno del contesto) una variabile con qualcosa che non è vuoto, quindi quello che otterremo sarà qualcosa di più lungo o al massimo della stessa lunghezza di quello che già abbiamo.
- Non tutte le grammatiche monotone sono sensibili al contesto
- Ogni grammatica monotona è equivalente a una grammatica sensibile al contesto (che genera lo stesso linguaggio)

> [!example] ![[Screenshot From 2026-03-05 10-26-52.png|300]]

**Sintesi:**
1. I linguaggi di tipo 1 costituiscono una sottoclasse propria dei linguaggi di tipo 0
2. Esistono algoritmi che risolvono il problema di ricognizione e il problema di parsing per grammatiche di tipo 1 (in generale sono molto costosi)
3. Esistono linguaggi di tipo 0 per cui un tale algoritmo non esiste

**Linguaggi di tipo 2 (Grammatiche non contestuali):** Una grammatica a struttura di frase si dice *non contestuale* o di *tipo 2* se le produzioni hanno la forma: $X \to \beta$, con $X \in N$, $\beta \in V^*$. Dunque una grammatica a struttura di frase si dice non contestuale se le produzioni hanno a sinistra una singola variabile. I linguaggi generati da grammatiche di tipo 2 si dicono *linguaggi di tipo 2* o anche *non contestuali*. Sono non contestuali perché la variabile può essere sostituita indipendentemente dal contesto in cui appare.

> [!example]
> ![[Screenshot From 2026-03-05 10-33-04.png|300]]

I linguaggi di tipo 2 costituiscono una sottoclasse propria dei linguaggi di tipo 1. Per esempio $L = \{ a^{n}b^{n}c^{n} | n > 0\}$ è un linguaggio di tipo 1 ma non è di tipo 2.
Gli algoritmi di ricognizione e parsing per le grammatiche di tipo 2 si comportano abbastanza bene, alcuni sono moderatamente efficienti . La maggior parte dei linguaggi di programmazione sono generati da una grammatica non contestuale.

**Linguaggi di tipo 3 Grammatiche regolari):** 

*Definizione:* una grammatica a struttura di frase si dice *regolare* o di *tipo 3* se le produzioni hanno la forma $X \to aY$ oppure $X \to a$, con $X,Y \in N$, $a \in \Sigma$. Dunque Le produzioni hanno una forma ancora più semplice: a sinistra sempre una variabile, come nelle grammatiche di tipo 2, ma a destra ho un terminale da solo o un terminale seguito da una variabile. I linguaggi generati da grammatiche di tipo 3 si dicono *linguaggi di tipo 3* o anche *regolari*. I linguaggi di tipo 3 costituiscono una sottoclasse dei linguaggi di tipo 2 perché le produzioni a sinistra hanno una singola variabile. Quindi un linguaggio generato da una grammatica di tipo 3, è anche un linguaggio generato da una grammatica di tipo 2. Anche in questo caso, però, ci sono linguaggi di tipo 2 che non sono linguaggi di tipo 3.

> [!example]
> ![[Pasted image 20260305103801.png|300]]

> [!info] Gli algoritmi di ricognizione per i linguaggi di tipo 3 esistono e sono estremamente efficienti, sono molto importanti per l'analisi lessicale ovvero l'analisi di token che costituiscono il nostro codice. Questo si ottiene attraverso la Teoria degli Automi a Stati Finiti.
### Automi a stati finiti deterministici
Un automa è un dispositivo che può raggiungere un certo numero finito di configurazioni interne. Il nastro di input è diviso in celle su cui si può scrivere l'input (scrivo la stringa in ingresso e la divido nei suoi caratteri) e a ogni passo il dispositivo:
- Legge il contenuto di una cella
- Passa in un nuovo stato (dipende esclusivamente dallo stato precedente e dal contenuto della cella letta)
- Passa a esaminare la cella successiva

Terminata la lettura del nastro, in base allo stato raggiunto (se è uno stato accettante o meno) , accetta o rifiuta la parola letta. Le parole accettate costituiscono un linguaggio di tipo 3 che, a sua volta, è sempre accettato da un automa a stati finiti. Questo algoritmo è lineare e in tempo reale, perché ad ogni cella letta so se tutto ciò che è venuto prima è stato accettato o meno.

*Definizione:* un *automa a stati finiti deterministico* è una quintupla $A = (Q, \Sigma, \delta, q_0, F)$ ove:
- $Q$ è un insieme finito, detto insieme degli stati = tutte le configurazioni che il dispositivo può assumere
- $\Sigma$ è un alfabeto detto alfabeto di input
- $\delta$: $Q \times \Sigma \to Q$ è la funzione di transizione  = legge con cui i nostri stati cambiano in base ad uno stato precedente e alla lettera scritta sul nastro
- $q_0 \in Q$ è lo stato iniziale
- $F \subseteq Q$ è l'insieme degli stati finiti = se al termine della computazione, il dispositivo si trova in uno stato che è compreso nell'insieme degli stati finali, la parola viene accettata

> [!example] ![[Screenshot From 2026-03-05 21-44-54.png|300]]

**Funzione di transizione estesa:** La funzione $\delta$ , affinché non lavori solo con le coppie *stato-lettera* ma anche con le coppie *stato-parola*, si estende a $\hat{\delta}$: $Q \times \Sigma^* \to Q$ ponendo:
- $\hat{\delta}(q, \epsilon) = q$, $\forall q \in Q$ 
- $\hat{\delta}(q,va) = \delta(\hat{\delta}(q,v),a)$, $\forall q \in Q$, $v \in \Sigma^*$, $a \in \Sigma$

Per ogni $q \in Q, u, v \in \Sigma^*$,  si ha $\hat{\delta}(q,uv) = \hat{\delta}(\hat{\delta}(q,u),v)$  

> [!example] ![[Screenshot From 2026-03-05 21-47-18.png|300]]

**Linguaggio accettato:**
*Definizione:* una parola $w \in \Sigma^*$ è accettata da $A$ se $\hat{\delta}(q_0, w) \in F$. L'insieme delle parole accettate da $A$ si dice linguaggio riconosciuto (o accettato) da $A$ e si denota come $L(A) = \{w \in \Sigma^* | \hat{\delta}(q_0,w) \in F\}$ 

**Grafo di un automa:** A ogni automa $A = (Q, \Sigma, \delta, q_0, F)$ associamo un grafo diretto con frecce etichettate:
- I vertici sono gli stati
- Le frecce sono le triple ($q, a, \delta (q, a)$) con $q \in Q , a \in \Sigma$ 
- Lo stato iniziale è denotato da una freccia entrante
- Gli stati finali sono identificati dal bordo doppio

> [!example] ![[Screenshot From 2026-03-05 21-53-02.png|300]]

- Per ogni stato $q$ e ogni lettera $a$ c'è esattamente una freccia con etichetta $a$ che esce da $q$
- C'è un cammino da uno stato $p$ a uno stato $q$ con etichetta $w$ $\iff$ $q = \hat{\delta}(p,w)$
- Le parole accettate dall'automa sono esattamente le etichette dei cammini da $q_0$ a uno stato finale
> [!example] ![[Screenshot From 2026-03-05 21-54-11.png|300]]
> 
> - ![[Screenshot From 2026-03-06 10-06-15.png|300]]
> - ![[Screenshot From 2026-03-06 10-06-27.png|300]]
> - ![[Screenshot From 2026-03-06 10-06-39.png|300]]
### Automi a stati finiti non deterministici
Abbiamo detto che un automa è rappresentato da un grafo e che per ogni stato $q$ e ogni lettera $a$ esiste esattamente un arco uscente da $q$ con etichetta $a$. Le parole accettate da $A$ sono le etichette dei cammini dallo stato iniziale a uno stato finale.
Se da questa definizione togliessi l'aspetto degli archi (cioè che da ogni vertice può uscire un solo arco per ogni etichetta possibile) otterrei: una maggiore libertà nella progettazione ma avrei più difficoltà a verificare se una parola è accettata o meno. La soluzione a ciò è un modello non-deterministico che mi permette di fare una conversione automatica al modello deterministico. Il modello non-deterministico:
- A ogni passo il dispositivo esegue una fra le transizioni possibili (purché esista)
- L'input è accettato se almeno una delle possibili computazioni ha successo (termina in uno stato finale)
- L'input è rifiutato se tutte le possibili computazioni falliscono (terminano in uno stato non finale o prima che il nastro sia esaurito)

> [!info] **Sintesi:** 
> Possiamo definire Automi a Stati finiti non Deterministici, dei meccanismi che non hanno la restrizione dei modelli deterministici: un meccanismo che può assumere un numero predefinito di configurazioni interne, che ha un nastro di input, ma che ad ogni passo non esegue più una sola transizione ben determinata ma dispone di un ventaglio di possibilità e può quindi eseguire varie transizioni (o nessuna). L'input viene accettato se tra le tante computazioni possibili c'è n'è almeno una che termina in uno stato finale. Quindi una parola viene rifiutata se tutte le computazioni terminano in uno stato non finale o terminano prima che il nastro sia esaurito. 

*Definizione:* un automa a stati finiti non deterministico è una quintupla $A = \langle Q, \Sigma, \delta, q_0, F \rangle$ dove:
- Q è l'insieme degli stati
- $\Sigma$ è l'alfabeto di input
- $\delta: Q \times \Sigma \to \wp(Q)$ è la funzione di transizione
- $q_0 \in Q$ è lo stato iniziale
- $F \subseteq Q$ è l'insieme degli stati finali

In sostanza vuol dire che nel grafo, nel quale non c'è più il limite al numero di frecce che escono da uno stato, c'è un cammino che ha come etichetta $w$, che parte da uno stato iniziale e arriva in uno stato finale. Una parola $w = a_1, a_2,...,a_n, (a_1,...a_n \in \Sigma, n \geq 0)$ è *accettata* da $A$ se esistono stati $q_1, q_2,...,q_n \in Q$ tali che $q_i \in \delta(q_{i-1}, a_i), 1 \leq i \leq n, q_n \in F$. L'insieme delle parole accettate da $A$ si dice *linguaggio accettato (o riconosciuto)* da $A$ e si denota con $L(A)$. A ogni automa non-deterministico $A = \langle Q, \Sigma, \delta, q_0, F \rangle$ associamo un grafo diretto con frecce etichettate:
- I vertici sono gli stati
- Le frecce sono le triple $(q,a,p)$ con $q \in Q$, $a \in \Sigma$, $p \in \delta(q,a)$
- Lo stato iniziale è denotato da una freccia entrante
- Gli stati finali sono identificati dal doppio bordo

Ci si convince facilmente a questo punto che:
- Una parola $w$ può essere etichetta di molti cammini uscenti dallo stato iniziale, o anche nessuno
- La parola $w$ è accettata $\iff$ c'è un cammino dallo stato iniziale a uno stato finale con etichetta $w$
- Questo non esclude però che vi possano essere altri cammini con etichetta $w$ che partono dallo stato iniziale e terminano in uno stato non finale

> [!example] 
> ![[Screenshot From 2026-03-06 11-03-22.png|300]]
> ![[Screenshot From 2026-03-06 11-03-35.png|300]]
  >![[Pasted image 20260306110549.png|300]]
  > 
> Questa computazione non è accettata perché non finisce in uno stato finale, ma questo non significa che la parola non sia accettata perché possono esserci altre computazioni (come la verde) che terminano in uno stato finale

Apparentemente questo tipo di progettazione non ci da nessun vantaggio, perché le computazioni da fare diventerebbero troppe. Supponiamo infatti di avere 2 cammini possibili per ogni scelta, e dovessimo leggere una parola di lunghezza 10, questo vorrebbe dire che dovremmo fare 210 computazioni. Proviamo però a ragionare come segue: se mi trovo in q0 e il mio input è $a$, posso andare o in q0 o in q1, ecc… Se mi trovo in q1 e l'input è $b$, andrò in q2,… Quindi ad ogni passo, in base alla colonna che avevo al passo precedente, e alla lettera di input, produco una colonna nuova. Questo procedimento è praticamente uguale a quello di un automa a stati deterministico, perché in base alla configurazione del mio dispositivo e alla lettera di input, determina la configurazione successiva. In questo caso la differenza è che abbiamo delle colonnine di stati.

> [!example] ![[Pasted image 20260306110859.png|300]]

**Teorema:** Sia $A$ un automa a stati finiti non deterministico, esiste effettivamente un equivalente automa a stati finiti deterministico $A'$ tale che $L (A) = L (A')$, cioè esiste l'automa deterministico e l'algoritmo per costruire questo automa a partire dall'automa non deterministico. Quindi, prima andremo a disegnare l'automa a stati finiti non deterministico (prima immagine con più percorsi in verde e rosso), cosa che risulta molto più facile rispetto al crearne uno deterministico, e poi sapremo che esiste l'algoritmo in grado di trasformare quel grafo nel grafo di un automa a stati finiti deterministico (seconda immagine con colonne). 

*Dimostrazione:* Gli stati dell'automa deterministico registreranno l'insieme gli stati raggiungibili dalle computazioni dell'automa non deterministico sul medesimo input $\rightarrow$ Accetto se fra questi ce n'è uno finale

**Costruzione:** Sia $A = \langle Q, \Sigma, \delta, q_0, F \rangle$ l'automa non deterministico. Definisco l'automa deterministico $A' = \langle Q', \Sigma, \delta', s_0, F' \rangle$ come segue: 
- L'insieme degli stati è l'insieme $Q' = \wp(Q)$ costituito dai sottoinsiemi di $Q$ (uno stato dell' $A'$ deterministico è in realtà un'insieme di stati dell' $A$ non deterministico)
- Lo stato iniziale è $s_0 = \{q_0\}$ (singleton costituito dallo stato iniziale dell' $A$ non deterministico)
- Gli stati finali sono tutti i sottoinsiemi di $Q$ che contengono almeno un elemento di $F$ (quindi almeno uno stato finale dell' $A$ di partenza), cioè: $F' = \{s \in \wp(Q) | s \cap F \neq \emptyset \}$ 
- La funzione di transizione $\delta'(s,a) = \bigcup\limits_{q \in s} \delta(q,a), s \in \wp(Q), a \in \Sigma$ (prende in input uno stato dell' $A$, cioè un insieme degli stati dell'automa non-deterministico, e una lettera $a$)

*Dimostrazione:* per ogni $w \in \Sigma^*$ risulta che $\hat{\delta'}(s_0,w) = \{q \in Q | \text{nel grafo di A c'è un cammino da}\ q_0 \ \text{a} \ q \ \text{con etichetta} \ w\}$ (dimostrazione per induzione sulla lunghezza di $w$). Pertanto sono equivalenti: 
- $w\ \text{è accettata da} \ A$
- $\hat{\delta'}(s_0,w) \ \text{contiene uno stato di} \ F$
- $\hat{\delta'}(s_0,w) \in F'$
- $w \ \text{è accettata da} \ A'$
$\implies L(A) = L(A')$ 

Come abbiamo accennato precedentemente, però, sarebbero presenti troppi stati da controllare (se $Q=n$, allora $Q'=2^n$). Per fortuna <u>non tutti questi stati sono necessari: possiamo definire uno stato necessario/accessibile come un cammino nel grafo dell'automa che, dallo stato iniziale, mi porta a q. Tutti gli altri stati si dicono inaccessibili.</u>

*Definizione:* sia $A = \langle Q, \Sigma, \delta, q_0, F \rangle$ un automa a stati finiti deterministico. Uno stato $q \in Q$ si dice *accessibile* se $q = \hat{\delta}(q_0,w)$ per qualche $w \in \Sigma^*$. Se io elimino tutti questi stati inaccessibili e restringo la funzione di transizione (delta cappello) solo agli stati accessibili, ottengo un automa con un numero di stati più piccolo, equivalente a quello di partenza, che accetta lo stesso linguaggio di quello di partenza.

> [!example] 
> ![[Screenshot From 2026-03-06 11-26-56.png|300]] 
> 
> Nell'esempio vediamo come lo stato q1 è stato eliminato perché, iniziando da q0 (stato di partenza), non c'è nessun cammino che mi porta in q1.

**Algoritmo di determinizzazione (per prevedere solo gli stati accessibili):**
1. Prima di tutto inserisco nella lista degli stati, lo stato iniziale dell'automa deterministico, che come abbiamo visto prima corrisponde all'insieme costituito dal solo stato iniziale dell'automa non deterministico di partenza: $s_0=\{q_0\}$ 
2. Considero ogni stato $r$ della lista (che per ora contiene il solo elemento $q0$, singleton) e ogni lettera $a \in \Sigma$ 
3. Calcolo il valore della funzione di transizione sulla coppia stato-lettera: $s= \delta'(r, a)$; Una volta calcolato ho due possibilità: o lo stato è già nella lista degli stati, oppure no.   
4. Se $s$ non è nella lista degli stati allora lo aggiungo alla lista degli stati. 
5. Se $s$ contiene uno stato finale di $A$, aggiungo $s$ alla lista degli stati finali. 

> [!example] 
> ![[Screenshot From 2026-03-06 11-32-20.png|300]]
> ![[Screenshot From 2026-03-06 11-32-32.png|300]]
> ![[Screenshot From 2026-03-06 11-33-08.png|300]]
> 
> Come possiamo notare, l'automa deterministico ha più stati rispetto a quelli che avevamo trovato noi manualmente nel primo esempio abac

> [!info] Caso difficile:
> ![[Pasted image 20260308174311.png|300]]
### Automa a stati finiti non deterministico con epsilon-transizioni
Un automa a stati finiti non deterministico con $\epsilon$-transizioni è una quintupla $A = \langle Q,\Sigma, \delta, q_0, F \rangle$, dove $Q,\Sigma, q_0, F$ sono come nell'automa deterministico, mentre la funzione $\delta: Q \times (\Sigma \cup \{\epsilon\}) \rightarrow \wp(Q)$ è la funzione di transizione, che ora non è definita solo sulle coppie *stato-lettera*, ma anche sulle coppie *stato-parolaVuota*, e restituisce un insieme di stati (potrebbe essere anche vuoto). <u>Questo modello mi permette di saltare da uno stato all'altro senza leggere il nastro di input</u> (anche per questo modello c'è un algoritmo che mi permette di passare da un automa di questo tipo ad un automa a stati finiti deterministico). Una parola $w \in \Sigma^*$ è *accettata* da $A$ se esistono $n\geq 0$, $a_1,...,a_n \in \Sigma \cup \{\epsilon\}$ e $q_1,..q_n \in Q$ tali che $w = a_1a_2...a_n$, $q_i \in \delta(q_{i-1},a_i)$, $1 \leq i \leq n$, $q_n \in F$ . In altre parole la parola è accettata dall'automa se esiste un cammino dallo stato iniziale a uno finale nel grafo con etichetta $w$, e se $w$ può essere scritta come una sequenza di lettere e parole vuote. L'insieme delle parole accettate da $A$ si dice *linguaggio accettato (o riconosciuto)* da $A$ e si denota con $L(A)$.

**Grafo dell'automa:** A ogni automa non deterministico con $\epsilon$-transizioni $A = \langle Q,\Sigma, \delta, q_0, F \rangle$ associamo un grafo diretto con frecce etichettate:
- I vertici sono gli stati
- Le frecce sono le triple $(q,a,p)$ con $q \in Q$, $a \in \Sigma \cup \{\epsilon\}$, $p \in \delta(q,a)$ 
- Lo stato iniziale è denotato da una freccia entrante
- Gli stati finali sono identificati dal bordo doppio

Ci si convince facilmente a questo punto che:
- Una parola $w$ può essere etichetta di molti cammini uscenti dallo stato iniziale, o anche nessuno
- La parola $w$ è accettata $\iff$ c'è un cammino dallo stato iniziale a uno stato finale con etichetta $w$
- Questo non esclude però che vi possano essere altri cammini con etichetta $w$ che partono dallo stato iniziale e terminano in uno stato non finale

> [!example] ![[Screenshot From 2026-03-12 10-16-28.png|300]]
> Nota che tra le colonne ora appare anche $\epsilon$, cosa che prima non era contemplata

**Teorema:** Sia $A$ un automa a stati finiti non deterministico con $\epsilon$-transizioni. Esiste effettivamente un automa a stati finiti deterministico $A'$ tale che $L(A) = L(A')$ (accetta lo stesso linguaggio)

*Dimostrazione:* Gli stati dell'automa deterministico registreranno l'insieme degli stati raggiungibili dalle computazioni dell'automa non deterministico sul medesimo input. Accetto la parola se fra questi ce ne è uno finale

**$\epsilon$-chiusura:** Se però con l'automa a stati finiti non deterministico le nostre computazioni si limitavano alla potenza-numero-di-lettere della parola, ora teoricamente posso avere anche infinite $\epsilon-$transizioni (la concatenazione di un miliardo di $\epsilon$, è sempre una parola vuota). Quindi bisogna utilizzare una tecnica per trattare queste $\epsilon-$transizioni: la $\epsilon-$chiusura di uno stato. 

*Definizione:* sia $A = \langle Q,\Sigma, \delta, q_0, F \rangle$ un automa a stati finiti non deterministico con $\epsilon$-transizioni, la $\epsilon$-chiusura di uno stato $q \in Q$ è il più piccolo sottoinsieme $C_{\epsilon}(q)$ di $Q$ tale che:
- $q \in C_{\epsilon}(q)$
- Per ogni $p \in C_{\epsilon}(q)$, $\delta(p,\epsilon) \subseteq C_{\epsilon}(q)$

La $\epsilon$-chiusura di un insieme di stati $S \subseteq Q$ è l'unione delle $\epsilon$-chiusure dei singoli stati di $S: C_{\epsilon}(S) = \bigcup\limits_{q \in S} C_{\epsilon}(q)$. In altre parole la $\epsilon-$chiusura di uno stato è l'insieme degli stati che da questo si possono raggiungere utilizzando solo $\epsilon-$transizioni (questo insieme deve contenere lo stato stesso).

> [!example] ![[Screenshot From 2026-03-12 10-31-58.png|300]]
> - $C_{\epsilon}(q_0)$ deve contenere $q_3$ perché con $\epsilon$-transizioni vado in $q_3$
> - $C_{\epsilon}(\{q_0,q_1\})$ è la $\epsilon$-chiusura della coppia di stati $q_0,q_1$ che si ottiene facendo l'unione delle chiusure dei due stati

**Calcolo della $\epsilon$-chiusura:**
- *Input*: un automa non deterministico con $\epsilon$-transizioni $A = \langle Q,\Sigma, \delta, q_0, F \rangle$ e uno stato $q \in Q$
- *Output*: la $\epsilon$ -chiusura di $q$

```
R <- (q) // inizialmente la lista contiene solo lo stato q
p <- puntatore al primo elemento di R
while p != NIL   
Prendiamo tutti gli elementi che si ottengono partendo dal primo elemento della lista e seguendo tutte le ε transizioni e le aggiungiamo alla lista, avendo l'accortezza di aggiungere solo quelle nuove, e non quelle che già appaiono in essa, quindi non creiamo ripetizioni.
	 appendi a R tutti gli elementi di delta (p, epsilon) \ R
	 p <- successivo(p) // passo all'elemento successivo
return R
```

> [!example] ![[Screenshot From 2026-03-12 10-35-50.png|300]]

**Costruzione:** Sia $A = \langle Q,\Sigma, \delta, q_0, F \rangle$ l'automa non deterministico con $\epsilon$-transizioni, definisco l'automa deterministico $A' = \langle Q',\Sigma, \delta', q_0, F' \rangle$ come segue:
- L'insieme degli stati è l'insieme $Q' = \wp(Q)$ costituito dai sottoinsiemi di $Q$ (cioè uno stato di $A'$  è un insieme di stati di $A$ )
- Lo stato iniziale è $s_0 = C_{\epsilon}(q_0)$ (stato iniziale è una $\epsilon$-chiusura, quindi non si tratta più di un singleton, ma di un insieme di stati)
- Gli stati finali sono tutti i sottoinsiemi di Q che contengono almeno un elemento di $F$, cioè: $F' = \{s \in \wp(Q) | s \cap F \neq \emptyset\}$ (sottoinsiemi che contengono almeno uno stato finale dell' $A$ precedente)
- La funzione di transizione è $\delta'(s,a): C_{\epsilon}(\bigcup\limits_{p \in s} \delta(p,a))$, $s \in \wp(Q)$, $a \in \Sigma$ (insieme di tutti quegli stati che posso raggiungere in $A$ partendo da uno qualunque degli stati di $s$ e poi seguendo la freccia con etichetta $a$, e poi tute le frecce che voglio con etichetta $\epsilon$)

*Dimostrazione:* Per ogni $w \in \Sigma^*$ risulta che $\hat{\delta'}(s_0,w) = \{q \in Q |$ nel grafo di $A$  c'è un cammino da $q_0$ a $q$ con etichetta $w$ $\}$. Pertanto sono equivalenti:
- $w$ è accettata da $A$
- $\hat{\delta'}(s_0,w)$ contiene uno stato di $F$
- $\hat{\delta'}(s_0,w) \in F'$
- $w$ è accettata da $A'$

Quindi $L(A) = L(A')$
Sostanzialmente  è sempre possibile costruire un automa deterministico che accetta il medesimo linguaggio, ed esiste inoltre un algoritmo che mi permette di trovare questo automa deterministico equivalente

**Algoritmo di determinazione:** 
- *Input*: un automa non deterministico con $\epsilon$ transizioni $A = \langle Q,\Sigma, \delta, q_0, F \rangle$
- *Output*: un automa deterministico equivalente $A' = \langle Q',\Sigma, \delta', q_0, F' \rangle$

<div style="display: flex; justify-content: center;">
  <img src="Screenshot From 2026-03-23 17-01-04.png" width="300">
</div>

> [!example] 
> - valuto le chiusure di ogni stato
> ![[Screenshot From 2026-03-12 11-25-10.png|300]]
> - valuto la funzione $\delta$ su questo automa
> ![[Screenshot From 2026-03-12 11-25-29.png|300]]
> 
> Il grafo dell'esempio sarà quindi:
> ![[Screenshot From 2026-03-12 11-29-04.png|300]]
> 
> Ricorda: ogni valutazione va fatta ricordando che "sono contemporaneamente in ognuno di questi stati". L'algoritmo che si può sfruttare prevede di: applicare uno degli input possibili per volta, raccogliere tutti i risultati conseguenti e infine allargare con $\epsilon$
### Espressioni regolari
La teoria dei linguaggi formali può essere vista come l'arte di rappresentare degli oggetti infiniti con una descrizione finita. Per costruire nuovi linguaggi partendo da linguaggi noti, possiamo applicare delle operazioni: 
1. unione
2. intersezione 
3. completamento
4. Concatenazione: $L_1L_2 = \{uv | u \in L_1, v \in L_2\}$. Esempio:
<div style="display: flex; justify-content: center;">
  <img src="Screenshot From 2026-03-14 16-48-13.png" width="400">
</div>
5. Potenza: $L^n = LL...L$ (n volte). concatenazione di *n* copie di un linguaggio con se stesso. Esempio:
<div style="display: flex; justify-content: center;">
  <img src="Screenshot From 2026-03-14 16-50-51.png" width="400">
</div>
6. Chiusura di Kleene: $L^* = \bigcup\limits_{n \geq 0}L^n = \{u_1u_2...u_n | n \geq 0, u_1,...,u_n \in L\}$. Prendo tutte le potenze di un linguaggio e ne faccio l'unione. A differenza delle prime 3 operazioni che se applicate a insiemi finiti mi ridanno insiemi finiti, la chiusura di Kleene mi dà un insieme infinito, a meno che l'insieme di partenza non sia o l'insieme vuoto, o l'insieme che contiene la sola parola vuota. Conterrà quindi tutte le parole che posso ottenere concatenando un numero arbitrario di parole di L. Esempio:
<div style="display: flex; justify-content: center;">
  <img src="Screenshot From 2026-03-14 16-51-56.png" width="400">
</div>

*Osservazione:* unione, concatenazione e chiusura di Kleene sono dette operazioni regolari

Quindi, i linguaggi finiti li sappiamo definire facendo una semplice lista delle parole; con questi linguaggi finiti posso poi fare le operazioni di unione e concatenazione ed ottenere nuovi linguaggi finiti, e con la chiusura di Kleene ottengo nuovi linguaggi infiniti. Con questi linguaggi infiniti posso fare altre operazioni di unione e concatenazione, ottenendo altri linguaggi sempre più complessi. Questo vuol dire che posso scrivere la definizione dei linguaggi infiniti semplicemente indicando le operazioni che eseguo per ottenerli: ovvero le espressioni regolari. <u>Le espressioni sono parole che costruisco partendo dalle lettere dell'alfabeto e dal simbolo ∅, seguendo le regole.</u> Rappresentano quindi una sequenza di operazioni fatte a partire dall'alfabeto iniziale.  

> [!info] 
> - $\emptyset =$ **linguaggio vuoto**, cioè non contiene nessuna parola
> - $\epsilon =$ **parola vuota**, cioè la parola di lunghezza 0
>> [!example] 
>> - $\emptyset^* = \{\epsilon\}$
>> - $\epsilon^* = \{\epsilon\}$

*Definizione:* sia $\hat\Sigma$ l'alfabeto ottenuto aggiungendo a $\Sigma$ le lettere $\emptyset, +, *, (,), . \implies$ si dicono *espressioni regolari* sull'alfabeto $\Sigma$ le parole sull'alfabeto $\hat\Sigma$ che si ottengono applicando un numero finito di volte le regole seguenti regole:
1. Ogni lettera $a \in \Sigma$ è un'espressione regolare, $\emptyset$ è un'espressione regolare
2. Se $E$ e $F$ sono espressioni regolari, allora $(E+F), (EF)$ e $E^*$ sono espressioni regolari.

A ogni espressione regolare è associato un linguaggio, detto *linguaggio denotato dall'espressione regolare* e definito dalle regole seguenti:
1. Per ogni $a \in \Sigma$, l'espressione regolare $a$ denota il linguaggio $\{a\}$; l'espressione regolare $\emptyset$ denota il linguaggio vuoto
2. Detti $L_E$ e $L_F$ i linguaggi denotati dalle espressioni regolari $E$ e $F$, i linguaggi denotati dalle espressioni regolari $E+F, (EF), E^*$ sono, rispettivamente, $L_E \cup L_F, L_EL_F, L_{E^*}$

*Definizione:* i linguaggi denotati da espressioni regolari si dicono *linguaggi regolari* (i linguaggi di tipo 0 sono detti anche linguaggi regolari)

*Osservazione:* la classe dei linguaggi regolari è la più piccola famiglia di linguaggi che:
- Contiene i linguaggi finiti
- è chiusa per le operazioni regolari

Perché se due linguaggi sono denotati da espressioni regolari, la loro unione sarà denotata anche lei da un'espressione regolare. Lo stesso per la concatenazione. La chiusura di Kleene sarà denotata dall'espressione regolare del linguaggio con la * vicino

*Osservazione:* possiamo omettere le parentesi, rispettando la priorità:
- Chiusura di Kleene
- Concatenazione
- Somma

> [!example] 
> - $a+b$: espressione regolare che denota il linguaggio composto dalle due parole a e b (**unione**)
> - $(a+b)^*$: **chiusura di Kleene**, tutte le concatenazioni di parole che appartengono all'insieme a, b. Cioè tutte le concatenazioni di lettere a, b . In sostanza è l'insieme delle parole sull'alfabeto (a, b) 
> - $\emptyset$: se concateno 0 parole dell'insieme vuoto (linguaggio che contiene solo la parola vuota)
> - $((a+b)(a+b)(a+b))$: ottengo tutte le parole di lunghezza 3, e poi ne faccio la chiusura di Kleene e ottengo  tutte le parole che hanno lunghezza multipla di 3
> - $(a+ab)^*$: ogni blocco può essere $a$ oppure $ab$ (parole ottenute concatenando blocchi di "$a$" o "$ab$")
> - $(a+b)^*abb$: otteniamo tutte le lettere che finiscono con $abb$ (**concatenazione**)
> - $(b^*ab^*ab^*)^*b^*ab^*$: dentro alla parentesi ho: qualche $b$ (anche 0), poi qualche $a$, poi qualche $b$, poi un altro $a$, poi un altro $b \implies$ 2 $a$ (con qualsiasi numero di $b$ intorno e in mezzo). Poi si ripete questo blocco quante volte volte si vuole (qui ho parole con un numero pari di $a$). Lo step finale implica di aggiungere un solo $a$ 
> - $(a+b+c)^*abac (a+b+c)^*$: tutte le parole su $\{a,b,c\}$, poi si ha $abac$ e infine ancora $(a+b+c)^*$
> - $((a+b+c),)^*(a+b+c) \& (a+b+c)+a+b+c$

**Teorema di Kleene:** un linguaggio è regolare $\iff$ è riconosciuto da un automa a stati finiti. 

Esiste una procedura effettiva che:
- Data un'espressione regolare, produce un automa a stati finiti che accetta il linguaggio denotato da tale espressione (*sintesi*)
- Dato un'automa a stati finiti, produce un'espressione regolare che denota il linguaggio accettato da tale  automa (*analisi*)

Utilizzeremo esclusivamente automi non deterministici con $\epsilon-$transizioni con un unico stato finale (non restrittivo)

Un'espressione regolare, dove $E$ e $F$ sono espressioni regolari, può essere:

 <div style="display: flex; align-items: flex-start; gap: 20px;">
   <div style="flex: 1;">
     <img src="Screenshot From 2026-03-13 11-05-00.png" style="width: 100%; border-radius: 8px;">
   </div>
  <div style="flex: 1.5;">
<li>una lettera, un insieme vuoto o una parola vuota (che corrisponde alla chiusura di Kleene dell'insieme vuoto)</li>
<li>(E+F), cioè l'unione</li>
<li>(EF), cioè la concatenazione</li>
<li> E^*, cioè la chiusura di Kleene</li>
   </div>
 </div>
 
 Quindi il caso base prevede che se $G = a$ allora si costruisce un automa che accetta solo $a$, se $G = \emptyset$ allora si costruisce l'automa che non accetta nulla
 
 <div style="display: flex; align-items: flex-start; gap: 20px;">
   <div style="flex: 1;">
     <img src="Screenshot From 2026-03-13 11-05-38.png" style="width: 100%; border-radius: 8px;">
   </div>
  <div style="flex: 1.5;">
<b>Unione:</b> avendo già costruito gli automi di E e F (gli ovali nella figura), si aggiunge un nuovo stato iniziale e un nuovo stato finale, e due ε transizioni che mi portano dallo stato iniziale nuovo agli stati iniziali vecchi, e due ε transizioni che mi portano dagli stati finali vecchi allo stato finale nuovo. Questo nuovo automa accetta l'unione dei linguaggi, accettati dai due automi E e F. Quindi tutti i cammini riusciti dallo stato iniziale a quello finale, avranno come etichetta (l'ovale) una parola accettata dall'automa E o F.
   </div>
 </div>

Siano $A_1 = \langle Q_1,\Sigma,\delta_1, i_1, \{f_1\}\rangle$ e $A_2 = \langle Q_2,\Sigma,\delta_2, i_2, \{f_2\}\rangle$ gli automi che riconoscono i linguaggi denotati dalle espressioni regolari $E$ e $F$. 
Supponiamo $Q_1 \cap Q_2 = \emptyset$. 
Costruiamo $A = \langle Q,\Sigma,\delta, i, \{f\}\rangle$ come segue:
- $Q = Q_1 \cup Q_2 \cup \{i,f\}$ ove $i,f$ sono due nuovi stati
- $i,f$ sono rispettivamente lo stato iniziale e l'unico stato finale
- Le frecce del grafo di $A$ sono le frecce dei grafi di $A_1$ e $A_2$ e, inoltre, $(i,\epsilon,i_1), (i,\epsilon,i_2), (f_1,\epsilon, f), (f_2,\epsilon, f)$

Allora $A$ accetta il linguaggio denotato da $(E+F)$

> [!example] ![[Screenshot From 2026-03-13 11-13-32.png|300]]

 <div style="display: flex; align-items: flex-start; gap: 20px;">
   <div style="flex: 1;">
     <img src="Screenshot From 2026-03-13 11-14-59.png" style="width: 100%; border-radius: 8px;">
   </div>
  <div style="flex: 1.5;">
<b>Concatenazione:</b> si connette con una epsilon-transizione lo stato finale del primo automa con lo stato iniziale del secondo automa. Lo stato iniziale del nuovo automa sarà quello del primo automa, mentre il suo stato finale sarà quello del secondo automa.
   </div>
 </div>

Siano $A_1 = \langle Q_1,\Sigma,\delta_1, i_1, \{f_1\}\rangle$ e $A_2 = \langle Q_2,\Sigma,\delta_2, i_2, \{f_2\}\rangle$ gli automi che riconoscono i linguaggi denotati dalle espressioni regolari $E$ e $F$. 
Supponiamo $Q_1 \cap Q_2 = \emptyset$. 
Costruiamo $A = \langle Q,\Sigma,\delta, i, \{f\}\rangle$ come segue:
- $Q = Q_1 \cup Q_2$
- Lo stato iniziale è quello di $A_1$ e quello finale è quello di $A_2$
- Le frecce del grafo di $A$ sono le frecce dei grafi di $A_1$ e $A_2$ con l'aggiunta della freccia $f_1, \epsilon, i_2$

Allora $A$ accetta il linguaggio denotato da $(EF)$
> [!example] 
> ![[Screenshot From 2026-03-13 11-17-52.png|300]]

 <div style="display: flex; align-items: flex-start; gap: 20px;">
   <div style="flex: 1;">
     <img src="Screenshot From 2026-03-13 11-18-53.png" style="width: 100%; border-radius: 8px;">
   </div>
  <div style="flex: 1.5;">
<b>Chiusura di Kleene:</b> abbiamo l'automa E che riconosce una certa espressione, e poi il nuovo automa che riconosce la chiusura di Kleene del linguaggio accettato da questo, ovvero che accetti tutte le parole che si ottengono concatenando 0 o più parole accettate dall'automa di partenza. Aggiungo una ε transizione che dallo stato finale mi riporta a quello iniziale: così da accettare tutte le parole che si ottengono concatenando una o più parole dell'automa di partenza. Mi manca però il riconoscimento della parola vuota, quindi aggiungo un nuovo stato iniziale e un nuovo stato finale e poi una \epsilon-transizione dallo stato iniziale nuovo a quello vecchio, una \epsilon-transizione dallo stato finale vecchio a quello nuovo, e una \epsilon-transizione dallo stato iniziale nuovo allo stato finale n\uovo.
   </div>
 </div>

Siano $A_1 = \langle Q_1,\Sigma,\delta_1, i_1, \{f_1\}\rangle$ l'automa che riconosce il linguaggio denotato dall'espressione regolare $E$. 
Costruiamo $A = \langle Q,\Sigma,\delta, i, \{f\}\rangle$ come segue:
- $Q = Q_1 \cup \{i,f\}$ ove $i,f$ sono due nuovi stati
- $i,f$ sono rispettivamente lo stato iniziale e l'unico stato finale
- Le frecce del grafo di $A$ sono le frecce del grafo di $A_1$ e, inoltre, $(f_1,\epsilon,i_1), (i,\epsilon,i_1), (f_1,\epsilon, f), (i,\epsilon, f)$

Allora $A$ accetta il linguaggio denotato da $E^*$
> [!example]
> ![[Pasted image 20260313112221.png|300]]

**L'algoritmo di sintesi**
*Proposizione:* data un'espressione regolare $G$, si può effettivamente costruire un'automa a stati finiti che riconosce il linguaggio denotato da $G$:
1. se $G$ è un'espressione regolare di base (lettera o insieme vuoto), allora restituisco l'automa corrispondente; 
2. se $G = (E + F)$, allora calcolo gli automi $A_1$ e $A_2$ corrispondenti a $E$ e $F$; costruisco l'automa dell'espressione $(E + F)$; 
3. se $G = (EF)$, allora calcolo gli automi $A_1$ e $A_2$ corrispondenti a $E$ e $F$; costruisco l'automa dell'espressione $(EF)$; 
4. se $G = E^*$ , allora calcolo l'automa $A_1$ corrispondente a $E$; costruisco l'automa dell'espressione $E^*$;

Esempio completo nelle slide del prof

**Dall'automa a stati finiti all'espressione regolare:** sia $A$ un automa a stati finiti e per semplicità supporremo che:
- Gli stati sono $q_1, q_2,…q_n$
- Lo stato iniziale è $q_1$

> [!tip] Avendo un automa A, vogliamo che: 
> abbia un unico stato finale e, per semplicità, non abbia frecce che entrano nello stato iniziale né frecce che escono dallo stato finale. Questo implica che un percorso uscito dallo stato iniziale non potrà più rientrarci (perché non ci sono frecce entranti in esso) così come non potrà passare per lo stato finale se non all'ultimo (perché non ci sono frecce uscenti da esso). Aggiungo quindi delle $\epsilon$-transizioni che mi portano dal nuovo stato iniziale al vecchio stato iniziale, e dai vecchi stati finali al nuovo stato finale. Questo nuovo automa accetterà le stesse parole dell'automa vecchio.
> Inoltre, supponiamo anche di dare dei nomi ai vari stati del nostro automa: $q_1$, $q_2$, …, $q_n$, dove lo stato
> iniziale corrisponde al $q_{n-1}$ (quindi al penultimo stato) e lo stato finale a $q_n$. Prendiamo dei linguaggi che in qualche modo approssimano il linguaggio che vogliamo cercare e cerchiamo delle espressioni regolari per questi linguaggi.
> 
>![[Pasted image 20260319171718.png|300]]

Per ogni $i, j, k$ con $0 \leq k < i, j \leq n$ (indici i e j compresi tra 1 e n, mentre k compreso tra 0 e n) consideriamo l'insieme $L_{ijk}$ contenente le seguenti parole: *le etichette dei cammini nel grafo di $A$ che iniziano in $q_i$, terminano in $q_j$, e attraversano solo stati di indice $\leq k$* (che attraversi soltanto, non che parta o arrivi in uno di questi stati)

*Osservazione:* $L(A)= \bigcup\limits_{q_j \in F} L_{1jn}$. Quindi, se riesco a trovare espressioni regolari per tutti gli $L_{ijk}$, avrò anche un'espressione regolare di $L(A)$

> [!example] ![[Pasted image 20260319110515.png|300]]
> - $q_2$ stato iniziale
> - $q_5$ stato finale
> - quindi devo prendere tutti i cammini che partono da $q_2$, che arrivano in $q_4$ e che attraversano soltanto $q_1 \implies L(A) = L_{n-1 \ n \ n-2}$ dato dal cammino $q_2 \xrightarrow[]{b} q_1 \xrightarrow[]{c} q_5$

Se riesco a trovare un'espressione regolare per tutti gli $L_{ijk}$, avrò un'espressione regolare per $L_{n-1 \ n \ n-2}$, che
è proprio il linguaggio accettato dall'automa. Ho quindi ottenuto l'espressione regolare corrispondente al
linguaggio. Quello che dobbiamo mostrare è un meccanismo che ci permetta di determinare delle espressioni regolari che denotano questi linguaggi $L_{ijk}$ Per ottenere tutti questi $L_{ijk}$, si procede per induzione su $k$.

Il linguaggio $L_{ij0}$ è denotato dall'espressione regolare $E_{ij0}$ ottenuta come segue (supponendo $k = 0$):
- La somma delle etichette delle frecce da $q_i$ a $q_j$, se ne esistono e $i\neq j$ 
- La somma delle etichette delle frecce da $q_i$ a $q_j$ più $\epsilon$ se ne esistono e $i = j$
- $\epsilon$ se non ci sono frecce da $q_i$ a $q_j$ e $i = j$
- $\emptyset$ se non ci sono frecce da $q_i$ a $q_j$ e $i \neq j$

> [!example] ![[Screenshot From 2026-03-19 11-07-41.png|300]]
> Dato che $k \leq 0$ ciò significa che i cammini devono attraversare 0 stati (cioè si tratta di frecce dirette)

Supponiamo $k\geq 1$. I cammini nel grafo $A$ con origine in $q_i$, termine in $q_j$ e che attraversano esclusivamente stati di indice $\leq k$, si possono ripartire in due categorie:
1. Cammini da $q_i$ a $q_j$ che attraversano solo stati di indice $\leq k-1$. Questi cammini per $q_k$ potrebbero non passare proprio, oppure possono passare per $q_1,q_2,...,q_{k-1}$. Per questi l'espressione regolare già l'abbiamo
<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319173347.png" width="300">
</div>
2. I cammini costituiti dalla concatenazioni di (cioè i cammini che passano una o più volte per $k$):
	- Un segmento da $q_i$ a $q_k$ (segmento dallo stato iniziale fino a $q_k$)
	- Zero o più cammini con origine e termine in $q_k$ (dei cammini che da $q_k$ arrivano al secondo passaggio con $q_k$, e così via, fino all'ultimo passaggio con $q_k$)
	- Un segmento finale con origine da $q_k$ a $q_j$ (un ultimo segmento che da $q_k$ arriva allo stato finale)

Tutti che attraversano solo stati di indice $\leq k-1$

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319174122.png" width="300">
</div>

**Passo induttivo:** ora dobbiamo trovare l'espressione $E_{ijk}$, cioè l'espressione per il linguaggio costituito dalle etichette di tutti i cammini delle forme precederti. Innanzitutto facciamo l'unione dei cammini del primo punto con i cammini del secondo punto. I cammini del secondo gruppo saranno a loro volta riportati come la concatenazione del primo segmento con i segmenti intermedi (concatenazione di 0 o più coppie di un espressione utilizzando la chiusura di Kleene) e l'ultimo segmento. In sostanza le espressioni regolari $E_{ijk}$ che denotano i linguaggi $L_{ijk}$ devono soddisfare le relazioni ricorsive $E_{ijk} = E_{ijk−1} + (E_{ikk−1})(Ek k k−1)(E_{kjk−1}), \ 1 \leq k < i, j \leq n$ (la prima parte contiene tutti i cammini che NON usano $q_k$, mentre la seconda parte contiene tutti i cammini che usano $q_k$ ALMENO una volta). Un'espressione regolare di $L(A)$ è data dalla somma delle $E_{1jn}$ con $q_j$ stato finale. In questo mode è possibile calcolare tutte le espressioni regolari, compresa $E_{n-1 \ n \ n-2}$ che denota il linguaggio accettato dall'automa (esempio nelle slide 8-11)

**Regola pratica:** Per fare tutto ciò utilizzeremo il <u>Metodo di eliminazione degli stati</u> partendo dal nostro automa. Questo corrisponde ad un grafo sulle cui etichette delle frecce ci sono espressioni particolari che possono essere o lettere, o $\epsilon$ o somma di lettere.
1. Ci riduciamo al caso in cui:
	- $A$ ha un unico stato finale 
	- Nel grafo di $A$ non ci sono frecce che entrano nello stato iniziale, ne frecce che escono dallo stato finale
2. Sostituiamo le etichette multiple con la loro somma (cioè se ci sono più frecce tra due stati, uniscile con +)
3. Selezioniamo uno stato $q_k$ che non sia né iniziale né finale; per ogni freccia $q_i \xrightarrow[]{E_{ik}} q_k$ che entra nello stato $q_k$ e ogni freccia $q_k \xrightarrow[]{E_{kj}} q_j$ che esce da tale stato $(i, j \neq k)$ 
	- Rimpiazziamo l'etichetta della freccia $q_i \xrightarrow[]{E_{ij}} q_j$ con $E_{ij} + E_{ik}E_{kk}^*E_{kj}$ (o $E_{ij} + E_{ik}E_{kj}$ qualora manchi la freccia da $q_k$ a $q_k$); 
	- se tale freccia è assente, la creiamo (con etichetta $E_{ik}E_{kk}^*E_{kj}$ , ovvero $E_{ik}E_{kj}$ ); 
	- rimuoviamo lo stato $q_k$ e tutte le frecce che entrano o escono da esso; 
4. Ripetiamo il passo precedente fino a ottenere un automa con due stati e una sola freccia; 
5. L'etichetta di tale freccia è l'espressione regolare cercata.

> [!info] Intuizione finale:
> - un'automa riconosce parole tramite cammini
> - ogni cammino ha un'etichetta (una parola)
> - eliminando gli stati, stai "riscrivendo" tutti i cammini possibili in una sola formula regolare 
> - alla fine tutta la complessità del grafo si concentra in un'espressione regole

> [!example] > ![[Pasted image 20260319181524.png]]
### Operazioni booleane
*Proposizione:* la classe dei linguaggi regolari è un'algebra Booleana (altra conseguenza del teorema di Kleene). Ciò significa che i linguaggi regolari si comportano come insiemi normali e puoi fare operazioni logiche come OR, AND e NOT (e il risultato è sempre un linguaggio regolare)

*Dimostrazione:* Ci sono delle operazioni che sono naturali sulle espressioni regolari e altre operazioni che sono naturali sugli automi a stati finiti. Come abbiamo detto in precedenza, infatti, la famiglia dei linguaggi regolari è chiusa per le operazioni di unione, concatenazione e chiusura di Kleene. È anche chiusa per complemento ed intersezione, quindi per tutte le operazioni booleane.
- *Unione:* L'unione dei linguaggi denotati dalle espressioni $E$ e $F$ è denotato dall'espressione $E$ + $F$. 
- *Complemento:* Per il Teorema di Kleene, un linguaggio regolare è accettato da un automa a stati finiti deterministico. Il complemento è accettato dall'automa ottenuto scambiando stati finali e non. Quindi, per il Teorema di Kleene, è regolare 
- *Intersezione:* L'intersezione si esprime per mezzo di unione e complemento: $L \cap M = \Sigma^* - ((\Sigma^* -L) \cup (\Sigma^* - M))$

*Osservazione:* anche la differenza insiemistica si esprime per mezzo di unione e completamento: $L - M = \Sigma^* - ((\Sigma^* -L) \cup M)$ 

*Osservazione:* Date due espressioni regolari, posso effettivamente trovare espressioni regolari per l'unione, l'intersezione e i complementi dei linguaggi da esse denotati (devo però passare obbligatoriamente per gli automi).

**Equivalenza di espressioni regolari:** Date due espressioni regolari ho il problema di decidere se queste sono equivalenti o meno

> [!example] $a^* + (a^*ba^*)^* = (a+b)^*$ descrivono lo stesso linguaggio

*Osservazione:* dati due insiemi $X,Y$ sono uguali se e solo se $(X-Y) \cup (Y-X) = \emptyset$ 

La soluzione al mio problema è: date le espressioni regolari $E,F$, voglio sapere se $L_E = L_F$
1. Calcolare automi deterministici che accettano i linguaggi $L_E$ e $L_F$ denotati da $E$ e $F$, rispettivamente; 
2. Calcolare un automa deterministico che accetta $(L_E − L_F ) \cup (L_F − L_E )$; 
3. Verificare se tutti gli stati finali sono inaccessibili. In tal caso NON esiste nessuna parola accettata e quindi il linguaggio è vuoto

> [!tip] 
> - se **non esiste parola che li distingue** allora sono uguali
> - se **esiste almeno una parola diversa** allora non sono uguali

> [!info] Sintesi: per confrontare due espressioni regolari:
> - costruisci gli automi
> - confronti i linguaggi
> - riduci il problema a:"esiste una parola che li distingue?"
### Automa minimo
Due parole sono uguali per un linguaggio se aggiungendo lo stesso suffisso a destra si comportano sempre nello stesso modo rispetto all'appartenenza a $L$. Da questa idea nasce l'automa minimo. Si tratta di un'automa deterministico con il minor numero possibile di stati. <u>L'automa di Nerode di un linguaggio regolare L è un automa deterministico che accetta L col minimo numero di stati possibile</u>. Quindi definisco prima la congruenza di Nerode, così da definire un'equivalenza cioè una partizione dell'insieme $\Sigma^*$ in un numero finito di classi di equivalenza. Costruisco in seguito un automa minimo partendo da queste classi di equivalenza, che accetta il linguaggio $L$. 

**Equivalenza:** insieme sul quale è definita una relazione **riflessiva** ($x \sim x$), **simmetrica** ($x \sim y \implies y \sim x$) e **transitiva** ($x \sim y, y \sim z \implies x \sim z$). Ogni equivalenza divide l'insieme in **classi di equivalenza** (ogni stato sostituisce un singolo stato nell'automa finale), e queste classi formano una **partizione**. Quindi un equivalenza corrisponde ad una partizione, ogni elemento appartiene ad una sola classe e classi diverse non si intersecano. In altre parole 2 stati sono equivalenti se e solo se, per qualsiasi stringa in ingresso, portano l'automa rispettivamente a stati entrambi accettanti o entrambi non accettanti.

*Osservazione:* ogni equivalenza su un insieme $S$ determina una partizione dell'insieme $S$ in classi di equivalenza. Viceversa, ogni partizione dell'insieme $S$ è associata a un'equivalenza di cui le parti sono classi di equivalenza

*Definizione:* sia $\Sigma$ un alfabeto. Una relazione di equivalenza $\sim$ su $\Sigma^*$ si dice una congruenza destra (risp. Sinistra) se per ogni $u, v \in \Sigma^*$ tali che $u \sim v$ e per ogni lettera $a \in \Sigma$, si ha $ua \sim va$ (risp. $au \ \sim \ av$). <u>Una relazione che sia contemporaneamente una congruenza destra e una congruenza sinistra si dice congruenza.</u> In altre parole qui non abbiamo solo un insieme, ma un insieme che ha un'operazione definita su di esso. Ogni volta che si va a concatenare una lettera a destra di due parole equivalenti ottengo due parole equivalenti (rispettivamente per la congruenza sinistra).

> [!info] Si considerino le relazione $\sim$ su $\Sigma^*$ definite da:
> 1. $u \sim v$ se $u$ e $v$ iniziano con la stessa lettera (o $u$ = $v$ = $\epsilon$)
> 2. $u \sim v$ se $u$ e $v$ terminano con la stessa lettera (o $u$ = $v$ = $\epsilon$)
> 3. $u \sim v$ se $u$ e $v$ hanno la stessa lunghezza
> 4. $u \sim v$ se le prime due lettere di $u$ sono una permutazione delle prime due lettere di $v$ (o $u = v$)
> 
> Per es. $abaa \sim bab, aababb \sim aaabaa \ \text{ma} \ aababb \not\sim baaabaa$

**Equivalenza di Nerode:**
*Definizione:* sia $L$ un linguaggio sull'alfabeto $\Sigma$. L'equivalenza di Nerode associata al linguaggio $L$ è la relazione $N_L$ definita in $\Sigma^*$ da: $u \ N_L \ v$ se per ogni $y \in \Sigma^*$*, $uy \in L \iff vy \in L$. Cioè, $u \ N_L \ v$ se $u$ e $v$ sono equivalenti e se hanno gli stessi completamenti a destra per entrare in $L$. 

Quindi diciamo che due parole sono equivalenti se le parole che mi permettono di completare la prima ($u$, a destra) dentro una parola di $L$, sono esattamente le stesse che mi permettono di completare $v$ in una parola di $L$. Questo teorema ci da una caratterizzazione vera e propria dei linguaggi regolari in forma algebrica: <u>un linguaggio è regolare se lo posso ottenere come unione di classi di una congruenza destra di indice finito.</u>

$N_L$ è la relazione che identifica precisamente gli stati dell'automa minimo:
- Due parole stanno nella stessa classe se portano allo stesso "stato logico" rispetto a $L$
- Ogni classe diventa uno stato dell'automa di Nerode

**Lemma:** sia $L \subseteq \Sigma^*$ un linguaggio
1. La relazione $N_L$ è una congruenza destra. Quindi se $uN_Lv$ allora per ogni lettera $a$: $uaN_Lva$

*Dimostrazione:* $N_L$ è un'equivalenza. Dobbiamo verificare che per ogni $u, v \in \Sigma^*$ , $a \in \Sigma$, se $u N_L v$, allora $ua \ N_L \ va$. Invero sia $u \ N_L \ v$ e $a \in \Sigma$. Dalla definizione di $N_L$ si ha $uy \in L \iff vy \in L$ per ogni $y \in \Sigma^*$ Sostituendo $y$ con $ay$, ottengo $uay \in L \iff vay \in L$ Quindi si ha anche $ua \ N_L \ va$. Cioè ogni parola è in un'equivalenza $N_L$ con sé stessa (prop. Riflessiva) perché i suoi completamenti in $L$ sono gli stessi suoi completamenti in $L$. Se $u$ e $v$ hanno gli stessi completamenti, allora $v$ e $u$ avranno gli 25 stessi completamenti (prop. Simmetrica). Se $u$ e $v$ hanno gli stessi completamenti, e $v$ e $w$ hanno gli stessi completamenti, allora $u$ e $w$ hanno gli stessi completamenti (prop. Transitiva).

2. $L$ è unione di classi di equivalenza di $N_L$. Se $uN_Lv$ allora prendendo $y = \epsilon$: $u \in L \iff v \in L$. Quindi una classe di equivalenza non può "tagliare" il linguaggio a metà:
	- O sta tutto dentro $L$
	- Oppure sta tutta fuori

*Dimostrazione:* Sia $u \ N_L \ v$. Allora $uy \in L \iff vy \in L$ per ogni $y \in \Sigma^*$ . Prendendo $y = \epsilon$ ottengo $u \in L \iff v \in L$. Quindi ogni classe di equivalenza o sta tutta dentro L o sta tutta nel complemento. Cioè il linguaggio $L$ si può decomporre in classi di equivalenza $N_L$. In altre parole, $N_L$ è un'equivalenza quindi ci determina quella partizione in classi di equivalenza dell'insieme $\Sigma^*$. Ci sono queste famiglie di parole che sono tutte equivalenti tra loro, ogni parola sta in una classe di equivalenza e classi equivalenza diverse non hanno elementi comuni. Una classe di equivalenza o è contenuta dentro $L$ (quindi è l'unione di alcune classi di equivalenza), oppure è interamente fuori da $L$. Quello che non può succedere è che una classe di equivalenza sia a cavallo: qualche elemento in $L$ e qualche elemento fuori dal linguaggio $L$.

**Teorema di Myhill-Nerode (1958):** Sia $L$ un linguaggio sull'alfabeto $\Sigma$, le seguenti proposizioni sono equivalenti:
1. $L$ è regolare
2. $L$ è unione di classi di una congruenza destra su $\Sigma^*$ si dice finito 
3. $N_L$ ha indice finito (indice = numero di classi di equivalenza dell'equivalenza stessa)

Il teorema dice che il numero di “situazioni diverse” che bisogna distinguere per riconoscerlo è finito. Queste “situazioni diverse” sono proprio le classi di Nerode. Quindi:
- se ci sono finitissime classi, posso costruire un automa con finitissimi stati;
- se le classi sono infinite, non esiste un DFA che lo riconosce.

Posso vedere queste classi come gli stati di un automa, in cui la transizione con la lettera 'a' mi porta da una classe che contiene le parole $u,v$ a una classe che contiene le parole $ua,va$

<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Pasted image 20260320102535.png" style="width: 100%; border-radius: 8px;">
  </div>
  <div style="flex: 1.5;">
  L'ovale  rappresenta l'insieme di tutte le parole sull'alfabeto Sigma. Viene partizionato iN_2 parti:
  <li>Linguaggio L, cioè le parole che stanno nel linguaggio L</li>
  <li>Le parole che stanno fuori da linguaggio L e che non appartengono ad esso</li>
  
  Sopra a questo ovale diviso in un numero finito di pezzi, ci costruiremo sopra l'automa del linguaggio L
  </div>
</div>

Dimostrazione circolare per dimostrare che le 3 condizioni suddette siano equivalenti: 1. $\implies$ 2. $\implies$ 3. $\implies$ 1.

**Dimostrazione 1. $\implies$ 2.:** 
- Sia $L$ regolare
- $L$ è accettata da un automa deterministico $A = \langle Q; \Sigma, \delta, q_0, F \rangle$ con tutti gli stati accessibili
- Per ogni stato $q \in Q$ sia $L_q = \{w \in \Sigma^* | \hat\delta(q_0,w) = q\}$
- Gli insiemi $L_q$ costituiscono una partizione di $\Sigma^*$
- $L$ è l'unione delle classi $L_q$ con $q \in F$
- Quindi basta verificare che l'equivalenza $\sim$ associata alla partizione è una congruenza destra:
	- Supponiamo $u \sim v, u, v \in \Sigma^*$
	- Allora $u, v \in L_q$ per qualche stato $q$
	- Quindi $\hat\delta(q_0,u) = \hat\delta(q_0,v) = q$
	- Pertanto per ogni lettera $a, \hat\delta(q_0,va) = \delta(q,a) = q'$
	- Cioè $ua, va \in L_q$, per cui $ua \sim va$

**Dimostrazione 2. $\implies$ 3.:** $L$ unione di classi di una congruenza destra di indice finito $\implies N_L$ ha indice finito
- Sia $L$ unione di classi di una congruenza destra $\sim$ su $\Sigma^*$ di indice finito
- Basta mostra che ogni classe di $N_L$ è unione di classi di $\sim$
- Cioè che se $u \sim v$ allora $u \ N_L \ v,u,v \in \Sigma^*$ 
	- Supponiamo $u \sim v,u,v \in \Sigma^*$
	- Dato che $\sim$ è una congruenza destra, si ha $ua \sim va$ per ogni $a \in \Sigma$
	- Più in generale si ha $uy \sim vy$ per ogni $y \in \Sigma^*$
	- Visto che $L$ è unione di classi di $\sim$ ne segue che $uy \in L \iff vy \in L$ e quindi $u \ N_L \ v$
- Abbiamo mostrato che ogni classe di $N_L$ è unione di classi di $\sim$ve quindi contiene almeno una classe di $\sim$
- Ne concludiamo che l'indice di $N_L$ è minore o uguale a quello di $\sim$ e pertanto è finito

**Dimostrazione 3. $\implies$ 1.:** $N_L$ ha indice finito $\implies L$ è regolare
- Per il teorema di Kleene basta costruire un automa che accetta $L$
- Sia $A$ l'automa ottenuto nel modo seguente (automa di Nerode di $L$)
	- Gli stati sono le classi dell'equivalenza $N_L$
	- La funzione di transizione $\delta$ è definita nel modo seguente: siano $q \in Q$ e $a \in \Sigma$. Poiché $q$ è una classe dell'equivalenza $N_L$ tutte le parole $ua$ con $u \in q$ appartengono a un'unica classe $p$ dell'equivalenza $N_L$. Poniamo dunque $\delta(q,a) = p$
	- Lo stato iniziale $q_0$ è la classe d'equivalenza della parola vuota
	- Gli stati finali sono le classi d'equivalenza incluse in $L$
- Per ogni $w \in \Sigma, \hat\delta(q_0,w)$ è la classe di equivalenza della parola $w$
- Poiché $L$ è unione di classi dell'equivalenza $N_L$ tale stato è finale se e solo se $w \in L$
- Se ne conclude che il linguaggio accettato da $A$ è proprio $L$

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260320112259.png" width="300">
</div>

*Proposizione:* Sia $L$ un linguaggio regolare. L'automa di Nerode di $L$ è un automa deterministico che accetta $L$ col minimo numero di stati possibile. 

Il punto finale è il seguente:
- ogni stato di un DFA corrisponde, a una classe di parole indistinguibili;
- Nerode identifica proprio queste classi;
- quindi il numero minimo di stati possibili è proprio il numero delle classi di $N_L$​

Per questo l'automa di Nerode è il **minimo**.

*Dimostrazione:* 
- Sia $A$ l'automa deterministico che accetta $L$ col minimo numero di stati possibile e sia n tale numero; 
- nella dimostrazione 1. $\implies$ 2. si è visto che $L$ è unione di classi di una congruenza destra $\sim$ di indice n; 
- nella dimostrazione 2. $\implies$ 3. si è visto che l'indice di $N_L$ è minore o uguale di quello di $\sim$; 
- quindi l'indice di $N_L$ è minore o uguale a n; 
- nella dimostrazione 3. $\implies$ 1. si è costruito un automa che accetta $L$ con un numero di stati pari all'indice di $N_L$ ; 
- quindi l'indice di $N_L$ è maggiore o uguale a n; 
- ne segue l'asserto

> [!info] Sintesi:
> - **Congruenza destra**: se $u \sim v$, allora $ua \sim va$
> - **Nerode**: $uN_Lv$ se $u$ e $v$ hanno gli stessi completamenti a destra per entrare in $L$
> - **Myhill-Nerode**: $L$ è regolare se e solo se $N_L$​ ha indice finito
> - **Automa minimo**: gli stati sono le classi di $N_L$

> [!example] Immagina di leggere una parola a sinistra verso destra.
> A un certo punto, due prefissi $u$ $v$ possono essere considerati “uguali” per il linguaggio se, qualunque cosa ci aggiungi dopo, il risultato finale sarà sempre:
> - accettato in entrambi i casi
> - oppure rifiutato in entrambi i casi
> 
> Se questo succede, non ha senso tenerli in stati separati. Quindi l'automa minimo “fonde” tutte le parole che hanno lo stesso futuro possibile rispetto a $L$

**Costruzione dell'automa di Nerode:** dato un'automa a stati finiti deterministico $A = \langle Q, \Sigma, \delta, q_0, f \rangle$, costruire un'automa deterministico equivalente ad $A$ col minimo numero di stati:
- Possiamo ridurci al caso in cui tutti gli stati siano accessibili
- Gli stati dell'automa di Nerode sono le classi di $N_L$
- Ma anche unione di insiemi $L_q$ con $q \in Q$ (cf. Dimostrazione del th di Nerode $ii) \implies iii)$)
- Insomma la partizione di $\Sigma^*$ nelle classi di $N_L$ induce una partizione di $Q$
- Tale partizione gode delle seguenti 3 proprietà:
	1. L'insieme $F$ è unione di classi
	2. Se $p$ e $q$ sono nella medesima classe, allora pr ogni lettera $a \in \Sigma$, $\delta(p,a)$ e $\delta(q,a)$ cadono in una stessa classe
	3. è la meno fine tra le partizioni che soddisfano le due condizioni precedenti
- Il problema della minimizzazione si ridurrà pertanto alla ricerca della partizione $Q$ che soddisfa tali condizioni 

**Algoritmo di minimizzazione:** La partizione di $Q$ che soddisfa le precedenti Condizioni 1–3 si ottiene per approssimazioni successive: costruiamo la sequenza di partizioni $\Pi_n$ di $Q$ definite come segue:
- La partizione $\Pi_0$ (stati finali separati da quelli non finali) ha le sole due classi $F$ e $Q ∖ F$; 
- per ogni $n \geq 0$, le classi di $\Pi_{n+1}$ si ottengono spezzando quelle di $\Pi_n$ in modo da "separare" le coppie di stati che non soddisfano la Condizione 2 (spezzo una classe in modo da separare le coppie di stati che non soddisfano la condizione 2 vista prima, cioè che hanno l'immagine ch va in due classi differenti)
- ci arrestiamo non appena $\Pi_{n+1} = \Pi_n$ (continuo a spezzare classi finché non posso spezzare più nulla e mi fermo, creando così il minor numero di classi possibile)

L'automa minimo di $L$ si ottiene nel modo seguente:
- gli stati sono le classi della partizione $\Pi$, 
- per ogni classe $\mathscr C$ e ogni lettera $a \in \Sigma, \delta'(\mathscr C, a)$ è la classe che contiene tutti gli stati $\delta(q, a)$ con $q \in \mathscr C$, 
- lo stato iniziale è la classe $\mathscr C_0$ che contiene lo stato iniziale $q_0$ di $A$ 
- gli stati finali sono le classi contenute in $F$. 

> [!example] 
> Voglio minimizzare il seguente automa deterministico. Innanzitutto divido l'insieme degli stati in due classi: stati finali e stati non finali: $\Pi_0 =(Q \backslash F, F) = (\{q_0.q_1,q_2,q_3,q_4\},\{q_5,q_6\}$. > 
> ![[Pasted image 20260326162039.png|300]]
> 
> Dopodiché vado a vedere le lettere (a,b) su ogni stato, in quale classe mi riportano.
> Lo stato $q_0$, per esempio, con la lettera "a" mi porta in $q_3$, il quale fa parte della classe 0 $\implies$ scrivo 0. Faccio la stessa cosa in ogni classe per ciascuno stato (ogni colonnina). Gli stati che hanno colonnine uguali e che si trovano nella stessa classe, devono rimanere insieme: quindi $q_0$ e $q_3$ hanno entrambi colonnina 0,0 e devono restare insieme $\implies \Pi_1 = (\{q_0,q_3\}, \{q_1,q_2\}, \{q_4\}, \{q_5,q_6\}$
>> [!info] Anche se $q_2$ e $q_5$ hanno la stesa colonnina rimangono separati perché si trovano già in classi diverse
> 
> ![[Pasted image 20260326100813.png|300]]
> 
> Ripeto lo stesso identico meccanismo di prima e mi ritrovo con 5 classi
> ![[Pasted image 20260326100838.png|300]]
>> [!info] Quelli nella classe 3 sono gli stati finali
> 
> Rifaccio la tabellina con queste nuove classi e mi rendo conto che non c'è più nulla da spezzare. Il partizionamento è finito. Gli stati dell'automa saranno le classi di equivalenza con indice: 0, 1, 2, 3, 4
> 
> ![[Pasted image 20260326100850.png|300]]
> 
> Il nuovo automa minimizzato è il seguente e avrà 5 stati invece dei 7 iniziali
> 
> ![[Pasted image 20260326100911.png|300]]

**Procedura di raffinamento:**
<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260326102140.png" width="300">
</div>

**Algoritmo di minimizzazione:**
- *Input:* un automa deterministico $A = \langle Q, \Sigma, \delta, q_0, f \rangle$
- *Output:* la partizione $\Pi$ relativa all'automa minimo

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260326102153.png" width="300">
</div>
### Lemma di iterazione (proprietà dei linguaggi regolari)

**Lemma:** sia $L$ un linguaggio regolare. Esiste un intero n tale che ogni parola $w \in L$ di lunghezza $|w| \geq n$ si può fattorizzare $w = xyz$ con $y \neq \epsilon$ e $xy^mz \in L$ per ogni $m \geq 0$. In altri termini se $L$ è un linguaggio regolare, in ogni parola di $L$ abbastanza lunga si può trovare un fattore non vuoto che può essere iterato (*pumped*) senza uscire dal linguaggio. Pertanto se un linguaggio non ha questa proprietà non può essere regolare. Ci sono, però, linguaggi che hanno questa proprietà e non sono regolari

> [!tip] In altre parole:
> Se prendiamo un linguaggio regolare, possiamo trovare un intero n per cui tutte le parole del linguaggio che siano più lunghe di questo intero n, possono essere scomposte in 3 pezzi: xyz. La parte centrale y non è la parola vuota, ma un fattore che può essere ripetuto n volte senza che si esca dal linguaggio. Un linguaggio che non supporta questa proprietà di iterazione (pumping), non può essere regolare, anche se esistono linguaggi che hanno la suddetta proprietà ma non sono regolari.

> [!example]
> 1. Mostriamo che $L = \{a^m b^m \mid m \geq 0\}$ non è regolare
> - per $m$ abbastanza grande si dovrebbe avere  
>     $a^m b^m = xyz, y \neq \epsilon, xy^n z \in L$ per ogni $n \geq 0$     
> - se $|x| \geq m$ allora    
>     ![[Pasted image 20260326163412.png|300]]
>     
>     - si ha $y = b^k$, con $k > 0$ e $xz = a^m b^{m-k})$>         
>     - ma allora $xy^0 z = xz \notin L$, il che è assurdo
>         
> - se $|z| \geq m$, si ragiona in maniera simmetrica
> - se $|x| < m$ e $|z| < m$ allora
> 	![[Pasted image 20260326164834.png|300]]
> 	- si ha $y = a^hb^k$ con $h,k > 0$
> 	- quindi $xy^2z = xa^hb^ka^hb^kz$ contiene il fattore $ba$
> 	- ma allora $xy^2z \notin L$, il che è assurdo
> - pertanto $L$ NON è regolare

*Dimostrazione:* 
- Per il teorema di Kleene, $L$ è accettato da un automa a stati finiti $A$
- Sia $n$ il numero degli stati
- Sia $w \in L$ e $|w| \geq n$. Scriviamo $w = a_1 a_2...a_k$, con $a_1,a_2,...,a_k \in A, k \geq n$ ($k$ numero di lettere, quindi abbiamo preso una parola che abbia lunghezza $\geq n$)
- Nel grafo di $A$ c'è un cammino $q_0 \xrightarrow[]{a_1} q_1 \xrightarrow[]{a_2} ... \xrightarrow[]{a_k} q_k \in F$ (c'è un cammino da uno stato iniziale ad uno finale che ah $w$ come etichetta)
- Dato che $k \geq n$, troveremo uno stato ripetuto $q_i = q_j, \ 0 \leq i < j \leq l$ (ci saranno almeno $n+1$ stati, quindi questi saranno più del numero di stati totale dell'automa. Ci saranno di conseguenza delle ripetizioni)
- Quindi, posto $x = a_1...a_i, y = a_{i+1}...a_j, z = a_{j+1}...a_k$ avremo $w = xyz, y \neq \epsilon$ e 

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260326165631.png" width="300">
</div>

- Pertanto $A$ accetta tutte le parole $xy^nz$ con $n\geq0$ 
- Cioè $xy^nz \in L$ per ogni $n\geq 0$
### Analisi lessicale
Come abbiamo visto inizialmente, l'analisi è la prima fase nella compilazione di un codice sorgente. Dobbiamo identificare gli elementi atomici che costituiscono il nostro codice e individuare a quale categoria lessicale fanno riferimento. Il programma sorgente viene dato in pasto all'Analizzatore lessicale, che restituirà un token al Parser ogni volta che questo deve procedere con la sua analisi sintattica. Contemporaneamente, all'invio del token, l'analizzatore scriverà nella tabella dei simboli gli attributi del token, che poi potranno essere recuperati dal parser.

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260411093231.png" width="300">
</div>

**Perché analisi lessicale e sintattica devono essere separate?** E' importante separare l'analisi lessicale da quella sintattica perché si ha un'enorme semplificazione nella progettazione del compilatore, una maggiore efficienza e portabilità (perché cambiando sistema, cambia la gestione dell'input) rispetto a che se fossero svolte insieme.

- *Token:* è un simbolo astratto che rappresenta un'unità lessicale (es.una parola chiave, un identificatore, una costante numerica,...). I simboli processati dal parser sono i token. A ogni token viene associato un pattern (classe/simbolo che rappresenta il ruolo del lessema)
- *Lessema:* è una sequenza di caratteri di un codice sergente associato a un token. L'analizzatore lessicale identifica i lessemi come istanze del token a cui sono associati (stringa reale del programma)
- *Pattern:* è una descrizione della forma che i lessemi possono avere per poter essere associati ad un determinato token. In pratica si tratta di un'espressione regolare (Il pattern è la regola che descrive tutti i lessemi validi per un certo token.)

Quindi il processo che viene seguito è: pattern $\rightarrow$ genera/riconosce $\rightarrow$ lessemi $\rightarrow$ classificati come $\rightarrow$ token. Preso come esempio y = x + 42;

| Lessema | Token      | Pattern                                    |
| ------- | ---------- | ------------------------------------------ |
| `y`     | IDENTIFIER | Letter (parola che inizia con una lettera) |
| `=`     | ASSIGN     | `=`                                        |
| `x`     | IDENTIFIER | Letter                                     |
| `+`     | PLUS       | `+`                                        |
| 42      | NUMBER     | `digit+`                                   |
| `;`     | SEMICOLON  | `;`                                        |
> [!info] Sintesi
> - i lessemi sono ciò che leggo nel codice
> - i token sono delle etichette astratte
> - i pattern sono delle regole (regex) che definiscono i lessemi

**Attributi del token:** se più di un lessema è associato al medesimo pattern è necessario tenere traccia di ulteriori informazioni sul particolare lessema letto dall'analizzatore lessicale. Tali informazioni sono registrate nella *tabella dei simboli*

| Token             | Descrizione informale              | Esempi            |
| ----------------- | ---------------------------------- | ----------------- |
| If                | `if`                               | `if`              |
| Else              | `else`                             | `else`            |
| Comparazione      | `<,>,<=,>=,==,!=`                  | `<,>,<=,>=,==,!=` |
| Id                | Lettera seguita da lettere e cifre | Valore,a,c 1,...  |
| Numero            | Costanti nuemriche                 | 3.14,0,...        |
| Stringa letterale | Qualunque cosa tra ""              | "Hello",...       |
**Classi di token (relative all'esempio):** 
1. Un token per ogni parola chiave
2. Token per gli operatori (uno per ciascuno o raggruppati)
3. Un unico token per gli identificatori
4. Uno o più token per le costanti
5. Un token per ogni segno di punteggiatura

> [!example] $E = M * C ** 2$
> 1. Id, puntatore alla riga di $E$ nella tabella dei simboli
> 2. Operatore assegnazione
> 3. Id, puntatore alla riga di $M$ nella tabella dei simboli
> 4. Operatore prodotto
> 5. Id, puntatore alla riga di $C$ nella tabella dei simboli
> 6. Operatore potenza
> 7. Numero, valore intero 2

**Come si svolge l'analisi lessicale:** I token sono identificati da espressioni regolari (*pattern*) e hanno un ordine di priorità (ad esempio: if potrebbe essere sia un id di una variabile chiamata if, sia il costrutto if che conosciamo tutti). L'analizzatore lessicale quindi va ad identificare il pezzo più lungo del testo da analizzare che possa essere considerato un lessema, dopodiché restituisce il primo token che corrisponde a questo lessema. Questo token sarà quello di massima priorità, ai cui pattern posso associare il lessema.

> [!example] `int interesse = 5;`
> - il più lungo lessema iniziale è `int`, che corrisponde ai pattern dei token int e identificatore
> - l'analizzatore restituisce int, perché ha priorità maggiore
> - alla successiva attivazione analizza `interesse = 5;`
> - anche questo inizia con il lessema `int` ma non è il più lungo
> - il più lungo lessema iniziale è interesse, che corrisponde al pattern del token identificatore
> - l'analizzatore restituisce identificatore
> - alla successiva attivazione analizza =5

Il problema dell'analisi lessicale ha come dati: 
- una sequenza finita di pattern
- una parola $w$ (il testo da analizzare). 

Avendo questi dati devo determinare: 
- il più lungo prefisso $u$ di $w$ che appartenga al linguaggio denotato dall'espressione regolare
- il più piccolo indice $j$ tale che $u$ appartiene al linguaggio denotato (quindi devo capire qual è l'espressione regolare con cui questa parola corrisponde e, se ce n'è più di una, scegliere quella con priorità maggiore). 

**Automa per l'analisi lessicale:** per ognuna delle espressioni $E_j, j = 1, ..., n$ costruiamo un'automa a stati finiti $A_j = \langle Q_j, \Sigma, \delta_j, \{i_j\}, \{f_j\} \rangle$ che accetta il linguaggio denotato da $E_j$. Possiamo supporre che tali automi non abbiano stati comuni, quindi costruiamo l'automa $A$ che accetta il linguaggio denotato da $F = E_1 + ... + E_n$ con:
- Gli stati, le transizioni e gli stati finali di $A_1, A_2,...,A_n$  
- Un nuovo stato iniziale $i_0$
- Delle $\epsilon-$transizioni da $i_0$ agli stati iniziali originali degli automi $A_1, A_2,...,A_n$ (a differenza della costruzione vista nel teorema di Kleene, qui non aggiungiamo un nuovo stato finale perché vogliamo vedere da quali di miei n automi la parola viene accettata, quindi a quale espressione regolare corrisponde)

Con la costruzione dei sottoinsiemi otteniamo un'automa deterministico $A_D = \langle P, \Sigma, \delta, q_0, F' \rangle$ che accetta il linguaggio denotato da $F$. Osserviamo che :
- Gli stati di $A_D$ sono sottoinsiemi di $Q = \{i_0\} \cup \displaystyle\bigcup_{j=1}^n Q_j$
- Si ha che $w$ sta nel linguaggio denotato da $E_j$ se e soltanto se $f_j \in \hat\delta(q_0,u)$

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260409154305.png" width="300">
</div>

Troviamo il più lungo prefisso $u$ di $w$ accettato da $A_D$ e il più piccolo indice $j$ tale che $f_j \in \hat\delta(q_0,u)$. Gli stati di $A_D$ sono sottoinsiemi di $Q$ tra i quali c'è $\emptyset$ che funge da 'pozzo'. L'ultimo transito da uno stato finale prima dell'ingresso nel pozzo individua il più lungo prefisso di $w$ accettato da $A_D$

**Algoritmo:** 
*Input:* l'automa deterministico $A_D$ e un puntatore start al test da analizzare
*Output:* la massima lunghezza di un prefisso del testo accettato da $A_D$ e il token corrispondente

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260409154738.png" width="300">
</div>

<div style="display: flex; flex-direction: column; align-items: center; gap: 10px;">
  <img src="Pasted image 20260409154920.png" style="max-width: 50%; height: auto;">
  <img src="Pasted image 20260409154935.png" style="max-width: 50%; height: auto;">
  <img src="Pasted image 20260409155003.png" style="max-width: 50%; height: auto;">
</div>

**Lex:** è un generatore di analizzatori lessicali

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260409155211.png" width="300">
</div>
### Grammatiche regolari

*Definizione:* una grammatica a struttura di frase $G = \langle V,\Sigma,P,S \rangle$ è di tipo 3 se tutte le produzioni hanno le forme: $X \rightarrow aY, X \rightarrow a$, con $X,Y \in N, a \in \Sigma$ (cioè una variabile mi può produrre o simbolo terminale oppure un simbolo terminale con un'altra variabile). Nel caso in cui $S$ non compaia nei lati destri delle produzioni è ammessa anche la produzione $S \rightarrow \epsilon$ 

*Osservazione:* la produzione $S \rightarrow \epsilon$ serve esclusivamente per generare la parola vuota

> [!example] 
> Sia $G = \langle V,\Sigma,P,S \rangle$, con $\Sigma = \{a,b\}, V = \Sigma \cup \{S,X\}$ e con le produzioni: $S \rightarrow aS, S \rightarrow bX, S \rightarrow b, X \rightarrow aX, X \rightarrow bS, X \rightarrow a$. Si ha per esempio, $abbab \in L(G)$. Tale grammatica è di tipo 3

**Teorema:** un linguaggio è regolare se e soltanto se è generato da una grammatica di tipo 3

*Osservazione:* data $G$, si può effettivamente costruire un automa a stati finiti che accetta il linguaggio generato da $G$. Viceversa, dato un automa a stati finiti $A$, si può effettivamente costruire una grammatica di tipo 3 che genera il linguaggio accettato da $A$

Sia $G = \langle V,\Sigma,P,S \rangle$ una grammatica di tipo 3, per semplicità supponiamo che non abbia la produzione $S \rightarrow \epsilon$. Costruiamo l'automa non deterministico A come segue:
- Gli stati sono le variabili della grammatica G e la parola vuota
- Per ogni produzione $X \rightarrow aY$, nel grafo di $A$ c'è la freccia $X \xrightarrow{a} Y$, per ogni produzione $X \rightarrow a$, nel grafo di $A$ c'è la freccia $X \xrightarrow{a} \epsilon$
- Lo stato iniziale è $S$
- L'unico stato finale è $\epsilon$

Si verifica che $L(G) = L(A)$

**Inclusione $L(A) \subseteq L(G)$:** Sia $w \in L(A)$, scriviamo $w = a_1a_2…a_n$, con $a_1,a_2,…,a_n \in \Sigma, n \geq 1$. Nel grafo di $A$ ci sarà un cammino $S \xrightarrow{a_1} X_1 \xrightarrow{a_2} X_2 \xrightarrow{a_3}... \xrightarrow{a_{n-1}} X_{n-1} \xrightarrow{a_n} \epsilon$. Pertanto, nella grammatica $G$ ci saranno le produzioni $S → a_1 X_1,\ X_1 → a_2 X_2,\ …,\ X_{n-1} → a_n$ e quindi la derivazione $S \implies a_1X_1 \implies a_1a_2X_2 \implies … \implies a_1a_2…a_n = w$. Ne segue che $w \in L(G)$

**Inclusione $L(G) \subseteq L(A)$:** Sia $w \in L(G)$, scriviamo $w = a_1a_2…a_n$, con $a_1,a_2,…,a_n \in \Sigma, n \geq 1$. La derivazione di $w$ in $G$ avrà necessariamente la forma $S \implies a_1X_1 \implies a_1a_2X_2 \implies … \implies a_1a_2…a_n = w$ con $S → a_1 X_1,\ X_1 → a_2 X_2,\ …,\ X_{n-1} → a_n$ produzioni di $G$. Pertanto nel grafo di $A$ ci sarà un cammino $S \xrightarrow{a_1} X_1 \xrightarrow{a_2} X_2 \xrightarrow{a_3}... \xrightarrow{a_{n-1}} X_{n-1} \xrightarrow{a_n} \epsilon$. Ne segue che $w \in L(A)$

Quindi $L(G) = L(A)$

**Dalla grammatica all'automa:** 

*Osservazione:* se $S \rightarrow \epsilon$ è una produzione di $G$ (e quindi $S$ non campare nei lati destri delle produzioni), allora occorrerà aggiungere $S$ all'insieme degli stati finali di $A$ 

**Dall'automa alla grammatica:** sia $A = \langle Q, \Sigma, \delta, q_0, F \rangle$ un automa a stati finiti privo di $\epsilon-$transizioni. Per semplicità supponiamo che $\epsilon \notin L(A)$. Costruiamo la grammatica $G$ nel modo seguente:
- Le variabili della grammatica $G$ sono gli stati dell'automa $A$ 
- Per ogni freccia $X \xrightarrow{a} Y$ nel grafo di $A$ aggiungiamo a $G$ la produzione $X → aY$, se $Y$ è uno stato finale, aggiungiamo a $G$ anche la produzione $X → a$
- Il simbolo iniziale è lo stato iniziale $S = q_0$

Si verifica che $L(G) = L(A)$ (dimostrazione analoga alla precedente)

*Osservazione:* se $\epsilon \in L(A)$ sarà necessario modificare la grammatica $F$ in modo da fargli generare anche la parola vuota. Se però S compare nel alto destro di qualche produzione, non è possibile aggiungere la produzione $S \rightarrow \epsilon$. Dovremo pertanto prima introdurre una nuova variabile $S'$ le cui produzioni avranno gli stessi lati destri delle produzioni di $S$, prendere $S'$ come simbolo iniziale al posto di $S$ e, infine, aggiungere la produzione $S' → \epsilon$

> [!example] 
> Se $G$ ha le produzioni $S → aS, S →bX, S →b, X → aX, X → bS, X → a$, allora $L(G)$ è accettata dall'automa
> 
> ![[Pasted image 20260409162414.png|300]]

> [!example]
> ![[Pasted image 20260409162548.png|300]]

*Definizione:* una grammatica si dice lineare se tutte le produzioni hanno la forma $X → uYv$ o $X → u$ con $X,Y \in N$ e $u,v \in \Sigma^*$

Una grammatica si dice lineare destra se tutte le produzioni hanno la forma $X → uY$ o $X → u$ con $X,Y \in N$ e $u,v \in \Sigma^*$. In altre parole, una grammatica è lineare se i lati destri delle produzioni contengono al più una variabile; è lineare destra se tali variabili compaiono solo alla fine di tale termine

Le grammatiche di tipo 3 considerate finora sono evidentemente lineari destre. Viceversa si può dimostrare che ogni grammatica lineare destra genera un linguaggio di tipo 3

**Teorema:** un linguaggio è accettato da un automa a stati finiti se e solo se è generato da una grammatica lineare destra

> [!example] 
> Sia $G$ la grammatica con le produzioni $S → aaX, S → \epsilon, X → bbS$, l'automa che accetta $L(G)$ si ottiene nel modo seguente (altrimenti non avrei delle etichette su ogni arco ma delle parole):
> 
> ![[Pasted image 20260409183732.png]]
> 
> La grammatica $G$, inoltre, è equivalente alla grammatica con le produzioni: $S' → \epsilon, S' → aY, S → aY, Y → aX, X → bZ, Z → bS, Z → b$
### Linguaggi non contestuali
Sono linguaggi di tipo 2, e lo strumento necessario che ci permette di eseguire l'analisi sintattica.

> [!example] Sia $P$ il linguaggio delle palindrome sull'alfabeto $\Sigma = \{a,b\}$
> 1. $\epsilon, a, b$ sono palindrome
> 2. Per ogni palindroma $w$, le parole $awa$ e $bwb$ sono palindrome
> 3. Ogni palindroma si ottiene su $\Sigma$ si ottiene applicando un numero finito di volte le regole 1 e 2:$P = \epsilon+a+b+aPa+bPb$
> 
> Abbiamo ottenuto un algoritmo per verificare, ricorsivamente, se una data parola è palindroma e un metodo per generare tutte le palindrome (grammatica): $S \rightarrow \epsilon, S \rightarrow a, S \rightarrow b, S \rightarrow aSa, S \rightarrow bSb$ 

*Definizione:* una grammatica si dice non contestuale o di tipo 2 (CFG) se tutte le produzioni sono della forma: $X \rightarrow \alpha$, con $X \in N, \alpha \in V^*$ (cioè se hanno a sinistra una singola variabile). Un linguaggio si dice non contestuale o di tipo 2 (CFL) se è generato da una grammatica non contestuale. 

*Osservazione:* ogni grammatica lineare destra, e in generale le grammatiche lineari, sono non contestuali, perché a sinistra delle produzioni c'è sempre una variabile da sola. Quindi, dato che sono generati da grammatiche regolari destra, i linguaggi regolari sono non contestuali. Essi costituiscono un sotto insieme dei linguaggi non contestuali. Non vale l'inverso.

> [!example]
> Sia $G$ la grammatica con le produzioni $S → \epsilon, S → aSb$. Si ha $S \implies aSb \implies aaSbb \implies … \implies a^nSb^n \implies a^nb^n$. Quindi $L(G) = \{a^nb^n | n \geq 0\}$. Pertanto $\{a^nb^n | n \geq 0\}$ è un linguaggio non contestuale. Invece non è un linguaggio regolare

*Proposizione:* la classe dei linguaggi regolari è strettamente inclusa in quella dei linguaggi non-contestuali

Fra tutti i linguaggi non contestuali ce n'è uno che rappresenta l'archetipo di questo tipo di linguaggi: il linguaggio Dick. Si usa raggruppare le produzioni con lo stesso lato sinistro, separando i lati destri con una barra verticale. Per esempio: $S → a|b|\emptyset | (S+S) | (SS) | S^*$

**Linguaggio di Dick:** $\Sigma_n = \{a_1,a_2,…,a_n,b_1,b-2,…,b_n\}$. Sia $G_n = \langle V,\Sigma_n, P, S \rangle$ la grammatica con unica variabile $S$ e produzioni $S → a_iSb_i, i = 1,2,…,n$, $S → SS$, $S → \epsilon$, il linguaggio $D_n$ generato  da $G_n$ è detto linguaggio di Dick

> [!example]
> Supponiamo $n = 2$ e poniamo $a_1 = (,b_1 =) , a_2 = [, b_2 =]$. Allora le parole $(),[], [()]([]) , (())[([]())]$ sono elementi di $D_2$, mentre $)($ non lo è. Sequenze di parentesi non ben bilanciate

> [!info]
> Costruiamo una grammatica che genera le espressioni regolari sull'alfabeto $\{a,b\}$
> 1. $\emptyset,a,b$ sono espressioni regolari
> 2. Se $E$ e $F$ sono espressioni regolari, allora lo sono anche $(E+F),(EF), E^*$
> 3. Tutte le espressioni regolari sull'alfabeto $\{a,b\}$ si ottengono applicando un numero finito di volte le due regole precedenti
> 
> Pertanto prendiamo $\Sigma = \{a,b,\emptyset,(,),+,^*\}$ e produzioni $S →a, S → b, S → \emptyset, S → (S+S), S → (SS), S → S^*$

*Notazione:* si usa raggruppare le produzioni con lo stesso lato sinistro, separando i lati destri con una barra verticale. Per esempio $S \rightarrow a|b|\emptyset| (S+S)|(SS)|S^*$

> [!example] $L = \{w \in \Sigma^* | |w|_a = |w|_b\}$ cioè stesso numero di $a$ e di $b$ in qualunque ordine
> Consideriamo i linguaggi $L_a = \{w \in \Sigma^* | |w|_a = |w|_b + 1\}$, $L_b = \{w \in \Sigma^* | |w|_b = |w|_a + 1\}$, rispettivamente una $a$ e una $b$ in più. Verifichiamo le identità: $L = aL_b + bL_a + \epsilon, L_a = LaL, L_b = LbL$
> - ogni $w \in L$ si scrive $w = au$ con $u \in L_b$ o $w = bu$ con $u \in L_a$ o $w = \epsilon$
> - ogni $w \in L_b$ si scrive $w = ubv$ con $u,v \in L$
> 
> ![[Pasted image 20260411094136.png|300]]
> 
> Quindi $L$ è generato dalla grammatica con produzioni $S → aB | bA | \epsilon, A → SaS, B → SbS$

**Alberi di derivazione:** sono alberi con radice e con etichette sui nodi e i figli di ciascun nodo sono ordinati (da sinistra a destra). Ciò determina un ordinamento totale delle foglie. Sono una rappresentazione di una derivazione di una parola in una grammatica non contestuale.

*Definizione:* $T$ è un albero di derivazione della grammatica $G$ se verifica le seguenti condizioni:
1. L'etichetta della radice è $S$
2. Le etichette dei nodi interni sono elementi si $N$
3. Le etichette delle foglie sono elementi di $\Sigma \cup \{\epsilon\}$ (simboli terminali o parola vuota)
4. Le foglie con etichetta $\epsilon$ sono 'figli unici'
5. Se in nodo con etichetta $X$ ha figli etichettati nell'ordine $\alpha_1,…\alpha_k$, allora X → $\alpha_1,…,\alpha_k$ è una produzione di $G$

La parola che si ottiene leggendo, nell'ordine, le etichette delle foglie è la parola associata a $T$

> [!example] Nell'esempio di prima: 
> ![[Pasted image 20260411094509.png|300]]

*Proposizione:* siano $G$ una grammatica non contestuale e $w \in \Sigma^*$. Si ha $w \in L(G)$ se e solo se esiste un albero di derivazione di $G$ associato a $w$

> [!example] $S → aB | bA | \epsilon, A → SaS, B → SbS$
> La parola $ba$ ha 2 derivazioni:
> 1. $S \implies bA \implies bSaS \implies baS \implies ba$
> 2. $S \implies bA \implies bSaS \implies bSa \implies ba$
> 
> Ma un solo albero:
> 
> ![[Pasted image 20260411094606.png|300]]

> [!example]
> Sia $G$ la grammatica con $\Sigma = \{0,1,*,+\}, N = \{S\}$ e produzioni $S → S + S | S * S | 0 | 1$, allora $L(G) = (\{0,1\},\{+,*\})^*\{0,1\}$. La parola $1+1*0$ ha 2 alberi di derivazione:
> 
> ![[Pasted image 20260411094639.png|300]]
> 
> Dal punto di vista 'semantico' il primo corrisponde alla 'computazione' di $1+(1*0)$ e il secondo alla 'computazione' di $(1+1)*0$

**Grammatiche ambigue:** 

*Definizione:* una grammatica si dice NON ambigua se a ogni parola del linguaggio generato corrisponde un'unico albero di derivazione. Al contrario, le grammatiche ambigue sono quelle grammatiche dove ad ogni parola del linguaggio generato corrispondono più alberi di derivazione.

> [!example] 
> Sia $G'$ la grammatica con $\Sigma = \{0,1,*,+\}, N = \{S,P,X\}$ e produzioni $S → P+S|P, P → X*P|X, X → 0|1$. È equivalente a quella dell'esempio precedente ma non ambigua
> 
> ![[Pasted image 20260411094747.png|300]] 

*Definizione:* un linguaggio non contestuale si dice inerentemente ambiguo se non è generato da nessuna grammatica non contestuale non ambigua

> [!example] 
> Sia $G$ la grammatica con le produzioni $S → XA | AY, X → aXb| ab, A → aA|a, Y → bYa |ba$. Si ha  $L = \{ a^mb^ma^n|m,n > 0\} \cup \{a^mb^na^n|m,n > 0\}$. $G$ è ambigua 
> 
> ![[Pasted image 20260411094824.png|300]]
> 
> $L$ è inerentemente ambiguo
### Semplificazioni nelle grammatiche non contestuali
Ho due problemi: 
- di **ricognizione**: stabilire se una parola fa parte del linguaggio generato
- di **parsing**: data una parola che è nel linguaggio, trovare una derivazione

Sono in generale molto complessi, ma si semplificano se si evitano:
1. $X → \epsilon$ ($\epsilon$-produzioni)
2. $X → Y$ (produzioni unarie)
3. Le variabili che non compaiono nelle forme sentenziali (inaccessibili)
4. Le variabili da cui non si derivano parole prive di variabili (improduttive)

**$\epsilon-$produzioni e produzioni unarie:** 

*Definizione:* sia $G$ una grammatica non contestuale. Le produzioni della forma $X → \epsilon$ si dicono $\epsilon-$produzioni. Le produzioni della forma $X → Y$, con $X, Y$ variabili si dicono produzioni unarie

L'ideale sarebbe avere a sinistra sempre una variabile, e a destra avere o un terminale, o una sequenza di almeno due lettere variabile-terminale.

> [!example] Sia $G$ la grammatica con $N = \{S,O,P,E\}, \Sigma = \{a,b,x\}$ e con le produzioni: $S → aOb, O → P |aOb| OO, P → x | E, E → \epsilon$
> - $E → \epsilon$ è una $\epsilon-$produzione 
> - $O → P$ e $P → E$ sono produzioni unarie

*Osservazione:* in assenza di $\epsilon-$produzioni e produzioni unarie, se $\alpha \implies \beta$, allora 
- $|\beta| > |\alpha|$
- $|\beta| = |\alpha|$ ma $\beta$ contiene un terminale in più di $\alpha$ 

Pertanto, una parola di lunghezza $n$ può essere generata in $2n − 1$ passi al più (perché ad ogni passo o la mia forma sentenziale si allunga, oppure una variabile viene sostituita da un terminale).

*Proposizione:* ogni linguaggio non contestuale è effettivamente generato da una grammatica non contestuale priva di $\epsilon-$produzioni, eccettuata, eventualmente, la produzione $S → \epsilon$, ove $S$ è il simbolo iniziale. Inoltre, si può assumere che il simbolo iniziale non compaia nei lati destri delle produzioni. Visto che le parole vuote ci servono nella generazione dei nostri linguaggi, dobbiamo mantenere la possibilità di effettuare una $\epsilon-$produzione. La possiamo mettere al primo passo $S → \epsilon$, e solo se $S$ non compare nei lati destri delle produzioni. Quindi questa produzione $S → \epsilon$, c'è, ma è come se non ci fosse perché viene utilizzata solo ed unicamente per generare la parola vuota. Se non la uso al primo passo, non la posso utilizzare più.

**Costruzione della grammatica:**

*Definizione:* una variabile $X$ di una grammatica $G$ si dice annullabile se $X \implies^* \epsilon$ (cioè se genera la parola vuota). Ciò vale non solo per le produzioni dirette come $X→\epsilon$, ma anche per le produzioni concatenate come $X → Y, Y → W, W → \epsilon$. Quando andiamo a costruire una grammatica possiamo avere diversi casi:
1. *$S$ non è annullabile:* quindi non è possibile generare la parola vuota con il linguaggio. A partire da $G$ costruiamo una nuova grammatica $G'$ nel modo seguente:
	1. Aggiungiamo a $G'$ tutte le produzioni ottenute cancellando dei lati destri delle produzioni di $G$ una o più occorrenze di variabili annullabili
	2. Cancelliamo le $\epsilon-$produzioni

*Asserto:* La grammatica $G'$ è equivalente a $G$. 

*Dimostrazione:* 
- Per tutte le nuove produzioni $A → \gamma$ si ha $A \implies *\gamma$ in $G$. Quindi $L(G') \subseteq L(G)$ 
- Per l'inclusione inversa, faremo vedere che per ogni albero di derivazione di $G$ c'è un albero di derivazione di $G'$ associato alla stessa parola. 

Procediamo per induzione sul numero dei nodi dell'albero.

<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Pasted image 20260423100401.png" style="width: 100%; border-radius: 8px;">
  </div>
  <div style="flex: 1.5;">
  Dove X1,X2,… sono le produzioni di S, mentre w1,w2,… sono le produzioni delle produzioni di S. Guardando l'albero della vecchia grammatica G, e poi quello nuovo della grammatica G', possiamo vedere come, anche se eliminiamo le epsilon-produzioni (che nel nuovo linguaggio non servono più), posso ancora generare tutto ciò che generavo prima</div>
</div>

2. *$S$ è annullabile:* a partire da $G$ costruiamo una nuova grammatica $G'$ nel modo seguente:
	1. Procedendo come nel caso precedente si ottiene una grammatica che genera le stesse parole di $G$ eccetto $\epsilon$
	2. Se $S$ non compare nei lati destri, basta aggiungere la produzione $S → \epsilon$
	3. Se $S$ compare nei lati destri si aggiungono una nuova variabile $S'$, che sarà il nuovo simbolo iniziale e le produzioni $S' → \epsilon$ e $S' → S$. In questo modo ho fatto scomparire il simbolo iniziale dei lati destri delle produzioni e posso aggiungere $\epsilon$

*Asserto:* La grammatica $G'$ è equivalente a $G$

**Ricerca delle variabili annullabili:** costruiamo gli insiemi sempre più grandi di variabili $W_n$ come segue: 
- $W_1 = \{X \in N |$ c'è una produzione $X → \epsilon$ in $P\}$,
- $W_n = \{X \in N |$ c'è una produzione $X → \gamma$ in $P$ con $\gamma \in W^*_{n-1}\}, n \geq 2$ 

> [!info] In altre parole:
> 1. all'inizio prendo l'insieme di tutte le
> variabili che producono $\epsilon$.
> 2. Poi costruisco un altro insieme in cui
> oltre a queste variabili, aggiungo tutte
> quelle produzioni che hanno un lato destra
> fatto interamente da variabili
> annullabili (che conosco già grazie al
> passo precedente).
> 3. Ripeto il secondo passo più volte fino a
> quando gli insiemi non cresceranno più
> e che ci indicheranno che tutte le
> variabili annullabili sono state trovate.

*Osservazione:* 
- $W_n$ contiene solo variabili annullabili; 
- $W_n$ contiene tutte le variabili che si annullano in $n$ passi; 
- $W_1 \subseteq W_2 \subseteq ... \subseteq W_n \subseteq ... \subseteq N$; 
- Pertanto, per qualche $n$, $W_n = W_{n-1} = W$ . $W$ è l'insieme delle variabili annullabili.

**Algoritmo:** inizialmente ho in input una grammatica non contestuale $G = \langle V, \Sigma, S, P \rangle$. In output ho l'insieme $W$ delle variabili annullabili;

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260423100813.png" width="300">
</div>

> [!example] 
> Sia $G$ la grammatica con $N = \{S,O,P,E\}$, $\Sigma = \{a,b,x\}$ e con le produzioni: $S → aOb, O → P | aOb | OO, P → x | E, E → \epsilon$
> 
> ![[Pasted image 20260423095536.png|300]]
> 
> Le variabili annullabili sono: $E,P,O$. Quindi $G$ è equivalente alla grammatica con le produzioni: $S → aOb | ab, O → P | aOb | OO | ab | O, P → x | E$

*Proposizione:* Ogni linguaggio non contestuale è effettivamente generato da una grammatica non contestuale priva di produzioni primarie e di $\epsilon -$produzioni, eccettuata, eventualmente, la produzione $S → \epsilon$ , ove $S$ è il simbolo iniziale. Inoltre, si può assumere che il simbolo iniziale non compaia nei lati destri delle produzioni.

**Costruzione della grammatica:** A partire da $G$ costruiamo una nuova grammatica $G'$ nel modo seguente: 
1. si eliminano le $\epsilon$-produzioni con la costruzione precedente, 
2. per ogni coppia di variabili $A, B$ tali che $A \implies *B$ ($A$ deriva $B$ in uno o più passi), si aggiungono i lati destri delle produzioni di $B$ a quelli delle produzioni di $A$ (tutte le cose che produce $B$ le aggiungo a quelle che produce $A$) 
3. si cancellano le produzioni unarie

*Asserto:* La grammatica $G'$ è equivalente a $G$. 

*Dimostrazione:* 
- Per tutte le nuove produzioni $A → \gamma$ si ha $A \implies *\gamma$ in $G$. Quindi $L(G') \subseteq L(G)$. 
- Per l'inclusione inversa, faremo vedere che a ogni albero di derivazione di $G$ c'è un albero di derivazione di $G'$ associato alla stessa parola. Procediamo per induzione sul numero dei nodi dell'albero.

1. $S$ ha più di un figlio: 

<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Pasted image 20260423101234.png" style="width: 100%; border-radius: 8px;">
  </div>
  <div style="flex: 1.5;">
  L'albero G' è uguale perché abbiamo tolto solo le produzioni unarie 
  </div>
</div>

2. $S$ ha solo un figlio: 

<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Screenshot From 2026-04-23 10-13-12.png" style="width: 100%; border-radius: 8px;">
  </div>
  <div style="flex: 1.5;">
  Abbiamo tolto tutte le produzioni unarie S → Y1, S → Y2, S → Y3 e abbiamo aggiunto tutto a destra di S. Quindi ora S le deriva direttamente
  </div>
</div>

**Ricerca delle derivazioni unarie:** Costruiamo gli insiemi di coppie di variabili $T_n$ con tutte le produzioni unarie dirette. Ad ogni passo, vado ad aggiungere all'insieme, tutte le coppie $X,Y$ per cui ci sono già $X,Z$ e $Z,Y$ nell'insieme. Perché so che da $X → Z$ e da $Z → Y$, allora da $X →Y$. Procediamo così fin quando posso aggiungere coppie all'insieme, altrimenti vorrebbe dire che abbiamo trovato tutte le derivazioni unarie.
- $T_1 = \{(X,Y) \in N \times N | \ \text{c'è la produzione} \ X → Y \ \text{in P}$
- $T_n = T_{n-1} \cup \{(X,Y \in N \times N | \ \text{esiste} \ Z \in N, (X,Z), (Z,Y) \in T_{n-1} \}, n \geq 2$ 

> [!info]
> $T_1$ è l'insieme delle coppie $X,Y$ per cui $Y$ si deriva da $X$ in un passo, $T_2$ in al massimo due passi, $T_3$ in al massimo tre passi e $T_n$ in al massimo $2n-1$ passi. Quindi il procedimento è abbastanza veloce perché la lunghezza massima sarà il numero delle variabili.

*Osservazioni:* 
- $T_n$ contiene le coppie di variabili $(X,Y)$ tali che $X \implies ^*Y$ in al più $2^{n-1}$ passi; 
- $T_1 \subseteq T_2 \subseteq ... \subseteq T_n \subseteq ... \subseteq N \times N$  
- Pertanto, per qualche $n, T_n = T_{n-1} = T$. $T = \{(X,Y) \in N \ times N | X \implies *T\}$ 
- si ha $A \implies *B$ se e solo se c'è un cammino da $A$ a $B$ nel grafo $(N, E)$ con insieme di archi $E = (X,Y) \in N \times N | X → Y \ \text{in P} \}$ 

> [!example] 
> Sia $G$ la grammatica con $N = \{S,O,P,E\}$, $\Sigma = \{a,b,x\}$ e con le produzioni: $S → aOb, O → P | aOb | OO | ab | O, P → x | E$. Il grafo delle produzioni unarie è:
> 
> ![[Pasted image 20260423102435.png|300]]
> 
> Quindi $G$ equivalente alla grammatica con produzioni: $S → aOb | ab, O → aOb | OO | ab | x, P → x$

**Variabili improduttive e variabili inaccessibili:** 

*Definizione:* Sia $G = \langle V, \Sigma, P, S \rangle$  una grammatica non contestuale. Una variabile $X$ si dice **produttiva** se esiste $w \in \Sigma^*$ tale che $X \implies ^*w$, **improduttiva** altrimenti. Una variabile $X$ si dice **accessibile** se compare in qualche forma sentenziale, **inaccessibile** altrimenti.

> [!example ]
> Sia $G$ la grammatica con $N = \{N, X, Y, Z\}$, $\Sigma = \{a,b\}$ e con le produzioni: $S → aXa |aa, X → bSb | Y , Z → a$.
> - Le variabili $S$ e $X$ sono produttive e accessibili, invero, $S \implies *aa, X \implies *baab, S \implies *S, S \implies *aXa$
> - $Y$ è accessibile ma improduttiva, 
> - $Z$ è produttiva ma inaccessibile.

*Osservazione:* Data una grammatica non contestuale $G$, la grammatica ottenuta cancellando le variabili improduttive o inaccessibili e le produzioni che le contengono è equivalente a $G$

> [!example] 
> La grammatica con le produzioni: $S → aXa | aa, X → bSb | Y, Z → a$ è equivalente a quella con le sole produzioni $S → aXa | aa, X → bSb$

**Ricerca delle variabili produttive:** Costruiamo gli insiemi di variabili $W_n$ come segue: 
- $W_1 = \{X \in N | \ \text{c'è una produzione} \ X → w \ \text{in P con} \ w \in \Sigma^*\}$
- $W_n = \{X \in N | \ \text{c'è una produzione} \ X → \gamma \ \text{in P con} \ \gamma \in (W_{n-1} \cup \Sigma)^*\}, n \geq 2$

In altre parole creiamo un insieme di variabili ($W_n$) e gli aggiungiamo tutte le variabili che sono direttamente produttive, cioè le variabili che producono subito tutti simboli terminali (comprese le $\epsilon-$produzioni). Dopodiché gli aggiungo tutte quelle variabili che hanno un lato destra costituito da terminali e variabili che già so essere produttive. Ripeto questo procedimento finché l'insieme non cresce più. Tutte le variabili che rimangono fuori da questo insieme sono improduttive e quindi possono essere buttate via.

*Osservazione:*
- $W_n$ contiene solo variabili annullabili; 
- $W_n$ contiene tutte le variabili che si annullano in $n$ passi; 
- $W_1 \subseteq W_2 \subseteq ... \subseteq W_n \subseteq ... \subseteq N$; 
- Pertanto, per qualche $n$, $W_n = W_{n-1} = W$. $W$ è l'insieme delle variabili annullabili.

L'algoritmo per le variabili produttive prende in input una grammatica non contestuale $G = \langle V, \Sigma, P, S \rangle$ e restituisce in output l'insieme $W$ delle variabili produttive e l'insieme $N$ delle variabili improduttive;

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260423103009.png" width="300">
</div>

**Ricerca delle variabili accessibili:** costruiamo gli insiemi di variabili $K_n$ come segue:
- $K_0 = \{S\}$
- $K_n = K_{n-1} \cup \{$variabili che compaiono nei lati destri delle produzioni delle $X \in K_{n-1} \}, n \geq 1$

In altre parole creiamo un insieme di variabili ($K_n$) e gli aggiungiamo $S$ che sappiamo essere accessibile. Poi aggiungeremo tutte le variabili che compaiono nei lati destri delle produzioni di $S$. In seguito aggiungeremo tutte le variabili che hanno nel lato sinistro delle variabili che sappiamo già essere  accessibili. Ripetiamo il procedimento finché l'insieme non smette di crescere.

*Osservazioni:*
- $K_n$ contiene le variabili che compaiono nelle forme sentenziali derivabili in $n$ passi dal simbolo iniziale
- $K_1 \subseteq K_2 \subseteq ... \subseteq K_n \subseteq ... \subseteq N$
- Pertanto, per qualche $n$, $K_n = K_{n-1} = K$. $W$ è l'insieme delle variabili produttive.
- Si ha $A \in K$ se e solo se c'è un cammino da $S$ a $A$ nel grafo $(N,E)$ con insieme di archi $E = \{(X,Y) \in N \times N | X → \alpha Y \beta$ in P $\}$ 

L'algoritmo per le variabili accessibili prende in input una grammatica non contestuale $G = \langle V, \Sigma, P, S \rangle$ e restituisce in output l'insieme $K$ delle variabili produttive e l'insieme $N$ delle variabili improduttive;

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260423110309.png" width="300">
</div>

**Procedura di riduzione:** 
1. determinare le variabili produttive, 
2. eliminare le variabili improduttive e le produzioni che contengono tali variabili, 
3. determinare le variabili accessibili della grammatica ottenuta, 
4. eliminare le variabili inaccessibili e le produzioni che contengono tali variabili.

*Osservazione:* Non è possibile eseguire prima i passi 3–4 e poi i passi 1–2: l'eliminazione di variabili improduttive può creare nuove variabili inaccessibili.

> [!example] 
> Sia $G$ la grammatica $N = \{E,F,T,R\}, \Sigma = \{+,*,-,(,)\}$, simbolo iniziale $E$ e produzioni: $E → E+E|T|F, F → E*E|(T)|a, T → E-T, R → E-F$
> 
> ![[Pasted image 20260423110536.png|300]]

> [!example]
> Sia $G$ la grammatica con variabili $S,O,P$, terminali $a,b,x$ e produzioni $S→aOb|ab, O → aOb|OO|ab|x, P → x$. Tutte le variabili sono produttive
> 
> ![[Pasted image 20260501123226.png|300]]
> 
> Si ottiene la grammatica con le produzioni $S → aOb|ab, O → aob|OO|ab|x$
### Forma normale di Chomsky

*Definizione:* una grammatica non contestuale si dice in forma normale di Chomsky se ha solo produzioni dei tipi
- $X → YZ$ con $X,Y,Z \in N$
- $X → a$ con $X \in N, a \in \Sigma$
- $S → \epsilon$ ma solo a condizione che $S$ non compaia nei lati destri delle produzioni

> [!example]
> La grammatica con le produzioni $S → AB| \epsilon, A → AA|a, B → BB|b$ è la forma normale di Chomsky. Genera $L = \{\epsilon\} \cup \{a^mb^n|m,n > 0\}$

*Osservazione:* negli alberi di derivazione di una grammatica in forma normale di Chomsky
- Le foglie dono figli unici
- I nodi interni costituiscono un albero binario completo

> [!example] 
> $S → AB|\epsilon, A → AA|a, B → BB|b$. Un albero di derivazione di $aaabb$ è 
> 
> ![[Pasted image 20260501123457.png|300]]

*Teorema:* ogni linguaggio non contestuale è effettivamente generato da una grammatica non contestuale in forma normale di Chomsky

*Procedura:* 
1. Eliminare le $\epsilon$-produzioni e produzioni unarie; restano solo produzioni $X → a$ con $X \in N$ e $a \in \Sigma$ e produzioni $X → \gamma$ con $|\gamma| \geq 2$
2. Ridursi al caso $\gamma \in N^*$ (caso in cui $\gamma$ contiene solo variabili). Per ogni terminale $a$ che compare, non da solo, in qualche lato destro
	- Introduco una nuova variabile $A$ e una nuova produzione $A → a$
	- Sostituisco tutte le occorrenze di $a$ nei lati destri delle produzioni con $A$ (salvo che a non sia l'intero lato destro)

> [!example] 
> Sia $G$ la grammatica con produzioni $S → aOb|ab, O → aOb|OO|ab|x$. Con la procedura indicata, otteniamo: $S → AOB|AB, O → AOB|OO|AB|x, A → a, B → b$ 

Restano solo produzioni $X → a$ con $X \in N$ e $a \in \Sigma$ e produzioni $X → X_1 →…→ X_k$ con $X,X_1,…,X_k \in N, k \geq 2$ (salvo l'eventuale $S → \epsilon$)
3. Ridursi al caso $k = 2$: per ogni produzione $X → X_1 →…→ X_k$ con $k \geq 3$
	- Introduco una nuova variabile $Z$ e sostituisco la produzione $X → X_1 →…→ X_k$ con la coppia di produzioni $X → X_1 →…→ X_{k-2}Z$, $Z → X_{k-1}X_K$
	- Ripeto il passo precedente finché necessario

> [!example]
> Sia $G$ la grammatica con produzioni $S → AOB|AB, O → AOB|OO|AB|x, A → a, B → b$. Con la procedura indicata, otteniamo: $A → AZ|AB, O → AZ|OO||AB|x, Z → OB, A → a, B → b$. La grammatica ottenuta, equivalentemente a $G$, è in forma normale di Chomsky

**Forma normale di Greibach:**

*Definizione:* una grammatica non contestuale $G$ si dice in forma normale di Greibach se ha solo produzioni dei tipi:
- $X → a\gamma$, con $a \in \Sigma$ e $\gamma \in N^*$
- $S → \epsilon$, ma solo a condizione che $S$ non compaia nei lati destri delle produzioni

> [!example]
> La grammatica con le produzioni $S → aSB|aB, B → b$ è in forma normale di Greibach. Genera $L = \{a^nb^n|n>0\}$

*Osservazione:* in una grammatica in forma normale di Greibach, ogni derivazione di una parola $w \in L(G)$ di lunghezza $n$ richiede esattamente $n$ passi

> [!example]
> $S → aSB|aB, B → b$. Una derivazione di $a^3 b^3$: $S \implies aSB \implies aaSBB \implies aaaBBB \implies aaabBB \implies aaabbB \implies aaabbb$

*Teorema:* ogni linguaggio non contestuale è effettivamente generato da una grammatica non contestuale in forma normale di Greibach
### Algoritmo di Cocke-Kasami-Younger
*Problema di ricognizione:* ho in input una grammatica $G = \langle V, \Sigma, P ,S \rangle$ e una parola $w \in \Sigma^*$. In output ho 'SI' se $w \in L(G)$, 'NO' altrimenti

*Problema di parsing:* ho in input una grammatica $G = \langle V, \Sigma, P ,S \rangle$ e una parola $w \in L(G)$. In output ho una derivazione $S \implies \alpha_1 \implies \alpha_2 \implies … \implies \alpha_n \implies w$ di $w$ in $G$ (o equivalentemente un albero di derivazione)

**Algoritmo CYK:** una soluzione moderatamente efficiente al problema di ricognizione per i linguaggi non contestuali e, con opportune modifiche, al problema di parsing:
- Ideale per le grammatiche in forma normale di Chomsky
- Ma ogni grammatica non contestuale si può ridurre in tale forma

L'idea di tale algoritmo è:
- Siano $G$ una grammatica in forma normale di Chomsky e $w$ una parola
- Scriviamo $w$ come stringa di lettere $w = a_1a_2…a_n$
- Cerchiamo di calcolare, per $0 \leq i < j \leq n$, tutte le variabili da cui si può derivare il fattore $a_{i+1}a_{i+2}…a_j$
- Per $i = 0, j = n$, avrò le variabili da cui si può derivare $w$; quindi $w \in L(G)$ se e solo se tra esse il simbolo iniziale $S$ è una di tali variabili

**Il calcolo:** Dobbiamo calcolare gli insiemi $N_{ij} = \{X \in N | X \overset{*}{\Rightarrow} a_{i+1}a_{i+2}…a_j\}, 0 \leq i < j \leq n$, cioè tutte le coppie di non terminali che possono generare il pezzo di stringa da posizione $i + 1$ a $j$ (insieme delle variabili che mi generano il fattore $a_{i+1}a_{i+2}…a_j$). Per questo sto riempendo una tabella triangolare con sottostringhe corte nelle celle piccole e sottostringhe lunghe in quelle grandi. 
- Per $j = i+1, N_{ij} = \{X \in N | X \overset{*}{\Rightarrow} a_j\} = \{X \in N | X →a_j$ in $P \}$. Ovvero guardo ogni lettera $a_j$ e metto dentro $N_{j-1,j}$ tutte le variabili che producono quella lettera
- Per $j \geq i + 2$, si ha $X \overset{*}{\Rightarrow} a_{i+1}a_{i+2}…a_j$ se e solo se $X → YZ, Y \overset{*}{\Rightarrow} a_{i+1}a_{i+2}…a_h, Z \overset{*}{\Rightarrow} a_{h+1}a_{h+2}…a_j$, per opportuni $Y,Z \in N$ e $i<h<j$. Otteniamo così $N_{ij} = \bigcup\limits^{j-1}_{h = i+1} \{X \in N | \exists Y \in N_{ih}, Z \in N_{hj}, X → YZ$ in $P\}$ 

Per ogni coppia di variabili $Y,Z$, poniamo $Y \circ Z = \{X \in N | X → YZ$ in $P\}$ e, più in generale, per ogni coppia di insiemi di variabili $H,K \subseteq N$, $H \circ K = \bigcup\limits_{\substack{Y \in H \\ Z \in K}} Y \circ Z$. Allora la nostra formula ricorsiva diventa: $N_{ij} = \bigcup\limits^{j-1}_{h = i+1} N_{ih} \circ N_{hj}$. Anzi se conveniamo che $N_{ij} = \emptyset$ per $0 \leq j \leq i \leq n, N_{ij} = \bigcup\limits_{h=0}^n N_{ih} \circ N_{hj}$

> [!info] N.B.
> L'operazione di $\circ$ indica la concatenazione e poi la sostituzione con il lato sinistro

> [!example] 
> ![[Pasted image 20260503174624.png|300]]

**Algoritmo CYK (on-line):** prende in input una grammatica in f.n. di Chomsky e una parola $w = a_1a_2…a_n$. In output da la matrice di ricognizione $N_{ij}$ associata a $w$. Ritorna 'TRUE' se $w \in L(G)$, 'FALSE' altrimenti

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260503174651.png" width="300">
</div>

Analisi:
- Calcolare $N_{ih} \circ N_{hj}$ ha costo 1
- L'istruzione 4 ha costo $j-i-1 \leq n$ e viene eseguita $1+2+…+ (n-1) = n(n-1)/2$ volte
- Il costo totale è quindi $O(n^3)$

**Parsing:** sia $0 \leq i < j \leq n$ e $X \in N_{i,j}$
- Posso determinare un indice $h$ e variabili $Y \in N_{i,h}$ , $Z \in N_{h,j}$ tali che $X → YZ$ è una produzione di $G$
- Quindi una derivazione di $a_{i+1}a_{i+2}…a_j$ da $X$ si ottiene concatenando
	- $X \implies YZ$
	- Una derivazione di $a_{i+1}a_{i+2}…a_h$ da $Y$
	- Una derivazione di $a_{h+1}a_{h+2}…a_j$ da $Z$

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260503174842.png" width="300">
</div>

> [!example]
> ![[Pasted image 20260503174824.png|300]]

**Caso generale:** se $G$ non è in forma normale di Chomsky
1. Si mette $G$ in forma normale di Chomsky con l'algoritmo standard
2. Si esegue il parsing nella nuova grammatica con l'algoritmo CYK
3. Riesaminando al contrario il passo 2, si risale al parsing nella grammatica originale
### Lemma di iterazione per i linguaggi non contestuali
Sia $L$ un linguaggio regolare. Esiste un intero $n$ tale che ogni parola $w \in L$ di lunghezza $|w| \geq n$ si può fattorizzare $w = xyz$ con $y \neq \epsilon$ e  $xy^kz \in L$ per ogni $k \geq 0$

> [!example] 
> $L = \{a^mb^m | m > 0\}$ non è regolare. Questo perché $y$ corrisponderebbe ad un blocco di 'a' o ad un blocco di 'b'. Iterando quindi quante volte voglio 'a' o 'b' non saranno più della stessa lunghezza ($y$ potrerbbe essere anche a cavallo tra 'a' o 'b' e mi ritroverei fuori dal linguaggio). Se però prendessi un 'u' dentro le 'a' e una 'v' dentro le 'b' (con |a| = |b|), potreri iterarle parallelamente e non uscire dal linguaggio. Invero, se $xyz \in L$ con $y \neq \epsilon$, allora $xyyz \notin L$. Però per $x = y= z = \epsilon, un = a ^m, v = b^m$, si ha $$a^mb^m = xuyvz, xu^nyv^nz 0 a^{mn}b^{mn} \in L$$Si possono iterare due blocchi parallelamente 

Sia $L$ un linguaggio non contestuale. Esiste un intero $n$ tale che ogni parola $w \in L$ di lunghezza $|w| \geq n$ si può fattorizzare $w = xuyvz$ con $uv \neq \epsilon$ e $$xu^kyv^kz \in L$$
per ogni $k \geq 0$

*Osservazione:* se $L$ è un linguaggio non contestuale, in ogni parola di $L$ abbastanza lunga si possono trovare due fattori non entrambi vuoti che possono essere iterati simultaneamente (**pumped**) senza uscire dal linguaggio. Pertanto, un linguaggio che non ha questa proprietà non è un linguaggio non contestuale (ma ci sono linguaggi che hanno questa proprietà e non sono linguaggi non contestuali)

> [!example] 
> Mostriamo che $L = \{a^mb^mc^m | m > 0\}$ non è contestuale. 
> - Altrimenti, per $w$ abbastanza grande si dovrebbe avere $a^mb^mc^m = xuyvz, uv \neq \epsilon, xu^kyv^kz \in L$ per ogni $k \geq 0$
> - se $|xu| > m$, allora 
> 	- $xu$ inizia con $a^mb$, come pure $xu^2yv^2z$
> 	- ma l'unica parola di $L$ che inizia con $a^mb$ è $a^mb^mc^m$
> 	- quindi $xuyvz = xu^2 yv^2 z = a^mb^mc^m$, il che è ASSURDO
> - se $|vz|>m$, si ragiona in maniera simmetrica
> 
> ![[Pasted image 20260501124354.png|300]]
> 
> - se $|xu| \leq m$ e $|vz| \leq m$, allora 
> 	- $u = a^i, v = c^j$, con $i+j > 0$ 
> 	- allora $xu^2yv^2z = a^{m+i}b^mc^{m+j} \notin L$, il che è ASSURDO
> 
> ![[Pasted image 20260501124959.png|300]]

**Idea della prova:** 
<div style="display: flex; justify-content: center;">
  <img src="Screenshot From 2026-05-01 12-50-37.png" width="300">
</div>

*Dimostrazione:* 
- Sia $L = L(G)$, con $G$ priva di $\epsilon-$produzioni e produzioni unarie
- Sia $n$ la massima lunghezza delle parole di $L$ che hanno un albero un albero di derivazione di altezza $\leq Card(N)$
- Sia $w \in L$, con $|w| > n$
- In un cammino di lunghezza $> Card(N)$ dalla radice a una foglia nell'albero di derivazione c'è una variabile che si ripete
- $S \implies^* xAz, A \implies^* uAv, A \implies^* y, w = xuyvz, uv \neq \epsilon$ 
- $S \implies^* xAz \implies^* xuAvz \implies^* xuuAvvz \implies^* ... \implies^* xu^kAv^kz \implies^* xu^kyv^kz$
- Quindi $xu^kyv^kz \in L$ per ogni $k \geq 0$

<div style="display: flex; justify-content: center;">
  <img src="Screenshot From 2026-05-01 12-56-24.png" width="300">
</div>

### Parsing top-down 
In genere l'operazione di parsing richiede un tempo notevole $O(n^3)$ e non è accettabile nel caso dovessimo fare la compilazione di un programma molto complesso (perché come sappiamo, raramente un programma non presenta errori alla prima compilazione, quindi dovremmo fare più tentativi e di conseguenza effettuare la compilazione e il parsing più volte). Il problema di parsing ha come input una grammatica $G = \langle V,\Sigma,P,S \rangle$ e una parola $w \in L(G)$ e in output una derivazione $S \implies \alpha_1 \implies \alpha_2 \implies … \alpha_n \implies w$ di $w \in G$ (o equivalentemente un albero di derivazione)

*Osservazione:* a ogni derivazione corrisponde un albero di derivazione, ma a un albero di derivazione possono corrispondere più derivazioni 

> [!example]
> ![[Pasted image 20260501130348.png|300]]
> 
> Sostituendo prima una variabile invece che un'altra, otteniamo sempre lo stesso albero di parsing. Sarebbe quindi bene mettersi d'accordo su un ordine da seguire per sostituire le variabili. Per fare ciò usiamo le derivazioni sinistra

*Definizione:* sia $G = \langle V,\Sigma,P,S \rangle$ una grammatica non contestuale e $\alpha,\beta \in V$. Scriveremo $\alpha \implies_L\beta$ se $\beta$ si ottiene da $\alpha$ sostituendo la variabile più a sinistra con il lato destro di una sua produzione. Una derivazione sinistra di una parola $w$ è una sequenza $S = \alpha_0 \xrightarrow[L]{} \alpha_1 \xrightarrow[L]{} … \xrightarrow[L]{} \alpha_n = w$ (con '=' rappresenta una conseguenza e '=>' rappresenta una derivazione)

*Proposizione:* c'è una corrispondenza biunivoca tra le derivazioni sinistre di una parola $w$ e gli alberi di derivazione associati a $w$. <u>Ogni albero di derivazione corrisponde ad una e una sola derivazione sinistra e viceversa</u>. Inoltre, se ho una derivazione sinistra, per costruire l'albero in realtà mi basta ricordare le produzioni che ho utilizzato, perché la variabile che ho sostituito già la so, è quella più a sinistra.

*Corollario:* data una grammatica non contestuale $G$, ogni parola di $L(G)$ si ottiene dal simbolo iniziale $S$ con una derivazione sinistra. Questo perché le parole si ottengono con una derivazione, alla derivazione corrisponde un albero, e all'albero corrisponde una derivazione sinistra.

**Endmaker:** sia $G = \langle V,\Sigma,P,S \rangle$ una grammatica non contestuale. Per motivi di efficienza, è utile modificare $G$ usando un marcatore di fine parola (**endmaker**). Per fare ciò aggiungo:
- Un nuovo simbolo iniziale $S'$
- Un nuovo terminale #
- La produzione $S' → S$#

E anziché il parsing di $w \in \Sigma^*$ farò il parsing di $w$#. Quindi il cancelletto lo genererò partendo da $S'→S$# e poi avrò $S$ da cui devo derivare $w$.

> [!example]
> ![[Pasted image 20260508100410.png|300]]

**Parsing predittivo:** 
<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260501130606.png" width="300">
</div>

Abbiamo l'input che dobbiamo derivare dal simbolo iniziale, quindi passo passo esaminiamo il nostro input da sinistra verso destra e iniziamo a costruire la derivazione dall'alto in basso del nostro albero di derivazione. Ad ogni passo avremo una parte di input già analizzato (quello alla sinistra delle tabelline dell'esempio precedente), l'analisi (ciò che avevamo derivato dal nostro simbolo iniziale che faceva match con quello di sopra), a destra abbiamo l'input ancora non analizzato, e sotto quello che ci avanzava dal quale dovevamo derivare la parte superiore (chiamata predizione). Le regole di questo procedimento sono:
- Esaminiamo sempre il simbolo più a sinistra della predizione
- Se è una variale, dobbiamo fare una predizione: va sostituito con uno dei suoi lati destri
- Se è un terminale, deve coincidere con la prima lettera dell'input da analizzare (simbolo di **look-ahead**). In caso affermativo, viene spostato nel lato 'analisi'; in caso negativo, l'esecuzione fallisce.

<u>Il fallimento non significa che non possiamo fare il parsing della parola, ma solo che la predizione che abbiamo fatto è sbagliata e che quindi dobbiamo tornare indietro e farne un'altra.</u>

*Osservazione:* la predizione ha una gestione LIFO

**Gestione del non-deterministico:** poiché la procedura è, in generale, non deterministica, è necessario esplorare tutte le possibili computazioni. Ci sono 2 possibilità:
1. *In ampiezza:* si esaminano parallelamente tutti i possibili sviluppi della computazione. Richiede molta memoria, per questo a volte è preferibile l'esplorazione in profondità.
2. *In profondità:* si esamina una computazione alla volta, salvo tornare indietro e provare un'altra strada in caso di fallimento (**backtracking**)

**Ampiezza:** dovremo tenere traccia:
- Della lettera che stiamo esaminando (simbolo di **look-ahead**)
- Delle produzioni utilizzate in ogni tentativo di parsing
- Delle predizioni relative a ciascun tentativo

> [!example] Esempio esplorazione in ampiezza:
>  ![[Pasted image 20260508100615.png|300]]

> [!example] Esempio esplorazione in profondità:
> ![[Pasted image 20260508100755.png|300]]

**Recursione destra:** 

> [!example]
> $S' → S$#, $S → Sb | a$ 
> 
> ![[Pasted image 20260501130947.png|300]]

*Osservazione:* <u>il metodo fallisce</u> se $X \implies^* X\alpha$ per qualche $X \in N, \alpha \in V$ (*recursione destra* = da una variabile posso derivare una forma sentenziale che comincia con la variabile stessa). Questo comporta che nella parte destra della nostra tabella non ci saranno mai solo terminali come prima lettera, e di conseguenza dovremo ripetere all'infinito gli stessi passaggi.

*Proposizione:* ogni grammatica non contestuale è equivalente a una grammatica priva di recursione destra. Per esempio, le grammatiche in forma normale di Greibach sono prive di recursione destra

**Backtracking (esplorazione in profondità):**

*Osservazione:* La ricerca in ampiezza richiede molta memoria. In qualche caso può essere preferibile la ricerca in profondità.

> [!example]
> $S' → S$#, $A → a|aA, C → c|cC, S → AB|DC, B → bc | bBc, D → ab | aDb$
> 
> ![[Pasted image 20260501132043.png|300]]
> 
> (questa è la fase finale, ricontrolla passaggi)

**Come realizzare un parser:**  
- Usare dei programmi che fanno questo lavoro
- Costruirci da soli il parser

Nella seconda opzione, la procedura più semplice da seguire è quella di realizzare un **parser a discesa ricorsiva**. Ovvero scrivere tante procedure quante sono le variabili della nostra grammatica. Interpretiamo la produzione $S → aB | bA$ come:
- $S$ ha successo se $a$ ha successo e poi $B$ ha successo
- oppure $b$ ha successo e poi $A$ ha successo 

È ideale per la programmazione di un parser per una grammatica specifica, in quanto si ottiene una procedura per ogni variabile. Dobbiamo tenere traccia:
- posizione corrente nella produzione (automatico), 
- l'input (globale), 
- posizione corrente nell'input (globale), 
- posizione nell'input al momento della chiamata (locale).

> [!example]
> ![[Pasted image 20260508101149.png|300]]

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260501132230.png" width="300">
</div>

### Automi a pila

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260507112643.png" width="300">
</div>

Nel capitolo precedente abbiamo visto come si fa il parsing top-down, dove costruiamo l'albero di parsing partendo dalla radice e producendo le foglie. In sostanza la struttura dati che ci serve è una sequenza di pile (stack) che contengono la cosiddetta predizione in cui leggiamo il simbolo che sta più a sinistra ed eventualmente lo andiamo a sostituire inserendo al posto della prima variabile, la predizione. 

Un **automa a pila** è un modello astratto di macchina (un automa a stati finiti con uno stack) che ci permette di definire i linguaggi non contestuali, esattamente come i linguaggi che sono accettati da questo nuovo modello di macchina astratta. È costituito da un'unità di controllo e un nastro di input (come l'automa a stati finiti) con in più uno stack. Ad ogni passo la macchina legge una lettera dal nastro di input e preleva (*pop*) il simbolo che sta in cima allo stack (se è vuoto, la macchina si blocca). Dopodiché, in base allo stato e ai due simboli letti, la macchina: assume un nuovo stato, scrive (*push*) in cima allo stack 0 o più lettere, e poi sposta la testina di lettura dal nastro di input sulla cella a destra. Sono ammesse delle $\epsilon-$transizioni in cui le operazioni di push e pop sono eseguite senza leggere il nastro di input. Quindi estrae e legge dalla pila, ma non tocca il nastro di input.

*Definizione:* un automa a pila è una settupla $A = \langle Q, \Sigma, \Gamma, \delta, q_0, Z_0, F \rangle$ dove:
- $Q$ è un insieme finito, detto insieme degli stati
- $\Sigma$ è un alfabeto, detto alfabeto di input
- $\Gamma$ è un alfabeto, detto alfabeto di pila
- $\delta: Q \times (\Sigma \cup \{\epsilon \}) \times \Gamma → \wp_F(Q \times \Gamma^*)$ è la funzione di transizione
- $q_0 \in Q$ è lo stato iniziale
- $Z_0 \in \Gamma$ è il simbolo iniziale della pila
- $F \subseteq Q$ è l'insieme degli stati finali 

*Funzionamento:* Inizialmente, l'automa è nello stato iniziale e la pila contiene il simbolo iniziale della pila. Supponiamo che a un dato istante della computazione l'automa si trovi nello stato $p$, riceva $a$ in input e $Z$ dall'operazione di pop. Se $\delta(p, a, Z)$ contiene una coppia $(q, \gamma)$, allora l'automa si può portare nello stato $q$ eseguendo il push di $\gamma$ sulla pila ($\gamma$ è quello che scrivo nella pila al posto della $Z$ che è stata estratta).

*Osservazione:* Modello non deterministico

> [!example]
> ![[Pasted image 20260507123330.png|300]]

*Definizione:* una configurazione interna dell'automa a pila $A$ è una coppia $(q,\gamma) \in Q \times \Gamma^*$. Istante per istante, la configurazione interna = stato + contenuto della pila

Se leggendo una lettera $a$ sul nastro di input l'automa può passare dalla configurazione $(q, \alpha)$ alla configurazione $(q ' , \alpha' )$, allora scriveremo $(q, \alpha) \overset{a}{\models} (q ' , \alpha' )$. Similmente, se l'automa può passare dalla configurazione $(q, \alpha)$ alla configurazione $(q ' , \alpha' )$ con una $\epsilon-$transizione, allora scriveremo $(q, \alpha) \overset{\epsilon}{\models}(q ' , \alpha' )$

*Definizione:* siano $(q,\alpha)$ e $(q',\alpha')$ due configurazioni interne e $a \in \Sigma \cup \{\epsilon\}$. Si ha $(q, \alpha) \overset{a}{\models} (q ' , \alpha' )$, se $\alpha = X\gamma, \alpha' = \beta\gamma$, $q \xrightarrow[]{a,\frac{X}{\beta}} q'$. Se $w = a_1a_2 ... a_n (a_1, ... , a_n \ \in \Sigma \cup {\epsilon})$ e $(q,\alpha) \overset{\Rightarrow{a_1}} (q_1, \alpha1) \overset{\Rightarrow{a_2}} ... \overset{\Rightarrow{a_n}} (qn , \alpha_n ) = (q ' , \alpha' )$ allora scriveremo sinteticamente $(q, \alpha) \overset{w}{\models} (q ' , \alpha' )$

*Osservazione:* Si ha $(q, \alpha) \overset{w}{\models} (q ' , \alpha' )$ se esiste una computazione che, leggendo $w$ sul nastro di input, porta l'automa dalla configurazione interna $(q, \alpha)$ alla configurazione interna $(q ' , \alpha' )$

*Definizione:* sia $A = \langle Q, \Sigma, \Gamma, \delta, q_0, Z_0, F \rangle$ un automa a pila e $w \in  \Sigma^*$. Diremo che la parola $w \in  \Sigma^*$ è accettata per lo stato finale da $A$ se si ha $(q_0, Z_0) \overset{w}{\models} (f,\gamma), f \in F, \gamma \in \Gamma^*$. Diremo che $w$ è accettata per pila vuota da $A$ se si ha $(q_0, Z_0) \overset{w}{\models} (p, \epsilon), p \in Q$. Infine diremo che $w$ è accettata per stato finale e pila vuota da $A$ se si ha $(q_0, Z_0) \overset{w}{\models} (f, \epsilon), f \in F$ 

> [!info] In altre parole:
> - Una parola viene **accettata per stato finale** se c'è una computazione che all'inizio parte con la mia parola sul nastro, stato iniziale, e simbolo iniziale della pila (soltanto nella prima); mano mano legge tutta la mia parola (consuma tutto il nastro di input) e alla fine si trova in uno stato finale che appartiene a $F$. Quindi c'è una computazione che mi parte da $(q_0, w, Z_0)$ e arriva allo stato finale, $\epsilon$ (nastro di input vuoto), $\gamma$
> - Una parola è **accettata per pila vuota** se c'è una computazione che parte con la solita configurazione iniziale $(q_0, w, Z_0)$ e termina in uno stato non-finale perché sia il nastro di input che la pila sono vuote. 
> - Una parola è **accettata per stato finale e pila vuota**, se si verificano entrambe le condizioni viste in precedenza
> 
> Ciò sta ad indicare che ogni automa riconosce 3 linguaggi che potrebbero essere diversi

**Metodi di accettazione:** l'insieme delle parole accettate per lo stato finale (risp., per pila vuota, per stato finale e pila vuota) sarà chiamato il linguaggio accettato da $A$ per stato finale (risp., per pila vuota, per stato finale e pila vuota). I linguaggi accettati da $A$ per stato finale, per pila vuota, per stato finale e pila vuota sono denotati rispettivamente con $L_F(A),L_P(A),L(A)$

> [!example]
> ![[Pasted image 20260507113340.png|300]]

> [!example]
> Problema: Realizzare un automa a pila che accetti per pila vuota le palindrome pari sull'alfabeto $\Sigma = \{a, b\}$.
>
> Strategia:
> 1. copiare nella pila la prima metà dell'input
> 2. verificare se ciò che rimane dell'input coincide con la pila letta al contrario
>
> ![[Pasted image 20260507113443.png|300]] ![[Pasted image 20260507113512.png|300]]
> 
> Il non-determinismo ci consente di non specificare in quale punto della computazione si passa dalla fase 1 alla fase 2.

> [!example] 
> Problema: Realizzare un automa a pila che accetti per pila vuota le palindrome pari sull'alfabeto $\Sigma = \{a, b\}$. Si può ottenere il medesimo risultato con un solo stato 
> 
> Strategia: 
> 1. copiare nella pila la prima metà dell'input 'sotto il simbolo iniziale $Z_0$', 
> 2. cancellare il simbolo iniziale Z 0 dalla cima della pila, 
> 3. verificare che ciò che rimane dell'input è uguale a ciò che abbiamo nella pila, letto in ordine inverso. 
> 
> ![[Pasted image 20260507223911.png|300]]

> [!example]
> Problema: Realizzare un automa a pila che accetti per pila vuota le palindrome sull'alfabeto $\Sigma = \{a, b, c\}$ che contengono una sola occorrenza della lettera $c$. 
> 
> Strategia:
> 1. copiare l'input nella pila finché non si legge la $c$, 
> 2. verificare che ciò che rimane dell'input è uguale a ciò che abbiamo nella pila, letto in ordine inverso. 
> 
> ![[Pasted image 20260507224126.png|300]] ![[Pasted image 20260507224205.png|300]] ![[Pasted image 20260507224222.png|300]]

**Automi a pila deterministici:** 

*Definizione:* Un automa a pila $A = \langle Q, \Sigma, \Gamma, \delta, q_0, Z_0, F \rangle$ è deterministico se per ogni $q \in Q, a \in \Sigma, Z \in \Gamma$, l'insieme $\delta(q, a, Z) \cup \delta(q, \epsilon, Z)$ contiene al più un elemento (cioè <u>quando per ogni combinazione di stato corrente, simbolo in input e simbolo in cima alla pila, le transizioni possibili sono uniche e non conflittuali</u>). Un linguaggio accettato per stato finale da un automa a pila deterministico si dirà deterministico. 

> [!example]
> Il linguaggio delle palindrome sull'alfabeto $\Sigma = \{a, b, c\}$ che contengono una sola occorrenza della lettera $c$ è deterministico. Il linguaggio delle palindrome sull'alfabeto $\Sigma = \{a, b\}$ è non-deterministico. 

*Osservazione:* Se $A$ è deterministico, per ogni configurazione interna $(q, \alpha)$ e ogni lettera $a$ dell'alfabeto di input esiste al più una configurazione interna $(q' , \alpha' )$ per cui $(q, \alpha) \overset{a}{\models} (q ' , \alpha' )$ oppure $(q, \alpha) \overset{\epsilon}{\models} (q ' , \alpha' )$. In tal caso, si può verificare in tempo lineare se una parola è accettata o meno dall'automa.  

**Metodi di accettazione:** 

*Proposizione:* Le tre seguenti classi di linguaggi coincidono: 
1. linguaggi accettati da automi a pila per pila vuota, 
2. linguaggi accettati da automi a pila per stato finale, 
3. linguaggi accettati da automi a pila per pila vuota e stato finale. 

**Pila vuota e stato finale:** 
*Lemma:* Sia $L$ il linguaggio accettato da un automa a pila per pila vuota e stato finale. Esiste effettivamente un automa a pila $A$ tale che $L = L(A) = L_P (A) = L_F (A)$.  

*Dimostrazione:* Sia $A_0$ l'automa a pila che accetta $L$ per pila vuota e stato finale. Costruiamo un automa a pila $A$ con un unico stato finale $f$ che si comporta nel modo seguente: 
- come primo passo, scrive un simbolo speciale # all'inizio della pila, 
- poi simula il funzionamento dell'automa $A_0$, 
- quando si trova in uno stato finale di $A_0$ con # sulla pila, una $\epsilon-$transizione permette di estrarre # dalla pila, in modo da svuotarla, e passare nel suo unico stato finale. 

In una computazione di $A$ la pila si svuota se e solo se raggiunge lo stato finale. Quindi $$L = L(A) = L_P (A) = L_F (A)$$
Inoltre ciò avviene se e solo se l'automa $A_0$ col medesimo input raggiunge uno stato finale con pila vuota. Pertanto $$L = L(A) = L_P (A) = L_F (A)$$

> [!example]
> ![[Pasted image 20260507225456.png|300]]

**Pila vuota:** 

*Lemma:* Sia $L$ il linguaggio accettato da un automa a pila per pila vuota. Esiste effettivamente un automa a pila $A$ tale che $$L = L(A) = L_P (A) = L_F (A)$$
*Dimostrazione:* Sia $A_0$ l'automa a pila che accetta $L$ per pila vuota. Posso supporre, senza perdita di generalità, che tutti gli stati siano finali. Allora $A_0$ accetta $L$ per pila vuota e stato finale. La conclusione segue dal lemma precedente.  

**Stato finale:** 

*Lemma:* Sia $L$ il linguaggio accettato da un automa a pila per stato finale. Esiste effettivamente un automa a pila $A$ tale che$$L = L(A) = L_P (A) = L_F (A)$$
*Dimostrazione:* Sia $A_0$ l'automa a pila che accetta $L$ per stato finale. Costruiamo un automa a pila $A$, aggiungendo ad $A_0$ lo stato e le istruzioni necessarie a svuotare la pila quando raggiunge uno stato finale: 
- un ulteriore stato finale $f$ , 
- $\epsilon-$transizioni da ogni stato finale di $A_0$ in $f$ , per qualunque simbolo estratto dalla pila, senza scrivere nulla nella pila, 
- $\epsilon-$transizioni da $f$ in $f$ , per qualunque simbolo estratto dalla pila, senza scrivere nulla nella pila, 

È facile verificare che $L = L(A)$. La conclusione segue da uno dei lemmi precedenti.  

> [!example]
> ![[Pasted image 20260507225916.png|300]]

**Il caso deterministico:** 

*Proposizione:* Le due seguenti classi di linguaggi coincidono: 
1. linguaggi accettati da automi a pila deterministici per pila vuota, 
2. linguaggi accettati da automi a pila deterministici per pila vuota e stato finale. 

Tale classe è inclusa strettamente nella classe dei linguaggi deterministici. 

> [!example]
> Il linguaggio $\{a n b m | n \geq m \geq 0\}$ è accettato per stato finale dall'automa a pila deterministico 
> 
> ![[Pasted image 20260507230112.png|300]]
> 
> ma non è accettato per pila vuota da nessun automa a pila deterministico

### Teorema di caratterizzazione dei linguaggi non contestuali
Questo teorema ha la funzione di far vedere che i linguaggi accettati dagli automi finiti, quel modello di macchina astratta introdotta nel capitolo precedente, coincidono con la classe dei linguaggi non contestuali (di tipo 2), generati da una grammatica non contestuale. Come per i linguaggi regolari avevamo l'automa a stati finiti che definiva il modello di computazione per verificare se una parola apparteneva o meno al linguaggio, così il modello di macchina per i linguaggi non contestuali è l'automa a pila.

*Teorema:* un linguaggio è non contestuale se e solo se è riconosciuto da un automa a pila.

*Proposizione:* sua $G$ una grammatica non contestuale. Esiste effettivamente un automa a pila con un solo stato che accetta per pila vuota il linguaggio generato da $G$

*Proposizione:* sia $A$ un automa a pila con un solo stato. Esiste effettivamente una grammatica non contestuale $G$ che genera il linguaggio accettato per pila vuota da $A$

*Proposizione:* sia $A$ un automa a pila. Esiste effettivamente automa a pila con un solo stato che accetta per pila vuota il medesimo linguaggio di $A$

**Dalla grammatica all'automa:** 

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260508131000.png" width="300">
</div>

Lo sappiamo già fare perché abbiamo già visto il parsing predittivo che non è altro che un automa a pila

*Procedura non deterministica:*
- Esaminiamo sempre il simbolo più a sinistra della predizione (in cima allo stack)
- Se è una variabile, dobbiamo fare una produzione: va sostituito con uno dei suoi lati destri (simile ad una $\epsilon-$transizione di un automa a pila)
- Se il simbolo a sinistra è un terminale deve coincidere con la prima lettera dell'input da analizzare (simbolo di *look-ahead*). In caso affermativo, viene spostato nel lato 'analisi'; in caso negativo, l'esecuzione fallisce

*Osservazione:* la predizione ha una gestione LIFO. Quindi, salvo la gestione dell'output, ho un automa a pila con un solo stato

**Costruzione dell'automa:** 

*Proposizione:* ogni linguaggio non contestuale p accettato per pila vuota da un automa pila con un solo stato

*Dimostrazione:* sia $L = L(G)$ con $G = \langle V, \Sigma, P, S \rangle$. Costruiamo $A$ con un solo stato e
- Alfabeto di input $\Sigma$
- Alfabeto di pila $V$
- Simbolo iniziale della pila $S$
- Transizioni 

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260508131231.png" width="300">
</div>

Verifichiamo che $L(G) = L_P(A)$ 

*Verifica:* 
<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260508131335.png" width="300">
</div>

Verifichiamo che $L_P(A) \subseteq L(G)$. Se $w \in L_P(A)$, allora avremo $$w = a_1a_2…a_n, S = \alpha_0 \overset{a_1}{\models} \alpha_1 \overset{a_2}{\models} … \overset{a_n}{\models} \alpha_n = \epsilon, a_1,…,a_n \in \Sigma \cup \{\epsilon\}$$
Dalla definizione di $A$, o $a_i = \epsilon$ e $\alpha_{i-1} \implies \alpha_i$ o $a_i \in \Sigma$ e $\alpha_{i-1} = a_i\alpha_i$. In ambo i casi $\alpha_{i-1} \implies^* a_i\alpha_i, i = 1,…,n$. Pertanto, $$S = \alpha_0 \overset{*}{\Rightarrow} a_1\alpha_1 \overset{*}{\Rightarrow} a_1 a_2\alpha_2 \overset{*}{\Rightarrow} … \overset{*}{\Rightarrow} a_1a_2...a_n\alpha_n = w$$ 
e quindi $w \in L(G)$  

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260508131656.png" width="300">
</div>

Verifichiamo che $L(G) \subseteq L_P(A)$. Ricordiamo che ogni parola di $L(G)$ si ottiene dal simbolo iniziale con una derivazione sinistra (cioè sostituendo sempre la variabile più a sinistra). Quindi se $w \in L(G)$, avremo $S \implies u_1X_1\alpha_1 \underset{L}{\Rightarrow} u_1u_2X_2\alpha_2 \underset{L}{\Rightarrow} … \underset{L}{\Rightarrow} u_1u_2…u_n = w, u_1, …, u_n \in \Sigma^*, X_1,…,X_{n-1} \in N$ e $\alpha_1,…,\alpha_{n-1} \in N^*$. Pertanto $X_{i-1}\alpha_{i-1} \implies_L u_iX_i\alpha_i$ e, per la definizione di $A$, $X_{i-1}\alpha_{i-1} \overset{\epsilon}{\models} u_iX_i\alpha_i$. Inoltre $u_{i-1}X_{i}\alpha_{i} \overset{u_i}{\models} X_i\alpha_i$. Quindi, $S \overset{\epsilon}{\models} u_1X_1 \overset{u_1}{\models} X_1\alpha_1 \overset{\epsilon}{\models} u_2X_2\alpha_2 \overset{\epsilon}{\models} … \overset{u_n}{\models} \epsilon$, sicché $w = u_1u_2…u_n \in L_P(A)$ 

> [!example]
> Sia $L = {a^nb^n | n \geq 0}$. $L$ è generato dalle produzioni $$S → aSb | \epsilon$$ L'automa a pila corrispondente avrà alfabeto di input $\{a, b\}$, alfabeto di pila $\{S, a, b\}$, simbolo iniziale della pila $S$ e la funzione di transizione
> 
> ![[Pasted image 20260508133302.png|300]]
> 
> Le $\epsilon-$produzioni sostituiscono ad ogni variabile in cima alla pila i lati destri delle produzioni. Per ognuna delle produzioni $A→B→b$ mettiamo una $\epsilon$ all'incrocio tra il lato sinistro e destro della produzione

**Dall'automa alla grammatica:** 

*Proposizione:* Sia $A$ un automa a pila con un solo stato. Esiste effettivamente una grammatica non contestuale $G$ che genera il linguaggio accettato per pila vuota da $A$. 

*Dimostrazione:* Sia $A = \langle \{q 0\}, \Sigma, \Gamma, \delta, q_0, Z_0, \{q_0\}\rangle$. Possiamo supporre $\Gamma \cap \Sigma = \emptyset$. Costruiamo $G = \langle V , \Sigma, P, Z_0\rangle$ con 
- vocabolario $V = \Sigma \cup \Gamma$, 
- alfabeto dei simboli terminali $\Sigma$, 
- simbolo iniziale $Z_0$, 
- la produzione $X → a\alpha$, per ogni transizione $q_0 \xrightarrow []{a,\frac{X}{\alpha}} q_0$. 

Verificheremo che $L(G) = L_P(A)$. 

*Verifica:* 
Sostituiamo ogni transizione $q_0 \xrightarrow[]{a,\frac{X}{\alpha}}q_0$ di $A$ ($a \in \Sigma, X \in \Gamma, \alpha \in \Gamma^*$) con la coppia di transizioni $$q_0 \xrightarrow[]{\epsilon,\frac{X}{a\alpha}}q_0 , q_0 \xrightarrow[]{a,\frac{a}{\epsilon}} q_0$$
Poiché $\Gamma \cap \Sigma = \emptyset$, la prima transizione deve essere necessariamente seguita dalla seconda. Pertanto, il linguaggio accettato dall'automa non cambia. Si verifica facilmente che, in base alla dimostrazione della proposizione precedente, questo nuovo automa è proprio l'automa a pila che accetta $L(G)$.  

> [!example]
> Sia $A$ l'automa a pila con un unico stato, alfabeto di input $\Sigma = \{a, b, \bar{a}, \bar{b}\}$, alfabeto di pila $\{S, A, B\}$, simbolo iniziale della pila $S$ e la funzione di transizione
> 
> ![[Pasted image 20260508134812.png|300]]
> 
> La grammatica corrispondente avrà le produzioni $S → aA | bB, A → aAA | bBA | \bar{a}, B → aAB | bBB | \bar{b}$. Se $S$ legge 'a', scrive $A$; se $S$ legge 'a barrato', si limita a fare il pop senza il push

**Eliminazione degli stati:** 

*Proposizione:* Sia $A$ un automa a pila. Esiste effettivamente un automa a pila con un solo stato che accetta per pila vuota il medesimo linguaggio di $A$. 

Idea della dimostrazione: Dato un automa a pila $A$, realizzare un automa a pila con un solo stato che simula $A$, registrando nella pila anche lo stato di $A$. 

Problema: Se $A$ non fa push, nella simulazione non posso registrare lo stato. 

Soluzione: Nell'automa a uno stato ogni cella della pila contiene 
- il simbolo di pila di $A$, 
- uno stato arbitrario (cioè tutti, per il non-determinismo), 
- lo stato arbitrario contenuto nella cella sottostante, 
- solo la cella top, al posto dello stato arbitrario contiene lo stato dell'automa simulato.  

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260508135350.png" width="300">
</div>

**Costruzione dell'automa:** Dato $A = \langle Q, \Sigma, \Gamma, \delta, q_0, Z_0, F \rangle$, definiamo $A'$ con: 
- un solo stato, 
- alfabeto di pila $\Gamma' = Q \times \Gamma \times Q$, 
- simbolo iniziale della pila $Z_0' = (q_0, Z_0, q_0)$, 
- la funzione di transizione è definita come segue: Per ogni transizione $p q_0 \xrightarrow[]{a,\frac{X}{X_n...X_1}}q$ in $A$, in $A'$ ci sono le transizioni $q_0 \xrightarrow[]{a, \frac{(p,X,s_1)}{(q,X_n,s_n)(s_n,X_{n-1},s_{n-1})...(S_2,X_1,s_1)}}q_0, a \in \Sigma \cup {\epsilon}, X , X_1, . . . , X_n \in N, p, q, s_1, . . . , s_n \in Q$, Per ogni transizione $p\xrightarrow[]{a,\frac{X}{\epsilon}}q$ in $A$, in $A'$ c'è la transizione $$q_0 \xrightarrow[]{a,\frac{(p,X,q)}{\epsilon}}q_0, a \in \Sigma \cup \{\epsilon\}, p,q \in Q$$

Si verifica che $L_P (A) = L_P (A' )$.

> [!info] N.B.
> Tutto ciò funziona per la classe di automi a pila con un solo stato. Perché esiste effettivamente un automa a pila con un solo stato, che accetta per pila vuota lo stesso linguaggio accettato da un automa qualunque con un qualsiasi numero di stati? Se ho la pila, lo stato lo scrivo sulla pila: l'automa ad ogni passo guarda lo stato, guarda il simbolo in cima alla pila e fa delle cose, e scrive lo stato nuovo sulla pila. Tuttavia, se l'automa non scrive nulla sulla pila, dove scrivo lo stato? Sulla pila non andrò a scrivere solo lo stato, ma scriverò anche cosa c'è scritto sotto. 

### Parsing top-down deterministico

> [!example]
> $S' → S$#, $S → aB|bA, A → a|aS|bAA, B → b|bS|aBB$. Vogliamo fare fil parsing di $aabb$#
> 
> ![[Pasted image 20260513154312.png|300]]

*Osservazione:* possiamo scartare le produzioni che iniziano con un terminale diverso dal simbolo di look-ahead

**La funzione FIRST:**

*Definizione:* sia $G = \langle V, \Sigma, P,S \rangle$ una grammatica non contestuale priva di $\epsilon-$produzioni, la funzione FIRST: $V^+ → \wp(\Sigma)$ è definita da $$ FIRST(\alpha) = \{a \in \Sigma | \alpha \overset{*}{\Rightarrow} a\beta, \beta \in V^*\}, \alpha \in V^*$$
In altri termini $FIRST(\alpha)$ è l'insieme dei terminali che possono comparire come prima lettera di una conseguenza di $\alpha$ 

*Problema:* cerchiamo una derivazione sinistra di $w$: $S$# = $\alpha_0 \underset{L}{\Rightarrow} \alpha_1 \underset{L}{\Rightarrow} \alpha_2 \underset{L}{\Rightarrow} … \underset{L}{\Rightarrow} \alpha_n = w$#. Fissato $i = 0,… n-1$, si avrà:$\alpha_i = uX\beta, \alpha_{i+1} = u\gamma\beta, X→ \gamma$ in $P$ con $u \in \Sigma^*, X \in N$ e $\beta \in V^*$, e anche $w$# = $uv, \gamma\beta \underset{L}{\overset{*}{\Rightarrow}} v, v \in (\Sigma \cup \{$#$\})^*$. Sia $a$ la prima lettera di $v$, allora si deve avere $a \in FIRST(\gamma)$ 

*Conclusione:* nel parsing dobbiamo scegliere sempre produzioni $X → \gamma$ tali che $FIRST(\gamma)$ contenga il simbolo di look-ahead

**Tabella di parsing:** Per ottimizzare il parsing, conviene pre-calcolare $FIRST(\alpha)$ per tutti i lati destri delle produzioni e costruire una tabella in cui a ogni coppia $(X , a) \in N \times \Sigma$ corrispondono le produzioni $X → \gamma$ tali che $a \in FIRST(\gamma)$. Una tale tabella prende il nome di **Tabella di parsing**. 

*Definizione:* Una grammatica non-contestuale $G$ priva di $\epsilon-$produzioni si dice di classe **LL(1)** se per ogni coppia di produzioni distinte $X → \alpha | \gamma$ si ha $FIRST(\alpha) \cap FIRST(\gamma) = \emptyset$. 

In altri termini, una grammatica LL(1) è una grammatica la cui tabella di parsing contiene al più una produzione in ogni ingresso. Questo assicura che il parsing può essere effettuato in maniera deterministica. 

*Osservazione:* Il termine LL(1) viene da: derivazione Left-most, input esaminato Left-to-right, 1 simbolo di look-ahead.

**Parsing LL(1)** ho un automa a pila deterministico che prende in input la tabella di parsing, la tabella di una grammatica non contestuale e una parola $w = a_1a_2…a_n$#. In output restituisce le produzioni di una derivazione sinistra di $w$

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260513155934.png" width="300">
</div>

**Generazione della tabella di parsing:**
1. se la prima lettera di $\alpha$ è un simbolo terminale $b$, allora $FIRST(\alpha) = \{b\}$; 
2. se la prima lettera di $\alpha$ è una variabile $X$ , allora per ogni produzione $X → \beta$, i terminali di $FIRST(\beta)$ sono anche terminali di $FIRST(\alpha)$, 
3. tutti gli elementi di $FIRST(\alpha)$ con $\alpha$ lato destro di qualche produzione, si trovano applicando un numero finito di volte le regole precedenti

*Calcolo di $FIRST(\alpha)$:* 
1. poniamo $FIRST(\alpha) = \emptyset$ per tutti i lati destri $\alpha$ delle produzioni di $G$. 
2. ricalcoliamo tali insiemi usando le regole 1–2 e i valori attuali degli insiemi $FIRST(\beta)$. 
3. iteriamo il passo precedente, arrestandoci solo quando nessuno degli insiemi $FIRST(\alpha)$ subisce modifiche.

> [!example]
> Si consideri la grammatica con le produzioni $S → FS|Q|aSbS, F → c, Q → q$
> 
> ![[Pasted image 20260513160409.png|300]] ![[Pasted image 20260513160426.png|300]]
> 
> ![[Pasted image 20260513160557.png|300]]
> 
> (immagine solo della situazione finale)

**$\epsilon-$produzioni:**

*Definizione:* sia $G = \langle V, \Sigma, P,S \rangle$ una grammatica non contestuale. Le funzioni $FIRST: V^* → \wp(\Sigma \cup \{\epsilon\})$ e $FOLLOW: N → \wp(\Sigma)$ sono definite da:$$FIRST(\alpha) = \{\alpha \in \Sigma | \alpha \overset{*}\Rightarrow a\beta, \beta \in V^*\}, \alpha \in V^*, \alpha \not{\overset{*}\Rightarrow} \epsilon$$ $$ FIRST(\alpha) = \{\alpha \in \Sigma | \alpha \overset{*}\Rightarrow a\beta, \beta \in V^*\} \cup \{\epsilon\}, \alpha \in V^*, \alpha \overset{*}\Rightarrow \epsilon$$ $$
 FOLLOW(X) = \{\alpha \in \Sigma | S\# \overset{*}\Rightarrow \beta Xa\gamma, \beta, \gamma \in V^*\}, X \in N$$
*Problema:* Cerchiamo una derivazione sinistra di $w$: $$S\# = \alpha_0 \underset{L}\Rightarrow \alpha_1 \underset{L}\Rightarrow  \alpha2 \underset{L}\Rightarrow · · · \underset{L}\Rightarrow \alpha_n = w\#$$
Fissato $i = 0, . . . , n − 1$, si avrà $\alpha_i = uX \beta , \alpha_{i+1} = u\gamma\beta , X → \gamma$ in $P$ , con $u \in \Sigma^* , X \in N$ e $\beta \in V^*$, e anche $w\# = uv , \gamma\beta \overset{*}{\underset{L}\Rightarrow} v , v \in (\Sigma \cup {\#})^*$ . Sia $a$ la prima lettera di $v$. Allora si deve avere $a \in FIRST(\gamma)$ oppure $\gamma \Rightarrow^* \epsilon$ e $\beta \Rightarrow^* v$. In quest'ultimo caso, $\epsilon \in FIRST(\gamma)$ e $a \in FOLLOW(X )$. 

*Conclusione:* Nel parsing, dobbiamo scegliere sempre produzioni $X → \gamma$ tali che il simbolo di look-ahead a soddisfi la condizione $a \in FIRST(\gamma)$ oppure $\epsilon \in FIRST(\gamma)$ e $a \in FOLLOW(X )$.

**Tabella di parsing:** In presenza di $\epsilon-$produzioni, nella tabella di parsing dovremo inserire per ogni coppia $(X , a) \in N \times \Sigma$, le produzioni $X → \gamma$ tali che $a \in FIRST(\gamma)$ oppure $\epsilon \in FIRST(\gamma)$ e $a \in FOLLOW(X )$. 

*Definizione:* Una grammatica non-contestuale $G$ si dice di classe LL(1) se per ogni coppia di produzioni distinte $X → \alpha | \gamma$ si ha 
1. $FIRST(\alpha) \cap FIRST(\gamma) = \emptyset$, 
2. se $\epsilon \in FIRST(\alpha)$, allora $\epsilon \notin FIRST(\gamma)$ e F$OLLOW(X ) \cap FIRST(\gamma) = \emptyset$. 

Anche in questo caso, una grammatica è LL(1) se e solo se la sua tabella di parsing contiene al più una produzione in ogni ingresso. 

**Calcolo di FIRST:** 
1. se la prima lettera di $\alpha$ è un simbolo terminale $b$, allora $FIRST(\alpha) = b$; 
2. se la prima lettera di $\alpha$ è una variabile $X$ , allora per ogni produzione $X → \beta$, i terminali di $FIRST(\beta)$ sono anche terminali di $FIRST(\alpha)$; 
3. se la prima lettera di $\alpha$ è una variabile annullabile, allora le due condizioni precedenti sono verificate anche dalla seconda lettera di $\alpha$; se quest'ultima è una variabile annullabile, allora anche dalla terza, e così via; 
4. se tutte le lettere di $\alpha$ sono variabili annullabili, allora $\epsilon \in FIRST(\alpha)$; 
5. Tutti gli elementi di $FIRST(\alpha)$ con $\alpha$ lato destro di qualche produzione, si trovano applicando un numero finito di volte le regole precedenti. 

*Calcolo di $FIRST(\alpha)$:* 
1. Poniamo $FIRST(\alpha) = \emptyset$ per tutti i lati destri $\alpha$ delle produzioni di $G$. 
2. Ricalcoliamo tali insiemi usando le regole 1–4 e i valori attuali degli insiemi $FIRST(\beta)$. 
3. Iteriamo il passo precedente, arrestandoci solo quando nessuno degli insiemi $FIRST(\alpha)$ subisce modifiche. 

> [!example]
> Si consideri la grammatica con le produzioni $S → FQ | aSbS, F → CF | \epsilon, C → c, Q → q$. L'unica variabile annullabile è $F$. Calcolo di $FIRST$ sui lati destri delle produzioni:
> 
> ![[Pasted image 20260514103916.png|300]]

**Calcolo di FOLLOW:** 
1. L'insieme $FOLLOW(S)$ contiene $\#$. 
2. Se c'è una produzione $X → \alpha Y \beta$ con $X , Y \in N, \alpha, \beta \in V^*$ , allora tutti i terminali contenuti in $FIRST(\beta)$ sono contenuti anche in $FOLLOW(Y )$. 
3. Se inoltre $\epsilon \in FIRST(\beta)$, cioè $\beta \Rightarrow^* \epsilon$, allora tutti i terminali contenuti in $FOLLOW(X )$ sono contenuti anche in $FOLLOW(Y )$. 
4. Tutti gli elementi di $FOLLOW(Y ), Y \in N$, si trovano applicando un numero finito di volte le regole precedenti. 

*Calcolo di $FOLLOW(Y )$:* 
1. per tutte le variabili $Y$ , poniamo $FOLLOW(Y )$ uguale a ciò che si ottiene applicando le regole 1. e 2. 
2. aggiungiamo a tali insiemi i terminali ottenuti usando la regola e i valori attuali degli insiemi $FOLLOW(X )$. 
3. iteriamo il passo precedente, arrestandoci solo quando nessuno degli insiemi $FOLLOW(Y )$ subisce modifiche.

> [!example]
> Si consideri la grammatica con le produzioni $S → FQ|aSbS, F → CF|\epsilon, C → c, Q →q$
> 
> ![[Pasted image 20260513161746.png|300]]

### Proprietà di chiusura dei linguaggi non contestuali
**Operazioni regolari:** 

*Proposizione:* La classe dei linguaggi non contestuali è chiusa per le operazioni regolari: unione, concatenazione e chiusura di Kleene. 

> [!example]
> Le grammatiche $G_1$ e $G_2$ che hanno, rispettivamente, le produzioni $S → aSb | ab e S → bSa | ba$, generano i linguaggi $L_1 = \{a^nb^n | n > 0\}$ e $L_2 = \{b^n a^n | n > 0\}$. Consideriamo la grammatica con le produzioni $S → S_1 | S_2, S_1 → aS_1b | ab, S_2 → bS_2a | ba$. Tale grammatica genera il linguaggio $L = L_1 \cup L_2$. Invece la grammatica con le produzioni $S → S_1 S_2, S_1 → aS_1 b | ab, S_2 → bS_2 a | ba$ genera il linguaggio $L = L_1 L_2 = \{a^n b^n+m a m \}$. Infine, la grammatica con le produzioni $S → S_1 S | \epsilon, S_1 → aS_1 b | ab$ genera il linguaggio $L = L_1^* = {\epsilon} \cup \{a^{n_1} b^{n_1} a^{n_2} b^{n^2} · · · a^{n_k} b^{n_k} | k, n_1, . . . , n_k > 0\}$.  

*Dimostrazione:* voglio esplicitare qual è la costruzione per realizzare la grammatica non contestuale che mi genera l'unione, la concatenazione o la chiusura di Kleene, date ,e grammatiche che mi generano i due linguaggi di partenza. Siano $L_1 = L(G_1), L_2 = L(G_2)$ con $G_1 = \langle V_1, \Sigma, P_1, S_1\rangle$ e $G_2 = \langle V_2, \Sigma, P_2, S_2\rangle$. Possiamo supporre $N_1 \cap N_2 = \emptyset$. Costruiamo la grammatica $G = \langle V , \Sigma, P, S\rangle$ con:
- $N = \{S\} \cup N_1 \cup N_2$, 
- $P = P_1 \cup P_2 \cup \{S → S_1 | S_2\}$. 

La grammatica $G$ genera $L_1 \cup L_2$. Costruiamo la grammatica $G = \langle V , \Sigma, P, S\rangle$ con 
- $N = \{S\} \cup N_1 \cup N_2$, 
- $P = P_1 \cup P_2 \cup {S → S_1 S_2}$. 

La grammatica $G$ genera $L_1 L_2$. Costruiamo la grammatica $G = \langle V , \Sigma, P, S\rangle$ con 
- $N = \{S\} \cup N_1$, 
- $P = P_1 \cup {S → S_1 S | \epsilon}$. 

La grammatica $G$ genera $L_1^*$  

> [!example]
> I linguaggi $L_1 = \{a^n b^n | n \geq 0\}, L_2 = \{b^n c^n | n \geq 0\}, a^* , c^*$ sono non contestuali. Per la proposizione precedente, lo saranno anche (otteniamo un linguaggio non contestuale) $$L_3 = L_1 c^* = \{a^n b^n c^m | m, n \geq 0\} \ , \ L_4 = a^*L_2 = \{a^m b^n c^n | m, n \geq 0\}$$ ma non la loro intersezione $$L_3 \cap L_4 = \{a^n b^n c^n | n \geq 0\}$$ *Pertanto:* La famiglia dei linguaggi non contestuali non è chiusa per intersezione né per complemento.  

**Intersezione con linguaggi regolari:** 
Se faccio l’intersezione non è detto che mi venga fuori un linguaggio non- contestuale (infatti quello dell’esempio non è non-contestuale). Però se uno è non contestuale e l’altro è regolare, sono certo che mi venga fuori un linguaggio non contestuale. 

*Proposizione:* L'intersezione di un linguaggio non contestuale con un linguaggio regolare è un linguaggio non contestuale. 

<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Pasted image 20260514112705.png" style="width: 100%; border-radius: 8px;">
  </div>
  <div style="flex: 1.5;">
  Dimostrazione: Un linguaggio regolare R è accettato da un automa a stati finiti A_1, Un linguaggio non contestuale L è accettato da un automa a pila A_2 con un solo stato, L'intersezione R ∩ L è accettata, per stato finale e pila vuota, da un automa a pila che simula A_1 e A_2 in parallelo.
  </div>
</div>

**Linguaggio di Dick:** Il linguaggio di Dick $D_n$ è il linguaggio sull'alfabeto $\Sigma_n = \{a_1, a_2, … , a_n , b_1, b_2,…, b_n \}$, generato dalla grammatica con le produzioni $S → a_iSb_iS , i = 1, 2, . . . n , S → \epsilon$ (non contestuale)

*Equivalentemente:* sequenze di parentesi ben bilanciate 

Accettato dall'automa a pila deterministico 
<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260514113127.png" width="300">
</div>

> [!example]
> Sia $R = a_1^* b_1^*$ . Allora $$D_1 \cap R = \{a^{n_1} b^{n_1} | n \geq 0\}$$ Si ha $D_1 \cap a_1^* b_1^* = \{a^{n_1} b^{n_1} | n \geq 0\}$ e, come sappiamo, questo linguaggio è non contestuale. Sia $R = (a_1 + a_2)^* (b_1 + b_2)^*$ . Allora il linguaggio $$L = D_2 \cap R$$ è l'insieme delle parole $uv$ con $u \in \{a_1, a_2\}^*$ e $v$ ottenuta da $u$ sostituendo le lettere $a_1$ e $a_2$ rispettivamente con $b_1$ e $b_2$ e capovolgendo la parola ottenuta.  

**Morfismi:** 

*Definizione:* Un morfismo $f : \Sigma^* → \Gamma^*$ è una funzione tale che 
- $f (\epsilon) = \{\epsilon\}$: la parola vuota vale la parola vuota
- $f (a_1 · · · a_n ) = f (a_1)· · · f (a_n )$, per ogni $a_1,… , a_n \in \Sigma, n > 0$: una parola che è la concatenazione dei n lettere, l'immagine si trova anche concatenando l'imagine delle lettere

Cioè è una funzione che definita sull'alfabeto $\Sigma$ mi restituisce parole sull'alfabeto $\Gamma$ (sostituisce un simbolo con una stringa)

*Osservazione:* Un morfismo è completamente determinato dalle immagini delle lettere. 

> [!example] 
> Sia $f : \{a, b\}^* → \{a, b\}^*$ il morfismo definito da $$f (a) = a \ , \ f (b) = bb$$ 
> Per es., si ha $f (bab) = bbabb$. Se $L = \{a^n b^n | n \geq 0\}$, allora $f (L) = \{a^n b^{2n} | n \geq 0\}$.  

**Chiusura:** 

*Proposizione:* La classe dei linguaggi non contestuali è chiusa per morfismi. Cioè se prendo un linguaggio non contestuale e un morfismo, faccio l'immagine di tutte le parole del linguaggio non contestuale, ottengo un nuovo linguaggio anche lui non contestuale 

> [!example] 
> Sia $L = \{a^n b^n | n \geq 0\}$ e $f : \{a, b\}^* → \{a, b\}^*$ il morfismo definito da $$f (a) = a, f (b) = bb$$ L è generato dalla grammatica con produzioni $$S → aSb | \epsilon$$ $f (L)$ è generato dalla grammatica con produzioni $$S → ASB | \epsilon, A → a, B → bb$$.  

**Teorema di rappresentazione:** 

> [!example] 
> Sia $D_2$ il linguaggio di Dick sull'alfabeto $\Sigma_2 = \{a_1, a_2, b_1, b_2\}, R = (a_1 + a_2)^* (b_1 + b_2)^*$ e $f : \Sigma^*_2 → \{a, b\}^*$ il morfismo definito da $$f (a_1) = f (b_1) = a, f (a_2) = f (b_2) = b$$ Allora $$f (D_2 \cap R)$$ è l'insieme delle palindrome pari sull'alfabeto $\Gamma$. 

*Teorema (Chomsky-Schützenberger):* Un linguaggio $L$ è non contestuale se e soltanto se esistono un intero $k > 0$, un linguaggio regolare $R$ e un morfismo $f$ tali che $$L = f (D_k \cap R)$$ 
con $D_k$ un linguaggio di Dick di parentesi bilanciate al quale impongo vincoli regolari ($R$) e di cui trasformi i simboli con un morfismo ($f$) 