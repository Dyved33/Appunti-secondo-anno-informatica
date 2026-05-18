## Parte 1 — BFS / DFS
### Domanda 1

Hai un grafo non orientato G=(V,E)G=(V,E)G=(V,E).

Vuoi verificare se è un albero.

1. Quali due proprietà devono essere verificate? Devo verificare che non ci siano cicli e che l'albero sia un'unica componente connessa
2. Quale visita useresti? Userei l'algoritmo DFS. Sia questo che il BFS controlla la presenza di cicli in base al verificarsi di archi backward o meno. Scelgo la DFS perché se lancio una DFS Visit e tengo traccia del numero di nodi visitati posso subito verificare se la dimensione della componente connessa è minore o uguale alla cardinalità dell'insieme G.V
3. Qual è la complessità con liste di adiacenza? O(V+E), sia per la DFS che per la BFS

---

### Domanda 2

Perché in un grafo non pesato la BFS trova sempre il cammino minimo in numero di archi? Perché la BFS modifica la distanza tra la sorgente e ogni nodo. Per fare ciò mette dentro ad una coda FIFO tutti i nodi adiacenti a quello che sta esplorando e fa altrettanto per ogni nodo del grafo. Alla fine restituirà la distanza minima dalla sorgente fino ad ogni nodo del grafo

---

### Domanda 3

In una DFS su grafo orientato:

- quando un arco è:
    - di albero: quando sto scoprendo un nuovo nodo. Formalmente ciò avviene quando arrivo ad un nodo di colore bianco
    - backward: quando arrivo ad un nodo grigio. Ciò significa che il nodo in questione è un predecessore del nodo da cui è partito l'arco e che quindi ho strutturato un ciclo
    - forward: quando il nodo al quale arrivo è nero e inoltre il tempo di scoperta del nodo di arrivo è minore del tempo di scoperta del nodo di partenza. Ciò significa che il nodo in questione è un antenato già scoperto e finito
    - cross: quando non c'è alcun tipo di relazione tra i due nodi che compongono l'arco. Tecnicamente un arco cross è qualsiasi arco che non sia di una delle tipologie precedenti. Si riconosce anche perché il tempo di scoperta del nodo di arrivo è maggiore del tempo di scoperta del nodo di partenza e il nodo di arrivo è già nero = già scoperto e finito

---

### Domanda 4

Dato un grafo non orientato connesso:come calcoli il diametro del grafo? Il diametro di un grafo è il minimo delle distanze massime che intercorrono tra ogni coppia di nodi 

1. Quale algoritmo usi? Se si parla di distanze bisogna utilizzare una BFS. Questo algoritmo, infatti, modifica la distanza tra ogni nodo e la sorgente 
2. Quante volte lo esegui? Lo eseguo una volta per ogni nodo dato che la valutazione che devo fare è per ogni nodo del mio grafo
3. Complessità finale? La complessità sarà V volte quella della BFS => O(V(V+E))

---

### Domanda 5

Hai un grafo non orientato e sconnesso. Quanti archi minimi servono per renderlo connesso? Serve un numero di archi pari al numero di componenti connesse meno 1. Questo perché mi è sufficiente usare un solo arco per connettere due componenti connesse distinte, 2 per connetterne 3, ecc. Non mi interessa di utilizzare archi particolari con pesi specifici, ma di connettere componenti connesse piuttosto

---

# Parte 2 — Ordinamento topologico e DAG

### Domanda 6

Quando esiste un ordinamento topologico? Quando non ho cicli in un grafo orientato. Infatti il topological sort si applica solo su grafi DAG, ovvero directed acyclic graphs

---

### Domanda 7

Durante una DFS per il topological sort trovi un arco backward. Cosa puoi concludere immediatamente? Che il topological sort non si può fare dato che il mio grafo ha un ciclo al suo interno

---

### Domanda 8

Perché l’algoritmo DAG-Shortest-Paths può rilassare ogni arco una sola volta? Si tratta di un algoritmo che cerca il cammino di costo minimo partendo da una sorgente singola per arrivare ad uno specifico nodo. Ho la possibilità di rilassare gli archi solo una volta perché all'inizio dell'algoritmo uso il topological sort per ordinare l'insieme dei vertici (ordine lineare e decrescente per i tempi di fine scoperta). Per questo motivo i nodi sono ordinati e basta un solo rilassamento per vedere il costo ottimale per raggiungerli dalla sorgente

---

# Parte 3 — Cammini minimi

### Domanda 9

Hai:

- grafo orientato
- pesato
- con pesi negativi possibili

Quale algoritmo useresti per:

1. trovare i cammini minimi?
2. rilevare cicli negativi?

