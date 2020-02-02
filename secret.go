package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mdp/sodiumbox"
)

type GithubPublickey struct {
	Key   string `json:"key"`
	KeyID string `json:"key_id"`
}

// Secret is the represenation of a single secret, node this contains on the name of the secret
// and not the actual decrypted value
type Secret struct {
	CreatedAt string `json:"created_at"`
	Name      string `json:"name"`
	UpdatedAt string `json:"updated_at"`
}

// Secrets contain all the secrets from a repository
type Secrets struct {
	TotalCount int      `json:"total_count"`
	Secrets    []Secret `json:"secrets"`
}

// GetPublickey retrieves the public key from a give repository in a given organizaion
func GetPublickey(org, repo, token string) *GithubPublickey {
	u := fmt.Sprintf("https://api.github.com/repos/%v/%v/actions/secrets/public-key", org, repo)

	req, err := http.NewRequest("GET", u, nil)
	req.Header.Add("Authorization", "token "+token)
	req.Header.Add("Accept", "application/vnd.github.v3.patch")
	if err != nil {
		log.Println(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	gpk := &GithubPublickey{}
	err = json.NewDecoder(resp.Body).Decode(gpk)
	if err != nil {
		log.Println(err)
	}
	return gpk
}

// GetSecrets get's all the secrets from a repository in a organization
func GetSecrets(org, repo, token string) *Secrets {
	u := fmt.Sprintf("https://api.github.com/repos/%v/%v/actions/secrets", org, repo)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("Authorization", "token "+token)
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}

	gpk := &Secrets{}
	err = json.NewDecoder(resp.Body).Decode(gpk)
	if err != nil {
		log.Println(err)
	}
	return gpk
}

// GetSecret finds a given secret in a repository
func GetSecret(org, repo, name, token string) *Secret {

	u := fmt.Sprintf("https://api.github.com/repos/%v/%v/actions/secrets/%v", org, repo, name)

	req, err := http.NewRequest("GET", u, nil)
	req.Header.Add("Authorization", "token "+token)
	req.Header.Add("Accept", "application/vnd.github.v3.patch")
	if err != nil {
		log.Println(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}

	gpk := &Secret{}
	err = json.NewDecoder(resp.Body).Decode(gpk)
	if err != nil {
		log.Println(err)
	}
	return gpk
}

// Updatesecret will create or update a given secret. Please note this calls PrivateKey in order
// to fetch the given public for that repostiory in the organization
func Updatesecret(org, repo, name, value, token string) error {
	publickey := GetPublickey(org, repo, token)
	decryptedKey, err := base64.StdEncoding.DecodeString(publickey.Key)
	if err != nil {
		log.Println(err)

		return err
	}

	pk := *new([32]byte)
	copy(pk[:], decryptedKey[0:32])

	sealedMsg, err := sodiumbox.Seal([]byte("secretmessage"), &pk)
	encryptedKey := base64.StdEncoding.EncodeToString(sealedMsg.Box)

	type EncSecret struct {
		EncryptedValue string `json:"encrypted_value"`
		KeyID          string `json:"key_id"`
	}

	payload := EncSecret{KeyID: publickey.KeyID, EncryptedValue: encryptedKey}
	pdata, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
		return err
	}

	u := fmt.Sprintf("https://api.github.com/repos/%v/%v/actions/secrets/%v", org, repo, name)
	req, err := http.NewRequest("PUT", u, bytes.NewBuffer(pdata))
	req.Header.Add("Authorization", "token "+token)
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.Header.Add("User-Agent", "Awesome-Octocat-App")
	req.Header.Add("Content-type", "application/json")
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	if resp.StatusCode == 204 || resp.StatusCode == 201 {
		return nil
	}
	return fmt.Errorf("Error updating secret: %v", err)
}

// DeleteSecret deletes a given secret in a repository
func DeleteSecret(org, repo, name, token string) error {

	u := fmt.Sprintf("https://api.github.com/repos/%v/%v/actions/secrets/%v", org, repo, name)

	req, err := http.NewRequest("DELETE", u, nil)
	req.Header.Add("Authorization", "token "+token)
	req.Header.Add("Accept", "application/vnd.github.v3.patch")
	if err != nil {
		log.Println(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("Error updating secret: %v", err)
	}
	return nil
}
