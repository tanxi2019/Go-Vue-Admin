import request from '@/utils/request'

// 创建
export function createDictDetails(data) {
  return request({
    url: '/api/dict/details/create',
    method: 'post',
    data
  })
}

// 列表
export function listDictDetails(params) {
  return request({
    url: '/api/dict/details/list',
    method: 'get',
    params
  })
}

// 更新
export function putDictDetails(data) {
  return request({
    url: '/api/dict/details/put',
    method: 'put',
    data
  })
}

// 删除
export function deleteDictDetails(data) {
  return request({
    url: '/api/dict/details/delete',
    method: 'delete',
    data
  })
}

// 批量删除
export function removeDictDetails(data) {
  return request({
    url: '/api/dict/details/remove',
    method: 'delete',
    data
  })
}

