#include <stdio.h>
#include <stdlib.h>

struct Node
{
    int data;
    struct Node *left;
    struct Node *right;
    char color;
} node;

struct Node *newNode(int key)
{
    struct Node *node = (struct Node *)malloc(sizeof(struct Node));
    node->data = key;
    node->left = NULL;
    node->right = NULL;
    node->color = 'R';
    return node;
}

struct Node *insert(struct Node *parent, int key)
{
    if (parent == NULL)
        return newNode(key); /* se non c'è un nodo ne creo uno nuovo */
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

/* RB_insert (T, z)
x = T.root
y = T.NIL
while x != T.NIL
	  y = x
	  if z.key < x.key
		   x = x.left
	  else
		   x = x.right 
 z.p = y
 
 if y == T.NIL
	  T.root = z
 elseif z.key < y.key
	  y.left = z
 else 
	  y.right = z
 z.left = T.NIL
 z.right = T.NIL
 z.color = RED
 RB_insert_fixup */