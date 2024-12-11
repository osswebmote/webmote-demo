package transform

import "math"

// SCREEN RESOLUTION :: FULL HD
const SCREEN_X = 1920
const SCREEN_Y = 1080

type Rotation struct {
	IsSet bool
	Roll  float64 // 왼기울기 - 오른기울기 + ( -90~90 0은 화면위 방향 )
	Pitch float64 // 위쪽 + 아래쪽 - ( -180~180 0은 화면 위 방향 )
	Yaw   float64 // 왼쪽 + 오른쪽 - ( 0~360 = 원점 보정 필요 )
}

type Coordinate struct {
	X float64
	Y float64
	R int // -1 = left, 0 = center, 1 = right
}

type Transform struct {
	IsCalibrated bool
	oYaw         float64 // origin yaw
	eYaw         float64 // end yaw
	hasOrigin    bool    // yaw 의 원점 좌표가 yaw 범위에 포함중인지?
	oPitch       float64 // origin pitch
	ePitch       float64 // end yaw
	oRoll        float64 // origin roll
	yawFov       float64 // screen fov of yaw
	pitchFov     float64 // screen fov of pitch
	LeftTopR     Rotation
	RightTopR    Rotation
	LeftBottomR  Rotation
	RightBottomR Rotation
}

func New() *Transform {
	return &Transform{}
}

func (t *Transform) Calibrate() {
	t.oYaw = t.LeftTopR.Yaw
	if t.LeftBottomR.Yaw > t.oYaw {
		t.oYaw = t.LeftBottomR.Yaw
	}

	t.eYaw = t.RightTopR.Yaw
	if t.RightBottomR.Yaw < t.eYaw {
		t.eYaw = t.RightBottomR.Yaw
	}

	if t.oYaw < t.eYaw { // 원점을 지남
		t.hasOrigin = true
		t.yawFov = 360 + t.oYaw - t.eYaw
	} else {
		t.yawFov = t.oYaw - t.eYaw
	}

	t.oPitch = t.LeftTopR.Pitch
	if t.RightTopR.Pitch > t.oPitch {
		t.oPitch = t.RightTopR.Pitch
	}

	t.ePitch = t.LeftBottomR.Pitch
	if t.RightBottomR.Pitch < t.ePitch {
		t.ePitch = t.RightBottomR.Pitch
	}

	t.pitchFov = t.oPitch - t.ePitch

	t.oRoll = (t.LeftTopR.Roll + t.RightTopR.Roll + t.LeftBottomR.Roll + t.RightBottomR.Roll) / 4
	t.IsCalibrated = true
}

func (t *Transform) ScreenCoordinate(cR Rotation) Coordinate {
	var yawDiff, pitchDiff float64
	// X == get Yaw Difference -> tan normalize
	// yawDiff 가 크면 오른쪽으로 많이 움직임
	if t.hasOrigin {
		if cR.Yaw < t.oYaw { // 왼쪽 ~ 원점
			yawDiff = t.oYaw - cR.Yaw
		} else if cR.Yaw > t.oYaw && cR.Yaw < t.oYaw+90 { // 왼쪽 + 90
			yawDiff = 0
		} else { // 원점 ~ 오른쪽
			yawDiff = t.oYaw + cR.Yaw - t.eYaw
		}
	} else {
		if cR.Yaw > t.oYaw { // 원점보다 더 왼쪽
			yawDiff = 0
		} else if cR.Yaw < t.eYaw { // 끝 보다 더 오른쪽
			yawDiff = t.yawFov
		} else { // fov 내부
			yawDiff = t.oYaw - cR.Yaw
		}
	}

	X := math.Tan(deg2rad(yawDiff)) / math.Tan(deg2rad(t.yawFov)) * SCREEN_X
	if X < 0 {
		X = 0
	} else if X > SCREEN_X {
		X = SCREEN_X
	}

	// Y == get Pitch Diff -> tan normalize
	// pitchDiff 가 크면 아래로 많이 움직임
	pitchDiff = t.oPitch - cR.Pitch
	Y := math.Tan(deg2rad(pitchDiff)) / math.Tan(deg2rad(t.pitchFov)) * SCREEN_Y
	if Y < 0 {
		Y = 0
	} else if Y > SCREEN_Y {
		Y = SCREEN_Y
	}

	if math.IsNaN(X) {
		X = 0
	}
	if math.IsNaN(Y) {
		Y = 0
	}

	roll := 0
	if cR.Roll < t.oRoll-30 {
		roll = -1
	} else if cR.Roll > t.oRoll+30 {
		roll = 1
	}

	return Coordinate{X, Y, roll}
}

func deg2rad(deg float64) float64 {
	return deg * math.Pi / 180
}
