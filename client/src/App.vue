<template>
  <div class="app">
    <!-- 登录页 -->
    <div v-if="!loggedIn" class="auth-page">
      <div class="auth-card">
        <div class="auth-logo">
          <h1>聊天室</h1>
          <p class="subtitle">登录账号，继续聊天</p>
        </div>

        <div v-if="authView === 'login'">
          <div class="field">
            <input v-model="loginForm.username" placeholder="用户名" @keyup.enter="doLogin" />
          </div>
          <div class="field">
            <input v-model="loginForm.password" type="password" placeholder="密码" @keyup.enter="doLogin" />
          </div>
          <p v-if="loginError" class="error">{{ loginError }}</p>
          <button class="btn-primary" @click="doLogin">登录</button>
          <p class="switch">没有账号？<a @click="authView='register'">注册</a></p>
        </div>

        <div v-else>
          <div class="field">
            <input v-model="regForm.username" placeholder="用户名" @keyup.enter="doRegister" />
          </div>
          <div class="field">
            <input v-model="regForm.password" type="password" placeholder="密码" @keyup.enter="doRegister" />
          </div>
          <div class="field">
            <input v-model="regForm.nickname" placeholder="昵称" @keyup.enter="doRegister" />
          </div>
          <p v-if="regError" class="error">{{ regError }}</p>
          <button class="btn-primary" @click="doRegister">注册</button>
          <p class="switch">已有账号？<a @click="authView='login'">登录</a></p>
        </div>
      </div>
    </div>

    <!-- 主界面 -->
    <div v-else class="layout">
      <div v-if="sidebarOpen" class="sidebar-overlay" @click="sidebarOpen = false"></div>

      <aside class="sidebar" :class="{ open: sidebarOpen }">
        <div class="sidebar-top">
          <div class="sidebar-brand">
            <div class="brand-icon" aria-hidden="true">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M21 15a4 4 0 0 1-4 4H8l-5 3V7a4 4 0 0 1 4-4h10a4 4 0 0 1 4 4z" />
                <path d="M8 9h8M8 13h5" />
              </svg>
            </div>
            <span class="brand-text">聊天室</span>
          </div>
          <nav class="sidebar-nav">
            <a class="nav-item" :class="{ active: view === 'chat' }" @click="switchView('chat')">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>
              聊天
            </a>
            <a class="nav-item" :class="{ active: view === 'profile' }" @click="switchView('profile')">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
              个人资料
            </a>
            <a class="nav-item" :class="{ active: view === 'about' }" @click="switchView('about')">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/></svg>
              关于
            </a>
          </nav>
        </div>

        <div class="sidebar-bottom">
          <div class="sidebar-user" @click="switchView('profile')">
            <div class="user-avatar">
              <img v-if="avatar" :src="avatar" :alt="nickname" />
              <span v-else>{{ nickname.charAt(0) }}</span>
            </div>
            <div class="user-meta">
              <span class="user-name">{{ nickname }}</span>
              <span class="user-status">在线</span>
            </div>
          </div>
          <button class="btn-logout" @click="logout" title="退出登录">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>
          </button>
        </div>
      </aside>

      <main class="main-content">
        <!-- 聊天页 -->
        <div v-if="view === 'chat'" class="chat-view">
          <div class="chat-header">
            <button class="menu-btn" @click="sidebarOpen = true">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="3" y1="6" x2="21" y2="6"/><line x1="3" y1="12" x2="21" y2="12"/><line x1="3" y1="18" x2="21" y2="18"/></svg>
            </button>
            <div class="chat-header-info">
              <h2>{{ currentRoom }}</h2>
              <span class="room-status">
                <span class="status-dot"></span>
                {{ users.length }} 人在线
              </span>
            </div>
            <button class="btn-leave" @click="disconnect" title="退出房间">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>
            </button>
          </div>

          <div class="messages" ref="msgBox">
            <div v-for="(m, i) in messages" :key="i" class="msg-row" :class="{ self: m.nick === nickname }">
              <div class="msg-avatar">
                <img v-if="m.avatar" :src="m.avatar" :alt="m.nick" />
                <span v-else>{{ m.nick.charAt(0) }}</span>
              </div>
              <div class="msg-content">
                <span class="msg-nick">{{ m.nick }}</span>
                <img v-if="m.mediaType === 'image' && m.mediaUrl" :src="m.mediaUrl" class="msg-image" @load="scrollToBottom" />
                <video v-else-if="m.mediaType === 'video' && m.mediaUrl" :src="m.mediaUrl" class="msg-video" controls @loadedmetadata="scrollToBottom"></video>
                <div class="msg-body" :class="{ 'no-border': m.mediaUrl && m.mediaType !== 'file' }">
                  <a v-if="m.mediaType === 'file' && m.mediaUrl" :href="m.mediaUrl" class="msg-file" target="_blank" download>
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/></svg>
                    <span>{{ m.content || '查看文件' }}</span>
                  </a>
                  <p v-if="m.content">{{ m.content }}</p>
                  <span class="msg-time">{{ formatTime(m.time) }}</span>
                </div>
              </div>
            </div>
            <div v-if="messages.length === 0" class="empty-msg">
              <p>房间里还没有消息</p>
            </div>
          </div>

          <div class="input-bar">
            <label class="btn-attach" title="发送文件">
              <input type="file" ref="fileInput" @change="uploadMedia" />
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21.44 11.05l-9.19 9.19a6 6 0 0 1-8.49-8.49l9.19-9.19a4 4 0 0 1 5.66 5.66l-9.2 9.19a2 2 0 0 1-2.83-2.83l8.49-8.48"/></svg>
            </label>
            <input v-model="text" placeholder="输入消息..." @keyup.enter="send" />
            <button class="btn-send" @click="send">发送</button>
          </div>
        </div>

        <!-- 个人资料页 -->
        <div v-else-if="view === 'profile'" class="profile-view">
          <div class="profile-header">
            <button class="menu-btn" @click="sidebarOpen = true">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="3" y1="6" x2="21" y2="6"/><line x1="3" y1="12" x2="21" y2="12"/><line x1="3" y1="18" x2="21" y2="18"/></svg>
            </button>
            <h2>个人资料</h2>
          </div>
          <div class="profile-content">
            <div class="profile-avatar-section">
              <div class="profile-avatar">
                <img v-if="avatar" :src="avatar" :alt="nickname" />
                <span v-else>{{ nickname.charAt(0) }}</span>
              </div>
              <label class="avatar-upload">
                <input type="file" accept="image/*" @change="uploadAvatar" />
                更换头像
              </label>
            </div>
            <div class="profile-form">
              <div class="field">
                <label>用户名</label>
                <input :value="username" disabled />
              </div>
              <div class="field">
                <label>邮箱</label>
                <input v-model="profileForm.email" placeholder="输入邮箱" />
              </div>
              <div class="field">
                <label>昵称</label>
                <input v-model="profileForm.nickname" placeholder="输入昵称" />
              </div>
              <div class="field">
                <label>简介</label>
                <textarea v-model="profileForm.bio" placeholder="写一段个人简介" rows="3"></textarea>
              </div>
              <p v-if="profileMsg" class="profile-msg" :class="{ error: profileError }">{{ profileMsg }}</p>
              <button class="btn-primary" @click="saveProfile">保存</button>
            </div>
          </div>
        </div>

        <!-- 关于页 -->
        <div v-else class="about-view">
          <div class="about-header">
            <button class="menu-btn" @click="sidebarOpen = true">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="3" y1="6" x2="21" y2="6"/><line x1="3" y1="12" x2="21" y2="12"/><line x1="3" y1="18" x2="21" y2="18"/></svg>
            </button>
            <h2>关于</h2>
          </div>
          <div class="about-content">
            <div class="about-card">
              <h3>聊天室</h3>
              <p class="about-desc">一个简单的多人在线聊天工具。</p>
              <div class="about-row"><span>版本</span><strong>1.0.0</strong></div>
              <div class="about-row"><span>数据存储</span><strong>本地 JSON</strong></div>
            </div>
            <div class="about-card">
              <h4>说明</h4>
              <p class="about-desc">聊天记录保存在服务器本地。退出登录会清除当前设备上的登录状态，但不会删除历史消息。</p>
            </div>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>

