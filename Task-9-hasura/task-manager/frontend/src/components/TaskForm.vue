<template>
  <form @submit.prevent="addTask">
    <input v-model="title" placeholder="Task Title" required />
    <input v-model="description" placeholder="Task Description" required />
    <button type="submit">Add Task</button>
  </form>
</template>

<script>
import { ref } from 'vue';
import { useTaskStore } from '../store/taskStore';

export default {
  setup() {
    const title = ref('');
    const description = ref('');
    const taskStore = useTaskStore();

    const addTask = () => {
      if (!title.value.trim() || !description.value.trim()) return;
      taskStore.addTask(title.value, description.value);
      title.value = '';
      description.value = '';
    };

    return { title, description, addTask };
  }
};
</script>
