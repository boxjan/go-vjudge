{{ template "layout/app.tpl" . }}

{{ define "main-content" }}
<div class="container">
    <div class="row justify-content-center">
        <div class="col-md-8">
            <div class="card">
                <div class="card-header">
                    <h4>Login</h4>
                    {{ if .url }}
                    <p> {{ "You must login first!" }} </p>
                    {{ end }}
                </div>

                <div class="card-body">
                    <form method="POST" action="{{ urlfor "AuthController.Login" }}">
                        {{ .csrf }}

                        <div class="form-group row">
                            <label for="identification" class="col-md-4 col-form-label text-md-right">{{ "Username" }}</label>

                            <div class="col-md-6">
                                {{ if map_get .error "identification" }}
                                    <input id="identification" type="text" class="form-control is-invalid" name="identification"
                                           value="{{ map_get .old "identification" }}" required autofocus>
                                {{ else }}
                                    <input id="identification" type="text" class="form-control" name="identification" required autofocus>
                                {{ end }}


                                {{ if map_get .error "identification" }}
                                <span class="invalid-feedback" role="alert">
                                        <strong>{{ map_get .error "identification" }}</strong>
                                    </span>
                                {{ end }}
                            </div>
                        </div>

                        <div class="form-group row">
                            <label for="password" class="col-md-4 col-form-label text-md-right">{{ "Password" }}</label>

                            <div class="col-md-6">
                                <input id="password" type="password" class="form-control" name="password" required>
                            </div>
                        </div>

{{/*                        <div class="form-group row">*/}}
{{/*                            <div class="col-md-6 offset-md-4">*/}}
{{/*                                <div class="form-check">--}}*/}}
{{/*                                    <input class="form-check-input" type="checkbox" name="remember" id="remember" {{ old('remember') ? 'checked' : '' }}>*/}}

{{/*                                    <label class="form-check-label" for="remember">*/}}
{{/*                                        {{ __('Remember Me') }}*/}}
{{/*                                    </label>*/}}
{{/*                                </div>*/}}
{{/*                            </div>*/}}
{{/*                        </div>*/}}

                        <div class="form-group row mb-0">
                            <div class="col-md-8 offset-md-4">
                                <button type="submit" class="btn btn-primary">
                                    {{ "Login" }}
                                </button>

{{/*                                <a class="btn btn-link" href="{{ urlfor ".AuthController.request" }}">*/}}
{{/*                                    {{ "Forgot Your Password?" }}*/}}
{{/*                                </a>*/}}
                            </div>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    </div>
</div>
{{ end }}

{{ define "extant-header" }}
{{ end }}

{{ define "extant-script" }}
{{ end }}