package renderer

import (
	// "fmt"
	"bytes"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/go-gl/mathgl/mgl32"
	"golang.org/x/mobile/gl"
	// "log"
	"strconv"
)

type OpenGLESShader struct {
	program           gl.Program
	name              string
	shaders           [6]gl.Shader
	uniform_locations map[string]gl.Uniform
	validated         bool
	gles              *gl.Context
}

func CreateOpenGLESShader(name string) (*OpenGLESShader, error) {
	shader := &OpenGLESShader{
		program:           gl.Program{false, 0},
		name:              name,
		shaders:           [6]gl.Shader{gl.Shader{0}, gl.Shader{0}, gl.Shader{0}, gl.Shader{0}, gl.Shader{0}, gl.Shader{0}},
		uniform_locations: make(map[string]gl.Uniform),
		validated:         false,
	}
	renderer, _ := gohome.Render.(*OpenGLESRenderer)
	shader.gles = &renderer.gles
	program := (*shader.gles).CreateProgram()
	if program.Value == 0 {
		return shader, &OpenGLESError{errorString: "Couldn't create shader program of " + name}
	} else {
		shader.program = program
		return shader, nil
	}
}

func getShaderTypeName(shader_type uint8) string {
	var shader_type_name string
	switch shader_type {
	case gohome.VERTEX:
		shader_type_name = "Vertex Shader"
	case gohome.FRAGMENT:
		shader_type_name = "Fragment Shader"
	case gohome.GEOMETRY:
		shader_type_name = "Geometry Shader"
	case gohome.TESSELLETION:
		shader_type_name = "Tesselletion Shader"
	case gohome.EVELUATION:
		shader_type_name = "Eveluation Shader"
	case gohome.COMPUTE:
		shader_type_name = "Compute Shader"
	}

	return shader_type_name
}

func bindAttributesFromFile(program gl.Program, src string, gles gl.Context) {

	var line bytes.Buffer
	var lineString string
	var attributeNames []string
	var curChar byte = ' '
	var curIndex uint32 = 0
	var curWordIndex uint32 = 0
	var curWord uint32 = 0
	var wordBuffer bytes.Buffer
	var wordsString []string
	var readWord bool = false

	for curIndex < uint32(len(src)) {
		for curChar = ' '; curChar != '\n'; curChar = src[curIndex] {
			line.WriteByte(curChar)
			curIndex++
			if curIndex == uint32(len(src)) {
				break
			}
		}

		lineString = line.String()
		readWord = false
		curWord = 0
		for curWordIndex = 0; curWordIndex < uint32(len(lineString)); curWordIndex++ {
			curChar = lineString[curWordIndex]
			if curChar == ' ' {
				if readWord {
					wordsString[curWord] = wordBuffer.String()
					wordBuffer.Reset()
					curWord++
					readWord = false
				}
			} else {
				if !readWord {
					readWord = true
					wordsString = append(wordsString, string(' '))
				}
				wordBuffer.WriteByte(curChar)
			}
		}
		if readWord {
			wordsString[curWord] = wordBuffer.String()
		}
		wordBuffer.Reset()
		line.Reset()
		if len(wordsString) >= 2 && wordsString[0] == "void" && wordsString[1] == "main()" {
			break
		} else if len(wordsString) >= 3 {
			if wordsString[0] == "attribute" {
				if wordsString[2][len(wordsString[2])-1] == ';' {
					wordsString[2] = wordsString[2][0 : len(wordsString[2])-1]
				}
				attributeNames = append(attributeNames, wordsString[2])
			}
		}
		wordsString = append(wordsString[len(wordsString):], wordsString[:0]...)
	}

	for i := 0; i < len(attributeNames); i++ {
		gles.BindAttribLocation(program, gl.Attrib{Value: uint(i)}, attributeNames[i])
	}
}

func CheckOpenGLESError(gles gl.Context, errorString string) *OpenGLESError {
	if err := gles.GetError(); err != gl.NO_ERROR {
		return &OpenGLESError{errorString: errorString + " ErrorCode: " + strconv.Itoa(int(err))}
	}
	return nil
}

func toGohomeShaderType(shader_type gl.Enum) uint8 {
	switch shader_type {
	case gl.VERTEX_SHADER:
		return gohome.VERTEX
	case gl.FRAGMENT_SHADER:
		return gohome.FRAGMENT
	default:
		return 255
	}
	return 255
}

func toShaderName(shader_type gl.Enum) string {
	return getShaderTypeName(toGohomeShaderType(shader_type))
}

