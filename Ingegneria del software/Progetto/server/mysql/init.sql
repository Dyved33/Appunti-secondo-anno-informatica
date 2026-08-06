CREATE DATABASE IF NOT EXISTS gialosai;
USE gialosai;

CREATE TABLE IF NOT EXISTS Utenti (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    ruolo VARCHAR(50) NOT NULL,
    domanda_sicurezza VARCHAR(255) NOT NULL,
    risposta_hash VARCHAR(255) NOT NULL,
    tentativi_falliti INT NOT NULL DEFAULT 0,
    bloccato BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS Categorie (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nome VARCHAR(100) NOT NULL UNIQUE,
    descrizione TEXT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS Prodotti (
    id INT AUTO_INCREMENT PRIMARY KEY,
    codice VARCHAR(100) NOT NULL UNIQUE,
    barcode VARCHAR(100) NOT NULL UNIQUE,
    nome VARCHAR(150) NOT NULL,
    descrizione TEXT NULL,
    categoria_id INT NULL,
    prezzo_acquisto DECIMAL(10,2) NOT NULL DEFAULT 0,
    prezzo_vendita DECIMAL(10,2) NOT NULL DEFAULT 0,
    soglia_minima INT NOT NULL DEFAULT 0,
    CONSTRAINT fk_prodotti_categoria
        FOREIGN KEY (categoria_id) REFERENCES Categorie(id)
        ON UPDATE CASCADE
        ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS Giacenze (
    id INT AUTO_INCREMENT PRIMARY KEY,
    prodotto_id INT NOT NULL UNIQUE,
    quantita INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_giacenze_prodotto
        FOREIGN KEY (prodotto_id) REFERENCES Prodotti(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS OrdiniEntrata (
    id INT AUTO_INCREMENT PRIMARY KEY,
    numero VARCHAR(50) NOT NULL UNIQUE,
    fornitore VARCHAR(150) NOT NULL,
    data DATETIME NOT NULL,
    stato VARCHAR(30) NOT NULL,
    note TEXT NULL,
    created_by INT NOT NULL,
    CONSTRAINT fk_ordine_entrata_created_by
        FOREIGN KEY (created_by) REFERENCES Utenti(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS RigheOrdineEntrata (
    id INT AUTO_INCREMENT PRIMARY KEY,
    ordine_entrata_id INT NOT NULL,
    prodotto_id INT NOT NULL,
    quantita INT NOT NULL,
    prezzo_unitario DECIMAL(10,2) NOT NULL DEFAULT 0,
    CONSTRAINT fk_righe_ordine_entrata_ordine
        FOREIGN KEY (ordine_entrata_id) REFERENCES OrdiniEntrata(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT fk_righe_ordine_entrata_prodotto
        FOREIGN KEY (prodotto_id) REFERENCES Prodotti(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS OrdiniUscita (
    id INT AUTO_INCREMENT PRIMARY KEY,
    numero VARCHAR(50) NOT NULL UNIQUE,
    cliente VARCHAR(150) NOT NULL,
    data DATETIME NOT NULL,
    stato VARCHAR(30) NOT NULL,
    note TEXT NULL,
    created_by INT NOT NULL,
    ddt_id INT NULL UNIQUE,
    CONSTRAINT fk_ordine_uscita_created_by
        FOREIGN KEY (created_by) REFERENCES Utenti(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS DDT (
    id INT AUTO_INCREMENT PRIMARY KEY,
    numero VARCHAR(50) NOT NULL UNIQUE,
    ordine_uscita_id INT NOT NULL UNIQUE,
    data DATETIME NOT NULL,
    mittente VARCHAR(150) NOT NULL,
    destinatario VARCHAR(150) NOT NULL,
    path_pdf VARCHAR(255) NOT NULL,
    CONSTRAINT fk_ddt_ordine_uscita
        FOREIGN KEY (ordine_uscita_id) REFERENCES OrdiniUscita(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS RigheOrdineUscita (
    id INT AUTO_INCREMENT PRIMARY KEY,
    ordine_uscita_id INT NOT NULL,
    prodotto_id INT NOT NULL,
    quantita INT NOT NULL,
    prezzo_unitario DECIMAL(10,2) NOT NULL DEFAULT 0,
    CONSTRAINT fk_righe_ordine_uscita_ordine
        FOREIGN KEY (ordine_uscita_id) REFERENCES OrdiniUscita(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT fk_righe_ordine_uscita_prodotto
        FOREIGN KEY (prodotto_id) REFERENCES Prodotti(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS Fatture (
    id INT AUTO_INCREMENT PRIMARY KEY,
    numero VARCHAR(50) NOT NULL UNIQUE,
    anno INT NOT NULL,
    tipo VARCHAR(30) NOT NULL,
    ordine_uscita_id INT NOT NULL,
    imponibile DECIMAL(10,2) NOT NULL DEFAULT 0,
    iva DECIMAL(10,2) NOT NULL DEFAULT 0,
    totale DECIMAL(10,2) NOT NULL DEFAULT 0,
    data DATETIME NOT NULL,
    path_pdf VARCHAR(255) NOT NULL,
    CONSTRAINT fk_fatture_ordine_uscita
        FOREIGN KEY (ordine_uscita_id) REFERENCES OrdiniUscita(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS Resi (
    id INT AUTO_INCREMENT PRIMARY KEY,
    ordine_uscita_id INT NOT NULL,
    prodotto_id INT NOT NULL,
    quantita INT NOT NULL,
    motivo TEXT NOT NULL,
    stato VARCHAR(30) NOT NULL,
    data DATETIME NOT NULL,
    created_by INT NOT NULL,
    CONSTRAINT fk_resi_ordine_uscita
        FOREIGN KEY (ordine_uscita_id) REFERENCES OrdiniUscita(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT fk_resi_prodotto
        FOREIGN KEY (prodotto_id) REFERENCES Prodotti(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,
    CONSTRAINT fk_resi_created_by
        FOREIGN KEY (created_by) REFERENCES Utenti(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS Sessioni (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    token_hash VARCHAR(64) NOT NULL UNIQUE,
    expires_at DATETIME NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_sessioni_user
        FOREIGN KEY (user_id) REFERENCES Utenti(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;