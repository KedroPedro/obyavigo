package mail

import (
	"crypto/aes"

	"fmt"

	"cmd/obyavigo/main.go/internal/config"
	"cmd/obyavigo/main.go/internal/secure"

	"github.com/wneessen/go-mail"
)

type Mail struct {
	client *mail.Client
	Cipher *secure.AESMail
}

func New(cfg *config.Config) (*Mail, error) {
	client, err := mail.NewClient(
		cfg.EMail.SMTPServer,
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(cfg.EMail.Username),
		mail.WithPassword(cfg.EMail.Password),
		mail.WithPort(cfg.EMail.Port),
		mail.WithTLSPortPolicy(mail.TLSMandatory),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create mail client: %w", err)
	}

	key := []byte(cfg.Secure.SecretKey)
	block, err := aes.NewCipher(key)
	cipher := secure.NewAESMail(block)

	return &Mail{
		client: client,
		Cipher: cipher,
	}, nil
}

func (m *Mail) SendRegConfirm(email string) (string, error) {
	msg := mail.NewMsg()

	if err := msg.From("obyavigo@gmail.com"); err != nil {
		return "", fmt.Errorf("failed to set From address: %w", err)
	}
	if err := msg.To(email); err != nil {
		return "", fmt.Errorf("failed to set To address: %w", err)
	}

	cipheredMail, err := m.Cipher.EncryptMail(email)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt mail: %w", err)
	}

	msg.Subject("Подтверждение регистрации аккаунта - Obyavigo")
	msg.SetBodyString(mail.TypeTextPlain,
		"Перейдите по этой ссылке для подтверждения регистрации: https://obyavigo.by/api/auth/confirm-email/"+cipheredMail)

	if err := m.client.DialAndSend(msg); err != nil {
		return "", fmt.Errorf("failed to send mail: %w", err)
	}

	return cipheredMail, nil
}

func (m *Mail) SendPasswordReset(email, token string) error {
	msg := mail.NewMsg()

	if err := msg.From("obyavigo@gmail.com"); err != nil {
		return fmt.Errorf("failed to set From address: %w", err)
	}
	if err := msg.To(email); err != nil {
		return fmt.Errorf("failed to set To address: %w", err)
	}

	msg.Subject("Восстановление пароля - Obyavigo")
	msg.SetBodyString(mail.TypeTextPlain,
<<<<<<< HEAD
		"Для восстановления пароля перейдите по ссылке: https://localhost:443/reset-password?token="+token+"\n\n")
=======
		"Для восстановления пароля перейдите по ссылке: https://obyavigo.by/reset-password?token="+token+"\n")
>>>>>>> c04cbe9c777b551f29c288a2c9d239c0b97177a5

	if err := m.client.DialAndSend(msg); err != nil {
		return fmt.Errorf("failed to send mail: %w", err)
	}

	return nil
}
