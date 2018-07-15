package gohome

import (
	"strconv"
	"sync"
)

type preloadedTexture struct {
	Name                 string
	Path                 string
	fileAlreadyPreloaded bool
}

type preloadedShader struct {
	Name                      string
	VertexShader              string
	FragmentShader            string
	GeometryShader            string
	TesselletionControlShader string
	EveluationShader          string
	ComputeShader             string
}

type preloadedLevel struct {
	Name                 string
	Path                 string
	LoadToGPU            bool
	fileAlreadyPreloaded bool
}

type preloadedTextureData struct {
	Tex      Texture
	img_data []byte
	width    int
	height   int
	path     string
}

type preloadedFont struct {
	Name string
	Path string
}

type preloadedFontData struct {
	preloadedFont
	font Font
}

type preloadedSound struct {
	Name string
	Path string
}

type preloadedSoundData struct {
	preloadedSound
	sound Sound
}

type preloadedMusic struct {
	Name string
	Path string
}

type preloadedMusicData struct {
	preloadedMusic
	music Music
}

type preloadedLevelData struct {
	Lvl  *Level
	path string
}

type preloadedShaderData struct {
	name     string
	contents [6]string
}

type preloadedLevelObject struct {
	Lvl    *Level
	Lvlobj LevelObject
}

type PreloadedMesh struct {
	Mesh      Mesh3D
	LoadToGPU bool
}

type alreadyPreloadedResource struct {
	Path string
	Name string
}

type preloader struct {
	preloadedTextures []preloadedTexture
	preloadedShaders  []preloadedShader
	preloadedLevels   []preloadedLevel
	preloadedFonts    []preloadedFont
	preloadedSounds	  []preloadedSound
	preloadedMusics	  []preloadedMusic

	preloadedShaderDataChan     chan preloadedShaderData
	preloadedLevelsChan         chan preloadedLevelData
	PreloadedModelsChan         chan *Model3D
	PreloadedMeshesChan         chan PreloadedMesh
	preloadedTextureDataChan    chan preloadedTextureData
	preloadedFontChan           chan preloadedFontData
	preloadedSoundChan			chan preloadedSoundData
	preloadedMusicChan			chan preloadedMusicData
	alreadyPreloadedTextureChan chan alreadyPreloadedResource
	alreadyPreloadedLevelChan   chan alreadyPreloadedResource
	exitChan                    chan bool
	exitLevelsChan              chan bool
	exitTexturesChan            chan bool
	exitShadersChan             chan bool
	exitFontsChan               chan bool
	exitSoundsChan				chan bool
	exitMusicsChan				chan bool

	preloadedTexturesToFinish []preloadedTextureData
	preloadedShadersToFinish  []preloadedShaderData
	PreloadedMeshesToFinish   []PreloadedMesh

	alreadyPreloadedTexturePathsToSet []alreadyPreloadedResource
	alreadyPreloadedLevelPathsToSet   []alreadyPreloadedResource
}

func (this *preloader) Init() {
}

func (this *preloader) loadPreloadedLevel(lvl *preloadedLevel, wg *sync.WaitGroup) {
	defer wg.Done()

	if !lvl.fileAlreadyPreloaded {
		name := lvl.Name
		path := lvl.Path

		level := ResourceMgr.loadLevel(name, path, true, lvl.LoadToGPU)
		if level != nil {
			preLevel := preloadedLevelData{
				level,
				lvl.Path,
			}
			this.preloadedLevelsChan <- preLevel
		}
	} else {
		this.alreadyPreloadedLevelChan <- alreadyPreloadedResource{
			lvl.Path,
			lvl.Name,
		}
	}
}

func (this *preloader) loadPreloadedLevels() {
	if len(this.preloadedLevels) == 0 {

	} else {
		var wg1 sync.WaitGroup
		wg1.Add(len(this.preloadedLevels))
		for i := 0; i < len(this.preloadedLevels); i++ {
			go this.loadPreloadedLevel(&this.preloadedLevels[i], &wg1)
		}
		wg1.Wait()
	}

	go func() {
		this.exitLevelsChan <- true
	}()

}

