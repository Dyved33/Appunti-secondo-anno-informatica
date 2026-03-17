###### Alberi binari
Le linked lists vanno bene per rappresentare delle relazioni lineari (non tutte le relazioni sono lineari). In particolare ci concentriamo sul problema di rappresentare un albero costruito a partire da una linked data structure. Ogni nodo dell'albero è un oggetto e supponiamo che ogni nodo dell'albero contenga un attributo chiave (key). L'altro attributo d'interesse sono i puntatori che connettono ogni nodo agli altri. Un **albero binario** è un grafo non orientato rappresentato con i puntatori, caratterizzato dal fatto che ogni nodo ha al massimo due figli, uno destro e uno sinistro:
 - Ogni nodo è definito con 3 puntatori, $p,\space left,\space right$, in cui $p$ è il puntatore al nodo padre; il nodo $left$ è il puntatore al figlio sinistro e il nodo $right$ è il puntatore al figlio destro.
 - Se il nodo $x$ non ha figli sinistri, il puntatore $left$ è inizializzato a $NIL$. Analogamente per il caso destro.
 - La radice ha il puntatore $p$ inizializzato a $NIL$.
 - La radice dell' intero albero $T$ è puntata dal parametro $T.root$.
 - Se $T.root = NIL$ allora l'intero albero è vuoto.
 - Le foglie dell'albero hanno i puntatori $left$ e $right$ inizializzati a $NIL$

<div style="display: flex; justify-content: center;">
  <img src="Screenshot From 2026-02-24 11-19-46.png" width="190">
</div>

È possibile estendere lo schema visto nell'albero binario a tutte le classi di alberi dove il numero di figlie di ogni nodo è al massimo un valore costante k. Per fare ciò è sufficiente sostituire i puntatori $left$ e $right$ con altri attributi come $child_1$, $child_2$, $...$ , $child_k$.

**Problemi principali:** Questo schema non funziona:
- In caso di variabile $k$ non limitata superiormente.
- Se $k$ è limitato da una costante molto grande ma la maggior parte dei nodi ha un numero basso di figli (spreco di memoria per i puntatori inutilizzati).

**Rappresentazione *Left-child*,*right-sibling:** Serve a generalizzare il concetto di albero, rappresenta un albero con un numero arbitrario di figli: 
- Ha la caratteristica si usare solo $O(n)$ spazio per ogni albero di $n$ nodi 
- Ogni nodo è definito con un puntatore $p$ al padre.
- Invece di avere un puntatore per ogni figlio, ogni nodo ha solo due puntatori
	- $left \_ child$ : figlio sinistro.
	- $right\_sibling$: fratello destro.
-  Se il nodo non ha figli allora $left\_child = NIL$
- Se il nodo non ha fratelli destri allora $right\_sibling = NIL$
###### Albero binario di ricerca
Si tratta di un albero binario nel quale le chiavi (cioè le informazioni contenute nei nodi) sono sempre contenute in modo da soddisfare le proprietà dell'albero.
**Proprietà dell'albero binario di ricerca:** Considerando un nodo $x$ all' interno di un *albero binario di ricerca:*
- se il nodo $y$  si trova nella parte sinistra del sotto-albero di $x$, allora $y.key \leq x.key$
- se il nodo $y$  si trova nella parte destra del sotto-albero di $x$, allora $y.key \geq x.key$

<div style="display: flex; justify-content: center;">
  <img src="Screenshot From 2026-02-24 11-25-52.png" width="190">
</div>

**In-order tree walk:** Data la 'proprietà degli alberi binari di ricerca, è possibile visualizzare tutte le chiavi dell'albero in maniera ordinata grazie a un semplice algoritmo ricorsivo.

~~~ 
inorderTreeWalk(x)
 if x!=NIl
    inorderTreeWalk(x.left)
    print x.key
    inorderTreeWalk(x.right)
~~~

L'algoritmo è così chiamato perché prima stampa tutti gli elementi di un sotto-albero partendo dalla sinistra della radice e concludendo alla sua destra.

**Pre-order tree walk:** Stampa la radice prima dei valori di ciascun sotto-albero

**Post-order tree walk:** Stampa la radice dopo i valori di ciascun sotto-albero

> [!info] Richiede $\Theta(n)$ in tempo per passare un albero binario di ricerca con n nodi. Dopo la chiamata iniziale la procedura chiama se stessa ricorsivamente 2 volte per ogni nodo nell'albero (una per ogni figlio)

Gli alberi binari di ricerca possono supportare gli algoritmi di richiesta di massimo, minimo, successore e predecessore esattamente come quello di ricerca. Queste operazioni richiedono un tempo $O(h)$ in un albero binario di ricerca di altezza h

**Ricerca:** Procedura di ricerca per un nodo con una chiave data in un albero binario di ricerca:

~~~
treeSearch(x,k)
 if x == NIL || k == x.key
    return x
 if k < x.key
    return treeSearch(x.left,k)
 else return treeSearch(x.right,k)
 
 iterativeTreeSearch(x,k)
  while x != NIL && k != x.key
   if k < x.key
    x = x.left
   else
    x = x.right
  return x
~~~

Dato un puntatore x alla radice di un sotto-albero e una chiave k, il metodo *treeSearch (x,k)* ritorna un puntatore al nodo in cui la chiave è contenuta. Se tale nodo non esiste il metodo ritorna $NIL$. Per fare una ricerca "globale" della chiave richiamo *treeSearch (T.root, k)*. Questa procedura inizia alla radice dell'albero e traccia un percorso scendendo lungo l'albero. Per ogni nodo che viene incontrato si confronta la chiave k con la *x.key* e se le due chiavi sono uguali la ricerca termina. Se la la *x.key* è più piccola di k la ricerca continua nel sotto-albero sinistro di x, altrimenti (*x.key* più grande di k) si continua nel sotto-albero destro di x.

**Massimo e minimo:** Per trovare il minimo in un albero di ricerca binario basta seguire i puntatori ai figli sinistri di ogni nodo dalla radice fino a $NIL$ (la procedura ritorna un puntatore al nodo contenente il minimo elemento). Per trovare il massimo si segue lo stesso principio scegliendo però il percorso fatto da sotto-alberi destri. 

