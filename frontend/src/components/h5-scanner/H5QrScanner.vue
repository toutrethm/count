<template>
  <view v-if="modelValue" class="overlay" @touchmove.stop.prevent>
    <view class="sheet">
      <view class="head">
        <text class="title">摄像头扫码</text>
        <button class="ghost" @click="close">关闭</button>
      </view>

      <view class="preview">
        <video ref="videoRef" class="video" autoplay muted playsinline></video>
        <canvas ref="canvasRef" class="canvas"></canvas>
      </view>

      <text class="status">{{ statusText }}</text>

      <view class="actions">
        <button class="primary" @click="restart">重新扫码</button>
        <button class="ghost" @click="close">取消</button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { nextTick, onBeforeUnmount, ref, watch } from 'vue'
import { createH5CameraScanSession, isH5CameraScanSupported } from '../../utils/h5-scanner'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue', 'detected', 'error'])

const videoRef = ref(null)
const canvasRef = ref(null)
const statusText = ref('请将二维码放入取景框')

let session = null

watch(
  () => props.modelValue,
  async (visible) => {
    if (visible) {
      await start()
    } else {
      stop()
    }
  },
  { immediate: true }
)

async function start() {
  stop()
  statusText.value = '正在打开摄像头...'

  if (!isH5CameraScanSupported()) {
    statusText.value = '当前浏览器不支持摄像头扫码'
    emit('error', new Error(statusText.value))
    return
  }

  await nextTick()
  if (!videoRef.value || !canvasRef.value) {
    statusText.value = '扫码视图未就绪'
    emit('error', new Error(statusText.value))
    return
  }

  session = createH5CameraScanSession({
    videoEl: videoRef.value,
    canvasEl: canvasRef.value,
  })

  statusText.value = '请将二维码对准取景框'
  session.promise
    .then(({ text }) => {
      emit('detected', text)
      close()
    })
    .catch((error) => {
      if (error?.code === 'H5_SCAN_CANCELLED') {
        return
      }
      statusText.value = error?.message || '扫码失败'
      emit('error', error)
    })
}

function stop() {
  if (!session) {
    return
  }
  session.stop()
  session = null
}

function close() {
  stop()
  emit('update:modelValue', false)
}

async function restart() {
  await start()
}

onBeforeUnmount(() => {
  stop()
})
</script>

<style scoped>
.overlay {
  position: fixed;
  inset: 0;
  z-index: 999;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 32rpx;
  box-sizing: border-box;
  background: rgba(15, 23, 42, 0.72);
}

.sheet {
  width: 100%;
  max-width: 720rpx;
  border-radius: 18rpx;
  background: #ffffff;
  padding: 24rpx;
  box-sizing: border-box;
}

.head,
.actions {
  display: flex;
  justify-content: space-between;
  gap: 16rpx;
  align-items: center;
}

.title {
  font-size: 32rpx;
  font-weight: 700;
}

.preview {
  position: relative;
  margin-top: 18rpx;
  border-radius: 16rpx;
  overflow: hidden;
  background: #0f172a;
  aspect-ratio: 3 / 4;
}

.video,
.canvas {
  width: 100%;
  height: 100%;
  display: block;
}

.canvas {
  position: absolute;
  left: 0;
  top: 0;
  width: 100%;
  height: 100%;
  opacity: 0;
  pointer-events: none;
}

.status {
  display: block;
  margin-top: 16rpx;
  color: #475569;
  font-size: 24rpx;
}

.actions {
  margin-top: 18rpx;
}

.ghost,
.primary {
  flex: 1;
}

.primary {
  background: #1f6f5b;
  color: #ffffff;
}

.ghost {
  background: #eff3f7;
}
</style>
