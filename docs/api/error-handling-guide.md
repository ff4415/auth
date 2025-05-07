# 错误处理与安全隐藏指南

## 🛡️ 四级错误安全隐藏体系

### 概览

YuzhaLink API采用四级错误分类和安全隐藏机制，确保在提供有用错误信息的同时保护系统安全。

| 级别 | 错误类型 | 用户可见性 | 日志记录 | 安全等级 |
|------|----------|------------|----------|----------|
| **🟢 Level-1** | 验证错误 | 完全可见 | INFO | 低风险 |
| **🔶 Level-2** | 业务错误 | 部分可见 | WARN | 中风险 |
| **🔶 Level-3** | 认证错误 | 统一消息 | ERROR | 中风险 |
| **🔴 Level-4** | 系统错误 | 完全隐藏 | ERROR | 高风险 |

## 🟢 Level-1: 验证错误 (用户输入错误)

### 特征
- 用户输入导致的错误
- 不涉及系统安全
- 需要向用户显示具体错误信息帮助修正

### 示例场景

#### 表单验证失败

**请求**:
```bash
curl -X POST "https://api.yuzhalink.com/api/v1/auth/devices" \
  -H "Authorization: Bearer valid_token" \
  -H "Accept-Language: zh-CN,zh;q=0.9" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "ph_sensor"
  }'
```

**响应** (中文用户):
```json
{
  "code": 400,
  "data": null,
  "message": "设备名称为必填项"
}
```

**响应** (英文用户):
```json
{
  "code": 400,
  "data": null,
  "message": "Device name is required"
}
```

#### 多字段验证错误

**请求**:
```bash
curl -X POST "https://api.yuzhalink.com/auth/signup" \
  -H "Accept-Language: en-US,en;q=0.9" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "invalid-email",
    "password": "123"
  }'
```

**响应**:
```json
{
  "code": 400,
  "data": null,
  "message": "Email format is invalid; Password does not meet security requirements"
}
```

### 开发者处理建议

```javascript
// 前端处理验证错误
try {
    const result = await api.createDevice(deviceData);
} catch (error) {
    if (error.code === 400) {
        // Level-1错误：显示具体错误信息
        showValidationError(error.message);
        highlightErrorFields(error.data?.fields);
    }
}

function showValidationError(message) {
    // 直接显示API返回的本地化错误消息
    document.getElementById('error-message').textContent = message;
    document.getElementById('error-container').style.display = 'block';
}
```

## 🔶 Level-2: HTTP业务错误

### 特征
- 业务逻辑相关错误
- 预定义的业务场景
- 显示用户友好的业务消息，隐藏技术细节

### 示例场景

#### 资源不存在

**请求**:
```bash
curl -X GET "https://api.yuzhalink.com/api/v1/auth/devices/nonexistent_device" \
  -H "Authorization: Bearer valid_token" \
  -H "Accept-Language: zh-CN,zh;q=0.9"
```

**响应**:
```json
{
  "code": 404,
  "data": null,
  "message": "设备不存在"
}
```

**内部日志**:
```json
{
  "level": "WARN",
  "timestamp": "2024-01-15T10:30:00Z",
  "message": "Device not found",
  "device_id": "nonexistent_device",
  "user_id": "user_123",
  "request_id": "req_456",
  "path": "/api/v1/auth/devices/nonexistent_device"
}
```

#### 权限不足

**请求**:
```bash
curl -X DELETE "https://api.yuzhalink.com/api/v1/auth/devices/device_001" \
  -H "Authorization: Bearer limited_user_token"
```

**响应**:
```json
{
  "code": 403,
  "data": null,
  "message": "权限不足"
}
```

#### 业务规则冲突

**请求** (设备正在使用中，无法删除):
```bash
curl -X DELETE "https://api.yuzhalink.com/api/v1/auth/devices/device_001" \
  -H "Authorization: Bearer admin_token"
```

**响应**:
```json
{
  "code": 409,
  "data": null,
  "message": "设备正在使用中，无法删除"
}
```

### 开发者处理建议

```javascript
// 业务错误处理
try {
    const result = await api.deleteDevice(deviceId);
    showSuccess(result.message);
} catch (error) {
    switch (error.code) {
        case 403:
            showError('您没有权限执行此操作');
            break;
        case 404:
            showError('设备不存在或已被删除');
            break;
        case 409:
            showError(error.message); // 显示具体的业务冲突信息
            break;
        default:
            showError('操作失败，请稍后重试');
    }
}
```

## 🔶 Level-3: 认证/OAuth错误

### 特征
- 认证和授权相关错误
- 统一为401错误码
- 不暴露具体的认证失败原因

