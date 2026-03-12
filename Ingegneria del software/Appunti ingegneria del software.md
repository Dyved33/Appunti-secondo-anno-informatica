### 26-02-25 lz.1
**Ingegneria del software oggi**
Quando si parla di software, si parla della struttura portante dei sistemi moderni
- trasporti
- comunicazione 
- finanza
E i problemi che genera quando fallisce sono pericolosamente alti
- interruzione dei servizi essenziali
- problema di sicurezza e privacy
- costi economici e perdita di fiducia degli utenti
Per **L'ingegneria del software** non si parla solo di programmare, ma serve un approccio **ingegneristico** per: 
- Progettare
- costruire
- mantenere
Sistemi complessi e duraturi

**Quando il software fallisce**
I problemi, oltre ad essere "bug fastidiosi", creano anche dei problemi reali
- sistemi di prenotazioni in down durante i picchi di accesso
- piattaforme di pagamento non disponibili
- errori nei sistemi di controllo 
- perdita o esposizione dei dati sensibili degli utenti
E in questi casi la risoluzione non consiste solo nella correzione di una riga ma può prevedere anche la ri-progettazione dell'intero sistema. L'**Ingegneria del software** nasce per ridurre la probabilità che questi scenari si verifichino.

**Progettare sistemi complessi**
Nei sistemi reali entrano in gioco diversi fattori
- molti operatori coinvolti
- vincoli tecnici e organizzativi
- grandi moli di dati
- interazione tra componenti eterogenei
- requisiti che evolvono
- tecnologia che cambia
- contesti d'uso diversificati
L’Ingegneria del Software nasce per gestire questa complessità e fornire strumenti, metodi e pratiche per affrontare lo sviluppo in modo sistematico.

**Programmazione $\neq$ ingegneria del software**
La scrittura del codice è solo una delle attività dello sviluppo software e rappresenta una porzione limitata del lavoro complessivo.
Prima dell’implementazione c’è la fase di analisi e comprensione del problema:
- definizione dei requisiti
- identificazione dei vincoli
- analisi del dominio.
Dopo l’implementazione inizia un percorso ancora più lungo:
- verifica del comportamento
- correzione dei difetti
- gestione dei cambiamenti
- evoluzione del sistema nel tempo.
L’Ingegneria del Software tiene insieme tutto questo, garantendo che il sistema rimanga comprensibile e modificabile nel tempo, indipendentemente da chi scrive il codice.

**Che cos'è l'ingegneria del software**
L’Ingegneria del Software è l’insieme di principi, metodi e strumenti per sviluppare e mantenere sistemi software affidabili e utili in modo economicamente sostenibile.
L'attività di sviluppo si suddivide principalmente in 4 processi:
- **Specifica**: capire che cosa il software deve fare e quali vincoli deve rispettare.
- **Sviluppo**: progettare l’architettura e le strutture interne, scrivere il codice.
- **Validazione**: verificare che il sistema soddisfi davvero i requisiti attesi.
- **Evoluzione**: adattare il software a nuove esigenze, correzioni, cambiamenti di contesto.
**Il costo reale del software è nella manutenzione**
Nei progetti reali, la fase  iniziale di sviluppo rappresenta solo una parte del costo complessivo. Ma è la **manutenzione** che assorbe la maggior parte delle risorse
- correzione dei bug
- aggiornamenti richiesti
- adattamenti a normative e tecnologie
- estensioni e miglioramenti continui
La sostenibilità del software dipende dalla qualità delle scelte fatte nelle fasi iniziali. Un buon progetto facilita l’evoluzione, mentre un progetto fragile la rende costosa e rischiosa.

**Ciclo di vita**
Il **ciclo di vita** fornisce un quadro concettuale per organizzare le attività
- raccolta dei requisiti
- progettazione soluzione
- implementazione
- verifica
- rilascio
- manutenzione
1. Le attività iniziali seguono un flusso lineare: requisiti → progettazione → implementazione.
2. Dopo la prima versione, il processo entra in un ciclo continuo di test, rilascio e manutenzione.
3. Le attività della seconda riga sono iterative: ogni modifica o bug-fix riattiva il ciclo.
4. La manutenzione genera spesso nuovi requisiti, chiudendo il loop verso la fase iniziale.

**Comprensione del problema**
La costruzione del software deve partire da una reale **comprensione** del problema da risolvere e delle esigenze degli utenti.
- Chi userà il sistema? 
- Con quali obiettivi? 
- In quali contesti? 
- Con quali vincoli?
Un’analisi affrettata porta a sviluppare software poco utile o costoso da correggere; dedicare tempo alla comprensione iniziale evita errori di fondo.

**Nessun progetto procede in linea retta**
Durante lo sviluppo emergono inevitabilmente nuove informazioni e nuove esigenze, che mettono in discussione decisioni prese in precedenza. Il processo reale è **iterativo**
- si torna sui requisiti
- si rivede il design
- si modifica il codice
- si aggiornano i test.
L’ingegneria del software non elimina queste iterazioni, ma fornisce strumenti per gestirle in modo controllato e produttivo.

**La qualità non è un optional**
Un sistema software non deve solo funzionare: deve funzionare bene, con:
- affidabilità
- efficienza
- semplicità d’uso 
- capacità di evolvere nel tempo
La qualità si manifesta su due livelli:
- **Front-end**: esperienza d'uso, chiarezza, prestazioni
- **Back-end**: architettura, struttura del codice, dipendenze, documentazione
Se la qualità interna è scarsa, il sistema diventa fragile e difficile da modificare, indipendentemente da come appare all’ esterno.

**Le figure presenti nello sviluppo di un progetto**
Un sistema software è il risultato del lavoro di molte figure diverse con ruoli differenti e competenze specifiche.
- **utenti e committenti**:  definiscono bisogni e obiettivi
- **analisti e progettisti**: traducono le esigenze in modelli e architetture
- **sviluppatori e tester**: costruiscono e verificano il sistema
- **manutentori**: garantiscono la continuità del servizio nel tempo
La collaborazione tra questi ruoli è un fattore determinante per il successo del progetto.
Ogni ruolo ha una prospettiva diversa sul sistema e l’ingegneria del software deve farle convivere in modo coerente. 
- **l’utente** si concentra sull'esperienza d’uso e sulla semplicità delle operazioni.
- **lo sviluppatore** guarda alla correttezza e alla struttura del codice.
- **l’architetto** pensa alla coerenza strutturale e alle dipendenze tra moduli;
- **il manutentore** si preoccupa della comprensibilità interna e della facilità di modifica.
Integrare queste prospettive è fondamentale per ottenere un sistema davvero efficace.

**Quando la comunicazione fallisce, fallisce il progetto**
Molti difetti e ritardi non derivano da problemi tecnici, ma da incomprensioni tra le persone coinvolte. 
- requisiti fraintesi o ambigui
- aspettative non esplicitate
- decisioni progettuali non documentate
- informazioni critiche che non vengono condivise
Una Comunicazione efficace è **parte integrante dell’ingegneria del software**, tanto quanto il codice.
 
 **Il rischio è parte naturale di ogni progetto**
