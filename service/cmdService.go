package service

// CmdService single flow config
type CmdService struct {
	MethodName string
	Input      []IO
	Output     []IO
	Cmd        []string
}
