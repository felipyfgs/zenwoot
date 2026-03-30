<script setup lang="ts">
const api = useApi()
const toast = useToast()

// Wizard state
const currentStep = ref(1)
const selectedChannel = ref<string | null>(null)
const inboxName = ref('')
const creating = ref(false)
const createdInboxId = ref<string | null>(null)
const qrData = ref('')
const qrLoading = ref(false)

// Steps definition
const steps = [
  { number: 1, title: 'Escolha o Canal', description: 'Escolha o provedor que você deseja integrar.' },
  { number: 2, title: 'Criar Caixa de Entrada', description: 'Configure sua caixa de entrada.' },
  { number: 3, title: 'Adicionar Agentes', description: 'Adicionar agentes à caixa de entrada criada.' },
  { number: 4, title: 'Então!', description: 'Está tudo pronto para começar!' }
]

// Available channels (only WhatsApp for now)
const channels = [
  { key: 'whatsapp', title: 'WhatsApp', description: 'Atenda seus clientes no WhatsApp', icon: 'i-simple-icons-whatsapp' }
]

function selectChannel(channel: string) {
  selectedChannel.value = channel
  currentStep.value = 2
}

async function createInbox() {
  if (!inboxName.value.trim()) {
    toast.add({ title: 'Por favor, digite um nome para a caixa de entrada', color: 'error' })
    return
  }

  creating.value = true
  try {
    const res = await api.post<{ inboxId: string }>('/inboxes', {
      name: inboxName.value,
      channelType: selectedChannel.value
    })
    createdInboxId.value = res.inboxId
    toast.add({ title: 'Caixa de entrada criada!', color: 'success' })
    currentStep.value = 3
  } catch (e: unknown) {
    toast.add({ title: e instanceof Error ? e.message : 'Falha ao criar', color: 'error' })
  } finally {
    creating.value = false
  }
}

function skipAgents() {
  currentStep.value = 4
  loadQR()
}

async function loadQR() {
  if (!createdInboxId.value) return
  qrLoading.value = true
  try {
    // First connect the inbox
    await api.post(`/inboxes/${createdInboxId.value}/connect`)
    // Then fetch QR
    const res = await api.get<{ qr: string }>(`/inboxes/${createdInboxId.value}/qr`)
    qrData.value = res?.qr || ''
  } catch {
    toast.add({ title: 'Falha ao buscar QR Code', color: 'error' })
  } finally {
    qrLoading.value = false
  }
}

function goBack() {
  if (currentStep.value > 1) {
    currentStep.value--
    if (currentStep.value === 1) {
      selectedChannel.value = null
    }
  } else {
    navigateTo('/settings/inboxes')
  }
}

function finish() {
  navigateTo('/settings/inboxes')
}
</script>

