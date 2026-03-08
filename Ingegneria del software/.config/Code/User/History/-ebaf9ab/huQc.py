import sqlite3
import os
import csv
import tempfile
import subprocess
from datetime import datetime, timedelta
from typing import List, Tuple, Optional
from collections import defaultdict  
import re
import shutil
# openpyxl per leggere .xlsx
try:
    import openpyxl
except Exception:
    openpyxl = None  # gestiremo l'assenza più sotto

DB_FILE = "database.db"

# Colonne ESATTAMENTE come richieste (nomi preservati)
CLIENTI_COLUMNS = [
    "CODICE","RAG_SOC","P_IVA","IND_VIA","COMUNE","PROV","CAP","MF","MARCA","LOGOTIPO",
    "MATRICOLA","SCAD_ASS","IMP_ASS","VF","IMP_RIP","TEL","COMMENTI","PASSAGGIO","DATAINST",
    "TIPOMF","DATA","GRUPPO","DATADIS","CONTA","DM","CF","SCANNER","SCANNERAGG","N_PINPAD",
    "N_PIN_BACK","PEC","FINEMF","CODPUNTOCA","RTREADY","CODUNIVOCO","GARANZIA","AZZ","QRCODE",
    "LICENZA","PSW111","PSWSCDIG","MXAZ"
]

# --- mappatura di possibili varianti di header verso il nome canonico in CLIENTI_COLUMNS ---
CLIENTI_SYNONYMS = {
    # varianti comuni di TIPOMF
    "TIPORT": "TIPOMF",
    "TIPO_RT": "TIPOMF",
    "TIPO-MF": "TIPOMF",
    # varianti comuni di FINEMF
    "FINE_MF": "FINEMF",
    "FINE-MF": "FINEMF",
    # eventuali altre varianti che il fornitore potrebbe usare
    # "VECCHIOCAMPONAME": "NOMECANONICO",
}

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

# -------------------------
# DB connection & init
# -------------------------

def get_connection():
    # abilitiamo foreign keys per sicurezza
    conn = sqlite3.connect(DB_FILE)
    conn.execute("PRAGMA foreign_keys = ON")
    return conn


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

    # INT A E: create simple table if not exists, poi aggiungi colonne mancanti in modo sicuro
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
        sql = f"SELECT {cols_select} FROM clienti WHERE " + " OR ".join(where_parts) + " ORDER BY RAG_SOC LIMIT ?"
        cur.execute(sql, params)
    else:
        sql = f"SELECT {cols_select} FROM clienti ORDER BY RAG_SOC LIMIT ?"
        cur.execute(sql, (limit,))
    rows = cur.fetchall()
    conn.close()
    return cols, rows


def _normalize_key(k: str) -> str:
    # rimuove spazi e caratteri non alfanumerici, tutto uppercase
    if k is None:
        return ""
    return re.sub(r'[^A-Z0-9]', '', k.upper().strip())

def insert_cliente_from_map(data: dict) -> int:
    """Inserisce un cliente usando l'ordine stabile CLIENTI_COLUMNS.
    Normalizza le chiavi di `data` per tollerare spazi, minuscole, underscore ecc.
    """
    if not data:
        raise ValueError("data vuoto")

    # normalizza il dict di input (la versione normale delle chiavi -> valore)
    normalized_input = {_normalize_key(k): v for k, v in data.items()}

    # assicurati che la tabella e le colonne esistano
    conn = get_connection()
    cur = conn.cursor()
    cur.execute("PRAGMA table_info(clienti)")
    existing = [r[1] for r in cur.fetchall()]

    # aggiungi eventuali colonne mancanti dalla lista standard
    for c in CLIENTI_COLUMNS:
        if c not in existing:
            cur.execute(f'ALTER TABLE clienti ADD COLUMN "{c}" TEXT')

    # prepara valori: cerca la colonna normalizzata nella mappa normalizzata
    vals = []
    for c in CLIENTI_COLUMNS:
        norm = _normalize_key(c)
        v = normalized_input.get(norm, "")
        # se vuoi forzare stringa: v = "" if v is None else str(v)
        vals.append("" if v is None else str(v))

    colnames_sql = ", ".join([f'"{c}"' for c in CLIENTI_COLUMNS])
    placeholders = ", ".join(["?"] * len(CLIENTI_COLUMNS))
    sql = f'INSERT INTO clienti ({colnames_sql}) VALUES ({placeholders})'
    cur.execute(sql, vals)
    nid = cur.lastrowid
    conn.commit()
    conn.close()
    return nid

def update_cliente_by_id(rowid: int, updates: dict):
    if not updates:
        return
    cols = list(updates.keys())
    _ensure_clienti_columns(cols)
    parts = ", ".join([f'"{c}"=?' for c in cols])
    vals = [str(updates[c]) for c in cols] + [rowid]
    conn = get_connection()
    cur = conn.cursor()
    cur.execute(f'UPDATE clienti SET {parts} WHERE id=?', vals)
    conn.commit()
    conn.close()


def delete_cliente(rowid: int):
    conn = get_connection()
    cur = conn.cursor()
    cur.execute("DELETE FROM macchine WHERE cliente_id=?", (rowid,))
    cur.execute("DELETE FROM clienti WHERE id=?", (rowid,))
    conn.commit()
    conn.close()


