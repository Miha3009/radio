package cloud

type Config struct {
	Bucket string `yaml:"bucket"`
	Key    string `yaml:"key"`
	Secret string `yaml:"secret"`
}
