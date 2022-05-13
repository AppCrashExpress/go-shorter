package shorter

import (
    "context"
    "log"
    "net"
    "os"
    "testing"

    "google.golang.org/grpc"
    "google.golang.org/grpc/test/bufconn"

    pb "github.com/AppCrashExpress/go-shorter/src/api"
    "github.com/AppCrashExpress/go-shorter/src/database"
)

const bufSize = 1024 * 1024

var (
    lis *bufconn.Listener
    ctx context.Context
    client pb.ShortnerClient
)

func bufDialer(context.Context, string) (net.Conn, error) {
    return lis.Dial()
}

func TestMain(m *testing.M) {
    lis = bufconn.Listen(bufSize)

    db := database.NewMemoryDatabase()
    grpcServer := grpc.NewServer()
    pb.RegisterShortnerServer(grpcServer, NewServer(db))
    go func() {
        if err := grpcServer.Serve(lis); err != nil {
            log.Fatalf("Server exited with error: %v", err)
        }
    }()
    
    ctx = context.Background()
    conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to dial bufnet: %v", err)
    }
    defer conn.Close()
    client = pb.NewShortnerClient(conn)

    os.Exit(m.Run())
}

func TestCreateAndGet(t *testing.T) {
    url := "https://google.com"

    pbShortUrl, err := client.CreateNew(ctx, &pb.LongURL{Lurl: url})
    if err != nil {
        t.Fatalf("client.CreateNew failed: %v", err)
    }

    shortUrl := pbShortUrl.Surl
    pbLongUrl, err := client.GetAssociated(ctx, &pb.ShortURL{Surl: shortUrl})
    if err != nil {
        t.Fatalf("client.GetAssociated failed: %v", err)
    }

    longUrl := pbLongUrl.Lurl
    if longUrl != url {
        t.Fatalf("Returned URL does not match initial (expected '%s', got '%s')", url, longUrl)
    }
}

func TestCreateTwice(t *testing.T) {
    url := "https://go.dev"

    _, err := client.CreateNew(ctx, &pb.LongURL{Lurl: url})
    if err != nil {
        t.Fatalf("client.CreateNew failed when shouldn't have: %v", err)
    }
    _, err = client.CreateNew(ctx, &pb.LongURL{Lurl: url})
    if err == nil {
        t.Fatalf("client.CreateNew didn't fail properly (nil error)")
    }
}

func TestGetNonexisting(t *testing.T) {
    url := "https://go.dev/doc"

    _, err := client.GetAssociated(ctx, &pb.ShortURL{Surl: url})
    if err == nil {
        t.Fatalf("client.GetAssociated didn't fail properly (nil error)")
    }
}