def duplicate_cliente(rowid: int, new_suffix: str = "_COPIA") -> Optional[int]:
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
        if col == "RAG_SOC" and val:
            val = f"{val}{new_suffix}"
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

def import_clients_from_csv_generic(path: str, delimiter: str = ";", encoding: str = "latin-1") -> int:
    import difflib, re
    if not os.path.exists(path):
        raise FileNotFoundError(path)

    # sniffing del delimitatore (legge un campione)
    with open(path, "r", encoding=encoding, errors="replace") as sf:
        sample = sf.read(8192)
        sf.seek(0)
        if delimiter is None or delimiter == "":
            try:
                sniffer = csv.Sniffer()
                dialect = sniffer.sniff(sample, delimiters=",;|\t")
                used_delim = dialect.delimiter
            except Exception:
                # fallback: prova i delimitatori più comuni
                for d in [",", ";", "\t", "|"]:
                    if sample.count(d) >= 1:
                        used_delim = d
                        break
                else:
                    used_delim = ","
        else:
            # se il delimiter passato sembra sbagliato (prima riga non separa), proviamo a sniffare
            first_line = sample.splitlines()[0] if sample else ""
            if first_line.count(delimiter) == 0 and (first_line.count(",") + first_line.count(";") + first_line.count("\t")) > 0:
                try:
                    sniffer = csv.Sniffer()
                    dialect = sniffer.sniff(sample, delimiters=",;|\t")
                    used_delim = dialect.delimiter
                except Exception:
                    used_delim = delimiter
            else:
                used_delim = delimiter

    # leggiamo tutto con DictReader (possibile costo in memoria per file grandi, ma più semplice)
    with open(path, newline='', encoding=encoding, errors='replace') as f:
        reader = csv.DictReader(f, delimiter=used_delim)
        # Header originali così come sono nel file (es. "TIPOMF,C,3")
        headers = reader.fieldnames or []

        # Header puliti (mantiene solo la parte prima della prima virgola)
        # "TIPOMF,C,3" → "TIPOMF"
        # "FINEMF,D" → "FINEMF"
        headers_clean = [h.split(",")[0].strip() for h in headers]

        if not headers:
            raise ValueError("CSV privo di intestazioni o delimiter non riconosciuto")
        rows = list(reader)  # elenco di dict delle righe lette

    # normalizzatore per stringhe: rimuove tutto tranne ascii alfanum, uppercase
    def normalize(s: str) -> str:
        if s is None:
            return ""
        return re.sub(r'[^0-9A-Za-z]', '', str(s)).upper()

        # costruzione mappa normalized -> canonical usando CLIENTI_SYNONYMS
    syn_norm = { normalize(k): v for k, v in CLIENTI_SYNONYMS.items() }
    
    # mapping header -> colonna CLIENTI_COLUMNS con fuzzy matching
    map_header_to_col = {}
    wanted = CLIENTI_COLUMNS.copy()
    wanted_norm = {wc: normalize(wc) for wc in wanted}

    # ==== NUOVO BLOCCO DI MAPPING: usa headers_clean per il matching ma conserva headers originali ====
    for h_orig, h_clean in zip(headers, headers_clean):
        hn = normalize(h_clean)  # es. "TIPOMF"
        matched = None

        # 1) match diretto / contains (preferibile)
        for wc, wn in wanted_norm.items():
            if hn == wn or wn in hn or hn in wn:
                matched = wc
                break
            
        # 1b) controllo synonyms (se la header normalizzata è in syn_norm -> mappa)
        if not matched and hn in syn_norm:
            matched = syn_norm[hn]

        # 2) fuzzy matching fallback
        if not matched:
            scores = [(difflib.SequenceMatcher(None, hn, wn).ratio(), wc) for wc, wn in wanted_norm.items()]
            best_score, best_wc = max(scores, key=lambda x: x[0])
            if best_score >= 0.60:  # soglia ragionevole
                matched = best_wc

        if matched:
            map_header_to_col[h_orig] = matched

    # (opzionale) debug: vedi quale mapping è stato trovato
    # print("map_header_to_col:", map_header_to_col)

    # Fallback posizionale (solo per campi critici) se ancora mancanti:
    crucial = ["RAG_SOC", "P_IVA", "IND_VIA"]
    mapped_values = set(map_header_to_col.values())
    missing_crucial = [c for c in crucial if c not in mapped_values]

    if missing_crucial and rows:
        # guardiamo la prima riga di dati per capire se l'ordine è standard
        first_row = rows[0]
        # ricaviamo i valori nella stessa sequenza di headers
        first_vals = [first_row.get(h, "") for h in headers]
        # se abbiamo almeno 4 colonne, proviamo a mappare per posizione: 0->CODICE,1->RAG_SOC,2->P_IVA,3->IND_VIA
        if len(headers) >= 4:
            pos_map = {0: "CODICE", 1: "RAG_SOC", 2: "P_IVA", 3: "IND_VIA"}
            for idx, col_name in pos_map.items():
                hdr = headers[idx]
                if hdr not in map_header_to_col and col_name in missing_crucial:
                    # semplice check: il valore della riga 0 nella colonna 'third' sembra numerico per P_IVA?
                    val = first_vals[idx] if idx < len(first_vals) else ""
                    if col_name == "P_IVA":
                        # se sembra un numero o contiene cifre, ok
                        if any(ch.isdigit() for ch in str(val)):
                            map_header_to_col[hdr] = col_name
                    else:
                        # per ragione sociale/indirizzo non facciamo controllo stringente
                        map_header_to_col[hdr] = col_name

    # Assicurati che tutte le colonne standard esistano nella tabella
    conn = get_connection()
    cur = conn.cursor()
    cur.execute("PRAGMA table_info(clienti)")
    existing = [r[1] for r in cur.fetchall()]
    for c in CLIENTI_COLUMNS:
        if c not in existing:
            cur.execute(f'ALTER TABLE clienti ADD COLUMN "{c}" TEXT')

    # Preparazione insert: usa l'ordine fisso CLIENTI_COLUMNS
    insert_cols = CLIENTI_COLUMNS
    placeholders = ", ".join(["?"] * len(insert_cols))
    colnames_sql = ", ".join([f'"{c}"' for c in insert_cols])
    sql = f'INSERT INTO clienti ({colnames_sql}) VALUES ({placeholders})'
    to_insert = []
    for r in rows:
        vals = []
        for c in insert_cols:
            # cerchiamo se esiste un header che mappa a c
            found = False
            for orig_h, mapped in map_header_to_col.items():
                if mapped == c:
                    v = r.get(orig_h, "")
                    if v is None:
                        v = ""
                    vals.append(str(v))
                    found = True
                    break
            if not found:
                vals.append("")
        to_insert.append(vals)

    if to_insert:
        cur.executemany(sql, to_insert)
    conn.commit()
    n = cur.rowcount if cur.rowcount is not None else len(to_insert)
    conn.close()
    return n

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

