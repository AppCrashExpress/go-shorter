package database

import "testing"

type Testcase struct {
    input    string;
    output   string;
    descript string;
}

func TestCreation(t *testing.T) {
    urls := []LongURL {
        "",
        "some url",
        "http://localhost.com/",
    }

    memoryDatabase := NewMemoryDatabase()
    
    for _, v := range urls {
        shortUrl, err := memoryDatabase.CreateUrl(v)
        if err != nil {
            t.Errorf("Database failed prematurely when adding '%s'", v)
            continue
        }

        longUrl, err := memoryDatabase.GetUrl(shortUrl)
        if err != nil {
            t.Errorf("Database failed prematurely while getting '%s'", longUrl)
            continue
        }

        if longUrl != v {
            t.Errorf("'%s' is not equal '%s'", v, longUrl)
        }
    }
}

func TestRepeat(t *testing.T) {
    url := LongURL("http://localhost.com/")
    memoryDatabase := NewMemoryDatabase()

    _, err := memoryDatabase.CreateUrl(url)
    if err != nil {
        t.Fatalf("Database failed prematurely when adding '%s'", url)
    }

    _, err = memoryDatabase.CreateUrl(url)
    if err == nil {
        t.Fatalf("'%s' added twice", url)
    }
}
