package main

import (
	"context"
	"fmt"
	"log"
	"math"
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
	evDev, err := gpadkm.GetAvailableEventDevices()
	if err != nil {
		log.Fatalf("error getting available event devices: %v", err)
	}
	if len(evDev) == 0 {
		log.Fatalf("no event devices found")
	}
	fmt.Printf("Select an event device:\n")
	for i, d := range evDev {
		fmt.Printf("  %d: %s\n", i, d.Name)
	}

	var evDevIdx int
	fmt.Printf("Enter the event device index [0-%d]: ", len(evDev)-1)
	fmt.Scanf("%d", &evDevIdx)
	if evDevIdx < 0 || evDevIdx >= len(evDev) {
		log.Fatalf("invalid event device index")
	}

	g, err := gpadkm.NewGamepad(0, evDev[evDevIdx].Path)
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
		uinput.KEY_SPACE,
		uinput.KEY_ESC,
		uinput.KEY_LEFTCTRL,
		uinput.KEY_ENTER,

		uinput.KEY_UP,
		uinput.KEY_DOWN,
		uinput.KEY_LEFT,
		uinput.KEY_RIGHT,

		uinput.KEY_V,
		uinput.KEY_C,
	})
	if err != nil {
		log.Fatalf("error creating keyboard emulator: %v", err)
	}

	println(
		"Gamepad Keyboard Mouse\n" +
			"  Left Stick: Move Mouse\n" +
			"  Right Stick: Scroll\n" +
			"  A: Space\n" +
			"  B: \"V\"\n" +
			"  X: \"C\"\n" +
			"  Y: Enter\n" +
			"  Left Bumper: Left Mouse Button\n" +
			"  Right Bumper: Right Mouse Button\n" +
			"  Left+Right Bumper: Middle Mouse Button\n" +
			"  Start: Escape\n" +
			"  Select: CTRL\n" +
			"  Button Logo: Super\n" +
			"  Right Trigger: Speed Up\n" +
			"  Pad Up/Down/Left/Right: Arrow Keys\n" +
			"  Left Stick Pressed + Right Stick Pressed: Toggle Enabled\n" +
			"  Press CTRL+C to exit\n",
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	log.Println("Press CTRL+C to exit")

	maxSpeed := MAX_SPEED
	moveX, moveY := 0., 0.
	scrollX, scrollY := 0, 0
	super := false
	escape := false
	lkm, rkm := false, false
	kup, kdown, kleft, kright := false, false, false, false
	lstickPressed, rstickPressed := false, false

	kspace := false
	enter := false
	lastScrollX, lastScrollY := time.Now(), time.Now()
	keyV, keyC := false, false
	ctrl := false

	rumbled := false
	enabled := true
	go func(ctx context.Context) {
		for {
			keymap := map[uinput.EventCode]bool{
				uinput.KEY_UP:    kup,
				uinput.KEY_DOWN:  kdown,
				uinput.KEY_LEFT:  kleft,
				uinput.KEY_RIGHT: kright,

				uinput.KEY_LEFTMETA: super,
				uinput.KEY_SPACE:    kspace,
				uinput.KEY_ESC:      escape,
				uinput.KEY_ENTER:    enter,

				uinput.KEY_V:        keyV,
				uinput.KEY_C:        keyC,
				uinput.KEY_LEFTCTRL: ctrl,
			}
			if lstickPressed && rstickPressed {
				log.Println("Toggle enabled")
				enabled = !enabled
				time.Sleep(500 * time.Millisecond)
			}
			if !enabled {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			select {
			case <-ctx.Done():
				log.Println("Exiting")
				return
			default:
				if lkm && rkm {
					if err := me.ButtonDown(uinput.BTN_MIDDLE); err != nil {
						log.Printf("error pressing button: %v", err)
					}
					goto skipLR
				} else {
					if err := me.ButtonUp(uinput.BTN_MIDDLE); err != nil {
						log.Printf("error releasing button: %v", err)
					}
				}

				if lkm {
					if err := me.ButtonDown(uinput.BTN_LEFT); err != nil {
						log.Printf("error pressing button: %v", err)
					}
				} else {
					if err := me.ButtonUp(uinput.BTN_LEFT); err != nil {
						log.Printf("error releasing button: %v", err)
					}
				}
				if rkm {
					if err := me.ButtonDown(uinput.BTN_RIGHT); err != nil {
						log.Printf("error pressing button: %v", err)
					}
				} else {
					if err := me.ButtonUp(uinput.BTN_RIGHT); err != nil {
						log.Printf("error releasing button: %v", err)
					}
				}

			skipLR:
				for k, v := range keymap {
					if v {
						if err := ke.KeyDown(k); err != nil {
							log.Printf("error pressing key: %v", err)
						}
					} else {
						if err := ke.KeyUp(k); err != nil {
							log.Printf("error releasing key: %v", err)
						}
					}

					if k == uinput.KEY_LEFTMETA && v {
						if !rumbled {
							err := g.Rumble(0x8000, 0x8000, 100)
							if err != nil {
								log.Printf("error rumbling: %v", err)
							}
							rumbled = true
						}
					} else if k == uinput.KEY_LEFTMETA && !v {
						rumbled = false
					}
				}

				moveX := int(moveX * maxSpeed)
				moveY := int(moveY * maxSpeed)
				if err := me.Move(moveX, moveY); err != nil {
					log.Printf("error moving mouse: %v", err)
				}
				if scrollX != 0 {
					if time.Since(lastScrollX) > time.Second/time.Duration(math.Abs(float64(scrollX))) {
						if err := me.Scroll(scrollX, 0); err != nil {
							log.Printf("error scrolling mouse: %v", err)
						}
						lastScrollX = time.Now()
					}
				}
				if scrollY != 0 {
					if time.Since(lastScrollY) > time.Second/time.Duration(math.Abs(float64(scrollY))) {
						if err := me.Scroll(0, -scrollY); err != nil {
							log.Printf("error scrolling mouse: %v", err)
						}
						lastScrollY = time.Now()
					}
				}
				time.Sleep(10 * time.Millisecond)
			}
		}
	}(ctx)

	for s := range g.Listen(ctx, LISTEN_TICK_INTERVAL) {
		log.Printf("state: %+v", s)
		scrollX = s.AxisData[3] / 5000
		scrollY = s.AxisData[4] / 5000

		kspace = s.Buttons&1 == 1
		lkm = s.Buttons&16 == 16
		rkm = s.Buttons&32 == 32
		super = s.Buttons&256 == 256
		escape = s.Buttons&128 == 128
		lstickPressed = s.Buttons&512 == 512
		rstickPressed = s.Buttons&1024 == 1024
		keyV = s.Buttons&2 == 2
		keyC = s.Buttons&4 == 4
		enter = s.Buttons&8 == 8
		ctrl = s.Buttons&64 == 64

		moveX, moveY = float64(s.AxisData[0])/(1<<15), float64(s.AxisData[1])/(1<<15)
		speedup := float64(s.AxisData[5]) / (1 << 15)
		if speedup > 0 {
			maxSpeed = MAX_SPEED + MAX_SPEED*speedup
		} else {
			maxSpeed = MAX_SPEED
		}

		kleft, kright = s.AxisData[6] < 0, s.AxisData[6] > 0
		kup, kdown = s.AxisData[7] < 0, s.AxisData[7] > 0
	}
	cancel()
	<-ctx.Done()
	log.Println("Exiting")
}

func absInt(i int) int {
	if i < 0 {
		return -i
	}
	return i
}