~~~
treeMinimum(x)
 while x.left != NIL
  x = x.left
 return x

treeMaximum(x)
 while x.right != NIL
  x = x.right
 return x
~~~

Le proprietà degli alberi di ricerca binaria mi garantiscono la correttezza di questi metodi.

**Successore e predecessore:** Il successore di un nodo x è il nodo con la chiave più grande della *x.key*. Si definisce un successore il successivo nodo visitato in un inorder tree walk. La struttura degli alberi binari di ricerca mi permette di determinare il successore senza comparare le chiavi dei nodi. Tale procedura ritorna il puntatore al nodo successivo se esiste, altrimenti ritorna $NIL$.  

~~~
treeSuccessor(x)
 if x.right != NIl
  return treeMinimum(x.right)
 else
  y = x.p
  while y != NIL && x == y.right
   x = y
   y = y.p
 return y
~~~

Tale procedura presenta 2 casi:
 - il sotto-albero destro di x è non vuoto => il successore di x è il nodo più a sinistra nel sotto-albero destro di x.
 - Il sotto-albero destro di x è vuoto e x ha un successore y => y è il più piccolo "genitore" di x il cui figlio (di y) è a sua volta un "genitore" di x
 L'algoritmo per trovare il successore di un nodo richiede $O(h)$ tempo per un albero di altezza h in quanto segue o un *path-up* o un *path-down* dell'albero. Anche la procedura per trovare un predecessore, dato che è simmetrica a quella per trovare il successore, ha costo $O(h)$

**Inserimento:** Sia la procedura di inserimento che quella di eliminazione causano delle modifiche all'albero (per rispettare le proprietà degli alberi binari di ricerca). La procedura per inserire un nodo dentro ad un albero binario di ricerca richiede un albero *T* e un nodo *z* la cui chiave deve essere già riempita con *z.left* = $NIL$ e *z.right* = $NIL$. 

> [!info] N.B. la procedura modificherà *T* e alcune proprietà di *z*

~~~
treeInsert(T,z)
 x = T.root
 y = NIL
 while x != NIL
  y = x
  if z.key < x.key
   x = x.left
  else
   x = x.right 
 z.p = y
 if y == NIL
  T.root = z
 elseif z.key < y.key
  y.left = z
 else 
  y.right = z
~~~

La procedura inizia alla radice dell'albero e il puntatore x traccia un percorso discendente cercando un $NIL$ da rimpiazzare con il nodo z. Questo algoritmo mantiene il puntatore  come "genitore" di x. Il while fa si che i due puntatori (a x e a y) si spostino verso il basso andando a destra e a sinistra in base al confronto di *z.key* e *x.key*, finché x non diventa $NIL$ (questo $NIL$ occupa la posizione dove verrà posizionato il nodo z = $NIL$ è un attributo destro o sinistro del nodo che diventerà "genitore" di z, oppure è la *T.root* se *T* è vuoto). Anche la procedura di inserimento costa $O(h)$ in tempo per un albero di altezza h. 

**Eliminazione:** Esistono 3 possibili strategie per eliminare un nodo da un albero di ricerca binaria:
1. *z* non ha "figli" (al massimo ne ha uno destro), quindi quando lo rimuovo modifico il "genitore" di *z* sostituendo quest'ultimo con $NIL$
<div style="display: flex; justify-content: center;">
  <img src="Screenshot From 2026-02-25 11-55-49.png" width="190">
