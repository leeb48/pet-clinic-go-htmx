CREATE TABLE owners (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    firstName VARCHAR(255) NOT NULL,
    lastName VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    state VARCHAR(255) NOT NULL,
    city VARCHAR(255) NOT NULL,
    phone VARCHAR(10) NOT NULL,
    email VARCHAR(255) NOT NULL,
    birthdate DATE NOT NULL,
    created DATETIME NOT NULL
);

CREATE TABLE pets (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    birthdate DATE NOT NULL,
    created DATETIME NOT NULL
);

CREATE TABLE petTypes (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL
);

ALTER TABLE
    petTypes
ADD
    CONSTRAINT petType_uc_name UNIQUE (name);

ALTER TABLE
    pets
ADD
    petTypeId INTEGER,
ADD
    CONSTRAINT FK_pets_petTypes_petTypeId FOREIGN KEY (petTypeId) REFERENCES petTypes (id);

ALTER TABLE
    pets
ADD
    ownerId INTEGER,
ADD
    CONSTRAINT FK_pets_owner_ownerId Foreign KEY (ownerId) REFERENCES owners(id);

INSERT INTO
    petTypes (name)
VALUES
    ('DOG');