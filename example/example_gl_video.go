package main

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/yourok/go-mpv/mpv"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	//Init gl
	err := glfw.Init()
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(640, 480, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	//init mpv
	m := mpv.Create()
	defer m.TerminateDestroy()

	m.SetOption("no-resume-playback", mpv.FORMAT_FLAG, true)
	m.SetOption("softvol", mpv.FORMAT_STRING, "yes")
	m.SetOption("volume", mpv.FORMAT_INT64, 20)
	m.SetOption("mute", mpv.FORMAT_FLAG, true)
	m.SetOptionString("hwdec", "auto")

	//cache
	m.SetOption("cache-default", mpv.FORMAT_INT64, 160)
	m.SetOption("cache-seek-min", mpv.FORMAT_INT64, 16)

	//GL
	m.SetOptionString("vo", "opengl-cb")

	err = m.Initialize()
	if err != nil {
		log.Println("Mpv init:", err.Error())
		return
	}
	//Set video file
	m.Command([]string{"loadfile", "http://techslides.com/demos/sample-videos/small.webm"})

	//Init mpv gl
	mgl := m.GetSubApiGL()
	if mgl == nil {
		return
	}

	mgl.InitGL()
	defer mgl.UninitGL()

	//Draw
	for !window.ShouldClose() {
		mgl.Draw(0, 640, 480)
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
