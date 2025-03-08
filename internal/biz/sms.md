# SMS 模块 DDD 设计文档

## 领域模型

SMS 模块采用 DDD（领域驱动设计）架构进行设计，清晰地分离了不同的职责层和关注点。

### 核心实体

- **SmsInfo**：短信实体，包含以下属性：
  - ID：短信唯一标识
  - From：发送方号码
  - To：接收方号码
  - Message：短信内容
  - CreatedAt：创建时间

## 架构层次

SMS 模块按照典型的 DDD 分层架构组织：

### 1. 领域层 (Domain Layer)

位于 `internal/biz` 目录，定义了核心业务逻辑和规则。

核心组件：
- **SmsInfo**：短信领域实体
- **SmsBiz**：短信业务领域接口，定义发送和查询短信的能力
- **SmsRepo**：短信存储仓储接口，定义保存和查询短信记录的能力

### 2. 基础设施层 (Infrastructure Layer)

位于 `internal/data` 目录，负责实现领域层定义的仓储接口。

核心组件：
- **Sms**：实现 `SmsRepo` 接口，提供对短信数据的访问
- **SmsCache**：提供短信数据的内存缓存能力，包括自动清理过期记录

### 3. 应用层 (Application Layer)

位于 `internal/service` 目录，负责协调领域层与外部接口的交互。

核心组件：
- **smsService**：实现 API 接口服务，处理请求并调用领域层业务逻辑

### 4. 接口层 (Interface Layer)

位于 `internal/server` 和 `api/sms` 目录，负责对外提供服务接口。

核心组件：
- **SmsService**：定义 HTTP/gRPC 服务接口
- **WebServer**：配置并启动 HTTP 服务器

## 依赖引用关系

```
                       depends on
┌─────────────────────┐  <───  ┌──────────────────┐
│     Interface       │        │    Application   │
│  (server, api/sms)  │        │   (service)      │
└─────────────────────┘        └──────────────────┘
                                       │
                              depends  │
                                 on    ▼
┌─────────────────────┐  <───  ┌──────────────────┐
│   Infrastructure    │depends │     Domain       │
│      (data)         │   on   │     (biz)        │
└─────────────────────┘        └──────────────────┘
```

### 详细引用关系

1. **领域层** (`internal/biz`) 定义了核心接口和实体：
   - `SmsInfo`：核心领域实体
   - `SmsBiz`：领域服务接口
   - `SmsRepo`：仓储接口
   
   **不依赖**其他层组件

2. **基础设施层** (`internal/data`) 实现了领域层定义的仓储接口：
   - `Sms` 实现了 `biz.SmsRepo` 接口
   - `SmsCache` 提供缓存功能
   
   **依赖**：
   - `biz.SmsRepo`：实现此接口
   - `biz.SmsInfo`：使用此实体

3. **应用层** (`internal/service`) 协调领域业务逻辑：
   - `smsService` 实现了 `sms.SmsServiceHTTPServer` 接口
   
   **依赖**：
   - `biz.SmsBiz`：调用业务逻辑
   - `biz.SmsInfo`：使用业务实体
   - `api/sms`：实现对外接口

4. **接口层** 提供对外服务：
   - `api/sms`：定义 Protobuf 服务接口
   - `internal/server`：配置和启动 Web 服务器
   
   **依赖**：
   - `service`：使用服务实现

5. **应用组装** (`internal/app/sms-bridge`)：
   - 使用依赖注入组装各层组件，由 Google Wire 自动生成
   
   依赖关系：
   ```
   WebServer ← SmsService ← SmsBiz ← SmsRepo ← SmsCache
   ```

## 数据流转过程

在短信发送过程中的数据流转：

1. 接口层接收 HTTP 请求 (`api/sms` 和 `internal/server`)
2. 应用层处理请求并创建领域实体 (`internal/service`)
3. 领域层执行核心业务逻辑 (`internal/biz`)
4. 基础设施层保存数据 (`internal/data`)

## 优势分析

1. **关注点分离**：各层职责明确，便于维护和扩展
2. **依赖倒置**：通过接口定义，高层不依赖低层实现细节
3. **可测试性**：接口化设计便于单元测试和模拟
4. **可扩展性**：可以轻松替换具体实现，如将内存缓存替换为数据库存储
