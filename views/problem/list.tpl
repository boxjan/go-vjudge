{{ template "layout/app.tpl" . }}

{{ define "extant-header" }}

{{ end }}

{{ define "main-content" }}
    <div class="container-fluid">
        <div class="el-card">
            <div class="el-card__header">
                <div class="title h2"> Problem List </div>
            </div>

            <div class="el-card__body">
                <form method="get">
                    <input hidden name="ori_oj" v-model="search.ori_oj">
                    <input hidden name="ori_id" v-model="search.ori_id">
                    <input hidden name="title" v-model="search.title">
                    <input hidden name="source" v-model="search.source">
                    <input hidden name="page" v-model="search.page">
                    <input type="submit" hidden class="form-control btn btn-primary" id="SubmitButton" value="Search">
                </form>
                <el-table :data="problems.list" style="width: 100%" stripe>
                    <el-table-column prop="ori_oj" label="OJ" width="120">
                        <template slot="header" slot-scope="scope">
                            <el-select v-model="search.ori_oj" placeholder="OJ" size="mini" clearable @change="handleSearch()">
                                <el-option v-for="item in supportOj" :key="item" :label="item.toUpperCase()" :value="item" ></el-option>
                            </el-select>
                        </template>
                    </el-table-column>
                    <el-table-column
                            prop="ori_id"
                            label="Problem ID"
                            width="90">
                        <template slot="header" slot-scope="scope">
                            <el-input v-model="search.ori_id" size="mini" placeholder="Problem Id" @change="handleSearch()"/>
                        </template>
                    </el-table-column>
                    <el-table-column
                            prop="title"
                            label="Title">
                        <template slot="header" slot-scope="scope">
                            <el-input v-model="search.title" size="mini" placeholder="Title" @change="handleSearch()"/>
                        </template>
                    </el-table-column>
                    <el-table-column
                            prop="update_at"
                            sortable
                            label="Update Time"
                            width="180">
                        <template slot-scope = "scope">
                            {{ "{{ timeFormat(scope.row.update_at) }}" }}
                        </template>
                    </el-table-column>
                    <el-table-column
                            prop="source"
                            label="Source">
                        <template slot="header" slot-scope="scope">
                            <el-input v-model="search.source" size="mini" placeholder="Source"/>
                        </template>
                    </el-table-column>
                    <el-table-column
                            fixed="right"
                            label="Do It"
                            width="90">
                        <template slot-scope="scope">
                            <el-button @click="handleDoIt(scope.row.id)" type="text" size="small"> Do It! </el-button>
                        </template>
                    </el-table-column>
                </el-table>

                <el-pagination
                  :page-size="15"
                  :pager-count="7"
                  :current-page.sync="search.page"
                  @current-change="handlePageChange()"
                  layout="prev, pager, next"
                  :total="problems.tol">
                </el-pagination>
            </div>
        </div>
    </div>
{{ end }}

{{ define "extant-script" }}
    <script>
        const mixin = {
            data : function () {
                return {
                    supportOj: {{ marshal .supportOJ }},
                    problems: {},
                    search: {{ marshal .search }},
                    problemApiUrl: {{ urlfor "ProblemController.ApiList" }},
                    willUpdate: 0,
                }
            },
            methods: {
                handleGetProblems() {
                    axios.get(this.problemApiUrl + window.location.search).then(
                        (response) => {
                            console.log(response.data);
                            this.problems = response.data
                        }
                    ).catch((error) => {
                        this.$message({
                            message: 'Get Data fail',
                            type: 'error',
                        });
                        console.log(error);
                    })
                },

                noDataRefresh() {
                    if (this.problems.list == null) {
                        this.handleGetProblems();
                        setTimeout(this.noDataRefresh, 3000)
                    }

                },

                timeFormat(string) {
                    const time = moment(string);
                    return time.format("LL");
                },

                handlePageChange() {
                    console.log(this.search);
                    $("#SubmitButton").click();
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

                handleDoIt(id) {
                    window.open('{{ urlfor "ProblemController.List" }}' + '/' + id + '/', '');
                },


            },

            mounted() {
                this.handleGetProblems();
                if (this.problems.list == null) {
                    setTimeout(this.noDataRefresh, 1000);
                }
            }
        }
    </script>
{{ end }}