func (this *preloader) loadPreloadedShader(s *preloadedShader, wg *sync.WaitGroup) {
	defer wg.Done()
	name := s.Name
	vertex_path := s.VertexShader
	fragment_path := s.FragmentShader
	geometry_path := s.GeometryShader
	tesselletion_control_path := s.TesselletionControlShader
	eveluation_path := s.EveluationShader
	compute_path := s.ComputeShader
	ResourceMgr.loadShader(name, vertex_path, fragment_path, geometry_path, tesselletion_control_path, eveluation_path, compute_path, true)
}

func (this *preloader) loadPreloadedShaders() {
	if len(this.preloadedShaders) == 0 {

	} else {
		var wg1 sync.WaitGroup
		wg1.Add(len(this.preloadedShaders))
		for i := 0; i < len(this.preloadedShaders); i++ {
			go this.loadPreloadedShader(&this.preloadedShaders[i], &wg1)
		}
		wg1.Wait()
	}

	go func() {
		this.exitShadersChan <- true
	}()
}

func (this *preloader) loadPreloadedTexture(tex *preloadedTexture, wg *sync.WaitGroup) {
	defer wg.Done()

	if !tex.fileAlreadyPreloaded {
		name := tex.Name
		path := tex.Path

		ResourceMgr.LoadTextureFunction(name, path, true)
	} else {
		this.alreadyPreloadedTextureChan <- alreadyPreloadedResource{
			tex.Path,
			tex.Name,
		}
	}
}

func (this *preloader) loadPreloadedTextures() {
	if len(this.preloadedTextures) == 0 {

	} else {
		var wg1 sync.WaitGroup
		wg1.Add(len(this.preloadedTextures))
		for i := 0; i < len(this.preloadedTextures); i++ {
			go this.loadPreloadedTexture(&this.preloadedTextures[i], &wg1)
		}
		wg1.Wait()
	}

	go func() {
		this.exitTexturesChan <- true
	}()
}

func (this *preloader) loadPreloadedFonts() {
	if len(this.preloadedFonts) == 0 {

	} else {
		var wg1 sync.WaitGroup
		wg1.Add(len(this.preloadedFonts))
		for i := 0; i < len(this.preloadedFonts); i++ {
			go this.loadPreloadedFont(this.preloadedFonts[i], &wg1)
		}
		wg1.Wait()
	}

	go func() {
		this.exitFontsChan <- true
	}()
}

func (this *preloader) loadPreloadedFont(f preloadedFont, wg *sync.WaitGroup) {
	defer wg.Done()

	ResourceMgr.loadFont(f.Name, f.Path, true)
}

func (this *preloader) loadPreloadedSounds() {
	if len(this.preloadedSounds) == 0 {

	} else {
		var wg1 sync.WaitGroup
		wg1.Add(len(this.preloadedSounds))
		for i := 0; i < len(this.preloadedSounds); i++ {
			go this.loadPreloadedSound(this.preloadedSounds[i], &wg1)
		}
		wg1.Wait()
	}

	go func() {
		this.exitSoundsChan <- true
	}()
}

func (this *preloader) loadPreloadedSound(f preloadedSound, wg *sync.WaitGroup) {
	defer wg.Done()

	ResourceMgr.loadSound(f.Name, f.Path, true)
}

func (this *preloader) loadPreloadedMusics() {
	if len(this.preloadedMusics) == 0 {

	} else {
		var wg1 sync.WaitGroup
		wg1.Add(len(this.preloadedMusics))
		for i := 0; i < len(this.preloadedMusics); i++ {
			go this.loadPreloadedMusic(this.preloadedMusics[i], &wg1)
		}
		wg1.Wait()
	}

	go func() {
		this.exitMusicsChan <- true
	}()
}

func (this *preloader) loadPreloadedMusic(f preloadedMusic, wg *sync.WaitGroup) {
	defer wg.Done()

	ResourceMgr.loadMusic(f.Name, f.Path, true)
}

