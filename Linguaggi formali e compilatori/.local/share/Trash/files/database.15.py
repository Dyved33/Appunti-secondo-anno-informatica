import sqlite3
import os
import csv
import tempfile
import subprocess
import calendar
from datetime import datetime, timedelta, date
from dateutil import parser
from typing import List, Tuple, Optional, Dict, Any
import re
import unicodedata
import traceback
import sys
from io import BytesIO
from PIL import Image
from reportlab.pdfbase.pdfmetrics import stringWidth
from math import isfinite

try:
    from reportlab.pdfgen import canvas
    from reportlab.lib.pagesizes import A4
    from reportlab.lib.units import mm
except Exception:
    canvas = None

try:
    from PyPDF2 import PdfReader, PdfWriter
except Exception:
    PdfReader = PdfWriter = None

try:
    from barcode import Code128
    from barcode.writer import ImageWriter
except Exception:
    Code128 = None


DB_FILE = "database.db"

# Colonne ESATTAMENTE come richieste (nomi preservati)
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

# Colonne per MISURART (normalizzate, tutte TEXT tranne id)
MISURART_COLUMNS = [
    "TIPO","CF","ALFASIGILL","NUMSIGILLO","DATAINIZIO","DATAFINE","ESITO",
    "PIVAESERCE","RAGSOCESE","COMUNESE","PROVINCIA","INDIRIZZO","CAPESE",
    "TIPOINT","MARCAMF","MODELLO","LOGOTIPO","MATRICOLA",
    "F1","F2","F3","F4","F5","F6","CARCTRL","DATAVP"
]

DATE_COLS = {"FINEMF", "GARANZIA", "DATAINST", "SCAD_ASS", "DATA", "DATADIS", "DATA_INTERVENTO", "DATAINIZIO", "DATAFINE"}

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
    Ritorna dict: { 'table_name': rows_updated, ... }
    Solleva ValueError se matricola vuota o data non valida.
    """
    if not matricola or not str(matricola).strip():
        raise ValueError("Devi fornire una matricola non vuota.")

    mnorm = _normalize_matricola_for_compare(matricola)
    if not mnorm:
        raise ValueError("Matricola non valida dopo normalizzazione.")

    # normalizza data
    iso_date = _parse_date_param_to_iso(applicazione_date)

    conn = get_connection()
    cur = conn.cursor()
    updated_counts = {}

    try:
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

            # procediamo solo se la tabella ha MATRICOLA e almeno una delle colonne DATADIST o VF
            if "MATRICOLA" not in col_names:
                continue
            has_datadist = "DATADIST" in col_names
            has_vf = "VF" in col_names

            if not has_datadist and not has_vf:
                # niente da aggiornare in questa tabella
                continue

            # costruiamo SET dinamico
            set_parts = []
            params = []

            if has_datadist:
                set_parts.append("DATADIST = ?")
                params.append(iso_date)

            if has_vf:
                # decide se VF è numerico (INTEGER/REAL) -> 0 altrimenti stringa vuota
                vf_type = col_types.get("VF", "")
                if any(t in vf_type for t in ("INT", "REAL", "NUM", "DECIMAL", "FLOAT", "DOUBLE")):
                    vf_zero = 0
                else:
                    vf_zero = ""
                set_parts.append("VF = ?")
                params.append(vf_zero)

            set_clause = ", ".join(set_parts)
            # where: confronto matricola case-insensitive (COLLATE NOCASE se presente)
            # usiamo la colonna raw MATRICOLA nella query ma parametri normalizzati
            # per compatibilità faremo WHERE lower(matricola) = ?
            where_clause = "lower(matricola) = ?"
            params.append(str(matricola).strip().lower())

            sql = f"UPDATE {tbl} SET {set_clause} WHERE {where_clause}"
            try:
                cur.execute(sql, params)
                cnt = cur.rowcount if hasattr(cur, "rowcount") else cur.execute("SELECT changes()").fetchone()[0]
                # sqlite cursors rowcount may be -1; meglio usare changes()
                # use total changes for this connection: compute difference
                # we'll fetch changes() AFTER executing all updates per table:
                # workaround: read changes() now:
                cur2 = conn.cursor()
                cur2.execute("SELECT changes()")
                cnt = cur2.fetchone()[0]
                # But changes() returns cumulative for connection; to be precise we'd have to track before/after.
                # Simpler: use cur.rowcount where available; else estimate via SELECT COUNT with same WHERE BEFORE update.
            except Exception:
                # se UPDATE fallisce (es. permessi o tipi) proviamo alternativa robusta:
                try:
                    # conteggio righe da aggiornare
                    sel_sql = f"SELECT COUNT(*) FROM {tbl} WHERE lower(matricola) = ?"
                    cur.execute(sel_sql, (mnorm,))
                    cnt = cur.fetchone()[0] or 0
                    # proviamo comunque ad eseguire l'UPDATE con valori testuali (fallback)
                    cur.execute(sql, params)
                except Exception:
                    cnt = 0

            # store count (num aggiornate rilevato)
            updated_counts[tbl] = cnt

        # ... dentro deactivate_machine, dopo aver determinato iso_date e prima del commit finale:
        try:
            # usa la stessa connessione per consistenza transazionale
            cnt_clienti = update_clienti_datadist_vf(matricola, iso_date, conn=conn)
            # sovrascrivi o integra il conteggio nella mappa di risultati
            updated_counts['clienti'] = cnt_clienti
        except Exception:
            # non vogliamo bloccare l'intera disattivazione per un errore sul solo clienti
            updated_counts['clienti'] = updated_counts.get('clienti', 0)
        
        # commit finale
        conn.commit()
        return updated_counts

    finally:
        try:
            conn.close()
        except Exception:
            pass


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

def import_clients_from_csv_generic(path: str, delimiter: str | None = None, encoding: str = "latin-1", batch_commit: int = 500) -> tuple[int,int]:
    """
    Import CSV in modo robusto:
      - rileva le colonne reali della tabella clienti usando PRAGMA
      - costruisce INSERT dinamico solo sulle colonne trovate
      - calcola GARANZIA e FINEMF prima dell'inserimento quando mancanti
      - ritorna (inserted, skipped)
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
        # Verifica che il file contenga le colonne fondamentali per un cliente
        required_columns = ["DATAINST", "MATRICOLA", "RAG_SOC"]
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

