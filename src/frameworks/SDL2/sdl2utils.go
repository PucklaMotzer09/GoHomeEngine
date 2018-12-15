package framework

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/PucklaMotzer09/GoHomeEngine/src/loaders/obj"
	"github.com/PucklaMotzer09/go-sdl2/sdl"
)

func sdlKeysTogohomeKeys(key sdl.Keycode) gohome.Key {
	switch key {
	case sdl.K_UNKNOWN:
		return gohome.KeyUnknown
	case sdl.K_SPACE:
		return gohome.KeySpace
	case sdl.K_COMMA:
		return gohome.KeyComma
	case sdl.K_MINUS:
		return gohome.KeyMinus
	case sdl.K_PERIOD:
		return gohome.KeyPeriod
	case sdl.K_SLASH:
		return gohome.KeySlash
	case sdl.K_0:
		return gohome.Key0
	case sdl.K_1:
		return gohome.Key1
	case sdl.K_2:
		return gohome.Key2
	case sdl.K_3:
		return gohome.Key3
	case sdl.K_4:
		return gohome.Key4
	case sdl.K_5:
		return gohome.Key5
	case sdl.K_6:
		return gohome.Key6
	case sdl.K_7:
		return gohome.Key7
	case sdl.K_8:
		return gohome.Key8
	case sdl.K_9:
		return gohome.Key9
	case sdl.K_SEMICOLON:
		return gohome.KeySemicolon
	case sdl.K_EQUALS:
		return gohome.KeyEqual
	case sdl.K_a:
		return gohome.KeyA
	case sdl.K_b:
		return gohome.KeyB
	case sdl.K_c:
		return gohome.KeyC
	case sdl.K_d:
		return gohome.KeyD
	case sdl.K_e:
		return gohome.KeyE
	case sdl.K_f:
		return gohome.KeyF
	case sdl.K_g:
		return gohome.KeyG
	case sdl.K_h:
		return gohome.KeyH
	case sdl.K_i:
		return gohome.KeyI
	case sdl.K_j:
		return gohome.KeyJ
	case sdl.K_k:
		return gohome.KeyK
	case sdl.K_l:
		return gohome.KeyL
	case sdl.K_m:
		return gohome.KeyM
	case sdl.K_n:
		return gohome.KeyN
	case sdl.K_o:
		return gohome.KeyO
	case sdl.K_p:
		return gohome.KeyP
	case sdl.K_q:
		return gohome.KeyQ
	case sdl.K_r:
		return gohome.KeyR
	case sdl.K_s:
		return gohome.KeyS
	case sdl.K_t:
		return gohome.KeyT
	case sdl.K_u:
		return gohome.KeyU
	case sdl.K_v:
		return gohome.KeyV
	case sdl.K_w:
		return gohome.KeyW
	case sdl.K_x:
		return gohome.KeyX
	case sdl.K_y:
		return gohome.KeyY
	case sdl.K_z:
		return gohome.KeyZ
	case sdl.K_LEFTBRACKET:
		return gohome.KeyLeftBracket
	case sdl.K_BACKSLASH:
		return gohome.KeyBackslash
	case sdl.K_RIGHTBRACKET:
		return gohome.KeyRightBracket
	case sdl.K_ESCAPE:
		return gohome.KeyEscape
	case sdl.K_RETURN:
		return gohome.KeyEnter
	case sdl.K_TAB:
		return gohome.KeyTab
	case sdl.K_BACKSPACE:
		return gohome.KeyBackspace
	case sdl.K_INSERT:
		return gohome.KeyInsert
	case sdl.K_DELETE:
		return gohome.KeyDelete
	case sdl.K_RIGHT:
		return gohome.KeyRight
	case sdl.K_LEFT:
		return gohome.KeyLeft
	case sdl.K_DOWN:
		return gohome.KeyDown
	case sdl.K_UP:
		return gohome.KeyUp
	case sdl.K_PAGEUP:
		return gohome.KeyPageUp
	case sdl.K_PAGEDOWN:
		return gohome.KeyPageDown
	case sdl.K_HOME:
		return gohome.KeyHome
	case sdl.K_END:
		return gohome.KeyEnd
	case sdl.K_CAPSLOCK:
		return gohome.KeyCapsLock
	case sdl.K_SCROLLLOCK:
		return gohome.KeyScrollLock
	case sdl.K_NUMLOCKCLEAR:
		return gohome.KeyNumLock
	case sdl.K_PRINTSCREEN:
		return gohome.KeyPrintScreen
	case sdl.K_PAUSE:
		return gohome.KeyPause
	case sdl.K_F1:
		return gohome.KeyF1
	case sdl.K_F2:
		return gohome.KeyF2
	case sdl.K_F3:
		return gohome.KeyF3
	case sdl.K_F4:
		return gohome.KeyF4
	case sdl.K_F5:
		return gohome.KeyF5
	case sdl.K_F6:
		return gohome.KeyF6
	case sdl.K_F7:
		return gohome.KeyF7
	case sdl.K_F8:
		return gohome.KeyF8
	case sdl.K_F9:
		return gohome.KeyF9
	case sdl.K_F10:
		return gohome.KeyF10
	case sdl.K_F11:
		return gohome.KeyF11
	case sdl.K_F12:
		return gohome.KeyF12
	case sdl.K_F13:
		return gohome.KeyF13
	case sdl.K_F14:
		return gohome.KeyF14
	case sdl.K_F15:
		return gohome.KeyF15
	case sdl.K_F16:
		return gohome.KeyF16
	case sdl.K_F17:
		return gohome.KeyF17
	case sdl.K_F18:
		return gohome.KeyF18
	case sdl.K_F19:
		return gohome.KeyF19
	case sdl.K_F20:
		return gohome.KeyF20
	case sdl.K_F21:
		return gohome.KeyF21
	case sdl.K_F22:
		return gohome.KeyF22
	case sdl.K_F23:
		return gohome.KeyF23
	case sdl.K_F24:
		return gohome.KeyF24
	case sdl.K_KP_0:
		return gohome.KeyKP0
	case sdl.K_KP_1:
		return gohome.KeyKP1
	case sdl.K_KP_2:
		return gohome.KeyKP2
	case sdl.K_KP_3:
		return gohome.KeyKP3
	case sdl.K_KP_4:
		return gohome.KeyKP4
	case sdl.K_KP_5:
		return gohome.KeyKP5
	case sdl.K_KP_6:
		return gohome.KeyKP6
	case sdl.K_KP_7:
		return gohome.KeyKP7
	case sdl.K_KP_8:
		return gohome.KeyKP8
	case sdl.K_KP_9:
		return gohome.KeyKP9
	case sdl.K_KP_DECIMAL:
		return gohome.KeyKPDecimal
	case sdl.K_KP_DIVIDE:
		return gohome.KeyKPDivide
	case sdl.K_KP_MULTIPLY:
		return gohome.KeyKPMultiply
	case sdl.K_KP_MEMSUBTRACT:
		return gohome.KeyKPSubtract
	case sdl.K_KP_MEMADD:
		return gohome.KeyKPAdd
	case sdl.K_KP_ENTER:
		return gohome.KeyKPEnter
	case sdl.K_KP_EQUALS:
		return gohome.KeyKPEqual
	case sdl.K_LSHIFT:
		return gohome.KeyLeftShift
	case sdl.K_LCTRL:
		return gohome.KeyLeftControl
	case sdl.K_LALT:
		return gohome.KeyLeftAlt
	case sdl.K_RSHIFT:
		return gohome.KeyRightShift
	case sdl.K_RCTRL:
		return gohome.KeyRightControl
	case sdl.K_RALT:
		return gohome.KeyRightAlt
	case sdl.K_MENU:
		return gohome.KeyMenu
	}

	return gohome.KeyUnknown
}

func sdlMouseButtonTogohomeKeys(mb uint8) gohome.Key {
	switch mb {
	case sdl.BUTTON_LEFT:
		return gohome.MouseButtonLeft
	case sdl.BUTTON_RIGHT:
		return gohome.MouseButtonRight
	case sdl.BUTTON_MIDDLE:
		return gohome.MouseButtonMiddle
	}

	return gohome.MouseButtonLast
}

func loadLevelOBJ(rsmgr *gohome.ResourceManager, name, path string, preloaded, loadToGPU bool) *gohome.Level {
	gohome.Framew.Log("LoadLevelOBJ in Framework")
	return loader.LoadLevelOBJ(rsmgr, name, path, preloaded, loadToGPU)
}

func loadLevelOBJString(rsmgr *gohome.ResourceManager, name, contents, fileName string, preloaded, loadToGPU bool) *gohome.Level {
	return loader.LoadLevelOBJString(rsmgr, name, contents, fileName, preloaded, loadToGPU)
}
