package main

const serviceName = "textDecorator"

var configs Configs

func init() {
	configs = Configs{
		"pre": NewServiceConfig(
			serviceName,
			"pre",
			WrappedPre,
			[]string{"text"},
			[]string{"text"},
		),
		"cowsay": NewServiceConfig(
			serviceName,
			"cowsay",
			WrappedCowsay,
			[]string{"text"},
			[]string{"text"},
		),
		"figlet": NewServiceConfig(
			serviceName,
			"figlet",
			WrappedFiglet,
			[]string{"text"},
			[]string{"text"},
		),
	}
}

func main() {
	// consume and publish forever
	ch := make(chan bool)
	ConsumeAndPublish(configs)
	<-ch
}