func compileOpenGLESShader(shader_name string, shader_type gl.Enum, src string, program gl.Program, gles gl.Context) (gl.Shader, error) {
	gles.GetError()
	shader := gles.CreateShader(shader_type)
	if err := CheckOpenGLESError(gles, "Couldn't create "+toShaderName(shader_type)+" of "+shader_name+":"); err != nil {
		return shader, err
	}
	gles.ShaderSource(shader, src)
	if err := CheckOpenGLESError(gles, "Couldn't source "+toShaderName(shader_type)+" of "+shader_name+":"); err != nil {
		return shader, err
	}
	gles.CompileShader(shader)
	if err := CheckOpenGLESError(gles, "Couldn't call compile function of "+toShaderName(shader_type)+" of "+shader_name+":"); err != nil {
		return shader, err
	}

	status := gles.GetShaderi(shader, gl.COMPILE_STATUS)
	if err := CheckOpenGLESError(gles, "Couldn't get compile status of "+toShaderName(shader_type)+" of "+shader_name+":"); err != nil {
		return shader, err
	}
	if status == gl.FALSE {
		logText := gles.GetShaderInfoLog(shader)
		if err := CheckOpenGLESError(gles, "Couldn't get info log of "+toShaderName(shader_type)+" of "+shader_name+":"); err != nil {
			return shader, err
		}

		return gl.Shader{0}, &OpenGLESError{errorString: logText}
	}
	gles.AttachShader(program, shader)
	if err := CheckOpenGLESError(gles, "Couldn't attach "+toShaderName(shader_type)+" of "+shader_name+":"); err != nil {
		return shader, err
	}
	if shader_type == gl.VERTEX_SHADER {
		bindAttributesFromFile(program, src, gles)
	}
	return shader, nil
}

func (s *OpenGLESShader) AddShader(shader_type uint8, src string) error {
	var err error
	var shaderName gl.Shader
	switch shader_type {
	case gohome.VERTEX:
		shaderName, err = compileOpenGLESShader(s.name, gl.VERTEX_SHADER, src, s.program, (*s.gles))
	case gohome.FRAGMENT:
		shaderName, err = compileOpenGLESShader(s.name, gl.FRAGMENT_SHADER, src, s.program, (*s.gles))
	case gohome.GEOMETRY:
		err = &OpenGLESError{errorString: "Geometry shader is not supported by OpenGLES"}
		// shaderName, err = compileOpenGLESShader(gl.GEOMETRY_SHADER, src, s.program, (*s.gles))
	case gohome.TESSELLETION:
		err = &OpenGLESError{errorString: "Tesselletion shader is not supported by OpenGLES"}
		// shaderName, err = compileOpenGLESShader(gl.TESS_CONTROL_SHADER, src, s.program, (*s.gles))
	case gohome.EVELUATION:
		err = &OpenGLESError{errorString: "Eveluation shader is not supported by OpenGLES"}
		// shaderName, err = compileOpenGLESShader(gl.TESS_EVALUATION_SHADER, src, s.program, (*s.gles))
	case gohome.COMPUTE:
		err = &OpenGLESError{errorString: "Compute shader is not supported by OpenGLES"}
		// shaderName, err = compileOpenGLESShader(gl.COMPUTE_SHADER, src, s.program, (*s.gles))
	}

	if err != nil {
		return &OpenGLESError{errorString: "Couldn't compile source of " + getShaderTypeName(shader_type) + " of " + s.name + ": " + err.Error()}
	}

	s.shaders[shader_type] = shaderName
	return nil
}

func (s *OpenGLESShader) deleteAllShaders() {
	for i := 0; i < 6; i++ {
		if s.shaders[i].Value != 0 {
			(*s.gles).DetachShader(s.program, s.shaders[i])
			(*s.gles).DeleteShader(s.shaders[i])
		}
	}
}

func (s *OpenGLESShader) Link() error {
	defer s.deleteAllShaders()

	(*s.gles).LinkProgram(s.program)

	status := (*s.gles).GetProgrami(s.program, gl.LINK_STATUS)
	if status == gl.FALSE {
		logText := (*s.gles).GetProgramInfoLog(s.program)

		return &OpenGLESError{errorString: "Couldn't link shader " + s.name + ": " + logText}
	}
	return nil
}

func (s *OpenGLESShader) Use() {
	(*s.gles).UseProgram(s.program)
}

func (s *OpenGLESShader) Unuse() {
	(*s.gles).UseProgram(gl.Program{false, 0})
}

func (s *OpenGLESShader) Setup() error {
	s.validate()
	return nil
}

func (s *OpenGLESShader) Terminate() {
	(*s.gles).DeleteProgram(s.program)
}

func (s *OpenGLESShader) getUniformLocation(name string) gl.Uniform {
	var loc gl.Uniform
	var ok bool
	if loc, ok = s.uniform_locations[name]; !ok {
		loc = (*s.gles).GetUniformLocation(s.program, name)
		s.uniform_locations[name] = loc
	}
	if !ok && loc.Value == -1 {
		gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_WARNING, "Shader", s.name, "Couldn't find uniform "+name)
	}
	return loc
}

