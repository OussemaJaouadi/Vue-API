<script setup lang="ts">
import {
  PhBracketsCurly,
  PhCaretDown,
  PhCheck,
  PhCode,
  PhCopy,
  PhTextT,
} from '@phosphor-icons/vue'
import { closeBrackets, closeBracketsKeymap } from '@codemirror/autocomplete'
import {
  defaultKeymap,
  history,
  historyKeymap,
  indentWithTab,
} from '@codemirror/commands'
import { html } from '@codemirror/lang-html'
import { json } from '@codemirror/lang-json'
import { xml } from '@codemirror/lang-xml'
import { yaml } from '@codemirror/lang-yaml'
import {
  bracketMatching,
  defaultHighlightStyle,
  foldGutter,
  foldKeymap,
  HighlightStyle,
  indentOnInput,
  syntaxHighlighting,
} from '@codemirror/language'
import {
  Compartment,
  EditorState,
  type Extension,
} from '@codemirror/state'
import {
  crosshairCursor,
  drawSelection,
  dropCursor,
  EditorView,
  highlightActiveLine,
  highlightActiveLineGutter,
  highlightSpecialChars,
  keymap,
  lineNumbers,
  rectangularSelection,
} from '@codemirror/view'
import { tags as t } from '@lezer/highlight'
import type { Component } from 'vue'
import type { BodyLanguage } from '~/composables/useWorkbench'

