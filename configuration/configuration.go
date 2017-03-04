package configuration

//noinspection ALL
type Configuration struct {
	Host         string `env:"HOST" envDefault:"irc.chat.twitch.tv"`
	Nickname     string `env:"NICKNAME"`
	Oauth        string `env:"OAUTH"`
	Channel      string `env:"CHANNEL"`
	UploaderHost string `env:"UPLOADER_HOST"`
	UploaderPort int    `env:"UPLOADER_PORT"`
}
