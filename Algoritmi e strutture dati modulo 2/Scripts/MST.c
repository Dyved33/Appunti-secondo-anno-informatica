// implementation of Kruskal's and Prim's algorithms for Minimum Spanning Tree
#include <stdio.h>
#include <stdlib.h>

struct Edge
{
    int src;
    int dest;
    int weight;
} edge;

static int findSet(int parent[], int x)
{
    if (parent[x] != x)
        parent[x] = findSet(parent, parent[x]);
    return parent[x];
}

static void unionSet(int parent[], int rank[], int a, int b)
{
    int rootA = findSet(parent, a);
    int rootB = findSet(parent, b);

    if (rootA == rootB)
        return;

    if (rank[rootA] < rank[rootB])
        parent[rootA] = rootB;
    else if (rank[rootA] > rank[rootB])
        parent[rootB] = rootA;
    else
    {
        parent[rootB] = rootA;
        rank[rootA]++;
    }
}

// edges e' gia' ordinato per peso crescente.
// mst deve avere spazio per almeno (numVertices - 1) archi.
// Ritorna quanti archi sono stati inseriti nel MST.
int kruskalMST(struct Edge edges[], int numEdges, int numVertices, struct Edge mst[])
{
    int *parent = (int *)malloc(numVertices * sizeof(int));
    int *rank = (int *)calloc(numVertices, sizeof(int));
    if (parent == NULL || rank == NULL)
    {
        free(parent);
        free(rank);
        return 0;
    }

    for (int v = 0; v < numVertices; v++)
        parent[v] = v;

    int mstSize = 0;

    for (int i = 0; i < numEdges && mstSize < numVertices - 1; i++)
    {
        int u = edges[i].src;
        int v = edges[i].dest;

        if (findSet(parent, u) != findSet(parent, v))
        {
            mst[mstSize++] = edges[i];
            unionSet(parent, rank, u, v);
        }
    }

    free(parent);
    free(rank);
    return mstSize;
}

/*
MST-Kruskal(G,w)
    A = 0
    for ogni vertice v in G.V
        Make_Set(v)
    crea una lista di archi in G.E
    ordina la lista di archi in ordine crescente in base al peso w
    for ogni arco (u,v) preso dalla lista ordinata
        if Find_Set(u) != Find_Set(v)
            A = A U {(u,v)}
            Union(u,v)
    return A */

// Ritorna quanti archi sono stati inseriti nel MST.
int primMST(struct Edge edges[], int numEdges, int numVertices, struct Edge mst[])
{
    int *inMST = (int *)calloc(numVertices, sizeof(int));
    if (inMST == NULL)
        return 0;
    inMST[0] = 1; // Inizia dal primo vertice
    int mstSize = 0;
    while (mstSize < numVertices - 1)
    {
        struct Edge minEdge = {0, 0, __INT_MAX__};
        for (int i = 0; i < numEdges; i++)
        {
            if ((inMST[edges[i].src] && !inMST[edges[i].dest]) ||
                (!inMST[edges[i].src] && inMST[edges[i].dest]))
            {
                if (edges[i].weight < minEdge.weight)
                    minEdge = edges[i];// Trova l'arco con peso minimo che collega un vertice in MST a uno fuori da MST
            }
        }
        mst[mstSize++] = minEdge;
        inMST[minEdge.src] = 1;
        inMST[minEdge.dest] = 1;
    }
    free(inMST);
    return mstSize;
}