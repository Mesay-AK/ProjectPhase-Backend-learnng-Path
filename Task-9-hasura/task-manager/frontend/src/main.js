import { createApp, provide, h } from 'vue';
import { ApolloProvider } from '@vue/apollo-option';
import { apolloClient } from './graphql/apolloClient';
import App from './App.vue';
import router from './router';
import { createPinia } from 'pinia';

const app = createApp({
  setup() {
    provide(ApolloProvider, { defaultClient: apolloClient });
  },
  render: () => h(App),
});

app.use(router);
app.use(createPinia());
app.mount('#app');
