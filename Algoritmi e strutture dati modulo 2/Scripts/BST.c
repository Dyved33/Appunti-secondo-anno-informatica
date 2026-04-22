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

// mi sposto nell'albero a destra o a sinistra del nodo che sto leggendo (in base alle proprietà degli alberi binari) finché non arrivo in fondo all'albero e lì inserisco il nodo
struct Node *insert(struct Node *parent, int key)
{
    if (parent == NULL)
        return newNode(key); // se non c'è un nodo padre ne creo uno nuovo (mi serve anche per inserire il nuovo nodo in fondo all'albero)
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

// per condurre la ricerca mi bastano le proprietà dell'albero
struct Node *search(struct Node *root, int key)
{
    if (root == NULL || root->data == key)
        return root;
    if (key < root->data)
        return search(root->left, key);
    else if (key > root->data)
        return search(root->right, key);
}

// per trovare il minimo nell'albero vado tutto a sinistra
struct Node *minimum(struct Node *root)
{
    while (root->left != NULL)
        root = root->left;
    return root;
}

// per trovare il massimo nell'albero vado tutto a destra
struct Node *maximum(struct Node *root)
{
    while (root->right != NULL)
        root = root->right;
    return root;
}

// cerco il maggiore dei nodi più piccoli del corrente
struct Node *Predecessor(struct Node *root, struct Node *curr)
{
    if (root == NULL || curr == NULL)
        return NULL;

    // Caso 1: c'è sottoalbero sinistro -> massimo del sottoalbero sinistro
    if (curr->left != NULL)
        return maximum(curr->left);

    // Caso 2: nessun sottoalbero sinistro -> risalgo dagli antenati
    struct Node *parent = root;
    while (parent != curr)
    {
        if (curr->data > parent->data)
        {
            parent = parent->right;
        }
        else if (curr->data < parent->data)
        {
            parent = parent->left;
        }
    }
    return parent;
}

// cerco il minore dei nodi più grandi del corrente
struct Node *Successor(struct Node *root, struct Node *curr)
{
    if (root == NULL || curr == NULL)
        return NULL;

    // Caso 1: c'è sottoalbero destro -> minimo del sottoalbero destro
    if (curr->left != NULL)
        return minimum(curr->right);

    // Caso 2: nessun sottoalbero destro -> risalgo dagli antenati
    struct Node *parent = root;
    while (parent != curr)
    {
        if (curr->data > parent->data)
        {
            parent = parent->right;
        }
        else if (curr->data < parent->data)
        {
            parent = parent->left;
        }
    }
    return parent;
}

struct Node *delete(struct Node *root, struct Node *curr)
{
    struct Node *tmp = NULL;
    struct Node *parent = root;
    struct Node *succ = NULL;

    tmp = search(root, curr);

    while (parent != curr)
    {
        if (curr->data > parent->data)
        {
            parent = parent->right;
        }
        else if (curr->data < parent->data)
        {
            parent = parent->left;
        }
    }

    if (tmp == NULL)
        return NULL;

    if (tmp->left == NULL && tmp->right == NULL)
        free(tmp);
    else if (tmp->left == NULL)
    {
        if (parent->data > tmp->data)
        {
            parent->left = tmp->right;
            free(tmp);
        }
        else
        {
            parent->right = tmp->right;
            free(tmp);
        }
    }
    else if (tmp->right == NULL)
    {
        if (parent->data > tmp->data)
        {
            parent->left = tmp->left;
            free(tmp);
        }
        else
        {
            parent->right = tmp->left;
            free(tmp);
        }
    }
    else if (tmp->left != NULL && tmp->right != NULL)
    {
        succ = Successor(root, curr);

        if (parent->data > succ->data)
        {
            parent->left = succ;
            succ->left = tmp->left;
            succ->right = tmp->right;
            free(tmp);
        }
        else if (parent->data < succ->data)
        {
            parent->right = succ;
            succ->left = tmp->left;
            succ->right = tmp->right;
            free(tmp);
        }
    }
}

/* treeDelete(T,z)
 if z.left == NIL
  transplant(T,z,z.right)
 elseif z.right == NIL
  transplant(T,z,z.left)
 else
  y = treeMinimum(z.right)

  if y != z.right
   transplant(T,y,y.right)
   y.right = z.right
   y.right = y

  transplant(T,z,y)
  y.left = z.left
  y.left.p = y */

  
  /* transplant(T,u,v)
 if u.p == NIL
  T.root = v
 elseif u == u.p.left
  u.p.left = v
 else 
  u.p.right = v

 if v != NIL
  v.p = u.p
 */