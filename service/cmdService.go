package service

// CmdService single flow config
type CmdService struct {
	Input  []IO
	Output []IO
	Cmd    []string
}

// NewCmd create new cmd
func NewCmd(cmdConfig CmdService) CommonService {
	var service CommonService
	service.Input = cmdConfig.Input
	service.Output = cmdConfig.Output
	return service
}
