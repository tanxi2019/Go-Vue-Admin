<template>
  <div>
    <!--
    @author:风很大
    @description: 表格组件
    @table:Array
    @setting:Object, setting:{checkbox:true,order:false },
    @select:全选自定义事件
    @time: 2022/1/20 0020
    -->
    <el-table
      v-loading="setting.loading"
      :data="table"
      row-key="id"
      @selection-change="handleSelectionChange">

      <el-table-column
        v-if="setting.checkbox"
        type="selection"
        align="center"
        width="50">
      </el-table-column>

      <el-table-column label="序号" v-if="setting.order" align="center" width="200">
        <template slot-scope="scope">
          <span>
            {{scope.$index +1 }}
          </span>
        </template>
      </el-table-column>
      <!--
      @author:风很大
      @description: 表格数据默认插槽
      @time: 2022/1/13 0013
      -->
      <slot></slot>
    </el-table>

  </div>
</template>

<script>
  export default {
    props: {
      table: {
        type: Array
      },
      setting: {
        type: Object
      }
    },
    methods: {
      // 全选
      handleSelectionChange(val) {
        this.$emit('select', val)
      }
    }
  }
</script>

