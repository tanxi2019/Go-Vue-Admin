<template>
  <div class="container">
    <h2>WebSocket</h2>
    <div class="msg">心跳包：{{ new Date().toLocaleTimeString() }} ：{{ heartBeat }}</div>
    <div class="msg">测试返回消息：{{ new Date().toLocaleTimeString() }} ：{{ msg }}</div>
    <el-button @click="handClickTest">测试连接：hello</el-button>
  </div>
</template>

<script>
export default {
  name: 'index',
  data() {
    return {
      msg: '',
      heartBeat:'',
      lockReconnect: false, // 是否真正建立连接
      timeout: 5 * 1000, // 5秒一次心跳
      timeoutObj: null, // 心跳心跳倒计时
      serverTimeoutObj: null, // 心跳倒计时
      timeoutNum: null,// 断开 重连倒计时
      reconnectTime: 5 * 1000 // 断开 重连倒计时
    }
  },
  created() {
    // 页面刚进入时开启长连接
    this.initWebSocket()
  },
  destroyed() {
    // 页面销毁时关闭长连接，必须是挂载方法。直接调用 this.close() 不生效
    this.$websocket.close()
    // 清除时间
    clearTimeout(this.timeoutObj)
    clearTimeout(this.serverTimeoutObj)
  },

  methods: {
    //初始化websocket
    initWebSocket() {
      this.$websocket = new WebSocket(process.env.VUE_APP_WEB_SOCKET+"/ws")
      this.$websocket.onopen = this.webSocketOnopen
      this.$websocket.onerror = this.webSocketOnerror
      this.$websocket.onmessage = this.webSocketOnmessage
      this.$websocket.onclose = this.webSocketClose
    },
    // 连接成功
    webSocketOnopen() {
      if (this.$websocket.readyState === 1) {
        console.log('连接成功', '状态：' + this.$websocket.readyState)
        // 开启心跳
        // this.start()
      }else {
        this.reconnect()
      }

    },
    // 链接错误
    webSocketOnerror(e) {
      console.log('连接发生错误', e.type, '状态：' + this.$websocket.readyState)
      // 重连
      this.reconnect()
    },
    // 接收消息
    webSocketOnmessage(e) {
      // let obj = e.data
      // console.log(obj,new Date().toLocaleTimeString())
      // this.msg = obj
      // this.reset()

      let obj = JSON.parse(e.data)
      console.log(obj)
      switch (obj.type) {
        case 'heartBeat':
          //收到服务器信息，心跳重置
          this.heartBeat = obj
          console.log('接收到的服务器消息：', obj.msg)
          // this.reset()
          break
        case 'hello':
          this.msg = obj
          console.log('接收到的服务器消息：', obj.msg)
      }
    },
    // 发送消息
    handClickTest:function() {
      if (this.$websocket.readyState === 1) {
        let data = {
          type:"hello"
        }
        this.$websocket.send(JSON.stringify(data))
      }
    },
    // 关闭
    webSocketClose(e) {
      if (e.code === 10000 && this.$websocket.readyState === 3) {
        console.log('手动关闭连接，不再重连', e.code, '状态：' + this.$websocket.readyState)
      } else {
        console.log('断线重连', e.code, '状态：' + this.$websocket.readyState)
        // 断线重连
        this.reconnect()
      }
    },
    //重新连接
    reconnect() {
      if (this.lockReconnect) {
        return
      }
      this.lockReconnect = true
      //没连接上会一直重连，设置延迟避免请求过多
      console.log('每5秒发起一次重连', new Date().toLocaleTimeString())
      this.timeoutNum && clearTimeout(this.timeoutNum)
      this.timeoutNum = setTimeout(() => {
        this.initWebSocket()//新连接
        this.lockReconnect = false
      }, this.reconnectTime)
    },
    // 重置心跳
    reset() {
      clearTimeout(this.timeoutObj)//清除时间
      clearTimeout(this.serverTimeoutObj)//清除时间
      console.log('每5秒发起一次心跳', new Date().toLocaleTimeString())
      this.start() //重启心跳
    },
    // 开启心跳
    start() {
      this.timeoutObj && clearTimeout(this.timeoutObj)
      this.serverTimeoutObj && clearTimeout(this.serverTimeoutObj)
      this.timeoutObj = setTimeout(() => {
        // 这里发送一个心跳，后端收到后，返回一个心跳消息，
        if (this.$websocket.readyState === 1) {
          // 如果连接正常
          let data = {
            type: "heartBeat"
          }
          this.$websocket.send(JSON.stringify(data))
        } else {
          // 断线重连
          this.reconnect()
        }
        this.serverTimeoutObj = setTimeout(() => {
          //超时关闭
          this.$websocket.close()
        }, this.timeout)

      }, this.timeout)
    }

  }
}
</script>

<style lang="scss" scoped>
.msg {
  padding: 20px;
  box-sizing: border-box;
}
</style>
