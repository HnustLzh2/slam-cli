package caller

import (
	configpkg "slam-cli/internal/pkg/config"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	defaultSLAMAPIBaseURL = "https://slam.byted.org/api/slam/v2"
	headerXJWTToken       = "X-Jwt-Token"
)

type Client struct {
	baseURL string
	token   string
	http    *resty.Client
	auth    configpkg.AuthConfig
}

type tokenFile struct {
	Region  string    `json:"region"`
	Token   string    `json:"token"`
	SavedAt time.Time `json:"saved_at"`
}


