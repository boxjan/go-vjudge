{{ template "layout/app.tpl" . }}

{{ define "extant-header" }}

{{ end }}


{{ define "main-content" }}
    <div class="container-fluid row">
        <div class="col-1"></div>
        <div class="el-card col-10">
            <div class="el-card__header">
                <div class="row title justify-content-center"><h2> {{ "{{problem.title}}" }} </h2></div>
                <div class="row justify-content-center">
                    <div class="m-1" ><h6>Time Limit: {{"{{problem.time_limit}}"}}</h6></div>
                    <div class="m-1"><h6>Memory Limit: {{"{{problem.memory_limit}}"}}</h6></div>
                </div>
                <div class="row justify-content-center">
                    <el-button plain size="small" type="info" v-on:click="jumpSubmitPage()">Submit</el-button>
                    <el-button plain size="small" type="info" v-on:click="jumpStatusPage()">Status</el-button>
                </div>

            </div>

            <div class="el-card__body">
                 <div class="card m-3">
                     <div class="card-title"><h3 class="align-self-lg-auto">Description</h3></div>
                     <div class="card-body" v-html="problem.description">
                     </div>
                 </div>
                <div class="card m-3">
                    <div class="card-title"><h3>Input</h3></div>
                    <div class="card-body" v-html="problem.input">
                    </div>
                </div>
                <div class="card m-3">
                    <div class="card-title"><h3>Output</h3></div>
                    <div class="card-body" v-html="problem.output">
                    </div>
                </div>
                <div class="card m-3">
                    <div class="card-title"><h3>Sample</h3></div>
                    <div class="card-body" v-html="problem.sample">
                    </div>
                </div>
                <div class="card m-3">
                    <div class="card-title"><h3>Hint</h3></div>
                    <div class="card-body" v-html="problem.hint">
                    </div>
                </div>
                <div class="card m-3">
                    <div class="card-title"><h3>Source</h3></div>
                    <div class="card-body" v-html="problem.source">
                    </div>
                </div>
            </div>
        </div>
    </div>
{{ end }}


{{ define "extant-script" }}
    <script>
        const mixin = {
            data: function () {
                return {
                    problem: {{ marshal .problem }},
                }
            },
            methods: {
                jumpSubmitPage() {
                    window.open('{{ urlfor "SolutionController.UseInRoute" }}'+ '/' + this.problem.id, '_self');
                },
                jumpStatusPage() {
                    window.open('{{ urlfor "SolutionController.List" }}' +
                        '?ori_oj=' + this.problem.ori_oj +
                        '&ori_id=' + this.problem.ori_id, '_self');
                },
            }
        }
    </script>
{{ end }}