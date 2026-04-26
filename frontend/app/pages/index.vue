<template>
  <div class="max-w-4xl mx-auto space-y-6">
    <div class="flex items-center gap-2">
      <div class="bg-primary/10 text-primary px-3 py-1 rounded font-bold text-sm">GET</div>
      <div class="flex-1 bg-muted px-3 py-1 rounded text-sm text-muted-foreground border">
        http://localhost:8080/healthz
      </div>
      <button class="bg-primary text-primary-foreground px-4 py-1 rounded text-sm font-medium hover:opacity-90 transition-opacity">
        Send
      </button>
    </div>

    <div class="border rounded-lg bg-card">
      <div class="border-b px-4 py-2 flex gap-4 text-xs font-medium text-muted-foreground">
        <div class="text-foreground border-b-2 border-primary pb-2 -mb-2">Params</div>
        <div>Headers</div>
        <div>Body</div>
        <div>Auth</div>
      </div>
      <div class="p-8 text-center text-muted-foreground text-sm">
        Request parameters will appear here.
      </div>
    </div>
    
    <div class="space-y-2">
      <h3 class="text-sm font-semibold">Response</h3>
      <div class="border rounded-lg bg-card min-h-[200px] p-4 flex items-center justify-center text-muted-foreground text-sm">
        Click Send to see a response
      </div>
    </div>
  </div>
</template>
