CREATE TABLE Tasks 
(
    Id serial not null unique,
    Method varchar(6),
    Url varchar(255),
    Headers varchar(255),
    Body varchar(255),
    Status varchar(255),
    HttpStatusCode varchar(255),
    Length varchar(255)
);