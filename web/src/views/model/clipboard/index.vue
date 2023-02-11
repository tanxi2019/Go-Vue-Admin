<template>
  <div class="app-container">
    <el-tabs v-model="activeName">
      <el-tab-pane label="函数形式" name="directly">
        <el-input v-model="inputData" placeholder="Please input" style="width:400px;max-width:100%;" />
        <el-button type="primary" icon="el-icon-document" @click="handleCopy(inputData,$event)">点击复制</el-button>
      </el-tab-pane>
      <el-tab-pane label="指令形式" name="v-directive">
        <el-input v-model="inputData" placeholder="Please input" style="width:400px;max-width:100%;" />
        <el-button v-clipboard:copy="inputData" v-clipboard:success="clipboardSuccess" type="primary" icon="el-icon-document">点击复制</el-button>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
  import clip from '@/utils/clipboard'
  import clipboard from '@/directive/clipboard'

  export default {
    name: 'ClipboardDemo',
    directives: {
      clipboard
    },
    data() {
      return {
        activeName: 'directly',
        inputData: '点击复制按钮触发的事件 🍍🍓🍌'
      }
    },
    methods: {
      handleCopy(text, event) {
        clip(text, event)
      },
      clipboardSuccess() {
        this.$message({
          message: 'Copy successfully',
          type: 'success',
          duration: 1500
        })
      }
    }
  }
</script>

