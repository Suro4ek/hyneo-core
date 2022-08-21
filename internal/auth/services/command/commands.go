package command

var commands = make(map[string]*Command)

func RegisterCommand(command *Command) {
	commands[command.Name] = command
}

func RegisterCommands() {
	RegisterCommand(Bind)
	RegisterCommand(UnBind)
	RegisterCommand(Account)
	RegisterCommand(ClearKeyboard)
	RegisterCommand(Accounts)
}

func GetCommands() map[string]*Command {
	return commands
}
