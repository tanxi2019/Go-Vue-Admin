import copy from './copy'

const install = function(Vue) {
  Vue.directive('copy', copy)
}

if (window.Vue) {
  window.debounce = copy
  Vue.use(install) // eslint-disable-line
}

copy.install = install
export default copy
