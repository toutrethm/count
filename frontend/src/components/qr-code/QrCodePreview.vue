<template>
  <view class="qr-wrap" :style="wrapStyle">
    <svg class="qr-svg" :viewBox="`0 0 ${moduleCount} ${moduleCount}`" aria-hidden="true">
      <rect x="0" y="0" :width="moduleCount" :height="moduleCount" fill="#ffffff" />
      <rect
        v-for="rect in rects"
        :key="`${rect.x}-${rect.y}`"
        :x="rect.x"
        :y="rect.y"
        width="1"
        height="1"
        fill="#111827"
      />
    </svg>
  </view>
</template>

<script setup>
import { computed } from 'vue'
import UQRCode from 'uqrcodejs'
import { buildQrRectangles } from '../../utils/qr-render'

const props = defineProps({
  value: {
    type: String,
    required: true,
  },
  size: {
    type: Number,
    default: 240,
  },
})

const modules = computed(() => {
  const qr = new UQRCode()
  qr.data = props.value || 'empty'
  qr.make()
  return qr.modules
})

const moduleCount = computed(() => modules.value.length || 1)
const rects = computed(() => buildQrRectangles(modules.value))
const wrapStyle = computed(() => `width:${props.size}px;height:${props.size}px;`)
</script>

<style scoped>
.qr-wrap {
  display: block;
  margin: 0 auto;
  padding: 16rpx;
  border: 1rpx solid #e5ebf2;
  background: #ffffff;
  box-sizing: content-box;
  -webkit-print-color-adjust: exact;
  print-color-adjust: exact;
}

.qr-svg {
  display: block;
  width: 100%;
  height: 100%;
  shape-rendering: crispEdges;
  -webkit-print-color-adjust: exact;
  print-color-adjust: exact;
}
</style>
