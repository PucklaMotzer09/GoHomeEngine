package gohome

import (
	// "fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"runtime"
	// "strconv"
	"strings"
)

const (
	VERTEX       uint8 = 0
	FRAGMENT     uint8 = 1
	GEOMETRY     uint8 = 2
	TESSELLETION uint8 = 3
	EVELUATION   uint8 = 4
	COMPUTE      uint8 = 5
)

type Shader interface {
	AddShader(shader_type uint8, src string) error
	Link() error
	Setup() error
	Terminate()
	Use()
	Unuse()
	SetUniformV2(name string, value mgl32.Vec2) error
	SetUniformV3(name string, value mgl32.Vec3) error
	SetUniformV4(name string, value mgl32.Vec4) error
	SetUniformF(name string, value float32) error
	SetUniformI(name string, value int32) error
	SetUniformUI(name string, value uint32) error
	SetUniformB(name string, value uint8) error
	SetUniformM2(name string, value mgl32.Mat2) error
	SetUniformM3(name string, value mgl32.Mat3) error
	SetUniformM4(name string, value mgl32.Mat4) error
	SetUniformMaterial(mat Material) error
	SetUniformLights(lightCollectionIndex int32) error
	GetName() string
}

type OpenGLShader struct {
	program             uint32
	name                string
	shaders             [6]uint32
	uniform_locations   map[string]int32
	attribute_locations map[string]uint32
}

func CreateOpenGLShader(name string) (*OpenGLShader, error) {
	shader := &OpenGLShader{
		program:             0,
		name:                name,
		shaders:             [6]uint32{0, 0, 0, 0, 0, 0},
		uniform_locations:   make(map[string]int32),
		attribute_locations: make(map[string]uint32),
	}
	program := gl.CreateProgram()
	if program == 0 {
		return shader, &OpenGLError{errorString: "Couldn't create shader program of " + name}
	} else {
		shader.program = program
		return shader, nil
	}
}

func getShaderTypeName(shader_type uint8) string {
	var shader_type_name string
	switch shader_type {
	case VERTEX:
		shader_type_name = "Vertex Shader"
	case FRAGMENT:
		shader_type_name = "Fragment Shader"
	case GEOMETRY:
		shader_type_name = "Geometry Shader"
	case TESSELLETION:
		shader_type_name = "Tesselletion Shader"
	case EVELUATION:
		shader_type_name = "Eveluation Shader"
	case COMPUTE:
		shader_type_name = "Compute Shader"
	}

	return shader_type_name
}

func compileOpenGLShader(shader_type uint32, src **uint8, program uint32) (uint32, error) {
	shader := gl.CreateShader(shader_type)
	gl.ShaderSource(shader, 1, src, nil)
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		logText := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(logText))

		return 0, &OpenGLError{errorString: logText}
	}
	gl.AttachShader(program, shader)

	return 0, nil
}

func (s *OpenGLShader) AddShader(shader_type uint8, src string) error {
	csource, free := gl.Strs(src + "\x00")
	defer free()
	var err error
	var shaderName uint32
	switch shader_type {
	case VERTEX:
		shaderName, err = compileOpenGLShader(gl.VERTEX_SHADER, csource, s.program)
	case FRAGMENT:
		shaderName, err = compileOpenGLShader(gl.FRAGMENT_SHADER, csource, s.program)
	case GEOMETRY:
		shaderName, err = compileOpenGLShader(gl.GEOMETRY_SHADER, csource, s.program)
	case TESSELLETION:
		shaderName, err = compileOpenGLShader(gl.TESS_CONTROL_SHADER, csource, s.program)
	case EVELUATION:
		shaderName, err = compileOpenGLShader(gl.TESS_EVALUATION_SHADER, csource, s.program)
	case COMPUTE:
		shaderName, err = compileOpenGLShader(gl.COMPUTE_SHADER, csource, s.program)
	}

	if err != nil {
		return &OpenGLError{errorString: "Couldn't compile source of " + getShaderTypeName(shader_type) + " of " + s.name + ": " + err.Error()}
	}

	s.shaders[shader_type] = shaderName

	return nil
}

func (s *OpenGLShader) deleteAllShaders() {
	for i := 0; i < 6; i++ {
		if s.shaders[i] != 0 {
			gl.DetachShader(s.program, s.shaders[i])
			gl.DeleteShader(s.shaders[i])
		}
	}
}

