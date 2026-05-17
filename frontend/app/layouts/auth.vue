<template>
  <main class="relative flex min-h-dvh overflow-hidden bg-muted/10">
    <div class="absolute inset-0 bg-[radial-gradient(circle_at_25%_20%,oklch(0.508_0.118_165.612/0.08),transparent_34%),radial-gradient(circle_at_50%_80%,oklch(0.967_0.001_286.375/0.55),transparent_30%)] dark:bg-[radial-gradient(circle_at_25%_20%,oklch(0.696_0.17_162.48/0.13),transparent_34%),radial-gradient(circle_at_50%_80%,oklch(0.274_0.006_286.033/0.55),transparent_30%)]" />
    <div class="auth-grid absolute inset-0 opacity-60 mask-[linear-gradient(to_right,black,transparent_68%)]" />

    <div class="absolute inset-0 hidden w-[calc(100%-450px)] items-center justify-center p-8 lg:flex xl:w-[calc(100%-550px)]">
      <div class="pointer-events-none absolute inset-0 bg-primary/5" />
      <div class="auth-ambient-a pointer-events-none absolute left-[10%] top-[10%] h-[50%] w-[50%] rounded-full bg-primary/10 blur-[120px]" />
      <div class="auth-ambient-b pointer-events-none absolute bottom-[10%] left-[30%] h-[50%] w-[50%] rounded-full bg-secondary/20 blur-[100px]" />

      <div class="relative z-10 flex w-full max-w-3xl flex-col items-start space-y-8">
        <div class="space-y-3">
          <AppLogo class="h-12 w-12 transition-transform duration-500 hover:-rotate-3 hover:scale-105" />
          <h1 class="font-heading text-4xl font-bold leading-snug tracking-tight text-foreground/90 lg:text-5xl">
            The API workbench <br> for modern teams.
          </h1>
        </div>

        <div class="auth-preview-in flex aspect-16/10 w-full max-w-2xl flex-col overflow-hidden rounded-xl border border-muted-foreground/20 bg-background/60 shadow-2xl shadow-primary/10 backdrop-blur-2xl transition duration-500 hover:-translate-y-1 hover:border-primary/30 hover:shadow-primary/20">
          <div class="flex h-10 items-center gap-2 border-b border-muted-foreground/20 bg-muted/30 px-4">
            <div class="flex gap-1.5">
              <div class="size-2.5 rounded-full bg-destructive/70" />
              <div class="size-2.5 rounded-full bg-amber-500/70" />
              <div class="auth-status-dot size-2.5 rounded-full bg-emerald-500/80" />
            </div>
            <div class="auth-shimmer mx-auto h-4 w-1/3 rounded-md bg-muted/60" />
          </div>
          <div class="relative flex flex-1">
            <div class="w-1/4 space-y-4 border-r border-muted-foreground/20 bg-background/40 p-4">
              <div class="auth-shimmer mb-6 h-3 w-3/4 rounded bg-muted/60" />
              <div class="h-2 w-full rounded bg-muted/50 transition-colors duration-500 hover:bg-primary/35" />
              <div class="h-2 w-5/6 rounded bg-muted/50 transition-colors duration-500 hover:bg-primary/35" />
              <div class="h-2 w-4/5 rounded bg-muted/50 transition-colors duration-500 hover:bg-primary/35" />
              <div class="h-2 w-full rounded bg-muted/50 transition-colors duration-500 hover:bg-primary/35" />
              <div class="h-2 w-2/3 rounded bg-muted/50 transition-colors duration-500 hover:bg-primary/35" />
            </div>

            <div class="flex flex-1 flex-col gap-4 p-5">
              <div class="auth-preview-card flex gap-2">
                <div class="grid h-8 w-20 place-items-center rounded-md bg-primary/20 text-[10px] font-bold text-primary">POST</div>
                <div class="auth-progress-line h-8 flex-1 rounded-md bg-muted/40" />
                <div class="h-8 w-20 rounded-md bg-primary/90 transition-transform duration-300 hover:translate-x-0.5" />
              </div>
              <div class="auth-preview-card-delayed h-32 w-full space-y-2 rounded-md border border-muted-foreground/10 bg-muted/30 p-4 shadow-sm">
                <div class="flex items-center justify-between">
                  <div class="h-2 w-1/4 rounded bg-muted/50" />
                  <div class="h-4 w-14 rounded bg-primary/15" />
                </div>
                <div class="h-2 w-1/2 rounded bg-muted/50" />
                <div class="auth-shimmer h-2 w-1/3 rounded bg-primary/40" />
                <div class="h-2 w-2/5 rounded bg-muted/50" />
              </div>
              <div class="w-full flex-1 space-y-2 rounded-md border-t border-muted-foreground/10 bg-muted/20 p-4">
                <div class="mb-3 flex items-center gap-2">
                  <div class="h-2 w-1/6 rounded bg-muted/40" />
                  <div class="ml-auto h-4 w-12 rounded bg-emerald-500/15" />
                </div>
                <div class="h-2 w-full rounded bg-muted/40" />
                <div class="h-2 w-5/6 rounded bg-muted/40" />
                <div class="auth-shimmer h-2 w-4/5 rounded bg-secondary/40" />
              </div>
            </div>

            <div class="auth-preview-card pointer-events-none absolute bottom-5 right-5 w-44 rounded-lg border border-primary/20 bg-background/80 p-3 text-[10px] shadow-xl shadow-primary/10 backdrop-blur-xl">
              <div class="mb-2 flex items-center justify-between font-mono uppercase tracking-wider text-muted-foreground">
                <span>event stream</span>
                <span class="auth-status-dot size-2 rounded-full bg-emerald-500" />
              </div>
              <div class="space-y-1 font-mono text-muted-foreground">
                <div class="h-1.5 w-full rounded bg-muted" />
                <div class="h-1.5 w-4/5 rounded bg-muted" />
                <div class="h-1.5 w-2/3 rounded bg-primary/35" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <section class="relative z-20 ml-auto flex w-full flex-col justify-center bg-background px-6 py-12 lg:w-112.5 lg:px-16 xl:w-137.5">
      <svg
        class="pointer-events-none absolute inset-y-0 right-full hidden h-full w-37.5 xl:w-50 text-background lg:block drop-shadow-[-15px_0_25px_rgba(0,0,0,0.05)] dark:drop-shadow-[-15px_0_25px_rgba(0,0,0,0.4)]"
        viewBox="0 0 100 100"
        preserveAspectRatio="none"
        fill="currentColor"
      >
        <path d="M100,0 L0,0 C40,0 80,100 100,100 Z" />
      </svg>

      <slot />
    </section>
  </main>
</template>
