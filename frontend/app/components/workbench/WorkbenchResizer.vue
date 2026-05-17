<script setup lang="ts">
import { useEventListener } from '@vueuse/core'

const props = defineProps<{
  orientation: 'horizontal' | 'vertical'
  modelValue: number
  min?: number
  max?: number
  reverse?: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: number]
}>()

const isDragging = ref(false)
const startPoint = ref(0)
const startValue = ref(0)
const trackSize = ref(1)

const startDrag = (e: MouseEvent) => {
  e.preventDefault()
  isDragging.value = true
  startPoint.value = props.orientation === 'horizontal' ? e.clientX : e.clientY
  startValue.value = props.modelValue
  trackSize.value = props.orientation === 'vertical'
    ? ((e.currentTarget as HTMLElement).parentElement?.getBoundingClientRect().height || window.innerHeight)
    : 1
  document.body.style.cursor = props.orientation === 'horizontal' ? 'col-resize' : 'row-resize'
  document.body.style.userSelect = 'none'
}

useEventListener('mousemove', (e: MouseEvent) => {
  if (!isDragging.value) return

  if (props.orientation === 'horizontal') {
    const delta = e.clientX - startPoint.value
    let newWidth = props.reverse ? startValue.value - delta : startValue.value + delta
    if (props.min && newWidth < props.min) newWidth = props.min
    if (props.max && newWidth > props.max) newWidth = props.max
    emit('update:modelValue', newWidth)
  } else {
    const delta = e.clientY - startPoint.value
    let newHeightPercent = startValue.value + (delta / trackSize.value) * 100
    if (props.min && newHeightPercent < props.min) newHeightPercent = props.min
    if (props.max && newHeightPercent > props.max) newHeightPercent = props.max
    emit('update:modelValue', newHeightPercent)
  }
})

useEventListener('mouseup', () => {
  if (isDragging.value) {
    isDragging.value = false
    document.body.style.cursor = ''
    document.body.style.userSelect = ''
  }
})
</script>

<template>
  <div
    class="relative flex items-center justify-center transition-colors hover:bg-primary/30"
    :class="[
      orientation === 'horizontal' ? 'w-1.5 h-full cursor-col-resize' : 'h-1.5 w-full cursor-row-resize',
      isDragging ? 'bg-primary/50' : 'bg-transparent'
    ]"
    @mousedown="startDrag"
  >
    <!-- Tactical Handle - centered within drag area -->
    <div 
      class="absolute bg-primary/40 transition-all group-hover:bg-primary/80"
      :class="[
        orientation === 'horizontal' 
          ? 'w-0.5 h-10 rounded-full' 
          : 'h-0.5 w-10 rounded-full',
        isDragging ? 'bg-primary scale-150' : ''
      ]"
    />
  </div>
</template>