func (s *OpenGLShader) Link() error {
	defer s.deleteAllShaders()

	gl.LinkProgram(s.program)

	var status int32
	gl.GetProgramiv(s.program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(s.program, gl.INFO_LOG_LENGTH, &logLength)

		logtext := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(s.program, logLength, nil, gl.Str(logtext))

		return &OpenGLError{errorString: "Couldn't link shader " + s.name + ": " + logtext}
	}

	return nil
}

func (s *OpenGLShader) Use() {
	gl.UseProgram(s.program)
}

func (*OpenGLShader) Unuse() {
	gl.UseProgram(0)
}

func (s *OpenGLShader) Setup() error {

	// var vao uint32
	// gl.CreateVertexArrays(1, &vao)
	// gl.BindVertexArray(vao)
	// defer gl.DeleteVertexArrays(1, &vao)
	// defer gl.BindVertexArray(0)
	// for i := 0; i < 100; i++ {
	// 	s.SetUniformI("pointLights["+strconv.Itoa(i)+"].shadowmap", 1)
	// }
	// gl.ValidateProgram(s.program)
	// var status int32
	// gl.GetProgramiv(s.program, gl.VALIDATE_STATUS, &status)
	// if status == gl.FALSE {
	// 	var logLength int32
	// 	gl.GetProgramiv(s.program, gl.INFO_LOG_LENGTH, &logLength)

	// 	logtext := strings.Repeat("\x00", int(logLength+1))
	// 	gl.GetProgramInfoLog(s.program, logLength, nil, gl.Str(logtext))

	// 	return &OpenGLError{errorString: "Couldn't validate shader " + s.name + ": " + logtext}
	// }

	if runtime.GOOS != "windows" {
		s.Use()

		var c int32
		var i uint32
		s.uniform_locations = make(map[string]int32)
		s.attribute_locations = make(map[string]uint32)
		gl.GetProgramiv(s.program, gl.ACTIVE_UNIFORMS, &c)
		for i = 0; i < uint32(c); i++ {
			var buf [256]byte
			gl.GetActiveUniform(s.program, i, 256, nil, nil, nil, &buf[0])
			loc := gl.GetUniformLocation(s.program, &buf[0])
			name := gl.GoStr(&buf[0])
			s.uniform_locations[name] = loc
		}
		gl.GetProgramiv(s.program, gl.ACTIVE_ATTRIBUTES, &c)
		for i = 0; i < uint32(c); i++ {
			var buf [256]byte
			gl.GetActiveAttrib(s.program, i, 256, nil, nil, nil, &buf[0])
			loc := gl.GetAttribLocation(s.program, &buf[0])
			name := gl.GoStr(&buf[0])
			s.attribute_locations[name] = uint32(loc)
		}

		s.Unuse()
	}

	return nil
}

func (s *OpenGLShader) Terminate() {
	gl.DeleteProgram(s.program)
}

