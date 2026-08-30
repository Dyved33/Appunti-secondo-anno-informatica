-- ========================================================
-- DATI DIMOSTRATIVI (ambiente di sviluppo/demo)
-- Eseguito DOPO init.sql dall'entrypoint di MySQL.
-- Le password sono hash bcrypt; credenziali documentate nel README.
-- ========================================================

USE gialosai;

-- ========================================================
-- UTENTI
-- ========================================================
INSERT INTO
Utenti (
id,
username,
password_hash,
ruolo,
tentativi_falliti,
bloccato
)
VALUES (
1,
'admin',
'$2a$10$J2aTcZKE9z7lr3fgWpkfVu8Rp66ardtfXGMLYRtsNx2arTXsm6mIK',
'admin',
0,
FALSE
),
(
2,
'magazziniere',
'$2a$10$qv3wmYCl5EUz9MEN1EkKWu50/3W.PO9CKEqXSuRJNsLemxptiLAiy',
'magazziniere',
0,
FALSE
),
(
3,
'contabile',
'$2a$10$6eK7WXcyBJoPoVhNI4LahuHx3WM7tjMjqpKxaH1R269sDIzGHEgV2',
'contabile',
0,
FALSE
),
(
4,
'mario.rossi',
'$2a$10$lBS2hnb4FA8NEjtFCuOuW.4IHAL.QZoxGqPVd6EVXwBnxv2caej7C',
'magazziniere',
0,
FALSE
),
(
5,
'laura.bianchi',
'$2a$10$Su2.Es0iu6giw5/R88rxZueQ11EEswmrvTQcketLHG1YBS7XPSacG',
'contabile',
0,
FALSE
),
(
6,
'giuseppe.verdi',
'$2a$10$lBS2hnb4FA8NEjtFCuOuW.4IHAL.QZoxGqPVd6EVXwBnxv2caej7C',
'magazziniere',
0,
FALSE
),
(
7,
'anna.neri',
'$2a$10$jMYr0RTKjJqZ8TvZNp8r0O5cGgye6NKu66QGgCVeYQiyuy53/7pb6',
'admin',
0,
FALSE
),
(
8,
'paolo.gialli',
'$2a$10$lBS2hnb4FA8NEjtFCuOuW.4IHAL.QZoxGqPVd6EVXwBnxv2caej7C',
'magazziniere',
0,
FALSE
),
(
9,
'elena.viola',
'$2a$10$Su2.Es0iu6giw5/R88rxZueQ11EEswmrvTQcketLHG1YBS7XPSacG',
'contabile',
0,
FALSE
),
(
10,
'stefano.bruni',
'$2a$10$lBS2hnb4FA8NEjtFCuOuW.4IHAL.QZoxGqPVd6EVXwBnxv2caej7C',
'magazziniere',
0,
FALSE
);

-- ========================================================
-- CATEGORIE
-- ========================================================
INSERT INTO
Categorie (id, nome, descrizione)
VALUES (
1,
'Elettronica',
'Smartphone, computer e componenti hardware'
),
(
2,
'Abbigliamento',
'Capi di vestiario, calzature e accessori'
),
(
3,
'Alimentari',
'Prodotti alimentari secchi e confezionati'
),
(
4,
'Cancelleria',
'Articoli per ufficio, carta e scrittura'
),
(
5,
'Utensileria',
'Attrezzi da lavoro e minuteria metallica'
),
(
6,
'Informatica',
'Periferiche pc, networking e cavetteria'
),
(
7,
'Elettrodomestici',
'Piccoli e grandi elettrodomestici da casa'
),
(
8,
'Igiene e Pulizia',
'Detergenti per ambienti e cura della persona'
),
(
9,
'Arredamento Ufficio',
'Sedie ergonomiche, scrivanie e lampade'
),
(
10,
'Packaging e Imballaggio',
'Scatole di cartone, nastro adesivo e pluriball'
);

