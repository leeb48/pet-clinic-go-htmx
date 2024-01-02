ALTER TABLE
    pets DROP FOREIGN KEY FK_pets_petTypes_petTypeId;

ALTER TABLE
    pets DROP COLUMN petTypeId;

ALTER TABLE
    pets DROP FOREIGN KEY FK_pets_owner_ownerId;

ALTER TABLE
    pets DROP COLUMN ownerId;