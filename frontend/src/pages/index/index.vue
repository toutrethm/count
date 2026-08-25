<template>
  <view class="page">
    <view class="summary">
      <view>
        <text class="eyebrow">工作台</text>
        <text class="title">{{ userLabel }}</text>
      </view>
      <text class="role">{{ isAdmin ? '管理员' : capabilityLabel }}</text>
    </view>

    <view v-if="isWorker" class="panel">
      <text class="panel-title">扫码记录工序</text>
      <view class="field-row">
        <input v-model="scanForm.orderNo" class="input" placeholder="订单号 / 二维码内容" />
        <button class="mini-button" @click="handleCameraScan">扫码</button>
      </view>
      <button class="secondary-button" @click="handlePreviewScan">查看订单内容</button>
      <button class="primary-button" @click="handleRecordScan">确认记录</button>
      <text class="hint">系统会按订单当前流程判断是否轮到你的能力，无需手动选择工序。</text>
    </view>

    <!-- #ifdef H5 -->
    <h5-qr-scanner
      v-if="showH5Scanner"
      v-model="showH5Scanner"
      @detected="handleH5ScanDetected"
      @error="handleH5ScanError"
    />
    <!-- #endif -->

    <view v-if="previewOrder" class="panel">
      <view class="preview-head">
        <view>
          <text class="panel-title">订单内容</text>
          <text class="preview-meta">订单号：{{ previewOrder.order_no }}</text>
        </view>
        <text :class="previewCanRecord ? 'done-tag' : 'wait-tag'">
          {{ previewCanRecord ? '可记录' : '未轮到你' }}
        </text>
      </view>

      <view class="current-process">
        <text class="current-title">本人对应工序</text>
        <template v-if="previewCurrentProcess">
          <text class="current-line">{{ processItemName(previewCurrentProcess) }} / {{ processName(previewCurrentProcess) }}</text>
          <text class="current-line">{{ formatCapabilityListLabel(previewCurrentProcess.station_role) }}</text>
        </template>
        <text v-else class="current-line">当前订单暂无可记录工序</text>
      </view>

      <view v-if="previewMatchingProcesses.length" class="process-list">
        <view v-for="process in previewMatchingProcesses" :key="process.id" class="process-item">
          <text class="process-name">{{ processItemName(process) }} / {{ processName(process) }}</text>
          <text class="process-meta">{{ formatCapabilityListLabel(process.station_role) }}</text>
        </view>
      </view>

      <view class="order-items">
        <view v-for="item in previewOrder.items || []" :key="item.id" class="order-item">
          <text class="order-item-name">{{ item.part_name }}</text>
          <text class="order-item-spec">{{ item.spec }}</text>
          <text class="order-item-qty">x{{ item.quantity }}</text>
        </view>
      </view>
    </view>

    <view v-if="isWorker" class="panel">
      <text class="panel-title">我的可扫工序</text>
      <view v-if="myProcesses.length" class="tag-list">
        <text v-for="item in myProcesses" :key="item.id" class="tag">{{ formatCapabilityListLabel(item.station_role) }}</text>
      </view>
      <text v-else class="empty">暂无匹配工序，请确认账号能力是否正确。</text>
    </view>

    <view v-if="isAdmin" class="panel">
      <text class="panel-title">录入订单</text>
      <view class="item-toolbar">
        <text class="hint">每行填写一个元件，系统自动生成订单号、二维码和生产流程。</text>
        <button class="mini-button add-button" @click="addOrderItem">新增元件</button>
      </view>

      <view v-for="(item, index) in orderItems" :key="index" class="item-card">
        <view class="item-head">
          <text class="item-title">元件 {{ index + 1 }}</text>
          <button v-if="orderItems.length > 1" class="remove-button" @click="removeOrderItem(index)">删除</button>
        </view>

        <picker :range="componentTypes" range-key="label" @change="(event) => onComponentTypeChange(event, index)">
          <view class="picker">{{ getComponentDefinition(item.typeKey).label }}</view>
        </picker>

        <input v-model="item.quantity" class="input" placeholder="数量" type="number" />

        <view class="dimension-grid">
          <view v-for="field in getComponentDefinition(item.typeKey).fields" :key="field.key" class="dimension-item">
            <text class="dimension-label">{{ field.label }}</text>
            <input v-model="item.dimensions[field.key]" class="input compact" :placeholder="field.label" />
          </view>
        </view>
      </view>

      <button class="primary-button" @click="handleCreateOrder">提交订单</button>
    </view>

    <view v-if="isAdmin && createdOrder" class="panel print-panel">
      <view class="print-toolbar no-print">
        <text class="panel-title">A5 打印预览</text>
        <button class="primary-button" @click="handlePrint">打印A5</button>
      </view>

      <view class="print-sheet">
        <view class="print-header">
          <view>
            <text class="print-title">订单打印单</text>
            <text class="print-subtitle">订单号：{{ createdOrder.order_no }}</text>
            <text class="print-subtitle">提交时间：{{ formatTime(createdOrder.created_at) }}</text>
          </view>
          <qr-code-preview :value="createdOrder.order_no" :size="128" />
        </view>

        <view class="print-order">
          <text class="print-line">二维码内容：{{ createdOrder.order_no }}</text>
          <text class="print-line">总数量：{{ totalCreatedQuantity }}</text>
        </view>

        <view class="print-table">
          <view class="print-row print-head-row">
            <text class="cell no">序号</text>
            <text class="cell name">元件</text>
            <text class="cell spec">尺寸</text>
            <text class="cell qty">数量</text>
          </view>
          <view v-for="(item, index) in createdOrder.items || []" :key="item.id || `${item.item_no}-${index}`" class="print-row">
            <text class="cell no">{{ item.item_no || index + 1 }}</text>
            <text class="cell name">{{ item.part_name }}</text>
            <text class="cell spec">{{ item.spec }}</text>
            <text class="cell qty">{{ item.quantity }}</text>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed, reactive, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import QrCodePreview from '../../components/qr-code/QrCodePreview.vue'
