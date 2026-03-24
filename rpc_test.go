package massalib

import (
	"context"
	"log"
	"testing"

	"github.com/KarpelesLab/massalib/massagrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestRPC(t *testing.T) {
	c, err := New("localhost:33037", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Skipf("failed to establish RPC: %s", err)
		return
	}

	st, err := c.GetStatus(context.Background())
	if err != nil {
		// can be a failed to connect to rpc
		t.Skipf("failed to get massa status: %s", err)
		return
	}
	log.Printf("connected to massa nodeid = %s version = %s", st.NodeId, st.Version)

	ch, clo, err := c.GetSlotTransfers(context.Background(), massagrpc.FinalityLevel_FINALITY_LEVEL_FINAL)
	if err != nil {
		t.Errorf("failed to request transfers: %s", err)
	}
	defer clo.Close()

	cnt := 0

	for tx := range ch {
		log.Printf("tx = %+v %+v", tx.Slot, tx.Transfers)
		cnt += 1
		if cnt > 4 {
			break
		}
	}
	log.Printf("end")
}
