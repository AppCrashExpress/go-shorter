CREATE TABLE urls (
    id        SERIAL        PRIMARY KEY, 
    short_url VARCHAR(256)  NOT NULL UNIQUE, 
    long_url  VARCHAR(1024) NOT NULL
);