#se voglio fare validazione obbligatoria, decommento questo blocco
#######################################################################################################
            # === VALIDAZIONE OBBLIGATORI PER IMPORT CLIENTI ===
            # lista campi obbligatori (stessa lista che usi per l'inserimento manuale)
            #required_fields = [
            #    "RAG_SOC","P_IVA","IND_VIA","COMUNE","PROV","CAP","MF","MARCA","LOGOTIPO",
            #    "MATRICOLA","SCAD_ASS","IMP_ASS","VF"
            #]
#
            ## usa _is_empty già definita nella funzione per decidere se un campo è "vuoto"
            #missing_required = [f for f in required_fields if _is_empty(val_for(f))]
            #if missing_required:
            #    skipped += 1
            #    # opzionale: log diagnostico
            #    # print(f"[IMPORT] riga {rows_processed} saltata: mancanti campi {missing_required}")
            #    continue
#######################################################################################################
            # Forza TIPOMF = "RT" per ogni riga importata (sovrascrive eventuale valore del CSV)
            if "TIPOMF" in columns_to_insert:
                try:
                    tip_idx = columns_to_insert.index("TIPOMF")
                    vals[tip_idx] = "RT"
                except Exception:
                    pass


            # controllo duplicato: serve avere rag, piva e matricola tutti presenti
            rag = val_for("RAG_SOC")
            piva = val_for("P_IVA")
            mat = val_for("MATRICOLA")

            # cliente_esiste(conn, rag, piva, mat) deve essere definita in questo module
            try:
                if cliente_esiste(conn, rag, piva, mat):
                    skipped += 1
                    # print(f"[IMPORT] duplicato saltato riga {rows_processed}: {rag} / {piva} / {mat}")
                    continue
            except Exception:
                # se il controllo duplicato fallisce, non blocchiamo l'import (log e procedi)
                print(f"[IMPORT WARNING] controllo duplicati fallito riga {rows_processed}")

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
        print(f"[IMPORT] file={path!r} delimiter={used_delim!r} headers={headers} rows_processed={rows_processed} inserted={inserted} skipped={skipped}")
        return inserted, skipped

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
    for fmt in ("%d/%m/%Y", "%d-%m-%Y", "%Y-%m-%d", "%Y/%m/%d", "%Y%m%d"):
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

