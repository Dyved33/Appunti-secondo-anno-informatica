import sqlite3
from sqlite3 import Connection
from datetime import datetime, timedelta

DB_NAME = "database.db"#assegno nome al database (se lo voglio cambiare lo faccio da qui)

def get_connection() -> Connection:#connessione a database
    return sqlite3.connect(DB_NAME)

def create_tables():#crazione tabelle con informazioni
    conn = get_connection()#connessione al database
    cur = conn.cursor()#cursore che mi permette di interagire con il database

    # Tabella clienti
    cur.execute("""
    CREATE TABLE IF NOT EXISTS clienti (
        piva TEXT PRIMARY KEY,
        ragione_sociale TEXT,
        indirizzo TEXT,
        comune TEXT,
        telefono TEXT,
        email TEXT,
        matricola TEXT,
        data_installazione TEXT,
        scadenza_garanzia TEXT,
        numero_azzeramenti INTEGER DEFAULT 0,
        scadenza_memoria TEXT,
        disattivazione INTEGER DEFAULT 0,
        vf INTEGER DEFAULT 0,
        qrcode TEXT
    )
    """)

    # Tabella interventi (INTAE)
    cur.execute("""
    CREATE TABLE IF NOT EXISTS interventi (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        matricola TEXT,
        data_intervento TEXT,
        cf_tecnico TEXT,
        codice_intervento INTEGER,
        note TEXT
    )
    """)

    # Tabella verifiche (MISURART)
    cur.execute("""
    CREATE TABLE IF NOT EXISTS verifiche (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        matricola TEXT,
        piva_cliente TEXT,
        data_intervento TEXT,
        tipo_intervento INTEGER,
        note TEXT,
        FOREIGN KEY (piva_cliente) REFERENCES clienti(piva)
    )
    """)

    conn.commit()#salvo le modifiche
    conn.close()#chiudo la connessione

# -------------------------------
# Funzioni clienti
# -------------------------------

def insert_cliente(piva, ragione_sociale, indirizzo=None, comune=None, telefono=None, email=None,
                   matricola=None, data_installazione=None, numero_azzeramenti=0, qrcode=None):
    """Inserisce un nuovo cliente e calcola scadenze automaticamente"""
    conn = get_connection()
    cur = conn.cursor()

    scadenza_garanzia = None
    scadenza_memoria = None
    if data_installazione:
        try:
            dt = datetime.strptime(data_installazione, "%Y-%m-%d")
            scadenza_garanzia = (dt + timedelta(days=365*2)).strftime("%Y-%m-%d")  # esempio: garanzia 2 anni
            scadenza_memoria = (dt + timedelta(days=numero_azzeramenti*365)).strftime("%Y-%m-%d")
        except:
            pass

    cur.execute("""
        INSERT OR IGNORE INTO clienti (piva, ragione_sociale, indirizzo, comune, telefono, email,
                                       matricola, data_installazione, scadenza_garanzia, numero_azzeramenti,
                                       scadenza_memoria, qrcode)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    """, (piva, ragione_sociale, indirizzo, comune, telefono, email,
          matricola, data_installazione, scadenza_garanzia, numero_azzeramenti,
          scadenza_memoria, qrcode))
    conn.commit()
    conn.close()

def update_cliente(piva, **kwargs):
    """Modifica i dati di un cliente"""
    conn = get_connection()
    cur = conn.cursor()

    fields = ", ".join([f"{k}=?" for k in kwargs.keys()])
    values = list(kwargs.values())
    values.append(piva)
    sql = f"UPDATE clienti SET {fields} WHERE piva=?"
    cur.execute(sql, values)

    conn.commit()
    conn.close()

def duplicate_cliente(piva):
    """Duplica tutte le informazioni di un cliente generando nuova PIVA (esempio aggiunge '_dup')"""
    conn = get_connection()
    cur = conn.cursor()

    cur.execute("SELECT * FROM clienti WHERE piva=?", (piva,))
    row = cur.fetchone()
    if row:
        new_piva = f"{piva}_dup"
        data = list(row)
        data[0] = new_piva  # sostituisci PIVA
        cur.execute("""
            INSERT INTO clienti (piva, ragione_sociale, indirizzo, comune, telefono, email,
                                  matricola, data_installazione, scadenza_garanzia, numero_azzeramenti,
                                  scadenza_memoria, disattivazione, vf, qrcode)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        """, data)
    conn.commit()
    conn.close()

# -------------------------------
# Funzioni interventi
# -------------------------------

def insert_intervento(matricola, data_intervento, cf_tecnico, codice_intervento, note=None):
    conn = get_connection()
    cur = conn.cursor()
    cur.execute("""
        INSERT INTO interventi (matricola, data_intervento, cf_tecnico, codice_intervento, note)
        VALUES (?, ?, ?, ?, ?)
    """, (matricola, data_intervento, cf_tecnico, codice_intervento, note))
    conn.commit()
    conn.close()

# -------------------------------
# Funzioni verifiche (MISURART)
# -------------------------------

def insert_verifica(matricola, piva_cliente, data_intervento, tipo_intervento, note=None):
    conn = get_connection()
    cur = conn.cursor()
    cur.execute("""
        INSERT INTO verifiche (matricola, piva_cliente, data_intervento, tipo_intervento, note)
        VALUES (?, ?, ?, ?, ?)
    """, (matricola, piva_cliente, data_intervento, tipo_intervento, note))
    conn.commit()
    conn.close()

# -------------------------------
# Funzioni ricerca / report
# -------------------------------

def find_cliente(ragione_sociale=None, comune=None, matricola=None):
    conn = get_connection()
    cur = conn.cursor()
    sql = "SELECT * FROM clienti WHERE 1=1"
    params = []
    if ragione_sociale:
        sql += " AND ragione_sociale LIKE ?"
        params.append(f"%{ragione_sociale}%")
    if comune:
        sql += " AND comune LIKE ?"
        params.append(f"%{comune}%")
    if matricola:
        sql += " AND matricola=?"
        params.append(matricola)
    cur.execute(sql, params)
    rows = cur.fetchall()
    conn.close()
    return rows

def verifiche_in_scadenza():
    """Restituisce le verifiche periodiche in scadenza (es. entro 30 giorni)"""
    conn = get_connection()
    cur = conn.cursor()
    oggi = datetime.today().strftime("%Y-%m-%d")
    cur.execute("""
        SELECT * FROM clienti
        WHERE scadenza_memoria <= ?
    """, (oggi,))
    rows = cur.fetchall()
    conn.close()
    return rows

def report_interventi_tipo(tipo):
    """Restituisce interventi filtrati per tipo"""
    conn = get_connection()
    cur = conn.cursor()
    cur.execute("SELECT * FROM interventi WHERE codice_intervento=?", (tipo,))
    rows = cur.fetchall()
    conn.close()
    return rows