def update_scadenze_clienti():
    conn = get_connection()
    cur = conn.cursor()

    cur.execute("""
        SELECT id, DATAINST, GARANZIA, FINEMF, MXAZ
        FROM clienti
    """)
    rows = cur.fetchall()

    for cid, data_inst, garanzia, finemf, mxaz in rows:
        
        # --- Calcolo scadenza garanzia
        nuova_gar = calcola_scadenza_garanzia(data_inst, garanzia)

        # --- Calcolo fine memoria fiscale
        nuova_finemf = finemf
        if not nuova_finemf:
            nuova_finemf = calcola_fine_memoria(data_inst, mxaz)

        cur.execute("""
            UPDATE clienti
            SET SCAD_ASS=?, FINEMF=?
            WHERE id=?
        """, (nuova_gar or "", nuova_finemf or "", cid))

    conn.commit()
    conn.close()

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
                # prepariamo records
                records = []
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
                    records.append(tuple(rec))
                # inserimento
                conn = get_connection()
                cur = conn.cursor()
                colnames_sql = ", ".join([f'"{c}"' for c in INTAE_COLUMNS])
                placeholders = ", ".join(["?"] * len(INTAE_COLUMNS))
                if records:
                    cur.executemany(f'INSERT INTO intae ({colnames_sql}) VALUES ({placeholders})', records)
                conn.commit()
                n = cur.rowcount if cur.rowcount is not None else len(records)
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
    conn = get_connection()
    cur = conn.cursor()
    colnames = ", ".join([f'"{c}"' for c in INTAE_COLUMNS])
    # cerchiamo matricola (case-insensitive)
    cur.execute(f"""
        SELECT id, {colnames}
        FROM intae
        WHERE lower(MATRICOLA)=? OR MATRICE LIKE ?
        ORDER BY DATA_INTERVENTO DESC
    """, (matricola.lower(), f"%{matricola}%"))
    rows = cur.fetchall()
    conn.close()
    return rows

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
    cols = list(updates.keys())
    # aggiungi colonne se necessarie (sicurezza)
    _ensure_clienti_columns(cols)  # la funzione aggiunge colonne, ma funziona per qualsiasi tabella? -> useremo ALTER manuale
    # costruzione SQL
    parts = ", ".join([f'"{c}"=?' for c in cols])
    vals = [str(updates[c]) for c in cols] + [record_id]
    conn = get_connection()
    cur = conn.cursor()
    cur.execute(f'UPDATE misurart SET {parts} WHERE id=?', vals)
    conn.commit()
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
        
def reset_VF_if_tipo_2():
    """
    Se un intervento misurart ha TIPOINT = '2',
    allora azzera VF delle macchine con la stessa matricola.
    """
    conn = get_connection()
    cur = conn.cursor()

    # prende tutte le matricole con interventi tipo 2
    cur.execute("SELECT MATRICOLA FROM misurart WHERE TIPOINT='2'")
    mats = [r[0] for r in cur.fetchall() if r[0]]

    for matricola in mats:
        cur.execute("UPDATE macchine SET vf='' WHERE matricola=?", (matricola,))

    conn.commit()
    conn.close()

# -------------------------
# Report: calcolo fine memoria (MXAZ interpreted as days)
# -------------------------

def _try_parse_date(s: Optional[str]) -> Optional[datetime]:
    if not s:
        return None
    s = str(s).strip().split()[0]
    fmts = ["%Y-%m-%d","%d-%m-%Y","%d/%m/%Y","%Y%m%d","%d%m%Y","%Y/%m/%d","%d.%m.%Y"]
    for f in fmts:
        try:
            return datetime.strptime(s, f)
        except Exception:
            pass
    return None


