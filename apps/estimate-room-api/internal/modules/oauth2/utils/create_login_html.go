package oauth2utils

import (
	"fmt"
	"html"
	"strings"
)

var formKeyMap = map[string]string{
	"ClientID":            "client_id",
	"RedirectURI":         "redirect_uri",
	"ResponseType":        "response_type",
	"Scopes":              "scopes",
	"State":               "state",
	"CodeChallenge":       "code_challenge",
	"CodeChallengeMethod": "code_challenge_method",
	"Nonce":               "nonce",
}

func CreateLoginHtml(params map[string]any) string {
	var hiddenFields strings.Builder
	for key, val := range params {
		formKey := key
		mapped, ok := formKeyMap[key]

		if ok {
			formKey = mapped
		}

		strVal := fmt.Sprintf("%v", val)
		hiddenFields.WriteString(`<input type="hidden" name="` + html.EscapeString(formKey) + `" value="` + html.EscapeString(strVal) + `"/>` + "\n")
	}

	return `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>Login</title>
	<style>
		body { font-family: Arial, sans-serif; background: #f2f2f2; display: flex; justify-content: center; align-items: center; height: 100vh; }
		.login-container { background: white; padding: 2rem 2.5rem; border-radius: 12px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); width: 100%; max-width: 400px; }
		h2 { margin-bottom: 1.5rem; text-align: center; }
		label { display: block; margin-bottom: 0.5rem; font-weight: bold; }
		input[type="email"], input[type="password"] { width: 100%; padding: 0.6rem; margin-bottom: 1rem; border: 1px solid #ccc; border-radius: 6px; box-sizing: border-box; }
		button { width: 100%; padding: 0.8rem; background-color: #007bff; color: white; border: none; border-radius: 6px; cursor: pointer; font-size: 1rem; }
		button:hover { background-color: #0056b3; }
	</style>
</head>
<body>
	<div class="login-container">
		<h2>Login to Continue</h2>
		<form method="POST" action="/api/v1/oauth2/login">
			` + hiddenFields.String() + `
			
			<label for="email">Email</label>
			<input type="email" name="email" id="email" required>

			<label for="password">Password</label>
			<input type="password" name="password" id="password" required minlength="6">

			<button type="submit">Sign In</button>
		</form>
	</div>
</body>
</html>
	`
}
