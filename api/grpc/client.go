package grpc

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/BYTE-Club-CCNY/Silvercord/api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// LLMClient wraps a gRPC connection to the Python LLM service
type LLMClient struct {
	conn   *grpc.ClientConn
	client pb.LLMServiceClient
}

// NewLLMClient creates a new gRPC client connection to the Python LLM service
func NewLLMClient(address string) (*LLMClient, error) {
	// Set up connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to LLM service at %s: %w", address, err)
	}

	client := pb.NewLLMServiceClient(conn)
	log.Printf("Connected to LLM service at %s", address)

	return &LLMClient{
		conn:   conn,
		client: client,
	}, nil
}

// ProcessQuery sends a query to the LLM service and returns the response
func (c *LLMClient) ProcessQuery(ctx context.Context, userID, command, query string) (string, error) {
	// Set a reasonable timeout for LLM processing
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	request := &pb.QueryRequest{
		UserId:  userID,
		Query:   query,
		Command: command,
	}

	response, err := c.client.ProcessRequest(ctx, request)
	if err != nil {
		return "", fmt.Errorf("failed to process query: %w", err)
	}

	if !response.Success {
		return "", fmt.Errorf("LLM service returned unsuccessful response")
	}

	return response.ResponseText, nil
}

// Close closes the gRPC connection
func (c *LLMClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