### 示例场景

#### JWT Token过期

**请求**:
```bash
curl -X GET "https://api.yuzhalink.com/api/v1/auth/profile" \
  -H "Authorization: Bearer expired_token"
```

**用户看到的响应**:
```json
{
  "code": 401,
  "data": null,
  "message": "未授权访问"
}
```

**内部详细日志**:
```json
{
  "level": "ERROR",
  "timestamp": "2024-01-15T10:35:00Z",
  "message": "JWT token expired",
  "error_type": "token_expired",
  "original_error": "Token expired at 2024-01-15T09:30:00Z",
  "token_issued_at": "2024-01-15T08:30:00Z",
  "token_expires_at": "2024-01-15T09:30:00Z",
  "user_id": "user_123",
  "request_id": "req_789",
  "client_ip": "203.0.113.10",
  "user_agent": "Mozilla/5.0...",
  "path": "/api/v1/auth/profile"
}
```

#### JWT签名无效

**请求**:
```bash
curl -X GET "https://api.yuzhalink.com/api/v1/auth/profile" \
  -H "Authorization: Bearer invalid_signature_token"
```

**用户响应** (统一处理):
```json
{
  "code": 401,
  "data": null,
  "message": "Unauthorized access"
}
```

#### OAuth错误

**请求**:
```bash
curl -X POST "https://api.yuzhalink.com/auth/oauth/callback" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "invalid_oauth_code",
    "state": "random_state"
  }'
```

**用户响应**:
```json
{
  "code": 401,
  "data": null,
  "message": "认证失败"
}
```

### 开发者处理建议

```javascript
// 认证错误统一处理
class APIClient {
    async request(endpoint, options = {}) {
        try {
            const response = await fetch(endpoint, options);
            const data = await response.json();
            
            if (data.code === 401) {
                // Level-3错误：统一认证失败处理
                this.handleAuthenticationError();
                throw new AuthError(data.message);
            }
            
            return data;
        } catch (error) {
            if (error instanceof AuthError) {
                throw error;
            }
            throw new APIError('Network error');
        }
    }
    
    handleAuthenticationError() {
        // 清除本地token
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        
        // 重定向到登录页
        window.location.href = '/login';
    }
}

class AuthError extends Error {
    constructor(message) {
        super(message);
        this.name = 'AuthError';
    }
}
```

## 🔴 Level-4: 系统内部错误

### 特征
- 系统级别的技术错误
- 完全隐藏错误详情
- 统一返回通用错误消息
- 详细错误信息仅记录在服务器日志

### 示例场景

#### 数据库连接失败

**请求**:
```bash
curl -X GET "https://api.yuzhalink.com/api/v1/auth/devices" \
  -H "Authorization: Bearer valid_token"
```

**用户看到的响应**:
```json
{
  "code": 500,
  "data": null,
  "message": "系统错误，请稍后重试"
}
```

**内部详细日志**:
```json
{
  "level": "ERROR",
  "timestamp": "2024-01-15T10:40:00Z",
  "message": "Database connection failed",
  "error_type": "database_connection_error",
  "original_error": "pq: connection refused to host 192.168.1.100:5432",
  "database_host": "192.168.1.100",
  "database_port": 5432,
  "connection_pool_status": "exhausted",
  "retry_attempts": 3,
  "request_id": "req_abc123",
  "user_id": "user_123",
  "user_language": "zh",
  "path": "/api/v1/auth/devices",
  "method": "GET",
  "query_params": "page=1&limit=10",
  "client_ip": "203.0.113.10",
  "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
  "stack_trace": "goroutine 42 [running]:\n...",
  "system_metrics": {
    "cpu_usage": "85%",
    "memory_usage": "92%",
    "disk_usage": "78%"
  }
}
```

#### 第三方服务异常

**场景**: ThingsBoard API调用失败

**用户响应**:
```json
{
  "code": 500,
  "data": null,
  "message": "Service temporarily unavailable"
}
```

**内部日志**:
```json
{
  "level": "ERROR",
  "timestamp": "2024-01-15T10:45:00Z",
  "message": "ThingsBoard API call failed",
  "error_type": "external_service_error",
  "service_name": "thingsboard",
  "service_endpoint": "http://thingsboard.local:8080/api/plugins/telemetry",
  "original_error": "dial tcp 10.0.1.50:8080: connection timed out",
  "timeout_duration": "30s",
  "retry_count": 2,
  "request_id": "req_def456"
}
```

#### 内存不足错误

**用户响应**:
```json
{
  "code": 500,
  "data": null,
  "message": "系统繁忙，请稍后重试"
}
```

