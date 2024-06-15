package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/ic-it/gpadkm"
	"github.com/ynsta/uinput"
)

const (
	LISTEN_TICK_INTERVAL = 10 * time.Millisecond
	MAX_SPEED            = 5.
)

func main() {
	g, err := gpadkm.NewGamepad(0)
	if err != nil {
		log.Fatalf("error creating gamepad: %v", err)
	}
	defer g.Close()

	me, err := gpadkm.NewMouseEmulator()
	if err != nil {
		log.Fatalf("error creating mouse emulator: %v", err)
	}
	defer me.Close()

	ke, err := gpadkm.NewKeyboardEmulator([]uinput.EventCode{
		uinput.KEY_LEFTMETA,
	})
	if err != nil {
		log.Fatalf("error creating keyboard emulator: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	log.Println("Press CTRL+C to exit")

	maxSpeed := MAX_SPEED
	moveX, moveY := 0., 0.
	scrollX, scrollY := 0, 0
	super := false
	leftMouse := false
	rightMouse := false
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if super {
					if err := ke.KeyDown(uinput.KEY_LEFTMETA); err != nil {
						log.Printf("error pressing key: %v", err)
					}
				} else {
					if err := ke.KeyUp(uinput.KEY_LEFTMETA); err != nil {
						log.Printf("error releasing key: %v", err)
					}
				}
				if leftMouse && rightMouse {
					if err := me.ButtonDown(uinput.BTN_MIDDLE); err != nil {
						log.Printf("error pressing button: %v", err)
					}
					goto skipLR
				} else {
					if err := me.ButtonUp(uinput.BTN_MIDDLE); err != nil {
						log.Printf("error releasing button: %v", err)
					}
				}

				if leftMouse {
					if err := me.ButtonDown(uinput.BTN_LEFT); err != nil {
						log.Printf("error pressing button: %v", err)
					}
				} else {
					if err := me.ButtonUp(uinput.BTN_LEFT); err != nil {
						log.Printf("error releasing button: %v", err)
					}
				}
				if rightMouse {
					if err := me.ButtonDown(uinput.BTN_RIGHT); err != nil {
						log.Printf("error pressing button: %v", err)
					}
				} else {
					if err := me.ButtonUp(uinput.BTN_RIGHT); err != nil {
						log.Printf("error releasing button: %v", err)
					}
				}

			skipLR:

				moveX := int(moveX * maxSpeed)
				moveY := int(moveY * maxSpeed)
				if err := me.Move(moveX, moveY); err != nil {
					log.Printf("error moving mouse: %v", err)
				}
				if err := me.Scroll(scrollX, -scrollY); err != nil {
					log.Printf("error scrolling mouse: %v", err)
				}
				time.Sleep(10 * time.Millisecond)
			}
		}
	}(ctx)

	for s := range g.Listen(ctx, LISTEN_TICK_INTERVAL) {
		log.Printf("state: %+v", s)
		if scrollX = s.AxisData[3]; scrollX > 0 {
			scrollX = 1
		} else if scrollX < 0 {
			scrollX = -1
		} else {
			scrollX = 0
		}

		if scrollY = s.AxisData[4]; scrollY > 0 {
			scrollY = 1
		} else if scrollY < 0 {
			scrollY = -1
		} else {
			scrollY = 0
		}

		if s.Buttons&16 == 16 {
			leftMouse = true
		} else {
			leftMouse = false
		}
		if s.Buttons&32 == 32 {
			rightMouse = true
		} else {
			rightMouse = false
		}
		if s.Buttons&256 == 256 {
			super = true
		} else {
			super = false
		}

		moveX, moveY = float64(s.AxisData[0])/(1<<15), float64(s.AxisData[1])/(1<<15)
		speedup := float64(s.AxisData[2]) / (1 << 15)
		if speedup > 0 {
			maxSpeed = MAX_SPEED + MAX_SPEED*speedup
		} else {
			maxSpeed = MAX_SPEED
		}
	}

	<-ctx.Done()
	log.Println("Exiting")
}