def generate_barcode_image(matricola: str, out_path: str) -> str:
    """Genera PNG Code128 e ritorna il percorso del file .png"""
    matricola = str(matricola or "").strip()
    if matricola == "":
        raise ValueError("Matricola vuota per barcode")
    base, ext = os.path.splitext(out_path)
    if not ext.lower().endswith(".png"):
        out_path = base + ".png"
    writer = ImageWriter()
    code = Code128(matricola, writer=writer)
    saved = code.save(base)  # ritorna base + '.png'
    return saved

def _fill_pdf_overlay(template_size: tuple, context: dict, barcode_path: str, field_positions: dict) -> BytesIO:
    """Riempie un overlay PDF con testo e barcode (versione semplificata)."""
    width_pt, height_pt = template_size
    buf = BytesIO()
    c = canvas.Canvas(buf, pagesize=(width_pt, height_pt))
    
    BARCODE_PADDING = mm_to_pt(1.5)
    
    for key, pos in field_positions.items():
        if key == "BARCODE":
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
    
    # Aggiungi barcode se presente
    if "BARCODE" in field_positions and barcode_path and os.path.exists(barcode_path):
        pos = field_positions["BARCODE"]
        if len(pos) == 4:
            x, y, box_w, box_h = pos
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
                
                c.drawImage(barcode_path, x_img, y_img, width=draw_w, height=draw_h, 
                           preserveAspectRatio=True, mask='auto')
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
            d = _parse_vf_date(vf_val)  # Usa il parser specifico per VF
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
        
        try:
            generate_barcode_image(matric, barcode_path)
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
            "MATRICOLA": matric,
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
                for fmt in ("%d/%m/%Y", "%d.%m.%Y", "%Y-%m-%d", "%Y/%m/%d", "%Y%m%d"):
                    try:
                        parsed_date = datetime.strptime(s, fmt).date()
                        break
                    except Exception:
                        parsed_date = None

        if not parsed_date:
            if len(unparsed) < 20:
                unparsed.append((rag, finemf_raw))
            continue

        if today <= parsed_date <= end_date:
            display_date = parsed_date.strftime("%d/%m/%Y")
            out.append((rag or "", display_date, mf or "", tel or "", codice or ""))

    if unparsed:
        print("[report_memorie_in_scadenza] alcune FINEMF non parsate (esempi):", unparsed[:10])

    conn.close()
    print(f"[report_memorie_in_scadenza] fallback, trovati {len(out)}")
    return out

