const process = require('process');
const appHelper = require('./app-helper.js');

const servicename = 'cities';
const hostname =  process.env[`${servicename}_hostname`.toUpperCase()] || '127.0.0.1';
const port = process.env[`${servicename}_port`.toUpperCase()] || 3000;

const serviceConfigs = [
    {
        pathname: 'getCities',
        fn: require(`${__dirname}/lib.js`).getCities,
        inputs: [],
    },
];

const server = appHelper.createServer(servicename, serviceConfigs)

server.listen(port, hostname, () => {
    console.log(`${servicename} running at http://${hostname}:${port}/`);
});
