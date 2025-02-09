package graphql
// GraphQL queries and mutations
const GetTasksQuery = `
query {
  tasks {
    id
    title
    description
    status
  }
}`

const AddTaskMutation = `
mutation($title: String!, $description: String!) {
  insert_tasks(objects: { title: $title, description: $description }) {
    returning {
      id
      title
      description
      status
    }
  }
}`

const UpdateTaskMutation = `
mutation($id: uuid!, $status: String!) {
  update_tasks_by_pk(pk_columns: { id: $id }, _set: { status: $status }) {
    id
    title
    description
    status
  }
}`

const DeleteTaskMutation = `
mutation($id: uuid!) {
  delete_tasks_by_pk(id: $id) {
    id
  }
}`
