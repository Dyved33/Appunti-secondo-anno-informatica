//implementation of Floyd-Warshall and Johnson's algorithms for All Pairs Shortest Path
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

struct *floydWarshall(struct Graph *graph)
{
    int **dist = (int **)malloc(graph->V * sizeof(int *));
    for (int i = 0; i < graph->V; i++)
        dist[i] = (int *)malloc(graph->V * sizeof(int));
    for (int i = 0; i < graph->V; i++)
        for (int j = 0; j < graph->V; j++)
            if (i == j)
                dist[i][j] = 0;
            else
                dist[i][j] = INT_MAX;
    for (int i = 0; i < graph->E; i++)
    {
        int u = graph->edges[i].src;
        int v = graph->edges[i].dest;
        int weight = graph->edges[i].weight;
        dist[u][v] = weight;
    }
    for (int k = 0; k < graph->V; k++)
        for (int i = 0; i < graph->V; i++)
            for (int j = 0; j < graph->V; j++)
                if (dist[i][k] != INT_MAX && dist[k][j] != INT_MAX && dist[i][k] + dist[k][j] < dist[i][j])
                    dist[i][j] = dist[i][k] + dist[k][j];
    return dist;
}

struct *johnson(struct Graph *graph)
{
    int *h = (int *)malloc(graph->V * sizeof(int));
    for (int i = 0; i < graph->V; i++)
        h[i] = 0;
    for (int i = 1; i < graph->V; i++)
        for (int j = 0; j < graph->E; j++)
        {
            int u = graph->edges[j].src;
            int v = graph->edges[j].dest;
            int weight = graph->edges[j].weight;
            if (h[u] != INT_MAX && h[u] + weight < h[v])
                h[v] = h[u] + weight;
        }
    for (int j = 0; j < graph->E; j++)
    {
        int u = graph->edges[j].src;
        int v = graph->edges[j].dest;
        int weight = graph->edges[j].weight;
        if (h[u] != INT_MAX && h[u] + weight < h[v])
        {
            printf("Graph contains negative weight cycle");
            return NULL;
        }
    }
    for (int i = 0; i < graph->E; i++)
    {
        int u = graph->edges[i].src;
        int v = graph->edges[i].dest;
        int weight = graph->edges[i].weight;
        graph->edges[i].weight += h[u] - h[v];
    }
    struct *dist = (struct *)malloc(graph->V * sizeof(struct));
    for (int i = 0; i < graph->V; i++)
        dist[i] = dijkstra(graph, i);
    for (int i = 0; i < graph->V; i++)
        for (int j = 0; j < graph->V; j++)
            if (dist[i][j] != INT_MAX)
                dist[i][j] += h[j] - h[i];
    return dist;
}