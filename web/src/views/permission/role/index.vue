<template>
  <div class="container">
    <div class="tabor">
      <el-form size="mini" :inline="true" :model="params" class="demo-form-inline">
        <el-form-item label="角色名称">
          <el-input v-model.trim="params.name" clearable placeholder="角色名称" @clear="search" />
        </el-form-item>
        <el-form-item label="关键字">
          <el-input v-model.trim="params.keyword" clearable placeholder="关键字" @clear="search" />
        </el-form-item>
        <el-form-item label="角色状态">
          <el-select v-model.trim="params.status" clearable placeholder="角色状态" @change="search" @clear="search">
            <el-option label="正常" :value="1" />
            <el-option label="禁用" :value="2" />
          </el-select>
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
        <el-table-column show-overflow-tooltip sortable prop="name" label="角色名称" />
        <el-table-column show-overflow-tooltip sortable prop="keyword" label="关键字" />
        <el-table-column show-overflow-tooltip sortable prop="sort" label="等级" />
        <el-table-column show-overflow-tooltip sortable prop="status" label="角色状态" >
          <template slot-scope="scope">
            <el-tag size="small" :type="scope.row.status === 1 ? 'success':'danger'" disable-transitions>{{ scope.row.status === 1 ? '正常':'禁用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip sortable prop="creator" label="创建人" />
        <el-table-column show-overflow-tooltip sortable prop="desc" label="说明" />
        <el-table-column label="操作" align="center" width="300px">
          <template slot-scope="scope">
            <!--
            @author:风很大
            @description: 操作  编辑 删除
            @time: 2021/12/22 0022
            -->
            <Button :but="Edit" @but="update(scope.row)"/>
            <el-button size="mini" icon="el-icon-key" type="text" @click="updateMenuPermission(scope.row.ID)">菜单权限</el-button>
            <el-button size="mini" icon="el-icon-key" type="text" @click="updateApiPermission(scope.row.ID)">数据权限</el-button>
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


      <el-dialog :title="dialogFormTitle" :visible.sync="dialogFormVisible" width="600px">
        <el-form ref="dialogForm" :inline="true" size="small" :model="dialogFormData" :rules="dialogFormRules" label-width="100px">
          <el-form-item label="角色名称" prop="name">
            <el-input v-model.trim="dialogFormData.name" placeholder="角色名称" style="width: 420px" />
          </el-form-item>
          <el-form-item label="关键字" prop="keyword">
            <el-input v-model.trim="dialogFormData.keyword" placeholder="关键字" style="width: 420px" />
          </el-form-item>
          <el-form-item label="角色状态" prop="status">
            <el-select v-model.trim="dialogFormData.status" placeholder="请选择角色状态" style="width: 180px">
              <el-option label="正常" :value="1" />
              <el-option label="禁用" :value="2" />
            </el-select>
          </el-form-item>
          <el-form-item label="等级" prop="sort">
            <el-input-number v-model.number="dialogFormData.sort" controls-position="right" :min="1" :max="999" />
          </el-form-item>
          <el-form-item label="说明" prop="desc">
            <el-input v-model.trim="dialogFormData.desc" style="width: 420px" type="textarea" placeholder="说明" show-word-limit maxlength="100" />
          </el-form-item>
        </el-form>
        <div slot="footer">
          <el-button size="mini" @click="cancelForm()">取 消</el-button>
          <el-button size="mini" :loading="submitLoading" type="primary" @click="submitForm()">确 定</el-button>
        </div>
      </el-dialog>
<!--菜单权限-->
      <el-dialog title="修改权限" :visible.sync="permsDialogVisible" width="600px" custom-class="perms-dialog">
        <el-tabs>
          <el-tab-pane>
            <span slot="label"> 菜单权限</span>
            <el-tree
              ref="roleMenuTree"
              :props="{children: 'children',label: 'title'}"
              :data="menuTree"
              show-checkbox
              node-key="ID"
              :default-checked-keys="defaultCheckedRoleMenu"
            />

          </el-tab-pane>
        </el-tabs>
        <div slot="footer">
          <el-button size="mini" :loading="permissionLoading" @click=" permsDialogVisible = false">取 消</el-button>
          <el-button size="mini" type="primary" @click="updateRoleMenusById">确 定</el-button>
        </div>
      </el-dialog>
<!--数据权限-->
    <el-dialog title="修改权限" :visible.sync="permsApiDialogVisible" width="580px" custom-class="perms-dialog">
      <el-tabs>
        <el-tab-pane>
                    <span slot="label">数据权限</span>
                    <el-tree
                      ref="roleApiTree"
                      :props="{children: 'children',label: 'desc'}"
                      :data="apiTree"
                      show-checkbox
                      node-key="ID"
                      :default-checked-keys="defaultCheckedRoleApi"
                    />

                  </el-tab-pane>
      </el-tabs>
      <div slot="footer">
        <el-button size="mini" :loading="permissionApiLoading" @click="permsApiDialogVisible = false">取 消</el-button>
        <el-button size="mini" type="primary" @click="updateRoleApisById">确 定</el-button>
      </div>
    </el-dialog>

  </div>
</template>

<script>
import { getRoles, createRole, updateRoleById, batchDeleteRoleByIds, getRoleMenusById, getRoleApisById, updateRoleMenusById, updateRoleApisById } from '@/api/permission/role'
import { getMenuTree } from '@/api/permission/menu'
import { getApiTree } from '@/api/system/api'
import Table from '@/components/Table'
import Dialog from '@/components/Dialog'
import Button from '@/components/Button'
import Pagination from '@/components/Pagination'

export default {
  name: 'Role',
  inject: ['reload'],
  components: { Button, Table, Dialog, Pagination },
  data() {
    return {
      // 查询参数
      params: {
        name: '',
        keyword: '',
        status: '',
        pageNum: 1,
        pageSize: 10
      },
      // 表格数据
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
      loading: false,

      // dialog对话框
      submitLoading: false,
      dialogFormTitle: '',
      dialogType: '',
      dialogFormVisible: false,
      dialogFormData: {
        name: '',
        keyword: '',
        status: 1,
        sort: 1,
        desc: ''
      },
      dialogFormRules: {
        name: [
          { required: true, message: '请输入角色名称', trigger: 'blur' },
          { min: 1, max: 20, message: '长度在 1 到 20 个字符', trigger: 'blur' }
        ],
        keyword: [
          { required: true, message: '请输入关键字', trigger: 'blur' },
          { min: 1, max: 20, message: '长度在 1 到 20 个字符', trigger: 'blur' }
        ],
        status: [
          { required: true, message: '请选择角色状态', trigger: 'change' }
        ],
        desc: [
          { required: false, message: '说明', trigger: 'blur' },
          { min: 0, max: 100, message: '长度在 0 到 100 个字符', trigger: 'blur' }
        ]
      },

      // 删除按钮弹出框
      popoverVisible: false,
      // 表格多选
      multipleSelection: [],

      // 修改菜单权限
      permsDialogVisible: false,
      permissionLoading: false,
      menuTree: [],
      defaultCheckedRoleMenu: [],
      // 修改数据权限
      permsApiDialogVisible: false,
      permissionApiLoading: false,
      apiTree: [],
      defaultCheckedRoleApi: [],

      // 被修改权限的角色ID
      roleId: 0
    }
  },
  created() {
    this.getTableData()
    this.getMenuTree()
    this.getApiTree()
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
        const { data } = await getRoles(this.params)
        this.table = data.roles
        this.pagination.total = data.total
      } finally {
        this.loading = false
      }
    },

    // 新增
    create() {
      this.dialogFormTitle = '新增角色'
      this.dialogType = 'create'
      this.dialogFormVisible = true
    },

    // 修改
    update(row) {
      this.dialogFormData.ID = row.ID
      this.dialogFormData.name = row.name
      this.dialogFormData.keyword = row.keyword
      this.dialogFormData.sort = row.sort
      this.dialogFormData.status = row.status
      this.dialogFormData.desc = row.desc

      this.dialogFormTitle = '修改角色'
      this.dialogType = 'update'
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
              const { msg } = await createRole(this.dialogFormData)
              message = msg
            } else {
              const { msg } = await updateRoleById(this.dialogFormData.ID, this.dialogFormData)
              message = msg
            }
          } finally {
            this.submitLoading = false
          }

          this.resetForm()
          this.getTableData()
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
        name: '',
        keyword: '',
        status: 1,
        sort: 1,
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
        const roleIds = []
        this.multipleSelection.forEach(x => {
          roleIds.push(x.ID)
        })
        let message = ''
        try {
          const { msg } = await batchDeleteRoleByIds({ roleIds: roleIds })
          message = msg
        } finally {
          this.loading = false
        }

        this.getTableData()
        this.$message({
          showClose: true,
          message: message,
          type: 'success'
        })
      }).catch(() => {
        this.$message({
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
    async singleDelete(id) {
      this.$confirm('此操作将永久删除, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async res => {
        this.loading = true
        const roleIds = []
        this.multipleSelection.forEach(x => {
          roleIds.push(x.ID)
        })
        let message = ''
        try {
          const { msg } = await batchDeleteRoleByIds({ roleIds: [id] })
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
          type: 'info',
          message: '已取消删除'
        })
      })
    },

    // 修改菜单权限
    async updateMenuPermission(roleId) {
      this.roleId = roleId
      this.permsDialogVisible = true
      await this.getMenuTree()
      await this.getRoleMenusById(roleId)
    },
    // 修改数据权限
    async updateApiPermission(roleId) {
      this.roleId = roleId
      this.permsApiDialogVisible = true
      await this.getApiTree()
      await this.getRoleApisById(roleId)
    },

    // 获取菜单树
    async getMenuTree() {
      this.menuTreeLoading = true
      try {
        const { data } = await getMenuTree()
        this.menuTree = data.menuTree
      } finally {
        this.menuTreeLoading = false
      }
    },

    // 获取接口树
    async getApiTree() {
      this.apiTreeLoading = true
      try {
        const { data } = await getApiTree()
        this.apiTree = data.apiTree
      } finally {
        this.apiTreeLoading = false
      }
    },

    // 获取角色的权限菜单
    async getRoleMenusById(roleId) {
      this.permissionLoading = true
      let rseData = []
      try {
        const { data } = await getRoleMenusById(roleId)
        rseData = data
      } finally {
        this.permissionLoading = false
      }

      const menus = rseData.menus
      const ids = []
      menus.forEach(x => { ids.push(x.ID) })
      this.defaultCheckedRoleMenu = ids
      this.$refs.roleMenuTree.setCheckedKeys(this.defaultCheckedRoleMenu)
    },

    // 获取角色的权限接口
    async getRoleApisById(roleId) {
      this.permissionLoading = true
      let resData = []
      try {
        const { data } = await getRoleApisById(roleId)
        resData = data
      } finally {
        this.permissionLoading = false
      }

      const apis = resData.apis
      const ids = []
      apis.forEach(x => { ids.push(x.ID) })
      this.defaultCheckedRoleApi = ids
      this.$refs.roleApiTree.setCheckedKeys(this.defaultCheckedRoleApi)
    },

    // 修改角色菜单
    async updateRoleMenusById() {
      this.permissionLoading = true
      let ids = this.$refs.roleMenuTree.getCheckedKeys()
      const idsHalf = this.$refs.roleMenuTree.getHalfCheckedKeys()
      ids = ids.concat(idsHalf)
      ids = [...new Set(ids)]

      try {
        await updateRoleMenusById(this.roleId, { menuIds: ids })
      } finally {
        this.permissionLoading = false
      }

      this.permsDialogVisible = false
      this.$message({
        showClose: true,
        message: '成功',
        type: 'success'
      })
    },

    // 修改角色接口
    async updateRoleApisById() {
      this.permissionApiLoading = true
      const ids = this.$refs.roleApiTree.getCheckedKeys(true)
      try {
        await updateRoleApisById(this.roleId, { apiIds: ids })
      } finally {
        this.permissionApiLoading = false
      }

      this.permsApiDialogVisible = false
      this.$message({
        showClose: true,
        message: '成功',
        type: 'success'
      })
    },

    // 取消修改角色权限
    cancelPermissionForm() {
      this.permsDialogVisible = false
      this.permsApiDialogVisible = false
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

<style lang="scss">
.table-box {
  background-color: #ffffff;
  padding: 15px 10px;
}
  .perms-dialog > .el-dialog__body{
    padding-top: 0;
    padding-bottom: 15px;
  }
</style>

