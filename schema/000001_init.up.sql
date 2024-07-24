CREATE TABLE Tasks 
(
    Id serial primary key,
    Method varchar(6),
    Url varchar(255),
    Headers varchar(255),
    Body varchar(255)
);

CREATE TABLE TaskStatus
(
    Id serial,
    Status varchar(255),
    HttpStatusCode varchar(255),
    Length varchar(255),
    FOREIGN KEY (Id) REFERENCES Tasks (Id)
)