<template>
  <div class="sidebar">
    <ul>
      <li>
        <router-link to="/">端口扫描器</router-link>
      </li>
      <!-- <li>
        <router-link to="/about">关于</router-link>
      </li> -->
    </ul>
  </div>
</template>

<script setup>
</script>

<style scoped>
.sidebar {
  position: fixed;
  top: 0;
  left: 0;
  width: 200px;
  height: 100vh;
  /* background-color: #333; */
  color: white;
  padding-top: 20px;
}

.sidebar ul {
  list-style-type: none;
  padding: 0;
}

.sidebar li {
  margin: 20px 0;
}

.sidebar a {
  color: white;
  text-decoration: none;
  padding: 10px;
  display: block;
}

.sidebar a:hover {
  background-color: #9c9a9a;
}
</style>