func (s *OpenGLShader) SetUniformV2(name string, value mgl32.Vec2) error {
	var loc int32
	var ok bool
	if loc, ok = s.uniform_locations[name]; !ok {
		loc = gl.GetUniformLocation(s.program, gl.Str(name+"\x00"))
		s.uniform_locations[name] = loc
	}
	if loc == -1 {
		return &OpenGLError{errorString: "Couldn't find uniform " + name + " in shader " + s.name}
	}
	gl.Uniform2f(loc, value[0], value[1])

	return nil
}
func (s *OpenGLShader) SetUniformV3(name string, value mgl32.Vec3) error {
	var loc int32
	var ok bool
	if loc, ok = s.uniform_locations[name]; !ok {
		loc = gl.GetUniformLocation(s.program, gl.Str(name+"\x00"))
		s.uniform_locations[name] = loc
	}
	if loc == -1 {
		return &OpenGLError{errorString: "Couldn't find uniform " + name + " in shader " + s.name}
	}

	gl.Uniform3f(loc, value[0], value[1], value[2])

	return nil
}
func (s *OpenGLShader) SetUniformV4(name string, value mgl32.Vec4) error {
	var loc int32
	var ok bool
	if loc, ok = s.uniform_locations[name]; !ok {
		loc = gl.GetUniformLocation(s.program, gl.Str(name+"\x00"))
		s.uniform_locations[name] = loc
	}
	if loc == -1 {
		return &OpenGLError{errorString: "Couldn't find uniform " + name + " in shader " + s.name}
	}

	gl.Uniform4f(loc, value[0], value[1], value[2], value[3])

	return nil
}
func (s *OpenGLShader) SetUniformF(name string, value float32) error {
	var loc int32
	var ok bool
	if loc, ok = s.uniform_locations[name]; !ok {
		loc = gl.GetUniformLocation(s.program, gl.Str(name+"\x00"))
		s.uniform_locations[name] = loc
	}
	if loc == -1 {
		return &OpenGLError{errorString: "Couldn't find uniform " + name + " in shader " + s.name}
	}

	gl.Uniform1f(loc, value)

	return nil
}
func (s *OpenGLShader) SetUniformI(name string, value int32) error {
	var loc int32
	var ok bool
	if loc, ok = s.uniform_locations[name]; !ok {
		loc = gl.GetUniformLocation(s.program, gl.Str(name+"\x00"))
		s.uniform_locations[name] = loc
	}
	if loc == -1 {
		return &OpenGLError{errorString: "Couldn't find uniform " + name + " in shader " + s.name}
	}

	gl.Uniform1i(loc, value)

	return nil
}
func (s *OpenGLShader) SetUniformUI(name string, value uint32) error {
	var loc int32
	var ok bool
	if loc, ok = s.uniform_locations[name]; !ok {
		loc = gl.GetUniformLocation(s.program, gl.Str(name+"\x00"))
		s.uniform_locations[name] = loc
	}
	if loc == -1 {
		return &OpenGLError{errorString: "Couldn't find uniform " + name + " in shader " + s.name}
	}

	gl.Uniform1ui(loc, value)

	return nil
}
func (s *OpenGLShader) SetUniformB(name string, value uint8) error {
	var loc int32
	var ok bool
	if loc, ok = s.uniform_locations[name]; !ok {
		loc = gl.GetUniformLocation(s.program, gl.Str(name+"\x00"))
		s.uniform_locations[name] = loc
	}
	if loc == -1 {
		return &OpenGLError{errorString: "Couldn't find uniform " + name + " in shader " + s.name}
	}

	gl.Uniform1i(loc, int32(value))

	return nil
}
func (s *OpenGLShader) SetUniformM2(name string, value mgl32.Mat2) error {
	var loc int32
	var ok bool
	if loc, ok = s.uniform_locations[name]; !ok {
		loc = gl.GetUniformLocation(s.program, gl.Str(name+"\x00"))
		s.uniform_locations[name] = loc
	}
	if loc == -1 {
		return &OpenGLError{errorString: "Couldn't find uniform " + name + " in shader " + s.name}
	}

	gl.UniformMatrix2fv(loc, 1, false, &value[0])

	return nil
}
func (s *OpenGLShader) SetUniformM3(name string, value mgl32.Mat3) error {
	var loc int32
	var ok bool
	if loc, ok = s.uniform_locations[name]; !ok {
		loc = gl.GetUniformLocation(s.program, gl.Str(name+"\x00"))
		s.uniform_locations[name] = loc
	}
	if loc == -1 {
		return &OpenGLError{errorString: "Couldn't find uniform " + name + " in shader " + s.name}
	}

	gl.UniformMatrix3fv(loc, 1, false, &value[0])

	return nil
}
func (s *OpenGLShader) SetUniformM4(name string, value mgl32.Mat4) error {
	var loc int32
	var ok bool
	if loc, ok = s.uniform_locations[name]; !ok {
		loc = gl.GetUniformLocation(s.program, gl.Str(name+"\x00"))
		s.uniform_locations[name] = loc
	}
	if loc == -1 {
		return &OpenGLError{errorString: "Couldn't find uniform " + name + " in shader " + s.name}
	}

	gl.UniformMatrix4fv(loc, 1, false, &value[0])

	return nil
}

