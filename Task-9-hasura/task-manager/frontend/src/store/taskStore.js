import { defineStore } from 'pinia';
import { useQuery, useMutation } from '@vue/apollo-composable';
import { GET_TASKS, ADD_TASK, DELETE_TASK } from '../graphql/tasks';

export const useTaskStore = defineStore('taskStore', {
  state: () => ({
    tasks: [],
  }),
  actions: {
    async fetchTasks() {
      const { result } = useQuery(GET_TASKS);
      this.tasks = result.value?.tasks || [];
    },
    async addTask(title, description) {
      const { mutate } = useMutation(ADD_TASK);
      await mutate({ title, description });
      this.fetchTasks();
    },
    async deleteTask(id) {
      const { mutate } = useMutation(DELETE_TASK);
      await mutate({ id });
      this.fetchTasks();
    }
  }
});

