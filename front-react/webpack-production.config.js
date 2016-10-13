module.exports = require("./webpack.config")({
    minify: true,
    backend: '/api/v1',
    engines: ['auth', 'reading', 'forum', 'ops', 'shop']
});