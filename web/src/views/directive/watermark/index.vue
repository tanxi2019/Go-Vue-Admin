<template>
  <div class="app-container"
       style="height: 900px;"
       v-watermark="watermark"
  >
  </div>
</template>

<script>
import watermark from '@/directive/watermark'
import { aesEncrypt, aesDecrypt, sign } from '@/utils/cryptoJS'

export default {
  name: 'index',
  directives: { watermark },
  data() {
    return {
      watermark: {
        text: '🍇🍇🍇🍓🍓🍓',
        font: '16px Microsoft JhengHei',
        color: 'rgba(180, 180, 180, 0.3)',
        time: new Date().toLocaleTimeString()
      }
    }
  },
  created() {
    let a = aesEncrypt(JSON.stringify(this.watermark))
    console.log('密文',a)
    let b = JSON.parse(aesDecrypt(a))
    console.log('明文',b)

  }
}
</script>

<style scoped>

</style>
