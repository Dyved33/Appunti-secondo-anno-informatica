import tkinter as tk
from tkinter import ttk, filedialog, messagebox
import os
import json
import shutil
from datetime import date, datetime
import traceback
import database
from database import CLIENTI_COLUMNS
from database import get_connection, _normalize_matricola_for_compare
import subprocess
import sys
import re

def _safe_call(fn, *args, **kwargs):
    """Wrapper per chiamate di funzioni con gestione errori automatica."""
    try:
        return fn(*args, **kwargs)
    except Exception as e:
        traceback.print_exc()
        messagebox.showerror("Errore", str(e))
        return None

def _format_date_string(digits_only, max_length=8):
    """
    Formatta una stringa di soli numeri come DD/MM/YYYY o DD/MM/YY.
    Ritorna la stringa formattata.
    """
    if max_length == 6:  # DD/MM/YY
        if len(digits_only) >= 6:
            return f"{digits_only[0:2]}/{digits_only[2:4]}/{digits_only[4:6]}"
        elif len(digits_only) >= 4:
            return f"{digits_only[0:2]}/{digits_only[2:4]}/{digits_only[4:]}"
        elif len(digits_only) >= 2:
            return f"{digits_only[0:2]}/{digits_only[2:]}"
        else:
            return digits_only
    else:  # DD/MM/YYYY (max_length=8 o default)
        if len(digits_only) >= 8:
            return f"{digits_only[0:2]}/{digits_only[2:4]}/{digits_only[4:8]}"
        elif len(digits_only) >= 5:
            return f"{digits_only[0:2]}/{digits_only[2:4]}/{digits_only[4:]}"
        elif len(digits_only) >= 3:
            return f"{digits_only[0:2]}/{digits_only[2:]}"
        else:
            return digits_only

def _validate_date_input(entry_widget, max_length=8):
    """
    Valida input data: accetta solo numeri (0-9), limita a max_length digit.
    Formatta automaticamente come DD/MM/YYYY (se max_length=8) o DD/MM/YY (se max_length=6).
    Il cursore rimane alla fine del campo.
    """
    value = entry_widget.get()
    
    # Estrai solo i numeri
    digits_only = re.sub(r'[^\d]', '', value)
    
    # Limita a max_length cifre
    if len(digits_only) > max_length:
        digits_only = digits_only[:max_length]
    
    # Formatta
    formatted = _format_date_string(digits_only, max_length)
    
    # Aggiorna il campo se diverso
    if formatted != value:
        entry_widget.delete(0, tk.END)
        entry_widget.insert(0, formatted)
        # Cursore alla fine
        try:
            entry_widget.icursor(len(formatted))
        except Exception:
            pass

def _create_date_entry(parent, width=20, max_length=8, with_var=False):
    """
    Crea un Entry per le date con formattazione automatica DD/MM/YYYY.
    
    Args:
        parent: widget genitore
        width: larghezza dell'entry
        max_length: numero massimo di cifre (8 per YYYY, 6 per YY)
        with_var: se True ritorna (widget, StringVar), altrimenti solo widget
    
    Returns:
        Entry widget se with_var=False, oppure tuple (widget, StringVar) se with_var=True
    """
    var = tk.StringVar() if with_var else None
    ent = ttk.Entry(parent, textvariable=var if with_var else None, width=width)
    
    def on_key_release(event):
        """Formatta la data nel formato DD/MM/YYYY o DD/MM/YY"""
        entry = event.widget
        value = entry.get()
        
        # Estrai solo numeri
        digits_only = re.sub(r'[^\d]', '', value)
        
        # Limita a max_length cifre
        if len(digits_only) > max_length:
            digits_only = digits_only[:max_length]
        
        # Formatta secondo max_length
        if max_length == 6:  # DD/MM/YY
            if len(digits_only) >= 6:
                formatted = f"{digits_only[0:2]}/{digits_only[2:4]}/{digits_only[4:6]}"
            elif len(digits_only) >= 4:
                formatted = f"{digits_only[0:2]}/{digits_only[2:4]}/{digits_only[4:]}"
            elif len(digits_only) >= 2:
                formatted = f"{digits_only[0:2]}/{digits_only[2:]}"
            else:
                formatted = digits_only
        else:  # DD/MM/YYYY (default)
            if len(digits_only) >= 8:
                formatted = f"{digits_only[0:2]}/{digits_only[2:4]}/{digits_only[4:8]}"
            elif len(digits_only) >= 5:
                formatted = f"{digits_only[0:2]}/{digits_only[2:4]}/{digits_only[4:]}"
            elif len(digits_only) >= 3:
                formatted = f"{digits_only[0:2]}/{digits_only[2:]}"
            else:
                formatted = digits_only
        
        # Aggiorna solo se diverso (evita flicker)
        if formatted != value:
            entry.delete(0, tk.END)
            entry.insert(0, formatted)
            # Cursore alla fine del campo
            entry.icursor(len(formatted))
    
    ent.bind('<KeyRelease>', on_key_release)
    return (ent, var) if with_var else ent



