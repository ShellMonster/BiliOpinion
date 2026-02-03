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

    this.client.interceptors.response.use(
      (response) => response,
      (error: AxiosError) => {
        if (error.response) {
          console.error('API Error:', error.response.data)
        } else if (error.request) {
          console.error('Network Error:', error.message)
        } else {
          console.error('Error:', error.message)
        }
        return Promise.reject(error)
      }
    )
  }

  async post<T>(url: string, data?: Record<string, unknown>): Promise<T> {
    const response = await this.client.post<T>(url, data)
    return response.data
  }
}

export const apiClient = new APIClient()
