<template>
  <div class="container">
    <div class="tabor">
      <el-form size="mini" :inline="true" :model="params" class="demo-form-inline">
        <el-form-item label="访问路径">
          <el-input v-model.trim="params.path" clearable placeholder="访问路径" @clear="search" />
        </el-form-item>
        <el-form-item label="所属类别">
          <el-input v-model.trim="params.category" clearable placeholder="所属类别" @clear="search" />
        </el-form-item>
        <el-form-item label="请求方法">
          <el-select v-model.trim="params.method" clearable placeholder="请求方式" @change="search" @clear="search">
            <el-option label="GET[获取资源]" value="GET" />
            <el-option label="POST[新增资源]" value="POST" />
            <el-option label="PUT[全部更新]" value="PUT" />
            <el-option label="PATCH[增量更新]" value="PATCH" />
            <el-option label="DELETE[删除资源]" value="DELETE" />
          </el-select>
        </el-form-item>
        <el-form-item label="创建人">
          <el-input v-model.trim="params.creator" clearable placeholder="创建人" @clear="search" />
        </el-form-item>
        <el-form-item>
          <Button :but="Select" @but="search"/>
          <Button :but="Add" @but="create"/>
          <Button :but="DeleteAll" @but="batchDelete"/>
          <Button :but="Refresh" @but="refresh()"/>
        </el-form-item>
      </el-form>

    </div>

    <div class="table-box"  >
      <!--
      @author:风很大
      @description: 表格数据
      @time: 2021/12/22 0022
       -->
      <Table
        :table="table"
        :pagination="pagination"
        :setting="setting"
        @select="handleSelectionChange"
        @page="handleCurrentChange"
      >
        <el-table-column show-overflow-tooltip sortable prop="path" label="访问路径" />
        <el-table-column show-overflow-tooltip sortable prop="category" label="所属类别" />
        <el-table-column show-overflow-tooltip sortable prop="method" label="请求方式">
          <template slot-scope="scope">
            <el-tag size="small" :type="scope.row.method | methodTagFilter" disable-transitions>{{ scope.row.method }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip sortable prop="creator" label="创建人" />
        <el-table-column show-overflow-tooltip sortable prop="desc" label="说明" />
        <el-table-column fixed="right" label="操作" align="center" width="250">
          <template slot-scope="scope">
            <Button :but="Edit" @but="update(scope.row)"/>
            <Button :but="Delete" @but="singleDelete(scope.row.ID)"/>
          </template>
        </el-table-column>
      </Table>
      <!--
      @author:风很大
      @description: 分页组件
      @time: 2022/1/17 0017
      -->
      <Pagination :pagination="pagination" @page="handleCurrentChange" @size="handleSizeChange"/>
    </div>


    <el-dialog :title="dialogFormTitle" :visible.sync="dialogFormVisible">
        <el-form ref="dialogForm" size="small" :model="dialogFormData" :rules="dialogFormRules" label-width="120px">
          <el-form-item label="访问路径" prop="path">
            <el-input v-model.trim="dialogFormData.path" placeholder="访问路径" />
          </el-form-item>
          <el-form-item label="所属类别" prop="category">
            <el-input v-model.trim="dialogFormData.category" placeholder="所属类别" />
          </el-form-item>
          <el-form-item label="请求方式" prop="method">
            <el-select v-model.trim="dialogFormData.method" placeholder="请选择请求方式">
              <el-option label="GET[获取资源]" value="GET" />
              <el-option label="POST[新增资源]" value="POST" />
              <el-option label="PUT[全部更新]" value="PUT" />
              <el-option label="PATCH[增量更新]" value="PATCH" />
              <el-option label="DELETE[删除资源]" value="DELETE" />
            </el-select>
          </el-form-item>
          <el-form-item label="说明" prop="desc">
            <el-input v-model.trim="dialogFormData.desc" type="textarea" placeholder="说明" show-word-limit maxlength="100" />
          </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer">
          <el-button size="mini" @click="cancelForm()">取 消</el-button>
          <el-button size="mini" :loading="submitLoading" type="primary" @click="submitForm()">确 定</el-button>
        </div>
      </el-dialog>

  </div>
</template>

<script>
import { getApis, createApi, updateApiById, batchDeleteApiByIds } from '@/api/system/api'
import Table from '@/components/Table'
import Dialog from '@/components/Dialog'
import Button from '@/components/Button'
import Pagination from '@/components/Pagination'

export default {
  name: 'Index',
  inject: ['reload'],
  components: { Button, Table, Dialog, Pagination },
  filters: {
    methodTagFilter(val) {
      if (val === 'GET') {
        return ''
      } else if (val === 'POST') {
        return 'success'
      } else if (val === 'PUT') {
        return 'info'
      } else if (val === 'PATCH') {
        return 'warning'
      } else if (val === 'DELETE') {
        return 'danger'
      } else {
        return 'info'
      }
    }
  },
  data() {
    return {
      // 查询参数
      params: {
        path: '',
        method: '',
        category: '',
        creator: '',
        pageNum: 1,
        pageSize: 10
      },
      // 按钮配置
      Add: {
        name: '新增',
        size: 'mini',
        type: 'primary',
        icon: 'el-icon-plus',
        plain: false,
        disabled: false,
        show: true
      },
      Select: {
        name: '查询',
        size: 'mini',
        type: 'success',
        icon: 'el-icon-search',
        plain: false,
        disabled: false,
        show: true
      },
      Detail: {
        name: '详情',
        size: 'mini',
        type: 'text',
        icon: 'el-icon-view',
        plain: false,
        disabled: false,
        show: true
      },
      Edit: {
        name: '编辑',
        size: 'mini',
        type: 'text',
        icon: 'el-icon-edit',
        plain: false,
        disabled: false,
        show: true
      },
      Delete: {
        name: '删除',
        size: 'mini',
        type: 'text',
        icon: 'el-icon-delete',
        plain: false,
        disabled: false,
        show: true
      },
      DeleteAll: {
        name: '批量删除',
        size: 'mini',
        type: 'danger',
        icon: 'el-icon-delete',
        plain: false,
        disabled: true,
        show: true
      },
      Refresh: {
        name: '刷新',
        size: 'mini',
        type: 'warning',
        icon: 'el-icon-refresh',
        circle: false,
        plain: false,
        disabled: false,
        show: true
      },
      Import: {
        name: '导入',
        size: 'mini',
        type: 'info',
        icon: 'el-icon-upload2',
        plain: true,
        disabled: false,
        show: true
      },
      Export: {
        name: '导出',
        size: 'mini',
        type: 'warning',
        icon: 'el-icon-download',
        plain: true,
        disabled: false,
        show: true

      },
      // 弹窗配置
      AddDialog: {
        title: '新增',
        dialog: false,
        width: '600px'
      },
      EditDialog: {
        title: '编辑',
        dialog: false,
        width: '600px'
      },
      DetailDialog: {
        title: '详情',
        dialog: false,
        width: '600px'
      },
      // 表格和分页配置
      pagination: {
        page: 1,
        size: 10,
        total: 0
      },
      table: [],
      setting: {
        checkbox: true,
        order: false,
        loading: false
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
        path: '',
        category: '',
        method: '',
        desc: ''
      },
      dialogFormRules: {
        path: [
          { required: true, message: '请输入访问路径', trigger: 'blur' },
          { min: 1, max: 100, message: '长度在 1 到 100 个字符', trigger: 'blur' }
        ],
        category: [
          { required: true, message: '请输入所属类别', trigger: 'blur' },
          { min: 1, max: 50, message: '长度在 1 到 50 个字符', trigger: 'blur' }
        ],
        method: [
          { required: true, message: '请选择请求方式', trigger: 'change' }
        ],
        desc: [
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
    // 刷新页面
    refresh() {
      this.reload()
    },
    // 查询
    search() {
      this.params.pageNum = 1
      this.getTableData()
    },

    // 获取表格数据
    async getTableData() {
      this.loading = true
      try {
        const { data } = await getApis(this.params)
        this.table = data.data
        this.pagination.total = data.total
      } finally {
        this.loading = false
      }
    },

    // 新增
    create() {
      this.dialogFormTitle = '新增接口'
      this.dialogType = 'create'
      this.dialogFormVisible = true
    },

    // 修改
    update(row) {
      this.dialogFormData.ID = row.ID
      this.dialogFormData.path = row.path
      this.dialogFormData.category = row.category
      this.dialogFormData.method = row.method
      this.dialogFormData.desc = row.desc

      this.dialogFormTitle = '修改接口'
      this.dialogType = 'update'
      this.dialogFormVisible = true
    },

    // 提交表单
    submitForm() {
      this.$refs['dialogForm'].validate(async valid => {
        if (valid) {
          let message = ''
          this.submitLoading = true
          try {
            if (this.dialogType === 'create') {
              const { msg } = await createApi(this.dialogFormData)
              message = msg
            } else {
              const { msg } = await updateApiById(this.dialogFormData.ID, this.dialogFormData)
              message = msg
            }
          } finally {
            this.submitLoading = false
          }

          this.resetForm()
          await this.getTableData()
          this.$message({
            showClose: true,
            message: message,
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
        path: '',
        category: '',
        method: '',
        desc: ''
      }
    },

    // 批量删除
    batchDelete() {
      this.$confirm('此操作将永久删除, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async res => {
        this.loading = true
        const apiIds = []
        this.multipleSelection.forEach(x => {
          apiIds.push(x.ID)
        })
        let message = ''
        try {
          const { msg } = await batchDeleteApiByIds({ apiIds: apiIds })
          message = msg
        } finally {
          this.loading = false
        }

        await this.getTableData()
        this.$message({
          showClose: true,
          message: message,
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
      this.DeleteAll.disabled = this.multipleSelection.length === 0

    },

    // 单个删除
    async singleDelete(Id) {
      this.$confirm('确认删除？')
        .then(async() => {
          this.loading = true
          let message = ''
          try {
            const { msg } = await batchDeleteApiByIds({ apiIds: [Id] })
            message = msg
          } finally {
            this.loading = false
          }

          await this.getTableData()
          this.$message({
            showClose: true,
            message: message,
            type: 'success'
          })
        })
        .catch(_ => {
          this.$message({
            showClose: true,
            type: 'info',
            message: '已取消删除'
          })
        })

    },

    // 分页
    handleSizeChange(val) {
      this.params.pageSize = val
      this.getTableData()
    },
    handleCurrentChange(val) {
      this.params.pageNum = val
      this.getTableData()
    }
  }
}
</script>

<style scoped>
.table-box {
  background-color: #ffffff;
  padding: 15px 10px;
}

  .delete-popover{
    margin-left: 10px;
  }
</style>