class App:
    def __init__(self, root):
        self.root = root
        self.root.title("Gestione Clienti e Interventi")
        self.root.geometry("1200x720")

        # lista in memoria dei clienti
        self.clienti = []

        # inizializza DB senza ricreare la tabella clienti (persistenza dati)
        db_path = os.path.abspath(database.DB_FILE)
        db_existed = os.path.exists(db_path)
        database.ensure_db_exists(recreate_clienti=False)

        # Configurazione backup solo al primo avvio con DB assente (bloccante)
        if not db_existed:
            self._require_backup_dir_first_time()

        # Backup giornaliero (senza forzare configurazione)
        self._ensure_backup_folder_and_daily_backup()

        # Mostra pop-up con memorie fiscali in scadenza
        self._show_expiring_memories_popup()

        # costruzione interfaccia
        self._build_ui()

    def _app_state_path(self):
        base_dir = os.path.dirname(os.path.abspath(__file__))
        return os.path.join(base_dir, "app_state.json")

    def _load_app_state(self):
        path = self._app_state_path()
        if not os.path.exists(path):
            return {}
        try:
            with open(path, "r", encoding="utf-8") as f:
                return json.load(f) or {}
        except Exception:
            return {}

    def _save_app_state(self, state: dict):
        path = self._app_state_path()
        try:
            with open(path, "w", encoding="utf-8") as f:
                json.dump(state or {}, f, indent=2)
        except Exception:
            pass
    def _require_backup_dir_first_time(self):
        """Obbliga la configurazione backup solo quando il DB non esiste ancora."""
        state = self._load_app_state()
        backup_dir = state.get("backup_dir")
        if backup_dir and os.path.isdir(backup_dir):
            return
        backup_dir = self._prompt_for_backup_dir()
        state["backup_dir"] = backup_dir
        state["last_backup_date"] = None
        self._save_app_state(state)

    def _prompt_for_backup_dir(self):
        while True:
            selected = filedialog.askdirectory(title="Seleziona cartella per backup", mustexist=True)
            if selected:
                return selected
            messagebox.showwarning(
                "Selezione richiesta",
                "Devi selezionare una cartella per continuare.",
                parent=self.root
            )

    def _ensure_backup_folder_and_daily_backup(self):
        """
        Verifica se e' stata configurata una cartella di backup e crea un backup giornaliero automatico.
        Cross-platform compatibile (Windows, Linux, macOS).
        """
        state = self._load_app_state()
        backup_dir = state.get("backup_dir")

        if not backup_dir or not os.path.isdir(backup_dir):
            print("[WARN] Cartella backup non configurata: skip backup giornaliero")
            return

        today_str = date.today().isoformat()
        if state.get("last_backup_date") == today_str:
            print(f"[INFO] Backup gia' eseguito oggi: {today_str}")
            return

        db_path = os.path.abspath(database.DB_FILE)
        if not os.path.exists(db_path):
            print(f"[WARN] Database non trovato: {db_path}")
            return

        try:
            ts = datetime.now().strftime("%d_%m_%Y_%H%M")
            backup_name = f"database_{ts}.db"
            dest_path = os.path.join(backup_dir, backup_name)

            if len(dest_path) > 250:
                print(f"[WARN] Path del backup troppo lungo per Windows: {len(dest_path)} caratteri")
                backup_name = f"db_{ts}.db"
                dest_path = os.path.join(backup_dir, backup_name)

            shutil.copy2(db_path, dest_path)
            state["last_backup_date"] = today_str
            state["last_backup_path"] = dest_path
            self._save_app_state(state)
            print(f"[INFO] Backup completato: {backup_name}")
        except PermissionError as e:
            print(f"[ERROR] Permessi insufficienti per backup: {e}")
            messagebox.showwarning("Backup", f"Attenzione: permessi insufficienti.\nCartella: {backup_dir}", parent=self.root)
        except Exception as e:
            print(f"[ERROR] Errore durante backup: {e}")
            messagebox.showwarning("Backup", f"Attenzione: errore durante backup.\nDettagli: {e}", parent=self.root)

    def _force_manual_backup(self):
        """Backup manuale con opzione di riconfigurare la cartella."""
        state = self._load_app_state()
        backup_dir = state.get("backup_dir")

        if backup_dir and os.path.isdir(backup_dir):
            if messagebox.askyesno(
                "Backup",
                "Vuoi riconfigurare la cartella di backup prima di procedere?",
                parent=self.root
            ):
                backup_dir = self._prompt_for_backup_dir()
                state["backup_dir"] = backup_dir
                state["last_backup_date"] = None
                self._save_app_state(state)
        else:
            backup_dir = self._prompt_for_backup_dir()
            state["backup_dir"] = backup_dir
            state["last_backup_date"] = None
            self._save_app_state(state)

        if not messagebox.askyesno(
            "Conferma Backup",
            f"Creare un backup manuale del database nella cartella:\n{backup_dir}?",
            parent=self.root
        ):
            return

        db_path = os.path.abspath(database.DB_FILE)
        if not os.path.exists(db_path):
            messagebox.showerror("Errore Backup", f"Database non trovato: {db_path}", parent=self.root)
            return

        try:
            ts = datetime.now().strftime("%d_%m_%Y_%H%M")
            backup_name = f"database_{ts}.db"
            dest_path = os.path.join(backup_dir, backup_name)

            if len(dest_path) > 250:
                print("[WARN] Path troppo lungo per Windows, uso nome abbreviato")
                backup_name = f"db_{ts}.db"
                dest_path = os.path.join(backup_dir, backup_name)

            shutil.copy2(db_path, dest_path)
            print(f"[INFO] Backup manuale completato: {backup_name}")
            messagebox.showinfo("Backup Completato", f"Backup salvato con successo:\n{backup_name}", parent=self.root)
        except PermissionError as e:
            print(f"[ERROR] Permessi insufficienti: {e}")
            messagebox.showerror("Errore Backup", f"Permessi insufficienti.\nCartella: {backup_dir}", parent=self.root)
        except Exception as e:
            print(f"[ERROR] Errore durante backup: {e}")
            messagebox.showerror("Errore Backup", f"Errore durante il backup:\n{e}", parent=self.root)

    def _show_expiring_memories_popup(self):
        """
        Verifica se ci sono memorie fiscali in scadenza entro 30 giorni e mostra un pop-up.
        Il popup è non-modal (permette di vedere dietro) ma sempre in primo piano.
        """
        try:
            memorie_in_scadenza = database.report_memorie_in_scadenza(within_days=30) or []
            if not memorie_in_scadenza:
                print("[INFO] Nessuna memoria in scadenza nei prossimi 30 giorni")
                return

            print(f"[INFO] Trovate {len(memorie_in_scadenza)} memorie in scadenza")

            popup = tk.Toplevel(self.root)
            popup.title("⚠️  ATTENZIONE - Memorie fiscali in scadenza")
            popup.transient(self.root)
            popup.geometry("800x450")
            
            # Forza il popup in primo piano e lo rende visibile
            popup.attributes("-topmost", True)
            popup.lift()
            popup.focus_force()
            popup.update()

            main_frame = ttk.Frame(popup, padding=12)
            main_frame.pack(fill="both", expand=True)

            # Titolo importante
            title_label = ttk.Label(
                main_frame,
                text="⚠️  Le seguenti memorie fiscali stanno scadendo nei prossimi 30 giorni:",
                font=("Helvetica", 11, "bold"),
                foreground="red"
            )
            title_label.pack(anchor="w", pady=(0, 12))

            # Frame con scrollbar per la lista
            list_frame = ttk.Frame(main_frame)
            list_frame.pack(fill="both", expand=True, pady=(0, 12))

            # Treeview per mostrare i dati in modo ordinato
            columns = ("Ragione Sociale", "Data Scadenza", "MF", "Telefono", "Codice")
            tree = ttk.Treeview(list_frame, columns=columns, show="headings", height=12)
            
            # Colonna scrollbar
            scrollbar = ttk.Scrollbar(list_frame, orient="vertical", command=tree.yview)
            tree.configure(yscrollcommand=scrollbar.set)
            
            # Intestazioni
            tree.heading("Ragione Sociale", text="Ragione Sociale")
            tree.heading("Data Scadenza", text="Data Scadenza")
            tree.heading("MF", text="MF")
            tree.heading("Telefono", text="Telefono")
            tree.heading("Codice", text="Codice")
            
            # Larghezze colonne
            tree.column("Ragione Sociale", width=280, anchor="w")
            tree.column("Data Scadenza", width=130, anchor="center")
            tree.column("MF", width=80, anchor="center")
            tree.column("Telefono", width=130, anchor="w")
            tree.column("Codice", width=100, anchor="center")
            
            # Inserisci i dati nella treeview
            # Le tuple sono: (RAG_SOC, FINEMF, MF, TEL, CODICE)
            for rag_soc, finemf, mf, tel, codice in memorie_in_scadenza:
                tree.insert("", "end", values=(rag_soc, finemf, mf, tel, codice))
            
            tree.pack(side="left", fill="both", expand=True)
            scrollbar.pack(side="right", fill="y")

            # Frame pulsanti
            btn_frame = ttk.Frame(main_frame)
            btn_frame.pack(fill="x", pady=(12, 0))

            def go_to_report():
                """Vai al report memorie scadenza"""
                popup.destroy()
                self.nb.select(3)  # Seleziona tab Report
                self.gui_memorie_scadenza(30)

            ttk.Button(btn_frame, text="📋 Vai al Report", command=go_to_report).pack(side="left", padx=6)
            ttk.Button(btn_frame, text="Chiudi", command=popup.destroy).pack(side="left", padx=6)
            
            # Chiudi automaticamente dopo 10 secondi se l'utente non interagisce
            def auto_close():
                try:
                    popup.destroy()
                except Exception:
                    pass
            
            popup.after(10000, auto_close)  # 10 secondi
            
        except Exception as e:
            print(f"[Errore pop-up memorie]: {e}")
            traceback.print_exc()


    def _build_ui(self):
        self.nb = ttk.Notebook(self.root)
        self.nb.pack(fill="both", expand=True)
        self.tab_clienti = ttk.Frame(self.nb)
        self.tab_intae = ttk.Frame(self.nb)
        self.tab_misurart = ttk.Frame(self.nb)
        self.tab_report = ttk.Frame(self.nb)
        self.nb.add(self.tab_clienti, text="Clienti")
        self.nb.add(self.tab_intae, text="INTAE")
        self.nb.add(self.tab_misurart, text="MISURART")
        self.nb.add(self.tab_report, text="Report")
        self._build_tab_clienti()
        self._build_tab_intae()
        self._build_tab_misurart()
        self._build_tab_report()

    # -------------------------
    # Tab Clienti
    # -------------------------
    def _build_tab_clienti(self):
        # top area: buttons + search stacked horizontally
        top = ttk.Frame(self.tab_clienti, padding=8)
        top.pack(fill="x")

        btn_fr = ttk.Frame(top)
        btn_fr.pack(side="left", fill="x", expand=False)
        ttk.Button(btn_fr, text="IMPORTA FILE", command=self.gui_import_any_clienti).pack(side="left", padx=4)
        ttk.Button(btn_fr, text="NUOVO", command=self.gui_new_cliente).pack(side="left", padx=4)
        ttk.Button(btn_fr, text="DUPLICA", command=self.gui_duplicate_cliente).pack(side="left", padx=4)
        # prima: ttk.Button(btn_fr, text="ELIMINA", command=self.gui_delete_cliente).pack(side="left", padx=4)
        tk.Button(btn_fr, text="ELIMINA", command=self.gui_delete_cliente, bg="#d9534f", fg="white", activebackground="#c9302c").pack(side="left", padx=4)
        ttk.Button(btn_fr, text="REPORT RT", command=self.gui_report_rt_per_model).pack(side="left", padx=4)
        ttk.Button(btn_fr, text="DISATTIVA MACCHINA", command=self.gui_deactivate_machine).pack(side="left", padx=4)
        ttk.Button(btn_fr, text="BACKUP", command=self._force_manual_backup).pack(side="left", padx=4)
        ttk.Button(btn_fr, text="REFRESH", command=lambda: self.load_clienti(show_all=True)).pack(side="left", padx=4)
        search = ttk.Frame(top)
        search.pack(side="right")
        ttk.Label(search, text="Cerca:").pack(side="left")
        self.client_search_var = tk.StringVar()
        ttk.Entry(search, textvariable=self.client_search_var, width=30).pack(side="left", padx=4)
        ttk.Button(search, text="CERCA", command=lambda: self.load_clienti(show_all=False)).pack(side="left", padx=4)

        # sotto i bottoni: area tabella clienti
        self.right_clienti_frame = ttk.Frame(self.tab_clienti, padding=8)
        self.right_clienti_frame.pack(fill="both", expand=True)

        # Treeview clienti
        self.tree = ttk.Treeview(
            self.right_clienti_frame,
            columns=CLIENTI_COLUMNS,
            show="headings"
        )
        self.tree.pack(fill="both", expand=True)

        # intestazioni e larghezza colonne
        for c in CLIENTI_COLUMNS:
            self.tree.heading(c, text=c)
            self.tree.column(c, width=100, anchor="center")

        # scrollbar verticale
        scrollbar = ttk.Scrollbar(self.right_clienti_frame, orient="vertical", command=self.tree.yview)
        self.tree.configure(yscrollcommand=scrollbar.set)
        scrollbar.pack(side="right", fill="y")

        # altre proprietà
        self.clienti_col_widths = {}
        self.selected_cliente_id = None
        self.last_clienti_cols = []

        # carica i dati dal DB e mostra la tabella
        self.load_clienti(show_all=True)

    def load_clienti(self, show_all=False, limit: int = 1000):
        term = None if show_all else self.client_search_var.get().strip()
        res = _safe_call(database.read_clienti_allcols, term, limit)
        if res is None:
            cols, rows = database.read_clienti_allcols(None, limit)
        else:
            cols, rows = res
        # salva colonne per riferimenti successivi
        self.last_clienti_cols = cols

        # --- prima di rimuovere la tree esistente, salva le larghezze attuali ---
        if self.tree is not None:
            # salva larghezze correnti
            self._save_clienti_col_widths(self.tree)

        # destroy existing tree area
        for w in self.right_clienti_frame.winfo_children():
            w.destroy()
        container = ttk.Frame(self.right_clienti_frame)
        container.pack(fill="both", expand=True)
        tree = ttk.Treeview(container, columns=cols, show="headings")
        vsb = ttk.Scrollbar(container, orient="vertical", command=tree.yview)
        hsb = ttk.Scrollbar(container, orient="horizontal", command=tree.xview)
        tree.configure(yscroll=vsb.set, xscroll=hsb.set)
        vsb.pack(side="right", fill="y")
        hsb.pack(side="bottom", fill="x")
        tree.pack(fill="both", expand=True)
        self.tree = tree

        # crea intestazioni e colonne riutilizzando larghezze salvate se presenti
        for c in cols:
            if c.lower() == "id":
                tree.heading(c, text=c)
            else:
                tree.heading(c, text=c, command=lambda _c=c: self._sort_clienti_by_column(_c))
            if c.lower() == "id":
                # nascondi id come prima
                tree.column(c, width=0, minwidth=0, stretch=False)
            else:
                # riprendi larghezza salvata se disponibile, altrimenti default 140
                width = self.clienti_col_widths.get(c, 140)
                tree.column(c, width=width, minwidth=80, stretch=False)

        for r in rows:
            tree.insert("", "end", values=r)

        # bind per aggiornare le larghezze quando l'utente rilascia il mouse
        # (dopo un ridimensionamento col cursore)
        tree.bind("<ButtonRelease-1>", lambda e, t=tree: self._save_clienti_col_widths(t))

        # IMPORTANTE: bind anche per il TreeviewSelect dopo aver ricreato l'albero
        tree.bind("<<TreeviewSelect>>", self.on_select_cliente)

    def on_select_cliente(self, event):
        sel = self.tree.selection()
        if not sel:
            return
        vals = self.tree.item(sel[0])["values"]
        cols = self.last_clienti_cols or database.get_clienti_columns()
        id_index = None
        for i, c in enumerate(cols):
            if c.lower() == "id":
                id_index = i
                break
        self.selected_cliente_id = vals[id_index] if id_index is not None and id_index < len(vals) else None
        # open editor window with all fields
        self.open_edit_clients_window(cols, vals)

    def _sort_clienti_by_column(self, col: str):
        tree = self.tree
        if tree is None:
            return
        is_date_col = col.upper() in getattr(database, "DATE_COLS", set())
        sort_state = getattr(self, "_clienti_sort_state", {})

        if col not in sort_state:
            descending = True if is_date_col else False
        else:
            descending = not sort_state[col]

        sort_state[col] = descending
        self._clienti_sort_state = sort_state

        def sort_key(iid):
            val = tree.set(iid, col)
            if is_date_col:
                try:
                    parsed = database._parse_date(val)
                    if isinstance(parsed, datetime):
                        return parsed.date()
                    if isinstance(parsed, date):
                        return parsed
                except Exception:
                    pass
                return date.min
            return str(val).casefold()

        items = list(tree.get_children(""))
        items.sort(key=sort_key, reverse=descending)
        for idx, iid in enumerate(items):
            tree.move(iid, "", idx)

    def open_edit_clients_window(self, cols, vals):
        """Apre finestra modale per modificare i dati di un cliente."""
        try:
            w = tk.Toplevel(self.root)
            w.title("Modifica cliente")
            w.geometry("900x700")  # Imposta dimensioni finestra
            w.transient(self.root)  # Associa al root
            w.update()  # Rendi la finestra visibile prima di grab_set()
            w.grab_set()  # Rendi modale
            w.lift()
            w.focus_force()
            
            # Frame principale che si espande
            frm = ttk.Frame(w, padding=8)
            frm.pack(fill="both", expand=True)
            frm.rowconfigure(0, weight=1)
            frm.columnconfigure(0, weight=1)
            
            # Canvas con scrollbar
            canvas = tk.Canvas(frm, bg="white", height=550, width=800)
            scrollbar = ttk.Scrollbar(frm, orient="vertical", command=canvas.yview)
            scroll_frame = ttk.Frame(canvas)
            scroll_frame.bind("<Configure>", lambda e: canvas.configure(scrollregion=canvas.bbox("all")))
            canvas.create_window((0,0), window=scroll_frame, anchor="nw")
            canvas.configure(yscrollcommand=scrollbar.set)
            
            # Grid del canvas e scrollbar
            canvas.grid(row=0, column=0, sticky="nsew")
            scrollbar.grid(row=0, column=1, sticky="ns")
            
            # Bind scroll wheel (compatibile con Windows/Mac/Linux)
            def _on_mousewheel(event):
                try:
                    if hasattr(event, 'delta'):
                        canvas.yview_scroll(int(-1*(event.delta/120)), "units")
                except Exception:
                    pass
            
            # Configura scroll per diversi sistemi operativi
            try:
                canvas.bind("<MouseWheel>", _on_mousewheel)  # Windows/Mac
            except Exception:
                pass
            try:
                canvas.bind("<Button-4>", lambda e: canvas.yview_scroll(-1, "units"))  # Linux scroll up
                canvas.bind("<Button-5>", lambda e: canvas.yview_scroll(1, "units"))   # Linux scroll down
            except Exception:
                pass
            
            # Crea campi di input per ciascuna colonna
            entries = {}
            rowi = 0
            for i, c in enumerate(cols):
                if c.lower() == "id":
                    continue
                try:
                    ttk.Label(scroll_frame, text=c).grid(row=rowi, column=0, sticky="w", padx=4, pady=2)
                    # Usa entry formattata per i campi data
                    if c.upper() in getattr(database, "DATE_COLS", set()):
                        ent = _create_date_entry(scroll_frame, width=70, max_length=8)
                    else:
                        ent = ttk.Entry(scroll_frame, width=70)
                    ent.grid(row=rowi, column=1, sticky="w", padx=4, pady=2)
                    val = vals[i] if i < len(vals) else ""
                    ent.insert(0, "" if val is None else str(val))
                    entries[c] = ent
                except Exception as e:
                    messagebox.showerror("Errore campo", f"Errore creando campo per {c}: {e}", parent=w)
                
                rowi += 1
            
            btn_fr = ttk.Frame(scroll_frame)
            btn_fr.grid(row=rowi+1, column=0, columnspan=2, pady=8)
            
            def do_save():
                # raccogli input
                updates = {c: entries[c].get().strip() for c in entries}

                # estrai matricola per controllo duplicato
                mat = updates.get("MATRICOLA", "").strip()
                
                # VALIDAZIONE: DATAINST è obbligatoria
                datainst_val = updates.get("DATAINST", "").strip()
                if not datainst_val:
                    messagebox.showerror("Errore", "DATAINST (Data Installazione) è obbligatoria", parent=w)
                    try:
                        w.lift()
                        w.focus_force()
                    except Exception:
                        pass
                    return

                try:
                    # Se è INSERT (nuovo cliente)
                    if not self.selected_cliente_id:
                        if mat:
                            # controllo duplicato: se esiste => popup e blocco
                            if database.matricola_exists(mat):
                                messagebox.showerror("Errore", f"Impossibile inserire: la matricola '{mat}' è già presente.", parent=w)
                                try:
                                    w.lift()
                                    w.focus_force()
                                except Exception:
                                    pass
                                return

                    else:
                        # Se è UPDATE: se la nuova matricola è presente in un'altra riga -> errore
                        if mat:
                            if database.matricola_exists(mat, exclude_id=self.selected_cliente_id):
                                messagebox.showerror("Errore", f"Impossibile aggiornare: la matricola '{mat}' è già presente su un altro cliente.", parent=w)
                                try:
                                    w.lift()
                                    w.focus_force()
                                except Exception:
                                    pass
                                return

                    # Normalizza / converte campi data in ISO prima di salvare (evita formati ambigui)
                    for k in list(updates.keys()):
                        if k.upper() in getattr(database, "DATE_COLS", set()):
                            v = updates[k].strip()
                            if v == "":
                                # lascia vuoto (DB accetta stringa vuota)
                                updates[k] = ""
                            else:
                                try:
                                    # _parse_date_param_to_iso accetta 'DDMMYYYY' e 'DD/MM/YYYY' e ritorna YYYY-MM-DD
                                    updates[k] = database._parse_date_param_to_iso(v)
                                except Exception as e:
                                    messagebox.showerror("Errore data", f"Formato data non valido per '{k}': {v}\nDettaglio: {e}", parent=w)
                                    try:
                                        w.lift()
                                        w.focus_force()
                                    except Exception:
                                        pass
                                    return

                    # chiamata DB (wrap _safe_call per gestire eccezioni)
                    if not self.selected_cliente_id:
                        nid = _safe_call(database.insert_clients_from_map, updates)
                        if nid:
                            msg = messagebox.showinfo("OK","Cliente inserito", parent=w)
                            try:
                                w.lift()
                                w.focus_force()
                            except Exception:
                                pass
                        else:
                            # insert fallito: _safe_call mostrerà errore se eccezione
                            return
                    else:
                        affected = _safe_call(database.update_clients_by_id, self.selected_cliente_id, updates)
                        if affected is None:
                            return
                        msg = messagebox.showinfo("OK","Cliente aggiornato", parent=w)
                        try:
                            w.lift()
                            w.focus_force()
                        except Exception:
                            pass

                    # Ricarica la vista clienti per mostrare immediatamente le modifiche
                    try:
                        self.load_clienti(show_all=True)
                    except Exception:
                        # fallback: ricarica manuale della lista interna e refresh della tabella
                        try:
                            self.carica_clienti_da_db()
                            self.refresh_tabella()
                        except Exception:
                            pass

                    # chiudi finestra solo dopo aver aggiornato la lista
                    w.destroy()

                except Exception as e:
                    messagebox.showerror("Errore salvataggio", str(e), parent=w)
                    try:
                        w.lift()
                        w.focus_force()
                    except Exception:
                        pass
                    return
            
            def do_preview_checklist():
                """Mostra anteprima della checklist senza stampare"""
                try:
                    # Raccogli i dati del cliente dal form
                    cliente_data = {c: entries[c].get().strip() for c in entries}
                    
                    # Genera PDF
                    pdf_path = _safe_call(database.generate_cliente_checklist_pdf, cliente_data)
                    
                    if pdf_path is None:
                        messagebox.showerror("Errore", "Errore nella generazione della checklist", parent=w)
                        return
                    
                    # Apri il PDF con l'applicazione predefinita (solo preview, non stampa)
                    if sys.platform.startswith("win"):
                        os.startfile(os.path.abspath(pdf_path))
                    elif sys.platform == "darwin":  # macOS
                        subprocess.Popen(["open", pdf_path])
                    else:  # Linux
                        subprocess.Popen(["xdg-open", pdf_path])
                except Exception as e:
                    messagebox.showerror("Errore", f"Errore nell'apertura dell'anteprima: {e}", parent=w)
            
            def do_print_checklist():
                """Stampa la checklist (2 copie per impostazione predefinita)"""
                try:
                    # Raccogli i dati del cliente dal form
                    cliente_data = {c: entries[c].get().strip() for c in entries}
                    
                    # Genera PDF
                    pdf_path = _safe_call(database.generate_cliente_checklist_pdf, cliente_data)
                    
                    if pdf_path is None:
                        messagebox.showerror("Errore", "Errore nella generazione della checklist", parent=w)
                        return
                    
                    # Stampa il PDF
                    _safe_call(database.print_pdf_file, pdf_path, copies=2)
                    messagebox.showinfo("Stampa", "Checklist inviata alla stampante (2 copie)", parent=w)
                except Exception as e:
                    messagebox.showerror("Errore", f"Errore nella stampa della checklist: {e}", parent=w)
     
            def do_cancel():
                w.destroy()
            
            ttk.Button(btn_fr, text="ANTEPRIMA CHECKLIST", command=do_preview_checklist).pack(side="left", padx=6)
            ttk.Button(btn_fr, text="STAMPA CHECKLIST", command=do_print_checklist).pack(side="left", padx=6)
            ttk.Button(btn_fr, text="SALVA", command=do_save).pack(side="left", padx=6)
            ttk.Button(btn_fr, text="ANNULLA", command=do_cancel).pack(side="left", padx=6)
            
            w.focus_force()
        
        except Exception as e:
            traceback.print_exc()
            messagebox.showerror("Errore", f"Errore nell'apertura della finestra: {e}", parent=self.root)

    def gui_new_cliente(self):
        cols = database.get_clienti_columns()
        self.selected_cliente_id = None
        # Crea default con MXAZ=2500
        default_vals = []
        for col in cols:
            if col == "MXAZ":
                default_vals.append("2500")
            else:
                default_vals.append("")
        self.open_edit_clients_window(cols, default_vals)

    def gui_update_cliente(self):
        sel = self.tree.selection()
        if not sel:
            messagebox.showerror("Errore","Seleziona un cliente", parent=self.root)
            return
        vals = self.tree.item(sel[0])["values"]
        cols = database.get_clienti_columns()
        # setta l'id selezionato correttamente
        id_index = None
        for i,c in enumerate(cols):
            if c.lower() == "id":
                id_index = i
                break
        self.selected_cliente_id = vals[id_index] if id_index is not None and id_index < len(vals) else None
        self.open_edit_clients_window(cols, vals)

    def gui_delete_cliente(self):
        if not self.selected_cliente_id:
            messagebox.showerror("Errore","Seleziona un cliente", parent=self.root)
            return
        if not messagebox.askyesno("Conferma","Eliminare cliente e macchine collegate?", parent=self.root):
            return
        _safe_call(database.delete_clients, self.selected_cliente_id)
        messagebox.showinfo("OK","Cliente eliminato", parent=self.root)
        self.load_clienti(show_all=True)
        self.selected_cliente_id = None

    def gui_duplicate_cliente(self):
        # Obbliga a selezionare una riga prima di duplicare
        sel = self.tree.selection()
        if not sel:
            messagebox.showerror("Errore","Seleziona un cliente dalla tabella", parent=self.root)
            return
        if not self.selected_cliente_id:
            messagebox.showerror("Errore","Seleziona un cliente", parent=self.root)
            return
        nid = _safe_call(database.duplicate_clients, self.selected_cliente_id)
        if nid:
            messagebox.showinfo("OK","Cliente duplicato", parent=self.root)
            self.load_clienti(show_all=True)

    def gui_import_csv_clienti(self, path: str | None = None):
        """
        Se path è None -> apre filedialog. Altrimenti usa il path passato.
        Gestisce le righe saltate mostrando un dialog se necessario.
        Ritorna un dict con {ok, inserted, skipped, ...} oppure None in caso di errore.
        """
        if path is None:
            path = filedialog.askopenfilename(title="Seleziona CSV", filetypes=[("CSV files","*.csv"),("All","*.*")])
            if not path:
                return None
        
        # try ; then ,
        result = _safe_call(database.import_clients_from_csv_generic, path, ";", "latin-1")
        if result is None or not result.get("ok"):
            result = _safe_call(database.import_clients_from_csv_generic, path, ",", "latin-1")
        
        return result

    def gui_import_xls_clienti(self, path: str | None = None):
        """
        Se path è None -> apre filedialog. Altrimenti usa il path passato.
        Gestisce il nuovo formato dict di ritorno con righe saltate.
        """
        if path is None:
            path = filedialog.askopenfilename(title="Seleziona XLS/XLSX", filetypes=[("Excel files","*.xls;*.xlsx"),("All","*.*")])
            if not path:
                return None
        
        result = _safe_call(database.import_excel_clienti_via_soffice, path)
        return result

    def gui_import_dbf_clienti(self, path: str | None = None):
        """
        Se path è None -> apre filedialog. Altrimenti usa il path passato.
        Gestisce il nuovo formato dict di ritorno con righe saltate.
        """
        if path is None:
            path = filedialog.askopenfilename(title="Seleziona DBF", filetypes=[("DBF files","*.dbf *.DBF"),("All","*.*")])
            if not path:
                return None
        result = _safe_call(database.import_clients_from_dbf, path)
        return result

    def gui_import_any_clienti(self):
        path = filedialog.askopenfilename(
            title="Seleziona file da importare",
            filetypes=[("CSV","*.csv"), ("Excel","*.xls *.xlsx"), ("DBF","*.dbf *.DBF"), ("Tutti i file","*.*")]
        )
        if not path:
            return

        ext = os.path.splitext(path)[1].lower()
        try:
            result = None
            if ext == ".csv":
                result = self.gui_import_csv_clienti(path)
            elif ext in (".xls", ".xlsx"):
                result = self.gui_import_xls_clienti(path)
            elif ext == ".dbf":
                result = self.gui_import_dbf_clienti(path)
            else:
                if messagebox.askyesno("Importazione", "Estensione non riconosciuta. Trattare come CSV?", parent=self.root):
                    self.root.lift()
                    self.root.focus_force()
                    result = self.gui_import_csv_clienti(path)
                elif messagebox.askyesno("Importazione", "Provare come file Excel (soffice)?", parent=self.root):
                    self.root.lift()
                    self.root.focus_force()
                    result = self.gui_import_xls_clienti(path)
                elif messagebox.askyesno("Importazione", "Provare come file DBF?", parent=self.root):
                    self.root.lift()
                    self.root.focus_force()
                    result = self.gui_import_dbf_clienti(path)
                else:
                    return

            # Aggiorna e ricarica
            _safe_call(database.aggiorna_finemf_e_garanzia)
            self.load_clienti(show_all=True)
            res = database.aggiorna_finemf_e_garanzia()
            print(res)

            try:
                self.load_clienti(show_all=True)
            except Exception as e:
                print("load_clienti failed:", e)
                try:
                    self.carica_clienti_da_db()
                    self.refresh_tabella()
                except Exception as e2:
                    print("fallback load failed:", e2)

            # NOTIFICA: mostra un solo popup riassuntivo
            try:
                msg = None
                if result is None:
                    msg = "Import annullato o nessun record inserito."
                elif isinstance(result, dict):
                    if result.get("ok"):
                        inserted = result.get("inserted", 0)
                        skipped = result.get("skipped", 0)
                        updated = result.get("updated", 0)
                        msg = f"Import completato: {inserted} inseriti"
                        if updated:
                            msg += f", {updated} aggiornati"
                        if skipped:
                            msg += f", {skipped} saltati"
                    else:
                        msg = f"Errore: {result.get('message', 'Errore sconosciuto')}"
                elif isinstance(result, tuple):
                    if len(result) == 2:
                        inserted, skipped = result
                        msg = f"Import completato: {inserted} inseriti, {skipped} saltati"
                    elif len(result) == 3:
                        inserted, skipped, updated = result
                        msg = f"Import completato: {inserted} inseriti, {updated} aggiornati, {skipped} saltati"
                    else:
                        msg = f"Import completato: {result}"
                else:
                    msg = f"Import completato: {result} record importati."
                if msg:
                    messagebox.showinfo("Importazione", msg, parent=self.root)
                    self.root.lift()
                    self.root.focus_force()
            except Exception as e:
                print("Errore mostrando popup:", e)

        except Exception as e:
            # mostra un errore esplicito (messagebox importato globalmente)
            messagebox.showerror("Errore importazione", f"Import fallito: {e}", parent=self.root)
            self.root.lift()
            self.root.focus_force()
                
    def _save_clienti_col_widths(self, tree):
        # salva larghezze correnti in self.clienti_col_widths
        if tree is None:
            return
        cols = self.last_clienti_cols or []
        for c in cols:
            try:
                # tree.column(c) ritorna un dict; prendiamo 'width'
                w = tree.column(c)['width']
                self.clienti_col_widths[c] = w
            except Exception:
                pass
    
    def gui_deactivate_machine(self):
        """
        Disattiva la macchina della riga selezionata nella tabella clienti.
        Estrae la matricola dalla riga selezionata, chiede la data (default oggi),
        aggiorna DATADIS e azzera VF.
        """
        # Verifica che una riga sia selezionata
        if not self.selected_cliente_id:
            messagebox.showwarning("Nessuna selezione", "Seleziona una riga dalla tabella clienti.", parent=self.root)
            return

        # Estrai i dati della riga selezionata
        sel = self.tree.selection()
        if not sel:
            messagebox.showwarning("Nessuna selezione", "Seleziona una riga dalla tabella clienti.", parent=self.root)
            return

        vals = self.tree.item(sel[0])["values"]
        cols = self.last_clienti_cols or database.get_clienti_columns()

        # Trova l'indice della colonna MATRICOLA
        matricola_index = None
        for i, c in enumerate(cols):
            if c == "MATRICOLA":
                matricola_index = i
                break

        if matricola_index is None or matricola_index >= len(vals):
            messagebox.showerror("Errore", "Impossibile trovare la MATRICOLA nella riga selezionata.", parent=self.root)
            return

        matricola = vals[matricola_index].strip()
        if not matricola:
            messagebox.showwarning("Errore", "La MATRICOLA nella riga selezionata è vuota.", parent=self.root)
            return

        # Crea dialog per la data
        top = tk.Toplevel(self.root)
        top.title("Disattiva Macchina")
        top.transient(self.root)
        top.grab_set()
        top.lift()
        top.focus_force()

        ttk.Label(top, text=f"Matricola: {matricola}").grid(row=0, column=0, columnspan=2, padx=6, pady=6, sticky="w")
        ttk.Label(top, text="Data (DDMMYYYY o DD/MM/YYYY) [vuoto = oggi]:").grid(row=1, column=0, padx=6, pady=6, sticky="w")
        ent_date = _create_date_entry(top, width=20, max_length=8)
        ent_date.grid(row=1, column=1, padx=6, pady=6)

        # funzione chiamata su conferma
        def do_deactivate():
            date_param = ent_date.get().strip() or None

            # Se non inserita una data, usa oggi
            if not date_param:
                date_param = date.today().strftime("%d/%m/%Y")

            # conferma
            if not messagebox.askyesno("Conferma", f"Confermi disattivazione matricola '{matricola}' (data={date_param})?", parent=top):
                return

            # chiama la funzione di backend
            try:
                result = database.deactivate_machine(matricola, date_param)
                messagebox.showinfo("Completato", "Macchina disattivata: DATADIS impostata e VF azzerato.", parent=top)
                # Ricarica la tabella clienti
                self.load_clienti(show_all=True)
            except Exception as e:
                messagebox.showerror("Errore", f"Errore durante la disattivazione:\n{e}", parent=top)
                try:
                    top.lift()
                    top.focus_force()
                except Exception:
                    pass
                return

            top.destroy()

        # bottoni
        btn_frame = ttk.Frame(top)
        btn_frame.grid(row=2, column=0, columnspan=2, pady=(6,10))
        ttk.Button(btn_frame, text="Annulla", command=top.destroy).pack(side="left", padx=6)
        ttk.Button(btn_frame, text="Conferma disattivazione", command=do_deactivate).pack(side="left", padx=6)

        # focus sul campo data
        ent_date.focus_set()

        # gestire Enter per conferma rapida
        def on_enter(event):
            do_deactivate()
        ent_date.bind("<Return>", on_enter)

    def gui_fill_and_print_checklist(self):
        """
        Usa template.pdf nella cartella dell'app, mostra preview 
        (apre il PDF generato con il viewer di sistema) e permette stampa.
        """
        dlg = tk.Toplevel(self.root)
        dlg.title("Seleziona periodo checklist")
        dlg.transient(self.root)
        dlg.grab_set()
        dlg.lift()
        dlg.focus_force()

        ttk.Label(dlg, text="Mese:").grid(row=0, column=0, padx=6, pady=6, sticky="w")
        months = [f"{i+1:02d} - {m}" for i, m in enumerate([
            "Gennaio","Febbraio","Marzo","Aprile","Maggio","Giugno",
            "Luglio","Agosto","Settembre","Ottobre","Novembre","Dicembre"
        ])]
        cmb_month = ttk.Combobox(dlg, values=months, state="readonly", width=20)
        cmb_month.grid(row=0, column=1, padx=6, pady=6)
        try:
            cmb_month.current(date.today().month - 1)
        except Exception:
            pass

        ttk.Label(dlg, text="Anno:").grid(row=1, column=0, padx=6, pady=6, sticky="w")
        ent_year = ttk.Entry(dlg, width=10)
        ent_year.grid(row=1, column=1, padx=6, pady=6, sticky="w")
        ent_year.insert(0, str(date.today().year))

        info_lbl = ttk.Label(dlg, text=f"Usando template: {os.path.join(os.getcwd(), 'template.pdf')}")
        info_lbl.grid(row=2, column=0, columnspan=2, padx=6, pady=6, sticky="w")

        def do_search():
            raw_m = cmb_month.get().strip()
            raw_y = ent_year.get().strip()
            try:
                month = int(raw_m.split()[0])
                year = int(raw_y)
            except Exception:
                messagebox.showerror("Errore", "Mese/anno non validi", parent=dlg)
                try:
                    dlg.lift()
                    dlg.focus_force()
                except Exception:
                    pass
                return

            tpl = os.path.join(os.getcwd(), "template.pdf")
            if not os.path.exists(tpl):
                messagebox.showerror("Template mancante", f"Non trovo {tpl}.", parent=dlg)
                try:
                    dlg.lift()
                    dlg.focus_force()
                except Exception:
                    pass
                return

            try:
                clients = database.get_clients_for_vf_month_year(month, year) or []
            except Exception as e:
                messagebox.showerror("Errore DB", f"Errore recupero clienti: {e}", parent=dlg)
                try:
                    dlg.lift()
                    dlg.focus_force()
                except Exception:
                    pass
                return

            if not clients:
                messagebox.showinfo("Nessun cliente", f"Nessun cliente trovato per {month:02d}/{year}", parent=dlg)
                return

            # Chiudi il dialog di input subito
            dlg.destroy()

            pv = tk.Toplevel(self.root)
            pv.title(f"Preview checklist {month:02d}/{year}")
            pv.geometry("1100x520")
            pv.transient(self.root)
            pv.grab_set()
            pv.lift()
            pv.focus_force()
            frame = ttk.Frame(pv)
            frame.pack(fill="both", expand=True, padx=6, pady=6)
            
            # **AGGIUNTA**: Colonna P_IVA e MARCA nella preview
            cols = ("RAG_SOC", "P_IVA", "INDIRIZZO", "COMUNE", "MATRICOLA", "MARCA", "VF")
            tree = ttk.Treeview(frame, columns=cols, show="headings", selectmode="extended")
            vsb = ttk.Scrollbar(frame, orient="vertical", command=tree.yview)
            hsb = ttk.Scrollbar(frame, orient="horizontal", command=tree.xview)
            tree.configure(yscrollcommand=vsb.set, xscrollcommand=hsb.set)
            tree.grid(row=0, column=0, sticky="nsew")
            vsb.grid(row=0, column=1, sticky="ns")
            hsb.grid(row=1, column=0, sticky="ew")
            frame.rowconfigure(0, weight=1)
            frame.columnconfigure(0, weight=1)
            
            for c in cols:
                tree.heading(c, text=c)
                tree.column(c, width=140)

            produced_map = {}
            for c in clients:
                vals = (
                    c.get("RAG_SOC", ""),
                    c.get("P_IVA", "") or c.get("PIVA", ""),           # <- NUOVO
                    c.get("IND_VIA", "") or c.get("INDIRIZZO", ""),
                    c.get("COMUNE", ""),
                    c.get("MATRICOLA", ""),
                    c.get("MARCA", ""),                                 # <- NUOVO
                    c.get("VF", "")
                )
                iid = tree.insert("", "end", values=vals)
                produced_map[iid] = c

            # pulsanti: Preview PDF, Stampa selezionati, Stampa tutti
            btns = ttk.Frame(pv); btns.pack(fill="x", pady=6)

            tempdir = None  # usato per preview/printing

            def ensure_generated_for_selection(selected_iids):
                nonlocal tempdir
                import tempfile
                if tempdir is None:
                    tempdir = tempfile.mkdtemp(prefix="checklist_")
                try:
                    files = database.generate_filled_pdfs_for_month(year, month, tpl, tempdir, auto_print=False)
                except Exception as e:
                    messagebox.showerror("Errore generazione", f"Errore: {e}", parent=pv)
                    try:
                        pv.lift()
                        pv.focus_force()
                    except Exception:
                        pass
                    return None
                if not selected_iids:
                    return files
                mats = set()
                for iid in selected_iids:
                    c = produced_map.get(iid, {})
                    mats.add((c.get("MATRICOLA","") or "").strip())
                selected_files = [p for p in files if any(m and m in os.path.basename(p) for m in mats)]
                if not selected_files:
                    selected_files = files
                return selected_files

            def do_preview_and_print():
                sel = tree.selection()
                if not sel:
                    messagebox.showwarning("Nessuna selezione", "Seleziona almeno un cliente", parent=pv)
                    return
                files = ensure_generated_for_selection(sel)
                if not files:
                    return
                for i, pdf_file in enumerate(files):
                    msg = f"File {i+1} di {len(files)}"
                    try:
                        if os.name == "nt":
                            os.startfile(pdf_file)
                        else:
                            subprocess.run(["xdg-open", pdf_file], check=False)
                    except Exception:
                        messagebox.showinfo("Preview", f"PDF: {pdf_file}\n({msg})", parent=pv)
                    if i < len(files) - 1:
                        if not messagebox.askyesno("Continua preview", f"Mostrare il file successivo?\n({msg})", parent=pv):
                            break
                cliente_list = []
                for iid in sel:
                    c = produced_map.get(iid, {})
                    rag_soc = c.get("RAG_SOC", "")
                    matricola = c.get("MATRICOLA", "")
                    cliente_list.append(f"{rag_soc} ({matricola})")
                clienti_str = "\n".join(cliente_list)
                msg = f"Stampare le checklist dei seguenti clienti in doppia copia?\n\n{clienti_str}"
                if not messagebox.askyesno("Conferma stampa", msg, parent=pv):
                    return
                for pdf_file in files:
                    try:
                        database.print_pdf_file(pdf_file, copies=2)
                    except Exception as e:
                        messagebox.showerror("Errore stampa", f"Errore stampando {pdf_file}:\n{e}", parent=pv)
                        return

            def do_print_all():
                files = ensure_generated_for_selection([])
                if not files:
                    return
                clienti_str = "\n".join([
                    f"{c.get('RAG_SOC','')} ({c.get('MATRICOLA','')})" for c in clients
                ])
                msg = f"Stampare le checklist di TUTTI i clienti in doppia copia?\n\n{clienti_str}"
                if not messagebox.askyesno("Conferma stampa", msg, parent=pv):
                    return
                for pdf_file in files:
                    try:
                        database.print_pdf_file(pdf_file, copies=2)
                    except Exception as e:
                        messagebox.showerror("Errore stampa", f"Errore stampando {pdf_file}:\n{e}", parent=pv)
                        return

            ttk.Button(btns, text="Preview e Stampa selezionati", command=do_preview_and_print).pack(side="left", padx=8)
            ttk.Button(btns, text="Stampa tutti", command=do_print_all).pack(side="left", padx=8)
            ttk.Button(btns, text="Chiudi", command=pv.destroy).pack(side="right", padx=6)

        # dialog buttons
        btnf = ttk.Frame(dlg)
        btnf.grid(row=4, column=0, columnspan=3, pady=8)
        ttk.Button(btnf, text="Annulla", command=dlg.destroy).pack(side="left", padx=6)
        ttk.Button(btnf, text="CERCA", command=do_search).pack(side="left", padx=6)
                
