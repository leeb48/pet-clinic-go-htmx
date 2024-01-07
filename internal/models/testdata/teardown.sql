ALTER TABLE
    pets DROP FOREIGN KEY FK_pets_petTypes_petTypeId;

ALTER TABLE
    pets DROP COLUMN petTypeId;

ALTER TABLE
    pets DROP FOREIGN KEY FK_pets_owner_ownerId;

ALTER TABLE
    pets DROP COLUMN ownerId;

ALTER TABLE
    petTypes DROP CONSTRAINT petType_uc_name;

DROP TABLE owners;

DROP TABLE pets;

DROP TABLE petTypes;