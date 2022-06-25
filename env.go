package main

import "os"

func envOr(env, def string) string {
	e := os.Getenv(env)
	if e == "" {
		return def
	}
	return e
}
