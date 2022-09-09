package keyboard

import "hyneo/internal/config"

//TODO реализовать клавиатуру здесь и возможно перенести в конфиг как то или в бд
//TODO подумать над кастомными клавиатурами и как они будут работать в боте
//TODO реализовать команды на сервере
var Keyboard = make([]config.KeyboardConfig, 0)

func Init(cfg []config.KeyboardConfig) {
	Keyboard = cfg
}
