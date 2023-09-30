package config

type Enviroment struct {
	Micro Micro
}

type Micro struct {
	Name string
	Port string
	Host string
}