# -------------------------
# Tab INTAE (import & search)
# -------------------------
    def _build_tab_intae(self):
        frm = ttk.Frame(self.tab_intae, padding=8)
        frm.pack(fill="x")
        ttk.Label(frm, text="Import & ricerca INTAE", font=("Arial",12)).pack(anchor="w")

        btn_fr = ttk.Frame(frm)
        btn_fr.pack(fill="x", pady=6)
        ttk.Button(btn_fr, text="IMPORTA FILE", command=self.gui_import_csv_intae).pack(side="left", padx=6)
        ttk.Button(btn_fr, text="CERCA PER MATRICOLA", command=self.search_intae_matricola).pack(side="left", padx=6)
        ttk.Button(btn_fr, text="REFRESH", command=self.refresh_intae).pack(side="left", padx=6)
        
        # SOLO colonne reali, senza id
        cols = database.INTAE_COLUMNS[:]    

        # Treeview con scrollbars
        container = ttk.Frame(self.tab_intae)
        container.pack(fill="both", expand=True)

        scroll_y = ttk.Scrollbar(container, orient="vertical")
        scroll_x = ttk.Scrollbar(container, orient="horizontal")

        self.tree_intae = ttk.Treeview(
            container,
            columns=cols,
            show="headings",
            yscrollcommand=scroll_y.set,
            xscrollcommand=scroll_x.set
        )

        scroll_y.config(command=self.tree_intae.yview)
        scroll_x.config(command=self.tree_intae.xview)

        scroll_y.pack(side="right", fill="y")
        scroll_x.pack(side="bottom", fill="x")
        self.tree_intae.pack(side="left", fill="both", expand=True)

        # intestazioni con larghezze adatte ai diversi campi
        col_widths = {
            "CODICE_INTERVENTO": 120,
            "PIVA": 120,
            "DATA_INTERVENTO": 150,
            "CF_TECNICO": 130,
            "STATO_CF": 100,
            "DISPOSITIVO": 100,
            "MATRICOLA": 130,
            "TIPOLOGIA": 100,
            "DATA_TRASMISSIONE": 150,
            "NOTA": 250,
            "VERSIONE_SW": 120,
            "DATA_RELEASE": 130,
        }
        
        for c in cols:
            self.tree_intae.heading(c, text=c.replace("_", " ").title())
            width = col_widths.get(c, 130)
            self.tree_intae.column(c, width=width, minwidth=80)

        self.refresh_intae()

    def refresh_intae(self):
        for row in self.tree_intae.get_children():
            self.tree_intae.delete(row)

        rows = _safe_call(database.read_intae) or []

        for r in rows:
            self.tree_intae.insert("", "end", iid=r[0], values=r[1:])


    def search_intae_matricola(self):
        """Ricerca interventi INTAE per matricola - gestisce matricole con/senza logotipo"""
        def do_search():
            mat = ent.get().strip()
            w.destroy()
            if not mat:
                return
            
            # Usa la funzione dedicata nel database per ricerca robusta
            filtered = _safe_call(database.cerca_interventi_intae, mat) or []
            
            # Aggiorna la tree con i risultati
            for r in self.tree_intae.get_children():
                self.tree_intae.delete(r)
            
            if not filtered:
                messagebox.showinfo("Nessun risultato", 
                                   f"Nessun intervento trovato per la matricola '{mat}'.\n\n"
                                   f"Suggerimento: Le matricole in INTAE possono avere un logotipo (es: '2CMST004032').\n"
                                   f"Prova a cercare solo parte della matricola (es: 'MST004032').",
                                   parent=self.root)
                return
            
            for row in filtered:
                iid = row[0] if row else None
                vals = row[1:] if len(row) > 1 else row
                if iid is not None:
                    self.tree_intae.insert("", "end", iid=iid, values=vals)
                else:
                    self.tree_intae.insert("", "end", values=vals)
            
            # Mostra messaggio con quante righe trovate
            messagebox.showinfo("Risultati", f"Trovati {len(filtered)} interventi per la matricola '{mat}'", parent=self.root)
        
        w = tk.Toplevel(self.root)
        w.title("Cerca Intervento INTAE")
        w.transient(self.root)
        w.grab_set()
        w.lift()
        w.focus_force()
        ttk.Label(w, text="Matricola:").pack(side="left", padx=6, pady=6)
        ent = ttk.Entry(w)
        ent.pack(side="left", padx=6, pady=6)
        ttk.Button(w, text="CERCA", command=do_search).pack(side="left", padx=6, pady=6)
        ent.focus_set()

    def gui_import_csv_intae(self):
        path = filedialog.askopenfilename(
            title="Seleziona file INTAE da importare",
            filetypes=[("CSV","*.csv"), ("Excel","*.xls *.xlsx"), ("DBF","*.dbf *.DBF"), ("Tutti i file","*.*")]
        )
        if not path:
            return

        ext = os.path.splitext(path)[1].lower()
        result = None

        try:
            if ext == ".csv":
                # Prova prima utf-8, poi latin-1, e delimitatori ; e ,
                result = _safe_call(database.import_csv_intae, path, ";", "utf-8")
                if not result or not result.get("ok"):
                    result = _safe_call(database.import_csv_intae, path, ",", "utf-8")
                if not result or not result.get("ok"):
                    result = _safe_call(database.import_csv_intae, path, ";", "latin-1")
                if not result or not result.get("ok"):
                    result = _safe_call(database.import_csv_intae, path, ",", "latin-1")
            elif ext in (".xls", ".xlsx"):
                messagebox.showinfo("Formato", "Convertire prima il file in CSV e riprovare.")
                return
            elif ext == ".dbf":
                result = _safe_call(database.import_intae_from_dbf, path)
            else:
                messagebox.showerror("Errore", "Formato file non supportato. Usare CSV o DBF.")
                return

            if result is not None and isinstance(result, dict):
                if result.get("ok"):
                    inserted = result.get('inserted', 0)
                    skipped = result.get('skipped', 0)
                    msg = f"Import INTAE completato: {inserted} inseriti"
                    if skipped:
                        msg += f", {skipped} saltati"
                    messagebox.showinfo("Import INTAE", msg, parent=self.root)
                    self.refresh_intae()
                else:
                    messagebox.showerror("Errore importazione", result.get("message", "Errore sconosciuto"), parent=self.root)
            else:
                messagebox.showerror("Errore", "Import non riuscito o formato non supportato.", parent=self.root)

        except Exception as e:
            messagebox.showerror("Errore importazione", f"Import fallito: {e}", parent=self.root)

