import axios from 'axios';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      if (window.location.pathname !== '/login') {
        window.location.href = '/login';
      }
    }
    return Promise.reject(error);
  }
);

export const authAPI = {
  register: (data) => api.post('/auth/register', data),
  login: (data) => api.post('/auth/login', data),
  profile: () => api.get('/auth/profile'),
};

export const emergencyAPI = {
  create: (data) => api.post('/emergencies', data),
  getActive: () => api.get('/emergencies/active'),
  getById: (id) => api.get(`/emergencies/${id}`),
  cancel: (id) => api.post(`/emergencies/${id}/cancel`),
  updateLocation: (id, data) => api.post(`/emergencies/${id}/location`, data),
  getLatestLocation: (id) => api.get(`/emergencies/${id}/location/latest`),
  getLocationHistory: (id) => api.get(`/emergencies/${id}/location/history`),
};

export const adminAPI = {
  listActive: () => api.get('/admin/emergencies'),
  search: (params) => api.get('/admin/emergencies/search', { params }),
  getDetails: (id) => api.get(`/admin/emergencies/${id}`),
  resolve: (id) => api.post(`/admin/emergencies/${id}/resolve`),
  stats: () => api.get('/admin/stats'),
};

export default api;
