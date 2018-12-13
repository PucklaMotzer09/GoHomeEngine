package renderer

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	gl "github.com/PucklaMotzer09/android-go/gles2"
	"image"
	"image/color"
	"strconv"
	"sync"
	"unsafe"
)

type OpenGLESTexture struct {
	width   int
	height  int
	name    string
	oglName uint32

	keyColor color.Color
	modColor color.Color
}

func (ogltex *OpenGLESTexture) bindingPoint() uint32 {
	return gl.TEXTURE_2D
}

func CreateOpenGLESTexture(name string) *OpenGLESTexture {
	tex := &OpenGLESTexture{
		name: name,
	}
	return tex
}

func printOGLTexture2DError(ogltex *OpenGLESTexture, data []byte, width, height int) {
	err := gl.GetError()
	if err != gl.NO_ERROR {
		var errString string

		if err == gl.INVALID_VALUE {
			if width < 0 {
				errString = "width is less than 0 "
			} else if width > gl.MAX_TEXTURE_SIZE {
				errString = "width is too large: " + strconv.Itoa(gl.MAX_TEXTURE_SIZE) + " "
			}
			if height < 0 {
				errString = "height is less than 0"
			} else if height > gl.MAX_TEXTURE_SIZE {
				errString = "height is too large: " + strconv.Itoa(gl.MAX_TEXTURE_SIZE)
			}
			if errString == "" {
				errString = "Invalid Value"
			}
		} else if err == gl.INVALID_ENUM {
			if ogltex.bindingPoint() != gl.TEXTURE_2D {
				errString = "target should be TEXTURE_2D"
			}
			if errString == "" {
				errString = "Invalid Enum"
			}

		} else if err == gl.INVALID_OPERATION {
			errString = "Invalid Operation"
		}

		gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_ERROR, "Texture", ogltex.GetName(), "Couldn't load data: "+strconv.Itoa(int(err))+": "+errString)
	}
}

