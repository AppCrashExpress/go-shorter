package shorter

import (
	"context"
	"log"

	pb "github.com/AppCrashExpress/go-shorter/src/api"
	"github.com/AppCrashExpress/go-shorter/src/database"
)

type URLServer struct {
    pb.UnimplementedShortnerServer
    db database.Database
}

func NewServer(db database.Database) *URLServer {
    return &URLServer{db: db}
}

func (srv *URLServer) CreateNew(ctx context.Context, pbLongUrl *pb.LongURL) (*pb.ShortURL, error) {
    longUrl := database.LongURL(pbLongUrl.Lurl)
    shortUrl, err := srv.db.CreateUrl(longUrl)

    pbShortUrl := &pb.ShortURL{Surl: string(shortUrl)}
    if err != nil {
        log.Printf("Server failed to create URL '%s': %v", longUrl, err)
    }

    return pbShortUrl, err
}

func (srv *URLServer) GetAssociated(ctx context.Context, pbShortUrl *pb.ShortURL) (*pb.LongURL, error) {
    shortUrl := database.ShortURL(pbShortUrl.Surl)
    longUrl, err := srv.db.GetUrl(shortUrl)

    pbLongUrl := &pb.LongURL{Lurl: string(longUrl)}
    if err != nil {
        log.Printf("Server failed to get short URL '%s': %v", shortUrl, err)
    }

    return pbLongUrl, err
}
