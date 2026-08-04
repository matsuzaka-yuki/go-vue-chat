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
            <a class="nav-item" :class="{ active: view === 'contacts' }" @click="switchView('contacts')">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
              联系人
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
              <img v-if="avatar" :src="assetUrl(avatar)" :alt="nickname" />
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
              <h2>{{ privateChatTo ? privateChatTo : currentRoom }}</h2>
              <span class="room-status" v-if="!privateChatTo">
                <span class="status-dot"></span>
                {{ users.length }} 人在线
              </span>
            </div>
            <button class="btn-leave" @click="handleExit" :title="privateChatTo ? '退出私聊' : '退出房间'">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>
            </button>
          </div>

          <div class="messages" ref="msgBox">
            <div v-for="(m, i) in messages" :key="i" class="msg-row" :class="{ self: m.nick === nickname }">
              <div class="msg-avatar" @click="showUserDetail(m)">
                <img v-if="m.avatar" :src="assetUrl(m.avatar)" :alt="m.nick" />
                <span v-else>{{ m.nick.charAt(0) }}</span>
              </div>
              <div class="msg-content">
                <span class="msg-nick">{{ m.nick }}</span>
                <img v-if="m.mediaType === 'image' && m.mediaUrl" :src="assetUrl(m.mediaUrl)" class="msg-image" @load="scrollToBottom" @click="openLightbox(assetUrl(m.mediaUrl))" />
                <video v-else-if="m.mediaType === 'video' && m.mediaUrl" :src="assetUrl(m.mediaUrl)" class="msg-video" controls @loadedmetadata="scrollToBottom"></video>
                <div class="msg-body" :class="{ 'no-border': m.mediaUrl && m.mediaType !== 'file' }">
                  <div v-if="m.mediaType === 'file' && m.mediaUrl" class="msg-file-card">
                    <div class="file-icon">{{ getFileExt(m.content) }}</div>
                    <div class="file-info" @click="openFile(m.mediaUrl)">
                      <strong>{{ m.content }}</strong>
                      <span>{{ m.fileSize ? formatFileSize(m.fileSize) : '点击下载' }}</span>
                    </div>
                    <a :href="assetUrl(m.mediaUrl)" target="_blank" download class="file-download">
                      <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
                    </a>
                  </div>
                  <p v-if="m.content && m.mediaType !== 'file'">{{ m.content }}</p>
                  <span class="msg-time">{{ formatTime(m.time) }}</span>
                </div>
              </div>
            </div>
            <div v-if="messages.length === 0" class="empty-msg">
              <p>房间里还没有消息</p>
            </div>
            <div v-if="uploading" class="uploading-indicator">
              <span class="spinner"></span>
              <span>正在上传...</span>
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
                <img v-if="avatar" :src="assetUrl(avatar)" :alt="nickname" />
                <span v-else>{{ nickname.charAt(0) }}</span>
              </div>
              <label class="avatar-upload">
                <input type="file" accept="image/*" @change="uploadAvatar" />
                更换头像
              </label>
            </div>
            <div class="profile-info">
              <h3>{{ nickname }}</h3>
              <p class="profile-username">@{{ username }}</p>
            </div>
            <div class="profile-stats">
              <div class="stat" @click="showRelationList('following', { username, nickname })"><strong>{{ profileFollowing }}</strong><span>关注</span></div>
              <div class="stat" @click="showRelationList('followers', { username, nickname })"><strong>{{ profileFollowers }}</strong><span>粉丝</span></div>
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

        <!-- 联系人页 -->
        <div v-else-if="view === 'contacts'" class="contacts-view">
          <div class="contacts-header">
            <button class="menu-btn" @click="sidebarOpen = true">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="3" y1="6" x2="21" y2="6"/><line x1="3" y1="12" x2="21" y2="12"/><line x1="3" y1="18" x2="21" y2="18"/></svg>
            </button>
            <h2>联系人</h2>
          </div>
          <div class="contacts-body">
            <div class="contacts-search">
              <svg class="search-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
              <input v-model="searchQuery" placeholder="搜索用户添加好友..." @input="searchUsers" />
              <div v-if="searchResults.length > 0" class="search-results">
                <div class="search-title">搜索结果</div>
                <div v-for="u in searchResults" :key="u.username" class="search-item">
                  <div class="search-avatar" @click="showContactDetail(u)">
                    <img v-if="u.avatar" :src="assetUrl(u.avatar)" />
                    <span v-else>{{ u.nickname.charAt(0) }}</span>
                  </div>
                  <div class="search-info" @click="showContactDetail(u)">
                    <strong>{{ u.nickname }}</strong>
                    <span>{{ u.username }}</span>
                  </div>
                  <button class="btn-add" v-if="u.username !== username && !isContact(u.username)" @click="doFollow(u)">关注</button>
                  <button class="btn-add followed" v-if="u.username !== username && isContact(u.username)" @click="doUnfollow(u)">已关注</button>
                </div>
                <div v-if="searchResults.length === 0 && searchQuery" class="search-title">未找到相关用户</div>
              </div>
            </div>
            <div class="contacts-section">
              <div class="section-header">
                <span>我的好友</span>
                <span class="section-count">{{ contacts.length }}</span>
              </div>
              <div class="contacts-list">
                <div v-for="c in contacts" :key="c.username" class="contact-item">
                  <div class="contact-avatar" @click="showContactDetail(c)">
                    <img v-if="c.avatar" :src="assetUrl(c.avatar)" />
                    <span v-else>{{ c.nickname.charAt(0) }}</span>
                  </div>
                  <div class="contact-info" @click="startPrivateChat(c)">
                    <strong>{{ c.nickname }}</strong>
                    <span>{{ c.username }}</span>
                  </div>
                  <button class="btn-chat" @click="startPrivateChat(c)">私聊</button>
                </div>
                <div v-if="contacts.length === 0" class="empty">暂无互关好友，搜索用户关注吧</div>
              </div>
            </div>
          </div>
        </div>

        <!-- 好友详情页 -->
        <div v-else-if="view === 'contact-detail'" class="profile-view">
          <div class="profile-header">
            <button class="menu-btn" @click="sidebarOpen = true">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="3" y1="6" x2="21" y2="6"/><line x1="3" y1="12" x2="21" y2="12"/><line x1="3" y1="18" x2="21" y2="18"/></svg>
            </button>
            <h2>用户资料</h2>
          </div>
          <div class="profile-content">
            <div class="profile-avatar-section">
              <div class="profile-avatar">
                <img v-if="contactDetail.avatar" :src="assetUrl(contactDetail.avatar)" />
                <span v-else>{{ (contactDetail.nickname || '').charAt(0) }}</span>
              </div>
            </div>
            <div class="profile-info">
              <h3>{{ contactDetail.nickname }}</h3>
              <p class="profile-username">@{{ contactDetail.username }}</p>
            </div>
            <div class="profile-stats">
              <div class="stat" @click="showRelationList('following', contactDetail)"><strong>{{ contactDetail.followingCount || 0 }}</strong><span>关注</span></div>
              <div class="stat" @click="showRelationList('followers', contactDetail)"><strong>{{ contactDetail.followersCount || 0 }}</strong><span>粉丝</span></div>
            </div>
            <p class="profile-bio">{{ contactDetail.bio || '暂无简介' }}</p>
            <div class="profile-actions">
              <button v-if="contactDetail.username !== username && !contactDetail.isFollowing" class="btn-primary" @click="doFollow(contactDetail)">关注</button>
              <button v-if="contactDetail.username !== username && contactDetail.isFollowing" class="btn-outline" @click="doUnfollow(contactDetail)">取消关注</button>
              <button v-if="contactDetail.username !== username && contactDetail.isMutual" class="btn-outline" @click="startPrivateChat(contactDetail)">发消息</button>
            </div>
          </div>
        </div>

        <!-- 关注/粉丝列表 -->
        <div v-else-if="view === 'relation-list'" class="profile-view">
          <div class="profile-header">
            <button class="menu-btn" @click="sidebarOpen = true">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="3" y1="6" x2="21" y2="6"/><line x1="3" y1="12" x2="21" y2="12"/><line x1="3" y1="18" x2="21" y2="18"/></svg>
            </button>
            <h2>{{ relationListType === 'following' ? '关注' : '粉丝' }}</h2>
          </div>
          <div class="relation-list">
            <div v-for="u in relationList" :key="u.username" class="contact-item">
              <div class="contact-avatar" @click="showContactDetail(u)">
                <img v-if="u.avatar" :src="assetUrl(u.avatar)" />
                <span v-else>{{ (u.nickname || '').charAt(0) }}</span>
              </div>
              <div class="contact-info" @click="showContactDetail(u)">
                <strong>{{ u.nickname }}</strong>
                <span>{{ u.username }}</span>
              </div>
            </div>
            <div v-if="relationList.length === 0" class="empty">暂无数据</div>
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

    <!-- 灯箱 -->
    <div v-if="lightboxUrl" class="lightbox" @click.self="closeLightbox">
      <button class="lightbox-close" @click="closeLightbox">&times;</button>
      <img :src="lightboxUrl" class="lightbox-img" @click.self="closeLightbox" />
    </div>
  </div>
