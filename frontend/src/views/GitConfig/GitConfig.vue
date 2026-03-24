<template>
  <div class="git-config-container">
    <RepoList />
    <ConfigEditor />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useGitConfigStore } from '@/stores/gitconfig'
import RepoList from './RepoList.vue'
import ConfigEditor from './ConfigEditor.vue'

const store = useGitConfigStore()

onMounted(async () => {
  await store.loadRepos()
  await store.loadKnownKeys()
})
</script>

<style scoped>
.git-config-container {
  display: flex;
  height: 100%;
  background: #f5f7fa;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #e8e8e8;
}
</style>