Ogni decisione progettuale e tecnologica introduce incertezza. 
- scelte tecnologiche premature o non verificate.
- requisiti instabili o in continuo cambiamento.
- integrazioni difficili o vincoli non previsti.
- stime errate di tempi e complessità.
L'ingegneria del software non elimina i rischi ma li gestisce in modo da correggerlo in tempo.

**La complessità è il vero avversario del software**
La difficoltà principale nello sviluppo del software non risiede nelle singole funzionalità, ma nell'interazione tra molte parti diverse.
La complessità non è solo quantitativa, ma anche qualitativa:
- dipendenze poco chiare
- eccezioni alla regola
- comportamenti impliciti
- interazioni inaspettate.
Senza un metodo rigoroso, la complessità rende ogni modifica rischiosa e tende a far degenerare il sistema nel tempo. Anche se le modifiche sembrano piccole e locali, la complessità interna del sistema tende a crescere in modo non lineare nel corso dei rilasci, soprattutto se mancano buone pratiche di progettazione, pattern, principi SOLIDI e test automatici. 
L’ingegneria del software offre strumenti per mantenere comprensibile un sistema complesso.
- architetture modulari, che dividono il sistema in parti coerenti
- modelli chiari, che permettono di ragionare sul dominio e sulle interazioni
- responsabilità ben definite, che evitano sovrapposizioni e ambiguità
- documentazione mirata, che supporta la comprensione senza appesantire.

**Obiettivo dell' ingegnere del software**
L'obbiettivo dell'ingegneria del software si pone come obbiettivo finale in un progetto la **riduzione dell'incertezza** durante lo sviluppo
- prevedere i tempi e i costi necessari
- valutare diverse alternative prima dell'implementazione
- rilevare problemi in anticipo
- mantenere il controllo anche in presenza di cambiamenti
L'approccio ingegneristico non elimina gli imprevisti, ma rende possibile affrontarli con metodo senza perdere stabilità.

---
### 26-02-26 lz.2
**Processi software**
Un' attività collettiva che coinvolge persone, ruoli e decisioni distribuite nel tempo. Quando si lavora con più persone, ognuna tende a seguire il proprio stile, creando
- fraintendimenti
- funzioni duplicate nel codice
- parti diverse del sistema che lavorano incoerentemente
Un **Processo software** fornisce una struttura comune per decidere come:
- procedere 
- coordinarsi
- verificare il tutto in maniera sistematica
L'obbiettivo è ridurre il caos tra i vari team, e le persone all'interno dei vari team. Si devono stilare delle regole per agire coerentemente e in maniera ordinata, sfruttando così al massimo la potenzialità del gruppo

**Dal ciclo di vita al software**
Nella pratica ci si chiede:
- in quale ordine vanno affrontate le attività
- quanto approfondire ciascuna fase
- quanto spesso tornare indietro per rivedere le decisioni già prese
I processi software rispondono proprio a queste domande:
- sono strategie organizzative per combinare le stesse fasi del ciclo di vita (waterfall, incrementale, spirale, agile, . . . )
- scegliamo il processo in base al contesto e al livello di rischio / incertezza

**Definizione (di processo)** 
Un insieme strutturato di attività, ruoli e artefatti che guidano la costruzione e l'evoluzione di un software. In un processo si specifica:
- chi è responsabile di questa cosa
- quali passi devono essere compiuti e in quale sequenza
- quali prodotti intermedi devono essere generati (documenti, modelli, codice, test, . . . )
L’obiettivo non è aggiungere burocrazia, ma garantire che il progetto proceda in modo coerente, tracciabile e comprensibile per tutti i soggetti coinvolti.

**Perché serve anche ai buoni programmatori**
Nei piccoli progetti individuali il programmatore può controllare tutto “nella propria testa”, ma questa illusione sparisce appena: 
- il numero di funzionalità cresce
- entra in gioco un team
- il software deve essere mantenuto nel tempo
Senza un processo, le decisioni restano implicite e si perdono facilmente:
- diventa difficile capire perché certe scelte sono state fatte
- diventa rischioso modificare parti del sistema già rilasciate
Un processo aiuta a rendere esplicito il modo di lavorare, in modo che qualità e conoscenza non
dipendano dai singoli individui ma dal team.

**Attività comuni a tutti i processi**
Qualunque sia il processo adottato, ritroveremo sempre alcune attività fondamentali:
- specifica: capire che cosa il sistema deve fare e quali vincoli deve rispettare;
- sviluppo: progettare e implementare la soluzione;
- validazione: verificare che il sistema soddisfi i requisiti e si comporti come atteso;
- evoluzione: modificare e mantenere il sistema in risposta a nuove esigenze.

**Perché esistono diversi modelli**
Non tutti i progetti sono uguali: alcuni hanno requisiti molto stabili e vincoli normativi forti; altri nascono in ambienti altamente dinamici, con idee ancora da esplorare. Applicare sempre lo stesso processo, senza guardare al contesto, porta facilmente a inefficienze:
- troppa rigidità dove servirebbe flessibilità
- troppa improvvisazione dove servirebbe pianificazione
I modelli di processo offrono quindi strategie diverse, che l’ingegnere del software deve conoscere per poter scegliere (e adattare) quella più adatta al progetto.

**Modelli  di processo**
Nel seguito vedremo tre modelli classici di processo software.
Non sono schemi rigidi e mutuamente esclusivi: nella pratica molti progetti reali combinano elementi di questi modelli, adattandoli al contesto. L’obiettivo non è “scegliere il vincitore”, ma capire quando e perché ciascun modello può essere più adatto.

**Modello waterfall**
Il modello waterfall (a cascata) organizza lo sviluppo come una sequenza di fasi ben definite, che procedono in modo ordinato e lineare. Ogni fase produce artefatti che diventano l’input per la fase successiva:
- progettazione 
- implementazione 
- validazione 
- rilascio e manutenzione
La logica è quella di minimizzare i ritorni indietro: si passa alla fase successiva solo quando la
precedente è considerata completa.

*Vantaggi*
La struttura sequenziale è facile da comunicare a tutti gli stakeholder e rende chiaro in quale fase
ci si trovi in ogni momento. La forte enfasi sulla documentazione iniziale favorisce:
- la tracciabilità delle decisioni
- la possibilità di svolgere audit e verifiche formali
- l’allineamento tra team tecnici e figure non tecniche.
È particolarmente adatto quando i requisiti sono ben noti e il contesto cambia lentamente, ad esempio in alcuni ambiti industriali regolamentati.

