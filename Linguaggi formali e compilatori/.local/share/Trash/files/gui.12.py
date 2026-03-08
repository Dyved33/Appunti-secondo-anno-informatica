import tkinter as tk
from tkinter import ttk, filedialog, messagebox
import traceback
import database
import os
from database import get_connection
from database import CLIENTI_COLUMNS

def _safe_call(fn, *args, **kwargs):
    try:
        return fn(*args, **kwargs)
    except Exception as e:
        traceback.print_exc()
        messagebox.showerror("Errore", str(e))
        return None

class App:
    def __init__(self, root):
        self.root = root
        self.root.title("Gestione Clienti e Interventi")
        self.root.geometry("1200x720")

        # lista in memoria dei clienti
        self.clienti = []

        # inizializza DB senza ricreare la tabella clienti (persistenza dati)
        database.ensure_db_exists(recreate_clienti=False)

        # costruzione interfaccia
        self._build_ui()

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
        ttk.Button(btn_fr, text="IMPORTA CLIENTI (CSV/XLSX)", command=self.gui_import_any_clienti).pack(side="left", padx=4)
        ttk.Button(btn_fr, text="NUOVO", command=self.gui_new_cliente).pack(side="left", padx=4)
        ttk.Button(btn_fr, text="DUPLICA", command=self.gui_duplicate_cliente).pack(side="left", padx=4)
        # prima: ttk.Button(btn_fr, text="ELIMINA", command=self.gui_delete_cliente).pack(side="left", padx=4)
        tk.Button(btn_fr, text="ELIMINA", command=self.gui_delete_cliente, bg="#d9534f", fg="white", activebackground="#c9302c").pack(side="left", padx=4)
        ttk.Button(btn_fr, text="REPORT RT", command=self.gui_report_rt_per_model).pack(side="left", padx=4)
        ttk.Button(btn_fr, text="DISATTIVA MACCHINA", command=self.gui_deactivate_machine).pack(side="left", padx=4)
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
            tree.heading(c, text=c)
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

    def open_edit_clients_window(self, cols, vals):
        w = tk.Toplevel(self.root)
        w.title("Modifica cliente")
        frm = ttk.Frame(w, padding=8); frm.pack(fill="both", expand=True)
        canvas = tk.Canvas(frm)
        scrollbar = ttk.Scrollbar(frm, orient="vertical", command=canvas.yview)
        scroll_frame = ttk.Frame(canvas)
        scroll_frame.bind("<Configure>", lambda e: canvas.configure(scrollregion=canvas.bbox("all")))
        canvas.create_window((0,0), window=scroll_frame, anchor="nw")
        canvas.configure(yscrollcommand=scrollbar.set, height=500)
        canvas.pack(side="left", fill="both", expand=True)
        scrollbar.pack(side="right", fill="y")
        entries = {}
        rowi = 0
        for i, c in enumerate(cols):
            if c.lower() == "id":
                continue
            ttk.Label(scroll_frame, text=c).grid(row=rowi, column=0, sticky="w", padx=4, pady=2)
            ent = ttk.Entry(scroll_frame, width=70)
            ent.grid(row=rowi, column=1, sticky="w", padx=4, pady=2)
            val = vals[i] if i < len(vals) else ""
            ent.insert(0, "" if val is None else str(val))
            entries[c] = ent
            
            # se è una colonna data, attacca il formatter (mostra DD/MM/YYYY mentre si digita)
            try:
                if c.upper() in getattr(database, "DATE_COLS", set()):
                    def make_on_date_input(entry_widget):
                        def on_date_input_event(event=None):
                            v = entry_widget.get()
                            digits = ''.join(ch for ch in v if ch.isdigit())[:8]
                            # format progressivo: DD / DDMM / DDMMYYYY
                            if len(digits) >= 5:
                                formatted = digits[:2] + '/' + digits[2:4] + '/' + digits[4:]
                            elif len(digits) >= 3:
                                formatted = digits[:2] + '/' + digits[2:]
                            else:
                                formatted = digits
                            # evita loop infinito: aggiorna solo se diverso
                            if formatted != v:
                                # salva cursor pos relativo (semplice heuristica)
                                try:
                                    pos = entry_widget.index(tk.INSERT)
                                except Exception:
                                    pos = None
                                entry_widget.delete(0, tk.END)
                                entry_widget.insert(0, formatted)
                                if pos is not None:
                                    try:
                                        # tenta riposizionare il cursore in fondo o vicino alla posizione originale
                                        entry_widget.icursor(min(len(formatted), pos + (1 if '/' in formatted else 0)))
                                    except Exception:
                                        pass
                        return on_date_input_event
                    on_date_input = make_on_date_input(ent)
                    ent.bind("<KeyRelease>", on_date_input)
            except Exception:
                pass
            
            rowi += 1
        btn_fr = ttk.Frame(scroll_frame); btn_fr.grid(row=rowi+1, column=0, columnspan=2, pady=8)
        
        def do_save():
            # raccogli input
            updates = {c: entries[c].get().strip() for c in entries}

            # estrai matricola per controllo duplicato
            mat = updates.get("MATRICOLA", "").strip()

            try:
                # Se è INSERT (nuovo cliente)
                if not self.selected_cliente_id:
                    if mat:
                        # controllo duplicato: se esiste => popup e blocco
                        if database.matricola_exists(mat):
                            messagebox.showerror("Errore", f"Impossibile inserire: la matricola '{mat}' è già presente.")
                            return

                else:
                    # Se è UPDATE: se la nuova matricola è presente in un'altra riga -> errore
                    if mat:
                        if database.matricola_exists(mat, exclude_id=self.selected_cliente_id):
                            messagebox.showerror("Errore", f"Impossibile aggiornare: la matricola '{mat}' è già presente su un altro cliente.")
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
                                messagebox.showerror("Errore data", f"Formato data non valido per '{k}': {v}\nDettaglio: {e}")
                                return

                # chiamata DB (wrap _safe_call per gestire eccezioni)
                if not self.selected_cliente_id:
                    nid = _safe_call(database.insert_clients_from_map, updates)
                    if nid:
                        messagebox.showinfo("OK","Cliente inserito")
                    else:
                        # insert fallito: _safe_call mostrerà errore se eccezione
                        return
                else:
                    affected = _safe_call(database.update_clients_by_id, self.selected_cliente_id, updates)
                    if affected is None:
                        return
                    messagebox.showinfo("OK","Cliente aggiornato")

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
                messagebox.showerror("Errore salvataggio", str(e))
                return
 
        def do_cancel():
            w.destroy()
        ttk.Button(btn_fr, text="SALVA", command=do_save).pack(side="left", padx=6)
        ttk.Button(btn_fr, text="ANNULLA", command=do_cancel).pack(side="left", padx=6)

    def gui_new_cliente(self):
        cols = database.get_clienti_columns()
        self.selected_cliente_id = None
        self.open_edit_clients_window(cols, [""] * len(cols))

    def gui_update_cliente(self):
        sel = self.tree.selection()
        if not sel:
            messagebox.showerror("Errore","Seleziona un cliente")
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
            messagebox.showerror("Errore","Seleziona un cliente")
            return
        if not messagebox.askyesno("Conferma","Eliminare cliente e macchine collegate?"):
            return
        _safe_call(database.delete_clients, self.selected_cliente_id)
        messagebox.showinfo("OK","Cliente eliminato")
        self.load_clienti(show_all=True)
        self.selected_cliente_id = None

    def gui_duplicate_cliente(self):
        if not self.selected_cliente_id:
            messagebox.showerror("Errore","Seleziona un cliente")
            return
        nid = _safe_call(database.duplicate_clients, self.selected_cliente_id)
        if nid:
            messagebox.showinfo("OK","Cliente duplicato")
            self.load_clienti(show_all=True)

    def gui_import_csv_clienti(self, path: str | None = None):
        """
        Se path è None -> apre filedialog. Altrimenti usa il path passato.
        Restituisce il numero di record importati oppure None in caso di annullamento / errore.
        """
        if path is None:
            path = filedialog.askopenfilename(title="Seleziona CSV", filetypes=[("CSV files","*.csv"),("All","*.*")])
            if not path:
                return None
        # try ; then ,
        n = _safe_call(database.import_clients_from_csv_generic, path, ";", "latin-1")
        if n is None:
            n = _safe_call(database.import_clients_from_csv_generic, path, ",", "latin-1")
        # Nota: non ricarichiamo la UI qui. Caller deciderà quando aggiornare.
        return n

    def gui_import_xls_clienti(self, path: str | None = None):
        """
        Se path è None -> apre filedialog. Altrimenti usa il path passato.
        Restituisce il numero di record importati oppure None.
        """
        if path is None:
            path = filedialog.askopenfilename(title="Seleziona XLS/XLSX", filetypes=[("Excel files","*.xls;*.xlsx"),("All","*.*")])
            if not path:
                return None
        n = _safe_call(database.import_excel_clienti_via_soffice, path)
        return n

    def gui_import_any_clienti(self):
        path = filedialog.askopenfilename(
            title="Seleziona file da importare",
            filetypes=[("CSV","*.csv"), ("Excel","*.xls *.xlsx"), ("Tutti i file","*.*")]
        )
        if not path:
            return

        ext = os.path.splitext(path)[1].lower()
        try:
            n = None
            if ext == ".csv":
                n = self.gui_import_csv_clienti(path)
            elif ext in (".xls", ".xlsx"):
                n = self.gui_import_xls_clienti(path)
            else:
                if messagebox.askyesno("Importazione", "Estensione non riconosciuta. Trattare come CSV?"):
                    n = self.gui_import_csv_clienti(path)
                elif messagebox.askyesno("Importazione", "Provare come file Excel (soffice)?"):
                    n = self.gui_import_xls_clienti(path)
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

            # NOTIFICA: qui mostriamo il popup sempre
            try:
                if n is None:
                    messagebox.showinfo("Importazione", "Import annullato o non è stato inserito alcun record.")
                else:
                    # se la funzione DB ritorna tuple (ad esempio (inserted, skipped)), gestiscila
                    if isinstance(n, tuple):
                        # esempio per (inserted, skipped)
                        if len(n) == 2:
                            inserted, skipped = n
                            messagebox.showinfo("Importazione", f"Import completato: {inserted} record inseriti. {skipped} righe saltate.")
                        elif len(n) == 3:
                            inserted, skipped, updated = n
                            messagebox.showinfo("Importazione", f"Import completato: {inserted} inseriti, {updated} aggiornati, {skipped} righe saltate.")
                        else:
                            messagebox.showinfo("Importazione", f"Import completato: {n}")
                    else:
                        messagebox.showinfo("Importazione", f"Import completato: {n} record importati.")
            except Exception as e:
                # se proprio messagebox fallisce, almeno logghiamo
                print("Errore mostrando popup:", e)

        except Exception as e:
            # mostra un errore esplicito (messagebox importato globalmente)
            messagebox.showerror("Errore importazione", f"Import fallito: {e}")
                
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
    
    def carica_clienti_da_db(self):
        try:
            conn = get_connection()
            cur = conn.cursor()
            cur.execute("SELECT * FROM clienti ORDER BY id DESC")
            self.clienti = cur.fetchall()
        finally:
            conn.close()

    def refresh_tabella(self):
        if not hasattr(self, "tree") or self.tree is None:
            return
        for row in self.tree.get_children():
            self.tree.delete(row)
        for c in self.clienti:
            self.tree.insert("", "end", values=c) 
            
    
    def gui_deactivate_machine(self):
        """
        Apri un dialog semplice per inserire matricola e data, chiedi conferma,
        quindi chiama deactivate_by_matricola(matricola, date).
        Incolla questo metodo nella tua classe GUI (es. dove hai gli altri handler).
        """

        # crea dialog
        top = tk.Toplevel(self.root if hasattr(self, "root") else self)
        top.title("Disattiva Macchina")
        top.grab_set()  # focus sul dialog

        ttk.Label(top, text="Matricola:").grid(row=0, column=0, padx=6, pady=6, sticky="w")
        ent_mat = ttk.Entry(top, width=30)
        ent_mat.grid(row=0, column=1, padx=6, pady=6)

        ttk.Label(top, text="Data (DDMMYYYY o DD/MM/YYYY) [vuoto = oggi]:").grid(row=1, column=0, padx=6, pady=6, sticky="w")
        ent_date = ttk.Entry(top, width=20)
        ent_date.grid(row=1, column=1, padx=6, pady=6)

        # helper per formattare la data in visuale durante la digitazione (opzionale)
        def on_date_input(event=None):
            v = ent_date.get()
            digits = ''.join(ch for ch in v if ch.isdigit())[:8]
            if len(digits) >= 5:
                ent_date.delete(0, tk.END)
                ent_date.insert(0, digits[:2] + '/' + digits[2:4] + '/' + digits[4:])
            elif len(digits) >= 3:
                ent_date.delete(0, tk.END)
                ent_date.insert(0, digits[:2] + '/' + digits[2:])
            else:
                ent_date.delete(0, tk.END)
                ent_date.insert(0, digits)

        ent_date.bind("<KeyRelease>", on_date_input)

        # funzione chiamata su conferma
        def do_deactivate():
            matricola = ent_mat.get().strip()
            date_param = ent_date.get().strip() or None
            if not matricola:
                messagebox.showwarning("Input mancante", "Inserisci una matricola.")
                return

            # conferma
            if not messagebox.askyesno("Conferma", f"Confermi disattivazione matricola '{matricola}' (data={date_param or 'oggi'})?"):
                return

            # chiama la funzione di backend (usa la wrapper che ti ho dato)
            try:
                # importa qui per evitare import circolare; presuppone che deactivate_by_matricola sia in database.py
                result = database.deactivate_machine(matricola, date_param)
                
                # ricarica l'elenco clienti per vedere subito le modifiche
                try:
                    self.load_clienti(show_all=True)
                except Exception:
                    try:
                        self.carica_clienti_da_db()
                        self.refresh_tabella()
                    except Exception:
                        pass

            except Exception as e:
                messagebox.showerror("Errore", f"Errore durante la disattivazione:\n{e}")
                return

            # DOPO result = deactivate_by_matricola(...)
            from database import get_connection, _normalize_matricola_for_compare

            # mostra messaggio di successo come prima...
            # ora recupera i valori aggiornati da clienti e aggiorna la UI
            try:
                conn = get_connection()
                cur = conn.cursor()
                # proviamo find-by-lower(matricola)
                cur.execute("SELECT DATADIS, VF, matricola FROM clienti WHERE lower(matricola) = ? LIMIT 1", (matricola.lower(),))
                row = cur.fetchone()
                if row:
                    datadist_val, vf_val, matricola_db = row[0], row[1], row[2]
                else:
                    # fallback: cerca normalizzando
                    cur.execute("SELECT rowid, matricola, DATADIS, VF FROM clienti")
                    datadist_val = vf_val = matricola_db = None
                    for rid, m_db, dd, vf in cur.fetchall():
                        if _normalize_matricola_for_compare(m_db) == _normalize_matricola_for_compare(matricola):
                            datadist_val, vf_val, matricola_db = dd, vf, m_db
                            break
            finally:
                try: conn.close()
                except: pass

            # ora aggiorna i widget: esempio se hai Label o Entry:
            if datadist_val is not None:
                # se hai un Entry chiamato self.entry_datadist
                try:
                    self.entry_datadist.delete(0, "end")
                    self.entry_datadist.insert(0, datadist_val)
                except Exception:
                    pass

            if vf_val is not None:
                try:
                    self.entry_vf.delete(0, "end")
                    self.entry_vf.insert(0, str(vf_val))
                except Exception:
                    pass

            # se usi un Treeview per la lista clienti, aggiorna solo la riga:
            # supponiamo self.tree_clients sia il Treeview e la colonna 0 contiene matricola
            if hasattr(self, "tree_clients"):
                for item in self.tree_clients.get_children():
                    vals = self.tree_clients.item(item, "values")
                    if vals:
                        # supponi che la matricola sia in posizione 0 della tupla values
                        if _normalize_matricola_for_compare(vals[0]) == _normalize_matricola_for_compare(matricola_db or matricola):
                            new_vals = list(vals)
                            # supponi che DATADIST sia colonna 5 e VF colonna 6 (adatta agli indici della tua tabella)
                            # se non sai gli indici, cerca nella tua configurazione Treeview
                            try:
                                new_vals[5] = datadist_val
                                new_vals[6] = vf_val
                            except Exception:
                                # se gli indici non combaciano, setta l'intera riga come preferisci
                                pass
                            self.tree_clients.item(item, values=new_vals)
                            break

            # mostra risultato
            if not result:
                messagebox.showinfo("Completato", "Operazione completata: nessuna riga aggiornata.")
            else:
                lines = "\n".join(f"{t}: {c}" for t, c in result.items())
                messagebox.showinfo("Completato", f"Disattivazione completata. Righe aggiornate:\n{lines}")

            top.destroy()

        # bottoni
        btn_frame = ttk.Frame(top)
        btn_frame.grid(row=2, column=0, columnspan=2, pady=(6,10))
        ttk.Button(btn_frame, text="Annulla", command=top.destroy).pack(side="left", padx=6)
        ttk.Button(btn_frame, text="Conferma disattivazione", command=do_deactivate).pack(side="left", padx=6)

        # focus sul campo matricola
        ent_mat.focus_set()

        # opzionale: gestire Enter per conferma rapida
        def on_enter(event):
            do_deactivate()
        ent_mat.bind("<Return>", on_enter)
        ent_date.bind("<Return>", on_enter)

            
