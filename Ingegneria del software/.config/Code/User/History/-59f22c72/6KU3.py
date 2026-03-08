import sqlite3
import os
import csv
import tempfile
import subprocess
from datetime import datetime, timedelta, date
from dateutil import parser
from typing import List, Tuple, Optional, Dict, Any
import re
import unicodedata
import traceback
import sys
from io import BytesIO
from PIL import Image

# PDF generation
try:
    from reportlab.pdfgen import canvas
    from reportlab.lib.pagesizes import A4
    from reportlab.lib.units import mm
    from reportlab.pdfbase.pdfmetrics import stringWidth
except Exception:
    canvas = None

# PDF manipulation
try:
    from PyPDF2 import PdfReader, PdfWriter
except Exception:
    PdfReader = PdfWriter = None

# Barcode generation
try:
    from barcode import Code128
    from barcode.writer import ImageWriter
except Exception:
    Code128 = None

# DBF file reading
try:
    from dbfread import DBF
except ImportError:
    DBF = None


DB_FILE = "database.db"

# ========================
# CLIENTI COLUMN DEFINITIONS
# ========================

CLIENTI_COLUMNS = [
    "RAG_SOC","P_IVA","IND_VIA","COMUNE","PROV","CAP","MF","MARCA","LOGOTIPO",
    "MATRICOLA","SCAD_ASS","IMP_ASS","VF","IMP_RIP","TEL","COMMENTI","PASSAGGIO","DATAINST",
    "TIPOMF","DATA","GRUPPO","DATADIS","CONTA","DM","CF","SCANNER","SCANNERAGG","N_PINPAD",
    "N_PIN_BACK","PEC","FINEMF","CODPUNTOCA","RTREADY","CODUNIVOCO","GARANZIA","AZZ","QRCODE",
    "LICENZA","PSW111","PSWSCDIG","MXAZ"
]

# Colonne possibili per INTAE (normalizzate, tutte TEXT)
INTAE_COLUMNS = [
    "CODICE_INTERVENTO", "PIVA", "DATA_INTERVENTO",
    "CF_TECNICO", "STATO_CF", "DISPOSITIVO", "MATRICOLA", "TIPOLOGIA",
    "DATA_TRASMISSIONE", "NOTA", "VERSIONE_SW", "DATA_RELEASE"
]

# Alias possibili per l'import INTAE (CSV/DBF)
INTAE_ALIASES = {
    "CODICE_INTERVENTO": [
        "CODICE INTERVENTO", "CODICE_INTERVENTO", "CODICE", "COD. INTERVENTO",
        "ID INTERVENTO", "ID_INTERVENTO", "COD INTERVENTO", "CODICEINTERVENTO"
    ],
    "PIVA": [
        "PIVA", "P.IVA", "P IVA", "PARTITA IVA", "PIVA LABORATORIO",
        "PIVA_LABORATORIO", "PIVA LAB", "PIVA LABORATORI", "PIVAESERCE"
    ],
    "DATA_INTERVENTO": [
        "DATA INTERVENTO", "DATA_INTERVENTO", "DATA", "DATA INTERV.", "DATAINTERVENTO"
    ],
    "CF_TECNICO": [
        "CF TECNICO", "CF_TECNICO", "CF", "CODICE FISCALE", "CODICE_FISCALE"
    ],
    "STATO_CF": [
        "STATO CF", "STATO_CF", "STATO", "STATO CODICE FISCALE"
    ],
    "DISPOSITIVO": [
        "DISPOSITIVO", "DEVICE", "APPARECCHIO", "STRUMENTO"
    ],
    "MATRICOLA": [
        "MATRICOLA", "MATR", "MATRICOLE"
    ],
    "TIPOLOGIA": [
        "TIPOLOGIA", "TIPO", "TIPOLOGIA INTERVENTO"
    ],
    "DATA_TRASMISSIONE": [
        "DATA TRASMISSIONE", "DATA_TRASMISSIONE", "TRASMISSIONE", "DATA INVIO"
    ],
    "NOTA": [
        "NOTA", "NOTE", "ANNOTAZIONI", "OSSERVAZIONI"
    ],
    "VERSIONE_SW": [
        "VERSIONE SW", "VERSIONE_SW", "VERSIONE SOFTWARE", "SW VERSION"
    ],
    "DATA_RELEASE": [
        "DATA RELEASE", "DATA_RELEASE", "RELEASE", "DATA VERSIONE"
    ],
}

def _normalize_header_key(h: str) -> str:
    return re.sub(r"[^A-Z0-9]+", "_", str(h or "").strip().upper()).strip("_")

def _intae_known_header_keys() -> set:
    keys = set()
    for col, aliases in INTAE_ALIASES.items():
        keys.add(_normalize_header_key(col))
        for a in aliases:
            keys.add(_normalize_header_key(a))
    return keys

def _build_intae_header_map(headers: list) -> dict:
    """Mappa header -> colonna INTAE usando alias e match per token."""
    alias_norm = {}
    for col, aliases in INTAE_ALIASES.items():
        alias_norm[col] = {_normalize_header_key(col)} | {_normalize_header_key(a) for a in aliases}

    mapping = {}
    used_cols = set()
    for h in headers:
        nh = _normalize_header_key(h)
        if not nh:
            continue
        best_col = None
        best_score = 0
        for col, alias_set in alias_norm.items():
            score = 0
            if nh in alias_set:
                score = 3
            else:
                ht = set(nh.split("_"))
                for a in alias_set:
                    at = set(a.split("_"))
                    score = max(score, len(ht & at))
            if score > best_score:
                best_score = score
                best_col = col
        if best_col and best_score >= 1 and best_col not in used_cols:
            mapping[h] = best_col
            used_cols.add(best_col)
    return mapping

def _intae_parse_date_value(raw) -> str:
    if raw is None:
        return ""
    s = str(raw).strip()
    if not s:
        return ""
    try:
        dt = _parse_date(s)
        if isinstance(dt, datetime):
            return dt.date().isoformat()
        if isinstance(dt, date):
            return dt.isoformat()
    except Exception:
        pass
    return s

def _intae_build_record(row_dict: dict, header_map: dict) -> dict:
    rec = {c: "" for c in INTAE_COLUMNS}
    for h, col in header_map.items():
        val = row_dict.get(h, "")
        rec[col] = str(val).strip() if val is not None else ""

    # normalizza date
    if rec.get("DATA_INTERVENTO"):
        rec["DATA_INTERVENTO"] = _intae_parse_date_value(rec["DATA_INTERVENTO"])
    if rec.get("DATA_TRASMISSIONE"):
        rec["DATA_TRASMISSIONE"] = _intae_parse_date_value(rec["DATA_TRASMISSIONE"])
    if rec.get("DATA_RELEASE"):
        rec["DATA_RELEASE"] = _intae_parse_date_value(rec["DATA_RELEASE"])

    return rec

# Colonne per MISURART (normalizzate, tutte TEXT tranne id)
MISURART_COLUMNS = [
    "TIPO","CF","ALFASIGILL","NUMSIGILLO","DATAINIZIO","DATAFINE","ESITO",
    "PIVAESERCE","RAGSOCESE","COMUNESE","PROVINCIA","INDIRIZZO","CAPESE",
    "TIPOINT","MARCAMF","MODELLO","LOGOTIPO","MATRICOLA",
    "F1","F2","F3","F4","F5","F6","CARCTRL","DATAVP"
]

DATE_COLS = {"FINEMF", "GARANZIA", "DATAINST", "SCAD_ASS", "DATA", "DATADIS", "DATA_INTERVENTO", "DATAINIZIO", "DATAFINE", "VF", "DATAVP"}

# -------------------------
# DB connection & init
# -------------------------

def get_connection():
    # abilitiamo foreign keys per sicurezza
    conn = sqlite3.connect(DB_FILE)
    conn.execute("PRAGMA foreign_keys = ON")
    conn.row_factory = row_factory_no_none
    return conn

def row_factory_no_none(cursor, row):
    return tuple("" if v is None else str(v) if not isinstance(v, str) else v for v in row)

def _create_clienti_table(cur):
    """Crea la tabella clienti con le colonne richieste (se non esiste)"""
    cols_sql = [ '"id" INTEGER PRIMARY KEY AUTOINCREMENT' ]
    for c in CLIENTI_COLUMNS:
        cols_sql.append(f'"{c}" TEXT')
    sql = f'CREATE TABLE IF NOT EXISTS clienti ({", ".join(cols_sql)});'
    cur.execute(sql)


def _add_missing_clienti_columns(cur):
    """Aggiunge le colonne mancanti (dalle CLIENTI_COLUMNS) con ALTER TABLE.
    Non rimuove colonne extra esistenti.
    """
    cur.execute("PRAGMA table_info(clienti)")
    existing = [r[1] for r in cur.fetchall()]
    for c in CLIENTI_COLUMNS:
        if c not in existing:
            cur.execute(f'ALTER TABLE clienti ADD COLUMN "{c}" TEXT')


def init_database(recreate_clienti: bool = False):
    """
    Crea il DB e le tabelle principali.
    Se recreate_clienti=True -> DROP TABLE clienti e creazione nuova tabella con colonne richieste.
    Di default (False) mantiene la tabella clienti esistente e aggiunge eventuali colonne mancanti.
    """
    conn = get_connection()
    cur = conn.cursor()

    # altre tabelle (non toccate se esistono)
    cur.execute("""
    CREATE TABLE IF NOT EXISTS clienti_clean (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        CODICE TEXT,
        RAG_SOC TEXT,
        P_IVA TEXT,
        IND_VIA TEXT,
        COMUNE TEXT,
        PROV TEXT,
        CAP TEXT,
        MF TEXT,
        MARCA TEXT,
        LOGOTIPO TEXT,
        SCAD_ASS TEXT,
        IMP_ASS TEXT,
        IMP_RIP TEXT,
        TEL TEXT,
        COMMENTI TEXT,
        PASSAGGIO TEXT,
        TIPOMF TEXT,
        DATA TEXT,
        GRUPPO TEXT,
        DATADIS TEXT,
        CONTA TEXT,
        DM TEXT,
        CF TEXT,
        SCANNER TEXT,
        SCANNERAGG TEXT,
        N_PINPAD TEXT,
        N_PIN_BACK TEXT,
        PEC TEXT,
        FINEMF TEXT,
        CODPUNTOCA TEXT,
        RTREADY TEXT,
        CODUNIVOCO TEXT,
        GARANZIA TEXT,
        AZZ TEXT,
        QRCODE TEXT,
        LICENZA TEXT,
        PSW111 TEXT,
        PSWSCDIG TEXT
    );
    """)
    cur.execute("""
    CREATE TABLE IF NOT EXISTS macchine (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        matricola TEXT UNIQUE,
        modello TEXT,
        descrizione TEXT,
        data_installazione TEXT,
        mxaz TEXT,
        cliente_id INTEGER,
        disattivazione TEXT,
        vf TEXT,
        scadenza_assistenza TEXT,
        FOREIGN KEY(cliente_id) REFERENCES clienti(id)
    );
    """)
    cur.execute("""
    CREATE TABLE IF NOT EXISTS operazioni (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        macchina_id INTEGER,
        data_intervento TEXT,
        tipo TEXT,
        codice INTEGER,
        note TEXT,
        operatore TEXT
    );
    """)

    # INTAE: create simple table if not exists, poi aggiungi colonne mancanti in modo sicuro
    cur.execute("""
    CREATE TABLE IF NOT EXISTS intae (
        id INTEGER PRIMARY KEY AUTOINCREMENT
    );
    """)
    # leggi le colonne esistenti UNA SOLA VOLTA
    cur.execute("PRAGMA table_info(intae)")
    existing_intae = [r[1] for r in cur.fetchall()]
    for c in INTAE_COLUMNS:
        if c not in existing_intae:
            try:
                cur.execute(f'ALTER TABLE intae ADD COLUMN "{c}" TEXT')
            except sqlite3.OperationalError as e:
                # ignora esplicitamente errori di colonna duplicata
                if 'duplicate column name' in str(e).lower():
                    continue
                else:
                    raise

    # MISURART: stessa logica
    cur.execute("""
    CREATE TABLE IF NOT EXISTS misurart (
        id INTEGER PRIMARY KEY AUTOINCREMENT
    );
    """)
    cur.execute("PRAGMA table_info(misurart)")
    existing_mis = [r[1] for r in cur.fetchall()]
    for c in MISURART_COLUMNS:
        if c not in existing_mis:
            try:
                cur.execute(f'ALTER TABLE misurart ADD COLUMN "{c}" TEXT')
            except sqlite3.OperationalError as e:
                if 'duplicate column name' in str(e).lower():
                    continue
                else:
                    raise

    # crea se non esiste e poi aggiunge colonne mancanti
    _create_clienti_table(cur)
    _add_missing_clienti_columns(cur)

    conn.commit()
    conn.close()

# -------------------------
# Utility columns (se serve aggiungere)
# -------------------------

def get_clienti_columns() -> List[str]:
    conn = get_connection()
    cur = conn.cursor()
    cur.execute("PRAGMA table_info(clienti)")
    cols = [r[1] for r in cur.fetchall()]
    conn.close()
    return cols


def _ensure_clienti_columns(cols: List[str]):
    """Aggiunge colonne mancanti (TEXT). Utility generale, ma di norma preferiamo
    far riferimento a CLIENTI_COLUMNS fissi.
    """
    if not cols:
        return
    conn = get_connection()
    cur = conn.cursor()
    cur.execute("PRAGMA table_info(clienti)")
    existing = [r[1] for r in cur.fetchall()]
    for c in cols:
        if c not in existing:
            cur.execute(f'ALTER TABLE clienti ADD COLUMN "{c}" TEXT')
    conn.commit()
    conn.close()

# -------------------------
# CRUD Clienti
# -------------------------

def read_clienti_allcols(filter_term: Optional[str] = None, limit: int = 1000) -> Tuple[List[str], List[Tuple]]:
    conn = get_connection()
    cur = conn.cursor()
    cur.execute("PRAGMA table_info(clienti)")
    cols = [r[1] for r in cur.fetchall()]
    cols_select = ", ".join([f'"{c}"' for c in cols])
    if filter_term:
        like = f"%{filter_term}%"
        where_parts = [f'"{c}" LIKE ?' for c in cols if c.lower() != "id"]
        params = [like] * len(where_parts) + [limit]
        sql = f"SELECT {cols_select} FROM clienti WHERE " + " OR ".join(where_parts) + " ORDER BY id DESC LIMIT ?"
        cur.execute(sql, params)
    else:
        sql = f"SELECT {cols_select} FROM clienti ORDER BY id DESC LIMIT ?"
        cur.execute(sql, (limit,))
    rows = cur.fetchall()
    conn.close()
    return cols, rows

def insert_clients_from_map(data: dict, run_post_update: bool = True) -> int:
    """
    Inserisce un cliente usando l'ordine stabile CLIENTI_COLUMNS.
    Se run_post_update è True, esegue l'aggiornamento automatico FINEMF/GARANZIA
    per l'id appena creato usando la funzione batch aggiorna_finemf_e_garanzia_per_id_list.

    Ritorna l'id (nid) della riga appena inserita.
    """
    if not data:
        raise ValueError("data vuoto")

    conn = get_connection()
    try:
        cur = conn.cursor()

        # crea tabella se non esiste (con id autoincrement)
        cur.execute("SELECT name FROM sqlite_master WHERE type='table' AND lower(name)='clienti'")
        if not cur.fetchone():
            cols_def = ", ".join([f'"{c}" TEXT' for c in CLIENTI_COLUMNS])
            cur.execute(f'CREATE TABLE clienti (id INTEGER PRIMARY KEY AUTOINCREMENT, {cols_def})')

        # aggiorna la lista di colonne esistenti
        cur.execute("PRAGMA table_info(clienti)")
        existing = [r[1] for r in cur.fetchall()]

        # aggiungi eventuali colonne mancanti dalla lista standard
        for c in CLIENTI_COLUMNS:
            if c not in existing:
                cur.execute(f'ALTER TABLE clienti ADD COLUMN "{c}" TEXT')
                existing.append(c)

        # === VALIDAZIONE: campi obbligatori per inserimento MANUALE ===
        required_fields = ["RAG_SOC","P_IVA","IND_VIA","COMUNE","PROV","CAP","MF",
                           "MARCA","LOGOTIPO","MATRICOLA","SCAD_ASS","IMP_ASS","VF"]
        missing = [f for f in required_fields if not str(data.get(f, "")).strip()]
        if missing:
            raise ValueError(f"Impossibile inserire cliente: mancanti campi obbligatori: {', '.join(missing)}")

        # --- CHECK: impedisci inserimento se la stessa matricola è già presente ---
        matricola_input = str(data.get("MATRICOLA", "")).strip()
        if matricola_input:
            if matricola_exists(matricola_input, conn=conn):
                raise ValueError(f"Impossibile inserire cliente: la matricola '{matricola_input}' è già presente nel database.")

        # --- CHECK DUPLICATO COMPLETO: saltare se TUTTI i campi della riga sono identici a una riga esistente ---
        # Costruisci WHERE clause per cercare righe con tutti i campi uguali
        where_parts = []
        where_vals = []
        for col in CLIENTI_COLUMNS:
            val = str(data.get(col, "")).strip()
            if val:
                where_parts.append(f'"{col}"=?')
                where_vals.append(val)
        
        if where_parts:
            where_clause = " AND ".join(where_parts)
            sql_check = f'SELECT COUNT(*) FROM clienti WHERE {where_clause}'
            cur.execute(sql_check, where_vals)
            if int(cur.fetchone()[0] or 0) > 0:
                raise ValueError("Riga completamente duplicata: tutti i campi corrispondono a una riga esistente nel database.")

        # prepariamo l'inserimento usando SOLO CLIENTI_COLUMNS (in ordine)
        colnames_sql = ", ".join([f'"{c}"' for c in CLIENTI_COLUMNS])
        placeholders = ", ".join(["?"] * len(CLIENTI_COLUMNS))
        vals = [str(data.get(c, "")) for c in CLIENTI_COLUMNS]

        # Forza TIPOMF sempre a "RT" (sovrascrive qualsiasi valore fornito)
        try:
            idx_tipomf = CLIENTI_COLUMNS.index("TIPOMF")
            vals[idx_tipomf] = "RT"
        except ValueError:
            # se per qualche ragione TIPOMF non presente nella lista, ignoriamo
            pass
        
        # Aggiungi default MXAZ=2500 se non presente
        try:
            idx_mxaz = CLIENTI_COLUMNS.index("MXAZ")
            if not vals[idx_mxaz] or str(vals[idx_mxaz]).strip() == "":
                vals[idx_mxaz] = "2500"
        except ValueError:
            pass

        sql = f'INSERT INTO clienti ({colnames_sql}) VALUES ({placeholders})'
        cur.execute(sql, vals)

        nid = cur.lastrowid

        conn.commit()
    except Exception:
        conn.rollback()
        raise
    finally:
        # esegui ricalcolo CODPUNTOCA per la tabella clienti
        try:
            calcola_codpuntocassa("clienti")
        except Exception:
            # non interrompere l'inserimento se il ricalcolo fallisce
            pass

        conn.close()

    # se richiesto, esegui l'aggiornamento batch per questo singolo id
    if run_post_update and nid:
        try:
            # aggiorna solo i campi vuoti per l'id appena creato
            aggiorna_finemf_e_garanzia([nid])
        except Exception:
            # non interrompiamo l'inserimento in caso di errore nell'aggiornamento;
            # volendo possiamo loggare l'eccezione qui.
            pass

    return nid

def update_clients_by_id(rowid: int, updates: dict, run_post_update: bool = True) -> int:
    """
    Aggiorna il cliente con id `rowid` usando le coppie chiave/valore in `updates`.
    Se run_post_update è True (default) chiama la funzione batch
    `aggiorna_finemf_e_garanzia_per_id_list([rowid])` per popolare FINEMF/GARANZIA
    se vuoti.

    Ritorna il numero di righe aggiornate (int).
    """
    if not updates:
        return 0

    updates["TIPOMF"] = "RT"
    
    # Assicuriamoci che le colonne esistano (funzione helper che già usi)
    cols = list(updates.keys())
    _ensure_clienti_columns(cols)

    parts = ", ".join([f'"{c}"=?' for c in cols])
    vals = [str(updates[c]) for c in cols] + [rowid]

    conn = get_connection()
    try:
        cur = conn.cursor()
        cur.execute(f'UPDATE clienti SET {parts} WHERE id=?', vals)
        affected = cur.rowcount if hasattr(cur, "rowcount") else 0
        conn.commit()
    except Exception:
        conn.rollback()
        raise
    finally:
        conn.close()

    # Esegui il post-update batch (aggiorna solo valori vuoti)
    if run_post_update and rowid:
        try:
            aggiorna_finemf_e_garanzia([rowid])
        except Exception:
            # non facciamo fallire l'update principale se il post-update fallisce;
            pass

    return int(affected or 0)

def delete_clients(rowid: int):
    conn = get_connection()
    cur = conn.cursor()
    cur.execute("DELETE FROM macchine WHERE cliente_id=?", (rowid,))
    cur.execute("DELETE FROM clienti WHERE id=?", (rowid,))
    conn.commit()
    conn.close()


