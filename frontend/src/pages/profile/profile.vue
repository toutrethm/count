<template>
  <view class="page">
    <view class="summary">
      <view>
        <text class="eyebrow">个人</text>
        <text class="title">{{ user?.username || '-' }}</text>
      </view>
    </view>

    <view class="panel">
      <view class="info-row">
        <text class="label">用户名</text>
        <text class="value">{{ user?.username || '-' }}</text>
      </view>
      <view class="info-row">
        <text class="label">电话</text>
        <text class="value">{{ user?.phone || '-' }}</text>
      </view>
    </view>

    <button class="logout" @click="handleLogout">退出登录</button>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { getMe } from '../../api/auth'
import { clearSession, getUser } from '../../utils/session'
import { redirectTo } from '../../utils/navigation'

const user = ref(getUser())

onShow(() => {
  if (!user.value) {
    redirectTo('/pages/login/login')
    return
  }
  loadProfile()
})

async function loadProfile() {
  const me = await getMe()
  user.value = me.user
}

function handleLogout() {
  clearSession()
  redirectTo('/pages/login/login')
}
</script>

<style scoped>
.page {
  min-height: 100vh;
  padding: 28rpx;
  box-sizing: border-box;
}

.summary {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24rpx;
  border-radius: 16rpx;
  background: #1f6f5b;
  color: #ffffff;
}

.eyebrow {
  display: block;
  font-size: 24rpx;
  opacity: 0.78;
}

.title {
  display: block;
  margin-top: 8rpx;
  font-size: 34rpx;
  font-weight: 700;
}

.panel {
  margin-top: 22rpx;
  padding: 24rpx;
  border-radius: 16rpx;
  background: #ffffff;
  border: 1rpx solid #e2e7ee;
}

.info-row {
  display: flex;
  justify-content: space-between;
  gap: 24rpx;
  padding: 18rpx 0;
  border-bottom: 1rpx solid #eef2f5;
}

.info-row:last-child {
  border-bottom: 0;
}

.label {
  color: #7a8494;
  font-size: 26rpx;
}

.value {
  color: #18212f;
  font-size: 26rpx;
  font-weight: 700;
}

.logout {
  margin-top: 28rpx;
  background: #c2412d;
  color: #ffffff;
}
</style>