-- ========================================================
-- PRODOTTI
-- ========================================================
INSERT INTO
Prodotti (
id,
codice,
barcode,
nome,
descrizione,
categoria_id,
prezzo_acquisto,
prezzo_vendita,
soglia_minima
)
VALUES (
1,
'PROD001',
'8001234500016',
'Monitor LED 24"',
'Monitor Full HD IPS 75Hz',
1,
95.00,
149.99,
10
),
(
2,
'PROD002',
'8001234500023',
'Tastiera Meccanica RGB',
'Layout Italiano, switch rossi',
6,
35.00,
59.90,
15
),
(
3,
'PROD003',
'8001234500030',
'Mouse Ottico Wireless',
'Mouse ergonomico 1600 DPI',
6,
8.50,
18.00,
20
),
(
4,
'PROD004',
'8001234500047',
'Caffè in Grani 1kg',
'Miscela 80% Arabica 20% Robusta',
3,
7.20,
14.50,
25
),
(
5,
'PROD005',
'8001234500054',
'Risanante Spray 500ml',
'Detergente igienizzante multiuso',
8,
1.80,
4.50,
30
),
(
6,
'PROD006',
'8001234500061',
'Sedia Ufficio Ergonomica',
'Supporto lombare e braccioli regolabili',
9,
85.00,
159.00,
5
),
(
7,
'PROD007',
'8001234500078',
'Cacciavite a Batteria 12V',
'Include set 20 inserti',
5,
22.00,
45.00,
8
),
(
8,
'PROD008',
'8001234500085',
'Risma Carta A4 80g 500ff',
'Carta bianca per stampanti e fotocopie',
4,
2.90,
5.20,
50
),
(
9,
'PROD009',
'8001234500092',
'Cavo HDMI 2.0 2m',
'Supporto 4K 60Hz con connettori placcati oro',
6,
2.10,
6.90,
40
),
(
10,
'PROD010',
'8001234500108',
'Scatola Cartone 40x30x30cm',
'Scatola ad alta resistenza per spedizioni',
10,
0.60,
1.80,
100
),
(
11,
'PROD011',
'8001234500115',
'Maglietta Cotone Nera L',
'100% cotone pettinato',
2,
3.50,
12.00,
30
),
(
12,
'PROD012',
'8001234500122',
'Bollitore Elettrico 1.7L',
'Corpo in acciaio inox 2200W',
7,
12.00,
24.90,
10
),
(
13,
'PROD013',
'8001234500139',
'Tablet 10" Wi-Fi',
'Display IPS 1280x800, 4GB RAM, 64GB',
1,
140.00,
219.00,
8
),
(
14,
'PROD014',
'8001234500146',
'Laptop 15.6" Full HD',
'Core i5, 8GB RAM, SSD 256GB',
1,
620.00,
899.00,
5
),
(
15,
'PROD015',
'8001234500153',
'Cuffie Bluetooth',
'Auricolari over-ear con cancellazione rumore',
1,
28.00,
54.90,
12
),
(
16,
'PROD016',
'8001234500160',
'Jeans Uomo Taglia 34',
'Denim elasticizzato, colore blu',
2,
15.00,
39.00,
20
),
(
17,
'PROD017',
'8001234500177',
'Scarpe Sportive Uomo',
'Tomaia in mesh, suola ammortizzata',
2,
22.00,
49.90,
15
),
(
18,
'PROD018',
'8001234500184',
'Olio Extravergine Oliva 1L',
'Spremuta a freddo, origine italiana',
3,
4.50,
8.90,
30
),
(
19,
'PROD019',
'8001234500191',
'Pasta di Semola 500g',
'Formato penne rigate',
3,
0.80,
1.40,
60
),
(
20,
'PROD020',
'8001234500207',
'Penne a Sfera blu 5pz',
'Scrittura scorrevole, punta media',
4,
1.20,
2.90,
50
),
(
21,
'PROD021',
'8001234500214',
'Cucitrice Metallica',
'Per ufficio, graffette 24/6',
4,
3.20,
7.50,
25
),
(
22,
'PROD022',
'8001234500221',
'Trapano Avvitatore 18V',
'Batteria 2Ah inclusa, mandrino 13mm',
5,
48.00,
89.00,
6
),
(
23,
'PROD023',
'8001234500238',
'Minifrigo 50L',
'Compressore silenzioso, doppio scomparto',
7,
85.00,
149.00,
4
),
(
24,
'PROD024',
'8001234500245',
'Detergente Pavimenti 1L',
'Profumazione agrumi',
8,
1.50,
3.90,
40
),
(
25,
'PROD025',
'8001234500252',
'Lampada da Scrivania LED',
'Regolazione luminosità e temperatura colore',
9,
18.00,
35.00,
10
),
(
26,
'PROD026',
'8001234500269',
'Sacchetti Imballo 100pz',
'Sacchetti di plastica trasparenti',
10,
1.10,
2.50,
80
);

