package database

type MemoryDatabase struct {
    urlMap map[ShortURL]LongURL;
}

func NewMemoryDatabase() *MemoryDatabase {
    return &MemoryDatabase{
        urlMap: make(map[ShortURL]LongURL),
    }
}

func (db *MemoryDatabase) CreateUrl(longUrl LongURL) (ShortURL, error) {
    shortUrl := ShortURL(ConvertUrl([]byte(longUrl)))

    _, ok := db.urlMap[shortUrl]
    if ok {
        return "", &AlreadyExists{longUrl}
    }

    db.urlMap[shortUrl] = longUrl
    return shortUrl, nil
}

func (db *MemoryDatabase) GetUrl(shortUrl ShortURL) (LongURL, error) {
    longUrl, ok := db.urlMap[shortUrl]
    if !ok {
        return "", &NotFound{shortUrl}
    }

    return longUrl, nil
}
