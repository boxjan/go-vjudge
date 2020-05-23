<nav class="navbar navbar-expand-md navbar-light navbar-laravel">
    <div class="container-fluid">
        <a class="navbar-brand" href="{{ urlfor "IndexController.Index" }}">
            {{ config "string" "appname" "go-vjudge" }}
        </a>
        <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarSupportedContent"
                aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="{{ "Toggle navigation" }}">
            <span class="navbar-toggler-icon"></span>
        </button>

        <div class="collapse navbar-collapse" id="navbarSupportedContent">
            <!-- Left Side Of Navbar -->
            <ul class="navbar-nav mr-auto">
                <li class="nav-item">
                    <a class="nav-link" id="nav-statistics" href="{{ urlfor "ProblemController.List" }}">
                        Problems
                    </a>
                </li>

                <li class="nav-item">
                    <a class="nav-link" id="nav-statistics" href="{{ urlfor "SolutionController.List" }}">
                        Solutions
                    </a>
                </li>
            </ul>

            <!-- Right Side Of Navbar -->
            <ul class="navbar-nav ml-auto">
                <!-- Authentication Links -->
                {{ if $.user }}
                <li class="nav-item dropdown">
                    <a id="navbarDropdown" class="nav-link dropdown-toggle" href="{{ urlfor "UserController.InfoPage" }}"
                       role="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false" v-pre>
                        Hello,
                        {{if .user.Nickname}}
                            {{ .user.Nickname }}
                        {{else}}
                            {{ .user.Username }}
                        {{end}}
                        <span class="caret"></span>
                    </a>

                    <div class="dropdown-menu dropdown-menu-right" aria-labelledby="navbarDropdown">
                        <a class="dropdown-item" href="{{ urlfor "UserController.ResetPasswordPage" }}">
                            {{ "Update Password" }}
                        </a>

                        <a class="dropdown-item" href="{{ urlfor "LoginController.Logout" }}"
                           onclick="document.getElementById('logout-form').submit();">
                            {{ "Logout" }}
                        </a>
                    </div>

                    <div class="dropdown-menu dropdown-menu-right" aria-labelledby="logout-button">


                        <form id="logout-form" action="{{ urlfor "AuthController.Logout" }}" method="POST" style="display: none;">
                            {{ .csrf }}
                        </form>
                    </div>
                </li>
                {{ else }}
                <li class="nav-item">
                    <a class="nav-link" href="{{ urlfor "LoginController.LoginPage" }}">Login</a>
                </li>
                <li class="nav-item">
                    <a class="nav-link" href="{{ urlfor "RegisterController.RegisterPage" }}">Register</a>
                </li>
                {{ end }}
            </ul>
        </div>
    </div>
</nav>