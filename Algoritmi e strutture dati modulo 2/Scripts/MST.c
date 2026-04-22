//implementation of Kruskal's and Prim's algorithms for Minimum Spanning Tree
#include <stdio.h>
#include <stdlib.h>

struct Edge
{
    int src;
    int dest;
    int weight;
} edge;

struct Graph
{
    int V;
    int E;
    struct Edge *edges;
} graph;

struct *kruskalMST(struct Graph *graph)
{
    struct Edge *result = (struct Edge *)malloc((graph->V - 1) * sizeof(struct Edge));
    int e = 0; // index for result[]
    int i = 0; // index for sorted edges
    qsort(graph->edges, graph->E, sizeof(struct Edge), compare);
    struct DisjointSet *ds = createDisjointSet(graph->V);
    while (e < graph->V - 1 && i < graph->E)
    {
        struct Edge next_edge = graph->edges[i++];
        int x = find(ds, next_edge.src);
        int y = find(ds, next_edge.dest);
        if (x != y)
        {
            result[e++] = next_edge;
            union(ds, x, y);
        }
    }
    return result;
}

struct primMST(struct Graph *graph)
{
    int *parent = (int *)malloc(graph->V * sizeof(int));
    int *key = (int *)malloc(graph->V * sizeof(int));
    int *mstSet = (int *)malloc(graph->V * sizeof(int));
    for (int i = 0; i < graph->V; i++)
    {
        key[i] = INT_MAX;
        mstSet[i] = 0;
    }
    key[0] = 0;
    parent[0] = -1;
    for (int count = 0; count < graph->V - 1; count++)
    {
        int u = minKey(key, mstSet, graph->V);
        mstSet[u] = 1;
        for (int v = 0; v < graph->V; v++)
            if (graph->adjMatrix[u][v] && !mstSet[v] && graph->adjMatrix[u][v] < key[v])
            {
                parent[v] = u;
                key[v] = graph->adjMatrix[u][v];
            }
    }
    return parent;
}