# -------------------------
# Tab INTAE (import & search)
# -------------------------
    def _build_tab_intae(self):
        frm = ttk.Frame(self.tab_intae, padding=8)
        frm.pack(fill="x")
        ttk.Label(frm, text="Import & ricerca INTAE", font=("Arial",12)).pack(anchor="w")

        btn_fr = ttk.Frame(frm)
        btn_fr.pack(fill="x", pady=6)
        ttk.Button(btn_fr, text="IMPORTA CSV INTAE", command=self.gui_import_csv_intae).pack(side="left", padx=6)
        ttk.Button(btn_fr, text="CERCA PER MATRICOLA", command=self.search_intae_prompt).pack(side="left", padx=6)
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

        # intestazioni
        for c in cols:
            self.tree_intae.heading(c, text=c.replace("_", " ").title())
            self.tree_intae.column(c, width=150, minwidth=80)

        self.refresh_intae()

    def refresh_intae(self):
        for row in self.tree_intae.get_children():
            self.tree_intae.delete(row)

        rows = _safe_call(database.read_intae) or []

        for r in rows:
            self.tree_intae.insert("", "end", iid=r[0], values=r[1:])

    def search_intae_prompt(self):
        def do_search():
            mat = ent.get().strip(); w.destroy()
            if not mat: return
            rows = _safe_call(database.cerca_interventi_intae, mat) or []
            for r in self.tree_intae.get_children(): 
                self.tree_intae.delete(r)
            for row in rows:
                # row = (id, col1, col2, ...). UI columns are only the col1.. => skip id
                try:
                    iid = row[0]
                    vals = row[1:]
                except Exception:
                    iid = None
                    vals = row
                # usa iid se presente (comodo per riferimenti futuri)
                if iid is not None:
                    self.tree_intae.insert("", "end", iid=iid, values=vals)
                else:
                    self.tree_intae.insert("", "end", values=vals)
        w = tk.Toplevel(self.root)
        ttk.Label(w, text="Matricola:").pack(side="left", padx=6, pady=6)
        ent = ttk.Entry(w); ent.pack(side="left", padx=6, pady=6)
        ttk.Button(w, text="CERCA", command=do_search).pack(side="left", padx=6, pady=6)

    def gui_import_csv_intae(self):
        path = filedialog.askopenfilename(title="CSV INTAE", filetypes=[("CSV files","*.csv"),("All","*.*")])
        if not path:
            return
        n = _safe_call(database.import_csv_intae, path, ";", "utf-8")
        if n is not None:
            messagebox.showinfo("Import", f"Importati {n} record in INTAE")
            self.refresh_intae()

