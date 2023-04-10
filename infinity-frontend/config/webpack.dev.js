const merge = require('webpack-merge');

const common = require('./webpack.common.js');

const isHot = process.argv.indexOf('--hot') !== -1;

const config = merge(common, {
  devServer: {
    historyApiFallback: true,
    compress: true,
    disableHostCheck: true,
    host: '0.0.0.0',
  },
});

// Add HMR as entry if enabled
if (isHot) {
  config.entry.main.unshift('react-hot-loader/patch');
}

module.exports = config;
