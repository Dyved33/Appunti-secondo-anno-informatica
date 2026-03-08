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
    "CODICE_INTERVENTO", "PIVA", "LABORATORIO", "DATA_INTERVENTO",
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
    return tuple("" if v is None else v for v in row)

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
        cnt_simple = cur.fetchone()[0] or 0

        rows_to_update_ids = []
        if cnt_simple > 0:
            # aggiorneremo con WHERE lower(matricola) = ?
            where_clause = "lower(matricola) = ?"
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
            # numero aggiornate: esegui SELECT COUNT per coerenza
            cur.execute("SELECT changes()")
            updated = cur.fetchone()[0] if cur.fetchone() else cur.rowcount
            # SQLite `changes()` can be fetched differently; safer: query matched rows count_simple
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
                placeholders = ",".join(["?"] * len(chunk))
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
    import csv, os, sqlite3

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
def normalize_fieldname(h: str) -> str:
    """Rimuove metadati dagli header tipo 'TIPOMF,C,3' -> 'TIPOMF'"""
    if h is None:
        return h
    return str(h).split(',')[0].strip()

def _normalize_date_for_db(raw) -> str | None:
    """
    Prova a parsare raw (string/date/datetime/tuple/list/dict) e ritorna ISO 'YYYY-MM-DD' oppure None.
    Usa _parse_date se possibile; accetta inoltre (day,month,year) come list/tuple o dict {'day':..,'month':..,'year':..}.
    """
    if raw is None:
        return None
    # se è già date/datetime
    if isinstance(raw, date) and not isinstance(raw, datetime):
        return raw.isoformat()
    if isinstance(raw, datetime):
        return raw.date().isoformat()

    # supporto per tuple/list come (day, month, year)
    if isinstance(raw, (list, tuple)) and len(raw) >= 3:
        dt = _parse_dmy_components(raw[0], raw[1], raw[2])
        if dt:
            return dt.isoformat()

    # supporto per dict con chiavi day/month/year o d/m/y
    if isinstance(raw, dict):
        keys = {k.lower(): k for k in raw.keys()}
        for trip in (("day","month","year"), ("d","m","y"), ("giorno","mese","anno")):
            if all(k in keys for k in trip):
                dt = _parse_dmy_components(raw[keys[trip[0]]], raw[keys[trip[1]]], raw[keys[trip[2]]])
                if dt:
                    return dt.isoformat()
    # se è stringa, lascia la logica precedente (ISO o _parse_date o vari formati)
    try:
        dt = _parse_date(raw)
        if dt:
            return dt.isoformat()
    except Exception:
        pass

    s = str(raw).strip()
    if s == "":
        return None
    for fmt in ("%d/%m/%Y", "%d-%m-%Y", "%Y-%m-%d", "%Y/%m/%d", "%Y%m%d"):
        try:
            from datetime import datetime as _dt
            return _dt.strptime(s.split()[0], fmt).date().isoformat()
        except Exception:
            pass
    return None


def value_date_with_fallback(raw):
    """
    Restituisce:
      - data in ISO se parsabile
      - altrimenti la stringa raw (se non vuota) — così non perdi il dato originale
      - altrimenti None
    """
    if raw is None:
        return None
    parsed = try_parse_date_raw(raw)
    if parsed:
        return parsed
    s = str(raw).strip()
    return s if s != "" and s.upper() != "NONE" else None


def _parse_dmy_components(day, month, year) -> date | None:
    """
    Supporta day/month/year passati come numeri o stringhe.
    Se year è a 2 cifre lo interpreta come 20xx se <50 altrimenti 19xx.
    """
    try:
        d = int(day)
        m = int(month)
        y = int(year)
        if y < 100:
            y = 2000 + y if y < 50 else 1900 + y
        return date(y, m, d)
    except Exception:
        return None


def try_parse_date_raw(raw):
    """Prova a parsare la data con _parse_date; se non parsabile ritorna None."""
    try:
        from database import _parse_date  # se esiste, preferiscilo
        dt = _parse_date(raw)
        return dt.isoformat() if dt else None
    except Exception:
        # fallback alla nostra versione interna
        try:
            dt = _parse_date(raw)
            return dt.isoformat() if dt else None
        except Exception:
            return None

def _parse_date(value) -> date | None:
    """
    Parsing robusto.
    - accetta date/datetime, int/float (Excel serial), stringhe
    - gestisce formati con / - . spazi e anche concatenati come 01012020 o 20200101
    - in caso di ambiguità prioritizza il formato europeo DD/MM/YYYY
    """
    if value is None:
        return None

    # tipo già data
    if isinstance(value, date) and not isinstance(value, datetime):
        return value
    if isinstance(value, datetime):
        return value.date()

    # Excel-style serial
    if isinstance(value, (int, float)):
        try:
            epoch = datetime(1899, 12, 30)
            return (epoch + timedelta(days=int(value))).date()
        except Exception:
            return None

    s = str(value).strip()
    if s == "":
        return None

    s0 = s.split()[0]  # ignora eventuale orario

    # 1) Pattern con separatori (/, -, .) - già presente nel tuo codice
    m = re.match(r'^(\d{1,2})[\/\-\.](\d{1,2})[\/\-\.](\d{2,4})$', s0)
    if m:
        p1 = int(m.group(1))
        p2 = int(m.group(2))
        # se secondo > 12 allora formato mm/dd/yyyy
        if p2 > 12 and p1 <= 12:
            fmts = ["%m/%d/%Y", "%m/%d/%y", "%m-%d-%Y", "%m-%d-%y", "%m.%d.%Y", "%m.%d.%y"]
            for f in fmts:
                try:
                    return datetime.strptime(s0, f).date()
                except Exception:
                    pass
        # se primo > 12 allora formato dd/mm/yyyy
        if p1 > 12 and p2 <= 12:
            fmts = ["%d/%m/%Y", "%d/%m/%y", "%d-%m-%Y", "%d-%m-%y", "%d.%m.%Y", "%d.%m.%y"]
            for f in fmts:
                try:
                    return datetime.strptime(s0, f).date()
                except Exception:
                    pass
        # ambigui -> prima europeo poi americano
        for f in ("%d/%m/%Y","%d/%m/%y","%m/%d/%Y","%m/%d/%y","%d-%m-%Y","%d-%m-%y","%m-%d-%Y","%m-%d-%y"):
            try:
                return datetime.strptime(s0, f).date()
            except Exception:
                pass

    # 2) Pattern con spazi: "1 1 2020" o "01 01 20"
    m_space = re.match(r'^(\d{1,2})\s+(\d{1,2})\s+(\d{2,4})$', s0)
    if m_space:
        p1 = int(m_space.group(1))
        p2 = int(m_space.group(2))
        year_part = m_space.group(3)
        s_space = f"{p1:02d}/{p2:02d}/{year_part}"
        # riusa logica sopra: proviamo dd/mm poi mm/dd
        for f in ("%d/%m/%Y","%d/%m/%y","%m/%d/%Y","%m/%d/%y"):
            try:
                return datetime.strptime(s_space, f).date()
            except Exception:
                pass

    # 3) Solo cifre (concatenate), es: 01012020 (DDMMYYYY), 20200101 (YYYYMMDD), 1012020 (DMMYYYY), 010120 (DDMMYY)
    if s0.isdigit():
        L = len(s0)
        # proviamo una serie di formati ordinati per robustezza
        digit_formats = []
        if L == 8:
            # può essere YYYYMMDD o DDMMYYYY -> proviamo YYYYMMDD poi DDMMYYYY? (preferisco provare YYYYMMDD se il primo 4 è anno plausibile)
            # controlliamo se i primi 4 costituiscono un anno plausibile (> 1900)
            first4 = int(s0[:4])
            if 1900 <= first4 <= 3000:
                digit_formats = ["%Y%m%d", "%d%m%Y"]
            else:
                digit_formats = ["%d%m%Y", "%Y%m%d"]
        elif L == 6:
            # DDMMYY o YYMMDD: proviamo DDMMYY poi YYMMDD
            digit_formats = ["%d%m%y", "%y%m%d"]
        elif L == 7:
            # es. DMMYYYY o YYYYMM D? proviamo alcuni candidati
            digit_formats = ["%d%m%Y", "%Y%m%d"]
        elif L in (4,5):
            # poco probabile, proviamo ddmy / dmyy: fallback coi parser comuni
            digit_formats = ["%d%m%Y","%d%m%y","%Y%m%d"]
        else:
            # fallback generico: proviamo alcuni formati utili
            digit_formats = ["%d%m%Y","%Y%m%d","%d%m%y","%y%m%d"]

        for f in digit_formats:
            try:
                return datetime.strptime(s0, f).date()
            except Exception:
                pass

    # 4) formati non slash già presenti nel tuo codice
    for f in ("%Y-%m-%d", "%Y/%m/%d", "%Y%m%d", "%d%m%Y", "%d.%m.%Y"):
        try:
            return datetime.strptime(s0, f).date()
        except Exception:
            pass

    # 5) try ISO fallback
    try:
        return datetime.fromisoformat(s0).date()
    except Exception:
        return None


