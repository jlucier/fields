package main

import (
	"fmt"
	"github.com/ojrac/opensimplex-go"
	"github.com/veandco/go-sdl2/sdl"
	"jlucier/engine"
	"math"
	"math/rand"
	"time"
)

const (
	windowWidth  = 1280
	windowHeight = 960
	cellSize     = 4
	numParticles = 2048
	maxAge       = 60 * 2
	noiseFactor  = 50
)

type Particle struct {
	pos engine.Vec2
	age int
}

func RandParticle() Particle {
	return Particle{
		age: rand.Intn(maxAge),
		pos: engine.Vec2{
			X: engine.Lerp(0, windowWidth, rand.Float64()),
			Y: engine.Lerp(0, windowHeight, rand.Float64()),
		},
	}
}

type Game struct {
	state     uint
	cellSize  uint
	noise     opensimplex.Noise
	field     []engine.Vec2
	particles []Particle
	text      *sdl.Texture
}

func (self *Game) Close() {
	self.text.Destroy()
}

func (self *Game) cellCenter(x uint, y uint) (uint, uint) {
	return self.cellSize*x + self.cellSize/2, self.cellSize*y + self.cellSize/2
}

func (self *Game) fieldShape() (uint, uint) {
	return windowWidth / self.cellSize, windowHeight / self.cellSize
}

func (self *Game) numVecs() int {
	return len(self.field)
}

func (self *Game) getFieldVec(x uint, y uint) *engine.Vec2 {
	nx, _ := self.fieldShape()
	return &self.field[x+y*nx]
}

func InitGame(app *engine.App, cellSize uint, noiseSeed int64) Game {
	app.Renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)

	cx := windowWidth / cellSize
	cy := windowHeight / cellSize
	ncells := cx * cy

	ptext, err := app.Renderer.CreateTexture(
		sdl.PIXELFORMAT_RGBA8888,
		sdl.TEXTUREACCESS_TARGET,
		windowWidth,
		windowHeight,
	)
	if err != nil {
		panic(err)
	}

	game := Game{
		0,
		cellSize,
		opensimplex.New(noiseSeed),
		make([]engine.Vec2, ncells),
		make([]Particle, numParticles),
		ptext,
	}

	// Initialize vectors
	for i := uint(0); i < cx; i++ {
		for j := uint(0); j < cy; j++ {
			nval := game.noise.Eval2(float64(i)/noiseFactor, float64(j)/noiseFactor) + 1
			ang := nval * math.Pi

			v := engine.Vec2FromAngle(ang)
			*game.getFieldVec(i, j) = v.Mul(nval)
		}
	}

	// Initialize particles
	for i := range game.particles {
		game.particles[i] = RandParticle()
	}

	return game
}

// Updates

func (self *Game) handleKeys(app *engine.App, t *sdl.KeyboardEvent) {
	keyCode := t.Keysym.Sym

	if t.State == sdl.RELEASED {
		switch string(keyCode) {
		case "q":
			app.Quit()
		case "p":
			if self.state == 0 {
				self.state = 1
			} else {
				self.state = 0
			}
		}
	}
}

func (self *Game) moveParticles() {
	for i := range self.particles {
		p := &self.particles[i]

		if p.age++; p.age > maxAge {
			// die
			*p = RandParticle()
		} else {
			// move
			cx := uint(p.pos.X / float64(self.cellSize))
			cy := uint(p.pos.Y / float64(self.cellSize))

			fv := self.getFieldVec(cx, cy)
			p.pos = p.pos.Add(*fv)
			p.pos.X = engine.Clamp(p.pos.X, 0, windowWidth-1)
			p.pos.Y = engine.Clamp(p.pos.Y, 0, windowHeight-1)
		}
	}
}

func (self *Game) fixedUpdate(t time.Time) {
	switch self.state {
	case 0:
		self.moveParticles()
	}
}

// Draw calls

func (self *Game) drawVecs(renderer *sdl.Renderer) {
	sdl.Do(func() {
		nx, ny := self.fieldShape()
		renderer.SetDrawColor(0, 255, 0, 255)

		for i := uint(0); i < nx; i++ {
			for j := uint(0); j < ny; j++ {
				px, py := self.cellCenter(i, j)

				pxVec := engine.Vec2{float64(px), float64(py)}
				v := self.getFieldVec(i, j)

				end := pxVec.Add(v.Mul(float64(self.cellSize / 2)))

				renderer.DrawLine(int32(px), int32(py), int32(end.X), int32(end.Y))
			}
		}
	})
}

func (self *Game) drawParticles(renderer *sdl.Renderer) {
	sdl.Do(func() {
		renderer.SetRenderTarget(self.text)
		self.text.SetBlendMode(sdl.BLENDMODE_BLEND)
		self.text.SetAlphaMod(20)
		renderer.Copy(self.text, nil, nil)

		// render current position
		renderer.SetDrawColor(0, 255, 255, 255)
		for _, p := range self.particles {
			px := int32(p.pos.X)
			py := int32(p.pos.Y)

			renderer.DrawPoint(px, py)
		}
	})

	sdl.Do(func() {
		renderer.SetRenderTarget(nil)
		self.text.SetBlendMode(sdl.BLENDMODE_NONE)
		renderer.Copy(self.text, nil, nil)
	})
}

func (self *Game) render(renderer *sdl.Renderer, window *sdl.Window) {
	// self.drawVecs(renderer)
	self.drawParticles(renderer)
}

func main() {
	seed := int64(0) // time.Now().UnixMicro()
	fmt.Println("seed", seed)

	app := engine.CreateApp("FlowField", windowWidth, windowHeight)
	defer app.Close()
	game := InitGame(app, cellSize, seed)
	defer game.Close()
	app.Run(engine.GameCallbacks{
		Render:      game.render,
		HandleKeys:  game.handleKeys,
		FixedUpdate: game.fixedUpdate,
	})
}
