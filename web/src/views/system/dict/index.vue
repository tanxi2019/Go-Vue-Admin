<template>
  <div class="container">
    <el-card class="container-card" shadow="always">
<!-- tabor     -->
      <el-form size="mini" :inline="true" :model="selectFrom"  class="demo-form-inline">
        <el-form-item label="ID" prop="age">
          <el-input v-model.trim="selectFrom.id" clearable placeholder="ID" @clear="hadleClickSelect" />
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input v-model.trim="selectFrom.name" clearable placeholder="名称" @clear="hadleClickSelect" />
        </el-form-item>
        <el-form-item label="关键词" prop="mobile">
          <el-input v-model.trim="selectFrom.keyword" clearable placeholder="关键词" @clear="hadleClickSelect" />
        </el-form-item>
        <el-form-item>
          <Button :but="Select" @but="hadleClickSelect"/>
          <Button :but="Add" @but="handClickAddDialog"/>
          <Button :but="DeleteAll" @but="hadleClickRemoveAll"/>
          <Button :but="Refresh" @but="refresh()"/>
        </el-form-item>

      </el-form>
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
        <el-table-column show-overflow-tooltip sortable prop="ID" label="ID" />
        <el-table-column show-overflow-tooltip sortable prop="name" label="名称" />
        <el-table-column show-overflow-tooltip sortable prop="keyword" label="关键词" />
        <el-table-column show-overflow-tooltip sortable label="状态" prop="hidden" align="center">
          <!--
          @author:风很大
          @description: 0:显示，1隐藏
          @time: 2022/1/18 0018
          -->
          <template slot-scope="scope">
            <el-switch v-model="scope.row.status!=0" :disabled="true"></el-switch>
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip sortable prop="desc" label="描述" />
        <el-table-column show-overflow-tooltip sortable prop="creator" label="创建人" />
        <el-table-column label="操作" align="center" width="300px">
          <template slot-scope="scope">
            <!--
            @author:风很大
            @description: 操作  编辑 删除
            @time: 2021/12/22 0022
            -->
            <Button :but="Edit" @but="handClickEditDialog(scope.row)"/>
            <Button :but="Detail" @but="handClickDetailDialog(scope.row)"/>
            <Button :but="Delete" @but="handleClickRemove(scope.row)"/>
          </template>
        </el-table-column>
      </Table>
      <!--
      @author:风很大
      @description: 分页组件
      @time: 2022/1/17 0017
      -->
      <Pagination :pagination="pagination" @page="handleCurrentChange" @size="handleSizeChange"/>
      <!--
      @author:风很大
      @description: 新增
      @time: 2022/1/13 0013
     -->
      <Dialog :dialog="AddDialog" @confirm="handClickAdd" >
        <template slot="dialog">
          <el-form ref="addForm" size="small" :model="addForm" :rules="addFormRules" label-width="80px">

            <el-form-item label="名称" prop="name">
              <el-input  v-model.trim="addForm.name" placeholder="名称" />
            </el-form-item>

            <el-form-item label="关键词" prop="keyword">
              <el-input v-model.trim="addForm.keyword" placeholder="关键词" />
            </el-form-item>

            <el-form-item label="排序" prop="sort">
              <el-input-number v-model.number="addForm.sort" controls-position="right" :min="1" :max="999" />
            </el-form-item>

            <el-form-item label="状态" prop="status">
              <el-switch v-model="addForm.status!=0" ></el-switch>

            </el-form-item>

            <el-form-item label="创建人" prop="创建人">
              <el-input v-model.trim="addForm.creator" placeholder="创建人" />
            </el-form-item>

            <el-form-item label="描述" prop="desc">
              <el-input v-model.trim="addForm.desc" type="textarea" placeholder="描述" show-word-limit maxlength="100" />
            </el-form-item>

          </el-form>
        </template>
      </Dialog>
      <!--
      @author:风很大
      @description: 编辑
      @time: 2022/1/13 0013
      -->
      <Dialog :dialog="EditDialog" @confirm="handClickEdit">
        <template slot="dialog">
          <el-form ref="editForm" size="small" :model="editForm" :rules="editFormRules" label-width="80px">

            <el-form-item label="名称" prop="name">
              <el-input  v-model.trim="editForm.name" placeholder="名称" />
            </el-form-item>

            <el-form-item label="关键词" prop="keyword">
              <el-input v-model.trim="editForm.keyword" placeholder="关键词" />
            </el-form-item>

            <el-form-item label="排序" prop="sort">
              <el-input-number v-model.number="editForm.sort" controls-position="right" :min="1" :max="999" />
            </el-form-item>

            <el-form-item label="状态" prop="status">
              <el-switch v-model="editForm.status"></el-switch>
            </el-form-item>

            <el-form-item label="创建人" prop="创建人">
              <el-input v-model.trim="editForm.creator" placeholder="创建人" />
            </el-form-item>

            <el-form-item label="描述" prop="desc">
              <el-input v-model.trim="editForm.desc" type="textarea" placeholder="描述" show-word-limit maxlength="100" />
            </el-form-item>

          </el-form>
        </template>
      </Dialog>
    </el-card>
  </div>