</div>
2. *z* ha un "figlio" (in tal caso l'unico "figlio" è sinistro), quindi pongo il "figlio" di *z* al posto di questo e modifico conseguentemente il "genitore" di *z* 
<div style="display: flex; justify-content: center;">
   <img src="Screenshot From 2026-02-25 11-56-25.png" width="190">
 </div>
3. *z* ha due "figli", quindi viene trovato il suo successore *y* (nel sotto-albero destro di *z*) e si sposta *y* affinché prenda il posto di *z*. Il resto del restante sotto-albero destro di *z* diventa il nuovo sotto-albero destro di *y* e il sotto-albero sinistro di *z* diventa il sotto-albero sinistro di *y*. *Y*, essendo il successore di *z*, non può avere un "figlio" sinistro e il "figlio" destro di *y* prende la posizione originaria di quest'ultimo (il resto del sotto-albero lo segue automaticamente). Di questo caso ci sono due possibili varianti:
    - *y* è il "figlio" destro di *z* 
<div style="display: flex; justify-content: center;">
  <img src="Screenshot From 2026-02-25 12-00-33.png" width="190">
 </div>
     - *y* si trova nel sotto-albero destro di *z*, ma non è il "figlio" destro
    <div style="display: flex; justify-content: center;">
  <img src="Screenshot From 2026-02-25 12-02-18.png" width="190">
</div>
La procedura di eliminazione di un nodo richiede come argomenti i puntatori a *T* e a *z*.

~~~
treeDelete(T,z)
 if z.left == NIL
  transplant(T,z,z.right)
 elseif z.right == NIL
  transplant(T,z,z.left)
 else
  y = treeMinimum(z.right)

  if y != z.right
   transplant(T,y,y.right)
   y.right = z.right
   y.right = y

  transplant(T,z,y)
  y.left = z.left
  y.left.p = y
~~~

Come parte del processo di eliminazione di un nodo, ho la necessità di spostare interi sotto-alberi dell'albero binario di ricerca. Posso fare ciò con la subroutine *trasnsplant*, la quale rimpiazza un sotto-albero "figlio" del proprio "genitore" con un altro sotto-albero. Nell'algoritmo quando rimpiazzo il sotto-albero con radice nel nodo *u* con il sotto-albero con radice nel nodo *v*, il "genitore" del nodo *u* diventa il genitore del nodo *v*.

~~~
transplant(T,u,v)
 if u.p == NIL
  T.root = v
 elseif u == u.p.left
  u.p.left = v
 else 
  u.p.right = v

 if v != NIL
  v.p = u.p
~~~

> [!info] N.B. Sia l'algoritmo di inserimento che quello di eliminazione richiedono tempo $O(h)$ per alberi di altezza h
###### Red-black trees
Il set di operazioni viste nella scorsa lezione sono veloci se l'altezza dell'albero di ricerca è piccola. Viceversa se l'altezza fosse un valore più grande le stesse operazioni sarebbero veloci tanto quanto una linked list. In questo contesto vengono introdotti gli alberi di ricerca bilanciati (**red-black trees**) che garantiscono delle operazioni che impiegano $O(log(n))$ di tempo nel caso peggiore.

**Proprietà:** 
-  Un red-black tree è un albero di ricerca binaria con un ulteriore informazione per ogni nodo: il colore del nodo, il quale può essere o rosso o nero.
-  Vincolando il colore del nodo ogni cammino semplice dalla radice alle foglie di un albero, ci si assicura che nessun cammino del genere sia più lungo del doppio di qualsiasi altro (dovrebbe voler dire che il cammino più lungo possibile è circa il doppio del cammino più corto possibile)
-  L'albero è approssimativamente  *bilanciato* (non ci sono foglie addossate a destra o a sinistra)
-  L'altezza di un red-black tree con n chiavi è al massimo $2*log(n+1)$ che ha costo $O(log(n))$

**Struttura:**
-  Ogni nodo dell'albero contiene gli attributi : colore, chiave, $left, right \ e\ p$
-  Se un figlio o un genitore di un nodo non esistono, il puntatore corrispondente (attributo del nodo) contiene il valore $NIL$ (questi valori possono essere pensati come i puntatori alle foglie)
- Un red-black tree è un albero di ricerca binaria che soddisfa le seguenti proprietà:
   1.  Ogni nodo o è nero o è rosso
   2. La radice è nera
   3. Ogni foglia ($NIL$) è nera
   4. Se un nodo è rosso $=>$ entrambi i suoi figli sono neri
   5. Per ogni nodo, ogni cammino semplice discendente dal nodo fino alle foglie contiene lo stesso numero di nodi neri
- Si usa una singola sentinella per rappresentare $NIL$.  La sentinella $T.NIL$ è un oggetto con gli stessi attributi di un normale nodo dell'albero. L'attributo colore della sentinella è il nero e gli altri attributi possono prendere valori arbitrari

<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Screenshot From 2026-03-02 18-40-20.png" style="width: 190%; border-radius: 8px;">
  </div>
  <div style="flex: 1.5;">
  Nel primo dei due alberi i numeri accanto ai nodi stanno ad indicare quanti nodi neri ci sono da quel nodo fino alle foglie (per rispettare la proprietà 5)
  </div>
</div>

La sentinella permette di trattare un figlio $NIL$ di un nodo *x*  come un nodo normale il cui padre sia *x*.  
Il numero di nodi neri in ogni cammino semplice discendente a partire (senza includere) un nodo *x*  fino alle foglie è denotato come *bh(x)* . A partire dalle proprietà sopra definite si può dire che un black-height è ben definito.  In generale un black-height di un red-black tree è un black-height della sua radice. Un red-black tree con *n* nodi interni ha altezza al massimo $2log(n+1)$

**Basic operations:** Le operazioni di base sono ricerca, minimo, massimo, successore, predecessore e richiedono tutti $O(log(n))$ di tempo. 

> [!info] Come gli alberi binari di ricerca nei casi di insert e di delete viene modificata la struttura dell'albero venendo meno alle proprietà di questi alberi. Per sistemare questo aspetto le proprietà del colore e del puntatore di alcuni nodi devono cambiare. La struttura dei puntatori cambia attraverso rotazione, che è un'operazione locale in un albero di ricerca che preserva le proprietà degli alberi di ricerca binaria. 

Esistono due tipi di rotazione: destra e sinistra. Entrambe richiedono $O(1)$ in tempo. Una rotazione sinistra di un nodo *x*, per esempio, richiede che il figlio destro del nodo (*y*) NON deve essere $T.NIL$. La rotazione cambia il sotto-albero originariamente con radice in *X* scambiando il collegamento tra *x* e *y* (come nell'immagine seguente).

<div style="display: flex; justify-content: center;">
  <img src="Screenshot From 2026-03-03 09-50-21.png" width="190">
</div>

**Inserimento:**
Per inserire un nodo in un red-black tree con n nodi interni in tempo $O(log(n))$ e mantenere le proprietà dell'albero è necessario modificare la procedura *treeInsert*. La procedura *RB_insert* inizia con l'inserimento del nodo *z* nell'albero *T*, come se fosse un normale albero di ricerca binaria, e colora questo nodo di rosso. Per garantire the le proprietà del red-black tree vengano preservate è necessario implementare la procedura *RB_insert_fixup* che ricolora i nodi e applica la rotazione.

```
RB_insert (T, z)
x = T.root
y = T.NIL
while x != T.NIL
	  y = x
	  if z.key < x.key
		   x = x.left
	  else
		   x = x.right 
 z.p = y
 
 if y == T.NIL
	  T.root = z
 elseif z.key < y.key
	  y.left = z
 else 
	  y.right = z
 z.left = T.NIL
 z.right = T.NIL
 z.color = RED
 RB_insert_fixup
```

**Differenze tra treeInsert e RB_insert:** 
- Tutte le istanze di $NIL$ in treeInsert sono rimpiazzate da $T.NIL$
- *z.left* e *z.right* sono importati a $T.NIL$ per mantenere la struttura dell'albero
- Alla fine il nodo *z* viene colorato di rosso
- Per evitare di violare le proprietà del red-black tree viene alla fine chiamato *RB_insert_fixup* (ogni nodo o è rosso o è nero, le foglie sono nero, il numero di nodi neri è lo stesso in ogni cammino semplice) => uniche proprietà che possono essere violate sono: radice deve essere nera e un nodo rosso non può avere un figlio rosso

> [!info] Ci riferiamo allo zio come il nodo che è il figlio (di qualsiasi colore) di un nodo padre (nero)

1. *y*, lo zio di *z*, è rosso: quando sia *z.p* sia *y* sono rossi. Poiché *z.p.p* è nero, questa proprietà può essere trasferita di un livello in basso sia a *z.p* sia a *y* così da sistemare il problema di avere sia *z* che *z.p* rossi
> [!example] ![[Screenshot From 2026-03-03 10-47-51.png|190]]
2. *y*, lo zio di *z*, è nero: 
   - *y* è nero e *z* è un figlio destro -> una rotazione sinistra su *z* trasforma la situazione nel prossimo caso
> [!example] ![[Screenshot From 2026-03-03 10-48-53 1.png|190]]

   - *y* è nero e *z* è un figlio sinistro -> cambia alcuni colori e applica una rotazione destra su *z.p.p*, il quale preserva la proprietà di non avere figli rossi (essendo un nodo rosso)
> [!example] ![[Screenshot From 2026-03-03 10-50-24.png|190]]
###### Hash table
Una hash table è una generalizzazione di un array. In un array gli elementi vengono immagazzinati con una chiave k alla posizione k. Data la chiave k, per trovare l'elemento la cui chiave è k, mi basta guardare alla k-esima posizione dell'array (*direct addressing*). Quando non voglio allocare un array con una posizione per ogni possibile chiave uso una hash table. Uso questa struttura anche nel caso in cui il numero di chiavi immagazzinate è più piccolo del numero delle possibili chiavi. Quando si parla di hash table ci si riferisce ad un array che tipicamente usa una dimensione proporzionale al numero di chiavi che vanno immagazzinate. Per questo motivo se ho la chiave k non uso questo valore come indice nell'array, invece definisco una *hash function* (funzione di k) e usiamo il suo valore per indicizzare l'array.

| Direct-address tables                                                        | Hash tables                                                                                                                                     |
| ---------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------- |
| Mantiene un set dinamico                                                     | Richiede meno spazio di un array dato che può succedere che l'insieme di chiavi usate sia più piccolo dell'insieme di tutte le chiavi possibili |
| Ogni elemento ha una chiave presa dall'insieme universo $U = \{0...m-1\}$    | Può ridurre la richiesta di storage fino a $\Theta(\|K\|)$, con K l'insieme di chiavi effettivamente usate                                      |
| Non esistono due elementi con la stessa chiave                               | Può ancora richiedere (nel caso medio, non nel peggiore) $O(1)$ in tempo per la ricerca                                                         |
| Se c'è un elemento *x* con chiave *k*  => *T[k]* contiene il puntatore a *x* | Si usa una funzione *h* per immagazzinare l'elemento nello slot $h(k)$                                                                          |
> [!info]  
> - chiamiamo *h* la **hash function** e *T* la **hash table**
> - $h: U \rightarrow \{0,...,m-1\}$, così $h(k)$ è uno slot in *T*

> [!warning] Quando due o più chiavi puntano allo stesso slot si ha una collisione:
> - può accadere quando ci sono più chiavi che slot $(|U| > m)$
> - dato un set K di chiavi con $|K| < m$ potrebbe come potrebbe non accadere (accade sicuramente se $|K| > m$)

**Hashing indipendente e uniforme:** Idealmente $h(k)$ verrebbe scelto casualmente e indipendentemente da $\{0,...,m-1\}$. Una volta scelto $h(k)$ ogni valutazione seguente di $h(k)$ deve ridare lo stesso risultato. L'idea che la hash function sia indipendente, uniforme o casuale è teoretico e non può essere implementato nella pratica (si usa per analizzare il comportamento della hash function e per trovare delle approssimazioni pratiche alla forma ideale)

<div style="display: flex; justify-content: center;">
  <img src="Screenshot From 2026-03-03 23-14-22.png" width="400">
</div>

**Collision resolution by chaining:** Come un divide-et-impera non ricorsivo posso usare la hash function per dividere gli n elementi in maniera casuale in m sottoinsiemi, ognuno dei quali con circa $\frac{n}{m}$ elementi. Successivamente posso gestire ciascun sottoinsieme in maniera indipendente come una linked list. Nei vari casi delle operazioni di inserimento, ricerca e cancellazione abbiamo:
- *insertion* = il caso peggiore richiede $O(1)$ in tempo
- *search* = il caso peggiore richiede una quantità di tempo proporzionale alla lunghezza della lista di elementi in ogni slot
- *deletion* = dato un puntatore x all'elemento da cancellare, in questo modo non devo cercare tale elemento

Data una chiave quanto ci metto a trovare l'elemento con quella chiave o per determinare che non c'è nessun elemento con quella chiave? Posso sfruttare un'analisi in termine di *fattore di carico*: $\alpha = \frac{n}{m}$ con:
- *n* = numero di elementi nella tabella
- *m* = numero di slot nella tabella = numero di linked lists (possibilmente vuote)
- Il fattore di carico è il numero medio di elementi per ogni linked list
- Posso avere $\alpha < 1, \alpha = 1, \alpha > 1$

Il caso medio dipende da quanto bene la hash function distribuisce le chiavi tra gli slot. Il caso peggiore, invece, si ha quando tutte le *n* chiavi sono inserite nello stesso slot $\implies$ per esempio il caso peggiore della ricerca è $\Theta(n)$. Se mi concentro sulla performance del caso medio dell'hashing with chaining:
- Assumo *hashing indipendente e uniforme* = ogni elemento dato può essere distribuito in ognuno degli *m* slot
- L'*hashing uniforme e indipendente* è universale = la probabilità che ogni 2 chiavi distinte collidano è $\frac{1}{m}$
- $j = 0,1,...,m-1$, si denota la lunghezza della lista $T[j]$ da $n_j$ , in questo modo $n = n_0 + n_1 + ... + n_{m-1}$
- Il valore atteso di $n_j$ è $E[n_j] = \alpha = \frac{n}{m}$
- Assumo che la hash function richieda $O(1)$ in tempo per essere eseguita, in questo modo il tempo richiesto per cercare l'elemento di chiave *k* dipende dalla lunghezza $n_{h(k)}$ della lista $T[h(k)]$

**Teorema:** una ricerca senza successo richiede nel caso medio $\Theta(1+\alpha)$ in tempo (dimostrazione si basa sul fatto che una ricerca senza successo arriva in fondo alla lista di lunghezza attesa $\alpha$ e a questo bisogna aggiungere il tempo per eseguire la hash function). Anche una ricerca con successo nel caso medio richiede $\Theta(1+\alpha)$ in tempo (**Teorema**). 

**Cosa rende buona una hash function?** Idealmente la hash function soddisfa l'*hashing indipendente e uniforme:* ogni chiave ha la stessa probabilità delle altre chiavi di essere distribuita in ognuno dei *m* slot. Nella pratica, però, non è possibile soddisfare tale assunto siccome non si conosce in anticipo la distribuzione di probabilità che le chiavi siano estratte, e le chiavi potrebbero non essere estratte indipendentemente (se si conosce la distribuzione delle chiavi si può avere un vantaggio da ciò)
Nella pratica le hash function assumono che le chiavi siano alternativamente:
- Un interno non negativo e corto che rientri in una *machine word* (32 o 64 bit di solito)
- Un vettore corto di interi non negativi, ciascuno dei quali rientra in un limite definito

**Static hashing:** Si tratta di una hash function fissa per la quale la casualità viene solo dalla "speranza per la distribuzione" delle chiavi
- *Division method:* $h(k) = k \ \text{mod} \ m$. Il vantaggio risiede nel fatto che è veloce e richiede solo un'operazione di divisione. Una buona scelta per $m$ è un numero primo non troppo vicino ad una potenza esatta di 2
> [!example] $m$ = 20 e $k$ = 91 $\implies$ $h(k) = 91 \ \text{mod} \ 20 = 11$
- *Multiplication method:* si sceglie la costante $A \ \text{nel range} \ 0 < A < 1 \ \rightarrow$ moltiplico la chiave $k \ \text{per} \ A \ \rightarrow$ estraggo la parte frazionaria da $k*A \rightarrow$ moltiplico la parte frazionaria per $m \rightarrow$  prendo il *floor* del risultato. Questo metodo è più lento del *division method*, ma il valore di $m$ non è critico (può essere scelto indipendentemente da A)

> [!info] $h(k) = \lfloor m(k*A \ \text{mod} \ 1)\rfloor$, dove $k*A \ \text{mod} \ 1 = k*A - \lfloor k*A \rfloor = \ \text{parte frazionaria di} \ k*A$ 

**Hashing casuale:** Un ipotetico avversario malevolo, conoscendo la funzione hash, potrebbe scegliere chiavi che collidano tutte nello stesso slot, producendo il caso pessimo. Qualsiasi funzione hash statica è vulnerabile a questo attacco.
La soluzione sta nello scegliere la funzione hash **casualmente**, indipendentemente dalle chiavi. Un caso particolare è l’**hashing universale**, che garantisce buone prestazioni medie quando si usa il chaining, indipendentemente dalle chiavi.
Si consideri una collezione finita $\mathscr{H}$  di funzioni hash che mappi l'universo $U$ di chiavi nel range $\{0,1,...,m-1\}$. $\mathscr{H}$ è universale  se per ogni coppia di chiavi $k_1, k_2 \in U$, il numero di funzioni hash $h \in \mathscr{H}$ con $h(k_1) = h(k_2) \leq \frac{|\mathscr{H}|}{m}$. In altre parole $\mathscr{H}$ è universale se, con una funzione hash *h* scelta casualmente da $\mathscr{H}$, la probabilità di una collisione tra due chiavi distinte non supera l'$\frac{1}{m}$ di probabilità di scegliere due slot casuali e indipendenti

**Proprietà hashing casuale:** Si possono ottenere limiti teorici sulle collisioni e sulle prestazioni attese. Considero una famiglia di hash function $\mathscr{H}$ su un dominio $U$ e con range $\{0,1,...,m-1\}$, le chiavi appartengono a $U$,i numeri degli slot sono in $\{0,1,...,m-1\}$, e una hash function *h* è presa casualmente da $\mathscr{H}$:
- Uniforme: la probabilità di prendere *h* tale che $h(k) = q$ è $\frac{1}{m}$
- Universale: per ogni chiave distinta $k_1, k_2$, la probabilità che $h(k_1) = h(k_2)$ è al massimo $\frac{1}{m}$
- $\epsilon$-universale: per ogni chiave distinta $k_1, k_2$, la probabilità che $h(k_1) = h(k_2)$ è al massimo $\epsilon$
- $d$-indipendente: per ogni chiave distinta $k_1, k_2$ e ogni slot $q_1,q_2,...q_d$, non necessariamente distinto, la probabilità che $h(k_i) = q_i$ è $\frac{1}{m^d}$ 

**Hashing basato su teoria dei numeri:** Alcune famiglie di funzioni hash utilizzano proprietà della teoria dei numeri per garantire buone distribuzioni. Scelto un numero primo *p* abbastanza grande da permettere che ogni possibile chiave appartenga a $\{0,1,...,p-1\}$ e si assuma che $p > m$. Si denoti $\mathbb{Z}_p = \{0,1,...,p-1\}$, $\mathbb{Z^*}_p = \{0,1,...,p-1\}$. Dato ogni $a \in \mathbb{Z^*}_p$ e ogni $b \in \mathbb{Z}_p$, si definisce $h_{ab}(k) = ((ak + b) \ \text{mod} \ p) \ \text{mod} \ m$
Dato *p* e *m*, la famiglia di hash function è $\mathscr{H}_{ab} = \{h_{ab} : a \in \mathbb{Z^*}_p, b \in \mathbb{Z}_p$ (ognuna delle quali mappa da $\mathbb{Z}_p$ a $\mathbb{Z^*}_p$). Questa famiglia di hash functions contiene $p(p-1)$ funzioni, una per ogni combinazioni di *a* e *b*.

**Teorema:** la famiglia $\mathscr{H}_{p m}$ definita sopra di hash functions è universale

**Open addressing:** Idea principale:
- Tutti gli elementi sono memorizzati **direttamente nella tabella**.   
- Ogni slot contiene un elemento oppure NIL.
- La tabella può riempirsi, ma il fattore di carico non può superare 1.

**Inserimento:** 
 1. Calcolare la prima posizione scelta.
 2. Se è libera, inserire l’elemento. 
 3. Se è occupata, calcolare la seconda posizione. 
 4. Continuare finché si trova uno slot libero.
Ogni elemento ha una propria sequenza di preferenze.

```
Hash-insert(T,k)
i = 0
repeat
	 q = h(k,i)
	 if Y[q] == NIL
		  T[q] = k
		  return q
	 else
		  i = i+1
until i == m
error "hash table overflow"   
```

**Ricerca:** Si segue la stessa sequenza di probe dell’inserimento, ma se si incontra uno slot vuoto, si conclude che l’elemento non è presente. Il vantaggio di tale opzione è che non servono puntatori (a differenza del chaining) e che lo spazio risparmiato può essere usato per aumentare la dimensione della tabella (slide 30 pdf 11 per avere più specifiche nella ricerca della chiave k).
```
Hah-Search(T,k)
i =0
repeat
	 q = h(k,i)
	 if T[q] == k
		  return q
	 i = i+1
until T[q] == NIL or i == m
return NIL  
```

**Cancellazione:** La cancellazione è più delicata rispetto al chaining. Non si può semplicemente mettere NIL, perché si interromperebbero le <u>sequenze di probe</u> = è l'ordine dei bucket (posizioni) controllati nella tabella hash fino a quando l'elemento viene trovato o viene trovata una posizione vuota (un probe o sonda è un tentativo di inserire o  cercare un elemento in una specifica posizione della tabella). Si usano marcatori speciali (slot DELETED) al posto di *NIL* per marcare uno slot come vuoto durante la cancellazione.
- Search dovrebbe trattare DELETED come se lo slot contenesse una chiave che non corrisponde a quella che viene cercata
- Insertion dovrebbe trattare DELETED come se lo slot fosse vuoto, così da poter essere riusato

Lo svantaggio di usare DELETED è che ora il tempo di ricerca non dipende più dal fattore di carico $\alpha$. Un caso speciale di open addressing, chiamato linear probing, evita di avere slot marcati con DELETED

**Sequenze di Probe:** La situazione ideale è quella di *independent uniform permutation hashing*: ogni chiave può avere ognuna delle $m!$ permutazioni di $\langle 0,1,...,m-1 \rangle$. Siccome è complesso implementare una vera independent uniform permutation hashing, posso approssimarla nei modi seguenti:
- **Double hashing:** Si usano due funzioni hash ausiliari $h_1, h_2$ e probe number i: $h(k,i)=(h_1(k)+i*h_2(k)) \ \text{mod} \ m$. Il primo probe va nello slot $h_1(k)$, successivamente i probe sono compensati dallo slot precedente da $h_2(k)$ mod $m$ $\implies$ la probe sequence dipende dalle chiavi in 2 modi. È necessario che $h_2(k)$ sia relativamente primo di *m* (nessun fattore comune a parte 1) per garantire che la probe sequence sia una permutazione di $\langle 0,1,...m-1 \rangle$. Questa soluzione riduce clustering (formazione di gruppi di elementi consecutivi nella tabella) e migliora la distribuzione.
  - Si può scegliere *m* affinché sia una potenza di 2 e $h_2$ affinché produca sempre un numero dispari
  - Si può fare in modo che *m* sia primo e avere $1 < h_2(k) < m$
- **Linear probing:** È un caso speciale del double hashing. Data una funzione hash ausiliaria $h'$, la probe sequence inizia allo slot $h'(k)$ e continua attraverso la tabella (da slot m-1 allo slot 0). Data una chiave k e un probe number $i, (0 \leq i \leq m)$: $h(k,i) = (h'(k) + i) \ \text{mod} \ m$. Il probe iniziale determina tutta la sequenza $\implies$ ogni *m* sequenza possibile. Semplice ma soffre di **primary clustering**.

**Analisi dell'open addressing:** Un 'analisi in termini di fattore di carico $\alpha$ implica che:
- Si assuma che la tabella non si riempe mai, quindi si ha sempre $0 \leq n < m \implies 0 \leq \alpha < 1$ 
- Si assuma un independent uniform permutation hashing
- Nessuna cancellazione
- In una ricerca con successo, ogni chiave è ugualmente probabile che sia cercata 

**Teorema:** il numero atteso di probe in una ricerca senza successo è al massimo $\frac{1}{(1-\alpha)}$

**Corollario:** il numero atteso di probe da inserire è al massimo $\frac{1}{(1-\alpha)}$

**Teorema:** il numero atteso di probe in una ricerca con successo è al massimo $\frac{1}{\alpha}*\ln(\frac{1}{1-\alpha})$
(su pdf ci sono le dimostrazioni di tutti e 3)

**Considerazioni pratiche:** Le CPU moderne influenzano il comportamento dell’hashing.
1. Gerarchie di memoria:
 - Le cache sono piccole e veloci.
 - Organizzate in cache line (es. 64 byte contigui).
 - È più veloce riusare una cache line che accedere alla memoria principale.
 - Iterare su dati contigui è vantaggioso.
2. Set di istruzioni avanzati: Alcune istruzioni hardware (usate per crittografia) possono essere sfruttate per calcolare hash efficienti.
3. Linear Probing: Nel modello RAM teorico il linear probing può performare male. In presenza di cache gerarchiche:
 - Funziona meglio perché probe consecutivi tendono a restare nella stessa cache line.

> [!info] Sintesi:
> Hashing diventa una funzione in due variabili $h(k,p)$ dove p indica il probe, ovvero il "tentativo" di inserimento nella hash table (in caso di collisione)
> 
> Hashing lineare: 
> Indice = $[h(k)+p] \ \text{mod} \ m$
> Facilmente "hackerabile" perché rende molto probabile che se c'è un blocco di elementi allocati il prossimo cada sul primo blocco libero
> 
> Hashing con moltiplicazione e somma: comunque facilmente hackerabili (ad es, m=16 se inserisco 26 + $k*m$ ottengo comunque collisioni sull'indice 10)
> 
> Hashing quadratico: meno prevedibile
###### Grafi
Dato un grafo $G = (V,E)$, in pseudocodice si rappresenta il set di vertici con $G.V$ e il set degli archi con $G.E$.
- $G$ può essere orientato o non orientato
- Due modi comuni per rappresentare i grafi per gli algoritmi sono: 
 1. *Liste di adiacenza*: array *Adj* di $|V|$ liste, una per vertice. La lista del vertice *u* ha tutti i vertici *v* in modo da avere che $(u,v) \in E$ (vale sia per grafi orientati che non orientati). In pseudocodice si denota tale array come un attributo (G.*Adj*)
 <div style="display: flex; justify-content: center;">
   <img src="Screenshot From 2026-03-11 14-55-01.png" width="300">
 </div>
   - Spazio: $\Theta(V+E)$
   - Tempo: per listare tutti i vertici adiacenti a *u*: $\Theta(\deg(u))$
   - Tempo: per determinare se $(u,v) \in E$: $O(\deg(u))$
 2. *Matrici di adiacenza*: $|V| \times |V|$, matrice $A = (a_{ij})$ con $a_{ij} =\begin{cases} 1 &\text{if} \ (i,j) \in E \\ 0 &\text{altrimenti} \end{cases}$. Queste richiedono:
   - Spazio: $\Theta(V^2)$
   - Tempo: per listare tutti i vertici adiacenti a *u*: $\Theta(V)$
   - Tempo: per determinare se $(u,v) \in E : \Theta(1)$
Per quanto riguarda l'espressione del tempo di esecuzione di un algoritmo, è spesso in termini sia di $|V|$ sia di $|E|$ (cardinalità).

**Breadth-first search:**
- In input ho un grafo $G = (V,E)$, sia orientato che non orientato, e un vertice sorgente $s \in V$. 
- In output ho:
	- $v.d$ = distanza (numero di archi più piccolo) da $s$ a $v$ per ogni $v \in V$
	- $v.\pi$ = è il predecessore di $v$ nel cammino più breve da $s$. $(u,v)$ è l'ultimo arco nel cammino più breve $s \rightsquigarrow v$. Il sotto-grafo predecessore contiene gli archi $(u,v)$ tali che $v.\pi = u$. Il sotto-grafo predecessore forma un albero chiamato *breadth-first tree*.

*Intuizione*: i breadth-first tree espandono la frontiera tra i vertici scoperti e non scoperti uniformemente attraverso l'ampiezza della frontiera. A partire da $s$ vengono scoperti i vertici in ondate
- Prima visita tutti i vertici distanti 1 arco da $s$
- Da qui visita tutti i vertici distanti 2 archi da $s$
- Etc.

Per mantenere il "fronte d'onda" si usa il metodo FIFO sulla coda $Q$:
- $v \in Q \iff$ l'onda ha visitato $v$ ma non è ancora uscito da $v$ 
- $Q$ contiene i vertici a distanza $k$, e possibilmente alcuni vertici a una distanza $k+1$. Perciò in ogni momento $Q$ contiene porzioni di due ondate consecutive

```text
BFS(G,s)
 for each vertex u ∈ G.V − {s}
	  u.color = WHITE
	  u.d = ∞
	  u.π = NIL
 s.color = GRAY
 s.d = 0
 s.π = NIL
 Q = ∅
 ENQUEUE(Q,s)
 while Q ≠ ∅
	  u = DEQUEUE(Q)
	  for each vertex v ∈ G.Adj[u] // cerca i vicini di u
		   if v.color == WHITE // v è stato scoperto?
			    v.color = GRAY
			    v.d = u.d + 1
			    v.π = u
			    ENQUEUE(Q,v) // v è sulla frontiera ora
  u.color = BLACK // u è dietro la frontiera ora
```

<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Screenshot From 2026-03-11 15-53-03.png" style="width: 190%; border-radius: 8px;">
  </div>
  <div style="flex: 1.5;">
  <li> gli archi più spessi sono nel sotto-grafo predecessore
  <li> le linee tratteggiate vanno ai vertici appena spcoperti
  <li> i vertici cerchiati 2 volte sono stati scoperti e sono in $Q$ in attesa di essere visitati
  <li> i vertici cerchiati con una linea più spessa sono stati scoperti, tolti dalla coda e visitati
  <li> BFS potrebbe non raggiungere tutti i vertici e richiede in tempo $O(V+E)$
  <li> Tempo = $O(V+E)$
  <li> $O (V)$ perché ogni vertice è messo in coda al massimo una volta
  <li> $O(E)$ perché ogni vertice è estratto dalla coda al massimo una volta e l'arco $(u,v)$ è esaminato solo quando $u$ è estratto dalla coda. In questo modo ogni arco viene esaminato al massimo 1 volta se orientato, al massimo 2 volte se non orientato
  </div>
</div>

Per stampare i vertici del cammino più corto da $s$ a $v$:
```Pseudocode
if v == s
	print s
elseif v.parent == NIL
	print "no path from" s "to" v "exists"
else PRINT-PATH(G,s,v.parent)
	print v
```

**Depth-first search:** 
- *Input:* un grafo $G = (V,E)$ sia ordinato che non ordinato. Non viene dato nessun vertice sorgente
- *Output:*
	- 2 marcatori del tempo su ogni vertice utili per algoritmi che vedremo poi:
		- $v.d =$ momento della scoperta del vertice 
		- $v.f =$ tempo finale
	- $v.\pi$ è il predecessore di $v$ nella **depth-first forest** con più di un **detpth-first tree**. Se $u = u.\pi$ allora $(u,v)$ è un ramo dell’albero

Questo algoritmo esplora ogni arco e inizia da vertici diversi se necessario. Appena viene scoperto un vertice inizia ad esplorare da questo (diversamente dal BFS che mette un vertice in coda così che venga esplorato successivamente). Man mano che il DFS progredisce, ogni vertice ottiene un colore:
- Bianco = vertice ancora non scoperto
- Grigio = vertice scoperto ma non finito (non si è finito di esplorare da quel vertice)
- Nero = vertice finito (è stato trovato tutto ciò che era raggiungibile da tale vertice)

Tempi di scoperta e di fine:
- Interi unici da 1 a $2|V|$
- Per ogni $v,v.d < v.f$

In altre parole $1 \leq v.d < v.f \leq 2|V|$

```Pseudocode
DFS(G)
	for each vertex u in G.V
		u.color = WHITE
		u.parent = NIL
	time = 0
	for each vertex u in G.V
		if u.color == WHITE
			DFS-Visit(G,u)
la ricerca globale fa partire una ricerca locale su ogni vertice per esplorare tutto il grafo
```

```Pseudocode
DFS-Visit
	time = time +1
	u.d = time
	u.color = GRAY
	for each vertex v in G.Adj[u]
		if v.color == WHITE
			v.parent = u
			DFS-Visit(G,v)
	time = time +1
	u.f = time
	u.color = BLACK
```

> [!example] ![[Screenshot From 2026-03-16 19-15-31.png|200]]
> Tempo = $\Theta(V+E)$
> - l’analisi è simile a quella per il BFS
> - $\Theta$ e non solo $O$ siccome è garantito che venga esaminato ogni vertice e ogni arco

Ogni depth-first tree è fatto da archi $(u, v)$ così che $u$ sia grigio e $v$ sia bianco quando $(u, v)$ è esplorato

**Teorema (delle parentesi):** per ogni $u,v$ vale esattamente una delle seguenti:
1. $u.d < u.f <v.d < v.f$ oppure $v.d < v.f < u.d < u.f$ (cioè gli intervalli $[u.d, u.f]$ e $[v.d,v.f]$ sono disgiunti) e ne $u$, ne $v$ è discendente dall’altro
2. $u.d < v.d <v.f < u.f$ e $v$ è discendente da $u$ ($v$ è scoperto dopo e finito prima di $u$)
3. $v.d < u.d < u.f < v.f$ e $u$ è discendente da $v$ ($u$ è scoperto dopo e finito prima di $v$)

Cioè $u.d < v.d <u.f < v.f$ non può accadere ($v$ è sia scoperto che finito dopo $u$)

**Classificazione degli archi:**
- **Tree edge:** dentro la depth-first forest, trovato esplorando $(u,v)$
- **Back edge:** $(u,v)$ con $u$ discendente da $v$
- **Forward edge:** $(u,v)$ con $v$ discendente da $u$, ma non è un tree edge
- **Cross edge:** ogni altro arco. Può andare in mezzo ai vertici dello stesso depth-first tree o in depth-first tree diversi

In un grafo non orientato, ci possono essere delle ambiguità siccome $(u,v)$ e $(v,u)$ sono lo stesso arco. Si classifica in base al primo tipo sopra descritto corrispondente

**Topological sort:** Un grafo orientato e senza cicli è buono per modellare i processi e le strutture che hanno un ordine particolare:
- $a > b$ e $b > c$ $\implies$ $a > c$
- Ma potrebbe avere $a$ e $b$ tali che ne $a > b$ ne $b > c$

Si può sempre ottenere un ordine totale da un ordine parziale (o $a>b$ o $b>a$ per ogni $a\neq b$). Difatti questo è ciò che fa un topological sort.

**Dag:** è un grafo orientato e aciclico. Un grafo $G$ è aciclico $\iff$ un DFS di $G$ non produce un arco nero (prova alla slide 21)

**Topological sort di un dag:** ordina linearmente i vertici così che se $(u,v) \in E$ allora $u$ appare da qualche parte prima di $v$:
1. Chiama DFS (G) per valutare i tempi di fine $v.f$ per ogni vertice $v$
2. Quando ogni vertice è finito, lo inserisce all’inizio della linked list
3. Ritorna la linked list dei vertici

Non necessita di ordinare in funzione dei tempi di fine:
- Può far uscire i vertici appena sono finiti e capisce che voglio il *reverse* della lista
- Oppure mette i vertici all’inizio della linked list appena sono finiti. Appena finito la lista contiene i vertici ordinati con ordine topologico 

**Tempo:** $\Theta(V+E)$
> [!example] ![[Screenshot From 2026-03-17 23-15-57.png|200]]

**Correttezza:** è necessario mostrare che se $(u,v) \in E$ allora $v.f < u.f$. Quando l’arco $(u,v)$ è esplorato, di che colore sono $u$ e $v$?
- $u$ è grigio
- Anche $v$ è grigio? No, perché in tal caso $v$ sarebbe “antenato” di $u$ $\implies (u,v)$ è un arco nero $\implies$ contraddizione del lemma precedente (i dag non hanno archi neri) 
- $v$ è bianco? In tal caso diventa “discendente” di $u$. In base al teorema delle parentesi $u.d < v.d<v.f<u.f$
- $v$ è nero? In tal caso $v$ sarebbe già finito. Siccome sto esplorando $(u,v)$. $u$ non è ancora finito, quindi $v.f<u.f$

**Componenti fortemente connesse:** dato un grafo ordinato $G=(V,E)$, una componente fortemente connessa (*SCC*) di $G$ è il set massimale di vertici $C \subseteq V$ così che per ogni $u, v \in C$, $u \rightsquigarrow v$ e $v \rightsquigarrow u$ 
> [!example] ![[Screenshot From 2026-03-17 23-27-30.png|200]]

*Nell’esempio:* L’algoritmo usa $G^T$=trasposta di $G$:
- $G^T = (V,E^T), E^T = \{(u,v) : (v,u) \in E\}$
- $G^T$ è $G$ con tutti gli archi invertiti
> [!info] posso creare $G^T$ in $\Theta(V+E)$ di tempo se uso le liste di adiacenza

**Osservazione:** $G$ e $G^T$ hanno gli stessi SCC ($U$ e $v$ sono raggiungibili l’uno dall’altro in $G$ se e solo so sono raggiungibili l’uno dall’altro in $G^T$)
- $G^{SCC} = (V^{SCC},E^{SCC})$ 
- $V^{SCC}$ ha un vertice per ogni SCC in $G$
- $E^{SCC}$ ha un arco se c’è un arco tra  i corrispondenti SCC in $G$

<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Screenshot From 2026-03-17 23-35-16.png" style="width: 190%; border-radius: 8px;">
  </div>
  <div style="flex: 1.5;">
Nel nostro esempio
  </div>
</div>

**Lemma:** $G^{SCC}$ è un dag. Più formalmente, siano $C$ e $C'$ due SCC distinti in $G$, siano $u,v \in C$, $u',v' \in C'$ e si supponga che ci sia un cammino $u \rightsquigarrow u'$  in $G$. Allora non ci potrà essere un cammino $v' \rightsquigarrow v$  in $G$

```Pseudocode
SCC(G)
	chiama DFS(G) per calcolare i tempi di fine u.f per ogni vertice u
	crea G^T
	chiama DFS(G^T), ma nel loop principale considera i vertici in ordine decrescente di u.f (come valutato nel primo DFS)
	restituisce i vertici in ogni albero del depth-first tree formato  nel secondo DFS come un SCC separato
(?)esegue DFS in G (?)
```