def duplicate_clients(rowid: int, new_suffix: str = "_COPIA") -> Optional[int]:
    conn = get_connection()
    cur = conn.cursor()
    cur.execute("SELECT * FROM clienti WHERE id=?", (rowid,))
    row = cur.fetchone()
    if not row:
        conn.close()
        return None
    cur.execute("PRAGMA table_info(clienti)")
    cols = [r[1] for r in cur.fetchall()]
    values = []
    cols_noid = []
    
    for i, col in enumerate(cols):
        if col == "id":
            continue
        val = row[i]
        # se duplico RAG_SOC aggiungo il suffisso
        if col == "RAG_SOC" and val:
            val = f"{val}{new_suffix}"
        # Forza TIPOMF su copia
        if col == "TIPOMF":
            val = "RT"
        values.append(val if val is not None else "")
        cols_noid.append(col)

    placeholders = ",".join(["?"] * len(cols_noid))
    colnames_sql = ",".join([f'"{c}"' for c in cols_noid])
    sql = f'INSERT INTO clienti ({colnames_sql}) VALUES ({placeholders})'
    cur.execute(sql, values)
    nid = cur.lastrowid
    conn.commit()
    conn.close()
    return nid

def sample_righe_vuoti(ids: Optional[List[int]] = None, limit: int = 20):
    """
    Restituisce una lista di tuple (id, RAG_SOC, MATRICOLA, DATAINST, MXAZ, FINEMF, GARANZIA)
    per le righe che hanno FINEMF o GARANZIA vuoti. Se ids forniti filtra solo quelli.
    """
    conn = get_connection()
    cur = conn.cursor()
    empty_cond = "(FINEMF IS NULL OR FINEMF = '' OR FINEMF = '0' OR GARANZIA IS NULL OR GARANZIA = '' OR GARANZIA = '0')"
    if ids:
        placeholders = ",".join("?" for _ in ids)
        q = f"SELECT id, RAG_SOC, MATRICOLA, DATAINST, MXAZ, FINEMF, GARANZIA FROM clienti WHERE id IN ({placeholders}) AND {empty_cond} LIMIT ?"
        cur.execute(q, tuple(ids) + (limit,))
    else:
        q = f"SELECT id, RAG_SOC, MATRICOLA, DATAINST, MXAZ, FINEMF, GARANZIA FROM clienti WHERE {empty_cond} LIMIT ?"
        cur.execute(q, (limit,))
    rows = cur.fetchall()
    conn.close()   


def _normalize_matricola_for_compare(s):
    """Rimuove caratteri non alfanumerici e fa lower-case per confronto."""
    return re.sub(r'[^0-9A-Za-z]', '', str(s or "")).lower().strip()

def _parse_date_param_to_iso(s):
    """
    Accetta:
      - None -> oggi (YYYY-MM-DD)
      - datetime.date / datetime.datetime -> ISO date (YYYY-MM-DD)
      - stringa 'DDMMYYYY' (8 cifre) -> ISO date
      - stringa 'YYYY-MM-DD' -> pass-through se valida
      - stringa 'DD/MM/YYYY' -> convert
    Ritorna stringa ISO 'YYYY-MM-DD'
    """
    if s is None:
        return date.today().isoformat()
    if isinstance(s, datetime):
        return s.date().isoformat()
    if isinstance(s, date):
        return s.isoformat()
    s_str = str(s).strip()
    # 8 cifre DDMMYYYY
    if re.fullmatch(r'\d{8}', s_str):
        day = int(s_str[0:2]); month = int(s_str[2:4]); year = int(s_str[4:8])
        try:
            return datetime(year, month, day).date().isoformat()
        except ValueError:
            raise ValueError(f"Data non valida: {s_str}")
    # DD/MM/YYYY
    try:
        if re.fullmatch(r'\d{2}/\d{2}/\d{4}', s_str):
            dt = datetime.strptime(s_str, "%d/%m/%Y")
            return dt.date().isoformat()
    except Exception:
        pass
    # YYYY-MM-DD
    try:
        dt = datetime.fromisoformat(s_str)
        return dt.date().isoformat()
    except Exception:
        pass
    # fallback: error
    raise ValueError(f"Formato data non riconosciuto: {s_str}")

def deactivate_machine(matricola, applicazione_date=None):
    """
    Disattiva una macchina identificata da `matricola`.
    Per tutte le tabelle che contengono la colonna MATRICOLA:
      - imposta DATADIST = applicazione_date (YYYY-MM-DD)
      - azzera VF (0 se colonna numerica, '' se testo)
    Inoltre, crea automaticamente un record in MISURART con i dati della disattivazione.
    Ritorna dict: { 'table_name': rows_updated, ..., 'misurart_created': True/False }
    Solleva ValueError se matricola vuota o data non valida.
    """
    if not matricola or not str(matricola).strip():
        raise ValueError("Devi fornire una matricola non vuota.")

    mnorm = _normalize_matricola_for_compare(matricola)
    if not mnorm:
        raise ValueError("Matricola non valida dopo normalizzazione.")

    # normalizza data
    iso_date = _parse_date_param_to_iso(applicazione_date)
    # Converti in formato DD/MM/YYYY per MISURART
    try:
        dt_obj = datetime.strptime(iso_date, "%Y-%m-%d").date()
        date_display = dt_obj.strftime("%d/%m/%Y")
    except Exception:
        date_display = iso_date

    conn = get_connection()
    cur = conn.cursor()
    updated_counts = {}
    misurart_created = False

    try:
        # Recupera i dati del cliente per compilare MISURART
        client_data = {}
        try:
            # Verifica che DATADIS sia aggiornato nel database
            cur.execute("""
                SELECT RAG_SOC, P_IVA, COMUNE, PROV, IND_VIA, CAP, MARCA, MF, LOGOTIPO
                FROM clienti 
                WHERE lower(trim(MATRICOLA)) = ?
                LIMIT 1
            """, (mnorm,))
            client_row = cur.fetchone()
            if client_row:
                client_data = {
                    "RAG_SOC": client_row[0] or "",
                    "P_IVA": client_row[1] or "",
                    "COMUNE": client_row[2] or "",
                    "PROV": client_row[3] or "",
                    "IND_VIA": client_row[4] or "",
                    "CAP": client_row[5] or "",
                    "MARCA": client_row[6] or "",
                    "MF": client_row[7] or "",
                    "LOGOTIPO": client_row[8] or ""
                }
        except Exception as e:
            print(f"[WARN] Impossibile recuperare dati cliente per MISURART: {e}")

        # prendi tutte le tabelle nel DB
        cur.execute("SELECT name FROM sqlite_master WHERE type='table'")
        tables = [r[0] for r in cur.fetchall()]

        for tbl in tables:
            # leggi schema della tabella
            try:
                cur.execute(f"PRAGMA table_info({tbl})")
                cols_info = cur.fetchall()  # (cid, name, type, notnull, dflt_value, pk)
                col_names = [c[1].upper() for c in cols_info]
                col_types = {c[1].upper(): (c[2] or "").upper() for c in cols_info}
            except Exception:
                continue

            # procediamo solo se la tabella ha MATRICOLA e almeno una delle colonne DATADIS o VF
            if "MATRICOLA" not in col_names:
                continue
            has_datadis = "DATADIS" in col_names or "DATADIST" in col_names  # Gestisci entrambe le varianti
            has_vf = "VF" in col_names

            if not has_datadis and not has_vf:
                # niente da aggiornare in questa tabella
                continue

            # costruiamo SET dinamico
            set_parts = []
            params = []

            if has_datadis:
                # Usa DATADIS se esiste, altrimenti DATADIST
                datadis_col = "DATADIS" if "DATADIS" in col_names else "DATADIST"
                set_parts.append(f"{datadis_col} = ?")
                params.append(iso_date)

            if has_vf:
                # decide se VF è numerico (INTEGER/REAL) -> 0 altrimenti stringa vuota
                vf_type = col_types.get("VF", "")
                if any(t in vf_type for t in ("INT", "REAL", "NUM", "DECIMAL", "FLOAT", "DOUBLE")):
                    set_parts.append("VF = 0")
                else:
                    set_parts.append("VF = ''")

            if not set_parts:
                continue

            sql_set = ", ".join(set_parts)
            params.append(mnorm)
            
            sql_update = f"""
                UPDATE {tbl}
                SET {sql_set}
                WHERE lower(MATRICOLA) = ?
            """

            try:
                cur.execute(sql_update, params)
                count = cur.rowcount
                updated_counts[tbl] = count
            except Exception as e:
                print(f"[WARN] Errore aggiornamento tabella {tbl}: {e}")
                continue

        # NUOVO: Crea record in MISURART per la disattivazione
        if client_data:
            print(f"[INFO] Disattivazione di {matricola}: DATADIS={date_display}")
            # Recupera i dati più recenti da INTAE per questa matricola
            intae_data = {}
            try:
                # Usa sia la matricola normalizzata che la matricola con logotipo
                search_patterns = [
                    f"%{mnorm}%",  # senza logotipo
                    f"%{str(matricola).strip()}%"  # con logotipo
                ]
                for pattern in search_patterns:
                    cur.execute("""
                        SELECT CODICE_INTERVENTO, DATA_INTERVENTO, CF_TECNICO, DISPOSITIVO, TIPOLOGIA, VERSIONE_SW
                        FROM intae 
                        WHERE upper(trim(MATRICOLA)) LIKE upper(?)
                        ORDER BY DATA_INTERVENTO DESC 
                        LIMIT 1
                    """, (pattern,))
                    intae_row = cur.fetchone()
                    if intae_row:
                        intae_data = {
                            "CODICE_INTERVENTO": intae_row[0] or "",
                            "DATA_INTERVENTO": intae_row[1] or "",
                            "CF_TECNICO": intae_row[2] or "",
                            "DISPOSITIVO": intae_row[3] or "",
                            "TIPOLOGIA": intae_row[4] or "",
                            "VERSIONE_SW": intae_row[5] or ""
                        }
                        break
            except Exception as e:
                print(f"[WARN] Impossibile recuperare dati INTAE per MISURART: {e}")
            
            try:
                misurart_data = {
                    "TIPO": "2",  # Tipo 2 = Disattivazione
                    "DATAINIZIO": date_display,
                    "DATAFINE": date_display,
                    "ESITO": "DISATTIVATO",
                    "PIVAESERCE": client_data.get("P_IVA", ""),
                    "RAGSOCESE": client_data.get("RAG_SOC", ""),
                    "COMUNESE": client_data.get("COMUNE", ""),
                    "PROVINCIA": client_data.get("PROV",  ""),
                    "INDIRIZZO": client_data.get("IND_VIA", ""),
                    "CAPESE": client_data.get("CAP", ""),
                    "MARCAMF": client_data.get("MARCA", ""),
                    "MODELLO": client_data.get("MF", ""),
                    "LOGOTIPO": client_data.get("LOGOTIPO", ""),
                    "MATRICOLA": str(matricola).strip().upper(),
                    "CF": intae_data.get("CF_TECNICO", ""),  # Verificatore da INTAE
                    # CODICE_INTERVENTO non è in MISURART_COLUMNS, non va inserito
                    "TIPOINT": "2"  # Tipo intervento 2 = Disattivazione
                }
                
                # Inserisci in MISURART
                cols_to_insert = [c for c in MISURART_COLUMNS if c in misurart_data]
                vals_to_insert = [misurart_data[c] for c in cols_to_insert]
                placeholders = ",".join("?" * len(cols_to_insert))
                cols_sql = ", ".join([f'"{c}"' for c in cols_to_insert])
                
                sql_insert = f"INSERT INTO misurart ({cols_sql}) VALUES ({placeholders})"
                cur.execute(sql_insert, vals_to_insert)
                misurart_created = True
                print(f"[INFO] Record MISURART creato per disattivazione matricola {matricola}")
                
            except Exception as e:
                print(f"[ERROR] Impossibile creare record MISURART: {e}")
                misurart_created = False

        conn.commit()
        updated_counts["misurart_created"] = misurart_created

    except Exception as e:
        conn.rollback()
        raise e
    finally:
        conn.close()

    return updated_counts


def update_clienti_datadist_vf(matricola, iso_date, conn=None):
    """
    Aggiorna nella tabella clienti le righe relative alla `matricola` impostando:
      - DATADIST = iso_date (es. '2025-12-21')
      - VF = 0 (se VF è numerico) oppure '' (se VF è testo)
    Ritorna il numero di righe aggiornate.
    Se viene passato `conn` (connessione sqlite attiva), la usa; altrimenti apre/chiude una connessione propria.
    """
    if not matricola or not str(matricola).strip():
        raise ValueError("Matricola richiesta")

    m_input = str(matricola).strip()
    m_norm = _normalize_matricola_for_compare(m_input)

    own_conn = False
    if conn is None:
        conn = get_connection()
        own_conn = True
    cur = conn.cursor()

    try:
        # controlla se esistono le colonne DATADIST / VF nella tabella clienti
        cur.execute("PRAGMA table_info(clienti)")
        cols = [(r[1].upper(), (r[2] or "").upper()) for r in cur.fetchall()]
        col_names = [c[0] for c in cols]
        if "DATADIST" not in col_names and "DATADIS" not in col_names:
            # se nessuna delle due colonne esiste non c'è niente da fare
            return 0

        # individua il nome esatto di DATADIST (gestisce eventuale DATADIS)
        datadist_col = "DATADIST" if "DATADIST" in col_names else "DATADIS"
        has_vf = "VF" in col_names

        # prima prova un confronto semplice case-insensitive
        m_lower = m_input.lower()
        cur.execute(f"SELECT COUNT(*) FROM clienti WHERE lower(matricola) = ?", (m_lower,))
        cnt_simple = int(cur.fetchone()[0] or 0)

        if cnt_simple > 0:
            # aggiorneremo con WHERE lower(matricola) = ?
            params = [iso_date]
            if has_vf:
                # decide valore azzeramento per tipo colonna VF
                vf_type = dict(cols).get("VF", "")
                if any(t in vf_type for t in ("INT", "REAL", "NUM", "DECIMAL", "FLOAT", "DOUBLE")):
                    vf_zero = 0
                else:
                    vf_zero = ""
                params.append(vf_zero)
            params.append(m_lower)
            sql = f"UPDATE clienti SET {datadist_col} = ?" + (", VF = ?" if has_vf else "") + " WHERE lower(matricola) = ?"
            cur.execute(sql, tuple(params))
            updated = cnt_simple
        else:
            # non trovato con confronto semplice: cerchiamo tutti i record e normalizziamo sul lato Python
            cur.execute("SELECT rowid, matricola FROM clienti")
            candidates = cur.fetchall()  # lista di (rowid, matricola)
            matched_ids = []
            for rowid, m_db in candidates:
                m_db_norm = _normalize_matricola_for_compare(m_db)
                if m_db_norm and m_db_norm == m_norm:
                    matched_ids.append(rowid)
            if not matched_ids:
                return 0
            
            # aggiorna per id (rowid) per sicurezza
            vf_type = dict(cols).get("VF", "")
            if any(t in vf_type for t in ("INT", "REAL", "NUM", "DECIMAL", "FLOAT", "DOUBLE")):
                vf_zero = 0
            else:
                vf_zero = ""
            
            # chunk update se troppi id
            updated = 0
            chunk_size = 200
            for i in range(0, len(matched_ids), chunk_size):
                chunk = matched_ids[i:i+chunk_size]
                placeholders = ",".join(["?"] * len(chunk))  # <- DEFINISCI QUI
                if has_vf:
                    sql = f"UPDATE clienti SET {datadist_col} = ?, VF = ? WHERE rowid IN ({placeholders})"
                    params = [iso_date, vf_zero] + chunk
                else:
                    sql = f"UPDATE clienti SET {datadist_col} = ? WHERE rowid IN ({placeholders})"
                    params = [iso_date] + chunk
                cur.execute(sql, tuple(params))
                updated += cur.rowcount if cur.rowcount is not None and cur.rowcount > 0 else len(chunk)

        if own_conn:
            conn.commit()
        return updated
    finally:
        if own_conn:
            try:
                conn.close()
            except Exception:
                pass

# -------------------------
# Import Clients CSV / XLS (robusto)
# -------------------------
def _detect_delimiter_and_dialect(sample: str, delimiters: str = ",;|\t") -> str:
    try:
        sniffer = csv.Sniffer()
        dialect = sniffer.sniff(sample, delimiters=delimiters)
        return dialect.delimiter
    except Exception:
        # fallback: try common delimiters manually
        for d in [",", ";", "\t", "|"]:
            if sample.count(d) >= 1:
                return d
    return ","  # ultima risorsa

def _normalize_str_for_compare(s):
    """Normalizza stringa: None -> '', rimuove accenti, spazi multipli e mette tutto MAIUSCOLO."""
    if s is None:
        return ""
    s = str(s).strip().upper()
    s = unicodedata.normalize("NFKD", s)
    s = "".join(ch for ch in s if not unicodedata.combining(ch))
    s = re.sub(r"\s+", " ", s)
    return s

def _normalize_piva(p):
    """Normalizza partita IVA: rimuove simboli e spazi e mette MAIUSCOLO."""
    if p is None:
        return ""
    return re.sub(r"\W+", "", str(p).upper())

def cliente_esiste(conn, rag_soc, p_iva, matricola):
    """
    Restituisce True se esiste nel DB una riga con:
      - RAG_SOC normalizzata uguale a rag_soc
      - P_IVA normalizzata uguale a p_iva
      - MATRICOLA normalizzata uguale a matricola
    Tutti e tre i campi devono corrispondere per considerare duplicato.
    """
    cur = conn.cursor()
    rag_n = _normalize_str_for_compare(rag_soc)
    piva_n = _normalize_piva(p_iva)
    mat_n = _normalize_str_for_compare(matricola)

    # Se uno dei tre non è presente (stringa vuota), consideriamo che non possiamo
    # fare il controllo "tutti e tre" e quindi non è duplicato secondo questa regola.
    if not (rag_n and piva_n and mat_n):
        return False

    sql = """
    SELECT 1 FROM clienti
    WHERE UPPER(RAG_SOC) = ?
      AND REPLACE(UPPER(P_IVA), ' ', '') = ?
      AND UPPER(MATRICOLA) = ?
    LIMIT 1
    """
    cur.execute(sql, (rag_n, piva_n, mat_n))
    return cur.fetchone() is not None

def matricola_exists(matricola: str, conn=None, exclude_id: Optional[int] = None) -> bool:
    """
    Restituisce True se esiste almeno una riga in clienti con la stessa matricola normalizzata.
    Se exclude_id fornito, ignora quella riga (utile per update).
    Accetta una connessione opzionale (conn) per riutilizzare la stessa transaction.
    """
    if not matricola or not str(matricola).strip():
        return False

    m_norm = _normalize_matricola_for_compare(matricola)
    if not m_norm:
        return False

    own_conn = False
    if conn is None:
        conn = get_connection()
        own_conn = True
    cur = conn.cursor()
    try:
        # prova rapido con lower(matricola) per limitare i risultati, poi confronto normalizzato in Python
        cur.execute("SELECT rowid, matricola FROM clienti WHERE lower(matricola) = ?", (str(matricola).strip().lower(),))
        rows = cur.fetchall()
        for row in rows:
            rid, m_db = row[0], row[1]
            if exclude_id is not None:
                try:
                    if int(rid) == int(exclude_id):
                        continue
                except Exception:
                    pass
            if _normalize_matricola_for_compare(m_db) == m_norm:
                return True

        # fallback: se non trovati con where, scansiona poche righe (rare)
        cur.execute("SELECT rowid, matricola FROM clienti")
        for rid, m_db in cur.fetchall():
            if exclude_id is not None:
                try:
                    if int(rid) == int(exclude_id):
                        continue
                except Exception:
                    pass
            if _normalize_matricola_for_compare(m_db) == m_norm:
                return True

        return False
    finally:
        if own_conn:
            try:
                conn.close()
            except Exception:
                pass

def validate_required_columns(file_headers: list, required_columns: list, file_type: str = "generico") -> dict:
    """
    Valida che i file di import contengono tutte le colonne richieste.
    Ritorna dict {ok, missing_columns, message}
    """
    def normalize_col(col: str) -> str:
        return str(col).strip().upper()
    
    normalized_file = {normalize_col(h): h for h in file_headers}
    normalized_required = [normalize_col(c) for c in required_columns]
    
    missing = []
    for req in normalized_required:
        if req not in normalized_file:
            missing.append(req)
    
    if missing:
        return {
            "ok": False,
            "missing_columns": missing,
            "message": f"[{file_type}] File di importazione non contiene le colonne obbligatorie: {', '.join(missing)}"
        }
    
    return {
        "ok": True,
        "missing_columns": [],
        "message": f"[{file_type}] Tutte le colonne obbligatorie sono presenti"
    }

