<template>
  <li>
    <span>{{ task.title }} - {{ task.status }}</span>
    <button @click="$emit('delete', task.id)">🗑 Delete</button>
  </li>
</template>

<script>
export default {
  props: ['task']
};
</script>
