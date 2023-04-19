package handlers

import (
	"context"
	"net/http"
)

var (
	client *http.Client

	ctx = context.Background()
)
