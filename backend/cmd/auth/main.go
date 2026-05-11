package main

func main() {
	// auth-service entrypoint. Implemented in M7 (real JWT login/logout/me).
	// Until then, the gateway's dev-auth middleware (M2) provides
	// X-User-ID / X-User-Role headers without bcrypt/JWT.
}
