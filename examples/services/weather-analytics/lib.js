/**
 * Get weathers
 * @param {Object[]} weathers - The weathers
 * @param {string} weathers[].main - condition.
 * @param {number} weathers[].temp - temperature.
 * @returns {main: string[], temp: number}
 */
function getAnalysis(weathers) {
    const weatherList = Object.keys(weathers).map((cityName) => weathers[cityName]);
    const cityCount = weatherList.length;
    // get mainMode
    const mainList = weatherList.map((weather) => weather.main);
    const mains = mainList.reduce((acc, main) => {
        if (!(main in acc)) {
            acc[main] = 0;
        }
        acc[main] ++;
        return acc;
    }, {});
    const maxMainCount = Object.keys(mains).reduce((max, mainName) => {
        if (max < mains[mainName]) {
            max = mains[mainName];
        }
        return max;
    }, 0);
    const mainMode = Object.keys(mains).filter((mainName) => mains[mainName] === maxMainCount);
    // get average temperature
    const tempList = weatherList.map((weather) => weather.temp);
    const tempAvg = tempList.reduce((acc, temp) => acc + temp, 0) / cityCount;
    return {main: mainMode, temp: tempAvg};
}

module.exports = {
    getAnalysis,
}

if (require.main === module) {
    console.log(getAnalysis({
        "king's landing": {
            "main": "clear",
            "temp": 281.34
        },
        "braavos": {
            "main": "clouds",
            "temp": 269.35
        },
        "volantis": {
            "main": "clouds",
            "temp": 294.21
        },
        "asshai": {
            "main": "rain",
            "temp": 235.43
        },
        "old valyria": {
            "main": "rain",
            "temp": 252.46
        },
        "free cities": {
            "main": "rain",
            "temp": 281.25
        },
        "qarth": {
            "main": "clouds",
            "temp": 243.32
        },
        "meereen": {
            "main": "clear",
            "temp": 284.32
        }
    }));
}
