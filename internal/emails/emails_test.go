package emails

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/masoncitemple4/certwatch/internal/defaults"
)

func TestSendEmail(t *testing.T) {

	err := godotenv.Load(defaults.DEFAULT_ENV_PATH)
	if err != nil {
		t.Errorf("error loading env file: %v", err)
	}

	// TODO: Might want to replace this with real emails
	// in the future.
	recipients := []string{
		"test1@example.com",
		"test2@example.com",
	}

	err = SendEmail(recipients, "GoLang Test Email", "Some really simple and boring text.")
	if err != nil {
		t.Errorf("error sending email: %v", err)
	}

}
