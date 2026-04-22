#include <stdio.h>
#include <stdlib.h>

struct Node
{
    int data;
    struct Node *left;
    struct Node *right;
} node;

struct Node *newNode(int key)
{
    struct Node *node = (struct Node *)malloc(sizeof(struct Node));
    node->data = key;
    node->left = NULL;
    node->right = NULL;
    return node;
}

struct Node *insert(struct Node *parent, int key)
{
    if (parent == NULL)
        return newNode(key); /* se non c'è un nodo padre ne creo uno nuovo */
    if (key < parent->data)
    {
        parent->left = insert(parent->left, key); /* se la chiave è minore del nodo padre, inserisco a sinistra */
    }
    else if (key > parent->data)
    {
        parent->right = insert(parent->right, key); /* se la chiave è maggiore del nodo padre, inserisco a destra */
    }
    return parent; /* ritorno il nodo padre */
}

struct Node *search(struct Node *root, int key)
{
    if (root == NULL || root->data == key)
        return root;
    if (key < root->data)
        return search(root->left, key);
    return search(root->right, key);
}

struct Node *minimum(struct Node *root)
{
    while (root->left != NULL)
        root = root->left;
    return root;
}

struct Node *maximum(struct Node *root)
{
    while (root->right != NULL)
        root = root->right;
    return root;
}

struct Node *Predecessor(struct Node *root, struct Node *curr)
{
    if (curr->left != NULL)
        return maximum(curr->left);
    struct Node *predecessor = NULL;
    while (root != NULL)
    {
        if (curr->data > root->data)
        {
            predecessor = root;
            root = root->right;
        }
        else
            root = root->left;
    }
    return predecessor;        
}

struct Node *Successor(struct Node *root, struct Node *curr)
{
    if (curr->right != NULL)
        return minimum(curr->right);
    struct Node *successor = NULL;
    while (root != NULL)
    {
        if (curr->data < root->data)
        {
            successor = root;
            root = root->left;
        }
        else
            root = root->right;
    }
    return successor;        
}

struct Node *delete(struct Node *root, int key)
{
    if (root == NULL)
        return root;
    if (key < root->data)
        root->left = delete(root->left, key);
    else if (key > root->data)
        root->right = delete(root->right, key);
    else
    {
        if (root->left == NULL)
        {
            struct Node *temp = root->right;
            free(root);
            return temp;
        }
        else if (root->right == NULL)
        {
            struct Node *temp = root->left;
            free(root);
            return temp;
        }
        struct Node *temp = minimum(root->right);
        root->data = temp->data;
        root->right = delete(root->right, temp->data);
    }
    return root;
}