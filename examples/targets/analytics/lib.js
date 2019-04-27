function getMean(arrayOfObj, key) {
    const data = arrayOfObj.map((obj) => obj[key]);
    const mean = data.reduce((total, element) => total + element, 0) / data.length;
    return mean;
}

function getMode(arrayOfObj, key) {
    const tabulation = arrayOfObj.reduce((tabulation, obj) => {
        const label = obj[key];
        if (!(label in tabulation)) {
            tabulation[label] = 0;
        }
        tabulation[label] ++;
        return tabulation;
    }, {});
    const labelList = Object.keys(tabulation);
    const countList = labelList.map((label) => tabulation[label]);
    const maxCount = Math.max(...countList);
    const modes = labelList.filter((label) => tabulation[label] === maxCount);
    return modes;
}

function getListFromMap(obj) {
    const keys = Object.keys(obj);
    const list = keys.map((key) => obj[key]);
    return list;
}

/**
 * Get population analysis
 * @param {Object[]} populations - The weathers
 * @param {string} weathers[].populations - population.
 * @returns {populationMean: number}
 */
function getPopulationAnalysis(populations) {
    const populationList = getListFromMap(populations);
    const populationMean = getMean(populationList, 'population');
    return {populationMean};
}

/**
 * Get weather analysis
 * @param {Object[]} weathers - The weathers
 * @param {string} weathers[].main - condition.
 * @param {number} weathers[].temp - temperature.
 * @returns {mainMode: string[], tempMean: number}
 */
function getWeatherAnalysis(weathers) {
    const weatherList = getListFromMap(weathers);
    const mainMode = getMode(weatherList, 'main');
    const tempMean = getMean(weatherList, 'temp');
    return {mainMode, tempMean};
}

module.exports = {
    getWeatherAnalysis,
    getPopulationAnalysis,
}