def calcola_fine_memoria(data_installazione: str, mxaz) -> Optional[str]:
    """
    Calcola la data di fine memoria fiscale: data_installazione + mxaz giorni.
    Ritorna stringa YYYY-MM-DD o None se non calcolabile.
    """
    if not data_installazione or mxaz is None or str(mxaz).strip() == "":
        return None
    # proviamo a interpretare mxaz come numero di giorni
    try:
        days = int(str(mxaz))
    except Exception:
        return None
    dt = _try_parse_date(data_installazione)
    if not dt:
        return None
    fin = dt + timedelta(days=days)
    return fin.strftime("%Y-%m-%d")

def _coerce_days(mxaz):
    """Prova a ottenere un numero di giorni da mxaz in modo permissivo."""
    if mxaz is None:
        return None
    s = str(mxaz).strip()
    if s == "":
        return None
    # se è già intero
    try:
        return int(s)
    except Exception:
        pass
    # estrai eventuali cifre
    digits = "".join(ch for ch in s if ch.isdigit())
    if digits:
        try:
            return int(digits)
        except Exception:
            pass
    # se contiene 'M' o 'ANN' interpretiamo come mesi (moltiplica 30)
    u = s.upper()
    if 'M' in u and any(ch.isdigit() for ch in s):
        # es. "12M" -> 12 mesi -> 360 giorni
        num = "".join(ch for ch in s if ch.isdigit())
        if num:
            try:
                return int(num) * 30
            except Exception:
                pass
    # fallback: non interpretabile
    return None


def populate_finemf_auto(verbose: bool = True, force: bool = False) -> int:
    """
    Popola clienti.FINEMF calcolando: macchine.data_installazione + mxaz (giorni).
    - Se FINEMF è già presente non la sovrascrive a meno che force=True.
    - Permissivo su mxaz: estrae cifre o interpreta 'M' come mesi.
    - Ritorna numero di clienti modificati.
    """
    conn = get_connection()
    cur = conn.cursor()

    updated_clients = set()
    errors = []

    # prendi tutte le macchine con cliente collegato
    try:
        cur.execute("SELECT id, matricola, data_installazione, mxaz, cliente_id FROM macchine WHERE cliente_id IS NOT NULL")
    except Exception as e:
        if verbose:
            print("populate_finemf_auto: errore SELECT macchine:", e)
        conn.close()
        return 0

    rows = cur.fetchall()
    if verbose:
        print(f"populate_finemf_auto: trovate {len(rows)} macchine con cliente_id non null")

    for mid, matricola, data_install, mxaz_machine, cliente_id in rows:
        if not cliente_id:
            continue

        # leggi FINEMF cliente
        try:
            cur.execute("SELECT FINEMF FROM clienti WHERE id=?", (cliente_id,))
            rf = cur.fetchone()
            finemf_cliente = rf[0] if rf and rf[0] is not None else None
        except Exception as e:
            errors.append((cliente_id, "SELECT FINEMF error", str(e)))
            continue

        if finemf_cliente and str(finemf_cliente).strip() != "" and not force:
            # già presente e non forzato -> salta
            continue

        # scegli mxaz: preferisci mxaz dalla macchina, altrimenti da cliente
        chosen_mxaz = mxaz_machine
        if not chosen_mxaz or str(chosen_mxaz).strip() == "":
            cur.execute("SELECT MXAZ FROM clienti WHERE id=?", (cliente_id,))
            r2 = cur.fetchone()
            chosen_mxaz = r2[0] if r2 and r2[0] is not None else None

        days = _coerce_days(chosen_mxaz)
        if days is None:
            errors.append((cliente_id, "mxaz non interpretabile", chosen_mxaz))
            continue

        # calcola fine memoria
        fin = calcola_fine_memoria(data_install, days)
        if not fin:
            errors.append((cliente_id, "data_install non parsabile o mancante", data_install))
            continue

        try:
            cur.execute("UPDATE clienti SET FINEMF=? WHERE id=?", (fin, cliente_id))
            updated_clients.add(cliente_id)
            if verbose:
                print(f"populate_finemf_auto: cliente {cliente_id} -> FINEMF={fin} (matr={matricola})")
        except Exception as e:
            errors.append((cliente_id, "UPDATE error", str(e)))
            continue

    # ------- Passata aggiuntiva: prova a popolare FINEMF direttamente dalla tabella clienti -------
    try:
        cur.execute("SELECT id, DATAINST, FINEMF, MXAZ FROM clienti")
        for cid, data_inst_cli, finemf_cli, mxaz_cli in cur.fetchall():
            # se FINEMF già presente e non forzato, salta
            if finemf_cli and str(finemf_cli).strip() != "" and not force:
                continue
            # se non abbiamo data di installazione, non possiamo calcolare
            if not data_inst_cli or str(data_inst_cli).strip() == "":
                continue
            # tenta di interpretare mxaz
            days_cli = _coerce_days(mxaz_cli)
            if days_cli is None:
                # non possiamo calcolare senza mxaz interpretabile
                continue
            fin_cli = calcola_fine_memoria(data_inst_cli, days_cli)
            if not fin_cli:
                # data non parsabile
                continue
            try:
                cur.execute("UPDATE clienti SET FINEMF=? WHERE id=?", (fin_cli, cid))
                updated_clients.add(cid)
                if verbose:
                    print(f"populate_finemf_auto (clienti): cliente {cid} -> FINEMF={fin_cli}")
            except Exception as e:
                errors.append((cid, "UPDATE client finemf error", str(e)))
    except Exception as e:
        if verbose:
            print("populate_finemf_auto: errore passata clienti:", e)

    conn.commit()
    conn.close()

    if verbose:
        print(f"populate_finemf_auto: aggiornati {len(updated_clients)} clienti")
        if errors:
            print("populate_finemf_auto: errori/parziali (max 20 mostrati):")
            for e in errors[:20]:
                print("  ", e)

    return len(updated_clients)

