package oauth

import "html/template"

// loginFormData carries the OAuth request parameters through the login form so
// the POST can rebuild the authorization request after the user authenticates.
// ErrorMessage is non-empty when re-rendering after a failed login attempt.
type loginFormData struct {
	ClientId      string
	RedirectUri   string
	State         string
	Scope         string
	CodeChallenge string
	Resource      string
	ClientName    string
	ErrorMessage  string
}

// loginTemplate renders the minimal, dependency-free login page shown at the
// authorize endpoint. All OAuth parameters are carried as hidden fields and
// values are auto-escaped by html/template.
var loginTemplate = template.Must(template.New("login").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Sign in to Receipt Wrangler</title>
<style>
  body { font-family: system-ui, -apple-system, sans-serif; background: #f5f5f5; margin: 0; display: flex; min-height: 100vh; align-items: center; justify-content: center; }
  .card { background: #fff; padding: 2rem; border-radius: 8px; box-shadow: 0 1px 4px rgba(0,0,0,.15); width: 320px; }
  h1 { font-size: 1.15rem; margin: 0 0 1rem; }
  p.sub { color: #555; font-size: .85rem; margin: 0 0 1.25rem; }
  label { display: block; font-size: .8rem; color: #333; margin: .75rem 0 .25rem; }
  input[type=text], input[type=password] { width: 100%; padding: .55rem; border: 1px solid #ccc; border-radius: 4px; box-sizing: border-box; }
  button { margin-top: 1.25rem; width: 100%; padding: .6rem; background: #3f51b5; color: #fff; border: 0; border-radius: 4px; font-size: .95rem; cursor: pointer; }
  .error { background: #fdecea; color: #b71c1c; padding: .5rem .75rem; border-radius: 4px; font-size: .8rem; margin-bottom: 1rem; }
</style>
</head>
<body>
  <form class="card" method="post" action="/oauth/authorize">
    <h1>Sign in to Receipt Wrangler</h1>
    <p class="sub">{{if .ClientName}}<strong>{{.ClientName}}</strong> is requesting access to your receipts.{{else}}An application is requesting access to your receipts.{{end}}</p>
    {{if .ErrorMessage}}<div class="error">{{.ErrorMessage}}</div>{{end}}
    <label for="username">Username</label>
    <input id="username" name="username" type="text" autocomplete="username" autofocus required>
    <label for="password">Password</label>
    <input id="password" name="password" type="password" autocomplete="current-password" required>
    <input type="hidden" name="client_id" value="{{.ClientId}}">
    <input type="hidden" name="redirect_uri" value="{{.RedirectUri}}">
    <input type="hidden" name="state" value="{{.State}}">
    <input type="hidden" name="scope" value="{{.Scope}}">
    <input type="hidden" name="code_challenge" value="{{.CodeChallenge}}">
    <input type="hidden" name="resource" value="{{.Resource}}">
    <button type="submit">Authorize</button>
  </form>
</body>
</html>`))
