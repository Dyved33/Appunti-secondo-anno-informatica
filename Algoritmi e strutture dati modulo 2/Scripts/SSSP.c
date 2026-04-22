//implementation of Bellman-Ford (and dag shortest path) and Dijkstra's algorithms for Single Source Shortest Path 
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

struct *bellmanFord(struct Graph *graph, int src)
{
    int *dist = (int *)malloc(graph->V * sizeof(int));
    for (int i = 0; i < graph->V; i++)
        dist[i] = INT_MAX;
    dist[src] = 0;
    for (int i = 1; i < graph->V; i++)
        for (int j = 0; j < graph->E; j++)
        {
            int u = graph->edges[j].src;
            int v = graph->edges[j].dest;
            int weight = graph->edges[j].weight;
            if (dist[u] != INT_MAX && dist[u] + weight < dist[v])
                dist[v] = dist[u] + weight;
        }
    for (int j = 0; j < graph->E; j++)
    {
        int u = graph->edges[j].src;
        int v = graph->edges[j].dest;
        int weight = graph->edges[j].weight;
        if (dist[u] != INT_MAX && dist[u] + weight < dist[v])
            printf("Graph contains negative weight cycle");
    }
    return dist;
}

struct *DAGshortestPath(struct Graph *graph, int src)
{
    int *dist = (int *)malloc(graph->V * sizeof(int));
    for (int i = 0; i < graph->V; i++)
        dist[i] = INT_MAX;
    dist[src] = 0;
    int *stack = (int *)malloc(graph->V * sizeof(int));
    int top = -1;
    for (int i = 0; i < graph->V; i++)
        if (!visited[i])
            topologicalSortUtil(graph, i, visited, stack, &top);
    while (top >= 0)
    {
        int u = stack[top--];
        if (dist[u] != INT_MAX)
            for (int v = 0; v < graph->V; v++)
                if (graph->adjMatrix[u][v] && dist[u] + graph->adjMatrix[u][v] < dist[v])
                    dist[v] = dist[u] + graph->adjMatrix[u][v];
    }
    return dist;
}

struct *dijkstra(struct Graph *graph, int src)
{
    int *dist = (int *)malloc(graph->V * sizeof(int));
    for (int i = 0; i < graph->V; i++)
        dist[i] = INT_MAX;
    dist[src] = 0;
    int *visited = (int *)malloc(graph->V * sizeof(int));
    for (int i = 0; i < graph->V; i++)
        visited[i] = 0;
    for (int count = 0; count < graph->V - 1; count++)
    {
        int u = minDistance(dist, visited, graph->V);
        visited[u] = 1;
        for (int v = 0; v < graph->V; v++)
            if (!visited[v] && graph->adjMatrix[u][v] && dist[u] != INT_MAX && dist[u] + graph->adjMatrix[u][v] < dist[v])
                dist[v] = dist[u] + graph->adjMatrix[u][v];
    }
    return dist;
}