/*
 * incrocio.h — Funzioni di supporto fornite dal docente
 *
 * Contiene:
 *   EstraiDirezione(int iAutomobile)
 *       Estrae casualmente la strada in cui l'automobile iAutomobile dovrà svoltare.
 *       La strada restituita è garantita essere diversa da iAutomobile.
 *
 *   GetNextCar(int *Direzioni)
 *       Determina quale automobile può attraversare l'incrocio secondo le regole
 *       di precedenza del codice della strada (precedenza a destra).
 *       Direzioni[i] == -1  →  l'auto i ha già attraversato (ignorata).
 *       Restituisce l'indice dell'auto che può passare, o -1 se nessuna.
 *
 * Tutte le funzioni sono dichiarate static inline per evitare symbol conflicts
 * quando l'header è incluso in più translation unit.
 */

#ifndef INCROCIO_H
#define INCROCIO_H

#include <stdlib.h>   /* rand() */

/* ────────────────────────────────────────────────────────────────────────────
 * EstraiDirezione
 *
 * Parametro:  iAutomobile — numero della strada da cui proviene l'auto (0‥3)
 * Ritorna:    numero della strada in cui dovrà svoltare  (≠ iAutomobile)
 * ────────────────────────────────────────────────────────────────────────── */
static inline int EstraiDirezione(int iAutomobile)
{
    int dir;
    do {
        dir = rand() % 4;
    } while (dir == iAutomobile);
    return dir;
}

/* ────────────────────────────────────────────────────────────────────────────
 * GetNextCar
 *
 * Implementa la regola della "precedenza a destra" dei quadrivi italiani:
 * un veicolo deve cedere il passo a chi proviene dalla sua destra.
 *
 * Convenzione strade (senso orario): 0 → 1 → 2 → 3 → 0
 * La strada alla DESTRA della strada i è  (i + 1) % 4
 *
 * Un'auto i può passare se la strada alla sua destra ((i+1)%4) è libera
 * (cioè Direzioni[(i+1)%4] == -1, l'auto di destra ha già attraversato).
 *
 * Caso di stallo (tutte e 4 le auto presenti):
 *   La regola ciclica genera un deadlock; si rompe cedendo la precedenza
 *   all'auto con indice più basso ancora in attesa.
 *
 * Parametro:  Direzioni — array di 4 interi (valore ≥ 0 = in attesa, -1 = passata)
 * Ritorna:    indice dell'auto che può attraversare, oppure -1 se nessuna
 * ────────────────────────────────────────────────────────────────────────── */
static inline int GetNextCar(int *Direzioni)
{
    int i;

    /* Prima passata: cerca un'auto che non abbia nessuno alla propria destra */
    for (i = 0; i < 4; i++) {
        if (Direzioni[i] == -1) continue;          /* già passata, salta      */
        int right = (i + 1) % 4;                   /* strada alla sua destra  */
        if (Direzioni[right] == -1) {              /* nessuno a destra → pass */
            return i;
        }
    }

    /*
     * Seconda passata (risoluzione stallo):
     * Tutte le auto in attesa hanno qualcuno alla propria destra.
     * Si dà precedenza alla prima auto disponibile (indice minore).
     */
    for (i = 0; i < 4; i++) {
        if (Direzioni[i] != -1) return i;
    }

    return -1;   /* nessuna auto in attesa */
}

#endif /* INCROCIO_H */
