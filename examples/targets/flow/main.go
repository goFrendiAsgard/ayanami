package main

var configs Configs

func init() {
	configs = Configs{
		"echo": SingleConfig{
			Input: StringDictionary{
				"trig.request.get./echo.out.form": "form",
			},
			Output: StringDictionary{
				"code":    "trig.response.get./echo.in.code",
				"content": "trig.response.get./echo.in.content",
			},
			Function: WrappedEcho,
		},
		// TODO: add (cowsay -> pre, figlet -> pre, and figlet -> cowsay -> pre)
	}
}

func main() {
	// consume and publish forever
	ch := make(chan bool)
	ConsumeAndPublish(configs)
	<-ch
}
