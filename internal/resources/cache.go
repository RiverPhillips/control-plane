package resources

type Listener struct {
	Name    string
	Routes  []string
	Address string
	Port    uint32
}

type Route struct {
	Name    string
	Prefix  string
	Service string
}

type Cluster struct {
	Name string
}
