### 26-02-23 lz.1
**Alberi binari**
Le linked lists vanno bene per rappresentare delle relazioni lineari (non tutte le relazioni sono lineari). In particolare ci concentriamo sul problema di rappresentare un albero costruito a partire da una linked data structure. Ogni nodo dell'albero è un oggetto e supponiamo che ogni nodo dell'albero contenga un attributo chiave (key). L'altro attributo d'interesse sono i puntatori che connettono ogni nodo agli altri. Un albero binario è un grafo non orientato rappresentato con i puntatori, caratterizzato dal fatto che ogni nodo ha al massimo due figli, uno destro e uno sinistro:
 - Ogni nodo è definito con 3 puntatori, $p,\space left,\space right$, in cui $p$ è il puntatore al nodo padre; il nodo $left$ è il puntatore al figlio sinistro e il nodo $right$ è il puntatore al figlio destro.
 - Se il nodo $x$ non ha figli sinistri, il puntatore $left$ è inizializzato a $NIL$. Analogamente per il caso destro.
 - La radice ha il puntatore $p$ inizializzato a $NIL$.
 - La radice dell' intero albero $T$ è puntata dal parametro $T.root$.
 - Se $T.root = NIL$ allora l'intero albero è vuoto.
 - Le foglie dell'albero hanno i puntatori $left$ e $right$ inizializzati a $NIL$
                  ![[Screenshot From 2026-02-24 11-19-46.png]]

**Alberi di classe maggiore di 2**
È possibile estendere lo schema visto nell'albero binario a tutte le classi di alberi dove il numero di figlie di ogni nodo è al massimo un valore costante k. Per fare ciò è sufficiente sostituire i puntatori $left$ e $right$ con altri attributi come $child_1$, $child_2$, $...$ , $child_k$.
 
**Problemi principali**
Questo schema non funziona:
- In caso di variabile $k$ non limitata superiormente.
- Se $k$ è limitato da una costante molto grande ma la maggior parte dei nodi ha un numero basso di figli (spreco di memoria per i puntatori inutilizzati).

**Rappresentazione *Left-child*,*right-sibling***
Serve a generalizzare il concetto di albero, rappresenta un albero con un numero arbitrario di figli: 
- Ha la caratteristica si usare solo $O(n)$ spazio per ogni albero di $n$ nodi 
- Ogni nodo è definito con un puntatore $p$ al padre.
- Invece di avere un puntatore per ogni figlio, ogni nodo ha solo due puntatori
	- $left \_ child$ : figlio sinistro.
	- $right\_sibling$: fratello destro.
-  Se il nodo non ha figli allora $left\_child = NIL$
- Se il nodo non ha fratelli destri allora $right\_sibling = NIL$

**Albero binario di ricerca**
Si tratta di un albero binario nel quele le chiavi (cioè le informazioni contenute nei nodi) sono sempre contenute in modo da soddisfare le proprietà dell'albero.

**Proprietà dell'albero binario di ricerca**
Considerando un nodo $x$ all' interno di un *albero binario di ricerca:*
- se il nodo $y$  si trova nella parte sinistra del sotto-albero di $x$, allora $y.key \leq x.key$
- se il nodo $y$  si trova nella parte destra del sotto-albero di $x$, allora $y.key \geq x.key$
                    ![[Screenshot From 2026-02-24 11-25-52.png]]

**In-order tree walk**
Data la 'proprietà degli alberi binari di ricerca, è possibile visualizzare tutte le chiavi dell'albero in maniera ordinata grazie a un semplice algoritmo ricorsivo.

~~~ 
inorderTreeWalk(x)
 if x!=NIl
    inorderTreeWalk(x.left)
    print x.key
    inorderTreeWalk(x.right)
~~~

L'algoritmo è così chiamato perché prima stampa tutti gli elementi di un sotto-albero partendo dalla sinistra della radice e concludendo alla sua destra.

---
### 26-02-24 lz.2
**Pre-order tree walk**
Stampa la radice prima dei valori di ciascun sotto-albero

**Post-order tree walk**
Stampa la radice dopo i valori di ciascun sotto-albero

