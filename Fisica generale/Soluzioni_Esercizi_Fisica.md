# Soluzioni — Esercizi di Fisica Generale (Simulazioni)

*Confronta questi passaggi con i tuoi, non solo i risultati finali: se un numero non torna, controlla prima le formule usate e i segni.*

---

## SIMULAZIONE 1

### Problema 1 — Due dischi su asta rotante

**Dati:** m₁=0,80 kg, m₂=1,4 kg, L=0,45 m, ω=6,5 rad/s

**(a) Posizione del centro di massa**
Per l'equilibrio rispetto al CM: m₁r₁ = m₂r₂, con r₁+r₂=L

r₁ = m₂L/(m₁+m₂) = (1,4·0,45)/2,2 = **0,286 m**
r₂ = m₁L/(m₁+m₂) = (0,80·0,45)/2,2 = **0,164 m**

**(b) Tensione dell'asta**
Ogni disco richiede una forza centripeta fornita dall'asta:
T = m₁ω²r₁ = 0,80·(6,5)²·0,286 = **9,68 N**
Verifica con l'altro disco: T = m₂ω²r₂ = 1,4·42,25·0,164 = 9,68 N ✓ (l'asta è in trazione/tensione)

**(c) Rottura dell'asta**
Sparita la forza centripeta, ciascun disco prosegue di moto rettilineo uniforme lungo la tangente alla traiettoria nell'istante della rottura (1ª legge di Newton), con velocità:
v₁ = ωr₁ = 6,5·0,286 = **1,86 m/s**
v₂ = ωr₂ = 6,5·0,164 = **1,06 m/s**
(direzioni tangenti e opposte, dato che i due dischi sono su lati opposti del CM)

---

### Problema 2 — Lastra infinita carica

**Dati:** σ=3,2·10⁻⁶ C/m², d=0,020 m, m=5,0·10⁻⁷ kg, q=−4,0·10⁻⁹ C

**(a) Campo e forza**
E = σ/(2ε₀) = 3,2·10⁻⁶/(2·8,854·10⁻¹²) = **1,81·10⁵ N/C** (diretto lontano dalla lastra, essendo σ>0)
F = |q|E = 4,0·10⁻⁹·1,81·10⁵ = **7,23·10⁻⁴ N**, diretta **verso la lastra** (carica negativa attratta da lastra positiva)

**(b) Velocità all'impatto**
a = F/m = 7,23·10⁻⁴/5,0·10⁻⁷ = 1446 m/s² (costante, campo uniforme)
v = √(2ad) = √(2·1446·0,020) = √57,84 = **7,61 m/s**

**(c) Seconda lastra con carica opposta**
Con due piani infiniti a cariche opposte (+σ e −σ), il campo:
- **fuori** dalle due lastre (a sinistra della prima e a destra della seconda): si annulla, **E = 0**
- **tra** le due lastre: i contributi si sommano, **E = σ/ε₀ = 3,62·10⁵ N/C**, diretto dalla lastra positiva verso quella negativa

La particella (che si trova tra le due lastre) sente quindi un campo **doppio** rispetto al punto (a): la forza su di essa raddoppia, F' ≈ 1,45·10⁻³ N, diretta ora verso la prima lastra (equilibrio impossibile, il moto continua accelerato verso la lastra positiva).

---

### Problema 3 — Filo infinito ed elettrone

**Dati:** I=8,5 A, v₀=2,1·10⁶ m/s (parallela al filo), r=3,0 mm

**(a) Campo magnetico**
B = μ₀I/(2πr) = (4π·10⁻⁷·8,5)/(2π·0,003) = **5,67·10⁻⁴ T**, tangente (azimutale) attorno al filo

**(b) Forza sull'elettrone**
Con v parallela al filo e B tangenziale (perpendicolari tra loro), |F| = |q|v₀B:
F = 1,60·10⁻¹⁹·2,1·10⁶·5,67·10⁻⁴ = **1,90·10⁻¹⁶ N**
Direzione: radiale, verso l'esterno o l'interno rispetto al filo a seconda del verso relativo di v e I (usa la regola della mano destra/sinistra per elettroni: F=qv×B con q negativo).

**(c) Lavoro della forza magnetica**
**No, il lavoro è sempre nullo.** La forza magnetica F=qv×B è sempre perpendicolare alla velocità, quindi non compie mai lavoro: cambia solo la direzione del moto, non il modulo della velocità (l'energia cinetica dell'elettrone resta costante).

---

## SIMULAZIONE 2

### Problema 1 — Blocco su guida circolare

**Dati:** m=2,0 kg, R=0,60 m, parte da fermo ad altezza 2R

**(a) Velocità nel punto più basso**
Conservazione dell'energia: mg(2R) = ½mv² → v² = 4gR = 4·9,807·0,60 = 23,54
v = **4,85 m/s**