def report_rt_for_clienti(month: int | None = None, year: int | None = None):
    """
    Versione aggiornata:
    - Legge clienti dalla tabella e filtra per VF nel mese/anno specificato
    - Se month/year sono None, restituisce tutti i clienti
    - Ritorna tuple: (RAG_SOC, IND_VIA, COMUNE, MATRICOLA, DATA, TIPOMF, DATAINST, AZZ, P_IVA, MARCA)
    """
    conn = get_connection()
    cur = conn.cursor()

    # 1) leggi una riga per ottenere i nomi delle colonne
    try:
        cur.execute("SELECT * FROM clienti LIMIT 1")
        cols = [d[0] for d in cur.description] if cur.description else []
    except Exception:
        print("[DEBUG][DB] Impossibile leggere schema clienti:", traceback.format_exc())
        cols = []

    print("[DEBUG][DB] colonne tabella clienti:", cols)

    # 2) trova la colonna VF (tu già lo facevi)
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
        print("[DEBUG][DB] Colonna 'VF' non trovata (heuristics). Nomi colonne:", cols)
        # fallback: comportamento precedente - restituisci prime 8 colonne (compatibilità)
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
            rows = cur.fetchall()
            conn.close()
            return [tuple('' if v is None else str(v) for v in row) for row in rows]
        except Exception:
            conn.close()
            print("[ERROR][DB] Errore fallback SELECT:", traceback.format_exc())
            return []

    print(f"[DEBUG][DB] Uso colonna VF: '{finemf_col}'  PIVA col: '{piva_col}'  id_col: '{id_col}' matricola_col: '{matricola_col}'")

    # 3) costruisci SELECT includendo la colonna VF e (se esiste) P_IVA
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

    print(f"[DEBUG][DB] righe lette: {len(all_rows)}")

    # 4) itera e filtra; per i match recupera P_IVA (già incluso) e MARCA (da tabella macchine)
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

        if sample_shown < SAMPLE_LIMIT:
            print(f"[DEBUG][DB] riga {idx} VF grezzo: '{finemf_raw}'  PIVA estratta: '{piva_val}'")
            sample_shown += 1

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
            # aggiungi alla lista (prime 8 valori + P_IVA + MARCA)
            matches.append(tuple(list(visible) + [piva_val, marca_val]))
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

                    matches.append(tuple(list(visible) + [piva_val, marca_val]))
                else:
                    # non coincide: ignora
                    pass
            except Exception:
                pass
        else:
            parsed_failures += 1

    print(f"[DEBUG][DB] parsed_success={parsed_success}, parsed_failures={parsed_failures}, matches={len(matches)}")
    for i, m in enumerate(matches[:10]):
        print(f"[DEBUG][DB] esempio match {i}: {m}")

    conn.close()
    return matches


# all'inizio del file, dopo gli import e le definizioni di connessione DB
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

# Posizioni di default (adatta ai tuoi template). Coordinate in punti (origine bottom-left).
FIELD_POSITIONS_DEFAULT = {
    "RAG_SOC": (mm_to_pt(85), mm_to_pt(250), 10, mm_to_pt(180)),
    "P_IVA": (mm_to_pt(85), mm_to_pt(235), 10, mm_to_pt(180)),     
    "INDIRIZZO": (mm_to_pt(85), mm_to_pt(225), 10, mm_to_pt(180)),
    "CAP_LOCALITA": (mm_to_pt(85), mm_to_pt(216), 10, mm_to_pt(180)),     
    "MARCA": (mm_to_pt(85), mm_to_pt(210), 10, mm_to_pt(180)),       
    "MF": (mm_to_pt(85), mm_to_pt(205), 10, mm_to_pt(180)),          
    "MATRICOLA": (mm_to_pt(85), mm_to_pt(200), 10, mm_to_pt(180)),
    
    "SCAD_ASS_LBL": (mm_to_pt(21), mm_to_pt(186), 8, mm_to_pt(30)),     # Etichetta a sinistra
    "SCAD_ASS": (mm_to_pt(37), mm_to_pt(186), 8, mm_to_pt(50)),         # Dato a destra
    
    "INST_LBL": (mm_to_pt(57), mm_to_pt(186), 8, mm_to_pt(30)),         # Etichetta a sinistra
    "DATAINST": (mm_to_pt(65), mm_to_pt(186), 8, mm_to_pt(50)),         # Dato a destra
    
    "GAR_LBL": (mm_to_pt(85), mm_to_pt(186), 8, mm_to_pt(30)),         # Etichetta a sinistra
    "GARANZIA": (mm_to_pt(94), mm_to_pt(186), 8, mm_to_pt(50)),        # Dato a destra
    
    "BARCODE": (mm_to_pt(140), mm_to_pt(275), mm_to_pt(50), mm_to_pt(20)),
}

reader = PdfReader("template.pdf")
page = reader.pages[0]
height_pt = float(page.mediabox.height)