def report_verifiche() -> List[Tuple]:
    """
    Report principale: per ogni macchina restituisce
    (RAG_SOC, COMUNE, MATRICOLA, VF, DATA_INSTALL, ULTIMA_VERIFICA, SCAD_ASS, FINEMF)
    """
    conn = get_connection()
    cur = conn.cursor()
    cur.execute("""
        SELECT m.id, m.matricola, m.data_installazione, m.mxaz, m.vf, c.RAG_SOC, c.COMUNE, c.P_IVA
        FROM macchine m
        LEFT JOIN clienti c ON m.cliente_id = c.id
    """)
    machines = cur.fetchall()
    result = []
    for m in machines:
        mac_id, matricola, data_install, mxaz, vf, rag, comune, piva = m

        # se matricola mancante skip (o puoi decidere di includere comunque)
        if not matricola:
            ultima_ver = None
        else:
            # prendi ultima data utile: preferisci DATAFINE poi DATAINIZIO
            cur.execute("""
                SELECT COALESCE(DATAFINE, DATAINIZIO)
                FROM misurart
                WHERE MATRICOLA = ? OR MATRICOLA LIKE ?
                ORDER BY COALESCE(DATAFINE, DATAINIZIO) DESC
                LIMIT 1
            """, (matricola, f"%{matricola}%"))
            r = cur.fetchone()
            ultima_ver = r[0] if r and r[0] else None

        # scadenza assistenza (eventuale campo nella tabella macchine)
        scad_ass = None
        cur.execute("SELECT scadenza_assistenza FROM macchine WHERE id=?", (mac_id,))
        rr = cur.fetchone()
        if rr and rr[0]:
            scad_ass = rr[0]

        # fine memoria: prefer FINEMF in clienti, poi calcolo da data_install + mxaz
        finemf = None
        cur.execute("SELECT FINEMF FROM clienti WHERE id=(SELECT cliente_id FROM macchine WHERE id=?)", (mac_id,))
        rf = cur.fetchone()
        if rf and rf[0]:
            finemf = rf[0]
        if not finemf and data_install:
            finemf = calcola_fine_memoria(data_install, mxaz)

        result.append((rag or "", comune or "", matricola or "", vf or "", data_install or "", ultima_ver or "", scad_ass or "", finemf or ""))
    conn.close()
    return result


def report_verifiche_in_scadenza(within_days: int = 30) -> List[Tuple]:
    """
    Restituisce le macchine la cui prossima scadenza di verifica (ultima_verifica + 365gg)
    cade entro 'within_days' giorni a partire da oggi.
    Ritorna tuple: (RAG_SOC, MATRICOLA, DATA_ULTIMA_VERIFICA, DATA_SCADENZA, GIORNI_RESIDUI)
    Nota: include solo quelle con 0 <= giorni_residui <= within_days.
    """
    conn = get_connection()
    cur = conn.cursor()
    today = datetime.now().date()

    cur.execute("""
        SELECT m.id, m.matricola, c.RAG_SOC
        FROM macchine m
        LEFT JOIN clienti c ON c.id = m.cliente_id
    """)

    results = []

    for mid, matricola, rag in cur.fetchall():
        if not matricola:
            continue

        # prendi ultima data utile (DATAFINE o DATAINIZIO)
        cur.execute("""
            SELECT COALESCE(DATAFINE, DATAINIZIO)
            FROM misurart
            WHERE MATRICOLA = ? OR MATRICOLA LIKE ?
            ORDER BY COALESCE(DATAFINE, DATAINIZIO) DESC
            LIMIT 1
        """, (matricola, f"%{matricola}%"))
        r = cur.fetchone()

        if not r or not r[0]:
            continue

        dt_obj = _try_parse_date(r[0])
        if not dt_obj:
            continue

        # calcolo scadenza (qui uso +365 giorni, adattabile)
        scadenza_date = (dt_obj + timedelta(days=365)).date()
        giorni = (scadenza_date - today).days

        # includi solo se tra 0 e within_days (=> prossime scadenze, non già scadute)
        if 0 <= giorni <= within_days:
            results.append((rag or "", matricola, r[0], scadenza_date.strftime("%Y-%m-%d"), giorni))

    conn.close()
    return results

# -------------------------
# Bootstrap helper
# -------------------------

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
    
    # popola FINEMF automaticamente (log su console)
    try:
        n = populate_finemf_auto(verbose=True, force=False)
        print(f"ensure_db_exists: populate_finemf_auto ha aggiornato {n} clienti")
    except Exception as e:
        print("ensure_db_exists: errore in populate_finemf_auto:", e)

