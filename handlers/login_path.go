package handlers

const loginPagePath = "/login"

func LoginErrorPath(code string) string {
	if code == "" {
		return loginPagePath
	}

	return loginPagePath + "?error=" + code
}