import { getMe } from '../../api/auth'
import { createOrder } from '../../api/order'
import { previewScanOrder, recordScan } from '../../api/scan'
import { getUser } from '../../utils/session'
import { scanQrCode } from '../../utils/qr'
import { redirectTo } from '../../utils/navigation'
import { formatCapabilityListLabel } from '../../utils/capability-labels'
import H5QrScanner from '../../components/h5-scanner/H5QrScanner.vue'
import {
  buildOrderItemPayload,
  createOrderItemDraft,
  getComponentDefinition,
  getComponentTypes,
} from '../../utils/order-item-form'

const user = ref(getUser())
const myProcesses = ref([])
const createdOrder = ref(null)
const componentTypes = getComponentTypes()
const showH5Scanner = ref(false)
const scanPreview = ref(null)
const scanRecordResult = ref(null)

const scanForm = reactive({
  orderNo: '',
})

const orderItems = ref([createOrderItemDraft()])

const isAdmin = computed(() => user.value?.role === 'admin')
const isWorker = computed(() => !isAdmin.value)
const userLabel = computed(() => `${user.value?.username || '-'} / ${user.value?.phone || '-'}`)
const capabilityLabel = computed(() => formatCapabilityListLabel(user.value?.station_role))
const previewOrder = computed(() => scanRecordResult.value?.order || scanPreview.value?.order)
const previewCurrentProcess = computed(() => {
  if (scanRecordResult.value?.record && previewOrder.value?.processes) {
    return findProcessByID(previewOrder.value.processes, scanRecordResult.value.record.order_process_id)
  }
  return scanPreview.value?.pending_process || null
})
const previewMatchingProcesses = computed(() => scanPreview.value?.matching_processes || [])
const previewCanRecord = computed(() => Boolean(scanPreview.value?.can_record) || Boolean(scanRecordResult.value?.record))
const totalCreatedQuantity = computed(() =>
  (createdOrder.value?.items || []).reduce((sum, item) => sum + Number(item.quantity || 0), 0)
)

onShow(() => {
  if (!user.value) {
    redirectTo('/pages/login/login')
    return
  }
  loadPageData()
})

async function loadPageData() {
  const me = await getMe()
  user.value = me.user
  myProcesses.value = me.processes || []
}

function onComponentTypeChange(event, index) {
  const typeKey = componentTypes[Number(event.detail.value || 0)]?.key || componentTypes[0].key
  const next = createOrderItemDraft(typeKey)
  const current = orderItems.value[index]
  if (!current) return
  current.typeKey = next.typeKey
  current.dimensions = next.dimensions
}

function addOrderItem() {
  orderItems.value.push(createOrderItemDraft())
}

function removeOrderItem(index) {
  if (orderItems.value.length <= 1) return
  orderItems.value.splice(index, 1)
}

async function handleCameraScan() {
  if (isBrowserRuntime()) {
    showH5Scanner.value = true
    return
  }

  const result = await scanQrCode()
  scanForm.orderNo = result.text
  await loadScanPreview(result.text)
}

async function handleH5ScanDetected(text) {
  scanForm.orderNo = text
  showH5Scanner.value = false
  await loadScanPreview(text)
}

function handleH5ScanError(error) {
  uni.showToast({ title: error?.message || '扫码失败', icon: 'none' })
}