-- ========================================================
-- GIACENZE
-- ========================================================
INSERT INTO
Giacenze (id, prodotto_id, quantita)
VALUES (1, 1, 5),
(2, 2, 10),
(3, 3, 15),
(4, 4, 60),
(5, 5, 200),
(6, 6, 12),
(7, 7, 25),
(8, 8, 350),
(9, 9, 95),
(10, 10, 500),
(11, 11, 75),
(12, 12, 18),
(13, 13, 12),
(14, 14, 3),
(15, 15, 20),
(16, 16, 35),
(17, 17, 8),
(18, 18, 50),
(19, 19, 120),
(20, 20, 80),
(21, 21, 30),
(22, 22, 5),
(23, 23, 6),
(24, 24, 45),
(25, 25, 9),
(26, 26, 150);

-- ========================================================
-- ORDINI DI ENTRATA
-- ========================================================
INSERT INTO
OrdiniEntrata (
id,
numero,
fornitore,
data,
note,
created_by
)
VALUES (
1,
'ORD-ENT-2026-001',
'TechDistruzione Srl',
'2026-01-10 09:30:00',
'Fornitura hardware inizio anno',
2
),
(
2,
'ORD-ENT-2026-002',
'Cartiera Industriale SpA',
'2026-01-15 11:00:00',
'Rifornimento cancelleria e imballi',
4
),
(
3,
'ORD-ENT-2026-003',
'Torrefazione Lombarda',
'2026-01-22 14:15:00',
'Lotto caffè mensile',
2
),
(
4,
'ORD-ENT-2026-004',
'Utensileria Bresciana',
'2026-02-02 10:00:00',
'Attrezzatura da laboratorio',
6
),
(
5,
'ORD-ENT-2026-005',
'Detergenti Italia Srl',
'2026-02-10 16:30:00',
'Igiene stock primo trimestre',
8
),
(
6,
'ORD-ENT-2026-006',
'Office Style Design',
'2026-02-18 08:45:00',
'Sedie ergonomiche e arredo',
2
),
(
7,
'ORD-ENT-2026-007',
'Textile Imports Srl',
'2026-03-01 12:20:00',
'Capi abbigliamento neutri',
4
),
(
8,
'ORD-ENT-2026-008',
'Global Electronics Co.',
'2026-03-10 15:10:00',
'Schede audio e cavetteria',
6
),
(
9,
'ORD-ENT-2026-009',
'Imballaggi Express',
'2026-03-15 09:00:00',
'Fornitura scatole cartone',
8
),
(
10,
'ORD-ENT-2026-010',
'ElettroDomus SpA',
'2026-03-20 11:30:00',
'Bollitori e piccoli elettrodomestici',
2
),
(
11,
'ORD-ENT-2026-011',
'TechElettro SpA',
'2026-04-02 09:00:00',
'Lotto pc, tablet e cuffie',
6
),
(
12,
'ORD-ENT-2026-012',
'Textile Imports Srl',
'2026-04-10 13:45:00',
'Abbigliamento e calzature primavera',
4
),
(
13,
'ORD-ENT-2026-013',
'Alimentari Rossi',
'2026-04-15 10:30:00',
'Derrate alimentari secche',
2
),
(
14,
'ORD-ENT-2026-014',
'Cartoleria Moderna',
'2026-04-20 15:00:00',
'Materiale di scrittura e cancelleria',
8
),
(
15,
'ORD-ENT-2026-015',
'Utensileria Bresciana',
'2026-04-25 08:30:00',
'Utensileria elettrica professionale',
6
),
(
16,
'ORD-ENT-2026-016',
'ElettroDomus SpA',
'2026-05-01 11:10:00',
'Piccoli elettrodomestici e imballi',
2
);

