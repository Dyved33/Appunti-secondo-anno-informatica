//implementation of bfs e dfs and topological sort
#include <stdio.h>
#include <stdlib.h>

struct Graph
{
    int V;
    int **adjMatrix;
} graph;

struct *bfs(struct Graph *graph, int start)
{
    int *visited = (int *)malloc(graph->V * sizeof(int));
    for (int i = 0; i < graph->V; i++)
        visited[i] = 0;
    int *queue = (int *)malloc(graph->V * sizeof(int));
    int front = 0, rear = 0;
    queue[rear++] = start;
    visited[start] = 1;
    while (front < rear)
    {
        int curr = queue[front++];
        printf("%d ", curr);
        for (int i = 0; i < graph->V; i++)
        {
            if (graph->adjMatrix[curr][i] == 1 && !visited[i])
            {
                queue[rear++] = i;
                visited[i] = 1;
            }
        }
    }
}

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