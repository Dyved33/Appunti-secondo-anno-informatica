// implementation of bfs e dfs and topological sort
#include <stdio.h>
#include <stdlib.h>

struct Node
{
    int d; // minimo numero di archi dalla sorgente al nodo
    int p;
    struct Node *next;
    char color;
} Node;

//suppongo di avere una linked list di nodi e la relativa lunghezza (len). suppongo inoltre che src sia il primo nodo
int *bfs(struct Node *src, struct Node *curr, int len)
{
    int queue[len];

    while (src->next != NULL)
    {
        src->color = 'W';
        src->d = __UINT8_MAX__;
        src->p = NULL;
    }

    src->color = 'G';
    src->d = 0;
    src->p = NULL;

    
}

//bianco = non visitato
//grigio = iniziato ma non finito
//nero = visitato

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


struct *dfs(struct Graph *graph, int start)
{
    int *visited = (int *)malloc(graph->V * sizeof(int));
    for (int i = 0; i < graph->V; i++)
        visited[i] = 0;
    dfsUtil(graph, start, visited);
}

int topologicalSort(struct Graph *graph)
{
    int *visited = (int *)malloc(graph->V * sizeof(int));
    for (int i = 0; i < graph->V; i++)
        visited[i] = 0;
    int *stack = (int *)malloc(graph->V * sizeof(int));
    int top = -1;
    for (int i = 0; i < graph->V; i++)
        if (!visited[i])
            topologicalSortUtil(graph, i, visited, stack, &top);
    while (top >= 0)
        printf("%d ", stack[top--]);
}