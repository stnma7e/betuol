package graphics

import (
	"image"
	"os"

	"github.com/go-gl/gl"

	"github.com/stnma7e/betuol/common"
)

func GlLoadTexture(filename string) gl.Texture {
	file, err := os.Open(filename)
	if err != nil {
		common.LogErr.Print(err)
	}
	defer file.Close()
	img, _, _ := image.Decode(file)
	rgbaImage, ok := img.(*image.NRGBA)
	if !ok {
		common.LogErr.Println("texture not in RGBA format")
		return gl.Texture(0)
	}

	tex := gl.GenTexture()
	tex.Bind(gl.TEXTURE_2D)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)

	width, height := img.Bounds().Dx(), img.Bounds().Dy()
	data := make([]byte, width*height*4)
	lineLen := width * 4
	dest := len(data) - lineLen
	for src := 0; src < len(rgbaImage.Pix); src += rgbaImage.Stride {
		copy(data[dest:dest+lineLen], rgbaImage.Pix[src:src+rgbaImage.Stride])
		dest -= lineLen
	}

	gl.TexImage2D(gl.TEXTURE_2D,
		0,
		4,
		width, height,
		0, gl.RGBA,
		gl.UNSIGNED_BYTE,
		data,
	)
	return tex
}
