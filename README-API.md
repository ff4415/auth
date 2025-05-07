# YuzhaLink多语言国际化API文档

## 🌐 概述

YuzhaLink API提供完整的多语言支持和智能错误处理机制，支持中英文无缝切换，适用于全球化应用开发。

**Base URL**: `https://api.yuzhalink.com`

**支持语言**: 
- `en` - English (默认)
- `zh` - 中文 (简体)

**支持国家**: US, CN, GB, CA, AU, HK, TW, SG

## 🚀 快速开始

### 1. 用户注册

```bash
curl -X POST "https://api.yuzhalink.com/auth/signup" \
  -H "Accept-Language: zh-CN,zh;q=0.9,en;q=0.8" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "StrongPass123!",
    "data": {
      "name": "用户姓名",
      "user_language": "zh",
      "country": "CN"
    }
  }'
```

**响应**:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "user_123",
    "email": "user@example.com",
    "user_metadata": {
      "user_language": "zh",
      "country": "CN"
    }
  }
}
```

### 2. 获取设备列表

```bash
curl -X GET "https://api.yuzhalink.com/api/v1/auth/devices" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

**响应** (自动中文):
```json
{
  "code": 0,
  "data": [
    {
      "id": "device_001",
      "name": "pH传感器工作站",
      "status": "online"
    }
  ],
  "message": "成功"
}
```

### 3. 临时切换语言

```bash
curl -X GET "https://api.yuzhalink.com/api/v1/auth/devices?lang=en" \
  -H "Authorization: Bearer chinese_user_token"
```

**响应** (强制英文):
```json
{
  "code": 0,
  "data": [
    {
      "id": "device_001", 
      "name": "pH Sensor Station",
      "status": "online"
    }
  ],
  "message": "Success"
}
```

## 🔧 语言检测机制

### 优先级顺序

1. **查询参数** `?lang=zh` (最高)
2. **自定义头** `X-Language: zh`
3. **JWT Claims** `user_metadata.user_language`
4. **Accept-Language** `Accept-Language: zh-CN,zh;q=0.9`
5. **系统默认** `en` (兜底)

### JavaScript SDK示例

```javascript
class YuzhaAPI {
    constructor(options = {}) {
        this.baseURL = 'https://api.yuzhalink.com';
        this.language = options.language || 'en';
        this.token = options.token;
    }
    
    async request(endpoint, options = {}) {
        const headers = {
            'Content-Type': 'application/json',
            'Accept-Language': this.getAcceptLanguageHeader(),
            ...options.headers
        };
        
        if (this.token) {
            headers.Authorization = `Bearer ${this.token}`;
        }
        
        const response = await fetch(`${this.baseURL}${endpoint}`, {
            ...options,
            headers
        });
        
        const data = await response.json();
        
        if (data.code !== 0) {
            throw new Error(data.message);
        }
        
        return data;
    }
    
    getAcceptLanguageHeader() {
        return this.language === 'zh' ? 
            'zh-CN,zh;q=0.9,en;q=0.8' : 
            'en-US,en;q=0.9';
    }
    
    // 用户管理
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
    
    // 语言管理
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
    
    // 临时语言切换
    withLanguage(language) {
        const tempAPI = Object.create(this);
        tempAPI.language = language;
        return tempAPI;
    }
}

// 使用示例
const api = new YuzhaAPI({ language: 'zh' });

// 注册用户
await api.signup('zhang@example.com', 'password123', {
    name: '张三',
    country: 'CN'
});

// 登录
await api.login('zhang@example.com', 'password123');

// 获取设备 (中文)
const devices = await api.getDevices();
console.log(devices.message); // "成功"

// 临时切换英文
const englishDevices = await api.withLanguage('en').getDevices();
console.log(englishDevices.message); // "Success"

// 永久切换语言
await api.updateLanguage('en', 'US');
```

## 🛡️ 错误处理

### 四级错误体系