def _parse_garanzia_to_months(garanzia_str) -> Optional[int]:
    """
    Interpreta GARANZIA:
    - se è un intero -> mesi
    - se contiene 'Y' o 'ANNI' -> moltiplica per 12
    - se contiene 'M' o 'MESI' -> usa il numero
    - altrimenti None
    """
    if not garanzia_str:
        return None
    s = str(garanzia_str).strip().upper()
    # solo numero -> mesi
    try:
        if s.isdigit():
            return int(s)
        # es. "1Y", "2ANNI"
        if 'Y' in s or 'ANNO' in s or 'ANNI' in s:
            num = ''.join(ch for ch in s if ch.isdigit())
            if num:
                return int(num) * 12
        if 'M' in s or 'MES' in s:
            num = ''.join(ch for ch in s if ch.isdigit())
            if num:
                return int(num)
    except Exception:
        return None
    return None

def calcola_scadenza_garanzia(data_installazione: str, garanzia=None) -> Optional[str]:
    """
    Calcola la data di fine garanzia: esattamente 1 anno (365 giorni)
    dopo la data di installazione.
    Restituisce stringa YYYY-MM-DD o None se non calcolabile.
    Nota: il parametro 'garanzia' è mantenuto per compatibilità ma non influisce sul calcolo.
    """
    if not data_installazione:
        return None
    dt = _try_parse_date(data_installazione)
    if not dt:
        return None
    fin = dt + timedelta(days=365)
    return fin.strftime("%Y-%m-%d")

def populate_garanzia_auto(verbose: bool = True, force: bool = False) -> int:
    """
    Popola automaticamente clienti.GARANZIA (default '12' mesi) e clienti.SCAD_ASS
    (calcolata come DATA_INSTALLAZIONE + 365 giorni) per i clienti che non hanno
    questi campi valorizzati. Se force=True sovrascrive i valori esistenti.
    Ritorna il numero di clienti aggiornati.
    """
    conn = get_connection()
    cur = conn.cursor()

    updated = 0
    errors = []

    # leggiamo tutti i clienti
    try:
        cur.execute("SELECT id, GARANZIA, SCAD_ASS, DATAINST FROM clienti")
        clienti = cur.fetchall()
    except Exception as e:
        if verbose:
            print("populate_garanzia_auto: errore SELECT clienti:", e)
        conn.close()
        return 0

    if verbose:
        print(f"populate_garanzia_auto: {len(clienti)} clienti trovati")

    for cid, gar, scad_ass, client_datainst in clienti:
        need_gar = (not gar or str(gar).strip() == "")
        need_scad = (not scad_ass or str(scad_ass).strip() == "")

        if not need_gar and not need_scad and not force:
            # nulla da fare
            continue

        # tenta di ottenere data di installazione da macchine collegate (se esiste)
        dt_install = None
        try:
            cur.execute("""
                SELECT data_installazione
                FROM macchine
                WHERE cliente_id = ? AND data_installazione IS NOT NULL AND trim(data_installazione) <> ''
                ORDER BY data_installazione ASC
                LIMIT 1
            """, (cid,))
            r = cur.fetchone()
            if r and r[0]:
                dt_install = r[0]
        except Exception as e:
            # non blocchiamo l'intero processo: proseguiamo con client_datainst
            if verbose:
                print(f"populate_garanzia_auto: errore ricerca macchine per cliente {cid}: {e}")

        # fallback: usa DATAINST nella tabella clienti
        if not dt_install and client_datainst and str(client_datainst).strip() != "":
            dt_install = client_datainst

        # se ancora non abbiamo data_installazione, non possiamo calcolare scadenza
        scad_val = None
        if dt_install:
            parsed = _try_parse_date(dt_install)
            if parsed:
                scad_val = (parsed + timedelta(days=365)).strftime("%Y-%m-%d")
            else:
                if verbose:
                    print(f"populate_garanzia_auto: data_install non parsabile per cliente {cid}: '{dt_install}'")

        # costruiamo l'update se necessario
        updates = {}
        if need_gar or force:
            # imposto GARANZIA a '12' (mesi) se mancante
            if need_gar or force:
                updates["GARANZIA"] = "12"

        if (need_scad or force) and scad_val:
            updates["SCAD_ASS"] = scad_val

        if updates:
            # assicurati colonne esistano (di solito già lo sono)
            cols_to_ensure = list(updates.keys())
            _ensure_clienti_columns(cols_to_ensure)
            # esegui update
            parts = ", ".join([f'"{k}" = ?' for k in updates.keys()])
            vals = [updates[k] for k in updates.keys()] + [cid]
            try:
                cur.execute(f'UPDATE clienti SET {parts} WHERE id=?', vals)
                updated += 1
                if verbose:
                    print(f"populate_garanzia_auto: aggiornato cliente id={cid}: {updates}")
            except Exception as e:
                errors.append((cid, str(e)))
                if verbose:
                    print(f"populate_garanzia_auto: errore update cliente {cid}: {e}")

    conn.commit()
    conn.close()

    if verbose:
        print(f"populate_garanzia_auto: clienti aggiornati = {updated}")
        if errors:
            print("populate_garanzia_auto: errori (max 20 mostrati):")
            for e in errors[:20]:
                print("  ", e)

    return updated

