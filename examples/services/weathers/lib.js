const fs = require('fs');

/**
 * Get weathers
 * @param {string[]} cityNames - Names of cities
 * @returns {main: string, temp: number}
 */
function getWeathers(cityNames) {
    const content = fs.readFileSync('./data.json');
    const data = JSON.parse(content);
    const result = cityNames.map((name) => {
        return {[name]: data[name]};
    });
    return result;
}

module.exports = {
    getWeathers,
};

if (require.main === module) {
    console.log(getWeathers(["volantis", "asshai"]));
}
