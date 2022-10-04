<template>
  <!--
  @author:风很大
  @description: 按钮组件封装, Add:{ name: '新增', size: 'mini',type: 'primary', icon: 'el-icon-plus',plain:true,disabled:true,show:true }
  @but:点击自定义事件
  @time: 2022/1/18 0018
  -->
  <el-button
    v-show="but.show"
    :circle="but.circle"
    :icon="but.icon"
    :type="but.type"
    :plain="but.plain"
    :size="but.size"
    :disabled="but.disabled"
    @click="handleClick">
    <template v-if="but.circle"></template>
    <template v-else-if="but.icon">{{but.name}}</template>
    <template v-else-if="but.name">{{but.name}}</template>
    <template v-else></template>
  </el-button>
</template>

<script>
  export default {
    props: {
      but: {
        type: Object
      },
    },
    name: 'index',
    methods: {
      handleClick: function() {
        this.$emit('but')
      }
    }
  }
</script>
