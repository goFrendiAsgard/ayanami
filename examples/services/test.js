const assert = require('assert').strict;
const cities = require('./cities/lib.js');
const populations = require('./populations/lib.js');
const weathers = require('./weathers/lib.js');
const analytics = require('./analytics/lib.js');

function testGetCities() {
    const expected = [
        "king's landing",
        "braavos",
        "volantis",
        "asshai",
        "old valyria",
        "free cities",
        "qarth",
        "meereen"
    ];
    const actual = cities.getCities();
    assert.deepStrictEqual(actual, expected);
}

function testGetPopulations() {
    const expected = [{"volantis":{"population":3000}},{"asshai":{"population":70000}}];
    const actual = populations.getPopulations(["volantis", "asshai"]);
    assert.deepStrictEqual(actual, expected);
}

function testGetWeathers() {
    const expected = [{"volantis":{"main":"clouds","temp":294.21}},{"asshai":{"main":"rain","temp":235.43}}];
    const actual = weathers.getWeathers(["volantis", "asshai"]);
    assert.deepStrictEqual(actual, expected);
}

function testGetPopulationAnalysis() {
    const expected = {"populationMean":31625};
    const actual = analytics.getPopulationAnalysis({
        "king's landing": {
            "population": 60000
        },
        "braavos": {
            "population": 4000
        },
        "volantis": {
            "population": 3000
        },
        "asshai": {
            "population": 70000
        },
        "old valyria": {
            "population": 4000
        },
        "free cities": {
            "population": 62000
        },
        "qarth": {
            "population": 40000
        },
        "meereen": {
            "population": 10000
        }
    });
    assert.deepStrictEqual(actual, expected);
}

function testGetWeatherAnalysis() {
    const expected = {"mainMode":["clouds","rain"],"tempMean":267.71000000000004};
    const actual = analytics.getWeatherAnalysis({
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
    });
    assert.deepStrictEqual(actual, expected);
}

if (require.main === module) {
    try {
        console.log("TEST STARTED");
        testGetCities();
        testGetPopulations();
        testGetWeathers();
        testGetWeatherAnalysis();
        testGetPopulationAnalysis();
        console.log("TEST DONE");
    } catch(err) {
        console.error(err);
    }
}
