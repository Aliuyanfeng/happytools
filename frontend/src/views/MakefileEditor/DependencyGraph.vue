<template>
  <div class="dependency-graph" ref="containerRef">
    <!-- Empty state -->
    <div v-if="!store.currentDoc || nodes.length === 0" class="empty-state">
      <ApartmentOutlined class="empty-icon" />
      <span>{{ t('makefileEditor.noTargets') }}</span>
    </div>

    <!-- SVG graph -->
    <svg
      v-else
      ref="svgRef"
      class="graph-svg"
      :viewBox="viewBox"
      @mousedown="onSvgMouseDown"
      @mousemove="onSvgMouseMove"
      @mouseup="onSvgMouseUp"
      @mouseleave="onSvgMouseUp"
      @wheel.prevent="onWheel"
    >
      <defs>
        <!-- Arrow marker for normal edges -->
        <marker
          id="arrow-normal"
          markerWidth="8"
          markerHeight="8"
          refX="7"
          refY="3"
          orient="auto"
        >
          <path d="M0,0 L0,6 L8,3 z" fill="#adb5bd" />
        </marker>
        <!-- Arrow marker for cycle edges -->
        <marker
          id="arrow-cycle"
          markerWidth="8"
          markerHeight="8"
          refX="7"
          refY="3"
          orient="auto"
        >
          <path d="M0,0 L0,6 L8,3 z" fill="#ff4d4f" />
        </marker>
      </defs>

      <g :transform="`translate(${pan.x}, ${pan.y}) scale(${zoom})`">
        <!-- Edges -->
        <g class="edges">
          <g v-for="edge in edges" :key="`${edge.from}-${edge.to}`">
            <line
              :x1="edge.x1"
              :y1="edge.y1"
              :x2="edge.x2"
              :y2="edge.y2"
              :stroke="edge.isCycle ? '#ff4d4f' : '#adb5bd'"
              stroke-width="1.5"
              :marker-end="edge.isCycle ? 'url(#arrow-cycle)' : 'url(#arrow-normal)'"
            />
          </g>
        </g>

        <!-- Nodes -->
        <g class="nodes">
          <g
            v-for="node in nodes"
            :key="node.name"
            :transform="`translate(${node.x}, ${node.y})`"
            class="node-group"
            @click.stop="handleNodeClick(node.name)"
            style="cursor: pointer"
          >
            <!-- Node background rect -->
            <rect
              :width="NODE_W"
              :height="NODE_H"
              :rx="6"
              :ry="6"
              :fill="nodeFill(node)"
              :stroke="nodeStroke(node)"
              stroke-width="2"
            />

            <!-- Phony badge dot -->
            <circle
              v-if="node.isPhony"
              :cx="NODE_W - 10"
              :cy="10"
              r="4"
              fill="#1677ff"
            />

            <!-- Cycle warning icon (⚠ -->
            <text
              v-if="node.isCycle"
              :x="10"
              :y="NODE_H - 8"
              font-size="11"
              fill="#ff4d4f"
            >⚠</text>

            <!-- Target name label -->
            <text
              :x="NODE_W / 2"
              :y="NODE_H / 2 + 4"
              text-anchor="middle"
              font-size="12"
              font-family="monospace"
              :fill="node.isCycle ? '#ff4d4f' : '#262626'"
            >{{ truncate(node.name, 16) }}</text>

            <!-- Tooltip title for cycle nodes -->
            <title v-if="node.isCycle">{{ cycleTooltip(node.name) }}</title>
            <!-- Tooltip for all nodes showing full name -->
            <title v-else>{{ node.name }}</title>
          </g>
        </g>
      </g>
    </svg>

    <!-- Legend -->
    <div v-if="store.currentDoc && nodes.length > 0" class="legend">
      <span class="legend-item">
        <span class="legend-dot phony" />
        {{ t('makefileEditor.phonyTarget') }}
      </span>
      <span class="legend-item">
        <span class="legend-dot regular" />
        {{ t('makefileEditor.regularTarget') }}
      </span>
      <span v-if="hasCycles" class="legend-item">
        <span class="legend-dot cycle" />
        {{ t('makefileEditor.cycleNode') }}
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ApartmentOutlined } from '@ant-design/icons-vue'
import { useMakefileEditorStore } from '../../stores/makefileEditor'

const { t } = useI18n()
const store = useMakefileEditorStore()

// ── Emits ──────────────────────────────────────────────────────────────────
const emit = defineEmits<{
  (e: 'select-target', name: string): void
}>()

// ── Constants ──────────────────────────────────────────────────────────────
const NODE_W = 140
const NODE_H = 40
const H_GAP = 60   // horizontal gap between nodes in same layer
const V_GAP = 80   // vertical gap between layers

