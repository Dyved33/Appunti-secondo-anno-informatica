###### **Compilatori e Interpreti**
La traduzione del codice sorgente in un linguaggio comprensibile dalla macchina può avvenire in modalità differenti (è anche possibile la traduzione da un linguaggio di programmazione a un altro):
* **Compilatore:** Esegue la traduzione in blocco *prima* dell'esecuzione. Prende in input il **codice sorgente** e, tramite il traduttore, restituisce il **codice oggetto** (oppure delle notifiche di errore).
* **Interprete:** Esegue la traduzione *istruzione per istruzione*, simultaneamente all'esecuzione del programma. Prende in input sia il **codice sorgente** sia l'**input** dell'utente, producendo direttamente l'**output**.

Esistono compilatori Ibridi (Es. Java) che combinano i due approcci:
1. Traduttore analizza il codice sorgente e lo traduce in un **programma intermedio** (ad esempio, il *Java bytecode*).
2. Successivamente, una **macchina virtuale** (ad esempio, la JVM) riceve gli input dell'utente e *interpreta* il codice intermedio per produrre l'output finale.
###### **Le Fasi della Compilazione: Dall'Analisi alla Sintesi**
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
###### **L'Architettura Software del Compilatore**
L'organizzazione interna del compilatore definisce come i moduli comunichino effettivamente tra loro durante il processo di traduzione.
* **Il Nucleo Centrale (Parser & Semantica):** la parte centrale di tutto il processo è l'**analizzatore sintattico**. In questa architettura, esso realizza contestualmente anche l'**analisi semantica statica**.
   > [!IMPORTANT] Avanzamento Parallelo
   > Dato che le due funzioni sono accorpate nel nucleo, l'analisi sintattica e l'analisi semantica avanzano parallelamente.
* **Gestione Input e Lettura:** Il sistema inizia con il modulo **gestore input**, che riceve il codice sorgente e può generare un "listato sorgente". Per formare i token, l'analizzatore lessicale deve leggere il codice un carattere alla volta invocando la funzione `readChar()` direttamente sul gestore input.
* **Dialogo Scanner-Parser:** L'analizzatore sintattico coordina il lavoro chiedendo i dati all'analizzatore lessicale tramite la chiamata alla funzione `prossimoToken()`. Questa funzione, implementata nello scanner, restituisce il `token` richiesto per continuare a costruire l'albero.
* **Gestione degli Errori:** Il **notificatore degli errori** riceve comunicazioni dirette e incrociate dal gestore input, dall'analizzatore lessicale e dal blocco sintattico/semantico ogni volta che viene rilevata un'anomalia. Alla fine, produce in output la "lista degli errori".
* **Generazione del Codice Finale:** Completata l'analisi senza errori bloccanti, è sempre l'analizzatore sintattico ad attivare la sintesi. Tramite chiamate a funzioni opportune, innesca il **generatore del codice**, che a sua volta emetterà il **codice oggetto** definitivo.
###### Il problema dell'infinito e il Teorema di Cantor
Per descrivere qualsiasi cosa, che sia la Divina Commedia o il codice genetico umano, usiamo delle parole. 
In informatica teorica, un **linguaggio è definito come un insieme di parole**.
* **Linguaggi Finiti:** Se un linguaggio è finito, la sua descrizione è banale poiché è possibile elencarne tutti gli elementi uno per uno. Un esempio è $L=\{a,ab,abb,abbb\}$.
* **Linguaggi Infiniti:** Il problema sorge quando cerchiamo di descrivere linguaggi con infinite parole.

Il **Teorema di Cantor** ci pone un limite matematico insuperabile: non esiste una funzione iniettiva dall'insieme delle parti $\mathcal{P}(\Sigma^{*})$ all'insieme delle parole $\Sigma^{*}$. Di conseguenza, non è possibile associare a *ciascun* linguaggio possibile una parola univoca (la sua "descrizione") che lo caratterizzi in modo perfetto. In altre parole, non possiamo inventare una "parola descrittiva" univoca per tutti i possibili linguaggi infiniti esistenti (linguaggio infinito non può essere descritto da una descrizione finita).
> [!NOTE] L'approccio Informatico
> Nonostante il Teorema di Cantor, in informatica non abbiamo bisogno di descrivere *tutti* i linguaggi infiniti. È comunque possibile fornire **descrizioni finite** per *alcuni* linguaggi infiniti che ci sono particolarmente utili. Ad esempio, possiamo usare una regola descrittiva semplice come "tutte le parole sull'alfabeto $\{a,b\}$ di lunghezza dispari", oppure la notazione formale $\{a^{n}b^{n}|n\ge0\}$.

Un costrutto fondamentale per generare infinità è la **chiusura di Kleene** di un linguaggio $L$ (indicata con $L^{*}$). Essa è costituita dalla parola vuota, dalle parole appartenenti a $L$ e da tutte le parole che si ottengono concatenando due o più parole di $L$. Se $L$ è un linguaggio finito, il linguaggio derivato $L^{*}$ è ben definito e, in generale, infinito.
###### Alfabeti , Parole e le loro Proprietà (Definizioni Formali):
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
###### Linguaggi e grammatiche
**Linguaggio Formale (Definizione Rigorosa):** Ogni sottoinsieme di $\Sigma^{*}$ si dice **linguaggio formale** (o semplicemente linguaggio) sull'alfabeto $\Sigma$. Possono esserci linguaggi limitati come $L_{1}=\{a,ab,abb\}$, linguaggi infiniti come $L_{2}=\{ab^{n}a|n\ge0\}$, l'insieme vuoto $L_{3}=\emptyset$ o l'intero universo $L_{4}=\Sigma^{*}$.

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
3.  **Il Motore dell'Infinito (Ricorsione):** $\langle lista Nomi \rangle$ può essere sostituito da $\langle nome \rangle$ oppure da $\langle nome \rangle, \langle lista Nomi \rangle$. 
> [!EXAMPLE] Come si genera l'infinito
    > Contenendo se stessa nella parte destra dell'assegnazione, può essere applicata all'infinito per aggiungere quante virgole e quanti nomi desideriamo, creando una lunghezza illimitata a partire da una singola regola.

4.  **Chiusura:** Il blocco $\langle nome \rangle \langle fine Lista \rangle$ si sostituisce rigorosamente con $e \langle nome \rangle$. Questo trucco garantisce che, alla fine del processo meccanico, l'ultimo elemento sia sempre preceduto dalla 'e' e non dalla virgola. 
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
**Definizione:** Due grammatiche si dicono equivalenti se generano lo stesso linguaggio.  
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
###### Gerarchia di Chomsky e tipi di linguaggio
Divide le grammatiche e i linguaggi generati da esse in 4 tipi (dai più complessi ai più semplici)

**Linguaggi di tipo 0:** 
**Definizione:** le grammatiche a struttura di frase (in base alle loro regolo possono scrivere intere frasi e non solo singoli caratteri) si dicono anche *grammatiche di tipo 0*. I linguaggi generati da grammatiche di tipo 0 si dicono *linguaggi di tipo 0* oppure *linguaggi ricorsivamente numerabili*. Fanno parte di questo livello le grammatiche a struttura di frase. Queste grammatiche hanno la massima espressività possibile: <u>se io ho una procedura effettiva infinita che mi genera una lista di parole, quella lista di parole è un linguaggio di tipo 0.</u>
> [!example] 
> ![[Screenshot From 2026-03-05 10-15-42.png|200]]

**Linguaggi di tipo 1 (grammatiche sensibili al contesto / contestuali):**
**Definizione:** una grammatica a struttura di frase si dice sensibile al contesto se le produzioni hanno la forma: $\alpha_1X\alpha_2 \to \alpha_1\beta\alpha_2$, con $X \in N$, $\alpha_1, \alpha_2, \beta \in V^*$ , $\beta \neq \epsilon$.  Ciò significa che le produzioni di queste grammatiche hanno il lato sx che contiene una variabile, mentre il lato dx si ottiene sostituendo quella variabile con una parola qualsiasi purché non vuota. Quindi mentre prima potevamo sostituire una parola con una parola, qui possiamo sostituire solo una lettera all'interno di una parola con una parola, che per di più non deve essere vuota (posso sostituire una variabile solo se è contenuta in un  contesto). I linguaggi generati da grammatiche di tipo *q* si dicono *linguaggi di tipo 1*  o anche *sensibili al contesto*.  A differenza dei linguaggi di tipo 0, per questo tipo di linguaggi esistono algoritmi che risolvono i problemi di parsing e ricognizione, ma sono molto costosi.
<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Screenshot From 2026-03-05 10-22-40.png" style="width: 100%; border-radius: 8px;">
  </div>
  <div style="flex: 1.5;">
  VN / VU è possibile se e solo so V è seguita da N / U (per questo dipende da contesto)
  </div>
