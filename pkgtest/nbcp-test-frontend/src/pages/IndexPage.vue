<template>
  <q-page class="flex flex-center">
    <q-card class="my-card" flat bordered>
      <!--<q-avatar size="100px" font-size="52px" color="positive" text-color="white" icon="check" />-->
      <q-card-section>
        <div class="text-h6">TEST</div>
      </q-card-section>

      <q-separator />

      <q-card-section>
        <div class="row flex-center q-pa-sm">
          <q-avatar size="100px" font-size="52px" color="positive" text-color="white" icon="check" />
        </div>
        <div class="row flex-center q-pa-sm">
          <div class="text-h6">OK</div>
        </div>
      </q-card-section>

      <q-separator />

      <q-card-actions>
        <q-btn flat round icon="event" />
        <q-btn flat @click="showAboutInNewWin">
          关于
        </q-btn>
        <q-btn flat color="primary" @click="test">
          Test
        </q-btn>
      </q-card-actions>
    </q-card>
  </q-page>
</template>

<script>
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'IndexPage',
  methods: {
    test(){
      console.log("test");
      console.log(this.$nbcp);
      this.$nbcp.test();
    },
    showAboutInNewWin(){
      this.$nbcp.rpc("new_window", {
        "name": "about",
        "url": this.$nbcp.getCurrentRootUrl() + "about",
        "hide_frame": true,
        "disable_maximize_btn": true,
        "no_resizable": true,
        "width": 600,
        "height": 300,
      })
    }
  }
})
</script>