func (this *preloader) finish(wg *sync.WaitGroup) {
	defer wg.Done()

	var done bool = false

	for true {
		select {
		case preLvl := <-this.preloadedLevelsChan:
			ResourceMgr.Levels[preLvl.Lvl.Name] = preLvl.Lvl
			ResourceMgr.resourceFileNames[preLvl.path] = preLvl.Lvl.Name
			ErrorMgr.Message(ERROR_LEVEL_LOG, "Level", preLvl.Lvl.Name, "Finished loading!")
		case tex := <-this.preloadedTextureDataChan:
			this.preloadedTexturesToFinish = append(this.preloadedTexturesToFinish, tex)
		case shader := <-this.preloadedShaderDataChan:
			this.preloadedShadersToFinish = append(this.preloadedShadersToFinish, shader)
		case mesh := <-this.PreloadedMeshesChan:
			this.PreloadedMeshesToFinish = append(this.PreloadedMeshesToFinish, mesh)
		case tex := <-this.alreadyPreloadedTextureChan:
			this.alreadyPreloadedTexturePathsToSet = append(this.alreadyPreloadedTexturePathsToSet, tex)
		case lvl := <-this.alreadyPreloadedLevelChan:
			this.alreadyPreloadedLevelPathsToSet = append(this.alreadyPreloadedLevelPathsToSet, lvl)
		case model := <-this.PreloadedModelsChan:
			ResourceMgr.Models[model.Name] = model
			ErrorMgr.Message(ERROR_LEVEL_LOG, "Model", model.Name, "Finished loading!")
		case font := <-this.preloadedFontChan:
			ResourceMgr.fonts[font.Name] = &font.font
			ResourceMgr.resourceFileNames[font.Path] = font.Name
			ErrorMgr.Log("Font", font.Name, "Finished Loading!")
		case sound := <-this.preloadedSoundChan:
			ResourceMgr.sounds[sound.Name] = sound.sound
			ResourceMgr.resourceFileNames[sound.Path] = sound.Name
			ErrorMgr.Log("Sound",sound.Name,"Finished Loading!")
		case music := <-this.preloadedMusicChan:
			ResourceMgr.musics[music.Name] = music.music
			ResourceMgr.resourceFileNames[music.Path] = music.Name
			ErrorMgr.Log("Music",music.Name,"Finished Loading!")
		case <-this.exitChan:
			done = true
		default:
		}
		if done {
			break
		}
	}
}

func (this *preloader) checkExit(wg *sync.WaitGroup) {
	defer wg.Done()

	var LevelsExit, texturesExit, shadersExit, fontsExit, soundsExit, musicsExit bool
	var done bool

	for true {
		select {
		case <-this.exitLevelsChan:
			LevelsExit = true
		case <-this.exitTexturesChan:
			texturesExit = true
		case <-this.exitShadersChan:
			shadersExit = true
		case <-this.exitFontsChan:
			fontsExit = true
		case <-this.exitSoundsChan:
			soundsExit = true
		case <-this.exitMusicsChan:
			musicsExit = true
		default:
			if LevelsExit && texturesExit && shadersExit && fontsExit && soundsExit && musicsExit {
				this.exitChan <- true
				done = true
			}
		}
		if done {
			break
		}
	}

}

func (this *preloader) finishTextures() {
	for i := 0; i < len(this.preloadedTexturesToFinish); i++ {
		tex := this.preloadedTexturesToFinish[i]
		tex.Tex.Load(tex.img_data, tex.width, tex.height, false)
		ResourceMgr.textures[tex.Tex.GetName()] = tex.Tex
		ResourceMgr.resourceFileNames[tex.path] = tex.Tex.GetName()
		ErrorMgr.Message(ERROR_LEVEL_LOG, "Texture", tex.Tex.GetName(), "Finished loading! W: "+strconv.Itoa(tex.width)+" H: "+strconv.Itoa(tex.height))
	}

	for i := 0; i < len(this.alreadyPreloadedTexturePathsToSet); i++ {
		ResourceMgr.textures[this.alreadyPreloadedTexturePathsToSet[i].Name] = ResourceMgr.textures[ResourceMgr.resourceFileNames[this.alreadyPreloadedTexturePathsToSet[i].Path]]
	}
}

func (this *preloader) finishShaders() {
	for i := 0; i < len(this.preloadedShadersToFinish); i++ {
		shader := this.preloadedShadersToFinish[i]
		s, err := Render.LoadShader(shader.name, shader.contents[VERTEX], shader.contents[FRAGMENT], shader.contents[GEOMETRY], shader.contents[TESSELLETION], shader.contents[EVELUATION], shader.contents[COMPUTE])
		if s != nil {
			ResourceMgr.shaders[shader.name] = s
			ErrorMgr.Message(ERROR_LEVEL_LOG, "Shader", s.GetName(), "Finished loading!")
		} else {
			ErrorMgr.MessageError(ERROR_LEVEL_ERROR, "Shader", s.GetName(), err)
		}
	}
}

