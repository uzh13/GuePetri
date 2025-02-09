package aggregate_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uzh13/GuePetri/pkg/petri/aggregate"
	"github.com/uzh13/GuePetri/pkg/petri/primitives/priority"
)

type StorageMock struct {
	q *priority.Queue[string, string]
}

func (s *StorageMock) Get(_ string) (*priority.Queue[string, string], error) {
	return s.q, nil
}

func TestBuilder(t *testing.T) {
	tests := []struct {
		name string
		q    *priority.Queue[string, string]
	}{
		{
			name: "Queue is nil",
			q:    nil,
		},
		{
			name: "Queue is not nil",
			q:    priority.NewPriorityQueue[string, string](),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &StorageMock{q: tt.q}
			target := aggregate.NewBuilder[string, string, string]("user1", storage)

			err := target.LoadState()
			if err != nil {
				t.Fatalf("failed to load state: %v", err)
			}

			result := target.Build("0")
			assert.NotNil(t, result.GetQueue())
		})
	}
}
