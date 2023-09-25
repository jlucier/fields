package engine

import (
	"github.com/veandco/go-sdl2/sdl"
	"os"
	"sync"
	"time"
)

const (
	fixedUpdateNs = 16_666_667
)

type App struct {
	Window          *sdl.Window
	Renderer        *sdl.Renderer
	running         bool
	runningMutex    sync.Mutex
	lastFixedUpdate time.Time
}

func CreateApp(title string, width int32, height int32) *App {
	window, err := sdl.CreateWindow(title,
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	return &App{
		window,
		renderer,
		true,
		sync.Mutex{},
		time.Time{},
	}
}

type GameCallbacks struct {
	Render      func(*sdl.Renderer, *sdl.Window)
	HandleKeys  func(*App, *sdl.KeyboardEvent)
	FixedUpdate func(time.Time)
}

func (self *App) Run(cbs GameCallbacks) {
	var exitcode int
	sdl.Main(func() {
		for self.running {
			// updates
			if time.Now().Sub(self.lastFixedUpdate) > fixedUpdateNs {
				self.lastFixedUpdate = time.Now()
				go cbs.FixedUpdate(self.lastFixedUpdate)
			}

			// input
			sdl.Do(func() {
				for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
					switch t := event.(type) {
					case *sdl.QuitEvent:
						self.runningMutex.Lock()
						self.running = false
						self.runningMutex.Unlock()
						exitcode = 0

					case *sdl.KeyboardEvent:
						cbs.HandleKeys(self, t)
					}
				}
			})

			// Render
			sdl.Do(func() {
				self.Renderer.SetDrawColor(0, 0, 0, 255)
				self.Renderer.Clear()
			})

			cbs.Render(self.Renderer, self.Window)

			sdl.Do(func() {
				self.Renderer.Present()
			})
		}
	})
	os.Exit(exitcode)
}

func (self *App) Quit() {
	self.runningMutex.Lock()
	self.running = false
	self.runningMutex.Unlock()
}

func (self *App) Close() {
	sdl.Do(func() {
		self.Window.Destroy()
	})
	sdl.Do(func() {
		self.Renderer.Destroy()
	})
}