**内部日志**:
```json
{
  "level": "CRITICAL",
  "timestamp": "2024-01-15T10:50:00Z",
  "message": "Out of memory error",
  "error_type": "system_resource_exhausted",
  "memory_usage": "98%",
  "available_memory": "128MB",
  "requested_memory": "500MB",
  "garbage_collection_stats": {...},
  "active_goroutines": 2500,
  "request_id": "req_ghi789"
}
```

### 开发者处理建议

```javascript
// Level-4错误处理
class RobustAPIClient {
    async request(endpoint, options = {}) {
        try {
            return await this.makeRequest(endpoint, options);
        } catch (error) {
            if (error.code === 500) {
                // Level-4错误：系统错误
                return this.handleSystemError(error, endpoint, options);
            }
            throw error;
        }
    }
    
    async handleSystemError(error, endpoint, options) {
        // 1. 显示用户友好的错误消息
        this.showUserFriendlyError(error.message);
        
        // 2. 记录客户端错误信息（用于调试）
        this.logClientError({
            endpoint,
            timestamp: new Date().toISOString(),
            userAgent: navigator.userAgent,
            error: error.message
        });
        
        // 3. 实施重试策略
        if (this.shouldRetry(endpoint)) {
            await this.delay(this.getRetryDelay());
            return this.request(endpoint, { ...options, _retry: true });
        }
        
        throw error;
    }
    
    showUserFriendlyError(message) {
        // 显示温和的错误消息
        const notification = document.createElement('div');
        notification.className = 'error-notification';
        notification.innerHTML = `
            <div class="error-icon">⚠️</div>
            <div class="error-text">${message}</div>
            <div class="error-suggestion">我们正在处理这个问题，请稍后重试</div>
        `;
        document.body.appendChild(notification);
        
        // 自动消失
        setTimeout(() => {
            document.body.removeChild(notification);
        }, 5000);
    }
    
    logClientError(errorInfo) {
        // 发送客户端错误信息到监控系统（不包含敏感信息）
        if (window.analytics) {
            window.analytics.track('API_Error', {
                endpoint: errorInfo.endpoint,
                timestamp: errorInfo.timestamp,
                userAgent: errorInfo.userAgent
            });
        }
    }
    
    shouldRetry(endpoint) {
        // 某些端点支持重试
        const retryableEndpoints = [
            '/api/v1/auth/devices',
            '/api/v1/auth/profile'
        ];
        return retryableEndpoints.includes(endpoint);
    }
    
    getRetryDelay() {
        // 指数退避重试
        return Math.min(1000 * Math.pow(2, this.retryCount || 0), 10000);
    }
    
    delay(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }
}
```

## 🔍 错误调试指南

### 开发环境调试

```javascript
// 开发环境错误调试工具
class DebugAPIClient extends APIClient {
    constructor(options = {}) {
        super(options);
        this.debugMode = options.debug || process.env.NODE_ENV === 'development';
    }
    
    async request(endpoint, options = {}) {
        if (this.debugMode) {
            console.group(`🔍 API Request: ${endpoint}`);
            console.log('Headers:', options.headers);
            console.log('Body:', options.body);
        }
        
        try {
            const result = await super.request(endpoint, options);
            
            if (this.debugMode) {
                console.log('✅ Success:', result.message);
                console.groupEnd();
            }
            
            return result;
        } catch (error) {
            if (this.debugMode) {
                console.error('❌ Error:', error);
                this.analyzeError(error);
                console.groupEnd();
            }
            throw error;
        }
    }
    
    analyzeError(error) {
        console.group('🔬 Error Analysis');
        
        switch (error.code) {
            case 400:
                console.log('📝 Validation Error - Check your input data');
                console.log('💡 Tip: Verify required fields and data formats');
                break;
            case 401:
                console.log('🔐 Authentication Error - Check your token');
                console.log('💡 Tip: Token may be expired or invalid');
                break;
            case 403:
                console.log('🚫 Permission Error - Insufficient privileges');
                console.log('💡 Tip: Contact admin for proper permissions');
                break;
            case 404:
                console.log('🔍 Not Found - Resource does not exist');
                console.log('💡 Tip: Check the resource ID');
                break;
            case 500:
                console.log('⚠️ System Error - Server-side issue');
                console.log('💡 Tip: Try again later or contact support');
                break;
        }
        
        console.groupEnd();
    }
}

// 使用调试客户端
const debugAPI = new DebugAPIClient({ debug: true });
```

### 生产环境监控

