package handle

import (
	"math/rand"
	"webmote/transform"
)

var transforms = make(map[string]*transform.Transform)

type Event struct {
	Event string `json:"event"`
	Data  struct {
		Alpha    float64 `json:"alpha"`
		Beta     float64 `json:"beta"`
		Gamma    float64 `json:"gamma"`
		Absolute bool    `json:"absolute"`
	} `json:"data"`
}

func Remove(id string) {
	if _, ok := transforms[id]; ok {
		delete(transforms, id)
	}
}

func WS(id string, e Event) (transform.Coordinate, bool) {
	var t *transform.Transform
	var ok bool
	if t, ok = transforms[id]; !ok {
		t = transform.New()
		transforms[id] = t
	}

	r := transform.Rotation{
		IsSet: true,
		Roll:  e.Data.Gamma,
		Pitch: e.Data.Beta,
		Yaw:   e.Data.Alpha,
	}

	switch e.Event {
	case "lt":
		t.LeftTopR = r
	case "rt":
		t.RightTopR = r
	case "lb":
		t.LeftBottomR = r
	case "rb":
		t.RightBottomR = r
	default:
		if !t.IsCalibrated {
			return transform.Coordinate{}, false
		}
		c := t.ScreenCoordinate(r)
		return c, true
	}

	if t.LeftTopR.IsSet && t.RightTopR.IsSet && t.LeftBottomR.IsSet && t.RightBottomR.IsSet {
		t.Calibrate()
	}

	return transform.Coordinate{}, false
}

func genRandomStr() string {
	var digits [6]byte
	for i := range digits {
		digits[i] = byte(rand.Intn(10) + 48) // 48은 '0'의 ASCII 코드
	}
	return string(digits[:])
}

func NewId() string {
	for {
		id := genRandomStr()
		if _, ok := transforms[id]; !ok {
			return id
		}
	}
}
