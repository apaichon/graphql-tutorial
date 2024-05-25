const express = require('express');
const { createProxyMiddleware } = require('http-proxy-middleware');

const app = express();
// Middleware to parse JSON bodies
app.use(express.json());

// Define different targets for different queries
const targets = {
    contacts: 'http://127.0.0.1:4000/graphql',
    users: 'http://127.0.0.1:4001/graphql',
    // Add more targets as needed
};

// Proxy middleware configuration
const proxyMiddleware = (req, res, next) => {
    // Determine the target based on the query in the request
    const target = getTarget(req.query.query);
    console.log('query', req.query)
     console.log('target', target)

    // Proxy the request to the target
    if (target) {
        const proxy = createProxyMiddleware({
            target,
            changeOrigin: true,
            /*pathRewrite: {
                '^/graphql': '', // Remove /graphql prefix
            },*/
            onProxyReq: (proxyReq, req, res) => {
                // Add Authorization header to the proxy request
                console.log('req')
                if (req.headers.authorization) {
                    proxyReq.setHeader('Authorization', req.headers.authorization);
                }
            },
            onProxyRes: (proxyRes, req, res) => {
                let body = '';
                proxyRes.on('data', (chunk) => {
                    body += chunk;
                });
                proxyRes.on('end', () => {
                    // Handle the response from the proxy server
                    // Here, you can process the response or aggregate it
                    // For simplicity, we will directly send it back to the client
                    console.log('response', body)
                    res.send(body);
                });
            },
        });
        proxy(req, res, next);
    } else {
        // No target found for the query
        res.status(404).send('Target not found for the query');
    }
};

// Function to determine the target based on the query
function getTarget(query) {
    // Implement your logic to map queries to targets
    // For example, you can use a switch statement or lookup table
    if (query === 'contacts') {
        return targets.contacts;
    } else if (query === 'users') {
        return targets.users;
    }
    // Add more conditions as needed
    return null; // Return null if no target found
}

// Mount the proxy middleware
app.use('/graphql', proxyMiddleware);

// Start the server
const PORT = process.env.PORT || 9009;
app.listen(PORT, () => {
    console.log(`GraphQL proxy server is running on port ${PORT}`);
});