<script>
const STORAGE_KEY = 'chat_user'

export default {
  data() {
    return {
      loggedIn: false,
      authView: 'login',
      username: '',
      nickname: '',
      email: '',
      bio: '',
      avatar: '',
      loginForm: { username: '', password: '' },
      regForm: { username: '', password: '', nickname: '' },
      loginError: '',
      regError: '',
      profileForm: { email: '', nickname: '', bio: '' },
      profileMsg: '',
      profileError: false,
      view: 'chat',
      sidebarOpen: false,
      room: '',
      currentRoom: '',
      text: '',
      connected: false,
      messages: [],
      users: [],
      ws: null,
    }
  },
  mounted() {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (saved) {
      try {
        const data = JSON.parse(saved)
        this.loginForm.username = data.username
        this.loginForm.password = data.password
        this.nickname = data.nickname
        this.username = data.username
        this.email = data.email || ''
        this.bio = data.bio || ''
        this.avatar = data.avatar || ''
        this.autoLogin(data.username, data.password)
      } catch {
        localStorage.removeItem(STORAGE_KEY)
      }
    }
  },
  methods: {
    async autoLogin(username, password) {
      try {
        const res = await fetch('/api/login', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ username, password }),
        })
        const data = await res.json()
        if (data.success) {
          this.applyUserData(data)
          this.loggedIn = true
          this.connect()
        } else {
          localStorage.removeItem(STORAGE_KEY)
        }
      } catch {
        localStorage.removeItem(STORAGE_KEY)
      }
    },
    async doLogin() {
      this.loginError = ''
      try {
        const res = await fetch('/api/login', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(this.loginForm),
        })
        const data = await res.json()
        if (data.success) {
          this.applyUserData(data)
          this.saveToLocal()
          this.loggedIn = true
          this.connect()
        } else {
          this.loginError = data.message
        }
      } catch {
        this.loginError = '网络错误，请重试'
      }
    },
    async doRegister() {
      this.regError = ''
      try {
        const res = await fetch('/api/register', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(this.regForm),
        })
        const data = await res.json()
        if (data.success) {
          this.loginForm.username = this.regForm.username
          this.loginForm.password = this.regForm.password
          this.authView = 'login'
          this.loginError = '注册成功，请登录'
        } else {
          this.regError = data.message
        }
      } catch {
        this.regError = '网络错误，请重试'
      }
    },
    applyUserData(data) {
      this.username = data.username
      this.nickname = data.nickname
      this.email = data.email || ''
      this.bio = data.bio || ''
      this.avatar = data.avatar || ''
      this.profileForm.email = data.email || ''
      this.profileForm.nickname = data.nickname
      this.profileForm.bio = data.bio || ''
    },
    saveToLocal() {
      localStorage.setItem(STORAGE_KEY, JSON.stringify({
        username: this.username,
        password: this.loginForm.password,
        nickname: this.nickname,
        email: this.email,
        bio: this.bio,
        avatar: this.avatar,
      }))
    },
    async saveProfile() {
      this.profileMsg = ''
      this.profileError = false
      try {
        const res = await fetch('/api/profile/update', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            username: this.username,
            email: this.profileForm.email,
            nickname: this.profileForm.nickname,
            bio: this.profileForm.bio,
          }),
        })
        const data = await res.json()
        if (data.success) {
          this.nickname = data.nickname
          this.email = data.email
          this.bio = data.bio
          this.profileForm.nickname = data.nickname
          this.profileForm.email = data.email
          this.profileForm.bio = data.bio
          this.saveToLocal()
          this.profileMsg = '保存成功'
        } else {
          this.profileMsg = data.message
          this.profileError = true
        }
      } catch {
        this.profileMsg = '网络错误，请重试'
        this.profileError = true
      }
    },
    async uploadAvatar(e) {
      const file = e.target.files[0]
      if (!file) return

      const form = new FormData()
      form.append('username', this.username)
      form.append('avatar', file)

      try {
        const res = await fetch('/api/avatar/upload', {
          method: 'POST',
          body: form,
        })
        const data = await res.json()
        if (data.success) {
          this.avatar = data.avatar
          this.saveToLocal()
        }
      } catch {
        // ignore
      }
    },
    async uploadMedia(e) {
      const file = e.target.files[0]
      if (!file) return
      const form = new FormData()
      form.append('file', file)
      try {
        const res = await fetch('/api/media/upload', { method: 'POST', body: form })
        const data = await res.json()
        if (data.success) {
          this.ws.send(JSON.stringify({
            type: 'message',
            content: data.mediaType === 'image' ? '' : file.name,
            room: this.currentRoom,
            mediaUrl: data.mediaUrl,
            mediaType: data.mediaType,
          }))
        }
      } catch {}
      this.$refs.fileInput.value = ''
    },
    scrollToBottom() {
      this.$nextTick(() => {
        const box = this.$refs.msgBox
        if (box) box.scrollTop = box.scrollHeight
      })
    },
    switchView(v) {
      this.view = v
      this.sidebarOpen = false
    },
    connect() {
      const room = this.room.trim() || '大厅'
      const wsUrl = `${location.protocol === 'https:' ? 'wss:' : 'ws:'}//${location.host}/ws?username=${encodeURIComponent(this.username)}&nick=${encodeURIComponent(this.nickname)}&room=${encodeURIComponent(room)}`
      this.ws = new WebSocket(wsUrl)

      this.ws.onopen = () => {
        this.connected = true
        this.currentRoom = room
      }

      this.ws.onmessage = (e) => {
        const msg = JSON.parse(e.data)
        switch (msg.type) {
          case 'message':
            this.messages.push(msg)
            this.$nextTick(() => {
              const box = this.$refs.msgBox
              if (box) box.scrollTop = box.scrollHeight
            })
            break
          case 'users':
            this.users = msg.users
            break
        }
      }

      this.ws.onclose = () => {
        this.connected = false
        this.messages = []
        this.users = []
      }
    },
    send() {
      if (!this.text.trim()) return
      this.ws.send(JSON.stringify({
        type: 'message',
        content: this.text,
        room: this.currentRoom,
      }))
      this.text = ''
    },
    disconnect() {
      this.ws.close()
      this.loggedIn = false
    },
    logout() {
      if (this.ws) this.ws.close()
      localStorage.removeItem(STORAGE_KEY)
      this.connected = false
      this.messages = []
      this.users = []
      this.loggedIn = false
    },
    formatTime(ts) {
      return new Date(ts).toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
    },
  },
}
</script>

