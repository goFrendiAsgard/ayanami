/**
 * Get available cities
 * @returns {string[]}
 */
function getCities() {
    const data = [
        "king's landing",
        "braavos",
        "volantis",
        "asshai",
        "old valyria",
        "free cities",
        "qarth",
        "meereen"
    ];
    return data;
}

module.exports = {
    getCities,
};
