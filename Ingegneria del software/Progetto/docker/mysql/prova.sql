-- ========================================================
-- DATI DIMOSTRATIVI (ambiente di sviluppo/demo)
-- Eseguito DOPO init.sql dall'entrypoint di MySQL.
-- Le password sono hash bcrypt; credenziali documentate nel README.
-- ========================================================

USE gialosai;

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
    );

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
    (12, 12, 18);

INSERT INTO
    OrdiniEntrata (
        id,
        numero,
        fornitore,
        data,
        stato,
        note,
        created_by
    )
VALUES (
        1,
        'ORD-ENT-2026-001',
        'TechDistruzione Srl',
        '2026-01-10 09:30:00',
        'evaso',
        'Fornitura hardware inizio anno',
        2
    ),
    (
        2,
        'ORD-ENT-2026-002',
        'Cartiera Industriale SpA',
        '2026-01-15 11:00:00',
        'evaso',
        'Rifornimento cancelleria e imballi',
        4
    ),
    (
        3,
        'ORD-ENT-2026-003',
        'Torrefazione Lombarda',
        '2026-01-22 14:15:00',
        'evaso',
        'Lotto caffè mensile',
        2
    ),
    (
        4,
        'ORD-ENT-2026-004',
        'Utensileria Bresciana',
        '2026-02-02 10:00:00',
        'evaso',
        'Attrezzatura da laboratorio',
        6
    ),
    (
        5,
        'ORD-ENT-2026-005',
        'Detergenti Italia Srl',
        '2026-02-10 16:30:00',
        'evaso',
        'Igiene stock primo trimestre',
        8
    ),
    (
        6,
        'ORD-ENT-2026-006',
        'Office Style Design',
        '2026-02-18 08:45:00',
        'evaso',
        'Sedie ergonomiche e arredo',
        2
    ),
    (
        7,
        'ORD-ENT-2026-007',
        'Textile Imports Srl',
        '2026-03-01 12:20:00',
        'annullato',
        'Capi abbigliamento neutri',
        4
    ),
    (
        8,
        'ORD-ENT-2026-008',
        'Global Electronics Co.',
        '2026-03-10 15:10:00',
        'annullato',
        'Schede audio e cavetteria',
        6
    ),
    (
        9,
        'ORD-ENT-2026-009',
        'Imballaggi Express',
        '2026-03-15 09:00:00',
        'confermato',
        'Fornitura scatole cartone',
        8
    ),
    (
        10,
        'ORD-ENT-2026-010',
        'ElettroDomus SpA',
        '2026-03-20 11:30:00',
        'confermato',
        'Bollitori e piccoli elettrodomestici',
        2
    );

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
    (13, 10, 12, 15, 12.00);

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
    );

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
    );

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
    (15, 10, 10, 150, 1.80);

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
        'Vendita',
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
        'Vendita',
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
        'Vendita',
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
        'Vendita',
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
        'Vendita',
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
        'Vendita',
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
        'Vendita',
        7,
        719.96,
        158.39,
        878.35,
        '2026-03-03 12:00:00',
        '/storage/fatture/fat_2026_007.pdf'
    ),
    (
        8,
        'FAT-2026-008',
        2026,
        'Vendita',
        8,
        520.00,
        114.40,
        634.40,
        '2026-03-09 13:30:00',
        '/storage/fatture/fat_2026_008.pdf'
    ),
    (
        9,
        'FAT-2026-009',
        2026,
        'Vendita',
        9,
        72.50,
        15.95,
        88.45,
        '2026-03-13 15:10:00',
        '/storage/fatture/fat_2026_009.pdf'
    ),
    (
        10,
        'FAT-2026-010',
        2026,
        'Vendita',
        10,
        270.00,
        59.40,
        329.40,
        '2026-03-19 09:45:00',
        '/storage/fatture/fat_2026_010.pdf'
    ),
    (
        11,
        'FAT-2026-011',
        2026,
        'Vendita',
        1,
        629.67,
        138.53,
        768.20,
        '2026-03-20 10:00:00',
        '/storage/fatture/fat_2026_011.pdf'
    ),
    (
        12,
        'FAT-2026-012',
        2026,
        'Vendita',
        2,
        617.98,
        135.96,
        753.94,
        '2026-03-21 09:30:00',
        '/storage/fatture/fat_2026_012.pdf'
    ),
    (
        13,
        'FAT-2026-013',
        2026,
        'Vendita',
        3,
        199.00,
        43.78,
        242.78,
        '2026-03-22 11:15:00',
        '/storage/fatture/fat_2026_013.pdf'
    ),
    (
        14,
        'FAT-2026-014',
        2026,
        'Vendita',
        4,
        130.00,
        28.60,
        158.60,
        '2026-03-23 14:00:00',
        '/storage/fatture/fat_2026_014.pdf'
    ),
    (
        15,
        'FAT-2026-015',
        2026,
        'Vendita',
        5,
        90.00,
        19.80,
        109.80,
        '2026-03-24 16:45:00',
        '/storage/fatture/fat_2026_015.pdf'
    );

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
        'in attesa',
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
        'in attesa',
        '2026-02-28 10:45:00',
        5
    );

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