*Limiti*
Nella realtà, i requisiti cambiano spesso durante lo sviluppo e raramente sono completi o corretti fin dall'inizio. Nel waterfall, tornare indietro a modificare fasi già “chiuse” è costoso e complicato:
- correggere un requisito errato dopo aver progettato e implementato il sistema può richiedere molto lavoro.
- Gli utenti vedono il sistema completo solo alla fine, quando è difficile correggere scelte di fondo: il rischio di scoprire tardi malintesi o mancanze è elevato.

*Costo del cambiamento*
- Nel modello waterfall gli errori scoperti tardi costano molto di più.
- Il costo del cambiamento cresce in modo non lineare man mano che si avanza nelle fasi

**Modello incrementale**
Nel modello incrementale, il sistema non viene progettato e sviluppato “tutto insieme”, ma costruito per passi successivi. Ogni incremento introduce un insieme coerente di funzionalità che, una volta completate, producono una versione utilizzabile del sistema: 
1. Una base minima funzionante
2. Estensioni e miglioramenti progressivi
L’utente può iniziare a usare presto una parte del sistema e fornire feedback che guida gli
incrementi successivi.

*Vantaggi*
Con l’approccio incrementale il sistema cresce per piccoli rilasci, ciascuno utilizzabile. Dopo ogni rilascio gli utenti possono provare il sistema e dare feedback immediato. Questo permette di:
- correggere la direzione molto prima
- ridurre il rischio complessivo
- concentrare lo sforzo sulle funzionalità che portano più valore

*Limiti*
La flessibilità dell’approccio incrementale non significa assenza di pianificazione: serve comunque una visione architetturale complessiva del sistema. Se ogni incremento viene progettato come un mini–progetto indipendente, si rischia di accumulare debito tecnico, cioè il “costo futuro” dovuto a soluzioni veloci o provvisorie prese oggi. ad esempio:
- codice duplicato
- strutture provvisorie mai riscritte
- moduli con interfacce incoerenti
- mancanza di test
- architettura confusa.
Come un debito finanziario, il debito tecnico genera “interessi”: nel tempo tutto diventa più lento, più fragile e più costoso da modificare. Per funzionare bene, gli incrementi devono essere pensati come parti di un disegno coerente e non come patch aggiunte una dopo l’altra.

**Modello a spirale**
Il modello a spirale rappresenta lo sviluppo come una serie di cicli successivi, ognuno dei quali
affronta rischi e incertezze. Ogni giro della spirale include:
- definizione degli obiettivi
- analisi dei rischi
- sviluppo di una soluzione (spesso parziale)
- valutazione prima di passare al ciclo successivo.
L’idea è quella di accettare che non tutto sia chiaro all'inizio e che servano esperimenti per esplorare alternative. Per prototipo intendiamo una versione semplificata del sistema (o di una sua parte), costruita velocemente per capire se un’idea funziona, senza preoccuparsi ancora di tutti i dettagli. 
I prototipi vengono usati per ridurre rischi specifici:
- provare una nuova tecnologia
- testare un’interfaccia con gli utenti
- verificare ipotesi sulle prestazioni.
Ogni giro della spirale produce nuova conoscenza che permette di prendere decisioni più informate, evitando sorprese nelle fasi avanzate.

**Confronto dei tre modelli visti**
- Il waterfall è adatto a contesti stabili, con requisiti chiari e forte enfasi sulla documentazione.
- L’approccio incrementale funziona bene quando vogliamo rilasciare valore progressivo e ottenere feedback continui.
- La spirale è indicata quando rischio e incertezza sono alti e servono cicli di esplorazione controllati.
Nella pratica, molti progetti combinano elementi di questi modelli, adattandoli al contesto.

| Waterfall               | Incrementale         | Spirale                             |
| ----------------------- | -------------------- | ----------------------------------- |
| Sequenziale, fasi fisse | Rilasci progressivi  | Cicli iterativi guidati dal rischio |
| Requisiti stabili       | Feedback frequente   | Gestione esplicita dell’incertezza  |
| Documentazione forte    | Adattamento continuo | Prototipi per ridurre rischi        |
| Poco flessibile         | Più flessibile       | Molto flessibile ma più costoso     |
**Perché servono metodi agili**
I modelli di processo più tradizionali, come il waterfall e altri approcci sequenziali, nascono in contesti in cui i requisiti sono relativamente stabili e i rilasci poco frequenti. Oggi molti progetti vivono in ambienti altamente dinamici: 
- richieste che cambiano rapidamente
- mercati competitivi
- aspettative di aggiornamenti continui
I metodi agili nascono per affrontare proprio questa realtà, introducendo iterazioni brevi, feedback costanti e collaborazione continua.
 
**Principi chiave dell'approccio agile**
L’attenzione si sposta da piani rigidi a una capacità di adattamento, in cui il cambiamento non è un problema da evitare ma un dato di fatto da gestire. Il software funzionante diventa la misura principale di avanzamento del progetto, più importante della quantità di documenti prodotti. Vengono valorizzate:
- le persone e le interazioni
- la comunicazione diretta
- il lavoro di team
- la collaborazione continua con gli stakeholder

**Triangolo del progetto: tradizionale vs agile**
Nei progetti tradizionali si tende a fissare lo scopo (tutte le funzionalità), lasciando variabili tempo e costo.
Nei metodi agili si fissano tempo (durata dello sprint) e costo (dimensione del team), rendendo flessibile lo scopo in base al valore e alle priorità.

**Scrum**
Scrum è un framework organizzativo: definisce come un team collabora, prende decisioni e rilascia valore in modo iterativo. Osserveremo in particolare:
- i **ruoli essenziali** (Product Owner, Scrum Master, Development Team)
- il **backlog** come lista ordinata di funzionalità
- gli **sprint** come cicli brevi che producono incrementi funzionanti.

**Ruoli fondamentali in Scrum**
Il **Product Owner** rappresenta gli interessi di utenti e stakeholders: 
- definisce la visione del prodotto
- stabilisce le priorità nel backlog
Lo **Scrum Master** facilita il lavoro del team:
- aiuta a rimuovere gli ostacoli
- promuove le pratiche agili e la collaborazione
Il **Development Team** 
- realizza concretamente gli incrementi del sistema

**Backlog e user stories**
Il **backlog** è l’elenco ordinato delle funzionalità, dei miglioramenti e dei bug da gestire sul sistema. Molti elementi del backlog vengono descritti come **user stories**: brevi frasi dal punto di vista dell’utente, **che chiariscono chi vuole cosa e per quale motivo**. 
Il **backlog** è un artefatto vivo: viene aggiornato continuamente **sulla base dei feedback e dei risultati degli sprint.**

