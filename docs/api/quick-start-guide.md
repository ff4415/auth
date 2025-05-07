# YuzhaLink API 快速入门指南

## 🚀 5分钟快速集成

### 步骤1: 获取API密钥

```bash
# 注册开发者账号
curl -X POST "https://api.yuzhalink.com/auth/signup" \
  -H "Accept-Language: zh-CN,zh;q=0.9,en;q=0.8" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "developer@yourcompany.com",
    "password": "YourStrongPassword123!",
    "data": {
      "name": "开发者姓名",
      "user_language": "zh",
      "country": "CN"
    }
  }'
```

### 步骤2: 基础API调用

```javascript
// 基础设置
const API_BASE = 'https://api.yuzhalink.com';
const token = 'your_access_token_here';

// 简单的API调用
async function callAPI(endpoint, options = {}) {
    const response = await fetch(`${API_BASE}${endpoint}`, {
        headers: {
            'Authorization': `Bearer ${token}`,
            'Accept-Language': 'zh-CN,zh;q=0.9,en;q=0.8',
            'Content-Type': 'application/json',
            ...options.headers
        },
        ...options
    });
    
    return response.json();
}

// 获取设备列表
const devices = await callAPI('/api/v1/auth/devices');
console.log(devices.message); // "成功" 或 "Success"
```

### 步骤3: 错误处理

```javascript
try {
    const result = await callAPI('/api/v1/auth/devices', {
        method: 'POST',
        body: JSON.stringify({
            name: '新设备',
            type: 'sensor'
        })
    });
    
    alert(result.message); // 显示本地化成功消息
} catch (error) {
    alert('请求失败，请稍后重试'); // 用户友好的错误提示
}
```

## 🎯 常用场景

### 场景1: 用户注册时的语言检测

```html
<!-- HTML表单 -->
<form id="signupForm">
    <input type="email" name="email" placeholder="邮箱地址" required>
    <input type="password" name="password" placeholder="密码" required>
    <select name="language">
        <option value="zh">中文</option>
        <option value="en">English</option>
    </select>
    <button type="submit">注册</button>
</form>

<script>
document.getElementById('signupForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const formData = new FormData(e.target);
    const data = Object.fromEntries(formData);
    
    try {
        const result = await fetch('/auth/signup', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                // 浏览器自动发送，系统自动检测
                'Accept-Language': navigator.language
            },
            body: JSON.stringify({
                email: data.email,
                password: data.password,
                data: {
                    user_language: data.language,
                    country: data.language === 'zh' ? 'CN' : 'US'
                }
            })
        });
        
        const result_data = await response.json();
        
        if (result_data.code === 0) {
            alert(result_data.message); // 本地化成功消息
            localStorage.setItem('token', result_data.access_token);
        } else {
            alert(result_data.message); // 本地化错误消息
        }
    } catch (error) {
        alert('网络错误，请稍后重试');
    }
});
</script>
```

### 场景2: 动态语言切换

```javascript
class LanguageSwitcher {
    constructor() {
        this.currentLanguage = localStorage.getItem('language') || 'en';
        this.token = localStorage.getItem('token');
    }
    
    async switchLanguage(newLanguage) {
        try {
            // 更新用户偏好
            const response = await fetch('/api/v1/auth/profile/language', {
                method: 'PUT',
                headers: {
                    'Authorization': `Bearer ${this.token}`,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    user_language: newLanguage,
                    country: newLanguage === 'zh' ? 'CN' : 'US'
                })
            });
            
            const result = await response.json();
            
            if (result.code === 0) {
                this.currentLanguage = newLanguage;
                localStorage.setItem('language', newLanguage);
                
                // 重新加载页面内容
                this.reloadPageContent();
                
                alert(result.message); // "语言偏好已更新" 或 "Language preference updated"
            }
        } catch (error) {
            alert('语言切换失败');
        }
    }
    
    async reloadPageContent() {
        // 重新获取设备列表
        const devices = await this.getDevices();
        this.updateDeviceList(devices);
    }
    
    async getDevices() {
        const response = await fetch('/api/v1/auth/devices', {
            headers: {
                'Authorization': `Bearer ${this.token}`,
                // 不需要设置Accept-Language，JWT中已包含用户偏好
            }
        });
        
        return response.json();
    }
    
    updateDeviceList(devicesData) {
        const container = document.getElementById('device-list');
        container.innerHTML = devicesData.data.devices.map(device => `
            <div class="device-item">
                <h3>${device.name}</h3>
                <p>状态: ${device.status}</p>
            </div>
        `).join('');
    }
}

// 使用示例
const languageSwitcher = new LanguageSwitcher();

// 添加语言切换按钮事件
document.getElementById('lang-zh').addEventListener('click', () => {
    languageSwitcher.switchLanguage('zh');
});

document.getElementById('lang-en').addEventListener('click', () => {
    languageSwitcher.switchLanguage('en');
});
```

