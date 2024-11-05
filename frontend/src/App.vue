<script setup>
import Sidebar from './components/nav/Sidebar.vue' // 新增
</script>

<template>
  <div>
    <Sidebar />
    <router-view>
    </router-view>
  </div>
</template>

<style>
</style>