def import_clients_from_csv_generic(path: str, delimiter: str | None = None, encoding: str = "latin-1", batch_commit: int = 500) -> dict:
    """
    Import CSV in modo robusto:
      - rileva le colonne reali della tabella clienti usando PRAGMA
      - costruisce INSERT dinamico solo sulle colonne trovate
      - calcola GARANZIA e FINEMF prima dell'inserimento quando mancanti
      - Ritorna dict con: {ok, inserted, skipped, skipped_rows, columns, message}
    """

    if not os.path.exists(path):
        raise FileNotFoundError(path)

    def _is_empty(v):
        return v is None or str(v).strip() == "" or str(v).strip() == "0"

    def normalize_header(h: str) -> str:
        # prende solo la parte prima della virgola
        return (h.split(",")[0].strip().upper() if h is not None else "")

    def normalize_value_for_db(raw):
        if raw is None:
            return None
        s = str(raw).strip()
        if s == "":
            return None
        return s

    # lettura sample per autodetect delimiter se non passato
    with open(path, "r", encoding=encoding, errors="replace") as sf:
        sample = sf.read(8192)
        sf.seek(0)
        if delimiter is None or delimiter == "":
            # preferenze: ; , \t |
            for d in (";", ",", "\t", "|"):
                if sample.count(d) > 0:
                    used_delim = d
                    break
            else:
                used_delim = ","
        else:
            used_delim = delimiter

    # apri reader
    with open(path, "r", encoding=encoding, errors="replace") as fh:
        reader = csv.DictReader(fh, delimiter=used_delim)
        headers = reader.fieldnames or []

        # ottieni colonne reali della tabella 'clienti' direttamente dal DB
        conn = get_connection()
        cur = conn.cursor()
        try:
            cur.execute("PRAGMA table_info(clienti)")
            pragma = cur.fetchall()  # (cid,name,type,notnull,dflt_value,pk)
            db_cols = [row[1] for row in pragma]
        except Exception:
            # se pragma fallisce, fallback a CLIENTI_COLUMNS se definita
            try:
                db_cols = CLIENTI_COLUMNS.copy()
            except Exception:
                db_cols = []
        if not db_cols:
            conn.close()
            raise RuntimeError("Impossibile determinare colonne tabella 'clienti' (PRAGMA fallito e CLIENTI_COLUMNS non definita)")

        # VALIDAZIONE: richiedi almeno la ragione sociale
        required_fields = ["RAG_SOC"]
        validation = validate_required_columns(headers, required_fields, "CLIENTI (import CSV)")
        if not validation["ok"]:
            conn.close()
            return {"ok": False, "inserted": 0, "skipped": 0, "skipped_rows": [], "columns": CLIENTI_COLUMNS, "message": validation["message"], "error": validation["message"]}

        # mappa headers del CSV ai nomi colonne DB (matching normalizzato)
        norm_to_db = {}
        norm_db = { normalize_header(c): c for c in db_cols }
        for h in headers:
            if h is None:
                continue
            nh = normalize_header(h)
            target = norm_db.get(nh)
            if target:
                norm_to_db[h] = target
            else:
                # prova a trovare con matching parziale
                for k_norm, k_real in norm_db.items():
                    if nh == k_norm or nh in k_norm or k_norm in nh:
                        norm_to_db[h] = k_real
                        break
                # altrimenti non mappata -> verrà ignorata

        # colonne DB che useremo (quelle mappate dal CSV)
        csv_to_db_cols = list(dict.fromkeys(norm_to_db.values()))  # preserva unicità e ordine di prima occorrenza
        if not csv_to_db_cols:
            # fallback: prova quick match by name (esatti)
            if all(h and h.upper() in [c.upper() for c in db_cols] for h in headers):
                csv_to_db_cols = [h for h in headers if h and h.upper() in [c.upper() for c in db_cols]]
            else:
                conn.close()
                raise RuntimeError(f"Nessuna colonna del CSV mappabile alle colonne DB. headers: {headers} db_cols: {db_cols}")

        # ordina colonne secondo ordine DB (per consistenza)
        columns_to_insert = sorted(csv_to_db_cols, key=lambda x: db_cols.index(x) if x in db_cols else 0)

        # === VALIDAZIONE COLONNE OBBLIGATORIE ===
        # Verifica che il file contenga almeno RAG_SOC
        required_columns = ["RAG_SOC"]
        missing_columns = [col for col in required_columns if col not in columns_to_insert]
        if missing_columns:
            conn.close()
            raise RuntimeError(f"Il file non contiene le colonne obbligatorie: {', '.join(missing_columns)}. "
                             f"Colonne trovate nel file: {', '.join(csv_to_db_cols)}")

        # Assicuriamoci che TIPOMF sia sempre presente nell'INSERT (lo impostiamo a "RT" per ogni riga)
        if "TIPOMF" not in columns_to_insert:
            columns_to_insert.append("TIPOMF")

        placeholders = ", ".join(["?"] * len(columns_to_insert))
        colnames_sql = ", ".join([f'"{c}"' for c in columns_to_insert])
        insert_sql = f'INSERT INTO clienti ({colnames_sql}) VALUES ({placeholders})'

        inserted = 0
        skipped = 0
        to_commit = 0
        rows_processed = 0
        skipped_rows = []  # Lista per raccogliere le righe saltate per errori
        duplicate_rows = []  # Lista per raccogliere i duplicati
        seen_in_file = {}  # Traccia righe già viste nel file corrente (chiave -> numero riga)


        for r in reader:
            rows_processed += 1

            # costruisci valori in ordine columns_to_insert
            vals = []
            for col in columns_to_insert:
                # trova header originale che mappa a questa colonna (prendo la prima chiave)
                header_keys = [h for h, mapped in norm_to_db.items() if mapped == col]
                if header_keys:
                    raw = r.get(header_keys[0], "")
                else:
                    raw = ""
                vals.append(normalize_value_for_db(raw))

            # Forza TIPOMF su ogni riga
            try:
                tip_idx = columns_to_insert.index("TIPOMF")
                vals[tip_idx] = "RT"
            except ValueError:
                pass

            # helper per recuperare valore dato nome colonna DB
            def val_for(colname):
                try:
                    idx = columns_to_insert.index(colname)
                    return vals[idx]
                except ValueError:
                    return None

            # Controlla se la riga è completamente vuota (tutti i campi NULL o vuoti)
            has_any_value = any(v is not None and str(v).strip() != "" for v in vals)
            
            if not has_any_value:
                # Riga completamente vuota, skippa sempre
                skipped += 1
                skipped_rows.append(vals)
                print(f"[IMPORT] Riga {rows_processed} completamente vuota, skippata")
                continue
            
            # Controllo duplicati: prima controlla se già visto nel file corrente
            # Crea una firma della riga basata sui campi non-nulli
            row_signature_parts = []
            for i, col in enumerate(columns_to_insert):
                val = vals[i]
                if val is not None and str(val).strip() != "":
                    row_signature_parts.append(f"{col}:{str(val).strip().lower()}")
            
            row_signature = "||".join(sorted(row_signature_parts))
            
            # Debug: mostra firma ogni 100 righe
            if rows_processed % 100 == 0:
                print(f"[IMPORT DEBUG] Riga {rows_processed}: firma={row_signature[:100]}... | duplicati trovati finora={len(duplicate_rows)}")
            
            # Se già visto nel file corrente, è un duplicato interno
            if row_signature in seen_in_file:
                skipped += 1
                duplicate_rows.append(vals)
                if len(duplicate_rows) <= 10:  # Mostra solo i primi 10 per non intasare
                    print(f"[IMPORT] Riga {rows_processed} DUPLICATO rispetto a riga {seen_in_file[row_signature]} nel file")
                continue
            
            # Registra la firma per rilevare duplicati successivi nel file
            seen_in_file[row_signature] = rows_processed
            
            # Controllo duplicati: confronta tutti i campi NON-NULL della riga nel DB
            # Costruisce WHERE dinamico basato sui campi che hanno valori
            where_parts = []
            where_vals = []
            
            for i, col in enumerate(columns_to_insert):
                val = vals[i]
                if val is not None and str(val).strip() != "":
                    # Campo con valore: aggiungi al WHERE
                    where_parts.append(f'LOWER(TRIM("{col}")) = LOWER(TRIM(?))')
                    where_vals.append(str(val).strip())
                else:
                    # Campo vuoto: verifica che sia NULL o vuoto anche nel DB
                    where_parts.append(f'("{col}" IS NULL OR TRIM("{col}") = "")')
            
            if where_parts:
                where_clause = " AND ".join(where_parts)
                check_sql = f"SELECT COUNT(*) FROM clienti WHERE {where_clause}"
                
                try:
                    cur.execute(check_sql, where_vals)
                    count = int(cur.fetchone()[0] or 0)
                    if count > 0:
                        # Riga già presente (duplicato esatto), skippa
                        skipped += 1
                        duplicate_rows.append(vals)  # Aggiungi ai duplicati, non agli skipped generici
                        print(f"[IMPORT] Riga {rows_processed} duplicata rispetto al database")
                        continue
                except Exception as e:
                    # Se il check fallisce, procedi comunque (evita blocchi)
                    print(f"[IMPORT WARN] Errore nel check duplicati riga {rows_processed}: {e}")
                    pass

            # Calcola GARANZIA se presente nella columns_to_insert e vuota
            if "GARANZIA" in columns_to_insert:
                gi = columns_to_insert.index("GARANZIA")
                if _is_empty(vals[gi]):
                    datainst_val = val_for("DATAINST")
                    durata_val = ""
                    for cand in ("DURATA_GAR", "DURATA", "DURATA_MESI", "DURATA_GARANZIA"):
                        v = val_for(cand)
                        if not _is_empty(v):
                            durata_val = v
                            break
                    try:
                        nuova_gar = calcola_scadenza_garanzia(datainst_val, durata_val)
                        if nuova_gar:
                            vals[gi] = nuova_gar
                    except Exception:
                        pass

            # Calcola FINEMF se presente e vuota
            if "FINEMF" in columns_to_insert:
                fi = columns_to_insert.index("FINEMF")
                if _is_empty(vals[fi]):
                    datainst_val = val_for("DATAINST")
                    mxaz_val = val_for("MXAZ")
                    try:
                        nuovo_fin = calcola_fine_memoria(datainst_val, mxaz_val)
                        if nuovo_fin:
                            vals[fi] = nuovo_fin
                    except Exception:
                        pass

            # Esegui insert parametrico
            try:
                cur.execute(insert_sql, vals)
                inserted += 1
                to_commit += 1
            except Exception as e:
                print(f"[IMPORT ERROR] riga {rows_processed} -> errore insert: {e}. sql: {insert_sql}, vals sample: {vals[:6]}")
                skipped += 1
                skipped_rows.append(vals)
                conn.rollback()
                continue

            # commit a blocchi
            if to_commit >= batch_commit:
                conn.commit()
                to_commit = 0

        # commit rimanenti
        if to_commit > 0:
            conn.commit()

        # opzionale: dopo l'import massivo esegui la funzione di aggiornamento per sicurezza
        try:
            _ = aggiorna_finemf_e_garanzia()
        except Exception:
            pass

        # opzionale: ricalcola CODPUNTOCA sia su clienti che su clienti_clean
        try:
            calcola_codpuntocassa("clienti")
        except Exception:
            pass
        try:
            calcola_codpuntocassa("clienti_clean")
        except Exception:
            pass

        conn.close()
        
        # Salva le righe saltate (per errori) in un file .txt nella stessa directory del file importato
        skipped_file_path = None
        duplicate_file_path = None
        
        if skipped_rows:
            try:
                import_dir = os.path.dirname(os.path.abspath(path))
                import_name = os.path.splitext(os.path.basename(path))[0]
                skipped_file_path = os.path.join(import_dir, f"{import_name}_SKIPPED_ROWS.txt")
                
                with open(skipped_file_path, "w", encoding="utf-8") as sf:
                    sf.write("RIGHE SALTATE DURANTE L'IMPORT (ERRORI)\n")
                    sf.write("=" * 80 + "\n")
                    sf.write(f"File: {path}\n")
                    sf.write(f"Data: {datetime.now().strftime('%d/%m/%Y %H:%M:%S')}\n")
                    sf.write(f"Saltate per errori: {len(skipped_rows)} righe\n")
                    sf.write("=" * 80 + "\n\n")
                    
                    # Scrivi i dati delle righe saltate in formato tabulare
                    if columns_to_insert:
                        # Header
                        sf.write("Colonne: " + " | ".join(columns_to_insert) + "\n")
                        sf.write("-" * 80 + "\n")
                        
                        # Righe
                        for row_vals in skipped_rows:
                            row_str = " | ".join(str(v) if v is not None else "" for v in row_vals)
                            sf.write(row_str + "\n")
                
                print(f"[IMPORT] File righe saltate salvato: {skipped_file_path}")
            except Exception as e:
                print(f"[IMPORT WARN] Errore nel salvataggio file righe saltate: {e}")
        
        # Salva i duplicati in un file .txt separato
        if duplicate_rows:
            try:
                import_dir = os.path.dirname(os.path.abspath(path))
                import_name = os.path.splitext(os.path.basename(path))[0]
                duplicate_file_path = os.path.join(import_dir, f"{import_name}_DUPLICATES.txt")
                
                with open(duplicate_file_path, "w", encoding="utf-8") as df:
                    df.write("RIGHE DUPLICATE DURANTE L'IMPORT\n")
                    df.write("=" * 80 + "\n")
                    df.write(f"File: {path}\n")
                    df.write(f"Data: {datetime.now().strftime('%d/%m/%Y %H:%M:%S')}\n")
                    df.write(f"Duplicate: {len(duplicate_rows)} righe\n")
                    df.write("=" * 80 + "\n\n")
                    
                    # Scrivi i dati dei duplicati in formato tabulare
                    if columns_to_insert:
                        # Header
                        df.write("Colonne: " + " | ".join(columns_to_insert) + "\n")
                        df.write("-" * 80 + "\n")
                        
                        # Righe
                        for row_vals in duplicate_rows:
                            row_str = " | ".join(str(v) if v is not None else "" for v in row_vals)
                            df.write(row_str + "\n")
                
                print(f"[IMPORT] File duplicati salvato: {duplicate_file_path}")
            except Exception as e:
                print(f"[IMPORT WARN] Errore nel salvataggio file duplicati: {e}")
        
        print(f"[IMPORT] file={path!r} delimiter={used_delim!r} headers={headers} rows_processed={rows_processed} inserted={inserted} skipped={skipped}")
        return {
            "ok": True,
            "inserted": inserted,
            "skipped": skipped,
            "skipped_rows": skipped_rows,
            "duplicate_rows": duplicate_rows,
            "skipped_file": skipped_file_path,
            "duplicate_file": duplicate_file_path,
            "columns": columns_to_insert,
            "message": f"Import completato: {inserted} importati, {skipped} saltati"
        }

