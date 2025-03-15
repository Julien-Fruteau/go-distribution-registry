package registry

type Registry struct {
	BaseUrl string
	Conf    Conf
}

type Conf struct {
	Host     string
	Scheme   string
	Username string
	Password string
	Mime     string
}

