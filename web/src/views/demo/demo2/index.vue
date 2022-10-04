<template>
  <div class="container">
    <el-card class="container-card" shadow="always">
<!-- tabor     -->
      <el-form size="mini" :inline="true" :model="selectFrom"  class="demo-form-inline">
        <el-form-item label="姓名" prop="name">
          <el-input v-model.trim="selectFrom.name" clearable placeholder="姓名" @clear="hadleClickSelect" />
        </el-form-item>
        <el-form-item label="年龄" prop="age">
          <el-input v-model.trim="selectFrom.age" clearable placeholder="年龄" @clear="hadleClickSelect" />
        </el-form-item>
        <el-form-item label="性别" prop="sex">
          <el-select v-model.trim="selectFrom.sex" clearable placeholder="性别" @change="hadleClickSelect" @clear="hadleClickSelect">
            <el-option label="男" value="1" />
            <el-option label="女" value="2" />
          </el-select>
        </el-form-item>
        <el-form-item label="手机号" prop="mobile">
          <el-input v-model.trim="selectFrom.mobile" clearable placeholder="手机号" @clear="hadleClickSelect" />
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
        <el-table-column show-overflow-tooltip sortable prop="name" label="姓名" />
        <el-table-column show-overflow-tooltip sortable prop="sex" label="性别" align="center">
          <template slot-scope="scope">
            <el-tag size="small" >{{ scope.row.sex === 2 ? "女" : "男"}}</el-tag>
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip sortable prop="age" label="年龄" />
        <el-table-column show-overflow-tooltip sortable prop="count" label="访问量" />
        <el-table-column show-overflow-tooltip sortable prop="mobile" label="手机号" />
        <el-table-column show-overflow-tooltip sortable prop="description" label="描述" />
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

            <el-form-item label="姓名" prop="name">
              <el-input ref="password" v-model.trim="addForm.name" placeholder="姓名" />
            </el-form-item>

            <el-form-item label="性别" prop="sex">
              <el-select v-model.trim="addForm.sex" placeholder="请选择性别" style="width:100%">
                <el-option label="男" :value="1" />
                <el-option label="女" :value="2" />
              </el-select>
            </el-form-item>

            <el-form-item label="年龄" prop="age">
              <el-input v-model.trim="addForm.age" placeholder="年龄" />
            </el-form-item>

            <el-form-item label="手机号" prop="mobile">
              <el-input v-model.trim="addForm.mobile" placeholder="手机号" />
            </el-form-item>

            <el-form-item label="说明" prop="description">
              <el-input v-model.trim="addForm.description" type="textarea" placeholder="说明" show-word-limit maxlength="100" />
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

            <el-form-item label="姓名" prop="name">
              <el-input ref="password" v-model.trim="editForm.name" placeholder="姓名" />
            </el-form-item>

            <el-form-item label="性别" prop="sex">
              <el-select v-model.trim="editForm.sex" placeholder="请选择性别" style="width:100%">
                <el-option label="男" :value="1" />
                <el-option label="女" :value="2" />
              </el-select>
            </el-form-item>

            <el-form-item label="年龄" prop="age">
              <el-input v-model.trim="editForm.age" placeholder="年龄" />
            </el-form-item>

            <el-form-item label="手机号" prop="mobile">
              <el-input v-model.trim="editForm.mobile" placeholder="手机号" />
            </el-form-item>

            <el-form-item label="说明" prop="description">
              <el-input v-model.trim="editForm.description" type="textarea" placeholder="说明" show-word-limit maxlength="100" />
            </el-form-item>

          </el-form>
        </template>
      </Dialog>
      <!--
    @author:风很大
    @description: 详情
    @time: 2022/1/13 0013
    -->
      <Dialog :dialog="DetailDialog" @confirm="handClickDetail">
        <template slot="dialog">
          <el-form ref="detailForm" size="small" label-width="80px">

            <el-form-item label="姓名" prop="name" >
              <el-input v-model.trim="detailForm.name" placeholder="姓名" :disabled="true" />
            </el-form-item>

            <el-form-item label="性别" prop="sex">
              <el-select v-model.trim="detailForm.sex" placeholder="请选择性别" :disabled="true" style="width:100%">
                <el-option label="男" :value="1" />
                <el-option label="女" :value="2" />
              </el-select>
            </el-form-item>

            <el-form-item label="年龄" prop="age">
              <el-input v-model.trim="detailForm.age" placeholder="年龄" :disabled="true" />
            </el-form-item>

            <el-form-item label="手机号" prop="mobile">
              <el-input v-model.trim="detailForm.mobile" placeholder="手机号" :disabled="true" />
            </el-form-item>

            <el-form-item label="访问量" prop="count">
              <el-input v-model.trim="detailForm.count" placeholder="访问量" :disabled="true" />
            </el-form-item>

            <el-form-item label="说明" prop="description">
              <el-input v-model.trim="detailForm.description" type="textarea" placeholder="说明" :disabled="true" show-word-limit maxlength="100" />
            </el-form-item>

          </el-form>
        </template>
      </Dialog>
    </el-card>
  </div>
</template>
<script>
import { listExample,removeExample,deleteExample,createExample,putExample,GetExample } from '@/api/example'
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
        name: '',
        age: null,
        sex: null,
        mobile: '',
        description: ''
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
        age: null,
        sex: 1,
        mobile: '',
        description: ''
      },
      // 编辑表单
      editForm: {
        id:'',
        name: '',
        age: 0,
        sex: null,
        mobile: '',
        description: ''
      },
      // 详情表单
      detailForm: {
        id:'',
        name: '',
        age: 0,
        sex: null,
        count:0,
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
      let table = await listExample({ page: this.pagination.page, size: this.pagination.size })
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
          deleteExample(data).then( res =>{
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
        name: this.selectFrom.name,
        age: this.selectFrom.age,
        sex: this.selectFrom.sex,
        mobile: this.selectFrom.mobile,
        description:this.selectFrom.description
      }
      listExample(data).then(res => {
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
      this.editForm.age = item.age
      this.editForm.sex = item.sex
      this.editForm.mobile = item.mobile
      this.editForm.description = item.description
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
          let data = {
            id:this.editForm.id,
            name:this.editForm.name,
            age: parseInt(this.editForm.age),
            sex: this.editForm.sex,
            mobile: this.editForm.mobile,
            description:this.editForm.description
          }
          putExample(data).then( res =>{
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
    handClickDetailDialog: async function(item) {
      this.DetailDialog.dialog = true

      let {code,data} = await GetExample({ id:item.ID })
      if (code === 200){
        this.detailForm.age = data.age
        this.detailForm.sex = data.sex
        this.detailForm.count = data.count
        this.detailForm.description = data.description
        this.detailForm.name = data.name
        this.detailForm.mobile = data.mobile
      }
    },
    /**
     * @Description: 详情确定
     * @author 风很大
     * @date 2022/1/20 0020
     */
    handClickDetail: function() {
      this.DetailDialog.dialog = false
      this.getTable()
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
            age: parseInt(this.addForm.age),
            sex: this.addForm.sex,
            mobile: this.addForm.mobile,
            description: this.addForm.description
          }
          createExample(data).then( res =>{
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
        const exampleIds = []
        this.multipleSelection.forEach(x => {
          exampleIds.push(x.ID)
        })

        const { msg } = await removeExample({ exampleIds: exampleIds })

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