**Sprint e incrementi**
Lo sviluppo procede per sprint, cioè intervalli di tempo brevi e fissi (ad esempio 2 settimane), con obiettivi limitati ma chiari. Alla fine di ogni sprint il team dovrebbe poter mostrare un incremento funzionante del sistema, anche se piccolo. Questo ritmo cadenzato consente:
- di valutare spesso i progressi
- di adattare il piano sulla base dell’esperienza
- di coinvolgere regolarmente gli stakeholder
Le funzionalità non vengono affrontate tutte insieme, ma organizzate in:
- epic → macro–aree funzionali
- user stories → unità più piccole e verificabili
- sprint → passi concreti che producono un incremento funzionante.
Questo rende visibile l’evoluzione del sistema e permette di: 
- adattare le priorità nel tempo
- collegare ogni incremento a valore per gli utenti
A partire dalle epic possiamo ricavare user stories più piccole e concrete, ciascuna con un obiettivo chiaro e verificabile. Le user stories vengono scritte dal punto di vista dell’utente, specificando:
- chi è l’attore
- che cosa vuole ottenere
- perché quella funzionalità è utile.
---
### 26-03-04 lz. 3
***Requisito:***
<u>Un requisito è una descrizione chiara e condivisa di qualcosa che il sistema deve fare oppure di un vincolo che deve rispettare.</u> Possiamo vederlo come una promessa tra chi richiede il sistema e chi lo realizza. Un requisito ben scritto:
- riduce le ambiguità e le interpretazioni personali (devono essere chiari); 
- permette di verificare in modo oggettivo se il sistema si comporta come atteso (se ho dei vincoli questi vengono rispettati); 
- guida direttamente la progettazione dei test. 
- Criterio guida: un requisito deve essere verificabile (e, quando possibile, anche misurabile). 

> [!example] 
> - Il sistema deve permettere allo studente di iscriversi a un appello d’esame.
> - L’accesso alla piattaforma deve essere possibile solo tramite credenziali istituzionali. Va

Un requisito può essere espresso a due livelli di dettaglio:
1. *Lato utente:* linguaggio semplice, orientato ai bisogni, agli obiettivi e al contesto -> comprensibile a chi usa il sistema
2. *Lato sistema:* descrizione più precisa, utile per progettare e implementare concretamente -> utilizzato da sviluppatori e progettisti

Non sono due categorie distinte, ma di due modi di esprimere la stessa necessità. Questi due livelli sono spesso indicati come *requisiti utente* e *requisiti di sistema*.

> [!example] 
> Versione lato utente 
> - «Il sistema deve permettere al docente di ottenere un riepilogo degli appelli d’esame.» 
> 
> Versione lato sistema 
> - Il sistema genera un elenco contenente: corso, data, iscritti, esiti. 
> - Il riepilogo può essere filtrato per periodo, corso o docente. 
> - Il risultato è disponibile sia a video sia in formato scaricabile.

I requisiti si dividono in due grandi famiglie:
1. *Requisiti funzionali:* descrivono cosa il sistema deve fare, cioè servizi, operazioni e comportamenti osservabili dagli utenti (quali servizi deve offrire il sistema)
2. *Requisiti non funzionali*: descrivono come il sistema deve funzionare, quali vincoli deve rispettare e quale livello di qualità deve garantire (in che condizioni questi servizi devono funzionare bene)

<u>Queste due categorie si completano in quanto i funzionali definiscono la funzione mentre i non funzionali definiscono qualità, limiti e contesto operativo</u>

***Requisiti funzionali:***
Sono i requisiti più intuitivi, descrivono le azioni che il sistema deve permettere e le interazioni possibili tra utenti e piattaforma. Riguardano tutto ciò che è osservabile durante l’uso:
- inserire, modificare, cancellare informazioni; 
- consultare elenchi o dati specifici; 
- avviare procedure, confermare scelte, scaricare documenti. 

I requisiti funzionali definiscono la trama del sistema: se mancasse anche solo una funzione essenziale, il sistema non potrebbe svolgere la sua attività principale. 

> [!example] 
> - Lo studente può iscriversi a un appello d’esame.
> - La segreteria può creare un nuovo corso e associarvi un docente.
> - Il docente può registrare i voti degli studenti iscritti.»

<u>Un requisito funzionale vago è una fonte di errore</u>. Un requisito vago porta inevitabilmente a: 
- interpretazioni diverse
- risultati incoerenti
- bisogno di riscrivere e correggere in corso d’opera

Un requisito formulato in modo generico non è verificabile: non indica con chiarezza quale risultato ci si aspetta e quindi non guida la definizione dei test. Ambiguità di questo tipo tendono ad accumularsi e generano debito tecnico: prima o poi sarà necessario rivedere codice e documentazione per riallinearli a una definizione più precisa

***Requisiti non funzionali:***
Non descrivono azioni specifiche del sistema, ma come il sistema deve comportarsi, quali vincoli deve rispettare e quale livello di qualità deve garantire. Questi riguardano aspetti trasversali che influenzano l’esperienza d’uso e l’affidabilità complessiva, come:
- prestazioni: tempi di risposta, capacità di gestire molti utenti contemporanei; 
- sicurezza: protezione dei dati, autenticazione, autorizzazioni;
- affidabilità: continuità del servizio, gestione degli errori, recupero dopo un guasto;
- accessibilità e usabilità: chiarezza dell’interfaccia, facilità d’uso, barriere ridotte;
- compatibilità: integrazione con sistemi già esistenti;
- rispetto di norme e regolamenti: ad esempio le regole sulla gestione dei dati personali. 

Questi requisiti incidono profondamente sulla riuscita del progetto: un sistema che “fa tutto” ma è lento, complesso da usare o poco sicuro risulta comunque inadeguato. Possiamo raggruppare i requisiti non funzionali in 3 famiglie (per non trascurare aspetti importanti ma meno visibili rispetto alle funzionalità):
1. *Legati al prodotto* = qualità interne del sistema (prestazioni, tempi di risposta, affidabilità, consumo di risorse). 
2. *Legati all'organizzazione* = vincoli del contesto in cui il sistema è sviluppato o usato (strumenti obbligatori, procedure interne, integrazioni con sistemi esistenti). 
3. *Legati all'ambiente esterno* =norme, regolamenti e standard tecnici che il sistema deve rispettare.

> [!example] 
> Requisiti non funzionali legati al prodotto 
> - Le operazioni di iscrizione a un appello devono avere un tempo di risposta inferiore a 2 secondi nella maggior parte dei casi.
> - Il sistema deve poter gestire accessi contemporanei da parte di migliaia di studenti durante i periodi di picco.
> 
> Requisiti legati all'organizzazione 
> - L’accesso alla piattaforma avviene esclusivamente tramite le credenziali istituzionali dell’ateneo.
> - Il sistema deve integrarsi con i servizi già presenti: anagrafe studenti, archivio docenti, gestione dei corsi.
> 
> Requisiti legati all'ambiente esterno 
> - Il trattamento dei dati degli studenti deve rispettare la normativa vigente sulla protezione dei dati personali.
> - La piattaforma deve garantire la conservazione dei verbali secondo le disposizioni ministeriali.

