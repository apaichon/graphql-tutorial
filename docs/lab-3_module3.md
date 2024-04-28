# Module 3: Define Schema

## Lab 2: Generate Schema
1. Create schema.graphql at root project.
```gql
# Define custom scalar types if needed
scalar Date

# Define types
type User {
  id: ID!
  username: String!
  email: String!
  createdAt: Date!
  posts: [Post!]!
}

type Post {
  id: ID!
  title: String!
  content: String!
  author: User!
  createdAt: Date!
}

# Define input types for mutations
input CreateUserInput {
  username: String!
  email: String!
}

input CreatePostInput {
  title: String!
  content: String!
  authorId: ID!
}

# Define query type
type Query {
  getUser(id: ID!): User
  getPost(id: ID!): Post
  getAllPosts: [Post!]!
}

# Define mutation type
type Mutation {
  createUser(input: CreateUserInput!): User!
  createPost(input: CreatePostInput!): Post!
}

```
2. Download github.com/99designs/gqlgen Library
```sh
 go get github.com/99designs/gqlgen
```

3. Build to binary
```sh
go build -o gqlgen github.com/99designs/gqlgen 
```
4. Copy to /usr/local/bin
```sh
sudo mv gqlgen /usr/local/bin
```
5. Create gqlgen.yml at root project.
```yml
schema:
  - schema.graphql
exec:
  package: model
  filename: generated_model.go
model:
  package: model
resolver:
  layout: follow-schema

```
6. Run gqlgen generate
```sh
gqlgen generate
```

