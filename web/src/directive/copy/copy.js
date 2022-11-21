import Vue from 'vue'

const debounce = {
  bind: function(el, { value }) {
    el.newvalue = value
    el.handler = function() {
      if (!el.newvalue) {
        return
      }
      let textarea = document.createElement('textarea')
      textarea.style.position = 'absolute'
      textarea.style.left = '-1000px'
      textarea.style.readOnly = 'readonly'
      textarea.value = el.newvalue
      document.body.appendChild(textarea)
      textarea.select()

      let result = document.execCommand('copy')
      if (result) {
        Vue.prototype.$message({
          message: '复制成功',
          type: 'success',
          duration: 1500
        })
      }
      document.body.removeChild(textarea)
    }
    el.addEventListener('click', el.handler)
  },
  inserted: function() {
    //插入节点
  },
  componentUpdated: function(el, { value }) {
    el.newvalue = value
  },
  unbind: function(el, { value }) {
    el.removeEventListener('click', el.handler)
  }
}

export default debounce
