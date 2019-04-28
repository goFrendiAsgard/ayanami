const http = require('http');
const url = require('url');

function respond(servicename, reqPathname, fn, inputs, req, res) {
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
        let fnInputs, fnOutput, statusCode;
        try {
            const postBody = JSON.parse(body);
            fnInputs = inputs.map((key) => postBody[key]);
            fnOutput = fn(...fnInputs);
            response.result = fnOutput;
            statusCode = 200;
            res.statusCode = statusCode;
            console.log(JSON.stringify({statusCode, servicename, reqPathname, body, fnInputs, fnOutput}, null, 2));
        } catch (err) {
            statusCode = 500;
            response.err = err.message
            res.err = err.message
            res.statusCode = statusCode;
            console.error(JSON.stringify({statusCode, servicename, reqPathname, body, fnInputs, fnOutput}, null, 2));
        }
        res.end(JSON.stringify(response));
    });
    return true;
}

function createServer(servicename, serviceConfigs) {
    const server = http.createServer((req, res) => {
        const reqUrl = url.parse(req.url, true);
        const reqPathname = reqUrl.pathname;
        // respond handler
        for (const serviceConfig of serviceConfigs) {
            const { pathname, fn, inputs } = serviceConfig;
            if (reqPathname === `/${pathname}` && req.method === 'POST') {
                return respond(servicename, reqPathname, fn, inputs, req, res);
            }
        }
        // 404 handler
        const err = new Error('Not found');
        const response = {
            result: {},
            err: err.message,
        }
        let body = '';
        req.on('data', function (chunk) {
            body += chunk;
        });
        req.on('end', function () {
            const statusCode = 404;
            console.error(JSON.stringify({statusCode, servicename, reqPathname, body, err}, null, 2));
            res.setHeader('Content-Type', 'application/json');
            res.err = err.message;
            res.statusCode = statusCode;
            res.end(JSON.stringify(response));
        });
        return false;
    });
    return server;
}

module.exports = {
    createServer,
    respond,
};
