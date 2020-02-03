package secret

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_publickey(t *testing.T) {
	org := "sorenmat"
	repo := "github-secret"

	token := os.Getenv("GITHUB_TOKEN")
	pk := GetPublickey(org, repo, token)
	assert.NotNil(t, pk)

	err := Updatesecret(org, repo, "my-super-secret", "testing", token)
	assert.NoError(t, err)

	secrets := GetSecrets(org, repo, token)
	assert.NotNil(t, secrets)

	secret := GetSecret(org, repo, "my-super-secret", token)
	assert.NotNil(t, secret)
	assert.Equal(t, "my-super-secret", secret.Name)

	//	err = DeleteSecret(org, repo, "my-super-secret", token)
	///	assert.NoError(t, err)

}