def import_excel_clienti_via_soffice(path: str) -> int:
    if not os.path.exists(path):
        raise FileNotFoundError(path)
    tmpdir = tempfile.mkdtemp()
    try:
        # convert to CSV
        cmd = ["soffice", "--headless", "--convert-to", "csv:Text - txt - csv (StarCalc)", path, "--outdir", tmpdir]
        subprocess.run(cmd, check=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        files = [f for f in os.listdir(tmpdir) if f.lower().endswith(".csv")]
        if not files:
            raise RuntimeError("Conversione con soffice non ha prodotto CSV")
        csv_path = os.path.join(tmpdir, files[0])

        # leggiamo un campione per sniffing e decidere delimiter
        with open(csv_path, "r", encoding="latin-1", errors="replace") as sf:
            sample = sf.read(8192)

        # rileviamo il delimitatore
        detected = _detect_delimiter_and_dialect(sample)
        # proviamo prima con quello rilevato, poi con ; poi con ,
        tried = []
        for delim in [detected, ";", ",", "\t", "|"]:
            if delim in tried:
                continue
            tried.append(delim)
            try:
                return import_clients_from_csv_generic(csv_path, delimiter=delim, encoding="latin-1")
            except Exception:
                # proviamo il prossimo delimiter
                continue
        # se nessuno funziona, solleviamo errore
        raise RuntimeError("Impossibile importare il CSV prodotto da soffice: controlla il file convertito")
    finally:
        # non cancelliamo temporanea per debug; puoi rimuovere la tmp se vuoi
        pass

# -------------------------
# parsing date 
# -------------------------
def _normalize_date_for_db(raw) -> str | None:
    """
    Prova a parsare raw (string/date/datetime) e ritorna ISO 'YYYY-MM-DD' oppure None.
    Usa _parse_date se disponibile.
    """
    if raw is None:
        return None
    # se è già date/datetime
    if isinstance(raw, date) and not isinstance(raw, datetime):
        return raw.isoformat()
    if isinstance(raw, datetime):
        return raw.date().isoformat()

    # se esiste _parse_date usalo
    try:
        dt = _parse_date(raw)
        if dt:
            return dt.isoformat()
    except Exception:
        pass

    # fallback: prova alcuni formati comuni
    s = str(raw).strip()
    if s == "":
        return None
    # Ordine di preferenza: DD/MM/YYYY prima, poi altri formati
    for fmt in ("%d/%m/%Y", "%d-%m-%Y", "%d.%m.%Y", "%Y-%m-%d", "%Y/%m/%d", "%Y%m%d"):
        try:
            return datetime.strptime(s.split()[0], fmt).date().isoformat()
        except Exception:
            pass
    return None


def save_all_cliente_fields(rowid: int, data: Dict[str, Any],
                            run_post_update: bool = True) -> Dict[str, Any]:
    """
    Aggiorna le colonne indicate in `data` per il cliente con id=rowid.
    Restituisce dict {ok, affected, error}.
    """
    if not isinstance(data, dict) or not rowid:
        return {"ok": False, "affected": 0, "error": "Parametri invalidi"}

    try:
        cols = list(data.keys())
        if cols:
            _ensure_clienti_columns(cols)

        conn = get_connection()
        cur = conn.cursor()

        set_parts = []
        params = []
        for col, val in data.items():
            set_parts.append(f'"{col}" = ?')
            
            if col.upper() in DATE_COLS:
                norm = _normalize_date_for_db(val)
                params.append(norm)                  
            elif val is None:
                params.append(None)
            elif isinstance(val, (date, datetime)):
                params.append(val.date().isoformat() if isinstance(val, datetime) else val.isoformat())
            else:
                sval = str(val).strip()
                params.append(None if sval == "" or sval.upper() == "NONE" else sval)

        if not set_parts:
            return {"ok": False, "affected": 0, "error": "Nessuna colonna da aggiornare"}

        params.append(rowid)
        sql = f'UPDATE clienti SET {", ".join(set_parts)} WHERE id = ?'
        cur.execute(sql, params)
        conn.commit()
        affected = conn.total_changes
        conn.close()

        if run_post_update:
            try:
                aggiorna_finemf_e_garanzia([rowid])
            except Exception:
                # non bloccare il salvataggio principale per errori nel post-update
                print("[save_all_cliente_fields] post-update error:", traceback.format_exc())

        return {"ok": True, "affected": affected, "error": None}

    except Exception as e:
        err = traceback.format_exc()
        print("[save_all_cliente_fields] ERRORE:", err)
        try:
            conn.rollback()
            conn.close()
        except Exception:
            pass
        return {"ok": False, "affected": 0, "error": str(e)}



# ---------- Helpers per compilazione PDF & barcode ----------

def _simple_wrap(text: str, max_width: float, font_name="Helvetica", font_size=10, c=None):
    if not text:
        return [""]
    if c is None:
        # fallback naive
        words = str(text).split()
        lines, cur = [], ""
        for w in words:
            if len(cur) + len(w) + 1 <= 40:
                cur = (cur + " " + w).strip()
            else:
                if cur: lines.append(cur)
                cur = w
        if cur: lines.append(cur)
        return lines
    lines=[]
    words = str(text).split()
    cur = ""
    for w in words:
        trial = (cur + " " + w).strip()
        if c.stringWidth(trial, font_name, font_size) <= max_width:
            cur = trial
        else:
            if cur: lines.append(cur)
            cur = w
    if cur: lines.append(cur)
    return lines

def generate_barcode_image(matricola: str, out_path: str, logo: str = "") -> str:
    """Genera PNG Code128 e ritorna il percorso del file .png
    Se logo è fornito, il barcode conterrà logo+matricola concatenati"""
    matricola = str(matricola or "").strip()
    if matricola == "":
        raise ValueError("Matricola vuota per barcode")

    if Code128 is None:
        raise RuntimeError(
            "Modulo barcode mancante. Su Windows installa: pip install python-barcode pillow"
        )
    
    # Se logo fornito, concatena con matricola
    barcode_content = f"{logo}{matricola}".strip() if logo else matricola
    
    base, ext = os.path.splitext(out_path)
    if not ext.lower().endswith(".png"):
        out_path = base + ".png"
    writer = ImageWriter()
    code = Code128(barcode_content, writer=writer)
    saved = code.save(base)  # ritorna base + '.png'
    return saved

def _fill_pdf_overlay(template_size: tuple, context: dict, barcode_path: str, field_positions: dict) -> BytesIO:
    """Riempie un overlay PDF con testo e barcode."""
    width_pt, height_pt = template_size
    buf = BytesIO()
    c = canvas.Canvas(buf, pagesize=(width_pt, height_pt))
    
    BARCODE_PADDING = mm_to_pt(1.5)
    
    for key, pos in field_positions.items():
        if key in ("BARCODE", "LOGOTIPO", "MATRICOLA"):
            continue
        
        # Supporta sia 4-tuple (x, y, font, max_w) sia 5-tuple (x, y, box_w, box_h, font)
        if len(pos) == 4:
            x, y, fsize, maxw = pos
            box_w = maxw
        elif len(pos) == 5:
            x, y, box_w, box_h, fsize = pos
        else:
            continue
        
        val = str(context.get(key, "") or "")
        lines = _simple_wrap(val, box_w, font_name="Helvetica", font_size=fsize, c=c)
        
        textobj = c.beginText()
        textobj.setFont("Helvetica", fsize)
        textobj.setLeading(fsize + 2)
        textobj.setTextOrigin(x, y)
        
        for ln in lines:
            textobj.textLine(ln)
        
        c.drawText(textobj)
    
    # Disegna LOGOTIPO e MATRICOLA insieme, concatenati, senza bold
    if "MATRICOLA" in field_positions:
        matric_pos = field_positions["MATRICOLA"]
        
        # Formato 4-tuple (x, y, font_size, max_width)
        if len(matric_pos) == 4:
            x, y, font_size, _ = matric_pos
            
            logotipo_val = str(context.get("LOGOTIPO", "") or "").strip()
            matricola_val = str(context.get("MATRICOLA", "") or "").strip()
            
            # Concatena logotipo e matricola
            display_text = f"{logotipo_val} {matricola_val}".strip() if logotipo_val else matricola_val
            
            if display_text:
                c.setFont("Helvetica", font_size)
                c.drawString(x, y, display_text)
    
    # Aggiungi BARCODE se presente (contiene gia logotipo + matricola)
    if "BARCODE" in field_positions:
        pos = field_positions["BARCODE"]
        if len(pos) == 4:
            x, y, box_w, box_h = pos
            if barcode_path and os.path.exists(barcode_path):
                try:
                    im = Image.open(barcode_path)
                    img_w, img_h = im.size
                    aspect = img_w / img_h if img_h != 0 else 1.0

                    avail_w = box_w - 2 * BARCODE_PADDING
                    avail_h = box_h - 2 * BARCODE_PADDING

                    if avail_w <= 0 or avail_h <= 0:
                        draw_w, draw_h = box_w, box_h
                        x_img, y_img = x, y
                    else:
                        if avail_w / avail_h > aspect:
                            draw_h = avail_h
                            draw_w = draw_h * aspect
                        else:
                            draw_w = avail_w
                            draw_h = draw_w / aspect
                        x_img = x + (box_w - draw_w) / 2.0
                        y_img = y + (box_h - draw_h) / 2.0

                    # Disegna il barcode (che contiene logotipo + matricola)
                    c.drawImage(barcode_path, x_img, y_img, width=draw_w, height=draw_h,
                               preserveAspectRatio=True, mask='auto')
                except Exception:
                    pass
            else:
                # Fallback QR se immagine barcode non disponibile (compatibilita Windows)
                try:
                    from reportlab.graphics.barcode import qr
                    from reportlab.graphics.shapes import Drawing
                    from reportlab.graphics import renderPDF

                    raw_value = str(context.get("BARCODE_VALUE") or context.get("MATRICOLA") or "").strip()
                    if raw_value:
                        qr_code = qr.QrCodeWidget(raw_value)
                        bounds = qr_code.getBounds()
                        qr_w = bounds[2] - bounds[0]
                        qr_h = bounds[3] - bounds[1]
                        if qr_w > 0 and qr_h > 0:
                            draw = Drawing(qr_w, qr_h)
                            draw.add(qr_code)
                            avail_w = max(1.0, box_w - 2 * BARCODE_PADDING)
                            avail_h = max(1.0, box_h - 2 * BARCODE_PADDING)
                            scale = min(avail_w / qr_w, avail_h / qr_h)
                            draw.scale(scale, scale)
                            x_img = x + (box_w - (qr_w * scale)) / 2.0
                            y_img = y + (box_h - (qr_h * scale)) / 2.0
                            renderPDF.draw(draw, c, x_img, y_img)
                except Exception:
                    pass
    
    c.showPage()
    c.save()
    buf.seek(0)
    return buf

def _merge_overlay_with_pdf(template_pdf_path: str, overlay_buf: BytesIO, out_pdf_path: str):
    """Mergia overlay PDF con template."""
    reader = PdfReader(template_pdf_path)
    overlay_reader = PdfReader(overlay_buf)
    writer = PdfWriter()
    for i, page in enumerate(reader.pages):
        base_page = page
        if i == 0:
            try:
                base_page.merge_page(overlay_reader.pages[0])
            except Exception:
                pass
        writer.add_page(base_page)
    with open(out_pdf_path, "wb") as f:
        writer.write(f)

def get_clients_for_vf_month_year(month: int, year: int):
    """Recupera clienti per mese/anno basato su VF."""
    conn = get_connection()
    cur = conn.cursor()
    cur.execute("PRAGMA table_info(clienti)")
    cols = [r[1] for r in cur.fetchall()]
    colnames_sql = ", ".join([f'"{c}"' for c in cols])
    cur.execute(f"SELECT {colnames_sql} FROM clienti")
    rows = cur.fetchall()
    out = []
    
    for row in rows:
        rec = {cols[i]: (row[i] if i < len(row) else "") for i in range(len(cols))}
        vf_val = rec.get("VF", "") or ""
        matched = False
        
        try:
            d = _parse_date(vf_val)  # Parse data VF
            if d and d.year == int(year) and d.month == int(month):
                matched = True
        except Exception:
            pass
        
        if matched:
            try:
                matricola = rec.get("MATRICOLA", "").strip()
                if matricola:
                    cur.execute("SELECT MARCA FROM macchine WHERE trim(matricola) = ? COLLATE NOCASE LIMIT 1",
                               (matricola,))
                    m_row = cur.fetchone()
                    if m_row and m_row[0]:
                        rec["MARCA"] = m_row[0]
            except Exception:
                pass
            out.append(rec)
    
    conn.close()
    return out

def get_clients_for_vf_date(vf_date: str):
    """
    Recupera clienti che hanno il campo VF uguale alla data passata.
    vf_date deve essere in formato DD/MM/YYYY, DD-MM-YYYY, o YYYY-MM-DD.
    Ritorna lista di dict con i dati dei clienti.
    """
    try:
        # Parsa la data di input
        target_date = _parse_date(vf_date)
        if not target_date:
            return []
    except Exception:
        return []
    
    conn = get_connection()
    cur = conn.cursor()
    cur.execute("PRAGMA table_info(clienti)")
    cols = [r[1] for r in cur.fetchall()]
    colnames_sql = ", ".join([f'"{c}"' for c in cols])
    cur.execute(f"SELECT {colnames_sql} FROM clienti")
    rows = cur.fetchall()
    out = []
    
    for row in rows:
        rec = {cols[i]: (row[i] if i < len(row) else "") for i in range(len(cols))}
        vf_val = rec.get("VF", "") or ""
        
        try:
            d = _parse_date(vf_val)
            if d and d == target_date:
                out.append(rec)
        except Exception:
            pass
    
    conn.close()
    return out

def fill_pdf_with_client_data(vf_date: str, client: dict, template_pdf_path: str, output_pdf_path: str):
    """
    Compila un singolo PDF di checklist con i dati di un cliente specifico.
    vf_date: data nel formato DD/MM/YYYY (usata per il report)
    client: dict con i dati del cliente
    template_pdf_path: percorso del template PDF
    output_pdf_path: percorso di output per il PDF compilato
    """
    if not os.path.exists(template_pdf_path):
        raise FileNotFoundError(f"Template PDF non trovato: {template_pdf_path}")
    
    matric = client.get("MATRICOLA", "") or ""
    safe_mat = re.sub(r'[^0-9A-Za-z_-]', '_', str(matric))
    temp_dir = os.path.dirname(output_pdf_path)
    os.makedirs(temp_dir, exist_ok=True)
    barcode_path = os.path.join(temp_dir, f"barcode_{safe_mat}.png")
    
    # Prepara logotipo + matricola per il barcode
    logotipo = client.get("LOGOTIPO", "").strip() or ""
    try:
        generate_barcode_image(matric, barcode_path, logo=logotipo)
    except Exception:
        barcode_path = None
    
    ctx = {
        "RAG_SOC": client.get("RAG_SOC", ""),
        "P_IVA": client.get("P_IVA", "") or client.get("PIVA", ""),
        "INDIRIZZO": client.get("IND_VIA", "") or client.get("INDIRIZZO", ""),
        "CAP_LOCALITA": f"{client.get('CAP', '')} {client.get('COMUNE', '')}".strip(),
        "MODELLO": client.get("MODELLO", ""),
        "MARCA": client.get("MARCA", ""),
        "MF": client.get("MF", ""),
        "MATRICOLA": matric,
        "LOGOTIPO": logotipo,
        "BARCODE_VALUE": f"{logotipo}{matric}".strip() if logotipo else matric,
        "SCAD_ASS_LBL": "SCAD.ASS:",
        "SCAD_ASS": client.get("SCAD_ASS", ""),
        "INST_LBL": "INST.:",
        "DATAINST": client.get("DATAINST", ""),
        "GAR_LBL": "GAR.:",
        "GARANZIA": client.get("GARANZIA", ""),
    }
    
    reader = PdfReader(template_pdf_path)
    first_page = reader.pages[0]
    media = first_page.mediabox
    width_pt = float(media.width)
    height_pt = float(media.height)
    
    field_positions = FIELD_POSITIONS_DEFAULT
    overlay = _fill_pdf_overlay((width_pt, height_pt), ctx, barcode_path, field_positions)
    _merge_overlay_with_pdf(template_pdf_path, overlay, output_pdf_path)

def generate_filled_pdfs_for_month(year: int, month: int, template_pdf_path: str, 
                                   output_dir: str, field_positions: dict = None, copies: int = 2, auto_print: bool = True):
    """Genera PDF compilati per il mese/anno. Se auto_print=True, li stampa automaticamente."""
    os.makedirs(output_dir, exist_ok=True)
    field_positions = field_positions or FIELD_POSITIONS_DEFAULT
    clients = get_clients_for_vf_month_year(month, year)
    produced = []
    
    for c in clients:
        matric = c.get("MATRICOLA", "") or ""
        safe_mat = re.sub(r'[^0-9A-Za-z_-]', '_', str(matric))
        barcode_path = os.path.join(output_dir, f"barcode_{safe_mat}.png")
        
        # Prepara logotipo + matricola per il barcode
        logotipo = c.get("LOGOTIPO", "").strip() or ""
        try:
            generate_barcode_image(matric, barcode_path, logo=logotipo)
        except Exception:
            barcode_path = None
        
        ctx = {
            "RAG_SOC": c.get("RAG_SOC", ""),
            "P_IVA": c.get("P_IVA", "") or c.get("PIVA", ""),
            "INDIRIZZO": c.get("IND_VIA", "") or c.get("INDIRIZZO", ""),
            "CAP_LOCALITA": f"{c.get('CAP', '')} {c.get('COMUNE', '')}".strip(),
            "MODELLO": c.get("MODELLO", ""),
            "MARCA": c.get("MARCA", ""),
            "MF": c.get("MF", ""),
            "MATRICOLA": f"{logotipo} {matric}".strip() if logotipo else matric,
            "LOGOTIPO": logotipo,
            "BARCODE_VALUE": f"{logotipo}{matric}".strip() if logotipo else matric,
            "SCAD_ASS_LBL": "SCAD.ASS:",
            "SCAD_ASS": c.get("SCAD_ASS", ""),
            "INST_LBL": "INST.:",
            "DATAINST": c.get("DATAINST", ""),
            "GAR_LBL": "GAR.:",
            "GARANZIA": c.get("GARANZIA", ""),
        }
        
        reader = PdfReader(template_pdf_path)
        first_page = reader.pages[0]
        media = first_page.mediabox
        width_pt = float(media.width)
        height_pt = float(media.height)
        
        overlay = _fill_pdf_overlay((width_pt, height_pt), ctx, barcode_path, field_positions)
        out_pdf = os.path.join(output_dir, f"check_{safe_mat}.pdf")
        _merge_overlay_with_pdf(template_pdf_path, overlay, out_pdf)
        
        # Stampa solo se auto_print è True
        if auto_print:
            try:
                print_pdf_file(out_pdf, copies=copies)
            except Exception as e:
                print(f"Avvertimento: impossibile stampare {out_pdf}: {e}")
        
        produced.append(out_pdf)
    
    return produced

def print_pdf_file(path: str, copies: int = 2):
    """Stampa un PDF con il numero di copie specificato (default: 2)."""
    if sys.platform.startswith("win"):
        try:
            for _ in range(copies):
                os.startfile(os.path.abspath(path), "print")
            return True
        except Exception as e:
            raise
    else:
        try:
            for _ in range(copies):
                subprocess.run(["lp", os.path.abspath(path)], check=True)
            return True
        except Exception as e:
            raise


def generate_cliente_checklist_pdf(cliente_data: dict) -> str:
    """
    Genera un PDF di checklist per il cliente usando template.pdf come base.
    Mostra i dati dei campi del cliente nelle posizioni definite in FIELD_POSITIONS_CHECKLIST.
    Aggiunge un barcode generato da MATRICOLA + LOGOTIPO.
    
    Args:
        cliente_data: dict con i campi del cliente (RAG_SOC, P_IVA, MATRICOLA, ecc.)
    
    Returns:
        str: path del file PDF generato
    """
    try:
        from reportlab.lib.pagesizes import A4
        from reportlab.lib.units import mm
        from reportlab.pdfgen import canvas as rl_canvas
        from PyPDF2 import PdfReader, PdfWriter
        from barcode import generate as barcode_generate
        import tempfile
        import os
        from io import BytesIO
        
        # Crea file temporaneo
        temp_dir = tempfile.mkdtemp(prefix="checklist_")
        pdf_path = os.path.join(temp_dir, "checklist_cliente.pdf")
        
        # Percorso al template
        app_dir = os.path.dirname(os.path.abspath(__file__))
        template_path = os.path.join(app_dir, "template.pdf")
        
        # Genero il dato per il codice a barre: MATRICOLA_LOGOTIPO
        matricola = cliente_data.get("MATRICOLA", "").strip()
        logotipo = cliente_data.get("LOGOTIPO", "").strip()
        barcode_data = f"{matricola}_{logotipo}" if matricola and logotipo else (matricola or "NOMAT")
        
        # Limita a 12 caratteri per code128
        barcode_data = barcode_data[:12]
        
        # Creo codice a barre (CODE128)
        barcode_path = None
        try:
            barcode_img = barcode_generate('code128', barcode_data, writer_options={'module_width': 0.5, 'module_height': 8})
            barcode_path = os.path.join(temp_dir, "barcode.png")
            barcode_img.save(barcode_path)
        except Exception as e:
            print(f"[WARN] Errore nella generazione barcode: {e}")
            barcode_path = None
        
        # Crea overlay con reportlab
        overlay_path = os.path.join(temp_dir, "overlay.pdf")
        c = rl_canvas.Canvas(overlay_path, pagesize=A4)
        width, height = A4
        
        # Disegna il barcode se è stato generato
        if barcode_path and "BARCODE" in FIELD_POSITIONS_CHECKLIST:
            try:
                barcode_pos = FIELD_POSITIONS_CHECKLIST["BARCODE"]
                barcode_x, barcode_y, barcode_width, barcode_height = barcode_pos
                c.drawImage(barcode_path, barcode_x, barcode_y, 
                           width=barcode_width, height=barcode_height)
            except Exception as e:
                print(f"[WARN] Errore nel disegno barcode: {e}")
        
        # Disegna i dati nella checklist usando FIELD_POSITIONS_CHECKLIST (escluso BARCODE)
        # Formato posizioni: "FIELD_NAME": (x_mm, y_mm, font_size, width_mm)
        # Campi terminanti con "_LBL" sono etichette fisse (cercano il testo in CHECKLIST_LABELS)
        for field_name, pos_data in FIELD_POSITIONS_CHECKLIST.items():
            if field_name == "BARCODE":
                continue  # Barcode già disegnato sopra
            
            x, y, font_size, max_width = pos_data
            
            # Ottieni il valore: da CHECKLIST_LABELS se è un'etichetta (_LBL), altrimenti da cliente_data
            if field_name.endswith("_LBL"):
                value = CHECKLIST_LABELS.get(field_name, "").strip()
                # Le etichette sono sempre bold
                c.setFont("Helvetica-Bold", font_size)
            else:
                value = cliente_data.get(field_name, "").strip()
                c.setFont("Helvetica", font_size)
            
            if value:  # Mostra solo se il valore esiste
                # Tronca il testo se supera la larghezza massima (approssimativamente)
                max_chars = int(max_width / (font_size * 0.5))
                if len(value) > max_chars:
                    value = value[:max_chars - 3] + "..."
                
                c.drawString(x, y, value)
        
        # Note a piè di pagina
        c.setFont("Helvetica-Oblique", 7)
        c.drawString(20*mm, 15*mm, f"Generato il: {datetime.now().strftime('%d/%m/%Y %H:%M:%S')}")
        
        # Salva overlay
        c.save()
        
        # Se esiste template.pdf, uniscilo con l'overlay
        if os.path.exists(template_path):
            try:
                reader = PdfReader(template_path)
                overlay_reader = PdfReader(overlay_path)
                writer = PdfWriter()
                
                # Prendi la prima pagina del template
                template_page = reader.pages[0]
                
                # Prendi la prima pagina dell'overlay
                overlay_page = overlay_reader.pages[0]
                
                # Unisci
                template_page.merge_page(overlay_page)
                writer.add_page(template_page)
                
                # Scrivi il risultato
                with open(pdf_path, "wb") as out_pdf:
                    writer.write(out_pdf)
            except Exception as e:
                print(f"[WARN] Errore nell'unione template con overlay: {e}")
                # Fallback: usa solo l'overlay
                with open(overlay_path, "rb") as f:
                    with open(pdf_path, "wb") as out:
                        out.write(f.read())
        else:
            # Fallback: usa solo l'overlay se template non esiste
            with open(overlay_path, "rb") as f:
                with open(pdf_path, "wb") as out:
                    out.write(f.read())
        
        return pdf_path
    
    except Exception as e:
        print(f"[ERROR] Errore nell'generazione della checklist PDF: {e}")
        raise


def report_memorie_in_scadenza(within_days: int = 30, db_path: str = "database.db") -> List[Tuple[str,str,str,str,str]]:
    """
    Restituisce tuple (RAG_SOC, FINEMF, MF, TEL, CODICE) per le memorie la cui FINEMF
    è entro 'within_days' giorni da oggi.
    Funzione robusta: se la tabella 'clienti' non contiene la colonna CODICE,
    cerca alternative (MATRICOLA, CODUNIVOCO, CODPUNTOCA) e le aliasa come CODICE.
    """
    conn = sqlite3.connect(db_path)
    conn.row_factory = sqlite3.Row
    cur = conn.cursor()

    today = date.today()
    end_date = today + timedelta(days=within_days)

    out = []

    # --- determinazione colonna identificativa (alias -> CODICE) ---
    try:
        cur.execute("PRAGMA table_info(clienti)")
        cols = [r[1] for r in cur.fetchall()]
    except Exception:
        cols = []

    # lista candidate nell'ordine di preferenza
    candidates = ("CODICE", "MATRICOLA", "CODUNIVOCO", "CODPUNTOCA", "P_IVA", "ID")
    id_col = None
    for c in candidates:
        if c in cols:
            id_col = c
            break

    if id_col:
        select_id = f'{id_col} AS CODICE'
    else:
        select_id = "'' AS CODICE"

    # 1) prova la via veloce in SQL sulla tabella 'clienti' (date() su FINEMF)
    try:
        sql = f"""
            SELECT RAG_SOC, FINEMF, MF, TEL, {select_id}
            FROM clienti
            WHERE FINEMF IS NOT NULL AND TRIM(FINEMF) != '' AND date(FINEMF) BETWEEN ? AND ?
            ORDER BY FINEMF
        """
        cur.execute(sql, (today.isoformat(), end_date.isoformat()))
        rows = cur.fetchall()
        if rows:
            for r in rows:
                out.append((r["RAG_SOC"] or "", r["FINEMF"] or "", r["MF"] or "", r["TEL"] or "", r["CODICE"] or ""))
            conn.close()
            print(f"[report_memorie_in_scadenza] SQL fast path, trovati {len(out)}")
            return out
    except Exception as e:
        # fallisce se FINEMF non è in formato riconosciuto da date() o altre ragioni
        print("[report_memorie_in_scadenza] SQL fast path fallito:", e)

    # 2) fallback: estrai tutte le righe con FINEMF non null e parsale in Python
    try:
        cur.execute(f"SELECT * FROM clienti WHERE FINEMF IS NOT NULL AND TRIM(FINEMF) != ''")
        rows = cur.fetchall()
    except Exception as e:
        print("[report_memorie_in_scadenza] Errore lettura fallback clienti:", e)
        conn.close()
        return []

    unparsed = []
    for row in rows:
        # leggi i campi con difesa
        rag = row["RAG_SOC"] if "RAG_SOC" in row.keys() else row.get("RAG_SOC", "")
        finemf_raw = row["FINEMF"] if "FINEMF" in row.keys() else row.get("FINEMF", "")
        mf = row["MF"] if "MF" in row.keys() else row.get("MF", "")
        tel = row["TEL"] if "TEL" in row.keys() else row.get("TEL", "")
        # ottieni codice anche in fallback: proviamo a leggere la colonna scelta prima o alcune alternative
        codice = ""
        if "CODICE" in row.keys():
            codice = row["CODICE"] or ""
        else:
            for c in ("MATRICOLA", "CODUNIVOCO", "CODPUNTOCA", "P_IVA", "id"):
                if c in row.keys():
                    codice = row[c] or ""
                    break

        # parsing data FINEMF (più robuste heuristics)
        parsed_date = None
        if finemf_raw:
            try:
                parsed_date = _parse_date(finemf_raw)
            except Exception:
                parsed_date = None

        # altre heuristics se _parse_date non esiste o ha fallito
        if not parsed_date:
            try:
                s = str(finemf_raw).strip().split()[0]
                parsed_date = datetime.fromisoformat(s).date()
            except Exception:
                parsed_date = None
            if not parsed_date:
                # prova formati comuni
                # Ordine di preferenza: DD/MM/YYYY prima
                for fmt in ("%d/%m/%Y", "%d-%m-%Y", "%d.%m.%Y", "%Y-%m-%d", "%Y/%m/%d", "%Y%m%d"):
                    try:
                        parsed_date = datetime.strptime(s, fmt).date()
                        break
                    except Exception:
                        parsed_date = None

        if not parsed_date:
            if len(unparsed) < 20:
                unparsed.append((rag, finemf_raw))
            continue

        # Assicura che parsed_date sia un date object (non datetime)
        if isinstance(parsed_date, datetime):
            parsed_date = parsed_date.date()

        if today <= parsed_date <= end_date:
            display_date = parsed_date.strftime("%d/%m/%Y")
            out.append((rag or "", display_date, mf or "", tel or "", codice or ""))

    if unparsed:
        print("[report_memorie_in_scadenza] alcune FINEMF non parsate (esempi):", unparsed[:10])

    conn.close()
    print(f"[report_memorie_in_scadenza] fallback, trovati {len(out)}")
    return out

def get_datavp_from_matricola(matricola: str):
    """
    Recupera la DATA_INTERVENTO più recente dalla tabella INTAE per una data matricola.
    La matricola potrebbe avere i primi 2 caratteri come logotipo.
    Ritorna la data in formato ISO (YYYY-MM-DD) o stringa vuota se non trovata.
    """
    if not matricola or str(matricola).strip() == "":
        return ""
    
    try:
        conn = get_connection()
        cur = conn.cursor()
        
        # Normalizza matricola
        mat_norm = str(matricola).strip().lower()
        
        # Prova match esatto
        cur.execute(
            "SELECT DATA_INTERVENTO FROM intae WHERE lower(MATRICOLA) = ? ORDER BY DATA_INTERVENTO DESC LIMIT 1",
            (mat_norm,)
        )
        row = cur.fetchone()
        
        # Se nessun match esatto, prova LIKE
        if not row:
            cur.execute(
                "SELECT DATA_INTERVENTO FROM intae WHERE lower(MATRICOLA) LIKE ? ORDER BY DATA_INTERVENTO DESC LIMIT 1",
                (f"%{mat_norm}%",)
            )
            row = cur.fetchone()
        
        conn.close()
        
        if row and row[0]:
            return str(row[0]).strip()
        return ""
    except Exception as e:
        print(f"[ERROR] get_datavp_from_matricola({matricola}): {e}")
        return ""


def report_rt_for_clienti(month: int | None = None, year: int | None = None):
    """
    - Legge clienti dalla tabella e filtra per VF nel mese/anno specificato
    - Se month/year sono None, restituisce tutti i clienti
    - Ritorna tuple: (RAG_SOC, IND_VIA, COMUNE, MATRICOLA, DATAVP, DATA, TIPOMF, DATAINST, AZZ, P_IVA, MARCA)
    """
    conn = get_connection()
    cur = conn.cursor()

    # 1) leggi una riga per ottenere i nomi delle colonne
    try:
        cur.execute("SELECT * FROM clienti LIMIT 1")
        cols = [d[0] for d in cur.description] if cur.description else []
    except Exception:
        cols = []

    # Estrai colonne VF (FINEMF), PIVA, ID e MATRICOLA per il filtraggio
    finemf_col = None
    piva_col = None
    id_col = None
    matricola_col = None

    if cols:
        lowered = [c.lower() for c in cols]
        # trova VF con priorità
        for target in ('vf',):
            for i, name in enumerate(lowered):
                if name == target:
                    finemf_col = cols[i]; break
            if finemf_col: break

        if not finemf_col:
            for key in ('vf', 'fine', 'scad', 'scaden', 'datafine', 'datainizio', 'data'):
                for i, name in enumerate(lowered):
                    if key in name:
                        finemf_col = cols[i]; break
                if finemf_col: break

        # trova PIVA con heuristics
        for i, name in enumerate(lowered):
            if 'p_iva' in name or 'piva' == name or 'partita' in name or 'iva' in name:
                piva_col = cols[i]; break

        # id e matricola (per join con macchine)
        for i, name in enumerate(lowered):
            if name == 'id':
                id_col = cols[i]
            if 'matricola' in name:
                matricola_col = cols[i]

    if not finemf_col:
        # Fallback: se colonna VF non trovata, restituisci configurazione di default
        sel = """
        SELECT
            COALESCE(RAG_SOC, ''),
            COALESCE(IND_VIA, ''),
            COALESCE(COMUNE, ''),
            COALESCE(MATRICOLA, ''),
            COALESCE(DATA, ''),
            COALESCE(TIPOMF, ''),
            COALESCE(DATAINST, ''),
            COALESCE(AZZ, '')
        FROM clienti
        ORDER BY UPPER(COALESCE(COMUNE, '')), UPPER(COALESCE(RAG_SOC, ''))
        """
        try:
            cur.execute(sel)
            all_rows = cur.fetchall()
            # Se month/year specificati, filtra anche in fallback
            if month is not None and year is not None:
                matches = []
                for row in all_rows:
                    data_raw = row[4] if len(row) > 4 else ''  # colonna DATA (5a colonna)
                    data_raw = '' if data_raw is None else str(data_raw).strip()
                    if data_raw:
                        parsed = _parse_date(data_raw)
                        if parsed and parsed.year == int(year) and parsed.month == int(month):
                            visible = tuple('' if v is None else str(v) for v in row[:8])
                            # Aggiungi DATAVP
                            matricola = visible[3] if len(visible) > 3 else ""
                            datavp = get_datavp_from_matricola(matricola)
                            # Inserisci DATAVP dopo MATRICOLA (indice 4)
                            visible = visible[:4] + (datavp,) + visible[4:]
                            matches.append(visible)
                conn.close()
                return matches
            else:
                conn.close()
                result = []
                for row in all_rows:
                    visible = tuple('' if v is None else str(v) for v in row)
                    # Aggiungi DATAVP
                    matricola = visible[3] if len(visible) > 3 else ""
                    datavp = get_datavp_from_matricola(matricola)
                    # Inserisci DATAVP dopo MATRICOLA (indice 4)
                    visible = visible[:4] + (datavp,) + visible[4:]
                    result.append(visible)
                return result
        except Exception:
            conn.close()
            print("[ERROR][DB] Errore fallback SELECT:", traceback.format_exc())
            return []

    # Costruisci SELECT per recuperare e filtrare i dati
    piva_expr = f"COALESCE({piva_col}, '')" if piva_col else "'' as piva_missing"
    sql = f"""
    SELECT
        COALESCE(RAG_SOC, ''),
        COALESCE(IND_VIA, ''),
        COALESCE(COMUNE, ''),
        COALESCE(MATRICOLA, ''),
        COALESCE(DATA, ''),
        COALESCE(TIPOMF, ''),
        COALESCE(DATAINST, ''),
        COALESCE(AZZ, ''),
        COALESCE({finemf_col}, '') as vf_col,
        {piva_expr} as piva_col
    FROM clienti
    ORDER BY UPPER(COALESCE(COMUNE, '')), UPPER(COALESCE(RAG_SOC, ''))
    """
    try:
        cur.execute(sql)
        all_rows = cur.fetchall()
    except Exception:
        print("[ERROR][DB] Errore SELECT con vf_col:", traceback.format_exc())
        conn.close()
        return []

    # Itera e filtra i risultati
    matches = []
    sample_shown = 0
    SAMPLE_LIMIT = 50
    parsed_failures = 0
    parsed_success = 0

    for idx, row in enumerate(all_rows):
        # row layout: 0..7 visible, 8=vf_col, 9=piva_col
        visible = tuple('' if v is None else str(v) for v in row[:8])
        finemf_raw = row[8] if len(row) > 8 else ''
        piva_val = row[9] if len(row) > 9 else ''
        finemf_raw = '' if finemf_raw is None else str(finemf_raw).strip()
        piva_val = '' if piva_val is None else str(piva_val).strip()

        if month is None and year is None:
            # se NON sono stati specificati filtri, aggiungi tutti i record
            marca_val = ''
            # costruiamo un dict minimale per cercare macchina
            client_row_dict = {}
            # prova a popolare matricola/id dal row se possibile
            # riga: visible[3] è MATRICOLA nella SELECT
            if visible and len(visible) > 3:
                client_row_dict['MATRICOLA'] = visible[3]
            # se avevi trovato id_col prima, possiamo provare a recuperare id con una query addizionale:
            if id_col:
                # proviamo a trovare l'id usando RAG_SOC+MATRICOLA come chiave (più robusto sarebbe includere id direttamente nella SELECT)
                try:
                    # cerca cliente completo con ragione sociale e matricola
                    rag = visible[0] if len(visible) > 0 else None
                    mat = visible[3] if len(visible) > 3 else None
                    if rag or mat:
                        conds = []
                        params = []
                        if rag:
                            conds.append("UPPER(COALESCE(RAG_SOC,'')) = UPPER(?)")
                            params.append(rag)
                        if mat:
                            conds.append("lower(trim(MATRICOLA)) = lower(trim(?))")
                            params.append(mat)
                            where = " AND ".join(conds)
                        cur.execute(f"SELECT * FROM clienti WHERE {where} LIMIT 1", tuple(params))
                        crow = cur.fetchone()
                        if crow:
                            info = cur.description
                            client_row_dict = {info[i][0]: crow[i] for i in range(len(crow))}
                except Exception:
                    pass

            # recupera marca tramite helper
            try:
                machine = get_primary_machine_for_client(cur, client_row_dict)
                if machine:
                    marca_val = machine.get('MARCA') or machine.get('marca') or machine.get('MODELLO') or ''
            except Exception:
                marca_val = ''
            
            # Aggiungi DATAVP
            matricola = visible[3] if len(visible) > 3 else ""
            datavp = get_datavp_from_matricola(matricola)
            # Inserisci DATAVP dopo MATRICOLA (indice 4)
            result_tuple = visible[:4] + (datavp,) + visible[4:] + (piva_val, marca_val)
            
            matches.append(result_tuple)
            continue

        if not finemf_raw:
            parsed_failures += 1
            continue

        parsed = _parse_date(finemf_raw)
        if parsed:
            parsed_success += 1
            try:
                # Converti month/year a int (potrebbero essere string da GUI)
                try:
                    month_int = int(month) if month else None
                    year_int = int(year) if year else None
                except (ValueError, TypeError):
                    month_int = None
                    year_int = None
                
                if month_int and year_int and parsed.year == year_int and parsed.month == month_int:
                    # recupera marca come sopra
                    marca_val = ''
                    client_row_dict = {}
                    if visible and len(visible) > 3:
                        client_row_dict['MATRICOLA'] = visible[3]
                    # try find full client row (optional)
                    try:
                        rag = visible[0] if len(visible) > 0 else None
                        mat = visible[3] if len(visible) > 3 else None
                        if rag or mat:
                            conds = []
                            params = []
                            if rag:
                                conds.append("UPPER(COALESCE(RAG_SOC,'')) = UPPER(?)")
                                params.append(rag)
                            if mat:
                                conds.append("lower(trim(MATRICOLA)) = lower(trim(?))")
                                params.append(mat)
                            where = " AND ".join(conds)
                            cur.execute(f"SELECT * FROM clienti WHERE {where} LIMIT 1", tuple(params))
                            crow = cur.fetchone()
                            if crow:
                                info = cur.description
                                client_row_dict = {info[i][0]: crow[i] for i in range(len(crow))}
                    except Exception:
                        pass

                    try:
                        machine = get_primary_machine_for_client(cur, client_row_dict)
                        if machine:
                            marca_val = machine.get('MARCA') or machine.get('marca') or machine.get('MODELLO') or ''
                    except Exception:
                        marca_val = ''

                    # Aggiungi DATAVP
                    matricola = visible[3] if len(visible) > 3 else ""
                    datavp = get_datavp_from_matricola(matricola)
                    # Inserisci DATAVP dopo MATRICOLA (indice 4)
                    result_tuple = visible[:4] + (datavp,) + visible[4:] + (piva_val, marca_val)
                    
                    matches.append(result_tuple)
                else:
                    # non coincide: ignora
                    pass
            except Exception:
                pass
        else:
            parsed_failures += 1

    conn.close()
    return matches


def get_primary_machine_for_client(cur, client_row_dict):
    """
    Restituisce dict macchina (o None) per il cliente. Cerca prima macchine con campo cliente_id,
    poi con matching matricola. Usa il cursor passato.
    """
    # cerca id cliente
    possible_ids = []
    for k in ("id", "ID", "Id"):
        if client_row_dict.get(k) is not None and str(client_row_dict.get(k)).strip() != "":
            possible_ids.append(client_row_dict.get(k))

    if possible_ids:
        for cid in possible_ids:
            try:
                cur.execute("SELECT * FROM macchine WHERE cliente_id = ? LIMIT 1", (cid,))
                row = cur.fetchone()
                if row:
                    info = cur.execute("PRAGMA table_info(macchine)").fetchall()
                    cols = [r[1] for r in info]
                    return {cols[i]: (row[i] if i < len(row) else None) for i in range(len(cols))}
            except Exception:
                pass

    # fallback matricola
    matricola = (client_row_dict.get("MATRICOLA") or client_row_dict.get("matricola") or "").strip()
    if matricola:
        try:
            cur.execute("SELECT * FROM macchine WHERE lower(trim(matricola)) = ? LIMIT 1", (matricola.lower(),))
            row = cur.fetchone()
            if row:
                info = cur.execute("PRAGMA table_info(macchine)").fetchall()
                cols = [r[1] for r in info]
                return {cols[i]: (row[i] if i < len(row) else None) for i in range(len(cols))}
        except Exception:
            pass

    return None


def build_ctx_from_client_and_machine(client: dict, machine: dict = None) -> dict:
    """
    Costruisce un ctx con alias multipli per campi come P_IVA, MARCA, MODELLO, ecc.
    (utile se in futuro dovrai creare preview PDF direttamente da qui).
    """
    if client is None:
        client = {}

    rag = client.get("RAG_SOC") or client.get("rag_soc") or client.get("RAGIONE_SOCIALE") or client.get("ragione_sociale") or client.get("NOME") or ""
    piva = client.get("P_IVA") or client.get("PIVA") or client.get("p_iva") or client.get("partita_iva") or ""
    indirizzo = client.get("INDIRIZZO") or client.get("IND_VIA") or ""
    cap = client.get("CAP") or ""
    comune = client.get("COMUNE") or ""
    matricola = client.get("MATRICOLA") or client.get("matricola") or ""

    marca = ""
    modello = ""
    if machine:
        marca = machine.get("MARCA") or machine.get("marca") or machine.get("MODELLO") or machine.get("model") or ""
        modello = machine.get("MODELLO") or machine.get("model") or ""

    if not marca:
        marca = client.get("MARCA") or client.get("marca") or ""

    ctx = {
        "RAG_SOC": rag, "rag_soc": rag, "RAGIONE_SOCIALE": rag,
        "P_IVA": piva, "PIVA": piva, "partita_iva": piva,
        "INDIRIZZO": indirizzo, "IND_VIA": indirizzo,
        "CAP": cap, "COMUNE": comune, "CAP_LOCALITA": f"{cap} {comune}".strip(),
        "MATRICOLA": matricola, "matricola": matricola,
        "MARCA": marca, "marca": marca, "MODELLO": modello, "MODEL": modello,
        "FULL_CLIENT_DISPLAY": f"{rag} - {piva}" if (rag or piva) else rag
    }

    return ctx


def mm_to_pt(v_mm: float) -> float:
    """Convert mm to points (reportlab uses points)."""
    return v_mm * 72.0 / 25.4

def pt_to_mm(v_pt: float) -> float:
    """Convert points to mm."""
    return v_pt * 25.4 / 72.0

# Posizioni di default (adatta ai tuoi template). Coordinate in punti (origine bottom-left).
FIELD_POSITIONS_DEFAULT = {
    "RAG_SOC": (mm_to_pt(85), mm_to_pt(250), 10, mm_to_pt(180)),
    "P_IVA": (mm_to_pt(85), mm_to_pt(235), 10, mm_to_pt(180)),     
    "INDIRIZZO": (mm_to_pt(85), mm_to_pt(225), 10, mm_to_pt(180)),
    "CAP_LOCALITA": (mm_to_pt(85), mm_to_pt(216), 10, mm_to_pt(180)),     
    "MARCA": (mm_to_pt(85), mm_to_pt(210), 10, mm_to_pt(180)),       
    "MF": (mm_to_pt(85), mm_to_pt(205), 10, mm_to_pt(180)),          
    "LOGOTIPO": (mm_to_pt(85), mm_to_pt(200), 10, mm_to_pt(30)),  # Posizione originale Y=200mm
    "MATRICOLA": (mm_to_pt(85), mm_to_pt(200), 10, mm_to_pt(140)),  # Sotto alle altre scritte X=85mm, Y=200mm
    
    "SCAD_ASS_LBL": (mm_to_pt(21), mm_to_pt(186), 8, mm_to_pt(30)),     # Etichetta a sinistra
    "SCAD_ASS": (mm_to_pt(37), mm_to_pt(186), 8, mm_to_pt(50)),         # Dato a destra
    
    "INST_LBL": (mm_to_pt(57), mm_to_pt(186), 8, mm_to_pt(30)),         # Etichetta a sinistra
    "DATAINST": (mm_to_pt(65), mm_to_pt(186), 8, mm_to_pt(50)),         # Dato a destra
    
    "GAR_LBL": (mm_to_pt(85), mm_to_pt(186), 8, mm_to_pt(30)),         # Etichetta a sinistra
    "GARANZIA": (mm_to_pt(94), mm_to_pt(186), 8, mm_to_pt(50)),        # Dato a destra
    
    "BARCODE": (mm_to_pt(140), mm_to_pt(275), mm_to_pt(50), mm_to_pt(20)),
}

# Posizioni per la checklist dei clienti nella tab Clienti (solo valori, senza etichette)
# Formato: "FIELD_NAME": (x_mm, y_mm, font_size, width_mm)
# Campi terminanti con "_LBL" sono etichette (non cercano valore in cliente_data, ma usano il testo in CHECKLIST_LABELS)
# Adatta facilmente le posizioni e i campi come preferisci
FIELD_POSITIONS_CHECKLIST = {
    "RAG_SOC": (85*mm, 250*mm, 10, 160*mm),         # Ragione Sociale
    "IND_VIA": (85*mm, 226*mm, 10, 160*mm),         # Indirizzo
    "P_IVA": (85*mm, 235*mm, 10, 160*mm),           # P.IVA
    "COMUNE": (85*mm, 242*mm, 10, 160*mm),          # Comune
    "CAP": (85*mm, 217*mm, 10, 160*mm),             # CAP
    "MATRICOLA": (85*mm, 200*mm, 10, 160*mm),       # Matricola
    "MARCA": (85*mm, 210*mm, 10, 160*mm),           # Marca
    "MF": (85*mm, 207*mm, 10, 160*mm),              # Modello
    "VF_LBL": (22*mm, 187*mm, 10, 40*mm),            # Etichetta "Data Verificazione:" a sinistra
    "VF": (55*mm, 187*mm, 10, 50*mm),               # Data Verificazione a destra dell'etichetta
    "FINEMF_LBL": (80*mm, 187*mm, 10, 50*mm),       # Etichetta "Fine Memoria Fiscale:" al centro
    "FINEMF": (120*mm, 187*mm, 10, 50*mm),          # Fine Memoria Fiscale a destra dell'etichetta
    "BARCODE": (140*mm, 260*mm, 50*mm, 20*mm),      # Barcode (x, y, width, height) in alto a destra
}

# Testo fisso per le etichette nella checklist
CHECKLIST_LABELS = {
    "VF_LBL": "Data Verificazione:",
    "FINEMF_LBL": "Fine Memoria Fiscale:",
}

# Carica il template PDF solo se PyPDF2 è disponibile
height_pt = 500  # Default height se template non è disponibile
if PdfReader is not None and os.path.exists("template.pdf"):
    try:
        reader = PdfReader("template.pdf")
        page = reader.pages[0]
        height_pt = float(page.mediabox.height)
    except Exception as e:
        print(f"[WARN] Impossibile caricare template.pdf: {e} - usando default")
        height_pt = 500

# ========================
# PDF HELPER FUNCTIONS
# ========================

def _parse_date(s):
    """
    Helper robusto per parsare date da varie fonti.
    Accetta string, date, datetime. Ritorna datetime o None.
    """
    if s is None:
        return None
    if isinstance(s, datetime):
        return s
    if isinstance(s, date):
        return datetime.combine(s, datetime.min.time())
    s_str = str(s).strip()
    if not s_str:
        return None
    
    # rimuovi spazi Unicode problematici
    s_str = re.sub(r'[\u00A0\u200B]', ' ', s_str).strip()
    
    # prova parser.parse (dateutil) con diversi flag
    for kwargs in ({'dayfirst': False, 'yearfirst': True, 'fuzzy': True},
                   {'dayfirst': False, 'yearfirst': False, 'fuzzy': True},
                   {'dayfirst': True, 'yearfirst': False, 'fuzzy': True}):
        try:
            return parser.parse(s_str, **kwargs)
        except Exception:
            pass
    
    # prova regex: YYYY-MM-DD 0?M 0?D
    m = re.search(r'([12][0-9]{3})[-/\\. ]0?([1-9])[-/\\. ]0?([1-9])', s_str)
    if m:
        try:
            y = int(m.group(1)); mo = int(m.group(2)); da = int(m.group(3))
            return datetime(y, mo, da)
        except Exception:
            pass
    
    # prova regex: DD/MM/YYYY
    m = re.search(r'0?([1-9]|[12][0-9]|3[01])[-/\\. ]0?([1-9]|1[0-2])[-/\\. ]([12][0-9]{3})', s_str)
    if m:
        try:
            da = int(m.group(1)); mo = int(m.group(2)); y = int(m.group(3))
            return datetime(y, mo, da)
        except Exception:
            pass
    
    # prova regex: MM/YYYY
    m = re.search(r'0?([1-9]|1[0-2])[-/\\. ]([12][0-9]{3})', s_str)
    if m:
        try:
            mo = int(m.group(1)); y = int(m.group(2))
            return datetime(y, mo, 1)
        except Exception:
            pass
    
    # prova regex: YYYYMM
    m = re.search(r'^([12][0-9]{3})0?([1-9]|1[0-2])$', s_str)
    if m:
        try:
            y = int(m.group(1)); mo = int(m.group(2))
            return datetime(y, mo, 1)
        except Exception:
            pass
    
    # formati comuni finali
    # Ordine di preferenza: DD/MM/YYYY prima
    for fmt in ("%d/%m/%Y", "%d-%m-%Y", "%d.%m.%Y", "%Y-%m-%d", "%Y/%m/%d", "%Y%m%d", "%d%m%Y"):
        try:
            return datetime.strptime(s_str.split()[0], fmt)
        except Exception:
            pass
    
    return None


def _format_date_for_tree(date_str):
    """Formatta una data ISO o raw per visualizzazione in tree (DD/MM/YYYY)."""
    if not date_str or str(date_str).strip() == "":
        return ""
    try:
        d = _parse_date(date_str)
        if d:
            return d.strftime("%d/%m/%Y")
    except Exception:
        pass
    return str(date_str)

def aggiorna_finemf_e_garanzia(ids: Optional[List[int]] = None):
    """
    Aggiorna FINEMF e GARANZIA per i clienti con id in `ids` (o tutti se None).
    Se i campi sono già compilati, non li sovrascrive.
    FINEMF = DATAINST + 2500 giorni (fissi)
    GARANZIA = DATAINST + 1 anno (12 mesi)
    Ritorna dict { id: { 'finemf': nuovo_valore, 'garanzia': nuovo_valore }, ... }
    """
    conn = get_connection()
    cur = conn.cursor()
    
    # query: prendi righe da aggiornare
    if ids:
        placeholders = ",".join("?" for _ in ids)
        sql = f"SELECT id, DATAINST, FINEMF, GARANZIA FROM clienti WHERE id IN ({placeholders})"
        cur.execute(sql, ids)
    else:
        sql = "SELECT id, DATAINST, FINEMF, GARANZIA FROM clienti"
        cur.execute(sql)
    
    rows = cur.fetchall()
    updates = {}
    
    for row in rows:
        rid, datainst, finemf, garanzia = row
        to_update = {}
        
        # calcola FINEMF se vuoto (DATAINST + 2500 giorni)
        if not finemf or str(finemf).strip() in ("", "0"):
            try:
                nuovo_fin = calcola_fine_memoria(datainst)
                if nuovo_fin:
                    to_update["FINEMF"] = nuovo_fin
            except Exception:
                pass
        
        # calcola GARANZIA se vuota (DATAINST + 1 anno)
        if not garanzia or str(garanzia).strip() in ("", "0"):
            try:
                nuovo_gar = calcola_scadenza_garanzia(datainst, 12)  # 12 mesi = 1 anno
                if nuovo_gar:
                    to_update["GARANZIA"] = nuovo_gar
            except Exception:
                pass
        
        if to_update:
            # esegui UPDATE
            set_clause = ", ".join(f"{k}=?" for k in to_update.keys())
            vals = list(to_update.values()) + [rid]
            cur.execute(f"UPDATE clienti SET {set_clause} WHERE id=?", vals)
            updates[rid] = to_update
    
    conn.commit()
    conn.close()
    return updates

def calcola_fine_memoria(datainst, mxaz=None):
    """
    Calcola data fine memoria come DATAINST + 2500 giorni.
    Il parametro mxaz è ignorato (mantiene compatibilità con le chiamate precedenti).
    Ritorna stringa ISO 'YYYY-MM-DD' oppure None.
    """
    if not datainst:
        return None
    try:
        datainst_iso = _parse_date_param_to_iso(datainst)
        d = datetime.fromisoformat(datainst_iso)
        # FINEMF = DATAINST + 2500 giorni (fissi)
        fine = d + timedelta(days=2500)
        return fine.date().isoformat()
    except Exception:
        return None

def calcola_scadenza_garanzia(datainst, durata_mesi):
    """
    Calcola scadenza garanzia dato DATAINST e durata in mesi.
    Ritorna stringa ISO 'YYYY-MM-DD' oppure None.
    """
    if not datainst or not durata_mesi:
        return None
    try:
        datainst_iso = _parse_date_param_to_iso(datainst)
        d = datetime.fromisoformat(datainst_iso)
        mesi = int(str(durata_mesi).strip())
        # aggiungi mesi
        anno = d.year
        mese = d.month + mesi
        while mese > 12:
            mese -= 12
            anno += 1
        scad = d.replace(year=anno, month=mese)
        return scad.date().isoformat()
    except Exception:
        return None

def calcola_codpuntocassa(table_name: str = "clienti"):
    """
    Ricalcola CODPUNTOCA per la tabella specificata.
    CODPUNTOCA = P_IVA + COMUNE (semplice concatenazione).
    Se non trova le colonne, ignora silenziosamente.
    """
    conn = get_connection()
    cur = conn.cursor()
    try:
        cur.execute(f"PRAGMA table_info({table_name})")
        cols = [r[1] for r in cur.fetchall()]
        if "CODPUNTOCA" not in cols or "P_IVA" not in cols or "COMUNE" not in cols:
            return
        
        cur.execute(f"SELECT id, P_IVA, COMUNE FROM {table_name}")
        rows = cur.fetchall()
        
        for rid, piva, comune in rows:
            codpuntoca = f"{piva or ''}{comune or ''}".strip()
            cur.execute(f"UPDATE {table_name} SET CODPUNTOCA=? WHERE id=?", (codpuntoca, rid))
        
        conn.commit()
    except Exception:
        pass
    finally:
        conn.close()

def ensure_db_exists(recreate_clienti: bool = False):
    """Alias per init_database (per compatibilità GUI)."""
    init_database(recreate_clienti)

def read_intae(limit: int = 1000):
    """Legge tutti i record INTAE (con id visibile come prima colonna)."""
    conn = get_connection()
    cur = conn.cursor()
    cur.execute("PRAGMA table_info(intae)")
    cols = [r[1] for r in cur.fetchall()]
    cols_sql = ", ".join([f'"{c}"' for c in cols])
    cur.execute(f"SELECT {cols_sql} FROM intae ORDER BY id DESC LIMIT ?", (limit,))
    rows = cur.fetchall()
    conn.close()
    # ritorna le righe come sono (id è già il primo elemento)
    return rows

def cerca_interventi_intae(matricola: str):
    """Cerca interventi INTAE per matricola con matching flessibile."""
    conn = get_connection()
    cur = conn.cursor()
    mat_input = str(matricola).strip().lower()
    
    # Prova match esatto
    cur.execute("SELECT * FROM intae WHERE lower(MATRICOLA) = ?", (mat_input,))
    rows = cur.fetchall()
    
    # Se no match esatto, prova LIKE
    if not rows:
        cur.execute("SELECT * FROM intae WHERE lower(MATRICOLA) LIKE ?", (f"%{mat_input}%",))
        rows = cur.fetchall()
    
    conn.close()
    return rows or []

def _import_intae_csv_generic(path: str, delimiter: str | None = None, encoding: str = "utf-8") -> dict:
    if not os.path.exists(path):
        raise FileNotFoundError(path)

    result = {"ok": False, "inserted": 0, "skipped": 0, "skipped_rows": [], "columns": INTAE_COLUMNS, "error": None, "message": ""}

    # detect delimiter se non fornito
    with open(path, "r", encoding=encoding, errors="replace") as sf:
        sample = sf.read(8192)
        sf.seek(0)
        if delimiter is None:
            for d in (";", ",", "\t", "|"):
                if sample.count(d) > 0:
                    delimiter = d
                    break
            else:
                delimiter = ","

    try:
        with open(path, "r", encoding=encoding, errors="replace") as fh:
            reader = csv.reader(fh, delimiter=delimiter)
            rows = list(reader)

        if not rows:
            result.update({"ok": True, "inserted": 0, "skipped": 0, "message": "File vuoto"})
            return result

        headers = [str(h or "").strip() for h in rows[0]]
        data_rows = rows[1:]

        def is_generic_header(h):
            return bool(re.fullmatch(r"N\d+", str(h or "").strip(), re.IGNORECASE)) or str(h or "").strip() == ""

        known_keys = _intae_known_header_keys()

        header_map = _build_intae_header_map(headers)

        # Se intestazioni generiche o mappatura scarsa, prova a usare la prima riga come header
        if headers and (sum(1 for h in headers if is_generic_header(h)) >= max(1, len(headers) // 2) or len(header_map) < 2):
            if data_rows:
                candidate_headers = [str(h or "").strip() for h in data_rows[0]]
                candidate_map = _build_intae_header_map(candidate_headers)
                candidate_hits = sum(1 for h in candidate_headers if _normalize_header_key(h) in known_keys)
                if len(candidate_map) > len(header_map) or candidate_hits >= 2:
                    headers = candidate_headers
                    data_rows = data_rows[1:]
                    header_map = candidate_map

        # Validazione RIGOROSA: minimo 8 colonne INTAE valide su 12 totali (67%)
        # Le colonne essenziali sono: CODICE_INTERVENTO, PIVA, DATA_INTERVENTO, CF_TECNICO, STATO_CF, DISPOSITIVO, MATRICOLA, TIPOLOGIA
        MIN_REQUIRED_INTAE_COLUMNS = 8
        essential_cols = {"CODICE_INTERVENTO", "PIVA", "DATA_INTERVENTO", "CF_TECNICO", "STATO_CF", "DISPOSITIVO", "MATRICOLA", "TIPOLOGIA"}
        found_cols = set(header_map.values()) if header_map else set()
        essential_found = found_cols & essential_cols
        
        if len(header_map) < MIN_REQUIRED_INTAE_COLUMNS or len(essential_found) < 6:
            result["error"] = f"File non valido per INTAE: trovate solo {len(found_cols)} colonne INTAE valide (minimo {MIN_REQUIRED_INTAE_COLUMNS} richieste).\nColonne trovate: {', '.join(sorted(found_cols)) if found_cols else 'nessuna'}\nColonne essenziali mancanti: {', '.join(sorted(essential_cols - essential_found))}"
            result["message"] = result["error"]
            return result

        conn = get_connection()
        cur = conn.cursor()
        inserted = 0
        skipped = 0
        skipped_rows = []

        for i, row in enumerate(data_rows, start=1):
            try:
                # costruisci dict header->valore
                if not row:
                    skipped += 1
                    skipped_rows.append({"rownum": i, "reason": "riga vuota", "data": row})
                    continue

                padded = list(row) + [""] * max(0, len(headers) - len(row))
                row_dict = {headers[idx]: padded[idx] for idx in range(len(headers))}

                rec = _intae_build_record(row_dict, header_map)

                # validazione minima
                if not rec.get("MATRICOLA") or not rec.get("DATA_INTERVENTO"):
                    skipped += 1
                    skipped_rows.append({"rownum": i, "reason": "manca MATRICOLA o DATA_INTERVENTO", "data": row_dict})
                    continue

                # previene duplicati MATRICOLA+DATA_INTERVENTO
                cur.execute(
                    "SELECT 1 FROM intae WHERE lower(ifnull(MATRICOLA,''))=? AND ifnull(DATA_INTERVENTO,'')=? LIMIT 1",
                    (rec.get("MATRICOLA", "").lower(), rec.get("DATA_INTERVENTO", ""))
                )
                if cur.fetchone():
                    skipped += 1
                    skipped_rows.append({"rownum": i, "reason": "duplicato (matricola+data)", "data": row_dict})
                    continue

                cols = [c for c in INTAE_COLUMNS if rec.get(c) is not None]
                vals = [rec.get(c, "") for c in cols]
                placeholders = ",".join(["?"] * len(cols))
                cols_join = ",".join([f'"{c}"' for c in cols])
                cur.execute(f'INSERT INTO intae ({cols_join}) VALUES ({placeholders})', tuple(vals))
                inserted += 1
                if inserted % 200 == 0:
                    conn.commit()
            except Exception as e:
                skipped += 1
                skipped_rows.append({"rownum": i, "reason": f"errore insert: {e}", "data": row})

        conn.commit()
        conn.close()
        result.update({"ok": True, "inserted": inserted, "skipped": skipped, "skipped_rows": skipped_rows, "message": f"Inserted {inserted}, skipped {skipped}"})
        return result
    except Exception as e:
        result["error"] = str(e)
        result["message"] = str(e)
        return result

def import_csv_intae(path: str, delimiter: str | None = None, encoding: str = "utf-8") -> dict:
    return _import_intae_csv_generic(path, delimiter, encoding)

def read_misurart_all(limit: int = 1000):
    """Legge tutti i record MISURART."""
    conn = get_connection()
    cur = conn.cursor()
    cur.execute("PRAGMA table_info(misurart)")
    cols = [r[1] for r in cur.fetchall()]
    cols_sql = ", ".join([f'"{c}"' for c in cols])
    cur.execute(f"SELECT {cols_sql} FROM misurart ORDER BY id DESC LIMIT ?", (limit,))
    rows = cur.fetchall()
    conn.close()
    return _misurart_sort_rows(rows, cols)

def _misurart_sort_rows(rows: list, cols: list) -> list:
    if not rows or not cols:
        return rows
    date_idx = None
    for c in ("DATAVP", "DATAINIZIO", "DATAFINE"):
        if c in cols:
            date_idx = cols.index(c)
            break
    if date_idx is None:
        return rows

    def sort_key(row):
        try:
            raw = row[date_idx] if date_idx < len(row) else ""
            raw = "" if raw is None else str(raw).strip()
            if raw:
                parsed = _parse_date(raw)
                if isinstance(parsed, datetime):
                    return parsed.date()
                if isinstance(parsed, date):
                    return parsed
        except Exception:
            pass
        return date.min

    return sorted(rows, key=sort_key, reverse=True)

def _normalize_misurart_dates(data: dict) -> dict:
    """Normalizza le date in MISURART da DD/MM/YYYY a YYYY-MM-DD."""
    normalized = dict(data)
    date_fields = {"DATAINIZIO", "DATAFINE", "DATAVP"}
    for field in date_fields:
        if field in normalized:
            val = normalized[field]
            if val and isinstance(val, str):
                val = val.strip()
                if "/" in val:
                    try:
                        normalized[field] = _parse_date_param_to_iso(val)
                    except Exception:
                        pass
    return normalized

def misurart_add_from_map(data: dict) -> int:
    """Inserisce un record MISURART da dict. Ritorna id oppure None."""
    if not data:
        return None
    
    # Normalizza le date prima di salvare
    data = _normalize_misurart_dates(data)
    
    conn = get_connection()
    cur = conn.cursor()
    
    # usa solo colonne che esistono
    cols = [c for c in MISURART_COLUMNS if c in data or str(data.get(c, "")).strip()]
    if not cols:
        conn.close()
        return None
    
    vals = [str(data.get(c, "")).strip() for c in cols]
    placeholders = ",".join("?" * len(cols))
    cols_sql = ", ".join([f'"{c}"' for c in cols])
    sql = f"INSERT INTO misurart ({cols_sql}) VALUES ({placeholders})"
    
    try:
        cur.execute(sql, vals)
        nid = cur.lastrowid
        conn.commit()
        return nid
    except Exception:
        conn.rollback()
        return None
    finally:
        conn.close()

def misurart_update(rowid: int, data: dict):
    """Aggiorna un record MISURART."""
    if not data or not rowid:
        return 0
    
    # Normalizza le date prima di salvare
    data = _normalize_misurart_dates(data)
    
    conn = get_connection()
    cur = conn.cursor()
    
    set_parts = [f'"{c}"=?' for c in data.keys()]
    vals = [str(data[c]).strip() for c in data.keys()] + [rowid]
    sql = f"UPDATE misurart SET {', '.join(set_parts)} WHERE id=?"
    
    try:
        cur.execute(sql, vals)
        affected = cur.rowcount
        conn.commit()
        return affected
    except Exception:
        conn.rollback()
        return 0
    finally:
        conn.close()

def misurart_delete(rowid: int):
    """Elimina un record MISURART."""
    conn = get_connection()
    cur = conn.cursor()
    cur.execute("DELETE FROM misurart WHERE id=?", (rowid,))
    conn.commit()
    conn.close()

def misurart_search_by_matricola(matricola: str):
    """Cerca record MISURART per matricola con matching flessibile."""
    conn = get_connection()
    cur = conn.cursor()
    cur.execute("PRAGMA table_info(misurart)")
    cols = [r[1] for r in cur.fetchall()]
    cols_sql = ", ".join([f'"{c}"' for c in cols])
    cur.execute(f"SELECT {cols_sql} FROM misurart")
    rows = cur.fetchall()
    conn.close()
    mat_input = _normalize_matricola_for_compare(matricola)
    if not mat_input:
        return []
    if "MATRICOLA" not in cols:
        return []
    mat_idx = cols.index("MATRICOLA")

    filtered = []
    for r in rows:
        raw = r[mat_idx] if mat_idx < len(r) else ""
        raw_norm = _normalize_matricola_for_compare(raw)
        if raw_norm and mat_input in raw_norm:
            filtered.append(r)

    return _misurart_sort_rows(filtered, cols)

def verify_and_register_misurart(matricola: str, codice: int = 3) -> dict:
    """
    Verifica una matricola in INTAE/clienti e prepara i dati per MISURART.
    Ritorna dict {'ok': bool, 'cliente_nome': str, 'data_intervento': str, 'operatore': str, 'error': str}
    """
    result = {
        "ok": False,
        "cliente_nome": "",
        "data_intervento": "",
        "operatore": "",
        "error": None
    }
    
    if not matricola or not str(matricola).strip():
        result["error"] = "Matricola richiesta"
        return result
    
    conn = get_connection()
    cur = conn.cursor()
    
    try:
        # cerca in clienti
        cur.execute("SELECT RAG_SOC, DATAINST FROM clienti WHERE lower(trim(MATRICOLA)) = ?", (str(matricola).lower(),))
        c_row = cur.fetchone()
        
        if c_row:
            result["cliente_nome"] = c_row[0] or ""
            result["data_intervento"] = c_row[1] or ""
        
        # cerca in INTAE per ricavare DATA_INTERVENTO se più recente
        cur.execute("SELECT DATA_INTERVENTO, CF_TECNICO FROM intae WHERE lower(trim(MATRICOLA)) = ? ORDER BY DATA_INTERVENTO DESC LIMIT 1", (str(matricola).lower(),))
        i_row = cur.fetchone()
        
        if i_row:
            if i_row[0] and not result["data_intervento"]:
                result["data_intervento"] = i_row[0]
            if i_row[1]:
                result["operatore"] = i_row[1]
        
        result["ok"] = True
    except Exception as e:
        result["error"] = str(e)
    finally:
        conn.close()
    
    return result

def report_cambi_dgfe(date_from: str, date_to: str):
    """
    Cerca tutte le matricole che hanno ricevuto interventi nel range di date specificato.
    Per ogni matricola, conta quanti interventi hanno CODICE_INTERVENTO = 2 (cambio DGFE).
    Ritorna lista di tuple: [(ragione_sociale, matricola, count), ...]
    
    LOGICA:
    1. Filtra tabella INTAE per DATA_INTERVENTO nell'intervallo specificato
    2. Filtra per CODICE_INTERVENTO = '2' (con TRIM per gestire spazi)
    3. Raggruppa per MATRICOLA (completa, con logotipo)
    4. Conta quante righe ci sono per ogni matricola (COUNT(*))
    5. Fa JOIN con CLIENTI per recuperare RAG_SOC (rimuovendo i primi 2 caratteri della matricola INTAE)
    
    NOTA: Se vedi un conteggio maggiore di quanto ti aspetti, verifica nel database INTAE
    se ci sono record duplicati per quella matricola nell'intervallo di date.
    """
    conn = get_connection()
    cur = conn.cursor()

    # helper di parsing locale 
    def _to_dateobj(s):
        iso = _parse_date(s)
        # _parse_date ritorna un datetime object, estrarre direttamente .date()
        if isinstance(iso, str):
            # Se è una stringa, parsarla
            iso_date_only = iso.split()[0] if ' ' in iso else iso
            return datetime.strptime(iso_date_only, "%Y-%m-%d").date()
        else:
            # Se è un datetime object, prendere direttamente la data
            return iso.date()

    # normalizza parametri
    try:
        d1 = _to_dateobj(date_from) if date_from else date(1900,1,1)
    except Exception:
        d1 = date(1900,1,1)
    try:
        d2 = _to_dateobj(date_to) if date_to else date.today()
    except Exception:
        d2 = date.today()

    d1_iso = d1.isoformat()
    d2_iso = d2.isoformat()
    
    print(f"[DEBUG DGFE] Cerco interventi dal {d1_iso} al {d2_iso}")

    try:
        # Step 1: Prendi tutte le righe da INTAE nell'intervallo di date
        query_intae = """
        SELECT MATRICOLA, CODICE_INTERVENTO, DATA_INTERVENTO
        FROM intae
        WHERE DATA_INTERVENTO >= ? AND DATA_INTERVENTO <= ?
        ORDER BY MATRICOLA, DATA_INTERVENTO
        """
        cur.execute(query_intae, (d1_iso, d2_iso))
        all_rows = cur.fetchall()
        
        print(f"\n[DEBUG DGFE] =================================================")
        print(f"[DEBUG DGFE] REPORT CAMBI DGFE")
        print(f"[DEBUG DGFE] Intervallo: dal {d1_iso} al {d2_iso}")
        print(f"[DEBUG DGFE] Totale righe INTAE nell'intervallo: {len(all_rows)}")
        print(f"[DEBUG DGFE] =================================================")
        
        if len(all_rows) == 0:
            print(f"[DEBUG DGFE] NESSUN intervento trovato nell'intervallo specificato!")
            conn.close()
            return []
        
        if all_rows:
            print(f"\n[DEBUG DGFE] --- TUTTE LE RIGHE TROVATE ---")
            for i, row in enumerate(all_rows, 1):
                mat, cod, data = row
                print(f"[DEBUG DGFE] #{i}: MAT='{mat}' COD='{cod}' DATA='{data}'")
        
        # Step 2: Conta per ogni matricola SOLO i codici '2'
        matricola_count = {}
        matricola_details = {}  # Memorizza dettagli completi per debug
        
        for matricola, codice, data_int in all_rows:
            if not matricola:
                print(f"[DEBUG DGFE] SALTATA riga con MATRICOLA vuota")
                continue
                
            mat_key = str(matricola).strip().upper()
            codice_str = str(codice).strip()
            
            # Inizializza strutture
            if mat_key not in matricola_count:
                matricola_count[mat_key] = 0
                matricola_details[mat_key] = []
            
            # Memorizza TUTTI i codici per questa matricola
            matricola_details[mat_key].append({
                'codice': codice_str,
                'data': data_int
            })
            
            # Incrementa SOLO se CODICE_INTERVENTO è esattamente '2'
            if codice_str == '2':
                matricola_count[mat_key] += 1
                print(f"[DEBUG DGFE] >>> CONTEGGIO: {mat_key} - codice '2' trovato (totale finora: {matricola_count[mat_key]})")
        
        # Step 3: Riepilogo dettagliato per matricola
        if matricola_details:
            print(f"\n[DEBUG DGFE] --- RIEPILOGO DETTAGLIATO PER MATRICOLA ---")
            for mat_key in sorted(matricola_details.keys()):
                details = matricola_details[mat_key]
                count_2 = matricola_count.get(mat_key, 0)
                codici_list = [d['codice'] for d in details]
                print(f"[DEBUG DGFE] {mat_key}:")
                print(f"[DEBUG DGFE]   - Interventi totali: {len(details)}")
                print(f"[DEBUG DGFE]   - Codici: {codici_list}")
                print(f"[DEBUG DGFE]   - Conteggio CODICE='2': {count_2}")
                for d in details:
                    print(f"[DEBUG DGFE]     * COD={d['codice']} DATA={d['data']}")
        
        # Step 4: Filtra solo matricole con almeno un intervento codice=2
        matricole_con_codice2 = [(mat, cnt) for mat, cnt in matricola_count.items() if cnt > 0]
        
        if matricole_con_codice2:
            print(f"\n[DEBUG DGFE] Matricole con almeno 1 intervento CODICE='2': {len(matricole_con_codice2)}")
            for mat, cnt in matricole_con_codice2:
                print(f"[DEBUG DGFE]   - {mat}: {cnt} intervento/i")
        else:
            print(f"\n[DEBUG DGFE] Nessuna matricola con interventi CODICE='2'")
        
        # Step 5: Per ogni matricola, recupera la ragione sociale da CLIENTI
        out = []
        if matricole_con_codice2:
            print(f"\n[DEBUG DGFE] --- RICERCA CLIENTI PER RAGIONE SOCIALE ---")
        for mat_display, cnt in matricole_con_codice2:
            # Rimuovi i primi 2 caratteri (logotipo) per cercare in CLIENTI
            mat_no_logo = mat_display[2:] if len(mat_display) > 2 else mat_display
            
            print(f"[DEBUG DGFE] Cerco in CLIENTI: MATRICOLA='{mat_no_logo}' (originale: '{mat_display}')")
            
            cur.execute(
                "SELECT RAG_SOC FROM clienti WHERE LOWER(TRIM(MATRICOLA)) = LOWER(TRIM(?)) LIMIT 1",
                (mat_no_logo,)
            )
            cliente_row = cur.fetchone()
            rag_display = cliente_row[0] if cliente_row else "N/D"
            
            print(f"[DEBUG DGFE]   -> RAG_SOC: '{rag_display}'")
            
            out.append((str(rag_display).strip(), mat_display, cnt))
        
        # Ordina per conteggio decrescente
        out.sort(key=lambda x: x[2], reverse=True)
        
        print(f"\n[DEBUG DGFE] --- RISULTATO FINALE ---")
        for rag, mat, cnt in out:
            print(f"[DEBUG DGFE] {rag} | {mat} | {cnt} cambi DGFE")
        
        print(f"[DEBUG DGFE] =================================================")
        print(f"[DEBUG DGFE] FINE REPORT\n")
        
        conn.close()
        return out
        
    except Exception as e:
        print(f"[DGFE Error] {e}")
        import traceback
        traceback.print_exc()
        conn.close()
        return []
        return out
        
    except Exception as e:
        print(f"[DGFE Error] {e}")
        conn.close()
        return []

def get_verification_data(matricola: str) -> dict:
    """
    Recupera i dati completi per la verifica di una macchina.
    Ricerca in CLIENTI e INTAE. Restituisce dict con dati cliente, dati INTAE e informazioni sulla verifica.
    
    Le matricole in INTAE hanno i primi 2 caratteri che rappresentano il logotipo.
    Per cercare in CLIENTI, rimuoviamo i primi 2 caratteri dalla matricola input.
    Per cercare in INTAE, usiamo la matricola completa.
    """
    result = {
        "ok": False,
        "error": None,
        "cliente": {},
        "intae_data": {},
        "verificatore": "",
        "data_intervento": "",
        "codice_intervento": ""
    }
    
    if not matricola or not str(matricola).strip():
        result["error"] = "Matricola richiesta"
        return result
    
    conn = get_connection()
    cur = conn.cursor()
    
    try:
        mat_normalized = str(matricola).lower().strip()
        
        # Per CLIENTI: la matricola è senza logotipo (rimuovi primi 2 caratteri se len > 2)
        mat_no_logo = mat_normalized[2:] if len(mat_normalized) > 2 else mat_normalized
        mat_like_no_logo = f"%{mat_no_logo}%"
        
        # Per INTAE: la matricola è completa (con logotipo)
        mat_like_full = f"%{mat_normalized}%"
        
        # Recupera dati cliente - Cerca con matricola senza logotipo
        cur.execute("SELECT * FROM clienti WHERE lower(trim(MATRICOLA)) = ? LIMIT 1", (mat_no_logo,))
        cols_client = [r[1] for r in cur.execute("PRAGMA table_info(clienti)").fetchall()]
        client_row = cur.fetchone()
        
        # Se non trova corrispondenza esatta, prova con LIKE
        if not client_row:
            cur.execute("SELECT * FROM clienti WHERE lower(MATRICOLA) LIKE ? LIMIT 1", (mat_like_no_logo,))
            client_row = cur.fetchone()
        
        if client_row:
            result["cliente"] = {cols_client[i]: client_row[i] for i in range(len(cols_client))}
        else:
            # Prova a cercare in INTAE se non trovato in CLIENTI (con matricola completa)
            cur.execute("""
                SELECT CODICE_INTERVENTO, DATA_INTERVENTO, CF_TECNICO, STATO_CF, 
                       DISPOSITIVO, TIPOLOGIA, DATA_TRASMISSIONE, NOTA, VERSIONE_SW, DATA_RELEASE
                FROM intae 
                WHERE lower(trim(MATRICOLA)) = ? 
                ORDER BY DATA_INTERVENTO DESC 
                LIMIT 1
            """, (mat_normalized,))
            
            intae_row = cur.fetchone()
            if intae_row:
                intae_cols = ["CODICE_INTERVENTO", "DATA_INTERVENTO", "CF_TECNICO", "STATO_CF",
                             "DISPOSITIVO", "TIPOLOGIA", "DATA_TRASMISSIONE", "NOTA", "VERSIONE_SW", "DATA_RELEASE"]
                result["intae_data"] = {intae_cols[i]: intae_row[i] for i in range(len(intae_cols))}
                result["verificatore"] = intae_row[2] or ""
                result["data_intervento"] = intae_row[1] or ""
                result["codice_intervento"] = intae_row[0] or ""
                result["ok"] = True
                conn.close()
                return result
            
            result["error"] = f"Matricola '{matricola}' non trovata in CLIENTI né in INTAE"
            conn.close()
            return result
        
        # Recupera ultime info da INTAE (con data più recente)
        # Cerca sempre con LIKE per trovare matricole con o senza logotipo
        search_pattern = f"%{mat_normalized}%"
        print(f"[DEBUG get_verification_data] Cerco INTAE con LIKE: {search_pattern} (matricola input: {matricola})")
        
        cur.execute("""
            SELECT CODICE_INTERVENTO, DATA_INTERVENTO, CF_TECNICO, STATO_CF, 
                   DISPOSITIVO, TIPOLOGIA, DATA_TRASMISSIONE, NOTA, VERSIONE_SW, DATA_RELEASE
            FROM intae 
            WHERE upper(trim(MATRICOLA)) LIKE upper(?)
            ORDER BY DATA_INTERVENTO DESC 
            LIMIT 1
        """, (search_pattern,))
        
        intae_row = cur.fetchone()
        
        if intae_row:
            intae_cols = ["CODICE_INTERVENTO", "DATA_INTERVENTO", "CF_TECNICO", "STATO_CF",
                         "DISPOSITIVO", "TIPOLOGIA", "DATA_TRASMISSIONE", "NOTA", "VERSIONE_SW", "DATA_RELEASE"]
            result["intae_data"] = {intae_cols[i]: intae_row[i] for i in range(len(intae_cols))}
            result["verificatore"] = intae_row[2] or ""  # CF_TECNICO
            result["data_intervento"] = intae_row[1] or ""  # DATA_INTERVENTO
            result["codice_intervento"] = intae_row[0] or ""  # CODICE_INTERVENTO
            print(f"[DEBUG get_verification_data] ✓ Trovato INTAE per matricola {matricola}: DATA_INTERVENTO={intae_row[1]}")
        else:
            print(f"[DEBUG get_verification_data] ✗ Nessun dato INTAE trovato per matricola {matricola}")
            # Tenta ricerca alternativa senza LIKE, cercando solo gli ultimi 10 caratteri
            alt_pattern = f"%{mat_normalized[-10:] if len(mat_normalized) > 10 else mat_normalized}%"
            print(f"[DEBUG get_verification_data] Tentativo alternativo con pattern: {alt_pattern}")
            cur.execute("""
                SELECT CODICE_INTERVENTO, DATA_INTERVENTO, CF_TECNICO, STATO_CF, 
                       DISPOSITIVO, TIPOLOGIA, DATA_TRASMISSIONE, NOTA, VERSIONE_SW, DATA_RELEASE
                FROM intae 
                WHERE upper(trim(MATRICOLA)) LIKE upper(?)
                ORDER BY DATA_INTERVENTO DESC 
                LIMIT 1
            """, (alt_pattern,))
            intae_row = cur.fetchone()
            if intae_row:
                intae_cols = ["CODICE_INTERVENTO", "DATA_INTERVENTO", "CF_TECNICO", "STATO_CF",
                             "DISPOSITIVO", "TIPOLOGIA", "DATA_TRASMISSIONE", "NOTA", "VERSIONE_SW", "DATA_RELEASE"]
                result["intae_data"] = {intae_cols[i]: intae_row[i] for i in range(len(intae_cols))}
                result["verificatore"] = intae_row[2] or ""
                result["data_intervento"] = intae_row[1] or ""
                result["codice_intervento"] = intae_row[0] or ""
                print(f"[DEBUG get_verification_data] ✓ Trovato con ricerca alternativa: DATA_INTERVENTO={intae_row[1]}")
        
        result["ok"] = True
        
    except Exception as e:
        result["error"] = str(e)
        result["ok"] = False
    finally:
        conn.close()
    
    return result

def update_verification_data(cliente_id: int, vf_date: str, azz_value: str, misurart_data: dict) -> dict:
    """
    Aggiorna VF (data verificazione) e AZZ in clienti, e inserisce dati in MISURART.
    Ritorna dict con risultato operazione.
    """
    result = {
        "ok": False,
        "error": None,
        "cliente_updated": False,
        "misurart_id": None
    }
    
    conn = get_connection()
    
    try:
        cur = conn.cursor()
        
        # Aggiorna VF e AZZ nella tabella clienti
        if vf_date or azz_value:
            updates = {}
            if vf_date:
                updates["VF"] = vf_date
            if azz_value:
                updates["AZZ"] = azz_value
            
            if updates:
                set_parts = [f'"{k}"=?' for k in updates.keys()]
                vals = list(updates.values()) + [cliente_id]
                sql = f"UPDATE clienti SET {', '.join(set_parts)} WHERE id=?"
                cur.execute(sql, vals)
                result["cliente_updated"] = cur.rowcount > 0
        
        # Inserisci in MISURART
        if misurart_data:
            cols = [c for c in misurart_data.keys() if misurart_data[c] is not None and str(misurart_data[c]).strip()]
            if cols:
                vals = [str(misurart_data[c]).strip() for c in cols]
                placeholders = ",".join("?" * len(cols))
                cols_sql = ", ".join([f'"{c}"' for c in cols])
                sql = f"INSERT INTO misurart ({cols_sql}) VALUES ({placeholders})"
                cur.execute(sql, vals)
                result["misurart_id"] = cur.lastrowid
        
        conn.commit()
        result["ok"] = True
        
    except Exception as e:
        conn.rollback()
        result["error"] = str(e)
    finally:
        conn.close()
    
    return result

# -------------------------
# DBF file handling
# -------------------------

def read_dbf_file(path: str) -> tuple[list, list]:
    """
    Legge un file DBF e ritorna (headers, rows).
    headers: lista di nomi colonne
    rows: lista di tuple (ogni riga)
    Solleva eccezione se il file non esiste o non può essere letto.
    """
    if DBF is None:
        raise RuntimeError("Libreria 'dbfread' non installata. Installa con: pip install dbfread")
    
    if not os.path.exists(path):
        raise FileNotFoundError(f"File DBF non trovato: {path}")
    
    try:
        dbf = DBF(path, encoding='latin-1', ignore_missing_memofile=True)
        # Estrai i nomi dei campi come sono nel DBF (preservando il case)
        field_names = [f.name for f in dbf.fields]
        headers = [f.upper() for f in field_names]
        
        rows = []
        for record in dbf:
            # Accedi ai campi usando i nomi originali del DBF (case-sensitive)
            row = tuple(str(record[field_names[i]]).strip() if record[field_names[i]] is not None else "" 
                       for i in range(len(field_names)))
            rows.append(row)
        return headers, rows
    except Exception as e:
        raise RuntimeError(f"Errore nella lettura del file DBF: {str(e)}")

def validate_import_headers(file_headers: list, target_columns: list, file_type: str = "generico") -> dict:
    """
    Valida che tutti gli header del file siano presenti nel database target.
    Ritorna dict {ok, missing_headers, message}
    """
    result = {
        "ok": True,
        "missing_headers": [],
        "message": ""
    }
    
    # Normalizza i nomi delle colonne per il confronto
    def normalize_col(col: str) -> str:
        return str(col).strip().upper()
    
    normalized_target = {normalize_col(c): c for c in target_columns}
    
    for header in file_headers:
        norm_header = normalize_col(header)
        if norm_header not in normalized_target:
            result["missing_headers"].append(header)
    
    if result["missing_headers"]:
        result["ok"] = False
        missing_str = ", ".join(result["missing_headers"])
        result["message"] = f"[{file_type}] Colonne non trovate nel database: {missing_str}"
    else:
        result["message"] = f"[{file_type}] Validazione header riuscita: tutte le colonne trovate"
    
    return result

def import_clients_from_dbf(path: str) -> dict:
    """
    Importa clienti da file DBF nella tabella clienti.
    Validazione migliorata: verifica che il file DBF contenga tutti i campi richiesti.
    Ritorna dict con {ok, inserted, skipped, skipped_rows, columns, message}.
    """
    if not os.path.exists(path):
        raise FileNotFoundError(path)
    
    # Leggi file DBF
    try:
        headers, rows = read_dbf_file(path)
    except Exception as e:
        return {
            "ok": False, 
            "inserted": 0, 
            "skipped": 0, 
            "skipped_rows": [], 
            "columns": [], 
            "message": f"Errore lettura file DBF: {e}",
            "error": f"Errore lettura file DBF: {e}"
        }
    
    # VALIDAZIONE STRUTTURALE: il file deve contenere almeno i campi fondamentali
    # Normalizza gli header per confronto (case-insensitive)
    def normalize_col(col: str) -> str:
        return str(col).strip().upper()
    
    headers_norm = {normalize_col(h): h for h in headers}
    
    # Campi richiesti OBBLIGATORIAMENTE nel DBF
    required_dbf_fields = ["RAG_SOC", "MATRICOLA", "MARCA"]
    missing_required = []
    for req in required_dbf_fields:
        if normalize_col(req) not in headers_norm:
            missing_required.append(req)
    
    if missing_required:
        return {
            "ok": False,
            "inserted": 0,
            "skipped": 0,
            "skipped_rows": [],
            "columns": [],
            "message": f"File DBF non valido: mancano i campi obbligatori {', '.join(missing_required)}",
            "error": f"Campi obbligatori mancanti: {', '.join(missing_required)}"
        }
    
    # Verifica che il file contenga ALMENO il 70% dei campi CLIENTI_COLUMNS
    matched_count = sum(1 for col in CLIENTI_COLUMNS if normalize_col(col) in headers_norm)
    coverage = matched_count / len(CLIENTI_COLUMNS) if CLIENTI_COLUMNS else 0
    
    if coverage < 0.70:
        return {
            "ok": False,
            "inserted": 0,
            "skipped": 0,
            "skipped_rows": [],
            "columns": [],
            "message": f"File DBF non valido: copertura campi insufficiente ({coverage*100:.0f}% trovato, richiesto almeno 70%)",
            "error": f"Copertura campi insufficiente: {matched_count}/{len(CLIENTI_COLUMNS)} campi trovati"
        }
    
    # Verifica che ci siano dati (almeno 1 record valido)
    if not rows or len(rows) == 0:
        return {
            "ok": False,
            "inserted": 0,
            "skipped": 0,
            "skipped_rows": [],
            "columns": [],
            "message": "File DBF vuoto: nessun record trovato",
            "error": "Nessun record trovato nel file"
        }
    
    # Procedi con import (simile a import_clients_from_csv_generic)
    conn = get_connection()
    cur = conn.cursor()
    
    try:
        cur.execute("PRAGMA table_info(clienti)")
        pragma = cur.fetchall()
        db_cols = [row[1] for row in pragma]
    except Exception:
        db_cols = CLIENTI_COLUMNS.copy()
    
    if not db_cols:
        conn.close()
        return {"ok": False, "inserted": 0, "skipped": 0, "skipped_rows": [], "columns": [], "message": "Impossibile determinare colonne tabella 'clienti'"}
    
    # Map DBF headers a DB columns - usa mappatura ESATTA per campi standard
    dbf_to_db = {}
    used_targets = set()
    
    # Prima mappa i campi con nome esatto (case-insensitive)
    for h in headers:
        h_norm = normalize_col(h)
        # Cerca corrispondenza esatta nelle colonne DB
        for target in db_cols:
            if normalize_col(target) == h_norm and target not in used_targets:
                dbf_to_db[h] = target
                used_targets.add(target)
                break
    
    # Poi mappa i restanti con token matching (solo se non già mappati)
    def tokens(s: str):
        return [t for t in re.split(r"[^A-Z0-9]+", normalize_col(s)) if t]

    for h in headers:
        if h in dbf_to_db:  # già mappato
            continue
        best_target = None
        best_score = 0
        h_tokens = set(tokens(h))
        if not h_tokens:
            continue
        for target in db_cols:
            if target in used_targets:
                continue
            t_tokens = set(tokens(target))
            score = len(h_tokens & t_tokens)
            if score > best_score:
                best_score = score
                best_target = target
        if best_score > 0 and best_target:
            dbf_to_db[h] = best_target
            used_targets.add(best_target)
    
    csv_to_db_cols = list(dict.fromkeys(dbf_to_db.values()))
    if not csv_to_db_cols:
        conn.close()
        return {"ok": False, "inserted": 0, "skipped": 0, "skipped_rows": [], "columns": [], "message": "Nessuna colonna del DBF mappabile alle colonne DB"}
    
    # Ordina secondo ordine DB
    columns_to_insert = sorted(csv_to_db_cols, key=lambda x: db_cols.index(x) if x in db_cols else 0)
    
    inserted = 0
    skipped = 0
    skipped_rows = []
    imported_ids = []  # traccia gli id importati per post-processing TIPOMF
    
    for row_data in rows:
        # Crea dict da headers
        row_dict = {headers[i]: row_data[i] for i in range(len(headers))}
        
        # Estrai valori per colonne mappate
        vals = []
        for col_db in columns_to_insert:
            # Trova l'header DBF corrispondente
            dbf_header = [h for h in headers if dbf_to_db.get(h) == col_db]
            if dbf_header:
                val = row_dict[dbf_header[0]].strip()
            else:
                val = ""
            vals.append(val or None)

        # helper per recuperare valore dato nome colonna DB
        def val_for(colname):
            try:
                idx = columns_to_insert.index(colname)
                return vals[idx]
            except ValueError:
                return None

        # MODIFICATO: NON saltare righe incomplete - importare TUTTE le righe
        
        # Calcola GARANZIA e FINEMF se mancanti
        if "GARANZIA" in columns_to_insert and (not vals[columns_to_insert.index("GARANZIA")] or vals[columns_to_insert.index("GARANZIA")] == ""):
            datainst_idx = columns_to_insert.index("DATAINST") if "DATAINST" in columns_to_insert else -1
            if datainst_idx >= 0 and vals[datainst_idx]:
                try:
                    dt = datetime.strptime(vals[datainst_idx], "%d/%m/%Y").date()
                    garanzia_date = dt + timedelta(days=365)
                    vals[columns_to_insert.index("GARANZIA")] = garanzia_date.strftime("%d/%m/%Y")
                except Exception:
                    pass
        
        if "FINEMF" in columns_to_insert and (not vals[columns_to_insert.index("FINEMF")] or vals[columns_to_insert.index("FINEMF")] == ""):
            datainst_idx = columns_to_insert.index("DATAINST") if "DATAINST" in columns_to_insert else -1
            if datainst_idx >= 0 and vals[datainst_idx]:
                try:
                    dt = datetime.strptime(vals[datainst_idx], "%d/%m/%Y").date()
                    finemf_date = dt + timedelta(days=365*6)  # 6 anni
                    vals[columns_to_insert.index("FINEMF")] = finemf_date.strftime("%d/%m/%Y")
                except Exception:
                    pass
        # MODIFICATO: NON fare controllo duplicati - importare TUTTE le righe
        # I vincoli di database handleranno gli errori di unicità se necessario
        
        # Controllo duplicati: skip se MATRICOLA già presente
        mat_idx = -1
        if "MATRICOLA" in columns_to_insert:
            mat_idx = columns_to_insert.index("MATRICOLA")
        
        if mat_idx >= 0 and mat_idx < len(vals) and vals[mat_idx]:
            mat_val = str(vals[mat_idx]).strip()
            try:
                cur.execute("SELECT 1 FROM clienti WHERE LOWER(TRIM(IFNULL(MATRICOLA,'')))=? LIMIT 1", (mat_val.lower(),))
                if cur.fetchone():
                    skipped += 1
                    skipped_rows.append(vals)
                    continue  # Skip duplicato
            except Exception:
                pass
        
        # Inserisci
        placeholders = ",".join("?" * len(columns_to_insert))
        cols_sql = ", ".join([f'"{c}"' for c in columns_to_insert])
        sql = f"INSERT INTO clienti ({cols_sql}) VALUES ({placeholders})"
        try:
            cur.execute(sql, vals)
            nid = cur.lastrowid
            imported_ids.append(nid)
            inserted += 1
        except Exception as e:
            print(f"Errore inserimento riga: {e}")
            skipped += 1
            skipped_rows.append(vals)
    
    conn.commit()
    
    # Post-processing: aggiungi "RT" a TIPOMF per i record importati
    if imported_ids:
        try:
            ids_placeholder = ",".join("?" * len(imported_ids))
            cur.execute(f"UPDATE clienti SET TIPOMF='RT' WHERE id IN ({ids_placeholder})", imported_ids)
            conn.commit()
        except Exception as e:
            print(f"Errore aggiornamento TIPOMF: {e}")
    
    conn.close()
    
    print(f"[DBF CLIENTI Import] Importati: {inserted}, Saltati: {skipped}")
    return {
        "ok": True,
        "inserted": inserted,
        "skipped": skipped,
        "skipped_rows": skipped_rows,
        "columns": columns_to_insert,
        "message": f"Import DBF completato: {inserted} inseriti, {skipped} saltati"
    }

def _format_import_date(value):
    """
    Formatta le date importate da 8 cifre consecutive (DDMMYYYY) a DD/MM/YYYY.
    Se il valore non è 8 cifre, lo ritorna così com'è.
    Se il valore è vuoto o None, lo ritorna così com'è.
    """
    if not value:
        return value
    s = str(value).strip()
    if not s:
        return s
    # Se sono 8 cifre numeriche, assumi DDMMYYYY e converti
    if re.fullmatch(r'\d{8}', s):
        try:
            day = s[0:2]
            month = s[2:4]
            year = s[4:8]
            # Valida la data (almeno che i numeri siano plausibili)
            int(year); int(month); int(day)  # raises ValueError se non validi
            return f"{day}/{month}/{year}"  # DD/MM/YYYY
        except (ValueError, IndexError):
            return s  # ritorna originale se fallisce
    return s  # ritorna originale se non 8 cifre

def import_csv_intae(path: str, delimiter: str | None = None, encoding: str = "utf-8"):
    """
    Import CSV -> tabella 'intae'.
    Ritorna dict: {ok, inserted, skipped, skipped_rows}
    """
    return _import_intae_csv_generic(path, delimiter, encoding)

def import_misurart_from_csv(path: str, delimiter: str = ";", encoding: str = "utf-8") -> dict:
    """
    Importa MISURART da file CSV con validazione header (permissivo con avviso).
    Se gli header non sono esattamente allineati, mostra avviso ma procede comunque.
    Ritorna dict {ok, inserted, skipped, skipped_rows, columns, error, message, warning}.
    """
    if not os.path.exists(path):
        raise FileNotFoundError(path)
    
    result = {"ok": False, "inserted": 0, "skipped": 0, "skipped_rows": [], "columns": MISURART_COLUMNS, "error": None, "message": "", "warning": None}
    
    try:
        with open(path, "r", encoding=encoding, errors="replace") as f:
            reader = csv.DictReader(f, delimiter=delimiter)
            headers = reader.fieldnames or []
            
            # Validazione header RIGOROSA: verifica che il file abbia la struttura MISURART
            norm_to_db = {normalize_col(c): c for c in MISURART_COLUMNS}
            matched_cols = [h for h in headers if normalize_col(h) in norm_to_db]
            
            # Colonne essenziali MISURART
            essential_misurart = {"TIPO", "CF", "MATRICOLA", "DATAINIZIO", "DATAFINE", "PIVAESERCE", "RAGSOCESE", "COMUNESE", "TIPOINT", "MARCAMF"}
            found_essential = {normalize_col(h) for h in matched_cols} & {normalize_col(c) for c in essential_misurart}
            
            MIN_REQUIRED_COLS = 10
            if len(matched_cols) < MIN_REQUIRED_COLS or len(found_essential) < 7:
                result["ok"] = False
                result["error"] = f"File CSV non valido per MISURART: trovate solo {len(matched_cols)} colonne valide (minimo {MIN_REQUIRED_COLS} richieste).\nColonne trovate: {', '.join(matched_cols[:10])}...\nColonne essenziali mancanti: {len(essential_misurart - found_essential)}"
                result["message"] = result["error"]
                return result
            
            # Procedi con import
            conn = get_connection()
            cur = conn.cursor()
            inserted = 0
            skipped = 0
            duplicates = 0
            skipped_rows = []
            
            # Normalizza mapping
            def normalize_col(col: str) -> str:
                return str(col).strip().upper()
            
            norm_to_db = {normalize_col(c): c for c in MISURART_COLUMNS}
            csv_to_db = {}
            
            for h in headers:
                nh = normalize_col(h)
                if nh in norm_to_db:
                    csv_to_db[h] = norm_to_db[nh]
            
            # Importa tutte le righe (no filtri su MATRICOLA o colonne vuote)
            rows_processed = 0
            for r in reader:
                rows_processed += 1

                # Costruisci mapping colonna->valore per questa riga
                mapped_row = {}
                for k, v in r.items():
                    db_col = csv_to_db.get(k)
                    if db_col:
                        # Formatta le date importate (da DDMMYYYY a DD/MM/YYYY)
                        if db_col in {"DATAINIZIO", "DATAFINE", "DATAVP"}:
                            v = _format_import_date(v)
                        mapped_row[db_col] = v

                # Filtra colonne con valori non vuoti
                cols = [c for c in MISURART_COLUMNS if c in mapped_row and str(mapped_row[c]).strip()]
                
                # Se nessuna colonna ha valore, salta senza contare come skip
                if not cols:
                    continue

                # Controllo duplicati: salta se la riga è ESATTAMENTE identica
                where_parts = []
                where_vals = []
                for c in cols:
                    v = str(mapped_row[c]).strip()
                    if v:
                        where_parts.append(f'"{c}"=?')
                        where_vals.append(v)
                
                is_duplicate = False
                if where_parts:
                    try:
                        cur.execute(f"SELECT 1 FROM misurart WHERE {' AND '.join(where_parts)} LIMIT 1", tuple(where_vals))
                        if cur.fetchone():
                            is_duplicate = True
                            skipped += 1
                            duplicates += 1
                            skipped_rows.append({"row": r, "reason": "Duplicato (riga identica già presente)"})
                    except Exception:
                        pass

                if is_duplicate:
                    continue

                # Inserisci la riga
                vals = [str(mapped_row[c]).strip() for c in cols]
                placeholders = ",".join("?" * len(cols))
                cols_sql = ", ".join([f'"{c}"' for c in cols])
                sql = f"INSERT INTO misurart ({cols_sql}) VALUES ({placeholders})"
                try:
                    cur.execute(sql, vals)
                    inserted += 1
                except Exception as e:
                    skipped += 1
                    skipped_rows.append({"row": r, "reason": str(e)})
            
            conn.commit()
            conn.close()
            
            result["ok"] = True
            result["inserted"] = inserted
            result["skipped"] = skipped
            result["skipped_rows"] = skipped_rows
            result["duplicates"] = duplicates
            result["message"] = f"Import MISURART completato: {inserted} righe importate" if inserted > 0 else "Nessuna riga importata"
    
    except Exception as e:
        result["error"] = str(e)
        result["message"] = str(e)
    
    return result

def import_misurart_from_xlsx(path: str) -> dict:
    """
    Importa MISURART da file XLSX (convertendo prima in CSV con soffice).
    Ritorna dict {ok, inserted, skipped, error, message}.
    """
    if not os.path.exists(path):
        raise FileNotFoundError(path)
    
    result = {"ok": False, "inserted": 0, "skipped": 0, "error": None, "message": ""}
    
    tmpdir = tempfile.mkdtemp()
    try:
        # Converti XLSX a CSV
        cmd = ["soffice", "--headless", "--convert-to", "csv:Text - txt - csv (StarCalc)", path, "--outdir", tmpdir]
        subprocess.run(cmd, check=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        files = [f for f in os.listdir(tmpdir) if f.lower().endswith(".csv")]
        if not files:
            result["error"] = "Conversione XLSX a CSV fallita"
            result["message"] = "Conversione XLSX a CSV fallita"
            return result
        
        csv_path = os.path.join(tmpdir, files[0])
        
        # Importa da CSV
        return import_misurart_from_csv(csv_path, delimiter=";", encoding="latin-1")
    
    except Exception as e:
        result["error"] = str(e)
        result["message"] = str(e)
        return result

def import_misurart_from_dbf(path: str) -> dict:
    """
    Importa MISURART da file DBF con validazione header.
    Ritorna dict {ok, inserted, skipped, error, message}.
    """
    if not os.path.exists(path):
        raise FileNotFoundError(path)
    
    result = {"ok": False, "inserted": 0, "skipped": 0, "error": None, "message": ""}
    
    try:
        # Leggi DBF
        headers, rows = read_dbf_file(path)
        
        # Validazione header RIGOROSA MISURART DBF
        def normalize_col(col: str) -> str:
            return str(col).strip().upper()
        
        norm_to_db = {normalize_col(c): c for c in MISURART_COLUMNS}
        
        # Controlla se gli header attuali sono compatibili
        matched_cols = [h for h in headers if normalize_col(h) in norm_to_db]
        matches = len(matched_cols)
        
        # Se pochi match e ci sono righe, prova con la prima riga come header
        if matches < 10 and rows:
            first_row = [str(v).strip() for v in rows[0]]
            first_row_matched = [v for v in first_row if v and normalize_col(v) in norm_to_db]
            if len(first_row_matched) > matches:
                headers = first_row
                rows = rows[1:]
                matched_cols = first_row_matched
                matches = len(matched_cols)
        
        # Validazione RIGOROSA: minimo 10 colonne valide
        essential_misurart = {"TIPO", "CF", "MATRICOLA", "DATAINIZIO", "DATAFINE", "PIVAESERCE", "RAGSOCESE", "COMUNESE", "TIPOINT", "MARCAMF"}
        found_essential = {normalize_col(h) for h in matched_cols} & {normalize_col(c) for c in essential_misurart}
        
        MIN_REQUIRED_COLS = 10
        if matches < MIN_REQUIRED_COLS or len(found_essential) < 7:
            result["ok"] = False
            result["error"] = f"File DBF non valido per MISURART: trovate solo {matches} colonne valide (minimo {MIN_REQUIRED_COLS} richieste).\nColonne trovate: {', '.join(matched_cols[:10])}...\nColonne essenziali mancanti: {len(essential_misurart - found_essential)}"
            result["message"] = result["error"]
            return result
        
        # Procedi con import
        conn = get_connection()
        cur = conn.cursor()
        inserted = 0
        skipped = 0
        duplicates = 0
        
        # Normalizza mapping (token-based, permissivo)
        def normalize_col(col: str) -> str:
            return str(col).strip().upper()

        def tokens(s: str):
            return [t for t in re.split(r"[^A-Z0-9]+", normalize_col(s)) if t]

        dbf_to_db = {}
        used_targets = set()
        for h in headers:
            best_target = None
            best_score = 0
            h_tokens = set(tokens(h))
            if not h_tokens:
                continue
            for target in MISURART_COLUMNS:
                t_tokens = set(tokens(target))
                score = len(h_tokens & t_tokens)
                if score > best_score and target not in used_targets:
                    best_score = score
                    best_target = target
            if best_score > 0 and best_target:
                dbf_to_db[h] = best_target
                used_targets.add(best_target)
        
        for row_data in rows:
            # Crea dict da headers
            row_dict = {headers[i]: row_data[i] for i in range(len(headers))}
            
            # Estrai valori mappati
            mapped_row = {}
            # Formatta le date importate (da DDMMYYYY a DD/MM/YYYY)
            for h in headers:
                if h in row_dict:
                    if h in {"DATAINIZIO", "DATAFINE", "DATAVP"}:
                        row_dict[h] = _format_import_date(row_dict[h])
            for dbf_col, db_col in dbf_to_db.items():
                val = str(row_dict.get(dbf_col, "")).strip()
                if val:
                    mapped_row[db_col] = val
            
            # Importa anche righe con poche colonne (basta 1)
            cols = [c for c in MISURART_COLUMNS if c in mapped_row and mapped_row[c]]
            if not cols:
                continue  # Salta solo righe completamente vuote

            # Controllo duplicati: salta se la riga è ESATTAMENTE identica
            where_parts = []
            where_vals = []
            for c in cols:
                v = str(mapped_row[c]).strip()
                if v:
                    where_parts.append(f'"{c}"=?')
                    where_vals.append(v)
            
            is_duplicate = False
            if where_parts:
                try:
                    cur.execute(f"SELECT 1 FROM misurart WHERE {' AND '.join(where_parts)} LIMIT 1", tuple(where_vals))
                    if cur.fetchone():
                        is_duplicate = True
                        skipped += 1
                        duplicates += 1
                except Exception:
                    pass

            if is_duplicate:
                continue

            # Inserisci la riga non duplicata
            vals = [mapped_row[c] for c in cols]
            placeholders = ",".join("?" * len(cols))
            cols_sql = ", ".join([f'"{c}"' for c in cols])
            sql = f"INSERT INTO misurart ({cols_sql}) VALUES ({placeholders})"
            
            try:
                cur.execute(sql, vals)
                inserted += 1
            except Exception as e:
                skipped += 1
        
        conn.commit()
        conn.close()
        
        result["ok"] = True
        result["inserted"] = inserted
        result["skipped"] = skipped
        result["duplicates"] = duplicates
        result["message"] = f"Import MISURART completato: {inserted} importati, {skipped} saltati"
    
    except Exception as e:
        result["error"] = str(e)
        result["message"] = str(e)
    
    return result

def import_intae_from_dbf(path: str) -> dict:
    """
    Importa INTAE da file DBF con validazione struttura.
    Se i campi DBF sono generici (N1, N2, ...), legge il primo record come intestazioni.
    Ritorna dict {ok, inserted, error, message}.
    """
    if not os.path.exists(path):
        raise FileNotFoundError(path)

    result = {"ok": False, "inserted": 0, "error": None, "message": ""}

    try:
        from dbfread import DBF
        dbf = DBF(path)
        
        # Estrai i nomi dei campi DBF (di solito N1, N2, ..., oppure i veri nomi)
        dbf_field_names = [f.name for f in dbf.fields]
        
        # Controlla se i campi sono generici (N1, N2, N3, ...)
        is_generic = all(f.name.startswith('N') and f.name[1:].isdigit() for f in dbf.fields)
        
        if not is_generic:
            # I nomi dei campi sono già i veri nomi, usa il mapping diretto
            headers = dbf_field_names
            rows_data = [list(record.values()) for record in dbf]
        else:
            # I campi sono generici: il primo record contiene le intestazioni
            records = list(dbf)
            if not records:
                result["message"] = "File DBF vuoto"
                return result
            
            # Estrai le intestazioni dal primo record
            first_record = records[0]
            headers = [str(first_record[name]).strip() for name in dbf_field_names]
            
            # I dati sono dal secondo record in poi
            rows_data = [list(record.values()) for record in records[1:]]
        
        # Mapping header -> colonne INTAE (robusto)
        header_to_intae = _build_intae_header_map(headers)

        # Se mappatura scarsa, prova a usare la prima riga come header
        if len(header_to_intae) < 2 and rows_data:
            candidate_headers = [str(v or "").strip() for v in rows_data[0]]
            candidate_map = _build_intae_header_map(candidate_headers)
            if len(candidate_map) > len(header_to_intae):
                headers = candidate_headers
                rows_data = rows_data[1:]
                header_to_intae = candidate_map

        # Validazione RIGOROSA: minimo 8 colonne INTAE valide su 12 totali (67%)
        MIN_REQUIRED_INTAE_COLUMNS = 8
        essential_cols = {"CODICE_INTERVENTO", "PIVA", "DATA_INTERVENTO", "CF_TECNICO", "STATO_CF", "DISPOSITIVO", "MATRICOLA", "TIPOLOGIA"}
        found_cols = set(header_to_intae.values()) if header_to_intae else set()
        essential_found = found_cols & essential_cols
        
        if len(header_to_intae) < MIN_REQUIRED_INTAE_COLUMNS or len(essential_found) < 6:
            result["error"] = f"File DBF non valido per INTAE: trovate solo {len(found_cols)} colonne INTAE valide (minimo {MIN_REQUIRED_INTAE_COLUMNS} richieste).\nColonne trovate: {', '.join(sorted(found_cols)) if found_cols else 'nessuna'}\nColonne essenziali mancanti: {', '.join(sorted(essential_cols - essential_found))}"
            result["message"] = result["error"]
            return result
        
        # Importa i record
        conn = get_connection()
        cur = conn.cursor()
        inserted = 0
        skipped = 0
        
        for row_data in rows_data:
            row_dict = {headers[i]: row_data[i] if i < len(row_data) else "" for i in range(len(headers))}
            rec = _intae_build_record(row_dict, header_to_intae)

            # NESSUN FILTRO - IMPORTA TUTTE LE RIGHE
            cols = [c for c in INTAE_COLUMNS if rec.get(c) is not None]
            vals = [rec.get(c, "") for c in cols]
            placeholders = ",".join("?" * len(cols))
            cols_sql = ", ".join([f'"{c}"' for c in cols])
            sql = f"INSERT INTO intae ({cols_sql}) VALUES ({placeholders})"

            try:
                cur.execute(sql, vals)
                inserted += 1
            except Exception as e:
                print(f"[INTAE DBF Import ERROR] {e}")
                skipped += 1
        
        conn.commit()
        conn.close()
        
        result["ok"] = True
        result["inserted"] = inserted
        result["skipped"] = skipped
        if skipped > 0:
            result["message"] = f"Import INTAE completato: {inserted} record inseriti, {skipped} errori."
        else:
            result["message"] = f"Import INTAE completato: {inserted} record inseriti."
        
    except Exception as e:
        result["error"] = str(e)
        result["message"] = str(e)
    
    return result