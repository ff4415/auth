# 多语言错误消息支持

这个包实现了多语言错误消息的支持，默认支持英文和中文，同时提供了隐藏内部错误详情的功能。

## 功能特性

1. **多语言支持**: 默认支持英文(en)和中文(zh)
2. **内部错误隐藏**: 自动隐藏数据库错误、SQL错误等内部技术细节
3. **灵活的语言检测**: 支持多种方式指定用户语言偏好
4. **优雅降级**: 未支持的语言自动降级到英文

## 语言检测优先级

系统按以下优先级检测用户的语言偏好：

1. **查询参数**: `?lang=zh` 或 `?lang=en`
2. **自定义请求头**: `X-Language: zh-CN`
3. **标准请求头**: `Accept-Language: zh-CN,zh;q=0.9,en;q=0.8`
4. **默认语言**: 英文 (en)

## 使用示例

### 1. 基本使用

```go
import "github.com/supabase/auth/internal/i18n"

// 从请求中获取用户语言偏好
userLang := i18n.GetLanguageFromRequest(r)

// 获取特定错误消息的本地化版本
message := i18n.GetMessage(userLang, "weak_password")
// 英文: "Password does not meet security requirements"
// 中文: "密码不符合安全要求"
```

### 2. 用户友好的错误消息

```go
// 将内部错误转换为用户友好的消息
userFriendlyMsg := i18n.GetUserFriendlyMessage(userLang, errorCode, originalError)

// 例如：数据库错误 "pq: duplicate key violates constraint" 
// 会被转换为 "内部服务器错误" (中文) 或 "Internal server error" (英文)
```

### 3. 客户端指定语言

客户端可以通过以下方式指定语言：

```bash
# 使用查询参数
curl "https://api.example.com/auth/signup?lang=zh"

# 使用自定义请求头
curl -H "X-Language: zh-CN" "https://api.example.com/auth/signup"

# 使用标准请求头
curl -H "Accept-Language: zh-CN,zh;q=0.9" "https://api.example.com/auth/signup"
```

## 支持的错误消息

### 通用错误
- `weak_password`: 弱密码错误
- `unexpected_failure`: 意外错误
- `unknown_error`: 未知错误
- `validation_failed`: 验证失败
- `bad_json`: JSON解析错误
- `internal_server_error`: 内部服务器错误

### HTTP状态错误
- `unauthorized`: 未授权 (401)
- `forbidden`: 禁止访问 (403)
- `not_found`: 未找到 (404)
- `too_many_requests`: 请求过于频繁 (429)
- `conflict`: 冲突 (409)
- `unprocessable_entity`: 无法处理的实体 (422)
- `bad_request`: 错误的请求 (400)

### 业务错误
- `duplicate_email`: 邮箱重复
- `duplicate_phone`: 手机号重复
- `captcha_failed`: 验证码失败
- `email_not_confirmed`: 邮箱未确认
- `phone_not_confirmed`: 手机号未确认
- `mfa_verification_failed`: 多因子认证失败
- `insufficient_aal`: 认证级别不足

## 错误隐藏策略

系统会自动隐藏以下类型的内部错误：

1. **数据库错误**: 包含 "sql", "database", "pq:" 等关键词的错误
2. **内部服务器错误**: 包含 "internal", "server error" 等关键词的错误
3. **技术实现细节**: 任何可能暴露系统架构的错误信息

这些错误在日志中会完整记录，但返回给用户的是友好的通用错误消息。

## 扩展新语言

要添加新语言支持，需要：

1. 在 `Language` 类型中添加新的语言常量
2. 在 `Messages` 映射中添加新语言的翻译
3. 在 `normalizeLanguage` 函数中添加语言检测逻辑

```go
const (
    LanguageEnglish Language = "en"
    LanguageChinese Language = "zh"
    LanguageFrench  Language = "fr"  // 新增法语
)

var Messages = map[Language]map[string]string{
    // ... 现有语言
    LanguageFrench: {
        "weak_password": "Le mot de passe ne répond pas aux exigences de sécurité",
        // ... 其他翻译
    },
}
```

## 注意事项

1. **日志记录**: 所有原始错误都会在服务器日志中完整记录，便于调试
2. **性能优化**: 消息映射在内存中缓存，查找效率很高
3. **向后兼容**: 如果翻译不存在，会自动降级到英文版本
4. **安全性**: 永远不会向用户暴露内部技术细节 