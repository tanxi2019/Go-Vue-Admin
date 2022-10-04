<template>
  <!--
  @author:风很大
  @description: 弹窗组件
  @dialog: Object， AddDialog:{ title:'新增',dialog:false, width:'600px' },
  @confirm:打开
  @time: 2022/1/20 0020
  -->
  <div class="dialog">
    <el-dialog
      :title="dialog.title"
      :width="dialog.width"
      :visible.sync="dialog.dialog">
      <slot name="dialog"></slot>
      <span slot="footer" class="dialog-footer">
    <el-button @click="dialog.dialog = false" size="small">取 消</el-button>
    <el-button type="primary" @click="handleClickConfirm" size="small">确 定</el-button>
  </span>
    </el-dialog>

  </div>
</template>

<script>
  export default {
    name: 'index',
    props: {
      dialog: {
        type: Object
      }
    },
    methods: {
      handleClickConfirm: function() {
        this.$emit('confirm')
      }
    }
  }
</script>

