ALTER TABLE
    owners DROP CONSTRAINT firstNameNotEmpty,
    DROP CONSTRAINT lastNameNotEmpty,
    DROP CONSTRAINT addrNotEmpty,
    DROP CONSTRAINT stateNotEmpty,
    DROP CONSTRAINT cityNotEmpty,
    DROP CONSTRAINT phoneNotEmpty,
    DROP CONSTRAINT emailNotEmpty;

ALTER TABLE
    petTypes DROP CONSTRAINT petTypeNotEmpty;

ALTER TABLE
    pets DROP CONSTRAINT petNameNotEmpty;