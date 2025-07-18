package worldgen

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"testing"

	"github.com/aquilax/go-perlin"
)

func TestTemperature(t *testing.T) {
	const (
		width   = 512
		height  = 512
		alpha   = 2.0 // persistence
		beta    = 2.0 // lacunarity
		octaves = 2
		scale   = 100.0
	)

	// Инициализируем Perlin
	tempP := perlin.NewPerlin(alpha, beta, octaves, 42)
	humidityP := perlin.NewPerlin(alpha, beta, octaves, 41)

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Нормализуем координаты и получаем значение шума
			temperatureNoise := tempP.Noise2D(float64(x)/scale, float64(y)/scale)
			humidityNoise := humidityP.Noise2D(float64(x)/scale, float64(y)/scale)

			// Преобразуем значение [-1, 1] в [0, 255]
			temperatureVal := uint8(math.Round((temperatureNoise + 1) * 127.5))
			humidityVal := uint8(math.Round((humidityNoise + 1) * 127.5))

			// Используем температуру и влажность для определения типа блока

			img.SetRGBA(x, y, color.RGBA{R: temperatureVal, G: humidityVal, B: 0, A: 255})
		}
	}

	// Сохраняем изображение
	file, err := os.Create("perlin_gray.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}

	println("Изображение сохранено в perlin_gray.png")
}
func TestMountains(t *testing.T) {
	const (
		width   = 5120
		height  = 5120
		alpha   = 1.25 // persistence
		beta    = 1.35 // lacunarity
		octaves = 3
		scale   = 0.004
	)

	// Инициализируем Perlin
	p := perlin.NewPerlin(alpha, beta, octaves, 42)

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Нормализуем координаты и получаем значение шума
			noise := p.Noise2D(float64(x)*scale, float64(y)*scale)

			// Преобразуем значение [-1, 1] в [0, 255]
			val := uint8(math.Round((noise + 1) * 127.5))

			// Используем температуру и влажность для определения типа блока

			img.SetRGBA(x, y, getColor(val))
		}
	}

	// Сохраняем изображение
	file, err := os.Create("mountains.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}

	println("Изображение сохранено в mountains.png")
}

func getColor(val uint8) color.RGBA {
	if val < 80 {
		return color.RGBA{R: val / 10, G: val / 10, B: val, A: 255}
	} else if val < 200 {
		return color.RGBA{R: val / 10, G: val, B: val / 10, A: 255}
	} else {
		return color.RGBA{R: val, G: val, B: val, A: 255}
	}
}
