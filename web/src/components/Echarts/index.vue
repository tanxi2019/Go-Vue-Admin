<template>
  <!--
  @author:风很大
  @description: Echarts
  @id: String
  @className: String
  @options: Object，参考echarts配置
  @time: 2022/1/20 0020
  -->
  <div :id="id" :class="className" ></div>
</template>
<script>
  import * as echarts from 'echarts'

  export default {
    name: 'index',
    props: {
      id: {
        type: String
      },
      options: {
        type: Object
      },
      className: {
        type: String
      }
    },
    mounted: function() {
      // 初始化
      let myChart = echarts.init(document.getElementById(this.id),'dark')
      myChart.setOption(this.options)
    }
  }
</script>

<style scoped>

</style>