// ── Refs ───────────────────────────────────────────────────────────────────
const containerRef = ref<HTMLDivElement | null>(null)
const svgRef = ref<SVGSVGElement | null>(null)

// Pan & zoom state
const pan = ref({ x: 20, y: 20 })
const zoom = ref(1)
const isPanning = ref(false)
const panStart = ref({ x: 0, y: 0 })
const panOrigin = ref({ x: 0, y: 0 })

// ── Computed: cycle set ────────────────────────────────────────────────────
const cycleNodeSet = computed<Set<string>>(() => {
  const s = new Set<string>()
  for (const cycle of store.cycleWarnings) {
    for (const name of cycle) s.add(name)
  }
  return s
})

const cycleEdgeSet = computed<Set<string>>(() => {
  const s = new Set<string>()
  for (const cycle of store.cycleWarnings) {
    for (let i = 0; i < cycle.length; i++) {
      const from = cycle[i]
      const to = cycle[(i + 1) % cycle.length]
      s.add(`${from}->${to}`)
    }
  }
  return s
})

const hasCycles = computed(() => cycleNodeSet.value.size > 0)

// ── Layout algorithm: topological layering ─────────────────────────────────
interface LayoutNode {
  name: string
  isPhony: boolean
  isCycle: boolean
  x: number
  y: number
}

interface LayoutEdge {
  from: string
  to: string
  x1: number
  y1: number
  x2: number
  y2: number
  isCycle: boolean
}

const nodes = computed<LayoutNode[]>(() => {
  const targets = store.currentDoc?.targets ?? []
  if (targets.length === 0) return []

  // Build adjacency: target -> deps
  const nameSet = new Set(targets.map(t => t.name))
  const inDegree = new Map<string, number>()
  const adjList = new Map<string, string[]>()

  for (const t of targets) {
    if (!inDegree.has(t.name)) inDegree.set(t.name, 0)
    adjList.set(t.name, [])
  }

  for (const t of targets) {
    for (const dep of (t.deps ?? [])) {
      if (nameSet.has(dep)) {
        adjList.get(t.name)!.push(dep)
        inDegree.set(dep, (inDegree.get(dep) ?? 0) + 1)
      }
    }
  }

  // Kahn's algorithm for layer assignment
  const layer = new Map<string, number>()
  const queue: string[] = []

  for (const [name, deg] of inDegree) {
    if (deg === 0) queue.push(name)
  }

  // BFS to assign layers
  const tempDeg = new Map(inDegree)
  let head = 0
  while (head < queue.length) {
    const cur = queue[head++]
    const curLayer = layer.get(cur) ?? 0
    for (const dep of (adjList.get(cur) ?? [])) {
      const newLayer = curLayer + 1
      if (!layer.has(dep) || layer.get(dep)! < newLayer) {
        layer.set(dep, newLayer)
      }
      tempDeg.set(dep, (tempDeg.get(dep) ?? 1) - 1)
      if (tempDeg.get(dep) === 0) queue.push(dep)
    }
  }

  // Nodes not reached (in cycles) get layer 0
  for (const t of targets) {
    if (!layer.has(t.name)) layer.set(t.name, 0)
  }

  // Group by layer
  const layerGroups = new Map<number, string[]>()
  for (const [name, l] of layer) {
    if (!layerGroups.has(l)) layerGroups.set(l, [])
    layerGroups.get(l)!.push(name)
  }

  // Position nodes
  const posMap = new Map<string, { x: number; y: number }>()
  const sortedLayers = [...layerGroups.keys()].sort((a, b) => a - b)

  for (const l of sortedLayers) {
    const group = layerGroups.get(l)!
    const totalW = group.length * NODE_W + (group.length - 1) * H_GAP
    const startX = 0
    group.forEach((name, i) => {
      posMap.set(name, {
        x: startX + i * (NODE_W + H_GAP),
        y: l * (NODE_H + V_GAP),
      })
    })
  }

  return targets.map(t => ({
    name: t.name,
    isPhony: t.isPhony,
    isCycle: cycleNodeSet.value.has(t.name),
    x: posMap.get(t.name)?.x ?? 0,
    y: posMap.get(t.name)?.y ?? 0,
  }))
})

