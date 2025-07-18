package worldgen

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/aquilax/go-perlin"
	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/window"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	ChunkSize   = 16
	ChunkRadius = 5
	VoxelSize   = float32(0.1)

	WorldMinY           = -64
	WorldMaxY           = 200
	PlayerMoveThreshold = ChunkSize * VoxelSize / 2
)

var PerlinScale = 0.01

type ChunkCoord struct {
	X, Z int
}

func TestMap3d() {
	app := app.App()
	scene := core.NewNode()

	cam := camera.New(1)
	cam.SetPosition(0, 5, 0)
	scene.Add(cam)

	app.Subscribe(window.OnWindowSize, func(evname string, ev interface{}) {
		width, height := app.GetSize()
		app.Gls().Viewport(0, 0, int32(width), int32(height))
		cam.SetAspect(float32(width) / float32(height))
	})
	app.Gls().ClearColor(0.5, 0.7, 1.0, 1.0)

	scene.Add(light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8))
	pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 100.0)
	pointLight.SetPosition(1, 20, 2)
	scene.Add(pointLight)

	var alpha float64 = 1.5
	var beta float64 = 2.0
	seed := rand.New(rand.NewSource(time.Now().UnixNano())).Int63()
	var p = perlin.NewPerlin(alpha, beta, 3, seed)

	chunkNodes := make(map[ChunkCoord]*core.Node)
	var lastPlayerChunk ChunkCoord

	// Оптимизированная функция генерации чанка
	generateChunk := func(cx, cz int) {
		coord := ChunkCoord{X: cx, Z: cz}
		if _, exists := chunkNodes[coord]; exists {
			return
		}

		geom := geometry.NewGeometry()
		positions := math32.NewArrayF32(0, 0)
		normals := math32.NewArrayF32(0, 0)
		indices := math32.NewArrayU32(0, 0)
		var offset uint32

		// Кубические вершины и нормали для одного вокселя
		cubeVertices := []float32{
			-0.5, -0.5, -0.5, 0.5, -0.5, -0.5, 0.5, 0.5, -0.5, -0.5, 0.5, -0.5, // передняя грань
			-0.5, -0.5, 0.5, 0.5, -0.5, 0.5, 0.5, 0.5, 0.5, -0.5, 0.5, 0.5, // задняя грань
			-0.5, -0.5, -0.5, -0.5, 0.5, -0.5, -0.5, 0.5, 0.5, -0.5, -0.5, 0.5, // левая грань
			0.5, -0.5, -0.5, 0.5, 0.5, -0.5, 0.5, 0.5, 0.5, 0.5, -0.5, 0.5, // правая грань
			-0.5, -0.5, -0.5, -0.5, -0.5, 0.5, 0.5, -0.5, 0.5, 0.5, -0.5, -0.5, // нижняя грань
			-0.5, 0.5, -0.5, -0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, -0.5, // верхняя грань
		}
		cubeNormals := []float32{
			0, 0, -1, 0, 0, -1, 0, 0, -1, 0, 0, -1, // передняя
			0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, // задняя
			-1, 0, 0, -1, 0, 0, -1, 0, 0, -1, 0, 0, // левая
			1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, // правая
			0, -1, 0, 0, -1, 0, 0, -1, 0, 0, -1, 0, // нижняя
			0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, // верхняя
		}
		cubeIndices := []uint32{
			0, 1, 2, 2, 3, 0, // передняя
			4, 5, 6, 6, 7, 4, // задняя
			8, 9, 10, 10, 11, 8, // левая
			12, 13, 14, 14, 15, 12, // правая
			16, 17, 18, 18, 19, 16, // нижняя
			20, 21, 22, 22, 23, 20, // верхняя
		}

		voxelMap := make(map[[3]int]bool)
		for x := 0; x < ChunkSize; x++ {
			for z := 0; z < ChunkSize; z++ {
				wx := cx*ChunkSize + x
				wz := cz*ChunkSize + z
				height := int(math.Round(p.Noise2D(float64(wx)*PerlinScale, float64(wz)*PerlinScale) * 70))
				if height > WorldMaxY {
					height = WorldMaxY
				}
				for y := height; y >= WorldMinY; y-- {
					voxelMap[[3]int{x, y, z}] = true
				}
			}
		}

		// Генерация геометрии с учетом видимых граней
		for x := 0; x < ChunkSize; x++ {
			for z := 0; z < ChunkSize; z++ {
				wx := float32(x) * VoxelSize
				wz := float32(z) * VoxelSize
				height := int(math.Round(p.Noise2D(float64(cx*ChunkSize+x)*PerlinScale, float64(cz*ChunkSize+z)*PerlinScale) * 70))
				if height > WorldMaxY {
					height = WorldMaxY
				}
				for y := height; y > WorldMinY; y-- {
					wy := float32(y) * VoxelSize
					// neighbors := [6]bool{
					// 	voxelMap[[3]int{x, y, z - 1}], // перед
					// 	voxelMap[[3]int{x, y, z + 1}], // зад
					// 	voxelMap[[3]int{x - 1, y, z}], // лево
					// 	voxelMap[[3]int{x + 1, y, z}], // право
					// 	voxelMap[[3]int{x, y - 1, z}], // низ
					// 	voxelMap[[3]int{x, y + 1, z}], // верх
					// }
					for face := 0; face < 6; face++ {
						for i := 0; i < 4; i++ {
							vertexIdx := face*12 + i*3
							positions.Append(
								cubeVertices[vertexIdx]*VoxelSize+wx,
								cubeVertices[vertexIdx+1]*VoxelSize+wy,
								cubeVertices[vertexIdx+2]*VoxelSize+wz,
							)
							normals.Append(
								cubeNormals[vertexIdx],
								cubeNormals[vertexIdx+1],
								cubeNormals[vertexIdx+2],
							)
						}
						for _, idx := range cubeIndices[face*6 : (face+1)*6] {
							indices.Append(idx + offset)
						}
						offset += 4
					}
				}
			}
		}

		// Проверяем, что геометрия не пуста
		if positions.Len() == 0 {
			return // Пропускаем пустой чанк
		}

		geom.AddVBO(gls.NewVBO(positions).AddAttrib(gls.VertexPosition))
		geom.AddVBO(gls.NewVBO(normals).AddAttrib(gls.VertexNormal))
		geom.SetIndices(indices)
		mat := material.NewStandard(math32.NewColor("DarkGreen"))
		chunkNode := core.NewNode()
		chunkNode.Add(graphic.NewMesh(geom, mat))
		chunkNode.SetPosition(float32(cx*ChunkSize)*VoxelSize, 0, float32(cz*ChunkSize)*VoxelSize)
		scene.Add(chunkNode)
		chunkNodes[coord] = chunkNode
	}

	updateChunks := func(playerPos math32.Vector3) {
		playerChunk := ChunkCoord{
			X: int(math.Floor(float64(playerPos.X / (VoxelSize * float32(ChunkSize))))),
			Z: int(math.Floor(float64(playerPos.Z / (VoxelSize * float32(ChunkSize))))),
		}
		if playerChunk == lastPlayerChunk {
			return
		}
		lastPlayerChunk = playerChunk

		newChunks := make(map[ChunkCoord]bool)
		for dx := -ChunkRadius; dx <= ChunkRadius; dx++ {
			for dz := -ChunkRadius; dz <= ChunkRadius; dz++ {
				newChunks[ChunkCoord{X: playerChunk.X + dx, Z: playerChunk.Z + dz}] = true
				generateChunk(playerChunk.X+dx, playerChunk.Z+dz)
			}
		}

		for coord, node := range chunkNodes {
			if !newChunks[coord] {
				scene.Remove(node)
				delete(chunkNodes, coord)
			}
		}
	}

	var yaw, pitch float32
	mouseDown := false
	lastMouseX, lastMouseY := 0, 0

	app.Subscribe(window.OnMouseDown, func(evname string, ev interface{}) {
		mouseDown = true
		app.IWindow.(*window.GlfwWindow).SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	})
	app.Subscribe(window.OnCursor, func(evname string, ev interface{}) {
		if mouseDown {
			mev := ev.(*window.CursorEvent)
			dx := float32(mev.Xpos - float32(lastMouseX))
			dy := float32(mev.Ypos - float32(lastMouseY))

			yaw += dx * -0.0005
			pitch += dy * -0.0005

			if pitch > math32.Pi/2 {
				pitch = math32.Pi / 2
			}
			if pitch < -math32.Pi/2 {
				pitch = -math32.Pi / 2
			}
			cam.SetRotation(pitch, yaw, 0)

			lastMouseX = int(mev.Xpos)
			lastMouseY = int(mev.Ypos)
		}
	})

	recreateWorld := func() {
		for _, node := range chunkNodes {
			scene.Remove(node)
		}
		chunkNodes = make(map[ChunkCoord]*core.Node)
		lastPlayerChunk = ChunkCoord{}
		seed = rand.New(rand.NewSource(time.Now().UnixNano())).Int63()
		p = perlin.NewPerlin(alpha, beta, 3, seed)
		fmt.Println("perlin", alpha, beta, PerlinScale)
	}

	label := gui.NewLabel("Тестовый текст")
	label.SetPosition(10, 10)
	label.SetFontSize(18)
	label.SetColor(math32.NewColor("black"))
	scene.Add(label)

	app.Subscribe(window.OnKeyDown, func(evname string, ev interface{}) {
		kev := ev.(*window.KeyEvent)
		updated := false

		step := 0.05
		scaleStep := 0.002
		switch kev.Key {
		case window.Key1:
			alpha -= step
			if alpha < 0.1 {
				alpha = 0.1
			}
			updated = true
		case window.Key2:
			alpha += step
			updated = true
		case window.Key3:
			beta -= step
			if beta < 0.1 {
				beta = 0.1
			}
			updated = true
		case window.Key4:
			beta += step
			updated = true
		case window.Key5:
			PerlinScale -= scaleStep
			updated = true
		case window.Key6:
			PerlinScale += scaleStep
			updated = true
		}

		if updated {
			recreateWorld()
		}
	})

	app.Run(func(r *renderer.Renderer, deltaTime time.Duration) {
		app.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)

		moveSpeed := float32(2)
		pos := cam.Position()
		dir := cam.Direction()

		if app.IWindow.(*window.GlfwWindow).GetKey(glfw.KeyW) == glfw.Press {
			pos.Add(dir.Clone().MultiplyScalar(moveSpeed * float32(deltaTime.Seconds())))
		}
		if app.IWindow.(*window.GlfwWindow).GetKey(glfw.KeyS) == glfw.Press {
			pos.Add(dir.Clone().MultiplyScalar(-moveSpeed * float32(deltaTime.Seconds())))
		}
		if app.IWindow.(*window.GlfwWindow).GetKey(glfw.KeyA) == glfw.Press {
			right := dir.Clone().Cross(math32.NewVector3(0, 1, 0)).Normalize()
			pos.Add(right.MultiplyScalar(-moveSpeed * float32(deltaTime.Seconds())))
		}
		if app.IWindow.(*window.GlfwWindow).GetKey(glfw.KeyD) == glfw.Press {
			right := dir.Clone().Cross(math32.NewVector3(0, 1, 0)).Normalize()
			pos.Add(right.MultiplyScalar(moveSpeed * float32(deltaTime.Seconds())))
		}

		cam.SetPositionVec(&pos)
		pointLight.SetPositionVec(&pos)
		updateChunks(pos)
		label.SetText(
			fmt.Sprintf("Alpha: %.2f\nBeta: %.2f\nFPS: %.2f", alpha, beta, 1.0/deltaTime.Seconds()),
		)
		r.Render(scene, cam)
	})
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
