package main

import (
    "os"
    "log"
    "fmt"
    "net"
    "google.golang.org/grpc"

    pb "github.com/AppCrashExpress/go-shorter/src/api"
    database "github.com/AppCrashExpress/go-shorter/src/database"
    shorter "github.com/AppCrashExpress/go-shorter/src"
)

func main() {
    envs := map[string]string {
        "POSTGRES_HOST": "",
        "POSTGRES_USER": "",
        "POSTGRES_PASSWORD": "",
        "POSTGRES_DB": "",
        "PORT": "",
    }
    inmemory := os.Getenv("INMEMORY")

    for key := range envs {
        envs[key] = os.Getenv(key)
        if envs[key] == "" {
            log.Fatalf("Environment variable %s is undefined\n", key)
        }
    }

    lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", envs["PORT"]))
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    var db database.Database
    if inmemory == "" {
        dbUrl := fmt.Sprintf("postgresql://%s:%s@%s:5432/%s",
            envs["POSTGRES_USER"], envs["POSTGRES_PASSWORD"], 
            envs["POSTGRES_HOST"], envs["POSTGRES_DB"])

        pg, err := database.NewPgDatabase(dbUrl, 50, 50)
        if err != nil {
            log.Fatalf("PostgreSQL failed to connect: %s\n", err)
        }

        defer func() {
            err := pg.Close()
            if err != nil {
                log.Fatalf("Database failed to close: %s\n", err)
            }
        }()

        db = pg

    } else {
        inm := database.NewMemoryDatabase()
        db = inm
    }
    
    log.Println("Starting server")
    grpcServer := grpc.NewServer()
    pb.RegisterShortnerServer(grpcServer, shorter.NewServer(db))
    grpcServer.Serve(lis)
}
