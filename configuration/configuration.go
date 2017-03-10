package configuration

//noinspection ALL
type Configuration struct {
	Host     string `env:"HOST" envDefault:"irc.chat.twitch.tv"`
	Nickname string `env:"NICKNAME"`
	Oauth    string `env:"OAUTH"`
	Channel  string `env:"CHANNEL"`

	UploaderHost string `env:"UPLOADER_HOST"`
	UploaderPort int    `env:"UPLOADER_PORT"`

	StorageHost string `env:"STORAGE_HOST"`

	MessagesBufferTime int     `env:"MESSAGES_BUFFER_TIME" envDefault:"25"`
	WindowSize         int     `env:"WINDOW_SIZE" envDefault:"10"`
	SpikeRate          float32 `env:"SPIKE_RATE" envDefault:"4"`
	SmoothRate         float32 `env:"SMOOTH_RATE" envDefault:"2"`
}
