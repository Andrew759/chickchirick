package config

type contextKey string

const (
	UserUserKey     contextKey = "validatedUser"
	UserPropertyKey contextKey = "validatedUserProperty"
)
