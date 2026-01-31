import axios, { type AxiosInstance, type AxiosError } from 'axios'

class APIClient {
  private client: AxiosInstance

  constructor() {
    this.client = axios.create({
      baseURL: 'http://localhost:8080/api',
      timeout: 30000,
      headers: {
        'Content-Type': 'application/json'
      }
    })

    // 响应拦截器 - 统一错误处理
    this.client.interceptors.response.use(
      (response) => response,
      (error: AxiosError) => {
        if (error.response) {
          // 服务器返回错误
          console.error('API Error:', error.response.data)
        } else if (error.request) {
          // 请求发送但无响应
          console.error('Network Error:', error.message)
        } else {
          // 其他错误
          console.error('Error:', error.message)
        }
        return Promise.reject(error)
      }
    )
  }

  // GET请求
  async get<T>(url: string, params?: any): Promise<T> {
    const response = await this.client.get<T>(url, { params })
    return response.data
  }

  // POST请求
  async post<T>(url: string, data?: any): Promise<T> {
    const response = await this.client.post<T>(url, data)
    return response.data
  }
}

export const apiClient = new APIClient()
