-- -----------------------------------------------------
-- Table USER
-- -----------------------------------------------------
CREATE TABLE "user" (
  ID SERIAL NOT NULL,
  NAME VARCHAR(50) NOT NULL,
  EMAIL VARCHAR(50) NOT NULL,
  CREATION_DATE TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UPDATE_DATE TIMESTAMP NULL,
  PRIMARY KEY (ID)
);

CREATE UNIQUE INDEX email_unique ON "user" (email);

-- -----------------------------------------------------
-- Table ACCOUNT
-- -----------------------------------------------------
CREATE TABLE account (
  ID SERIAL NOT NULL,
  NUMBER VARCHAR(50) NOT NULL,
  BALANCE FLOAT(15) NOT NULL,
  CURRENCY VARCHAR(2) NOT NULL,
  OWNER VARCHAR(50) NOT NULL,
  CREATED_AT TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UPDATED_AT TIMESTAMP NULL,
  PRIMARY KEY (ID),
  UNIQUE(NUMBER, OWNER)
);

ALTER TABLE account ADD CONSTRAINT fk_owner_user_id foreign key (OWNER) REFERENCES "user" (email);

-- -----------------------------------------------------
-- Table MOVEMENT
-- -----------------------------------------------------
CREATE TABLE movement (
  ID SERIAL NOT NULL,
  ORIGIN VARCHAR(50) NOT NULL,
  DESTINATION VARCHAR(50) NOT NULL,
  AMOUNT FLOAT(15) NOT NULL,
  CREATED_AT TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (ID)
);