-- ========================================================
-- RIGHE ORDINI DI ENTRATA
-- ========================================================
INSERT INTO
RigheOrdineEntrata (
id,
ordine_entrata_id,
prodotto_id,
quantita,
prezzo_unitario
)
VALUES (1, 1, 1, 20, 95.00),
(2, 1, 2, 30, 35.00),
(3, 2, 8, 200, 2.90),
(4, 2, 10, 400, 0.60),
(5, 3, 4, 50, 7.20),
(6, 4, 7, 15, 22.00),
(7, 5, 5, 150, 1.80),
(8, 6, 6, 10, 85.00),
(9, 7, 11, 50, 3.50),
(10, 8, 3, 50, 8.50),
(11, 8, 9, 60, 2.10),
(12, 9, 10, 300, 0.60),
(13, 10, 12, 15, 12.00),
(14, 11, 14, 10, 620.00),
(15, 11, 13, 20, 140.00),
(16, 11, 15, 40, 28.00),
(17, 12, 16, 100, 15.00),
(18, 12, 17, 40, 22.00),
(19, 13, 18, 80, 4.50),
(20, 13, 19, 200, 0.80),
(21, 14, 20, 150, 1.20),
(22, 14, 21, 50, 3.20),
(23, 15, 22, 20, 48.00),
(24, 16, 23, 12, 85.00),
(25, 16, 24, 100, 1.50),
(26, 16, 25, 30, 18.00),
(27, 16, 26, 200, 1.10);

-- ========================================================
-- ORDINI DI USCITA
-- ========================================================
INSERT INTO
OrdiniUscita (
id,
numero,
cliente,
data,
stato,
note,
created_by,
ddt_id
)
VALUES (
1,
'ORD-USC-2026-001',
'Acme Corporation Srl',
'2026-01-18 10:00:00',
'evaso',
'Spedizione urgente con corriere espresso',
3,
NULL
),
(
2,
'ORD-USC-2026-002',
'Studio Associato Rossi',
'2026-01-25 14:30:00',
'evaso',
'Fornitura postazioni lavoro',
5,
NULL
),
(
3,
'ORD-USC-2026-003',
'Ristorante La Perla',
'2026-02-01 16:00:00',
'evaso',
'Fornitura caffè e detersivi',
3,
NULL
),
(
4,
'ORD-USC-2026-004',
'Libreria Centrale',
'2026-02-08 09:15:00',
'evaso',
'Carta risme e cancelleria varia',
9,
NULL
),
(
5,
'ORD-USC-2026-005',
'Officina Meccanica Bianchi',
'2026-02-14 11:45:00',
'evaso',
'Utensileria elettrica',
5,
NULL
),
(
6,
'ORD-USC-2026-006',
'Hotel Bellavista',
'2026-02-22 15:00:00',
'evaso',
'Elettrodomestici e detergenti',
3,
NULL
),
(
7,
'ORD-USC-2026-007',
'Beta Solutions SpA',
'2026-03-02 10:30:00',
'annullato',
'Attrezzatura pc per nuovi assunti',
9,
NULL
),
(
8,
'ORD-USC-2026-008',
'Copisteria Moderno',
'2026-03-08 12:00:00',
'annullato',
'Ordine carta A4 massivo',
5,
NULL
),
(
9,
'ORD-USC-2026-009',
'Bar Sport',
'2026-03-12 17:10:00',
'confermato',
'Consumabili e caffè',
3,
NULL
),
(
10,
'ORD-USC-2026-010',
'Gamma Logistica',
'2026-03-18 08:30:00',
'confermato',
'Scatole e materiale da imballo',
9,
NULL
),
(
11,
'ORD-USC-2026-011',
'Casa Elettronica',
'2026-04-05 09:45:00',
'evaso',
'Tablet e cuffie per punto vendita',
3,
NULL
),
(
12,
'ORD-USC-2026-012',
'Moda & Stile',
'2026-04-12 16:20:00',
'evaso',
'Capo abbigliamento e calzature',
5,
NULL
),
(
13,
'ORD-USC-2026-013',
'Supermercato Fresco',
'2026-04-18 10:00:00',
'confermato',
'Derrate alimentari in attesa caricamento',
3,
NULL
),
(
14,
'ORD-USC-2026-014',
'Uffici & Co',
'2026-04-25 11:30:00',
'annullato',
'Ordine cancellato dal cliente',
9,
NULL
),
(
15,
'ORD-USC-2026-015',
'Cantiere Edile Nord',
'2026-05-02 14:00:00',
'evaso',
'Utensileria per squadra operai',
5,
NULL
),
(
16,
'ORD-USC-2026-016',
'Ristorante Al Porto',
'2026-05-06 10:15:00',
'confermato',
'Consumabili e detergenti',
3,
NULL
);