***Requisiti verificabili:***
<u>Un requisito è utile solo se possiamo verificarlo in modo oggettivo</u>. Spesso nelle fasi iniziali emergono affermazioni molto generali: «facile da usare», «veloce», «sicuro», «intuitivo». Questi obiettivi sono utili per capire la direzione desiderata, ma non dicono concretamente come il sistema deve comportarsi. In questa forma non sono verificabili: due persone possono interpretarli in modo anche molto diverso. Occorre quindi un passaggio successivo: trasformare gli obiettivi in requisiti misurabili, cioè descrizioni che permettono di verificare in modo oggettivo se il sistema rispetta ciò che promette.

> [!example]  
> Da: «Il sistema deve essere veloce.» 
> 
>  A: «Il sistema deve rispondere entro 2 secondi per almeno il 95% delle richieste di iscrizione agli appelli.»

Perché questa trasformazione è importante? 
- Permette di verificare il comportamento del sistema in modo oggettivo. 
- Evita interpretazioni diverse tra utenti, progettisti e sviluppatori. 
- Guida scelte architetturali e tecniche coerenti.

È del tutto normale che, nelle prime fasi, i requisiti non siano perfetti: alcuni mancano semplicemente perché nessuno ci ha ancora pensato, altri risultano in conflitto (ad esempio: «facile da usare» ma anche «richiedere password molto complesse»), altri ancora riflettono punti di vista diversi tra studenti, docenti, segreteria, tecnici. Questa situazione non è un difetto del progetto, ma una caratteristica naturale dei sistemi reali: le esigenze si scoprono poco alla volta e spesso emergono contraddizioni da risolvere. • L’analisi dei requisiti serve proprio a: 
- chiarire cosa si intende davvero
- verificare che le richieste siano realistiche e compatibili
- negoziare priorità e compromessi tra stakeholder con bisogni diversi

I requisiti non emergono spontaneamente dal codice o dalla tecnologia: occorre andare a cercarli parlando con chi userà davvero il sistema. Per raccoglierli si usano strumenti diversi, spesso complementari:
- colloqui con gli utenti per capire bisogni, difficoltà e aspettative; 
- osservazione delle attività reali per vedere come i processi funzionano oggi, senza filtri; 
- analisi di procedure, moduli e documenti esistenti per cogliere vincoli già presenti nell'organizzazione; 
- prototipi e dimostrazioni preliminari per far emergere idee e problemi che non verrebbero fuori solo parlando. 

Raccogliere requisiti è un’attività progressiva: solo con un confronto continuo con gli stakeholder si costruisce una visione chiara e condivisa del sistema. Una volta raccolti i bisogni, serve però un modo semplice e comprensibile per rappresentarli: da qui nasce l’uso delle user stories.

***User stories:***
Una user story descrive un’esigenza dal punto di vista dell’utente, in modo semplice e comprensibile anche a chi non è tecnico. Risponde a tre domande fondamentali:
- chi usa il sistema
- che cosa vuole ottenere
- perché quell'azione è importante. 

La forma tipica è molto essenziale:  «Come (tipo di utente), voglio (azione) per (obiettivo/beneficio).». Le user stories non descrivono i dettagli tecnici, ma chiariscono l’intenzione: servono ad allineare utenti e sviluppatori e a popolare un backlog comprensibile e prioritizzabile. Per essere davvero utili, dovrebbero:
- avere un chiaro valore per l’utente
- essere abbastanza piccole da stare in uno sprint
- risultare verificabili tramite criteri di accettazione espliciti.

> [!info]
> - **User story**: frase breve orientata a chi-cosa-perché, utile per raccogliere e prioritizzare bisogni nel backlog.
> - **Use case**: descrizione più strutturata dell’interazione, che useremo più avanti nella parte di modellazione UML.

***Criteri di accettazione:***
Ogni user story deve essere accompagnata da uno o più criteri di accettazione, cioè condizioni che stabiliscono con chiarezza quando la funzionalità può essere considerata davvero completata. I criteri di accettazione servono a: 
- eliminare ambiguità su cosa ci si aspetta dal sistema; 
- definire un confine preciso tra “fatto” e “non ancora fatto”; 
- costituire una base solida per i test che verificheranno il corretto comportamento del sistema. 

<u>In altre parole: la user story racconta l’intenzione, i criteri di accettazione ne specificano la verifica concreta</u>. Alcuni criteri possono essere espressi anche in forma strutturata, ad esempio con frasi del tipo «Dato che . . . Quando . . . Allora . . . », facilitando la verifica sistematica del comportamento del sistema.

> [!example]
> User story 
> - «Come studente, voglio iscrivermi a un appello d’esame per poter sostenere la prova nella data che scelgo.» 
>
> Criteri di accettazione 
> - Lo studente vede solo gli appelli disponibili per il proprio corso. 
> - Se i posti sono esauriti, il sistema mostra un messaggio chiaro e comprensibile.
> - Dopo l’iscrizione, lo studente visualizza immediatamente la conferma nella propria area personale. 
> - Il docente può consultare l’elenco aggiornato degli iscritti all'appello.

***Backlog:***
Le user stories raccolte diventano elementi del backlog che guida lo sviluppo: una raccolta ordinata di ciò che il sistema dovrà offrire. Non tutti gli elementi hanno lo stesso peso: occorre stabilire delle priorità. Alcune funzionalità vanno affrontate prima perché: 
- portano più valore agli utenti; 
- comportano rischi tecnici o organizzativi da affrontare presto;
- dipendono da altre funzionalità che devono essere pronte. 

Il backlog è un artefatto vivo: cambia, si aggiorna, viene riorganizzata man mano che emergono nuove informazioni o nuove esigenze.
I requisiti non funzionali non devono essere trascurati: anche loro entrano a pieno titolo nel backlog. Alcuni diventano attività specifiche, da pianificare come vere e proprie funzionalità: 
- migliorare i tempi di risposta nelle operazioni più frequenti; 
- introdurre o potenziare un meccanismo di sicurezza; 
- rendere l’interfaccia più accessibile o più semplice da usare. 

Altri, invece, sono vincoli trasversali: non compaiono come singole attività, ma devono essere rispettati in ogni funzionalità sviluppata (ad esempio sicurezza, normativa sui dati, compatibilità con sistemi esistenti). Ignorarli nelle prime fasi porta quasi sempre a interventi più costosi, più complessi e più rischiosi nelle fasi successive. 

---
### 26-03-05 lz. 4
<u>**Dominio** = Per dominio si intende l'insieme di: persone, oggetti, regole e processi reali che il software deve rappresentare e supportare.</u>

***Modellazione concettuale:***
La modellazione concettuale è l'attività che mira a:
- identificare i concetti importanti del dominio; 
- chiarire quali responsabilità ha ciascun concetto; 
- descrivere le relazioni tra questi concetti.

È una descrizione logica del problema, indipendente dalla soluzione tecnica. Nella pratica: stiamo cercando di costruire un vocabolario condiviso che tutti possano capire (utenti, analisti, sviluppatori).