function isBrowserRuntime() {
  return typeof window !== 'undefined' && typeof document !== 'undefined'
}

async function handlePreviewScan() {
  if (!scanForm.orderNo) {
    uni.showToast({ title: '请输入订单号', icon: 'none' })
    return
  }
  await loadScanPreview(scanForm.orderNo)
}

async function loadScanPreview(orderNo) {
  const value = String(orderNo || '').trim()
  if (!value) return

  try {
    scanRecordResult.value = null
    scanPreview.value = await previewScanOrder(value)
    uni.showToast({ title: '已读取订单', icon: 'success' })
  } catch (error) {
    scanPreview.value = null
    scanRecordResult.value = null
    uni.showToast({ title: previewErrorTitle(error), icon: 'none' })
  }
}

async function handleRecordScan() {
  if (!scanForm.orderNo) {
    uni.showToast({ title: '请输入订单号', icon: 'none' })
    return
  }

  try {
    const result = await recordScan({
      qr_token: scanForm.orderNo,
    })
    scanRecordResult.value = result
    scanPreview.value = {
      order: result.order,
      pending_process: findProcessByID(result.order?.processes || [], result.record?.order_process_id),
      matching_processes: scanPreview.value?.matching_processes || [],
      can_record: false,
    }
    scanForm.orderNo = ''
    uni.showToast({ title: '已记录', icon: 'success' })
  } catch (error) {
    uni.showToast({ title: previewErrorTitle(error), icon: 'none' })
  }
}

function previewErrorTitle(error) {
  const message = error?.data?.message || error?.message || ''
  if (message.includes('permission')) return '当前账号不能查看该订单'
  if (message.includes('not found')) return '订单不存在'
  if (message.includes('limit')) return '订单扫码次数已满'
  return '未轮到当前能力或订单不存在'
}

function findProcessByID(processes = [], id) {
  return processes.find((item) => item.id === id) || null
}

function processName(item) {
  return item?.process?.name || `工序 ${item?.process_id || '-'}`
}

function processItemName(item) {
  return item?.order_item?.part_name || `元件 ${item?.order_item_id || '-'}`
}

async function handleCreateOrder() {
  const payloadItems = []
  for (let i = 0; i < orderItems.value.length; i += 1) {
    const item = orderItems.value[i]
    if (!Number(item.quantity) || Number(item.quantity) <= 0) {
      uni.showToast({ title: `第 ${i + 1} 个元件数量不正确`, icon: 'none' })
      return
    }
    const definition = getComponentDefinition(item.typeKey)
    const missingField = definition.fields.find((field) => !String(item.dimensions[field.key] || '').trim())
    if (missingField) {
      uni.showToast({ title: `第 ${i + 1} 个元件缺少${missingField.label}`, icon: 'none' })
      return
    }
    payloadItems.push(buildOrderItemPayload(item, i))
  }

  try {
    const result = await createOrder({
      items: payloadItems,
      status: 'draft',
    })
    createdOrder.value = result.item
    orderItems.value = [createOrderItemDraft()]
    uni.showToast({ title: '订单已创建', icon: 'success' })
  } catch (error) {
    uni.showToast({ title: '创建失败', icon: 'none' })
  }
}

function handlePrint() {
  if (typeof window !== 'undefined') {
    window.print()
  }
}

