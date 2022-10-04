<template>
  <div class="container">
    <el-card class="container-card" shadow="always">
<!--顶部功能区-->
      <el-form size="mini" :inline="true" :model="params" class="demo-form-inline">
        <el-form-item label="姓名">
          <el-input v-model.trim="params.name" clearable placeholder="姓名" @clear="search" />
        </el-form-item>
        <el-form-item label="年龄">
          <el-input v-model.trim="params.age" clearable placeholder="年龄" @clear="search" />
        </el-form-item>
        <el-form-item label="性别">
          <el-select v-model.trim="params.sex" clearable placeholder="性别" @change="search" @clear="search">
            <el-option label="男" value="1" />
            <el-option label="女" value="2" />
          </el-select>
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model.trim="params.mobile" clearable placeholder="手机号" @clear="search" />
        </el-form-item>
        <el-form-item>
          <el-button :loading="loading" icon="el-icon-search" type="success" @click="search">查询</el-button>
        </el-form-item>
        <el-form-item>
          <el-button :loading="loading" icon="el-icon-plus" type="primary" @click="create">新增</el-button>
        </el-form-item>
        <el-form-item>
          <el-button :disabled="multipleSelection.length === 0" :loading="loading" icon="el-icon-delete" type="danger" @click="batchDelete">批量删除</el-button>
        </el-form-item>
      </el-form>
<!--表格功能区-->
      <el-table v-loading="loading" :data="tableData" border  style="width: 100%" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" align="center" />
        <el-table-column show-overflow-tooltip sortable prop="ID" label="ID" />
        <el-table-column show-overflow-tooltip sortable prop="name" label="姓名" />
        <el-table-column show-overflow-tooltip sortable prop="sex" label="性别" align="center">
          <template slot-scope="scope">
            <el-tag size="small" >{{ scope.row.sex === 2 ? "女" : "男"}}</el-tag>
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip sortable prop="age" label="年龄" />
        <el-table-column show-overflow-tooltip sortable prop="mobile" label="手机号" />
        <el-table-column show-overflow-tooltip sortable prop="description" label="描述" />
        <el-table-column fixed="right" label="操作" align="center" width="250">
          <template slot-scope="scope">
            <el-tooltip content="编辑" effect="dark" placement="top">
              <el-button size="mini" icon="el-icon-edit" type="text" @click="update(scope.row)">编辑</el-button>
            </el-tooltip>
            <el-tooltip class="delete-popover" content="删除" effect="dark" placement="top">
              <el-popconfirm title="确定删除吗？" @onConfirm="singleDelete(scope.row.ID)">
                <el-button slot="reference" size="mini" icon="el-icon-delete" type="text">删除</el-button>
              </el-popconfirm>
            </el-tooltip>
          </template>
        </el-table-column>
      </el-table>

      <!--分页功能区-->
      <el-pagination
        :current-page="params.page"
        :page-size="params.size"
        :total="total"
        :page-sizes="[1, 5, 10, 30]"
        layout="total, prev, pager, next, sizes"
        background
        style="margin-top: 10px;float:right;margin-bottom: 10px;"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
<!--弹框功能区-->
      <el-dialog :title="dialogFormTitle" :visible.sync="dialogFormVisible" width="30%">
        <el-form ref="dialogForm" size="small" :model="dialogFormData" :rules="dialogFormRules" label-width="100px">

          <el-form-item label="姓名" prop="name">
            <el-input ref="password" v-model.trim="dialogFormData.name" placeholder="姓名" />
          </el-form-item>

          <el-form-item label="性别" prop="sex">
            <el-select v-model.trim="dialogFormData.sex" placeholder="请选择性别" style="width:100%">
              <el-option label="男" :value="1" />
              <el-option label="女" :value="2" />
            </el-select>
          </el-form-item>

          <el-form-item label="年龄" prop="age">
            <el-input v-model.trim="dialogFormData.age" placeholder="年龄" />
          </el-form-item>

          <el-form-item label="手机号" prop="mobile">
            <el-input v-model.trim="dialogFormData.mobile" placeholder="手机号" />
          </el-form-item>

          <el-form-item label="说明" prop="description">
            <el-input v-model.trim="dialogFormData.description" type="textarea" placeholder="说明" show-word-limit maxlength="100" />
          </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer">
          <el-button size="mini" @click="cancelForm()">取 消</el-button>
          <el-button size="mini" :loading="submitLoading" type="primary" @click="submitForm()">确 定</el-button>
        </div>
      </el-dialog>
    </el-card>
  </div>
</template>

<script>
import { listExample,removeExample,createExample,putExample } from '@/api/example'


