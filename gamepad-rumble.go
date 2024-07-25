package gpadkm

/*
#include <fcntl.h>
#include <linux/input.h>
#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <errno.h>

// Constants for the force feedback effect
#define FF_RUMBLE 0x50
#define EV_FF 0x15

int _rumble(int device, int strong, int weak, int length) {
	struct ff_effect effect;
	memset(&effect, 0, sizeof(effect));

	effect.type = FF_RUMBLE;
	effect.id = -1; // -1 indicates we are creating a new effect
	effect.u.rumble.strong_magnitude = strong; // Strong rumble magnitude
	effect.u.rumble.weak_magnitude = weak;   // Weak rumble magnitude
	effect.replay.length = length;               // Length of the effect in ms
	effect.replay.delay = 0;                   // Delay before starting the effect

	if (ioctl(device, EVIOCSFF, &effect) < 0) {
		return errno;
	}

	struct input_event playEvent;
	memset(&playEvent, 0, sizeof(playEvent));

	playEvent.type = EV_FF;
	playEvent.code = effect.id;
	playEvent.value = 1; // Start the effect

	if (write(device, &playEvent, sizeof(playEvent)) < 0) {
		return errno;
	}

	// Let the effect play for its duration
	usleep(effect.replay.length * 1000);

	// Stop the effect
	playEvent.value = 0; // Stop the effect

	if (write(device, &playEvent, sizeof(playEvent)) < 0) {
		return errno;
	}

	return 0;
}
*/
import "C"
import (
	"fmt"
	"os"
)

func rumble(devicePath string, strong, weak, length int) error {
	fd, err := os.OpenFile(devicePath, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer fd.Close()

	if errno := C._rumble(C.int(fd.Fd()), C.int(strong), C.int(weak), C.int(length)); errno != 0 {
		return fmt.Errorf("error: %d", errno)
	}

	return nil
}
