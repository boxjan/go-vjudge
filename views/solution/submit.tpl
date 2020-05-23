{{ template "layout/app.tpl" . }}

{{ define "extant-header" }}

{{ end }}

{{ define "main-content" }}
    <div class="container">
        <div class="el-card">
            <form method="post">
                {{ .csrf }}
                <input hidden="hidden" value="{{ .problem_id }}" name="problem_id">
                <div class="el-card__header">

                </div>
                <div class="el-card__body">
                    <div class="form-group row">
                        <label for="problem" class="col-md-3 col-form-label text-md-right">{{ "Problem" }}</label>
                        <div class="col-md-4">
                            <input id="problem" type="text" class="form-control" disabled
                                   value="{{ .problem_ori_oj }} - {{ .problem_title }}" required>
                        </div>
                    </div>
                    <div class="form-group row">
                        <label for="Lang" class="col-md-3 col-form-label text-md-right">{{ "Language" }}</label>
                        <div class="col-md-4">
                            <select id="problem" type="text" class="form-control" name="language" required>
                                {{ range $lang := .language }}
                                    {{ if eq $lang $.last_lang }}
                                    <option value="{{ $lang }}" selected> {{ $lang }} </option>
                                    {{ else }}
                                    <option value="{{ $lang }}"> {{ $lang }} </option>
                                    {{ end }}
                                {{ end }}

                            </select>
                        </div>
                    </div>
                    <div class="form-group row">
                        <label for="code" class="col-md-3 col-form-label text-md-right">{{ "Code" }}</label>
                        <div class="col-md-6">
                            {{ if map_get .error "code" }}
                            <textarea name="code" id="code" class="form-control is-invalid" style="height: 450px;max-height: 450px">{{ map_get .old "code" }}</textarea>
                            {{ else }}
                            <textarea name="code" id="code" class="form-control" style="height: 450px;max-height: 450px"></textarea>
                            {{ end }}
                            {{ if map_get .error "code" }}
                                <span class="invalid-feedback" role="alert">
                                        <strong>{{ map_get .error "identification" }}</strong>
                                    </span>
                            {{ end }}
                        </div>
                    </div>
                    <div class="form-group row mb-0">
                        <div class="col-md-8 offset-md-3">
                            <button type="submit" class="btn btn-primary">
                                Submit
                            </button>
                        </div>
                    </div>
                </div>
            </form>

        </div>
    </div>
{{ end }}

{{ define "extant-script" }}

{{ end }}