Richiede $\Theta(n)$ in tempo per passare un albero binario di ricerca con n nodi. Dopo la chiamata iniziale la procedura chiama se stessa ricorsivamente 2 volte per ogni nodo nell'albero (una per ogni figlio)

Gli alberi binari di ricerca possono supportare gli algoritmi di richiesta di massimo, minimo, successore e predecessore esattamente come quello di ricerca. Queste operazioni richiedono un tempo $O(h)$ in un albero binario di ricerca di altezza h

**Ricerca:**
Procedura di ricerca per un nodo con una chiave data in un albero binario di ricerca:
 
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

Dato un puntatore x alla radice di un sotto-albero e una chiave k, il metodo *treeSearch(x,k)* ritorna un puntatore al nodo in cui la chiave è contenuta. Se tale nodo non esiste il metodo ritorna $NIL$. Per fare una ricerca "globale" della chiave richiamo *treeSearch(T.root,k)*. Questa procedura inizia alla radice dell'albero e traccia un percorso scendendo lungo l'albero. Per ogni nodo che viene incontrato si confronta la chiave k con la *x.key* e se le due chiavi sono uguali la ricerca termina. Se la la *x.key* è più piccola di k la ricerca continua nel sotto-albero sinistro di x, altrimenti (*x.key* più grande di k) si continua nel sotto-albero destro di x.

**Massimo e minimo:**
Per trovare il minimo in un albero di ricerca binario basta seguire i puntatori ai figli sinistri di ogni nodo dalla radice fino a $NIL$ (la procedura ritorna un puntatore al nodo contenente il minimo elemento). Per trovare il massimo si segue lo stesso principio scegliendo però il percorso fatto da sotto-alberi destri.
 
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

**Successore e predecessore:**
Il successore di un nodo x è il nodo con la chiave più grande della *x.key*. Si definisce un successore il successivo nodo visitato in un inorder tree walk. La struttura degli alberi binari di ricerca mi permette di determinare il successore senza comparare le chiavi dei nodi. Tale procedura ritorna il puntatore al nodo successivo se esiste, altrimenti ritorna $NIL$. 
 
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
 - il sotto-albero destro di x è vuoto e x ha un successore y => y è il più piccolo "genitore" di x il cui figlio (di y) è a sua volta un "genitore" di x
 L'algoritmo per trovare il successore di un nodo richiede $O(h)$ tempo per un albero di altezza h in quanto segue o un *path-up* o un *path-down* dell'albero. Anche la procedura per trovare un predecessore, dato che è simmetrica a quella per trovare il successore, ha costo $O(h)$

**Inserimento ed eliminazione:**
Sia la procedura di inserimento che quella di eliminazione causano delle modifiche all'albero (per rispettare le proprietà degli alberi binari di ricerca). La procedura per inserire un nodo dentro ad un albero binario di ricerca richiede un albero *T* e un nodo *z* la cui chiave deve essere già riempita con *z.left* = $NIL$ e *z.right* = $NIL$. 

N.B. la procedura modificherà *T* e alcune proprietà di *z*

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

Esistono 3 possibili strategie per eliminare un nodo da un albero di ricerca binaria:
1. *z* non ha "figli" (al massimo ne ha uno destro), quindi quando lo rimuovo modifico il "genitore" di *z* sostituendo quest'ultimo con $NIL$
                ![[Screenshot From 2026-02-25 11-55-49.png]]
