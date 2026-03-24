package massalib

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/KarpelesLab/massalib/massagrpc"
	"google.golang.org/grpc"
)

type RPC struct {
	cPub *grpc.ClientConn
	pub  massagrpc.PublicServiceClient
}

// New returns a new connection to a massa node
//
// grpc://localhost:33037
func New(target string, opts ...grpc.DialOption) (*RPC, error) {
	cPub, err := grpc.NewClient(target, opts...)
	if err != nil {
		return nil, err
	}
	cl := &RPC{
		cPub: cPub,
		pub:  massagrpc.NewPublicServiceClient(cPub),
	}

	return cl, nil
}

// Public exposes the raw grpc public interface
func (rpc *RPC) Public() massagrpc.PublicServiceClient {
	return rpc.pub
}

// GetStatus returns the status of the node
func (rpc *RPC) GetStatus(ctx context.Context) (*massagrpc.PublicStatus, error) {
	res, err := rpc.pub.GetStatus(ctx, &massagrpc.GetStatusRequest{})
	if err != nil {
		return nil, err
	}
	return res.Status, nil
}

func (rpc *RPC) GetSlotTransfers(ctx context.Context, finality massagrpc.FinalityLevel) (chan *massagrpc.NewSlotTransfersResponse, io.Closer, error) {
	// note: this requires massa to be compiled with feature massa-node/execution-trace
	bidi, err := rpc.pub.NewSlotTransfers(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("new slottransfers failed: %w", err)
	}
	req := &massagrpc.NewSlotTransfersRequest{
		FinalityLevel: finality,
	}
	if err := bidi.Send(req); err != nil {
		return nil, nil, fmt.Errorf("failed to set finality level: %w", err)
	}

	// when CloseSend is called, massa closes the pipe
	// https://github.com/massalabs/massa/blob/main/massa-grpc/src/stream/new_slot_transfers.rs#L142

	ch := make(chan *massagrpc.NewSlotTransfersResponse)
	go func() {
		defer close(ch)
		for {
			resp, err := bidi.Recv()
			if err == io.EOF {
				log.Printf("got EOF")
				return
			}
			if err != nil {
				log.Printf("MASSA error on recv: %s", err)
				return
			}
			ch <- resp
		}
	}()

	return ch, bidiCloser{bidi}, nil
}

func (rpc *RPC) SendOperations(ctx context.Context, op ...[]byte) ([]string, error) {
	bidi, err := rpc.pub.SendOperations(ctx)
	if err != nil {
		return nil, err
	}
	defer bidi.CloseSend()

	if err := bidi.Send(&massagrpc.SendOperationsRequest{Operations: op}); err != nil {
		return nil, err
	}
	// wait for response, just in case
	ids, err := bidi.Recv()
	if err != nil {
		return nil, err
	}
	return ids.GetOperationIds().OperationIds, nil
}

type bidiCloser struct {
	grpc.ClientStream
}

func (b bidiCloser) Close() error {
	return b.ClientStream.CloseSend()
}
