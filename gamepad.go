package gpadkm

import (
	"context"
	"sync"
	"time"

	"github.com/0xcafed00d/joystick"
)

// Gamepad wraps the joystick package
type Gamepad struct {
	joystick joystick.Joystick
	sync.Mutex
}

// NewGamepad creates a new Gamepad
func NewGamepad(jsid int) (*Gamepad, error) {
	js, err := joystick.Open(jsid)
	if err != nil {
		return nil, err
	}

	return &Gamepad{joystick: js}, nil
}

// Close closes the Gamepad
func (g *Gamepad) Close() {
	g.Lock()
	defer g.Unlock()

	g.joystick.Close()
}

// Listen returns a channel of joystick unique states
func (g *Gamepad) Listen(ctx context.Context, tickInterval time.Duration) <-chan joystick.State {
	g.Lock()
	defer g.Unlock()

	state := make(chan joystick.State)
	go func() {
		defer close(state)
		axCount := g.joystick.AxisCount()
		axData := make([]int, axCount)
		buttons := uint32(0)
	mainloop:
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(tickInterval):
			}
			s, err := g.joystick.Read()
			if err != nil {
				return
			}
			for i := 0; i < axCount; i++ {
				if s.AxisData[i] != axData[i] {
					goto update
				}
			}
			if s.Buttons != buttons {
				goto update
			}
			continue mainloop

		update:
			for i := 0; i < axCount; i++ {
				axData[i] = s.AxisData[i]
			}
			buttons = s.Buttons
			select {
			case state <- s:
			case <-ctx.Done():
				return
			}
		}
	}()
	return state
}
