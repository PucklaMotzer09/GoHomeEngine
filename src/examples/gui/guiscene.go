package main

import (
	"fmt"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
)

type GUIScene struct {
	btn  gohome.Button
	btn1 gohome.Button
}

func (this *GUIScene) Init() {
	gohome.Init2DShaders()
	this.btn.Init(gohome.Framew.WindowGetSize().Mul(0.25), "")
	this.btn.Transform.Origin = [2]float32{0.5, 0.5}
	this.btn.PressCallback = func(btn *gohome.Button) {
		fmt.Println("Pressed 0")
	}
	this.btn1.Init(gohome.Framew.WindowGetSize().Mul(0.75), "")
	this.btn1.Transform.Origin = [2]float32{0.5, 0.5}
	this.btn1.PressCallback = func(btn *gohome.Button) {
		fmt.Println("Pressed 1")
	}
}

func (this *GUIScene) Update(delta_time float32) {
}

func (this *GUIScene) Terminate() {

}