# --- funzione di conversione robusta: non eseguirla al momento dell'import ---
def convert_bottom_to_top(field_positions: dict, page_height_pt: float) -> dict:
    """
    Converti un dizionario con y misurate dal bottom in y_from_top e
    ritorna il nuovo dict in formato esteso (x, y_from_top, box_w, box_h, font) per i campi testuali
    e (x, y_from_top, box_w, box_h) per il BARCODE.
    Accetta vari formati in input:
      - text old: (x, y_from_bottom, font, max_w)
      - text new: (x, y_from_bottom, box_w, box_h, font)
      - barcode: (x, y_from_bottom, box_w, box_h)
    """
    out = {}
    for k, v in field_positions.items():
        if k == "BARCODE":
            if len(v) == 4:
                x, y_bottom, box_w, box_h = v
                y_from_top = page_height_pt - y_bottom
                out[k] = (x, y_from_top, box_w, box_h)
            else:
                raise ValueError(f"BARCODE atteso 4 tuple, trovato {len(v)} per {k}: {v}")
            continue

        # non-BARCODE: può essere 4-tuple (x, y, font, maxw) o 5-tuple (x, y, box_w, box_h, font)
        if len(v) == 4:
            # old format: (x, y_bottom, font, max_w)
            x, y_bottom, font, max_w = v
            y_from_top = page_height_pt - y_bottom
            # stima box_h = font + padding (6pt). Puoi regolare dopo.
            box_h = font + 6
            box_w = max_w
            out[k] = (x, y_from_top, box_w, box_h, font)
        elif len(v) == 5:
            # new format already: (x, y_bottom, box_w, box_h, font)
            x, y_bottom, box_w, box_h, font = v
            y_from_top = page_height_pt - y_bottom
            out[k] = (x, y_from_top, box_w, box_h, font)
        else:
            raise ValueError(f"Formato posizione non riconosciuto per {k}: {v}")
    return out
# --- end convert_bottom_to_top ---

#FIELD_POSITIONS_DEFAULT_top = convert_bottom_to_top(FIELD_POSITIONS_DEFAULT)
#print(FIELD_POSITIONS_DEFAULT_top)

# -------------------------
# Helper functions per GUI
# -------------------------

def _parse_vf_date(s):
    """
    Parser specifico per la colonna VF.
    Assume il formato DD/MM/YYYY (giorno/mese/anno) per evitare ambiguità.
    Se il mese è > 12, scambia giorno e mese (auto-correzione).
    Ritorna datetime o None.
    """
    if s is None or not str(s).strip():
        return None
    
    s_str = str(s).strip()
    
    # Prova innanzitutto con il formato DD/MM/YYYY
    try:
        parts = s_str.split("/")
        if len(parts) == 3:
            dd, mm, yyyy = int(parts[0]), int(parts[1]), int(parts[2])
            
            # Auto-correzione: se mese > 12, scambia giorno e mese
            if mm > 12 and dd <= 12:
                dd, mm = mm, dd
            
            return datetime(yyyy, mm, dd)
    except (ValueError, IndexError):
        pass
    
    # Fallback: usa il parser generico
    return _parse_date(s)

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
    for fmt in ("%d/%m/%Y", "%d.%m.%Y", "%Y-%m-%d", "%Y/%m/%d", "%d-%m-%Y", "%d%m%Y"):
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

