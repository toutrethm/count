<template>
  <view class="page">
    <view class="panel">
      <text class="title">工人注册</text>
      <input v-model="form.username" class="input" placeholder="真实姓名" />
      <input v-model="form.phone" class="input" placeholder="电话" />
      <input v-model="form.password" class="input" password placeholder="密码" />
      <button class="primary-button" @click="handleRegister">注册</button>
      <button class="ghost-button" @click="goLogin">返回登录</button>
      <text class="hint">注册后由管理员分配可扫码工序。</text>
    </view>
  </view>
</template>

<script setup>
import { reactive } from 'vue'
import { registerWorker } from '../../api/auth'
import { navigateTo, redirectTo } from '../../utils/navigation'

const form = reactive({
  username: '',
  phone: '',
  password: '',
})

async function handleRegister() {
  if (!form.username || !form.phone || !form.password) {
    uni.showToast({ title: '请完整填写信息', icon: 'none' })
    return
  }

  try {
    await registerWorker({
      username: form.username,
      phone: form.phone,
      password: form.password,
    })
    uni.showToast({ title: '注册成功', icon: 'success' })
    redirectTo('/pages/login/login')
  } catch (error) {
    uni.showToast({ title: '注册失败', icon: 'none' })
  }
}

function goLogin() {
  navigateTo('/pages/login/login')
}
</script>

<style scoped>
.page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 32rpx;
  box-sizing: border-box;
}

.panel {
  width: 100%;
  padding: 32rpx;
  border-radius: 16rpx;
  background: #ffffff;
  border: 1rpx solid #e2e7ee;
}

.title {
  display: block;
  margin-bottom: 28rpx;
  font-size: 38rpx;
  font-weight: 700;
}

.input {
  width: 100%;
  box-sizing: border-box;
  height: 88rpx;
  margin-bottom: 18rpx;
  padding: 0 24rpx;
  border: 1rpx solid #d7dfe7;
  border-radius: 12rpx;
  background: #fafbfc;
}

.primary-button,
.ghost-button {
  margin-top: 10rpx;
}

.primary-button {
  background: #1f6f5b;
  color: #ffffff;
}

.ghost-button {
  border: 1rpx solid #d7dfe7;
  background: #ffffff;
}

.hint {
  display: block;
  margin-top: 18rpx;
  color: #7a8494;
  font-size: 24rpx;
}
</style>
