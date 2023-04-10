const path = require('path');
const webpack = require('webpack');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const PreloadWebpackPlugin = require('preload-webpack-plugin');
const ExtractTextWebpackPlugin = require('extract-text-webpack-plugin');
const DotenvWebpackPlugin = require('dotenv-webpack');
const autoprefixer = require('autoprefixer');
const LodashWebpackPlugin = require('lodash-webpack-plugin');

const projectRoot = path.resolve(__dirname, '..'); // Project root

const isHot = process.argv.indexOf('--hot') !== -1;

// Define env vars
const NODE_ENV = process.env.NODE_ENV || 'development';
const __DEV__ = NODE_ENV === 'development';
const __TEST__ = NODE_ENV === 'test';
const __PROD__ = NODE_ENV === 'production';

const StyleExtractPlugin = new ExtractTextWebpackPlugin({
  filename: 'style-[contenthash:20].css',
  allChunks: true,
  disable: isHot, // Disable and fallback to style-loader if HMR is enabled
});

module.exports = {
  resolve: {
    modules: ['src', 'node_modules'],
  },
  entry: {
    main: ['./src/main.js'], // This is an array because we want to be able to add HMR in dev config
  },
  output: {
    path: path.resolve(projectRoot, 'dist'),
    filename: '[name].js',
    sourceMapFilename: '[file].map',
    chunkFilename: '[name].js',
    publicPath: '/',
  },
  module: {
    rules: [
      {
        // CSS/SCSS loader
        test: /\.s?css$/,
        exclude: /node_modules/,
        use: StyleExtractPlugin.extract({
          fallback: 'style-loader',
          // Use ETWP to put CSS in its own file
          use: [
            {
              loader: 'css-loader', // Load css and turn into modules
              options: {
                modules: true,
                localIdentName: '__[local]_[hash:base64:5]',
                sourceMap: true,
              },
            },
            {
              loader: 'postcss-loader', // Used to add autoprefixes
              options: {
                sourceMap: true,
                plugins: () =>
                  autoprefixer({
                    browsers: ['last 5 versions'],
                  }),
              },
            },
            'sass-loader', // Used to parse SCSS
            {
              loader: 'sass-resources-loader', // Used to inject SCSS variables, thus making them "global"
              options: {
                resources: path.resolve(projectRoot, 'src/styles/resources/_*.scss'),
                sourceMap: true,
              },
            },
          ],
        }),
      },
      // Use babel-loader for js
      {
        test: /\.jsx?$/,
        loader: 'babel-loader',
        exclude: /node_modules/,
      },
      {
        test: /\.(jpg|png)$/,
        oneOf: [
          {
            use: {
              loader: 'url-loader',
              options: {
                limit: 5000,
                name: '[name].[ext]',
              },
            },
          },
          {
            include: path.resolve(projectRoot, 'public/assets/images'),
            use: {
              loader: 'file-loader',
              options: {
                context: 'src',
                name: '[name].[ext]',
              },
            },
          },
        ],
      },
      {
        test: /\.(ttf|woff|woff2|eot)$/,
        oneOf: [
          {
            use: {
              loader: 'url-loader',
              options: {
                limit: 5000,
                name: '[hash].[ext]',
              },
            },
          },
          {
            include: path.resolve(projectRoot, 'public/assets/fonts'),
            use: {
              loader: 'file-loader',
              options: {
                name: '[name].[ext]',
              },
            },
          },
        ],
      },
      {
        test: /\.(mp4|ogg|webm)$/,
        oneOf: [
          {
            use: {
              loader: 'url-loader',
              options: {
                limit: 5000,
                name: '[hash].[ext]',
              },
            },
          },
          {
            include: path.resolve(projectRoot, 'public/assets/videos'),
            use: {
              loader: 'file-loader',
              options: {
                name: '[name].[ext]',
              },
            },
          },
        ],
      },
      {
        test: /\.(mp3)$/,
        oneOf: [
          {
            use: {
              loader: 'url-loader',
              options: {
                limit: 5000,
                name: '[hash].[ext]',
              },
            },
          },
          {
            include: path.resolve(projectRoot, 'public/assets/audio'),
            use: {
              loader: 'file-loader',
              options: {
                name: '[name].[ext]',
              },
            },
          },
        ],
      },
      {
        test: /\.(svg)$/,
        oneOf: [
          {
            exclude: path.resolve(projectRoot, 'public/assets/images'),
            use: {
              loader: 'svg-inline-loader',
              options: {
                removeTags: true,
                removingTags: ['title'],
              },
            },
          },
          {
            include: path.resolve(projectRoot, 'public/assets/images'),
            use: {
              loader: 'url-loader',
              options: {
                limit: 5000,
                name: '[name].[ext]',
              },
            },
          },
        ],
      },
    ],
  },
  plugins: [
    new HtmlWebpackPlugin({
      // Create a html with the build files from template
      template: './src/index.html',
    }),
    // Extract loaded CSS and save to a file
    StyleExtractPlugin,
    new DotenvWebpackPlugin(),
    new webpack.IgnorePlugin(/^\.\/locale$/, /moment$/),
    new webpack.HashedModuleIdsPlugin(), // So vendor caching works correctly
    // Extract commonly used components to main bundle
    new webpack.optimize.CommonsChunkPlugin({
      name: 'main',
      children: true,
      minChunks: 2,
    }),
    // Bundle vendor libraries
    new webpack.optimize.CommonsChunkPlugin({
      name: 'vendor',
      minChunks(module) {
        return module.context && module.context.indexOf('node_modules') >= 0;
      },
    }),
    // Extract manifest
    new webpack.optimize.CommonsChunkPlugin({ name: 'manifest', minChunks: Infinity }),
    // Inject env vars
    new webpack.DefinePlugin({
      'process.env.NODE_ENV': JSON.stringify(NODE_ENV),
      __DEV__,
      __TEST__,
      __PROD__,
    }),
    new LodashWebpackPlugin({
      collections: true,
      shorthands: true,
    }),
    new PreloadWebpackPlugin({
      rel: 'preload',
      fileWhitelist: [/.*main.*\.js/, /.*vendor.*\.js/, /.*.css/, /.*manifest.*\.js/],
      include: 'allChunks',
    }),
  ],
  node: {
    fs: 'empty',
  },
};