### 场景3: 移动应用集成 (React Native)

```javascript
// api.js
import AsyncStorage from '@react-native-async-storage/async-storage';

class MobileAPI {
    constructor() {
        this.baseURL = 'https://api.yuzhalink.com';
        this.initializeLanguage();
    }
    
    async initializeLanguage() {
        // 1. 检查存储的用户偏好
        const savedLang = await AsyncStorage.getItem('user_language');
        if (savedLang) {
            this.language = savedLang;
            return;
        }
        
        // 2. 使用设备语言
        import { getLocales } from 'expo-localization';
        const locales = getLocales();
        const deviceLang = locales[0].languageCode;
        this.language = deviceLang === 'zh' ? 'zh' : 'en';
    }
    
    async request(endpoint, options = {}) {
        const token = await AsyncStorage.getItem('access_token');
        
        const headers = {
            'Content-Type': 'application/json',
            'Accept-Language': this.language === 'zh' ? 
                'zh-CN,zh;q=0.9,en;q=0.8' : 'en-US,en;q=0.9',
            ...options.headers
        };
        
        if (token) {
            headers.Authorization = `Bearer ${token}`;
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
    
    async signup(email, password, userData = {}) {
        const result = await this.request('/auth/signup', {
            method: 'POST',
            body: JSON.stringify({
                email,
                password,
                data: {
                    user_language: this.language,
                    country: this.language === 'zh' ? 'CN' : 'US',
                    ...userData
                }
            })
        });
        
        await AsyncStorage.setItem('access_token', result.access_token);
        await AsyncStorage.setItem('user_language', this.language);
        
        return result;
    }
}

// 在React Native组件中使用
import React, { useState, useEffect } from 'react';
import { View, Text, Button, Alert } from 'react-native';

const DeviceScreen = () => {
    const [devices, setDevices] = useState([]);
    const [loading, setLoading] = useState(true);
    const api = new MobileAPI();
    
    useEffect(() => {
        loadDevices();
    }, []);
    
    const loadDevices = async () => {
        try {
            const result = await api.request('/api/v1/auth/devices');
            setDevices(result.data.devices);
            setLoading(false);
        } catch (error) {
            Alert.alert('错误', error.message);
            setLoading(false);
        }
    };
    
    const createDevice = async () => {
        try {
            const result = await api.request('/api/v1/auth/devices', {
                method: 'POST',
                body: JSON.stringify({
                    name: '新设备',
                    type: 'sensor',
                    location: '实验室'
                })
            });
            
            Alert.alert('成功', result.message);
            loadDevices(); // 重新加载列表
        } catch (error) {
            Alert.alert('错误', error.message);
        }
    };
    
    return (
        <View>
            <Button title="添加设备" onPress={createDevice} />
            {devices.map(device => (
                <View key={device.id}>
                    <Text>{device.name}</Text>
                    <Text>状态: {device.status}</Text>
                </View>
            ))}
        </View>
    );
};
```

## 🔧 开发工具

### Postman集合