def import_csv_intae(path: str, delimiter: str = ";", encoding: str = "utf-8") -> dict:
    """Importa CSV in tabella INTAE. Ritorna dict con chiavi 'inserted' e 'skipped'."""
    if not os.path.exists(path):
        raise FileNotFoundError(path)
    
    # Mapping dei nomi colonne CSV ai nomi INTAE_COLUMNS
    csv_to_intae_mapping = {
        "Codice intervento": "CODICE_INTERVENTO",
        "PIVA Laboratorio": "PIVA",  # Estrae solo PIVA, ignora Laboratorio
        "Data intervento": "DATA_INTERVENTO",
        "CF Tecnico": "CF_TECNICO",
        "Stato CF": "STATO_CF",
        "Dispositivo": "DISPOSITIVO",
        "Matricola": "MATRICOLA",
        "Tipologia": "TIPOLOGIA",
        "Data Trasmissione": "DATA_TRASMISSIONE",
        "Nota": "NOTA",
        "Versione SW": "VERSIONE_SW",
        "Data release": "DATA_RELEASE",
        "LABORATORIO": "LABORATORIO",  # se presente come colonna separata
    }
    
    conn = get_connection()
    cur = conn.cursor()
    inserted = 0
    skipped = 0
    
    try:
        with open(path, "r", encoding=encoding, errors="replace") as f:
            reader = csv.DictReader(f, delimiter=delimiter)
            for row in reader:
                # Mappa le colonne CSV alle colonne INTAE
                mapped_row = {}
                for csv_col, intae_col in csv_to_intae_mapping.items():
                    if csv_col in row:
                        val = str(row[csv_col]).strip() if row[csv_col] else ""
                        if val:
                            mapped_row[intae_col] = val
                
                # Se non trovato tramite mapping, prova normalizzazione diretta
                for csv_col, val in row.items():
                    if not val or not str(val).strip():
                        continue
                    normalized_key = csv_col.upper().replace(" ", "_")
                    if normalized_key in INTAE_COLUMNS and normalized_key not in mapped_row:
                        mapped_row[normalized_key] = str(val).strip()
                
                # Filtra: almeno CODICE_INTERVENTO deve essere presente
                if "CODICE_INTERVENTO" not in mapped_row or not mapped_row["CODICE_INTERVENTO"]:
                    skipped += 1
                    continue
                
                # inserisci solo colonne che esistono in INTAE_COLUMNS e hanno valori
                cols = [c for c in INTAE_COLUMNS if c in mapped_row and mapped_row[c]]
                if not cols:
                    skipped += 1
                    continue
                
                # CONTROLLO DUPLICATI: P_IVA + MATRICOLA + DATA_INTERVENTO
                piva = mapped_row.get("PIVA", "").strip()
                matricola = mapped_row.get("MATRICOLA", "").strip()
                data_intervento = mapped_row.get("DATA_INTERVENTO", "").strip()
                
                if piva and matricola and data_intervento:
                    # Controlla se esiste già
                    cur.execute(
                        "SELECT COUNT(*) FROM intae WHERE lower(PIVA)=? AND lower(MATRICOLA)=? AND DATA_INTERVENTO=?",
                        (piva.lower(), matricola.lower(), data_intervento)
                    )
                    if int(cur.fetchone()[0] or 0) > 0:
                        skipped += 1
                        continue  # Salta riga duplicata
                
                vals = [mapped_row[c] for c in cols]
                placeholders = ",".join("?" * len(cols))
                cols_sql = ", ".join([f'"{c}"' for c in cols])
                sql = f"INSERT INTO intae ({cols_sql}) VALUES ({placeholders})"
                cur.execute(sql, vals)
                inserted += 1
        conn.commit()
    finally:
        conn.close()
    
    # Log del risultato
    if skipped > 0:
        print(f"[INTAE Import] Importati: {inserted}, Saltati (duplicati o righe incomplete): {skipped}")
    
    return {"inserted": inserted, "skipped": skipped}

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
    return rows

def misurart_add_from_map(data: dict) -> int:
    """Inserisce un record MISURART da dict. Ritorna id oppure None."""
    if not data:
        return None
    
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
    mat_input = str(matricola).strip().lower()
    
    # Prova match esatto
    cur.execute("SELECT * FROM misurart WHERE lower(MATRICOLA) = ?", (mat_input,))
    rows = cur.fetchall()
    
    # Se no match esatto, prova LIKE
    if not rows:
        cur.execute("SELECT * FROM misurart WHERE lower(MATRICOLA) LIKE ?", (f"%{mat_input}%",))
        rows = cur.fetchall()
    
    conn.close()
    return rows or []

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