# -------------------------
# Tab MISURART
# -------------------------
    def _build_tab_misurart(self):
        frm = ttk.Frame(self.tab_misurart, padding=8)
        frm.pack(fill="x")
        ttk.Label(frm, text="Registrazioni MISURART", font=("Arial",12)).pack(anchor="w")

        # Form inserimento
        form = ttk.LabelFrame(self.tab_misurart, text="Nuova / Modifica registrazione", padding=8)
        form.pack(fill="x", padx=10, pady=10)

        self.misurart_fields = {}
        cols = database.MISURART_COLUMNS[:]   # solo colonne reali

        for i, key in enumerate(cols):
            r = i // 2
            c = (i % 2) * 2
            ttk.Label(form, text=key).grid(row=r, column=c, sticky="w", padx=4, pady=2)
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
            self.tree_misurart.heading(c, text=c.replace("_", " ").title())
            self.tree_misurart.column(c, width=140, minwidth=80)

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
        # raccogli dati dalla mappa di entry
        data = {k: self.misurart_fields[k].get().strip() for k in self.misurart_fields}
        # validazioni minime (se vuoi aggiungerne altre fallo qui)
        # ad esempio: matricola o datainizio possono essere obbligatorie a tua scelta
        nid = _safe_call(database.misurart_add_from_map, data)
        if nid:
            messagebox.showinfo("OK", "Registrazione inserita")
            self.misurart_load_all()

    def misurart_on_select(self, event):
        sel = self.tree_misurart.selection()
        if not sel:
            self.misurart_selected_id = None
            return
        values = self.tree_misurart.item(sel[0])["values"]
        self.misurart_selected_id = values[0]
        # riempi i campi: valori sono in ordine MISURART_COLUMNS
        for i, key in enumerate(database.MISURART_COLUMNS):
            val = values[i+1] if len(values) > i+1 else ""
            self.misurart_fields[key].delete(0, tk.END)
            self.misurart_fields[key].insert(0, val or "")

    def misurart_update(self):
        if not getattr(self, "misurart_selected_id", None):
            messagebox.showerror("Errore", "Seleziona una registrazione da aggiornare.")
            return
        data = {k: self.misurart_fields[k].get().strip() for k in self.misurart_fields}
        _safe_call(database.misurart_update, self.misurart_selected_id, data)
        messagebox.showinfo("OK", "Registrazione aggiornata")
        self.misurart_load_all()

    def misurart_delete(self):
        if not getattr(self, "misurart_selected_id", None):
            messagebox.showerror("Errore", "Seleziona una registrazione da eliminare.")
            return
        if not messagebox.askyesno("Conferma", "Eliminare la registrazione?"):
            return
        _safe_call(database.misurart_delete, self.misurart_selected_id)
        messagebox.showinfo("OK", "Registrazione eliminata")
        self.misurart_load_all()
        self.misurart_selected_id = None

    def misurart_search(self):
        term = self.misurart_search_entry.get().strip()
        if not term:
            messagebox.showerror("Errore", "Inserisci una matricola da cercare")
            return
        rows = _safe_call(database.misurart_search_by_matricola, term) or []
        for r in self.tree_misurart.get_children(): self.tree_misurart.delete(r)
        for row in rows:
            self.tree_misurart.insert("", "end", values=row)

    # -------------------------
    # Tab Report
    # -------------------------
    def _build_tab_report(self):
        frm = ttk.Frame(self.tab_report, padding=8); frm.pack(fill="x")
        ttk.Label(frm, text="Report", font=("Arial", 14)).pack(anchor="w")
        btn_fr = ttk.Frame(frm); btn_fr.pack(fill="x", pady=6)
        ttk.Button(btn_fr, text="MEMORIE IN SCADENZA (30gg)", command=lambda: self.gui_memorie_scadenza(30)).pack(side="left", padx=6)
        ttk.Button(btn_fr, text="CHECKLIST CLIENTI", command=self.gui_report_checklist_database_clienti).pack(side="left", padx=6)
        ttk.Button(btn_fr, text="CAMBI DGFE", command=self.gui_report_cambi_dgfe_prompt).pack(side="left", padx=6)
        ttk.Button(btn_fr, text="CONFRONTO INTAE/MISURART", command=self.gui_confronto_intae_misurart).pack(side="left", padx=6)

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

    def gui_report_rt_per_model(self):
        """
        Robust version: forza colonne/headings, seleziona il tab che contiene il Treeview (se presente),
        inserisce righe e porta in vista la prima riga. Mostra un popup che conferma il numero di righe.
        Usa _safe_call(database.report_rt_for_clienti) come nella tua app.
        """
        import traceback
        try:
            print("DEBUG: pulsante report premuto")
            rows = _safe_call(database.report_rt_for_clienti) or []
            print("DEBUG: righe ricevute =", len(rows))

            # intestazioni (stesso ordine delle tuple restituite dalla funzione DB)
            headers = ("Cliente", "Indirizzo", "Comune", "Matricola", "Data VP", "Tipo", "Data Inst.", "Azz")

            # --- Se il tree è in un Notebook, seleziona il tab che lo contiene (ricerca ricorsiva) ---
            try:
                def find_notebook_ancestor(widget):
                    w = widget
                    while True:
                        parent = getattr(w, "master", None)
                        if parent is None:
                            return None
                        if isinstance(parent, ttk.Notebook):
                            return parent
                        w = parent

