<template>
  <view class="page">
    <view class="summary">
      <view>
        <text class="eyebrow">记录</text>
        <text class="title">{{ isAdmin ? '全部扫码流水' : '我的工资流水' }}</text>
      </view>
      <button class="refresh" @click="loadPageData">刷新</button>
    </view>

    <view class="metrics">
      <view class="metric">
        <text class="metric-label">今日</text>
        <text class="metric-value">{{ formatMoney(todayWage) }}</text>
      </view>
      <view class="metric">
        <text class="metric-label">本月</text>
        <text class="metric-value">{{ formatMoney(monthWage) }}</text>
      </view>
      <view class="metric">
        <text class="metric-label">累计</text>
        <text class="metric-value">{{ formatMoney(totalWage) }}</text>
      </view>
    </view>

    <view v-if="recordTabs.length > 1" class="sub-tabs">
      <button
        v-for="item in recordTabs"
        :key="item.key"
        class="sub-tab"
        :class="{ active: activeTab === item.key }"
        @click="activeTab = item.key"
      >
        {{ item.label }}
      </button>
    </view>

    <view v-if="activeTab === 'scan'" class="statement-list">
      <view v-if="wageMonths.length">
        <view v-for="month in wageMonths" :key="month.key" class="month-card">
          <view class="month-head">
            <view>
              <text class="month-title">{{ month.label }}</text>
              <text class="month-subtitle">{{ month.days.length }} 天 · {{ month.orderCount }} 个订单</text>
            </view>
            <view class="month-total">
              <text class="amount-label">当月累计工资</text>
              <text class="amount-value">{{ formatMoney(month.totalWage) }}</text>
            </view>
          </view>

          <view v-for="day in month.days" :key="day.key" class="day-section">
            <view class="day-head">
              <view>
                <text class="day-title">{{ day.label }}</text>
                <text class="day-week">{{ day.weekday }} · {{ day.orders.length }} 个订单</text>
              </view>
              <text class="day-total">{{ formatMoney(day.totalWage) }}</text>
            </view>

            <view class="order-list">
              <view v-for="order in day.orders" :key="order.key" class="statement-row">
                <view class="statement-main">
                  <text class="order-no">{{ order.orderNo }}</text>
                  <text class="order-desc">{{ orderSummary(order) }}</text>
                  <view class="record-tags">
                    <text v-for="record in order.records" :key="record.id" class="record-tag">
                      {{ recordProcessName(record) }} · {{ recordUserName(record) }}
                    </text>
                  </view>
                </view>
                <view class="statement-side">
                  <text class="order-amount">{{ formatMoney(order.totalWage) }}</text>
                  <text class="order-time">{{ formatOrderTime(order) }}</text>
                </view>
              </view>
            </view>
          </view>
        </view>
      </view>
      <text v-else class="empty">暂无扫码记录</text>
    </view>

    <view v-if="isAdmin && activeTab === 'status'" class="panel">
      <text class="panel-title">订单流程状态</text>
      <view v-if="orders.length">
        <view v-for="order in orders" :key="order.id" class="order-card">
          <view class="order-head">
            <view>
              <text class="record-title">{{ order.order_no }}</text>
              <text class="record-meta">元件数 {{ (order.items || []).length }}</text>
            </view>
            <text class="status-tag">{{ order.status || 'draft' }}</text>
          </view>
          <view class="process-list">
            <view v-for="process in sortedProcesses(order.processes)" :key="process.id" class="process-item">
              <view>
                <text class="process-name">{{ processItemName(process) }} / {{ processName(process) }}</text>
                <text class="process-meta">{{ formatCapabilityListLabel(process.station_role) }}</text>
              </view>
              <text :class="hasScan(process.id) ? 'done' : 'pending'">
                {{ hasScan(process.id) ? '已扫码' : '待扫码' }}
              </text>
            </view>
          </view>
        </view>
      </view>
      <text v-else class="empty">暂无订单</text>
    </view>

    <view v-if="isAdmin && activeTab === 'order'" class="panel">
      <text class="panel-title">订单记录</text>
      <view v-if="orders.length">
        <view v-for="order in orders" :key="order.id" class="order-card">
          <view class="order-head">
            <view>
              <text class="record-title">{{ order.order_no }}</text>
              <text class="record-meta">{{ formatTime(order.created_at) }}</text>
            </view>
            <text class="status-tag">{{ order.status || 'draft' }}</text>
          </view>
          <text class="record-line">二维码：{{ order.qr_token }}</text>
          <text class="record-line">总数量：{{ order.quantity }}</text>
          <text class="record-line">订单总工资：{{ formatMoney(orderTotalWage(order)) }}</text>
          <text class="record-line">元件数：{{ (order.items || []).length }}</text>
          <view v-if="order.items && order.items.length" class="order-items">
            <view v-for="item in order.items" :key="item.id" class="order-item">
              <text class="order-item-name">{{ item.part_name }}</text>
              <text class="order-item-spec">{{ item.spec }}</text>
              <text class="order-item-qty">x{{ item.quantity }}</text>
            </view>
          </view>
          <view v-if="order.processes && order.processes.length" class="order-process-tags">
            <text v-for="process in sortedProcesses(order.processes)" :key="process.id" class="tag">
              {{ processItemName(process) }} / {{ processName(process) }} / {{ formatCapabilityListLabel(process.station_role) }}
            </text>
          </view>
        </view>
      </view>
      <text v-else class="empty">暂无订单记录</text>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { getMe } from '../../api/auth'
