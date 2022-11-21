<template>
  <div class="container">
   <div class="tabor">
     <el-form size="mini" :inline="true" :model="params" class="demo-form-inline">
       <el-form-item label="用户名">
         <el-input v-model.trim="params.username" clearable placeholder="用户名" @clear="search" />
       </el-form-item>
       <el-form-item label="昵称">
         <el-input v-model.trim="params.nickname" clearable placeholder="昵称" @clear="search" />
       </el-form-item>
       <el-form-item label="状态">
         <el-select v-model.trim="params.status" clearable placeholder="状态" @change="search" @clear="search">
           <el-option label="正常" value="1" />
           <el-option label="禁用" value="2" />
         </el-select>
       </el-form-item>
       <el-form-item label="手机号">
         <el-input v-model.trim="params.mobile" clearable placeholder="手机号" @clear="search" />
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
        <el-table-column show-overflow-tooltip sortable prop="username" label="用户名" />
        <el-table-column show-overflow-tooltip sortable prop="nickname" label="昵称" />
        <el-table-column show-overflow-tooltip sortable prop="status" label="状态" align="center">
          <template slot-scope="scope">
            <el-tag size="small" :type="scope.row.status === 1 ? 'success':'danger'" disable-transitions>{{ scope.row.status === 1 ? '正常':'禁用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip sortable prop="mobile" label="手机号" />
        <el-table-column show-overflow-tooltip sortable prop="creator" label="创建人" />
        <el-table-column show-overflow-tooltip sortable prop="introduction" label="说明" />
        <el-table-column label="操作" align="center" width="300px">
          <template slot-scope="scope">
            <!--
            @author:风很大
            @description: 操作  编辑 删除
            @time: 2021/12/22 0022
            -->
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

      <el-dialog :title="dialogFormTitle" :visible.sync="dialogFormVisible" width="30%">
        <el-form ref="dialogForm" size="small" :model="dialogFormData" :rules="dialogFormRules" label-width="100px">
          <el-form-item label="用户名" prop="username">
            <el-input ref="password" v-model.trim="dialogFormData.username" placeholder="用户名" />
          </el-form-item>
          <el-form-item :label="dialogType === 'create' ? '新密码':'重置密码'" prop="password">
            <el-input v-model.trim="dialogFormData.password" autocomplete="off" :type="passwordType" :placeholder="dialogType === 'create' ? '新密码':'重置密码'" />
            <span class="show-pwd" @click="showPwd">
              <svg-icon :icon-class="passwordType === 'password' ? 'eye' : 'eye-open'" />
            </span>
          </el-form-item>
          <el-form-item label="角色" prop="roleIds">
            <el-select v-model.trim="dialogFormData.roleIds" multiple placeholder="请选择角色" style="width:100%">
              <el-option
                v-for="item in roles"
                :key="item.ID"
                :label="item.name"
                :value="item.ID"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="状态" prop="status">
            <el-select v-model.trim="dialogFormData.status" placeholder="请选择状态" style="width:100%">
              <el-option label="正常" :value="1" />
              <el-option label="禁用" :value="2" />
            </el-select>
          </el-form-item>
          <el-form-item label="昵称" prop="nickname">
            <el-input v-model.trim="dialogFormData.nickname" placeholder="昵称" />
          </el-form-item>
          <el-form-item label="手机号" prop="mobile">
            <el-input v-model.trim="dialogFormData.mobile" placeholder="手机号" />
          </el-form-item>
          <el-form-item label="说明" prop="introduction">
            <el-input v-model.trim="dialogFormData.introduction" type="textarea" placeholder="说明" show-word-limit maxlength="100" />
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
import { getUsers, createUser, updateUserById, batchDeleteUserByIds } from '@/api/permission/user'
import { getRoles } from '@/api/permission/role'
import Table from '@/components/Table'
import Dialog from '@/components/Dialog'
import Button from '@/components/Button'
import Pagination from '@/components/Pagination'

export default {
  name: 'User',
  inject: ['reload'],
  components: { Button, Table, Dialog, Pagination },
  data() {
    var checkPhone = (rule, value, callback) => {
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
        username: '',
        nickname: '',
        status: '',
        mobile: '',
        pageNum: 1,
        pageSize: 10
      },
      // 表格数据
      loading: false,
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
      // 角色
      roles: [],

      passwordType: 'password',
      // dialog对话框
      submitLoading: false,
      dialogFormTitle: '',
      dialogType: '',
      dialogFormVisible: false,
      dialogFormData: {
        username: '',
        password: '',
        nickname: '',
        status: 1,
        mobile: '',
        avatar: '',
        introduction: '',
        roleIds: ''
      },
      dialogFormRules: {
        username: [
          { required: true, message: '请输入用户名', trigger: 'blur' },
          { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
        ],
        password: [
          { required: false, message: '请输入密码', trigger: 'blur' },
          { min: 6, max: 30, message: '长度在 6 到 30 个字符', trigger: 'blur' }
        ],
        nickname: [
          { required: false, message: '请输入昵称', trigger: 'blur' },
          { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
        ],
        mobile: [
          { required: true, validator: checkPhone, trigger: 'blur' }
        ],
        status: [
          { required: true, message: '请选择状态', trigger: 'change' }
        ],
        introduction: [
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
    this.getRoles()
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
        const { data } = await getUsers(this.params)
        this.table = data.users
        this.pagination.total = data.total
      } finally {
        this.loading = false
      }
    },

    // 获取角色数据
    async getRoles() {
      const res = await getRoles()
      this.roles = res.data.roles
    },

    // 新增
    create() {
      this.dialogFormTitle = '新增用户'
      this.dialogType = 'create'
      this.dialogFormVisible = true
    },

    // 修改
    update(row) {
      this.dialogFormData.ID = row.ID
      this.dialogFormData.username = row.username
      this.dialogFormData.password = ''
      this.dialogFormData.nickname = row.nickname
      this.dialogFormData.status = row.status
      this.dialogFormData.mobile = row.mobile
      this.dialogFormData.introduction = row.introduction
      this.dialogFormData.roleIds = row.roleIds

      this.dialogFormTitle = '修改用户'
      this.dialogType = 'update'
      this.passwordType = 'password'
      this.dialogFormVisible = true
    },

    // 提交表单
    submitForm() {
      this.$refs['dialogForm'].validate(async valid => {
        if (valid) {
          this.submitLoading = true

          let message = ''
          try {
            if (this.dialogType === 'create') {
              const { msg } = await createUser(this.dialogFormData)
              message = msg
            } else {
              const { msg } = await updateUserById(this.dialogFormData.ID, this.dialogFormData)
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
        username: '',
        password: '',
        nickname: '',
        status: 1,
        mobile: '',
        avatar: '',
        introduction: '',
        roleIds: ''
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
        const userIds = []
        this.multipleSelection.forEach(x => {
          userIds.push(x.ID)
        })
        let message = ''
        try {
          const { msg } = await batchDeleteUserByIds({ userIds: userIds })
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
      this.$confirm('此操作将永久删除, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async res => {
        this.loading = true
        const userIds = []
        this.multipleSelection.forEach(x => {
          userIds.push(x.ID)
        })
        let message = ''
        try {
          const { msg } = await batchDeleteUserByIds({ userIds: [Id] })
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

    showPwd() {
      if (this.passwordType === 'password') {
        this.passwordType = ''
      } else {
        this.passwordType = 'password'
      }
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
  .show-pwd {
    position: absolute;
    right: 10px;
    top: 3px;
    font-size: 16px;
    color: #889aa4;
    cursor: pointer;
    user-select: none;
  }
</style>
