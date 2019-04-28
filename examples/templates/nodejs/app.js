const process = require('process');
const appHelper = require('./app-helper.js');

const servicename = 'analytics';
const hostname =  process.env[`${servicename}_hostname`.toUpperCase()] || '127.0.0.1';
const port = process.env[`${servicename}_port`.toUpperCase()] || 3000;

const serviceConfigs = [
    {
        pathname: 'getPopulationAnalysis',
        fn: require(`${__dirname}/lib.js`).getPopulationAnalysis,
        inputs: ['populations'],
    },
    {
        pathname: 'getWeatherAnalysis',
        fn: require(`${__dirname}/lib.js`).getWeatherAnalysis,
        inputs: ['weathers'],
    },
];

const server = appHelper.createServer(servicename, serviceConfigs)

server.listen(port, hostname, () => {
    console.log(`${servicename} running at http://${hostname}:${port}/`);
});
