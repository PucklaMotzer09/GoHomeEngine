package renderer

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	gl "github.com/PucklaMotzer09/android-go/gles2"
	"unsafe"
)

type OpenGLES2Shape3DInterface struct {
	Name   string
	vbo    uint32
	loaded bool

	points      []gohome.Shape3DVertex
	drawMode    uint32
	numVertices uint32
	lineWidth   float32
}

func (this *OpenGLES2Shape3DInterface) Init() {
	this.loaded = false
}

func (this *OpenGLES2Shape3DInterface) AddPoints(points []gohome.Shape3DVertex) {
	if this.loaded {
		gohome.ErrorMgr.Warning("Shape3DInterface", this.Name, "It has already been loaded to the GPU! You can't add any vertices anymore!")
		return
	}
	gohome.Framew.Log("Adding points:", points)
	this.points = append(this.points, points...)
}

func (this *OpenGLES2Shape3DInterface) GetPoints() []gohome.Shape3DVertex {
	return this.points
}

func (this *OpenGLES2Shape3DInterface) SetDrawMode(drawMode uint8) {
	switch drawMode {
	case gohome.DRAW_MODE_POINTS:
		this.drawMode = gl.POINTS
	case gohome.DRAW_MODE_LINES:
		this.drawMode = gl.LINES
	case gohome.DRAW_MODE_TRIANGLES:
		this.drawMode = gl.TRIANGLES
	default:
		this.drawMode = gl.POINTS
	}
}

func (this *OpenGLES2Shape3DInterface) SetPointSize(size float32) {
}

func (this *OpenGLES2Shape3DInterface) SetLineWidth(width float32) {
	this.lineWidth = width
}

func (this *OpenGLES2Shape3DInterface) attributePointer() {
	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, gohome.SHAPE3DVERTEXSIZE, gl.PtrOffset(0))
	gl.VertexAttribPointer(1, 4, gl.FLOAT, gl.FALSE, gohome.SHAPE3DVERTEXSIZE, gl.PtrOffset(3*4))

	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
}

func (this *OpenGLES2Shape3DInterface) Load() {
	if this.loaded {
		return
	}

	this.numVertices = uint32(len(this.points))
	if this.numVertices == 0 {
		gohome.ErrorMgr.Error("Shape3DInterface", this.Name, "No Vertices have been added!")
		return
	}

	var buf [1]uint32
	gl.GenBuffers(1, buf[:])
	this.vbo = buf[0]

	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, int(gohome.SHAPE3DVERTEXSIZE*this.numVertices), unsafe.Pointer(&this.points[0]), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	this.loaded = true
}

func (this *OpenGLES2Shape3DInterface) Render() {
	hasLoaded := this.loaded
	if !hasLoaded {
		this.Load()
	}

	if this.numVertices == 0 {
		gohome.ErrorMgr.Error("Shape3DInterface", this.Name, "No Vertices have been added!")
		return
	}

	gl.LineWidth(this.lineWidth)
	this.attributePointer()
	gl.GetError()
	gl.DrawArrays(this.drawMode, 0, int32(this.numVertices))
	handleOpenGLError("Shape3DInterface", this.Name, "RenderError: ")
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	gl.LineWidth(1.0)

	if !hasLoaded {
		this.Terminate()
	}
}
func (this *OpenGLES2Shape3DInterface) Terminate() {
	var buf [1]uint32
	buf[0] = this.vbo
	gl.DeleteBuffers(1, buf[:])

	this.numVertices = 0
	this.loaded = false
	this.points = this.points[:0]
}