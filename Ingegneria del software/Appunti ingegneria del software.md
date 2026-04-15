###### Introduzione alla SE moderna
**Ingegneria del software oggi:** Quando si parla di software, si parla della struttura portante dei sistemi moderni
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

**Quando il software fallisce:** I problemi, oltre ad essere "bug fastidiosi", creano anche dei problemi reali
- sistemi di prenotazioni in down durante i picchi di accesso
- piattaforme di pagamento non disponibili
- errori nei sistemi di controllo 
- perdita o esposizione dei dati sensibili degli utenti

E in questi casi la risoluzione non consiste solo nella correzione di una riga ma può prevedere anche la ri-progettazione dell'intero sistema. L'**Ingegneria del software** nasce per ridurre la probabilità che questi scenari si verifichino.

**Progettare sistemi complessi:** Nei sistemi reali entrano in gioco diversi fattori
- molti operatori coinvolti
- vincoli tecnici e organizzativi
- grandi moli di dati
- interazione tra componenti eterogenei
- requisiti che evolvono
- tecnologia che cambia
- contesti d'uso diversificati

L’Ingegneria del Software nasce per gestire questa complessità e fornire strumenti, metodi e pratiche per affrontare lo sviluppo in modo sistematico.

**Programmazione $\neq$ ingegneria del software:** La scrittura del codice è solo una delle attività dello sviluppo software e rappresenta una porzione limitata del lavoro complessivo.
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

**Che cos'è l'ingegneria del software:** L’Ingegneria del Software è l’insieme di principi, metodi e strumenti per sviluppare e mantenere sistemi software affidabili e utili in modo economicamente sostenibile.
L'attività di sviluppo si suddivide principalmente in 4 processi:
- **Specifica:** capire che cosa il software deve fare e quali vincoli deve rispettare.
- **Sviluppo:** progettare l’architettura e le strutture interne, scrivere il codice.
- **Validazione:** verificare che il sistema soddisfi davvero i requisiti attesi.
- **Evoluzione:** adattare il software a nuove esigenze, correzioni, cambiamenti di contesto.

**Il costo reale del software è nella manutenzione**
Nei progetti reali, la fase  iniziale di sviluppo rappresenta solo una parte del costo complessivo. Ma è la **manutenzione** che assorbe la maggior parte delle risorse
- correzione dei bug
- aggiornamenti richiesti
- adattamenti a normative e tecnologie
- estensioni e miglioramenti continui

La sostenibilità del software dipende dalla qualità delle scelte fatte nelle fasi iniziali. Un buon progetto facilita l’evoluzione, mentre un progetto fragile la rende costosa e rischiosa.

**Ciclo di vita:** Il **ciclo di vita** fornisce un quadro concettuale per organizzare le attività
- raccolta dei requisiti
- progettazione soluzione
- implementazione
- verifica
- rilascio
- manutenzione

1. Le attività iniziali seguono un flusso lineare: requisiti $\rightarrow$ progettazione $\rightarrow$ implementazione.
2. Dopo la prima versione, il processo entra in un ciclo continuo di test, rilascio e manutenzione.
3. Le attività della seconda riga sono iterative: ogni modifica o bug-fix riattiva il ciclo.
4. La manutenzione genera spesso nuovi requisiti, chiudendo il loop verso la fase iniziale.

**Comprensione del problema:** La costruzione del software deve partire da una reale **comprensione** del problema da risolvere e delle esigenze degli utenti.
- Chi userà il sistema? 
- Con quali obiettivi? 
- In quali contesti? 
- Con quali vincoli?

Un’analisi affrettata porta a sviluppare software poco utile o costoso da correggere; dedicare tempo alla comprensione iniziale evita errori di fondo.

**Nessun progetto procede in linea retta:** Durante lo sviluppo emergono inevitabilmente nuove informazioni e nuove esigenze, che mettono in discussione decisioni prese in precedenza. Il processo reale è **iterativo**
- si torna sui requisiti
- si rivede il design
- si modifica il codice
- si aggiornano i test.

L’ingegneria del software non elimina queste iterazioni, ma fornisce strumenti per gestirle in modo controllato e produttivo.

**La qualità non è un optional:** Un sistema software non deve solo funzionare: deve funzionare bene, con:
- affidabilità
- efficienza
- semplicità d’uso 
- capacità di evolvere nel tempo

La qualità si manifesta su due livelli:
- **Front-end:** esperienza d'uso, chiarezza, prestazioni
- **Back-end:** architettura, struttura del codice, dipendenze, documentazione

Se la qualità interna è scarsa, il sistema diventa fragile e difficile da modificare, indipendentemente da come appare all’ esterno.

**Le figure presenti nello sviluppo di un progetto:** Un sistema software è il risultato del lavoro di molte figure diverse con ruoli differenti e competenze specifiche.
- **utenti e committenti:**  definiscono bisogni e obiettivi
- **analisti e progettisti:** traducono le esigenze in modelli e architetture
- **sviluppatori e tester:** costruiscono e verificano il sistema
- **manutentori:** garantiscono la continuità del servizio nel tempo

La collaborazione tra questi ruoli è un fattore determinante per il successo del progetto.
Ogni ruolo ha una prospettiva diversa sul sistema e l’ingegneria del software deve farle convivere in modo coerente. 
- **l’utente** si concentra sull'esperienza d’uso e sulla semplicità delle operazioni.
- **lo sviluppatore** guarda alla correttezza e alla struttura del codice.
- **l’architetto** pensa alla coerenza strutturale e alle dipendenze tra moduli;
- **il manutentore** si preoccupa della comprensibilità interna e della facilità di modifica.

Integrare queste prospettive è fondamentale per ottenere un sistema davvero efficace.

**Quando la comunicazione fallisce, fallisce il progetto:** Molti difetti e ritardi non derivano da problemi tecnici, ma da incomprensioni tra le persone coinvolte. 
- requisiti fraintesi o ambigui
- aspettative non esplicitate
- decisioni progettuali non documentate
- informazioni critiche che non vengono condivise

Una Comunicazione efficace è **parte integrante dell’ingegneria del software**, tanto quanto il codice.
 
 **Il rischio è parte naturale di ogni progetto:** Ogni decisione progettuale e tecnologica introduce incertezza. 
- scelte tecnologiche premature o non verificate.
- requisiti instabili o in continuo cambiamento.
- integrazioni difficili o vincoli non previsti.
- stime errate di tempi e complessità.

L'ingegneria del software non elimina i rischi ma li gestisce in modo da correggerlo in tempo.

**La complessità è il vero avversario del software:** La difficoltà principale nello sviluppo del software non risiede nelle singole funzionalità, ma nell'interazione tra molte parti diverse.
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

**Obiettivo dell' ingegnere del software:** L'obbiettivo dell'ingegneria del software si pone come obbiettivo finale in un progetto la **riduzione dell'incertezza** durante lo sviluppo
- prevedere i tempi e i costi necessari
- valutare diverse alternative prima dell'implementazione
- rilevare problemi in anticipo
- mantenere il controllo anche in presenza di cambiamenti

L'approccio ingegneristico non elimina gli imprevisti, ma rende possibile affrontarli con metodo senza perdere stabilità.

###### Processi software
Un' attività collettiva che coinvolge persone, ruoli e decisioni distribuite nel tempo. Quando si lavora con più persone, ognuna tende a seguire il proprio stile, creando
- fraintendimenti
- funzioni duplicate nel codice
- parti diverse del sistema che lavorano incoerentemente

Un **Processo software** fornisce una struttura comune per decidere come:
- procedere 
- coordinarsi
- verificare il tutto in maniera sistematica

L'obbiettivo è ridurre il caos tra i vari team, e le persone all'interno dei vari team. Si devono stilare delle regole per agire coerentemente e in maniera ordinata, sfruttando così al massimo la potenzialità del gruppo

**Dal ciclo di vita al software:** Nella pratica ci si chiede:
- in quale ordine vanno affrontate le attività
- quanto approfondire ciascuna fase
- quanto spesso tornare indietro per rivedere le decisioni già prese

I processi software rispondono proprio a queste domande:
- sono strategie organizzative per combinare le stesse fasi del ciclo di vita (waterfall, incrementale, spirale, agile, ... )
- scegliamo il processo in base al contesto e al livello di rischio / incertezza

**Definizione (di processo):** <u>Un insieme strutturato di attività, ruoli e artefatti che guidano la costruzione e l'evoluzione di un software</u>. In un processo si specifica:
- chi è responsabile di questa cosa
- quali passi devono essere compiuti e in quale sequenza
- quali prodotti intermedi devono essere generati (documenti, modelli, codice, test, ... )

L’obiettivo non è aggiungere burocrazia, ma garantire che il progetto proceda in modo coerente, tracciabile e comprensibile per tutti i soggetti coinvolti.

**Perché serve anche ai buoni programmatori?** Nei piccoli progetti individuali il programmatore può controllare tutto “nella propria testa”, ma questa illusione sparisce appena: 
- il numero di funzionalità cresce
- entra in gioco un team
- il software deve essere mantenuto nel tempo

Senza un processo, le decisioni restano implicite e si perdono facilmente:
- diventa difficile capire perché certe scelte sono state fatte
- diventa rischioso modificare parti del sistema già rilasciate

Un processo aiuta a rendere esplicito il modo di lavorare, in modo che qualità e conoscenza non
dipendano dai singoli individui ma dal team.

**Attività comuni a tutti i processi:** Qualunque sia il processo adottato, ritroveremo sempre alcune attività fondamentali:
- specifica: capire che cosa il sistema deve fare e quali vincoli deve rispettare;
- sviluppo: progettare e implementare la soluzione;
- validazione: verificare che il sistema soddisfi i requisiti e si comporti come atteso;
- evoluzione: modificare e mantenere il sistema in risposta a nuove esigenze.

**Perché esistono diversi modelli?** Non tutti i progetti sono uguali: alcuni hanno requisiti molto stabili e vincoli normativi forti; altri nascono in ambienti altamente dinamici, con idee ancora da esplorare. Applicare sempre lo stesso processo, senza guardare al contesto, porta facilmente a inefficienze:
- troppa rigidità dove servirebbe flessibilità
- troppa improvvisazione dove servirebbe pianificazione

I modelli di processo offrono quindi strategie diverse, che l’ingegnere del software deve conoscere per poter scegliere (e adattare) quella più adatta al progetto.

**Modelli  di processo:** Nel seguito vedremo tre modelli classici di processo software.
Non sono schemi rigidi e mutuamente esclusivi: nella pratica molti progetti reali combinano elementi di questi modelli, adattandoli al contesto. L’obiettivo non è “scegliere il vincitore”, ma capire quando e perché ciascun modello può essere più adatto.

**Modello waterfall:** Il modello waterfall (a cascata) organizza lo sviluppo come una sequenza di fasi ben definite, che procedono in modo ordinato e lineare. Ogni fase produce artefatti che diventano l’input per la fase successiva:
- progettazione 
- implementazione 
- validazione 
- rilascio e manutenzione

La logica è quella di minimizzare i ritorni indietro: si passa alla fase successiva solo quando la
precedente è considerata completa.

*Vantaggi:* La struttura sequenziale è facile da comunicare a tutti gli stakeholder e rende chiaro in quale fase ci si trovi in ogni momento. La forte enfasi sulla documentazione iniziale favorisce:
- la tracciabilità delle decisioni
- la possibilità di svolgere audit e verifiche formali
- l’allineamento tra team tecnici e figure non tecniche.

È particolarmente adatto quando i requisiti sono ben noti e il contesto cambia lentamente, ad esempio in alcuni ambiti industriali regolamentati.

*Limiti:* Nella realtà, i requisiti cambiano spesso durante lo sviluppo e raramente sono completi o corretti fin dall'inizio. Nel waterfall, tornare indietro a modificare fasi già “chiuse” è costoso e complicato:
- correggere un requisito errato dopo aver progettato e implementato il sistema può richiedere molto lavoro.
- Gli utenti vedono il sistema completo solo alla fine, quando è difficile correggere scelte di fondo: il rischio di scoprire tardi malintesi o mancanze è elevato.

*Costo del cambiamento:* 
- Nel modello waterfall gli errori scoperti tardi costano molto di più.
- Il costo del cambiamento cresce in modo non lineare man mano che si avanza nelle fasi

**Modello incrementale:** Nel modello incrementale, il sistema non viene progettato e sviluppato “tutto insieme”, ma costruito per passi successivi. Ogni incremento introduce un insieme coerente di funzionalità che, una volta completate, producono una versione utilizzabile del sistema: 
1. Una base minima funzionante
2. Estensioni e miglioramenti progressivi

L’utente può iniziare a usare presto una parte del sistema e fornire feedback che guida gli
incrementi successivi.

*Vantaggi:* Con l’approccio incrementale il sistema cresce per piccoli rilasci, ciascuno utilizzabile. Dopo ogni rilascio gli utenti possono provare il sistema e dare feedback immediato. Questo permette di:
- correggere la direzione molto prima
- ridurre il rischio complessivo
- concentrare lo sforzo sulle funzionalità che portano più valore

*Limiti:* La flessibilità dell’approccio incrementale non significa assenza di pianificazione: serve comunque una visione architetturale complessiva del sistema. Se ogni incremento viene progettato come un mini–progetto indipendente, si rischia di accumulare debito tecnico, cioè il “costo futuro” dovuto a soluzioni veloci o provvisorie prese oggi. ad esempio:
- codice duplicato
- strutture provvisorie mai riscritte
- moduli con interfacce incoerenti
- mancanza di test
- architettura confusa.

Come un debito finanziario, il debito tecnico genera “interessi”: nel tempo tutto diventa più lento, più fragile e più costoso da modificare. Per funzionare bene, gli incrementi devono essere pensati come parti di un disegno coerente e non come patch aggiunte una dopo l’altra.

**Modello a spirale:** Il modello a spirale rappresenta lo sviluppo come una serie di cicli successivi, ognuno dei quali
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

**Confronto dei tre modelli visti:** - Il waterfall è adatto a contesti stabili, con requisiti chiari e forte enfasi sulla documentazione.
- L’approccio incrementale funziona bene quando vogliamo rilasciare valore progressivo e ottenere feedback continui.
- La spirale è indicata quando rischio e incertezza sono alti e servono cicli di esplorazione controllati.

Nella pratica, molti progetti combinano elementi di questi modelli, adattandoli al contesto.

| Waterfall               | Incrementale         | Spirale                             |
| ----------------------- | -------------------- | ----------------------------------- |
| Sequenziale, fasi fisse | Rilasci progressivi  | Cicli iterativi guidati dal rischio |
| Requisiti stabili       | Feedback frequente   | Gestione esplicita dell’incertezza  |
| Documentazione forte    | Adattamento continuo | Prototipi per ridurre rischi        |
| Poco flessibile         | Più flessibile       | Molto flessibile ma più costoso     |
**Perché servono metodi agili?** I modelli di processo più tradizionali, come il waterfall e altri approcci sequenziali, nascono in contesti in cui i requisiti sono relativamente stabili e i rilasci poco frequenti. Oggi molti progetti vivono in ambienti altamente dinamici: 
- richieste che cambiano rapidamente
- mercati competitivi
- aspettative di aggiornamenti continui

I metodi agili nascono per affrontare proprio questa realtà, introducendo iterazioni brevi, feedback costanti e collaborazione continua.
 
**Principi chiave dell'approccio agile:** L’attenzione si sposta da piani rigidi a una capacità di adattamento, in cui il cambiamento non è un problema da evitare ma un dato di fatto da gestire. Il software funzionante diventa la misura principale di avanzamento del progetto, più importante della quantità di documenti prodotti. Vengono valorizzate:
- le persone e le interazioni
- la comunicazione diretta
- il lavoro di team
- la collaborazione continua con gli stakeholder

**Triangolo del progetto: tradizionale vs agile:** Nei progetti tradizionali si tende a fissare lo scopo (tutte le funzionalità), lasciando variabili tempo e costo.
Nei metodi agili si fissano tempo (durata dello sprint) e costo (dimensione del team), rendendo flessibile lo scopo in base al valore e alle priorità.

**Scrum:** Scrum è un framework organizzativo: definisce come un team collabora, prende decisioni e rilascia valore in modo iterativo. Osserveremo in particolare:
- i **ruoli essenziali** (Product Owner, Scrum Master, Development Team)
- il **backlog** come lista ordinata di funzionalità
- gli **sprint** come cicli brevi che producono incrementi funzionanti.

**Ruoli fondamentali in Scrum:** Il **Product Owner** rappresenta gli interessi di utenti e stakeholders: 
- definisce la visione del prodotto
- stabilisce le priorità nel backlog

Lo **Scrum Master** facilita il lavoro del team:
- aiuta a rimuovere gli ostacoli
- promuove le pratiche agili e la collaborazione

Il **Development Team** 
- realizza concretamente gli incrementi del sistema

**Backlog e user stories:** Il **backlog** è l’elenco ordinato delle funzionalità, dei miglioramenti e dei bug da gestire sul sistema. Molti elementi del backlog vengono descritti come **user stories:** brevi frasi dal punto di vista dell’utente, **che chiariscono chi vuole cosa e per quale motivo**. 
Il **backlog** è un artefatto vivo: viene aggiornato continuamente **sulla base dei feedback e dei risultati degli sprint.**

**Sprint e incrementi:** Lo sviluppo procede per sprint, cioè intervalli di tempo brevi e fissi (ad esempio 2 settimane), con obiettivi limitati ma chiari. Alla fine di ogni sprint il team dovrebbe poter mostrare un incremento funzionante del sistema, anche se piccolo. Questo ritmo cadenzato consente:
- di valutare spesso i progressi
- di adattare il piano sulla base dell’esperienza
- di coinvolgere regolarmente gli stakeholder

Le funzionalità non vengono affrontate tutte insieme, ma organizzate in:
- epic $\rightarrow$ macro–aree funzionali
- user stories $\rightarrow$ unità più piccole e verificabili
- sprint $\rightarrow$ passi concreti che producono un incremento funzionante.

Questo rende visibile l’evoluzione del sistema e permette di: 
- adattare le priorità nel tempo
- collegare ogni incremento a valore per gli utenti

A partire dalle epic possiamo ricavare user stories più piccole e concrete, ciascuna con un obiettivo chiaro e verificabile. Le user stories vengono scritte dal punto di vista dell’utente, specificando:
- chi è l’attore
- che cosa vuole ottenere
- perché quella funzionalità è utile.
###### Requisiti funzionali e non funzionali
**Requisito:** <u>Un requisito è una descrizione chiara e condivisa di qualcosa che il sistema deve fare oppure di un vincolo che deve rispettare.</u> Possiamo vederlo come una promessa tra chi richiede il sistema e chi lo realizza. Un requisito ben scritto:
- riduce le ambiguità e le interpretazioni personali (devono essere chiari); 
- permette di verificare in modo oggettivo se il sistema si comporta come atteso (se ho dei vincoli questi vengono rispettati); 
- guida direttamente la progettazione dei test. 
- Criterio guida: un requisito deve essere verificabile (e, quando possibile, anche misurabile). 

> [!example] 
> - Il sistema deve permettere allo studente di iscriversi a un appello d’esame.
> - L’accesso alla piattaforma deve essere possibile solo tramite credenziali istituzionali.

Un requisito può essere espresso a due livelli di dettaglio:
1. *Lato utente:* linguaggio semplice, orientato ai bisogni, agli obiettivi e al contesto $\rightarrow$ comprensibile a chi usa il sistema
2. *Lato sistema:* descrizione più precisa, utile per progettare e implementare concretamente -> utilizzato da sviluppatori e progettisti