import { listOrders } from '../../api/order'
import { listAllScanRecords, listMyScanRecords } from '../../api/scan'
import { getUser } from '../../utils/session'
import { redirectTo } from '../../utils/navigation'
import { getRecordTabs } from '../../utils/tab-config'
import { formatCapabilityListLabel } from '../../utils/capability-labels'

const user = ref(getUser())
const scanRecords = ref([])
const orders = ref([])
const activeTab = ref('scan')

const isAdmin = computed(() => user.value?.role === 'admin')
const recordTabs = computed(() => getRecordTabs(user.value?.role))
const currentDayKey = formatDateKey(new Date())
const currentMonthKey = formatMonthKey(new Date())
const orderWageMap = computed(() => buildOrderWageMap(scanRecords.value))
const totalWage = computed(() => sumWage(scanRecords.value))
const todayWage = computed(() => sumWage(scanRecords.value.filter((item) => getDayKey(item.scanned_at) === currentDayKey)))
const monthWage = computed(() => sumWage(scanRecords.value.filter((item) => getMonthKey(item.scanned_at) === currentMonthKey)))
const wageMonths = computed(() => buildWageMonths(scanRecords.value))

onShow(() => {
  if (!user.value) {
    redirectTo('/pages/login/login')
    return
  }
  activeTab.value = recordTabs.value.some((item) => item.key === activeTab.value) ? activeTab.value : 'scan'
  loadPageData()
})

async function loadPageData() {
  const me = await getMe()
  user.value = me.user
  scanRecords.value = isAdmin.value
    ? (await listAllScanRecords()).items || []
    : (await listMyScanRecords()).items || []
  orders.value = isAdmin.value ? (await listOrders()).items || [] : []
}

function buildOrderWageMap(records) {
  const map = new Map()
  for (const item of records) {
    const key = orderKey(item)
    map.set(key, (map.get(key) || 0) + recordWage(item))
  }
  return map
}

function buildWageMonths(records) {
  const monthMap = new Map()
  const sortedRecords = [...records].sort((left, right) => toTime(right.scanned_at) - toTime(left.scanned_at))

  for (const record of sortedRecords) {
    const time = toDate(record.scanned_at)
    if (!time) continue

    const monthKey = formatMonthKey(time)
    const dayKey = formatDateKey(time)
    let month = monthMap.get(monthKey)
    if (!month) {
      month = {
        key: monthKey,
        label: formatMonthLabel(monthKey),
        totalWage: 0,
        daysMap: new Map(),
      }
      monthMap.set(monthKey, month)
    }

    let day = month.daysMap.get(dayKey)
    if (!day) {
      day = {
        key: dayKey,
        label: formatDayLabel(dayKey),
        weekday: formatWeekday(time),
        totalWage: 0,
        ordersMap: new Map(),
      }
      month.daysMap.set(dayKey, day)
    }

    const orderKeyValue = orderKey(record)
    let order = day.ordersMap.get(orderKeyValue)
    if (!order) {
      order = {
        key: orderKeyValue,
        orderNo: record.order?.order_no || `订单 ${record.order_id}`,
        totalWage: 0,
        scanCount: 0,
        records: [],
      }
      day.ordersMap.set(orderKeyValue, order)
    }

    const wage = recordWage(record)
    order.totalWage += wage
    order.scanCount += 1
    order.records.push(record)
    day.totalWage += wage
    month.totalWage += wage
  }

  return [...monthMap.values()]
    .map((month) => {
      const days = [...month.daysMap.values()]
        .map((day) => ({
          key: day.key,
          label: day.label,
          weekday: day.weekday,
          totalWage: day.totalWage,
          orders: [...day.ordersMap.values()].sort((left, right) => toTime(right.records[0].scanned_at) - toTime(left.records[0].scanned_at)),
        }))
        .sort((left, right) => right.key.localeCompare(left.key))
      return {
        key: month.key,
        label: month.label,
        totalWage: month.totalWage,
        orderCount: days.reduce((sum, day) => sum + day.orders.length, 0),
        days,
      }
    })
    .sort((left, right) => right.key.localeCompare(left.key))
}

