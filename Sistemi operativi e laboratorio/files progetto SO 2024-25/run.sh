#!/bin/bash
# =============================================================================
# run.sh — Compila e avvia l'applicazione "Incrocio"
#
# Uso:
#   chmod +x run.sh
#   ./run.sh
#
# Per terminare l'applicazione: premi CTRL+C
# (lo script invia SIGTERM al processo incrocio, che gestisce la terminazione)
# =============================================================================

set -e   # interrompi se un comando fallisce

# ── Colori per output ─────────────────────────────────────────────────────
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'  # No Color

echo -e "${YELLOW}══════════════════════════════════════════${NC}"
echo -e "${YELLOW}   Applicazione Incrocio — a.a. 2024-25   ${NC}"
echo -e "${YELLOW}══════════════════════════════════════════${NC}"
echo ""

# ── Compilazione ──────────────────────────────────────────────────────────
echo -e "${YELLOW}[1/3] Compilazione...${NC}"

gcc -Wall -Wextra -o incrocio   incrocio.c   -lrt -lpthread
echo -e "  ${GREEN}✓${NC} incrocio compilato"

gcc -Wall -Wextra -o garage     garage.c     -lrt -lpthread
echo -e "  ${GREEN}✓${NC} garage compilato"

gcc -Wall -Wextra -o automobile automobile.c -lrt -lpthread
echo -e "  ${GREEN}✓${NC} automobile compilato"

echo ""

# ── Pulizia file di log precedenti ───────────────────────────────────────
echo -e "${YELLOW}[2/3] Pulizia file di log precedenti...${NC}"
rm -f incrocio.txt auto.txt
echo -e "  ${GREEN}✓${NC} incrocio.txt e auto.txt rimossi"
echo ""

# ── Avvio processi ────────────────────────────────────────────────────────
echo -e "${YELLOW}[3/3] Avvio dell'applicazione...${NC}"
echo ""

# Avvia incrocio in background; è lui che crea le risorse IPC
./incrocio &
INCROCIO_PID=$!
echo -e "  ${GREEN}►${NC} Processo INCROCIO avviato (PID=${INCROCIO_PID})"

# Piccola pausa per lasciare che incrocio inizializzi le risorse IPC
# prima che garage provi ad aprirle
sleep 0.5

# Avvia garage in background
./garage &
GARAGE_PID=$!
echo -e "  ${GREEN}►${NC} Processo GARAGE avviato   (PID=${GARAGE_PID})"

echo ""
echo -e "${YELLOW}── Applicazione in esecuzione ─────────────────${NC}"
echo -e "   Premi ${RED}CTRL+C${NC} per terminare."
echo ""

# ── Gestione CTRL+C / SIGTERM ─────────────────────────────────────────────
# La funzione cleanup viene invocata quando lo script riceve SIGINT o SIGTERM.
cleanup_app() {
    echo ""
    echo -e "${YELLOW}── Terminazione richiesta ──────────────────────${NC}"

    # Invia SIGTERM al processo incrocio: è lui che gestisce la cascata
    # di terminazione verso garage e le auto.
    if kill -0 "${INCROCIO_PID}" 2>/dev/null; then
        echo -e "  Invio SIGTERM a incrocio (PID=${INCROCIO_PID})..."
        kill -TERM "${INCROCIO_PID}"
    fi

    # Attende la terminazione di entrambi i processi principali
    wait "${INCROCIO_PID}" 2>/dev/null || true
    wait "${GARAGE_PID}"   2>/dev/null || true

    # ── Verifica dei file di output ──────────────────────────────────────
    echo ""
    echo -e "${YELLOW}── Verifica file di output ─────────────────────${NC}"

    echo ""
    echo "  Contenuto di incrocio.txt:"
    if [ -s incrocio.txt ]; then
        cat incrocio.txt | sed 's/^/    /'
    else
        echo "    (vuoto o non creato)"
    fi

    echo ""
    echo "  Contenuto di auto.txt:"
    if [ -s auto.txt ]; then
        cat auto.txt | sed 's/^/    /'
    else
        echo "    (vuoto o non creato)"
    fi

    echo ""
    if [ -f incrocio.txt ] && [ -f auto.txt ]; then
        if diff -q incrocio.txt auto.txt > /dev/null 2>&1; then
            echo -e "  ${GREEN}✓ I file incrocio.txt e auto.txt sono identici. Test SUPERATO.${NC}"
        else
            echo -e "  ${RED}✗ I file incrocio.txt e auto.txt sono DIVERSI! Differenze:${NC}"
            diff incrocio.txt auto.txt || true
        fi
    else
        echo -e "  ${YELLOW}⚠  Uno o entrambi i file non esistono.${NC}"
    fi

    echo ""
    echo -e "${YELLOW}══════════════════════════════════════════${NC}"
    echo -e "${GREEN}   Applicazione terminata correttamente.  ${NC}"
    echo -e "${YELLOW}══════════════════════════════════════════${NC}"
    exit 0
}

trap cleanup_app SIGINT SIGTERM

# Attendi la terminazione naturale dei processi (non avviene mai, l'app
# gira indefinitamente finché non si preme CTRL+C)
wait "${INCROCIO_PID}" 2>/dev/null || true
wait "${GARAGE_PID}"   2>/dev/null || true

# Se arriviamo qui senza CTRL+C (es. incrocio è crashato)
cleanup_app
