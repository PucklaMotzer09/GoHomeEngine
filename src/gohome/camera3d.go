package gohome

import (
	// "fmt"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	LOOK_DIRECTION_MAGNIFIER float32 = 100.0
)

type Camera3D struct {
	Position      mgl32.Vec3
	LookDirection mgl32.Vec3
	Up            mgl32.Vec3
	Tilt          float32
	rotation      mgl32.Vec2

	oldPosition      mgl32.Vec3
	oldLookDirection mgl32.Vec3
	oldTilt          float32
	oldUp            mgl32.Vec3
	oldrotation      mgl32.Vec2

	MaxRotation mgl32.Vec2
	MinRotation mgl32.Vec2

	viewMatrix        mgl32.Mat4
	inverseViewMatrix mgl32.Mat4
}

func (cam *Camera3D) Init() {
	cam.LookDirection = [3]float32{0.0, 0.0, -1.0}
	cam.MaxRotation = [2]float32{89.9, 370.0}
	cam.MinRotation = [2]float32{-89.9, -370.0}
	cam.Up = [3]float32{0.0, 1.0, 0.0}
	cam.CalculateViewMatrix()
}

func (cam *Camera3D) valuesChanged() bool {
	return cam.Position != cam.oldPosition || cam.LookDirection != cam.oldLookDirection || cam.Tilt != cam.oldTilt || cam.Up != cam.oldUp
}

func (cam *Camera3D) CalculateViewMatrix() {
	if cam.valuesChanged() {
		center := cam.Position.Add(cam.LookDirection.Mul(LOOK_DIRECTION_MAGNIFIER))
		cam.viewMatrix = mgl32.LookAtV(cam.Position, center, cam.Up)
		cam.inverseViewMatrix = cam.viewMatrix.Inv()
	} else {
		return
	}

	cam.oldPosition = cam.Position
	cam.oldLookDirection = cam.LookDirection
	cam.oldTilt = cam.Tilt
	cam.oldUp = cam.Up
}

func (cam *Camera3D) GetViewMatrix() mgl32.Mat4 {
	return cam.viewMatrix
}

func (cam *Camera3D) GetInverseViewMatrix() mgl32.Mat4 {
	return cam.inverseViewMatrix
}

func (cam *Camera3D) SetRotation(rot mgl32.Vec2) {
	if rot[0] > 360.0 {
		rot[0] = rot[0] - 360.0
	} else if rot[0] < -360.0 {
		rot[0] = -360.0 - rot[0]
	}
	if rot[1] > 360.0 {
		rot[1] = rot[1] - 360.0
	} else if rot[1] < -360.0 {
		rot[1] = -360.0 - rot[1]
	}

	rot[0] = mgl32.Clamp(rot[0], cam.MinRotation[0], cam.MaxRotation[0])
	rot[1] = mgl32.Clamp(rot[1], cam.MinRotation[1], cam.MaxRotation[1])

	RX := mgl32.Rotate3DX(mgl32.DegToRad(rot[0]))
	RY := mgl32.Rotate3DY(mgl32.DegToRad(rot[1]))
	matrix := RY.Mat4().Mul4(RX.Mat4())

	temp := [4]float32{0.0, 0.0, -1.0, 1.0}

	temp[0] = /*matrix.At(0, 0)*0.0 + matrix.At(0, 1)*0.0 +*/ matrix.At(0, 2)*-1.0 + matrix.At(0, 3)*1.0
	temp[1] = /*matrix.At(1, 0)*0.0 + matrix.At(1, 1)*0.0 +*/ matrix.At(1, 2)*-1.0 + matrix.At(1, 3)*1.0
	temp[2] = /*matrix.At(2, 0)*0.0 + matrix.At(2, 1)*0.0 +*/ matrix.At(2, 2)*-1.0 + matrix.At(2, 3)*1.0
	temp[3] = /*matrix.At(3, 0)*0.0 + matrix.At(3, 1)*0.0 +*/ matrix.At(3, 2)*-1.0 + matrix.At(3, 3)*1.0

	cam.LookDirection = [3]float32{temp[0] / temp[3], temp[1] / temp[3], temp[2] / temp[3]}

	cam.rotation = rot
}

func (cam *Camera3D) AddRotation(rot mgl32.Vec2) {
	cam.SetRotation(cam.rotation.Add(rot))
}

func (cam *Camera3D) AddPositionRelative(pos mgl32.Vec3) {
	if pos.Len() == 0.0 {
		return
	}

	cam.CalculateViewMatrix()
	matrix := cam.GetInverseViewMatrix()
	var worldPos mgl32.Vec3

	worldPos[0] = matrix.At(0, 0)*pos[0] + matrix.At(0, 1)*pos[1] + matrix.At(0, 2)*pos[2] + matrix.At(0, 3)*1.0
	worldPos[1] = matrix.At(1, 0)*pos[0] + matrix.At(1, 1)*pos[1] + matrix.At(1, 2)*pos[2] + matrix.At(1, 3)*1.0
	worldPos[2] = matrix.At(2, 0)*pos[0] + matrix.At(2, 1)*pos[1] + matrix.At(2, 2)*pos[2] + matrix.At(2, 3)*1.0

	worldPos = worldPos.Sub(cam.Position)

	cam.Position = cam.Position.Add(worldPos.Normalize().Mul(pos.Len()))
}
