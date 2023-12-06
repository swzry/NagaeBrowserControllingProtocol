<template>
  <q-layout view="hHh lpR fFf">

    <q-header class="bg-primary text-white">
      <q-bar class="q-electron-drag">
        <q-icon name="web_asset" />

        <div class="q-ml-md cursor-pointer non-selectable">
          App
          <q-menu auto-close>
            <q-list dense style="min-width: 100px">
              <q-item clickable>
                <q-item-section>About</q-item-section>
              </q-item>
              <q-item clickable @click="openDevTool">
                <q-item-section>Open Dev Tool</q-item-section>
              </q-item>
              <q-separator />
              <q-item clickable @click="setTitleBarDisplay(true)">
                <q-item-section>Show Title Bar</q-item-section>
              </q-item>
              <q-item clickable @click="setTitleBarDisplay(false)">
                <q-item-section>Hide Title Bar</q-item-section>
              </q-item>
              <q-separator />
              <q-item clickable @click="closeApp">
                <q-item-section>Exit</q-item-section>
              </q-item>
            </q-list>
          </q-menu>
        </div>

        <q-space />

        <div>NBCP Test Frontend</div>

        <q-space />

        <q-btn dense flat icon="minimize" @click="minimize" />
        <q-btn dense flat icon="crop_square" @click="toggleMaximize" />
        <q-btn dense flat icon="close" @click="closeWindow" />
      </q-bar>
    </q-header>

    <q-page-container>
      <router-view />
    </q-page-container>

    <q-footer class="bg-grey-8 text-white">
      <q-bar dense>
        <q-icon name="check" />
        <div>NBCP Test Frontend</div>
      </q-bar>
    </q-footer>

  </q-layout>
</template>

<script>
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'App',
  methods:{
    minimize(){
      this.$nbcp.rpc("minimize", {"name": this.$nbcp.getWindowName()})
    },
    toggleMaximize(){
      this.$nbcp.rpc("toggle_maximize", {"name": this.$nbcp.getWindowName()})
    },
    closeWindow(){
      this.$nbcp.rpc("close_window", {"name": this.$nbcp.getWindowName()})
    },
    closeApp(){
      this.$nbcp.endSession()
    },
    openDevTool(){
      this.$nbcp.rpc("open_dev_tool", {"name": this.$nbcp.getWindowName()})
    },
    setTitleBarDisplay(val){
      this.$nbcp.rpc("set_titlebar_display", {"name": this.$nbcp.getWindowName(), "value": val})
    }
  }
})
</script>
