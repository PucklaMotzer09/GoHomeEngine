package gohome

import (
	"github.com/go-gl/mathgl/mgl32"
)

type TweenPosition2D struct {
	Destination mgl32.Vec2
	Time float32
	TweenType uint8

	transform *TransformableObject2D
	velocity mgl32.Vec2
	elapsedTime float32
}

func (this *TweenPosition2D) Start(parent interface{}) {
	parent2D,ok := parent.(TweenableObject2D)
	if ok {
		this.transform = parent2D.GetTransform2D()
	} else {
		this.transform = nil
	}

	if this.transform != nil {
		this.velocity = this.Destination.Sub(this.transform.Position).Mul(1.0/this.Time)
	}
	this.elapsedTime = 0.0
}
func (this *TweenPosition2D) Update(delta_time float32) bool {
	if this.transform == nil {
		return true
	}
	this.elapsedTime += delta_time

	this.transform.Position = this.transform.Position.Add(this.velocity.Mul(delta_time))

	if this.elapsedTime >= this.Time {
		return true
	}

	return false
}

func (this *TweenPosition2D) GetType() uint8 {
	return this.TweenType
}
func (this *TweenPosition2D) End() {
	if this.transform != nil {
		this.transform.Position = this.Destination
	}
}
func (this *TweenPosition2D) Reset() {
	if this.transform != nil {
		this.transform.Position = this.Destination.Sub(this.velocity.Mul(this.Time))
		this.elapsedTime = 0.0
	}
}

type TweenRotation2D struct {
	Destination float32
	Time float32
	TweenType uint8

	transform *TransformableObject2D
	velocity float32
	elapsedTime float32
}

func (this *TweenRotation2D) Start(parent interface{}) {
	parent2D,ok := parent.(TweenableObject2D)
	if ok {
		this.transform = parent2D.GetTransform2D()
	} else {
		this.transform = nil
	}

	if this.transform != nil {
		this.velocity = (this.Destination - this.transform.Rotation)  / this.Time
	}
	this.elapsedTime = 0.0
}
func (this *TweenRotation2D) Update(delta_time float32) bool {
	if this.transform == nil {
		return true
	}
	this.elapsedTime += delta_time

	this.transform.Rotation += this.velocity * delta_time

	if this.elapsedTime >= this.Time {
		return true
	}

	return false
}

func (this *TweenRotation2D) GetType() uint8 {
	return this.TweenType
}
func (this *TweenRotation2D) End() {
	if this.transform != nil {
		this.transform.Rotation = this.Destination
	}
}
func (this *TweenRotation2D) Reset() {
	if this.transform != nil {
		this.transform.Rotation = this.Destination - this.velocity * this.Time
		this.elapsedTime = 0.0
	}
}

type TweenWait struct {
	Time float32
	TweenType uint8

	elapsedTime float32
}

func (this *TweenWait) Start(parent interface{}) {
	this.elapsedTime = 0.0
}
func (this *TweenWait) Update(delta_time float32) bool {
	this.elapsedTime += delta_time

	if this.elapsedTime >= this.Time {
		return true
	}

	return false
}

func (this *TweenWait) GetType() uint8 {
	return this.TweenType
}
func (this *TweenWait) End() {

}
func (this *TweenWait) Reset() {
	this.elapsedTime = 0.0
}

type BlinkableObject interface {
	SetVisible()
	SetInvisible()
	IsVisible() bool
}

type TweenBlink struct {
	Amount uint32
	Time float32
	TweenType uint8

	timeForOneBlink float32
	elapsedTime float32
	oneBlinkElapsedTime float32
	previousVisible bool
	parent BlinkableObject
}

func (this *TweenBlink) Start(parent interface{}) {
	this.elapsedTime = 0.0
	this.timeForOneBlink = this.Time/float32(this.Amount)
	if parent != nil {
		this.parent = parent.(BlinkableObject)
		if this.parent != nil {
			this.previousVisible = this.parent.IsVisible()
		}
	}
}
func (this *TweenBlink) Update(delta_time float32) bool {
	if this.parent == nil {
		return true
	}

	this.elapsedTime += delta_time
	this.oneBlinkElapsedTime += delta_time

	if this.oneBlinkElapsedTime >= this.timeForOneBlink/2.0 {
		if this.parent.IsVisible() {
			this.parent.SetInvisible()
		} else {
			this.parent.SetVisible()
		}
		this.oneBlinkElapsedTime = 0.0
	}

	if this.elapsedTime >= this.Time {
		return true
	}

	return false
}

