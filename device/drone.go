package device

import (
	"errors"
	"io" //  It provides basic interfaces to I/O primitives
	"log"
	"os/exec"
	"strconv" // Package strconv implements conversions to and from string
	"sync/atomic"
	"tello/model"
	"tello/service"
	"time" //For time related operation

	"gobot.io/x/gobot"                     // Gobot Framework.
	"gobot.io/x/gobot/platforms/dji/tello" // DJI Tello package.
	"gocv.io/x/gocv"                       // GoCV package to access the OpenCV library.
)

// Frame size constant.
const (
	frameX    = 960
	frameY    = 720
	frameSize = frameX * frameY * 3
)

var lastcmdTime uint64

func Start() {
	lastcmdTime = service.GetTimestamp()
	// Driver: Tello Driver
	drone := tello.NewDriver("8890")

	// OpenCV window to watch the live video stream from Tello.
	window := gocv.NewWindow("Tello")

	//FFMPEG command to convert the raw video from the drone.
	ffmpeg := exec.Command("ffmpeg", "-hwaccel", "auto", "-hwaccel_device", "opencl", "-i", "pipe:0",
		"-pix_fmt", "bgr24", "-s", strconv.Itoa(frameX)+"x"+strconv.Itoa(frameY), "-f", "rawvideo", "pipe:1")
	ffmpegIn, _ := ffmpeg.StdinPipe()
	ffmpegOut, _ := ffmpeg.StdoutPipe()

	initDrone := func() {
		//Starting FFMPEG.
		if err := ffmpeg.Start(); err != nil {
			log.Println(err)
			return
		}
		// Event: Listening the Tello connect event to start the video streaming.
		_ = drone.On(tello.ConnectedEvent, func(data interface{}) {
			log.Println("Connected to Tello.")
			_ = drone.StartVideo()
			_ = drone.SetVideoEncoderRate(tello.VideoBitRateAuto)
			_ = drone.SetExposure(0)
			go StartController(control, drone)
			//For continued streaming of video.
			gobot.Every(10*time.Millisecond, func() {
				_ = drone.StartVideo()
			})
		})

		//Event: Piping the video data into the FFMPEG function.
		drone.On(tello.VideoFrameEvent, func(data interface{}) {
			pkt := data.([]byte)
			if _, err := ffmpegIn.Write(pkt); err != nil {
				log.Println(err)
			}
		})
	}
	//Robot: Tello Drone
	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{drone},
		initDrone,
	)

	// calling Start(false) lets the Start routine return immediately without an additional blocking goroutine
	robot.Start(false)

	// now handle video frames from ffmpeg stream in main thread, to be macOs friendly
	for {
		if err := upstream(ffmpegOut, window); err != nil {
		}
	}
}

func control(drone *tello.Driver, command *model.Command) {
	current := service.GetTimestamp()
	if command != nil && current > atomic.LoadUint64(&lastcmdTime)+100 {
		lastcmdTime = current
		log.Printf("command %s value %d \n ", command.Name, command.Val)
	} else {
		return
	}
	switch command.Name {
	case model.TAKEOFF:
		_ = drone.TakeOff()
	case model.LAND:
		_ = drone.Land()
	case model.UP:
		_ = drone.Up(command.Val)
	case model.DOWN:
		_ = drone.Down(command.Val)
	case model.LEFT:
		_ = drone.Left(command.Val)
	case model.RIGHT:
		_ = drone.Right(command.Val)
	case model.FORWARD:
		_ = drone.Forward(command.Val)
	case model.BACKWARD:
		_ = drone.Backward(command.Val)
	case model.CLOCKWISE:
		_ = drone.Clockwise(command.Val)
	case model.COUNTER_CLOCKWISE:
		_ = drone.CounterClockwise(command.Val)
	case model.HOVER:
		drone.Hover()
	}
}

func upstream(ffmpegOut io.ReadCloser, window *gocv.Window) error {
	buf := make([]byte, frameSize)
	if _, err := io.ReadFull(ffmpegOut, buf); err != nil {
		log.Println(err)
		return nil
	}
	img, _ := gocv.NewMatFromBytes(frameY, frameX, gocv.MatTypeCV8UC3, buf)
	if img.Empty() {
		return nil
	}

	window.IMShow(img)
	if window.WaitKey(1) >= 0 {
		return errors.New("window wait more than 1")
	}
	return nil
}
