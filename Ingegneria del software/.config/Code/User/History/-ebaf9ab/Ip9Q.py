import sqlite3
import os
import csv
import tempfile
import subprocess
from datetime import datetime, timedelta, date
from typing import List, Tuple, Optional
from collections import defaultdict  # aggiungi in cima del file se non presente

DB_FILE = "database.db"

# Colonne ESATTAMENTE come richieste (nomi preservati)
CLIENTI_COLUMNS = [
    "CODICE","RAG_SOC","P_IVA","IND_VIA","COMUNE","PROV","CAP","MF","MARCA","LOGOTIPO",
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
        sql = f"SELECT {cols_select} FROM clienti WHERE " + " OR ".join(where_parts) + " ORDER BY RAG_SOC LIMIT ?"
        cur.execute(sql, params)
    else:
        sql = f"SELECT {cols_select} FROM clienti ORDER BY RAG_SOC LIMIT ?"
        cur.execute(sql, (limit,))
    rows = cur.fetchall()
    conn.close()
    return cols, rows


def insert_cliente_from_map(data: dict) -> int:
    """Inserisce un cliente usando l'ordine stabile CLIENTI_COLUMNS.
    Campi non presenti in `data` vengono impostati a stringa vuota.
    Eventuali chiavi extra in `data` vengono ignorate per mantenere schema pulito.
    """
    if not data:
        raise ValueError("data vuoto")
    # assicurati che la tabella e le colonne esistano
    conn = get_connection()
    cur = conn.cursor()
    cur.execute("PRAGMA table_info(clienti)")
    existing = [r[1] for r in cur.fetchall()]
    # aggiungi eventuali colonne mancanti dalla lista standard
    for c in CLIENTI_COLUMNS:
        if c not in existing:
            cur.execute(f'ALTER TABLE clienti ADD COLUMN "{c}" TEXT')
    # prepariamo l'inserimento usando SOLO CLIENTI_COLUMNS
    colnames_sql = ", ".join([f'"{c}"' for c in CLIENTI_COLUMNS])
    placeholders = ", ".join(["?"] * len(CLIENTI_COLUMNS))
    vals = [str(data.get(c, "")) for c in CLIENTI_COLUMNS]
    sql = f'INSERT INTO clienti ({colnames_sql}) VALUES ({placeholders})'
    cur.execute(sql, vals)
    nid = cur.lastrowid
    conn.commit()
    conn.close()
    # -- aggiornamenti automatici se i campi sono vuoti --
    if nid:
        aggiorna_finemf_per_cliente_se_vuoto(nid)
        aggiorna_garanzia_per_cliente_se_vuoto(nid)
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
    # -- aggiornamenti automatici se i campi sono vuoti --
    aggiorna_finemf_per_cliente_se_vuoto(rowid)
    aggiorna_garanzia_per_cliente_se_vuoto(rowid)

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
        headers = reader.fieldnames or []
        if not headers:
            raise ValueError("CSV privo di intestazioni o delimiter non riconosciuto")
        rows = list(reader)  # elenco di dict delle righe lette

    # normalizzatore per stringhe: rimuove tutto tranne ascii alfanum, uppercase
    def normalize(s: str) -> str:
        if s is None:
            return ""
        return re.sub(r'[^0-9A-Za-z]', '', str(s)).upper()

    # mapping header -> colonna CLIENTI_COLUMNS con fuzzy matching
    map_header_to_col = {}       # cleaned_header -> CLIENTI_COLUMN
    header_lookup = {}           # cleaned_header -> original_header (come letto da DictReader)
    wanted = CLIENTI_COLUMNS.copy()
    wanted_norm = {wc: normalize(wc) for wc in wanted}

    for h in headers:
        # usa la versione "pulita" dell'header: prima parte prima della virgola
        clean_h = h.split(",")[0].strip()
        hn = normalize(clean_h)

        # match SOLO se clean header coincide esattamente con una colonna
        if clean_h in wanted:
            map_header_to_col[clean_h] = clean_h
            header_lookup[clean_h] = h
            continue
        
        # oppure match case-insensitive
        for wc in wanted:
            if clean_h.upper() == wc.upper():
                map_header_to_col[clean_h] = wc
                header_lookup[clean_h] = h
                break

    # --- fallback posizionale (solo per campi critici) ---
    crucial = ["RAG_SOC", "P_IVA", "IND_VIA"]
    mapped_values = set(map_header_to_col.values())
    missing_crucial = [c for c in crucial if c not in mapped_values]

    if missing_crucial and rows:
        first_row = rows[0]
        first_vals = [first_row.get(h, "") for h in headers]
        if len(headers) >= 4:
            pos_map = {0: "CODICE", 1: "RAG_SOC", 2: "P_IVA", 3: "IND_VIA"}
            for idx, col_name in pos_map.items():
                hdr = headers[idx]
                if hdr not in map_header_to_col and col_name in missing_crucial:
                    val = first_vals[idx] if idx < len(first_vals) else ""
                    if col_name == "P_IVA":
                        if any(ch.isdigit() for ch in str(val)):
                            clean_h = hdr.split(",")[0].strip()
                            map_header_to_col[clean_h] = col_name
                            header_lookup[clean_h] = hdr
                    else:
                        clean_h = hdr.split(",")[0].strip()
                        map_header_to_col[clean_h] = col_name
                        header_lookup[clean_h] = hdr

    # --- DEBUG: stampiamo headers, mapping e prima riga (rimuovere in produzione) ---
    try:
        print("DEBUG: CSV headers as read:", headers)
        print("DEBUG: header_lookup (clean->original):", header_lookup)
        print("DEBUG: map_header_to_col (clean->CLIENTI_COL):", map_header_to_col)
        if rows:
            print("DEBUG: prima riga keys:", list(rows[0].keys()))
            # mostra i valori per FINEMF e TIPOMF se presenti
            for k in ("FINEMF", "TIPOMF"):
                # cerchiamo il clean header che mappa a k
                found = [ch for ch, mapped in map_header_to_col.items() if mapped == k]
                if found:
                    orig = header_lookup.get(found[0])
                    print(f"DEBUG: valore prima riga per {k} (orig header '{orig}') ->", rows[0].get(orig))
                else:
                    print(f"DEBUG: {k} non mappato nella prima riga")
    except Exception:
        pass

    # --- apriamo connessione DB (mancava nella tua patch!) ---
    conn = get_connection()
    cur = conn.cursor()

    # Preparazione insert: usa l'ordine fisso CLIENTI_COLUMNS
    insert_cols = CLIENTI_COLUMNS
    placeholders = ", ".join(["?"] * len(insert_cols))
    colnames_sql = ", ".join([f'"{c}"' for c in insert_cols])
    sql = f'INSERT INTO clienti ({colnames_sql}) VALUES ({placeholders})'
    to_insert = []
    for r in rows:
        vals = []
        for c in insert_cols:
            # cerchiamo se esiste un clean header che mappa a c
            v_found = ""
            for clean_h, mapped in map_header_to_col.items():
                if mapped == c:
                    orig_h = header_lookup.get(clean_h, clean_h)
                    # proviamo più possibilità per estrarre il valore dalla riga
                    v = r.get(orig_h, "")
                    if (v is None or v == ""):
                        # fallback: prova con la clean header se presente nelle keys
                        if clean_h in r:
                            v = r.get(clean_h, "")
                        else:
                            # ancora fallback: prova stripped key (senza spazi/BOM)
                            alt = orig_h.strip('\ufeff').strip()
                            v = r.get(alt, "")
                    if v is None:
                        v = ""
                    v_found = str(v)
                    break
            vals.append(v_found)
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

def aggiorna_finemf_per_cliente_se_vuoto(cid: int):
    conn = get_connection()
    cur = conn.cursor()
    cur.execute("SELECT DATAINST, MXAZ, FINEMF FROM clienti WHERE id=?", (cid,))
    row = cur.fetchone()
    if not row:
        conn.close()
        return
    datainst, mxaz, finemf = row

    if finemf not in (None, "", "0"):
        conn.close()
        return  # già presente → NON sovrascrivere

    nuovo_fin = calcola_fine_memoria(datainst, mxaz)
    if nuovo_fin:
        cur.execute("UPDATE clienti SET FINEMF=? WHERE id=?", (nuovo_fin, cid))
        conn.commit()
    conn.close()


def aggiorna_garanzia_per_cliente_se_vuoto(cid: int):
    conn = get_connection()
    cur = conn.cursor()
    cur.execute("SELECT DATAINST, GARANZIA FROM clienti WHERE id=?", (cid,))
    row = cur.fetchone()
    if not row:
        conn.close()
        return
    datainst, gar = row

    if gar not in (None, "", "0"):
        conn.close()
        return  # già presente → NON sovrascrivere

    nuova_gar = calcola_scadenza_garanzia(datainst, gar)
    if nuova_gar:
        cur.execute("UPDATE clienti SET GARANZIA=? WHERE id=?", (nuova_gar, cid))
        conn.commit()
    conn.close()

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
    if not data_installazione or mxaz is None or str(mxaz).strip() == "":
        return None
    try:
        days = int(str(mxaz))
    except Exception:
        return None
    dt = _try_parse_date(data_installazione)
    if not dt:
        return None
    fin = dt + timedelta(days=days)
    return fin.strftime("%Y-%m-%d")

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

def calcola_scadenza_garanzia(data_installazione: str, garanzia) -> Optional[str]:
    """
    Calcola la data di fine garanzia dato data_installazione e GARANZIA.
    - Se GARANZIA è interpretabile viene usata (mesi).
    - Se GARANZIA non è interpretabile o è vuota -> si assume garanzia base di 1 anno (365 giorni).
    Restituisce stringa YYYY-MM-DD o None se data_installazione non è parsabile.
    """
    if not data_installazione:
        return None
    dt = _try_parse_date(data_installazione)
    if not dt:
        return None

    # prova a interpretare garanzia come mesi
    months = _parse_garanzia_to_months(str(garanzia) if garanzia is not None else "")
    if months is None:
        # se non si riesce a interpretare, assumiamo 1 anno
        days_to_add = 365
    else:
        # conversione mesi -> giorni con approssimazione 30 giorni/mese
        days_to_add = int(months) * 30

    fin = dt + timedelta(days=days_to_add)
    return fin.strftime("%Y-%m-%d")

def _format_date_for_tree(d):
    """Accetta str/date/datetime e restituisce una stringa 'YYYY-MM-DD' o ''."""
    if not d:
        return ""
    if isinstance(d, str):
        # prova a parsare, altrimenti ritorna la stringa così com'è
        try:
            dt = _try_parse_date(d)
            if isinstance(dt, datetime):
                return dt.date().isoformat()
            if isinstance(dt, date):
                return d  # già in formato stringa ma non parsed, ritornala
        except Exception:
            return d
    if isinstance(d, datetime):
        return d.date().isoformat()
    if isinstance(d, date):
        return d.isoformat()
    return str(d)


def report_memorie_in_scadenza(within_days: int = 30, include_expired: bool = False) -> List[Tuple]:
    """
    Restituisce macchine la cui 'fine memoria fiscale' è entro 'within_days' giorni a partire da oggi.
    Usa in ordine: clienti.FINEMF, macchine.FINEMF (se presente), altrimenti calcola da data_installazione + mxaz.
    Ritorna lista di tuple: (cliente_rag, matricola, data_install, finemf_str, giorni_residui, vf, id_macchina)
    Parametri:
      - within_days: numero di giorni da oggi entro cui considerare la scadenza
      - include_expired: se True include anche macchine scadute (days_left < 0)
    """
    conn = get_connection()
    cur = conn.cursor()
    today = datetime.now().date()

    # Controllo se la colonna FINEMF esiste nelle macchine (una volta sola, fuori dal loop)
    # Stampa struttura tabella CLIENTI
    print("DB:", os.path.abspath("database.db"))

    cur.execute("SELECT * FROM macchine")
    print("macchine rows:", cur.fetchall())

    cur.execute("SELECT * FROM clienti")
    print("clienti rows:", cur.fetchall())

    mac_cols = [r[1].upper() for r in cur.fetchall()]
    has_fin_m = "FINEMF" in mac_cols

    # Carico tutti i record (in una sola query)

    rows = cur.fetchall()
    print("rows", rows)
    out = []
    for row in rows:
        print("days_left")
        mid, matricola, data_install, mxaz, vf, rag, finemf_cliente = row

        # Preferiamo FINEMF del cliente
        finemf = finemf_cliente

        # Se non c'è, proviamo a leggere FINEMF da macchine (solo se esiste la colonna)
        if not finemf and has_fin_m:
            cur.execute("SELECT FINEMF FROM macchine WHERE id=?", (mid,))
            r2 = cur.fetchone()
            if r2 and r2[0]:
                finemf = r2[0]

        # Se ancora nulla, calcoliamo da data_installazione + mxaz (se possibile)
        if not finemf and data_install:
            finemf = calcola_fine_memoria(data_install, mxaz)

        if not finemf:
            # nessuna data valida -> salto
            continue

        # Normalizziamo finemf in oggetto datetime/date
        dt = _try_parse_date(finemf)
        # _try_parse_date dovrebbe restituire datetime o None
        if not dt:
            continue

        # ricava la data come oggetto date
        if isinstance(dt, datetime):
            finemf_date = dt.date()
        elif isinstance(dt, date):
            finemf_date = dt
        else:
            # fallback: prova a convertire da stringa con isoformat
            try:
                finemf_date = datetime.fromisoformat(str(dt)).date()
            except Exception:
                continue
            
        days_left = (finemf_date - today).days

        # filtro: include_expired decide come trattare i negativi
        if include_expired:
            cond = days_left <= within_days
        else:
            cond = 0 <= days_left <= within_days

        if cond:
            out.append((
                rag or "",
                matricola or "",
                _format_date_for_tree(data_install),
                _format_date_for_tree(finemf_date),
                int(days_left),
                vf or "",
                mid
            ))

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
    Confronta dati tra intae e misurart usando la colonna join_on (default MATRICOLA).
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