<template>
  <div>
    <h2>Task List</h2>
    <ul>
      <TaskItem v-for="task in tasks" :key="task.id" :task="task" @delete="deleteTask(task.id)" />
    </ul>
  </div>
</template>

<script>
import { useTaskStore } from '../store/taskStore';
import TaskItem from './TaskItem.vue';
import { onMounted } from 'vue';

export default {
  components: { TaskItem },
  setup() {
    const taskStore = useTaskStore();
    onMounted(taskStore.fetchTasks);
    return { tasks: taskStore.tasks, deleteTask: taskStore.deleteTask };
  }
};
</script>
