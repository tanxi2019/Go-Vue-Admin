<template>
  <div class="app-container">
    <span class="text">้ฒๆๆไปค ๐๐๐๐๐๐</span>
    <el-button type="primary" v-throttle="throttle">่ๆตๆ้ฎ (1็งๅๆง่ก)</el-button>
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
      console.log('ๆถ้ด', new Date().toLocaleTimeString())
      console.log(this.copyText)
      Message.success('ๆๆฏ่ๆตๆ้ฎ่งฆๅ็ไบไปถ ๐๐๐')
    },
  }
}
</script>

<style scoped>

</style>
