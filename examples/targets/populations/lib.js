const fs = require('fs');

/**
 * Get statistics
 * @param {string[]} cityNames - Names of cities
 * @returns {population: number}
 */
function getPopulations(cityNames) {
    const content = fs.readFileSync(`${__dirname}/data.json`);
    const data = JSON.parse(content);
    const result = cityNames.map((name) => {
        return {[name]: data[name]};
    });
    return result;
}

module.exports = {
    getPopulations,
};
