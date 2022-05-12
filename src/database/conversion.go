package database

import (
	"crypto/sha1"
	"encoding/base64"
)

const sUrlMaxLen = 10

// Hacky algortihm
// Conversion to SHA1 insures uniqness on smallest changes (e.g. appended /1 to end)
// Conversion to Base64 is closest to needed base63

func ConvertUrl(longUrl []byte) []byte {
    sha1Url :=  sha1.Sum(longUrl)

	base64Url := make([]byte, base64.StdEncoding.EncodedLen(len(sha1Url)))
    base64.URLEncoding.Encode(base64Url, sha1Url[:])

    resUrl := make([]byte, 0, sUrlMaxLen)
    for _, v := range base64Url {
        if v == '-' {
            continue
        }
        if len(resUrl) == sUrlMaxLen {
            break
        }

        resUrl = append(resUrl, v)
    }

    return resUrl
}