def report_memorie_in_scadenza(within_days: int = 30) -> List[Tuple]:
    """
    Restituisce macchine la cui 'fine memoria fiscale' (FINEMF) è entro 'within_days' giorni a partire da oggi.
    Usa in ordine: clienti.FINEMF, macchine.FINEMF (se presente), altrimenti calcola da data_installazione + mxaz.
    Ritorna lista di tuple: (cliente_rag, matricola, data_install, finemf, giorni_residui, vf, id_macchina)
    """
    conn = get_connection()
    cur = conn.cursor()
    today = datetime.now().date()
    cur.execute("""
        SELECT m.id, m.matricola, m.data_installazione, m.mxaz, m.vf, c.RAG_SOC, c.FINEMF
        FROM macchine m
        LEFT JOIN clienti c ON m.cliente_id = c.id
    """)
    out = []
    for row in cur.fetchall():
        mid, matricola, data_install, mxaz, vf, rag, finemf_cliente = row
        # prefer FINEMF del cliente, poi (se disponibile) calcolo
        finemf = finemf_cliente
        if not finemf:
            # prova FINEMF in macchine (colonna non garantita ma la cerchiamo)
            cur.execute("PRAGMA table_info(macchine)")
            mac_cols = [r[1] for r in cur.fetchall()]
            if "FINEMF" in mac_cols:
                cur.execute("SELECT FINEMF FROM macchine WHERE id=?", (mid,))
                r2 = cur.fetchone()
                if r2 and r2[0]:
                    finemf = r2[0]
        if not finemf and data_install:
            finemf = calcola_fine_memoria(data_install, mxaz)
        if not finemf:
            continue
        dt = _try_parse_date(finemf)
        if not dt:
            continue
        days_left = (dt.date() - today).days
        if days_left <= within_days:
            out.append((rag or "", matricola or "", data_install or "", finemf, days_left, vf or "", mid))
    conn.close()
    return out


def report_rt_per_model() -> List[Tuple[str,int,int]]:
    """
    Costruisce un report: per ogni modello conta:
      - totale registrazioni MISURART per quel modello
      - quanti 'RT' (se trovabili) oppure conta su RTREADY in clienti per macchine di quel modello
    Ritorna lista di tuple: (model, count_total, count_rt)
    """
    conn = get_connection()
    cur = conn.cursor()
    # Prima: conta per modello in misurart
    cur.execute("SELECT MODELLO, COUNT(*) FROM misurart GROUP BY MODELLO")
    totals = {r[0] or "": r[1] for r in cur.fetchall()}
    # Conta RT provando a trovare campo ESITO che contenga 'RT' oppure cerco 'RT' in other columns
    cur.execute("PRAGMA table_info(misurart)")
    mis_cols = [r[1] for r in cur.fetchall()]
    count_rt = defaultdict(int)
    if "ESITO" in mis_cols:
        cur.execute("SELECT MODELLO, ESITO FROM misurart")
        for mod, es in cur.fetchall():
            if es and "RT" in es.upper():
                count_rt[mod or ""] += 1
    # fallback: usa RTREADY in clienti+macchine join
    cur.execute("""
        SELECT m.modello, COUNT(c.RTREADY)
        FROM macchine m
        LEFT JOIN clienti c ON m.cliente_id = c.id
        WHERE c.RTREADY IS NOT NULL AND c.RTREADY <> ''
        GROUP BY m.modello
    """)
    for mod, cnt in cur.fetchall():
        count_rt[mod or ""] += cnt or 0
    conn.close()
    # unisci i set di modelli
    keys = set(totals.keys()) | set(count_rt.keys())
    res = []
    for k in sorted(keys):
        res.append((k, totals.get(k, 0), count_rt.get(k, 0)))
    return res


def disattiva_macchina_by_matricola(matricola: str, operatore: str = "SISTEMA") -> bool:
    """
    Esegue disattivazione:
     - crea una riga in operazioni con tipo 'DIS' codificato come '2' come richiesto
     - aggiorna macchine.disattivazione con la data odierna e azzera campo vf
     - se la macchina è collegata a un cliente, aggiorna il campo clienti.DATADIS con la data odierna
    Ritorna True se trovata e aggiornata almeno una macchina, False altrimenti.
    """
    if not matricola:
        return False
    conn = get_connection()
    cur = conn.cursor()
    # trova macchina
    cur.execute("SELECT id, cliente_id FROM macchine WHERE lower(matricola)=? OR matricola LIKE ? COLLATE NOCASE", (matricola.lower(), f"%{matricola}%"))
    r = cur.fetchone()
    if not r:
        conn.close()
        return False
    mid, cliente_id = r
    now_s = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    # inserisci operazione (tipo '2', codice 2)
    cur.execute("INSERT INTO operazioni (macchina_id, data_intervento, tipo, codice, note, operatore) VALUES (?,?,?,?,?,?)",
                (mid, now_s, "DIS", 2, "Disattivazione imputata via procedura", operatore))
    # aggiorna macchine
    cur.execute("UPDATE macchine SET disattivazione=?, vf=? WHERE id=?", (now_s, "", mid))
    # aggiorna cliente.DATADIS se esiste
    if cliente_id:
        # assicurati che la colonna DATADIS esista
        cur.execute("PRAGMA table_info(clienti)")
        cols = [r2[1] for r2 in cur.fetchall()]
        if "DATADIS" in cols:
            cur.execute("UPDATE clienti SET DATADIS=? WHERE id=?", (now_s.split()[0], cliente_id))
    conn.commit()
    conn.close()
    return True


