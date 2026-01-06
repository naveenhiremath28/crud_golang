package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Helper structs for Vault API
type vaultRequest struct {
	Plaintext  string `json:"plaintext,omitempty"`  // Base64 encoded for encrypt
	Ciphertext string `json:"ciphertext,omitempty"` // For decrypt
}

type vaultResponse struct {
	Data struct {
		Ciphertext string `json:"ciphertext,omitempty"`
		Plaintext  string `json:"plaintext,omitempty"`
	} `json:"data"`
}

func (s *Service) encryptField(payload string, vaultEntityID string) (string, error) {
	if payload == "" {
		return "", nil
	}

	// 1. Base64 Encode
	encodedPayload := base64.StdEncoding.EncodeToString([]byte(payload))

	// 2. Prepare Vault Request
	reqBody := vaultRequest{
		Plaintext: encodedPayload,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/v1/transit/encrypt/%s", s.Config.VaultURL, vaultEntityID)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("X-Vault-Token", s.Config.VaultToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("vault encrypt failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// 4. Parse Response
	var vResp vaultResponse
	if err := json.Unmarshal(bodyBytes, &vResp); err != nil {
		return "", fmt.Errorf("failed to decode vault response: %w, body: %s", err, string(bodyBytes))
	}

	if vResp.Data.Ciphertext == "" {
		return "", errors.New("vault returned empty ciphertext")
	}

	return vResp.Data.Ciphertext, nil
}

func (s *Service) decryptField(ciphertext string, vaultEntityID string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	// 1. Prepare Vault Request
	reqBody := vaultRequest{
		Ciphertext: ciphertext,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// 2. Call Vault Decrypt API
	// URL: {VAULT_URL}/decrypt/{VaultEntityID}
	url := fmt.Sprintf("%s/v1/transit/decrypt/%s", s.Config.VaultURL, vaultEntityID)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("X-Vault-Token", s.Config.VaultToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("vault decrypt failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// 3. Parse Response
	var vResp vaultResponse
	if err := json.Unmarshal(bodyBytes, &vResp); err != nil {
		return "", fmt.Errorf("failed to decode vault response: %w, body: %s", err, string(bodyBytes))
	}

	if vResp.Data.Plaintext == "" {
		// Valid case? Maybe if input was empty, but we checked that.
		return "", nil
	}

	// 4. Base64 Decode
	decodedBytes, err := base64.StdEncoding.DecodeString(vResp.Data.Plaintext)
	if err != nil {
		return "", err
	}

	return string(decodedBytes), nil
}
