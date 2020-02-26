package auto

import (
	"github.com/nomeaning777/exec-with-secret"
	"log"
)

func init() {
	if err := secret.InjectSecretToEnvironment(); err != nil {
		log.Fatal(err)
	}
}