Userei di certo l'algoritmo di Bellman-Ford per via della possibile presenza di archi di peso negativo. Con questo algoritmo riesco a compiere un rilassamento su ogni arco per ogni nodo presente nel grafo. Alla fine dell'algoritmo si effettua un controllo che consiste nel valutare per ogni arco del grafo se anche dopo il rilassamento il costo del nodo di arrivo è maggiore del costo del nodo di partenza + il peso dell'arco e se ciò si verifica allora siamo in presenza di un ciclo

---

### Domanda 10

Perché Dijkstra non funziona con pesi negativi? Perché l'algoritmo di Dijkstra lavora scegliendo sempre l'arco di costo minimo raggiungibile da nodo in cui si trova. Inoltre se un nodo è già stato valutato, l'algoritmo non si preoccupa di modificarlo e per questo se trovo un arco di peso negativo che raggiunge un nodo già "chiuso" il peso per raggiungerlo non viene modificato

---

### Domanda 11

Confronta:

- Bellman-Ford: vuole un grafo pesato e orientato. Ammette archi di costo negativo ma non cicli. Ha una complessità di O(V*E)
- Dijkstra: vuole un grafo pesato e orientato ma senza archi di costo negativo. Ha una complessità di O(Elog(V))
- DAG-Shortest-Paths: vuole un grafo pesato, orientato e senza cicli. Ammette ogni tipo di peso e ha complessità di O(V+E)

Indicando:

- tipo di grafo richiesto
- pesi ammessi
- complessità

---

### Domanda 12

Nel Metodo Tropicale, cosa rappresenta la matrice $D^{(k)}$? Rappresenta la matrice delle adiacenze che rappresenta il grafo. Contiene: 0 se i = j, infinito se i e j non sono raggiungibili e altrimenti il peso che serve per raggiungere j da i

---

# Parte 4 — MST

### Domanda 13

Perché Kruskal funziona? Qual è l’idea greedy fondamentale? L'algoritmo di Kruskal prende gli archi del grafo in ordine di peso crescente e li aggiunge al MST se e solo se sono dei light edges. Ciò significa che gli archi in questione sono i più leggeri tra tutti gli archi rimasti e collegano un vertice che appartiene al MST con uno che non vi appartiene. L'algoritmo si ferma quando il MST ha V-1 archi. Funziona perché facendo questo controllo (sul rappresentante dell'insieme cui appartiene ogni vertice) evita di avere dei cicli e si assicura di ottenere l'albero di copertura minimo possibile

---

### Domanda 14

Hai un grafo con pesi: $w(e)∈{1,2,3}w(e)\in\{1,2,3\}w(e)∈{1,2,3}$ Perché puoi migliorare l’ordinamento degli archi con un Counting Sort? Perché così li ordinerei dal più piccolo al più grande e l'algoritmo che usa questo insieme ordinato si limita a fare il controllo di validità dell'arco senza preoccuparsi dei pesi

---

### Domanda 15

Un MST e un albero dei cammini minimi dalla sorgente coincidono sempre? Se no, perché? No, perché il MST vuole connettere tutti i vertici del grafo con gli archi più leggeri possibili. Il cammino minimo da una sorgente ad uno specifico nodo vuole solo raggiungere il suo obiettivo e non si preoccupa di passare per ogni nodo del grafo

---

# Parte 5 — Complessità

### Domanda 16

Completa la tabella:
la professoressa è stata vaga e non sempre coerente: a volte scriveva O e a volte $\Theta$ per BFS e DFS. Per questi due so che vengono visitati sempre tutti i nodi ma penso che il simbolo dipenda da se sono certo o meno di attraversare tutti i nodi

| Algoritmo         | Complessità                           |
| ----------------- | ------------------------------------- |
| BFS               | O(V+E)                                |
| DFS               | O(V+E)                                |
| Bellman-Ford      | O(V*E)                                |
| Dijkstra con heap | O(Elog(V))                            |
| Floyd-Warshall    | $\Theta(n^3)$ con n il numero di nodi |
| Kruskal           | O(Elog(V))                            |

---

# Parte 6 — Domande “trabocchetto”

### Domanda 17

Un grafo può:

- avere cicli
- ma NON avere archi backward in DFS?

No, non è possibile. Gli algoritmi in questione vanno avanti nel grafo solo attraverso nodi bianchi e se ne incontrano uno grigio sanno di aver trovato un ciclo

---

### Domanda 18

In un grafo non orientato, possono esistere archi forward o cross durante DFS? Perché? No, ci possono essere solo archi di tipo Tree e backward. Questo perché la DFS prosegue solo in presenza di nodi bianchi, se NON è costretta a visitare un nodo e questo è nero si ferma

---

### Domanda 19

Se un grafo è un DAG, può avere cicli di peso negativo? No, perché la definizione stessa di DAG implica che il grafo sia aciclico

---

### Domanda 20

Perché Floyd-Warshall è un algoritmo di programmazione dinamica? Perché si occupa di valutare il cammino minimo per ogni coppia di archi nel grafo allungando questo cammino di un arco per ogni volta che viene eseguito un prodotto tra matrici