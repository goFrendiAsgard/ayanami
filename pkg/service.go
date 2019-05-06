package pkg

// Generator ...
type Generator interface {
	Generate() error
}

// AyanamiGenerator ...
type AyanamiGenerator struct {
	Services []Service `json:"services"`
	Flows    []Flow    `json:"flows"`
}

// Generate ...
func (generator AyanamiGenerator) Generate() error {
	return nil
}

// NewGenerator ...
func NewGenerator(jsonString string) (Generator, error) {
	var generator AyanamiGenerator
	return generator, nil
}

// Service ...
type Service struct {
	Name      string     `json:"name"`      // service & directory name
	Template  string     `json:"template"`  // template used to generate service (e.g: nodejs)
	Functions []Function `json:"functions"` // functions
}

// Flow ...
type Flow struct {
	Name  string `json:"name"`
	Nodes []Node `json:"nodes"`
}

// Edges ...
type Edges struct {
	SrcNode Function
	SrcOut  string
	DstNode Function
	DstIn   string
}

// Node ...
type Node struct {
	FunctionName string `json:"function"` // servicename.functionname || flowname.functionname
	Function     Function
}

// Function ...
type Function struct {
	ServiceName string
	Name        string   `json:"name"` // function name (e.g: getAnalytics)
	Inputs      []string `json:"inputs"`
	Outputs     []string `json:"outputs"`
	Namespace   string   `json:"namespace"` // namespace for import the function into the service (e.g: lib.js/analytics.js)
}
