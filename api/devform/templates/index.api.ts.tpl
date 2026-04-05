import type { BaseResponse } from '@/utils/http/types'
import { http } from '@/utils/http'

/** 由在线表单生成 — [[.Description]] */
export interface [[.EntityName]]Row {
  id: number
  createTime: string
  updateTime: string
[[range .Fields -]]
  [[.JSONName]]: [[.TSType]]
[[end]]
}

export interface [[.EntityName]]ListSearchParams {
[[range .QueryFieldsTS -]]
  [[.JSONName]]?: [[.TSType]]
[[end]]
}

export class Api {
  static list(params: ApiUtil.WithPaginationParams<[[.EntityName]]ListSearchParams>): Promise<BaseResponse<ApiUtil.PaginationResponse<[[.EntityName]]Row>>> {
    return http.get('/api/[[.RouteGroup]]/list', { params })
  }

  static query(id: number | string): Promise<
    BaseResponse<[[.EntityName]]Row>
  > {
    return http.get(`/api/[[.RouteGroup]]/query/${id}`)
  }

  static add(data: Omit<[[.EntityName]]Row, 'id' | 'createTime' | 'updateTime'>) {
    return http.post('/api/[[.RouteGroup]]/add', data)
  }

  static edit(data: Pick<[[.EntityName]]Row, 'id'> & Partial<Omit<[[.EntityName]]Row, 'createTime' | 'updateTime'>>) {
    return http.put('/api/[[.RouteGroup]]/edit', data)
  }

  static del(id: number | string) {
    return http.delete(`/api/[[.RouteGroup]]/delete/${id}`)
  }
}