func (s *OpenGLShader) SetUniformMaterial(mat Material) error {
	var err error
	rnd, _ := Render.(*OpenGLRenderer)
	var diffBind int32 = 0
	var specBind int32 = 0
	var normBind int32 = 0
	var boundTextures uint32

	if mat.DiffuseTexture != nil {
		diffBind = int32(rnd.CurrentTextureUnit)
		rnd.CurrentTextureUnit++
		mat.DiffuseTexture.Bind(uint32(diffBind))
		mat.diffuseTextureLoaded = 1
		boundTextures++
	} else {
		mat.diffuseTextureLoaded = 0
	}

	if mat.SpecularTexture != nil {
		specBind = int32(rnd.CurrentTextureUnit)
		rnd.CurrentTextureUnit++
		mat.SpecularTexture.Bind(uint32(specBind))
		mat.specularTextureLoaded = 1
		boundTextures++
	} else {
		mat.specularTextureLoaded = 0
	}

	if mat.NormalMap != nil {
		normBind = int32(rnd.CurrentTextureUnit)
		rnd.CurrentTextureUnit++
		mat.NormalMap.Bind(uint32(normBind))
		mat.normalMapLoaded = 1
		boundTextures++
	} else {
		mat.normalMapLoaded = 0
	}

	rnd.CurrentTextureUnit -= boundTextures

	if err = s.SetUniformV3(MATERIAL_UNIFORM_NAME+"."+MATERIAL_DIFFUSE_COLOR_UNIFORM_NAME, colorToVec3(mat.DiffuseColor)); err != nil {
		// return err
	}
	if err = s.SetUniformV3(MATERIAL_UNIFORM_NAME+"."+MATERIAL_SPECULAR_COLOR_UNIFORM_NAME, colorToVec3(mat.SpecularColor)); err != nil {
		// return err
	}
	if err = s.SetUniformI(MATERIAL_UNIFORM_NAME+"."+MATERIAL_DIFFUSE_TEXTURE_UNIFORM_NAME, diffBind); err != nil {
		// return err
	}
	if err = s.SetUniformI(MATERIAL_UNIFORM_NAME+"."+MATERIAL_SPECULAR_TEXTURE_UNIFORM_NAME, specBind); err != nil {
		// return err
	}
	if err = s.SetUniformI(MATERIAL_UNIFORM_NAME+"."+MATERIAL_NORMALMAP_UNIFORM_NAME, normBind); err != nil {
		// return err
	}
	if err = s.SetUniformF(MATERIAL_UNIFORM_NAME+"."+MATERIAL_SHINYNESS_UNIFORM_NAME, mat.Shinyness); err != nil {
		// return err
	}

	if err = s.SetUniformB(MATERIAL_UNIFORM_NAME+"."+MATERIAL_DIFFUSE_TEXTURE_LOADED_UNIFORM_NAME, mat.diffuseTextureLoaded); err != nil {
		// return err
	}
	if err = s.SetUniformB(MATERIAL_UNIFORM_NAME+"."+MATERIAL_SPECULAR_TEXTURE_LOADED_UNIFORM_NAME, mat.specularTextureLoaded); err != nil {
		// return err
	}
	if err = s.SetUniformB(MATERIAL_UNIFORM_NAME+"."+MATERIAL_NORMALMAP_LOADED_UNIFORM_NAME, mat.normalMapLoaded); err != nil {
		// return err
	}

	return err
}

func (s *OpenGLShader) SetUniformLights(lightCollectionIndex int32) error {
	if lightCollectionIndex == -1 || lightCollectionIndex > int32(len(LightMgr.lightCollections)-1) {
		var err error
		if err = s.SetUniformUI(NUM_POINT_LIGHTS_UNIFORM_NAME, 0); err != nil {
			// return err
		}
		if err = s.SetUniformUI(NUM_DIRECTIONAL_LIGHTS_UNIFORM_NAME, 0); err != nil {
			// return err
		}
		if err = s.SetUniformUI(NUM_SPOT_LIGHTS_UNIFORM_NAME, 0); err != nil {
			// return err
		}

		if err = s.SetUniformV3(AMBIENT_LIGHT_UNIFORM_NAME, mgl32.Vec3{1.0, 1.0, 1.0}); err != nil {
			// return err
		}
		return nil
	}

	lightColl := LightMgr.lightCollections[lightCollectionIndex]

	var err error
	if err = s.SetUniformUI(NUM_POINT_LIGHTS_UNIFORM_NAME, uint32(len(lightColl.PointLights))); err != nil {
		// return err
	}
	if err = s.SetUniformUI(NUM_DIRECTIONAL_LIGHTS_UNIFORM_NAME, uint32(len(lightColl.DirectionalLights))); err != nil {
		// return err
	}
	if err = s.SetUniformUI(NUM_SPOT_LIGHTS_UNIFORM_NAME, uint32(len(lightColl.SpotLights))); err != nil {
		// return err
	}

	if err = s.SetUniformV3(AMBIENT_LIGHT_UNIFORM_NAME, colorToVec3(lightColl.AmbientLight)); err != nil {
		// return err
	}

	var i uint32
	for i = 0; i < uint32(len(lightColl.PointLights)); i++ {
		if err = lightColl.PointLights[i].SetUniforms(s, i); err != nil {
			// return err
		}
	}
	for i = 0; i < uint32(len(lightColl.DirectionalLights)); i++ {
		if err = lightColl.DirectionalLights[i].SetUniforms(s, i); err != nil {
			// return err
		}
	}
	for i = 0; i < uint32(len(lightColl.SpotLights)); i++ {
		if err = lightColl.SpotLights[i].SetUniforms(s, i); err != nil {
			// return err
		}
	}

	return err
}

func (s *OpenGLShader) GetName() string {
	return s.name
}