</template>

<script>
const STORAGE_KEY = 'chat_user'

const API_BASE = import.meta.env.VITE_API_BASE || '/api'
const WS_BASE = import.meta.env.VITE_WS_BASE || ''
const ASSET_BASE = import.meta.env.VITE_ASSET_BASE || ''

function assetUrlFun(path) {
  if (!path) return ''
  if (!ASSET_BASE) return path
  return ASSET_BASE + path
}

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
      privateChatTo: '',
      privateChatUsername: '',
      text: '',
      connected: false,
      messages: [],
      users: [],
      contacts: [],
      contactDetail: {},
      profileFollowing: 0,
      profileFollowers: 0,
      relationList: [],
      relationListType: 'following',
      searchQuery: '',
      searchResults: [],
      ws: null,
      lightboxUrl: '',
      uploading: false,
    }
  },
  mounted() {
    document.addEventListener('keydown', this.onKeyDown)
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
        const res = await fetch(`${API_BASE}/login`, {
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
        const res = await fetch(`${API_BASE}/login`, {
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
        const res = await fetch(`${API_BASE}/register`, {
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
        const res = await fetch(`${API_BASE}/profile/update`, {
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
        const res = await fetch(`${API_BASE}/avatar/upload`, {
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
      this.uploading = true
      const form = new FormData()
      form.append('file', file)
      try {
        const res = await fetch(`${API_BASE}/media/upload`, { method: 'POST', body: form })
        const data = await res.json()
        if (data.success) {
          this.ws.send(JSON.stringify({
            type: 'message',
            content: data.mediaType === 'image' ? '' : file.name,
            room: this.currentRoom,
            mediaUrl: data.mediaUrl,
            mediaType: data.mediaType,
            fileSize: file.size,
          }))
        }
      } catch {}
      this.uploading = false
      this.$refs.fileInput.value = ''
    },
    assetUrl(path) {
      return assetUrlFun(path)
    },
    getFileExt(name) {
      const ext = (name || '').split('.').pop() || 'FILE'
      return ext.length <= 4 ? ext.toUpperCase() : 'FILE'
    },
    formatFileSize(bytes) {
      if (!bytes && bytes !== 0) return ''
      if (bytes < 1024) return bytes + ' B'
      if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
      return (bytes / 1024 / 1024).toFixed(1) + ' MB'
    },
    openFile(url) {
      window.open(assetUrl(url), '_blank')
    },
    scrollToBottom() {
      this.$nextTick(() => {
        const box = this.$refs.msgBox
        if (box) box.scrollTop = box.scrollHeight
      })
    },
    async loadContacts() {
      try {
        const [fwRes, frRes] = await Promise.all([
          fetch(`${API_BASE}/following?username=${encodeURIComponent(this.username)}`),
          fetch(`${API_BASE}/followers?username=${encodeURIComponent(this.username)}`),
        ])
        const fw = await fwRes.json()
        const fr = await frRes.json()
        const following = (fw.users || []).map(u => u.username)
        const followers = (fr.users || []).map(u => u.username)
        this.contacts = (fw.users || []).filter(u => followers.includes(u.username))
      } catch {}
    },
    async loadProfileStats() {
      try {
        const [fwRes, frRes] = await Promise.all([
          fetch(`${API_BASE}/following?username=${encodeURIComponent(this.username)}`),
          fetch(`${API_BASE}/followers?username=${encodeURIComponent(this.username)}`),
        ])
        const fw = await fwRes.json()
        const fr = await frRes.json()
        this.profileFollowing = (fw.users || []).length
        this.profileFollowers = (fr.users || []).length
      } catch {}
    },
    async showRelationList(type, user) {
      this.relationListType = type
      this.relationList = []
      this.view = 'relation-list'
      try {
        const res = await fetch(`${API_BASE}/${type}?username=${encodeURIComponent(user.username)}`)
        const data = await res.json()
        this.relationList = data.users || []
      } catch {}
      this.sidebarOpen = false
    },
    async addContact(contactUsername) {
      try {
        const res = await fetch(`${API_BASE}/contacts/add`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ username: this.username, contact: contactUsername }),
        })
        const data = await res.json()
        if (data.success) {
          this.loadContacts()
          this.searchQuery = ''
          this.searchResults = []
        }
      } catch {}
    },
    async searchUsers() {
      if (!this.searchQuery.trim()) { this.searchResults = []; return }
      try {
        const res = await fetch(`${API_BASE}/users/search?q=${encodeURIComponent(this.searchQuery)}`)
        const data = await res.json()
        this.searchResults = data.users || []
      } catch {}
    },
    isContact(username) {
      return this.contacts.some(c => c.username === username)
    },
    async showContactDetail(contact) {
      this.contactDetail = { ...contact }
      try {
        const [relRes, fwRes, frRes] = await Promise.all([
          fetch(`${API_BASE}/relation?username=${this.username}&target=${contact.username}`),
          fetch(`${API_BASE}/following?username=${contact.username}`),
          fetch(`${API_BASE}/followers?username=${contact.username}`),
        ])
        const rel = await relRes.json()
        const fw = await fwRes.json()
        const fr = await frRes.json()
        this.contactDetail.isFollowing = rel.following
        this.contactDetail.isMutual = rel.mutual
        this.contactDetail.followingCount = (fw.users || []).length
        this.contactDetail.followersCount = (fr.users || []).length
      } catch {}
      this.view = 'contact-detail'
      this.sidebarOpen = false
    },
    async showUserDetail(msg) {
      if (msg.nick === this.nickname) {
        this.view = 'profile'
        this.sidebarOpen = false
        return
      }
      try {
        const res = await fetch(`${API_BASE}/users/search?q=${encodeURIComponent(msg.nick)}`)
        const data = await res.json()
        if (data.success && data.users.length > 0) {
          const u = data.users[0]
          await this.showContactDetail({ username: u.username, nickname: u.nickname, avatar: u.avatar, bio: u.bio || '' })
        }
      } catch {}
      this.sidebarOpen = false
    },
    async doFollow(target) {
      try {
        const res = await fetch(`${API_BASE}/follow`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ username: this.username, target: target.username }),
        })
        const data = await res.json()
        if (data.success) {
          this.contactDetail.isFollowing = true
          this.contactDetail.isMutual = data.mutual
          this.contactDetail.followersCount = (this.contactDetail.followersCount || 0) + 1
          this.loadContacts()
        }
      } catch {}
    },
    async doUnfollow(target) {
      try {
        await fetch(`${API_BASE}/unfollow`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ username: this.username, target: target.username }),
        })
        this.contactDetail.isFollowing = false
        this.contactDetail.isMutual = false
        this.contactDetail.followersCount = Math.max(0, (this.contactDetail.followersCount || 0) - 1)
        this.loadContacts()
      } catch {}
    },
    startPrivateChat(contact) {
      this.privateChatTo = contact.nickname || contact.username
      this.privateChatUsername = contact.username
      this.messages = []
      const room = 'private:' + [this.username, contact.username].sort().join(':')
      this.currentRoom = room
      this.view = 'chat'
      this.sidebarOpen = false
      this.connect(room)
    },
    leavePrivateChat() {
      this.privateChatTo = ''
      this.privateChatUsername = ''
      this.messages = []
      this.currentRoom = '大厅'
      this.connect()
    },
    handleExit() {
      if (this.privateChatTo) {
        this.leavePrivateChat()
      } else {
        this.disconnect()
      }
    },
    switchView(v) {
      this.view = v
      this.sidebarOpen = false
      if (v === 'profile') this.loadProfileStats()
    },
    connect(room) {
      const r = room || this.room.trim() || '大厅'
      const host = WS_BASE || location.host
      const wsUrl = `${location.protocol === 'https:' ? 'wss:' : 'ws:'}//${host}/ws?username=${encodeURIComponent(this.username)}&nick=${encodeURIComponent(this.nickname)}&room=${encodeURIComponent(r)}`
      if (this.ws) this.ws.close()
      this.ws = new WebSocket(wsUrl)

      this.ws.onopen = () => {
        this.connected = true
        this.currentRoom = r
        this.loadContacts()
      }

      this.ws.onmessage = (e) => {
        const msg = JSON.parse(e.data)
        switch (msg.type) {
          case 'message':
            if (msg.room && msg.room.startsWith('private:') && !this.privateChatTo) {
              const parts = msg.room.split(':')
              const other = parts[1] === this.username ? parts[2] : parts[1]
              this.privateChatUsername = other
              const contact = this.contacts.find(c => c.username === other)
              this.privateChatTo = contact ? contact.nickname : other
            }
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
      const msg = {
        type: this.privateChatTo ? 'private' : 'message',
        content: this.text,
        room: this.currentRoom,
      }
      if (this.privateChatUsername) msg.to = this.privateChatUsername
      this.ws.send(JSON.stringify(msg))
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
    openLightbox(url) {
      this.lightboxUrl = url
    },
    closeLightbox() {
      this.lightboxUrl = ''
    },
    onKeyDown(e) {
      if (e.key === 'Escape') this.closeLightbox()
    },
  },
  beforeUnmount() {
    document.removeEventListener('keydown', this.onKeyDown)
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
.msg-file-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 4px 0;
}
.file-icon {
  width: 40px;
  height: 40px;
  background: #f0f0ee;
  border: 1px solid #e2e2de;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 10px;
  font-weight: 700;
  color: #666;
  flex-shrink: 0;
}
.file-info {
  flex: 1;
  min-width: 0;
  cursor: pointer;
}
.file-info strong {
  display: block;
  font-size: 13px;
  color: #1a1a1a;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.file-info span { font-size: 11px; color: #888; }
.msg-row.self .file-info strong { color: #fff; }
.msg-row.self .file-icon { background: #333; border-color: #444; color: #fff; }
.file-download {
  color: #999;
  display: flex;
  align-items: center;
  flex-shrink: 0;
}
.file-download:hover { color: #1a1a1a; }
.msg-row.self .file-download { color: rgba(255,255,255,0.6); }
.msg-row.self .file-download:hover { color: #fff; }
.uploading-indicator {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  font-size: 13px;
  color: #888;
  align-self: flex-start;
}
.spinner {
  width: 14px;
  height: 14px;
  border: 2px solid #e0e0e0;
  border-top-color: #1a1a1a;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }
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
  padding: 40px 24px;
  max-width: 400px;
  margin: 0 auto;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
}
.profile-avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 20px;
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
.profile-info { text-align: center; margin-bottom: 16px; }
.profile-info h3 { font-size: 20px; font-weight: 700; }
.profile-username { font-size: 13px; color: #888; margin-top: 4px; }
.profile-stats {
  display: flex;
  gap: 40px;
  justify-content: center;
  margin-bottom: 16px;
}
.stat { text-align: center; cursor: pointer; }
.stat:hover strong { color: #1a1a1a; }
.stat strong { display: block; font-size: 18px; font-weight: 700; }
.stat span { font-size: 12px; color: #888; }
.relation-list {
  flex: 1;
  overflow-y: auto;
  padding: 16px 20px;
  max-width: 480px;
  margin: 0 auto;
  width: 100%;
}
.profile-bio {
  font-size: 14px;
  color: #444;
  line-height: 1.5;
  text-align: center;
  margin-bottom: 24px;
  max-width: 320px;
}
.profile-actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
  width: 100%;
  max-width: 320px;
}
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

/* Contacts */
.contacts-view { display: flex; flex-direction: column; height: 100vh; }
.contacts-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 20px;
  background: #f9f9f7;
  border-bottom: 1px solid #deded9;
}
.contacts-header h2 { font-size: 16px; font-weight: 600; }
.contacts-body { flex: 1; overflow-y: auto; padding: 20px; max-width: 720px; margin: 0 auto; width: 100%; }
.contacts-search { margin-bottom: 24px; position: relative; }
.contacts-search .search-icon {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: #999;
}
.contacts-search input {
  width: 100%;
  padding: 12px 16px 12px 36px;
  border: 1px solid #cececa;
  border-radius: 10px;
  font-size: 14px;
  outline: none;
  background: #fff;
  transition: border-color 0.2s;
}
.contacts-search input:focus { border-color: #1a1a1a; }
.contacts-search input::placeholder { color: #aaa; }
.search-results {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: #fff;
  border: 1px solid #e0e0e0;
  border-radius: 10px;
  margin-top: 6px;
  z-index: 10;
  max-height: 320px;
  overflow-y: auto;
  box-shadow: 0 4px 16px rgba(0,0,0,0.08);
}
.search-title {
  padding: 10px 14px 6px;
  font-size: 11px;
  font-weight: 600;
  color: #999;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}
.search-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 14px;
  cursor: pointer;
}
.search-item:hover { background: #f5f5f5; }
.search-avatar, .contact-avatar {
  width: 40px;
  height: 40px;
  background: #d4d4d0;
  color: #555;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
  overflow: hidden;
  flex-shrink: 0;
}
.search-avatar img, .contact-avatar img { width: 100%; height: 100%; object-fit: cover; }
.search-info { flex: 1; min-width: 0; }
.search-info strong { display: block; font-size: 14px; }
.search-info span { font-size: 12px; color: #888; }
.btn-add {
  padding: 6px 14px;
  background: #1a1a1a;
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 12px;
  cursor: pointer;
  flex-shrink: 0;
  transition: all 0.15s;
}
.btn-add:hover { background: #333; }
.btn-add.followed { background: transparent; color: #666; border: 1px solid #cececa; }
.btn-add.followed:hover { border-color: #e53935; color: #e53935; }
.contacts-section { margin-top: 4px; }
.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 4px 10px;
  font-size: 13px;
  font-weight: 600;
  color: #666;
}
.section-count {
  background: #ececea;
  color: #888;
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 10px;
}
.contacts-list { display: flex; flex-direction: column; gap: 4px; }
@media (min-width: 900px) {
  .contact-item { padding: 12px 14px; }
  .contacts-body { max-width: 900px; }
  .contacts-list { display: grid; grid-template-columns: repeat(2, 1fr); gap: 8px; }
}
.contact-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  border-radius: 10px;
  cursor: pointer;
  transition: background 0.15s;
}
.contact-item:hover { background: #f0f0ee; }
.contact-info { flex: 1; min-width: 0; }
.contact-info strong { display: block; font-size: 14px; }
.contact-info span { font-size: 12px; color: #888; }
.btn-chat {
  padding: 6px 14px;
  background: transparent;
  color: #666;
  border: 1px solid #cececa;
  border-radius: 8px;
  font-size: 12px;
  cursor: pointer;
  flex-shrink: 0;
  transition: all 0.15s;
}
.btn-chat:hover { border-color: #1a1a1a; color: #1a1a1a; }
.btn-outline {
  width: 100%;
  padding: 12px;
  background: transparent;
  border: 1px solid #cececa;
  border-radius: 7px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.15s;
}
.btn-outline:hover { border-color: #1a1a1a; color: #1a1a1a; }
.contact-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 8px;
  cursor: pointer;
}
.contact-item:hover { background: #f0f0ee; }
.contact-info { flex: 1; }
.contact-info strong { display: block; font-size: 14px; }
.contact-info span { font-size: 12px; color: #888; }
.back-room { font-size: 12px; color: #888; cursor: pointer; text-decoration: underline; }

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

/* Lightbox */
.lightbox {
  position: fixed;
  inset: 0;
  z-index: 200;
  background: rgba(0,0,0,0.85);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
  cursor: zoom-out;
}
.lightbox-img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
  border-radius: 4px;
  cursor: default;
}
.lightbox-close {
  position: absolute;
  top: 16px;
  right: 20px;
  background: none;
  border: none;
  color: #fff;
  font-size: 32px;
  cursor: pointer;
  opacity: 0.7;
  line-height: 1;
}
.lightbox-close:hover { opacity: 1; }

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
