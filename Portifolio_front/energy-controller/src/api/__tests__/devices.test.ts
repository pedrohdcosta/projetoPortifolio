import { describe, it, expect, vi, beforeEach } from 'vitest'
import { testDeviceConnection } from '../devices'
import api from '../axios'

// Mock the axios instance
vi.mock('../axios', () => ({
  default: {
    get: vi.fn(),
  },
}))

describe('devices API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('testDeviceConnection', () => {
    it('should call the correct endpoint with device id', async () => {
      const mockResponse = { data: { power: 150.5 } }
      vi.mocked(api.get).mockResolvedValue(mockResponse)

      const deviceId = 123
      const result = await testDeviceConnection(deviceId)

      expect(api.get).toHaveBeenCalledWith('/devices/123/read')
      expect(result).toEqual({ power: 150.5 })
    })

    it('should return power value from the API response', async () => {
      const mockResponse = { data: { power: 250.75 } }
      vi.mocked(api.get).mockResolvedValue(mockResponse)

      const result = await testDeviceConnection(42)

      expect(result.power).toBe(250.75)
    })

    it('should throw error when API call fails', async () => {
      const mockError = new Error('Network error')
      vi.mocked(api.get).mockRejectedValue(mockError)

      await expect(testDeviceConnection(1)).rejects.toThrow('Network error')
    })

    it('should handle device not configured error', async () => {
      const mockError = {
        response: {
          status: 400,
          data: { error: 'device not configured for live reads' }
        }
      }
      vi.mocked(api.get).mockRejectedValue(mockError)

      await expect(testDeviceConnection(1)).rejects.toMatchObject({
        response: {
          data: { error: 'device not configured for live reads' }
        }
      })
    })

    it('should handle authentication errors', async () => {
      const mockError = {
        response: {
          status: 401,
          data: { error: 'unauthorized' }
        }
      }
      vi.mocked(api.get).mockRejectedValue(mockError)

      await expect(testDeviceConnection(1)).rejects.toMatchObject({
        response: {
          data: { error: 'unauthorized' }
        }
      })
    })
  })
})
