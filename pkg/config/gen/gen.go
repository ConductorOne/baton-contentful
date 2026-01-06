package main

import (
	cfg "github.com/conductorone/baton-contentful/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/config"
)

func main() {
	config.Generate("contentful", cfg.Config)
}
