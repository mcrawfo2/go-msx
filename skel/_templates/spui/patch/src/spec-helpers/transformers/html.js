/* eslint-env node */
module.exports = {
  process(content) {
    return 'module.exports = ' + JSON.stringify(content) + ';';
  }
};
