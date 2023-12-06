<template>
  <q-page class="flex flex-center">
    <q-card class="my-card bg-warning">
      <q-card-section>
        <div class="text-h6">404</div>
        <div class="text-subtitle2">Page Not Found</div>
      </q-card-section>

      <q-card-actions>
        <q-btn flat to="/" >Back to Homepage</q-btn>
      </q-card-actions>
    </q-card>
  </q-page>
</template>

<script>
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'ErrorNotFound'
})
</script>
