package database

import (
    "database/sql"
    "errors"

	_ "github.com/jackc/pgx/v4/stdlib"
    "github.com/jackc/pgconn"
)

type LinkPair struct {
    id       int
    shortUrl string
    longUrl  string
}

type PgDatabase struct {
    db *sql.DB
}

func NewPgDatabase(dbUrl string, idleConn, maxConn int) (*PgDatabase, error) {
    db, err := sql.Open("pgx", dbUrl)
    if err != nil {
        return nil, err
    }

    err = db.Ping()
    if err != nil {
        return nil, err
    }

    db.SetMaxIdleConns(idleConn)
    db.SetMaxOpenConns(maxConn)

    return &PgDatabase{db}, nil
}

func (pg *PgDatabase) Close() error {
    return pg.db.Close()
}

func (pg *PgDatabase) CreateUrl(longUrl LongURL) (ShortURL, error) {
    shortUrl := ShortURL(ConvertUrl([]byte(longUrl)))
    _, err := pg.db.Exec("INSERT INTO urls (short_url, long_url) VALUES ($1, $2)", shortUrl, longUrl)

    if err != nil {
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) && pgErr.Code == "23505" {
            return "", &AlreadyExists{longUrl}
        }

        return "", err
    }

    return shortUrl, nil
}

func (pg *PgDatabase) GetUrl(shortUrl ShortURL) (LongURL, error) {
    row := pg.db.QueryRow("SELECT long_url FROM urls WHERE short_url = $1", shortUrl)
    
    var longUrl string
    if err := row.Scan(&longUrl); err != nil {
        if err == sql.ErrNoRows {
            return "", &NotFound{shortUrl}
        }

        return "", err
    }

    return LongURL(longUrl), nil
}