# -------------------------
# Tab MISURART
# -------------------------
    def _build_tab_misurart(self):
        frm = ttk.Frame(self.tab_misurart, padding=8)
        frm.pack(fill="x")
        ttk.Label(frm, text="Registrazioni MISURART", font=("Arial",12)).pack(anchor="w")

        # --- BLOCCO VERIFICA MATRICOLA (NUOVO) ---
        verify_fr = ttk.LabelFrame(
            self.tab_misurart,
            text="Verifica strumento da INTAE",
            padding=8
        )
        verify_fr.pack(fill="x", padx=10, pady=6)

        ttk.Label(verify_fr, text="Matricola:").grid(row=0, column=0, sticky="w", padx=4, pady=2)
        self.verify_matricola_entry = ttk.Entry(verify_fr, width=25)
        self.verify_matricola_entry.grid(row=0, column=1, padx=4, pady=2)

        ttk.Button(
            verify_fr,
            text="VERIFICA E COMPILA",
            command=self.gui_verifica_matricola
        ).grid(row=0, column=2, padx=10)
        # --- FINE BLOCCO VERIFICA MATRICOLA ---
        
        # Form inserimento
        form = ttk.LabelFrame(self.tab_misurart, text="Nuova / Modifica registrazione", padding=8)
        form.pack(fill="x", padx=10, pady=10)

        self.misurart_fields = {}
        cols = database.MISURART_COLUMNS[:]   # solo colonne reali

        for i, key in enumerate(cols):
            r = i // 2
            c = (i % 2) * 2
            ttk.Label(form, text=key).grid(row=r, column=c, sticky="w", padx=4, pady=2)
            # Usa entry data formattata per i campi data
            if key in ("DATAINIZIO", "DATAFINE", "DATAVP"):
                ent = _create_date_entry(form, width=30, max_length=8)
            else:
                ent = ttk.Entry(form, width=30)
            ent.grid(row=r, column=c+1, padx=4, pady=2)
            self.misurart_fields[key] = ent

        # Bottoni
        btns = ttk.Frame(form)
        btns.grid(row=(len(cols)//2)+2, column=0, columnspan=4, pady=8)
        ttk.Button(btns, text="INSERISCI", command=self.misurart_add).grid(row=0, column=0, padx=6)
        ttk.Button(btns, text="AGGIORNA", command=self.misurart_update).grid(row=0, column=1, padx=6)
        ttk.Button(btns, text="ELIMINA", command=self.misurart_delete).grid(row=0, column=2, padx=6)

        # Ricerca
        search_fr = ttk.Frame(self.tab_misurart)
        search_fr.pack(fill="x", pady=6)
        ttk.Label(search_fr, text="Cerca matricola:").pack(side="left")
        self.misurart_search_entry = ttk.Entry(search_fr, width=30)
        self.misurart_search_entry.pack(side="left", padx=6)
        ttk.Button(search_fr, text="CERCA", command=self.misurart_search).pack(side="left", padx=6)
        ttk.Button(search_fr, text="MOSTRA TUTTI", command=self.misurart_load_all).pack(side="left", padx=6)
        ttk.Button(search_fr, text="IMPORTA FILE", command=self.gui_import_misurart_file).pack(side="left", padx=6)

        # Tabelle con scrollbars
        tree_cols = cols[:]   # nessun ID visibile

        container = ttk.Frame(self.tab_misurart)
        container.pack(fill="both", expand=True)

        scroll_y = ttk.Scrollbar(container, orient="vertical")
        scroll_x = ttk.Scrollbar(container, orient="horizontal")

        self.tree_misurart = ttk.Treeview(
            container,
            columns=tree_cols,
            show="headings",
            yscrollcommand=scroll_y.set,
            xscrollcommand=scroll_x.set
        )

        scroll_y.config(command=self.tree_misurart.yview)
        scroll_x.config(command=self.tree_misurart.xview)

        scroll_y.pack(side="right", fill="y")
        scroll_x.pack(side="bottom", fill="x")
        self.tree_misurart.pack(side="left", fill="both", expand=True)

        for c in tree_cols:
            self.tree_misurart.heading(c, text=c.replace("_", " ").title(), anchor="w")
            # stretch=False evita il ridimensionamento automatico delle altre colonne
            self.tree_misurart.column(c, width=140, minwidth=20, stretch=False, anchor="w")

        self.tree_misurart.bind("<<TreeviewSelect>>", self.misurart_on_select)
        self.misurart_selected_id = None

        self.misurart_load_all()

    def misurart_load_all(self):
        for row in self.tree_misurart.get_children():
            self.tree_misurart.delete(row)

        rows = _safe_call(database.read_misurart_all) or []

        for r in rows:
            self.tree_misurart.insert("", "end", iid=r[0], values=r[1:])

    def misurart_add(self):
        # Non permettere inserimento se una riga è già selezionata
        if getattr(self, "misurart_selected_id", None):
            messagebox.showerror("Errore", "Deseleziona la riga prima di inserire una nuova registrazione", parent=self.root)
            return
        
        # raccogli dati dalla mappa di entry
        data = {k: self.misurart_fields[k].get().strip() for k in self.misurart_fields}
        
        # Controlla che almeno un campo sia stato riempito
        if not any(data.values()):
            messagebox.showwarning("Errore", "Inserisci almeno un campo per creare una registrazione", parent=self.root)
            try:
                self.root.lift()
                self.root.focus_force()
            except Exception:
                pass
            return
        
        nid = _safe_call(database.misurart_add_from_map, data)
        if nid:
            messagebox.showinfo("OK", "Registrazione inserita", parent=self.root)
            # Pulisci i campi e ricarica
            for key in self.misurart_fields:
                self.misurart_fields[key].delete(0, tk.END)
            self.misurart_selected_id = None
            self.misurart_load_all()

    def misurart_on_select(self, event):
        sel = self.tree_misurart.selection()
        if not sel:
            self.misurart_selected_id = None
            return
        item_id = sel[0]
        self.misurart_selected_id = item_id  # L'ID è memorizzato come iid nel treeview
        values = self.tree_misurart.item(item_id)["values"]
        # riempi i campi: valori sono in ordine MISURART_COLUMNS (senza ID)
        for i, key in enumerate(database.MISURART_COLUMNS):
            val = values[i] if len(values) > i else ""
            self.misurart_fields[key].delete(0, tk.END)
            self.misurart_fields[key].insert(0, val or "")

    def misurart_update(self):
        if not getattr(self, "misurart_selected_id", None):
            messagebox.showerror("Errore", "Seleziona una registrazione da aggiornare.", parent=self.root)
            return
        data = {k: self.misurart_fields[k].get().strip() for k in self.misurart_fields}
        _safe_call(database.misurart_update, self.misurart_selected_id, data)
        messagebox.showinfo("OK", "Registrazione aggiornata", parent=self.root)
        # Pulisci i campi e ricarica
        for key in self.misurart_fields:
            self.misurart_fields[key].delete(0, tk.END)
        self.misurart_selected_id = None
        self.misurart_load_all()

    def misurart_delete(self):
        if not getattr(self, "misurart_selected_id", None):
            messagebox.showerror("Errore", "Seleziona una registrazione da eliminare.", parent=self.root)
            return
        if not messagebox.askyesno("Conferma", "Eliminare la registrazione?", parent=self.root):
            return
        _safe_call(database.misurart_delete, self.misurart_selected_id)
        messagebox.showinfo("OK", "Registrazione eliminata", parent=self.root)
        # Pulisci i campi dopo eliminazione
        for key in self.misurart_fields:
            self.misurart_fields[key].delete(0, tk.END)
        self.misurart_load_all()
        self.misurart_selected_id = None

    def misurart_search(self):
        term = self.misurart_search_entry.get().strip()
        if not term:
            messagebox.showerror("Errore", "Inserisci una matricola da cercare", parent=self.root)
            return
        rows = _safe_call(database.misurart_search_by_matricola, term) or []
        for r in self.tree_misurart.get_children(): self.tree_misurart.delete(r)
        for row in rows:
            if row:
                self.tree_misurart.insert("", "end", iid=row[0], values=row[1:])

    def gui_import_misurart_file(self):
        """
        Apre dialog per importare file MISURART (CSV, XLSX, DBF).
        """
        path = filedialog.askopenfilename(
            title="Seleziona file MISURART da importare",
            filetypes=[("CSV","*.csv"), ("Excel","*.xls *.xlsx"), ("DBF","*.dbf *.DBF"), ("Tutti i file","*.*")]
        )
        if not path:
            return
        
        ext = os.path.splitext(path)[1].lower()
        result = None
        
        try:
            if ext == ".csv":
                result = _safe_call(database.import_misurart_from_csv, path, ";", "utf-8")
                if result is None:
                    # Riprova con altro encoding
                    result = _safe_call(database.import_misurart_from_csv, path, ";", "latin-1")
            elif ext in (".xls", ".xlsx"):
                result = _safe_call(database.import_misurart_from_xlsx, path)
            elif ext == ".dbf":
                result = _safe_call(database.import_misurart_from_dbf, path)
            else:
                if messagebox.askyesno("Importazione", "Estensione non riconosciuta. Trattare come CSV?"):
                    result = _safe_call(database.import_misurart_from_csv, path, ";", "utf-8")
                    if result is None:
                        result = _safe_call(database.import_misurart_from_csv, path, ";", "latin-1")
                elif messagebox.askyesno("Importazione", "Provare come file Excel?"):
                    result = _safe_call(database.import_misurart_from_xlsx, path)
                elif messagebox.askyesno("Importazione", "Provare come file DBF?"):
                    result = _safe_call(database.import_misurart_from_dbf, path)
                else:
                    return
            
            if result is not None:
                if isinstance(result, dict):
                    # Se c'è un warning, mostra un avviso ma consente di continuare
                    if result.get("warning"):
                        messagebox.showwarning("Avviso", result.get("warning"), parent=self.root)
                    
                    if not result.get("ok"):
                        messagebox.showerror("Errore importazione", result.get("message", "Errore sconosciuto"), parent=self.root)
                        return
                    inserted = result.get("inserted", 0)
                    
                    msg = f"Importati: {inserted} record MISURART"
                    messagebox.showinfo("Import MISURART", msg, parent=self.root)

                    self.misurart_load_all()
                else:
                    messagebox.showinfo("Import", f"Import completato", parent=self.root)
                    self.misurart_load_all()
        
        except Exception as e:
            messagebox.showerror("Errore", f"Errore durante l'importazione:\n{str(e)}", parent=self.root)

    def gui_verifica_matricola(self):
        """
        Apre una finestra completa di verifica matricola con:
        - Visualizzazione dati cliente e macchina
        - Visualizzazione dati INTAE (verificatore e data)
        - Modifica campi
        - Inserimento tipo intervento
        - Compilazione automatica MISURART
        - Aggiornamento VF e AZZ in clienti
        """
        matricola_input = self.verify_matricola_entry.get().strip()
        if not matricola_input:
            messagebox.showerror("Errore", "Inserisci una matricola", parent=self.root)
            return
        
        # Recupera i dati completi dalla DB
        verif_data = database.get_verification_data(matricola_input)
        
        if not verif_data.get("ok"):
            messagebox.showerror("Errore", verif_data.get("error", "Errore sconosciuto"), parent=self.root)
            return
        
        cliente = verif_data.get("cliente", {})
        intae = verif_data.get("intae_data", {})
        cliente_id = cliente.get("id")
        
        if not cliente_id:
            messagebox.showerror("Errore", "ID cliente non trovato", parent=self.root)
            return
        
        # Crea finestra di verifica
        w = tk.Toplevel(self.root)
        w.title("Verifica Macchina")
        w.geometry("1000x700")
        w.transient(self.root)
        w.grab_set()
        w.lift()
        w.focus_force()
        
        # ===== CONTENITORE PRINCIPALE CON SCROLLBAR =====
        main_frame = ttk.Frame(w)
        main_frame.pack(fill="both", expand=True, padx=8, pady=8)
        
        canvas = tk.Canvas(main_frame)
        scrollbar = ttk.Scrollbar(main_frame, orient="vertical", command=canvas.yview)
        scroll_frame = ttk.Frame(canvas)
        scroll_frame.bind("<Configure>", lambda e: canvas.configure(scrollregion=canvas.bbox("all")))
        canvas.create_window((0, 0), window=scroll_frame, anchor="nw")
        canvas.configure(yscrollcommand=scrollbar.set)
        
        # Bind mouse wheel scroll
        def _on_mousewheel(event):
            canvas.yview_scroll(int(-1*(event.delta/120)), "units")
        
        canvas.bind("<MouseWheel>", _on_mousewheel)
        canvas.bind("<Button-4>", lambda e: canvas.yview_scroll(-1, "units"))
        canvas.bind("<Button-5>", lambda e: canvas.yview_scroll(1, "units"))
        
        canvas.pack(side="left", fill="both", expand=True)
        scrollbar.pack(side="right", fill="y")
        
        # ===== SEZIONE 1: DATI CLIENTE =====
        cliente_lbl = ttk.LabelFrame(scroll_frame, text="Dati Cliente", padding=6)
        cliente_lbl.pack(fill="x", pady=6)
        
        client_fields_editable = {}
        client_row = 0
        
        client_display_fields = [
            ("RAG_SOC", "Ragione Sociale", True),
            ("P_IVA", "P.IVA", False),
            ("INDIRIZZO", "Indirizzo", False),
            ("COMUNE", "Comune", False),
            ("PROV", "Provincia", False),
            ("CAP", "CAP", False),
            ("MF", "MF", False),
            ("MARCA", "Marca", False),
            ("MATRICOLA", "Matricola", False),
            ("DATAINST", "Data Installazione", False),
        ]
        
        for field_key, label_text, editable in client_display_fields:
            ttk.Label(cliente_lbl, text=label_text + ":").grid(row=client_row, column=0, sticky="w", padx=4, pady=2)
            
            if editable:
                ent = ttk.Entry(cliente_lbl, width=40)
                ent.insert(0, cliente.get(field_key, "") or "")
                ent.grid(row=client_row, column=1, sticky="w", padx=4, pady=2)
                client_fields_editable[field_key] = ent
            else:
                lbl = ttk.Label(cliente_lbl, text=cliente.get(field_key, "") or "")
                lbl.grid(row=client_row, column=1, sticky="w", padx=4, pady=2)
            
            client_row += 1
        
        # ===== SEZIONE 2: DATI INTAE (VERIFICATORE) =====
        intae_lbl = ttk.LabelFrame(scroll_frame, text="Dati Verificazione (da INTAE)", padding=6)
        intae_lbl.pack(fill="x", pady=6)
        
        intae_fields = [
            ("CODICE_INTERVENTO", "Codice Intervento"),
            ("DATA_INTERVENTO", "Data Intervento"),
            ("CF_TECNICO", "Verificatore (CF)"),
            ("STATO_CF", "Stato CF"),
            ("DISPOSITIVO", "Dispositivo"),
            ("TIPOLOGIA", "Tipologia"),
            ("VERSIONE_SW", "Versione SW"),
        ]
        
        intae_row = 0
        for field_key, label_text in intae_fields:
            ttk.Label(intae_lbl, text=label_text + ":").grid(row=intae_row, column=0, sticky="w", padx=4, pady=2)
            lbl = ttk.Label(intae_lbl, text=intae.get(field_key, "") or "")
            lbl.grid(row=intae_row, column=1, sticky="w", padx=4, pady=2)
            intae_row += 1
        
        # ===== SEZIONE 3: CAMPI PER AGGIORNAMENTO VF E AZZ =====
        update_lbl = ttk.LabelFrame(scroll_frame, text="Aggiornamento Verificazione", padding=6)
        update_lbl.pack(fill="x", pady=6)
        
        # VF precompilato: prendi l'ultima DATA_INTERVENTO da INTAE relativa alla matricola
        vf_default = ""
        print(f"\n{'='*80}")
        print(f"[DEBUG GUI VERIFICA] Inizio calcolo VF")
        print(f"[DEBUG GUI VERIFICA] intae dict completo: {intae}")
        print(f"[DEBUG GUI VERIFICA] intae keys: {list(intae.keys())}")
        print(f"[DEBUG GUI VERIFICA] intae.get('DATA_INTERVENTO'): {intae.get('DATA_INTERVENTO')}")
        print(f"{'='*80}\n")
        
        if intae and intae.get("DATA_INTERVENTO"):
            vf_raw = str(intae.get("DATA_INTERVENTO", "")).strip()
            print(f"[DEBUG GUI VERIFICA] DATA_INTERVENTO trovata: '{vf_raw}'")
            if vf_raw:
                try:
                    # Se vf_raw è già in formato ISO (YYYY-MM-DD), usa direttamente
                    if "-" in vf_raw and len(vf_raw) == 10:
                        dt = datetime.strptime(vf_raw, "%Y-%m-%d")
                        print(f"[DEBUG GUI VERIFICA] Formato riconosciuto come ISO YYYY-MM-DD")
                    else:
                        # Altrimenti usa _parse_date che può gestire vari formati
                        print(f"[DEBUG GUI VERIFICA] Formato non ISO, usando _parse_date")
                        parsed = database._parse_date(vf_raw)
                        # _parse_date può ritornare datetime object o stringa
                        if isinstance(parsed, datetime):
                            dt = parsed
                        elif isinstance(parsed, str):
                            # È una stringa, prova a parsarla
                            dt = datetime.strptime(parsed.split()[0], "%Y-%m-%d")
                        else:
                            dt = None
                    
                    if dt:
                        vf_default = dt.strftime("%d/%m/%Y")
                        print(f"[DEBUG GUI VERIFICA] vf_default calcolato: '{vf_default}'")
                except Exception as e:
                    print(f"[DEBUG GUI VERIFICA] ✗ Errore durante parsing: {e}")
                    import traceback
                    traceback.print_exc()
                    vf_default = ""
        else:
            print(f"[DEBUG GUI VERIFICA] ✗ DATA_INTERVENTO NON trovata o intae è vuoto")
        
        ttk.Label(update_lbl, text="VF (Data Verificazione):").grid(row=0, column=0, sticky="w", padx=4, pady=2)
        vf_entry, vf_var = _create_date_entry(update_lbl, width=20, max_length=8, with_var=True)
        vf_entry.grid(row=0, column=1, sticky="w", padx=4, pady=2)
        
        # Precompila il campo VF usando insert() (come per AZZ) - soluzione affidabile
        # Se vf_default è vuoto, usa oggi come fallback
        if not vf_default:
            print(f"[DEBUG GUI VERIFICA] DATA_INTERVENTO vuoto, uso oggi come fallback")
            vf_default = date.today().strftime("%d/%m/%Y")
        
        print(f"[DEBUG GUI VERIFICA] *** PRECOMPILAZIONE VF con '{vf_default}' ***")
        
        # Usa after() per assicurare che il widget sia renderizzato prima di inserire il valore
        def fill_vf():
            vf_entry.delete(0, tk.END)  # Pulisci prima
            vf_entry.insert(0, vf_default)
            print(f"[DEBUG GUI VERIFICA] ✓ Campo VF riempito. Contenuto: '{vf_entry.get()}'")
        
        w.after(50, fill_vf)  # Aspetta 50ms per far renderizzare il widget
        print(f"{'='*80}\n")
        
        ttk.Label(update_lbl, text="AZZ (Azzeramenti):").grid(row=1, column=0, sticky="w", padx=4, pady=2)
        azz_entry = ttk.Entry(update_lbl, width=20)
        azz_entry.insert(0, cliente.get("AZZ", "") or "")
        azz_entry.grid(row=1, column=1, sticky="w", padx=4, pady=2)
        azz_var = None  # AZZ non usa StringVar
        
        ttk.Label(update_lbl, text="TIPOINT (Tipo Intervento):").grid(row=2, column=0, sticky="w", padx=4, pady=2)
        tipoint_combo = ttk.Combobox(update_lbl, values=("1", "3"), state="readonly", width=10)
        tipoint_combo.set("3")  # default 3
        tipoint_combo.grid(row=2, column=1, sticky="w", padx=4, pady=2)
        
        # ===== SEZIONE 4: CAMPI MISURART =====
        misurart_lbl = ttk.LabelFrame(scroll_frame, text="Dati MISURART", padding=6)
        misurart_lbl.pack(fill="x", pady=6)
        
        misurart_fields_editable = {}
        misurart_row = 0
        
        # Precompila dati MISURART da cliente e INTAE - TUTTI i campi disponibili
        misurart_defaults = {
            "TIPO": "3",  # Verificazione periodica
            "CF": intae.get("CF_TECNICO", "") or "",  # Da INTAE - Codice Fiscale Verificatore
            "DATAINIZIO": vf_default,  # Stessa data VF
            "DATAFINE": vf_default,  # Stessa data VF
            "DATAVP": vf_default,  # Data Verificazione
            "ESITO": "POSITIVO",  # Default positivo
            "PIVAESERCE": cliente.get("P_IVA", "") or "",  # P.IVA Esercente
            "RAGSOCESE": cliente.get("RAG_SOC", "") or "",  # Ragione Sociale Esercente
            "COMUNESE": cliente.get("COMUNE", "") or "",  # Comune Esercente
            "PROVINCIA": cliente.get("PROV", "") or "",  # Provincia Esercente
            "INDIRIZZO": cliente.get("IND_VIA", "") or cliente.get("INDIRIZZO", "") or "",  # Indirizzo Esercente
            "CAPESE": cliente.get("CAP", "") or "",  # CAP Esercente
            "MARCAMF": cliente.get("MARCA", "") or "",  # Marca Macchina da CLIENTI
            "MODELLO": cliente.get("MF", "") or "",  # Modello da CLIENTI
            "LOGOTIPO": cliente.get("LOGOTIPO", "") or "",  # Logotipo da CLIENTI
            "MATRICOLA": matricola_input,  # Matricola
            "TIPOINT": "3",  # Tipo Intervento (3 = Verificazione)
        }
        
        # Campi MISURART compilabili - SENZA i campi F1+ e CARCTRL e CODICE_INTERVENTO (non esiste in MISURART)
        misurart_editable_fields = [
            "TIPO", "CF", "MATRICOLA",
            "DATAINIZIO", "DATAFINE", "DATAVP", "VERSIONE_SW", "DATA_RELEASE",
            "ESITO", "ALFASIGILL", "NUMSIGILLO",
            "PIVAESERCE", "RAGSOCESE", "COMUNESE", "PROVINCIA", "INDIRIZZO", "CAPESE",
            "MARCAMF", "MODELLO", "LOGOTIPO", "TIPOINT"
        ]
        
        for field_key in misurart_editable_fields:
            ttk.Label(misurart_lbl, text=field_key + ":").grid(row=misurart_row, column=0, sticky="w", padx=4, pady=2)
            # Usa entry data formattata per i campi data
            if field_key in ("DATAINIZIO", "DATAFINE", "DATAVP"):
                ent = _create_date_entry(misurart_lbl, width=30, max_length=8)
            else:
                ent = ttk.Entry(misurart_lbl, width=30)
            
            # Precompila con valore di default se presente
            default_value = misurart_defaults.get(field_key, "")
            if default_value:
                ent.insert(0, default_value)
            
            ent.grid(row=misurart_row, column=1, sticky="w", padx=4, pady=2)
            misurart_fields_editable[field_key] = ent
            misurart_row += 1
        
        # ===== BOTTONI DI AZIONE =====
        button_frame = ttk.Frame(scroll_frame)
        button_frame.pack(fill="x", pady=12)
        
        def on_confirm():
            """Conferma e salva la verifica"""
            try:
                # Raccogli VF e AZZ aggiornati
                vf_date = vf_entry.get().strip()
                azz_value = azz_entry.get().strip()
                
                # Se VF non è stata inserita, usa la data di intervento INTAE
                if not vf_date and intae.get("DATA_INTERVENTO"):
                    vf_date = intae.get("DATA_INTERVENTO")
                
                # Se VF è in formato DD/MM/YYYY, convertila in ISO
                if vf_date and "/" in vf_date:
                    try:
                        vf_date = database._parse_date_param_to_iso(vf_date)
                    except Exception:
                        pass
                
                # Raccogli TUTTI i dati MISURART dagli entry widget
                # Gli entry hanno già i defaults precompilati, quindi contengono sia i defaults sia le modifiche utente
                misurart_data = {}
                for field_key, entry_widget in misurart_fields_editable.items():
                    val = entry_widget.get().strip()
                    if val:
                        # Converti le date DD/MM/YYYY in ISO se necessario
                        if field_key in ("DATAINIZIO", "DATAFINE", "DATAVP") and "/" in val:
                            try:
                                val = database._parse_date_param_to_iso(val)
                            except Exception:
                                pass
                        misurart_data[field_key] = val
                
                # Forza coerenza: VF -> DATAINIZIO/DATAFINE/DATAVP e TIPOINT da combo
                if vf_date:
                    misurart_data["DATAINIZIO"] = vf_date
                    misurart_data["DATAFINE"] = vf_date
                    misurart_data["DATAVP"] = vf_date

                # Assicurati che alcuni campi essenziali siano presenti
                if "MATRICOLA" not in misurart_data:
                    misurart_data["MATRICOLA"] = matricola_input
                misurart_data["TIPOINT"] = str(tipoint_combo.get())
                
                # Chiama il DB per aggiornare e inserire
                result = database.update_verification_data(
                    cliente_id=cliente_id,
                    vf_date=vf_date,
                    azz_value=azz_value,
                    misurart_data=misurart_data
                )
                
                if not result.get("ok"):
                    messagebox.showerror("Errore", result.get("error", "Errore sconosciuto durante il salvataggio"), parent=w)
                    return
                
                # Ricarica le tabelle GUI
                try:
                    self.load_clienti(show_all=True)
                except Exception:
                    pass
                
                try:
                    self.misurart_load_all()
                except Exception:
                    pass
                
                msg = f"Verifica completata con successo!\n\n"
                if result.get("cliente_updated"):
                    msg += f"- Cliente aggiornato (VF: {vf_date or 'odierno'}, AZZ: {azz_value or 'non modificato'})\n"
                if result.get("misurart_id"):
                    msg += f"- MISURART registrato (ID: {result.get('misurart_id')})\n"
                
                messagebox.showinfo("Successo", msg, parent=w)
                w.destroy()
                # Svuota il campo matricola dopo la verifica completata
                self.verify_matricola_entry.delete(0, tk.END)
                
            except Exception as e:
                traceback.print_exc()
                messagebox.showerror("Errore", f"Errore durante la verifica: {e}", parent=w)
        
        def on_cancel():
            w.destroy()
        
        ttk.Button(button_frame, text="Annulla", command=on_cancel).pack(side="left", padx=6)
        ttk.Button(button_frame, text="Conferma Verifica", command=on_confirm).pack(side="left", padx=6)
        
        w.focus_force()
    

    # -------------------------
    # Tab Report
    # -------------------------
    def _build_tab_report(self):
        frm = ttk.Frame(self.tab_report, padding=8); frm.pack(fill="x")
        ttk.Label(frm, text="Report", font=("Arial", 14)).pack(anchor="w")
        btn_fr = ttk.Frame(frm); btn_fr.pack(fill="x", pady=6)
        ttk.Button(btn_fr, text="MEMORIE IN SCADENZA (30gg)", command=lambda: self.gui_memorie_scadenza(30)).pack(side="left", padx=6)
        ttk.Button(btn_fr, text="CHECKLIST", command=self.gui_checklist_by_vf_date).pack(side="left", padx=6)
        ttk.Button(btn_fr, text="CAMBI DGFE", command=self.gui_report_cambi_dgfe_prompt).pack(side="left", padx=6)

        # tree report generico (lo riempiamo dinamicamente in base al report)
        cols = ("Ragione sociale", "Comune, matricola", "VF", "Data installazione", "Data ultima verifica", "Scadenza assistenza", "Data fine memoria fiscale")
        self.tree_report = ttk.Treeview(self.tab_report, columns=cols, show="headings")
        for c in cols:
            self.tree_report.heading(c, text=c)
            self.tree_report.column(c, width=160)
        self.tree_report.pack(fill="both", expand=True, padx=8, pady=8)

    # --- handlers GUI (aggiungi al file gui.py, sotto i metodi esistenti) ---

    def _clear_and_set_tree_columns(self, headers):
        # rimuove tutte le righe
        self.tree_report.delete(*self.tree_report.get_children())

        # imposta le colonne
        self.tree_report["columns"] = headers
        self.tree_report["show"] = "headings"   # ← FONDAMENTALE

        for c in headers:
            self.tree_report.heading(c, text=c)
            self.tree_report.column(c, width=160, anchor="center")

        print("TREE ESISTE:", self.tree_report)
        
    def gui_memorie_scadenza(self, within_days: int = 30):
        """
        Mostra nel tree i clienti la cui FINEMF è compresa tra oggi e oggi+within_days.
        """
        print("BOTTONE PREMUTO, giorni =", within_days)
        rows = database.report_memorie_in_scadenza(within_days)
        print("RIGHE TROVATE:", len(rows))
        
        headers = ("Ragione sociale", "Fine memoria", "MF", "TEL", "CODICE")
        self._clear_and_set_tree_columns(headers)
        for r in rows:
            rag, finemf_iso, mf, tel, codice = r
            display_date = database._format_date_for_tree(finemf_iso)
            values = (rag, display_date, mf, tel, codice)
            self.tree_report.insert("", "end", values=values)

    def gui_checklist_by_vf_date(self):
        """
        Funzione unica per gestire le checklist:
        1. Chiede mese/anno VF
        2. Mostra l'elenco dei clienti con VF nel mese/anno
        3. Permette di selezionare clienti
        4. Pulsante preview per i selezionati
        5. Opzione di stampare selezionati o tutti in doppia copia
        """
        # Dialog per inserire mese/anno
        dlg = tk.Toplevel(self.root)
        dlg.title("Checklist per Mese/Anno VF")
        dlg.transient(self.root)
        dlg.grab_set()
        dlg.lift()
        dlg.focus_force()

        ttk.Label(dlg, text="Mese:").pack(side="left", padx=6, pady=6)
        months = [f"{i+1:02d} - {m}" for i, m in enumerate([
            "Gennaio","Febbraio","Marzo","Aprile","Maggio","Giugno",
            "Luglio","Agosto","Settembre","Ottobre","Novembre","Dicembre"
        ])]
        cmb_month = ttk.Combobox(dlg, values=months, state="readonly", width=18)
        cmb_month.pack(side="left", padx=6, pady=6)
        try:
            cmb_month.current(date.today().month - 1)
        except Exception:
            pass

        ttk.Label(dlg, text="Anno:").pack(side="left", padx=6, pady=6)
        ent_year = ttk.Entry(dlg, width=8)
        ent_year.pack(side="left", padx=6, pady=6)
        ent_year.insert(0, str(date.today().year))

        def do_load_clients():
            raw_month = cmb_month.get().strip()
            raw_year = ent_year.get().strip()

            month_num = None
            if raw_month:
                try:
                    month_num = int(raw_month.split()[0])
                except Exception:
                    month_num = None

            year_num = None
            if raw_year:
                try:
                    year_num = int(raw_year)
                except Exception:
                    year_num = None

            if not month_num or not year_num:
                messagebox.showerror("Errore", "Inserisci mese e anno validi")
                return

            dlg.destroy()

            # Carica i clienti
            try:
                clients = database.get_clients_for_vf_month_year(month_num, year_num) or []
            except Exception as e:
                messagebox.showerror("Errore", f"Errore recupero clienti: {e}")
                return

            vf_label = f"{month_num:02d}/{year_num}"

            if not clients:
                messagebox.showinfo("Nessun cliente", f"Nessun cliente trovato con VF nel periodo {vf_label}")
                return
            
            # Crea finestra con lista di clienti
            win = tk.Toplevel(self.root)
            win.title(f"Checklist - VF {vf_label} (Totale: {len(clients)})")
            win.geometry("900x620")
            win.transient(self.root)
            win.grab_set()
            win.lift()
            win.focus_force()

            # Treeview con checkbox (simulato tramite selezione)
            tree_cols = ("RAG_SOC", "MATRICOLA", "MARCA", "INDIRIZZO")
            tree = ttk.Treeview(win, columns=tree_cols, show="headings", selectmode="extended")
            
            for col in tree_cols:
                tree.heading(col, text=col)
                if col == "RAG_SOC":
                    tree.column(col, width=320, anchor="w", stretch=True)
                elif col == "INDIRIZZO":
                    tree.column(col, width=320, anchor="w", stretch=True)
                elif col == "MATRICOLA":
                    tree.column(col, width=140, anchor="w", stretch=False)
                else:
                    tree.column(col, width=120, anchor="w", stretch=False)
            
            # Scrollbar direttamente nella finestra
            vsb = ttk.Scrollbar(win, orient="vertical", command=tree.yview)
            hsb = ttk.Scrollbar(win, orient="horizontal", command=tree.xview)
            tree.configure(yscrollcommand=vsb.set, xscrollcommand=hsb.set)
            
            tree.grid(row=0, column=0, sticky="nsew", padx=0, pady=0)
            vsb.grid(row=0, column=1, sticky="ns")
            hsb.grid(row=1, column=0, sticky="ew")
            
            # Configura il grid per espandere
            win.grid_rowconfigure(0, weight=1)
            win.grid_columnconfigure(0, weight=1)
            
            # Popola la tree
            client_map = {}
            for c in clients:
                vals = (
                    c.get("RAG_SOC", ""),
                    c.get("MATRICOLA", ""),
                    c.get("MARCA", ""),
                    c.get("IND_VIA", "") or c.get("INDIRIZZO", "")
                )
                iid = tree.insert("", "end", values=vals)
                client_map[iid] = c
            
            # Seleziona tutti per default
            for iid in tree.get_children():
                tree.selection_add(iid)
            
            # Frame pulsanti (in grid sotto la tree)
            btn_frame = ttk.Frame(win)
            btn_frame.grid(row=2, column=0, columnspan=2, sticky="ew", padx=8, pady=8)
            
            # Checkbox per "Seleziona tutti"
            select_all_var = tk.BooleanVar(value=True)
            chk = ttk.Checkbutton(
                btn_frame,
                text="Seleziona tutti",
                variable=select_all_var,
                command=lambda: [tree.selection_set(tree.get_children()) if select_all_var.get() else tree.selection_remove(tree.get_children())]
            )
            chk.pack(side="left", padx=6)
            
            def do_preview():
                sel = tree.selection()
                if not sel:
                    messagebox.showwarning("Nessuna selezione", "Seleziona un cliente")
                    return
                if len(sel) > 1:
                    messagebox.showinfo("Seleziona uno solo", "Per la preview seleziona **un solo** cliente. Per stampare puoi selezionare più clienti.")
                    return
                iid = sel[0]
                client = client_map.get(iid)
                if not client:
                    messagebox.showerror("Errore", "Cliente non trovato nella mappa")
                    return
                # genera un singolo pdf e aprilo
                import tempfile
                tempdir = tempfile.mkdtemp(prefix="check_preview_")
                tpl = os.path.join(os.getcwd(), "template.pdf")
                if not os.path.exists(tpl):
                    messagebox.showerror("Template mancante", f"Non trovo {tpl}.")
                    return
                safe_mat = re.sub(r'[^0-9A-Za-z_-]', '_', str(client.get("MATRICOLA","")))
                pdf_out = os.path.join(tempdir, f"checklist_{safe_mat}.pdf")
                try:
                    database.fill_pdf_with_client_data(vf_label, client, tpl, pdf_out)
                    if os.name == "nt":
                        os.startfile(pdf_out)
                    else:
                        subprocess.run(["xdg-open", pdf_out], check=False)
                except Exception as e:
                    messagebox.showerror("Errore generazione PDF", str(e))

            
            def do_stampa():
                """Stampa in doppia copia"""
                # Chiedi se stampare selezionati o tutti
                res = messagebox.askyesno("Stampa", "Stampare i SELEZIONATI?\n\nSì = Selezionati\nNo = TUTTI i clienti")
                
                if res:
                    # Stampa selezionati
                    sel = tree.selection()
                    if not sel:
                        messagebox.showwarning("Nessuna selezione", "Seleziona almeno un cliente")
                        return
                    selected_clients = [client_map[iid] for iid in sel if iid in client_map]
                else:
                    # Stampa tutti
                    selected_clients = clients
                
                # Conferma stampa
                client_list = "\n".join([f"• {c.get('RAG_SOC', '')} ({c.get('MATRICOLA', '')})" for c in selected_clients])
                msg = f"Stampare in doppia copia i seguenti {len(selected_clients)} clienti?\n\n{client_list}"
                
                if not messagebox.askyesno("Conferma stampa", msg, parent=win):
                    return
                
                tpl = os.path.join(os.getcwd(), "template.pdf")
                if not os.path.exists(tpl):
                    messagebox.showerror("Template mancante", f"Non trovo {tpl}.", parent=win)
                    return
                
                # Genera e stampa
                import tempfile
                tempdir = tempfile.mkdtemp(prefix="checklist_print_")
                try:
                    for client in selected_clients:
                        matric = client.get("MATRICOLA", "") or ""
                        safe_mat = re.sub(r'[^0-9A-Za-z_-]', '_', str(matric))
                        pdf_out = os.path.join(tempdir, f"checklist_{safe_mat}.pdf")
                        
                        try:
                            # Genera PDF
                            database.fill_pdf_with_client_data(vf_label, client, tpl, pdf_out)
                            
                            # Stampa in doppia copia
                            database.print_pdf_file(pdf_out, copies=2)
                        except Exception as e:
                            messagebox.showerror("Errore", f"Errore stampando {matric}: {e}")
                            return
                    
                    messagebox.showinfo("Stampa", f"Job di stampa inviato ({len(selected_clients)} file x 2 copie)", parent=win)
                except Exception as e:
                    messagebox.showerror("Errore", f"Errore: {e}", parent=win)
                    try:
                        win.lift()
                        win.focus_force()
                    except Exception:
                        pass
            
            ttk.Button(btn_frame, text="Chiudi", command=win.destroy).pack(side="right", padx=6)
            ttk.Button(btn_frame, text="Stampa", command=do_stampa).pack(side="right", padx=6)
            ttk.Button(btn_frame, text="Anteprima selezionati", command=do_preview).pack(side="right", padx=6)
        
        ttk.Button(dlg, text="Annulla", command=dlg.destroy).pack(side="left", padx=6, pady=6)
        ttk.Button(dlg, text="CARICA", command=do_load_clients).pack(side="left", padx=6, pady=6)


    def gui_report_rt_per_model(self):
        try:
            # costruiamo il dialog
            w = tk.Toplevel(self.root)
            w.title("Filtro Report RT - mese/anno")
            w.transient(self.root)
            w.grab_set()
            w.lift()
            w.focus_force()

            ttk.Label(w, text="Mese:").grid(row=0, column=0, padx=6, pady=6, sticky="w")
            months = [f"{i+1:02d} - {m}" for i, m in enumerate([
                "Gennaio","Febbraio","Marzo","Aprile","Maggio","Giugno",
                "Luglio","Agosto","Settembre","Ottobre","Novembre","Dicembre"
            ])]
            cmb_month = ttk.Combobox(w, values=months, state="readonly", width=18)
            cmb_month.grid(row=0, column=1, padx=6, pady=6)
            try:
                cmb_month.current(date.today().month - 1)
            except Exception:
                pass

            ttk.Label(w, text="Anno:").grid(row=1, column=0, padx=6, pady=6, sticky="w")
            ent_year = ttk.Entry(w, width=10)
            ent_year.grid(row=1, column=1, padx=6, pady=6, sticky="w")
            ent_year.insert(0, str(date.today().year))

            def do_search():
                # parse input
                raw_month = cmb_month.get().strip()
                raw_year = ent_year.get().strip()

                month_num = None
                if raw_month:
                    try:
                        month_num = int(raw_month.split()[0])
                    except Exception:
                        month_num = None

                year_num = None
                if raw_year:
                    try:
                        year_num = int(raw_year)
                    except Exception:
                        year_num = None

                try:
                    w.destroy()
                except Exception:
                    pass

                # chiamiamo la funzione DB
                try:
                    if month_num is not None and year_num is not None:
                        rows = database.report_rt_for_clienti(month_num, year_num) or []
                    else:
                        rows = database.report_rt_for_clienti() or []
                except TypeError as te:
                    # firma diversa: fallo sapere e richiama senza parametri
                    print("[WARN] report_rt_for_clienti firma diversa:", te)
                    try:
                        rows = database.report_rt_for_clienti() or []
                    except Exception as e:
                        traceback.print_exc()
                        messagebox.showerror("Errore DB", f"Errore chiamando il database: {e}", parent=self.root)
                        return
                except Exception as e:
                    traceback.print_exc()
                    messagebox.showerror("Errore DB", f"Errore chiamando il database: {e}", parent=self.root)
                    return

                # AGGIUNTO: Controlla se i risultati sono None o vuoti
                print(f"[DEBUG][UI] righe restituite da DB: {0 if rows is None else len(rows)}")

                if rows is None or len(rows) == 0:
                    messagebox.showinfo("Nessun risultato", "Nessun cliente trovato per il periodo selezionato.", parent=self.root)
                    return

                # intestazioni (stesso ordine delle tuple restituite dalla funzione DB)
                headers = ("Cliente", "Indirizzo", "Comune", "Matricola", "Data VP", "Data", "Tipo", "Data Inst.", "Azz")

                # crea una nuova finestra dedicata al report
                title_suffix = f"{month_num:02d}/{year_num}" if (month_num is not None and year_num is not None) else "Tutti"
                rep_win = tk.Toplevel(self.root)
                rep_win.title(f"Report RT - {title_suffix}")
                rep_win.geometry("1100x520")
                rep_win.transient(self.root)
                rep_win.grab_set()
                rep_win.lift()
                rep_win.focus_force()

                # frame con grid per tree + scrollbars
                frame = ttk.Frame(rep_win)
                frame.pack(fill="both", expand=True, padx=8, pady=8)
                frame.rowconfigure(0, weight=1)
                frame.columnconfigure(0, weight=1)

                # tree + scrollbars
                tree = ttk.Treeview(frame, columns=headers, show="headings")
                vsb = ttk.Scrollbar(frame, orient="vertical", command=tree.yview)
                hsb = ttk.Scrollbar(frame, orient="horizontal", command=tree.xview)
                tree.configure(yscrollcommand=vsb.set, xscrollcommand=hsb.set)

                tree.grid(row=0, column=0, sticky="nsew")
                vsb.grid(row=0, column=1, sticky="ns")
                hsb.grid(row=1, column=0, sticky="ew")

                # intestazioni e larghezze (adatta se vuoi altri valori)
                for h in headers:
                    tree.heading(h, text=h)
                    if h == "Cliente":
                        wcol = 250
                    elif h == "Indirizzo":
                        wcol = 200
                    elif h == "Comune":
                        wcol = 100
                    elif h == "Matricola":
                        wcol = 110
                    elif h in ("Data VP", "Data", "Data Inst."):
                        wcol = 100
                    elif h == "Tipo":
                        wcol = 70
                    else:
                        wcol = 60
                    try:
                        tree.column(h, width=wcol, anchor="w", stretch=False)
                    except Exception:
                        pass

                # inserisci i record restituiti (normalizza lunghezza)
                inserted = 0
                first_iid = None
                for r in rows:
                    try:
                        vals = tuple('' if v is None else str(v) for v in r)
                        if len(vals) < len(headers):
                            vals = vals + tuple("" for _ in range(len(headers) - len(vals)))
                        elif len(vals) > len(headers):
                            vals = vals[:len(headers)]
                        iid = tree.insert("", "end", values=vals)
                        if first_iid is None:
                            first_iid = iid
                        inserted += 1
                    except Exception as _e:
                        print("[ERROR][UI] inserimento riga nel report separato:", _e)

                # etichetta conteggio e pulsanti utili
                status_frame = ttk.Frame(rep_win)
                status_frame.pack(fill="x", padx=8, pady=(0,8))
                _count_lbl = ttk.Label(status_frame, text=f"Record mostrati: {inserted}")
                _count_lbl.pack(side="left")

                # Funzione helper per generare il PDF
                def generate_report_pdf():
                    """Genera il PDF del report RT e ritorna il path"""
                    import tempfile
                    
                    # Genera PDF temporaneo con i dati del report
                    temp_dir = tempfile.mkdtemp(prefix="report_rt_")
                    pdf_path = os.path.join(temp_dir, "report_rt.pdf")
                    
                    # Crea PDF con reportlab
                    from reportlab.lib.pagesizes import A4, landscape
                    from reportlab.lib.units import mm
                    from reportlab.pdfgen import canvas as rl_canvas
                    
                    c = rl_canvas.Canvas(pdf_path, pagesize=landscape(A4))
                    width, height = landscape(A4)
                    
                    # Titolo
                    c.setFont("Helvetica-Bold", 14)
                    title_txt = f"Report RT - {title_suffix}"
                    c.drawString(20*mm, height - 15*mm, title_txt)
                    
                    # Header tabella
                    y = height - 25*mm
                    c.setFont("Helvetica-Bold", 8)
                    # Posizioni X per 9 colonne
                    x_positions = [10*mm, 60*mm, 110*mm, 150*mm, 175*mm, 205*mm, 230*mm, 255*mm, 275*mm]
                    
                    for i, h in enumerate(headers):
                        if i < len(x_positions):
                            c.drawString(x_positions[i], y, h)
                    
                    # Linea sotto header
                    y -= 3*mm
                    c.line(10*mm, y, width - 10*mm, y)
                    y -= 5*mm
                    
                    # Dati - con testo troncato
                    c.setFont("Helvetica", 7)
                    max_chars = {0: 30, 1: 25, 2: 15, 3: 12, 4: 12, 5: 12, 6: 10, 7: 12, 8: 8}  # Max char per colonna
                    
                    for row in rows:
                        if y < 20*mm:  # Nuova pagina
                            c.showPage()
                            y = height - 20*mm
                            c.setFont("Helvetica", 7)
                        
                        vals = tuple('' if v is None else str(v) for v in row)
                        for i, val in enumerate(vals[:len(headers)]):
                            # Tronca il testo se più lungo del massimo per la colonna
                            max_len = max_chars.get(i, 20)
                            truncated = val[:max_len] if len(val) > max_len else val
                            if len(x_positions) > i:
                                c.drawString(x_positions[i], y, truncated)
                        y -= 5*mm
                    
                    c.save()
                    return pdf_path
                
                # bottone per anteprima report
                def do_preview_report():
                    """Mostra anteprima del report RT (senza stampare)"""
                    try:
                        pdf_path = generate_report_pdf()
                        
                        # Apri il PDF nel visualizzatore predefinito
                        if sys.platform.startswith("win"):
                            os.startfile(os.path.abspath(pdf_path))
                        elif sys.platform == "darwin":  # macOS
                            subprocess.Popen(["open", pdf_path])
                        else:  # Linux
                            subprocess.Popen(["xdg-open", pdf_path])
                    except Exception as e:
                        import traceback
                        traceback.print_exc()
                        messagebox.showerror("Errore", f"Errore durante l'apertura dell'anteprima: {e}", parent=rep_win)
                
                # bottone per stampare il report
                def do_print_report():
                    """Stampa il report RT (singola copia) con testo troncato per le colonne lunghe"""
                    try:
                        pdf_path = generate_report_pdf()
                        
                        # Stampa il PDF
                        database.print_pdf_file(pdf_path, copies=1)
                        messagebox.showinfo("Stampa", "Report inviato alla stampante", parent=rep_win)
                        
                    except Exception as e:
                        import traceback
                        traceback.print_exc()
                        messagebox.showerror("Errore", f"Errore durante la stampa: {e}", parent=rep_win)
                
                _preview_btn = ttk.Button(status_frame, text="Anteprima", command=do_preview_report)
                _preview_btn.pack(side="right", padx=6)
                
                _print_btn = ttk.Button(status_frame, text="Stampa", command=do_print_report)
                _print_btn.pack(side="right", padx=6)

                # bottone per chiudere la finestra
                _close_btn = ttk.Button(status_frame, text="Chiudi", command=rep_win.destroy)
                _close_btn.pack(side="right")

                # seleziona e porta in vista il primo elemento
                try:
                    if first_iid:
                        tree.selection_set(first_iid)
                        tree.focus(first_iid)
                        tree.see(first_iid)
                except Exception:
                    pass

                # metti a fuoco la nuova finestra (non bloccante)
                try:
                    rep_win.transient(self.root)
                    rep_win.focus_force()
                except Exception:
                    pass

                print(f"[DEBUG][UI] inserite righe nella finestra report: {inserted}")
                # ----------------- FINE: report in finestra separata -----------------

            # bottoni
            btn_fr = ttk.Frame(w)
            btn_fr.grid(row=2, column=0, columnspan=2, pady=8)
            ttk.Button(btn_fr, text="Annulla", command=w.destroy).pack(side="left", padx=6)
            ttk.Button(btn_fr, text="CERCA", command=do_search).pack(side="left", padx=6)

            ent_year.bind("<Return>", lambda e: do_search())
            cmb_month.bind("<Return>", lambda e: do_search())
            ent_year.focus_set()

        except Exception as e:
            traceback.print_exc()
            try:
                messagebox.showerror("Errore Report RT", f"Errore durante la generazione del report:\n{str(e)}", parent=self.root)
            except Exception as e2:
                print(f"[ERROR] Impossibile mostrare messagebox: {e2}")
                print(f"[ERROR] Errore originale: {e}")

    def gui_report_cambi_dgfe_prompt(self):
        def do_search():
            d1_raw = ent1.get().strip()
            d2_raw = ent2.get().strip()
            
            # Converti numeri DDMMYY o DDMMYYYY in formato DD/MM/YYYY
            try:
                if d1_raw:
                    if "/" in d1_raw:
                        parts = d1_raw.split("/")
                        if len(parts[2]) == 2:  # YY -> aggiungi 20
                            d1_raw = f"{parts[0]}/{parts[1]}/20{parts[2]}"
                        # Altrimenti già YYYY
                    elif len(d1_raw) == 6:  # DDMMYY -> DD/MM/20YY
                        d1_raw = f"{d1_raw[0:2]}/{d1_raw[2:4]}/20{d1_raw[4:6]}"
                    elif len(d1_raw) == 8:  # DDMMYYYY -> DD/MM/YYYY
                        d1_raw = f"{d1_raw[0:2]}/{d1_raw[2:4]}/{d1_raw[4:8]}"
                
                if d2_raw:
                    if "/" in d2_raw:
                        parts = d2_raw.split("/")
                        if len(parts[2]) == 2:  # YY -> aggiungi 20
                            d2_raw = f"{parts[0]}/{parts[1]}/20{parts[2]}"
                        # Altrimenti già YYYY
                    elif len(d2_raw) == 6:  # DDMMYY -> DD/MM/20YY
                        d2_raw = f"{d2_raw[0:2]}/{d2_raw[2:4]}/20{d2_raw[4:6]}"
                    elif len(d2_raw) == 8:  # DDMMYYYY -> DD/MM/YYYY
                        d2_raw = f"{d2_raw[0:2]}/{d2_raw[2:4]}/{d2_raw[4:8]}"
            except Exception:
                pass
            
            w.destroy()
            rows = _safe_call(database.report_cambi_dgfe, d1_raw, d2_raw) or []
            headers = ("Cliente", "Matricola", "Interventi")
            self._clear_and_set_tree_columns(headers)
            if not rows:
                messagebox.showinfo("Nessun risultato", 
                                   "Nessun intervento trovato per le matricole nell'intervallo specificato.")
                return
            
            # Le righe arrivano già con (ragione_sociale, matricola, count)
            for r in rows:
                rag = r[0] if len(r) > 0 else "N/D"
                mat = r[1] if len(r) > 1 else ""
                cnt = str(r[2]) if len(r) > 2 else "0"
                self.tree_report.insert("", "end", values=(rag, mat, cnt))
        
        w = tk.Toplevel(self.root)
        w.title("Cambi DGFE")
        w.transient(self.root)
        w.grab_set()
        w.lift()
        w.focus_force()
        ttk.Label(w, text="Data from (DDMMYY o DDMMYYYY):").pack(side="left", padx=4)
        ent1 = _create_date_entry(w, width=12, max_length=8)
        ent1.pack(side="left", padx=4)
        
        ttk.Label(w, text="Data to (DDMMYY o DDMMYYYY):").pack(side="left", padx=4)
        ent2 = _create_date_entry(w, width=12, max_length=8)
        ent2.pack(side="left", padx=4)
        
        ttk.Button(w, text="CERCA", command=do_search).pack(side="left", padx=6)

    def enable_tree_autosave(self, debounce_ms: int = 800):
        """Attiva editing inline + autosave sulla tree."""
        self._tree_edit_after_id = None
        self._tree_edit_debounce_ms = debounce_ms
        self.tree_report.bind("<Double-1>", lambda e: self._on_tree_double_click(e))

    def _on_tree_double_click(self, event):
        tv = self.tree_report
        region = tv.identify("region", event.x, event.y)
        if region != "cell":
            return
        rowid = tv.identify_row(event.y)
        col = tv.identify_column(event.x)
        if not rowid or not col:
            return

        x, y, width, height = tv.bbox(rowid, col)
        if width <= 0 or height <= 0:
            return

        col_idx = int(col.replace("#", "")) - 1
        values = tv.item(rowid, "values")
        try:
            cur_value = values[col_idx]
        except Exception:
            cur_value = ""

        entry = tk.Entry(tv)
        entry.insert(0, cur_value)
        entry.place(x=x, y=y, width=width, height=height)
        entry.focus_set()
        entry._tree_rowid = rowid
        entry._tree_col_idx = col_idx
        entry._tree_col_name = None
        
        if hasattr(self, "tree_headers") and len(self.tree_headers) > col_idx:
            entry._tree_col_name = self.tree_headers[col_idx]
        else:
            entry._tree_col_name = f"col_{col_idx}"

        def normalize_date_input(s):
            s = s.strip()
            if not s:
                return ""
            try:
                if "/" in s:
                    parts = s.split("/")
                    if len(parts) == 3:
                        if len(parts[2]) == 2:
                            return f"{parts[0]}/{parts[1]}/20{parts[2]}"
                        return f"{parts[0]}/{parts[1]}/{parts[2]}"
                elif len(s) == 6:
                    return f"{s[0:2]}/{s[2:4]}/20{s[4:6]}"
                elif len(s) == 8:
                    return f"{s[0:2]}/{s[2:4]}/{s[4:8]}"
            except Exception:
                pass
            return s

        def do_search():
            d1_raw = ent1.get().strip()
            d2_raw = ent2.get().strip()
            d1_norm = normalize_date_input(d1_raw)
            d2_norm = normalize_date_input(d2_raw)
            if not d1_norm or not d2_norm:
                messagebox.showerror("Errore", "Inserire entrambe le date in formato valido.")
                return
            w.destroy()
            try:
                rows = _safe_call(database.report_cambi_dgfe, d1_norm, d2_norm) or []
            except Exception as e:
                messagebox.showerror("Errore", f"Errore nella ricerca DGFE: {e}")
                return
            headers = ("Cliente", "Matricola", "Interventi")
            self._clear_and_set_tree_columns(headers)
            if not rows:
                messagebox.showinfo("Nessun risultato", 
                                   "Nessun intervento trovato per le matricole nell'intervallo specificato.")
                return
            # se manca ragione sociale, prova lookup su clienti
            rag_lookup = {}
            try:
                conn = database.get_connection()
                cur = conn.cursor()
                cur.execute("SELECT MATRICOLA, RAG_SOC FROM clienti")
                for m, rag in cur.fetchall():
                    if m:
                        key = _normalize_matricola_for_compare(m)
                        if key:
                            rag_lookup[key] = rag or ""
                conn.close()
            except Exception:
                rag_lookup = {}

            for r in rows:
                rag = r[0] if len(r) > 0 else ""
                mat = r[1] if len(r) > 1 else ""
                if not rag and mat:
                    rag = rag_lookup.get(_normalize_matricola_for_compare(mat), "")
                self.tree_report.insert("", "end", values=(rag, mat, str(r[2]) if len(r) > 2 else ""))
        w = tk.Toplevel(self.root)
        w.title("Cambi DGFE")
        w.transient(self.root)
        w.grab_set()
        w.lift()
        w.focus_force()
        ttk.Label(w, text="Data from (DDMMYY o DDMMYYYY):").pack(side="left", padx=4)
        ent1 = _create_date_entry(w, width=12, max_length=8)
        ent1.pack(side="left", padx=4)
        ttk.Label(w, text="Data to (DDMMYY o DDMMYYYY):").pack(side="left", padx=4)
        ent2 = _create_date_entry(w, width=12, max_length=8)
        ent2.pack(side="left", padx=4)
        ttk.Button(w, text="CERCA", command=do_search).pack(side="left", padx=6)
        # except Exception:
        #     print("[_on_tree_entry_commit] errore generico in commit entry.")
        #     return

        if getattr(self, "_tree_edit_after_id", None):
            try:
                self.after_cancel(self._tree_edit_after_id)
            except Exception:
                pass

        def do_save(db_id):
            try:
                data = {}
                if db_id is None:
                    print("[Autosave] db_id non definito, salvataggio saltato.")
                    return
                res = database.save_all_cliente_fields(db_id, data, run_post_update=True)
                if not res.get("ok", False):
                    print("[Autosave] errore salvataggio:", res.get("error"))
            except Exception as e:
                print("[Autosave] eccezione:", e)

        # Calcola db_id dal rowid selezionato
        db_id = None
        try:
            tags = self.tree_report.item(rowid, "tags")
            if tags:
                for t in tags:
                    try:
                        db_id = int(t)
                        break
                    except Exception:
                        if isinstance(t, str) and t.startswith("id:"):
                            try:
                                db_id = int(t.split(":", 1)[1])
                                break
                            except Exception:
                                pass
        except Exception:
            pass
        if db_id is not None:
            self._tree_edit_after_id = self.after(self._tree_edit_debounce_ms, lambda: do_save(db_id))


# Funzione di avvio
def avvia_gui():
    root = tk.Tk()
    app = App(root)
    root.mainloop()


if __name__ == "__main__":
    avvia_gui()