# gui.py
import tkinter as tk
from tkinter import ttk, messagebox, filedialog
import csv
from datetime import datetime

import database

DATE_FMT = "%Y-%m-%d"  # convenzione: YYYY-MM-DD

class AvviaGUI:
    def __init__(self):
        self.root = tk.Tk()
        self.root.title("Gestione Misuratori Fiscali")
        self.root.geometry("1000x600")
        self._build_ui()

    def run(self):
        self.root.mainloop()

    def _build_ui(self):
        # notebook (tabs)
        notebook = ttk.Notebook(self.root)
        notebook.pack(fill="both", expand=True)

        # tabs
        self.tab_clienti = ttk.Frame(notebook)
        self.tab_intae = ttk.Frame(notebook)
        self.tab_misurart = ttk.Frame(notebook)
        self.tab_report = ttk.Frame(notebook)

        notebook.add(self.tab_clienti, text="Clienti")
        notebook.add(self.tab_intae, text="INTAE (Interventi AE)")
        notebook.add(self.tab_misurart, text="MISURART (Verifiche)")
        notebook.add(self.tab_report, text="Report & Confronti")

        # build contents
        self._build_tab_clienti()
        self._build_tab_intae()
        self._build_tab_misurart()
        self._build_tab_report()

    # ---------------------------
    # TAB CLIENTI
    # ---------------------------
    def _build_tab_clienti(self):
        frame_top = ttk.Frame(self.tab_clienti, padding=8)
        frame_top.pack(fill="x")

        # Form per aggiungere/modificare cliente
        form = ttk.LabelFrame(frame_top, text="Aggiungi / Modifica Cliente", padding=8)
        form.pack(side="left", fill="x", expand=True, padx=5, pady=5)

        labels = ["Ragione Sociale", "Indirizzo", "Comune", "Matricola", "Modello", "Data Installazione (YYYY-MM-DD)", "MXAZ (anni)", "QRCODE"]
        self.client_entries = {}
        for i, lab in enumerate(labels):
            ttk.Label(form, text=lab).grid(row=i, column=0, sticky="w", pady=2)
            ent = ttk.Entry(form, width=40)
            ent.grid(row=i, column=1, sticky="w", pady=2)
            self.client_entries[lab] = ent

        btn_frame = ttk.Frame(form)
        btn_frame.grid(row=len(labels), column=0, columnspan=2, pady=6)

        ttk.Button(btn_frame, text="Aggiungi Cliente", command=self._on_add_cliente).grid(row=0, column=0, padx=4)
        ttk.Button(btn_frame, text="Modifica Cliente Selezionato", command=self._on_modifica_cliente).grid(row=0, column=1, padx=4)
        ttk.Button(btn_frame, text="Duplica Cliente Selezionato", command=self._on_duplica_cliente).grid(row=0, column=2, padx=4)

        # Ricerca e lista
        right_frame = ttk.Frame(self.tab_clienti, padding=8)
        right_frame.pack(fill="both", expand=True)

        search_frame = ttk.Frame(right_frame)
        search_frame.pack(fill="x", pady=4)

        ttk.Label(search_frame, text="Cerca (ragione sociale / comune / matricola):").pack(side="left")
        self.client_search = ttk.Entry(search_frame, width=40)
        self.client_search.pack(side="left", padx=6)
        ttk.Button(search_frame, text="Cerca", command=self._on_search_clienti).pack(side="left", padx=4)
        ttk.Button(search_frame, text="Carica Tutti", command=self._on_load_all_clienti).pack(side="left", padx=4)

        # Treeview risultati
        cols = ("id", "ragione_sociale","indirizzo", "comune", "matricola", "modello", "data_installazione", "mxaz", "disattivazione")
        self.tree_clienti = ttk.Treeview(right_frame, columns=cols, show="headings", selectmode="browse")
        for c in cols:
            self.tree_clienti.heading(c, text=c)
            self.tree_clienti.column(c, width=120, anchor="w")
        self.tree_clienti.pack(fill="both", expand=True, pady=6)
        self.tree_clienti.bind("<<TreeviewSelect>>", self._on_select_cliente)

        # inizializza lista
        self._on_load_all_clienti()

    def _on_add_cliente(self):
        # raccogli campi
        try:
            rag = self.client_entries["Ragione Sociale"].get().strip()
            indirizzo = self.client_entries["Indirizzo"].get().strip()
            comune = self.client_entries["Comune"].get().strip()
            matricola = self.client_entries["Matricola"].get().strip()
            modello = self.client_entries["Modello"].get().strip()
            data_inst = self.client_entries["Data Installazione (YYYY-MM-DD)"].get().strip()
            mxaz = int(self.client_entries["MXAZ (anni)"].get().strip() or 0)
            qrcode = self.client_entries["QRCODE"].get().strip()
        except Exception as e:
            messagebox.showerror("Errore", f"Dati non validi: {e}")
            return

        if not rag or not matricola:
            messagebox.showerror("Errore", "Ragione sociale e matricola sono obbligatori")
            return

        # chiama DB
        try:
            scad_g, scad_mem = database.aggiungi_cliente(rag, indirizzo, comune, matricola, modello, data_inst, mxaz, qrcode)
            messagebox.showinfo("OK", f"Cliente aggiunto. Garanzia fino a: {scad_g.date()}, memoria fiscale fine appross: {scad_mem.date()}")
            self._on_load_all_clienti()
        except Exception as e:
            messagebox.showerror("Errore DB", str(e))

    def _on_search_clienti(self):
        term = self.client_search.get().strip()
        if not term:
            messagebox.showinfo("Ricerca", "Inserisci un termine di ricerca")
            return
        risultati = database.cerca_cliente(term)
        self._populate_tree(self.tree_clienti, risultati)

    def _on_load_all_clienti(self):
        risultati = database.cerca_cliente("")  # cerca vuoto -> restituisce tutti (LIKE '%%')
        self._populate_tree(self.tree_clienti, risultati)

    def _populate_tree(self, tree, rows):
        # pulisci
        for r in tree.get_children():
            tree.delete(r)
        for row in rows:
            # row corrisponde alla riga della tabella clienti:
            # (id, ragione_sociale, indirizzo, comune, matricola, modello, data_installazione, mxaz, disattivazione, qrcode)
            # attenzione: ordini colonne coerenti con DB
            tree.insert("", "end", values=row)

    def _on_select_cliente(self, event):
        sel = self.tree_clienti.selection()
        if not sel:
            return
        item = self.tree_clienti.item(sel[0])["values"]
        # popola il form con i dati (occhio agli indici secondo il DB)
        # DB definisce: id, ragione_sociale, indirizzo, comune, matricola, modello, data_installazione, mxaz, disattivazione, qrcode
        try:
            self.client_entries["Ragione Sociale"].delete(0, tk.END); self.client_entries["Ragione Sociale"].insert(0, item[1])
            self.client_entries["Indirizzo"].delete(0, tk.END); self.client_entries["Indirizzo"].insert(0, item[2] or "")
            self.client_entries["Comune"].delete(0, tk.END); self.client_entries["Comune"].insert(0, item[3] or "")
            self.client_entries["Matricola"].delete(0, tk.END); self.client_entries["Matricola"].insert(0, item[4] or "")
            self.client_entries["Modello"].delete(0, tk.END); self.client_entries["Modello"].insert(0, item[5] or "")
            self.client_entries["Data Installazione (YYYY-MM-DD)"].delete(0, tk.END); self.client_entries["Data Installazione (YYYY-MM-DD)"].insert(0, item[6] or "")
            self.client_entries["MXAZ (anni)"].delete(0, tk.END); self.client_entries["MXAZ (anni)"].insert(0, item[7] or 0)
            self.client_entries["QRCODE"].delete(0, tk.END); self.client_entries["QRCODE"].insert(0, item[9] or "")
        except Exception:
            pass

    def _on_modifica_cliente(self):
        sel = self.tree_clienti.selection()
        if not sel:
            messagebox.showinfo("Seleziona", "Seleziona un cliente dalla lista prima di modificare")
            return
        item = self.tree_clienti.item(sel[0])["values"]
        id_cliente = item[0]
        # prepara kwargs in base ai campi compilati
        kwargs = {}
        mapping = {
            "ragione_sociale": "Ragione Sociale",
            "indirizzo": "Indirizzo",
            "comune": "Comune",
            "matricola": "Matricola",
            "modello": "Modello",
            "data_installazione": "Data Installazione (YYYY-MM-DD)",
            "mxaz": "MXAZ (anni)",
            "qrcode": "QRCODE"
        }
        for dbfield, label in mapping.items():
            val = self.client_entries[label].get().strip()
            if val != "":
                # tipo conversion per mxaz
                if dbfield == "mxaz":
                    try:
                        kwargs[dbfield] = int(val)
                    except:
                        kwargs[dbfield] = 0
                else:
                    kwargs[dbfield] = val
        try:
            database.modifica_cliente(id_cliente, **kwargs)
            messagebox.showinfo("OK", "Cliente aggiornato")
            self._on_load_all_clienti()
        except Exception as e:
            messagebox.showerror("Errore DB", str(e))

    def _on_duplica_cliente(self):
        sel = self.tree_clienti.selection()
        if not sel:
            messagebox.showinfo("Seleziona", "Seleziona un cliente dalla lista prima di duplicare")
            return
        item = self.tree_clienti.item(sel[0])["values"]
        id_cliente = item[0]
        ok = database.duplica_cliente(id_cliente)
        if ok:
            messagebox.showinfo("OK", "Cliente duplicato")
            self._on_load_all_clienti()
        else:
            messagebox.showerror("Errore", "Impossibile duplicare il cliente")

    # ---------------------------
    # TAB INTAE
    # ---------------------------
    def _build_tab_intae(self):
        top = ttk.Frame(self.tab_intae, padding=8)
        top.pack(fill="x")

        ttk.Label(top, text="Import INTAE (CSV) e ricerca interventi", font=("Arial", 12)).pack(anchor="w")

        btn_frame = ttk.Frame(top)
        btn_frame.pack(fill="x", pady=6)

        ttk.Button(btn_frame, text="Importa CSV INTAE", command=self._on_import_intae).pack(side="left", padx=6)
        ttk.Button(btn_frame, text="Cerca interventi per matricola", command=self._on_search_intae_prompt).pack(side="left", padx=6)

        # Treeview visualizzazione interventi INTAE
        cols = ("id", "matricola", "data_intervento", "codice_intervento", "cf_tecnico")
        self.tree_intae = ttk.Treeview(self.tab_intae, columns=cols, show="headings")
        for c in cols:
            self.tree_intae.heading(c, text=c)
            self.tree_intae.column(c, width=140)
        self.tree_intae.pack(fill="both", expand=True, padx=8, pady=8)

    def _on_import_intae(self):
        path = filedialog.askopenfilename(title="Seleziona file CSV INTAE", filetypes=[("CSV","*.csv"),("All","*.*")])
        if not path:
            return

        lista = []
        # Assumiamo CSV con intestazioni: matricola,data,codice,cf
        # Se il tuo CSV ha nomi diversi, adattare qui la mappatura.
        try:
            with open(path, newline='', encoding="utf-8") as f:
                reader = csv.DictReader(f)
                for r in reader:
                    # normalize keys (lowercase)
                    # try find fields in common variants
                    def getfield(d, names):
                        for n in names:
                            if n in d and d[n] != "":
                                return d[n]
                        return ""
                    matricola = getfield(r, ["matricola","MATRICOLA","Matricola"])
                    data = getfield(r, ["data","DATA","Data","data_intervento","DATA_INTERVENTO"])
                    codice_raw = getfield(r, ["codice","CODICE","Codice","codice_intervento"])
                    cf = getfield(r, ["cf","CF","cf_tecnico","CF_TECNICO"])
                    # try convert codice to int
                    try:
                        codice = int(codice_raw)
                    except:
                        codice = None
                    lista.append({"matricola": matricola, "data": data, "codice": codice, "cf": cf})
            database.importa_intae_da_csv(lista)
            messagebox.showinfo("Import", f"Import completato ({len(lista)} righe).")
        except Exception as e:
            messagebox.showerror("Errore import", str(e))

    def _on_search_intae_prompt(self):
        def do_search():
            matricola = entry.get().strip()
            win.destroy()
            if matricola:
                res = database.cerca_interventi_intae(matricola)
                self._populate_tree(self.tree_intae, res)
        win = tk.Toplevel(self.root)
        win.title("Cerca INTAE per matricola")
        ttk.Label(win, text="Matricola:").pack(side="left", padx=6, pady=6)
        entry = ttk.Entry(win)
        entry.pack(side="left", padx=6, pady=6)
        ttk.Button(win, text="Cerca", command=do_search).pack(side="left", padx=6, pady=6)

    # ---------------------------
    # TAB MISURART
    # ---------------------------
    def _build_tab_misurart(self):
        top = ttk.Frame(self.tab_misurart, padding=8)
        top.pack(fill="x")

        ttk.Label(top, text="Registrazione Verifiche / Disattivazioni", font=("Arial", 12)).pack(anchor="w")

        frm = ttk.Frame(top)
        frm.pack(fill="x", pady=6)
        ttk.Label(frm, text="Matricola:").grid(row=0, column=0, padx=4, pady=2, sticky="w")
        self.mis_matricola = ttk.Entry(frm, width=30); self.mis_matricola.grid(row=0, column=1, padx=4, pady=2)
        ttk.Label(frm, text="Data (YYYY-MM-DD):").grid(row=0, column=2, padx=4, pady=2, sticky="w")
        self.mis_data = ttk.Entry(frm, width=20); self.mis_data.grid(row=0, column=3, padx=4, pady=2)

        ttk.Button(frm, text="Registra Verifica (tipo verif)", command=self._on_registra_verifica).grid(row=1, column=0, padx=4, pady=6)
        ttk.Button(frm, text="Registra Disattivazione (tipo 2)", command=self._on_registra_disattivazione).grid(row=1, column=1, padx=4, pady=6)
        ttk.Button(frm, text="Cerca interventi MISURART per matricola", command=self._on_search_misurart_prompt).grid(row=1, column=2, padx=4, pady=6)

        # Tree misurart
        cols = ("id", "matricola", "data_intervento", "tipo")
        self.tree_misurart = ttk.Treeview(self.tab_misurart, columns=cols, show="headings")
        for c in cols:
            self.tree_misurart.heading(c, text=c)
            self.tree_misurart.column(c, width=150)
        self.tree_misurart.pack(fill="both", expand=True, padx=8, pady=8)

    def _on_registra_verifica(self):
        mat = self.mis_matricola.get().strip()
        data = self.mis_data.get().strip() or datetime.now().strftime(DATE_FMT)
        if not mat:
            messagebox.showerror("Errore", "Inserisci matricola")
            return
        try:
            database.registra_verificazione(mat, data)
            messagebox.showinfo("OK", "Verifica registrata")
        except Exception as e:
            messagebox.showerror("DB error", str(e))

    def _on_registra_disattivazione(self):
        mat = self.mis_matricola.get().strip()
        data = self.mis_data.get().strip() or datetime.now().strftime(DATE_FMT)
        if not mat:
            messagebox.showerror("Errore", "Inserisci matricola")
            return
        try:
            database.registra_disattivazione(mat, data)
            messagebox.showinfo("OK", "Disattivazione registrata e campo disattivazione aggiornato")
        except Exception as e:
            messagebox.showerror("DB error", str(e))

    def _on_search_misurart_prompt(self):
        def do_search():
            matricola = entry.get().strip()
            win.destroy()
            if matricola:
                res = database.cerca_interventi_misurart(matricola)
                self._populate_tree(self.tree_misurart, res)
        win = tk.Toplevel(self.root)
        win.title("Cerca MISURART per matricola")
        ttk.Label(win, text="Matricola:").pack(side="left", padx=6, pady=6)
        entry = ttk.Entry(win)
        entry.pack(side="left", padx=6, pady=6)
        ttk.Button(win, text="Cerca", command=do_search).pack(side="left", padx=6, pady=6)

    # ---------------------------
    # TAB REPORT & CONFRONTI
    # ---------------------------
    def _build_tab_report(self):
        top = ttk.Frame(self.tab_report, padding=8)
        top.pack(fill="x")

        ttk.Label(top, text="Report e Confronti", font=("Arial", 12)).pack(anchor="w")

        frm = ttk.Frame(top)
        frm.pack(fill="x", pady=6)

        ttk.Button(frm, text="Verifiche / Memorie in scadenza (30 gg)", command=self._on_report_scadenze).grid(row=0, column=0, padx=6, pady=4)
        ttk.Button(frm, text="RT per modello (INTAE tipo 3)", command=self._on_rt_per_modello).grid(row=0, column=1, padx=6, pady=4)
        ttk.Button(frm, text="Confronto INTAE vs MISURART (matricola)", command=self._on_confronto_prompt).grid(row=0, column=2, padx=6, pady=4)

        # area testo risultati
        self.txt_report = tk.Text(self.tab_report, height=20)
        self.txt_report.pack(fill="both", expand=True, padx=8, pady=8)

    def _on_report_scadenze(self):
        rows = database.report_memorie_fiscali_in_scadenza(30)
        self.txt_report.delete("1.0", tk.END)
        if not rows:
            self.txt_report.insert(tk.END, "Nessuna memoria fiscale in scadenza nei prossimi 30 giorni.\n")
            return
        for r in rows:
            # r: row completa della tabella clienti
            self.txt_report.insert(tk.END, f"ID:{r[0]} - {r[1]} - Matricola:{r[4]} - Installazione:{r[6]} - MXAZ:{r[7]}\n")

    def _on_rt_per_modello(self):
        # cerca interventi INTAE di tipo 3 e raggruppa per modello leggendo la tabella clienti
        intae = database.report_interventi_tipo(3) if hasattr(database, "report_interventi_tipo") else database.cerca_interventi_intae  # fallback
        # implement minimal grouping using DB queries (if function exists in DB preferirla)
        # Here we do a manual join-like mapping:
        all_intae = database.cerca_interventi_intae("")  # if empty returns none in our impl, but try...
        # safer: query DB directly
        conn = database.get_connection()
        cur = conn.cursor()
        cur.execute("SELECT i.id, i.matricola, i.data_intervento, i.codice_intervento, c.modello FROM intae i LEFT JOIN clienti c ON i.matricola = c.matricola WHERE i.codice_intervento = 3")
        rows = cur.fetchall()
        conn.close()

        if not rows:
            self.txt_report.delete("1.0", tk.END)
            self.txt_report.insert(tk.END, "Nessun RT (codice 3) trovato.\n")
            return

        # raggruppa per modello
        gruppi = {}
        for r in rows:
            modello = r[4] or "Sconosciuto"
            gruppi.setdefault(modello, 0)
            gruppi[modello] += 1
        self.txt_report.delete("1.0", tk.END)
        for m, cnt in gruppi.items():
            self.txt_report.insert(tk.END, f"{m}: {cnt} interventi (RT)\n")

    def _on_confronto_prompt(self):
        def do_confronto():
            mat = entry.get().strip()
            win.destroy()
            if not mat:
                return
            res = database.confronta_intae_misurart(mat)
            self.txt_report.delete("1.0", tk.END)
            self.txt_report.insert(tk.END, "Interventi presenti solo su INTAE:\n")
            for it in res["solo_intae"]:
                self.txt_report.insert(tk.END, f"  {it}\n")
            self.txt_report.insert(tk.END, "\nInterventi presenti solo su MISURART:\n")
            for it in res["solo_misurart"]:
                self.txt_report.insert(tk.END, f"  {it}\n")

        win = tk.Toplevel(self.root)
        win.title("Confronto INTAE vs MISURART per matricola")
        ttk.Label(win, text="Matricola:").pack(side="left", padx=6, pady=6)
        entry = ttk.Entry(win)
        entry.pack(side="left", padx=6, pady=6)
        ttk.Button(win, text="Confronta", command=do_confronto).pack(side="left", padx=6, pady=6)

    # ---------------------------
    # UTILITY
    # ---------------------------
    def _populate_tree(self, tree, rows):
        # clear
        for r in tree.get_children():
            tree.delete(r)
        for row in rows:
            tree.insert("", "end", values=row)

