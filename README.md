# github-secret
Go library for accessing Githubs secrets API


## Build
`go build`

## Install
`go install`


## Godoc

```
FUNCTIONS

func DeleteSecret(org, repo, name, token string) error
    DeleteSecret deletes a given secret in a repository

func Updatesecret(org, repo, name, value, token string) error
    Updatesecret will create or update a given secret. Please note this calls
    PrivateKey in order to fetch the given public for that repostiory in the
    organization


TYPES

type GithubPublickey struct {
	Key   string `json:"key"`
	KeyID string `json:"key_id"`
}

func GetPublickey(org, repo, token string) *GithubPublickey
    GetPublickey retrieves the public key from a give repository in a given
    organizaion

type Secret struct {
	CreatedAt string `json:"created_at"`
	Name      string `json:"name"`
	UpdatedAt string `json:"updated_at"`
}
    Secret is the represenation of a single secret, node this contains on the
    name of the secret and not the actual decrypted value

func GetSecret(org, repo, name, token string) *Secret
    GetSecret finds a given secret in a repository

type Secrets struct {
	TotalCount int      `json:"total_count"`
	Secrets    []Secret `json:"secrets"`
}
    Secrets contain all the secrets from a repository

func GetSecrets(org, repo, token string) *Secrets
    GetSecrets get's all the secrets from a repository in a organization

```