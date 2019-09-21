const path = require('path');
const webpack = require('webpack');
require('@babel/register');

const config = {
  entry: {
    game: './frontend/scripts/game.js',
    lobby: './frontend/scripts/lobby.js',
  },
  output: {
    path: `${__dirname}/dist`,
    filename: '[name].js',
  },
  module: {
    rules: [
      {
        test: /\.js$/,
        exclude: /node_modules/,
        use: ['babel-loader'],
      },
      {
        test: /\.css$/,
        use: ['style-loader', 'css-loader'],
      },
    ],
  },
  plugins: [
    new webpack.HotModuleReplacementPlugin(),
  ],
  resolve: {
    modules: [
      path.resolve('./src'),
      path.resolve('./node_modules'),
    ],
  },
  devServer: {
    contentBase: `${__dirname}/dist`,
    compress: true,
    hot: true,
  },
  watch: false,
  devtool: 'source-map',
};

module.exports = config;
