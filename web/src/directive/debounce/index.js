import debounce from './debounce'

const install = function(Vue) {
  Vue.directive('debounce', debounce)
}

if (window.Vue) {
  window.debounce = debounce
  Vue.use(install)
}

debounce.install = install
export default debounce
