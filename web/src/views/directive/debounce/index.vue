<template>
  <div class="app-container">
    <span class="text">้ฒๆๆไปค ๐๐๐๐๐๐</span>
    <el-button type="primary" v-debounce="debounce">้ฒๆๆ้ฎ (0.5็งๅๆง่ก)</el-button>
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
      console.log('ๆถ้ด', new Date().toLocaleTimeString())
      console.log(this.copyText)
      Message.success('ๆๆฏ้ฒๆๆ้ฎ่งฆๅ็ไบไปถ ๐๐๐')
    }
  }
}
</script>

<style scoped>

</style>
