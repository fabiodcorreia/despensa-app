CREATE TABLE Location (
  Id VARCHAR(10) PRIMARY KEY NOT NULL, 
  Name VARCHAR(24) NOT NULL UNIQUE
); 

CREATE TABLE Item (
  Id VARCHAR(10) PRIMARY KEY NOT NULL, 
  Name VARCHAR(24) NOT NULL UNIQUE
); 

CREATE TABLE ItemStored ( 
  LocationId VARCHAR(10) NOT NULL, 
  ItemId VARCHAR(10) NOT NULL, 
  Quantity INTEGER NOT NULL DEFAULT 0,
  StoredAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (LocationId, ItemId), 
  FOREIGN KEY (LocationId) REFERENCES Location(Id) ON DELETE CASCADE, 
  FOREIGN KEY (ItemId) REFERENCES Item(Id) ON DELETE CASCADE
);
