import request from '@/utils/request'

// 创建
export function createDict(data) {
  return request({
    url: '/api/dict/create',
    method: 'post',
    data
  })
}

// 列表
export function listDict(params) {
  return request({
    url: '/api/dict/list',
    method: 'get',
    params
  })
}

// 更新
export function putDict(data) {
  return request({
    url: '/api/dict/put',
    method: 'put',
    data
  })
}

// 删除
export function deleteDict(data) {
  return request({
    url: '/api/dict/delete',
    method: 'delete',
    data
  })
}

// 批量删除
export function removeDict(data) {
  return request({
    url: '/api/dict/remove',
    method: 'delete',
    data
  })
}

