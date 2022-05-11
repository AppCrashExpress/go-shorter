package database

import "fmt"

type LongURL  string;
type ShortURL string;

type AlreadyExists struct {
    longUrl LongURL;
}
func (m *AlreadyExists) Error() string {
    return fmt.Sprintf("URL '%s' already exists", m.longUrl)
}

type NotFound struct {
    shortUrl ShortURL;
}
func (m *NotFound) Error() string {
    return fmt.Sprintf("URL '%s' not found", m.shortUrl)
}

type Database interface {
    CreateUrl(LongURL) (ShortURL, error);
    GetUrl(ShortURL) (LongURL, error);
}