</div>
> [!example]
> ![[Screenshot From 2026-03-05 10-22-40.png|200]]

**Grammatiche monotone (sempre di tipo 1):**
**Definizione:** una grammatica si dice *monotona* se tutte le produzioni hanno la forma: $\alpha \to \beta$ con $|\alpha| \leq |\beta|$ . Cioè una grammatica in cui tutte le produzioni hanno il lato dx che ha lunghezza maggiore o uguale a quella del lato sx. Quindi non c’è più l’obbligo di sostituire singole variabili, ma c’è la sola regola che il lato dx non può essere più corto di quello sx. Inoltre: 
- Le grammatiche sensibili al contesto sono monotone, perché in esse si sostituisce (all'interno del contesto) una variabile con qualcosa che non è vuoto, quindi quello che otterremo sarà qualcosa di più lungo o al massimo della stessa lunghezza di quello che già abbiamo.
- Non tutte le grammatiche monotone sono sensibili al contesto
- Ogni grammatica monotona è equivalente a una grammatica sensibile al contesto (che genera lo stesso linguaggio)
> [!example] ![[Screenshot From 2026-03-05 10-26-52.png|200]]

**Sintesi:**
1. I linguaggi di tipo 1 costituiscono una sottoclasse propria dei linguaggi di tipo 0
2. Esistono algoritmi che risolvono il problema di ricognizione e il problema di parsing per grammatiche di tipo 1 (in generale sono molto costosi)
3. Esistono linguaggi di tipo 0 per cui un tale algoritmo non esiste

**Linguaggi di tipo 2 (Grammatiche non contestuali):** Una grammatica a struttura di frase si dice *non contestuale* o di *tipo 2* se le produzioni hanno la forma: $X \to \beta$, con $X \in N$, $\beta \in V^*$. Dunque una grammatica a struttura di frase si dice non contestuale se le produzioni hanno a sx una singola variabile. I linguaggi generati da grammatiche di tipo 2 si dicono *linguaggi di tipo 2* o anche *non contestuali*. Sono non contestuali perché la variabile può essere sostituita indipendentemente dal contesto in cui appare.
> [!example]
> ![[Screenshot From 2026-03-05 10-33-04.png|200]]

I linguaggi di tipo 2 costituiscono una sottoclasse propria dei linguaggi di tipo 1. Per esempio $L = \{ a^{n}b^{n}c^{n} | n > 0\}$ è un linguaggio di tipo 1 ma non è di tipo 2.
Gli algoritmi di ricognizione e parsing per le grammatiche di tipo 2 si comportano abbastanza bene, alcuni sono moderatamente efficienti . La maggior parte dei linguaggi di programmazione sono generati da una grammatica non contestuale.

**Linguaggi di tipo 3 Grammatiche regolari):** 
**Definizione:** una grammatica a struttura di frase si dice *regolare* o di *tipo 3* se le produzioni hanno la forma $X \to aY$ oppure $X \to a$, con $X,Y \in N$, $a \in \Sigma$. Dunque Le produzioni hanno una forma ancora più semplice: a sx sempre una variabile, come nelle grammatiche di tipo 2, ma a dx ho un terminale da solo o un terminale seguito da una variabile. I linguaggi generati da grammatiche di tipo 3 si dicono *linguaggi di tipo 3* o anche *regolari*. I linguaggi di tipo 3 costituiscono una sottoclasse dei linguaggi di tipo 2 perché le produzioni a sx hanno una singola variabile. Quindi un linguaggio generato da una grammatica di tipo 3, è anche un linguaggio generato da una grammatica di tipo 2. Anche in questo caso, però, ci sono linguaggi di tipo 2 che non sono linguaggi di tipo 3.
> [!example]
> ![[Pasted image 20260305103801.png|200]]

> [!info] Gli algoritmi di ricognizione per i linguaggi di tipo 3 esistono e sono estremamente efficienti, sono molto importanti per l’analisi lessicale ovvero l’analisi di token che costituiscono il nostro codice. Questo si ottiene attraverso la Teoria degli Automi a Stati Finiti.
###### Automi a stati finiti deterministici
Un automa è un dispositivo che può raggiungere un certo numero finito di configurazioni interne. Il nastro di input è diviso in celle su cui si può scrivere l'input (scrivo la stringa in ingresso e la divido nei suoi caratteri) e a ogni passo il dispositivo:
- Legge il contenuto di una cella
- Passa in un nuovo stato (dipende esclusivamente dallo stato precedente e dal contenuto della cella letta)
- Passa a esaminare la cella successiva

Terminata la lettura del nastro, in base allo stato raggiunto (se è uno stato accettante o meno) , accetta o rifiuta la parola letta. Le parole accettate costituiscono un linguaggio di tipo 3 che, a sua volta, è sempre accettato da un automa a stati finiti. Questo algoritmo è lineare e in tempo reale, perché ad ogni cella letta so se tutto ciò che è venuto prima è stato accettato o meno.

**Definizione:** un *automa a stati finiti deterministico* è una quintupla $A = (Q, \Sigma, \delta, q_0, F)$ ove:
- $Q$ è un insieme finito, detto insieme degli stati = tutte le configurazioni che il dispositivo può assumere
- $\Sigma$ è un alfabeto detto alfabeto di input
- $\delta$: $Q \times \Sigma \to Q$ è la funzione di transizione  = legge con cui i nostri stati cambiano in base ad uno stato precedente e alla lettera scritta sul nastro
- $q_0 \in Q$ è lo stato iniziale
- $F \subseteq Q$ è l'insieme degli stati finiti = se al termine della computazione, il dispositivo si trova in uno stato che è compreso nell'insieme degli stati finali, la parola viene accettata
> [!example] ![[Screenshot From 2026-03-05 21-44-54.png|200]]

**Funzione di transizione estesa:** La funzione $\delta$ , affinché non lavori solo con le coppie *stato-lettera* ma anche con le coppie *stato-parola*, si estende a $\hat{\delta}$: $Q \times \Sigma^* \to Q$ ponendo:
- $\hat{\delta}(q, \epsilon) = q$, $\forall q \in Q$ 
- $\hat{\delta}(q,va) = \delta(\hat{\delta}(q,v),a)$, $\forall q \in Q$, $v \in \Sigma^*$, $a \in \Sigma$

Per ogni $q \in Q, u, v \in \Sigma^*$,  si ha $\hat{\delta}(q,uv) = \hat{\delta}(\hat{\delta}(q,u),v)$  
> [!example] ![[Screenshot From 2026-03-05 21-47-18.png|200]]

**Linguaggio accettato:**
**Definizione:** una parola $w \in \Sigma^*$ è accettata da $A$ se $\hat{\delta}(q_0, w) \in F$. L'insieme delle parole accettate da $A$ si dice linguaggio riconosciuto (o accettato) da $A$ e si denota come $L(A) = \{w \in \Sigma^* | \hat{\delta}(q_0,w) \in F\}$ 

**Grafo di un automa:** A ogni automa $A = (Q, \Sigma, \delta, q_0, F)$ associamo un grafo diretto con frecce etichettate:
- I vertici sono gli stati
- Le frecce sono le triple ($q, a, \delta (q, a)$) con $q \in Q , a \in \Sigma$ 
- Lo stato iniziale è denotato da una freccia entrante
- Gli stati finali sono identificati dal bordo doppio
> [!example] ![[Screenshot From 2026-03-05 21-53-02.png|200]]