**(b) Reazione normale nel punto più basso**
N − mg = mv²/R → N = m(g + v²/R) = m(g+4g) = 5mg
N = 5·2,0·9,807 = **98,1 N**

**(c) Distacco dalla guida**
Attenzione: partendo esattamente da fermo all'altezza 2R (il punto più alto di un anello di raggio R), questo è **il caso limite classico del "giro della morte"**: con questa energia iniziale, N=0 esattamente nel punto più alto e N>0 in ogni altro punto del percorso. Per conservazione dell'energia, risalendo dall'altro lato il blocco riguadagna la stessa quota con la stessa velocità (in modulo) punto per punto — quindi **non si stacca mai dalla guida**: risale fino a 2R con v=0 dall'altra parte, in perfetta simmetria. Se invece l'altezza di partenza fosse stata *inferiore* a 2R, il blocco si sarebbe staccato prima di raggiungere la sommità (dove N si annullerebbe per primo).

---

### Problema 2 — Circuito a ponte con due pile

**Dati:** V=9,0 V (due pile), R=15,0 Ω (tre resistori identici); topologia: due rami (uno per pila) verso un nodo centrale, più un terzo resistore dal nodo centrale a massa (schema tipo Problema 2 dell'Esame del 4/6/2025)

Ponendo la massa (i due poli comuni delle pile) a 0 V, il nodo "in alto" della pila di sinistra è a +9 V, quello della pila di destra (polarità invertita rispetto alla prima) è a −9 V. Chiamando V_A il potenziale del nodo centrale, applico Kirchhoff (nodi):

(9−V_A)/R + (−9−V_A)/R + (0−V_A)/R = 0
9 − V_A − 9 − V_A − V_A = 0 → **V_A = 0 V**

**(a) Correnti**
- Ramo pila sinistra: I = (9−0)/15 = **0,60 A**, entrante nel nodo
- Ramo pila destra: I = (0−(−9))/15 = **0,60 A**, uscente dal nodo verso la pila destra
- Ramo centrale (verso massa): I = (0−0)/15 = **0 A**

(Risultato tipico di un ponte bilanciato: il ramo centrale non porta corrente perché il nodo si trova esattamente al potenziale medio delle due pile.)

**(b) Potenza totale**
P = V·I per ciascuna pila = 9,0·0,60 = 5,4 W; totale = **10,8 W**
Verifica coi resistori: P=I²R sui due rami attivi = 0,60²·15 = 5,4 W ciascuno, 0 W sul ramo centrale → totale 10,8 W ✓

**(c) Con una pila invertita**
Ora entrambi i "nodi pila" sono a +9 V:
(9−V_A)/R + (9−V_A)/R + (0−V_A)/R = 0 → 18 = 3V_A → **V_A = 6 V**
- Correnti nei due rami pila: (9−6)/15 = **0,20 A** ciascuno, entranti nel nodo
- Corrente nel ramo centrale: (6−0)/15 = **0,40 A**, uscente verso massa

Le due pile ora "collaborano" (invece di scambiarsi corrente a vicenda) e alimentano insieme il resistore centrale.

---

### Problema 3 — Spira e campo variabile

**Dati:** l=0,12 m (A=0,0144 m²), R=0,18 Ω, B(t)=0,40t² T

**(a) F.e.m. e corrente a t=2,0 s**
dB/dt = 0,80t → a t=2s: dB/dt = 1,60 T/s
|ε| = A·dB/dt = 0,0144·1,60 = **0,0230 V (23,0 mV)**
I = ε/R = 0,0230/0,18 = **0,128 A**

**(b) Verso (Legge di Lenz)**
Se B è crescente (ad es. uscente dal piano della spira), la corrente indotta si oppone all'aumento del flusso, quindi crea un campo indotto opposto (entrante): per la regola della mano destra, la corrente indotta scorre **in verso orario** (vista dal lato da cui B esce).

**(c) Potenza dissipata**
P = I²R = 0,128²·0,18 = **2,95·10⁻³ W (2,95 mW)**

---

## SIMULAZIONE 3

### Problema 1 — Urto elastico, masse uguali

**Dati:** m₁=m₂=0,17 kg, v₁=2,5 m/s, v₂=−1,2 m/s

**(a) Velocità finali**
Per un urto elastico tra masse **uguali**, le velocità si scambiano semplicemente:
v₁' = **−1,2 m/s**, v₂' = **2,5 m/s**

**(b) Verifica energia**
Eₖ,i = ½·0,17·(2,5²+1,2²) = 0,085·7,69 = 0,654 J
Eₖ,f = ½·0,17·((−1,2)²+2,5²) = 0,654 J ✓ conservata

**(c) Urto completamente anelastico**
Conservazione quantità di moto: p = 0,17·2,5 + 0,17·(−1,2) = 0,221 kg·m/s
v_f = p/(m₁+m₂) = 0,221/0,34 = 0,65 m/s
Eₖ,f = ½·0,34·0,65² = 0,0718 J
**ΔE dissipata = 0,654 − 0,0718 ≈ 0,582 J**

---

### Problema 2 — Condensatore con dielettrico (batteria collegata)

**Dati:** A=2,50·10⁻³ m², d=8,0·10⁻⁴ m, V₀=60,0 V, εr=4,2

**(a) Prima del dielettrico**
C₀ = ε₀A/d = 8,854·10⁻¹²·2,50·10⁻³/8,0·10⁻⁴ = **2,77·10⁻¹¹ F**
Q₀ = C₀V₀ = **1,66·10⁻⁹ C**
U₀ = ½C₀V₀² = **4,98·10⁻⁸ J**

**(b) Con dielettrico (batteria sempre collegata, V resta 60 V)**
C' = εrC₀ = 4,2·2,77·10⁻¹¹ = **1,16·10⁻¹⁰ F**
Q' = C'V₀ = **6,97·10⁻⁹ C**
U' = ½C'V₀² = **2,09·10⁻⁷ J**

**(c) Lavoro della batteria e commento**
ΔQ = Q'−Q₀ = 5,31·10⁻⁹ C, fornita dalla batteria a V₀ costante:
W_batteria = V₀·ΔQ = **3,19·10⁻⁷ J**
ΔU_condensatore = U'−U₀ = **1,59·10⁻⁷ J**

La batteria fornisce più energia (3,19·10⁻⁷ J) di quanta ne resti immagazzinata nel condensatore (1,59·10⁻⁷ J): la differenza (≈1,59·10⁻⁷ J) corrisponde al lavoro meccanico legato all'inserimento del dielettrico, che viene "risucchiato" dentro le armature dalla forza elettrostatica attrattiva. Se il dielettrico è libero, questa energia diventa energia cinetica; se viene inserito lentamente e trattenuto, un agente esterno deve compiere lavoro negativo pari a questa quantità.

---

### Problema 3 — Selettore di velocità (Forza di Lorentz)

**Dati:** E=1,50·10³ N/C ĵ, B=0,220 T k̂, particella positiva non deflessa

**(a) Velocità della particella**
Condizione di equilibrio: qE + q(v×B) = 0 → E = −(v×B)

Con v = v_x î: v×B = v_x B(î×k̂) = −v_x B ĵ
Quindi: E = v_x B ĵ → v_x = E/B = 1500/0,220 = **6,82·10³ m/s**, diretta lungo **+x**

Questo è il classico **selettore di velocità**: solo le particelle con v=E/B attraversano indeflesse, indipendentemente dalla carica.

**(b) Energia cinetica del protone**
Eₖ = ½m_p v² = 0,5·1,673·10⁻²⁷·(6,82·10³)² = **3,89·10⁻²⁰ J** (≈0,24 eV)

**(c) Se la particella fosse un elettrone**
La condizione di equilibrio E+v×B=0 **non dipende dal segno o dal valore della carica** (q si semplifica nell'equazione qE+qv×B=0, purché q≠0). Quindi **anche l'elettrone** attraverserebbe la regione indeflesso, con la stessa identica velocità v=6,82·10³ m/s in direzione +x. Il selettore di velocità funziona per qualunque carica, non solo per una specifica massa o segno.

*(Nota aggiuntiva, non richiesta esplicitamente: se il campo elettrico venisse raddoppiato lasciando invariata v, l'equilibrio si romperebbe: la forza netta diventerebbe F=qE₀ residua, deviando la particella nella direzione di E.)*

---

## Riepilogo rapido dei risultati numerici

| # | Risultato |
|---|---|
| S1-P1 | r₁=0,286m, r₂=0,164m; T≈9,68N; v₁≈1,86 m/s, v₂≈1,06 m/s |
| S1-P2 | E≈1,81·10⁵ N/C; v≈7,61 m/s; tra le lastre E raddoppia a 3,62·10⁵ N/C |
| S1-P3 | B≈5,67·10⁻⁴ T; F≈1,90·10⁻¹⁶ N; lavoro nullo |
| S2-P1 | v≈4,85 m/s; N≈98,1 N; nessun distacco (caso limite) |
| S2-P2 | I=0,60A (rami), 0A (centrale); P=10,8W; invertendo: 0,20A/0,20A/0,40A |
| S2-P3 | ε≈23,0 mV; I≈0,128A; P≈2,95 mW |
| S3-P1 | v₁'=−1,2 m/s, v₂'=2,5 m/s; ΔE≈0,582 J (urto anelastico) |
| S3-P2 | C₀≈27,7 pF; U₀≈49,8 nJ; U'≈209 nJ; W_batt≈319 nJ |
| S3-P3 | v≈6,82·10³ m/s (+x); Eₖ≈3,89·10⁻²⁰ J; stesso risultato per l'elettrone |
