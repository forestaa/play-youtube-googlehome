const path = require('path')

module.exports = {
  mode: 'development',
  entry: [
    './src/typescript/index.tsx'
  ],
  output: {
    path: path.resolve(__dirname + 'dist'),
    filename: 'bundle.js'
  },
  devServer: {
    contentBase: './dist'
  },
  module: {
    rules: [
        {
            test: /\.tsx?$/,
            exclude: /(node_modules)/,
            use: [
              {
                loader: 'ts-loader',
                options: {
                  configFile: "tsconfig.json"
                },
              },
              {
                loader: 'tslint-loader',
                options: {
                  emitErrors: true
                },
              }

            ]
        },
    ]
  },
  resolve: {
    extensions: ['.tsx', '.ts', '.js']
  }
};