- Per ogni stato $q$ e ogni lettera $a$ c'è esattamente una freccia con etichetta $a$ che esce da $q$
- C'è un cammino da uno stato $p$ a uno stato $q$ con etichetta $w$ $\iff$ $q = \hat{\delta}(p,w)$
- Le parole accettate dall'automa sono esattamente le etichette dei cammini da $q_0$ a uno stato finale
> [!example] ![[Screenshot From 2026-03-05 21-54-11.png|200]]
> 
> - ![[Screenshot From 2026-03-06 10-06-15.png|200]]
> - ![[Screenshot From 2026-03-06 10-06-27.png|200]]
> - ![[Screenshot From 2026-03-06 10-06-39.png|200]]
###### Automi a stati finiti non deterministici
Abbiamo detto che un automa è rappresentato da un grafo e che per ogni stato $q$ e ogni lettera $a$ esiste esattamente un arco uscente da $q$ con etichetta $a$. Le parole accettate da $A$ sono le etichette dei cammini dallo stato iniziale a uno stato finale.
Se da questa definizione togliessi l'aspetto degli archi (cioè che da ogni vertice può uscire un solo arco) otterrei: una maggiore libertà nella progettazione ma avrei più difficoltà a verificare se una parola è accettata o meno. La soluzione a ciò è un modello non-deterministico che mi permette di fare una conversione automatica al modello deterministico.  Il modello non-deterministico:
- A ogni passo il dispositivo esegue una fra le transizioni possibili (purché esista)
- L'input è accettato se almeno una delle possibili computazioni ha successo (termina in uno stato finale)
- L'input è rifiutato se tutte le possibili computazioni falliscono (terminano in uno stato non finale o prima che il nastro sia esaurito)

**Sintesi:** Possiamo definire Automi a Stati finiti non Deterministici, dei meccanismi che non hanno la restrizione dei modelli deterministici: un meccanismo che può assumere un numero predefinito di configurazioni interne, che ha un nastro di input, ma che ad ogni passo non esegue più una sola transizione ben determinata ma dispone di un ventaglio di possibilità e può quindi eseguire varie transizioni (o nessuna). L’input viene accettato se tra le tante computazioni possibili c’è n’è almeno una che termina in uno stato finale. Quindi una parola viene rifiutata se tutte le computazioni terminano in uno stato non finale o terminano prima che il nastro sia esaurito. 

**Definizione:** un automa a stati finiti non deterministico è una quintupla $A = \langle Q, \Sigma, \delta, q_0, F \rangle$ dove:
- Q è l'insieme degli stati
- $\Sigma$ è l'alfabeto di input
- $\delta: Q \times \Sigma \to \wp(Q)$ è la funzione di transizione
- $q_0 \in Q$ è lo stato iniziale
- $F \subseteq Q$ è l'insieme degli stati finali

In sostanza vuol dire che nel grafo, nel quale non c’è più il limite al numero di frecce che escono da uno stato, c’è un cammino che ha come etichetta $w$, che parte da uno stato iniziale e arriva in uno stato finale. Una parola $w = a_1, a_2,...,a_n, (a_1,...a_n \in \Sigma, n \geq 0)$ è *accettata* da $A$ se esistono stati $q_1, q_2,...,q_n \in Q$ tali che $q_i \in \delta(q_{i-1}, a_i), 1 \leq i \leq n, q_n \in F$. L'insieme delle parole accettate da $A$ si dice *linguaggio accettato (o riconosciuto)* da $A$ e si denota con $L(A)$. A ogni automa non-deterministico $A = \langle Q, \Sigma, \delta, q_0, F \rangle$ associamo un grafo diretto con frecce etichettate:
- I vertici sono gli stati
- Le frecce sono le triple $(q,a,p)$ con $q \in Q$, $a \in \Sigma$, $p \in \delta(q,a)$
- Lo stato iniziale è denotato da una freccia entrante
- Gli stati finali sono identificati dal doppio bordo

Ci si convince facilmente a questo punto che:
- Una parola *w* può essere etichetta di molti cammini uscenti dallo stato iniziale, o anche nessuno
- La parola *w* è accettata $\iff$ c'è un cammino dallo stato iniziale a uno stato finale con etichetta *w*
- Questo non esclude però che vi possano essere altri cammini con etichetta *w* che partono dallo stato iniziale e terminano in uno stato non finale
> [!example] 
> ![[Screenshot From 2026-03-06 11-03-22.png|200]]
> ![[Screenshot From 2026-03-06 11-03-35.png|200]]
  >![[Pasted image 20260306110549.png|200]]
  > 
> Questa computazione non è accettata perché non finisce in uno stato finale, ma questo non significa che la parola non sia accettata perché possono esserci altre computazioni (come la verde) che terminano in uno stato finale

Apparentemente questo tipo di progettazione non ci da nessun vantaggio, perché le computazioni da fare diventerebbero troppe. Supponiamo infatti di avere 2 cammini possibili per ogni scelta, e dovessimo leggere una parola di lunghezza 10, questo vorrebbe dire che dovremmo fare 210 computazioni. Proviamo però a ragionare come segue: se mi trovo in q0 e il mio input è $a$, posso andare o in q0 o in q1, ecc… Se mi trovo in q1 e l’input è $b$, andrò in q2,… Quindi ad ogni passo, in base alla colonna che avevo al passo precedente, e alla lettera di input, produco una colonna nuova. Questo procedimento è praticamente uguale a quello di un automa a stati deterministico, perché in base alla configurazione del mio dispositivo e alla lettera di input, determina la configurazione successiva. In questo caso la differenza è che abbiamo delle colonnine di stati.
> [!example] ![[Pasted image 20260306110859.png|200]]

**Teorema:** Sia $A$ un automa a stati finiti non deterministico, esiste effettivamente un equivalente automa a stati finiti deterministico $A’$ tale che $L (A) = L (A’)$, cioè esiste l’automa deterministico e l’algoritmo per costruire questo automa a partire dall'automa non deterministico. Quindi, prima andremo a disegnare l’automa a stati non deterministico (prima immagine con più percorsi in verde e rosso), cosa che risulta molto più facile rispetto al crearne uno deterministico, e poi sapremo che esiste l’algoritmo in grado di trasformare quel grafo nel grafo di un automa a stati deterministico (seconda immagine con colonne). 
**Dimostrazione:** Gli stati dell’automa deterministico registreranno l’insieme gli stati raggiungibili dalle computazioni dell’automa non deterministico sul medesimo input $\rightarrow$ Accetto se fra questi ce n’è uno finale