func (this *preloader) finishMeshes() {
	for i := 0; i < len(this.PreloadedMeshesToFinish); i++ {
		mesh := this.PreloadedMeshesToFinish[i]
		if mesh.LoadToGPU {
			mesh.Mesh.Load()
		}
		ErrorMgr.Message(ERROR_LEVEL_LOG, "Mesh", mesh.Mesh.GetName(), "Finished loading! V: "+strconv.Itoa(int(mesh.Mesh.GetNumVertices()))+" I: "+strconv.Itoa(int(mesh.Mesh.GetNumIndices())))
	}
}

func (this *preloader) finishLevels() {
	for i := 0; i < len(this.alreadyPreloadedLevelPathsToSet); i++ {
		ResourceMgr.Levels[this.alreadyPreloadedLevelPathsToSet[i].Name] = ResourceMgr.Levels[ResourceMgr.resourceFileNames[this.alreadyPreloadedLevelPathsToSet[i].Path]]
	}
}

func (this *preloader) finishData() {
	this.finishLevels()
	this.finishTextures()
	this.finishShaders()
	this.finishMeshes()
}

func (this *preloader) openChannels() {
	this.preloadedShaderDataChan = make(chan preloadedShaderData)
	this.preloadedLevelsChan = make(chan preloadedLevelData)
	this.PreloadedModelsChan = make(chan *Model3D)
	this.PreloadedMeshesChan = make(chan PreloadedMesh)
	this.alreadyPreloadedTextureChan = make(chan alreadyPreloadedResource)
	this.alreadyPreloadedLevelChan = make(chan alreadyPreloadedResource)
	this.preloadedTextureDataChan = make(chan preloadedTextureData)
	this.preloadedFontChan = make(chan preloadedFontData)
	this.preloadedSoundChan = make(chan preloadedSoundData)
	this.preloadedMusicChan = make(chan preloadedMusicData)
	this.exitChan = make(chan bool)
	this.exitLevelsChan = make(chan bool)
	this.exitTexturesChan = make(chan bool)
	this.exitShadersChan = make(chan bool)
	this.exitFontsChan = make(chan bool)
	this.exitSoundsChan = make(chan bool)
	this.exitMusicsChan = make(chan bool)
}

func (this *preloader) closeChannels() {
	close(this.preloadedShaderDataChan)
	close(this.preloadedLevelsChan)
	close(this.PreloadedModelsChan)
	close(this.PreloadedMeshesChan)
	close(this.preloadedTextureDataChan)
	close(this.alreadyPreloadedTextureChan)
	close(this.alreadyPreloadedLevelChan)
	close(this.preloadedFontChan)
	close(this.preloadedSoundChan)
	close(this.preloadedMusicChan)
	close(this.exitChan)
	close(this.exitLevelsChan)
	close(this.exitTexturesChan)
	close(this.exitShadersChan)
	close(this.exitFontsChan)
}

func (this *preloader) clearSlices() {
	this.preloadedTextures = this.preloadedTextures[:0]
	this.preloadedShaders = this.preloadedShaders[:0]
	this.preloadedLevels = this.preloadedLevels[:0]
	this.preloadedTexturesToFinish = this.preloadedTexturesToFinish[:0]
	this.preloadedShadersToFinish = this.preloadedShadersToFinish[:0]
	this.PreloadedMeshesToFinish = this.PreloadedMeshesToFinish[:0]
	this.alreadyPreloadedTexturePathsToSet = this.alreadyPreloadedTexturePathsToSet[:0]
	this.alreadyPreloadedLevelPathsToSet = this.alreadyPreloadedLevelPathsToSet[:0]
	this.preloadedFonts = this.preloadedFonts[:0]
	this.preloadedSounds = this.preloadedSounds[:0]
	this.preloadedMusics = this.preloadedMusics[:0]
}

func (this *preloader) loadPreloadedResources() {
	this.openChannels()
	var wg sync.WaitGroup
	wg.Add(2)

	go this.checkExit(&wg)
	go this.finish(&wg)

	go this.loadPreloadedLevels()
	go this.loadPreloadedShaders()
	go this.loadPreloadedTextures()
	go this.loadPreloadedFonts()
	go this.loadPreloadedSounds()
	go this.loadPreloadedMusics()

	wg.Wait()

	this.closeChannels()
	this.finishData()
	this.clearSlices()
}
