// Floyd-Warshall e Johnson per il problema degli All Pairs Shortest Path
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
};

static int **allocMatrix(int n)
{
    int **matrix = (int **)malloc(n * sizeof(int *));
    if (matrix == NULL)
        return NULL;

    for (int i = 0; i < n; i++)
    {
        matrix[i] = (int *)malloc(n * sizeof(int));
        if (matrix[i] == NULL)
            return NULL;
    }

    return matrix;
}

static int minDistance(int dist[], int visited[], int n)
{
    int best = INT_MAX;
    int index = -1;

    for (int i = 0; i < n; i++)
        if (!visited[i] && dist[i] <= best)
        {
            best = dist[i];
            index = i;
        }

    return index;
}

static void bellmanFordSimple(struct Graph *graph, int src, int dist[])
{
    for (int i = 0; i < graph->V; i++)
        dist[i] = INT_MAX;
    dist[src] = 0;

    for (int i = 0; i < graph->V - 1; i++)
        for (int j = 0; j < graph->E; j++)
        {
            int u = graph->edges[j].src;
            int v = graph->edges[j].dest;
            int w = graph->edges[j].weight;

            if (dist[u] != INT_MAX && dist[u] + w < dist[v])
                dist[v] = dist[u] + w;
        }
}

static int **dijkstraAllPairs(struct Graph *graph, int **adj, int *h)
{
    int **dist = allocMatrix(graph->V);
    if (dist == NULL)
        return NULL;

    for (int s = 0; s < graph->V; s++)
    {
        int *best = (int *)malloc(graph->V * sizeof(int));
        int *visited = (int *)calloc(graph->V, sizeof(int));
        if (best == NULL || visited == NULL)
        {
            free(best);
            free(visited);
            return NULL;
        }

        for (int i = 0; i < graph->V; i++)
            best[i] = INT_MAX;
        best[s] = 0;

        for (int count = 0; count < graph->V; count++)
        {
            int u = minDistance(best, visited, graph->V);
            if (u == -1)
                break;

            visited[u] = 1;

            for (int v = 0; v < graph->V; v++)
                if (!visited[v] && adj[u][v] != INT_MAX && best[u] != INT_MAX && best[u] + adj[u][v] < best[v])
                    best[v] = best[u] + adj[u][v];
        }

        for (int v = 0; v < graph->V; v++)
        {
            if (best[v] == INT_MAX)
                dist[s][v] = INT_MAX;
            else
                dist[s][v] = best[v] + h[v] - h[s];
        }

        free(best);
        free(visited);
    }

    return dist;
}

// Floyd-Warshall: prova tutte le possibili tappe intermedie tra i nodi.
int **floydWarshall(struct Graph *graph)
{
    int **dist = allocMatrix(graph->V);
    if (dist == NULL)
        return NULL;

    for (int i = 0; i < graph->V; i++)
        for (int j = 0; j < graph->V; j++)
            dist[i][j] = (i == j) ? 0 : INT_MAX;

    for (int i = 0; i < graph->E; i++)
    {
        int u = graph->edges[i].src;
        int v = graph->edges[i].dest;
        int w = graph->edges[i].weight;
        dist[u][v] = w;
    }

    for (int k = 0; k < graph->V; k++)
        for (int i = 0; i < graph->V; i++)
            for (int j = 0; j < graph->V; j++)
                if (dist[i][k] != INT_MAX && dist[k][j] != INT_MAX && dist[i][k] + dist[k][j] < dist[i][j])
                    dist[i][j] = dist[i][k] + dist[k][j];

    return dist;
}

// Johnson: usa Bellman-Ford per i pesi potenziali, poi Dijkstra da ogni sorgente.
int **johnson(struct Graph *graph)
{
    int *h = (int *)malloc(graph->V * sizeof(int));
    int **adj = allocMatrix(graph->V);
    if (h == NULL || adj == NULL)
    {
        free(h);
        return NULL;
    }

    for (int i = 0; i < graph->V; i++)
        h[i] = 0;

    bellmanFordSimple(graph, 0, h);

    for (int i = 0; i < graph->V; i++)
        for (int j = 0; j < graph->V; j++)
            adj[i][j] = INT_MAX;

    for (int i = 0; i < graph->E; i++)
    {
        int u = graph->edges[i].src;
        int v = graph->edges[i].dest;
        int w = graph->edges[i].weight;
        adj[u][v] = w + h[u] - h[v];
    }

    int **dist = dijkstraAllPairs(graph, adj, h);
    free(h);
    return dist;
}