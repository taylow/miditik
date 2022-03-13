package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"gitlab.com/gomidi/midi"
)

// decodeMidi decodes a midi.Message into a Midi struct
func decodeMidi(msg midi.Message) (*Midi, error) {
	parts := strings.Split(msg.String(), " ")
	fmt.Printf("Parts: %#v\n", parts)

	var channel string
	var key string
	var value string
	action := parts[0]
	midi := &Midi{}

	switch action {
	case "channel.NoteOn":
		midi.action = NoteOn
		channel = parts[2]
		key = parts[4]
	case "channel.NoteOff":
		midi.action = NoteOff
		channel = parts[2]
		key = parts[4]
	case "channel.Pitchbend":
		midi.action = Pitchbend
		value = parts[6]
	default:
		return nil, fmt.Errorf("invalid action: %s", action)
	}

	if channel != "" {
		pChan, err := strconv.ParseInt(channel, 0, 8)
		if err != nil {
			return nil, fmt.Errorf("could not parse channel int: %v", err)
		}
		midi.channel = int(pChan)
	}

	if key != "" {
		pKey, err := strconv.ParseInt(key, 0, 8)
		if err != nil {
			return nil, fmt.Errorf("could not parse key int: %v", err)
		}
		midi.key = int(pKey)
	}

	if value != "" {
		pValue, err := strconv.ParseInt(value, 0, 32)
		if err != nil {
			return nil, fmt.Errorf("could not parse value int: %v", err)
		}
		midi.value = int(pValue)
	}

	return midi, nil
}

// calculateMidiFrequencies calculates the frequencies for each midi value from min to max
func calculateMidiFrequencies(min, max int) map[int]float64 {
	midiMap := make(map[int]float64)
	for i := min; i < max; i++ {
		midiMap[i] = math.Pow(2, (float64(i)-49)/12.0) * 440
	}
	return midiMap
}

// randomBeeps iterates a map of midi keys and plays the keys in the randomised map order through the Beeper
func randomBeeps(beeper Beeper, midiMap map[int]float64) error {
	for k := range midiMap {
		err := beeper.Beep(k, 0)
		if err != nil {
			return err
		}
	}
	beeper.NoBeep()
	return nil
}

// must checks if err != nil and panics if not nil
func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}

// remove removes a value from a slice of ints and returns its index
func remove(slice []int, remove int) ([]int, int) {
	for i, v := range slice {
		if v == remove {
			return append(slice[:i], slice[i+1:]...), i
		}
	}
	return slice, 0
}
