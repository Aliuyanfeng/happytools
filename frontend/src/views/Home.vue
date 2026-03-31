<template>
  <div class="home">
    <!-- 背景光晕 -->
    <div class="aurora">
      <div class="aurora-blob a1" />
      <div class="aurora-blob a2" />
    </div>

    <!-- 左侧：展示区 -->
    <div class="panel-left">
      <!-- 图片 -->
      <div class="img-wrap">
        <img src="@/assets/images/leftimage-blue.png" alt="banner" class="banner-img" />
        <div class="img-glow" />
      </div>

      <!-- 文字说明 -->
      <div class="desc-block">
        <h2 class="desc-title">{{ t('home.welcome') }}<br /><span class="desc-accent">{{ t('home.subtitle') }}</span></h2>
        <div class="feature-tags">
          <span v-for="f in features" :key="f.icon" class="ftag">
            <span class="ftag-icon">{{ f.icon }}</span>{{ f.label }}
          </span>
        </div>
      </div>

      <!-- 时钟 -->
      <div class="clock-row">
        <span class="clock-dot" />
        <span class="clock-time">{{ currentTime }}</span>
        <span class="clock-greet">· {{ greetingText }}</span>
      </div>
    </div>

    <!-- 右侧：功能卡片 -->
    <div class="panel-right">
      <div class="right-header">
        <span class="right-title">{{ t('home.featureEntry') }}</span>
        <span class="right-count">{{ visibleModules.length }} 个工具</span>
      </div>
      <div class="card-grid">
        <div
          v-for="(module, i) in visibleModules"
          :key="module.id"
          class="card"
          :class="`c-${module.theme}`"
          :style="{ animationDelay: `${i * 50}ms` }"
          @mouseenter="onEnter"
          @mouseleave="onLeave"
          @click="go(module.path)"
        >
          <div class="card-spot" />
          <!-- 右侧大号装饰图标 -->
          <div class="card-deco">
            <component :is="getIconComponent(module.icon)" />
          </div>
          <!-- 左侧内容 -->
          <div class="card-icon-wrap">
            <component :is="getIconComponent(module.icon)" />
          </div>
          <span class="card-name">{{ t(module.nameKey) }}</span>
          <ArrowRightOutlined class="card-arrow" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  DashboardOutlined, CheckCircleOutlined, ToolOutlined,
  SafetyOutlined, ApartmentOutlined, CalendarOutlined,
  BranchesOutlined, FileTextOutlined, ArrowRightOutlined,
} from '@ant-design/icons-vue'
import { modules } from '@/config/modules'
import { useSettingsStore } from '@/stores/settings'

const { t } = useI18n()
const router = useRouter()
const settingsStore = useSettingsStore()

const visibleModules = computed(() =>
  modules.filter(m => !settingsStore.hiddenModules.includes(m.id))
)

const iconMap: Record<string, any> = {
  DashboardOutlined, CheckCircleOutlined, ToolOutlined,
  SafetyOutlined, ApartmentOutlined, CalendarOutlined,
  BranchesOutlined, FileTextOutlined,
}
function getIconComponent(n: string) { return iconMap[n] || DashboardOutlined }
function go(path: string) { router.push(path) }

const features = [
  { icon: '⚡', label: t('home.lightFast') },
  { icon: '🔒', label: t('home.secure') },
  { icon: '🎯', label: t('home.simple') },
]

// 3D 倾斜 + 鼠标光斑
function onEnter(e: MouseEvent) {
  const card = e.currentTarget as HTMLElement
  const spot = card.querySelector('.card-spot') as HTMLElement

  function move(ev: MouseEvent) {
    const r = card.getBoundingClientRect()
    const x = ev.clientX - r.left
    const y = ev.clientY - r.top
    spot.style.setProperty('--sx', `${x}px`)
    spot.style.setProperty('--sy', `${y}px`)
    spot.style.opacity = '1'
    const rx = ((y / r.height) - 0.5) * -10
    const ry = ((x / r.width)  - 0.5) *  10
    card.style.transform = `translateY(-6px) scale(1.06) rotateX(${rx}deg) rotateY(${ry}deg)`
  }
  card.addEventListener('mousemove', move)
  ;(card as any)._move = move
}

