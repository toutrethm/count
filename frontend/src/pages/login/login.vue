<template>
  <view class="page">
    <view class="panel">
      <text class="title">工厂计件登录</text>
      <input v-model="form.account" class="input" placeholder="手机号 / 管理员账号" />
      <input v-model="form.password" class="input" password placeholder="密码" />
      <button class="primary-button" @click="handleLogin">登录</button>
      <button class="ghost-button" @click="goRegister">注册工人</button>
    </view>
  </view>
</template>

<script setup>
import { reactive, onMounted } from 'vue'
import { login } from '../../api/auth'
import { getToken, setSession } from '../../utils/session'
import { navigateTo, switchTabTo } from '../../utils/navigation'

const form = reactive({
  account: '',
  password: '',
})

onMounted(() => {
  if (getToken()) {
    switchTabTo('/pages/index/index')
  }
})

async function handleLogin() {
  if (!form.account || !form.password) {
    uni.showToast({ title: '请填写账号和密码', icon: 'none' })
    return
  }

  try {
    const result = await login(form)
    setSession(result.token, result.user)
    switchTabTo('/pages/index/index')
  } catch (error) {
    uni.showToast({ title: '登录失败', icon: 'none' })
  }
}

function goRegister() {
  navigateTo('/pages/register/register')
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
  height: 88rpx;
  margin-bottom: 18rpx;
  padding: 0 24rpx;
  box-sizing: border-box;
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
</style>