Non sono due categorie distinte, ma di due modi di esprimere la stessa necessità. Questi due li
velli sono spesso indicati come *requisiti utente* e *requisiti di sistema*.

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
2. *Requisiti non funzionali:* descrivono come il sistema deve funzionare, quali vincoli deve rispettare e quale livello di qualità deve garantire (in che condizioni questi servizi devono funzionare bene)

<u>Queste due categorie si completano in quanto i funzionali definiscono la funzione mentre i non funzionali definiscono qualità, limiti e contesto operativo</u>

**Requisiti funzionali:** Sono i requisiti più intuitivi, descrivono le azioni che il sistema deve permettere e le interazioni possibili tra utenti e piattaforma. Riguardano tutto ciò che è osservabile durante l’uso:
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

**Requisiti non funzionali:** Non descrivono azioni specifiche del sistema, ma come il sistema deve comportarsi, quali vincoli deve rispettare e quale livello di qualità deve garantire. Questi riguardano aspetti trasversali che influenzano l’esperienza d’uso e l’affidabilità complessiva, come:
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

**Requisiti verificabili:** <u>Un requisito è utile solo se possiamo verificarlo in modo oggettivo</u>. Spesso nelle fasi iniziali emergono affermazioni molto generali: «facile da usare», «veloce», «sicuro», «intuitivo». Questi obiettivi sono utili per capire la direzione desiderata, ma non dicono concretamente come il sistema deve comportarsi. In questa forma non sono verificabili: due persone possono interpretarli in modo anche molto diverso. Occorre quindi un passaggio successivo: trasformare gli obiettivi in requisiti misurabili, cioè descrizioni che permettono di verificare in modo oggettivo se il sistema rispetta ciò che promette.

> [!example]  
> Da: «Il sistema deve essere veloce.» 
> 
>  A: «Il sistema deve rispondere entro 2 secondi per almeno il 95% delle richieste di iscrizione agli appelli.»

Perché questa trasformazione è importante? 
- Permette di verificare il comportamento del sistema in modo oggettivo. 
- Evita interpretazioni diverse tra utenti, progettisti e sviluppatori. 
- Guida scelte architetturali e tecniche coerenti.

È del tutto normale che, nelle prime fasi, i requisiti non siano perfetti: alcuni mancano semplicemente perché nessuno ci ha ancora pensato, altri risultano in conflitto (ad esempio: «facile da usare» ma anche «richiedere password molto complesse»), altri ancora riflettono punti di vista diversi tra studenti, docenti, segreteria, tecnici. Questa situazione non è un difetto del progetto, ma una caratteristica naturale dei sistemi reali: le esigenze si scoprono poco alla volta e spesso emergono contraddizioni da risolvere. L’analisi dei requisiti serve proprio a: 
- chiarire cosa si intende davvero
- verificare che le richieste siano realistiche e compatibili
- negoziare priorità e compromessi tra stakeholder con bisogni diversi

I requisiti non emergono spontaneamente dal codice o dalla tecnologia: occorre andare a cercarli parlando con chi userà davvero il sistema. Per raccoglierli si usano strumenti diversi, spesso complementari:
- colloqui con gli utenti per capire bisogni, difficoltà e aspettative; 
- osservazione delle attività reali per vedere come i processi funzionano oggi, senza filtri; 
- analisi di procedure, moduli e documenti esistenti per cogliere vincoli già presenti nell'organizzazione; 
- prototipi e dimostrazioni preliminari per far emergere idee e problemi che non verrebbero fuori solo parlando. 

Raccogliere requisiti è un’attività progressiva: solo con un confronto continuo con gli stakeholder si costruisce una visione chiara e condivisa del sistema. Una volta raccolti i bisogni, serve però un modo semplice e comprensibile per rappresentarli: da qui nasce l’uso delle user stories.

**User stories:** Una user story descrive un’esigenza dal punto di vista dell’utente, in modo semplice e comprensibile anche a chi non è tecnico. Risponde a tre domande fondamentali:
- chi usa il sistema
- che cosa vuole ottenere
- perché quell'azione è importante. 

La forma tipica è molto essenziale:  «Come (tipo di utente), voglio (azione) per (obiettivo/beneficio).». Le user stories non descrivono i dettagli tecnici, ma chiariscono l’intenzione: servono ad allineare utenti e sviluppatori e a popolare un backlog comprensibile e prioritizzabile. Per essere davvero utili, dovrebbero:
- avere un chiaro valore per l’utente
- essere abbastanza piccole da stare in uno sprint
- risultare verificabili tramite criteri di accettazione espliciti.

> [!info]
> - **User story:** frase breve orientata a chi-cosa-perché, utile per raccogliere e prioritizzare bisogni nel backlog.
> - **Use case:** descrizione più strutturata dell’interazione, che useremo più avanti nella parte di modellazione UML.

**Criteri di accettazione:** Ogni user story deve essere accompagnata da uno o più criteri di accettazione, cioè condizioni che stabiliscono con chiarezza quando la funzionalità può essere considerata davvero completata. I criteri di accettazione servono a: 
- eliminare ambiguità su cosa ci si aspetta dal sistema; 
- definire un confine preciso tra “fatto” e “non ancora fatto”; 
- costituire una base solida per i test che verificheranno il corretto comportamento del sistema. 

<u>In altre parole: la user story racconta l’intenzione, i criteri di accettazione ne specificano la verifica concreta</u>. Alcuni criteri possono essere espressi anche in forma strutturata, ad esempio con frasi del tipo «Dato che … Quando … Allora … », facilitando la verifica sistematica del comportamento del sistema.

> [!example]
> User story 
> - «Come studente, voglio iscrivermi a un appello d’esame per poter sostenere la prova nella data che scelgo.» 
>
> Criteri di accettazione 
> - Lo studente vede solo gli appelli disponibili per il proprio corso. 
> - Se i posti sono esauriti, il sistema mostra un messaggio chiaro e comprensibile.
> - Dopo l’iscrizione, lo studente visualizza immediatamente la conferma nella propria area personale. 
> - Il docente può consultare l’elenco aggiornato degli iscritti all'appello.

**Backlog:** Le user stories raccolte diventano elementi del backlog che guida lo sviluppo: una raccolta ordinata di ciò che il sistema dovrà offrire. Non tutti gli elementi hanno lo stesso peso: occorre stabilire delle priorità. Alcune funzionalità vanno affrontate prima perché: 
- portano più valore agli utenti; 
- comportano rischi tecnici o organizzativi da affrontare presto;
- dipendono da altre funzionalità che devono essere pronte. 

Il backlog è un artefatto vivo: cambia, si aggiorna, viene riorganizzata man mano che emergono nuove informazioni o nuove esigenze.
I requisiti non funzionali non devono essere trascurati: anche loro entrano a pieno titolo nel backlog. Alcuni diventano attività specifiche, da pianificare come vere e proprie funzionalità: 
- migliorare i tempi di risposta nelle operazioni più frequenti; 
- introdurre o potenziare un meccanismo di sicurezza; 
- rendere l’interfaccia più accessibile o più semplice da usare. 

Altri, invece, sono vincoli trasversali: non compaiono come singole attività, ma devono essere rispettati in ogni funzionalità sviluppata (ad esempio sicurezza, normativa sui dati, compatibilità con sistemi esistenti). Ignorarli nelle prime fasi porta quasi sempre a interventi più costosi, più complessi e più rischiosi nelle fasi successive. 
###### Modellazione concettuale      
<u><b>Dominio</b> = Per dominio si intende l'insieme di: persone, oggetti, regole e processi reali che il software deve rappresentare e supportare.</u>

**Modellazione concettuale:** La modellazione concettuale è l'attività che mira a:
- identificare i concetti importanti del dominio; 
- chiarire quali responsabilità ha ciascun concetto; 
- descrivere le relazioni tra questi concetti.

È una descrizione logica del problema, indipendente dalla soluzione tecnica. Nella pratica: stiamo cercando di costruire un vocabolario condiviso che tutti possano capire (utenti, analisti, sviluppatori).

**Concetti:**<u>Per concetti intendiamo gli elementi significativi del dominio</u> (es. Studente, Corso, Appello), non ancora classi, tabelle o dettagli implementativi. I requisiti che abbiamo scritto non sono solo “liste di funzionalità”: contengono già molti suggerimenti sul mondo che il sistema deve rappresentare. Leggendo frasi come:  «Lo studente si iscrive a un appello d’esame.», «Il docente registra il voto.», «Lo studente visualizza la propria carriera.» possiamo chiederci: quali entità compaiono? quali azioni vengono svolte? quali relazioni tra concetti sono implicite? La modellazione concettuale parte spesso proprio da questo lavoro di “lettura attenta” dei requisiti. Dal punto di vista concettuale, nei requisiti ci interessano in particolare: 
- i sostantivi (studente, appello, corso, voto, verbale, carriera...) che possono suggerire entità del dominio; 
- i verbi (iscriversi, registrare, consultare, creare, chiudere...) che possono suggerire azioni o relazioni; 
- le espressioni di vincolo (solo se, fino a, entro...) che possono nascondere regole del dominio. 

Non è una trasformazione meccanica, ma una guida: ci aiuta a individuare i concetti davvero importanti.

**Entità concettuali:** Un’entità concettuale è un elemento del mondo che vogliamo rappresentare in modo esplicito nel sistema:
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
- Studente $\rightarrow$ Iscrizione: crea una nuova iscrizione a un appello; 
- Docente $\rightarrow$ Voto: registra o aggiorna l’esito di uno studente; 
- Segreteria $\rightarrow$ Corso/Appello: crea o aggiorna i dati organizzativi. 

Queste relazioni sono più vicine al comportamento del sistema: hanno tipicamente precondizioni, effetti e controlli.

**Regole del dominio:** <u>Le regole del dominio esprimono vincoli che devono valere sempre, indipendentemente dalla tecnologia</u>. 

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
> Oltre alle entità del dominio, è importante chiarire chi interagisce con il sistema: studenti, docenti, segreteria, eventuali amministratori di sistema. Dobbiamo chiederci: che cosa sta dentro UniManager e che cosa sta fuori?  
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

**Rappresentare la struttura concettuale:** Dopo la vista delle azioni, passiamo alla vista strutturale del dominio con una mappa concettuale: 
- nodi per le entità (Studente, Corso, Appello, Iscrizione, Voto ...);
- collegamenti che indicano le relazioni principali (ad esempio “Studente è iscritto a Appello tramite Iscrizione”). 

Usiamo una rappresentazione informale: 
- niente simboli complicati, 
- etichette in italiano, 
- solo gli elementi davvero utili per capire il dominio. 

Questa mappa sarà la base del futuro modello del dominio in notazione più formale. Come prima approssimazione, possiamo disegnare: 
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

**Checklist rapido:** Prima di proseguire, facciamo un controllo rapido: 
- *Nomi chiari:* ogni entità ha un nome comprensibile anche a chi non è tecnico? 
- *Relazioni esplicite:* i collegamenti principali sono scritti e coerenti con i requisiti? 
- *Niente dettagli tecnici:* stiamo evitando classi, tabelle DB, API e scelte implementative? 

Se una risposta è “no”, conviene correggere subito: nelle fasi successive gli errori costano di più. Una buona modellazione concettuale: 
- riduce le ambiguità, 
- aiuta a prendere decisioni progettuali più consapevoli, 
- facilita il dialogo tra persone con competenze diverse. 

Al contrario, se si salta questa fase: 
- si rischiano fraintendimenti tra utenti e sviluppatori; 
- si introducono vincoli tecnici non necessari; 
- diventa più difficile correggere gli errori in seguito. 

Investire tempo qui fa risparmiare tempo nelle fasi successive.

**Errori tipici:** Alcuni errori ricorrenti: 
- confondere subito entità concettuali e classi di implementazione; 
- introdurre dettagli tecnici (database, interfacce grafiche...) troppo presto; 
- creare entità “che fanno tutto” e relazioni poco chiare; 
- usare nomi vaghi o poco significativi (dati1, oggetto, elemento...).

È importante ricordare che in questa fase lo scopo non è “scrivere codice”, ma capire bene il problema.
###### Use case diagram      
**UML:** Iniziamo ad usare una notazione standard per rappresentare i modelli: UML (unified modelling language) che è un linguaggio visivo condiviso. UML è un linguaggio di modellazione standard, non un linguaggio di programmazione. Fornisce un insieme di diagrammi per descrivere:
- la struttura statica del sistema
- il comportamento dinamico
- le interazioni tra parti del sistema e attori esterni. 

Non obbliga a un particolare processo (waterfall, agile,...), ma può essere usato in tutti. 
UML nasce negli anni ’90 per unificare diverse notazioni di analisi e progettazione. I principali contributori iniziali sono Grady Booch, Ivar Jacobson e James Rumbaugh. La prima standardizzazione avviene in ambito OMG (Object Management Group). Nel tempo UML ha avuto più versioni: la famiglia 1.x e poi 2.x, con estensioni e chiarimenti della notazione. Oggi UML resta un riferimento diffuso per comunicare modelli in modo condiviso tra team diversi. UML mette a disposizione diversi tipi di diagrammi. In questo corso vedremo soprattutto: 
- Use Case Diagram — vista del sistema dal punto di vista degli attori esterni; 
- Modello di dominio — modello concettuale del dominio, spesso rappresentato con un class diagram semplificato; 
- Class Diagram — struttura statica del dominio (classi, attributi, associazioni); 
- Sequence Diagram — interazioni nel tempo tra oggetti/attori; 
- Activity Diagram — flusso di attività, decisioni e responsabilità di processo. 

UML include anche altri diagrammi (ad esempio State Machine, Component, Deployment). Inizieremo dall'Use Case Diagram, perché: 
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

**Attore:** <u>Un attore è un ruolo esterno che interagisce con il sistema</u>:
- una persona (studente, docente), 
- un altro sistema (sistemi di autenticazione, basi dati esterne), 
- un dispositivo. 

È importante ricordare che: 
- non rappresenta una singola persona, ma un tipo di utilizzatore; 
- è sempre fuori dal sistema. 

In UniManager i nostri attori principali sono: Studente, Docente, Segreteria.

**Caso d'uso:** Un caso d’uso rappresenta un obiettivo significativo per un attore non è un singolo “click” ne una funzione tecnica (“query sul DB”), ma un risultato che ha senso dal punto di vista dell’utente. 

> [!example] Esempi tipici in UniManager: 
> -  Iscrizione all’appello, 
> - Registrazione voti, 
> - Consultazione libretto, 
> - Creazione appello

Nella notazione UML, i casi d’uso sono rappresentati come ellissi con un nome. Nei diagrammi completi i casi d’uso stanno dentro il confine del sistema; negli esempi minimali di queste slide il rettangolo può essere omesso per semplicità visiva. Nella lezione precedente abbiamo letto frasi come: 
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

> [!tip] ![[Screenshot From 2026-03-12 15-15-14.png|200]]

**Relazione *include* (obbligatorio):** <u>Include rappresenta un comportamento comune e obbligatorio tra più casi d’uso</u>. Un caso d’uso A include B se: 
- ogni volta che si esegue A, a un certo punto è necessario eseguire B; 
- B è un “pezzo riusabile” usato da più casi d’uso. 

A livello grafico: 
- la relazione è una freccia tratteggiata con l’etichetta «include»; 
- la freccia punta verso il caso d’uso incluso (da A verso B). 

Serve quando una parte del comportamento è identica e va eseguita sempre prima o durante un caso d’uso principale.

> [!example] Esempio uso include
>  ![[Screenshot From 2026-03-12 15-18-15.png|200]]

**Relazione *extend* (opzionale):** <u>Extend rappresenta una variante opzionale di un caso d’uso</u>. Un caso d’uso B estende A se: 
- A descrive il flusso “normale”; 
- B aggiunge passi opzionali o condizioni particolari. 

A livello grafico: 
- la relazione è una freccia tratteggiata con etichetta «extend»; 
- la freccia punta verso il caso d’uso esteso (da B verso A). 
- Serve quando un comportamento aggiuntivo può avvenire, ma non è obbligatorio.

> [!example] Esempio uso extend
> ![[Screenshot From 2026-03-12 15-21-24.png|200]]

**Generalizzazione tra attori:** È possibile organizzare gli attori in una gerarchia, serve quando più attori condividono gli stessi casi d’uso. Significato: 
- l’attore “padre” rappresenta un ruolo più generale; 
- l’attore “figlio” eredita le sue interazioni. 

A livello grafico (UML): 
- si usa una linea continua con una freccia triangolare non piena; 
- la freccia punta verso l’attore più generale (da Specializzato $\rightarrow$ Generale). 

Va usata solo se semplifica il modello o riduce ripetizioni. In alcuni casi può essere utile introdurre un attore più generale: Utente come ruolo comune. Attori più specifici: 
- Studente, 
- Docente. 

La generalizzazione permette di evitare ripetizioni: 
- i casi d’uso comuni a tutti gli utenti possono collegarsi a Utente; 
- i casi specifici rimangono ai ruoli specializzati.

**Note e commenti:** In UML possiamo aggiungere note/commenti per chiarire assunzioni e scelte. Le note aiutano a esplicitare dettagli utili senza appesantire i casi d’uso. Sono particolarmente utili quando una regola non è rappresentabile in modo diretto con soli attori e relazioni. Una buona pratica è quella di scrivere note brevi e mirate, collegate solo ai punti davvero ambigui.

> [!example] ![[Screenshot From 2026-03-12 15-27-37.png|200]]

**Package e raggruppamenti:** In UML possiamo raggruppare casi d’uso in package/aree funzionali. Ciò è utile quando il diagramma cresce e vogliamo separare aree funzionali. La motivazione principale è la leggibilità: meno incroci visivi, più chiarezza su responsabilità e perimetri. I package non cambiano il comportamento dei casi d’uso: organizzano la vista del modello.

> [!example] ![[Screenshot From 2026-03-12 15-27-47.png|200]]

**Molteplicità sulle relazioni:** In alcuni casi possiamo esplicitare la molteplicità sulle relazioni (senza esagerare per rimanere il più asciutti possibile, altrimenti si rischia solo rumore). Va usata solo quando aggiunge chiarezza al significato del modello. Nei Use Case Diagram è una notazione possibile ma non centrale (più comune nei Class Diagram). 

> [!example]
> - Esempio 1 (1–*): uno Studente può effettuare molte iscrizioni nel tempo. 
> - Esempio 2 (1–0..*): un Docente può registrare zero o più voti. 

