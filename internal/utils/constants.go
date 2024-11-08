package utils

import (
    "net/http"
    "time"
)

const (
    DefaultTimeout = 10 * time.Second
)

var (
    SupportedMethods = []string{"GET", "POST"}
    DefaultHTTPClient = &http.Client{
        Timeout: DefaultTimeout,
    }
)

