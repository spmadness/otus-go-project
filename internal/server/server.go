package server

import (
	"fmt"
	"log"
	"net"

	"github.com/spmadness/otus-go-project/internal/scraper"
	"github.com/spmadness/otus-go-project/internal/server/pb"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedMonitoringServiceServer

	app    Application
	server *grpc.Server
	port   int
}

type Application interface {
	Scrapers() map[scraper.MetricType]scraper.Scraper
	Scraper(code scraper.MetricType) (scraper.Scraper, error)
}

func NewServer(app Application, port int) *Server {
	return &Server{
		app:  app,
		port: port,
	}
}

func (s *Server) Start() {
	lsn, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Println(err.Error())
	}

	s.server = grpc.NewServer()
	pb.RegisterMonitoringServiceServer(s.server, s)

	log.Printf("starting grpc server on %s \n", lsn.Addr().String())

	if err = s.server.Serve(lsn); err != nil {
		log.Println(err.Error())
	}
}

func (s *Server) Stop() {
	log.Println("stopping grpc server...")
	s.server.Stop()
}
