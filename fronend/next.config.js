/* eslint-disable @typescript-eslint/no-var-requires */
const path = require('path')

/** @type {import('next').NextConfig} */
module.exports = {
  trailingSlash: true,
  reactStrictMode: false,
  webpack: (config, { buildId, dev, isServer, defaultLoaders, webpack }) => {
    config.resolve.alias = {
      ...config.resolve.alias,
      apexcharts: path.resolve(__dirname, './node_modules/apexcharts-clevision')
    }

    // Modify webpack settings to ignore certain warnings
    config.stats = config.stats || {}
    config.stats.warningsFilter = warnings => {
      const filteredWarnings = warnings.filter(
        // Use regex or string to identify warnings to ignore
        warning => !/no-unused-vars/.test(warning) && !/lines-around-comment/.test(warning)
      )
      return filteredWarnings
    }

    return config
  },

  // NOTE: env variables goes here
  env: {
    REMOTE_CONTROL_SERVER_URL: process.env.REMOTE_CONTROL_SERVER_URL,
    CLIENT_NAME: process.env.CLIENT_NAME || 'test',
    TIME_ZONE: process.env.TIME_ZONE || 'UTC',
    OFFSET_HOURS: process.env.OFFSET_HOURS || '-6'
  }
}