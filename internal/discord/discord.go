package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"go.mlcdf.fr/sally/build"
)

type Client struct {
	WebhookURL string
}

// Webhook is the webhook object sent to discord
type Webhook struct {
	Username  string  `json:"username"`
	AvatarURL string  `json:"avatar_url"`
	Content   string  `json:"content"`
	Embeds    []Embed `json:"embeds"`
}

// Embed is the embed object
type Embed struct {
	Author      Author  `json:"author"`
	Title       string  `json:"title"`
	URL         string  `json:"url"`
	Description string  `json:"description"`
	Color       int64   `json:"color"`
	Fields      []Field `json:"fields"`
	Thumbnail   Image   `json:"thumbnail"`
	Image       Image   `json:"image"`
	Footer      Footer  `json:"footer"`
	TimeStamp   string  `json:"timestamp"`
}

// Author is the author object
type Author struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	IconURL string `json:"icon_url"`
}

// Field is the field object inside an embed
type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

// Footer is the footer of the embed
type Footer struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url"`
}

// Image is an image possibly contained inside the embed
type Image struct {
	URL string `json:"url"`
}

func NewWebhook() *Webhook {
	webhook := &Webhook{Username: build.String(), Embeds: make([]Embed, 0)}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "(unknown)"
	}
	webhook.Embeds[0].Footer.Text = hostname
	return webhook
}

func (c *Client) PostInfo(webhook *Webhook) error {
	webhook.Embeds[0].Color = 2201331
	return c.Post(webhook)
}

func (c *Client) PostError(webhook *Webhook) error {
	webhook.Embeds[0].Color = 15092300
	return c.Post(webhook)
}

func (c *Client) PostSuccess(webhook *Webhook) error {
	webhook.Embeds[0].Color = 5747840
	return c.Post(webhook)
}

func (c *Client) Post(webhook *Webhook) error {
	payload, err := json.Marshal(webhook)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %s", err)
	}

	res, err := http.Post(c.WebhookURL, "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to post to webhook: %s", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed read body response : %s", err)
	}

	if res.StatusCode >= 400 || err != nil {
		return fmt.Errorf("failed to post to webhook: reason=%s status=%s", body, res.Status)
	}
	return nil
}

func (c *Client) Write(p []byte) (n int, err error) {
	w := NewWebhook()
	w.Embeds[0] = Embed{
		Description: string(p),
	}

	return len(p), c.PostError(w)
}