-- ========================================================
-- RIGHE ORDINI DI USCITA
-- ========================================================
INSERT INTO
RigheOrdineUscita (
id,
ordine_uscita_id,
prodotto_id,
quantita,
prezzo_unitario
)
VALUES (1, 1, 1, 3, 149.99),
(2, 1, 2, 3, 59.90),
(3, 2, 6, 2, 159.00),
(4, 2, 1, 2, 149.99),
(5, 3, 4, 10, 14.50),
(6, 3, 5, 12, 4.50),
(7, 4, 8, 25, 5.20),
(8, 5, 7, 2, 45.00),
(9, 6, 12, 3, 24.90),
(10, 6, 5, 20, 4.50),
(11, 7, 1, 4, 149.99),
(12, 7, 3, 4, 18.00),
(13, 8, 8, 100, 5.20),
(14, 9, 4, 5, 14.50),
(15, 10, 10, 150, 1.80),
(16, 11, 13, 3, 219.00),
(17, 11, 15, 6, 54.90),
(18, 12, 16, 10, 39.00),
(19, 12, 17, 5, 49.90),
(20, 13, 18, 20, 8.90),
(21, 13, 19, 50, 1.40),
(22, 14, 20, 30, 2.90),
(23, 15, 22, 4, 89.00),
(24, 16, 4, 20, 14.50),
(25, 16, 5, 30, 4.50);

-- ========================================================
-- DDT
-- ========================================================
INSERT INTO
DDT (
id,
numero,
ordine_uscita_id,
data,
mittente,
destinatario,
path_pdf
)
VALUES (
1,
'DDT-2026-001',
1,
'2026-01-18 11:30:00',
'Gialosai SpA',
'Acme Corporation Srl',
'/storage/ddt/ddt_2026_001.pdf'
),
(
2,
'DDT-2026-002',
2,
'2026-01-25 15:45:00',
'Gialosai SpA',
'Studio Associato Rossi',
'/storage/ddt/ddt_2026_002.pdf'
),
(
3,
'DDT-2026-003',
3,
'2026-02-01 17:15:00',
'Gialosai SpA',
'Ristorante La Perla',
'/storage/ddt/ddt_2026_003.pdf'
),
(
4,
'DDT-2026-004',
4,
'2026-02-08 10:30:00',
'Gialosai SpA',
'Libreria Centrale',
'/storage/ddt/ddt_2026_004.pdf'
),
(
5,
'DDT-2026-005',
5,
'2026-02-14 13:00:00',
'Gialosai SpA',
'Officina Meccanica Bianchi',
'/storage/ddt/ddt_2026_005.pdf'
),
(
6,
'DDT-2026-006',
6,
'2026-02-22 16:30:00',
'Gialosai SpA',
'Hotel Bellavista',
'/storage/ddt/ddt_2026_006.pdf'
),
(
7,
'DDT-2026-007',
11,
'2026-04-05 10:30:00',
'Gialosai SpA',
'Casa Elettronica',
'/storage/ddt/ddt_2026_007.pdf'
),
(
8,
'DDT-2026-008',
12,
'2026-04-12 17:00:00',
'Gialosai SpA',
'Moda & Stile',
'/storage/ddt/ddt_2026_008.pdf'
),
(
9,
'DDT-2026-009',
15,
'2026-05-02 15:00:00',
'Gialosai SpA',
'Cantiere Edile Nord',
'/storage/ddt/ddt_2026_009.pdf'
);

