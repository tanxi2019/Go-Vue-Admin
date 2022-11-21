const debounce = {
  bind: function(el, { value }) {
    el.newValue = value
    el.__handleClick__ = function() {
      let timer = null
      // 闭包 函数里return出函数
      return function() {
        // 判断timer的值
        if (timer != null) {
          // 清除定时器
          clearTimeout(timer)
        }
        // 定时器
        timer = setTimeout(() => {
          // 回调函数，apply() 改变 this 指向。
          el.newValue()
          timer = null
        }, 500)
      }
    }
    el.addEventListener('click', el.__handleClick__())
    //插入节点
  },
  inserted() {},
  componentUpdated: function(el) {
    el.newvalue = value
  },
  unbind: function(el) {
    el.removeEventListener('click', el.__handleClick__())
  }
}

export default debounce
