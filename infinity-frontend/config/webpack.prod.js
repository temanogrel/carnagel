const path = require('path');
const merge = require('webpack-merge');
const UglifyJsWebpackPlugin = require('uglifyjs-webpack-plugin');
const SWPrecacheWebpackPlugin = require('sw-precache-webpack-plugin');
const FileManagerPlugin = require('filemanager-webpack-plugin');

const { BundleAnalyzerPlugin } = require('webpack-bundle-analyzer');
const { BannerPlugin } = require('webpack');

const common = require('./webpack.common.js');

const projectRoot = path.resolve(__dirname, '..');

module.exports = merge(common, {
  devtool: 'source-map', // Add source maps
  output: {
    // Add chunkhash to file names in prod
    filename: '[name]-[chunkhash].js',
    chunkFilename: '[name]-[chunkhash].js',
    // ... also, put build files in static/build
    path: path.resolve(projectRoot, 'dist', 'static', 'build'),
    publicPath: '/static/build/',
  },
  plugins: [
    // UglifyJS
    new UglifyJsWebpackPlugin({ sourceMap: true, uglifyOptions: { output: { comments: false } } }),
    // Add header text to all js files
    new BannerPlugin(`© 2018 - ${new Date().getFullYear()} Camtube, ALL RIGHTS RESERVED ([hash])`),
    // Clean /dist onStart and move index.html to root onEnd
    new FileManagerPlugin({
      onStart: {
        delete: ['./dist'],
      },
      onEnd: {
        move: [{ source: './dist/static/build/index.html', destination: './dist/index.html' }],
        copy: [
          { source: './public/assets', destination: './dist/public/assets' },
          { source: './public/files', destination: './dist' },
        ],
      },
    }),
    // Generate service worker
    new SWPrecacheWebpackPlugin({
      cacheId: 'camtube-sw',
      filepath: path.resolve(projectRoot, 'dist', 'sw.js'),
      minify: true,
      // Remove local path prefixes
      stripPrefix: `${process.cwd()}/dist`,
      // Ignore index.html
      staticFileGlobsIgnorePatterns: [/build\/index\.html$/, /\.map$/],
      // Do not cache bust files with chunkhash
      dontCacheBustUrlsMatching: /-\w{20}/,

      // todo: review this later because this is insanely large
      maximumFileSizeToCacheInBytes: 10 * 1024 * 1024,
    }),
    // Create a bundle analyzer report
    new BundleAnalyzerPlugin({
      analyzerMode: 'static',
      reportFilename: '../../../report.html',
      openAnalyzer: false,
    }),
  ],
  performance: {
    hints: 'warning',
    maxAssetSize: 1000000,
    maxEntrypointSize: 1500000,
  },
});