**Vincoli:** In UML possiamo esprimere vincoli con testo tra parentesi graffe. I vincoli servono per esplicitare regole di business non evidenti dal solo diagramma (riguardo l'implementazione). Sono utili quando una condizione è decisiva per il dominio (scadenze, prerequisiti, autorizzazioni). Un vincolo deve essere verificabile: formulazioni vaghe riducono il valore del modello.

> [!example] 
> ![[Screenshot From 2026-03-12 15-36-05.png|200]] 
> ![[Screenshot From 2026-03-12 15-36-23.png|200]]

**Good practices:**
- Usare i casi d’uso per descrivere obiettivi, non schermate o click. 
- Mantenere un livello di dettaglio gestibile: pochi casi d’uso macro, non decine di micro-casi. 
- Evitare di usare include/extend in modo eccessivo o artificiale. 
- Annotare, quando serve, qualche breve descrizione testuale per chiarire le intenzioni. 
- Ricordare che il Use Case Diagram è uno strumento di comunicazione con gli stakeholder: deve essere leggibile anche da chi non è tecnico.

**Descrizione testuale:** Il diagramma dei casi d’uso dà una vista sintetica, ma spesso non basta. Per i casi d’uso più importanti conviene scrivere una scheda testuale con: 
- Nome del caso d’uso; 
- Attori coinvolti; 
- Breve descrizione (obiettivo principale); 
- Precondizioni (cosa deve essere vero prima); 
- Postcondizioni (cosa deve essere vero dopo); 
- Flusso principale (passi “normali”); 
- eventuali flussi alternativi / eccezioni. 

L’idea è avere un livello di dettaglio sufficiente per: 
- chiarire il comportamento atteso, 
- preparare i diagrammi successivi (sequence, activity... ).

> [!example] ![[Screenshot From 2026-03-12 15-40-44.png|200]]

**Da use case a domain model:** Nella lezione precedente abbiamo visto: 
- attori esterni al sistema; 
- casi d’uso come obiettivi significativi; 
- relazioni include/extend e generalizzazioni. 

Il diagramma dei casi d’uso risponde alla domanda: «Chi usa il sistema e per quali obiettivi?» Oggi cambiamo prospettiva: non più “chi usa il sistema”, ma “che cosa esiste nel mondo del sistema”.

Requisiti, user stories e casi d’uso contengono storie: 
- «Lo studente si iscrive a un appello.» 
- «Il docente registra un voto.» 
- «La segreteria crea un appello.» 

Da queste storie possiamo estrarre le entità concettuali del dominio: Studente, Docente, Corso, Appello, Iscrizione, Voto…

<div style="display: flex; justify-content: center;">
  <img src="Screenshot From 2026-03-18 14-54-16.png" width="300">
</div>

**Domain model:** È una rappresentazione concettuale del dominio del problema. Mostra: 
- i concetti rilevanti (entità, oggetti, ruoli); 
- le relazioni tra questi concetti; 
- alcuni attributi essenziali che li caratterizzano. 

È una fotografia del mondo a cui il sistema si riferisce: 
- non contiene dettagli tecnici; 
- non parla di database, API, interfacce grafiche; 
- non contiene metodi o algoritmi.

Non è:
1. un diagramma delle classi implementativo: 
	- non rappresenta direttamente il codice Java/C++/Python; 
	- non specifica metodi, eccezioni, pattern di design. 
2. il modello del database: 
	- non specifica chiavi primarie, indici, tabelle tecniche; 
	- non è ottimizzato per performance o storage. 
3. un diagramma dei casi d’uso: 
	- non descrive scenari o obiettivi, 
	- ma gli oggetti coinvolti in quegli scenari.

Aiuta a capire il problema prima di pensare alla soluzione. È un linguaggio condiviso tra: 
- sviluppatori, 
- analisti, 
- domain expert, 
- stakeholder non tecnici. 

È un termometro della comprensione del dominio: 
- se il modello è confuso, probabilmente non abbiamo capito bene il problema; 
- se è chiaro e coerente, il passaggio alla progettazione sarà più semplice.

Punto di partenza: descrizioni testuali, requisiti, casi d’uso:
- evidenziare i sostantivi significativi (candidati concetti); 
- evidenziare i verbi principali (candidati per associazioni o responsabilità); 
- discutere con gli stakeholder cosa è davvero importante. 

Domanda guida: «Se togliamo questo concetto, capiamo ancora il dominio?»

Un concetto è un buon candidato se: 
- è nominato spesso nei requisiti; 
- ha regole di business proprie; 
- gli stakeholder lo riconoscono come entità reale del loro mondo; 
- può cambiare nel tempo e queste modifiche sono importanti. 

**Attributi essenziali:** Gli attributi descrivono informazioni importanti su un concetto. Dovrebbero essere: 
- pochi (minimalismo); 
- stabili nel tempo; 
- non derivabili facilmente da altri; 
- rilevanti per le regole di business. 

> [!example] Esempio:
> - Studente: matricola, nome, cognome; 
> - meno adatti come attributi nel Domain Model: colore preferito, ultima data di login.

**Errori frequenti degli attributi:**
- Confondere il Domain Model con il database: inserire id tecnici, timestamp di audit, campi di log…
- Aggiungere attributi troppo tecnici: flag di caching, token di sessione, dettagli di autenticazione. 
- Inserire attributi derivabili: numero di esami superati (derivabile dalle iscrizioni e dai voti) o media voti (derivabile dai singoli voti).

**Associazioni tra concetti:** Le associazioni rappresentano legami concettuali tra elementi del modello. Rispondono alla domanda: «Questi due concetti hanno una relazione significativa nel dominio?» 

> [!example] Esempi: 
> - uno Studente è iscritto a un Appello; 
> - un Appello appartiene a un Corso; 
> - un Docente tiene un Corso.

**Molteplicità:** La molteplicità esprime quante istanze di un concetto possono essere collegate a un altro. Lettura rapida: 
- 1: esattamente una; 
- 0..1: opzionale; 
- 0..:* zero o molte; 
- 1..:* almeno una

<div style="display: flex; justify-content: center;">
  <img src="Screenshot From 2026-03-18 15-05-46.png" width="200">
</div>

Nelle associazioni, spesso ci interessa la molteplicità: 
- uno studente può avere molte iscrizioni; 
- un appello riguarda un solo corso; 
- un corso può avere più appelli. 

Le molteplicità non sono numerini casuali: 
- esprimono vincoli reali del dominio; 
- aiutano a scoprire regole di business nascoste.

> [!tip] A partire dalle lezioni precedenti, possiamo individuare alcuni concetti chiave: 
> - Studente 
> - Docente 
> - Corso 
> - Appello 
> - Iscrizione 
> - (eventualmente) Voto, Carriera, Piano di studi…
> 
> Ogni concetto cattura un nucleo di significato del dominio universitario.

> [!example] Rappresentazione concettuale del dominio di UniManager: 
> - concetti principali, attributi essenziali, relazioni e molteplicità; 
> - nessun dettaglio di implementazione (tipi, metodi, framework…)
> ![[Screenshot From 2026-03-18 15-09-32.png|200]]

**Domain model come strumento di verifica:** Possiamo usare il Domain Model per fare domande al domain expert: 
- «Esistono appelli senza corso?»; 
- «Uno studente può iscriversi più volte allo stesso appello?»; 
- «Può esistere un corso senza docente associato?». 

Ogni risposta può portare a correggere o raffinare il modello:
- aggiungere o rimuovere un concetto; 
- modificare una molteplicità; introdurre una nuova associazione.

**Da domain model a class diagram:** Il Domain Model è concettuale e indipendente dalla tecnologia. Il Class Diagram inizia a guardare verso il software: 
- introduce il concetto di classe; 
- prepara il terreno per attributi e metodi; 
- sarà la base per codice e architettura. 

Domain Model: 
- concetti del mondo reale; 
- attributi essenziali per capire il problema; 
- relazioni concettuali tra entità. 

Class Diagram: 
- classi del sistema software; 
- attributi che saranno effettivamente memorizzati; 
- relazioni che guideranno la struttura del codice. 

Spesso i nomi coincidono, ma il punto di vista è diverso.

<div style="display: flex; justify-content: center;">
  <img src="Screenshot From 2026-03-18 15-13-44.png" width="200">
</div>
 
 Il class diagram ci aiuta a non fissare troppo presto una soluzione tecnica. Mantiene un livello di astrazione sufficientemente alto, ma già orientato al software. Permette di: 
 - verificare la struttura del dominio; 
 - discutere con il team prima di scrivere codice; 
 - cambiare idea senza riscrivere tutto.

**Da domain model alle classi: checklist:** Cosa portiamo nel primo Class Diagram: 
- il nome (chiaro e significativo); 
- gli attributi essenziali del dominio; 
- le relazioni con altre classi (ispirate al Domain Model).

Cosa lasciamo fuori in questa fase: 
- metodi e logica di dettaglio; 
- classi tecniche (Repository, DAO = Data Access Object, Controller, helper…); 
- Repository/DAO: accesso e persistenza dei dati; Controller: coordinamento tra input e logica applicativa;
- dettagli legati a framework specifici. 

Obiettivo: restare vicini al dominio, ma con una struttura già utile al design software.

**Notazione UML: visibilità degli attributi:** In un Class Diagram possiamo indicare la visibilità: 
- + public — visibile a tutti 
- private — visibile solo all’oggetto stesso  
- `#` protected — visibile a classe e sottoclassi 
- ~ package — visibile all’interno del package

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260318151844.png" width="200">
</div>

**Attributo derivato:** Un attributo è derivato se può essere calcolato da altri dati. Notazione: si usa il simbolo / davanti al nome.

<div style="display: flex; justify-content: center;">
  <img src="Screenshot From 2026-03-18 15-19-42.png" width="200">
</div>

**Attributo statico:** La notazione UML prevede attributi di classe (statici). Notazione: nome sottolineato.

<div style="display: flex; justify-content: center;">
  <img src="Screenshot From 2026-03-18 15-20-50.png" width="200">
</div>

**Ereditarietà:** La generalizzazione si indica con una freccia a triangolo vuoto verso la superclasse. Le sottoclassi ereditano attributi e operazioni.

<div style="display: flex; justify-content: center;">
  <img src="Screenshot From 2026-03-18 15-22-00.png" width="200">
</div>

**Aggregazione e composizione:** Sono due tipi particolari di associazione UML. In entrambi i casi c’è un rapporto tra un tutto e una parte. 
- *Aggregazione* (rombo vuoto): la parte può esistere anche senza quel tutto. 
- *Composizione* (rombo pieno): la parte dipende dal tutto per esistere.

*Aggregazione:* la parte resta indipendente dal tutto. 
- nell’esempio: Corso (lato con rombo) è il tutto, Docente è la parte aggregata; 
- In breve: se elimino un Corso, il Docente resta.

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260318152430.png" width="200">
</div>

*Composizione:* la parte dipende dal tutto. 
- nell’esempio: Appello (lato con rombo pieno) è il tutto, Iscrizione è la parte contenuta; 
- In breve: se elimino un Appello, l’Iscrizione non esiste più.

<div style="display: flex; justify-content: center;">
  <img src="Screenshot From 2026-03-18 15-25-22.png" width="200">
</div>

**Class diagram di *Unimanager:*** Traduzione del Domain Model in una rappresentazione orientata al software: 
- classi con attributi tipizzati; 
- relazioni che suggeriscono i futuri collegamenti tra oggetti; 
- in questo diagramma mostriamo solo associazioni.

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260318152700.png" width="200">
</div>

**Errori tipici nel passaggio a class diagram:** Aggiungere subito: 
- metodi per ogni scenario immaginabile; 
- classi per ogni schermata dell’interfaccia grafica; 
- dettagli di persistenza (DB, ORM, ecc.). 
- ORM: meccanismo per mappare oggetti ↔ database;

Confondere classi del dominio con: 
- classi di infrastruttura, 
- classi di supporto (logger, helper, utility…). 

Dimenticare che il Class Diagram deve essere ancora vicino al dominio, non al framework.

**Struittura statica vs comportamento dinamico:** Il Domain Model e il Class Diagram descrivono la struttura statica: 
- quali entità esistono; 
- come sono collegate; 
- quali informazioni portano con sé. 

Non dicono ancora: 
- come gli oggetti collaborano; 
- quali messaggi si scambiano; 
- in che ordine avvengono le operazioni. 

###### Sequence diagram e activity diagram
Nella lezione precedente abbiamo definito la struttura del sistema: 
- **Domain Model:** concetti del dominio e relazioni principali;
- **Class Diagram:** classi software e legami tra classi. 

Oggi passiamo al livello dinamico: 
- come gli oggetti collaborano durante uno scenario; 
- in quale ordine avvengono azioni e decisioni. 

Useremo due viste complementari: 
- **Sequence Diagram** per le interazioni tra partecipanti; 
- **Activity Diagram** per il flusso del processo.

###### Sequence diagram
Si utilizza per chiarire come uno scenario viene eseguito nel tempo. Risponde a tre domande operative: 
- chi interagisce con chi; 
- in quale ordine avvengono i messaggi; 
- chi prende in carico ogni responsabilità. 

È particolarmente utile per scenari di caso d’uso con esito positivo e varianti di errore. Non è il diagramma migliore per descrivere l’intero processo: quello è il ruolo degli Activity Diagram

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319145158.png" width="200">
</div>

**Lettura:** La lettura corretta segue due assi: 
- verticale: il tempo scorre dall’alto verso il basso; 
- orizzontale: partecipanti diversi coinvolti nello scenario. 

Ogni freccia rappresenta un messaggio tra lifeline e quindi un passo osservabile del comportamento. La lettura consigliata interpreta il diagramma come una traccia esecutiva: 
- primo messaggio, eventuale risposta, passo successivo; 
- in presenza di rami, capire quale condizione seleziona il percorso. 

**Lifeline:** Le lifeline identificano i partecipanti coinvolti nello scenario. Qui distinguiamo attore, interfaccia e servizio: ruoli diversi, responsabilità diverse. L’obiettivo è rendere esplicito chi partecipa prima di parlare dei dettagli dei messaggi. Quando usarlo: all’inizio, per fissare i confini dello scenario. Errore tipico: introdurre partecipanti non necessari al caso in esame.

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319145433.png" width="200">
</div>

**Attivazione:** La barra di attivazione indica il periodo in cui il partecipante sta eseguendo lavoro. Nel diagramma si vede chiaramente quando *svc:AppelloService* è attivo e quando restituisce il controllo a *ui:Portale*. La freccia tratteggiata *esito(OK)* è un return message: rappresenta la risposta alla chiamata, non un messaggio asincrono. Serve a leggere meglio durata logica e nidificazione delle chiamate. Quando usarlo: quando la durata logica del lavoro è rilevante. Errore tipico: disegnare attivazioni ovunque senza valore informativo.

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319145618.png" width="200">
</div>

**Messaggio sincrono:** Il chiamante invia il messaggio e attende il completamento. Nel nostro esempio *ui:Portale* invia la richiesta a *svc:AppelloService* e resta in attesa dell’esito. Il messaggio sincrono rappresenta una dipendenza temporale immediata tra chiamata e risultato. Quando usarlo: per chiamate con dipendenza immediata dal risultato. Errore tipico: usare sync anche per notifiche/eventi che non bloccano.

**Messaggio asincrono:** Il chiamante invia il messaggio ma prosegue senza blocco. Nel nostro esempio *svc:AppelloService* invia una notifica a *noti:Notifier* senza attendere risposta immediata. Qui il focus è la disaccoppiata temporale tra chi invia e chi riceve. Quando usarlo: in integrazioni event-driven o code di messaggi. Errore tipico: usare async senza chiarire chi gestisce l’esito.

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319145841.png" width="200">
</div>

**Return message:** La freccia tratteggiata rappresenta il risultato della richiesta. Qui *svc:AppelloService* restituisce a *ui:Portale* l’esito della richiesta (OK/errore). Nelle slide useremo il return solo quando porta informazione utile. Quando usarlo: se il valore di ritorno guida una decisione successiva. Errore tipico: mostrare tutti i return e appesantire il diagramma.

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319145618.png" width="200">
</div>

**Create:** Mostra che un oggetto nasce durante l’interazione, non esiste dall’inizio. In UniManager, *svc:AppelloService* crea una nuova *iscr:Iscrizione* quando i vincoli sono soddisfatti. Questo chiarisce ownership e momento di istanziazione. Quando usarlo: se il momento di creazione è parte del comportamento. Errore tipico: creare oggetti tecnici ovunque invece di riusarli.

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319150059.png" width="200">
</div>

**Destroy:** Rappresenta la fine del ciclo di vita di un oggetto temporaneo. In UniManager, *svc:AppelloService* può eliminare una *bozza:Iscrizione* scaduta/non confermata. Utile per risorse effimere e cleanup esplicito. Quando usarlo: per risorse temporanee con fine vita esplicita. Errore tipico: distruggere entità che dovrebbero restare persistenti.

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319150212.png" width="200">
</div>

**Self-call:** Un partecipante può invocare una propria operazione interna. In UniManager, *svc:AppelloService* richiama internamente *validaVincoliIscrizione()* prima di proseguire. Evitiamo eccesso di self-call: solo quelle che chiariscono il comportamento. Quando usarlo: per evidenziare una sotto-fase interna significativa. Errore tipico: dettagliare micro-metodi irrilevanti per lo scenario.

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319150320.png" width="200">
</div>

**Frammento alt:** `alt` modella rami alternativi mutuamente esclusivi. Nel nostro esempio, *svc: AppelloService* sceglie il ramo in base ai prerequisiti dello studente. È la forma naturale per success/failure nello stesso scenario. Quando usarlo: per rami alternativi che non possono coesistere. Errore tipico: guardie sovrapposte o non mutuamente esclusive.

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319150439.png" width="200">
</div>

**Frammento opt:** `opt` esprime un comportamento opzionale (if senza else). Il flusso principale resta valido anche senza il blocco opzionale. Tipico per notifiche, log o passi accessori. Quando usarlo: per comportamenti accessori non sempre presenti. Errore tipico: usare opt per passi in realtà obbligatori.

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319150545.png" width="200">
</div>

**Frammento loop:** `loop` rappresenta una ripetizione controllata da condizione. Qui la stessa operazione si applica a più elementi. Il blocco rappresenta lo stesso micro-flusso ripetuto secondo un criterio esplicito. Quando usarlo: su operazioni ripetute con criterio esplicito. Errore tipico: loop senza condizione chiara o troppo generica.

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319150650.png" width="200">
</div>

**Frammento break:** `break` interrompe il normale scenario in caso eccezionale. Quando il break scatta, il resto del flusso principale non prosegue. Utile per errori critici e stop immediato dell’operazione. Quando usarlo: per uscita anticipata in condizioni eccezionali. Errore tipico: usarlo al posto di un normale ramo alt.

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319150855.png" width="200">
</div>

**Frammento par:** `par` descrive rami concorrenti eseguibili in parallelo. Nel mini-esempio salvataggio e notifica sono indipendenti. Aiuta a discutere latenza e sincronizzazione senza entrare nel codice. Quando usarlo: per attività indipendenti eseguibili in parallelo. Errore tipico: segnare par su rami che in realtà si bloccano.

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319151013.png" width="200">
</div>

**Frammento critical:** `critical` delimita una sezione che non deve essere interrotta da altre interazioni concorrenti. Serve per proteggere aggiornamenti sensibili da interferenze. Lo usiamo per mettere in evidenza vincoli di mutua esclusione. Quando usarlo: su sezioni che richiedono mutua esclusione. Errore tipico: marcare blocchi lunghi senza reale criticità.

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319151148.png" width="200">
</div>

**Frammento  ref:** `ref` rimanda a un’interazione definita altrove, in modo analogo alla chiamata di una funzione. Riduce duplicazione quando lo stesso sotto-scenario ricorre spesso. Mantiene il diagramma principale leggibile e modulare. Quando usarlo: per riuso di sotto-scenari comuni. Errore tipico: nascondere dentro ref passaggi fondamentali del caso.

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319151259.png" width="200">
</div>

**Frammento assert:** `assert` indica una traccia che deve valere obbligatoriamente. Nel diagramma usiamo un blocco etichettato assert per evidenziare il vincolo. È utile per evidenziare protocolli o vincoli forti. Se la traccia reale diverge, siamo fuori specifica. Quando usarlo: per specificare sequenze obbligatorie di protocollo. Errore tipico: usarlo senza una regola verificabile a supporto.

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319151417.png" width="200">
</div>

**Frammento neg:** `neg` rappresenta una traccia vietata. Nel diagramma usiamo un blocco etichettato neg per evidenziare il contro-esempio. Qui evidenziamo un comportamento che il sistema non deve ammettere. È didatticamente utile per discutere contro-esempi. Quando usarlo: per rendere espliciti comportamenti vietati. Errore tipico: confondere neg con un ramo alternativo valido.

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319151532.png" width="200">
</div>

**Frammenti strict vs seq:** 
- `strict`: ordine globale totale tra eventi (no interleaving). 
- `seq`: ordine parziale (weak); preserva solo i vincoli locali. Con `seq` è possibile interleaving tra lifeline diverse. `seq` è il comportamento di default nei diagrammi UML. 

Uso: 
- `strict`: protocolli o scenari con ordine rigidamente definito. 
- `seq`: sistemi concorrenti, per evitare over-constraining. 

Errore tipico: usare strict quando basta seq, irrigidendo il modello.

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319151725.png" width="200">
</div>

###### Activity diagram
Con i Sequence Diagram abbiamo visto chi parla con chi e in quale ordine. Ora cambiamo punto di vista: 
- non più il dettaglio dei messaggi tra oggetti; 
- ma il flusso complessivo del processo. 

L’Activity Diagram è utile per analizzare: 
- passi operativi end-to-end; 
- decisioni e guardie; 
- parallelismi e sincronizzazioni

Uso l’activity diagram: 
- Per descrivere flussi di attività e regole decisionali
- Focus sul  processo, non sui singoli messaggi

**Lettura di base:** 
- Nodo iniziale/finale, azioni e archi
- Decisioni con guardie, parallelismo con fork/join

**Initial e final:** Ogni Activity parte da un nodo iniziale e termina in uno finale. Notazione grafica: 
- nodo iniziale: cerchio pieno (nero); 
- nodo finale: cerchio pieno con bordo esterno (bersaglio).

Questa cornice minima definisce inizio e fine del processo. Anche negli esempi semplici è bene mantenerla esplicita. Quando usarlo: sempre, per delimitare chiaramente il processo. Errore tipico: final implicito che lascia il flusso ambiguo.

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319152240.png" width="200">
</div>

**Action e control flow:** Le azioni sono i passi di lavoro; gli archi definiscono l’ordine. Notazione grafica: 
- azione: rettangolo con angoli arrotondati; 
- control flow: freccia piena tra un’azione e la successiva.

Qui il flusso è lineare: un passo prepara il successivo. È la base su cui poi aggiungiamo decisioni e parallelismo. Quando usarlo: per modellare la sequenza operativa principale. Errore tipico: azioni troppo grandi che nascondono decisioni

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319152431.png" width="200">
</div>

**Decision e merge:** Decision separa i rami; merge li ricompone. Notazione grafica: 
- decision: rombo con un ingresso e più uscite; 
- merge: rombo con più ingressi e una sola uscita. 

Le guardie devono essere chiare e possibilmente mutualmente esclusive. Nel diagramma si vede il classico esito positivo/negativo. Quando usarlo: quando un punto del processo biforca i casi. Errore tipico: dimenticare il merge e lasciare rami scollegati

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319152602.png" width="200">
</div>

**Guardie:** Le guardie esplicitano la condizione logica su ogni ramo. Notazione grafica: 
- ogni uscita dal nodo decision ha una condizione tra parentesi quadre, ad esempio (posti > 0); 
- le guardie devono coprire i casi rilevanti e non sovrapporsi. 

Etichette come (posti > 0) evitano ambiguità interpretative. Sono fondamentali per verificare completezza dei casi. Quando usarlo: per esplicitare condizioni precise sui rami. Errore tipico: guardie vaghe (es. ’ok’) non verificabili.

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319152729.png" width="200">
</div>

**Fork e join:** Fork apre esecuzioni concorrenti; join le sincronizza. Notazione grafica: 
- fork: barra nera spessa con un ingresso e più uscite; 
- join: barra nera spessa con più ingressi e una sola uscita. 

Utile quando due attività sono indipendenti ma devono convergere. Qui salvataggio e notifica procedono in parallelo. Quando usarlo: per parallelismo reale con sincronizzazione finale. Errore tipico: fork senza join quando la convergenza è richiesta.

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319152844.png" width="200">
</div>

**Swimlane:** Le swimlane assegnano responsabilità organizzative ai passi. Notazione grafica:
- il diagramma è diviso in partizioni (corsie) etichettate;
- ogni azione è collocata nella corsia del ruolo/sistema responsabile. 

Si capisce subito chi fa cosa (studente vs sistema). Ottimo strumento di comunicazione con stakeholder non tecnici. Quando usarlo: per chiarire ownership tra ruoli/sistemi. Errore tipico: cambiare lane senza criterio di responsabilità.

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319152956.png" width="200">
</div>

**Send signal:** L’attività emette un segnale verso l’esterno. Notazione grafica: azione di invio evento, con flusso che prosegue senza attesa di risposta immediata. È un punto di integrazione event-driven nel processo. Qui l’evento comunica che gli esiti sono disponibili. Quando usarlo: quando il processo pubblica un evento esterno. Errore tipico: usare send signal per normali chiamate sincrone.

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319153053.png" width="200">
</div>

**Receive signal:** Il flusso resta in attesa di un segnale in ingresso. Notazione grafica: azione di ricezione evento, che abilita la prosecuzione del flusso. La prosecuzione dipende da un evento esterno. Utile per modellare processi reattivi e asincroni. Quando usarlo: per stati in attesa di input/evento esterno. Errore tipico: non specificare chi emette il segnale atteso.

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319153154.png" width="200">
</div>

**Time event:** Il processo reagisce a una scadenza temporale. Notazione grafica: azione di attesa temporale, che abilita una transizione al verificarsi della scadenza. Rende esplicite le transizioni guidate dall’orologio. Nel nostro esempio, il flusso prosegue quando viene raggiunta la scadenza delle iscrizioni. Quando usarlo: per transizioni innescate dal tempo. Errore tipico: modellare tempi con testo libero invece di eventi

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319153249.png" width="200">
</div>

> [!example] Esempio con send, receive e time
> ![[Pasted image 20260319153413.png|200]]
> 
> - Lo scenario combina in un unico flusso invio evento, attesa evento e transizione temporale. 
> - Sequenza logica: pubblicazione esiti, notifica esterna, conferma ricevuta, attesa scadenza, archiviazione finale.
> - Il diagramma chiarisce come eventi esterni e tempo guidano avanzamento e chiusura del processo.

**Loop structure:** Un blocco viene ripetuto finché la condizione resta vera. Qui il ciclo copre tutti gli elementi da elaborare. Esplicitare il loop evita di nascondere lavoro ripetitivo. Quando usarlo: per cicli di processo su insiemi di elementi. Errore tipico: non esplicitare condizione di uscita del ciclo.

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319153508.png" width="200">
</div>

**Interruptable region:** Una regione può essere interrotta da un evento di annullamento. Il diagramma evidenzia il ramo di uscita anticipata. Fondamentale per processi lunghi con comando stop. Quando usarlo: su processi lunghi annullabili dall’utente. Errore tipico: annullamento modellato ma senza cleanup logico.

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260319153603.png" width="200">
</div>

**Confronto finale:**
- Sequence Diagram: 
	- risponde a: chi interagisce con chi, in quale ordine, con quali messaggi; 
	- è ideale per analizzare responsabilità tra componenti e interfacce. 
- Activity Diagram: 
	- risponde a: quali passi compongono il processo, con quali decisioni/parallelismi; 
	- è ideale per descrivere il flusso end-to-end anche a livello di business. 
- Criterio di scelta: 
	- usa Sequence quando il rischio è distribuire male la logica tra oggetti; 
	- usa Activity quando il rischio è perdere la visione complessiva del processo. 

###### Principi di architettura software
Nelle tre lezioni UML precedenti abbiamo usato cinque diagrammi per capire che cosa deve fare il sistema: 
1. Use Case Diagram per funzionalità e attori; 
2. Domain Model  per i concetti del dominio;
3. Class Diagram  per la struttura software; 
4. Sequence Diagram per i messaggi nel tempo; 
5. Activity Diagram per il flusso di processo. 

Questa analisi ci dà una buona fotografia del comportamento e del dominio di UniManager. Il passo successivo è chiederci: come organizziamo il software in modo che questi elementi siano implementati in maniera ordinata, comprensibile e mantenibile?

L’architettura software è il ponte tra: l’analisi (requisiti, casi d’uso, modelli UML), e l’implementazione concreta (classi, moduli, codice). <u>Se l’analisi ci dice cosa deve succedere, l’architettura ci dice come organizzare il sistema perché quelle funzionalità possano essere realizzate in modo coerente</u>. Un’architettura chiara rende più semplice: 
- assegnare il lavoro a più sviluppatori, 
- prevedere l’impatto dei cambiamenti, 
- discutere le scelte progettuali usando una mappa condivisa.

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260325132937.png" width="200">
</div>

> [!info] Parliamo di architettura in questo momento del corso perché: 
> - abbiamo già compreso il dominio di UniManager (gestione di corsi, appelli, iscrizioni e voti), 
> - abbiamo identificato gli attori e i casi d’uso principali di questo dominio (Studente, Docente, Segreteria), 
> - abbiamo un modello concettuale stabile degli elementi centrali del dominio (Studente, Corso, Appello, Iscrizione). 
>  
 È il momento di decidere come strutturare il codice in modo che rispecchi questa comprensione, evitando di arrivare all’implementazione con un insieme di file disorganizzati. Le scelte architetturali prese all’inizio sono difficili da cambiare e condizionano l’evoluzione futura del sistema: per questo è importante renderle esplicite.

**Che cos’è l’architettura software?** Possiamo definire l’architettura software come: 
- la struttura di alto livello del sistema, 
- la suddivisione in componenti principali (moduli, servizi, sottosistemi), 
- le relazioni e le dipendenze tra questi componenti. 

L’architettura non entra nei dettagli di ogni singolo metodo, ma definisce il contenitore in cui questi dettagli andranno a collocarsi. Una buona architettura aiuta chiunque, anche non autore del codice, a farsi un’idea del sistema “a colpo d’occhio”.

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260325133013.png" width="200">
</div>

Possiamo considerare l’architettura software a due livelli: 
1. *in piccolo:* riguarda la struttura interna di un singolo sistema applicativo, come UniManager. Ci chiediamo come suddividere il software in componenti, servizi, moduli. 
2. *in grande:* riguarda sistemi più complessi che coinvolgono molte applicazioni, servizi distribuiti, componenti remoti o cloud. Qui le decisioni architetturali includono protocolli, distribuzione, scalabilità. 

Nel corso ci concentreremo soprattutto sull’architettura in piccolo, perché è la base su cui costruire soluzioni più complesse. Il passaggio all’architettura “in grande” avviene spesso in modo naturale, man mano che un sistema cresce e integra componenti esterni. 

> [!example]  Esempio concreto di architettura “in grande”: 
> - una piattaforma e-commerce come Amazon, con sistemi distinti per catalogo, ordini, pagamenti, spedizioni e raccomandazioni, 
> - applicazioni dedicate per utenti finali e operatori (web, app mobile, app per venditori e logistica), 
> - servizi distribuiti e comunicazione tra sistemi tramite API/eventi.

L’architettura è anche l’insieme delle decisioni strutturali che: 
- sono difficili da cambiare una volta implementate, 
- influenzano molte parti del sistema, 
- impattano su prestazioni, manutenibilità, riusabilità.

> [!example]  Alcuni esempi di decisioni architetturali: 
> - separare o meno la logica di business dall’accesso ai dati, 
> - usare o meno livelli distinti per presentazione, logica, persistenza, 
> - organizzare il codice per funzionalità (Studenti, Corsi, ... ) o per tipo di responsabilità (domain = entità e regole di business, service = logica applicativa/casi d’uso, repository = accesso ai dati e persistenza). 

Queste decisioni definiscono confini e regole di dipendenza:
- chi può chiamare chi, 
- dove può risiedere una regola di business, 
- quali componenti possono conoscere dettagli tecnologici. 

> [!info] In UniManager, per esempio, decidere che la CLI parli solo con i Service e non direttamente con i Repository è una decisione architetturale, non di dettaglio.

**Progettare l’architettura:** significa prendere una serie di decisioni che influenzano profondamente il sistema, per esempio: 
- come suddividere il software in componenti, 
- quale forma far assumere ai confini tra i componenti, 
- come far comunicare i diversi moduli tra loro, 
- quali dipendenze ammettere e quali vietare, 
- quanto isolare il dominio dai dettagli tecnologici. 

Ogni decisione architetturale ha un costo e un beneficio: l’obiettivo è trovare l’equilibrio che rende il sistema comprensibile, efficace e manutenibile. L’architettura emerge quindi non come una singola scelta, ma come un insieme coerente di decisioni che lavorano insieme.

**Architettura vs design di dettaglio:** È utile distinguere l’architettura dal design di dettaglio: 
- l’architettura definisce i grandi blocchi e le loro relazioni, 
- il design di dettaglio decide come questi blocchi sono implementati internamente. 

Le due dimensioni non sono in conflitto: 
- l’architettura crea i confini (chi fa cosa), 
- il design di dettaglio popola questi confini con classi, metodi e algoritmi specifici. 

> [!example]  Esempi tipici di design di dettaglio: 
> - firma dei metodi e scelta dei parametri (iscriviStudente(matricola, idAppello)); 
> - struttura delle classi interne a un modulo e responsabilità dei singoli oggetti; 
> - algoritmo di ordinamento/ricerca adottato in un servizio;
> - gestione puntuale di eccezioni, validazioni e messaggi di errore.

**Strumento per il team:** In un progetto reale, l’architettura non è solo un concetto tecnico, ma un strumento di comunicazione: 
- permette a sviluppatori, progettisti e responsabili di parlare del sistema con un linguaggio comune,
- aiuta a spiegare rapidamente “come è fatto” il software a chi entra nel team. 

Sapere in quale componente collocare una nuova funzionalità riduce discussioni inutili e favorisce un’evoluzione ordinata del codice. Un’architettura condivisa rende più semplice anche la revisione del codice, perché si può valutare se una modifica rispetta o meno i confini stabiliti.

**Se non ci fosse l'architettura:** Se l’architettura non viene pensata o viene lasciata “emergere” senza guida, il rischio è che il sistema si trasformi in una *“big ball of mud”:* 
- file molto grandi e confusi, 
- logica di business mescolata all’accesso ai dati, 
- interfaccia utente che conosce dettagli interni del dominio. 

In questo scenario ogni modifica è potenzialmente pericolosa, nessuno sa esattamente che cosa rompe cambiando una riga e si generano “zone del codice” che nessuno vuole toccare.

**Obiettivi concreti dell'architettura:** Una buona architettura si pone obiettivi molto pratici:
- ridurre la complessità percepita del sistema, 
- facilitare i cambiamenti nel tempo, isolandoli dove possibile, 
- favorire la testabilità, separando ciò che è facilmente verificabile da ciò che dipende dall’ambiente, 
- supportare la crescita del progetto senza doverlo riscrivere da zero ogni anno. 

Nel caso di UniManager, vogliamo che aggiungere una funzionalità (es. un nuovo tipo di report) non significhi toccare codice in dieci posti diversi.

**Architettura e requisiti non funzionali:** L’architettura non serve solo a organizzare il codice: è lo strumento principale per soddisfare molti requisiti non funzionali. Aspetti come: 
- prestazioni (tempi di risposta, carico), 
- affidabilità e tolleranza ai guasti, 
- scalabilità, 
- sicurezza, 
- manutenibilità ed estendibilità, 

dipendono molto più dalle scelte architetturali che dal codice delle singole classi

Un’architettura ben pensata permette di prevedere come il sistema si comporterà sotto stress e come potrà evolvere in futuro. Per questo motivo, le decisioni architetturali vengono prese molto presto nel progetto e orientano molte attività successive.

**Decisioni architetturali e trade-off:** In architettura non esistono scelte “gratis”: quasi ogni decisione porta vantaggi su alcuni aspetti e costi su altri.

> [!example]  Esempi tipici di trade-off: 
> - più livelli e confini $\rightarrow$ maggiore manutenibilità, ma maggiore complessità iniziale; 
> - moduli molto indipendenti $\rightarrow$ migliore evoluzione, ma più interfacce da progettare; 
> - accesso ai dati centralizzato $\rightarrow$ coerenza, ma possibile collo di bottiglia. 
> 

Per questo una decisione architetturale va sempre motivata rispetto a: 
- requisiti prioritari del progetto, 
- vincoli tecnici e organizzativi, 
- costo di implementazione e di evoluzione.  

**Documentare le decisioni architetturali:** Anche in progetti didattici conviene rendere esplicite le decisioni architetturali, evitando che restino solo “nella testa” di chi le ha proposte. Un formato minimo efficace può contenere: 
- *contesto:* quale problema dobbiamo risolvere; 
- *decisione:* quale scelta adottiamo; 
- *alternative:* quali opzioni abbiamo scartato e perché; 
- *conseguenze:* effetti attesi, rischi, impatti. 

Buona pratica operativa: 
- una decisione per documento, titolo chiaro e data; 
- aggiornamento del documento quando la decisione cambia. 

> [!example] Esempio in UniManager:
>  “i comandi CLI invocano solo i Service, senza accesso diretto ai Repository”. Questo riduce ambiguità nel team e facilita revisioni future.

**Tre principi fondamentali:** Per lavorare sull’architettura useremo come filo conduttore tre principi semplici ma molto potenti: 
- *modularità:* dividere il sistema in parti coese, 
- *separazione delle responsabilità:* evitare componenti “tuttofare”, 
- *layering:* organizzare il sistema in livelli con ruoli distinti. 

Questi principi non sono legati a un linguaggio specifico, e valgono anche oltre l’OS, come criterio generale di organizzazione. 

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260325150441.png" width="200">
</div>

**Vocabolario operativo: CLI/UI, Service, Domain, Repository:** Fissiamo i quattro concetti che useremo nel resto del corso: 
- *CLI/UI (Command-Line Interface / User Interface):* riceve input utente e mostra output, senza contenere regole di business; 
- *Service:* coordina i casi d’uso applicativi (sequenza dei passi, validazioni applicative, gestione esiti/errore) e orchestra le interazioni tra Domain e Repository; 
- *Domain:* rappresenta entità e regole del dominio (vincoli, invarianti, significato dei dati), indipendenti da UI e persistenza; 
- *Repository:* si occupa di accesso ai dati e persistenza (lettura/scrittura). 

Regola guida: 
- la *regola di business* vive tra Service e Domain; 
- CLI/UI e Repository sono livelli di interazione, non sede della logica di dominio. 

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260325150627.png" width="200">
</div>

> [!example] Esempio di una decisione in UniManager 
> - Problema: dove collocare la regola “uno studente non può iscriversi due volte allo stesso appello”? 
> - Possibili scelte: 
> 	- nella CLI/UI (scelta errata: vincolo di dominio in livello di presentazione), 
> 	- nel Repository (scelta parziale: verifica tecnica sui dati ma non governo del caso d’uso), 
> 	- in Service + Domain (scelta coerente con le responsabilità architetturali). 
> - Decisione architetturale coerente: 
> 	- la CLI/UI invia la richiesta di iscrizione al Service; 
> 	- il Service applica la regola di business usando concetti e vincoli di Domain; 
> 	- il Repository viene usato solo per leggere/scrivere lo stato persistente; 
> 	- l’esito (OK/errore) torna alla CLI/UI senza esporre dettagli di persistenza.  

**Primo principio- modularità:** La modularità consiste nel dividere il sistema in parti che:
- hanno un ruolo chiaro, 
- comunicano tramite interfacce ben definite, 
- sono ragionabili e testabili quasi indipendentemente. 

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260325150733.png" width="200">
</div>

**Perché la modularità è essenziale?** La modularità riduce la porzione di sistema da analizzare per prendere una decisione progettuale: 
- possiamo concentrarci su un modulo alla volta, 
- il resto del sistema viene visto come un “servizio” con alcune funzioni disponibili. 

Un sistema modulare: 
- si presta meglio al lavoro in team (ogni gruppo segue un modulo), 
- facilita il riuso (un modulo può essere sfruttato in altri contesti), 
- rende più prevedibili gli impatti di ogni modifica. 

> [!example] ![[Pasted image 20260325150841.png|200]]

**Segnali di buona e cattiva modularità:** 
- *Buona modularità:* nomi dei moduli intuitivi e stabili nel tempo, interfacce semplici (pochi metodi, significativi), dipendenze chiare, idealmente in una sola direzione. 
- *Scarsa modularità:* moduli che cambiano continuamente scopo, componenti che conoscono dettagli interni l’uno dell’altro, dipendenze incrociate difficili da districare. 

> [!example] 
> ![[Pasted image 20260325151046.png|200]] 

**Modularità funzionale in UniManager:** Guardando ai requisiti di UniManager, emergono in modo naturale alcune aree funzionali: 
- gestione degli Studenti, 
- gestione dei Corsi, 
- gestione degli Appelli e delle Iscrizioni, 
- gestione della persistenza delle informazioni. 

Possiamo quindi immaginare di avere, almeno concettualmente, un modulo per ciascuna di queste aree, che si occupa di quella parte del dominio. Questa è una prima forma di modularità, che poi andremo a raffinare con il principio della separazione delle responsabilità.

> [!example] ![[Pasted image 20260325151146.png|200]]

**Secondo principio- separazione delle responsabilità:** La separazione delle responsabilità invita a evitare componenti che fanno “un po’ di tutto”. Ogni parte del sistema dovrebbe avere una responsabilità principale ben riconoscibile: 
- rappresentare dati del dominio, 
- applicare regole di business, 
- comunicare con l’utente, 
- accedere al database. 

Quando questi ruoli si mescolano, il codice diventa più difficile da capire, da testare e da modificare. 

> [!example] ![[Pasted image 20260325151234.png|200]] 

**Perché separare le responsabilità:** Una separazione chiara delle responsabilità offre diversi vantaggi: 
- rende più semplice capire “a chi compete” un certo cambiamento, 
- riduce il rischio di introdurre effetti collaterali non desiderati, 
- permette di sostituire una parte (es. la persistenza) mantenendo le altre stabili. 

Al contrario, una classe o un modulo che gestisce: 
- logica del dominio, 
- interfaccia utente, 
- accesso ai file

è destinato a diventare un punto critico del sistema.  

**Tipi di responsabilità in UniManager:**
In questa slide adottiamo una vista per responsabilità tecniche (non per aree funzionali): 
- è una prospettiva diversa da “modularità funzionale” (Studenti, Corsi, Appelli), 
- serve a chiarire chi fa cosa tra i componenti tecnici. 

In UniManager possiamo individuare almeno quattro tipi di responsabilità: 
- dominio: rappresentare entità e regole del mondo universitario, 
- logica applicativa: orchestrare i casi d’uso, 
- persistenza: leggere e scrivere dati, 
- interfaccia: interagire con l’utente. 

Separare queste responsabilità ci aiuta a progettare componenti più puliti e facili da testare. 

> [!example] ![[Pasted image 20260325151358.png|200]]

**Separazione in pratica- entità, servizi, repository, CLI:** Possiamo tradurre queste responsabilità in quattro famiglie di componenti: 
- *entità del dominio* (es. Studente, Corso, Appello), che contengono dati e regole del dominio; 
- *servizi applicativ*i (es. AppelloService), che implementano le operazioni richieste dai casi d’uso; 
- *repository* (es. StudenteRepository), che si occupano di salvare e recuperare i dati; 
- *componenti CLI/UI*, che gestiscono input e output verso l’utente. 

In questo modo, ogni famiglia ha una responsabilità precisa e riconoscibile.  

**Terzo principio- layering:** Il terzo principio riguarda il layering, cioè l’idea di organizzare il sistema in livelli di astrazione:
- un livello più vicino all’utente (interfaccia), 
- livelli intermedi di logica, 
- un livello più vicino alle tecnologie di base (file, database, rete). 

Ogni livello offre servizi al livello superiore e sfrutta i servizi del livello inferiore, secondo regole di dipendenza ben precise. 

> [!example] ![[Pasted image 20260325151530.png|200]]

**Perché usare i layer:** I layer servono a tenere sotto controllo le dipendenze: 
- il codice dell’interfaccia utente non dovrebbe conoscere i dettagli del database, 
- la logica di business non dovrebbe dipendere da come i dati sono memorizzati. 

I benefici principali sono: 
- maggiore chiarezza nella struttura, 
- confini più netti per i test, 
- facilità nel sostituire un livello (es. cambiare DB) mantenendo stabili gli altri.

> [!example] ![[Pasted image 20260325151648.png|200]]
 
**Un layering semplice per UniManager:** Per UniManager possiamo immaginare quattro livelli principali:
- *CLI/UI:* comandi testuali o schermate che l’utente utilizza; 
- *Service:* implementazione dei casi d’uso; 
- *Domain:* entità e regole del dominio universitario; 
- *Repository:* accesso ai dati su file o database. 

Questi livelli formano una “pila” dove ciascuno ha un ruolo distinto. 

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260325151828.png" width="200">
</div>

**Regole di comunicazione tra i layer:** Per mantenere l’architettura pulita, è utile stabilire alcune regole: 
- la CLI/UI parla solo con il livello Service, 
- i Service parlano con il Domain e i Repository, 
- i Repository non devono dipendere dall’interfaccia utente, 
- il Domain non dovrebbe conoscere dettagli di persistenza o di interfaccia. 

In questo modo, le dipendenze hanno una direzione chiara e il dominio rimane al centro dell’architettura. 

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260325151933.png" width="200">
</div>

**Checklist rapida di qualità architetturale:** Domande operative da porre al modello architetturale: 
- ogni componente ha una responsabilità principale chiara? 
- le dipendenze rispettano i confini tra layer? 
- una modifica tipica richiede tocchi localizzati o diffusi? 
- il dominio resta isolato da UI e persistenza? 
- i test possono concentrarsi su un livello alla volta? 

Se molte risposte sono negative, l’architettura va rifinita prima di ampliare il codice.  

**UniManager- una prima suddivisione in pacchetti:** Possiamo tradurre questi livelli in una prima proposta di struttura del progetto: 
- cli, per l’interfaccia testuale o grafica, 
- service, per i servizi applicativi, 
- domain, per le entità e le regole del dominio, 
- repository, per l’accesso ai dati.

Questa suddivisione ha il vantaggio di: 
- rendere subito visibile la separazione delle responsabilità, 
- dare una mappa chiara del sistema fin dai primi commit. 

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260325152042.png" width="200">
</div>

> [!example]  Iscrizione a un appello attraverso i layer. Immaginiamo l’operazione “Studente si iscrive a un appello”:
> - la CLI raccoglie l’input (matricola, appello) e invoca il servizio opportuno; 
> - il Service verifica le regole (esistenza di studente e appello, vincoli sulle iscrizioni) usando il Domain; 
> - il Service chiede al Repository di aggiornare la persistenza; 
> - il risultato torna verso la CLI, che mostra all’utente conferma o eventuali errori. 
>
>La responsabilità di ciascun passo è chiara e localizzata nel layer corretto. 
> ![[Pasted image 20260325152144.png|200]]
###### Stili architetturali
Nella lezione precedente abbiamo visto: 
- che cos'è l’architettura software, 
- perché è importante pensarla esplicitamente, 
- tre principi chiave: modularità, separazione delle responsabilità, layering. 

Oggi facciamo un passo avanti: vediamo come questi principi possono essere combinati in forme ricorrenti, cioè in stili architetturali. Obiettivo: dare un vocabolario per riconoscere le strutture tipiche dei sistemi e per discutere le scelte senza basarsi solo su tecnologie o mode del momento.  

**Che cos'è uno stile architetturale?** Uno stile architetturale descrive:
- un insieme di componenti tipici, 
- le regole di organizzazione (es. livelli, ruoli), 
- le modalità di comunicazione tra i componenti, 
- i vincoli sulle dipendenze. 

Non è una tecnologia o un framework, è un modello concettuale che può essere implementato in molti modi diversi. Serve per ragionare in termini di forme globali del sistema, non di singole classi. 

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260326150435.png" width="200">
</div>

**Perché esistono più stili?** I sistemi software devono soddisfare esigenze molto diverse: interattività, prestazioni, scalabilità, affidabilità, integrazione con altri sistemi, ... Ogni stile è nato per rispondere meglio a alcuni di questi bisogni. Non esiste uno stile “migliore in assoluto”: ogni scelta comporta vantaggi e svantaggi, capire questi compromessi è parte del mestiere dell’ingegnere del software.  

In questa lezione ci concentriamo su quattro stili molto diffusi: 
1. *Layered:* organizzazione a livelli (es. presentazione, logica, dati). 
2. *Client–server:* separazione tra chi chiede un servizio e chi lo fornisce. 
3. *Repository:* più componenti che condividono una base dati comune. 
4. *Service-based:* sistema composto da servizi indipendenti che espongono interfacce.  

**Architettura layered:** Lo stile layered organizza il sistema in livelli di astrazione: 
- ogni livello ha un ruolo specifico, 
- i livelli comunicano tra loro secondo regole definite. 

Tipicamente si distinguono livelli come: 
- presentazione verso l’utente, 
- logica di business / servizi applicativi, 
- accesso ai dati e infrastruttura di base.

L’idea chiave è tenere separati i diversi tipi di responsabilità e controllare le dipendenze tra le parti del sistema. In pratica, layered risponde bene alla domanda: dove deve vivere una certa logica, per evitare che UI, regole di business e persistenza si mescolino. Non impone una tecnologia specifica: è uno schema organizzativo che può essere applicato in modo simile in stack diversi. 

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260326150612.png" width="200">
</div>

Ogni layer offre servizi a quello sopra e usa i servizi di quello sotto. Le dipendenze seguono una direzione prevalente: dall’alto verso il basso. 

**Layered- vantaggi:**
1. *Chiarezza strutturale:* 
	- è facile spiegare “come è fatto” il sistema, 
	- i confini tra responsabilità sono visibili. 
2. *Facilità di manutenzione:*
	- cambiare il database può richiedere modifiche solo nel layer di accesso ai dati, 
	- cambiare UI dovrebbe toccare principalmente il layer di presentazione. 
3. *Testabilità:* 
	- si possono testare i servizi senza coinvolgere né UI né DB reali (mock). 
4. *Evoluzione controllata:* 
	- una modifica in un layer tende a propagarsi meno agli altri, se le interfacce restano stabili. 
5. *Onboarding del team:* 
	- la struttura a livelli aiuta a capire rapidamente dove collocare nuova logica.  

**Layered: limiti e compromessi:** 
1. *Possibile overhead:* 
	- le chiamate attraversano più layer, con qualche costo in termini di prestazioni. 
2. *Rischio di rigidità:* 
	- se i confini tra layer sono troppo rigidi, diventa difficile gestire casi “trasversali”. 
3. *Necessità di disciplina:* 
	- se i layer iniziano a “bucarsi” (UI che accede direttamente al DB), la struttura si degrada in fretta.  
4. *Rischio di stratificazione artificiale:* 
	- in sistemi piccoli, troppi layer possono introdurre complessità non necessaria. 
5. *Possibile duplicazione di modelli:*  
	- spesso gli stessi dati attraversano più layer con trasformazioni ripetitive.  

> [!info] Per UniManager possiamo riprendere il layering già introdotto: 
> - CLI / UI: interazione con l’utente (comandi, menu). 
> - Service: implementazione dei casi d’uso (iscrizione, creazione appelli, ... ). 
> - Domain: entità e regole del dominio universitario. 
> - Repository: accesso ai dati su file o database.
> 
> ![[Pasted image 20260326150653.png|200]]

> [!example] Iscrizione a un appello nei layer 
> - Scenario «Studente si iscrive a un appello»: 
> 	1. CLI/UI chiede matricola e appello e invoca AppelloService. 
> 	2. Service controlla i vincoli usando oggetti del Domain. 
> 	3. Service chiede al Repository di aggiornare la persistenza. 
> 	4. CLI/UI mostra il risultato (successo o errore). 
> - Il flusso attraversa i layer in modo ordinato, con responsabilità chiare
> ![[Pasted image 20260326150719.png|200]]

**Architettura client–server:** Lo stile client–server separa: 
- il client, che richiede servizi, 
- il server, che fornisce servizi. 

Tipicamente: 
- il server è centrale, mantiene dati e logica, 
- i client sono distribuiti, interagiscono con il server via rete. 

È il modello concettuale alla base di molte applicazioni web e sistemi distribuiti. La separazione client/server rende esplicito un confine importante: l’esperienza utente sta lato client, mentre regole e dati condivisi stanno lato server. Questo stile è spesso il primo passo quando un sistema locale deve diventare multiutente.  

<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Pasted image 20260326150805.png" style="width: 100%; border-radius: 8px;">
  </div>
  <div style="flex: 1.5;">
  I client inviano richieste (request), il server risponde (response).  
  </div>
</div>

**Client–serverà- vantaggi:** 
- *Centralizzazione:* 
	- dati e logica critici restano concentrati sul server, 
	- più facile aggiornare e controllare il sistema. 
- *Multi-client:* 
	- lo stesso server può servire molti client diversi (desktop, web, mobile). 
- *Sicurezza:* 
	- le regole di accesso ai dati sono gestite lato server. 
- *Governance semplificata:* 
	- aggiornamenti di logica e policy avvengono centralmente, senza distribuire patch su ogni client. 
- *Riuso dei servizi:* 
	- più canali (web, mobile, desktop) possono sfruttare lo stesso backend.  

**Client–server- limiti e compromessi:** 
- *Single point of failure:* 
	- se il server è giù, tutti i client sono bloccati. 
- *Scalabilità:*
	- se il numero di client cresce molto, il server può diventare un collo di bottiglia. 
- *Complessità di rete:* 
	- bisogna gestire latenza, errori di comunicazione, autenticazione, ... 
- *Dipendenza dalla connettività:* 
	- in assenza di rete, il client può perdere funzionalità critiche. 
- *Scelte di stato/sessione:* 
	- la gestione dello stato lato server o lato client impatta scalabilità e semplicità.  

Al momento, UniManager è pensato come applicazione locale (es. CLI). Potremmo però immaginare un’evoluzione: 
- un server UniManager che espone servizi (API), 
- diversi client: 
	- applicazione web per gli studenti, 
	- applicazione interna per la segreteria, 
	- app mobile per consultare il libretto. 
	
In questo scenario, lo stile client–server diventa il modello concettuale di riferimento.  

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260326151002.png" width="200">
</div>

**Architettura repository:** Nello stile repository, più componenti condividono un’unica base dati comune. Il repository è il “cuore” dove risiedono le informazioni condivise. Intorno ad esso ci sono componenti che: 
- elaborano dati, 
- offrono funzionalità, 
- reagiscono ai cambiamenti nel repository. 

 La priorità di questo stile è la coerenza informativa: più moduli diversi lavorano sullo stesso patrimonio dati. È adatto quando il valore principale del sistema sta nei dati condivisi più che nel comportamento di un singolo componente.  

**Repository:** Nello stile repository, tutti i componenti ruotano attorno a un’unica base dati condivisa. Ogni componente legge e aggiorna il repository, ma rimane indipendente dagli altri. 

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260326151052.png" width="200">
</div>

**Repository- vantaggi:** 
- *Consistenza dei dati:* 
	- una sola fonte di verità per le informazioni condivise. 
- *Condivisione:* 
	- sottosistemi diversi possono cooperare accedendo agli stessi dati. 
- *Evoluzione dei componenti:* 
	- è possibile sostituire o aggiornare i componenti attorno al repository mantenendo stabile la base dati (entro certi limiti). 
- *Tracciabilità:* 
	- avere un punto dati centrale facilita audit, report e analisi trasversali. 
- *Integrazione:* 
	- nuovi componenti possono essere aggiunti collegandosi al repository condiviso.  

**Repository- limiti e compromessi:** 
- *Punto critico di fallimento:* 
	- se il repository si blocca o viene corrotto, l’intero sistema è in difficoltà. 
- *Scalabilità e prestazioni:* 
	- molte richieste concorrenti possono creare colli di bottiglia. 
- *Dipendenza forte dal modello dati:* 
	- cambiare lo schema del repository può impattare molti componenti. 
- *Accoppiamento semantico:* 
	- componenti nominalmente separati possono diventare dipendenti dagli stessi dettagli di schema. 
- *Migrazioni delicate:* 
	- evolvere il repository richiede coordinamento tra molti moduli e finestre di rilascio ben pianificate.  

In UniManager abbiamo concettualmente una base dati universitaria: studenti, corsi, appelli, iscrizioni, voti, ... Per ora la implementiamo in modo semplice (file / strutture in memoria), ma concettualmente è un repository informativo condiviso. In una versione più evoluta: servizi diversi (es. gestione esami, reportistica, analisi statistiche) potrebbero accedere allo stesso repository.  

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260326151202.png" width="200">
</div>

**Architettura service-based:** Nello stile service-based, il sistema è costituito da servizi: 
- ciascun servizio è un componente indipendente, 
- espone una interfaccia ben definita, 
- può essere sviluppato, distribuito e aggiornato separatamente. 

I servizi comunicano tra loro tramite chiamate di rete (es. API). Ogni servizio presidia una responsabilità funzionale chiara: questo riduce l’impatto dei cambiamenti e favorisce evoluzione incrementale. L’attenzione si sposta dai confini tra classi ai contratti tra servizi (endpoint, payload, versioni).  

**Da client–server a service-based:** Lo stile service-based può essere visto come una evoluzione del client–server: 
- invece di un unico server “monolitico”, abbiamo più servizi specializzati, 
- ogni servizio può avere i propri dati e responsabilità.

> [!example] In UniManager:
>  - servizio StudenteService per la gestione studenti, 
>  - servizio AppelloService per appelli e voti, 
>  - servizio ReportingService per report e statistiche.  

<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Pasted image 20260326151415.png" style="width: 100%; border-radius: 8px;">
  </div>
  <div style="flex: 1.5;">
  <b>Service-based:</b>
  <li> Il client/UI parla con più servizi indipendenti, ciascuno con una responsabilità ben definita.</li>
  <li> I servizi possono anche collaborare tra loro tramite chiamate di rete (linee tratteggiate).</li>
  </div>
</div>

**Service-based- vantaggi:** 
- *Scalabilità:* 
	- ogni servizio può essere scalato indipendentemente (più istanze del solo servizio critico). 
- *Autonomia dei team:* 
	- team diversi possono sviluppare e rilasciare servizi separati. 
- *Evoluzione tecnologica:* 
	- servizi diversi possono usare tecnologie diverse, purché rispettino i contratti di interfaccia. 
- *Isolamento dei guasti:* 
	- un problema in un servizio può avere impatto limitato, se i fallback sono progettati bene. 
- *Rilasci mirati:* 
	- è possibile aggiornare un singolo servizio senza rilasciare l’intero sistema.  

**Service-based- limiti e compromessi:** 
- *Maggiore complessità complessiva:* 
	- bisogna gestire comunicazione tra servizi, fallimenti parziali, versioni delle API. 
- *Osservabilità e debugging più difficili:* 
	- un flusso di richiesta può attraversare molti servizi diversi. 
- *Coerenza dei dati:* 
	- se ogni servizio ha il proprio storage, garantire consistenza diventa più complesso. 
- *Transazioni distribuite:* 
	- operazioni che coinvolgono più servizi richiedono sequenze coordinate e azioni di annullamento in caso di errore. 
- *Costo operativo:* 
	- servono strumenti maturi per monitoraggio, tracing, deployment e gestione API.  


<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Pasted image 20260326151541.png" style="width: 100%; border-radius: 8px;">
  </div>
  <div style="flex: 1.5;">
  Per il corso manteniamo UniManager semplice e monolitico (layered). Ma possiamo usare il modello service-based per: 
<li> discutere come dividere il sistema in aree funzionali indipendenti, </li>
<li> ragionare su possibili evoluzioni future (es. servizi separati per studenti, esami, analytics), </li>
<li> introdurre l’idea di contratti tra componenti (interfacce ben definite).</li>
  </div>
</div>

**Stili a confronto** Ogni stile risponde meglio a una domanda diversa: 
- *Layered:* come separare responsabilità interne e tenere ordinate le dipendenze? 
- *Client–server:* come distribuire il sistema tra client e server condiviso? 
- *Repository:* come far cooperare componenti diversi sugli stessi dati? 
- *Service-based:* come far evolvere parti del sistema in modo indipendente? 

Sintesi operativa: 
- priorità a manutenibilità interna e chiarezza: *Layered*; 
- priorità a multiutente via rete: *Client–server*; 
- priorità a consistenza e patrimonio dati comune: *Repository*; 
- priorità a scalabilità organizzativa e rilasci indipendenti: *Service-based*. 

Nei sistemi reali la scelta è spesso ibrida, ma i confini devono restare espliciti.  

Ogni stile ha impatti diversi su: 
- *prestazioni*, 
- *scalabilità*, 
- *affidabilità*, 
- *manutenibilità*,
- *evolutività*. 

> [!example] 
> - layered aiuta molto la manutenibilità e la testabilità; 
> - client–server e service-based introducono nuove sfide di prestazioni e affidabilità di rete; 
> - repository facilita la consistenza, ma crea possibili colli di bottiglia.  

**Come scegliere uno stile (o più di uno):** Nella pratica, i sistemi reali sono spesso una combinazione di stili:
- un’applicazione web può essere client–server e internamente layered, 
- più servizi possono condividere un repository o avere ognuno il proprio. 

Alcune domande utili: 
- Quanto è importante la manutenibilità interna? (layered) 
- Il sistema è distribuito o locale? (client–server / service-based) 
- I dati condivisi sono il vero cuore del sistema? (repository) 
- Ci sono team diversi che devono lavorare in modo indipendente? (service-based) 

Sintesi per UniManager in questa fase: 
- base layered con repository centrale dei dati,
- client–server e service-based come possibili evoluzioni successive. 

###### Principi SOLID
Finora ci siamo concentrati su: 
- principi architetturali generali (modularità, separazione delle responsabilità, layering), 
- stili architetturali (layered, client–server, repository, service-based). 

Oggi scendiamo di livello: come progettare le classi e le loro dipendenze, in modo coerente con l’architettura scelta. L'obiettivo è quello di collegare le decisioni locali (sulle classi) alla struttura globale del sistema

**Perché servono principi di design?** Il codice non deve solo funzionare oggi: deve poter evolvere senza degradarsi. Quando il progetto cresce, emergono problemi ricorrenti: 
- classi con troppe responsabilità e confini poco chiari; 
- dipendenze rigide che propagano gli effetti di ogni modifica; 
- regressioni frequenti e test costosi da mantenere. 

In queste condizioni, il costo principale non è scrivere nuova logica, ma capire e cambiare quella esistente. I principi di design (come *SOLID*) servono come criteri operativi per ridurre accoppiamento, aumentare coesione e migliorare testabilità. 

**Che cos’è SOLID?** SOLID è un acronimo che indica cinque principi: 
1. **Single Responsibility Principle (SRP):** principio di responsabilità singola 
2. **Open/Closed Principle (OCP):** principio aperto/chiuso 
3. **Liskov Substitution Principle (LSP):** principio di sostituzione di Liskov 
4. **Interface Segregation Principle (ISP):** principio di segregazione delle interfacce 
5. **Dependency Inversion Principle (DIP):** principio di inversione delle dipendenze 

Non sono leggi assolute, ma linee guida per prendere decisioni migliori. 

**SRP — Single Responsibility Principle:** Una classe dovrebbe avere un solo motivo per cambiare. Per “motivo per cambiare” si intende una famiglia omogenea di modifiche (stesso tipo di responsabilità). Quando una classe mescola responsabilità diverse: 
- aumenta l’accoppiamento interno; 
- crescono effetti collaterali e regressioni; 
- test e manutenzione diventano più costosi. 

> [!example] Esempio in UniManager: 
> - una classe che contiene insieme regole d’esame, persistenza e messaggi UI viola SRP; 
> - queste responsabilità vanno separate in componenti distinti

In UniManager conviene distinguere almeno queste responsabilità: 
- business: regole del dominio (vincoli su appelli, prerequisiti, esiti); 
- persistenza: lettura/scrittura su file o database; 
- presentazione: messaggi verso CLI/UI e formati di output; 
- cross-cutting: log, audit, monitoraggio. 

Criterio pratico per rilevare una violazione SRP: “quante ragioni indipendenti ho per modificare questa classe?”. Se la risposta è più di una, serve separare. L'obiettivo è quello di avere classi piccole ma coese, con confini chiari e test più mirati.  

<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Pasted image 20260401162605.png" style="width: 100%; border-radius: 8px;">
  </div>
  <div style="flex: 1.5;">
  GestoreEsami mescola: 
	  <li> business (regole sugli appelli), </li>
	  <li> persistenza (salvataggio su file), </li>
	  <li> presentazione (stampa su console). </li>
  Cambia il formato di salvataggio? ⇒ modifico GestoreEsami. Cambia il modo di comunicare con l’utente? ⇒ modifico ancora GestoreEsami. Risultato: classe fragile, difficile da testare e da capire.
  </div>
</div>

<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Pasted image 20260401162702.png" style="width: 100%; border-radius: 8px;">
  </div>
  <div style="flex: 1.5;">
  <b>Prima: </b>
  <li>UI → GestoreEsami (regole + salvataggio su file + stampa).</li>
  <li> Un’unica classe con responsabilità miste (business, persistenza, presentazione). </li>
  <b>Dopo: </b>
  <li> UI → ServizioEsami (solo regole di business), </li>
  <li> ServizioEsami delega a ArchivioEsamiSuFile la persistenza, </li>
  <li> la UI (CLI) si occupa dei messaggi verso l’utente (stampe, errori, conferme). </li>
  <b>Effetti:</b> 
  <li> ogni classe ha un motivo principale per cambiare, </li>
  <li> il comportamento si capisce per “blocchi” separati, </li>
  <li> il design è pronto per introdurre OCP e DIP.</li>
  </div>
</div>

**OCP — Open/Closed Principle:** Il software dovrebbe essere aperto all’estensione, ma chiuso alla modifica. Interpretazione operativa: 
- nuovo comportamento = aggiungere nuove classi/composizioni; 
- codice stabile = evitare modifiche continue agli stessi metodi centrali. 

OCP non significa “mai modificare nulla”, ma ridurre le modifiche invasive nelle parti già consolidate.

Esempio in UniManager: aggiungere un nuovo criterio di ordinamento senza riscrivere ServizioEsami.  

OCP nasce a causa di un problema ricorrente:  il requisito resta lo stesso, ma cambiano le varianti (nuovi casi, nuove regole, nuove modalità). Se ogni variante richiede di riaprire lo stesso componente centrale: 
- aumenta il rischio di regressioni; 
- i test da rieseguire crescono a ogni modifica; 
- quel componente diventa un collo di bottiglia del team. 

OCP nasce per gestire la variabilità in modo più stabile: 
- la variabilità va isolata in estensioni dedicate; 
- il nucleo resta stabile e riusabile. 

<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Pasted image 20260401163352.png" style="width: 100%; border-radius: 8px;">
  </div>
  <div style="flex: 1.5;">
Metodo "monolitico" con tre comportamenti diversi. Ogni nuovo criterio richiede di: 
<li>modificare ServizioEsami,</li>
<li>aggiungere un nuovo blocco if/else. </li>
Questo viola OCP: 
<li>il codice non è chiuso alla modifica,</li>
<li>cresce nel tempo e diventa fragile.</li>
La logica di ordinamento è tutta concentrata qui, difficile da testare in isolamento. 
  </div>
</div>

Perché questa versione viola OCP? Ogni nuovo criterio di ordinamento richiede: 
- di modificare ServizioEsami, 
- di aggiungere un altro blocco if/else, 
- di rischiare di rompere casi già funzionanti. 

Criterio diagnostico: se una nuova variante obbliga a riaprire sempre lo stesso metodo, il design non è OCP. Il metodo cresce continuamente: code smell di variabilità concentrata. 

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260401163939.png" width="200">
</div>

Separiamo il criterio di ordinamento dal servizio. OrdinatoreAppelli descrive l’operazione in astratto. Ogni criterio (per data, per corso, per iscritti, . . . ) diventa: 
- una classe separata, 
- aggiungibile senza modificare codice esistente. 

ServizioEsami non cambia più:
- riceve un ordinatore dall’esterno, 
- lo usa senza sapere come funziona. 

Questo rispetta pienamente OCP: 
- aperto a nuovi criteri: basta aggiungere una classe; 
- chiuso alle modifiche: nessun if/else da aggiornare.  

Condizioni operative: La soluzione resta OCP solo se il punto di estensione è stabile: 
- interfaccia chiara (OrdinatoreAppelli), 
- servizio dipendente dall’astrazione, non dai criteri concreti. 

Per introdurre una nuova variante: 
- si aggiunge una classe (OrdinaPer...), 
- si seleziona il criterio in composizione/iniezione, 
- ServizioEsami resta invariato. 

Checklist rapida OCP: 
- una nuova variante richiede solo codice nuovo? 
- i test esistenti restano validi senza riscritture massive? 
- il comportamento variabile è confinato in componenti dedicati?  

**LSP — Liskov Substitution Principle:** Le sottoclassi dovrebbero poter sostituire le superclassi senza cambiare il comportamento atteso del programma. In pratica: 
- il codice client deve poter usare la sottoclasse al posto della superclasse senza accorgersene; 
- stessa chiamata, stesso significato, stesso tipo di esito. 

Segnali tipici di violazione: 
- la sottoclasse rifiuta casi che la superclasse accetta;
- la sottoclasse introduce errori/eccezioni inattese; 
- il client deve aggiungere controlli sul tipo concreto per evitare bug.  

In molti sistemi OO usiamo ereditarietà e polimorfismo per riuso. Senza LSP, il polimorfismo diventa fragile: 
- il chiamante si aspetta un comportamento della superclasse, 
- alcune sottoclassi rompono quelle attese. 

Effetti pratici: 
- if speciali per distinguere sottoclassi, 
- bug difficili da vedere in revisione, 
- test incompleti perché dipendono dal tipo concreto.  

<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Pasted image 20260401164427.png" style="width: 100%; border-radius: 8px;">
  </div>
  <div style="flex: 1.5;">
La superclasse promette un contratto semplice: 
<li>input valido, </li>
<li>output booleano. </li>
La sottoclasse introduce un comportamento incompatibile: lancia eccezioni inattese. Un chiamante polimorfico può rompersi pur rispettando il contratto della base.  
  </div>
</div>

Perché questa versione viola LSP? La sottoclasse VerificatoreSoloInCorso restringe il contratto: 
- aggiunge una precondizione implicita (appello aperto), 
- cambia la modalità di errore (eccezione invece di false). 

Chi usa il tipo base non è più protetto: serve conoscere il tipo concreto per evitare crash. Se serve un comportamento diverso, è meglio modellare una strategia, non una sottoclasse incompatibile. 

<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Pasted image 20260401164800.png" style="width: 100%; border-radius: 8px;">
  </div>
  <div style="flex: 1.5;">
Tutte le implementazioni rispettano lo stesso contratto: 
<li>nessuna eccezione inattesa, </li>
<li>risultato sempre booleano. </li>
Le varianti cambiano la regola, non il tipo di comportamento osservabile. Il chiamante resta stabile e sostituibile. 
  </div>
</div>

La soluzione resta LSP solo se il contratto è uniforme: 
- stessi input ammessi (niente precondizioni più forti nelle sottoclassi), 
- stesso tipo di output/errore osservabile dal client. 

Le varianti devono cambiare la regola, non il significato del comportamento atteso. Checklist rapida LSP: 
- il client funziona uguale con qualunque implementazione? 
- servono if sul tipo concreto per evitare crash? 
- una sottoclasse introduce eccezioni o vincoli non previsti?  

**ISP — Interface Segregation Principle:** È meglio avere più interfacce piccole e specifiche che un’unica interfaccia “gigante” che fa di tutto. In pratica: 
- ogni client dovrebbe vedere solo le operazioni del proprio ruolo; 
- le funzionalità non pertinenti vanno in interfacce separate. 

Violazioni tipiche: 
- interfacce con molti metodi opzionali, 
- classi client che devono implementare metodi “vuoti” o fittizi.  

Problema ricorrente: 
- un’unica interfaccia viene usata da ruoli diversi (studente, docente, segreteria), 
- col tempo accumula metodi eterogenei. 

Effetti pratici: 
- i client dipendono da operazioni che non usano; 
- compaiono metodi vuoti o UnsupportedOperationException; 
- una modifica locale impatta molti componenti non coinvolti. 

ISP nasce per separare le interfacce per ruolo/scenario e ridurre l’accoppiamento inutile.  

<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Pasted image 20260401165152.png" style="width: 100%; border-radius: 8px;">
  </div>
  <div style="flex: 1.5;">
ClientStudente dipende da metodi solo docenti/segreteria. Implementa metodi non pertinenti con eccezioni fittizie. Ogni cambiamento dell’interfaccia impatta tutti i client.
  </div>
</div>

La violazione nasce da un problema di modellazione: PortaleDidattico mescola responsabilità di ruoli diversi in un unico contratto. Effetto sui client: 
- ogni client dipende anche da metodi che non dovrebbe conoscere; 
- compaiono implementazioni fittizie (UnsupportedOperationException). 

Conseguenza pratica: modifiche pensate per un ruolo impattano anche client di altri ruoli. Correzione: separare l’interfaccia per responsabilità funzionali (Studente, Docente, Segreteria).  

<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Pasted image 20260401165415.png" style="width: 100%; border-radius: 8px;">
  </div>
  <div style="flex: 1.5;">
Ogni client dipende solo da ciò che usa davvero.  I contratti sono più piccoli e chiari.  Le evoluzioni su un ruolo non forzano modifiche sugli altri.
  </div>
</div> 

La soluzione resta ISP se le interfacce sono coerenti con i ruoli: 
- interfacce distinte per Studente, Docente, Segreteria; 
- ogni client implementa solo metodi pertinenti. 

Quando aggiungi una funzionalità: 
- modifichi solo l’interfaccia del ruolo coinvolto; 
- eviti impatti trasversali sui client non interessati. 

Checklist rapida ISP: 
- compaiono metodi vuoti o UnsupportedOperationException? 
- un client dipende da operazioni che non usa? 
- una modifica locale forza cambi in client di altri ruoli

**DIP — Dependency Inversion Principle:** I moduli di alto livello e quelli di basso livello dovrebbero dipendere da astrazioni; i dettagli concreti dovrebbero dipendere da tali astrazioni. In pratica: 
- alto livello = logica applicativa/business, più astratta (es. ServizioEsami); 
- basso livello = dettagli tecnici, più concreti (es. ArchivioEsamiSuFile); 
- entrambi si collegano tramite un’astrazione comune (es. ArchivioEsami). 

Effetto: il business resta stabile anche se cambiano tecnologie e componenti infrastrutturali.  

Problema ricorrente: 
- i servizi applicativi creano direttamente classi infrastrutturali (new Archivio..., new ClientHttp...); 
- quindi il livello alto dipende direttamente dal livello basso. 

Effetti pratici: 
- cambiare tecnologia richiede modifiche nel business; 
- i test diventano lenti e fragili perché dipendono da risorse reali; 
- cresce l’accoppiamento tra livelli che dovrebbero restare separati. 

DIP nasce per invertire questa dipendenza: 
- il business dipende da interfacce stabili, 
- i dettagli concreti si agganciano a quelle interfacce.  

Mappa dei ruoli (in UniManager): 
- alto livello: ServizioEsami (regole applicative); 
- basso livello: ArchivioEsamiSuFile (I/O su file); 
- astrazione: ArchivioEsami. 

Regola DIP: ServizioEsami deve dipendere da ArchivioEsami, non da ArchivioEsamiSuFile. Se dipende dal dettaglio concreto: cambiare persistenza (file, DB, memoria) richiede modifiche nel servizio. DIP mantiene pulito il confine tra business e infrastruttura.  

<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Pasted image 20260401170554.png" style="width: 100%; border-radius: 8px;">
  </div>
  <div style="flex: 1.5;">
ServizioEsami dipende da ArchivioEsamiSuFile, una classe concreta. La scelta del file ("esami.txt") è codificata dentro il servizio. Problemi: 
<li>per usare un database devo modificare ServizioEsami;</li>
<li>per fare test devo scrivere davvero sul file;</li>
<li>la logica di business è legata a un dettaglio di infrastruttura.</li>
Qui DIP è violato: un modulo di alto livello dipende da un modulo di basso livello.
  </div>
</div> 

Nel nostro esempio: 
- alto livello: ServizioEsami; 
- basso livello: ArchivioEsamiSuFile. 

Violazione: l’alto livello dipende direttamente dal basso livello concreto. Il servizio crea direttamente la dipendenza: new ArchivioEsamiSuFile(...) dentro il costruttore. Effetto: ogni variazione infrastrutturale entra nel codice business. Criterio diagnostico: se per cambiare tecnologia devo modificare un servizio di dominio/applicazione, DIP non è rispettato.  

<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Pasted image 20260401171025.png" style="width: 100%; border-radius: 8px;">
    <img src="Pasted image 20260401171107.png" style="width: 100%; border-radius: 8px;">
  </div>
  <div style="flex: 1.5;">
Spostiamo l’interfaccia ArchivioEsami nel dominio (modulo più stabile). ServizioEsami dipende solo da ArchivioEsami, non da classi concrete. Le implementazioni concrete (SuFile, SuDb, InMemoria, . . . ) dipendono dall’astrazione, non il contrario. Abbiamo invertito la dipendenza: 
<li>il modulo di alto livello (service) non conosce più i dettagli di persistenza;</li>
<li>i dettagli dipendono dall’astrazione comune.</li>
  </div>
</div> 

La soluzione resta DIP se: 
- il servizio dipende da ArchivioEsami (astrazione), non da ArchivioEsamiSuFile; 
- la scelta dell’implementazione avviene all’esterno (composizione/iniezione). 

Impatto diretto sui test: 
- si può usare un’implementazione finta: 
	- ArchivioEsamiInMemoria, 
	- oppure un mock creato con un framework di test. 
- i test non richiedono file reali, database o configurazioni esterne. 

Checklist rapida DIP: 
- il servizio conosce classi concrete infrastrutturali? 
- la dipendenza concreta viene istanziata dentro il servizio? 
- è possibile testare il servizio senza risorse esterne?  

**I principi SOLID: come lavorano insieme:** 
- **SRP:** mantiene coese le classi (una responsabilità principale per classe). 
- **OCP:** consente nuove varianti senza modificare continuamente il nucleo stabile. 
- **LSP:** rende affidabile la sostituzione tra superclasse e sottoclasse. 
- **ISP:** evita interfacce troppo ampie e dipendenze inutili lato client. 
- **DIP:** separa business e infrastruttura tramite astrazioni. 

Insieme, questi principi supportano: 
- manutenibilità, 
- estendibilità, 
- testabilità.  

Applicare SOLID non significa: 
- creare mille micro-classi inutili, 
- astrarre tutto “per principio”. 

Significa: 
- riconoscere code smell (classi sovraccariche, dipendenze rigide), 
- avere strumenti concettuali per migliorare il design, 
- trovare un equilibrio tra semplicità e flessibilità.  

**SOLID e refactoring:** Spesso non si progetta “perfettamente” al primo colpo: SOLID aiuta a migliorare in modo incrementale. Workflow tipico di refactoring guidato da SOLID: 
- individuare code smell e punti di accoppiamento eccessivo; 
- chiarire responsabilità (SRP) e interfacce (ISP); 
- isolare variabilità (OCP) e dipendenze infrastrutturali (DIP); 
- verificare sostituibilità corretta tra tipi correlati (LSP). 

Obiettivo: migliorare il design senza riscrivere tutto da zero.  

**Mini-checklist per UniManager:** 
1. SRP: 
	- ci sono classi che fanno sia business che I/O? 
	- posso separare persistenza, log, presentazione? 
2. OCP: 
	- ci sono grandi if/switch su tipo o modalità? 
	- posso spostare logiche variabili in classi dedicate, evitando di toccare sempre lo stesso metodo? 
3. LSP: 
	- una sottoclasse può sostituire la superclasse senza rompere il client? 
	- compaiono controlli sul tipo concreto o eccezioni inattese? 
4.  ISP: 
	- i client dipendono da metodi che non usano? 
	- ci sono metodi vuoti o UnsupportedOperationException? 
5. DIP: 
	- i servizi dipendono da classi concrete di repository? 
	- posso introdurre interfacce nel dominio e implementazioni nei layer esterni?

###### Design patterns essenziali: strategy e factory method
Nella lezione precedente abbiamo visto tutti e cinque i principi SOLID. Questi principi ci dicono come dovrebbe essere un buon design: 
- classi con ruoli chiari; 
- codice estendibile senza modificare quello già scritto; 
- dipendenze dirette da astrazioni e non da dettagli concreti. 

Domande naturali: 
- quali soluzioni concrete possiamo usare per mettere in pratica questi principi? 
- come riconoscere schemi ricorrenti di buona progettazione? 

In questa lezione vedremo due pattern che realizzano nella pratica questi principi: 
- **Strategy** soprattutto a supporto di OCP, 
- **Factory Method** molto legato a **DIP** (e a OCP sulla creazione).  

**Design pattern:** I design pattern nascono da problemi che si ripetono: 
- varianti di comportamento gestite con `if/else` crescenti; 
- creazione di oggetti sparsa in molti punti del codice. 

Un pattern fornisce: 
- una struttura riusabile; 
- un linguaggio condiviso nel team; 
- una soluzione coerente con principi come OCP e DIP. 

Criterio pratico: 
- usare un pattern quando semplifica evoluzione e manutenzione; 
- evitarlo quando aggiunge astrazione senza variabilità reale.  

Oggi ci concentriamo su due pattern fondamentali e trasversali: 
- *Strategy:* rendere intercambiabili i comportamenti (es. diversi criteri di ordinamento degli appelli in UniManager), in linea con OCP (nuovi comportamenti senza modificare il contesto). 
- *Factory Method:* centralizzare e strutturare la creazione degli oggetti (es. creazione di diverse tipologie di ExamSession in UniManager), aiutando a rispettare DIP (dipendere da ExamSession, non dalle sue sottoclassi). 

L’idea non è memorizzare definizioni, ma: 
- riconoscere quando un problema è adatto a questi pattern; 
- saperli usare per mantenere il codice pulito, estensibile e coerente con i principi SOLID visti nella lezione precedente.  

**Strategy:** contesto tipico:
- una classe deve eseguire un certo comportamento; 
- esistono varie alternative per implementarlo; 
- vorremmo poter cambiare questo comportamento senza modificare la classe. 

Strategy propone di: 
- estrarre il comportamento in un’interfaccia separata; 
- avere più strategie concrete intercambiabili; 
- far sì che il contesto usi sempre l’interfaccia, non le implementazioni concrete. 

Effetto: 
- Strategy permette di iniettare comportamenti diversi senza toccare il codice del contesto. 
- È uno strumento concreto per applicare OCP e ridurre dipendenze rigide verso implementazioni specifiche.  

> [!example] grandi `if/else` per i comportamenti 
> Nell'esempio di UniManager vogliamo ordinare gli appelli: 
> - per data; 
> - per corso; 
> - per numero di iscritti. 
> 
> ![[Pasted image 20260413182522.png|200]]
> 
>Problemi: ogni volta che vogliamo aggiungere un nuovo criterio dobbiamo modificare il metodo: il codice non è pensato per essere esteso facilmente. Il metodo, inoltre, cresce sempre di più ed è difficile da mantenere e testare.  
>
>Idea: 
> - definire un’interfaccia SortStrategy; 
> - avere implementazioni diverse per ogni criterio; 
> - lasciare che un ExamSessionService usi sempre una SortStrategy.
> 
> ![[Pasted image 20260413182647.png|200]]
> 
> ![[Pasted image 20260413182709.png|200]]
> 
> ![[Pasted image 20260413182751.png|200]]
> 
> Il contesto non conosce le implementazioni concrete, usa sempre e solo l’interfaccia SortStrategy.
> 
> Strategie intercambiabili a runtime:
> 
> ![[Pasted image 20260413182906.png|200]]
> 
> Il codice che usa ExamSessionService può: 
> - cambiare la strategia in base a una scelta dell’utente; 
> - selezionare la strategia da configurazione; 
> - sostituire le strategie senza modificare il servizio.
> 
> Secondo esempio Strategy: politiche di ammissione:
> 
> ![[Pasted image 20260413183043.png|200]]
> 
> - Il contesto cambia politica senza modificare il servizio. 
> - Nuova regola: nuova classe, non nuovo `if/else`.
> 
> Terzo esempio Strategy: costo spedizione (e-commerce):
> 
> ![[Pasted image 20260413183137.png|200]]
> 
> - Cambia la politica di spedizione, non il CheckoutService. 
> - Stessa struttura del pattern, dominio completamente diverso.

**Strategy: vantaggi principali (collegamento a SOLID):** 
- *OCP (Open/Closed):* 
	- per aggiungere un nuovo criterio basta creare una nuova SortStrategy; 
	- il contesto non si modifica: si estende senza toccare codice esistente. 
- *DIP (Dependency Inversion):* 
	- il contesto dipende dall’interfaccia SortStrategy, non dalle classi concrete; 
	- riduce dipendenze rigide e rende la sostituzione delle strategie più semplice. 
- *SRP (Single Responsibility):* 
	- il contesto non deve “conoscere” tutti i criteri di ordinamento; 
	- ogni criterio è confinato in una classe dedicata con un compito chiaro. 
- *Testabilità:* 
	- strategie testabili in isolamento; nel contesto si possono usare strategie fittizie.  
	
**Strategy: errori tipici (rischi su SOLID):** 
- *Introdurre Strategy senza varianti reali:* 
	- aggiunge classi senza beneficio: aumenta complessità senza migliorare OCP/DIP. 
- *Strategie troppo grandi o “onnivore”:* 
	- una Strategy dovrebbe incapsulare un comportamento specifico; 
	- se dentro finisce “mezza applicazione”, si perde SRP. 
- *Contesto che continua a fare if/else sulle strategie:* 
	- es. controlli sul tipo concreto (instanceof) o switch sul criterio; 
	- si indebolisce OCP e spesso anche DIP. 
- *Dipendenze concrete nelle strategie:* 
	- se ogni Strategy istanzia direttamente servizi o repository concreti, si rischia una violazione di DIP (dipendenze dai dettagli).  
	
Con Strategy abbiamo visto come rendere flessibili i comportamenti di un oggetto: 
- criteri di ordinamento intercambiabili; 
- possibilità di cambiare logica senza toccare il codice del servizio. 

Ma in un sistema reale esiste un’altra fonte comune di rigidità: la creazione degli oggetti. Spesso troviamo: 
- molti new sparsi per il codice; 
- istanze create in modi diversi e difficili da uniformare; 
- dipendenze concrete “nascoste” dentro i costruttori. 

Domanda successiva naturale:  possiamo rendere flessibile anche il modo in cui gli oggetti vengono creati?  

**Factory Method:** Definiamo: 
- un prodotto astratto (es. ExamSession); 
- una o più sottoclassi concrete (es. WrittenSession, OralSession); 
- una classe Creator che espone un metodo di fabbrica (il Factory Method) per creare il prodotto. 

Il codice che usa il Creator: 
- lavora solo con il tipo astratto; 
- non conosce i dettagli delle classi concrete. 

Effetto: 
- la decisione su quale sottoclasse istanziare è concentrata nel Factory Method; 
- il resto del codice può dipendere solo dal tipo astratto: questo va nella direzione di DIP e OCP.  

> [!example] UniManager: creazione di ExamSession 
> Immaginiamo di avere varie tipologie di appello: 
> - WrittenSession (scritto); 
> - OralSession (orale); 
> - MixedSession (scritto + orale, non mostrata nei diagrammi per brevità). 
> - 
> Senza Factory Method: 
> - il codice farebbe new WrittenSession(...) o new OralSession(...) in molti punti diversi; 
> - cambiare il modo in cui creiamo le sessioni richiederebbe modifiche in tante classi. 
> 
> Con Factory Method: 
> - concentriamo la creazione in un creator (o pochi creator); 
> - il resto del sistema dipende solo da ExamSession.  
> 
> ![[Pasted image 20260413183852.png|200]]
> 
> Factory Method in Java: prodotto astratto:
> 
> ![[Pasted image 20260413183915.png|200]]
> 
> - ExamSession è il tipo astratto comune usato dal client. 
> - WrittenSession/OralSession cambiano il dettaglio concreto ma rispettano lo stesso contratto. 
> - Il client lavora sul tipo astratto, non sulla sottoclasse concreta.
> 
> Factory Method in Java: Creator e ConcreteCreator
> 
> ![[Pasted image 20260413184114.png|200]]
> 
> - ExamSessionCreator definisce il flusso comune e delega il punto variabile a createSession(). 
> - I concrete creator (WrittenSessionCreator, OralSessionCreator) decidono quale prodotto concreto istanziare. 
> - La creazione resta estendibile: nuovo tipo di appello = nuovo creator, senza modificare il client.
> 
> Uso del Factory method
> 
> ![[Pasted image 20260413184201.png|200]]
> 
> - Per cambiare il tipo di sessione basta fornire un ExamSessionCreator diverso. 
> - Nessuna modifica al codice di ExamManagementService.  
> 
> Secondo esempio Factory Method: nuova tipologia di appello
>  
> ![[Pasted image 20260413184445.png|200]]  
> 
> Terzo esempio Factory Method: parser documenti 
> 
> ![[Pasted image 20260413184506.png|200]]
> 
> - Prodotto: Documento; concreti: DocumentoCsv/DocumentoJson. 
> - Flusso comune in parse(), varia solo la creazione del prodotto concreto. 
> - Nuovo formato (XmlParser) senza toccare il client.  

**Factory Method: vantaggi principali (collegamento a SOLID):** 
- *DIP (Dependency Inversion):* 
	- il client lavora con ExamSession e ExamSessionCreator, non con WrittenSession/OralSession; 
	- la scelta della classe concreta è confinata nel Factory Method. 
- *OCP (Open/Closed):* 
	- per aggiungere un nuovo tipo di sessione si estende con nuove sottoclassi/creator; 
	- il codice che usa il prodotto non va modificato. 
- *SRP (Single Responsibility):* 
	- la creazione non è sparsa nel codice: è centralizzata nel Creator; 
	- il servizio si concentra sulla logica applicativa, non sui dettagli costruttivi. 
- *Coerenza del processo di creazione:* 
	- eventuali passi comuni (validazioni, log, notifiche) stanno in un punto solo.  
	
**Factory Method: errori tipici (rischi su SOLID):** 
- *Creator che diventa un “God object”:* 
	- se nel Creator finiscono troppe responsabilità oltre alla creazione, si rischia di violare SRP. 
- *Il client continua a conoscere le classi concrete:* 
	- se compaiono comunque new WrittenSession(...) in giro, si perde il beneficio su DIP. 
- *Esplosione di sottoclassi senza bisogno reale:* 
	- troppi Creator “micro” possono aumentare complessità senza vero guadagno; 
	- meglio introdurlo quando la creazione varia davvero.  
	
**Strategy vs Factory Method:** 
- *Strategy:* 
	- incapsula un comportamento variabile; 
	- il contesto delega l’esecuzione di un’operazione a una Strategy; 
	- esempio UniManager: criterio di ordinamento degli appelli. 
- *Factory Method:*
	- incapsula la creazione di oggetti (prodotti); 
	- il client delega la scelta della classe concreta al Factory Method; 
	- esempio UniManager: creazione di diverse tipologie di appello (ExamSession). 
- *In comune:* 
	- favoriscono OCP: estendere (nuove strategie / nuovi prodotti) senza modificare i client; 
	- favoriscono DIP: dipendenza da astrazioni (SortStrategy, ExamSession, Creator); 
	- aiutano SRP: separano responsabilità (comportamenti / creazione) dal resto della logica.

###### Design Patterns essenziali: Adapter e Observer
Abbiamo visto: 
- Strategy: rendere intercambiabili i comportamenti; 
- Factory Method: centralizzare e strutturare la creazione degli oggetti. 

Questi pattern agivano soprattutto: 
- “dentro” al sistema; 
- su logiche di dominio e creazione. 

Ora affrontiamo due problemi diversi ma altrettanto comuni: 
- come parlare con il mondo esterno senza contaminare il dominio interno (Adapter); 
- come reagire a eventi interni senza creare dipendenze rigide tra classi (Observer). 

Obiettivo: riconoscere quando questi pattern emergono “da soli” mentre il sistema cresce e si collega ad altri componenti.  

**Adapter:** 
Il nostro sistema è costruito attorno a interfacce interne pulite. Le sorgenti esterne (CSV, JSON, API) espongono interfacce diverse. Effetti tipici: 
- formati incompatibili; 
- strutture dati grezze; 
- gestione I/O o eccezioni non allineata al nostro dominio. 

Obiettivo: integrare queste sorgenti senza “sporcare” il codice del dominio.  

> [!example] Esempi dal caso UniManager: 
> Libreria LegacyCsvLibrary: 
> - restituisce solo List; 
> - formati di data non uniformi; 
> - API con metodi statici e dettagli I/O esposti. 
> 
> Libreria o servizio JSON: 
> - restituisce strutture grezze (Map); 
> - formati temporali diversi (ISO, timestamp, timezone); 
> - gestione di errori di rete e codici HTTP. 
> 
Queste differenze non devono “trapelare” dentro il dominio.
  
**Adapter:** Definiamo una interfaccia Target interna (es. ExamSessionImporter). Per ogni sorgente (CSV, JSON, API,…) realizziamo un Adapter dedicato. Ogni Adapter: 
- implementa l’interfaccia Target; 
- incapsula la libreria esterna (‘Adaptee’); 
- converte i dati grezzi in oggetti di dominio. 

Il dominio usa sempre il Target e resta indipendente dalla tecnologia esterna.  

**UniManager: CSV e JSON come sorgenti esterne:** UniManager vuole ottenere sempre una List. Le fonti possono essere diverse: 
- file CSV letti da LegacyCsvLibrary; 
- file JSON o API con LegacyJsonLibrary. 

Problema: senza Adapter, il servizio dovrebbe conoscere formati, parsing, eccezioni, convenzioni di entrambe. Soluzione: due Adapter: 
- CsvImporterAdapter; 
- JsonImporterAdapter; 

che implementano la stessa interfaccia ExamSessionImporter.  

In UniManager il dominio vede sempre e solo: 
- un’interfaccia (ExamSessionImporter), 
- un tipo di ritorno (List). 

Integrazioni diverse (CSV, JSON, API) sono intercambiabili. Possiamo aggiungere nuovi formati senza toccare il dominio. Tutta la logica di mapping e parsing resta confinata negli Adapter.  

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260414094231.png" width="200">
</div>

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260414094437.png" width="200">
</div>

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260414094507.png" width="200">
</div>

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260414094537.png" width="200">
</div>

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260414094622.png" width="200">
</div>

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260414094654.png" width="200">
</div>

*Vantaggi:*
- *DIP (Dependency Inversion):* 
	- ExamSessionService dipende dall’interfaccia ExamSessionImporter; 
	- le librerie esterne restano “dietro” gli Adapter. 
- *OCP (Open/Closed):* 
	- per supportare un nuovo formato (XML, API diversa,…) basta aggiungere un nuovo Adapter; 
	- il servizio che importa non cambia. 
- *SRP (Single Responsibility):* 
	- parsing, conversione e mapping sono confinati nell’Adapter; 
	- il dominio e i servizi non si riempiono di dettagli I/O e formati. 
- *Testabilità e manutenzione:* 
	- Adapter testabile in isolamento; 
	- il servizio può usare importer fittizi nei test.  
	
*Errori tipici:* 
- Adapter che fa troppo: se oltre a convertire dati inizia a contenere logica di dominio, si viola SRP. 
- Service che dipende comunque dalla libreria esterna: se il dominio usa tipi della libreria (es. String[] o Map ovunque), si indebolisce DIP. 
- If/else sui formati dentro il servizio: scegliere CSV/JSON con switch nel servizio porta a modifiche frequenti: peggiora OCP.  

**Observer:** Con Adapter abbiamo gestito l’integrazione con l’esterno. Ora il problema è interno: come reagire a un evento? Il contesto tipico è: 
- una classe genera un evento (es. nuovo appello); 
- più componenti vogliono essere informati. 

Soluzione ingenua: la classe chiama direttamente tutti gli altri (sendEmail(), updateUi(), writeLog(),…). Conseguenze: 
- la classe finisce per conoscere troppi dettagli su altri componenti; 
- il codice diventa più difficile da mantenere e capire; 
- aggiungere un nuovo comportamento richiede modifiche al codice esistente.  

Il pattern Observer separa: 
- il Subject, che produce gli eventi; 
- gli Observer, che vogliono essere notificati. 

Meccanismo: 
- gli Observer si registrano presso il Subject; 
- quando qualcosa cambia, il Subject chiama un unico metodo (notifyObservers()) che notifica tutti gli iscritti; 
- il Subject non conosce i dettagli di cosa faranno gli Observer. 

Effetto: 
- possiamo aggiungere/rimuovere osservatori senza modificare il Subject; 
- il sistema diventa più estendibile e testabile.  

> [!example] UniManager: notificare la creazione di un appello 
> In UniManager, la creazione di un nuovo ExamSession deve attivare più reazioni: 
> - aggiornare la vista, 
> - inviare notifiche agli studenti, 
> - scrivere un log. 
> 
> Il servizio che crea l’appello non deve conoscere questi dettagli. Con Observer: 
> - ExamEventPublisher funge da Subject; 
> - EmailNotifier, UiRefresher, LoggingObserver sono Observer che reagiscono indipendentemente.  

<div style="display: flex; justify-content: center;">
  <img src="Screenshot From 2026-04-14 09-53-59.png" width="200">
</div>

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260414095536.png" width="200">
</div>

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260414095611.png" width="200">
</div>

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260414095640.png" width="200">
</div>

*Vantaggi:* 
- *DIP (Dependency Inversion):* 
	- il Subject dipende dall’interfaccia ExamEventObserver, non dagli observer concreti; 
	- le reazioni (email, log, UI, . . . ) restano dettagli sostituibili. 
- *OCP (Open/Closed):* 
	- nuove reazioni all’evento si aggiungono registrando un nuovo Observer; 
	- il codice che genera l’evento non va modificato. 
- *SRP (Single Responsibility):* 
	- il servizio si concentra sulla logica principale (creare l’appello); 
	- le reazioni all’evento stanno in classi separate, ognuna con un compito specifico. 
- *Testabilità:* 
	- nei test si possono usare observer fake e verificare che la notifica avvenga correttamente.  
	
*Errori tipici:* 
- Subject che conosce gli observer concreti: se istanzia direttamente new EmailNotifier(), perde DIP. 
- Observer con responsabilità miste: se un singolo observer gestisce email + log + UI, si perde SRP. 
- Notifiche “a cascata” difficili da controllare: 
	- se gli observer generano altri eventi senza regole chiare, il sistema diventa fragile; 
	- spesso si finisce a “patchare” codice esistente, andando contro OCP.  
	
**Adapter vs Observer:** 
- *Adapter:*
	- risolve problemi di incompatibilità di interfacce; 
	- tipico per integrare librerie, servizi, formati dati esterni; 
	- es. UniManager: import di appelli da CSV/JSON con un’interfaccia interna pulita. 
- *Observer:*
	- risolve problemi di notifica di eventi; 
	- tipico per disaccoppiare chi produce un evento da chi reagisce; 
	- es. UniManager: notifiche agli studenti e aggiornamento viste quando viene creato un nuovo appello. 
- *In comune:* 
	- proteggono il codice da dipendenze rigide; 
	- spostano la complessità nei confini del sistema (integrazione, eventi); 
	- rendono il design più modulare e evolvibile.  
	
**Design Pattern essenziali:** 

| Pattern        | Problema che risolve                                                                 | Esempi in UniManager                                                       |
| -------------- | ------------------------------------------------------------------------------------ | -------------------------------------------------------------------------- |
| Strategy       | Comportamenti che possono variare nel tempo; evita lunghi blocchi `if/else.`         | Ordinamento degli appelli secondo criteri diversi (data, corso, iscritti). |
| Factory Method | Creazione di oggetti che può richiedere logiche diverse; centralizza la costruzione. | Creazione di diverse tipologie di ExamSession (Written, Oral).             |
| Adapter        | Integrazione di librerie/formati esterni incompatibili con il dominio interno.       | Import di appelli da CSV/JSON tramite ExamSessionImporter.                 |
| Observer       | Notificare cambiamenti senza accoppiamento rigido tra componenti.                    | Avvisare studenti, viste e log quando viene creato un nuovo appello.       |
I pattern affrontano quattro problemi distinti: comportamenti, creazione, integrazione ed eventi. Insieme costituiscono la base per un design evolvibile.

> [!info] Sintesi:
> - Adapter e Observer affrontano due problemi ricorrenti diversi: integrazione con esterni ed eventi interni tra componenti. 
> - Adapter protegge il dominio da formati e librerie incompatibili, spostando la complessità ai confini del sistema. 
> - Observer consente di reagire a cambiamenti senza legare direttamente chi produce l’evento a tutti i componenti interessati. 
> - Entrambi i pattern favoriscono estendibilità e riduzione dell’accoppiamento, in linea con OCP e DIP. 
> - Il beneficio dei pattern dipende da come vengono applicati: migliorano il design se rispondono a una variabilità reale. 
> - Gli esempi su UniManager mostrano come questi pattern preparino il sistema a integrazioni e notifiche future.

###### Qualità del design: code smells
**Perché parlare di qualità interna?** Un software può funzionare ed essere comunque difficile da mantenere e da far crescere nel tempo. La qualità interna determina in modo diretto: 
- quanto è semplice correggere bug senza introdurre effetti collaterali; 
- quanto è facile aggiungere nuove funzionalità senza “toccare tutto”; 
- quanto si capisce il codice dopo settimane/mesi (anche da parte di chi non lo ha scritto). 

Nella pratica, la qualità interna è ciò che rende un progetto: 
- evolvibile invece che fragile, 
- leggibile invece che opaco, 
- modulare invece che intrecciato, 
- più testabile e quindi più sicuro da modificare.  

**Qualità esterna vs qualità interna:** 
- *Qualità esterna:* ciò che percepisce l’utente finale (funzionalità offerte, prestazioni, affidabilità, esperienza d’uso). 
- *Qualità interna:* proprietà strutturali del codice (chiarezza, organizzazione, modularità, dipendenze sane). 

Effetto tipico nei progetti reali: oggi funziona → domani costa troppo cambiarlo o estenderlo. In questa lezione ci concentriamo in particolare su: 
- segnali che indicano un design che sta lentamente peggiorando; 
- strumenti concettuali per riconoscerli e intervenire in anticipo.  

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260415160609.png" width="200">
</div>

**Come leggere il design:** 
- *Coesione:* quanto gli elementi di uno stesso modulo/classe stanno bene insieme e lavorano per un unico scopo. Se c'è alta coesione allora avrò una classe focalizzata, più facile da capire e testare. 
- *Accoppiamento:* quanto un modulo/classe dipende dagli altri. Se ho basso accoppiamento allora ci saranno meno dipendenze rigide, cambiamenti più locali. 

Regola pratica: alta coesione + basso accoppiamento.  

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260415160835.png" width="200">
</div>

**Coesione:** misura la “compattezza” delle responsabilità assegnate a un modulo o a una classe. Un modulo realmente coeso: 
- ha un obiettivo chiaro e facilmente esplicitabile; 
- contiene responsabilità dello stesso tipo e tra loro correlate; 
- si può descrivere in modo naturale con una singola frase. 

Segnali tipici di bassa coesione: 
- una classe che offre molte funzionalità non correlate tra loro; 
- metodi che lavorano su concetti diversi e che non condividono un filo logico comune; 
- nomi generici (Manager, Utils, Helper) che cercano di “coprire tutto” senza un ruolo preciso.  

**Accoppiamento:** misura quanto un modulo dipende da altri moduli e dalle loro decisioni interne. Un accoppiamento elevato significa: 
- conoscere e sfruttare dettagli interni di altri componenti; 
- vincolarsi a implementazioni concrete invece che a contratti astratti; 
- amplificare l’effetto di una modifica anche lontano dal punto di intervento. 

Segnali tipici di alto accoppiamento: 
- accesso diretto a campi interni e strutture dati di altre classi; 
- catene di chiamate lunghe (train wreck); 
- dipendenze diffuse “ovunque” verso componenti instabili o tecnologici.  

Nei progetti reali emergono spesso due problemi ricorrenti, che tendono a rafforzarsi a vicenda se non vengono controllati: 
1. classi che crescono troppo nel tempo, accumulando responsabilità diverse → bassa coesione; 
2. componenti sempre più intrecciati tra loro, con dipendenze diffuse → alto accoppiamento. 

Una regola pratica utile per diagnosticare la situazione è la seguente: 
- se per cambiare una parte A è necessario capire anche B, C, D, allora l’accoppiamento è probabilmente eccessivo; 
- se una classe ha “troppe ragioni per cambiare” in contesti diversi, allora la sua coesione è probabilmente insufficiente. 

L’obiettivo del buon design è mantenere questi due aspetti in equilibrio, così che ogni modulo possa evolvere in modo il più possibile indipendente.

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260415161143.png" width="200">
</div>

Spesso il problema non è che il codice sia “sbagliato” o non funzionante, ma che risulti scomodo da leggere, capire e modificare. Situazioni tipiche in cui emerge questa sensazione: 
- metodi troppo lunghi, con molti passaggi e livelli di dettaglio mescolati; 
- classi che crescono a ogni nuova feature senza una struttura chiara; 
- logiche di dominio mescolate con dettagli tecnici o di presentazione. 

Questa sensazione di disagio è in realtà preziosa, non è una prova matematica di errore, ma è un segnale da prendere sul serio e verificare usando criteri progettuali e strumenti concettuali.  

**Code smells:** sono indicatori di possibili problemi nella struttura e nell’organizzazione del design. Non sono bug né errori di funzionamento: 
- il programma continua a funzionare correttamente; 
- ma il design mostra segnali di tensione e perdita di chiarezza. 

Sono particolarmente utili perché: 
- aiutano a individuare il degrado del design prima che diventi evidente o ingestibile; 
- suggeriscono punti del codice in cui il refactoring può portare il maggior beneficio.  

**Smells e refactoring:** 
- Smell → ipotesi di problema nel design: qualcosa suggerisce che la struttura non sia ottimale. 
- Refactoring → intervento disciplinato per migliorare la struttura senza cambiare il comportamento osservabile. 

L'idea chiave è che gli smells non dicono “cambia qui” in modo automatico, ma dicono “qui vale la pena guardare meglio e fare una verifica”. In pratica, il flusso di lavoro tipico è: 
- si osserva il codice e si riconoscono possibili segnali; 
- si formula una diagnosi (quale problema progettuale c’è sotto?); 
- si applicano refactoring mirati, piccoli e controllati.  

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260415172553.png" width="200">
</div>

Gli smell che vedremo sono:
- *Large Class:* troppe responsabilità nella stessa classe. 
- *Long Method:* metodo troppo lungo e difficile da manutenere. 
- *Feature Envy:* metodo che usa più dati di un’altra classe che della propria. 
- *Data Clumps:* gruppi di parametri che ricorrono sempre insieme. 
- *God Object:* classe centrale che coordina, decide e conosce troppo. 
- *Duplicated Code:* stessa logica copiata in più punti. 
- *Primitive Obsession:* uso eccessivo di tipi primitivi al posto di concetti di dominio. 
- *Divergent Change:* una classe cambia per motivi eterogenei. 
- *Shotgun Surgery:* una modifica richiede piccoli cambi sparsi in molti file.  

**Large Class:** Una Large Class è una classe che accumula troppe responsabilità e, di conseguenza, troppe ragioni diverse per cambiare nel tempo. Effetti tipici che emergono in questi casi: 
- bassa coesione, perché le responsabilità non sono omogenee; 
- maggiore difficoltà di test, poiché è complesso isolare i comportamenti; 
- modifiche rischiose, perché molte funzionalità “passano tutte di lì”. 

Questo smell si accompagna spesso a: 
- nomi generici (`*Manager`, `*Service`) che non chiariscono il ruolo reale; 
- una crescita “per accumulo” di metodi, senza una struttura progettuale chiara. 

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260415172849.png" width="200">
</div>

> [!example]
> ![[Pasted image 20260415173015.png|200]] 
> 
> - Una sola classe copre dominio, persistenza e presentazione.
> - Ogni nuova richiesta tende a finire nello stesso punto.

<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Pasted image 20260415173252.png" style="width: 100%; border-radius: 8px;">
  </div>
  <div style="flex: 1.5;">
  <b>Large class (refactoring)</b>
  <li>Obiettivo: separare per responsabilità e per motivo di cambiamento.</li>
  <li>Decomposizione esplicita: dominio, persistenza e presentazione.</li> 
  <li>Effetto: classi più coese e modifiche più locali.</li>
  </div>
</div>

**Long Method:** Un Long Method è un metodo che cresce oltre misura e tende a concentrare troppa logica in un unico punto del codice. Tipicamente mescola tra loro: 
- livelli di astrazione diversi (alto livello e dettagli operativi); 
- logica di business e dettagli tecnici o infrastrutturali; 
- casi normali ed eccezioni gestiti nello stesso flusso lineare. 

Effetti tipici di questo smell: 
- è difficile capire “che cosa fa” davvero il metodo nel suo insieme; 
- è difficile riusare o testare singole parti della logica; 
- aumenta la probabilità di introdurre bug quando si modifica anche solo un passaggio. 

Refactoring tipico e consigliato: **Extract Method**.

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260415173426.png" width="200">
</div>

> [!example]
> ![[Pasted image 20260415173549.png|200]]
>- Validazione, regole, persistenza, audit e output nello stesso metodo.

<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Pasted image 20260415173748.png" style="width: 100%; border-radius: 8px;">
  </div> 
  <div style="flex: 1.5;">
  <b>Long method (refactoring)</b>
  <li>Refactoring applicato: <b>Extract method</b> </li>
  <li>Ogni step ha nome esplicito e responsabilità unica</li> 
  <li>Il metodo principale torna leggibile a livello di processo</li>
  </div>
</div>

**Feature Envy:** Un Feature Envy si verifica quando un metodo utilizza più dati e comportamenti di un’altra classe che non quelli della classe in cui è definito. Indicazione progettuale tipica: 
- la responsabilità è probabilmente posizionata nel posto sbagliato; 
- il modello concettuale risulta sbilanciato o mal distribuito. 

Sintomi comuni di questo smell: 
- catene di accessi del tipo `a.getB().getC().getX()`, che rivelano una conoscenza eccessiva della struttura interna altrui; 
- un metodo che “conosce troppo” lo stato e i dettagli di un altro oggetto. 

Refactoring tipici: **Move Method** oppure **Extract Class**. 

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260415174018.png" width="200">
</div>

> [!example]
> ![[Pasted image 20260415174120.png|200]]
> - Il metodo usa quasi solo dati interni di ExamSession, Course, Room, Teacher. 
> - La responsabilità di formattazione è nel posto sbagliato.  

<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Pasted image 20260415174223.png" style="width: 100%; border-radius: 8px;">
  </div>  
  <div style="flex: 1.5;">
  <b>Long method (refactoring)</b>
  <li>Refactoring applicato: <b>Move method</b> </li>
  <li>Il comportamento torna vicino ai dati che usa davvero</li> 
  </div>
</div>

**Data Clumps:** I Data Clumps sono gruppi di variabili che compaiono sempre insieme e vengono passati in giro come un unico “pacchetto”, anche se nel codice restano separati. Situazioni tipiche: 
- gli stessi parametri ripetuti in molti metodi e costruttori; 
- le stesse variabili che “viaggiano insieme” tra classi e servizi. 

Questo smell spesso indica che nel modello manca un concetto: un’entità o, più spesso, un Value Object che rappresenti quel gruppo in modo esplicito. Il refactoring tipico è quello di introdurre un **Value Object** (es. StudentId, ExamInfo, ExamDescriptor,…). 

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260415174407.png" width="200">
</div>

> [!example]
> ![[Pasted image 20260415174452.png|200]]
> - Lo stesso blocco di parametri si ripete in più metodi. 
> - Manca un concetto esplicito nel modello.  

<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Pasted image 20260415174703.png" style="width: 100%; border-radius: 8px;">
  </div>  
  <div style="flex: 1.5;">
  <b>Data clumps (refactoring)</b>
  <li>Refactoring applicato: <b>Introduce value object</b> </li>
  <li>Firma più corta, semantica più chiara, validazioni centralizzabili</li> 
  <li>In questo esempio introduciamo anche StudentDescriptor per accorpare i dati studente che viaggiano insieme</li>
  </div>
</div>

**God Object:** Un God Object è una classe “onnipotente” che finisce per diventare il centro operativo del sistema e concentra decisioni che dovrebbero essere distribuite. In pratica è una classe che: 
- coordina tutto (orchestrazione di molte parti), 
- decide tutto (regole e policy sparse), 
- sa tutto (conoscenza di dettagli interni di molti moduli). 

Tipicamente produce effetti molto negativi: 
- accoppia tra loro parti del sistema che dovrebbero restare separate; 
- diventa un collo di bottiglia cognitivo: per capire il sistema bisogna passare da lì; 
- impedisce modularità e testabilità, perché è difficile isolarla e simulare dipendenze. 

È spesso la forma estrema di Large Class + responsabilità non distribuite (tutto viene “assorbito” in un punto). 

<div style="display: flex; justify-content: center;">
  <img src="Pasted image 20260415174855.png" width="200">
</div> 

> [!example]
> ![[Pasted image 20260415174953.png|200]]
> - Una sola classe concentra dominio, I/O, reportistica e presentazione. 
> - Ogni modulo finisce per dipendere da questo punto centrale.  

<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Pasted image 20260415175116.png" style="width: 100%; border-radius: 8px;">
  </div>   
  <div style="flex: 1.5;">
  <b>God object (refactoring)</b>
  <li>I ruoli sono separati in servizi distinti</li>
  <li>Il caso d'uso non passa più da una classe onnicomprensiva</li> 
  </div>
</div>

**Duplicated Code:** Il Duplicated Code compare quando la stessa logica è copiata in più punti. Effetti tipici: 
- ogni correzione va ripetuta in più file/metodi; 
- alto rischio di inconsistenza tra copie quasi uguali. 

Indizio operativo: “copia e incolla” usato per introdurre varianti minime. Refactoring tipici: **Extract Method, Extract Class, Template Method**.  

> [!example]
> ![[Pasted image 20260415175258.png|200]]
> - Le prime tre regole sono duplicate. 
> - La variante reale è solo l’ultima condizione.  

<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Pasted image 20260415175424.png" style="width: 100%; border-radius: 8px;">
  </div>   
  <div style="flex: 1.5;">
  <b>Duplicated code (refactoring)</b>
  <li>Regole comuni in un punto solo</li>
  <li>Le variabili restano piccole e controllabili</li> 
  </div>
</div>

**Primitive Obsession:** Riguarda l'uso eccessivo di tipi primitivi (String, int, double) al posto di concetti di dominio. Effetti tipici: 
- validazioni sparse e duplicate; 
- parametri poco autoesplicativi; 
- errori di scambio tra campi semanticamente diversi. 

Refactoring tipici: **Replace Data Value with Object, Introduce Value Object**.  

> [!example]
> ![[Pasted image 20260415175614.png|200]]
> - I parametri sono primitivi, ma rappresentano concetti di dominio. 
> - Le regole di validazione non sono incapsulate.  

<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Pasted image 20260415175734.png" style="width: 100%; border-radius: 8px;">
  </div>   
  <div style="flex: 1.5;"> 
  <b>Primitive obsession (refactoring)</b>
  <li>I vincoli diventano parte del modello</li>
  <li>La firma del metodo è più espressiva e robusta</li> 
  </div>
</div>

**Divergent Change:** Divergent Change si ha quando una singola classe viene modificata per motivi diversi e non correlati. Il segnale tipico è che ogni nuova richiesta (report, validazione, persistenza, UI) tocca sempre la stessa classe. L'effetto è quello di forte instabilità locale e regressioni frequenti. Il refactoring tipico: separare responsabilità per asse di cambiamento.  

> [!example]
> ![[Pasted image 20260415175911.png|200]]
> - Quattro motivi di cambiamento nella stessa classe

<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Pasted image 20260415180102.png" style="width: 100%; border-radius: 8px;">
  </div>   
  <div style="flex: 1.5;"> 
  <b>Divergent change (refactoring)</b>
  <li>Ogni classe di cambiamento ha una classe dedicata</li>
  <li>Le modifiche eterogenee non ricadono più nello stesso file</li> 
  </div>
</div>

**Shotgun Surgery:** Shotgun Surgery è il caso opposto del Divergent Change: una modifica concettualmente unica richiede tanti cambi minori in classi diverse. Segnali tipici: 
- pull request con molti file toccati per una feature piccola; 
- rischio elevato di dimenticare un punto. 

Refactoring tipico: accorpare comportamento correlato e chiarire i confini.  

> [!example]
>  - Nuova regola: “uno studente non può iscriversi a due appelli nello stesso slot orario”. 
>  - Implementazione dispersa:
>  
> ![[Pasted image 20260415180214.png|200]]
> - Una sola regola di dominio impatta troppi punti diversi.

<div style="display: flex; align-items: flex-start; gap: 20px;">
  <div style="flex: 1;">
    <img src="Pasted image 20260415180453.png" style="width: 100%; border-radius: 8px;">
  </div>   
  <div style="flex: 1.5;"> 
  <b>Shotgun surgery (refactoring)</b>
  <li>La regola è centralizzata in una policy unica</li>
  <li>I moduli delegano: una modifica tocca pochi punti</li> 
  </div>
</div>

> [!info] Riepilogo:
> - La qualità interna del design incide direttamente sulla capacità di mantenere, testare ed evolvere il software nel tempo. 
> - Coesione e accoppiamento sono criteri pratici per leggere la “salute” del codice oltre il solo comportamento corretto. 
> - I code smells non sono prove definitive di errore, ma segnali utili per individuare zone a rischio del design. 
> - Riconoscere presto gli smells aiuta a intervenire con refactoring prima che il degrado diventi architetturale. 
> - I miglioramenti di design spesso non cambiano il comportamento esterno, ma riducono fragilità e costo delle modifiche future. 
> - Gli esempi su UniManager mostrano come collegare osservazione degli smells e refactoring concreto.