</template>
<script>
import { listDict,removeDict,deleteDict,createDict,putDict } from '@/api/system/dict'
import Table from '@/components/Table'
import Dialog from '@/components/Dialog'
import Button from '@/components/Button'
import Pagination from '@/components/Pagination'

export default {
  name:"index",
  inject:['reload'],
  components: { Button, Table, Dialog, Pagination },
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
      multipleSelection:[],
      // 查询表单
      selectFrom: {
        id: null,
        name: '',
        keyword: '',
      },
      // 添加表单验证
      addFormRules: {
        name: [
          { required: true, message: '请输入姓名', trigger: 'blur' },
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
      // 编辑表单验证
      editFormRules: {
        name: [
          { required: true, message: '请输入姓名', trigger: 'blur' },
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
      // 详情表单验证
      detailFormRules: {
        name: [
          { required: true, message: '请输入姓名', trigger: 'blur' },
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
      // 添加表单
      addForm: {
        name: '',
        keyword: '',
        desc: '',
        status:1,
        sort: 999,
        creator: ''
      },
      // 编辑表单
      editForm: {
        id:'',
        name: '',
        keyword: '',
        desc: '',
        sort: '',
        status:'',
        creator: ''
      },
      // 详情表单
      detailForm: {
        id:'',
        name: '',
        age: 0,
        sex: null,
        mobile: '',
        description: ''
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
        circle:false,
        plain:false,
        disabled:false,
        show:true
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
        show: true,

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
        loading:false
      }
    }
  },
  created() {
    this.getTable()
  },
  methods: {
    // 刷新页面
    refresh() {
      this.reload();
    },
    /**
     * @Description: 获取表格数据
     * @author 风很大
     * @page 页码
     * @size 页面大小
     * @date 2022/1/17 0017
     */
    getTable: async function() {
      this.setting.loading = true
      let table = await listDict({ page: this.pagination.page, size: this.pagination.size })
      let {code,data} = table
      if (code === 200){
        this.table = data.data
        this.pagination.total = data.total
        setTimeout(() => {
          this.setting.loading = false
        }, 1000);
      }

    },
    /**
     * @Description: 分页点击事件
     * @author 风很大
     * @date 2022/1/17 0017
     */
    handleCurrentChange: function(currentPage) {
      this.pagination.page = currentPage
      this.getTable()
    },
    /**
     * @Description: 控制每页条数
     * @author 风很大
     * @date 2022/1/17 0017
     */
    handleSizeChange:function(val) {
      this.pagination.pageSize = val
      this.getTable()
    },
    /**
     * @Description: 删除点击事件
     * @author 风很大
     * @date 2022/1/17 0017
     */
    handleClickRemove: function(item) {
      this.$confirm('确认删除？')
        .then(() => {
          let data = { id: item.ID }
          deleteDict(data).then( res =>{
            this.$message({
              message: res.msg,
              type: 'success'
            });
            this.getTable()
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
    /**
     * @Description: 查询事件
     * @author 风很大
     * @date 2022/1/17 0017
     */
    hadleClickSelect: function() {
      this.setting.loading = true
      let data = {
        page:this.pagination.page,
        id: this.selectFrom.id,
        name: this.selectFrom.name,
        keyword: this.selectFrom.keyword,

      }
      listDict(data).then(res => {
        this.table = res.data.data // 赋值
        this.pagination.total = res.data.total
        setTimeout(() => {
          this.setting.loading = false
        }, 1000);
      })
    },
    /**
     * @Description: 编辑事件
     * @author 风很大
     * @date 2022/1/13 0013
     */
    handClickEditDialog: function(item) {

      this.EditDialog.dialog = true
      this.editForm.id = item.ID
      this.editForm.name = item.name
      this.editForm.keyword = item.keyword
      this.editForm.sort = item.sort
      this.editForm.status = item.status
      this.editForm.creator = item.creator
      this.editForm.desc = item.desc
      this.$nextTick(()=>{
        this.$refs.editForm.clearValidate()
      })
    },

    /**
     * @Description: 编辑确定事件
     * @author 风很大
     * @date 2022/1/13 0013
     */
    handClickEdit: function() {
      this.$refs.editForm.validate((valid) => {
        if (valid) {
          this.EditDialog.dialog = false
          putDict(this.editForm).then( res =>{
            this.getTable()
            this.$message({
              message: res.msg,
              type: 'success'
            });

          })
        } else {
          console.log('error submit!!')
          return false
        }
      })

    },
    /**
     * @Description: 详情
     * @author 风很大
     * @date 2022/1/20 0020
     */
    handClickDetailDialog: function(item) {
      this.$router.push({
        path:`/system/dict/${item.keyword}`,
        params: {keyword:item.keyword}
      }) // 动态路由
    },
    /**
     * @Description: 详情确定
     * @author 风很大
     * @date 2022/1/20 0020
     */
    handClickDetail: function() {
      this.DetailDialog.dialog = false
      console.log('详情确定', this.detailForm)
    },
    /**
     * @Description: 新增弹框
     * @author 风很大
     * @date 2022/1/13 0013
     */
    handClickAddDialog: function() {
      this.AddDialog.dialog = true
      this.$nextTick(()=>{
        this.$refs.addForm.clearValidate()
      })
    },

    /**
     * @Description: 新增确定事件
     * @author 风很大
     * @date 2022/1/13 0013
     */
    handClickAdd: function() {
      this.$refs.addForm.validate((valid) => {
        if (valid) {
          this.AddDialog.dialog = false
          let data = {
            name: this.addForm.name,
            keyword: this.addForm.keyword,
            sort: this.addForm.sort,
            creator: this.addForm.creator,
            desc: this.addForm.desc
          }
          createDict(data).then( res =>{
            this.$message({
              message: res.msg,
              type: 'success'
            });
            this.getTable()
          })

        } else {
          console.log('error submit!!')
          return false
        }
      })
    },
    /**
     * @Description: 全选事件
     * @author 风很大
     * @date 2022/1/17 0017
     */
    handleSelectionChange: function(val) {
      this.multipleSelection = val
      this.DeleteAll.disabled =  this.multipleSelection.length === 0
    },
    /**
     * @Description: 全部删除事件
     * @author 风很大
     * @date 2022/1/17 0017
     */
    hadleClickRemoveAll: function() {
      this.$confirm('此操作将永久删除, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        const dictIds = []
        this.multipleSelection.forEach(x => {
          dictIds.push(x.ID)
        })

        const { msg } = await removeDict({ dictIds: dictIds })

        await this.getTable()
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
    /**
     * @Description: 导入事件
     * @author 风很大
     * @date 2022/1/17 0017
     */
    hadleClickImport: function() {
      console.log('导入')
    },
    /**
     * @Description: 导出事件
     * @author 风很大
     * @date 2022/1/17 0017
     */
    hadleClickExport: function() {
      console.log('导出')
    }
  }
}
</script>

<style lang="scss">
.container-card{
  margin: 10px;
}
.table-box {
  width: 100%;
  margin: 15px 0;
}

.but-box {
  padding-top: 8px;
}

.search {
  height: 33px;
  display: inline-block;
}

.login-form {
  border-radius: 6px;
  background: #ffffff;

  .search-input {
    height: 33px;
    width: 300px;
    margin-right: 10px;

    input {
      height: 33px;
    }
  }

  .input-icon {
    height: 39px;
    width: 14px;
    margin-left: 2px;
  }
}
</style>

