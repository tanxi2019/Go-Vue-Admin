<template>
  <div class="app-container">
    <el-input v-model="copyText" placeholder="Please input" style="width:400px;max-width:100%;" />
    <el-button  type="primary" icon="el-icon-document" v-copy="copyText">点击复制</el-button>
  </div>
</template>

<script>
import copy from '@/directive/copy'

export default {
  name: 'index',
  directives: { copy },
  data() {
    return {
      copyText: '点击复制按钮触发的事件 🍍🍓🍌'
    }
  }
}
</script>

<style scoped>

</style>