def report_checklist_from_database_clienti() -> List[Tuple]:
    """
    Restituisce tutte le righe e colonne della tabella 'clienti'.
    Serve come checklist dei campi realmente presenti nel database clienti.
    """
    conn = get_connection()
    cur = conn.cursor()

    # lista tabelle
    cur.execute("SELECT name FROM sqlite_master WHERE type='table'")
    tables = [r[0].lower() for r in cur.fetchall()]

    # la tabella corretta è 'clienti'
    if "clienti" not in tables:
        conn.close()
        return []

    # leggo tutte le righe
    cur.execute("SELECT * FROM clienti")
    rows = cur.fetchall()

    # leggo le colonne
    cur.execute("PRAGMA table_info(clienti)")
    cols = [r[1] for r in cur.fetchall()]

    conn.close()
    return [(cols, rows)]



def report_cambi_dgfe(data_from: str, data_to: str) -> List[Tuple]:
    """
    Restituisce le operazioni/registrazioni che menzionano 'dgfe' (case-insensitive)
    tra due date (data formato libero: f-parse con _try_parse_date).
    Se non trova colonne con DGFE esplicito, cerca nei campi note di operazioni e in MISURART.* campi testuali.
    Ritorna lista di tuple (sorgente, id, data, campo, valore)
    """
    d1 = _try_parse_date(data_from)
    d2 = _try_parse_date(data_to)
    if not d1 or not d2:
        return []
    start = d1
    end = d2 + timedelta(days=1)
    conn = get_connection()
    cur = conn.cursor()
    out = []
    # cerca in operazioni.note
    cur.execute("SELECT id, macchina_id, data_intervento, note FROM operazioni WHERE note IS NOT NULL")
    for oid, mid, dt_s, note in cur.fetchall():
        if not note:
            continue
        if "dgfe" in note.lower():
            dt_obj = _try_parse_date(dt_s)
            if not dt_obj:
                continue
            if start <= dt_obj <= end:
                out.append(("operazioni", oid, dt_s, "note", note))
    # cerca in misurart: tutte le colonne testuali
    cur.execute("PRAGMA table_info(misurart)")
    mis_cols = [r[1] for r in cur.fetchall()]
    text_cols = mis_cols  # sono tutte TEXT per come definiamo
    cur.execute(f"SELECT id, {', '.join(text_cols)} FROM misurart")
    for row in cur.fetchall():
        rid = row[0]
        values = row[1:]
        # proviamo a estrarre la data di riferimento: DATAINIZIO/DATAFINE
        # controllo semplice: se qualunque campo contiene 'dgfe'
        joined = " ".join([str(v or "") for v in values]).lower()
        if "dgfe" in joined:
            # proviamo a prendere DATAVP/DATAINIZIO come data se esiste
            date_candidate = None
            try:
                if "DATAINIZIO" in mis_cols:
                    idx = mis_cols.index("DATAINIZIO")
                    date_candidate = _try_parse_date(values[idx])
                if not date_candidate and "DATAFINE" in mis_cols:
                    idx2 = mis_cols.index("DATAFINE")
                    date_candidate = _try_parse_date(values[idx2])
            except Exception:
                date_candidate = None
            if not date_candidate:
                # fallback: includi senza controllo data
                out.append(("misurart", rid, "", "multiple", "contains dgfe"))
            else:
                if start <= date_candidate <= end:
                    out.append(("misurart", rid, date_candidate.strftime("%Y-%m-%d"), "multiple", "contains dgfe"))
    conn.close()
    return out


def confronto_intae_misurart(join_on: str = "MATRICOLA") -> List[Tuple]:
    """
    Confronta dati tra intae e misurart usando la colonna join_on (default MATRICO LA).
    Restituisce elenco di tuple (valore_join, count_intae, count_misurart).
    Serve come controllo rapido di divergenze.
    """
    conn = get_connection()
    cur = conn.cursor()
    # verifica che le colonne esistano
    cur.execute("PRAGMA table_info(intae)")
    in_cols = [r[1] for r in cur.fetchall()]
    cur.execute("PRAGMA table_info(misurart)")
    mi_cols = [r[1] for r in cur.fetchall()]

    if join_on not in in_cols and join_on not in mi_cols:
        conn.close()
        return []

    # conta per valore in intae
    cur.execute(f"SELECT {join_on}, COUNT(*) FROM intae WHERE {join_on} IS NOT NULL AND {join_on}<>'' GROUP BY {join_on}")
    counts_intae = {r[0]: r[1] for r in cur.fetchall()}
    cur.execute(f"SELECT {join_on}, COUNT(*) FROM misurart WHERE {join_on} IS NOT NULL AND {join_on}<>'' GROUP BY {join_on}")
    counts_mis = {r[0]: r[1] for r in cur.fetchall()}
    keys = set(counts_intae.keys()) | set(counts_mis.keys())
    out = []
    for k in sorted(keys):
        out.append((k, counts_intae.get(k, 0), counts_mis.get(k, 0)))
    conn.close()
    return out