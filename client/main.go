package main

import (
    "os"
    "context"
    "fmt"
    "log"
    "time"

    "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

    pb "github.com/AppCrashExpress/go-shorter/src/api"
)

func createNew(client pb.ShortnerClient, longUrl string) string {
    pbLongUrl := &pb.LongURL{Lurl: longUrl}
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    pbShortUrl, err := client.CreateNew(ctx, pbLongUrl)
    if err != nil {
        log.Fatalf("client.CreateNew failed: %v", err)
    }

    return pbShortUrl.Surl
}

func GetAssociated(client pb.ShortnerClient, shortUrl string) string {
    pbShortUrl := &pb.ShortURL{Surl: shortUrl}
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    pbLongUrl, err := client.GetAssociated(ctx, pbShortUrl)
    if err != nil {
        log.Fatalf("client.GetAssociated failed: %v", err)
    }

    return pbLongUrl.Lurl
}

func main() {
    envs := map[string]string {
        "PORT": "",
    }

    for key := range envs {
        envs[key] = os.Getenv(key)
        if envs[key] == "" {
            log.Fatalf("Environment variable %s is undefined\n", key)
        }
    }

    serverUrl := fmt.Sprintf("localhost:%s", envs["PORT"])
    conn, err := grpc.Dial(serverUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("Did not connect: %v", err)
    }
    defer conn.Close()
    c := pb.NewShortnerClient(conn)

    url := "https://google.com"

    shortUrl := createNew(c, url)
    longUrl  := GetAssociated(c, shortUrl)

    fmt.Printf("%s == %s: %v\n", url, longUrl, url == longUrl)
}
