# YuzhaLink多语言国际化API使用文档

## 📋 目录

- [1. API概览](#1-api概览)
- [2. 认证与语言检测](#2-认证与语言检测)
- [3. 用户管理API](#3-用户管理api)
- [4. 设备管理API](#4-设备管理api)
- [5. 系统配置API](#5-系统配置api)
- [6. 错误处理](#6-错误处理)
- [7. SDK集成示例](#7-sdk集成示例)
- [8. 最佳实践](#8-最佳实践)

## 1. API概览

### 🌐 基础信息

**Base URL**: `https://api.yuzhalink.com`

**支持的语言**: 
- `en` - English (默认)
- `zh` - 中文

**支持的国家/地区**:
- `US` - 美国 (默认)
- `CN` - 中国
- `GB` - 英国
- `CA` - 加拿大
- `AU` - 澳大利亚
- `HK` - 香港
- `TW` - 台湾
- `SG` - 新加坡

### 🔄 语言检测优先级

| 优先级 | 方法 | 示例 | 说明 |
|--------|------|------|------|
| **1** | 查询参数 | `?lang=zh` | 最高优先级，适用于临时切换 |
| **2** | 自定义头 | `X-Language: zh` | 应用程序级语言设置 |
| **3** | JWT Claims | `user_metadata.user_language` | 用户个人偏好设置 |
| **4** | Accept-Language | `Accept-Language: zh-CN,zh;q=0.9` | 浏览器自动检测 |
| **5** | 系统默认 | `en` | 兜底方案 |

### 📦 统一响应格式

```json
{
  "code": 0,           // 业务状态码：0=成功，其他=错误
  "data": {...},       // 响应数据
  "message": "Success" // 本地化消息
}
```

## 2. 认证与语言检测

### 🔐 用户注册

**接口**: `POST /auth/signup`

**请求示例**:
```bash
curl -X POST "https://api.yuzhalink.com/auth/signup" \
  -H "Accept-Language: zh-CN,zh;q=0.9,en;q=0.8" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "zhang@example.com",
    "password": "StrongPass123!",
    "data": {
      "name": "张三",
      "user_language": "zh",
      "country": "CN"
    }
  }'
```

**成功响应**:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 3600,
  "user": {
    "id": "user_123",
    "email": "zhang@example.com",
    "user_metadata": {
      "name": "张三",
      "user_language": "zh",
      "country": "CN",
      "timezone": "Asia/Shanghai",
      "date_format": "YYYY年MM月DD日",
      "number_format": "zh-CN"
    },
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

**验证错误响应** (中文):
```json
{
  "code": 400,
  "data": null,
  "message": "密码不符合安全要求"
}
```

**验证错误响应** (英文):
```json
{
  "code": 400,
  "data": null,
  "message": "Password does not meet security requirements"
}
```

### 🔑 用户登录

**接口**: `POST /auth/login`

**请求示例**:
```bash
curl -X POST "https://api.yuzhalink.com/auth/login" \
  -H "Accept-Language: en-US,en;q=0.9" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "zhang@example.com",
    "password": "StrongPass123!"
  }'
```

**成功响应**:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 3600,
  "user": {
    "id": "user_123",
    "email": "zhang@example.com",
    "user_metadata": {
      "user_language": "zh",
      "country": "CN"
    }
  }
}
```

## 3. 用户管理API

### 👤 获取用户资料

**接口**: `GET /api/v1/auth/profile`

**请求示例**:
```bash
curl -X GET "https://api.yuzhalink.com/api/v1/auth/profile" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -H "Accept-Language: zh-CN,zh;q=0.9"
```

**成功响应** (中文用户):
```json
{
  "code": 0,
  "data": {
    "id": "user_123",
    "email": "zhang@example.com",
    "name": "张三",
    "user_metadata": {
      "user_language": "zh",
      "country": "CN",
      "timezone": "Asia/Shanghai",
      "date_format": "YYYY年MM月DD日"
    },
    "last_login": "2024-01-15T10:30:00Z"
  },
  "message": "成功"
}
```

### 🔄 更新用户语言偏好

**接口**: `PUT /api/v1/auth/profile/language`

**请求示例**:
```bash
curl -X PUT "https://api.yuzhalink.com/api/v1/auth/profile/language" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -H "Content-Type: application/json" \
  -d '{
    "user_language": "en",
    "country": "US",
    "timezone": "America/New_York"
  }'
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "user_language": "en",
    "country": "US", 
    "timezone": "America/New_York",
    "updated_at": "2024-01-15T11:00:00Z"
  },
  "message": "Language preference updated successfully"
}
```

### 🌍 获取支持的语言和国家列表

**接口**: `GET /api/v1/public/locales`

**请求示例**:
```bash
curl -X GET "https://api.yuzhalink.com/api/v1/public/locales?lang=zh"
```

**响应**:
```json
{
  "code": 0,
  "data": {
    "languages": [
      {
        "code": "en",
        "name": "English",
        "native_name": "English"
      },
      {
        "code": "zh", 
        "name": "Chinese",
        "native_name": "中文"
      }
    ],
    "countries": [
      {
        "code": "US",
        "name": "美国",
        "timezone": "America/New_York"
      },
      {
        "code": "CN",
        "name": "中国", 
        "timezone": "Asia/Shanghai"
      },
      {
        "code": "GB",
        "name": "英国",
        "timezone": "Europe/London"
      }
    ]
  },
  "message": "成功"
}
```

## 4. 设备管理API

### 📱 获取设备列表

**接口**: `GET /api/v1/auth/devices`

**查询参数**:
- `lang` (可选): 强制指定响应语言
- `page` (可选): 页码，默认1
- `limit` (可选): 每页数量，默认10

**请求示例** (中文用户):
```bash
curl -X GET "https://api.yuzhalink.com/api/v1/auth/devices?page=1&limit=5" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "devices": [
      {
        "id": "device_001",
        "name": "pH传感器工作站",
        "type": "ph_sensor",
        "status": "online",
        "location": "实验室A",
        "last_seen": "2024-01-15T10:25:00Z",
        "telemetry": {
          "ph": 7.12,
          "temperature": 25.3,
          "timestamp": "2024-01-15T10:24:45Z"
        }
      },
      {
        "id": "device_002", 
        "name": "温度监控设备",
        "type": "temperature_sensor",
        "status": "offline",
        "location": "实验室B",
        "last_seen": "2024-01-15T08:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 5,
      "total": 12,
      "total_pages": 3
    }
  },
  "message": "成功"
}
```

**临时切换到英文**:
```bash
curl -X GET "https://api.yuzhalink.com/api/v1/auth/devices?lang=en" \
  -H "Authorization: Bearer chinese_user_token"
```

**响应** (强制英文):
```json
{
  "code": 0,
  "data": {
    "devices": [
      {
        "id": "device_001",
        "name": "pH Sensor Station",
        "type": "ph_sensor",
        "status": "online"
      }
    ]
  },
  "message": "Success"
}
```

### 🆕 创建设备

**接口**: `POST /api/v1/auth/devices`

**请求示例**:
```bash
curl -X POST "https://api.yuzhalink.com/api/v1/auth/devices" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -H "Content-Type: application/json" \
  -d '{
    "name": "新的pH传感器",
    "type": "ph_sensor",
    "location": "实验室C",
    "config": {
      "measurement_interval": 30,
      "alert_threshold": 8.5
    }
  }'
```

**成功响应** (中文):
```json
{
  "code": 0,
  "data": {
    "id": "device_003",
    "name": "新的pH传感器",
    "type": "ph_sensor",
    "status": "pending",
    "created_at": "2024-01-15T11:00:00Z"
  },
  "message": "设备创建成功"
}
```

**验证错误** (必填字段缺失):
```json
{
  "code": 400,
  "data": null,
  "message": "设备名称为必填项"
}
```

### 🔧 更新设备信息

**接口**: `PUT /api/v1/auth/devices/{device_id}`

**请求示例**:
```bash
curl -X PUT "https://api.yuzhalink.com/api/v1/auth/devices/device_001" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -H "Content-Type: application/json" \
  -d '{
    "name": "pH传感器工作站-升级版",
    "location": "实验室A-2号位置"
  }'
```

**成功响应**:
```json
{
  "code": 0,
  "data": {
    "id": "device_001",
    "name": "pH传感器工作站-升级版",
    "location": "实验室A-2号位置",
    "updated_at": "2024-01-15T11:30:00Z"
  },
  "message": "设备信息更新成功"
}
```

## 5. 系统配置API

### ⚙️ 获取系统信息

**接口**: `GET /api/v1/public/system/info`

**请求示例**:
```bash
curl -X GET "https://api.yuzhalink.com/api/v1/public/system/info" \
  -H "Accept-Language: zh-CN,zh;q=0.9"
```

**响应**:
```json
{
  "code": 0,
  "data": {
    "version": "1.0.0",
    "build": "2024.01.15",
    "status": "healthy",
    "supported_languages": ["en", "zh"],
    "supported_countries": ["US", "CN", "GB", "CA", "AU", "HK", "TW", "SG"],
    "features": {
      "multi_language": true,
      "error_hiding": true,
      "jwt_language_detection": true
    }
  },
  "message": "系统运行正常"
}
```

### 🏥 健康检查

**接口**: `GET /api/v1/public/health`

**请求示例**:
```bash
curl -X GET "https://api.yuzhalink.com/api/v1/public/health"
```

**响应**:
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T11:45:00Z",
  "services": {
    "database": "healthy",
    "redis": "healthy", 
    "i18n_service": "healthy"
  },
  "message": "All services operational"
}
```

## 6. 错误处理

### 🛡️ 错误级别说明

| 错误级别 | HTTP状态码 | 业务代码 | 说明 | 用户可见信息 |
|----------|------------|----------|------|-------------|
| **Level-1** | 400 | 400 | 验证错误 | 具体验证失败原因 |
| **Level-2** | 400-499 | 对应HTTP码 | 业务错误 | 预定义业务消息 |
| **Level-3** | 401 | 401 | 认证错误 | 统一认证失败消息 |
| **Level-4** | 500 | 500 | 系统错误 | 通用系统错误消息 |

### 📋 常见错误码

#### 认证相关错误

```bash
# JWT Token无效
curl -X GET "https://api.yuzhalink.com/api/v1/auth/profile" \
  -H "Authorization: Bearer invalid_token"
```

**响应** (中文用户):
```json
{
  "code": 401,
  "data": null,
  "message": "未授权访问"
}
```

**响应** (英文用户):
```json
{
  "code": 401,
  "data": null,
  "message": "Unauthorized access"
}
```

#### 权限不足错误

```bash
# 普通用户访问管理员接口
curl -X GET "https://api.yuzhalink.com/api/admin/users" \
  -H "Authorization: Bearer user_token"
```

**响应**:
```json
{
  "code": 403,
  "data": null,
  "message": "权限不足"
}
```

#### 资源不存在错误

```bash
# 获取不存在的设备
curl -X GET "https://api.yuzhalink.com/api/v1/auth/devices/nonexistent_device" \
  -H "Authorization: Bearer valid_token"
```

**响应**:
```json
{
  "code": 404,
  "data": null,
  "message": "设备不存在"
}
```

#### 系统内部错误

```bash
# 数据库连接失败等系统错误
```

**响应** (技术细节完全隐藏):
```json
{
  "code": 500,
  "data": null,
  "message": "系统错误，请稍后重试"
}
```

**开发者日志** (仅后台可见):
```json
{
  "level": "ERROR",
  "timestamp": "2024-01-15T11:50:00Z",
  "error_type": "database_connection_error",
  "original_error": "pq: connection refused to host 192.168.1.100:5432",
  "request_id": "req_12345",
  "user_id": "user_123",
  "path": "/api/v1/auth/devices",
  "method": "GET",
  "user_language": "zh",
  "client_ip": "203.0.113.10"
}
```

## 7. SDK集成示例

### 🚀 JavaScript SDK

```javascript
class YuzhaLinkAPI {
    constructor(options = {}) {
        this.baseURL = options.baseURL || 'https://api.yuzhalink.com';
        this.language = options.language || this.detectBrowserLanguage();
        this.token = options.token;
    }
    
    detectBrowserLanguage() {
        const lang = navigator.language.toLowerCase();
        return lang.startsWith('zh') ? 'zh' : 'en';
    }
    
    getHeaders(customHeaders = {}) {
        const headers = {
            'Content-Type': 'application/json',
            'Accept-Language': this.getAcceptLanguageHeader(),
            ...customHeaders
        };
        
        if (this.token) {
            headers.Authorization = `Bearer ${this.token}`;
        }
        
        return headers;
    }
    
    getAcceptLanguageHeader() {
        switch (this.language) {
            case 'zh':
                return 'zh-CN,zh;q=0.9,en;q=0.8';
            case 'en':
                return 'en-US,en;q=0.9';
            default:
                return 'en-US,en;q=0.9';
        }
    }
    
    async request(endpoint, options = {}) {
        const url = `${this.baseURL}${endpoint}`;
        const config = {
            headers: this.getHeaders(options.headers),
            ...options
        };
        
        try {
            const response = await fetch(url, config);
            const data = await response.json();
            
            if (data.code !== 0) {
                throw new APIError(data.code, data.message, data);
            }
            
            return data;
        } catch (error) {
            if (error instanceof APIError) {
                throw error;
            }
            throw new APIError(500, '网络请求失败', { originalError: error });
        }
    }
    
    // 用户认证
    async signup(email, password, userData = {}) {
        return this.request('/auth/signup', {
            method: 'POST',
            body: JSON.stringify({
                email,
                password,
                data: {
                    user_language: this.language,
                    ...userData
                }
            })
        });
    }
    
    async login(email, password) {
        const result = await this.request('/auth/login', {
            method: 'POST',
            body: JSON.stringify({ email, password })
        });
        
        this.token = result.access_token;
        return result;
    }
    
    // 用户管理
    async getProfile() {
        return this.request('/api/v1/auth/profile');
    }
    
    async updateLanguage(language, country) {
        const result = await this.request('/api/v1/auth/profile/language', {
            method: 'PUT',
            body: JSON.stringify({
                user_language: language,
                country: country
            })
        });
        
        this.language = language;
        return result;
    }
    
    // 设备管理
    async getDevices(options = {}) {
        const params = new URLSearchParams();
        if (options.page) params.set('page', options.page);
        if (options.limit) params.set('limit', options.limit);
        if (options.forceLanguage) params.set('lang', options.forceLanguage);
        
        const query = params.toString() ? `?${params.toString()}` : '';
        return this.request(`/api/v1/auth/devices${query}`);
    }
    
    async createDevice(deviceData) {
        return this.request('/api/v1/auth/devices', {
            method: 'POST',
            body: JSON.stringify(deviceData)
        });
    }
    
    async updateDevice(deviceId, updateData) {
        return this.request(`/api/v1/auth/devices/${deviceId}`, {
            method: 'PUT',
            body: JSON.stringify(updateData)
        });
    }
    
    // 临时语言切换
    withLanguage(language) {
        const tempAPI = Object.create(this);
        tempAPI.language = language;
        return tempAPI;
    }
}

class APIError extends Error {
    constructor(code, message, data) {
        super(message);
        this.code = code;
        this.data = data;
    }
}

// 使用示例
const api = new YuzhaLinkAPI({
    baseURL: 'https://api.yuzhalink.com',
    language: 'zh'
});

// 用户注册
try {
    const result = await api.signup('zhang@example.com', 'StrongPass123!', {
        name: '张三',
        country: 'CN'
    });
    console.log('注册成功:', result.message);
} catch (error) {
    console.error('注册失败:', error.message);
}

// 登录
try {
    await api.login('zhang@example.com', 'StrongPass123!');
    console.log('登录成功');
} catch (error) {
    console.error('登录失败:', error.message);
}

// 获取设备列表
try {
    const devices = await api.getDevices({ page: 1, limit: 10 });
    console.log('设备列表:', devices.data.devices);
} catch (error) {
    console.error('获取设备失败:', error.message);
}

// 临时切换到英文
try {
    const englishDevices = await api.withLanguage('en').getDevices();
    console.log('English devices:', englishDevices.message); // "Success"
} catch (error) {
    console.error('Failed to get devices:', error.message);
}

// 更新用户语言偏好
try {
    await api.updateLanguage('en', 'US');
    console.log('Language updated to English');
} catch (error) {
    console.error('Language update failed:', error.message);
}
```

### 🐍 Python SDK

```python
import requests
import json
from typing import Optional, Dict, Any

class YuzhaLinkAPI:
    def __init__(self, base_url: str = "https://api.yuzhalink.com", 
                 language: str = "en", token: Optional[str] = None):
        self.base_url = base_url
        self.language = language
        self.token = token
        self.session = requests.Session()
    
    def get_headers(self, custom_headers: Optional[Dict] = None) -> Dict[str, str]:
        headers = {
            'Content-Type': 'application/json',
            'Accept-Language': self._get_accept_language_header()
        }
        
        if self.token:
            headers['Authorization'] = f'Bearer {self.token}'
        
        if custom_headers:
            headers.update(custom_headers)
        
        return headers
    
    def _get_accept_language_header(self) -> str:
        if self.language == 'zh':
            return 'zh-CN,zh;q=0.9,en;q=0.8'
        else:
            return 'en-US,en;q=0.9'
    
    def request(self, endpoint: str, method: str = 'GET', 
                data: Optional[Dict] = None, params: Optional[Dict] = None) -> Dict[str, Any]:
        url = f"{self.base_url}{endpoint}"
        headers = self.get_headers()
        
        try:
            response = self.session.request(
                method=method,
                url=url,
                headers=headers,
                json=data,
                params=params
            )
            
            result = response.json()
            
            if result.get('code', 0) != 0:
                raise APIError(result.get('code', 500), result.get('message', 'Unknown error'))
            
            return result
        
        except requests.RequestException as e:
            raise APIError(500, f"Network error: {str(e)}")
    
    # 用户认证
    def signup(self, email: str, password: str, user_data: Optional[Dict] = None) -> Dict[str, Any]:
        data = {
            'email': email,
            'password': password,
            'data': {
                'user_language': self.language,
                **(user_data or {})
            }
        }
        
        result = self.request('/auth/signup', 'POST', data)
        self.token = result.get('access_token')
        return result
    
    def login(self, email: str, password: str) -> Dict[str, Any]:
        data = {'email': email, 'password': password}
        result = self.request('/auth/login', 'POST', data)
        self.token = result.get('access_token')
        return result
    
    # 用户管理
    def get_profile(self) -> Dict[str, Any]:
        return self.request('/api/v1/auth/profile')
    
    def update_language(self, language: str, country: str) -> Dict[str, Any]:
        data = {
            'user_language': language,
            'country': country
        }
        result = self.request('/api/v1/auth/profile/language', 'PUT', data)
        self.language = language
        return result
    
    # 设备管理
    def get_devices(self, page: int = 1, limit: int = 10, 
                   force_language: Optional[str] = None) -> Dict[str, Any]:
        params = {'page': page, 'limit': limit}
        if force_language:
            params['lang'] = force_language
        
        return self.request('/api/v1/auth/devices', params=params)
    
    def create_device(self, device_data: Dict[str, Any]) -> Dict[str, Any]:
        return self.request('/api/v1/auth/devices', 'POST', device_data)
    
    def update_device(self, device_id: str, update_data: Dict[str, Any]) -> Dict[str, Any]:
        return self.request(f'/api/v1/auth/devices/{device_id}', 'PUT', update_data)
    
    # 临时语言切换
    def with_language(self, language: str) -> 'YuzhaLinkAPI':
        temp_api = YuzhaLinkAPI(self.base_url, language, self.token)
        temp_api.session = self.session
        return temp_api

class APIError(Exception):
    def __init__(self, code: int, message: str):
        self.code = code
        self.message = message
        super().__init__(f"API Error {code}: {message}")

# 使用示例
if __name__ == "__main__":
    # 初始化API客户端
    api = YuzhaLinkAPI(language='zh')
    
    try:
        # 用户注册
        result = api.signup('zhang@example.com', 'StrongPass123!', {
            'name': '张三',
            'country': 'CN'
        })
        print(f"注册成功: {result['message']}")
        
        # 获取设备列表
        devices = api.get_devices(page=1, limit=5)
        print(f"设备数量: {len(devices['data']['devices'])}")
        print(f"响应消息: {devices['message']}")
        
        # 临时切换到英文
        english_devices = api.with_language('en').get_devices()
        print(f"English message: {english_devices['message']}")
        
        # 更新语言偏好
        api.update_language('en', 'US')
        print("语言偏好已更新为英文")
        
    except APIError as e:
        print(f"API错误: {e.message}")
    except Exception as e:
        print(f"其他错误: {str(e)}")
```

## 8. 最佳实践

### ✅ 推荐做法

#### 1. **语言检测策略**
```javascript
// 优先级使用建议
const api = new YuzhaLinkAPI();

// 方式1: 让系统自动检测（推荐）
await api.getDevices(); // 基于JWT、Accept-Language自动检测

// 方式2: 临时切换语言
await api.getDevices({ forceLanguage: 'en' }); // 查询参数优先级最高

// 方式3: 永久更新用户偏好
await api.updateLanguage('zh', 'CN'); // 更新JWT Claims
```

#### 2. **错误处理**
```javascript
// 统一错误处理
try {
    const result = await api.createDevice(deviceData);
    showSuccess(result.message); // 显示本地化成功消息
} catch (error) {
    if (error.code === 400) {
        showValidationError(error.message); // 显示验证错误
    } else if (error.code === 401) {
        redirectToLogin(); // 重定向到登录页
    } else {
        showGenericError(error.message); // 显示通用错误
    }
}
```

#### 3. **用户体验优化**
```javascript
// 智能语言检测和缓存
class SmartLanguageAPI extends YuzhaLinkAPI {
    constructor(options = {}) {
        super(options);
        this.loadSavedLanguage();
    }
    
    loadSavedLanguage() {
        // 1. 检查本地存储
        const saved = localStorage.getItem('user_language');
        if (saved) {
            this.language = saved;
            return;
        }
        
        // 2. 检查浏览器语言
        this.language = this.detectBrowserLanguage();
    }
    
    async updateLanguage(language, country) {
        const result = super.updateLanguage(language, country);
        
        // 保存到本地存储
        localStorage.setItem('user_language', language);
        localStorage.setItem('user_country', country);
        
        return result;
    }
}
```

### ❌ 避免的做法

#### 1. **不要硬编码语言**
```javascript
// ❌ 错误做法
function showMessage() {
    alert("操作成功"); // 硬编码中文
}

// ✅ 正确做法
function showMessage(message) {
    alert(message); // 使用API返回的本地化消息
}
```

#### 2. **不要忽略错误处理**
```javascript
// ❌ 错误做法
const devices = await api.getDevices(); // 没有错误处理

// ✅ 正确做法
try {
    const devices = await api.getDevices();
    // 处理成功逻辑
} catch (error) {
    // 处理错误
}
```

#### 3. **不要频繁切换语言**
```javascript
// ❌ 错误做法 - 频繁创建新实例
for (let i = 0; i < requests.length; i++) {
    const tempAPI = new YuzhaLinkAPI({ language: 'en' });
    await tempAPI.getDevices();
}

// ✅ 正确做法 - 复用临时实例
const englishAPI = api.withLanguage('en');
for (let i = 0; i < requests.length; i++) {
    await englishAPI.getDevices();
}
```

### 🔧 调试技巧

#### 1. **语言检测调试**
```bash
# 查看完整的请求头
curl -v -X GET "https://api.yuzhalink.com/api/v1/auth/profile" \
  -H "Authorization: Bearer TOKEN" \
  -H "Accept-Language: zh-CN,zh;q=0.9,en;q=0.8"
```

#### 2. **错误追踪**
```javascript
// 添加请求ID追踪
class DebugAPI extends YuzhaLinkAPI {
    async request(endpoint, options = {}) {
        const requestId = this.generateRequestId();
        console.log(`[${requestId}] Request: ${endpoint}`);
        
        try {
            const result = await super.request(endpoint, options);
            console.log(`[${requestId}] Success:`, result.message);
            return result;
        } catch (error) {
            console.error(`[${requestId}] Error:`, error.message);
            throw error;
        }
    }
    
    generateRequestId() {
        return Math.random().toString(36).substr(2, 9);
    }
}
```

### 📊 性能优化建议

#### 1. **缓存策略**
```javascript
// 缓存支持的语言和国家列表
class CachedAPI extends YuzhaLinkAPI {
    constructor(options = {}) {
        super(options);
        this.localesCache = null;
        this.cacheExpiry = null;
    }
    
    async getLocales() {
        const now = Date.now();
        
        // 检查缓存是否有效（1小时过期）
        if (this.localesCache && this.cacheExpiry > now) {
            return { data: this.localesCache, from_cache: true };
        }
        
        // 获取新数据
        const result = await this.request('/api/v1/public/locales');
        this.localesCache = result.data;
        this.cacheExpiry = now + 3600000; // 1小时
        
        return result;
    }
}
```

#### 2. **请求合并**
```javascript
// 批量设备状态查询
async function getBatchDeviceStatus(deviceIds) {
    return api.request('/api/v1/auth/devices/batch/status', {
        method: 'POST',
        body: JSON.stringify({ device_ids: deviceIds })
    });
}
```

---

## 📞 技术支持

**API文档版本**: v1.0.0  
**最后更新**: 2024-01-15  
**技术支持**: api-support@yuzhalink.com  
**状态页面**: https://status.yuzhalink.com

### 📋 常见问题

**Q: 如何检查我的JWT Token中包含的语言偏好？**  
A: 可以解码JWT Token查看`user_metadata.user_language`字段。

**Q: 为什么我设置了Accept-Language但响应还是英文？**  
A: 请检查JWT Token中是否包含`user_language`字段，JWT中的语言偏好优先级更高。

**Q: 如何永久更改用户的语言偏好？**  
A: 使用`PUT /api/v1/auth/profile/language`接口更新用户的语言偏好。

**Q: 系统错误为什么没有详细信息？**  
A: 出于安全考虑，系统级错误会隐藏技术细节，详细信息记录在服务器日志中。 