const props = withDefaults(defineProps<{
  modelValue: string
  language?: BodyLanguage
  editable?: boolean
  label?: string
}>(), {
  language: 'json',
  editable: false,
  label: 'Body',
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
  'update:language': [value: BodyLanguage]
}>()

const copied = ref(false)
const editorRoot = ref<HTMLElement>()
let editorView: EditorView | undefined

const languageCompartment = new Compartment()
const editableCompartment = new Compartment()

const languages: Array<{
  value: BodyLanguage
  label: string
  description: string
  icon: Component
  accentClass: string
}> = [
  { value: 'json', label: 'JSON', description: 'objects and arrays', icon: PhBracketsCurly, accentClass: 'bg-sky-500' },
  { value: 'xml', label: 'XML', description: 'tags and attrs', icon: PhCode, accentClass: 'bg-violet-500' },
  { value: 'html', label: 'HTML', description: 'markup preview', icon: PhCode, accentClass: 'bg-orange-500' },
  { value: 'yaml', label: 'YAML', description: 'indented config', icon: PhTextT, accentClass: 'bg-amber-500' },
  { value: 'text', label: 'Text', description: 'plain response', icon: PhTextT, accentClass: 'bg-muted-foreground' },
]

const languageExtension = (language: BodyLanguage): Extension => {
  switch (language) {
    case 'json':
      return json()
    case 'xml':
      return xml()
    case 'html':
      return html()
    case 'yaml':
      return yaml()
    default:
      return []
  }
}

const editableExtension = (editable: boolean): Extension => [
  EditorView.editable.of(editable),
  EditorState.readOnly.of(!editable),
]

const foldMarker = (open: boolean) => {
  const marker = document.createElement('span')
  marker.className = 'cm-workbench-fold-marker'
  marker.textContent = open ? 'v' : '>'
  return marker
}

const workbenchHighlightStyle = HighlightStyle.define([
  { tag: t.propertyName, color: 'oklch(0.52 0.17 253)', fontWeight: '700' },
  { tag: t.string, color: 'oklch(0.53 0.18 16)' },
  { tag: t.number, color: 'oklch(0.58 0.14 70)' },
  { tag: t.bool, color: 'oklch(0.52 0.18 302)', fontWeight: '700' },
  { tag: t.atom, color: 'oklch(0.55 0.02 286)', fontStyle: 'italic' },
  { tag: t.keyword, color: 'oklch(0.45 0.13 165)', fontWeight: '700' },
  { tag: t.operator, color: 'oklch(0.50 0.02 286)' },
  { tag: t.variableName, color: 'oklch(0.47 0.13 250)' },
  { tag: t.typeName, color: 'oklch(0.50 0.15 290)', fontWeight: '700' },
  { tag: t.className, color: 'oklch(0.48 0.15 165)', fontWeight: '700' },
  { tag: t.labelName, color: 'oklch(0.48 0.15 165)', fontWeight: '700' },
  { tag: t.comment, color: 'oklch(0.55 0.02 286)', fontStyle: 'italic' },
  { tag: t.meta, color: 'oklch(0.52 0.14 70)' },
  { tag: t.regexp, color: 'oklch(0.55 0.15 25)' },
  { tag: t.invalid, color: 'oklch(0.58 0.22 27)', textDecoration: 'underline' },
])

const workbenchEditorTheme = EditorView.theme({
  '&': {
    background: 'transparent',
    color: 'var(--foreground)',
    fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace',
    fontSize: '12px',
    height: '100%',
  },
  '&.cm-focused': {
    outline: 'none',
  },
  '.cm-scroller': {
    fontFamily: 'inherit',
    lineHeight: '1.5rem',
    overflow: 'auto',
  },
  '.cm-content': {
    minHeight: '100%',
    padding: '1rem',
  },
  '.cm-line': {
    padding: '0 0.125rem',
  },
  '.cm-gutters': {
    background: 'color-mix(in oklch, var(--muted) 18%, transparent)',
    borderRight: '1px solid color-mix(in oklch, var(--primary) 10%, transparent)',
    color: 'color-mix(in oklch, var(--muted-foreground) 45%, transparent)',
  },
  '.cm-lineNumbers .cm-gutterElement': {
    minWidth: '2.75rem',
    padding: '0 0.75rem 0 0.5rem',
  },
  '.cm-foldGutter .cm-gutterElement': {
    cursor: 'pointer',
    padding: '0 0.25rem',
  },
  '.cm-workbench-fold-marker': {
    color: 'color-mix(in oklch, var(--primary) 70%, transparent)',
    display: 'inline-block',
    fontSize: '9px',
    fontWeight: '900',
    lineHeight: '1.5rem',
    width: '0.75rem',
  },
  '.cm-activeLine': {
    background: 'color-mix(in oklch, var(--primary) 7%, transparent)',
  },
  '.cm-activeLineGutter': {
    background: 'color-mix(in oklch, var(--primary) 10%, transparent)',
    color: 'var(--primary)',
  },
  '.cm-selectionBackground, &.cm-focused .cm-selectionBackground': {
    background: 'color-mix(in oklch, var(--primary) 22%, transparent)',
  },
  '.cm-cursor': {
    borderLeftColor: 'var(--primary)',
  },
  '.cm-foldPlaceholder': {
    background: 'color-mix(in oklch, var(--primary) 10%, transparent)',
    border: '1px solid color-mix(in oklch, var(--primary) 20%, transparent)',
    borderRadius: '0',
    color: 'var(--primary)',
    margin: '0 0.25rem',
    padding: '0 0.35rem',
  },
}, { dark: false })

const editorExtensions = () => [
  lineNumbers(),
  foldGutter({ markerDOM: foldMarker }),
  highlightSpecialChars(),
  history(),
  drawSelection(),
  dropCursor(),
  rectangularSelection(),
  crosshairCursor(),
  highlightActiveLine(),
  highlightActiveLineGutter(),
  indentOnInput(),
  bracketMatching(),
  closeBrackets(),
  syntaxHighlighting(defaultHighlightStyle, { fallback: true }),
  syntaxHighlighting(workbenchHighlightStyle),
  workbenchEditorTheme,
  languageCompartment.of(languageExtension(props.language)),
  editableCompartment.of(editableExtension(props.editable)),
  keymap.of([
    indentWithTab,
    ...closeBracketsKeymap,
    ...defaultKeymap,
    ...historyKeymap,
    ...foldKeymap,
  ]),
  EditorView.updateListener.of((update) => {
    if (!update.docChanged) return
    emit('update:modelValue', update.state.doc.toString())
  }),
]

const createEditor = () => {
  if (!editorRoot.value || editorView) return

  editorView = new EditorView({
    parent: editorRoot.value,
    state: EditorState.create({
      doc: props.modelValue,
      extensions: editorExtensions(),
    }),
  })
}

const destroyEditor = () => {
  editorView?.destroy()
  editorView = undefined
}

const formatJson = () => {
  try {
    emit('update:modelValue', JSON.stringify(JSON.parse(props.modelValue), null, 2))
  } catch {
    // Invalid JSON stays untouched.
  }
}

const minifyJson = () => {
  try {
    emit('update:modelValue', JSON.stringify(JSON.parse(props.modelValue)))
  } catch {
    // Invalid JSON stays untouched.
  }
}

const copyCode = async () => {
  if (!import.meta.client) return
  await navigator.clipboard.writeText(props.modelValue)
  copied.value = true
  window.setTimeout(() => {
    copied.value = false
  }, 1200)
}

watch(() => props.modelValue, (value) => {
  if (!editorView || value === editorView.state.doc.toString()) return

  editorView.dispatch({
    changes: {
      from: 0,
      to: editorView.state.doc.length,
      insert: value,
    },
  })
})

watch(() => props.language, (language) => {
  editorView?.dispatch({
    effects: languageCompartment.reconfigure(languageExtension(language)),
  })
})

watch(() => props.editable, (editable) => {
  editorView?.dispatch({
    effects: editableCompartment.reconfigure(editableExtension(editable)),
  })
})

onMounted(createEditor)
onBeforeUnmount(destroyEditor)
</script>

<template>
  <div class="flex h-full min-h-0 flex-col overflow-hidden bg-card">
    <div class="flex h-8 shrink-0 items-center justify-between border-b bg-muted/20 px-2">
      <div class="flex items-center gap-2 font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground">
        <component :is="language === 'json' ? PhBracketsCurly : PhTextT" class="size-3 text-primary" />
        {{ label }}
      </div>

      <div class="flex items-center gap-1">
        <UiTooltip v-if="editable && language === 'json'">
          <UiTooltipTrigger as-child>
            <button
              class="h-5 border border-border bg-background px-1.5 font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground transition-colors hover:border-primary/40 hover:text-primary"
              type="button"
              @click="formatJson"
            >
              Format
            </button>
          </UiTooltipTrigger>
          <UiTooltipContent side="bottom">Pretty print JSON</UiTooltipContent>
        </UiTooltip>

        <UiTooltip v-if="editable && language === 'json'">
          <UiTooltipTrigger as-child>
            <button
              class="h-5 border border-border bg-background px-1.5 font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground transition-colors hover:border-primary/40 hover:text-primary"
              type="button"
              @click="minifyJson"
            >
              Minify
            </button>
          </UiTooltipTrigger>
          <UiTooltipContent side="bottom">Compact JSON</UiTooltipContent>
        </UiTooltip>

        <UiTooltip>
          <UiTooltipTrigger as-child>
            <button
              class="grid size-5 place-items-center border border-border bg-background text-muted-foreground transition-colors hover:border-primary/40 hover:text-primary"
              type="button"
              @click="copyCode"
            >
              <PhCheck v-if="copied" class="size-3 text-primary" />
              <PhCopy v-else class="size-3" />
            </button>
          </UiTooltipTrigger>
          <UiTooltipContent side="bottom">{{ copied ? 'Copied' : 'Copy code' }}</UiTooltipContent>
        </UiTooltip>

        <UiDropdownMenu>
          <UiDropdownMenuTrigger as-child>
            <button
              aria-label="Syntax mode"
              class="flex h-5 items-center gap-1.5 border border-border bg-background px-1.5 font-mono text-[9px] font-black uppercase tracking-widest text-muted-foreground transition-colors hover:border-primary/40 hover:text-foreground"
              type="button"
            >
              {{ language }}
              <PhCaretDown class="size-2.5 opacity-50" />
            </button>
          </UiDropdownMenuTrigger>
          <UiDropdownMenuContent align="end" class="w-44 rounded-none border-2 p-1">
            <UiDropdownMenuItem
              v-for="item in languages"
              :key="item.value"
              class="group flex cursor-pointer items-center gap-2 rounded-none border border-transparent px-2 py-1.5 font-mono text-[10px] transition-colors focus:bg-primary/8"
              :class="language === item.value ? 'border-primary/20 bg-primary/10 text-foreground' : 'text-muted-foreground hover:text-foreground'"
              @click="emit('update:language', item.value)"
            >
              <span class="flex size-5 shrink-0 items-center justify-center border border-border bg-background">
                <component :is="item.icon" class="size-3" />
              </span>
              <span class="min-w-0 flex-1">
                <span class="flex items-center gap-1.5">
                  <span class="font-black uppercase tracking-widest">{{ item.label }}</span>
                  <span class="h-1.5 w-1.5" :class="item.accentClass" />
                </span>
                <span class="block truncate text-[8px] font-bold uppercase tracking-wider text-muted-foreground/60">
                  {{ item.description }}
                </span>
              </span>
              <PhCheck v-if="language === item.value" class="size-3 shrink-0 text-primary" />
            </UiDropdownMenuItem>
          </UiDropdownMenuContent>
        </UiDropdownMenu>
      </div>
    </div>

    <div ref="editorRoot" class="min-h-0 flex-1 overflow-hidden" />
  </div>
</template>
