const fs = require('fs');

/**
 * Get available cities
 * @returns {string[]}
 */
function getCities() {
    const content = fs.readFileSync(`${__dirname}/data.json`);
    const data = JSON.parse(content);
    return data;
}

module.exports = {
    getCities,
};
