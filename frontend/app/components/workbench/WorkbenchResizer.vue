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
    class="group relative flex items-center justify-center transition-all duration-300 select-none z-20"
    :class="[
      orientation === 'horizontal' ? 'w-1.5 h-full cursor-col-resize hover:bg-primary/10' : 'h-1.5 w-full cursor-row-resize hover:bg-primary/10',
      isDragging ? 'bg-primary/20' : 'bg-transparent'
    ]"
    @mousedown="startDrag"
  >
    <!-- Technical Handle -->
    <div 
      class="absolute bg-primary/50 transition-all duration-300"
      :class="[
        orientation === 'horizontal' 
          ? 'w-[2px] h-12 opacity-20 group-hover:opacity-100 group-hover:h-20' 
          : 'h-[2px] w-12 opacity-20 group-hover:opacity-100 group-hover:w-20',
        isDragging ? 'opacity-100' : ''
      ]"
    />
    
    <div 
      v-if="isDragging"
      class="absolute bg-primary/40"
      :class="orientation === 'horizontal' ? 'w-px h-full' : 'h-px w-full'"
    />
  </div>
</template>