**Concetti:**
Per concetti intendiamo gli elementi significativi del dominio (es. Studente, Corso, Appello), non ancora classi, tabelle o dettagli implementativi. I requisiti che abbiamo scritto non sono solo “liste di funzionalità”: contengono già molti suggerimenti sul mondo che il sistema deve rappresentare. Leggendo frasi come:  «Lo studente si iscrive a un appello d’esame.», «Il docente registra il voto.», «Lo studente visualizza la propria carriera.» possiamo chiederci: quali entità compaiono? quali azioni vengono svolte? quali relazioni tra concetti sono implicite? La modellazione concettuale parte spesso proprio da questo lavoro di “lettura attenta” dei requisiti.

Dal punto di vista concettuale, nei requisiti ci interessano in particolare: 
- i sostantivi (studente, appello, corso, voto, verbale, carriera...) che possono suggerire entità del dominio; 
- i verbi (iscriversi, registrare, consultare, creare, chiudere...) che possono suggerire azioni o relazioni; 
- le espressioni di vincolo (solo se, fino a, entro...) che possono nascondere regole del dominio. 

Non è una trasformazione meccanica, ma una guida: ci aiuta a individuare i concetti davvero importanti.

**Entità concettuali:**
Un’entità concettuale è un elemento del mondo che vogliamo rappresentare in modo esplicito nel sistema:
- esiste indipendentemente dalla tecnologia che useremo; 
- ha un significato riconoscibile anche da chi non è tecnico. 

> [!example] Esempi generici
> - persona, ordine, prodotto, pagamento, prenotazione
> - studente, corso, esame, voto, appello

L’errore da evitare è confondere subito entità concettuali con strutture di codice.

> [!tip] Ricorda:
> Nel nostro caso di studio, una prima lista di entità concettuali naturali è: 
> - Studente = la persona che segue i corsi e sostiene esami
> - Corso = rappresenta l'insegnamento (nome, docente, CFU,...)
> - Appello = occasione concreta per sostenere l'esame
> - Iscrizione = atto con cui lo studente si registra ad uno specifico appello
> - Voto = esito ottenuto in quell'appello
> 
> In questo modo posso rappresentare situazioni reali come:
> - uno studente si può iscrivere più volte
> - può sostenere lo stesso corso in appelli diversi
> - un appello può avere molti iscritti e molti voti
> 
> L’idea non è ancora “definire le classi”, ma capire: 
> - che cosa rappresenta ciascuna di queste entità; 
> - perché è utile distinguerle; 
> - che tipo di informazioni ci aspettiamo che custodiscano.

