const fs = require('fs');

/**
 * Get available cities
 * @returns {string[]}
 */
function getCities() {
    const content = fs.readFileSync('./data.json');
    const data = JSON.parse(content);
    return data;
}

module.exports = {
    getCities,
};

if (require.main === module) {
    console.log(getCities());
}
