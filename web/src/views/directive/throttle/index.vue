<template>
  <div class="app-container">
    <span class="text">防抖指令 🍇🍇🍇🍓🍓🍓</span>
    <el-button type="primary" v-throttle="throttle">节流按钮 (1秒后执行)</el-button>
  </div>
</template>

<script>
import { Message } from 'element-ui'
import throttle from '@/directive/throttle'


export default {
  name: 'index',
  directives: { throttle },
  data() {
    return {
      copyText: '66666666'
    }
  },
  methods: {
    throttle: function() {
      console.log('时间', new Date().toLocaleTimeString())
      console.log(this.copyText)
      Message.success('我是节流按钮触发的事件 🍍🍓🍌')
    },
  }
}
</script>

<style scoped>

</style>