<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
body {
  font-family: -apple-system, BlinkMacSystemFont, 'PingFang SC', 'Microsoft YaHei', sans-serif;
  background: #f4f4f2;
  color: #202020;
  height: 100vh;
}
.app { height: 100vh; }

/* Auth */
.auth-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  background: #f4f4f2;
}
.auth-card {
  width: 100%;
  max-width: 360px;
  padding: 24px;
}
.auth-logo { margin-bottom: 28px; }
.auth-logo h1 { font-size: 26px; font-weight: 650; }
.subtitle { color: #888; font-size: 14px; margin-top: 6px; }
.field { margin-bottom: 14px; }
.field label { display: block; font-size: 13px; color: #666; margin-bottom: 4px; }
.field input, .field textarea {
  width: 100%;
  padding: 14px 16px;
  border: 1px solid #cececa;
  border-radius: 7px;
  font-size: 14px;
  outline: none;
  transition: border-color 0.2s;
  background: #fff;
  font-family: inherit;
}
.field textarea { resize: vertical; }
.field input:focus, .field textarea:focus { border-color: #1a1a1a; background: #fff; }
.field input:disabled { background: #f0f0ee; color: #888; }
.error { color: #e53935; font-size: 13px; margin-bottom: 12px; text-align: center; }
.btn-primary {
  width: 100%;
  padding: 14px;
  background: #1a1a1a;
  color: #fff;
  border: none;
  border-radius: 7px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}
.btn-primary:hover { background: #333; }
.switch { text-align: center; margin-top: 18px; font-size: 13px; color: #888; }
.switch a { color: #1a1a1a; cursor: pointer; font-weight: 600; }

/* Layout */
.layout { display: flex; height: 100vh; }

/* Sidebar */
.sidebar-overlay {
  display: none;
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.3);
  z-index: 99;
}
.sidebar {
  width: 236px;
  background: #ececea;
  border-right: 1px solid #d9d9d5;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  flex-shrink: 0;
}
.sidebar-brand {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 20px 20px 24px;
  border-bottom: 1px solid #d9d9d5;
}
.brand-icon {
  width: 28px;
  height: 28px;
  color: #292929;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.brand-text { font-size: 16px; font-weight: 650; }
.sidebar-nav { padding: 12px 10px; display: flex; flex-direction: column; gap: 4px; }
.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 14px;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  color: #666;
  cursor: pointer;
  transition: all 0.15s;
  text-decoration: none;
}
.nav-item:hover { background: #e2e2df; color: #1a1a1a; }
.nav-item.active { background: #dcdcd8; color: #1a1a1a; font-weight: 600; }
.sidebar-bottom {
  padding: 16px;
  border-top: 1px solid #d9d9d5;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.sidebar-user { display: flex; align-items: center; gap: 10px; cursor: pointer; }
.user-avatar {
  width: 36px;
  height: 36px;
  background: #1a1a1a;
  color: #fff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
  overflow: hidden;
  flex-shrink: 0;
}
.user-avatar img { width: 100%; height: 100%; object-fit: cover; }
.user-meta { display: flex; flex-direction: column; }
.user-name { font-size: 14px; font-weight: 600; }
.user-status { font-size: 11px; color: #4caf50; }
.btn-logout {
  width: 36px;
  height: 36px;
  border: none;
  background: transparent;
  border-radius: 6px;
  cursor: pointer;
  color: #999;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s;
}
.btn-logout:hover { background: #f5f5f5; color: #e53935; }

/* Main */
.main-content { flex: 1; display: flex; flex-direction: column; min-width: 0; }

/* Chat */
.chat-view { display: flex; flex-direction: column; height: 100vh; }
.chat-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 20px;
  background: #f9f9f7;
  border-bottom: 1px solid #deded9;
  flex-shrink: 0;
}
.menu-btn {
  display: none;
  width: 36px;
  height: 36px;
  border: none;
  background: transparent;
  border-radius: 6px;
  cursor: pointer;
  color: #666;
  align-items: center;
  justify-content: center;
}
.menu-btn:hover { background: #f5f5f5; }
.chat-header-info { flex: 1; }
.chat-header-info h2 { font-size: 16px; font-weight: 600; }
.room-status {
  font-size: 12px;
  color: #888;
  display: flex;
  align-items: center;
  gap: 6px;
}
.status-dot {
  width: 6px;
  height: 6px;
  background: #4caf50;
  border-radius: 50%;
}
.btn-leave {
  width: 36px;
  height: 36px;
  border: none;
  background: transparent;
  border-radius: 6px;
  cursor: pointer;
  color: #999;
  display: flex;
  align-items: center;
  justify-content: center;
}
.btn-leave:hover { background: #f5f5f5; color: #e53935; }

/* Messages */
.messages {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 6px;
  background: #f9f9f7;
}
.empty-msg {
  margin: auto;
  text-align: center;
  color: #bbb;
}
.empty-msg p { font-size: 14px; }
.msg-row {
  display: flex;
  gap: 8px;
  max-width: min(76%, 680px);
  align-self: flex-start;
}
.msg-row.self {
  align-self: flex-end;
  flex-direction: row-reverse;
}
.msg-avatar {
  width: 32px;
  height: 32px;
  background: #d4d4d0;
  color: #555;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  flex-shrink: 0;
  overflow: hidden;
  margin-top: 2px;
}
.msg-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.msg-row.self .msg-avatar {
  background: #1a1a1a;
  color: #fff;
}
.msg-content {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.msg-row.self .msg-content {
  align-items: flex-end;
}
.msg-nick {
  font-size: 11px;
  color: #999;
  padding: 0 4px;
}
.msg-row.self .msg-nick {
  display: none;
}
.msg-body {
  background: #fff;
  padding: 8px 12px;
  border: 1px solid #e2e2de;
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.msg-body.no-border {
  border: none;
  background: transparent;
  padding: 4px 0;
}
.msg-row.self .msg-body {
  background: #1a1a1a;
  border-color: #1a1a1a;
}
.msg-row.self .msg-body.no-border {
  background: transparent;
  border: none;
}
.msg-body p {
  font-size: 14px;
  line-height: 1.5;
  word-break: break-word;
  color: #202020;
}
.msg-row.self .msg-body p {
  color: #fff;
}
.msg-time {
  font-size: 10px;
  color: #bbb;
  text-align: right;
  line-height: 1;
}
.msg-row.self .msg-time {
  color: rgba(255,255,255,0.5);
}

/* Input */
.input-bar {
  display: flex;
  padding: 14px 20px;
  gap: 8px;
  background: #f9f9f7;
  border-top: 1px solid #deded9;
  align-items: center;
}
.btn-attach {
  width: 36px;
  height: 36px;
  border: 1px solid #cececa;
  border-radius: 8px;
  background: #fff;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #666;
  flex-shrink: 0;
  transition: all 0.15s;
}
.btn-attach:hover { border-color: #1a1a1a; color: #1a1a1a; }
.btn-attach input { display: none; }
.input-bar input {
  flex: 1;
  padding: 12px 16px;
  border: 1px solid #cececa;
  border-radius: 8px;
  font-size: 14px;
  outline: none;
  transition: border-color 0.2s;
  background: #fff;
}
.input-bar input:focus { border-color: #1a1a1a; background: #fff; }
.btn-send {
  padding: 10px 20px;
  background: #1a1a1a;
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
  flex-shrink: 0;
}
.btn-send:hover { background: #333; }

/* Message media */
.msg-image {
  width: 100%;
  max-height: 400px;
  object-fit: contain;
  border-radius: 8px;
  display: block;
  background: #f0f0ee;
  cursor: pointer;
}
.msg-video {
  max-width: 100%;
  max-height: 300px;
  border-radius: 8px;
  display: block;
}
.msg-file {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 0;
  color: #1a1a1a;
  text-decoration: none;
  font-size: 13px;
}
.msg-row.self .msg-file { color: #fff; }
.msg-file:hover { text-decoration: underline; }

/* Profile */
.profile-view { display: flex; flex-direction: column; height: 100vh; }
.profile-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 20px;
  background: #f9f9f7;
  border-bottom: 1px solid #deded9;
}
.profile-header h2 { font-size: 16px; font-weight: 600; }
.profile-content {
  flex: 1;
  overflow-y: auto;
  padding: 32px 20px;
  max-width: 480px;
  margin: 0 auto;
  width: 100%;
}
.profile-avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 32px;
}
.profile-avatar {
  width: 80px;
  height: 80px;
  background: #1a1a1a;
  color: #fff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  font-weight: 600;
  overflow: hidden;
  margin-bottom: 12px;
}
.profile-avatar img { width: 100%; height: 100%; object-fit: cover; }
.avatar-upload {
  font-size: 13px;
  color: #666;
  cursor: pointer;
  padding: 6px 16px;
  border: 1px solid #cececa;
  border-radius: 6px;
  transition: all 0.15s;
}
.avatar-upload:hover { border-color: #1a1a1a; color: #1a1a1a; }
.avatar-upload input { display: none; }
.profile-form { width: 100%; }
.profile-msg { font-size: 13px; margin-bottom: 12px; text-align: center; }
.profile-msg.error { color: #e53935; }

/* About */
.about-view { display: flex; flex-direction: column; height: 100vh; }
.about-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 20px;
  background: #f9f9f7;
  border-bottom: 1px solid #deded9;
}
.about-header h2 { font-size: 16px; font-weight: 600; }
.about-content {
  flex: 1;
  overflow-y: auto;
  padding: 24px 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  max-width: 680px;
  margin: 0 auto;
  width: 100%;
}
.about-card {
  padding: 4px 0 22px;
  border-bottom: 1px solid #deded9;
}
.about-card h3 { font-size: 18px; font-weight: 700; margin-bottom: 8px; }
.about-card h4 { font-size: 14px; font-weight: 600; margin-bottom: 12px; color: #666; }
.about-desc { font-size: 14px; color: #666; line-height: 1.6; }
.about-row { display: flex; justify-content: space-between; padding: 13px 0; border-top: 1px solid #e3e3df; font-size: 13px; }
.about-row:first-of-type { margin-top: 20px; }
.about-row span { color: #777; }
.about-row strong { font-weight: 500; }

/* Mobile */
@media (max-width: 768px) {
  .sidebar-overlay { display: block; }
  .sidebar {
    position: fixed;
    left: -280px;
    top: 0;
    bottom: 0;
    z-index: 100;
    transition: left 0.25s ease;
    box-shadow: 2px 0 8px rgba(0,0,0,0.12);
  }
  .sidebar.open { left: 0; }
  .menu-btn { display: flex; }
  .msg-row { max-width: 90%; }
  .messages { padding: 14px 12px; }
  .input-bar { padding: 10px 12px max(10px, env(safe-area-inset-bottom)); }
  .auth-card { padding: 24px 8px; }
}
</style>