def report_cambi_dgfe(data_from: str, data_to: str) -> List[Tuple]:
    """
    Report interventi INTAE tra due date.
    Ritorna liste di tuple (cliente, matricola, count_interventi).
    """
    conn = get_connection()
    cur = conn.cursor()
    
    try:
        d1 = _parse_date_param_to_iso(data_from) if data_from else "1900-01-01"
        d2 = _parse_date_param_to_iso(data_to) if data_to else date.today().isoformat()
    except Exception:
        d1, d2 = "1900-01-01", date.today().isoformat()
    
    try:
        sql = """
            SELECT c.RAG_SOC, i.MATRICOLA, COUNT(*)
            FROM intae i
            LEFT JOIN clienti c ON lower(trim(c.MATRICOLA)) = lower(trim(i.MATRICOLA))
            WHERE i.DATA_INTERVENTO BETWEEN ? AND ?
            GROUP BY i.MATRICOLA
            ORDER BY COUNT(*) DESC
        """
        cur.execute(sql, (d1, d2))
        rows = cur.fetchall()
        conn.close()
        # Converte COUNT(*) (terzo elemento) da stringa a int
        result = [(r[0], r[1], int(r[2] or 0)) for r in rows]
        return result or []
    except Exception:
        conn.close()
        return []

def get_verification_data(matricola: str) -> dict:
    """
    Recupera i dati completi per la verifica di una macchina.
    Ricerca in CLIENTI e INTAE. Restituisce dict con dati cliente, dati INTAE e informazioni sulla verifica.
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
        mat_like = f"%{mat_normalized}%"
        
        # Recupera dati cliente - Cerca prima con match esatto, poi con LIKE
        cur.execute("SELECT * FROM clienti WHERE lower(trim(MATRICOLA)) = ? LIMIT 1", (mat_normalized,))
        cols_client = [r[1] for r in cur.execute("PRAGMA table_info(clienti)").fetchall()]
        client_row = cur.fetchone()
        
        # Se non trova corrispondenza esatta, prova con LIKE
        if not client_row:
            cur.execute("SELECT * FROM clienti WHERE lower(MATRICOLA) LIKE ? LIMIT 1", (mat_like,))
            client_row = cur.fetchone()
        
        if client_row:
            result["cliente"] = {cols_client[i]: client_row[i] for i in range(len(cols_client))}
        else:
            # Prova a cercare in INTAE se non trovato in CLIENTI
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
        cur.execute("""
            SELECT CODICE_INTERVENTO, DATA_INTERVENTO, CF_TECNICO, STATO_CF, 
                   DISPOSITIVO, TIPOLOGIA, DATA_TRASMISSIONE, NOTA, VERSIONE_SW, DATA_RELEASE
            FROM intae 
            WHERE lower(trim(MATRICOLA)) = ? 
            ORDER BY DATA_INTERVENTO DESC 
            LIMIT 1
        """, (mat_normalized,))
        
        intae_row = cur.fetchone()
        
        # Se non trovato con match esatto, prova con LIKE
        if not intae_row:
            cur.execute("""
                SELECT CODICE_INTERVENTO, DATA_INTERVENTO, CF_TECNICO, STATO_CF, 
                       DISPOSITIVO, TIPOLOGIA, DATA_TRASMISSIONE, NOTA, VERSIONE_SW, DATA_RELEASE
                FROM intae 
                WHERE lower(MATRICOLA) LIKE ?
                ORDER BY DATA_INTERVENTO DESC 
                LIMIT 1
            """, (mat_like,))
            intae_row = cur.fetchone()
        
        if intae_row:
            intae_cols = ["CODICE_INTERVENTO", "DATA_INTERVENTO", "CF_TECNICO", "STATO_CF",
                         "DISPOSITIVO", "TIPOLOGIA", "DATA_TRASMISSIONE", "NOTA", "VERSIONE_SW", "DATA_RELEASE"]
            result["intae_data"] = {intae_cols[i]: intae_row[i] for i in range(len(intae_cols))}
            result["verificatore"] = intae_row[2] or ""  # CF_TECNICO
            result["data_intervento"] = intae_row[1] or ""  # DATA_INTERVENTO
            result["codice_intervento"] = intae_row[0] or ""  # CODICE_INTERVENTO
        
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