function onLeave(e: MouseEvent) {
  const card = e.currentTarget as HTMLElement
  const spot = card.querySelector('.card-spot') as HTMLElement
  card.removeEventListener('mousemove', (card as any)._move)
  spot.style.opacity = '0'
  card.style.transform = ''
}

// 时钟
const now = ref(new Date())
let timer: ReturnType<typeof setInterval>
onMounted(() => { timer = setInterval(() => { now.value = new Date() }, 1000) })
onUnmounted(() => clearInterval(timer))

const pad = (n: number) => n.toString().padStart(2, '0')
const currentTime = computed(() =>
  `${pad(now.value.getHours())}:${pad(now.value.getMinutes())}:${pad(now.value.getSeconds())}`
)
const greetingText = computed(() => {
  const h = now.value.getHours()
  if (h < 6)  return 'Good Night'
  if (h < 12) return 'Good Morning'
  if (h < 18) return 'Good Afternoon'
  return 'Good Evening'
})
</script>

<style scoped>
/* ── 根容器 ── */
.home {
  position: relative;
  width: 100%;
  height: 100%;
  background: #f0f4ff;
  display: flex;
  overflow: hidden;
}

/* ── 背景光晕 ── */
.aurora { position: absolute; inset: 0; pointer-events: none; z-index: 0; }
.aurora-blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(90px);
  animation: blob 20s ease-in-out infinite alternate;
}
.a1 {
  width: 50vw; height: 50vw;
  background: radial-gradient(circle, rgba(139,92,246,0.15), transparent 70%);
  top: -20%; left: -10%;
}
.a2 {
  width: 40vw; height: 40vw;
  background: radial-gradient(circle, rgba(6,182,212,0.12), transparent 70%);
  bottom: -15%; right: 30%;
  animation-delay: -10s;
}
@keyframes blob {
  0%   { transform: translate(0,0) scale(1); }
  50%  { transform: translate(4%,6%) scale(1.08); }
  100% { transform: translate(-3%,3%) scale(0.95); }
}

/* ══════════════════════════════════════
   左侧展示区
══════════════════════════════════════ */
.panel-left {
  position: relative;
  z-index: 1;
  width: 42%;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 20px;
  padding: 24px 28px;
  border-right: 1px solid rgba(99,102,241,0.1);
  background: rgba(255,255,255,0.45);
  backdrop-filter: blur(12px);
}

/* 图片容器 */
.img-wrap {
  position: relative;
  width: 100%;
  max-width: 340px;
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 12px 48px rgba(99,102,241,0.18);
  animation: float 5s ease-in-out infinite;
}
@keyframes float {
  0%,100% { transform: translateY(0); }
  50%      { transform: translateY(-8px); }
}

.banner-img {
  width: 100%;
  height: auto;
  display: block;
  border-radius: 20px;
}

.img-glow {
  position: absolute;
  inset: 0;
  border-radius: 20px;
  background: linear-gradient(135deg, rgba(99,102,241,0.08) 0%, transparent 60%);
  pointer-events: none;
}

/* 文字说明 */
.desc-block {
  width: 100%;
  max-width: 340px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  align-items: center;
  text-align: center;
}

.desc-title {
  font-size: clamp(20px, 2.2vw, 28px);
  font-weight: 700;
  color: #1e1b4b;
  margin: 0;
  line-height: 1.4;
  animation: title-in 0.6s 0.1s ease both;
}
@keyframes title-in {
  from { opacity: 0; transform: translateY(8px); }
  to   { opacity: 1; transform: none; }
}