export default {
  name: 'User',
  data() {
    let checkPhone = (rule, value, callback) => {
      if (!value) {
        return callback(new Error('手机号不能为空'))
      } else {
        const reg = /^1[3|4|5|7|8][0-9]\d{8}$/
        if (reg.test(value)) {
          callback()
        } else {
          return callback(new Error('请输入正确的手机号'))
        }
      }
    }
    return {
      // 查询参数
      params: {
        page: 1,
        size: 10
      },
      // 表格数据
      tableData: [],
      total: 0,
      loading: false,

      // dialog对话框
      submitLoading: false,
      dialogFormTitle: '',
      dialogType: '',
      dialogFormVisible: false,
      dialogFormData: {
        name: '',
        age: 0,
        sex: 1,
        mobile: '',
        description: ''
      },
      dialogFormRules: {
        name: [
          { required: true, message: '请输入姓名', trigger: 'blur' },
          { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
        ],
        age: [
          { required: true, message: '请输入年龄', trigger: 'blur' }
        ],
        mobile: [
          { required: true, validator: checkPhone, trigger: 'blur' }
        ],
        sex: [
          { required: true, message: '请选择性别', trigger: 'change' }
        ],
        description: [
          { required: false, message: '说明', trigger: 'blur' },
          { min: 0, max: 100, message: '长度在 0 到 100 个字符', trigger: 'blur' }
        ]
      },

      // 删除按钮弹出框
      popoverVisible: false,
      // 表格多选
      multipleSelection: []
    }
  },
  created() {
    this.getTableData()
  },
  methods: {
    // 查询
    search() {
      this.params.page = 1
      this.getTableData()
    },

    // 获取表格数据
    getTableData() {
      this.loading = true

      listExample(this.params).then(res => {
        this.tableData = res.data.data
        this.total = res.data.total
        this.loading = false

      })
    },
    // 新增
    create() {

      this.dialogFormData.name = ""
      this.dialogFormData.age =  0
      this.dialogFormData.sex =  1
      this.dialogFormData.mobile =  ""
      this.dialogFormData.description = ""
      this.dialogFormTitle = '新增用户'
      this.dialogType = 'create'
      this.dialogFormVisible = true
    },

    // 修改
    update(row) {
      this.dialogFormData.ID = row.ID
      this.dialogFormData.name = row.name
      this.dialogFormData.age = row.age
      this.dialogFormData.sex = row.sex
      this.dialogFormData.mobile = row.mobile
      this.dialogFormData.description = row.description

      this.dialogFormTitle = '修改用户'
      this.dialogType = 'update'
      this.dialogFormVisible = true

    },

    // 提交表单
    submitForm() {
      this.$refs['dialogForm'].validate(async valid => {
        if (valid) {
          this.submitLoading = true
          let msg = ''
          try {
            if (this.dialogType === 'create') {
              let data = {
                name: this.dialogFormData.name ,
                age: this.dialogFormData.age ,
                sex:this.dialogFormData.sex ,
                mobile:this.dialogFormData.mobile ,
                description:this.dialogFormData.description
              }
              const { message } = await createExample(data)
              msg = message
            } else {
              const { message } = await putExample(this.dialogFormData)
              msg = message
            }
          } finally {
            this.submitLoading = false
          }

          this.resetForm()
          this.getTableData()
          this.$message({
            showClose: true,
            message: msg,
            type: 'success'
          })
        } else {
          this.$message({
            showClose: true,
            message: '表单校验失败',
            type: 'error'
          })
          return false
        }
      })
    },

    // 提交表单
    cancelForm() {
      this.resetForm()
    },

    resetForm() {
      this.dialogFormVisible = false
      this.$refs['dialogForm'].resetFields()
      this.dialogFormData = {
        name: '',
        age: 0,
        sex: 1,
        mobile: '',
        description: ''
      }
    },

    // 批量删除
    batchDelete() {
      this.$confirm('此操作将永久删除, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        this.loading = true
        const exampleIds = []
        this.multipleSelection.forEach(x => {
          exampleIds.push(x.ID)
        })
        let msg = ''
        try {
          const { message } = await removeExample({ exampleIds: exampleIds })
          msg = message
        } finally {
          this.loading = false
        }

        this.getTableData()
        this.$message({
          showClose: true,
          message: msg,
          type: 'success'
        })
      }).catch(() => {
        this.$message({
          showClose: true,
          type: 'info',
          message: '已取消删除'
        })
      })
    },

    // 表格多选
    handleSelectionChange(val) {
      this.multipleSelection = val
    },

    // 单个删除
    async singleDelete(Id) {
      this.loading = true
      let msg = ''
      try {
        const { message } = await removeExample({ exampleIds: [Id] })
        msg = message
      } finally {
        this.loading = false
      }

      this.getTableData()
      this.$message({
        showClose: true,
        message: msg,
        type: 'success'
      })
    },


    // 分页
    handleSizeChange(val) {
      this.params.size = val
      this.getTableData()
    },
    handleCurrentChange(val) {
      this.params.page = val
      this.getTableData()
    }
  }
}
</script>

<style scoped>
.container-card{
  margin: 10px;
}

.delete-popover{
  margin-left: 10px;
}

</style>
