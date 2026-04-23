// implementation of bfs e dfs and topological sort
#include <stdio.h>
#include <stdlib.h>
#include <limits.h>

struct Node
{
    int d; // minimo numero di archi dalla sorgente al nodo
    struct Node *p;
    int f; // tempo di completamento
    struct Node *next;
    char color;
} Node;

// suppongo di avere una linked list di nodi e la relativa lunghezza (len). suppongo inoltre che src sia il primo nodo
int *bfs(struct Node *src, struct Node *curr, int len)
{
    if (src == NULL || curr == NULL || len <= 0)
        return NULL;

    struct Node **queue = (struct Node **)malloc(len * sizeof(struct Node *));
    if (queue == NULL)
        return NULL;

    int front = 0;
    int rear = 0;

    // Inizializza tutti i nodi: non visitati, distanza infinita, padre nullo.
    for (struct Node *n = curr; n != NULL; n = n->next)
    {
        n->color = 'W';
        n->d = INT_MAX;
        n->p = NULL;
    }

    src->color = 'G';
    src->d = 0;
    src->p = NULL;
    queue[rear++] = src;

    while (front < rear)
    {
        struct Node *u = queue[front++];

        // Assunzione: i vicini di u sono nella lista concatenata u->next.
        for (struct Node *v = u->next; v != NULL; v = v->next)
        {
            if (v->color == 'W')
            {
                v->color = 'G';
                v->d = (u->d == INT_MAX) ? INT_MAX : u->d + 1;
                v->p = u;

                if (rear < len)
                    queue[rear++] = v;
            }
        }
        queue[rear++] = u; // Aggiungi u alla coda dopo aver esplorato i suoi vicini
        u->color = 'B';
    }
    return queue; // Restituisce la coda con i nodi visitati in ordine di esplorazione
}

// bianco = non visitato
// grigio = iniziato ma non finito
// nero = visitato

/* BFS(G,s)
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
  u.color = BLACK // u è dietro la frontiera ora */

void dfs(struct Node *src)
{
    if (src == NULL)
        return;

    // inizializza tutti i nodi: non visitati, padre nullo
    for (struct Node *n = src; n != NULL; n = n->next)
    {
        n->color = 'W';
        n->p = NULL;
        n->f = 0;
    }

    int time = 0;
    for (struct Node *n = src; n != NULL; n = n->next)
    {
        if (n->color == 'W')
        {
            time++;
            n->d = time;
            n->color = 'G';

            for (struct Node *v = n->next; v != NULL; v = v->next)
            {
                if (v->color == 'W')
                {
                    v->p = n;
                    dfs(v);
                }
            }
            time++;
            n->f = time;
            n->color = 'B';
        }
    }
}

/* DFS(G)
    for each vertex u in G.V
        u.color = WHITE
        u.parent = NIL
    time = 0
    for each vertex u in G.V
        if u.color == WHITE
            DFS-Visit(G,u) */

/* DFS-Visit
    time = time +1
    u.d = time
    u.color = GRAY
    for each vertex v in G.Adj[u]
        if v.color == WHITE
            v.parent = u
            DFS-Visit(G,v)
    time = time +1
    u.f = time
    u.color = BLACK */

// suppongo di avere una linked list di len nodi
int *topologicalSort(struct Node *src, int len)
{
    dfs(src);
    int *sorted = (int *)malloc(len * sizeof(int));
    if (sorted == NULL)
        return NULL;
    int index = 0;
    for (struct Node *n = src; n != NULL; n = n->next)
    {
        sorted[index++] = n->d; // ordina i nodi in base al tempo di completamento
    }
    return sorted;
}