package renderer

import (
	// "fmt"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"golang.org/x/mobile/gl"
	"image/color"
	"log"
)

type OpenGLESRenderTexture struct {
	Name         string
	fbo          gl.Framebuffer
	rbo          gl.Renderbuffer
	depthBuffer  bool
	shadowMap    bool
	cubeMap      bool
	textures     []gohome.Texture
	prevViewport gohome.Viewport
	viewport     gohome.Viewport
	gles         gl.Context
}

func CreateOpenGLESRenderTexture(name string, width, height, textures uint32, depthBuffer, shadowMap, cubeMap bool) *OpenGLESRenderTexture {
	rt := &OpenGLESRenderTexture{}
	render, _ := gohome.Render.(*OpenGLESRenderer)
	rt.gles = render.gles
	rt.Create(name, width, height, textures, depthBuffer, shadowMap, cubeMap)

	return rt
}

func (this *OpenGLESRenderTexture) loadTextures(width, height, textures uint32, cubeMap bool) {
	var i uint32
	for i = 0; i < textures; i++ {
		var ogltex *OpenGLESTexture
		var oglcubemap *OpenGLESCubeMap
		var texture gohome.Texture
		if cubeMap {
			oglcubemap = CreateOpenGLESCubeMap(this.Name)
			texture = oglcubemap
		} else {
			ogltex = CreateOpenGLESTexture(this.Name)
			texture = ogltex
		}
		texture.Load(nil, int(width), int(height), this.shadowMap)
		if cubeMap {
			this.gles.BindTexture(gl.TEXTURE_CUBE_MAP, oglcubemap.oglName)
		} else {
			this.gles.BindTexture(ogltex.bindingPoint(), ogltex.oglName)
		}
		if this.shadowMap {
			if cubeMap {
				this.gles.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_CUBE_MAP_POSITIVE_X, oglcubemap.oglName, 0)
			} else {
				this.gles.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, ogltex.bindingPoint(), ogltex.oglName, 0)
			}
		} else {
			if cubeMap {
				this.gles.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0+gl.Enum(i), gl.TEXTURE_CUBE_MAP_POSITIVE_X, oglcubemap.oglName, 0)
			} else {
				this.gles.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0+gl.Enum(i), ogltex.bindingPoint(), ogltex.oglName, 0)
			}
		}
		if !cubeMap {
			texture.SetFiltering(gohome.FILTERING_LINEAR)
		}
		if cubeMap {
			this.gles.BindTexture(gl.TEXTURE_CUBE_MAP, gl.Texture{0})
		} else {
			this.gles.BindTexture(ogltex.bindingPoint(), gl.Texture{0})
		}
		this.textures = append(this.textures, texture)
	}
}

func (this *OpenGLESRenderTexture) loadRenderBuffer(width, height uint32) {
	if this.depthBuffer {
		this.rbo = this.gles.CreateRenderbuffer()
		this.gles.BindRenderbuffer(gl.RENDERBUFFER, this.rbo)
		this.gles.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH24_STENCIL8, int(width), int(height))
		this.gles.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.RENDERBUFFER, this.rbo)
		this.gles.BindRenderbuffer(gl.RENDERBUFFER, gl.Renderbuffer{0})
	}
}

func (this *OpenGLESRenderTexture) Create(name string, width, height, textures uint32, depthBuffer, shadowMap, cubeMap bool) {
	if textures == 0 {
		textures = 1
	}

	this.Name = name
	this.shadowMap = shadowMap
	this.depthBuffer = depthBuffer && !shadowMap
	this.cubeMap = cubeMap

	this.fbo = this.gles.CreateFramebuffer()

	this.gles.BindFramebuffer(gl.FRAMEBUFFER, this.fbo)

	this.loadRenderBuffer(width, height)
	this.loadTextures(width, height, textures, cubeMap)
	if this.gles.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		log.Println("Error creating gohome.RenderTexture: Framebuffer is not complete")
		return
	}
	this.gles.BindFramebuffer(gl.FRAMEBUFFER, gl.Framebuffer{0})

	this.viewport = gohome.Viewport{
		0,
		0, 0,
		int(width), int(height),
	}
}

func (this *OpenGLESRenderTexture) Load(data []byte, width, height int, shadowMap bool) error {
	return &OpenGLESError{errorString: "The Load method of RenderTexture is not used!"}
}