```json
{
    "info": {
        "name": "YuzhaLink API - 多语言",
        "description": "YuzhaLink多语言API测试集合"
    },
    "variable": [
        {
            "key": "baseUrl",
            "value": "https://api.yuzhalink.com"
        },
        {
            "key": "token",
            "value": "{{access_token}}"
        }
    ],
    "item": [
        {
            "name": "认证",
            "item": [
                {
                    "name": "注册用户 (中文)",
                    "request": {
                        "method": "POST",
                        "header": [
                            {
                                "key": "Accept-Language",
                                "value": "zh-CN,zh;q=0.9,en;q=0.8"
                            },
                            {
                                "key": "Content-Type",
                                "value": "application/json"
                            }
                        ],
                        "url": "{{baseUrl}}/auth/signup",
                        "body": {
                            "mode": "raw",
                            "raw": "{\n  \"email\": \"test@example.com\",\n  \"password\": \"StrongPass123!\",\n  \"data\": {\n    \"name\": \"测试用户\",\n    \"user_language\": \"zh\",\n    \"country\": \"CN\"\n  }\n}"
                        }
                    }
                },
                {
                    "name": "用户登录",
                    "request": {
                        "method": "POST",
                        "header": [
                            {
                                "key": "Content-Type",
                                "value": "application/json"
                            }
                        ],
                        "url": "{{baseUrl}}/auth/login",
                        "body": {
                            "mode": "raw",
                            "raw": "{\n  \"email\": \"test@example.com\",\n  \"password\": \"StrongPass123!\"\n}"
                        }
                    }
                }
            ]
        },
        {
            "name": "设备管理",
            "item": [
                {
                    "name": "获取设备列表 (自动语言)",
                    "request": {
                        "method": "GET",
                        "header": [
                            {
                                "key": "Authorization",
                                "value": "Bearer {{token}}"
                            }
                        ],
                        "url": "{{baseUrl}}/api/v1/auth/devices"
                    }
                },
                {
                    "name": "获取设备列表 (强制英文)",
                    "request": {
                        "method": "GET",
                        "header": [
                            {
                                "key": "Authorization",
                                "value": "Bearer {{token}}"
                            }
                        ],
                        "url": {
                            "raw": "{{baseUrl}}/api/v1/auth/devices?lang=en",
                            "host": ["{{baseUrl}}"],
                            "path": ["api", "v1", "auth", "devices"],
                            "query": [
                                {
                                    "key": "lang",
                                    "value": "en"
                                }
                            ]
                        }
                    }
                }
            ]
        }
    ]
}
```

### 测试脚本

```bash
#!/bin/bash
# test-multilingual-api.sh

BASE_URL="https://api.yuzhalink.com"
EMAIL="test@example.com"
PASSWORD="StrongPass123!"

echo "🧪 测试多语言API..."

# 1. 注册中文用户
echo "📝 注册中文用户..."
SIGNUP_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/signup" \
  -H "Accept-Language: zh-CN,zh;q=0.9,en;q=0.8" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$EMAIL\",
    \"password\": \"$PASSWORD\",
    \"data\": {
      \"name\": \"测试用户\",
      \"user_language\": \"zh\",
      \"country\": \"CN\"
    }
  }")

echo "注册响应: $SIGNUP_RESPONSE"

# 提取token
TOKEN=$(echo $SIGNUP_RESPONSE | jq -r '.access_token')

if [ "$TOKEN" != "null" ]; then
    echo "✅ 注册成功，Token: ${TOKEN:0:20}..."
    
    # 2. 测试自动语言检测
    echo "🌐 测试自动语言检测..."
    DEVICES_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/auth/devices" \
      -H "Authorization: Bearer $TOKEN")
    
    MESSAGE=$(echo $DEVICES_RESPONSE | jq -r '.message')
    echo "设备列表响应消息: $MESSAGE"
    
    # 3. 测试强制语言切换
    echo "🔄 测试强制语言切换..."
    ENGLISH_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/auth/devices?lang=en" \
      -H "Authorization: Bearer $TOKEN")
    
    EN_MESSAGE=$(echo $ENGLISH_RESPONSE | jq -r '.message')
    echo "英文响应消息: $EN_MESSAGE"
    
    # 4. 测试创建设备 (验证错误)
    echo "❌ 测试验证错误..."
    ERROR_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/devices" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d "{
        \"type\": \"sensor\"
      }")
    
    ERROR_MESSAGE=$(echo $ERROR_RESPONSE | jq -r '.message')
    echo "验证错误消息: $ERROR_MESSAGE"
    
else
    echo "❌ 注册失败"
fi

echo "🏁 测试完成"
```

## 📚 更多资源

- [完整API文档](./multilingual-api-guide.md)
- [错误处理指南](./error-handling-guide.md)
- [SDK下载](./sdks/)
- [示例项目](./examples/)

## 💡 小贴士

1. **性能优化**: 使用JWT中的语言偏好可以减少语言检测开销
2. **用户体验**: 注册时自动检测用户语言，减少手动设置
3. **错误处理**: 始终显示API返回的本地化错误消息
4. **缓存策略**: 缓存支持的语言和国家列表，减少API调用 