func (this *TweenBlink) GetType() uint8 {
	return this.TweenType
}
func (this *TweenBlink) End() {
	if this.parent != nil {
		if this.previousVisible {
			this.parent.SetVisible()
		} else {
			this.parent.SetInvisible()
		}
	}
}
func (this *TweenBlink) Reset() {
	this.elapsedTime = 0.0
	this.oneBlinkElapsedTime = 0.0
	if this.parent != nil {
		if this.previousVisible {
			this.parent.SetVisible()
		} else {
			this.parent.SetInvisible()
		}
	}
}

type TweenScale2D struct {
	Destination mgl32.Vec2
	Time float32
	TweenType uint8

	elapsedTime float32
	velocity mgl32.Vec2
	transform *TransformableObject2D
}

func (this *TweenScale2D) Start(parent interface{}) {
	this.elapsedTime = 0.0
	this.transform = nil

	if parent != nil {
		parent2D,ok := parent.(TweenableObject2D)
		if ok {
			this.transform = parent2D.GetTransform2D()
		} else {
			return
		}
	} else {
		return
	}

	this.velocity = this.Destination.Sub(this.transform.Scale).Mul(1.0/this.Time)
}

func (this *TweenScale2D) Update(delta_time float32) bool {
	if this.transform == nil {
		return true
	}

	this.transform.Scale = this.transform.Scale.Add(this.velocity.Mul(delta_time))

	this.elapsedTime += delta_time
	if this.elapsedTime >= this.Time {
		return true
	}

	return false
}

func (this *TweenScale2D) End() {
	if this.transform == nil {
		return
	}

	this.transform.Scale = this.Destination
	this.elapsedTime = 0.0
}

func (this *TweenScale2D) GetType() uint8 {
	return this.TweenType
}

func (this *TweenScale2D) Reset() {
	if this.transform == nil {
		return
	}
	this.elapsedTime = 0.0
	this.transform.Scale = this.Destination.Sub(this.velocity.Mul(this.Time))
}

type TweenRegion2D struct {
	Destination TextureRegion
	Time float32
	TweenType uint8

	startRegion  TextureRegion
	startSize	 mgl32.Vec2
	parent *Sprite2D
	elapsedTime float32
}

func (this *TweenRegion2D) Start(parent interface{}) {
	this.elapsedTime = 0.0

	if parent != nil {
		this.parent = parent.(*Sprite2D)
		if this.parent != nil {
			this.startRegion = this.parent.TextureRegion
			this.parent.TextureRegion = this.Destination
			if this.parent.Transform != nil {
				this.startSize = this.parent.Transform.Size
				this.parent.Transform.Size = [2]float32{this.Destination.Width(),this.Destination.Height()}
			}
		}
	}
}

func (this *TweenRegion2D) Update(delta_time float32) bool {
	if this.parent == nil {
		return true
	}

	this.elapsedTime += delta_time
	if this.elapsedTime >= this.Time {
		return true
	}

	return false
}

func (this *TweenRegion2D) End() {
	if this.parent == nil {
		return
	}

	this.elapsedTime = 0.0
}

func (this *TweenRegion2D) GetType() uint8 {
	return this.TweenType
}

func (this *TweenRegion2D) Reset() {
	if this.parent == nil {
		return
	}
	this.elapsedTime = 0.0
	this.parent.TextureRegion = this.startRegion
	if this.parent.Transform != nil {
	   this.parent.Transform.Size = this.startSize
	}
}

type TweenTexture2D struct {
	Destination Texture
	Time float32
	TweenType uint8

	elapsedTime float32
	parent *Sprite2D
	startTexture Texture
	startSize mgl32.Vec2
}

func (this *TweenTexture2D) Start(parent interface{}) {
	this.elapsedTime = 0.0

	if parent != nil {
		this.parent = parent.(*Sprite2D)
		if this.parent != nil {
			this.startTexture = this.parent.Texture
			this.parent.Texture = this.Destination
			if this.parent.Transform != nil {
				this.startSize = this.parent.Transform.Size
				this.parent.Transform.Size = [2]float32{float32(this.Destination.GetWidth()),float32(this.Destination.GetHeight())}
			}
		}
	}
}

func (this *TweenTexture2D) Update(delta_time float32) bool {
	if this.parent == nil {
		return true
	}

	this.elapsedTime += delta_time
	if this.elapsedTime >= this.Time {
		return true
	}

	return false
}

