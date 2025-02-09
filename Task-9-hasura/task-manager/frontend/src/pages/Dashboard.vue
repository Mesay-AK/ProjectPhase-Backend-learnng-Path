<template>
  <div>
    <h1>Dashboard</h1>
    <TaskForm />
    <TaskList />
  </div>
</template>

<script>
import TaskList from '../components/TaskList.vue';
import TaskForm from '../components/TaskForm.vue';

export default {
  components: { TaskList, TaskForm }
};
</script>