-- ========================================================
-- FATTURE  (note di credito incluse)
-- imponibile, iva (22%) e totale calcolati sulle righe degli ordini evasi
-- ========================================================
INSERT INTO
Fatture (
id,
numero,
anno,
tipo,
ordine_uscita_id,
imponibile,
iva,
totale,
data,
path_pdf
)
VALUES (
1,
'FAT-2026-001',
2026,
'fattura',
1,
629.67,
138.53,
768.20,
'2026-01-19 10:00:00',
'/storage/fatture/fat_2026_001.pdf'
),
(
2,
'FAT-2026-002',
2026,
'fattura',
2,
617.98,
135.96,
753.94,
'2026-01-26 09:30:00',
'/storage/fatture/fat_2026_002.pdf'
),
(
3,
'FAT-2026-003',
2026,
'fattura',
3,
199.00,
43.78,
242.78,
'2026-02-02 11:15:00',
'/storage/fatture/fat_2026_003.pdf'
),
(
4,
'FAT-2026-004',
2026,
'fattura',
4,
130.00,
28.60,
158.60,
'2026-02-09 14:00:00',
'/storage/fatture/fat_2026_004.pdf'
),
(
5,
'FAT-2026-005',
2026,
'fattura',
5,
90.00,
19.80,
109.80,
'2026-02-15 16:45:00',
'/storage/fatture/fat_2026_005.pdf'
),
(
6,
'FAT-2026-006',
2026,
'fattura',
6,
164.70,
36.23,
200.93,
'2026-02-23 10:20:00',
'/storage/fatture/fat_2026_006.pdf'
),
(
7,
'FAT-2026-007',
2026,
'fattura',
11,
986.40,
217.01,
1203.41,
'2026-04-06 09:00:00',
'/storage/fatture/fat_2026_007.pdf'
),
(
8,
'FAT-2026-008',
2026,
'fattura',
12,
639.50,
140.69,
780.19,
'2026-04-13 10:30:00',
'/storage/fatture/fat_2026_008.pdf'
),
(
9,
'FAT-2026-009',
2026,
'fattura',
15,
356.00,
78.32,
434.32,
'2026-05-03 09:15:00',
'/storage/fatture/fat_2026_009.pdf'
),
(
10,
'NC-2026-001',
2026,
'nota_credito',
3,
199.00,
43.78,
242.78,
'2026-02-08 10:00:00',
'/storage/fatture/nc_2026_001.pdf'
),
(
11,
'NC-2026-002',
2026,
'nota_credito',
4,
130.00,
28.60,
158.60,
'2026-02-20 12:00:00',
'/storage/fatture/nc_2026_002.pdf'
);

-- ========================================================
-- RESI
-- ========================================================
INSERT INTO
Resi (
id,
ordine_uscita_id,
prodotto_id,
quantita,
motivo,
stato,
data,
created_by
)
VALUES (
1,
1,
2,
1,
'Tasto barra spaziatrice difettoso',
'accettato',
'2026-01-22 11:00:00',
3
),
(
2,
3,
5,
2,
'Flacone danneggiato durante il trasporto',
'rifiutato',
'2026-02-05 09:30:00',
5
),
(
3,
4,
8,
5,
'Confezione aperta e carte sgualcite',
'da_verificare',
'2026-02-12 14:10:00',
3
),
(
4,
5,
7,
1,
'Cacciavite non tiene la carica',
'accettato',
'2026-02-20 16:00:00',
9
),
(
5,
6,
12,
1,
'Bollitore non riscalda l acqua',
'da_verificare',
'2026-02-28 10:45:00',
5
),
(
6,
6,
5,
5,
'Confezioni di detergente bagnate',
'da_verificare',
'2026-03-04 09:00:00',
3
),
(
7,
11,
13,
1,
'Touchscreen rigato in spedizione',
'accettato',
'2026-04-10 15:30:00',
9
),
(
8,
12,
17,
1,
'Taglia non corrispondente all ordine',
'rifiutato',
'2026-04-16 10:20:00',
5
),
(
9,
3,
4,
2,
'Caffè in grani con odore di muffa',
'da_verificare',
'2026-04-20 11:00:00',
3
);

-- ========================================================
-- SESSIONI (solo riferimento; i token reali vengono emessi al login)
-- ========================================================
INSERT INTO
Sessioni (
id,
user_id,
token_hash,
expires_at
)
VALUES (
1,
1,
'e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855',
'2026-08-22 18:00:00'
),
(
2,
2,
'f2d81a260dea8a100dd517984e53c56a7523d96942a834b9cdc249bd4e6c7952',
'2026-08-22 18:00:00'
),
(
3,
3,
'a12b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b',
'2026-08-22 18:00:00'
);
