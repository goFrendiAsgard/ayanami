const http = require('http');
const url = require('url');

const servicename = 'cities';
const hostname = '127.0.0.1';
const port = 3000;
const serviceConfigs = [
    {
        pathname: 'getCities',
        fn: require(`${__dirname}/lib.js`).getCities,
        inputs: [],
    }
];

function respond(pathname, fn, inputs, req, res) {
    const reqUrl = url.parse(req.url, true);
    const reqPathname = reqUrl.pathname;
    if (reqPathname != `/${pathname}` && req.method === 'POST') {
        return false;
    }
    const response = {
        result: {},
        err: null,
    }
    // set body
    let body = '';
    req.on('data', function (chunk) {
        body += chunk;
    });
    // request completed
    req.on('end', function () {
        res.setHeader('Content-Type', 'application/json');
        try {
            const postBody = JSON.parse(body);
            const fnInputs = inputs.map((key) => postBody[key]);
            const fnOutput = fn(...fnInputs);
            response.result = fnOutput;
            res.statusCode = 200;
            console.log(res.statusCode, {reqPathName, fnInputs, fnOutput})
        } catch (err) {
            res.statusCode = 500;
            res.err = err
            console.error(res.statusCode, {reqPathName, body, err});
        }
        res.end(JSON.stringify(response));
    });
    return true;
}

const server = http.createServer((req, res) => {
    const reqUrl = url.parse(req.url, true);
    const reqPathname = reqUrl.pathname;

    for (const serviceConfig of serviceConfigs) {
        // TODO: continue this
    }

    if (reqPathname != `/${pathname}` && req.method === 'POST') {
        return false;
    }

    // getCities()
    if(respond('getCities', require(`${__dirname}/lib.js`).getCities, [], req, res)) {
        return true;
    }

    res.statusCode = 404;
    res.end();

});

server.listen(port, hostname, () => {
    console.log(`Server running at http://${hostname}:${port}/`);
});