<template>
  <div class="min-h-screen bg-default">
    <!-- Header -->
    <div class="border-b border-default px-6 py-4">
      <div class="flex flex-wrap items-center gap-2 sm:gap-3">
        <UDashboardSidebarCollapse />

        <button class="flex items-center gap-2 text-sm text-muted hover:text-default" @click="goBack">
          <UIcon name="i-lucide-chevron-left" class="size-4" />
          <span>Anterior</span>
        </button>

        <UIcon name="i-lucide-chevron-right" class="size-4 text-muted" />

        <h1 class="text-sm font-semibold text-highlighted">
          Caixas de Entrada
        </h1>
      </div>
    </div>

    <!-- Main content -->
    <div class="px-4 py-6 sm:px-6 lg:px-8">
      <div class="mx-auto overflow-hidden rounded-2xl border border-default">
        <div class="flex w-full max-w-[1600px] flex-col lg:flex-row">
          <!-- Sidebar with steps -->
          <div class="w-full shrink-0 border-b border-default p-4 lg:w-80 lg:border-r lg:border-b-0 lg:p-6">
            <div class="flex lg:flex-col gap-4 lg:gap-6 overflow-x-auto lg:overflow-visible">
              <div v-for="step in steps" :key="step.number" class="flex items-start gap-3 lg:gap-4 shrink-0">
                <div
                  class="size-8 rounded-full flex items-center justify-center shrink-0 text-sm font-medium"
                  :class="[
                    currentStep > step.number ? 'bg-primary text-white' : '',
                    currentStep === step.number ? 'bg-primary text-white' : '',
                    currentStep < step.number ? 'bg-muted/20 text-muted' : ''
                  ]"
                >
                  <UIcon v-if="currentStep > step.number" name="i-lucide-check" class="size-4" />
                  <span v-else>{{ step.number }}</span>
                </div>
                <div class="hidden lg:block">
                  <p
                    class="text-sm font-medium"
                    :class="currentStep >= step.number ? 'text-primary' : 'text-muted'"
                  >
                    {{ step.title }}
                  </p>
                  <p class="text-xs text-muted mt-0.5">
                    {{ step.description }}
                  </p>
                </div>
              </div>
            </div>
          </div>

          <!-- Content area -->
          <div class="min-w-0 flex-1 px-4 py-6 sm:px-6 lg:px-8 lg:py-8">
            <!-- Step 1: Channel Selection -->
            <div v-if="currentStep === 1">
              <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
                <button
                  v-for="channel in channels"
                  :key="channel.key"
                  class="flex flex-col items-center gap-3 p-6 rounded-xl border border-default bg-elevated hover:border-primary hover:bg-primary/5 transition-colors text-left"
                  @click="selectChannel(channel.key)"
                >
                  <div class="size-12 rounded-xl bg-muted/10 flex items-center justify-center">
                    <UIcon :name="channel.icon" class="size-6 text-muted" />
                  </div>
                  <div class="text-center">
                    <p class="font-medium text-highlighted">
                      {{ channel.title }}
                    </p>
                    <p class="text-xs text-muted mt-1">
                      {{ channel.description }}
                    </p>
                  </div>
                </button>

                <!-- Coming soon placeholders -->
                <div class="flex flex-col items-center gap-3 p-6 rounded-xl border border-dashed border-muted/30 opacity-50">
                  <div class="size-12 rounded-xl bg-muted/10 flex items-center justify-center">
                    <UIcon name="i-lucide-message-square" class="size-6 text-muted" />
                  </div>
                  <div class="text-center">
                    <p class="font-medium text-muted">
                      Site
                    </p>
                    <p class="text-xs text-muted mt-1">
                      Em breve
                    </p>
                  </div>
                </div>

                <div class="flex flex-col items-center gap-3 p-6 rounded-xl border border-dashed border-muted/30 opacity-50">
                  <div class="size-12 rounded-xl bg-muted/10 flex items-center justify-center">
                    <UIcon name="i-simple-icons-telegram" class="size-6 text-muted" />
                  </div>
                  <div class="text-center">
                    <p class="font-medium text-muted">
                      Telegram
                    </p>
                    <p class="text-xs text-muted mt-1">
                      Em breve
                    </p>
                  </div>
                </div>

                <div class="flex flex-col items-center gap-3 p-6 rounded-xl border border-dashed border-muted/30 opacity-50">
                  <div class="size-12 rounded-xl bg-muted/10 flex items-center justify-center">
                    <UIcon name="i-lucide-mail" class="size-6 text-muted" />
                  </div>
                  <div class="text-center">
                    <p class="font-medium text-muted">
                      E-Mail
                    </p>
                    <p class="text-xs text-muted mt-1">
                      Em breve
                    </p>
                  </div>
                </div>
              </div>
            </div>

            <!-- Step 2: Inbox Form -->
            <div v-else-if="currentStep === 2" class="w-full max-w-5xl">
              <div class="rounded-2xl border border-default bg-elevated p-4 sm:p-6 lg:p-8">
                <div class="grid gap-6 lg:grid-cols-[minmax(0,1fr)_auto] lg:items-end">
                  <UFormField label="Nome da Caixa de Entrada" name="name">
                    <UInput
                      v-model="inboxName"
                      placeholder="Por favor, digite um nome para caixa de entrada"
                      class="w-full"
                      size="lg"
                    />
                  </UFormField>

                  <div class="flex w-full lg:w-auto lg:justify-end">
                    <UButton
                      :loading="creating"
                      :disabled="!inboxName.trim()"
                      size="lg"
                      class="w-full lg:w-auto"
                      @click="createInbox"
                    >
                      Criar canal do WhatsApp
                    </UButton>
                  </div>
                </div>
              </div>
            </div>

            <!-- Step 3: Add Agents (optional) -->
            <div v-else-if="currentStep === 3" class="w-full max-w-5xl">
              <div class="rounded-2xl border border-default bg-elevated px-4 py-8 text-center sm:px-6 lg:px-8">
                <UIcon name="i-lucide-users" class="size-16 text-muted mx-auto mb-4" />
                <h3 class="text-lg font-medium text-highlighted mb-2">
                  Adicionar Agentes
                </h3>
                <p class="text-sm text-muted mb-6">
                  Você pode adicionar agentes agora ou fazer isso depois nas configurações.
                </p>
                <div class="flex gap-3 justify-center">
                  <UButton variant="ghost" color="neutral" @click="skipAgents">
                    Pular por agora
                  </UButton>
                  <UButton @click="skipAgents">
                    Continuar
                  </UButton>
                </div>
              </div>
            </div>

            <!-- Step 4: Success + QR -->
            <div v-else-if="currentStep === 4" class="w-full max-w-5xl">
              <div class="py-8 text-center">
                <div class="size-16 rounded-full bg-success/10 flex items-center justify-center mx-auto mb-4">
                  <UIcon name="i-lucide-check" class="size-8 text-success" />
                </div>
                <h3 class="text-lg font-medium text-highlighted mb-2 text-center">
                  Caixa de entrada criada!
                </h3>
                <p class="text-sm text-muted mb-6 text-center">
                  Escaneie o QR Code abaixo com seu WhatsApp para conectar.
                </p>

                <div class="bg-elevated rounded-xl p-6 border border-default">
                  <div v-if="qrLoading" class="flex flex-col items-center gap-4 py-8">
                    <UIcon name="i-lucide-loader-circle" class="size-12 animate-spin text-muted" />
                    <p class="text-sm text-muted">
                      Gerando QR Code...
                    </p>
                  </div>
                  <div v-else-if="qrData" class="flex flex-col items-center gap-4">
                    <img
                      :src="qrData"
                      class="size-64 object-contain rounded-lg"
                      alt="QR Code"
                    >
                    <p class="text-xs text-muted">
                      Abra o WhatsApp no seu telefone, vá em Dispositivos Conectados e escaneie este código.
                    </p>
                  </div>
                  <div v-else class="py-8">
                    <p class="text-sm text-muted">
                      Não foi possível gerar o QR Code. Tente novamente nas configurações da caixa de entrada.
                    </p>
                  </div>
                </div>

                <div class="mt-6">
                  <UButton size="lg" @click="finish">
                    Concluir Configuração
                  </UButton>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