func (s *OpenGLESShader) SetUniformV2(name string, value mgl32.Vec2) {
	loc := s.getUniformLocation(name)
	if loc.Value != -1 {
		(*s.gles).Uniform2f(loc, value[0], value[1])
	}
}
func (s *OpenGLESShader) SetUniformV3(name string, value mgl32.Vec3) {
	loc := s.getUniformLocation(name)
	if loc.Value != -1 {
		(*s.gles).Uniform3f(loc, value[0], value[1], value[2])
	}
}
func (s *OpenGLESShader) SetUniformV4(name string, value mgl32.Vec4) {
	loc := s.getUniformLocation(name)
	if loc.Value != -1 {
		(*s.gles).Uniform4f(loc, value[0], value[1], value[2], value[3])
	}
}
func (s *OpenGLESShader) SetUniformIV2(name string, value []int32) {
	loc := s.getUniformLocation(name)
	if loc.Value != -1 {
		(*s.gles).Uniform2i(loc, int(value[0]), int(value[1]))
	}
}
func (s *OpenGLESShader) SetUniformIV3(name string, value []int32) {
	loc := s.getUniformLocation(name)
	if loc.Value != -1 {
		(*s.gles).Uniform3i(loc, value[0], value[1], value[2])
	}
}
func (s *OpenGLESShader) SetUniformIV4(name string, value []int32) {
	loc := s.getUniformLocation(name)
	if loc.Value != -1 {
		(*s.gles).Uniform4i(loc, value[0], value[1], value[2], value[3])
	}
}
func (s *OpenGLESShader) SetUniformF(name string, value float32) {
	loc := s.getUniformLocation(name)
	if loc.Value != -1 {
		(*s.gles).Uniform1f(loc, value)
	}
}
func (s *OpenGLESShader) SetUniformI(name string, value int32) {
	loc := s.getUniformLocation(name)
	if loc.Value != -1 {
		(*s.gles).Uniform1i(loc, int(value))
	}
}
func (s *OpenGLESShader) SetUniformUI(name string, value uint32) {
	s.SetUniformI(name, int32(value))
}
func (s *OpenGLESShader) SetUniformB(name string, value uint8) {
	s.SetUniformI(name, int32(value))
}
func (s *OpenGLESShader) SetUniformM2(name string, value mgl32.Mat2) {
	loc := s.getUniformLocation(name)
	if loc.Value != -1 {
		(*s.gles).UniformMatrix2fv(loc, value[:])
	}
}
func (s *OpenGLESShader) SetUniformM3(name string, value mgl32.Mat3) {
	loc := s.getUniformLocation(name)
	if loc.Value != -1 {
		(*s.gles).UniformMatrix3fv(loc, value[:])
	}
}
func (s *OpenGLESShader) SetUniformM4(name string, value mgl32.Mat4) {
	loc := s.getUniformLocation(name)
	if loc.Value != -1 {
		(*s.gles).UniformMatrix4fv(loc, value[:])
	}
}

