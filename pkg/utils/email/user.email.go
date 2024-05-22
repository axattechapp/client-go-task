package email

import (
	"client_task/pkg/common/config"
	"fmt"
	"strconv"

	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

type EmailConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Sender   string `json:"sender"`
}

func LoadEmailConfig() (*EmailConfig, error) {
	err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	portString := viper.GetString("EMAIL_PORT")
	port, err := strconv.Atoi(portString)
	if err != nil {
		fmt.Printf("Error converting EMAIL_PORT to int: %v\n", err)
		return nil, err
	}

	return &EmailConfig{
		Host:     viper.GetString("EMAIL_HOST"),
		Port:     port,
		Username: viper.GetString("EMAIL_USERNAME"),
		Password: viper.GetString("EMAIL_PASSWORD"),
		Sender:   viper.GetString("EMAIL_SENDER"),
	}, nil
}

type gomailSender struct {
	config *EmailConfig
}

func NewGomailSender(config *EmailConfig) *gomailSender {
	return &gomailSender{config: config}
}

func SendEmail(config *EmailConfig, to string, subject, body string) error {

	m := gomail.NewMessage()
	m.SetHeader("From", config.Sender)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)

	return d.DialAndSend(m)

}