def _parse_duration_to_months(s) -> int | None:
    """
    Accetta: int/float/str.
    - '12' => 12 mesi
    - '2y','2 anni' => 24 mesi
    - '36m','36 mesi' => 36 mesi
    - '365' (numero >=365) interpretato come giorni -> months = days//30
    """
    if s is None:
        return None

    if isinstance(s, (int, float)):
        v = int(s)
        # heuristic: if large (>365) treat as days
        if v >= 365:
            return max(1, v // 30)
        return v

    t = str(s).strip().lower()
    if t == "":
        return None

    t = t.replace(",", "").replace("+", "").replace(".", "").strip()

    # anni
    m = re.match(r"^(\d+)\s*(y|yr|year|anni|anno)$", t)
    if m:
        return int(m.group(1)) * 12
    # mesi
    m = re.match(r"^(\d+)\s*(m|mo|mes|mese|mesi)$", t)
    if m:
        return int(m.group(1))
    # giorni
    m = re.match(r"^(\d+)\s*(d|day|giorni|giorno)$", t)
    if m:
        days = int(m.group(1))
        return max(1, days // 30)
    # solo numero
    if t.isdigit():
        v = int(t)
        if v >= 365:
            return max(1, v // 30)
        return v

    return None


def _add_months(dt: date, months: int) -> date:
    year = dt.year + (dt.month - 1 + months) // 12
    month = (dt.month - 1 + months) % 12 + 1
    day = dt.day
    last_day = calendar.monthrange(year, month)[1]
    if day > last_day:
        day = last_day
    return date(year, month, day)


def _parse_days(s) -> int | None:
    """
    Interpreta s come numero di giorni.
    - accetta int/float -> giorni
    - accetta str: '365', '365d', '365 giorni'
    - ritorna None se non interpretabile
    """
    if s is None:
        return None
    # se già numero
    if isinstance(s, (int, float)):
        try:
            return int(s)
        except Exception:
            return None
    t = str(s).strip().lower()
    if t == "":
        return None
    t = t.replace(",", "").replace("+", "").replace(".", "").strip()
    m = re.match(r"^(\d+)\s*(d|day|giorni|giorno)?$", t)
    if m:
        return int(m.group(1))
    # fallback: se è solo cifra
    if t.isdigit():
        return int(t)
    return None

# -------------------------
# calcolo fine memoria
# -------------------------
def calcola_fine_memoria(datainst, mxaz) -> str | None:
    parsed = _parse_date(datainst)
    if not parsed:
        return None

    try:
        days = int(str(mxaz).strip())
    except Exception:
        return None

    scadenza = parsed + timedelta(days=days)
    return scadenza.isoformat()


# -------------------------
# calcolo garanzia (default 12 mesi)
# -------------------------
def calcola_scadenza_garanzia(datainst, durata=None) -> str | None:
    """
    GARANZIA = DATAINST + durata (se durata fornita, la interpretiamo in mesi; altrimenti 12 mesi).
    Se DATAINST non è parsabile, ritorna None.
    """
    parsed = _parse_date(datainst)
    if not parsed:
        return None

    # se fornisci durata la interpretiamo come mesi (sosteniamo diversi formati in _parse_duration_to_months)
    months = _parse_duration_to_months(durata) if durata is not None else None
    if months is None:
        months = 12

    scadenza = _add_months(parsed, months)
    return scadenza.isoformat()

# -------------------------
# calcolo contatore punti cassa
# -------------------------
def calcola_codpuntocassa(table: str = "clienti", column: str = "CODPUNTOCA") -> int:
    """
    Assegna un contatore partendo dal basso (id DESC) con 1 nella riga con id più alto.
    Sovrascrive eventuali valori esistenti. Ritorna il numero di righe aggiornate.
    Esempio: calcola_codpuntocassa("clienti"), calcola_codpuntocassa("clienti_clean").
    """
    conn = get_connection()
    cur = conn.cursor()

    # verifica esistenza tabella
    cur.execute("SELECT name FROM sqlite_master WHERE type='table' AND name=?", (table,))
    if not cur.fetchone():
        conn.close()
        return 0

    # assicura colonna esista
    cur.execute(f"PRAGMA table_info({table})")
    existing_cols = [r[1] for r in cur.fetchall()]
    if column not in existing_cols:
        cur.execute(f'ALTER TABLE {table} ADD COLUMN "{column}" TEXT')

    # recupera gli id dal basso verso l'alto (ORDER BY id DESC)
    cur.execute(f"SELECT id FROM {table} ORDER BY id ASC")
    rows = cur.fetchall()
    if not rows:
        conn.close()
        return 0

    updated = 0
    try:
        i = 1
        for (rid,) in rows:
            cur.execute(f'UPDATE {table} SET "{column}" = ? WHERE id = ?', (str(i), rid))
            i += 1
            updated += 1
        conn.commit()
    except Exception:
        conn.rollback()
        raise
    finally:
        conn.close()
    return updated


def aggiorna_finemf_e_garanzia(ids: list[int] | None = None) -> dict:
    conn = get_connection()
    cur = conn.cursor()
    where_empty = "(FINEMF IS NULL OR TRIM(FINEMF) = '' OR FINEMF = '0' OR GARANZIA IS NULL OR TRIM(GARANZIA) = '' OR GARANZIA = '0')"
    if ids:
        placeholders = ",".join("?" for _ in ids)
        sql = f"SELECT id, DATAINST, MXAZ, FINEMF, GARANZIA FROM clienti WHERE id IN ({placeholders}) AND {where_empty}"
        cur.execute(sql, tuple(ids))
    else:
        sql = f"SELECT id, DATAINST, MXAZ, FINEMF, GARANZIA FROM clienti WHERE {where_empty}"
        cur.execute(sql)

    rows = cur.fetchall()
    updated_finemf = []
    updated_garanzia = []
    failed = []

    try:
        for cid, datainst, mxaz, finemf, gar in rows:
            before_changes = conn.total_changes

            # GARANZIA
            try:
                if gar is None or str(gar).strip() == "" or str(gar).strip() == "0":
                    nuova_g = calcola_scadenza_garanzia(datainst, None)
                    if nuova_g:
                        cur.execute(
                            "UPDATE clienti SET GARANZIA = ? WHERE id = ? AND (GARANZIA IS NULL OR TRIM(GARANZIA) = '' OR GARANZIA = '0')",
                            (nuova_g, cid)
                        )
                        if conn.total_changes > before_changes:
                            updated_garanzia.append(cid)
                            before_changes = conn.total_changes
                    else:
                        failed.append((cid, "GARANZIA non calcolabile", {"DATAINST": datainst, "MXAZ": mxaz}))
            except Exception as e:
                failed.append((cid, f"Errore GARANZIA: {e}", {"DATAINST": datainst, "MXAZ": mxaz}))

            # FINEMF
            try:
                if finemf is None or str(finemf).strip() == "" or str(finemf).strip() == "0":
                    nuovo_fin = calcola_fine_memoria(datainst, mxaz)
                    if nuovo_fin:
                        cur.execute(
                            "UPDATE clienti SET FINEMF = ? WHERE id = ? AND (FINEMF IS NULL OR TRIM(FINEMF) = '' OR FINEMF = '0')",
                            (nuovo_fin, cid)
                        )
                        if conn.total_changes > before_changes:
                            updated_finemf.append(cid)
                            before_changes = conn.total_changes
                    else:
                        # dettaglio del motivo di fallimento
                        motivo = []
                        if _parse_date(datainst) is None:
                            motivo.append("DATAINST non parsabile")
                        if _parse_days(mxaz) is None:
                            motivo.append("MXAZ non parsabile/assente")
                        failed.append((cid, "FINEMF non calcolabile", {"DATAINST": datainst, "MXAZ": mxaz, "motivo": ";".join(motivo)}))
            except Exception as e:
                failed.append((cid, f"Errore FINEMF: {e}", {"DATAINST": datainst, "MXAZ": mxaz}))

        conn.commit()
    except Exception as e:
        conn.rollback()
        conn.close()
        return {"error": str(e)}
    finally:
        conn.close()

    return {
        "rows_checked": len(rows),
        "updated_finemf": updated_finemf,
        "updated_garanzia": updated_garanzia,
        "failed": failed,
    }

#per controllare problemi legati all'import
#def inspect_db_sample(): 
#    conn = get_connection()
#    cur = conn.cursor()
#    cur.execute("PRAGMA table_info(clienti)")
#    cols = [r[1] for r in cur.fetchall()]
#    print("DB columns:", cols)
#    cur.execute("SELECT id, FINEMF, TIPOMF, DATAINST, MXAZ, GARANZIA FROM clienti LIMIT 10")
#    rows = cur.fetchall()
#    for row in rows:
#        print("ROW:", row)
#    conn.close()

# -------------------------
# INTAE
# -------------------------
def import_csv_intae(path: str, delimiter: str = ";", encoding: str = "utf-8") -> int:
    """
    Importa CSV per intae. Rileva delimiter automaticamente e mappa header
    alle colonne INTAE_COLUMNS. Valori mancanti -> stringa vuota.
    """
    import csv, os
    if not os.path.exists(path):
        raise FileNotFoundError(path)

    # sniffing del delimiter
    with open(path, "r", encoding=encoding, errors="replace") as sf:
        sample = sf.read(8192)
        sf.seek(0)
        try:
            sniffer = csv.Sniffer()
            dialect = sniffer.sniff(sample, delimiters=",;|\t")
            detected = dialect.delimiter
        except Exception:
            # fallback semplice
            detected = ";" if ";" in sample else ("," if "," in sample else ("\t" if "\t" in sample else ","))

    tried = []
    # proviamo con detected poi con qualche fallback
    for delim in [detected, delimiter, ";", ",", "\t", "|"]:
        if not delim or delim in tried:
            continue
        tried.append(delim)
        try:
            with open(path, newline='', encoding=encoding, errors='replace') as f:
                reader = csv.DictReader(f, delimiter=delim)
                headers = reader.fieldnames or []
                if not headers:
                    raise ValueError("CSV privo di intestazioni")
                # mapping header -> INTAE_COLUMNS
                map_header_to_col = {}
                wanted = INTAE_COLUMNS.copy()
                # normalizzazione semplice
                def norm(s): return "".join(ch for ch in (s or "").upper() if ch.isalnum())
                wanted_norm = {w: norm(w) for w in wanted}
                for h in headers:
                    hn = norm(h)
                    matched = None
                    for wc, wn in wanted_norm.items():
                        if wn in hn or hn in wn or wn == hn:
                            matched = wc
                            break
                    if matched:
                        map_header_to_col[h] = matched

                                # costruzione records ma SKIP se (matricola + data_intervento) già presenti nel DB
                records = []

                # normalize helper per matricola
                def _norm_mat(s):
                    return re.sub(r'[^0-9A-Za-z]', '', (s or "")).lower().strip()

                # normalizza la data/intervento includendo l'orario quando presente
                def _norm_date_for_compare(s):
                    if s is None:
                        return ""
                    s_str = str(s).strip()
                    if s_str == "":
                        return ""
                    try:
                        d = _try_parse_date(s_str)  # funzione esistente nel file
                    except Exception:
                        d = None
                    if d:
                        # se è un datetime, tolgo i microsecondi per stabilità e uso ISO completo
                        try:
                            if hasattr(d, "hour"):
                                # datetime-like
                                if getattr(d, "microsecond", 0):
                                    d = d.replace(microsecond=0)
                            return d.isoformat()
                        except Exception:
                            # fallback a stringa pulita
                            return s_str
                    # fallback: ritorna la stringa originale normalizzata
                    return s_str

                # apri conn e leggi coppie esistenti per intae (matricola, data_intervento)
                conn = get_connection()
                cur = conn.cursor()
                existing_pairs = set()
                try:
                    # leggo info colonne tabella intae
                    cur.execute("PRAGMA table_info(intae)")
                    cols_info = [r[1].lower() for r in cur.fetchall()]

                    mat_col = None
                    date_col = None
                    if "matricola" in cols_info:
                        mat_col = "matricola"
                    if "data_intervento" in cols_info:
                        date_col = "data_intervento"

                    # se abbiamo almeno matricola, carichiamo le coppie esistenti (matricola, data_intervento)
                    if mat_col:
                        if date_col:
                            cur.execute(f"SELECT {mat_col}, {date_col} FROM intae")
                            for row in cur.fetchall():
                                try:
                                    mm = row[0]
                                    dd = row[1] if len(row) > 1 else None
                                    mnorm = _norm_mat(mm)
                                    dnorm = _norm_date_for_compare(dd)
                                    if mnorm:
                                        existing_pairs.add((mnorm, dnorm))
                                except Exception:
                                    continue
                        else:
                            # se non c'è la colonna data_intervento, carichiamo solo matricole (data vuota)
                            cur.execute(f"SELECT {mat_col} FROM intae")
                            for (mm,) in cur.fetchall():
                                try:
                                    mnorm = _norm_mat(mm)
                                    if mnorm:
                                        existing_pairs.add((mnorm, ""))
                                except Exception:
                                    continue
                except Exception:
                    # se qualcosa fallisce lasciamo existing_pairs vuoto (importerà tutto)
                    existing_pairs = set()
                finally:
                    try:
                        conn.close()
                    except Exception:
                        pass

                # ora scansioniamo il CSV e aggiungiamo solo righe non duplicate
                for r in reader:
                    rec = []
                    for col in INTAE_COLUMNS:
                        # cerca l'header originale che mappa a questa colonna
                        found = False
                        for orig_h, mapped in map_header_to_col.items():
                            if mapped == col:
                                v = r.get(orig_h, "")
                                rec.append((v or "").strip())
                                found = True
                                break
                        if not found:
                            rec.append("")  # campo assente -> vuoto

                    # verifica duplicato locale: matricola + data_intervento (incluso orario)
                    try:
                        mat_idx = INTAE_COLUMNS.index("MATRICOLA") if "MATRICOLA" in INTAE_COLUMNS else None
                        data_idx = INTAE_COLUMNS.index("DATA_INTERVENTO") if "DATA_INTERVENTO" in INTAE_COLUMNS else None

                        mat_val = rec[mat_idx] if (mat_idx is not None and mat_idx < len(rec)) else ""
                        dt_val = rec[data_idx] if (data_idx is not None and data_idx < len(rec)) else None

                        mnorm = _norm_mat(mat_val)
                        dnorm = _norm_date_for_compare(dt_val)

                        # se entrambe presenti e la coppia è già in existing_pairs -> skip
                        if mnorm and dnorm and (mnorm, dnorm) in existing_pairs:
                            # riga duplicata: saltiamo (non la aggiungiamo ai records)
                            continue
                        else:
                            # se non duplicata la aggiungiamo e, per evitare duplicates tra le righe importate,
                            # aggiungiamo anche alla set existing_pairs
                            records.append(tuple(rec))
                            if mnorm:
                                existing_pairs.add((mnorm, dnorm))
                    except Exception:
                        # in caso di eccezione nella verifica, aggiungiamo comunque la riga (strategia conservativa)
                        records.append(tuple(rec))


                # inserimento
                conn = get_connection()
                try:
                    cur = conn.cursor()
                    colnames_sql = ", ".join([f'"{c}"' for c in INTAE_COLUMNS])
                    placeholders = ", ".join(["?"] * len(INTAE_COLUMNS))
                    if records:
                        # usiamo il context manager della connessione per commit/rollback automatico
                        with conn:
                            cur.executemany(
                                f'INSERT INTO intae ({colnames_sql}) VALUES ({placeholders})',
                                records
                            )
                    # cur.rowcount non è affidabile su sqlite, meglio contare i records
                    n = len(records)
                    return n
                except Exception:
                    # con `with conn:` il rollback è automatico, ma qui lo mettiamo per sicurezza
                    conn.rollback()
                    raise
                finally:
                    conn.close()
                return n
        except Exception:
            # try next delimiter
            continue

    raise RuntimeError("Impossibile importare il CSV INTAE (delimiter non riconosciuto o file corrotto)")

def read_intae(limit: int = 500) -> list:
    conn = get_connection()
    cur = conn.cursor()
    colnames = ", ".join([f'"{c}"' for c in INTAE_COLUMNS])
    cur.execute(f"""
        SELECT id, {colnames}
        FROM intae
        ORDER BY DATA_INTERVENTO DESC
        LIMIT ?
    """, (limit,))
    rows = cur.fetchall()
    conn.close()
    return rows

def cerca_interventi_intae(matricola: str) -> list:
    """
    Cerca per matricola (case-insensitive). Se la colonna MATRICOla non esiste
    restituisce lista vuota anziché lanciare un'eccezione.
    """
    conn = get_connection()
    try:
        cur = conn.cursor()
        # controllo che la colonna esista davvero nella tabella
        cur.execute("PRAGMA table_info(intae)")
        cols = [row[1].upper() for row in cur.fetchall()]  # row[1] è il nome della colonna
        if "MATRICOLA" not in cols:
            return []

        colnames = ", ".join([f'"{c}"' for c in INTAE_COLUMNS])
        # usiamo COLLATE NOCASE per ricerca case-insensitive anche più semplice
        cur.execute(f"""
            SELECT id, {colnames}
            FROM intae
            WHERE lower("MATRICOLA") = lower(?) OR "MATRICOLA" LIKE ? COLLATE NOCASE
            ORDER BY DATA_INTERVENTO DESC
        """, (matricola, f"%{matricola}%"))
        rows = cur.fetchall()
        return rows
    finally:
        conn.close()

# -------------------------
# MISURART
# -------------------------

def misurart_add_from_map(data: dict) -> int:
    """
    Inserisce una registrazione misurart usando la mappa data (chiavi corrispondono a MISURART_COLUMNS).
    Campi mancanti -> stringa vuota.
    """
    conn = get_connection()
    cur = conn.cursor()
    # assicurati colonne esistenti già fatte in init_database
    colnames_sql = ", ".join([f'"{c}"' for c in MISURART_COLUMNS])
    placeholders = ", ".join(["?"] * len(MISURART_COLUMNS))
    vals = [str(data.get(c, "")) for c in MISURART_COLUMNS]
    cur.execute(f'INSERT INTO misurart ({colnames_sql}) VALUES ({placeholders})', vals)
    nid = cur.lastrowid
    conn.commit()

    # Se il TIPO è '2', azzera VF del cliente associato alla stessa matricola
    try:
        tipo_val = data.get("TIPO", "")
        matricola_val = data.get("MATRICOLA", "") or data.get("MATRICA", "")
        if str(tipo_val).strip() == "2" and matricola_val:
            _clear_vf_for_matricola(matricola_val)
    except Exception:
        pass

    conn.close()
    return nid

def misurart_read_all(limit: int = 500) -> list:
    conn = get_connection()
    cur = conn.cursor()
    colnames = ", ".join([f'"{c}"' for c in MISURART_COLUMNS])
    cur.execute(f"SELECT id, {colnames} FROM misurart ORDER BY DATAINIZIO DESC LIMIT ?", (limit,))
    rows = cur.fetchall()
    conn.close()
    return rows

def misurart_update(record_id: int, updates: dict):
    if not updates:
        return
    # Assicuriamoci le colonne esistano nella tabella misurart (utile se l'update aggiunge colonne)
    conn = get_connection()
    cur = conn.cursor()
    cur.execute("PRAGMA table_info(misurart)")
    existing = [r[1] for r in cur.fetchall()]
    for c in updates.keys():
        if c not in existing:
            cur.execute(f'ALTER TABLE misurart ADD COLUMN "{c}" TEXT')
    # costruzione SQL per update
    cols = list(updates.keys())
    parts = ", ".join([f'"{c}"=?' for c in cols])
    vals = [str(updates[c]) for c in cols] + [record_id]
    cur.execute(f'UPDATE misurart SET {parts} WHERE id=?', vals)
    conn.commit()

    # Se TIPO è stato impostato a '2' (nella update) -> troviamo la matricola (se fornita nell'updates o leggendo la riga)
    try:
        tipo_was_updated = "TIPO" in updates
        tipo_val = updates.get("TIPO", "")
        matricola_val = updates.get("MATRICOLA", "") or updates.get("MATRICA", "")
        if tipo_was_updated and str(tipo_val).strip() == "2":
            # se matricola non era nella update, leggiamola dal DB per questo record
            if not matricola_val:
                cur.execute("SELECT MATRICA, MATRICO LA, MATRICOLA FROM misurart WHERE id=?", (record_id,))
                row = cur.fetchone()
                if row:
                    # prova a prendere il primo valore non-null nella riga (compatibilità nomi)
                    for v in row:
                        if v:
                            matricola_val = v
                            break
            if matricola_val:
                _clear_vf_for_matricola(matricola_val)
    except Exception:
        pass

    conn.close()

def misurart_delete(record_id: int):
    conn = get_connection()
    cur = conn.cursor()
    cur.execute("DELETE FROM misurart WHERE id=?", (record_id,))
    conn.commit()
    conn.close()

def misurart_search_by_matricola(matricola: str) -> list:
    conn = get_connection()
    cur = conn.cursor()
    colnames = ", ".join([f'"{c}"' for c in MISURART_COLUMNS])
    cur.execute(f"SELECT id, {colnames} FROM misurart WHERE MATRICA LIKE ? OR MATRICOLA LIKE ? ORDER BY DATAINIZIO DESC", (f"%{matricola}%", f"%{matricola}%"))
    rows = cur.fetchall()
    conn.close()
    return rows

def read_misurart_all():
    """Ritorna tutte le righe della tabella MISURART ordinate per id."""
    conn = get_connection()
    cur = conn.cursor()
    try:
        cur.execute(f"SELECT id, {', '.join(MISURART_COLUMNS)} FROM misurart ORDER BY id ASC")
        return cur.fetchall()
    except Exception as e:
        print("Errore in read_misurart_all:", e)
        return []
    finally:
        conn.close()

def _clear_vf_for_matricola(matricola: str) -> int:
    """
    Cerca una matricola normalizzata (rimuove caratteri non alfanumerici).
    - Prima cerca in macchine.matricola (ottimale, prende cliente_id).
    - Se non trova macchine corrispondenti cerca in clienti.MATRICOLA.
    Aggiorna clienti.VF='' per i clienti trovati e azzera anche macchine.VF per le macchine trovate.
    Ritorna il numero di clienti aggiornati.
    """
    import re
    if not matricola:
        return 0
    target_raw = str(matricola).strip()
    if not target_raw:
        return 0

    def norm(s):
        return re.sub(r'[^0-9A-Za-z]', '', (s or "")).lower()

    target = norm(target_raw)
    conn = get_connection()
    cur = conn.cursor()

    # 1) cerca macchine normalizzando in Python (più robusto)
    cur.execute("SELECT id, matricola, cliente_id FROM macchine")
    mac_rows = cur.fetchall()
    matching_machine_ids = []
    cliente_ids = set()
    for mid, mval, cid in mac_rows:
        if mval is None:
            continue
        if norm(mval) == target or target in norm(mval) or norm(mval) in target:
            matching_machine_ids.append(mid)
            if cid:
                try:
                    cid_int = int(cid)
                    if cid_int > 0:
                        cliente_ids.add(cid_int)
                except Exception:
                    pass

    # 2) se non ho clienti trovati da macchine, cerco in clienti.MATRICOLA
    if not cliente_ids:
        cur.execute("SELECT id, MATRICOLA FROM clienti")
        for cid, mat in cur.fetchall():
            if mat and norm(mat) == target or (mat and target in norm(mat)) or (mat and norm(mat) in target):
                try:
                    if int(cid) > 0:
                        cliente_ids.add(int(cid))
                except Exception:
                    pass

    if not cliente_ids and not matching_machine_ids:
        conn.close()
        return 0

    updated_clients = 0
    if cliente_ids:
        placeholders = ", ".join(["?"] * len(cliente_ids))
        sql = f"UPDATE clienti SET VF='' WHERE id IN ({placeholders})"
        cur.execute(sql, list(cliente_ids))
        # rowcount può essere -1 su alcune build; usiamo fallback
        updated_clients = cur.rowcount if (cur.rowcount and cur.rowcount > 0) else len(cliente_ids)

    # azzera VF anche nelle macchine trovate (se ci sono)
    if matching_machine_ids:
        placeholders2 = ", ".join(["?"] * len(matching_machine_ids))
        cur.execute(f"UPDATE macchine SET VF='' WHERE id IN ({placeholders2})", matching_machine_ids)

    conn.commit()
    conn.close()
    return updated_clients

# -------------------------
# Bootstrap helper
# -------------------------

def _try_parse_date(s: Optional[str]) -> Optional[date]:
    """
    Ritorna datetime.date o None.
    Versione robusta che estrae sottostringhe date-like e prova più strategie.
    """
    if not s:
        return None
    raw = str(s)
    # normalizza spazi invisibili e trim
    raw = raw.replace('\u00A0', ' ').replace('\u200B', '').strip()
    if raw == "":
        return None

    # prova diretta su raw (ISO)
    try:
        return datetime.fromisoformat(raw).date()
    except Exception:
        pass

    # rimuovi orari come "10:00" e parentesi, e token timezone
    cleaned = re.sub(r'\b\d{1,2}:\d{2}(?::\d{2})?\b', ' ', raw)
    cleaned = re.sub(r'[\(\)\[\]]', ' ', cleaned)
    cleaned = re.sub(r'\b(?:utc|cet|cest|gmt|[+\-]?\d{1,2}:?\d{0,2})\b', ' ', cleaned, flags=re.I)
    cleaned = cleaned.strip()

    # 1) prova formati espliciti sull'intera stringa pulita
    fmts = [
        "%d/%m/%Y", "%d-%m-%Y", "%d.%m.%Y",
        "%d/%m/%y", "%d-%m-%y", "%d.%m.%y",
        "%Y/%m/%d", "%Y-%m-%d", "%Y.%m.%d",
        "%Y%m%d", "%d%m%Y"
    ]
    for cand in (cleaned, raw):
        for f in fmts:
            try:
                return datetime.strptime(cand, f).date()
            except Exception:
                pass

    # 2) normalizza separatori e riprova
    norm = re.sub(r'[\/\.]', '-', cleaned)
    for f in ("%d-%m-%Y", "%d-%m-%y", "%Y-%m-%d"):
        try:
            return datetime.strptime(norm, f).date()
        except Exception:
            pass

    # 3) estrai la prima sottostringa "num sep num sep num" e prova quella
    m = re.search(r'(\d{1,4}\D+\d{1,2}\D+\d{1,4})', cleaned)
    if m:
        sub = m.group(1).strip()
        sub_norm = re.sub(r'[\/\.]', '-', sub)
        # prova patterns sui sottostringhe
        for f in ("%d-%m-%Y", "%d-%m-%y", "%Y-%m-%d", "%d-%m-%Y"):
            try:
                return datetime.strptime(sub_norm, f).date()
            except Exception:
                pass
        # heuristica numerica
        mm = re.match(r'^\s*(\d{1,4})\D+(\d{1,2})\D+(\d{1,4})\s*$', sub)
        if mm:
            a, b, c = mm.group(1), mm.group(2), mm.group(3)
            try:
                if len(a) == 4:
                    y, mth, d = int(a), int(b), int(c)
                elif len(c) == 4:
                    d, mth, y = int(a), int(b), int(c)
                else:
                    d, mth, y = int(a), int(b), int(c)
                    if y < 100:
                        y = 2000 + y if y < 50 else 1900 + y
                return date(y, mth, d)
            except Exception:
                pass

    # 4) sostituzioni mesi italiani -> inglesi e fallback a dateutil
    italian_to_english_months = {
        r'\bgennario\b': 'january', r'\bgen\b': 'jan',
        r'\bfebbraio\b': 'february', r'\bfeb\b': 'feb',
        r'\bmarzo\b': 'march', r'\bmar\b': 'mar',
        r'\baprile\b': 'april', r'\bapr\b': 'apr',
        r'\bmaggio\b': 'may', r'\bmag\b': 'may',
        r'\bgiugno\b': 'june', r'\bgiu\b': 'jun',
        r'\bluglio\b': 'july', r'\blug\b': 'jul',
        r'\bagosto\b': 'august', r'\bago\b': 'aug',
        r'\bsettembre\b': 'september', r'\bset\b': 'sep',
        r'\bottobre\b': 'october', r'\bott\b': 'oct',
        r'\bnovembre\b': 'november', r'\bnov\b': 'nov',
        r'\bdicembre\b': 'december', r'\bdic\b': 'dec'
    }
    token_for_parser = cleaned or raw
    for pat, en in italian_to_english_months.items():
        token_for_parser = re.sub(pat, en, token_for_parser, flags=re.I)

    try:
        return parser.parse(token_for_parser, dayfirst=True, fuzzy=True).date()
    except Exception:
        pass

    return None

def ensure_db_exists(recreate_clienti: bool = False):
    if not os.path.exists(DB_FILE):
        init_database(recreate_clienti=recreate_clienti)
        return

    conn = get_connection()
    cur = conn.cursor()
    # leggi una sola volta l'output di PRAGMA table_info
    cur.execute("PRAGMA table_info(clienti)")
    rows = cur.fetchall()
    existing = [r[1] for r in rows] if rows else []

    for c in CLIENTI_COLUMNS:
        if c not in existing:
            try:
                cur.execute(f'ALTER TABLE clienti ADD COLUMN \"{c}\" TEXT')
            except sqlite3.OperationalError as e:
                # ignora errori di colonna duplicata (possibile race / stato inconsistente)
                if 'duplicate column name' in str(e).lower():
                    continue
                else:
                    raise
    conn.commit()
    conn.close()

# -------------------------
# Report
# -------------------------

def report_memorie_in_scadenza(within_days: int = 30, db_path: str = "database.db") -> List[Tuple[str,str,str,str,str]]:
    conn = sqlite3.connect(db_path)
    conn.row_factory = sqlite3.Row
    cur = conn.cursor()

    today = date.today()
    end_date = today + timedelta(days=within_days)

    out = []

    # 1) prova la via veloce in SQL sulla tabella 'clienti' (date() su FINEMF)
    try:
        cur.execute("""
            SELECT RAG_SOC, FINEMF, MF, TEL, CODICE
            FROM clienti
            WHERE FINEMF IS NOT NULL AND TRIM(FINEMF) != '' AND date(FINEMF) BETWEEN ? AND ?
            ORDER BY FINEMF
        """, (today.isoformat(), end_date.isoformat()))
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

    # 2) fallback: leggi tutte le righe non vuote da 'clienti' e parse in Python usando _parse_date se disponibile
    cur.execute("SELECT RAG_SOC, FINEMF, MF, TEL, CODICE FROM clienti WHERE FINEMF IS NOT NULL AND TRIM(FINEMF) != ''")
    rows = cur.fetchall()
    print("[report_memorie_in_scadenza] righe con FINEMF non vuota (fetched):", len(rows))

    # preferiamo usare la tua funzione di parsing interna se esiste (_parse_date)
    try:
        from database import _parse_date as _local_parse_date
    except Exception:
        _local_parse_date = None

    unparsed = []
    for row in rows:
        rag, finemf_raw, mf, tel, codice = row["RAG_SOC"], row["FINEMF"], row["MF"], row["TEL"], row["CODICE"]
        parsed_date = None

        if _local_parse_date:
            try:
                parsed_date = _local_parse_date(finemf_raw)
            except Exception:
                parsed_date = None

        # fallback ad una semplice heuristica se _parse_date non esiste o fallisce
        if not parsed_date:
            import re
            s = str(finemf_raw).strip()
            s = s.split()[0]  # togli eventuale time part
            # prova fromisoformat
            try:
                from datetime import datetime
                parsed_date = datetime.fromisoformat(s).date()
            except Exception:
                # prova formati comuni
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

        if today <= parsed_date <= end_date:
            display_date = parsed_date.strftime("%d/%m/%Y")
            out.append((rag or "", display_date, mf or "", tel or "", codice or ""))

    conn.close()
    print("[report_memorie_in_scadenza] parsed-in-range examples:", out[:10])
    print("[report_memorie_in_scadenza] non parsati (esempi):", unparsed[:10])
    print("[report_memorie_in_scadenza] totale trovati in range:", len(out))
    return out

def report_rt_for_clienti(month: int | None = None, year: int | None = None):
    """
    Versione diagnostica e robusta:
    - individua dinamicamente la colonna 'finemf' tra le colonne della tabella clienti
    - prova diversi metodi di parsing e stampa info di debug su console
    - filtra per month/year se forniti, altrimenti ritorna tutto (compatibilità)
    Restituisce liste di tuple (prime 8 colonne come stringhe) come la GUI si aspetta.
    """
    import traceback
    from dateutil import parser as _parser

    conn = get_connection()
    cur = conn.cursor()

    # 1) leggi una riga per ottenere i nomi delle colonne
    try:
        cur.execute("SELECT * FROM clienti LIMIT 1")
        cols = [d[0] for d in cur.description] if cur.description else []
    except Exception:
        # tabella non presente o altro errore: fallback a comportamento precedente
        print("[DEBUG][DB] Impossibile leggere schema clienti:", traceback.format_exc())
        cols = []

    print("[DEBUG][DB] colonne tabella clienti:", cols)

    # 2) prova a trovare la colonna finemf con heuristics
    finemf_col = None
    if cols:
        lowered = [c.lower() for c in cols]
        # priorità: 'finemf' esatta, poi qualsiasi contenente 'finem', 'fine', 'scaden', 'datafine'
        for target in ('finemf',):
            for i, name in enumerate(lowered):
                if name == target:
                    finemf_col = cols[i]; break
            if finemf_col:
                break
        if not finemf_col:
            for key in ('finem', 'fine', 'scad', 'scaden', 'datafine', 'datainizio', 'data'):
                for i, name in enumerate(lowered):
                    if key in name:
                        finemf_col = cols[i]; break
                if finemf_col:
                    break

    if not finemf_col:
        print("[DEBUG][DB] Colonna 'FINEMF' non trovata (heuristics). Nomi colonne:", cols)
        # fallback: prova comunque a selezionare le colonne standard usate dalla GUI
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
            # normalizza
            return [tuple('' if v is None else str(v) for v in row) for row in rows]
        except Exception:
            conn.close()
            print("[ERROR][DB] Errore fallback SELECT:", traceback.format_exc())
            return []

    print(f"[DEBUG][DB] Uso colonna FINEMF: '{finemf_col}'")

    # 3) recupera tutte le righe compresa la colonna finemf dinamica
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
        COALESCE({finemf_col}, '')
    FROM clienti
    ORDER BY UPPER(COALESCE(COMUNE, '')), UPPER(COALESCE(RAG_SOC, ''))
    """
    try:
        cur.execute(sql)
        all_rows = cur.fetchall()
    except Exception:
        print("[ERROR][DB] Errore SELECT con finemf_col:", traceback.format_exc())
        conn.close()
        return []
    conn.close()

    print(f"[DEBUG][DB] righe lette: {len(all_rows)}")

    # 4) funzione helper per provare a parsare FINEMF con più strategie
    import re
    def try_parse_date(s):
        if not s:
            return None
        s = str(s).strip()
        # rimuovi caratteri invisibili
        s = re.sub(r'[\u00A0\u200B]', ' ', s).strip()
        # tentativo con dateutil (prima con yearfirst True, poi False; dayfirst True/False)
        for kwargs in ({'dayfirst': False, 'yearfirst': True, 'fuzzy': True},
                       {'dayfirst': False, 'yearfirst': False, 'fuzzy': True},
                       {'dayfirst': True, 'yearfirst': False, 'fuzzy': True}):
            try:
                d = _parser.parse(s, **kwargs)
                return d
            except Exception:
                pass
        # fallback a regex per formati comuni
        # YYYY-MM-DD, YYYY/MM/DD
        m = re.search(r'([12][0-9]{3})[-/\\. ]0?([1-9])[-/\\. ]0?([1-9])', s)
        if m:
            try:
                y = int(m.group(1)); mo = int(m.group(2)); da = int(m.group(3))
                from datetime import datetime
                return datetime(y, mo, da)
            except Exception:
                pass
        # DD/MM/YYYY o D/M/YYYY
        m = re.search(r'0?([1-9]|[12][0-9]|3[01])[-/\\. ]0?([1-9]|1[0-2])[-/\\. ]([12][0-9]{3})', s)
        if m:
            try:
                da = int(m.group(1)); mo = int(m.group(2)); y = int(m.group(3))
                from datetime import datetime
                return datetime(y, mo, da)
            except Exception:
                pass
        # MM/YYYY or MM-YYYY
        m = re.search(r'0?([1-9]|1[0-2])[-/\\. ]([12][0-9]{3})', s)
        if m:
            try:
                mo = int(m.group(1)); y = int(m.group(2))
                from datetime import datetime
                return datetime(y, mo, 1)
            except Exception:
                pass
        # YYYYMM or YYYY
        m = re.search(r'^([12][0-9]{3})0?([1-9]|1[0-2])$', s)
        if m:
            try:
                y = int(m.group(1)); mo = int(m.group(2))
                from datetime import datetime
                return datetime(y, mo, 1)
            except Exception:
                pass
        return None

    # 5) applica filtro e logga i campioni per capire cosa succede
    matches = []
    sample_shown = 0
    SAMPLE_LIMIT = 50
    parsed_failures = 0
    parsed_success = 0

    for idx, row in enumerate(all_rows):
        visible = tuple('' if v is None else str(v) for v in row[:8])
        finemf_raw = row[8] if len(row) > 8 else ''
        finemf_raw = '' if finemf_raw is None else str(finemf_raw).strip()

        # mostra i primi SAMPLE_LIMIT valori per debug
        if sample_shown < SAMPLE_LIMIT:
            print(f"[DEBUG][DB] riga {idx} FINEMF grezzo: '{finemf_raw}'")
            sample_shown += 1

        if month is None or year is None:
            matches.append(visible)
            continue

        if not finemf_raw:
            parsed_failures += 1
            continue

        parsed = try_parse_date(finemf_raw)
        if parsed:
            parsed_success += 1
            try:
                if parsed.year == int(year) and parsed.month == int(month):
                    matches.append(visible)
                else:
                    # non coincide: puoi loggare se vuoi
                    pass
            except Exception:
                # conversione anno/mese fallita
                pass
        else:
            parsed_failures += 1

    print(f"[DEBUG][DB] parsed_success={parsed_success}, parsed_failures={parsed_failures}, matches={len(matches)}")
    # show up to 10 matched examples
    for i, m in enumerate(matches[:10]):
        print(f"[DEBUG][DB] esempio match {i}: {m}")

    return matches

def report_checklist_clienti_selected():
    """
    Restituisce solo le colonne richieste (se presenti) e le loro righe:
    RAG_SOC, MATRICOLA, FINEMF, GARANZIA.
    Ritorna [(cols, rows)] per essere compatibile col codice GUI esistente.
    """
    conn = get_connection()
    cur = conn.cursor()

    # controllo esistenza tabella
    cur.execute("SELECT name FROM sqlite_master WHERE type='table'")
    tables = [r[0].lower() for r in cur.fetchall()]
    if "clienti" not in tables:
        conn.close()
        return []

    # prendo le colonne effettive dalla tabella
    cur.execute("PRAGMA table_info(clienti)")
    real_cols = [r[1] for r in cur.fetchall()]
    real_cols_lower = [c.lower() for c in real_cols]

    wanted = ["RAG_SOC", "MATRICOLA", "FINEMF", "GARANZIA"]
    # mantengo l'ordine wanted ma uso i nomi reali (rispetto a come sono in DB)
    selected_cols = []
    for w in wanted:
        try:
            idx = real_cols_lower.index(w.lower())
            selected_cols.append(real_cols[idx])
        except ValueError:
            # colonna non presente -> la salto
            pass

    if not selected_cols:
        conn.close()
        return [([], [])]

    # costruisco la query con le colonne effettive
    q = "SELECT " + ", ".join([f'"{c}"' for c in selected_cols]) + " FROM clienti"
    cur.execute(q)
    rows = cur.fetchall()
    conn.close()
    return [(selected_cols, rows)]

from typing import List, Tuple, Dict
from datetime import date as _date

def report_cambi_dgfe(data_from: str, data_to: str) -> List[Tuple[str, str, int]]:
    """
    Conta gli interventi (righe in 'intae') per ciascuna matricola la cui data
    di intervento è compresa tra data_from e data_to (inclusi).
    Ritorna lista di tuple: (ragione_sociale_cliente_or_empty, matricola, numero_interventi)

    - usa _try_parse_date per normalizzare le date (se non sono parseabili solleva ValueError)
    - è tollerante sui nomi di colonna (cerca matricola/data in intae)
    - risale il nome cliente tramite macchine -> clienti se possibile, altrimenti cerca clienti.matricola
    - non ripete matricole nella lista di output (ogni matricola appare una volta)
    """
    # parse input dates (richiediamo che vengano parseate)
    d1 = _try_parse_date(data_from)
    d2 = _try_parse_date(data_to)
    if not d1 or not d2:
        raise ValueError(f"Date non parseabili: '{data_from}' -> {d1}, '{data_to}' -> {d2}")

    # assicuriamoci che start <= end
    start, end = (d1, d2) if d1 <= d2 else (d2, d1)

    conn = get_connection()
    cur = conn.cursor()
    try:
        # controlla che intae esista
        cur.execute("SELECT name FROM sqlite_master WHERE type='table'")
        tables = {r[0].lower() for r in cur.fetchall()}
        if "intae" not in tables:
            return []

        # trova colonne di intae
        cur.execute("PRAGMA table_info(intae)")
        info = cur.fetchall()
        cols = [r[1] for r in info]
        cols_lower = [c.lower() for c in cols]

        # determina colonna matricola
        matricola_col = None
        for cand in ("matricola", "matricola_macchina", "mat"):
            if cand in cols_lower:
                matricola_col = cols[cols_lower.index(cand)]
                break
        # fallback: prima colonna se non troviamo nulla (meno ideale ma evita eccezioni)
        if not matricola_col and len(cols) > 0:
            matricola_col = cols[0]

        # determina colonna data-intervento (preferenze)
        date_col = None
        for cand in ("data_intervento", "dataintervento", "data", "datavp"):
            if cand in cols_lower:
                date_col = cols[cols_lower.index(cand)]
                break
        # fallback: prima colonna che contiene 'data' nel nome
        if not date_col:
            for i, c in enumerate(cols_lower):
                if "data" in c:
                    date_col = cols[i]
                    break

        # se non troviamo almeno la matricola, non possiamo procedere
        if not matricola_col:
            return []

        # leggiamo le righe (prendiamo matricola e data_col se disponibile)
        select_cols = [matricola_col]
        if date_col:
            select_cols.append(date_col)
        q = f"SELECT {', '.join(select_cols)} FROM intae"
        cur.execute(q)
        rows = cur.fetchall()

        # conta interventi per matricola nel range
        counts: Dict[str, int] = {}
        for row in rows:
            mat_raw = row[0]
            if mat_raw is None:
                continue
            mat = str(mat_raw).strip()
            if mat == "":
                continue

            # se abbiamo data_col, parsala e filtra
            if date_col:
                raw_dt = row[1] if len(row) > 1 else None
                if raw_dt is None:
                    continue
                dt = _try_parse_date(raw_dt)
                if not dt:
                    # salta le righe con data non parsata (opzione conservativa)
                    continue
                if not (start <= dt <= end):
                    continue
            # conta la matricola
            counts[mat] = counts.get(mat, 0) + 1

        if not counts:
            return []

        # helper per trovare ragione sociale partendo da matricola
        def _find_client_for_matricola(matricola: str) -> str:
            m = matricola.strip()
            if m == "":
                return ""
            try:
                # se esiste la tabella macchine, usala per risalire al cliente
                if "macchine" in tables and "clienti" in tables:
                    cur.execute("PRAGMA table_info(macchine)")
                    mac_cols = [r[1].lower() for r in cur.fetchall()]
                    if "matricola" in mac_cols:
                        cur.execute(
                            "SELECT cliente_id FROM macchine WHERE trim(matricola)=? COLLATE NOCASE LIMIT 1",
                            (m,)
                        )
                        r = cur.fetchone()
                        if r and r[0] is not None:
                            cid = r[0]
                            # trova campo nome in clienti
                            cur.execute("PRAGMA table_info(clienti)")
                            cli_info = cur.fetchall()
                            mapping = {x[1].lower(): x[1] for x in cli_info}
                            for cand in ("rag_soc", "ragione_sociale", "ragionesociale", "nome"):
                                if cand in mapping:
                                    colreal = mapping[cand]
                                    cur.execute(f"SELECT {colreal} FROM clienti WHERE id=? LIMIT 1", (cid,))
                                    cres = cur.fetchone()
                                    if cres and cres[0]:
                                        return str(cres[0])
                # fallback: se clienti ha colonna matricola, cerca direttamente
                if "clienti" in tables:
                    cur.execute("PRAGMA table_info(clienti)")
                    cli_info = cur.fetchall()
                    cli_cols_lower = [x[1].lower() for x in cli_info]
                    mapping = {x[1].lower(): x[1] for x in cli_info}
                    if "matricola" in cli_cols_lower:
                        cur.execute(
                            "SELECT * FROM clienti WHERE trim(matricola)=? COLLATE NOCASE LIMIT 1",
                            (m,)
                        )
                        cres = cur.fetchone()
                        if cres:
                            # trova campo nome se esiste
                            for cand in ("rag_soc", "ragione_sociale", "ragionesociale", "nome"):
                                if cand in mapping:
                                    idx = list(mapping.keys()).index(cand)
                                    if idx < len(cres) and cres[idx]:
                                        return str(cres[idx])
            except Exception:
                # in caso di errore, fallback a stringa vuota
                return ""
            return ""

        # costruisci risultato: (nome_cliente, matricola, count)
        res: List[Tuple[str, str, int]] = []
        for mat, cnt in counts.items():
            name = _find_client_for_matricola(mat) or ""
            res.append((name, mat, int(cnt)))

        # ordina: per numero interventi discendente poi matricola
        res.sort(key=lambda x: (-x[2], x[1]))
        return res

    finally:
        conn.close()

def confronto_intae_misurart(join_on: str = "MATRICOLA") -> list:
    """
    Confronta le ultime date per matricola tra INTAE (DATA_INTERVENTO) e MISURART (DATAVP).
    Ritorna SOLO le matricole per le quali esistono date in entrambe le tabelle e le due date sono diverse.
    Output: lista di tuple (ragione_sociale, matricola, data_intae_iso, data_misurart_iso)
    """
    conn = get_connection()
    cur = conn.cursor()

    # raccogli mappe come prima (ultima data parsed per matricola)
    def _build_map(table: str, matricola_candidates, date_candidates):
        m = {}
        try:
            cur.execute("SELECT name FROM sqlite_master WHERE type='table'")
            tables = {r[0].lower() for r in cur.fetchall()}
            if table.lower() not in tables:
                return m
            cur.execute(f"PRAGMA table_info({table})")
            info = cur.fetchall()
            cols = [r[1] for r in info]
            cols_lower = [c.lower() for c in cols]
            # trova matricola e data
            mat_col = None
            for cand in matricola_candidates:
                if cand.lower() in cols_lower:
                    mat_col = cols[cols_lower.index(cand.lower())]
                    break
            dt_col = None
            for cand in date_candidates:
                if cand.lower() in cols_lower:
                    dt_col = cols[cols_lower.index(cand.lower())]
                    break
            if not mat_col or not dt_col:
                return m
            cur.execute(f"SELECT {mat_col}, {dt_col} FROM {table} WHERE {mat_col} IS NOT NULL AND TRIM({mat_col})<>''")
            for mat_raw, dt_raw in cur.fetchall():
                if mat_raw is None:
                    continue
                mat = str(mat_raw).strip()
                if mat == "":
                    continue
                dt = _try_parse_date(dt_raw)
                if not dt:
                    continue
                prev = m.get(mat)
                if prev is None or dt > prev:
                    m[mat] = dt
        except Exception:
            pass
        return m

    matricola_candidates = ["MATRICOLA", "matricola", "matricola_macchina", "mat"]
    intae_date_candidates = ["DATA_INTERVENTO", "data_intervento", "dataintervento", "data"]
    mis_date_candidates = ["DATAVP", "datavp", "data", "datavp"]

    intae_map = _build_map("intae", matricola_candidates, intae_date_candidates)
    mis_map = _build_map("misurart", matricola_candidates, mis_date_candidates)

    # helper client name (come prima, macchine->clienti fallback)
    cur.execute("SELECT name FROM sqlite_master WHERE type='table'")
    tables = {r[0].lower() for r in cur.fetchall()}

    def _client_name_for_matricola(matricola: str) -> str:
        mkey = (matricola or "").strip()
        if not mkey:
            return ""
        try:
            if "macchine" in tables and "clienti" in tables:
                cur.execute("PRAGMA table_info(macchine)")
                mac_cols = [r[1].lower() for r in cur.fetchall()]
                if "matricola" in mac_cols:
                    cur.execute("SELECT cliente_id FROM macchine WHERE trim(matricola)=? COLLATE NOCASE LIMIT 1", (mkey,))
                    row = cur.fetchone()
                    if row and row[0] is not None:
                        cid = row[0]
                        cur.execute("PRAGMA table_info(clienti)")
                        cli_info = cur.fetchall()
                        mapping = {x[1].lower(): x[1] for x in cli_info}
                        for cand in ("rag_soc", "ragione_sociale", "ragionesociale", "nome"):
                            if cand in mapping:
                                colreal = mapping[cand]
                                cur.execute(f"SELECT {colreal} FROM clienti WHERE id=? LIMIT 1", (cid,))
                                cres = cur.fetchone()
                                if cres and cres[0]:
                                    return str(cres[0])
            # fallback: clienti.matricola
            if "clienti" in tables:
                cur.execute("PRAGMA table_info(clienti)")
                cli_info = cur.fetchall()
                mapping = {x[1].lower(): x[1] for x in cli_info}
                if "matricola" in mapping:
                    cur.execute("SELECT * FROM clienti WHERE trim(matricola)=? COLLATE NOCASE LIMIT 1", (mkey,))
                    cres = cur.fetchone()
                    if cres:
                        for cand in ("rag_soc", "ragione_sociale", "ragionesociale", "nome"):
                            if cand in mapping:
                                idx = list(mapping.keys()).index(cand)
                                if idx < len(cres) and cres[idx]:
                                    return str(cres[idx])
        except Exception:
            return ""
        return ""

    # consideriamo solo matricole presenti in entrambe le mappe
    inter = set(intae_map.keys()) & set(mis_map.keys())
    out = []
    for k in sorted(inter):
        di = intae_map.get(k)
        dm = mis_map.get(k)
        if not di or not dm:
            continue
        if di == dm:
            # se sono uguali, NON le vogliamo
            continue
        # rappresentazioni ISO
        si = di.isoformat()
        sm = dm.isoformat()
        client_name = _client_name_for_matricola(k) or ""
        out.append((client_name, k, si, sm))

    conn.close()
    return out

# -------------------------
# Autosave quando modifico un campo
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
            from datetime import datetime as _dt
            return _dt.strptime(s.split()[0], fmt).date().isoformat()
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
