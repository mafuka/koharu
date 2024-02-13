package bot

const AUTH_URL = "https://bots.qq.com/app/getAppAccessToken"

type AuthConfig struct {
	APPID     uint   `yaml:"app-id"`
	APPSecret string `yaml:"app-secret"`
}
