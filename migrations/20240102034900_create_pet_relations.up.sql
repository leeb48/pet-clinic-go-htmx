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
    CONSTRAINT FK_pets_owner_ownerId Foreign KEY (ownerId) REFERENCES owners(id) ON DELETE CASCADE ON UPDATE CASCADE;