const edges = computed<LayoutEdge[]>(() => {
  const targets = store.currentDoc?.targets ?? []
  const posMap = new Map<string, { x: number; y: number }>()
  for (const n of nodes.value) posMap.set(n.name, { x: n.x, y: n.y })

  const nameSet = new Set(targets.map(t => t.name))
  const result: LayoutEdge[] = []

  for (const t of targets) {
    const fromPos = posMap.get(t.name)
    if (!fromPos) continue
    for (const dep of (t.deps ?? [])) {
      if (!nameSet.has(dep)) continue
      const toPos = posMap.get(dep)
      if (!toPos) continue

      const isCycle = cycleEdgeSet.value.has(`${t.name}->${dep}`)

      // Connect bottom-center of source to top-center of target
      result.push({
        from: t.name,
        to: dep,
        x1: fromPos.x + NODE_W / 2,
        y1: fromPos.y + NODE_H,
        x2: toPos.x + NODE_W / 2,
        y2: toPos.y,
        isCycle,
      })
    }
  }
  return result
})

// ── ViewBox: auto-fit to content ───────────────────────────────────────────
const viewBox = computed(() => {
  if (nodes.value.length === 0) return '0 0 400 300'
  const maxX = Math.max(...nodes.value.map(n => n.x + NODE_W)) + 40
  const maxY = Math.max(...nodes.value.map(n => n.y + NODE_H)) + 40
  return `0 0 ${maxX} ${maxY}`
})

// ── Node styling helpers ───────────────────────────────────────────────────
function nodeFill(node: LayoutNode): string {
  if (node.isCycle) return '#fff1f0'
  if (node.isPhony) return '#e6f4ff'
  return '#fafafa'
}

function nodeStroke(node: LayoutNode): string {
  if (node.isCycle) return '#ff4d4f'
  if (node.isPhony) return '#1677ff'
  return '#d9d9d9'
}

function truncate(s: string, max: number): string {
  return s.length > max ? s.slice(0, max - 1) + '...' : s
}

function cycleTooltip(name: string): string {
  const cycles = store.cycleWarnings.filter(c => c.includes(name))
  if (cycles.length === 0) return name
  return `${t('makefileEditor.cycleWarning')}: ${cycles.map(c => c.join(' -> ')).join('; ')}`
}

// ── Pan & zoom handlers ────────────────────────────────────────────────────
function onSvgMouseDown(e: MouseEvent) {
  if (e.button !== 0) return
  isPanning.value = true
  panStart.value = { x: e.clientX, y: e.clientY }
  panOrigin.value = { x: pan.value.x, y: pan.value.y }
}

function onSvgMouseMove(e: MouseEvent) {
  if (!isPanning.value) return
  pan.value = {
    x: panOrigin.value.x + (e.clientX - panStart.value.x),
    y: panOrigin.value.y + (e.clientY - panStart.value.y),
  }
}

function onSvgMouseUp() {
  isPanning.value = false
}

function onWheel(e: WheelEvent) {
  const delta = e.deltaY > 0 ? 0.9 : 1.1
  zoom.value = Math.min(3, Math.max(0.2, zoom.value * delta))
}

// ── Node click ─────────────────────────────────────────────────────────────
function handleNodeClick(name: string) {
  store.setSelectedTarget(name)
  emit('select-target', name)
}

// ── Reset pan/zoom when doc changes ───────────────────────────────────────
watch(
  () => store.currentDoc,
  () => {
    pan.value = { x: 20, y: 20 }
    zoom.value = 1
  },
)

// ── Keyboard shortcut: reset view with Escape ──────────────────────────────
function onKeyDown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    pan.value = { x: 20, y: 20 }
    zoom.value = 1
  }
}

onMounted(() => window.addEventListener('keydown', onKeyDown))
onUnmounted(() => window.removeEventListener('keydown', onKeyDown))
</script>

<style scoped>
.dependency-graph {
  position: relative;
  width: 100%;
  height: 100%;
  min-height: 300px;
  background: #f8f9fa;
  border-radius: 8px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.graph-svg {
  flex: 1;
  width: 100%;
  height: 100%;
  cursor: grab;
  user-select: none;
}

.graph-svg:active {
  cursor: grabbing;
}

.node-group {
  transition: opacity 0.15s;
}

.node-group:hover rect {
  filter: brightness(0.95);
}

.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10px;
  color: #bfbfbf;
  font-size: 13px;
}

.empty-icon {
  font-size: 36px;
  color: #d9d9d9;
}

.legend {
  display: flex;
  gap: 16px;
  padding: 6px 12px;
  background: rgba(255, 255, 255, 0.85);
  border-top: 1px solid #f0f0f0;
  font-size: 11px;
  color: #595959;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 5px;
}

.legend-dot {
  display: inline-block;
  width: 12px;
  height: 12px;
  border-radius: 3px;
  border: 2px solid;
}

.legend-dot.phony {
  background: #e6f4ff;
  border-color: #1677ff;
}

.legend-dot.regular {
  background: #fafafa;
  border-color: #d9d9d9;
}

.legend-dot.cycle {
  background: #fff1f0;
  border-color: #ff4d4f;
}
</style>
