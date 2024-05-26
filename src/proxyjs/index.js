const express = require('express');
const axios = require('axios');
const { parse, visit } = require('graphql/language');
const { print } = require('graphql');

const app = express();
app.use(express.json());

let token =""
let headers = {}
async function login () {
    const data = await axios.post('http://localhost:4000/login', {
        user_name: "puppy",
        password: 'P@ssw0rd'
    })
    // console.log("token", data.data.token)
    token = data.data.token
    headers = {
        'Authorization': token // Include 'Bearer ' before the token
      };
}

login()

async function executeQuery(query, variables = {}) {
    const contactQuery = [];
    const userQuery = [];
  
    // Parse the GraphQL query
    const ast = parse(query);
  
    // Visit the AST and separate the queries for contacts and users
    visit(ast, {
      Field(node) {
        if (node.name.value === 'contacts') {
          contactQuery.push(node);
        } else if (node.name.value === 'data') {
          userQuery.push(node);
        }
      },
    });

    // console.log('contactQuery:',print(contactQuery))
  
    const contactsResponse = contactQuery.length > 0
      ? await axios.post('http://localhost:4000/graphql', {
          query: '{' +print({ kind: 'Document', definitions: contactQuery }) +'}',
          variables,
        }, { headers })
      : { data: { data: null } };
  
    const usersResponse = userQuery.length > 0
      ? await axios.post('http://localhost:4001/graphql', {
          query: '{' +print({ kind: 'Document', definitions: userQuery }) + '}',
          variables,
        }, { headers })
      : { data: { data: null } };
  
    return {
      data: {
        ...contactsResponse.data.data,
        ...usersResponse.data.data,
      },
    };
  }

app.post('/graphql', async (req, res) => {
    const { query, variables } = req.body;

    executeQuery(query, variables)
        .then(response => res.json (response.data))
        .catch(error => res.status(500).json({ message: `Error: ${error}` }));
});

app.listen(9009, () => console.log('Proxy server listening on port 9009'));