{{ template "layout/app.tpl" . }}

{{ define "main-content" }}
    <div class="container">
        <div class="row justify-content-center">
            <div class="col-md-8">
                <div class="card">
                    <div class="card-header">Register</div>

                    <div class="card-body">
                        <form method="POST" action="{{ urlfor "RegisterController.RegisterPage" }}">
                            {{ .csrf }}

                            <div class="form-group row">
                                <label for="username" class="col-md-4 col-form-label text-md-right">Username</label>

                                <div class="col-md-6">
                                    {{ if map_get .error "username" }}
                                        <input id="username" type="text" class="form-control is-invalid"
                                               name="username" value="{{ map_get .old "username" }}" required autofocus>
                                        <span class="invalid-feedback" role="alert">
                                        <strong>{{ map_get .error "username" }}</strong>
                                    </span>
                                    {{ else }}
                                        <input id="username" type="text" class="form-control"
                                               name="username" value="{{ map_get .old "username" }}" required autofocus>
                                    {{ end }}
                                </div>
                            </div>

                            <div class="form-group row">
                                <label for="nickname" class="col-md-4 col-form-label text-md-right"> Nickname </label>

                                <div class="col-md-6">
                                    {{ if map_get .error "nickname" }}
                                        <input id="nickname" type="text" class="form-control is-invalid"
                                               name="nickname" value="{{ map_get .old "nickname" }}" required autofocus>
                                        <span class="invalid-feedback" role="alert">
                                        <strong>{{ map_get .error "nickname" }}</strong>
                                    </span>
                                    {{ else }}
                                        <input id="nickname" type="text" class="form-control"
                                               name="nickname" value="{{ map_get .old "nickname" }}" required autofocus>
                                    {{ end }}
                                </div>
                            </div>

                            <div class="form-group row">
                                <label for="email" class="col-md-4 col-form-label text-md-right"> E-Mail Address </label>
                                <div class="col-md-6">

                                    {{ if map_get .error "email" }}
                                        <input id="email" type="text" class="form-control is-invalid"
                                               name="email" value="{{ map_get .old "email" }}" required autofocus>
                                        <span class="invalid-feedback" role="alert">
                                        <strong>{{ map_get .error "email" }}</strong>
                                    </span>
                                    {{ else }}
                                        <input id="email" type="text" class="form-control"
                                               name="email" value="{{ map_get .old "email" }}" required autofocus>
                                    {{ end }}

                                </div>
                            </div>

                            <div class="form-group row">
                                <label for="password" class="col-md-4 col-form-label text-md-right"> Password </label>

                                <div class="col-md-6">
                                    {{ if map_get .error "password" }}
                                        <input id="password" type="password" class="form-control is-invalid"
                                               name="password" required autofocus>

                                        <span class="invalid-feedback" role="alert">
                                        <strong>{{ map_get .error "password" }}</strong>
                                    </span>
                                    {{ else }}
                                        <input id="password" type="password" class="form-control"
                                               name="password" required autofocus>
                                    {{ end }}
                                </div>
                            </div>

                            <div class="form-group row">
                                <label for="password-confirm" class="col-md-4 col-form-label text-md-right"> Confirm Password </label>

                                <div class="col-md-6">
                                    <input id="password-confirm" type="password" class="form-control" name="password_confirmation" required>
                                </div>
                            </div>

                            <div class="form-group row">
                                <label for="school" class="col-md-4 col-form-label text-md-right">School</label>

                                <div class="col-md-6">
                                    {{ if map_get .error "school" }}
                                        <input id="school" type="text" class="form-control is-invalid"
                                               name="school" value="{{ map_get .old "school" }}" required autofocus>
                                        <span class="invalid-feedback" role="alert">
                                        <strong>{{ map_get .error "school" }}</strong>
                                    </span>
                                    {{ else }}
                                        <input id="school" type="text" class="form-control"
                                               name="school" value="{{ map_get .old "school" }}" required autofocus>
                                    {{ end }}
                                </div>
                            </div>

                            <div class="form-group row mb-0">
                                <div class="col-md-6 offset-md-4">
                                    <button type="submit" class="btn btn-primary">
                                        {{ "Register" }}
                                    </button>
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