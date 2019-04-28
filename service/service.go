package service

// Generator ...
type Generator interface {
	generate() error
}

// Service ...
type Service struct {
	Name      string     // service & directory name
	Template  string     // template used to generate service (e.g: nodejs)
	Functions []Function // functions
}

// Function ...
type Function struct {
	Name      string   // function name (e.g: getAnalytics)
	Inputs    []string // input names
	Namespace string   // namespace (e.g: lib.js/analytics.js)
}
