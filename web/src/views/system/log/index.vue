<template>
  <div class="container">
      <div class="tabor">
        <el-form size="mini" :inline="true" :model="params" class="demo-form-inline">
          <el-form-item label="请求人">
            <el-input v-model.trim="params.username" clearable placeholder="请求人" @clear="search" />
          </el-form-item>
          <el-form-item label="IP地址">
            <el-input v-model.trim="params.ip" clearable placeholder="IP地址" @clear="search" />
          </el-form-item>
          <el-form-item label="请求路径">
            <el-input v-model.trim="params.path" clearable placeholder="请求路径" @clear="search" />
          </el-form-item>
          <el-form-item label="请求状态">
            <el-input v-model.trim="params.status" clearable placeholder="请求状态" @clear="search" />
          </el-form-item>
          <el-form-item>
            <Button :but="Select" @but="search"/>
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
        <el-table-column show-overflow-tooltip sortable prop="username" label="请求人" />
        <el-table-column show-overflow-tooltip sortable prop="ip" label="IP地址" />
        <el-table-column show-overflow-tooltip sortable prop="path" label="请求路径" />
        <el-table-column show-overflow-tooltip sortable prop="status" label="请求状态" align="center">
          <template slot-scope="scope">
            <el-tag size="small" :type="scope.row.status | statusTagFilter" disable-transitions>{{ scope.row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip sortable prop="startTime" label="发起时间">
          <template slot-scope="scope">
            {{ parseGoTime(scope.row.startTime) }}
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip sortable prop="timeCost" label="请求耗时(ms)" align="center">
          <template slot-scope="scope">
            <el-tag size="small" :type="scope.row.timeCost | timeCostTagFilter" disable-transitions>{{ scope.row.timeCost }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip sortable prop="desc" label="说明" />
        <el-table-column fixed="right" label="操作" align="center" width="150">
          <template slot-scope="scope">
            <el-button :loading="loading" size="mini" icon="el-icon-delete" type="text" @click="singleDelete(scope.row.ID)">删除</el-button>
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
  </div>
</template>

<script>
import { getOperationLogs, batchDeleteOperationLogByIds } from '@/api/system/operationLog'
import { parseGoTime } from '@/utils'
import Table from '@/components/Table'
import Button from '@/components/Button'
import Pagination from '@/components/Pagination'

export default {
  name: 'Index',
  inject: ['reload'],
  components: { Button, Table, Pagination },
  filters: {
    statusTagFilter(val) {
      if (val === 200) {
        return 'success'
      } else if (val === 400) {
        return 'warning'
      } else if (val === 401) {
        return 'danger'
      } else if (val === 403) {
        return 'danger'
      } else if (val === 500) {
        return 'danger'
      } else {
        return 'info'
      }
    },
    timeCostTagFilter(val) {
      if (val <= 200) {
        return 'success'
      } else if (val > 200 && val <= 1000) {
        return ''
      } else if (val > 1000 && val <= 2000) {
        return 'warning'
      } else {
        return 'danger'
      }
    }
  },
  data() {
    return {
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
      // 查询参数
      params: {
        username: '',
        ip: '',
        path: '',
        status: '',
        pageNum: 1,
        pageSize: 10
      },
      // 表格和分页配置
      pagination: {
        page: 1,
        size: 10,
        total: 0
      },
      // 表格数据
      table: [],
      setting: {
        checkbox: true,
        order: false,
        loading: false
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
    parseGoTime,
    // 刷新页面
    refresh() {
      this.reload()
    },

    // 查询
    search() {
      this.getTableData()
    },

    // 获取表格数据
    async getTableData() {
      this.loading = true
      try {
        const { data } = await getOperationLogs(this.params)
        this.table = data.data
        this.pagination.total = data.total
      } finally {
          this.loading = false
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
        const operationLogIds = []
        this.multipleSelection.forEach(x => {
          operationLogIds.push(x.ID)
        })
        let message = ''
        try {
          const { msg } = await batchDeleteOperationLogByIds({ operationLogIds: operationLogIds })
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
      }).then(async () => {
        this.loading = true
        const operationLogIds = []
        this.multipleSelection.forEach(x => {
          operationLogIds.push(x.ID)
        })
        let message = ''
        try {
          const { msg } = await batchDeleteOperationLogByIds({ operationLogIds: [Id] })
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
</style>