function orderKey(item) {
  return item.order?.order_no || String(item.order_id)
}

function recordWage(item) {
  return Number(item.wage_amount || 0)
}

function sumWage(records) {
  return records.reduce((sum, item) => sum + recordWage(item), 0)
}

function toDate(value) {
  if (!value) return null
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? null : date
}

function toTime(value) {
  const date = toDate(value)
  return date ? date.getTime() : 0
}

function getDayKey(value) {
  const date = toDate(value)
  return date ? formatDateKey(date) : ''
}

function getMonthKey(value) {
  const date = toDate(value)
  return date ? formatMonthKey(date) : ''
}

function formatDateKey(value) {
  const date = value instanceof Date ? value : toDate(value)
  if (!date) return ''
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

function formatMonthKey(value) {
  const date = value instanceof Date ? value : toDate(value)
  if (!date) return ''
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  return `${year}-${month}`
}

function formatMonthLabel(key) {
  if (!key) return '-'
  const [year, month] = key.split('-')
  return `${year}年${Number(month)}月`
}

function formatDayLabel(key) {
  if (!key) return '-'
  const [, month, day] = key.split('-')
  return `${Number(month)}月${Number(day)}日`
}

function formatWeekday(date) {
  return ['周日', '周一', '周二', '周三', '周四', '周五', '周六'][date.getDay()]
}

function orderSummary(order) {
  const first = order.records[0]
  const item = first?.order_item || first?.order_process?.order_item
  const itemText = item ? `${item.part_name || '-'} ${item.spec || ''}`.trim() : '订单扫码'
  return `${itemText} · ${order.scanCount} 笔`
}

function formatOrderTime(order) {
  const first = order.records[0]
  if (!first?.scanned_at) return '-'
  const text = formatTime(first.scanned_at)
  return text.includes(' ') ? text.split(' ')[1] : text
}

function recordProcessName(item) {
  return item.process?.name || `工序 ${item.process_id}`
}

function recordUserName(item) {
  return item.user?.username || `用户 ${item.user_id}`
}

function sortedProcesses(processes = []) {
  return [...processes].sort((left, right) => (left.sort || 0) - (right.sort || 0))
}

function processName(item) {
  return item.process?.name || `工序 ${item.process_id}`
}

function processItemName(item) {
  return item.order_item?.part_name || `元件 ${item.order_item_id || '-'}`
}

function hasScan(orderProcessId) {
  return scanRecords.value.some((item) => item.order_process_id === orderProcessId)
}

function orderTotalWage(order) {
  return orderWageMap.value.get(order.order_no || String(order.id)) || 0
}

function formatMoney(value) {
  return `${Number(value || 0).toFixed(2)} 元`
}

function formatTime(value) {
  if (!value) return '-'
  return String(value).replace('T', ' ').replace(/\.\d+(Z)?$/, '').replace(/Z$/, '')
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

.refresh {
  height: 64rpx;
  line-height: 64rpx;
  padding: 0 24rpx;
  border-radius: 999rpx;
  background: rgba(255, 255, 255, 0.16);
  color: #ffffff;
  font-size: 24rpx;
}

.metrics {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12rpx;
  margin-top: 22rpx;
}

.metric {
  padding: 18rpx 16rpx;
  border-radius: 14rpx;
  background: #ffffff;
  border: 1rpx solid #e2e7ee;
}

.metric-label,
.amount-label,
.month-subtitle,
.day-week,
.order-desc,
.order-time,
.record-meta,
.record-line,
.empty {
  display: block;
  color: #7a8494;
  font-size: 24rpx;
}

.metric-value {
  display: block;
  margin-top: 8rpx;
  font-size: 30rpx;
  font-weight: 700;
  color: #18212f;
}

.sub-tabs {
  display: flex;
  gap: 12rpx;
  margin-top: 22rpx;
  padding: 8rpx;
  border-radius: 14rpx;
  background: #e9eef3;
}

.sub-tab {
  flex: 1;
  height: 68rpx;
  line-height: 68rpx;
  border-radius: 10rpx;
  background: transparent;
  color: #596575;
  font-size: 26rpx;
}

.sub-tab.active {
  background: #ffffff;
  color: #1f6f5b;
  font-weight: 700;
}

.statement-list {
  margin-top: 22rpx;
}

.month-card {
  margin-bottom: 22rpx;
  overflow: hidden;
  border-radius: 16rpx;
  background: #ffffff;
  border: 1rpx solid #e2e7ee;
}

.month-head {
  display: flex;
  justify-content: space-between;
  gap: 18rpx;
  align-items: center;
  padding: 24rpx;
  background: #f8fafc;
  border-bottom: 1rpx solid #e7edf3;
}

.month-title {
  display: block;
  font-size: 34rpx;
  font-weight: 800;
  color: #18212f;
}

.month-subtitle {
  margin-top: 8rpx;
}

.month-total {
  flex: 0 0 220rpx;
  text-align: right;
}

.amount-label {
  font-size: 22rpx;
}

.amount-value {
  display: block;
  margin-top: 8rpx;
  color: #1f6f5b;
  font-size: 34rpx;
  font-weight: 800;
}

.day-section {
  padding: 20rpx 24rpx 8rpx;
  border-bottom: 1rpx solid #eef2f5;
}

.day-section:last-child {
  border-bottom: 0;
}

.day-head,
.statement-row,
.order-head,
.process-item {
  display: flex;
  justify-content: space-between;
  gap: 16rpx;
  align-items: center;
}

.day-title {
  display: block;
  font-size: 28rpx;
  font-weight: 800;
  color: #18212f;
}

.day-week {
  margin-top: 6rpx;
}

.day-total {
  color: #1f6f5b;
  font-size: 28rpx;
  font-weight: 800;
}

.order-list {
  margin-top: 14rpx;
}

.statement-row {
  min-height: 112rpx;
  padding: 18rpx 0;
  border-top: 1rpx solid #f0f3f6;
}

.statement-main {
  min-width: 0;
  flex: 1;
}

.statement-side {
  flex: 0 0 180rpx;
  text-align: right;
}

.order-no {
  display: block;
  color: #18212f;
  font-size: 28rpx;
  font-weight: 800;
  word-break: break-all;
}

.order-desc {
  margin-top: 8rpx;
}

.order-amount {
  display: block;
  color: #18212f;
  font-size: 28rpx;
  font-weight: 800;
}

.order-time {
  margin-top: 8rpx;
}

.record-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8rpx;
  margin-top: 10rpx;
}

.record-tag,
.status-tag,
.tag {
  padding: 8rpx 14rpx;
  border-radius: 999rpx;
  background: #eff3f7;
  color: #516070;
  font-size: 22rpx;
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

.order-card {
  padding: 18rpx 0;
  border-bottom: 1rpx solid #eef2f5;
}

.order-card:last-child {
  border-bottom: 0;
}

.record-title {
  display: block;
  font-size: 28rpx;
  font-weight: 700;
  color: #18212f;
}

.record-meta,
.record-line {
  margin-top: 8rpx;
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

.process-name,
.process-meta {
  display: block;
}

.process-name {
  color: #18212f;
  font-weight: 700;
}

.process-meta {
  margin-top: 4rpx;
  color: #7a8494;
}

.done {
  color: #1f6f5b;
  font-weight: 700;
}

.pending {
  color: #a45b18;
  font-weight: 700;
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

.order-item-name {
  font-weight: 700;
  color: #18212f;
}

.order-item-spec,
.order-item-qty {
  color: #5f6a7a;
}

.order-process-tags {
  margin-top: 14rpx;
  display: flex;
  flex-wrap: wrap;
  gap: 10rpx;
}
</style>
