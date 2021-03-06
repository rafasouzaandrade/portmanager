package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "portservice/proto"
	"portservice/server"
	"runtime"
	"time"
)

const (
	port = ":50051"
)

func bToKb(b uint64) string {
	return fmt.Sprintf("%vKb", b / 1024)
}

func PrintMemUsage() {
	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <- ticker.C:
				var m runtime.MemStats
				runtime.ReadMemStats(&m)
				log.Printf("Alloc = %v, HeapAlloc = %v, TotalAlloc = %v, Sys = %v, NumGC = %v", bToKb(m.Alloc), bToKb(m.HeapAlloc	), bToKb(m.TotalAlloc), bToKb(m.Sys), m.NumGC)
			case <- quit:
				ticker.Stop()
				return
			}
		}
	}()

}

func main() {
	PrintMemUsage()
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPortDomainServer(s, &server.PortServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
