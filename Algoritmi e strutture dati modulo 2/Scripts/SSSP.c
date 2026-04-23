// Single Source Shortest Path: Bellman-Ford, DAG shortest path, Dijkstra
#include <stdio.h>
#include <stdlib.h>
#include <limits.h>

struct Edge
{
    int src;
    int dest;
    int weight;
};

struct Graph
{
    int V;
    int E;
    struct Edge *edges;
    int **adjMatrix;
};

static int minDistance(int dist[], int visited[], int V)
{
    int min = INT_MAX;
    int minIndex = -1;

    for (int i = 0; i < V; i++)
        if (!visited[i] && dist[i] <= min)
        {
            min = dist[i];
            minIndex = i;
        }

    return minIndex;
}

static void topoSortUtil(struct Graph *graph, int v, int visited[], int stack[], int *top)
{
    visited[v] = 1;

    for (int i = 0; i < graph->V; i++)
        if (graph->adjMatrix[v][i] && !visited[i])
            topoSortUtil(graph, i, visited, stack, top);

    stack[++(*top)] = v;
}

// Bellman-Ford: funziona anche con pesi negativi.
int *bellmanFord(struct Graph *graph, int src)
{
    int *dist = (int *)malloc(graph->V * sizeof(int));
    if (dist == NULL)
        return NULL;

    for (int i = 0; i < graph->V; i++)
        dist[i] = INT_MAX;
    dist[src] = 0;

    for (int i = 0; i < graph->V - 1; i++)
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
        {
            printf("Graph contains negative weight cycle\n");
            break;
        }
    }

    return dist;
}

// Shortest path in a DAG: prima ordino i nodi topologicamente, poi rilasso gli archi.
int *DAGshortestPath(struct Graph *graph, int src)
{
    int *dist = (int *)malloc(graph->V * sizeof(int));
    int *visited = (int *)calloc(graph->V, sizeof(int));
    int *stack = (int *)malloc(graph->V * sizeof(int));

    if (dist == NULL || visited == NULL || stack == NULL)
    {
        free(dist);
        free(visited);
        free(stack);
        return NULL;
    }

    for (int i = 0; i < graph->V; i++)
        dist[i] = INT_MAX;
    dist[src] = 0;

    int top = -1;
    for (int i = 0; i < graph->V; i++)
        if (!visited[i])
            topoSortUtil(graph, i, visited, stack, &top);

    while (top >= 0)
    {
        int u = stack[top--];

        if (dist[u] == INT_MAX)
            continue;

        for (int v = 0; v < graph->V; v++)
            if (graph->adjMatrix[u][v] && dist[u] + graph->adjMatrix[u][v] < dist[v])
                dist[v] = dist[u] + graph->adjMatrix[u][v];
    }

    free(visited);
    free(stack);
    return dist;
}

// Dijkstra: funziona con pesi non negativi.
int *dijkstra(struct Graph *graph, int src)
{
    int *dist = (int *)malloc(graph->V * sizeof(int));
    int *visited = (int *)calloc(graph->V, sizeof(int));

    if (dist == NULL || visited == NULL)
    {
        free(dist);
        free(visited);
        return NULL;
    }

    for (int i = 0; i < graph->V; i++)
        dist[i] = INT_MAX;
    dist[src] = 0;

    for (int count = 0; count < graph->V - 1; count++)
    {
        int u = minDistance(dist, visited, graph->V);

        if (u == -1)
            break;

        visited[u] = 1;

        for (int v = 0; v < graph->V; v++)
            if (!visited[v] && graph->adjMatrix[u][v] && dist[u] != INT_MAX && dist[u] + graph->adjMatrix[u][v] < dist[v])
                dist[v] = dist[u] + graph->adjMatrix[u][v];
    }

    free(visited);
    return dist;
}