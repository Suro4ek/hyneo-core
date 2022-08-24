package command

var commands = make(map[string]*Command)

func RegisterCommand(command *Command) {
	commands[command.Name] = command
}

func RegisterCommands() {
	RegisterCommand(Bind)
	RegisterCommand(UnBind)
	RegisterCommand(Account)
	RegisterCommand(Accounts)
	RegisterCommand(Kick)
}

func GetCommands() map[string]*Command {
	return commands
}
