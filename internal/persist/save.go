package persist

import (
	"encoding/json"
	"os"

	"github.com/popolque/firstbitengi/internal/model"
)

func Save(state model.GameState, path string) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func Load(path string) (*model.GameState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var state model.GameState
	err = json.Unmarshal(data, &state)
	if err != nil {
		return nil, err
	}
	state.Sanitize()
	return &state, nil
}