**Costruzione:** Sia $A = \langle Q, \Sigma, \delta, q_0, F \rangle$ l'automa non deterministico. Definisco l'automa deterministico $A' = \langle Q', \Sigma, \delta', s_0, F' \rangle$ come segue: 
- L'insieme degli stati è l'insieme $Q' = \wp(Q)$ costituito dai sottoinsiemi di $Q$ (uno stato dell' $A'$ deterministico è in realtà un'insieme di stati dell' $A$ non deterministico)
- Lo stato iniziale è $s_0 = \{q_0\}$ (singleton costituito dallo stato iniziale dell' $A$ non deterministico)
- Gli stati finali sono tutti i sottoinsiemi di $Q$ che contengono almeno un elemento di $F$ (quindi almeno uno stato finale dell' $A$ di partenza), cioè: $F' = \{s \in \wp(Q) | s \cap F \neq \emptyset \}$ 
- La funzione di transizione $\delta'(s,a) = \displaystyle\bigcup_{q \in s} \delta(q,a), s \in \wp(Q), a \in \Sigma$ (prende in input uno stato dell' $A$, cioè un insieme degli stati dell'automa non-deterministico, e una lettera $a$)

**Dimostrazione:** per ogni $w \in \Sigma^*$ risulta che $\hat{\delta'}(s_0,w) = \{q \in Q | \text{nel grafo di A c'è un cammino da}\ q_0 \ \text{a} \ q \ \text{con etichetta} \ w\}$ (dimostrazione per induzione sulla lunghezza di $w$). Pertanto sono equivalenti: 
- $w\ \text{è accettata da} \ A$
- $\hat{\delta'}(s_0,w) \ \text{contiene uno stato di} \ F$
- $\hat{\delta'}(s_0,w) \in F'$
- $w \ \text{è accettata da} \ A'$
$\implies L(A) = L(A')$ 

Come abbiamo accennato precedentemente, però, sarebbero presenti troppi stati da controllare (se $Q=n$, allora $Q’=2^n$). Per fortuna <u>non tutti questi stati sono necessari: possiamo definire uno stato necessario/accessibile come un cammino nel grafo dell’automa che, dallo stato iniziale, mi porta a $q$. Tutti gli altri stati si dicono inaccessibili.</u>

**Definizione:** sia $A = \langle Q, \Sigma, \delta, q_0, F \rangle$ un automa a stati finiti deterministico. Uno stato $q \in Q$ si dice *accessibile* se $q = \hat{\delta}(q_0,w)$ per qualche $w \in \Sigma^*$.
Se io elimino tutti questi stati inaccessibili e restringo la funzione di transizione (delta cappello) solo agli stati accessibili, ottengo un automa con un numero di stati più piccolo, equivalente a quello di partenza, che accetta lo stesso linguaggio di quello di partenza.
> [!example] 
> ![[Screenshot From 2026-03-06 11-26-56.png|200]] 
> 
> Nell’esempio vediamo come lo stato q1 è stato eliminato perché, iniziando da q0 (stato di partenza), non c’è nessun cammino che mi porta in q1.

**Algoritmo di determinizzazione (per prevedere solo gli stati accessibili):**
1. Prima di tutto inserisco nella lista degli stati, lo stato iniziale dell’automa deterministico, che come abbiamo visto prima corrisponde all'insieme costituito dal solo stato iniziale dell’automa non deterministico di partenza: $s_0=\{q_0\}$ 
2. Considero ogni stato $r$ della lista (che per ora contiene il solo elemento $q0$, singleton) e ogni lettera $a \in \Sigma$ 
3. Calcolo il valore della funzione di transizione sulla coppia stato-lettera: $s= \delta'(r, a)$; Una volta calcolato ho due possibilità: o lo stato è già nella lista degli stati, oppure no.   
4. Se $s$ non è nella lista degli stati allora lo aggiungo alla lista degli stati. 
5. Se $s$ contiene uno stato finale di $A$, aggiungo $s$ alla lista degli stati finali. 
> [!example] 
> ![[Screenshot From 2026-03-06 11-32-20.png|200]]
> ![[Screenshot From 2026-03-06 11-32-32.png|200]]
> ![[Screenshot From 2026-03-06 11-33-08.png|200]]
> 
> Come possiamo notare, l’automa deterministico ha più stati rispetto a quelli che avevamo trovato noi manualmente nel primo esempio abac (pag.11), successivamente vedremo infatti un altro algoritmo in grado di trovare partendo da un automa non deterministico, un automa deterministico con il minimo numero di stati possibile (automa minimo).

> [!info] Caso difficile:
> ![[Pasted image 20260308174311.png|200]]
###### Automa a stati finiti non deterministico con epsilon-transizioni
Un automa a stati finiti non deterministico con $\epsilon$-transizioni è una quintupla $A = \langle Q,\Sigma, \delta, q_0, F \rangle$, dove $Q,\Sigma, q_0, F$ sono come nell'automa deterministico, mentre la funzione $\delta: Q \times (\Sigma \cup \{\epsilon\}) \rightarrow \wp(Q)$ è la funzione di transizione, che ora non è definita solo sulle coppie *stato-lettera*, ma anche sulle coppie *stato-parolaVuota*, e restituisce un insieme di stati (potrebbe essere anche vuoto). Questo modello mi permette di saltare da uno stato all'altro senza leggere il nastro di input (anche per questo modello c'è un algoritmo che mi permette di passare da un automa di questo tipo ad un automa a stati finiti deterministico). Una parola $w \in \Sigma^*$ è *accettata* da $A$ se esistono $n\geq 0$, $a_1,...,a_n \in \Sigma \cup \{\epsilon\}$ e $q_1,..q_n \in Q$ tali che $w = a_1a_2...a_n$, $q_i \in \delta(q_{i-1},a_i)$, $1 \leq i \leq n$, $q_n \in F$ . In altre parole la parola è accettata dall'automa se esiste un cammino dallo stato iniziale a uno finale nel grafo con etichetta $w$, e se $w$ può essere scritta come una sequenza di lettere e parole vuote. L'insieme delle parole accettate da $A$ si dice *linguaggio accettato (o riconosciuto)* da $A$ e si denota con $L(A)$.

**Grafo dell'automa:** A ogni automa non deterministico con $\epsilon$-transizioni $A = \langle Q,\Sigma, \delta, q_0, F \rangle$ associamo un grafo diretto con frecce etichettate:
- I vertici sono gli stati
- Le frecce sono le triple $(q,a,p)$ con $q \in Q$, $a \in \Sigma \cup \{\epsilon\}$, $p \in \delta(q,a)$ 
- Lo stato iniziale è denotato da una freccia entrante
- Gli stati finali sono identificati dal bordo doppio

Ci si convince facilmente a questo punto che:
- Una parola *w* può essere etichetta di molti cammini uscenti dallo stato iniziale, o anche nessuno
- La parola *w* è accettata $\iff$ c'è un cammino dallo stato iniziale a uno stato finale con etichetta *w*
- Questo non esclude però che vi possano essere altri cammini con etichetta *w* che partono dallo stato iniziale e terminano in uno stato non finale
> [!example] ![[Screenshot From 2026-03-12 10-16-28.png|200]]
> Nota che tra le colonne ora appare anche $\epsilon$, cosa che prima non era contemplata

**Teorema:** Sia $A$ un automa a stati finiti non deterministico con $\epsilon$-transizioni. Esiste effettivamente un automa a stati finiti deterministico $A'$ tale che $L(A) = L(A')$ (accetta lo stesso linguaggio)

**Dimostrazione:** Gli stati dell’automa deterministico registreranno l’insieme degli stati raggiungibili dalle computazioni dell’automa non deterministico sul medesimo input. Accetto la parola se fra questi ce uno finale

**$\epsilon$-chiusura:** Se però con l’automa a stati finiti non deterministico le nostre computazioni si limitavano alla potenza-numero-di-lettere della parola, ora teoricamente posso avere anche infinite ε-transizioni (la concatenazione di un miliardo di ε, è sempre una parola vuota). Quindi bisogna utilizzare una tecnica per trattare queste ε-transizioni: la ε-chiusura di uno stato. 

**Definizione:** sia $A = \langle Q,\Sigma, \delta, q_0, F \rangle$ un automa a stati finiti non deterministico con $\epsilon$-transizioni, la $\epsilon$-chiusura di uno stato $q \in Q$ è il più piccolo sottoinsieme $C_{\epsilon}(q)$ di $Q$ tale che:
- $q \in C_{\epsilon}(q)$
- Per ogni $p \in C_{\epsilon}(q)$, $\delta(p,\epsilon) \subseteq C_{\epsilon}(q)$

La $\epsilon$-chiusura di un insieme di stati $S \subseteq Q$ è l'unione delle $\epsilon$ chiusure dei singoli stati di $S: C_{\epsilon}(S) = \bigcup_{q \in S} C_{\epsilon}(q)$.  In altre parole la ε-chiusura di uno stato è l’insieme degli stati che da questo si possono raggiungere utilizzando solo ε- transizioni (questo insieme deve contenere lo stato stesso).
> [!example] ![[Screenshot From 2026-03-12 10-31-58.png|200]]
> - $C_{\epsilon}(q_0)$ deve contenere $q_3$ perché con $\epsilon$-transizioni vado in $q_3$
> - $C_{\epsilon}(\{q_0,q_1\})$ è la $\epsilon$-chiusura della coppia di stati $q_0,q_1$ che si ottiene facendo l’unione delle chiusure dei due stati

**Calcolo della $\epsilon$-chiusura:**
- *Input*: un automa non deterministico con $\epsilon$-transizioni $A = \langle Q,\Sigma, \delta, q_0, F \rangle$ e uno stato $q \in Q$
- *Output*: la $\epsilon$ -chiusura di $q$
```
R <- (q) // inizialmente la lista contiene solo lo stato q
p <- puntatore al primo elemento di R
while p != NIL   // Prendiamo tutti gli elementi che si ottengono partendo dal primo elemento della // lista e seguendo tutte le ε transizioni e le aggiungiamo alla lista, avendo l’accortezza di        // aggiungere solo quelle nuove, e non quelle che già appaiono in essa, quindi non creiamo             // ripetizioni.)
	 appendi a R tutti gli elementi di $\delta$ (p, $\epsilon$) \ R
	 p <- successivo(p) // passo all'elemento successivo
return R
```
> [!example] ![[Screenshot From 2026-03-12 10-35-50.png|200]]

**Costruzione:** Sia $A = \langle Q,\Sigma, \delta, q_0, F \rangle$ l'automa non deterministico con $\epsilon$-transizioni, definisco l'automa deterministico $A' = \langle Q',\Sigma, \delta', q_0, F' \rangle$ come segue:
- L'insieme degli stati è l'insieme $Q' = \wp(Q)$ costituito dai sottoinsiemi di $Q$ (cioè uno stato di $A'$  è un insieme di stati di $A$ )
- Lo stato iniziale è $s_0 = C_{\epsilon}(q_0)$ (stato iniziale è una $\epsilon$ -chiusura, quindi non si tratta più di un singleton, ma di un insieme di stati)
- Gli stati finali sono tutti i sottoinsiemi di Q che contengono almeno un elemento di $F$, cioè: $F' = \{s \in \wp(Q) | s \cap F \neq \emptyset\}$ (sottoinsiemi che contengono almeno uno stato finale dell' $A$ precedente)
- La funzione di transizione è $\delta'(s,a): C_{\epsilon}(\bigcup_{p \in s} \delta(p,a))$, $s \in \wp(Q)$, $a \in \Sigma$ (insieme di tutti quegli stati che posso raggiungere in $A$ partendo da uno qualunque degli stati di $s$ e poi seguendo la freccia con etichetta a, e poi tute le frecce che voglio con etichetta $\epsilon$)

**Dimostrazione:** Per ogni $w \in \Sigma^*$ risulta che $\hat{\delta'}(s_0,w) = \{q \in Q |$ nel grafo di $A$  c'è un cammino da $q_0$ a $q$ con etichetta $w$ $\}$. Pertanto sono equivalenti:
- $w$ è accettata da $A$
- $\hat{\delta'}(s_0,w)$ contiene uno stato di $F$
- $\hat{\delta'}(s_0,w) \in F'$
- $w$ è accettata da $A'$

Quindi $L(A) = L(A')$
Sostanzialmente  è sempre possibile costruire un automa deterministico che accetta il medesimo linguaggio, ed esiste inoltre un algoritmo che mi permette di trovare questo automa deterministico equivalente

**Algoritmo di determinazione:** 
- *Input*: un automa non deterministico con $\epsilon$ transizioni $A = \langle Q,\Sigma, \delta, q_0, F \rangle$
- *Output*: un automa deterministico equivalente $A' = \langle Q',\Sigma, \delta', q_0, F' \rangle$

n ← 1  // numero degli stati  
j ← 0  
$s_0 \leftarrow C_{\epsilon}(q_0)$  
if $s_0 \cap F \neq \emptyset$ then $F'$ else $F' \leftarrow \emptyset$
while $j < n$
	for $a \in \Sigma$
		$s \leftarrow \bigcup_{q \in s_j}\delta(q,a)$
		$s \leftarrow C_{\epsilon}(s)$
		if $s = s_i$ per qualche $i < n$
			 $\delta'(j,a) \leftarrow i$
	    else
	         $s_n \leftarrow s$
	         $\delta'(j,a) \leftarrow n$
	         if $s \cap F \neq \emptyset$ then $F' \leftarrow F \cup \{n\}$
	         $n \leftarrow n+1$
	$j \leftarrow j+1$
	
> [!example] 
> - valuto le chiusure di ogni stato
> ![[Screenshot From 2026-03-12 11-25-10.png|200]]
> - valuto la funzione $\delta$ su questo automa
> ![[Screenshot From 2026-03-12 11-25-29.png|200]]
> 
> Il grafo dell’esempio sarà quindi:
> ![[Screenshot From 2026-03-12 11-29-04.png|200]]
###### Espressioni regolari
La teoria dei linguaggi formali può essere vista come l’arte di rappresentare degli oggetti infiniti con una descrizione finita. Per costruire nuovi linguaggi partendo da linguaggi noti, possiamo applicare delle operazioni: 
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
6. Chiusura di Kleene: $L^* = \bigcup\limits_{n \geq 0}L^n = \{u_1u_2...u_n | n \geq 0, u_1,...,u_n \in L\}$. Prendo tutte le potenze di un linguaggio e ne faccio l’unione. A differenza delle prime 3 operazioni che se applicate a insiemi finiti mi ridanno insiemi finiti, la chiusura di Kleene mi dà un insieme infinito, a meno che l’insieme di partenza non sia o l’insieme vuoto, o l’insieme che contiene la sola parola vuota. Conterrà quindi tutte le parole che posso ottenere concatenando un numero arbitrario di parole di L. Esempio:
<div style="display: flex; justify-content: center;">
  <img src="Screenshot From 2026-03-14 16-51-56.png" width="400">
</div>

**Osservazione:** unione, concatenazione e chiusura di Kleene sono dette operazioni regolari

Quindi, i linguaggi finiti li sappiamo definire facendo una semplice lista delle parole; con questi linguaggi finiti posso poi fare le operazioni di unione e concatenazione ed ottenere nuovi linguaggi finiti, e con la chiusura di Kleene ottengo nuovi linguaggi infiniti. Con questi linguaggi infiniti posso fare altre operazioni di unione e concatenazione, ottenendo altri linguaggi sempre più complessi. Questo vuol dire che posso scrivere la definizione dei linguaggi infiniti semplicemente indicando le operazioni che eseguo per ottenerli: ovvero le espressioni regolari. <u>Le espressioni sono parole che costruisco partendo dalle lettere dell’alfabeto e dal simbolo ∅, seguendo le regole.</u> Rappresentano quindi una sequenza di operazioni fatte a partire dall’alfabeto iniziale. 

**Definizione:** sia $\hat\Sigma$ l'alfabeto ottenuto aggiungendo a $\Sigma$ le lettere $\emptyset, +, *, (,), . \implies$ si dicono *espressioni regolari* sull'alfabeto $\Sigma$ le parole sull'alfabeto $\hat\Sigma$ che si ottengono applicando un numero finito di volte le regole seguenti regole:
1. Ogni lettera $a \in \Sigma$ è un'espressione regolare, $\emptyset$ è un'espressione regolare
2. Se $E$ e $F$ sono espressioni regolari, allora $(E+F), (EF)$ e $E^*$ sono espressioni regolari.

A ogni espressione regolare è associato un linguaggio, detto *linguaggio denotato dall'espressione regolare* e definito dalle regole seguenti:
1. Per ogni $a \in \Sigma$, l'espressione regolare $a$ denota il linguaggio $\{a\}$; l'espressione regolare $\emptyset$ denota il linguaggio vuoto
2. Detti $L_E$ e $L_F$ i linguaggi denotati dalle espressioni regolari $E$ e $F$, i linguaggi denotati dalle espressioni regolari $E+F, (EF), E^*$ sono, rispettivamente, $L_E \cup L_F, L_EL_F, L_{E^*}$

**Definizione:** i linguaggi denotati da espressioni regolari si dicono *linguaggi regolari* (i linguaggi di tipo 0 sono detti anche linguaggi regolari)

**Osservazione:** la classe dei linguaggi regolari è la più piccola famiglia di linguaggi che:
- Contiene i linguaggi finiti
- è chiusa per le operazioni regolari

Perché se due linguaggi sono denotati da espressioni regolari, la loro unione sarà denotata anche lei da un’espressione regolare. Lo stesso per la concatenazione. La chiusura di Kleene sarà denotata dall’espressione regolare del linguaggio con la * vicino

**Osservazione:** possiamo omettere qualche parentesi, rispettando la priorità:
- Chiusura di Kleene
- Concatenazione
- Somma
> [!example] Se si omettono le parentesi, la priorità delle operazioni sono: chiusura di Kleene, concatenazione e poi somma.
> - $a+b$: espressione regolare che denota il linguaggio composto dalle due parole a e b
> - $(a+b)^*$: chiusura di Kleene, tutte le concatenazioni di parole che appartengono all'insieme a, b. Cioè tutte le concatenazioni di lettere a, b . In sostanza è l'insieme delle parole sull'alfabeto (a, b)
> - $\emptyset$: se concateno 0 parole dell'insieme vuoto (linguaggio che contiene solo la parola vuota)
> - $((a+b)(a+b)(a+b))$: ottengo tutte le parole di lunghezza 3, e poi ne faccio la chiusura di Kleene e ottengo  tutte le parole che hanno lunghezza multipla di 3
> - $(a+ab)^*$:
> - $(a+b)^*abb$: otteniamo tutte le lettere che finiscono con $abb$
> - $(b^*ab^*ab^*)^*b^*ab^*$:
> - $(a+b+c)^*abac (a+b+c)^*$:
> - $((a+b+c),)^*(a+b+c) \& (a+b+c)+a+b+c$

**Teorema di Kleene:** un linguaggio è regolare $\iff$ è riconosciuto da un automa a stati finiti. 

Esiste una procedura effettiva che:
- Data un'espressione regolare, produce un automa a stati finiti che accetta il linguaggio denotato da tale espressione (*sintesi*)
- Dato un'automa a stati finiti, produce un'espressione regolare che denota il linguaggio accettato da tale  automa (*analisi*)

Utilizzeremo esclusivamente automi non deterministici con 𝜀-transizioni con un unico stato finale (non restrittivo)

Un’espressione regolare, dove E e F sono espressioni regolari, può essere:

 <div style="display: flex; align-items: flex-start; gap: 20px;">
   <div style="flex: 1;">
     <img src="Screenshot From 2026-03-13 11-05-00.png" style="width: 100%; border-radius: 8px;">
   </div>
  <div style="flex: 1.5;">
<li>una lettera, un insieme vuoto o una parola vuota (che corrisponde alla chiusura di Kleene dell’insieme vuoto)</li>
<li>(E+F), cioè l’unione</li>
<li>(EF), cioè la concatenazione</li>
<li>E^*, cioè la chiusura di Kleene</li>
   </div>
 </div>
 
 <div style="display: flex; align-items: flex-start; gap: 20px;">
   <div style="flex: 1;">
     <img src="Screenshot From 2026-03-13 11-05-38.png" style="width: 100%; border-radius: 8px;">
   </div>
  <div style="flex: 1.5;">
<b>Unione:</b> avendo già costruito gli automi di E e F (gli ovali nella figura), si aggiunge un nuovo stato iniziale e un nuovo stato finale, e due ε transizioni che mi portano dallo stato iniziale nuovo agli stati iniziali vecchi, e due ε transizioni che mi portano dagli stati finali vecchi allo stato finale nuovo. Questo nuovo automa accetta l’unione dei linguaggi, accettati dai due automi E e F. Quindi tutti i cammini riusciti dallo stato iniziale a quello finale, avranno come etichetta (l’ovale) una parola accettata dall’automa E o F.
   </div>
 </div>

Siano $A_1 = \langle Q_1,\Sigma,\delta_1, i_1, \{f_1\}\rangle$ e $A_2 = \langle Q_2,\Sigma,\delta_2, i_2, \{f_2\}\rangle$ gli automi che riconoscono i linguaggi denotati dalle espressioni regolari $E$ e $F$. 
Supponiamo $Q_1 \cap Q_2 = \emptyset$. 
Costruiamo $A = \langle Q,\Sigma,\delta, i, \{f\}\rangle$ come segue:
- $Q = Q_1 \cup Q_2 \cup \{i,f\}$ ove $i,f$ sono due nuovi stati
- $i,f$ sono rispettivamente lo stato iniziale e l'unico stato finale
- Le frecce del grafo di $A$ sono le frecce dei grafi di $A_1$ e $A_2$ e, inoltre, $(i,\epsilon,i_1), (i,\epsilon,i_2), (f_1,\epsilon, f), (f_2,\epsilon, f)$

Allora $A$ accetta il linguaggio denotato da $(E+F)$
> [!example] ![[Screenshot From 2026-03-13 11-13-32.png|200]]

 <div style="display: flex; align-items: flex-start; gap: 20px;">
   <div style="flex: 1;">
     <img src="Screenshot From 2026-03-13 11-14-59.png" style="width: 100%; border-radius: 8px;">
   </div>
  <div style="flex: 1.5;">
<b>Concatenazione:</b> si connette con una ε-transizione lo stato finale del primo automa con lo stato iniziale del secondo automa. Lo stato iniziale del nuovo automa sarà quello del primo automa, mentre il suo stato finale sarà quello del secondo automa.
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
> ![[Screenshot From 2026-03-13 11-17-52.png|200]]

 <div style="display: flex; align-items: flex-start; gap: 20px;">
   <div style="flex: 1;">
     <img src="Screenshot From 2026-03-13 11-18-53.png" style="width: 100%; border-radius: 8px;">
   </div>
  <div style="flex: 1.5;">
<b>Chiusura di Kleene:</b> abbiamo l’automa E che riconosce una certa espressione, e poi il nuovo automa che riconosce la chiusura di Kleene del linguaggio accettato da questo, ovvero che accetti tutte le parole che si ottengono concatenando 0 o più parole accettate dall’automa di partenza. Aggiungo una ε transizione che dallo stato finale mi riporta a quello iniziale: così da accettare tutte le parole che si ottengono concatenando una o più parole dell’automa di partenza. Mi manca però il riconoscimento della parola vuota, quindi aggiungo un nuovo stato iniziale e un nuovo stato finale e poi una ε-transizione dallo stato iniziale nuovo a quello vecchio, una ε-transizione dallo stato finale vecchio a quello nuovo, e una ε-transizione dallo stato iniziale nuovo allo stato finale nuovo.
   </div>
 </div>

Siano $A_1 = \langle Q_1,\Sigma,\delta_1, i_1, \{f_1\}\rangle$ l'automa che riconosce il linguaggio denotato dall'espressione regolare $E$. 
Costruiamo $A = \langle Q,\Sigma,\delta, i, \{f\}\rangle$ come segue:
- $Q = Q_1 \cup \{i,f\}$ ove $i,f$ sono due nuovi stati
- $i,f$ sono rispettivamente lo stato iniziale e l'unico stato finale
- Le frecce del grafo di $A$ sono le frecce del grafo di $A_1$ e, inoltre, $(i,\epsilon,i_1), (i,\epsilon,i_2), (f_1,\epsilon, f), (f_2,\epsilon, f)$

Allora $A$ accetta il linguaggio denotato da $E^*$
> [!example]
> ![[Pasted image 20260313112221.png|200]]

**L'algoritmo di sintesi**
**Proposizione:** data un'espressione regolare $G$, si può effettivamente costruire un'automa a stati finiti che riconosce il linguaggio denotato da $G$:
1. se $G$ è un’espressione regolare di base (lettera o insieme vuoto), allora restituisco l’automa corrispondente; 
2. se $G = (E + F)$, allora calcolo gli automi $A_1$ e $A_2$ corrispondenti a $E$ e $F$; costruisco l’automa dell’espressione $(E + F)$; 
3. se $G = (EF)$, allora calcolo gli automi $A_1$ e $A_2$ corrispondenti a $E$ e $F$; costruisco l’automa dell’espressione $(EF)$; 
4. se $G = E^*$ , allora calcolo l’automa $A_1$ corrispondente a $E$; costruisco l’automa dell’espressione $E^*$;

Esempio completo nelle slide del prof

**Dall’automa a stati finiti all’espressione regolare:** sia $A$ un automa a stati finiti e per semplicità supporremo che:
- Gli stati sono $q_1, q_2,…q_n$
- Lo stato iniziale è $q_1$

> [!tip] Avendo un automa A, vogliamo che: 
> abbia un unico stato finale e, per semplicità, non abbia frecce che entrano nello stato iniziale né frecce che escono dallo stato finale. Questo implica che un percorso uscito dallo stato iniziale non potrà più rientrarci (perché non ci sono frecce entranti in esso) così come non potrà passare per lo stato finale se non all’ultimo (perché non ci sono frecce uscenti da esso). Aggiungo quindi delle $\epsilon$ transizioni che mi portano dal nuovo stato iniziale al vecchio stato iniziale, e dai vecchi stati finali al nuovo stato finale. Questo nuovo automa accetterà le stesse parole dell’automa vecchio.
> Inoltre, supponiamo anche di dare dei nomi ai vari stati del nostro automa: $q_1$, $q_2$, …, $q_n$, dove lo stato
> iniziale corrisponde al $q_{n-1}$ (quindi al penultimo stato) e lo stato finale a $q_n$. Prendiamo dei linguaggi che in qualche modo approssimano il linguaggio che vogliamo cercare e cerchiamo delle espressioni regolari per questi linguaggi.
> 
>![[Pasted image 20260319171718.png|300]]

Per ogni $i, j, k$ con $0 \leq k < i, j \leq n$ (indici i e j compresi tra 1 e n, mentre k compreso tra 0 e n) consideriamo l’insieme $L_{ijk}$ contenente le seguenti parole: *le etichette dei cammini nel grafo di $A$ che iniziano in $q_i$, terminano in $q_j$, e attraversano solo stati di indice $\leq k$* (che attraversi soltanto, non che parta o arrivi in uno di questi stati)

**Osservazione:** $L(A)= \bigcup\limits_{q_j \in F} L_{1jn}$. Quindi, se riesco a trovare espressioni regolari per tutti gli $L_{ijk}$, avrò anche un’espressione regolare di $L(A)$

> [!example] ![[Pasted image 20260319110515.png|200]]
> - $q_4$ stato finale
> - $q_5$ stato finale
> - quindi devo prendere tutti i cammini che partono da $q_2$, che arrivano in $q_5$ e che attraversano soltanto $q_1 \implies L(A) = L_{n-1 \ n \ n-2}$ dato dal cammino $q_2 \xrightarrow[]{b} q_1 \xrightarrow[]{c} q_5$

Se riesco a trovare un’espressione regolare per tutti gli $L_{ijk}$, avrò un’espressione regolare per $L_{n-1 \ n \ n-2}$, che
è proprio il linguaggio accettato dall’automa. Ho quindi ottenuto l’espressione regolare corrispondente al
Linguaggio. Quello che dobbiamo mostrare è un meccanismo che ci permetta di determinare delle espressioni regolari che denotano questi linguaggi $L_{ijk}$ Per ottenere tutti questi $L_{ijk}$, si procede per induzione su $k$.

Il linguaggio $L_{ij0}$ è denotato dall’espressione regolare $E_{ij0}$ ottenuta come segue (supponendo $k = 0$):
- La somma delle etichette delle frecce da $q_i$ a $q_j$, se ne esistono e $i\neq j$ 
- La somma delle etichette delle frecce da $q_i$ a $q_j$ più $\epsilon$ se ne esistono e $i = j$
- $\epsilon$ se non ci sono frecce da $q_i$ a $q_j$ e $i = j$

> [!example] ![[Screenshot From 2026-03-19 11-07-41.png|200]]
> Dato che $k \leq 0$ ciò significa che i cammini devono attraversare 0 stati (cioè si tratta di frecce dirette)

Supponiamo $k\geq 1$. I cammini nel grafo $A$ con origine in $q_i$, termine in $q_j$ e che attraversano esclusivamente stati di indice $\leq k$, si possono ripartire in due categorie:
1. Cammini da $q_i$ a $q_j$ che attraversano solo stati di indice $\leq k-1$. Questi cammini per $q_k$ potrebbero no passare proprio, oppure possono passare per $q_1,q_2,...,q_{k-1}$. Per questi l’espressione regolare già l’abbiamo
<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319173347.png" width="200">
</div>
2. I cammini costituiti dalla concatenazioni di (cioè i cammini che passano una o più volte per $k$):
	- Un segmento da $q_i$ a $q_k$ (segmento dallo stato iniziale fino a $q_k$)
	- Zero o più cammini con origine e termine in $q_k$ (dei cammini che da $q_k$ arrivano al secondo passaggio con $q_k$, e così via, fino all’ultimo passaggio con $q_k$)
	- Un segmento finale con origine da $q_k$ a $q_j$ (un ultimo segmento che da $q_k$ arriva allo stato finale)

Tutti che attraversano solo stati di indice $\leq k-1$

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319174122.png" width="200">
</div>

**Passo induttivo:** ora dobbiamo trovare l’espressione $E_{ijk}$, cioè l’espressione per il linguaggio costituito dalle etichette di tutti i cammini delle forme precederti. Innanzitutto facciamo  l’unione dei cammini del primo punto con i cammini del secondo punto. I cammini del secondo gruppo saranno a loro volta riportati come la concatenazione del primo segmento con i segmenti intermedi (concatenazione di 0 o più coppi[]()e di un 0 espressione utilizzando la chiusura di Kleene) e l’ultimo segmento. In sostanza le espressioni regolari $E_{ijk}$ che denotano i linguaggi $L_{ijk}$ devono soddisfare le relazioni ricorsive $E_{ijk} = E_{ijk−1} + (E_{ikk−1})(Ek k k−1)(E_{kjk−1}), \ 1 \leq k < i, j \leq n$. Un’espressione regolare di $L(A)$ è data dalla somma delle $E_{1jn}$ con $q_j$ stato finale. In questo mode è possibile calcolare tutte le espressioni regolari, compresa $E_{n-1 \ n \ n-2}$ che denota il linguaggio accettato dall’automa (esempio nelle slide 8-11)

**Regola pratica:** Per fare tutto ciò utilizzeremo il Metodo di eliminazione degli stati partendo dal nostro automa che corrisponde ad un grafo sulle cui etichette delle frecce ci sono espressioni particolari che possono essere o lettere, o ε o somma di lettere.
1. ci riduciamo al caso in cui:
	- $A$ ha un unico stato finale 
	- Nel grafo di $A$ non ci sono frecce che entrano nello stato iniziale, ne frecce che escono dallo stato finale
2. Sostituiamo le etichette multiple con la loro somma 
3. Selezioniamo uno stato $q_k$ che non sia né iniziale né finale; per ogni freccia $q_i \xrightarrow[]{E_{ik}} q_k$ che entra nello stato $q_k$ e ogni freccia $q_k \xrightarrow[]{E_{kj}} q_j$ che esce da tale stato $(i, j \neq k)$ 
	- Rimpiazziamo l’etichetta della freccia $q_i \xrightarrow[]{E_{ij}} q_j$ con $E_{ij} + E_{ik}E_{kk}^*E_{kj}$ (o $E_{ij} + E_{ik}E_{kj}$ qualora manchi la freccia da $q_k$ a $q_k$); 
	- se tale freccia è assente, la creiamo (con etichetta $E_{ik}E_{kk}^*E_{kj}$ , ovvero $E_{ik}E_{kj}$ ); 
	- rimuoviamo lo stato $q_k$ e tutte le frecce che entrano o escono da esso; 
4. Ripetiamo il passo precedente fino a ottenere un automa con due stati e una sola freccia; 
5. L’etichetta di tale freccia è l’espressione regolare cercata.

> [!example] > ![[Pasted image 20260319181524.png]]
###### Operazioni booleane
**Proposizione:** la classe dei linguaggi regolari è un’algebra Booleana (altra conseguenza del teorema di Kleene)

**Dimostrazione:** Ci sono delle operazioni che sono naturali sulle espressioni regolari e altre operazioni che sono naturali sugli automi a stati finiti. Come abbiamo detto in precedenza, infatti, la famiglia dei linguaggi regolari è chiusa per le operazioni di unione, concatenazione e chiusura di Kleene. È anche chiusa per complemento ed intersezione, quindi per tutte le operazioni booleane.
- *Unione:* L’unione dei linguaggi denotati dalle espressioni $E$ e $F$ è denotato dall’espressione $E$ + $F$. 
- *Complemento:* Per il Teorema di Kleene, un linguaggio regolare è accettato da un automa a stati finiti deterministico. Il complemento è accettato dall’automa ottenuto scambiando stati finali e non. Quindi, per il Teorema di Kleene, è regolare 
- *Intersezione:* L’intersezione si esprime per mezzo di unione e complemento: $L \cap M = \Sigma^* - ((\Sigma^* -L) \cup (\Sigma^* - M))$

**Osservazione:** anche la differenza insiemistica si esprime per messo di unione e completamento: $L - M = \Sigma^* - ((\Sigma^* -L) \cup M)$ 

**Osservazione:** Date due espressioni regolari, posso effettivamente trovare espressioni regolari per l’unione, l’intersezione e i complementi dei linguaggi da esse denotati.

**Equivalenza di espressioni regolari:** $a^* + (a^*ba^*)^* = (a+b)^*$. Date due espressioni regolari ho il problema di decidere se queste sono equivalenti o meno

**Osservazioni:** dati due insiemi $X,Y$ sono uguali se e solo se $(X-Y) \cup (Y-X) = \emptyset$ 

La soluzione al mio problema è: date le espressioni regolari $E,F$
1. Calcolare automi deterministici che accettano i linguaggi $L_E$ e $L_F$ denotati da $E$ e $F$, rispettivamente; 
2. Calcolare un automa deterministico che accetta $(L_E − L_F ) \cup (L_F − L_E )$; 
3. Verificare se tutti gli stati finali sono inaccessibili.

###### Automa minimo
Si tratta di un’automa deterministico con il minor numero possibile di stati.
L’automa di Nerode di un linguaggio regolare L è un automa deterministico che accetta L col minimo numero di stati possibile. Quindi definisco prima la congruenza di Nerode, così da definire un’equivalenza cioè una partizione dell’insieme Σ* in un numero finito di classi di equivalenza. Costruisco in seguito un automa minimo partendo da queste classi di equivalenza, che accetta il linguaggio L. 


Equivalenza: insieme sul quale è definita una relazione riflessiva (cioè ogni elemento è equivalente a se stesso), simmetrica (se x è equivalente a y, allora y è equivalente a x) e transitiva (se abbiamo tre elementi e il primo è equivalente al secondo, il secondo è equivalente al terzo, allora il primo è equivalente al terzo). C’è una corrispondenza tra la nozione di equivalenza, quindi un insieme di coppie che soddisfa le proprietà citate sopra, e la nozione di partizione su un insieme, cioè la divisione di un insieme in sottoinsiemi che sono a due a due disgiunti (privi di elementi in comune), e che uniti mi danno tutto l’insieme. Quindi se ho un’equivalenza su un insieme, le classi di equivalenza (cioè l’insieme di tutti gli elementi che sono equivalenti l’uno all’altro), costituiscono una partizione sull'insieme. Allo stesso modo se ho una partizione, gli insiemi che la costituiscono sono classi di un’opportuna equivalenza. 


**Definizione:** un’equivalenza è una relazione riflessiva, simmetrica e transitiva

**Osservazione:** ogni equivalenza su un insieme $S$ determina una partizione dell’insieme $S$ in classi di equivalenza. Viceversa, ogni partizione dell’insieme $S$ è associata a un’equivalenza di cui le parti sono classi di equivalenza

**Definizione:** sia $\Sigma$ un alfabeto. Una relazione di equivalenza $\sim$ su $\Sigma^*$ si dice una congruenza destra (risp. Sinistra) se per ogni $u, v \in \Sigma^*$ tali che $u \sim v$ e per ogni lettera $a \in \Sigma$, si ha $ua \sim va$ (risp. $au \ \sim \ av$). Una relazione che sia contemporaneamente una congruenza destra e una congruenza sinistra si dice congruenza. In altre parole qui non abbiamo solo un insieme, ma un insieme che ha un’operazione definita su di esso. Ogni volta che si va a concatenare a dx di due parole equivalenti, una lettera, ottengo due parole equivalenti (rispettivamente per la congruenza sx).

> [!example] Si considerino le relazione $\sim$ su $\Sigma^*$ definite da:
> 1. $u \sim v$ se $u$ e $v$ iniziano con la stessa lettera (o $u$ = $v$ = $\epsilon$)
> 2. $u \sim v$ se $u$ e $v$ terminano con la stessa lettera (o $u$ = $v$ = $\epsilon$)
> 3. $u \sim v$ se $u$ e $v$ hanno la stessa lunghezza
> 4. $u \sim v$ se $u$ e  $v$ le prime due lettere di u sono una permutazione delle prime due lettere di $v$ (o $u = v$)
> 
> Per es. $abaa \sim bab, aababb \sim aaabaa \ \text{ma} \ aababb \not\sim baaabaa$

**Equivalenza di Nerode:**

**Definizione:** sia $L$ un linguaggio sull’alfabeto $\Sigma$. L’equivalenza di Nerode associata al linguaggio $L$ è la relazione $N_L$ definita in $\Sigma^*$ da: $u \ N_L \ v$ se per ogni $y \in \Sigma^*$*, $uy \in L \iff vy \in L$. Cioè, $u \ N_L \ v$ se $u$ e $v$ hanno gli stessi completamenti a destra in $L$

**Lemma:** sia $L \subseteq \Sigma^*$ un linguaggio
1. La relazione $N_L$ è una congruenza destra

**Dimostrazione:** $N_L$ è un’equivalenza. Dobbiamo verificare che per ogni $u, v \in \Sigma^*$ , $a \in \Sigma$, se $u N_L v$, allora $ua \ N_L \ va$. Invero sia $u \ N_L \ v$ e $a \in \Sigma$. Dalla definizione di $N_L$ si ha $uy \in L \iff vy \in L$ per ogni $y \in \Sigma^*$ Sostituendo $y$ con $ay$, ottengo $uay \in L \iff vay \in L$ Quindi si ha anche $ua \ N_L \ va$.

2. $L$ è unione di classi di equivalenza di $N_L$

**Dimostrazione:** Sia $u \ N_L \ v$. Allora $uy \in L \iff vy \in L$ per ogni $y \in \Sigma^*$ . Prendendo $y = \epsilon$ ottengo $u \in L \iff v \in L$. Quindi ogni classe di equivalenza o sta tutta dentro L o sta tutta nel complemento.

**Teorema di Myhill-Nerode (1958):** Sia $L$ un linguaggio sull’alfabeto $\Sigma$, le seguenti proposizioni sono equivalenti:
1. $L$ è regolare
2. $L$ è unione di classi di una congruenza destra su $\Sigma^*$ si dice finito
3. $N_L$ ha indice finito

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260320102535.png" width="200">
</div>

Dimostrazione circolare: 1. $\implies$ 2. $\implies$ 3. $\implies$ 1.

**Dimostrazione 1. $\implies$ 2.:** 
- Sia $L$ regolare
- $L$ è accettata da un automa deterministico $A = \langle Q; \Sigma, \delta, q_0, F \rangle$ con tutti gli stati accessibili
- Per ogni stato $q \in Q$ sia $L_q = \{w \in \Sigma^* | \hat\delta(q_0,w) = q\}$
- Gli insiemi $L_q$ costituiscono una partizione di $\Sigma^*$
- $L$ è l’unione delle classi $L_q$ con $q \in F$
- Quindi basta verificare che l’equivalenza $\sim$ associata alla partizione è una congruenza destra:
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
- Ne concludiamo che l’indice di $N_L$ è minore o uguale a quello di $\sim$ e pertanto è finito

**Dimostrazione 3. $\implies$ 1.:** $N_L$ ha indice finito $\implies L$ è regolare
- Per il teorema di Kleene basta costruire un automa che accetta $L$
- Sia $A$ l’automa ottenuto nel modo seguente (automa di Nerode di $L$)
	- Gli stati sono le classi dell’equivalenza $N_L$
	- La funzione di transizione $\delta$ è definita nel modo seguente: siano $q \in Q$ e $a \in \Sigma$. Poiché $q$ è una classe dell’equivalenza $N_L$ tutte le parole $ua$ con $u \in q$ appartengono a un’unica classe $p$ dell’equivalenza $N_L$. Poniamo dunque $\delta(q,a) = p$
	- Lo stato iniziale $q_0$ è la classe d’equivalenza della parola vuota
	- Gli stati finali sono le classi d’equivalenza incluse in $L$
- Per ogni $w \in \Sigma, \hat\delta(q_0,w)$ è la classe di equivalenza della parola $w$
- Poiché $L$ è unione di classi dell’equivalenza $N_L$ tale stato è finale se e solo se $w \in L$
- Se ne conclude che il linguaggio accettato da $A$ è proprio $L$

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260320112259.png" width="200">
</div>

**Proposizione:** Sia $L$ un linguaggio regolare. L’automa di Nerode di $L$ è un automa deterministico che accetta $L$ col minimo numero di stati possibile. 

**Dimostrazione:** 
- Sia $A$ l’automa deterministico che accetta $L$ col minimo numero di stati possibile e sia n tale numero; 
- nella dimostrazione 1. $\implies$ 2. si è visto che $L$ è unione di classi di una congruenza destra $\sim$ di indice n; 
- nella dimostrazione 2. $\implies$ 3. si è visto che l’indice di $N_L$ è minore o uguale di quello di $\sim$; 
- quindi l’indice di $N_L$ è minore o uguale a n; 
- nella dimostrazione 3. $\implies$ 1. si è costruito un automa che accetta $L$ con un numero di stati pari all’indice di $N_L$ ; 
- quindi l’indice di $N_L$ è maggiore o uguale a n; 
- ne segue l’asserto