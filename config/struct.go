package config

type CStruct struct {
	Local LocalStruct `yaml:"Local"`
	Gitea GiteaStruct `yaml:"Gitea"`
	Cos   CosStruct   `yaml:"Cos"`
}

type LocalStruct struct {
	Dir        string `yaml:"Dir"`
	MaxFileNum int    `yaml:"MaxFileNum"`
}

type GiteaStruct struct {
	BinPath string `yaml:"BinPath"`
	WorkDir string `yaml:"WorkDir"`
	Cron    string `yaml:"Cron"`
}

type CosStruct struct {
	Region string `yaml:"Region"`
	Bucket string `yaml:"Bucket"`
	Secret struct {
		Id  string `yaml:"ID"`
		Key string `yaml:"Key"`
	} `yaml:"Secret"`
	Path       string `yaml:"Path"`
	MaxFileNum int    `yaml:"MaxFileNum"`
}
