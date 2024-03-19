package server

import (
	"fmt"
	"time"

	"github.com/spmadness/otus-go-project/internal/scraper"
	"github.com/spmadness/otus-go-project/internal/server/pb"
)

func (s *Server) GetMetrics(request *pb.Request, server pb.MonitoringService_GetMetricsServer) error {
	fmt.Println("new grpc connection")

	scr, err := s.app.Scraper(scraper.MetricType(request.GetType()))
	if err != nil {
		return err
	}

	scr.AddConnection()
	defer scr.RemoveConnection()

	ctx := server.Context()

	durationWait := time.Second * time.Duration(request.GetM())
	durationTicker := time.Second * time.Duration(request.GetN())

	time.Sleep(durationWait)

	ticker := time.NewTicker(durationTicker)

	err = sendSnapshot(request.GetM(), scr, server)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			fmt.Println("closing grpc connection...")
			ticker.Stop()
			return nil

		case <-ticker.C:
			err = sendSnapshot(request.GetM(), scr, server)
			if err != nil {
				return err
			}
		}
	}
}

func sendSnapshot(seconds int64, scraper scraper.Scraper, server pb.MonitoringService_GetMetricsServer) error {
	fmt.Println("sending data...")
	data, err := scraper.GetSnapshot(seconds)
	if err != nil {
		return err
	}

	for _, d := range data {
		result := &pb.Result{
			Data: d,
		}
		err = server.Send(result)
		if err != nil {
			return err
		}
	}

	return nil
}