#RICORDA: ho fatto in modo che se richiamo questa funzione viene aperto il report nella tab report
#per rimuovere questa feature bsta togliere il blocco qui sotto
                nb = find_notebook_ancestor(self.tree_report)
                if nb is not None:
                    for tab_id in nb.tabs():
                        try:
                            tab_widget = nb.nametowidget(tab_id)
                        except Exception:
                            tab_widget = None

                        def is_ancestor(ancestor, widget):
                            w = widget
                            while True:
                                if w is None:
                                    return False
                                if w == ancestor:
                                    return True
                                w = getattr(w, "master", None)

                        if tab_widget is not None and is_ancestor(tab_widget, self.tree_report):
                            try:
                                nb.select(tab_id)
                                nb.update_idletasks()
                            except Exception:
                                pass
                            break
            except Exception:
                traceback.print_exc()

            # --- Svuota e imposta colonne/headings ---
            try:
                try:
                    self.tree_report.delete(*self.tree_report.get_children())
                except Exception:
                    pass

                self.tree_report["columns"] = headers
                self.tree_report["show"] = "headings"

                for h in headers:
                    self.tree_report.heading(h, text=h)
                    if h == "Cliente":
                        w = 250
                    elif h == "Indirizzo":
                        w = 240
                    elif h == "Comune":
                        w = 130
                    elif h == "Matricola":
                        w = 120
                    elif h in ("Data VP", "Data Inst."):
                        w = 120
                    elif h == "Tipo":
                        w = 80
                    else:
                        w = 80
                    try:
                        self.tree_report.column(h, width=w, anchor="w", stretch=False)
                    except Exception:
                        pass
            except Exception:
                traceback.print_exc()

            # --- Inserisci le righe ---
            first_item_id = None
            for r in rows:
                vals = tuple('' if v is None else str(v) for v in r)
                if len(vals) < len(headers):
                    vals = vals + tuple("" for _ in range(len(headers) - len(vals)))
                iid = self.tree_report.insert("", "end", values=vals)
                if first_item_id is None:
                    first_item_id = iid

            # --- Porta in vista e seleziona la prima riga ---
            try:
                if first_item_id:
                    self.tree_report.selection_set(first_item_id)
                    self.tree_report.focus(first_item_id)
                    self.tree_report.see(first_item_id)
            except Exception:
                traceback.print_exc()

            # Forza refresh GUI
            try:
                self.tree_report.update_idletasks()
                self.tree_report.update()
            except Exception:
                pass

            print("DEBUG: popolato tree con", len(rows), "righe")
            if len(rows) > 0:
                print("DEBUG: prima riga:", rows[0])

            # popup di conferma
            #try:
            #    from tkinter import messagebox
            #    messagebox.showinfo("Report", f"Report generato: {len(rows)} righe.")
            #except Exception:
            #    pass

        except Exception:
            print("DEBUG: errore in gui_report_rt_per_model:")
            traceback.print_exc()
            try:
                from tkinter import messagebox
                messagebox.showerror("Errore", "Errore durante la generazione del report. Controlla la console.")
            except Exception:
                pass

    def gui_report_checklist_database_clienti(self):
        res = _safe_call(database.report_checklist_clienti_selected) or []

        if not res:
            messagebox.showinfo("Checklist", "Nessun dato trovato.")
            return

        cols, rows = res[0]

        if not cols:
            messagebox.showinfo("Checklist", "Colonne non presenti nel DB.")
            return

        # RIUSO la tree_report già visibile (niente tab, niente notebook)
        self._clear_and_set_tree_columns(cols)

        for r in rows:
            self.tree_report.insert("", "end", values=r)

    def gui_report_cambi_dgfe_prompt(self):
        def do_search():
            d1_raw = ent1.get().strip(); d2_raw = ent2.get().strip()
            w.destroy()
            rows = _safe_call(database.report_cambi_dgfe, d1_raw, d2_raw) or []
            headers = ("Cliente", "Matricola", "Interventi")
            self._clear_and_set_tree_columns(headers)
            if not rows:
                messagebox.showinfo("Nessun risultato", "Nessun intervento trovato per le matricole nell'intervallo specificato.")
                return
            for r in rows:
                self.tree_report.insert("", "end", values=(r[0], r[1], str(r[2])))
        w = tk.Toplevel(self.root)
        ttk.Label(w, text="Data from (es. 01/01/2024):").pack(side="left", padx=4)
        ent1 = ttk.Entry(w); ent1.pack(side="left", padx=4)
        ttk.Label(w, text="Data to:").pack(side="left", padx=4)
        ent2 = ttk.Entry(w); ent2.pack(side="left", padx=4)
        ttk.Button(w, text="CERCA", command=do_search).pack(side="left", padx=6)

    def gui_confronto_intae_misurart(self):
        rows = _safe_call(database.confronto_intae_misurart) or []
        headers = ("Cliente", "Matricola", "Data INTAE", "Data MISURART")
        self._clear_and_set_tree_columns(headers)
        for r in rows:
            # r è (ragione_sociale, matricola, data_intae_iso, data_mis_iso)
            disp_intae = database._format_date_for_tree(r[2]) if r[2] else ""
            disp_mis = database._format_date_for_tree(r[3]) if r[3] else ""
            self.tree_report.insert("", "end", values=(r[0], r[1], disp_intae, disp_mis))

    # -------------------------
    # Autosave quando modifico un campo
    # -------------------------

    def enable_tree_autosave(self, debounce_ms: int = 800):
        """
        Attiva l'editing inline + autosave sulla tree. Chiamare una volta dopo
        che self.tree_report è stato creato e configurato.
        debounce_ms: tempo (ms) per raggruppare salvataggi rapidi (debounce).
        """
        self._tree_edit_after_id = None
        self._tree_edit_debounce_ms = debounce_ms
        # bind per iniziare edit su doppio click
        self.tree_report.bind("<Double-1>", lambda e: self._on_tree_double_click(e))

    def _on_tree_double_click(self, event):
        # identifica item e column
        tv = self.tree_report
        region = tv.identify("region", event.x, event.y)
        if region != "cell":
            return
        rowid = tv.identify_row(event.y)
        col = tv.identify_column(event.x)  # es. '#1', '#2', ...
        if not rowid or not col:
            return

        # bbox della cella
        x, y, width, height = tv.bbox(rowid, col)
        if width <= 0 or height <= 0:
            return

        # col index (0-based)
        col_idx = int(col.replace("#", "")) - 1
        # ottieni valore attuale e testo
        values = tv.item(rowid, "values")
        try:
            cur_value = values[col_idx]
        except Exception:
            cur_value = ""

        # crea Entry sopra la cella
        entry = tk.Entry(tv)
        entry.insert(0, cur_value)
        entry.place(x=x, y=y, width=width, height=height)
        entry.focus_set()
        # salva informazioni per callback
        entry._tree_rowid = rowid
        entry._tree_col_idx = col_idx
        entry._tree_col_name = None
        # se hai header memorizzati in self.tree_headers (tuple) prendi corrispondente
        if hasattr(self, "tree_headers") and len(self.tree_headers) > col_idx:
            entry._tree_col_name = self.tree_headers[col_idx]
        else:
            # fallback: userà index come chiave
            entry._tree_col_name = f"col_{col_idx}"

        # binding
        entry.bind("<Return>", lambda e, w=entry: self._on_tree_entry_commit(w))
        entry.bind("<FocusOut>", lambda e, w=entry: self._on_tree_entry_commit(w))
        entry.bind("<Escape>", lambda e, w=entry: self._on_tree_entry_cancel(w))

    def _on_tree_entry_cancel(self, entry_widget):
        try:
            entry_widget.destroy()
        except Exception:
            pass

    def _on_tree_entry_commit(self, entry_widget):
        """
        Commit immediato e autosave con debounce: aggiorna il valore visivo,
        costruisce il dict completo dei campi della riga e chiama save_all_cliente_fields.
        """
        try:
            new_val = entry_widget.get()
            rowid = entry_widget._tree_rowid
            col_idx = entry_widget._tree_col_idx
            col_name = entry_widget._tree_col_name
        except Exception:
            entry_widget.destroy()
            return

        # aggiorna visivamente la tree
        values = list(self.tree_report.item(rowid, "values"))
        # assicurati la tupla abbia abbastanza colonne
        if col_idx >= len(values):
            # allunga con stringhe vuote
            values += [""] * (col_idx - len(values) + 1)
        values[col_idx] = new_val
        self.tree_report.item(rowid, values=tuple(values))

        # rimuovi l'Entry
        entry_widget.destroy()

        # costruiamo il dict completo dei campi da salvare:
        # Mappiamo headers <-> valori. Assumiamo self.tree_headers = ("RAG_SOC","FINEMF","MF","TEL","CODICE",...)
        if hasattr(self, "tree_headers") and len(self.tree_headers) >= len(values):
            data = {self.tree_headers[i]: values[i] for i in range(len(values))}
        else:
            # fallback: usa nomi col_0, col_1...
            data = {f"col_{i}": values[i] for i in range(len(values))}

        # deduci rowid db: preferiamo avere iid==db_id o tag con id
        try:
            db_id = int(rowid)  # se hai settato iid=str(db_id)
        except Exception:
            # prova tags
            tags = self.tree_report.item(rowid, "tags")
            db_id = None
            if tags:
                for t in tags:
                    try:
                        db_id = int(t)
                        break
                    except Exception:
                        # tag like 'id:123'?
                        if isinstance(t, str) and t.startswith("id:"):
                            try:
                                db_id = int(t.split(":",1)[1])
                                break
                            except Exception:
                                pass
            if db_id is None:
                # non riesco a dedurre id: abortiamo salvataggio DB
                print("[_on_tree_entry_commit] impossibile dedurre db_id per row:", rowid)
                return

        # debounce: se l'utente modifica spesso, aspettiamo un piccolo intervallo prima di salvare
        if getattr(self, "_tree_edit_after_id", None):
            try:
                self.after_cancel(self._tree_edit_after_id)
            except Exception:
                pass

        def do_save():
            try:
                res = database.save_all_cliente_fields(db_id, data, run_post_update=True)
                if not res.get("ok", False):
                    print("[Autosave] errore salvataggio:", res.get("error"))
                else:
                    # opzionale: ricarica report solo se vuoi aggiornare la lista
                    # self.gui_memorie_scadenza(30)
                    pass
            except Exception as e:
                print("[Autosave] eccezione:", e)

        self._tree_edit_after_id = self.after(self._tree_edit_debounce_ms, do_save)

    

# Avvio GUI
def avvia_gui():
    root = tk.Tk()
    app = App(root)
    root.mainloop()

if __name__ == "__main__":
    avvia_gui()