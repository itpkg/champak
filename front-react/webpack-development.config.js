module.exports = require("./webpack.config")({
    backend: 'http://localhost:8080',
    engines: ['auth', 'reading', 'forum', 'ops', 'shop']
});
