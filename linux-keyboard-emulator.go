package gpadkm

import (
	"context"
	"sync"
	"time"

	"github.com/ynsta/uinput"
)

// KeyboardEmulator is a struct that holds the uinput device and the context
type KeyboardEmulator struct {
	wd uinput.WriteDevice
	ui uinput.UInput
	sync.Mutex
}

// NewKeyboardEmulator creates a new KeyboardEmulator
func NewKeyboardEmulator(
	keys []uinput.EventCode,
) (*KeyboardEmulator, error) {
	var ke KeyboardEmulator
	ke.wd.Open()
	if err := ke.ui.Init(
		&ke.wd,
		"gpadkm-keyboard-emulator",
		0x1234, // vendor
		0x5678, // product
		0x1,    // version
		keys,
		[]uinput.EventCode{},
		[]uinput.AxisSetup{},
		true,
	); err != nil {
		return nil, err
	}
	return &ke, nil
}

// Close closes the uinput device
func (ke *KeyboardEmulator) Close() {
	ke.Lock()
	defer ke.Unlock()

	ke.wd.Close()
}

// PressKey presses a key for a duration
func (ke *KeyboardEmulator) PressKey(ctx context.Context, key uinput.EventCode, duration time.Duration) error {
	ke.Lock()
	defer ke.Unlock()

	if err := ke.ui.KeyEvent(key, 1); err != nil {
		return err
	}
	if err := ke.ui.SynEvent(); err != nil {
		return err
	}
	select {
	case <-time.After(duration):
	case <-ctx.Done():
	}
	if err := ke.ui.KeyEvent(key, 0); err != nil {
		return err
	}
	if err := ke.ui.SynEvent(); err != nil {
		return err
	}
	return nil
}

// KeyDown presses a key
func (ke *KeyboardEmulator) KeyDown(key uinput.EventCode) error {
	ke.Lock()
	defer ke.Unlock()

	if err := ke.ui.KeyEvent(key, 1); err != nil {
		return err
	}
	if err := ke.ui.SynEvent(); err != nil {
		return err
	}
	return nil
}

// KeyUp releases a key
func (ke *KeyboardEmulator) KeyUp(key uinput.EventCode) error {
	ke.Lock()
	defer ke.Unlock()

	if err := ke.ui.KeyEvent(key, 0); err != nil {
		return err
	}
	if err := ke.ui.SynEvent(); err != nil {
		return err
	}
	return nil
}
