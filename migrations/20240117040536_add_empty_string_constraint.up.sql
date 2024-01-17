ALTER TABLE
    owners
ADD
    CONSTRAINT firstNameNotEmpty CHECK (firstName <> ''),
ADD
    CONSTRAINT lastNameNotEmpty CHECK (lastName <> ''),
ADD
    CONSTRAINT addrNotEmpty CHECK (address <> ''),
ADD
    CONSTRAINT stateNotEmpty CHECK (state <> ''),
ADD
    CONSTRAINT cityNotEmpty CHECK (city <> ''),
ADD
    CONSTRAINT phoneNotEmpty CHECK (phone <> ''),
ADD
    CONSTRAINT emailNotEmpty CHECK (email <> '');

ALTER TABLE
    petTypes
ADD
    CONSTRAINT petTypeNotEmpty CHECK (name <> '');

ALTER TABLE
    pets
ADD
    CONSTRAINT petNameNotEmpty CHECK (name <> '');