2. *z* ha un "figlio" (in tal caso l'unico "figlio" è sinistro), quindi pongo il "figlio" di *z* al posto di questo e modifico conseguentemente il "genitore" di *z*
                ![[Screenshot From 2026-02-25 11-56-25.png]]
3. *z* ha due "figli", quindi viene trovato il suo successore *y* (nel sotto-albero destro di *z*) e si sposta *y* affinché prenda il posto di *z*. Il resto del restante sotto-albero destro di *z* diventa il nuovo sotto-albero destro di *y* e il sotto-albero sinistro di *z* diventa il sotto-albero sinistro di *y*. *Y*, essendo il successore di *z*, non può avere un "figlio" sinistro e il "figlio" destro di *y* prende la posizione originaria di quest'ultimo (il resto del sotto-albero lo segue automaticamente). Di questo caso ci sono due possibili varianti:
    - *y* è il "figlio" destro di *z* 
                ![[Screenshot From 2026-02-25 12-00-33.png]]
    - *y* si trova nel sotto-albero destro di *z*, ma non è il "figlio" destro
    ![[Screenshot From 2026-02-25 12-02-18.png]]
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
 
Come parte del processo di eliminazione di un nodo, ho la necessità di spostare interi sotto-alberi dell'albero binario di ricerca. posso fare ciò con la subroutine *trasnsplant*, la quale rimpiazza un sotto-albero "figlio" del proprio "genitore" con un altro sotto-albero. Nell'algoritmo quando rimpiazzo il sotto-albero con radice nel nodo *u* con il sotto-albero con radice nel nodo *v*, il "genitore" del nodo *u* diventa il genitore del nodo *v*.
 
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
 
N.B. Sia l'algoritmo di inserimento che quello di eliminazione richiedono tempo $O(h)$ per alberi di altezza h

---
### 26-03-02/03 lz.3/4
Il set di operazioni viste nella scorsa lezione sono veloci se l'altezza dell'albero di ricerca è piccola. Viceversa se l'altezza fosse un valore più grande le stesse operazioni sarebbero veloci tanto quanto una linked list. In questo contesto vengono introdotti gli alberi di ricerca bilanciati (**red-black trees**) che garantiscono delle operazioni che impiegano $O(log(n))$ di tempo nel caso peggiore.

***Proprietà:***
-  Un red-black tree è un albero di ricerca binaria con un ulteriore informazione per ogni nodo: il colore del nodo, il quale può essere o rosso o nero.
-  Vincolando il colore del nodo ogni cammino semplice dalla radice alle foglie di un albero, ci si assicura che nessun cammino del genere sia più lungo del doppio di qualsiasi altro  (dovrebbe voler dire che il cammino più lungo possibile è circa il doppio del cammino più corto possibile)
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

> [!example] ![[Screenshot From 2026-03-02 18-40-20.png]]
> Nel primo dei due alberi i numeri accanto ai nodi stanno ad indicare quanti nodi neri ci sono da quel nodo fino alle foglie (per rispettare la proprietà 5)

La sentinella permette di trattare un figlio $NIL$ di un nodo *x*  come un nodo normale il cui padre sia *x*.  
Il numero di nodi neri in ogni cammino semplice discendente a partire (senza includere) un nodo *x*  fino alle foglie è denotato come *bh(x)* . A partire dalle proprietà sopra definite si può dire che un black-height è ben definito.  In generale un black-height di un red-black tree è un black-height della sua radice. Un red-black tree con *n* nodi interni ha altezza al massimo $2log(n+1)$

***Basic operations:***
Le operazioni di base sono ricerca, minimo, massimo, successore, predecessore e richiedono tutti $O(log(n))$ di tempo. 

> [!info] Come gli alberi binari di ricerca nei casi di insert e di delete viene modificata la struttura dell'albero venendo meno alle proprietà di questi alberi. Per sistemare questo aspetto le proprietà del colore e del puntatore di alcuni nodi devono cambiare. La struttura dei puntatori cambia attraverso rotazione, che è un'operazione locale in un albero di ricerca che preserva le proprietà degli alberi di ricerca binaria. 

Esistono due tipi di rotazione: destra e sinistra. Entrambe richiedono $O(1)$ in tempo. Una rotazione sinistra di un nodo *x*, per esempio, richiede che il figlio destro del nodo (*y*) NON deve essere $T.NIL$. La rotazione cambia il sotto-albero originariamente con radice in *X* scambiando il collegamento tra *x* e *y* (come nell'immagine seguente).
                       ![[Screenshot From 2026-03-03 09-50-21.png]]
***Inserimento:***
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

***Differenze tra treeInsert e RB_insert***
- Tutte le istanze di $NIL$ in treeInsert sono rimpiazzate da $T.NIL$
- *z.left* e *z.right* sono importati a $T.NIL$ per mantenere la struttura dell'albero
- Alla fine il nodo *z* viene colorato di rosso
- Per evitare di violare le proprietà del red-black tree viene alla fine chiamato *RB_insert_fixup* (ogni nodo o è rosso o è nero, le foglie sono nero, il numero di nodi neri è lo stesso in ogni cammino semplice) => uniche proprietà che possono essere violate sono: radice deve essere nera e un nodo rosso non può avere un figlio rosso

> [!info] Ci riferiamo allo zio come il nodo che è il figlio (di qualsiasi colore) di un nodo padre (nero)

1. *y*, lo zio di *z*, è rosso: quando sia *z.p* sia *y* sono rossi. Poiché *z.p.p* è nero, questa proprietà può essere trasferita di un livello in basso sia a *z.p* sia a *y* così da sistemare il problema di avere sia *z* che *z.p* rossi
> [!example] ![[Screenshot From 2026-03-03 10-47-51.png]]
2. *y*, lo zio di *z*, è nero: 
   - *y* è nero e *z* è un figlio destro -> una rotazione sinistra su *z* trasforma la situazione nel prossimo caso
> [!example] ![[Screenshot From 2026-03-03 10-48-53 1.png]]

   - *y* è nero e *z* è un figlio sinistro -> cambia alcuni colori e applica una rotazione destra su *z.p.p*, il quale preserva la proprietà di non avere figli rossi (essendo un nodo rosso)
> [!example] ![[Screenshot From 2026-03-03 10-50-24.png]]

***Hash table:***
Una hash table è una generalizzazione di un array. In un array gli elementi vengono immagazzinati con una chiave k alla posizione k. Data la chiave k, per trovare l'elemento la cui chiave è k, mi basta guardare alla k-esima posizione dell'array (*direct addressing*). Quando non voglio allocare un array con una posizione per ogni possibile chiave uso una hash table. Uso questa struttura anche nel caso in cui il numero di chiavi immagazzinate è più piccolo del numero delle possibili chiavi. Quando si parla di hash table ci si riferisce ad un array che tipicamente usa una dimensione proporzionale al numero di chiavi che vanno immagazzinate. Per questo motivo se ho la chiave k non uso questo valore come indice nell'array, invece definisco una *hash function* (funzione di k) e usiamo il suo valore per indicizzare l'array.

| Direct-address tables                                                        | Hash tables                                                                                                                                               |
| ---------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Mantiene un set dinamico                                                     | Richiede meno spazio di un array dato che può succedere che l'insieme di chiavi usate sia inferiore più piccolo dell'insieme di tutte le chiavi possibili |
| Ogni elemento ha una chiave presa dall'insieme universo $U = \{0...m-1\}$    | Può ridurre la richiesta di storage fino a $\Theta(\|K\|)$, con K l'insieme di chiavi effettivamente usate                                                |
| Non esistono due elementi con la stessa chiave                               | Può ancora richiedere (nel caso medio, non nel peggiore) $O(1)$ in tempo per la ricerca                                                                   |
| Se c'è un elemento *x* con chiave *k*  => *T[k]* contiene il puntatore a *x* | Si usa una funzione *h* per immagazzinare l'elemento nello slot *h(k)*                                                                                    |
> [!info]  
> - chiamiamo *h* la **hash function** e *T* la **hash table**
> - $h: U \rightarrow \{0,...,m-1\}$, così *h(k)* è uno slot in *T*

> [!warning] Quando due o più chiavi puntano allo stesso slot si ha una collisione:
> - può accadere quando ci sono più chiavi che slot $(|U| > m)$
> - dato un set K di chiavi con $|K| < m$ potrebbe come potrebbe non accadere (accade sicuramente se $|K| > m$)

***Hashing indipendente e uniforme:***
Idealmente *h(k)* verrebbe scelto casualmente e indipendentemente da $\{0,...,m-1\}$. Una volta scelto *h(k)* ogni valutazione seguente di *h(k)* deve ridare lo stesso risultato. L'idea che la hash function sia indipendente, uniforme o casuale è teoretico e non può essere implementato nella pratica (si usa per analizzare il comportamento della hash function e per trovare delle approssimazioni pratiche alla forma ideale)

![[Screenshot From 2026-03-03 23-14-22.png]]