func (s *OpenGLESShader) SetUniformMaterial(mat gohome.Material) {
	var diffBind int32 = 0
	var specBind int32 = 0
	var normBind int32 = 0
	var boundTextures uint32

	maxtextures := gohome.Render.GetMaxTextures()

	if mat.DiffuseTexture != nil {
		diffBind = int32(gohome.Render.NextTextureUnit())
		if diffBind >= maxtextures-1 {
			diffBind = 0
			mat.DiffuseTextureLoaded = 0
			gohome.Render.DecrementTextureUnit(1)
		} else {
			mat.DiffuseTexture.Bind(uint32(diffBind))
			mat.DiffuseTextureLoaded = 1
			boundTextures++
		}
	} else {
		mat.DiffuseTextureLoaded = 0
	}

	if mat.SpecularTexture != nil {
		specBind = int32(gohome.Render.NextTextureUnit())
		if specBind >= maxtextures-1 {
			specBind = 0
			mat.SpecularTextureLoaded = 0
			gohome.Render.DecrementTextureUnit(1)
		} else {
			mat.SpecularTexture.Bind(uint32(specBind))
			mat.SpecularTextureLoaded = 1
			boundTextures++
		}
	} else {
		mat.SpecularTextureLoaded = 0
	}

	if mat.NormalMap != nil {
		normBind = int32(gohome.Render.NextTextureUnit())
		if normBind >= maxtextures-1 {
			normBind = 0
			mat.NormalMapLoaded = 0
			gohome.Render.DecrementTextureUnit(1)
		} else {
			mat.NormalMap.Bind(uint32(normBind))
			mat.NormalMapLoaded = 1
			boundTextures++
		}
	} else {
		mat.NormalMapLoaded = 0
	}

	gohome.Render.DecrementTextureUnit(boundTextures)

	s.SetUniformV3(gohome.MATERIAL_UNIFORM_NAME+"."+gohome.MATERIAL_DIFFUSE_COLOR_UNIFORM_NAME, gohome.ColorToVec3(mat.DiffuseColor))
	s.SetUniformV3(gohome.MATERIAL_UNIFORM_NAME+"."+gohome.MATERIAL_SPECULAR_COLOR_UNIFORM_NAME, gohome.ColorToVec3(mat.SpecularColor))
	s.SetUniformF(gohome.MATERIAL_UNIFORM_NAME+"."+gohome.MATERIAL_SHINYNESS_UNIFORM_NAME, mat.Shinyness)

	s.SetUniformB(gohome.MATERIAL_UNIFORM_NAME+"."+gohome.MATERIAL_DIFFUSE_TEXTURE_LOADED_UNIFORM_NAME, mat.DiffuseTextureLoaded)
	s.SetUniformB(gohome.MATERIAL_UNIFORM_NAME+"."+gohome.MATERIAL_SPECULAR_TEXTURE_LOADED_UNIFORM_NAME, mat.SpecularTextureLoaded)
	s.SetUniformB(gohome.MATERIAL_UNIFORM_NAME+"."+gohome.MATERIAL_NORMALMAP_LOADED_UNIFORM_NAME, mat.NormalMapLoaded)
	s.SetUniformI(gohome.MATERIAL_UNIFORM_NAME+gohome.MATERIAL_DIFFUSE_TEXTURE_UNIFORM_NAME, diffBind)
	s.SetUniformI(gohome.MATERIAL_UNIFORM_NAME+gohome.MATERIAL_SPECULAR_TEXTURE_UNIFORM_NAME, specBind)
	s.SetUniformI(gohome.MATERIAL_UNIFORM_NAME+gohome.MATERIAL_NORMALMAP_UNIFORM_NAME, normBind)
}

func (s *OpenGLESShader) SetUniformLights(lightCollectionIndex int32) {
	if lightCollectionIndex == -1 || lightCollectionIndex > int32(len(gohome.LightMgr.LightCollections)-1) {
		s.SetUniformI(gohome.NUM_POINT_LIGHTS_UNIFORM_NAME, 0)
		s.SetUniformI(gohome.NUM_DIRECTIONAL_LIGHTS_UNIFORM_NAME, 0)
		s.SetUniformI(gohome.NUM_SPOT_LIGHTS_UNIFORM_NAME, 0)

		s.SetUniformV3(gohome.AMBIENT_LIGHT_UNIFORM_NAME, mgl32.Vec3{1.0, 1.0, 1.0})
	}

	lightColl := gohome.LightMgr.LightCollections[lightCollectionIndex]

	s.SetUniformI(gohome.NUM_POINT_LIGHTS_UNIFORM_NAME, int32(len(lightColl.PointLights)))
	s.SetUniformI(gohome.NUM_DIRECTIONAL_LIGHTS_UNIFORM_NAME, int32(len(lightColl.DirectionalLights)))
	s.SetUniformI(gohome.NUM_SPOT_LIGHTS_UNIFORM_NAME, int32(len(lightColl.SpotLights)))

	s.SetUniformV3(gohome.AMBIENT_LIGHT_UNIFORM_NAME, gohome.ColorToVec3(lightColl.AmbientLight))

	var i uint32
	for i = 0; i < uint32(len(lightColl.PointLights)); i++ {
		lightColl.PointLights[i].SetUniforms(s, i)
	}
	for i = 0; i < uint32(len(lightColl.DirectionalLights)); i++ {
		lightColl.DirectionalLights[i].SetUniforms(s, i)
	}
	for i = 0; i < uint32(len(lightColl.SpotLights)); i++ {
		lightColl.SpotLights[i].SetUniforms(s, i)
	}
}

func (s *OpenGLESShader) GetName() string {
	return s.name
}

func (s *OpenGLESShader) validate() bool {
	if s.validated {
		return true
	}
	s.Use()
	maxtextures := gohome.Render.GetMaxTextures()
	for i := 0; i < 31; i++ {
		s.SetUniformI("pointLights["+strconv.Itoa(i)+"].shadowmap", maxtextures-1)
	}
	s.Unuse()
	s.validated = true
	(*s.gles).ValidateProgram(s.program)
	status := (*s.gles).GetProgrami(s.program, gl.VALIDATE_STATUS)
	if status == gl.FALSE {
		logText := (*s.gles).GetProgramInfoLog(s.program)
		s.validated = false
		gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_ERROR, "Shader", s.name, "Couldn't validate: "+logText)
		return false
	}

	return true
}

func (s *OpenGLESShader) AddAttribute(name string, location uint32) {
	(*s.gles).BindAttribLocation(s.program, gl.Attrib{uint(location)}, name)
}