| 级别 | 类型 | 示例 | 用户可见 |
|------|------|------|----------|
| **Level-1** | 验证错误 | 必填字段缺失 | ✅ 具体错误 |
| **Level-2** | 业务错误 | 设备不存在 | ✅ 业务消息 |
| **Level-3** | 认证错误 | Token过期 | ⚠️ 统一消息 |
| **Level-4** | 系统错误 | 数据库故障 | ❌ 完全隐藏 |

### 错误处理示例

```javascript
try {
    const result = await api.createDevice({
        type: 'sensor'
        // 缺少必填的name字段
    });
} catch (error) {
    switch (error.code) {
        case 400:
            // Level-1: 显示具体验证错误
            alert(error.message); // "设备名称为必填项"
            break;
        case 401:
            // Level-3: 认证失败，重定向登录
            window.location.href = '/login';
            break;
        case 500:
            // Level-4: 系统错误，温和提示
            alert('系统暂时不可用，请稍后重试');
            break;
    }
}
```

## 📱 常用API端点

### 认证相关

| 方法 | 端点 | 说明 |
|------|------|------|
| POST | `/auth/signup` | 用户注册 |
| POST | `/auth/login` | 用户登录 |
| POST | `/auth/logout` | 用户登出 |

### 用户管理

| 方法 | 端点 | 说明 |
|------|------|------|
| GET | `/api/v1/auth/profile` | 获取用户资料 |
| PUT | `/api/v1/auth/profile/language` | 更新语言偏好 |

### 设备管理

| 方法 | 端点 | 说明 |
|------|------|------|
| GET | `/api/v1/auth/devices` | 获取设备列表 |
| POST | `/api/v1/auth/devices` | 创建设备 |
| GET | `/api/v1/auth/devices/{id}` | 获取设备详情 |
| PUT | `/api/v1/auth/devices/{id}` | 更新设备 |
| DELETE | `/api/v1/auth/devices/{id}` | 删除设备 |

### 系统信息

| 方法 | 端点 | 说明 |
|------|------|------|
| GET | `/api/v1/public/health` | 健康检查 |
| GET | `/api/v1/public/locales` | 支持的语言列表 |
| GET | `/api/v1/public/system/info` | 系统信息 |

## 🔍 调试技巧

### 1. 查看语言检测结果

```bash
# 查看完整请求头
curl -v -X GET "https://api.yuzhalink.com/api/v1/auth/profile" \
  -H "Authorization: Bearer TOKEN" \
  -H "Accept-Language: zh-CN,zh;q=0.9,en;q=0.8"
```

### 2. 测试不同语言

```bash
# 测试中文
curl -X GET "https://api.yuzhalink.com/api/v1/auth/devices?lang=zh" \
  -H "Authorization: Bearer TOKEN"

# 测试英文  
curl -X GET "https://api.yuzhalink.com/api/v1/auth/devices?lang=en" \
  -H "Authorization: Bearer TOKEN"
```

### 3. 验证错误测试

```bash
# 测试验证错误
curl -X POST "https://api.yuzhalink.com/api/v1/auth/devices" \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"type": "sensor"}' # 缺少name字段
```

## 📚 更多文档

- [完整API参考](./docs/api/multilingual-api-guide.md)
- [快速入门指南](./docs/api/quick-start-guide.md)  
- [错误处理指南](./docs/api/error-handling-guide.md)
- [SDK下载](./docs/api/sdks/)

## 💡 最佳实践

1. **自动语言检测**: 注册时基于Accept-Language自动设置用户语言
2. **JWT语言存储**: 用户语言偏好存储在JWT Claims中，减少检测开销
3. **优雅错误处理**: 根据错误级别显示适当的用户消息
4. **临时语言切换**: 使用查询参数`?lang=`进行临时语言切换
5. **缓存策略**: 缓存语言检测结果和翻译消息

## 🆘 技术支持

- **邮箱**: api-support@yuzhalink.com
- **文档**: https://docs.yuzhalink.com
- **状态**: https://status.yuzhalink.com
- **GitHub**: https://github.com/yuzhalink/api-docs

---

**API版本**: v1.0.0  
**文档更新**: 2024-01-15  
**兼容性**: 向后兼容 