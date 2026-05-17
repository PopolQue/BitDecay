package persist

import (
	"time"

	"github.com/popolque/firstbitengi/internal/model"
)

type Autosaver struct {
	snapshotCh chan model.GameState
	path       string
	interval   time.Duration
	lastSave   time.Time
}

func NewAutosaver(path string, interval time.Duration) *Autosaver {
	return &Autosaver{
		snapshotCh: make(chan model.GameState, 1),
		path:       path,
		interval:   interval,
		lastSave:   time.Now(),
	}
}

func (a *Autosaver) Start() {
	go func() {
		for snap := range a.snapshotCh {
			_ = Save(snap, a.path)
		}
	}()
}

func (a *Autosaver) MaybeSnapshot(state *model.GameState) {
	if time.Since(a.lastSave) >= a.interval {
		a.lastSave = time.Now()
		
		// Deep clone state
		snap := *state
		snap.Hardware = make(map[string]int)
		for k, v := range state.Hardware {
			snap.Hardware[k] = v
		}
		snap.Upgrades = make(map[string]bool)
		for k, v := range state.Upgrades {
			snap.Upgrades[k] = v
		}

		select {
		case a.snapshotCh <- snap:
		default:
			// Previous save still in progress, skip this one
		}
	}
}
