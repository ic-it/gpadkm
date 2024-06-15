//go:build linux
// +build linux

package gpadkm

import (
	"sync"

	"github.com/ynsta/uinput"
)

// MouseEmulator is a struct that holds the uinput device and the context
type MouseEmulator struct {
	wd uinput.WriteDevice
	ui uinput.UInput
	sync.Mutex
}

// NewMouseEmulator creates a new MouseEmulator
func NewMouseEmulator() (*MouseEmulator, error) {
	var me MouseEmulator
	me.wd.Open()
	if err := me.ui.Init(
		&me.wd,
		"gpadkm-mouse-emulator",
		0x1234, // vendor
		0x5678, // product
		0x1,    // version
		[]uinput.EventCode{
			uinput.BTN_LEFT,
			uinput.BTN_RIGHT,
			uinput.BTN_MIDDLE,
		},
		[]uinput.EventCode{
			uinput.ABS_X,
			uinput.ABS_Y,

			uinput.REL_X,
			uinput.REL_Y,
			uinput.REL_WHEEL,
			uinput.REL_HWHEEL,
		},
		[]uinput.AxisSetup{},
		false,
	); err != nil {
		return nil, err
	}
	return &me, nil
}

// Close closes the uinput device
func (me *MouseEmulator) Close() {
	me.Lock()
	defer me.Unlock()

	me.wd.Close()
}

// MoveTo moves the mouse to x and y
func (me *MouseEmulator) MoveTo(x, y int) error {
	me.Lock()
	defer me.Unlock()

	if err := me.ui.AbsEvent(uinput.ABS_X, uinput.EventValue(x)); err != nil {
		return err
	}
	if err := me.ui.AbsEvent(uinput.ABS_Y, uinput.EventValue(y)); err != nil {
		return err
	}
	if err := me.ui.SynEvent(); err != nil {
		return err
	}
	return nil
}

// Move moves the mouse by x and y
func (me *MouseEmulator) Move(x, y int) error {
	me.Lock()
	defer me.Unlock()

	if err := me.ui.RelEvent(uinput.REL_X, uinput.EventValue(x)); err != nil {
		return err
	}
	if err := me.ui.RelEvent(uinput.REL_Y, uinput.EventValue(y)); err != nil {
		return err
	}
	if err := me.ui.SynEvent(); err != nil {
		return err
	}
	return nil
}

// Scroll scrolls the mouse wheel
func (me *MouseEmulator) Scroll(x, y int) error {
	me.Lock()
	defer me.Unlock()

	if err := me.ui.RelEvent(uinput.REL_WHEEL, uinput.EventValue(y)); err != nil {
		return err
	}
	if err := me.ui.RelEvent(uinput.REL_HWHEEL, uinput.EventValue(x)); err != nil {
		return err
	}
	if err := me.ui.SynEvent(); err != nil {
		return err
	}
	return nil
}

// ButtonDown presses the button
func (me *MouseEmulator) ButtonDown(button uinput.EventCode) error {
	me.Lock()
	defer me.Unlock()

	if err := me.ui.KeyEvent(button, 1); err != nil {
		return err
	}
	if err := me.ui.SynEvent(); err != nil {
		return err
	}
	return nil
}

// ButtonUp releases the button
func (me *MouseEmulator) ButtonUp(button uinput.EventCode) error {
	me.Lock()
	defer me.Unlock()

	if err := me.ui.KeyEvent(button, 0); err != nil {
		return err
	}
	if err := me.ui.SynEvent(); err != nil {
		return err
	}
	return nil
}
