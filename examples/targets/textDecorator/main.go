package main

const serviceName = "textDecorator"

var configs SrvcConfigs

func init() {
	configs = SrvcConfigs{
		"pre": SrvcNewServiceConfig(
			serviceName,
			"pre",
			[]string{"text"},
			[]string{"text"},
			WrappedPre,
		),
		"cowsay": SrvcNewServiceConfig(
			serviceName,
			"cowsay",
			[]string{"text"},
			[]string{"text"},
			WrappedCowsay,
		),
		"figlet": SrvcNewServiceConfig(
			serviceName,
			"figlet",
			[]string{"text"},
			[]string{"text"},
			WrappedFiglet,
		),
	}
}

func main() {
	// consume and publish forever
	ch := make(chan bool)
	SrvcConsumeAndPublish(configs)
	<-ch
}