func (ogltex *OpenGLESTexture) Load(data []byte, width, height int, shadowMap bool) error {
	ogltex.width = width
	ogltex.height = height
	var tex [1]uint32
	gl.GenTextures(1, tex[:])
	ogltex.oglName = tex[0]

	gl.BindTexture(ogltex.bindingPoint(), ogltex.oglName)

	gl.TexParameterf(ogltex.bindingPoint(), gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameterf(ogltex.bindingPoint(), gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameterf(ogltex.bindingPoint(), gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameterf(ogltex.bindingPoint(), gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameterf(ogltex.bindingPoint(), gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	var ptr unsafe.Pointer
	if data == nil {
		ptr = nil
	} else {
		ptr = unsafe.Pointer(&data[0])
	}
	gl.GetError()
	if shadowMap {
		gl.TexImage2D(ogltex.bindingPoint(), 0, gl.DEPTH_COMPONENT, int32(ogltex.width), int32(ogltex.height), 0, gl.DEPTH_COMPONENT, gl.FLOAT, ptr)
	} else {
		gl.TexImage2D(ogltex.bindingPoint(), 0, gl.RGBA, int32(ogltex.width), int32(ogltex.height), 0, gl.RGBA, gl.UNSIGNED_BYTE, ptr)
	}
	printOGLTexture2DError(ogltex, data, width, height)

	gl.GenerateMipmap(ogltex.bindingPoint())

	gl.BindTexture(ogltex.bindingPoint(), 0)

	return nil

}

func loadImageData(img_data *[]byte, img image.Image, start_width, end_width, max_width, max_height uint32, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	var y uint32
	var x uint32
	var r, g, b, a uint32
	var color color.Color
	for x = start_width; x < max_width && x < end_width; x++ {
		for y = 0; y < max_height; y++ {
			color = img.At(int(x), int(y))
			r, g, b, a = color.RGBA()
			(*img_data)[(x+y*max_width)*4+0] = byte(float64(r) / float64(0xffff) * float64(255.0))
			(*img_data)[(x+y*max_width)*4+1] = byte(float64(g) / float64(0xffff) * float64(255.0))
			(*img_data)[(x+y*max_width)*4+2] = byte(float64(b) / float64(0xffff) * float64(255.0))
			(*img_data)[(x+y*max_width)*4+3] = byte(float64(a) / float64(0xffff) * float64(255.0))
		}
	}
}

func (ogltex *OpenGLESTexture) LoadFromImage(img image.Image) error {

	width := img.Bounds().Size().X
	height := img.Bounds().Size().Y

	img_data := make([]byte, width*height*4)

	var wg1 sync.WaitGroup
	var i float32
	deltaWidth := float32(width) / float32(gohome.NUM_GO_ROUTINES_TEXTURE_LOADING)
	wg1.Add(int(gohome.NUM_GO_ROUTINES_TEXTURE_LOADING + 1))
	for i = 0; i <= float32(gohome.NUM_GO_ROUTINES_TEXTURE_LOADING); i++ {
		go loadImageData(&img_data, img, uint32(i*deltaWidth), uint32((i+1)*deltaWidth), uint32(width), uint32(height), &wg1)
	}
	wg1.Wait()

	ogltex.Load(img_data, width, height, false)

	return nil
}

func toTextureUnit(unit uint32) uint32 {
	return gl.TEXTURE0 + unit
}

func (ogltex *OpenGLESTexture) Bind(unit uint32) {
	newUnit := toTextureUnit(unit)
	gl.GetError()
	gl.ActiveTexture(newUnit)
	handleOpenGLError("Texture", ogltex.name, "glActiveTexture in Bind with "+strconv.Itoa(int(unit)))
	gl.BindTexture(ogltex.bindingPoint(), ogltex.oglName)
	handleOpenGLError("Texture", ogltex.name, "glBindTexture in Bind")
}
func (ogltex *OpenGLESTexture) Unbind(unit uint32) {
	newUnit := toTextureUnit(unit)
	gl.GetError()
	gl.ActiveTexture(newUnit)
	handleOpenGLError("Texture", ogltex.name, "glActiveTexture in Unbind with "+strconv.Itoa(int(unit)))
	gl.BindTexture(ogltex.bindingPoint(), 0)
	handleOpenGLError("Texture", ogltex.name, "glBindTexture in Unbind")
}
func (ogltex *OpenGLESTexture) GetWidth() int {
	return ogltex.width
}
func (ogltex *OpenGLESTexture) GetHeight() int {
	return ogltex.height
}

func (ogltex *OpenGLESTexture) Terminate() {
	var tex [1]uint32
	tex[0] = ogltex.oglName
	gl.DeleteTextures(1, tex[:])
}

func (ogltex *OpenGLESTexture) SetFiltering(filtering uint32) {
	gl.BindTexture(ogltex.bindingPoint(), ogltex.oglName)
	var filter int32
	if filtering == gohome.FILTERING_NEAREST {
		filter = gl.NEAREST
	} else if filtering == gohome.FILTERING_LINEAR {
		filter = gl.LINEAR
	} else {
		filter = gl.NEAREST
	}
	gl.TexParameteri(ogltex.bindingPoint(), gl.TEXTURE_MIN_FILTER, filter)
	gl.TexParameteri(ogltex.bindingPoint(), gl.TEXTURE_MAG_FILTER, filter)

	gl.BindTexture(ogltex.bindingPoint(), 0)
}

func (ogltex *OpenGLESTexture) SetWrapping(wrapping uint32) {
	gl.BindTexture(ogltex.bindingPoint(), ogltex.oglName)
	var wrap int32
	if wrapping == gohome.WRAPPING_REPEAT {
		wrap = gl.REPEAT
	} else if wrapping == gohome.WRAPPING_CLAMP_TO_EDGE {
		wrap = gl.CLAMP_TO_EDGE
	} else if wrapping == gohome.WRAPPING_MIRRORED_REPEAT {
		wrap = gl.MIRRORED_REPEAT
	} else {
		wrap = gl.CLAMP_TO_EDGE
	}
	gl.TexParameteri(ogltex.bindingPoint(), gl.TEXTURE_WRAP_S, wrap)
	gl.TexParameteri(ogltex.bindingPoint(), gl.TEXTURE_WRAP_T, wrap)

	gl.BindTexture(ogltex.bindingPoint(), 0)
}

func (ogltex *OpenGLESTexture) SetBorderColor(col color.Color) {
	gohome.ErrorMgr.Error("Texture", ogltex.name, "SetBorderColor does not work in OpenGLES 2.0")
}

func (ogltex *OpenGLESTexture) SetBorderDepth(depth float32) {
	gohome.ErrorMgr.Error("Texture", ogltex.name, "SetBorderDepth does not work in OpenGLES 2.0")

}

func (ogltex *OpenGLESTexture) SetKeyColor(col color.Color) {
	ogltex.keyColor = col
}

func (ogltex *OpenGLESTexture) GetKeyColor() color.Color {
	return ogltex.keyColor
}

func (ogltex *OpenGLESTexture) SetModColor(col color.Color) {
	ogltex.modColor = col
}

func (ogltex *OpenGLESTexture) GetModColor() color.Color {
	return ogltex.modColor
}

func (ogltex *OpenGLESTexture) GetName() string {
	return ogltex.name
}

func (ogltex *OpenGLESTexture) GetData() (data []byte, width int, height int) {
	width, height = ogltex.GetWidth(), ogltex.GetHeight()
	gohome.ErrorMgr.Error("Texture", ogltex.name, "GetData does not work in OpenGLES 2.0")
	return
}