func (this *OpenGLESRenderTexture) GetName() string {
	return this.Name
}

func (this *OpenGLESRenderTexture) SetAsTarget() {
	this.gles.BindFramebuffer(gl.FRAMEBUFFER, this.fbo)
	gohome.Render.ClearScreen(&gohome.Color{0, 0, 0, 0})
	this.prevViewport = gohome.Render.GetViewport()
	gohome.Render.SetViewport(this.viewport)
}

func (this *OpenGLESRenderTexture) UnsetAsTarget() {
	this.gles.BindFramebuffer(gl.FRAMEBUFFER, gl.Framebuffer{0})
	gohome.Render.SetViewport(this.prevViewport)
}

func (this *OpenGLESRenderTexture) Blit(rtex gohome.RenderTexture) {
	// var width int32
	// var height int32
	// var x int32
	// var y int32
	// if rtex != nil {
	// 	rtex.SetAsTarget()
	// 	width = int32(rtex.GetWidth())
	// 	height = int32(rtex.GetHeight())
	// 	x = 0
	// 	y = 0
	// } else {
	// 	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	// 	width = int32(this.prevViewport.Width)
	// 	height = int32(this.prevViewport.Height)
	// 	x = int32(this.prevViewport.X)
	// 	y = int32(this.prevViewport.Y)
	// }
	// gl.BindFramebuffer(gl.READ_FRAMEBUFFER, this.fbo)
	// gl.BlitFramebuffer(0, 0, int32(this.GetWidth()), int32(this.GetHeight()), x, y, width, height, gl.COLOR_BUFFER_BIT|gl.DEPTH_BUFFER_BIT|gl.STENCIL_BUFFER_BIT, gl.NEAREST)
	// if rtex != nil {
	// 	rtex.UnsetAsTarget()
	// } else {
	// 	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, 0)
	// 	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	// }

	log.Println("BlitFramebuffer is not supported by OpenGLES")
}

func (this *OpenGLESRenderTexture) Bind(unit uint32) {
	this.BindIndex(0, unit)
}

func (this *OpenGLESRenderTexture) Unbind(unit uint32) {
	this.UnbindIndex(0, unit)
}

func (this *OpenGLESRenderTexture) BindIndex(index, unit uint32) {
	this.textures[index].Bind(unit)
}

func (this *OpenGLESRenderTexture) UnbindIndex(index, unit uint32) {
	this.textures[index].Unbind(unit)
}

func (this *OpenGLESRenderTexture) GetWidth() int {
	return this.textures[0].GetWidth()
}

func (this *OpenGLESRenderTexture) GetHeight() int {
	return this.textures[0].GetHeight()
}

func (this *OpenGLESRenderTexture) Terminate() {
	defer this.gles.DeleteFramebuffer(this.fbo)
	if this.depthBuffer {
		defer this.gles.DeleteRenderbuffer(this.rbo)
	}
	for i := 0; i < len(this.textures); i++ {
		defer this.textures[i].Terminate()
	}
	this.textures = append(this.textures[:0], this.textures[len(this.textures):]...)
}

func (this *OpenGLESRenderTexture) ChangeSize(width, height uint32) {
	if uint32(this.GetWidth()) != width || uint32(this.GetHeight()) != height {
		textures := uint32(len(this.textures))
		this.Terminate()
		this.Create(this.Name, width, height, textures, this.depthBuffer, this.shadowMap, this.cubeMap)
	}
}

func (this *OpenGLESRenderTexture) SetFiltering(filtering uint32) {
	for i := 0; i < len(this.textures); i++ {
		this.textures[i].SetFiltering(filtering)
	}
}

func (this *OpenGLESRenderTexture) SetWrapping(wrapping uint32) {
	for i := 0; i < len(this.textures); i++ {
		this.textures[i].SetWrapping(wrapping)
	}
}

func (this *OpenGLESRenderTexture) SetBorderColor(col color.Color) {
	for i := 0; i < len(this.textures); i++ {
		this.textures[i].SetBorderColor(col)
	}
}

func (this *OpenGLESRenderTexture) SetBorderDepth(depth float32) {
	for i := 0; i < len(this.textures); i++ {
		this.textures[i].SetBorderDepth(depth)
	}
}