func (this *TweenTexture2D) End() {
	if this.parent == nil {
		return
	}

	this.elapsedTime = 0.0
}

func (this *TweenTexture2D) GetType() uint8 {
	return this.TweenType
}

func (this *TweenTexture2D) Reset() {
	if this.parent == nil {
		return
	}
	this.elapsedTime = 0.0
	this.parent.Texture = this.startTexture
	if this.parent.Transform != nil {
		this.parent.Transform.Size = this.startSize
	}
}

type TweenPosition3D struct {
	Destination mgl32.Vec3
	Time float32
	TweenType uint8

	transform *TransformableObject3D
	velocity mgl32.Vec3
	elapsedTime float32
}

func (this *TweenPosition3D) Start(parent interface{}) {
	if parent != nil {
		parent3D,ok := parent.(TweenableObject3D)
		if ok {
			this.transform = parent3D.GetTransform3D()
		} else {
			this.transform = nil
		}
	}


	if this.transform != nil {
		this.velocity = this.Destination.Sub(this.transform.Position).Mul(1.0/this.Time)
	}
	this.elapsedTime = 0.0
}
func (this *TweenPosition3D) Update(delta_time float32) bool {
	if this.transform == nil {
		return true
	}
	this.elapsedTime += delta_time

	this.transform.Position = this.transform.Position.Add(this.velocity.Mul(delta_time))

	if this.elapsedTime >= this.Time {
		return true
	}

	return false
}

func (this *TweenPosition3D) GetType() uint8 {
	return this.TweenType
}
func (this *TweenPosition3D) End() {
	if this.transform != nil {
		this.transform.Position = this.Destination
	}
}
func (this *TweenPosition3D) Reset() {
	if this.transform != nil {
		this.transform.Position = this.Destination.Sub(this.velocity.Mul(this.Time))
		this.elapsedTime = 0.0
	}
}

type TweenRotation3D struct {
	Destination mgl32.Quat
	Time float32
	TweenType uint8

	transform *TransformableObject3D
	start mgl32.Quat
	elapsedTime float32
}

func (this *TweenRotation3D) Start(parent interface{}) {
	parent3D,ok := parent.(TweenableObject3D)
	if ok {
		this.transform = parent3D.GetTransform3D()
	} else {
		this.transform = nil
	}

	this.elapsedTime = 0.0
	if this.transform != nil {
		this.start = this.transform.Rotation
	}
}
func (this *TweenRotation3D) Update(delta_time float32) bool {
	if this.transform == nil {
		return true
	}
	this.elapsedTime += delta_time

	this.transform.Rotation = mgl32.QuatSlerp(this.start,this.Destination,this.elapsedTime/this.Time)

	if this.elapsedTime >= this.Time {
		return true
	}

	return false
}

func (this *TweenRotation3D) GetType() uint8 {
	return this.TweenType
}
func (this *TweenRotation3D) End() {
	if this.transform != nil {
		this.transform.Rotation = this.Destination
	}
}
func (this *TweenRotation3D) Reset() {
	if this.transform != nil {
		this.transform.Rotation = this.start
	}
	this.elapsedTime = 0.0
}

type TweenScale3D struct {
	Destination mgl32.Vec3
	Time float32
	TweenType uint8

	elapsedTime float32
	velocity mgl32.Vec3
	transform *TransformableObject3D
}

func (this *TweenScale3D) Start(parent interface{}) {
	this.elapsedTime = 0.0
	this.transform = nil

	if parent != nil {
		parent3D,ok := parent.(TweenableObject3D)
		if ok {
			this.transform = parent3D.GetTransform3D()
		} else {
			return
		}
	} else {
		return
	}

	this.velocity = this.Destination.Sub(this.transform.Scale).Mul(1.0/this.Time)
}

func (this *TweenScale3D) Update(delta_time float32) bool {
	if this.transform == nil {
		return true
	}

	this.transform.Scale = this.transform.Scale.Add(this.velocity.Mul(delta_time))

	this.elapsedTime += delta_time
	if this.elapsedTime >= this.Time {
		return true
	}

	return false
}

func (this *TweenScale3D) End() {
	if this.transform == nil {
		return
	}

	this.transform.Scale = this.Destination
	this.elapsedTime = 0.0
}

func (this *TweenScale3D) GetType() uint8 {
	return this.TweenType
}

func (this *TweenScale3D) Reset() {
	if this.transform == nil {
		return
	}
	this.elapsedTime = 0.0
	this.transform.Scale = this.Destination.Sub(this.velocity.Mul(this.Time))
}