.desc-accent {
  background: linear-gradient(135deg, #6366f1, #ec4899);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.feature-tags {
  display: flex;
  gap: 8px;
  flex-wrap: nowrap;
  justify-content: center;
  animation: title-in 0.6s 0.2s ease both;
}

.ftag {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  border-radius: 10px;
  background: rgba(255,255,255,0.7);
  border: 1px solid rgba(99,102,241,0.12);
  font-size: 11px;
  color: #475569;
  white-space: nowrap;
  font-weight: 500;
}
.ftag-icon { font-size: 13px; }

/* 时钟行 */
.clock-row {
  display: flex;
  align-items: center;
  gap: 8px;
  animation: title-in 0.6s 0.3s ease both;
}
.clock-dot {
  width: 6px; height: 6px;
  border-radius: 50%;
  background: #6366f1;
  animation: ping 2s ease-in-out infinite;
  box-shadow: 0 0 0 0 rgba(99,102,241,0.5);
}
@keyframes ping {
  0%,100% { box-shadow: 0 0 0 0 rgba(99,102,241,0.5); }
  50%      { box-shadow: 0 0 0 6px rgba(99,102,241,0); }
}
.clock-time {
  font-size: 13px;
  font-weight: 700;
  color: #1e1b4b;
  font-variant-numeric: tabular-nums;
  letter-spacing: 2px;
}
.clock-greet { font-size: 12px; color: #94a3b8; white-space: nowrap; }

/* ══════════════════════════════════════
   右侧功能卡片
══════════════════════════════════════ */
.panel-right {
  position: relative;
  z-index: 1;
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  padding: 20px 20px 16px;
  gap: 12px;
  overflow: hidden;
}

.right-header {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  flex-shrink: 0;
}
.right-title { font-size: 15px; font-weight: 700; color: #1e1b4b; white-space: nowrap; }
.right-count { font-size: 11px; color: #cbd5e1; white-space: nowrap; }

/* 等高均匀网格 */
.card-grid {
  flex: 1;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  grid-auto-rows: 1fr;
  gap: 10px;
  overflow: hidden;
}

/* 卡片入场 */
@keyframes card-in {
  from { opacity: 0; transform: translateY(14px) scale(0.96); }
  to   { opacity: 1; transform: none; }
}

/* 卡片基础：左对齐，图标左上，文字左下 */
.card {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: space-between;
  padding: 14px 14px 12px;
  border-radius: 16px;
  cursor: pointer;
  overflow: hidden;
  transform-style: preserve-3d;
  transition: transform 0.28s cubic-bezier(0.23,1,0.32,1), box-shadow 0.28s ease;
  animation: card-in 0.4s cubic-bezier(0.23,1,0.32,1) both;
  border: 1px solid rgba(255,255,255,0.75);
  min-height: 0;
}
.card:hover { transform: translateY(-4px) scale(1.02); }

/* 右侧大号半透明装饰图标 */
.card-deco {
  position: absolute;
  right: -6px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 64px;
  line-height: 1;
  opacity: 0.08;
  pointer-events: none;
  transition: opacity 0.28s, transform 0.28s cubic-bezier(0.34,1.56,0.64,1);
}
.card:hover .card-deco {
  opacity: 0.15;
  transform: translateY(-50%) scale(1.12) rotate(-8deg);
}

/* 左上：图标 */
.card-icon-wrap {
  width: 38px; height: 38px;
  border-radius: 11px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  flex-shrink: 0;
  transition: transform 0.28s cubic-bezier(0.34,1.56,0.64,1);
  position: relative;
  z-index: 1;
}
.card:hover .card-icon-wrap { transform: scale(1.12) rotate(5deg); }

/* 左下：文字 */
.card-name {
  font-size: 18px;
  font-weight: 700;
  white-space: nowrap;
  letter-spacing: 0.2px;
  position: relative;
  z-index: 1;
  transition: transform 0.2s ease;
}
.card:hover .card-name { transform: translateX(2px); }

/* 右下角箭头（hover 出现） */
.card-arrow {
  position: absolute;
  bottom: 10px;
  right: 12px;
  font-size: 14px;
  opacity: 0;
  transform: translate(-3px, 3px);
  transition: opacity 0.2s, transform 0.2s;
  z-index: 1;
}
.card:hover .card-arrow { opacity: 1; transform: translate(0, 0); }

/* 鼠标光斑 */
.card-spot {
  position: absolute;
  width: 120px; height: 120px;
  border-radius: 50%;
  transform: translate(calc(var(--sx,50%) - 60px), calc(var(--sy,50%) - 60px));
  pointer-events: none;
  opacity: 0;
  transition: opacity 0.2s;
}

/* ── 主题色 ── */
.c-blue   { background: linear-gradient(145deg,#eef2ff,#e0e7ff); box-shadow:0 2px 10px rgba(99,102,241,.1); }
.c-green  { background: linear-gradient(145deg,#ecfdf5,#d1fae5); box-shadow:0 2px 10px rgba(16,185,129,.1); }
.c-purple { background: linear-gradient(145deg,#faf5ff,#f3e8ff); box-shadow:0 2px 10px rgba(168,85,247,.1); }
.c-orange { background: linear-gradient(145deg,#fff7ed,#ffedd5); box-shadow:0 2px 10px rgba(249,115,22,.1); }
.c-red    { background: linear-gradient(145deg,#fff1f2,#ffe4e6); box-shadow:0 2px 10px rgba(239,68,68,.1); }
.c-cyan   { background: linear-gradient(145deg,#ecfeff,#cffafe); box-shadow:0 2px 10px rgba(6,182,212,.1); }

.c-blue:hover   { box-shadow:0 10px 30px rgba(99,102,241,.25); }
.c-green:hover  { box-shadow:0 10px 30px rgba(16,185,129,.25); }
.c-purple:hover { box-shadow:0 10px 30px rgba(168,85,247,.25); }
.c-orange:hover { box-shadow:0 10px 30px rgba(249,115,22,.25); }
.c-red:hover    { box-shadow:0 10px 30px rgba(239,68,68,.25); }
.c-cyan:hover   { box-shadow:0 10px 30px rgba(6,182,212,.25); }

.c-blue   .card-icon-wrap { background:rgba(99,102,241,.14); color:#4f46e5; }
.c-green  .card-icon-wrap { background:rgba(16,185,129,.14); color:#059669; }
.c-purple .card-icon-wrap { background:rgba(168,85,247,.14); color:#9333ea; }
.c-orange .card-icon-wrap { background:rgba(249,115,22,.14); color:#ea580c; }
.c-red    .card-icon-wrap { background:rgba(239,68,68,.14);  color:#dc2626; }
.c-cyan   .card-icon-wrap { background:rgba(6,182,212,.14);  color:#0891b2; }

.c-blue   .card-name  { color:#3730a3; } .c-blue   .card-arrow { color:#4f46e5; } .c-blue   .card-deco { color:#4f46e5; }
.c-green  .card-name  { color:#065f46; } .c-green  .card-arrow { color:#059669; } .c-green  .card-deco { color:#059669; }
.c-purple .card-name  { color:#6b21a8; } .c-purple .card-arrow { color:#9333ea; } .c-purple .card-deco { color:#9333ea; }
.c-orange .card-name  { color:#9a3412; } .c-orange .card-arrow { color:#ea580c; } .c-orange .card-deco { color:#ea580c; }
.c-red    .card-name  { color:#991b1b; } .c-red    .card-arrow { color:#dc2626; } .c-red    .card-deco { color:#dc2626; }
.c-cyan   .card-name  { color:#155e75; } .c-cyan   .card-arrow { color:#0891b2; } .c-cyan   .card-deco { color:#0891b2; }

.c-blue   .card-spot { background:radial-gradient(circle,rgba(99,102,241,.18),transparent 70%); }
.c-green  .card-spot { background:radial-gradient(circle,rgba(16,185,129,.18),transparent 70%); }
.c-purple .card-spot { background:radial-gradient(circle,rgba(168,85,247,.18),transparent 70%); }
.c-orange .card-spot { background:radial-gradient(circle,rgba(249,115,22,.18),transparent 70%); }
.c-red    .card-spot { background:radial-gradient(circle,rgba(239,68,68,.18), transparent 70%); }
.c-cyan   .card-spot { background:radial-gradient(circle,rgba(6,182,212,.18), transparent 70%); }
</style>
