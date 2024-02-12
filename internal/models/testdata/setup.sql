CREATE TABLE owners (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    firstName VARCHAR(255) NOT NULL CHECK (firstName <> ''),
    lastName VARCHAR(255) NOT NULL CHECK (lastName <> ''),
    address VARCHAR(255) NOT NULL CHECK (address <> ''),
    state VARCHAR(255) NOT NULL CHECK (state <> ''),
    city VARCHAR(255) NOT NULL CHECK (city <> ''),
    phone VARCHAR(10) NOT NULL CHECK (phone <> ''),
    email VARCHAR(255) NOT NULL CHECK (email <> ''),
    birthdate DATE NOT NULL,
    created DATETIME NOT NULL,
    modifiedDate DATETIME
);

CREATE TABLE petTypes (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL CHECK (name <> ''),
    CONSTRAINT petType_uc_name UNIQUE (name)
);

CREATE TABLE pets (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL CHECK (name <> ''),
    birthdate DATE NOT NULL,
    created DATETIME NOT NULL,
    petTypeId INTEGER NOT NULL,
    ownerId INTEGER NOT NULL,
    modifiedDate DATETIME,
    CONSTRAINT FK_pets_petTypes_petTypeId FOREIGN KEY (petTypeId) REFERENCES petTypes (id),
    CONSTRAINT FK_pets_owner_ownerId Foreign KEY (ownerId) REFERENCES owners(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Insert test owner
INSERT INTO
    owners (
        firstName,
        lastName,
        address,
        state,
        city,
        phone,
        email,
        birthdate,
        created
    )
VALUES
    (
        'Bong',
        'LEE',
        '123 St',
        'NV',
        'Las Vegas',
        '1112223333',
        'test@test.com',
        '1991-12-12',
        UTC_TIMESTAMP()
    );

INSERT INTO
    petTypes (name)
VALUES
    ('DOG'),
    ('CAT');

-- Insert test pets
INSERT INTO
    pets (name, birthdate, petTypeId, ownerId, created)
VALUES
    ('Mango', "2020-01-01", 1, 1, UTC_TIMESTAMP()),
    ('Acorn', UTC_TIMESTAMP(), 1, 1, UTC_TIMESTAMP());