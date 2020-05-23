{{ template "layout/app.tpl" . }}

{{ define "extant-header" }}

{{ end }}

{{ define "main-content" }}
    <div class="container-fluid">
        <div class="el-card">
            <div class="el-card__header">
                <div class="title h2"> Solution List </div>
            </div>

            <div class="el-card__body">
                <form method="get">
                    <input hidden name="ori_oj" v-model="search.ori_oj">
                    <input hidden name="ori_id" v-model="search.ori_id">
                    <input hidden name="user" v-model="search.user">
                    <input hidden name="status" v-model="search.status">
                    <input hidden name="page" v-model="search.page">
                    <input type="submit" hidden class="form-control btn btn-primary" id="SubmitButton" value="Search">
                </form>
                <el-table :data="solutions.list" style="width: 100%" stripe>
                    <el-table-column prop="solution_id" label="Run ID" width="120">
                    </el-table-column>

                    <el-table-column prop="user" label="User" width="360">
                        <template slot="header" slot-scope="scope">
                            <el-input v-model="search.user" size="mini" placeholder="Username" @change="handleSearch()"/>
                        </template>
                    </el-table-column>

                    <el-table-column prop="ori_oj" label="OJ" width="120">
                        <template slot="header" slot-scope="scope">
                            <el-select v-model="search.ori_oj" placeholder="OJ" size="mini" clearable @change="handleSearch()">
                                <el-option v-for="item in supportOj" :key="item" :label="item.toUpperCase()" :value="item" ></el-option>
                            </el-select>
                        </template>
                    </el-table-column>

                    <el-table-column prop="ori_id" label="Problem ID" width="90">
                        <template slot="header" slot-scope="scope">
                            <el-input v-model="search.ori_id" size="mini" placeholder="Problem Id" @change="handleSearch()"/>
                        </template>
                        <template slot-scope="scope">
                            <el-button @click="handleOpenProblem(scope.row.problem_id)" type="text" size="small">
                                {{ "{{scope.row.ori_id}}" }}
                            </el-button>
                        </template>
                    </el-table-column>

                    <el-table-column prop="status" label="Status" width="270">
                        <template slot="header" slot-scope="scope">
                            <el-select :value="search.status" placeholder="Status" size="mini" clearable @change="handleSearch()">
                                <el-option v-for="item in status" :key="item.key" :label="item.value" :value="item.key"></el-option>
                            </el-select>
                        </template>
                        <template slot-scope="scope">
                            <el-tooltip  v-if="scope.row.status > 0" class="item" :content=scope.row.raw_status placement="bottom-start">
                                <el-button @click="handleGetErrorDetail(scope.row.solution_id)" type="text" size="small">
                                    {{ "{{ handlerGetStatusText(scope.row.status) }}" }}
                                </el-button>
                            </el-tooltip>
                            <el-button v-else @click="handleGetErrorDetail(scope.row.solution_id)" type="text" size="small">
                                {{ "{{ handlerGetStatusText(scope.row.status) }}" }}
                            </el-button>

                        </template>
                    </el-table-column>

                    <el-table-column prop="memory_used" label="Memory Used(kb)" width="150">
                    </el-table-column>

                    <el-table-column prop="time_used" label="Time Used(ms)" width="150">
                    </el-table-column>

                    <el-table-column prop="land" label="Language" width="120">
                       <template slot-scope="scope">
                           <el-button @click="handleGetCode(scope.row.solution_id)" type="text" size="small">
                               {{ "{{ scope.row.lang }}" }}
                           </el-button>
                        </template>
                    </el-table-column>

                    <el-table-column prop="length" label="length" width="120">
                    </el-table-column>

                    <el-table-column
                            prop="submit_time"
                            sortable
                            label="Update Time"
                            width="180">
                        <template slot-scope = "scope">
                            {{ "{{ timeFormat(scope.row.submit_time) }}" }}
                        </template>
                    </el-table-column>
                </el-table>

                <el-pagination
                        :page-size="15"
                        :pager-count="7"
                        :current-page.sync="search.page"
                        @current-change="handlePageChange()"
                        layout="prev, pager, next"
                        :total="solutions.tol">
                </el-pagination>
            </div>
        </div>
    </div>
{{ end }}

{{ define "extant-script" }}
    <script>
        const mixin = {
            data: function () {
                return {
                    supportOj: {{ marshal .supportOJ }},
                    solutions: {
                        tol: 0,
                        cur: 0
                    },
                    search: {{ marshal .search }},
                    solutionsApi: {{ urlfor "SolutionController.ApiList" }},
                    willUpdate: 0,
                    status: [
                        {key: -5, value: "waiting"},
                        {key: -4, value: "pending"},
                        {key: -3, value: "queuing"},
                        {key: -2, value: "compiling"},
                        {key: -1, value: "running"},
                        {key: 0, value: "accept"},
                        {key: 1, value: "presentation error"},
                        {key: 2, value: "wrong answer"},
                        {key: 3, value: "time limit exceeded"},
                        {key: 4, value: "memory limit exceed"},
                        {key: 5, value: "output limit exceed"},
                        {key: 6, value: "runtime error"},
                        {key: 7, value: "compile error"},
                        {key: 8, value: "submit error"},
                        {key: 9, value: "system error"},
                        {key: 10, value: "other failed"},
                    ]
                }
            },
            methods: {
                handleGetStatus() {
                    axios.get(this.solutionsApi + window.location.search).then(
                        (response) => {
                            console.log(response.data);
                            this.solutions = response.data
                        }
                    ).catch((error) => {
                        this.$message({
                            message: 'Get Data fail',
                            type: 'error',
                        });
                        console.log(error);
                    });
                    setTimeout(this.handleGetStatus, 3000);
                },
                handleGetCode(solutionId) {
                    axios.get(this.solutionsApi + '/' + solutionId + '/code').then(
                        (response) => {
                            console.log(response.data);
                        }
                    ).catch((error) => {
                        console.log(error);
                    });
                },


                handleGetErrorDetail(solutionId) {
                    axios.get(this.solutionsApi + '/' + solutionId + '/error').then(
                        (response) => {
                            console.log(response.data);
                            this.$alert('<pre>' + response.data + '</pre>', '', {
                                dangerouslyUseHTMLString: true
                            });
                        }
                    ).catch((error) => {
                        console.log(error);
                    });
                },

                handleOpenProblem(problemId) {
                    window.open('{{ urlfor "ProblemController.List" }}' + '/' + problemId + '/', '');
                },

                handlerGetStatusText(status) {
                    switch (status) {
                        case -6: return "starting";
                        case -5: return "waiting";
                        case -4: return "pending";
                        case -3: return "queuing";
                        case -2: return "compiling";
                        case -1: return "running";
                        case 0: return "accept";
                        case 1: return "presentation error";
                        case 2: return "wrong answer";
                        case 3: return "time limit exceeded";
                        case 4: return "memory limit exceed";
                        case 5: return "output limit exceed";
                        case 6: return "runtime error";
                        case 7: return "compile error";
                        case 8: return "submit error";
                        case 9: return "system error";
                        case 10: return "other failed";
                    }
                },
                timeFormat(string) {
                    const time = moment(string);
                    return time.format("lll");
                },

                async handleSearch() {
                    console.log(this.search);
                    const s = ++this.willUpdate;
                    await sleep(160).then(() => {
                        if (this.willUpdate === s) {
                            $("#SubmitButton").click();
                        }
                    });
                },
                handlePageChange() {
                    console.log(this.search);
                    $("#SubmitButton").click();
                },
            },

            mounted() {
                this.handleGetStatus();
            }
        }
    </script>
{{ end }}