package query

import (
	"testing"
)

func TestMockQueryClient(t *testing.T) {
	mockQueryClient := MockQueryClient{}

	resp, err := mockQueryClient.GetMiners()

	if err != nil {
		t.Errorf("GetMiners(): %v", err)
	}

	if len(resp.Data) != 2 {
		t.Errorf("expected two miner data: %v", resp.Data)
	}
}
