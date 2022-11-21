<template>
  <div class="app-container">
    <span class="text">防抖指令 🍇🍇🍇🍓🍓🍓</span>
    <el-button type="primary" v-debounce="debounce">防抖按钮 (0.5秒后执行)</el-button>
  </div>
</template>

<script>
import { Message } from 'element-ui'
import debounce from '@/directive/debounce'


export default {
  name: 'index',
  directives: { debounce},
  data() {
    return {
      copyText: '66666666'
    }
  },
  methods: {
    debounce: function() {
      console.log('时间', new Date().toLocaleTimeString())
      console.log(this.copyText)
      Message.success('我是防抖按钮触发的事件 🍍🍓🍌')
    }
  }
}
</script>

<style scoped>

</style>
