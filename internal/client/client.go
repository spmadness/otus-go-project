package client

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spmadness/otus-go-project/internal/server/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OutputMode int

const (
	OutputRaw OutputMode = iota
	OutputTable
)

type Client struct {
	host string
	port int64
	mode OutputMode

	conn *grpc.ClientConn
}

func NewClient(host string, port int64, mode OutputMode) *Client {
	return &Client{
		host: host,
		port: port,
		mode: mode,
	}
}

func (c *Client) Connect() error {
	var err error

	log.Println("connecting to server...")

	c.conn, err = grpc.Dial(
		fmt.Sprintf("%s:%d", c.host, c.port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC server: %w", err)
	}
	return nil
}

func (c *Client) Disconnect() error {
	err := c.conn.Close()
	log.Println("closing connection...")
	if err != nil {
		return err
	}

	return err
}

func (c *Client) Start(m int64, n int64, statType int64) error {
	var err error
	client := pb.NewMonitoringServiceClient(c.conn)

	req := &pb.Request{
		N:    n,
		M:    m,
		Type: pb.MetricType(statType),
	}

	stream, err := client.GetMetrics(context.Background(), req)
	if err != nil {
		return fmt.Errorf("grpc call error: %w", err)
	}

	err = c.readData(stream)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) readData(stream pb.MonitoringService_GetMetricsClient) error {
	for {
		response, err := stream.Recv()
		if err != nil {
			return fmt.Errorf("data read error: %w", err)
		}

		if c.mode == OutputTable {
			data := strings.Split(response.Data, "\n")
			c.drawTable(data)
		}
		if c.mode == OutputRaw {
			log.Println(response.Data)
		}
	}
}

func (c *Client) drawTable(data []string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	for k, line := range data {
		if k == 0 {
			t.SetTitle(line)
			continue
		}
		if k == 1 {
			header := strings.Fields(line)
			r := table.Row{}
			for _, val := range header {
				r = append(r, val)
			}
			t.AppendHeader(r)
			continue
		}
		row := strings.Fields(line)
		r := table.Row{}
		for _, val := range row {
			r = append(r, val)
		}
		t.AppendRow(r)
	}
	t.Render()
}