A ogni entità possiamo associare alcune responsabilità (dominio delle operazioni consentite ad un'entità), cioè:
- quali informazioni custodisce
- in quali interazioni è coinvolta
- quali decisioni del dominio “passano da lei”. 

Definire bene le responsabilità riduce il rischio di “entità onnipotenti” o confuse.

> [!example] 
> - lo Studente è responsabile dei propri dati anagrafici e della propria carriera
> - l'appello è responsabile di data, ora, aula, lista potenziale degli iscritti
> - iscrizione è responsabile del legame tra uno specifico studente e uno specifico appello

Le entità non vivono isolate: ciò che rende interessante un dominio sono spesso le relazioni tra i concetti. In questa fase distinguiamo due famiglie utili: 
- relazioni che descrivono la struttura del dominio; 
- relazioni che descrivono le azioni che avvengono nel dominio. 

Capire bene le relazioni è essenziale per rispondere a domande tipiche: 
- «Quali appelli può vedere questo studente?» 
- «Quali studenti si sono iscritti a questo appello?» 
- «Quali voti ha ottenuto uno studente in un certo arco di tempo?»

Alcune relazioni descrivono la struttura stabile del dominio: 
- uno studente ha una carriera composta da esami e voti; 
- un corso è tenuto da uno o più docenti; 
- un appello è relativo a un certo corso e a una certa data.

> [!example]
> - Corso – Docente: un corso è tenuto da uno o più docenti. 
> - Corso – Appello: un corso ha uno o più appelli nel tempo. 
> - Studente – Iscrizione: uno studente può avere più iscrizioni. 
> - Iscrizione – Appello: ogni iscrizione si riferisce a uno specifico appello. 

Queste relazioni esistono “a prescindere” dal singolo utilizzo del sistema: descrivono il modo in cui è organizzato il mondo universitario. Sono la base della nostra mappa concettuale del dominio.

Altre relazioni descrivono azioni/eventi che avvengono nel dominio: 
- Studente → Iscrizione: crea una nuova iscrizione a un appello; 
- Docente → Voto: registra o aggiorna l’esito di uno studente; 
- Segreteria → Corso/Appello: crea o aggiorna i dati organizzativi. 

Queste relazioni sono più vicine al comportamento del sistema: hanno tipicamente precondizioni, effetti e controlli.

**Regole del dominio:**
Le regole del dominio esprimono vincoli che devono valere sempre, indipendentemente dalla tecnologia. 

> [!example] in UniManager
> - uno studente può iscriversi solo ad appelli del proprio corso; 
> - non ci si può iscrivere dopo la chiusura delle iscrizioni; 
> - il voto può essere registrato solo per studenti regolarmente iscritti all’appello; 
> - la verbalizzazione chiude la modifica libera dei voti. 

Queste regole guidano sia la mappa concettuale sia i controlli che progetteremo nelle fasi successive.

> [!tip] Caratteristiche iscrizione
> - una data in cui è stata compiuta; 
> - un eventuale stato (attiva, ritirata, annullata ...); 
> - può essere usata per verificare la regolarità di uno studente in un appello.

> [!tip] Entità voto
> - è legato a uno studente e a un appello specifico; 
> - può avere uno stato (proposto, verbalizzato, rifiutato ...); 
> - può essere oggetto di visualizzazione, controllo, statistiche.

> [!info] Confine del sistema
> Oltre alle entità del dominio, è importante chiarire chi interagisce con il sistema: studenti, docenti, segreteria, eventuali amministratori di sistema. Dobbiamo chiederci: che cosa sta dentro UniManager e che cosa sta fuori? • 
>> [!example] Elementi esterni:
>> - sistemi di autenticazione dell’ateneo; 
>> - basi dati istituzionali di anagrafe
>> - altri sistemi amministrativi
> 

**Rappresentazione contesto in modo informale:**
Prima ancora di usare notazioni formali, possiamo disegnare schemi semplici: 
- un riquadro per il sistema UniManager; 
- attorno, gli attori principali (studenti, docenti, segreteria); 
- frecce che indicano i tipi di interazione (iscrizione, consultazione, gestione corsi ...). 

Questi schemi servono a: 
- chiarire chi parla con chi; 
- evitare di attribuire al sistema responsabilità che non ha; 
- distinguere ciò che dobbiamo progettare da ciò che esiste già.

Prima della mappa concettuale è utile esplicitare alcune azioni tipiche del sistema, ad esempio: 
- lo studente trova gli appelli disponibili; 
- lo studente si iscrive; 
- il docente registra i voti; 
- lo studente consulta la propria carriera. 

Possiamo rappresentarle in modo semplice, ad esempio con: 
- una sequenza di passi, 
- frecce tra attori e sistema, 
- brevi descrizioni delle interazioni principali.

**Rappresentare la struttura concettuale:**
Dopo la vista delle azioni, passiamo alla vista strutturale del dominio con una mappa concettuale: 
- nodi per le entità (Studente, Corso, Appello, Iscrizione, Voto ...);
- collegamenti che indicano le relazioni principali (ad esempio “Studente è iscritto a Appello tramite Iscrizione”). 

Usiamo una rappresentazione informale: 
- niente simboli complicati, 
- etichette in italiano, 
- solo gli elementi davvero utili per capire il dominio. 

Questa mappa sarà la base del futuro modello del dominio in notazione più formale.

Come prima approssimazione, possiamo disegnare: 
- Studente, Corso, Appello come entità principali; 
- una relazione che lega Corso e Appello; 
- collegamenti tra Studente e Appello attraverso l’idea di iscrizione

Questa versione iniziale è volutamente semplice: 
- ci permette di parlare, discutere, correggere; 
- mette in evidenza eventuali concetti mancanti

Una volta condivisa la prima versione, possiamo arricchirla con: 
- l’entità Iscrizione come collegamento esplicito; 
- l’entità Voto come esito dell’iscrizione; 
- eventuali altri concetti emersi (sessione, verbale, carriera ...).

La mappa concettuale è uno strumento vivo: 
- può cambiare quando scopriamo nuovi bisogni; 
- può essere semplificata se ci accorgiamo di avere concetti inutilmente complicati.

Passiamo ora da una mappa minima a una mappa più ricca, aggiungendo tre elementi: 
- Iscrizione come collegamento esplicito tra Studente e Appello; 
- Voto come esito associato all'iscrizione; 
- alcuni concetti opzionali (sessione, verbale, carriera) per mostrare come il modello può evolvere. 

Obiettivo della versione estesa: 
- aumentare chiarezza e capacità esplicativa del modello
- senza introdurre ancora dettagli implementativi.

Una rappresentazione più completa del dominio, dove compaiono: 
- l’entità Iscrizione come collegamento esplicito tra Studente e Appello; 
- l’entità Voto come esito dell’iscrizione; 
- alcuni concetti aggiuntivi possibili (sessione, verbale, carriera) per mostrare come la mappa possa evolvere

**Checklist rapido:**
Prima di proseguire, facciamo un controllo rapido: 
- Nomi chiari: ogni entità ha un nome comprensibile anche a chi non è tecnico? 
- Relazioni esplicite: i collegamenti principali sono scritti e coerenti con i requisiti? 
- Niente dettagli tecnici: stiamo evitando classi, tabelle DB, API e scelte implementative? 

Se una risposta è “no”, conviene correggere subito: nelle fasi successive gli errori costano di più.

Una buona modellazione concettuale: 
- riduce le ambiguità, 
- aiuta a prendere decisioni progettuali più consapevoli, 
- facilita il dialogo tra persone con competenze diverse. 

Al contrario, se si salta questa fase: 
- si rischiano fraintendimenti tra utenti e sviluppatori; 
- si introducono vincoli tecnici non necessari; 
- diventa più difficile correggere gli errori in seguito. 

Investire tempo qui fa risparmiare tempo nelle fasi successive.

**Errori tipici:**
Alcuni errori ricorrenti: 
- confondere subito entità concettuali e classi di implementazione; 
- introdurre dettagli tecnici (database, interfacce grafiche...) troppo presto; 
- creare entità “che fanno tutto” e relazioni poco chiare; 
- usare nomi vaghi o poco significativi (dati1, oggetto, elemento...).

È importante ricordare che in questa fase lo scopo non è “scrivere codice”, ma capire bene il problema.

---
### 26-03-12 lz. 5
Iniziamo oggi ad usare una notazione standard per rappresentare i modelli: UML (unified modelling language) che è un linguaggio visivo condiviso

UML è un linguaggio di modellazione standard, non un linguaggio di programmazione. Fornisce un insieme di diagrammi per descrivere:
- la struttura statica del sistema
- il comportamento dinamico
- le interazioni tra parti del sistema e attori esterni. 
Non obbliga a un particolare processo (waterfall, agile,...), ma può essere usato in tutti. 

UML nasce negli anni ’90 per unificare diverse notazioni di analisi e progettazione. I principali contributori iniziali sono Grady Booch, Ivar Jacobson e James Rumbaugh. La prima standardizzazione avviene in ambito OMG (Object Management Group). Nel tempo UML ha avuto più versioni: la famiglia 1.x e poi 2.x, con estensioni e chiarimenti della notazione. Oggi UML resta un riferimento diffuso per comunicare modelli in modo condiviso tra team diversi.

UML mette a disposizione diversi tipi di diagrammi. In questo corso vedremo soprattutto: 
- Use Case Diagram — vista del sistema dal punto di vista degli attori esterni; 
- Modello di dominio — modello concettuale del dominio, spesso rappresentato con un class diagram semplificato; 
- Class Diagram — struttura statica del dominio (classi, attributi, associazioni); 
- Sequence Diagram — interazioni nel tempo tra oggetti/attori; 
- Activity Diagram — flusso di attività, decisioni e responsabilità di processo. 
UML include anche altri diagrammi (ad esempio State Machine, Component, Deployment). 

Inizieremo dall'Use Case Diagram, perché: 
- è vicino ai requisiti e alle user stories, 
- riguarda chi usa il sistema e per quali obiettivi.

> [!info] I casi d’uso descrivono gli obiettivi significativi per gli utenti.

Abbiamo già ragionato su: 
- chi interagisce con un sistema software; 
- che cosa sta dentro e che cosa sta fuori dal sistema. 
Nel Use Case Diagram UML: 
- gli attori sono ruoli esterni; 
- il sistema è rappresentato da un rettangolo che racchiude i casi d’uso; 
- le interazioni diventano associazioni.

**Attore**
Un attore è un ruolo esterno che interagisce con il sistema:
- una persona (studente, docente), 
- un altro sistema (sistemi di autenticazione, basi dati esterne), 
- un dispositivo. 
È importante ricordare che: 
- non rappresenta una singola persona, ma un tipo di utilizzatore; 
- è sempre fuori dal sistema. 
In UniManager i nostri attori principali sono: Studente, Docente, Segreteria.

**Caso d'uso**
Un caso d’uso rappresenta un obiettivo significativo per un attore non è un singolo “click” ne una funzione tecnica (“query sul DB”), ma un risultato che ha senso dal punto di vista dell’utente. 

> [!example] Esempi tipici in UniManager: 
> -  Iscrizione all’appello, 
> - Registrazione voti, 
> - Consultazione libretto, 
> - Creazione appello

Nella notazione UML, i casi d’uso sono rappresentati come ellissi con un nome. Nei diagrammi completi i casi d’uso stanno dentro il confine del sistema; negli esempi minimali di queste slide il rettangolo può essere omesso per semplicità visiva.

Nella lezione precedente abbiamo letto frasi come: 
- «Lo studente si iscrive a un appello d’esame.» 
- «Il docente registra il voto.» 
- «Lo studente visualizza la propria carriera.» 
Dal punto di vista dei casi d’uso: 
- attore: chi compie l’azione (Studente, Docente,...); 
- obiettivo: che cosa vuole ottenere (essere iscritto, vedere la carriera, registrare il voto); 
- sistema: che parte del lavoro è di UniManager. 
Il Use Case Diagram sintetizza queste relazioni in forma visuale.
Dai requisiti e dalle user stories emergono quattro casi d’uso fondamentali: 
- Iscrizione all’appello; 
- Creazione appello; 
- Registrazione voti; 
- Consultazione libretto. 
Questi casi d’uso coinvolgono i tre attori principali del sistema: Studente, Docente, Segreteria. Sono il punto di partenza per costruire il primo Use Case Diagram di UniManager.

> [!tip] ![[Screenshot From 2026-03-12 15-15-14.png|400]]

**Relazione *include* (obbligatorio)**
Include rappresenta un comportamento comune e obbligatorio tra più casi d’uso. Un caso d’uso A include B se: 
- ogni volta che si esegue A, a un certo punto è necessario eseguire B; 
- B è un “pezzo riusabile” usato da più casi d’uso. 
A livello grafico: 
- la relazione è una freccia tratteggiata con l’etichetta «include»; 
- la freccia punta verso il caso d’uso incluso (da A verso B). 
Serve quando una parte del comportamento è identica e va eseguita sempre prima o durante un caso d’uso principale.

> [!example] Esempio uso include
>  ![[Screenshot From 2026-03-12 15-18-15.png|400]]

**Relazione *extend* (opzionale)**
Extend rappresenta una variante opzionale di un caso d’uso. Un caso d’uso B estende A se: 
- A descrive il flusso “normale”; 
- B aggiunge passi opzionali o condizioni particolari. 
A livello grafico: 
- la relazione è una freccia tratteggiata con etichetta «extend»; 
- la freccia punta verso il caso d’uso esteso (da B verso A). 
- Serve quando un comportamento aggiuntivo può avvenire, ma non è obbligatorio.

> [!example] Esempio uso extend
> ![[Screenshot From 2026-03-12 15-21-24.png|400]]

**Generalizzazione tra attori**
È possibile organizzare gli attori in una gerarchia, serve quando più attori condividono gli stessi casi d’uso. Significato: 
- l’attore “padre” rappresenta un ruolo più generale; 
- l’attore “figlio” eredita le sue interazioni. 
A livello grafico (UML): 
- si usa una linea continua con una freccia triangolare non piena; 
- la freccia punta verso l’attore più generale (da Specializzato → Generale). 
Va usata solo se semplifica il modello o riduce ripetizioni. • In alcuni casi può essere utile introdurre un attore più generale: Utente come ruolo comune. Attori più specifici: 
- Studente, 
- Docente. 
La generalizzazione permette di evitare ripetizioni: 
- i casi d’uso comuni a tutti gli utenti possono collegarsi a Utente; 
- i casi specifici rimangono ai ruoli specializzati.

**Note e commenti**
In UML possiamo aggiungere note/commenti per chiarire assunzioni e scelte. Le note aiutano a esplicitare dettagli utili senza appesantire i casi d’uso. Sono particolarmente utili quando una regola non è rappresentabile in modo diretto con soli attori e relazioni. Una buona pratica è quella di scrivere note brevi e mirate, collegate solo ai punti davvero ambigui.

> [!example] ![[Screenshot From 2026-03-12 15-27-37.png|400]]

**Package e raggruppamenti**
In UML possiamo raggruppare casi d’uso in package/aree funzionali. Ciò è utile quando il diagramma cresce e vogliamo separare aree funzionali. La motivazione principale è la leggibilità: meno incroci visivi, più chiarezza su responsabilità e perimetri. I package non cambiano il comportamento dei casi d’uso: organizzano la vista del modello.

> [!example] ![[Screenshot From 2026-03-12 15-27-47.png|400]]

**Molteplicità sulle relazioni**
In alcuni casi possiamo esplicitare la molteplicità sulle relazioni (senza esagerare per rimanere il più asciutti possibile, altrimenti si rischia solo rumore). Va usata solo quando aggiunge chiarezza al significato del modello. Nei Use Case Diagram è una notazione possibile ma non centrale (più comune nei Class Diagram). 

> [!example]
> - Esempio 1 (1–*): uno Studente può effettuare molte iscrizioni nel tempo. 
> - Esempio 2 (1–0..*): un Docente può registrare zero o più voti. 

**Vincoli**
In UML possiamo esprimere vincoli con testo tra parentesi graffe. I vincoli servono per esplicitare regole di business non evidenti dal solo diagramma (riguardo l'implementazione). Sono utili quando una condizione è decisiva per il dominio (scadenze, prerequisiti, autorizzazioni). Un vincolo deve essere verificabile: formulazioni vaghe riducono il valore del modello.

> [!example] ![[Screenshot From 2026-03-12 15-36-05.png|400]]
![[Screenshot From 2026-03-12 15-36-23.png|400]]

**Good practices**
- Usare i casi d’uso per descrivere obiettivi, non schermate o click. 
- Mantenere un livello di dettaglio gestibile: pochi casi d’uso macro, non decine di micro-casi. 
- Evitare di usare include/extend in modo eccessivo o artificiale. 
- Annotare, quando serve, qualche breve descrizione testuale per chiarire le intenzioni. 
- Ricordare che il Use Case Diagram è uno strumento di comunicazione con gli stakeholder: deve essere leggibile anche da chi non è tecnico.

**Descrizione testuale**
Il diagramma dei casi d’uso dà una vista sintetica, ma spesso non basta. Per i casi d’uso più importanti conviene scrivere una scheda testuale con: 
- Nome del caso d’uso; 
- Attori coinvolti; 
- Breve descrizione (obiettivo principale); 
- Precondizioni (cosa deve essere vero prima); 
- Postcondizioni (cosa deve essere vero dopo); 
- Flusso principale (passi “normali”); 
- eventuali flussi alternativi / eccezioni. 
L’idea è avere un livello di dettaglio sufficiente per: 
- chiarire il comportamento atteso, 
- preparare i diagrammi successivi (sequence, activity. . . ).

> [!example] ![[Screenshot From 2026-03-12 15-40-44.png|400]]

---
