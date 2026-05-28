<h1 align="center">📦 GiaLoSai</h1>

<p align="center"><b>Sistema integrato per la gestione del magazzino, contabilità aziendale e monitoraggio delle vendite.</b></p>

<p align="center">
  <img src="https://img.shields.io/badge/Language-Go-00ADD8?style=flat&logo=go" alt="Go">
  <img src="https://img.shields.io/badge/Status-In%20Development-yellow" alt="Status">
</p>

---

## 🎯 Overview 

**GiaLoSai** nasce nell'ottica di poter fornire a un'azienda uno strumento utile, leggero e puntuale per il controllo del magazzino, una contabilità precisa e una gestione efficiente dei flussi di vendita.

---

## 🚀 Quickstart

Di seguito, i passaggi per configurare, compilare ed eseguire il progetto in locale utilizzando il `Makefile` incluso.

### Prerequisiti
Prima di iniziare, assicurati di avere installato sul tuo sistema:
* **Go** (versione 1.16 o superiore)
* **Make** (disponibile tramite il pacchetto `build-essential` su sistemi Debian/Ubuntu)
* **Git**

### 1. Clonazione del Repository

#### **Via HTTPS**
```bash
 git clone [https://github.com/tommasogarofani/ingsw_25-26_group-11.git](https://github.com/tommasogarofani/ingsw_25-26_group-11.git)
 cd ingsw_25-26_group-11
```

#### **Via SSH**
```bash
 git clone git@github.com:tommasogarofani/ingsw_25-26_group-11.git
 cd ingsw_25-26_group-11
```

### 2. Compilazione ed Esecuzione
Grazie al `Makefile` incluso, è possibile gestire l'intero ciclo di vita dell'applicazione con comandi immediati dal terminale:

* **Per compilare ed avviare l'applicazione direttamente:**
```bash
 make run
```
*(Questo comando creerà la cartella bin/, compilerà il codice sorgente generando l'eseguibile e lo avvierà automaticamente).*

* **Solo per compilare (senza avviare):**
```bash
  make compile
```
*(L'eseguibile verrà generato all'interno della cartella ./bin).*

* **Pulizia dei file generati:**
  Se si desidera rimuovere i file binari e riportare la directory allo stato iniziale pulito:
```bash
  make clean
```

---

## 💡 Usage

L'applicazione è progettata per adattarsi a diversi profili aziendali e flussi di lavoro:

* **Magazziniere:** Gestisce l'inventario monitorando in tempo reale i flussi di merci tra entrate (carico) e uscite (scarico) di prodotti.
* **Contabile:** Gestisce il ciclo di fatturazione in modo puntuale, verificando la corrispondenza tra bolle di acquisto, bolle di vendita e documenti di trasporto.
* **Proprietario / Amministratore:** Dispone di una panoramica globale sull'andamento dell'attività per avere metriche e controllo aziendale sempre sotto controllo.

---

## ⚙️ Configuration

Il software è pronto all'uso e non richiede configurazioni iniziali. Eventuali parametri predefiniti possono essere modificati all'interno dei file di configurazione dedicati nella directory principale.

---

## 🧪 Test

Il progetto include una suite di test automatizzati per verificare l'integrità del codice. Di seguito comando dedicato del `Makefile`:
```bash
 make test
```

---

## 🤝 Contributing

Se si vuole contribuire al progetto, segui le linee guida del team:

1. **Branching:** Non è mai consentito lavorare direttamente sul branch `main`. E' necessario creare un branch dedicato per ogni feature o fix (es. `feature/nome-funzionalita` o `bugfix/nome-bug`).
2. **Pull Requests:** Una volta completato il lavoro, aprire una Pull Request verso `main`. Il codice deve essere revisionato e approvato da almeno un altro membro del gruppo prima del merge.
3. **Codice e Test:** Assicurarsi che il codice sia formattato correttamente con `go fmt` e che tutti i test passino (`make test`) prima di effettuare il push.
4. **Aggiornamento:** Quando si lavora al progetto è OBBLIGATORIO aggiornare tutte le parti della documentazione necessaria (README, documenazione... .tex, architettura.md, ecc.) e i test sui componenti singoli o di integrazione.

---

## 👥 Credits (Group 11)

Questo progetto è sviluppato come attività collaborativa per il corso di **Ingegneria del Software** (A.A. 2025/2026).

* **Tommaso Garofani** - [@tommasogarofani](https://github.com/tommasogarofani)
* **David Cardinali** - [@Dyved33](https://github.com/Dyved33)
* **Luca Rosadi** - [@iperblastorm](https://github.com/iperblastorm)

---