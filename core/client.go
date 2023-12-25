package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/xerrors"
	"io"
	"net/http"
	"sync"
	"time"
)

// Client represents a client for Mirai HTTP API.
type Client struct {
	ID      int
	Key     string
	Address string

	session string
	mu      sync.RWMutex // RW lock for session

	httpClient *http.Client
}

type ClientConfig struct {
	ID      int    `yaml:"id"`
	Key     string `yaml:"key"`
	Address string `yaml:"address"`
	Timeout int    `yaml:"timeout"`
}

func NewClient(cfg ClientConfig) *Client {
	c := &Client{
		ID:         cfg.ID,
		Key:        cfg.Key,
		Address:    cfg.Address,
		session:    "",
		httpClient: &http.Client{Timeout: time.Duration(cfg.Timeout) * time.Second},
	}
	return c
}

func (c *Client) SetSession(session string) {
	c.mu.Lock()
	c.session = session
	c.mu.Unlock()
}

func (c *Client) GetSession() string {
	c.mu.RLock()
	session := c.session
	c.mu.RUnlock()
	return session
}

// Code represents an error or success code from Mirai HTTP API.
type Code uint

// Defined Mirai HTTP API response codes.
const (
	Success          Code = 0
	InvalidVerifyKey Code = 1
	BotNotFound      Code = 2
	InvalidSession   Code = 3
	Unauthenticated  Code = 4
	TargetNotFound   Code = 5
	FileNotFound     Code = 6
	NoPermission     Code = 10
	BotMuted         Code = 20
	MessageTooLong   Code = 30
	BadRequest       Code = 400
)

// String provides the string representation of the Code.
func (c Code) String() string {
	switch c {
	case Success:
		return "Success"
	case InvalidVerifyKey:
		return "Invalid VerifyKey"
	case BotNotFound:
		return "Bot Not Found"
	case InvalidSession:
		return "Invalid Session"
	case Unauthenticated:
		return "Unauthenticated"
	case TargetNotFound:
		return "Target Not Found"
	case FileNotFound:
		return "File Not Found"
	case NoPermission:
		return "No Permission"
	case BotMuted:
		return "Bot Muted"
	case MessageTooLong:
		return "Message Too Long"
	case BadRequest:
		return "Bad Request"
	default:
		return fmt.Sprintf("Unknown Code: %d", c)
	}
}

// POST sends a POST request to the Mirai HTTP API and parses the response.
func (c *Client) POST(endpoint string, body interface{}, result interface{}) error {
	url := c.Address + endpoint

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return xerrors.Errorf("%w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return xerrors.Errorf("%w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	session := c.GetSession()
	if session != "" {
		req.Header.Set("Authorization", "session "+session)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return xerrors.Errorf("Failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return xerrors.Errorf("%w", err)
	}

	err = json.Unmarshal(respBytes, result)
	if err != nil {
		return xerrors.Errorf("%w", err)
	}

	return nil
}

// Verify uses Key to authenticate the Client and returns a session.
func (c *Client) Verify() (string, error) {
	in := struct {
		VerifyKey string `json:"verifyKey"`
	}{c.Key}

	out := struct {
		Code    `json:"code"`
		Session string `json:"session"`
	}{}

	err := c.POST("/verify", in, &out)
	if err != nil {
		return "", xerrors.Errorf("Failed to request Verify API: %w", err)
	}

	if out.Code != Success {
		return "", xerrors.Errorf(out.Code.String())
	}

	return out.Session, nil
}

func (c *Client) Bind(session string) error {
	in := struct {
		SessionKey string `json:"sessionKey"`
		QQ         int    `json:"qq"`
	}{session, c.ID}

	out := struct {
		Code `json:"code"`
	}{}

	err := c.POST("/bind", in, &out)
	if err != nil {
		return xerrors.Errorf("Failed to request Bind API: %w", err)
	}

	if out.Code != Success {
		return xerrors.Errorf(out.Code.String())
	}

	c.SetSession(session)
	return nil
}

func (c *Client) Release(session string) error {
	in := struct {
		SessionKey string `json:"sessionKey"`
		QQ         int    `json:"qq"`
	}{session, c.ID}

	out := struct {
		Code `json:"code"`
	}{}

	err := c.POST("/release", in, &out)
	if err != nil {
		return xerrors.Errorf("Failed to request Release API: %w", err)
	}

	if out.Code != Success {
		return xerrors.Errorf(out.Code.String())
	}

	c.SetSession("")
	return nil
}
