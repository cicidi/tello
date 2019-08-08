package model

const (
	TAKEOFF           = "takeoff"
	LAND              = "land"
	UP                = "up"
	DOWN              = "down"
	LEFT              = "left"
	RIGHT             = "right"
	FORWARD           = "forward"
	BACKWARD          = "backward"
	CLOCKWISE         = "CLOCKWISE"
	COUNTER_CLOCKWISE = "COUNTER_CLOCKWISE"
	HOVER             = "HOVER"
)

type Command struct {
	Name string
	Val  int
}