```javascript
// 生产环境错误监控
class MonitoredAPIClient extends APIClient {
    constructor(options = {}) {
        super(options);
        this.errorMetrics = {
            level1: 0, // 验证错误
            level2: 0, // 业务错误
            level3: 0, // 认证错误
            level4: 0  // 系统错误
        };
    }
    
    async request(endpoint, options = {}) {
        const startTime = Date.now();
        
        try {
            const result = await super.request(endpoint, options);
            
            // 记录成功指标
            this.recordMetric('api_success', {
                endpoint,
                duration: Date.now() - startTime
            });
            
            return result;
        } catch (error) {
            // 记录错误指标
            this.recordError(error, endpoint, Date.now() - startTime);
            throw error;
        }
    }
    
    recordError(error, endpoint, duration) {
        // 分类记录错误
        let level;
        if (error.code === 400) level = 'level1';
        else if (error.code >= 400 && error.code < 500) level = 'level2';
        else if (error.code === 401) level = 'level3';
        else level = 'level4';
        
        this.errorMetrics[level]++;
        
        // 发送到监控系统
        this.recordMetric('api_error', {
            endpoint,
            error_code: error.code,
            error_level: level,
            duration,
            user_language: this.language
        });
    }
    
    recordMetric(event, data) {
        // 发送到分析平台
        if (window.analytics) {
            window.analytics.track(event, data);
        }
        
        // 发送到错误监控服务
        if (window.Sentry && event === 'api_error') {
            window.Sentry.captureMessage(`API Error: ${data.endpoint}`, {
                level: data.error_level === 'level4' ? 'error' : 'warning',
                extra: data
            });
        }
    }
    
    getErrorMetrics() {
        return { ...this.errorMetrics };
    }
}
```

## 📊 错误处理最佳实践

### 1. 用户体验优化

```javascript
// 优雅的错误处理UI
class ErrorHandler {
    static show(error) {
        switch (error.code) {
            case 400:
                return this.showValidationError(error.message);
            case 401:
                return this.showAuthError();
            case 403:
                return this.showPermissionError();
            case 404:
                return this.showNotFoundError();
            case 500:
                return this.showSystemError();
            default:
                return this.showGenericError();
        }
    }
    
    static showValidationError(message) {
        // 显示具体的验证错误，帮助用户修正
        return {
            type: 'warning',
            title: '输入错误',
            message: message,
            action: '请检查输入内容',
            autoHide: false
        };
    }
    
    static showAuthError() {
        // 认证错误：引导用户重新登录
        return {
            type: 'error',
            title: '认证过期',
            message: '您的登录已过期，请重新登录',
            action: '立即登录',
            redirect: '/login'
        };
    }
    
    static showSystemError() {
        // 系统错误：温和的提示，不恐慌用户
        return {
            type: 'info',
            title: '服务暂时不可用',
            message: '我们正在修复这个问题，请稍后重试',
            action: '稍后重试',
            autoRetry: true
        };
    }
}
```

### 2. 错误恢复策略

```javascript
// 智能错误恢复
class ResilientAPIClient extends APIClient {
    async callWithFallback(primaryEndpoint, fallbackEndpoint, options = {}) {
        try {
            return await this.request(primaryEndpoint, options);
        } catch (error) {
            if (error.code === 500 && fallbackEndpoint) {
                console.warn('Primary endpoint failed, trying fallback...');
                return await this.request(fallbackEndpoint, options);
            }
            throw error;
        }
    }
    
    async getDevicesWithFallback() {
        // 主要端点失败时使用缓存端点
        return this.callWithFallback(
            '/api/v1/auth/devices',
            '/api/v1/auth/devices/cached'
        );
    }
}
```

### 3. 多语言错误消息

```javascript
// 错误消息本地化
class LocalizedErrorHandler {
    constructor(language = 'en') {
        this.language = language;
        this.messages = {
            'en': {
                'network_error': 'Network connection failed',
                'retry_suggestion': 'Please check your connection and try again',
                'contact_support': 'If the problem persists, contact support'
            },
            'zh': {
                'network_error': '网络连接失败',
                'retry_suggestion': '请检查网络连接并重试',
                'contact_support': '如果问题持续存在，请联系客服'
            }
        };
    }
    
    getMessage(key) {
        return this.messages[this.language]?.[key] || 
               this.messages['en'][key] || 
               key;
    }
    
    handleNetworkError() {
        return {
            title: this.getMessage('network_error'),
            message: this.getMessage('retry_suggestion'),
            footer: this.getMessage('contact_support')
        };
    }
}
```

这个错误处理指南确保了在提供良好用户体验的同时，保护系统的安全性和稳定性。开发者应该根据错误级别采用相应的处理策略，为用户提供准确、有用且安全的错误信息。 