function formatTime(value) {
  if (!value) return '-'
  return String(value).replace('T', ' ').replace(/\.\d+(Z)?$/, '')
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

.eyebrow,
.hint {
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

.role {
  padding: 8rpx 14rpx;
  border-radius: 999rpx;
  background: rgba(255, 255, 255, 0.18);
  font-size: 24rpx;
}

.panel {
  margin-top: 22rpx;
  padding: 24rpx;
  border-radius: 16rpx;
  background: #ffffff;
  border: 1rpx solid #e2e7ee;
}

.panel-title {
  display: block;
  margin-bottom: 18rpx;
  font-size: 30rpx;
  font-weight: 700;
}

.input,
.picker {
  width: 100%;
  box-sizing: border-box;
  border: 1rpx solid #d7dfe7;
  border-radius: 12rpx;
  background: #fafbfc;
}

.input {
  height: 84rpx;
  padding: 0 20rpx;
  margin-bottom: 14rpx;
}

.input.compact {
  margin-bottom: 0;
}

.picker {
  min-height: 84rpx;
  padding: 24rpx 20rpx;
  margin-bottom: 14rpx;
}

.field-row {
  display: flex;
  gap: 12rpx;
  align-items: center;
}

.field-row .input {
  flex: 1;
}

.mini-button {
  flex: 0 0 148rpx;
  height: 84rpx;
  line-height: 84rpx;
  background: #eff3f7;
}

.primary-button {
  background: #1f6f5b;
  color: #ffffff;
}

.secondary-button {
  margin-bottom: 14rpx;
  background: #eff3f7;
  color: #1f6f5b;
}

.preview-head,
.process-item {
  display: flex;
  justify-content: space-between;
  gap: 16rpx;
  align-items: center;
}

.preview-meta,
.current-line,
.process-meta {
  display: block;
  margin-top: 8rpx;
  color: #7a8494;
  font-size: 24rpx;
}

.done-tag,
.wait-tag {
  padding: 8rpx 14rpx;
  border-radius: 999rpx;
  font-size: 22rpx;
  font-weight: 700;
}

.done-tag {
  background: #e7f5ef;
  color: #1f6f5b;
}

.wait-tag {
  background: #fff4e6;
  color: #a45b18;
}

.current-process {
  margin-top: 16rpx;
  padding: 16rpx;
  border-radius: 12rpx;
  background: #f8fafc;
}

.current-title,
.process-name,
.order-item-name {
  display: block;
  font-weight: 700;
  color: #18212f;
}

.process-list {
  margin-top: 14rpx;
  border-radius: 12rpx;
  background: #f6f8fa;
  overflow: hidden;
}

.process-item {
  padding: 14rpx 16rpx;
  border-bottom: 1rpx solid #e8edf2;
  font-size: 24rpx;
}

.process-item:last-child {
  border-bottom: 0;
}

.order-items {
  margin-top: 14rpx;
  display: flex;
  flex-direction: column;
  gap: 10rpx;
}

.order-item {
  display: grid;
  grid-template-columns: 1.2fr 1.8fr 0.6fr;
  gap: 12rpx;
  padding: 12rpx 14rpx;
  border-radius: 10rpx;
  background: #f8fafc;
  font-size: 24rpx;
}

.order-item-spec,
.order-item-qty {
  color: #5f6a7a;
}

.tag-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12rpx;
}

.tag {
  padding: 10rpx 16rpx;
  border-radius: 999rpx;
  background: #eff3f7;
  color: #2f3b4a;
  font-size: 24rpx;
}

.empty {
  color: #7a8494;
  font-size: 26rpx;
}

.item-toolbar,
.item-head,
.print-toolbar,
.print-header {
  display: flex;
  justify-content: space-between;
  gap: 16rpx;
  align-items: center;
}

.item-toolbar {
  margin-bottom: 16rpx;
}

.add-button,
.remove-button {
  background: #eff3f7;
}

.item-card {
  margin-bottom: 18rpx;
  padding: 18rpx;
  border-radius: 14rpx;
  background: #f8fafc;
  border: 1rpx solid #e5ebf2;
}

.item-title {
  font-size: 28rpx;
  font-weight: 700;
}

.dimension-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12rpx;
}

.dimension-item {
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.dimension-label {
  font-size: 24rpx;
  color: #5f6a7a;
}

.print-panel {
  overflow: hidden;
}

.print-sheet {
  width: 100%;
  max-width: 148mm;
  min-height: 210mm;
  margin: 0 auto;
  padding: 10mm;
  box-sizing: border-box;
  background: #ffffff;
  color: #111827;
}

.print-title {
  display: block;
  font-size: 32rpx;
  font-weight: 700;
}

.print-subtitle,
.print-line {
  display: block;
  margin-top: 6rpx;
  font-size: 22rpx;
}

.print-order {
  margin-top: 16rpx;
}

.print-table {
  margin-top: 18rpx;
  border: 1rpx solid #dbe3ea;
}

.print-row {
  display: flex;
  border-top: 1rpx solid #dbe3ea;
}

.print-row:first-child {
  border-top: 0;
}

.print-head-row {
  background: #f4f6f8;
  font-weight: 700;
}

.cell {
  padding: 10rpx 8rpx;
  box-sizing: border-box;
  font-size: 22rpx;
  border-left: 1rpx solid #dbe3ea;
  word-break: break-all;
}

.cell:first-child {
  border-left: 0;
}

.cell.no {
  width: 12%;
}

.cell.name {
  width: 23%;
}

.cell.spec {
  width: 48%;
}

.cell.qty {
  width: 16%;
}

@media print {
  .no-print {
    display: none !important;
  }

  .page {
    padding: 0;
    background: #ffffff;
  }

  .panel:not(.print-panel),
  .summary {
    display: none !important;
  }

  .print-panel {
    margin: 0;
    padding: 0;
    border: 0;
  }

  .print-sheet {
    max-width: none;
    min-height: auto;
    padding: 8mm;
  }
}
</style>
