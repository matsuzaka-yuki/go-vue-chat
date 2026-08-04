# 聊天室 (Chat Room)

一个基于 **Go + WebSocket + Vue 3** 的多人在线聊天室，支持注册登录、实时聊天、私聊、关注/粉丝、图片与文件发送等完整功能。

## 技术栈

| 端 | 技术 |
| --- | --- |
| 后端 | Go 1.21、gorilla/websocket |
| 前端 | Vue 3、Vite、pnpm |
| 存储 | JSON 文件（`server/data/` 目录） |

## 功能特性

- **账号系统**：注册、登录、自动登录（localStorage 持久化）
- **个人资料**：头像上传、昵称、邮箱、简介，关注/粉丝数
- **社交系统**：关注/取关、粉丝列表、互关（好友）后才可私聊
- **实时聊天**：多房间大厅，消息实时推送，在线用户列表
- **私聊**：互关好友之间一对一聊天，仅双方可见
- **多媒体消息**：图片、视频、文件发送，图片灯箱预览
- **历史消息**：加载最近 50 条历史记录，持久化到本地
- **响应式**：移动端抽屉式侧栏，适配手机

## 目录结构

```
chat/
├── server/                 # Go 后端
│   ├── main.go             # HTTP 服务、API 路由
│   ├── hub.go              # 房间管理、消息广播、私聊路由
│   ├── client.go           # WebSocket 连接处理
│   ├── store.go            # JSON 持久化（用户/消息/关注关系）
│   ├── go.mod              # Go 依赖
│   └── data/               # 运行时数据（自动生成）
│       ├── users.json      # 用户信息
│       ├── messages.json   # 消息记录
│       ├── relations.json  # 关注关系
│       ├── avatars/        # 头像图片
│       └── media/          # 聊天媒体文件
└── client/                 # Vue 前端
    ├── vite.config.js      # Vite 配置（含代理）
    ├── package.json
    └── src/
        └── App.vue         # 主界面
```

## 快速开始

### 环境要求
- Go 1.21+
- Node.js 18+
- pnpm（或 npm）

### 1. 启动后端

```bash
cd server
go mod tidy      # 首次运行，下载依赖
go run .
```

后端默认监听 `http://localhost:8086`。

### 2. 启动前端

```bash
cd client
pnpm install     # 首次运行，安装依赖
pnpm dev
```

前端默认运行在 `http://localhost:3000`，已配置代理将 `/api`、`/ws`、`/avatars`、`/media` 转发到后端。

### 3. 打开使用

浏览器访问 `http://localhost:3000`。

> 注：如果使用 pnpm 遇到 PowerShell 执行策略报错，可改用 `npx pnpm`，或先执行 `Set-ExecutionPolicy -Scope CurrentUser RemoteSigned`。

## 使用说明

### 注册与登录
1. 打开页面，点击「注册」创建账号（用户名、密码、昵称）
2. 注册成功后自动切换回登录页，登录后进入聊天室
3. 关闭浏览器重开会自动保持登录状态

### 个人资料
- 点击左侧「个人资料」进入，可修改邮箱、昵称、简介
- 点击「更换头像」上传图片（支持 jpg/png/gif）
- 上方显示关注数和粉丝数，点击可查看列表

### 聊天大厅
- 进入后默认在「大厅」房间，右侧输入框发送消息
- 点击输入框左侧的 📎 图标可发送图片、视频或文件
- 点击消息旁的头像可查看该用户资料

### 关注好友
- 点击「联系人」页面，搜索用户名添加关注
- 对方关注你后即成为**互关好友**，出现在「我的好友」列表
- 点击好友头像可查看资料，点击「私聊」进入一对一聊天

### 私聊
- 仅**互关好友**之间可以私聊
- 私聊时右上角退出按钮可返回大厅
- 支持文字、图片、文件发送

## API 接口

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/api/register` | 注册 |
| POST | `/api/login` | 登录 |
| GET | `/api/profile` | 获取用户资料 |
| POST | `/api/profile/update` | 更新资料 |
| POST | `/api/avatar/upload` | 上传头像 |
| POST | `/api/media/upload` | 上传聊天媒体 |
| POST | `/api/follow` | 关注用户 |
| POST | `/api/unfollow` | 取消关注 |
| GET | `/api/relation` | 查询关注关系 |
| GET | `/api/following` | 关注列表 |
| GET | `/api/followers` | 粉丝列表 |
| GET | `/api/users/search` | 搜索用户 |
| WS | `/ws` | WebSocket 实时通信 |

## 数据存储

所有数据以 JSON 文件形式保存在 `server/data/` 目录，无需额外数据库。删除该目录即可重置数据（注意：会同时清除所有用户和消息）。

## 常见问题

**Q: 端口被占用怎么办？**
后端可用环境变量改端口：`PORT=9090 go run .`，前端改 `vite.config.js` 中的代理目标。

**Q: 中文文件名头像/文件 404？**
旧版本曾用中文文件名，现已改为哈希命名。如果历史数据残留旧路径，删除 `server/data/avatars/` 和 `server/data/media/` 下的旧文件并重新上传即可。