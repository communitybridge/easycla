'use strict'

/**
 * New Relic agent configuration.
 *
 * See lib/config.defaults.js in the agent distribution for a more complete
 * description of configuration variables and their potential values.
 */

var newRelicLicenseKey = process.env['NEWRELIC_LICENSE'];
var newRelicAppName = process.env['NEWRELIC_APP_NAME'];
var newRelicLabels = process.env['NEWRELIC_APP_LABELS'];

exports.config = {
  /**
   * Array of application names.
   */
  app_name: [newRelicAppName],
  /**
   * Your New Relic license key.
   */
  license_key: newRelicLicenseKey,
  /**
   * Your New Relic labels.
   * Specify your labels as objects or a semicolon-delimited string of
   * colon-separated pairs (for example, Server:One;Data Center:Primary).
   */
  labels: newRelicLabels,
  logging: {
    /**
     * Level at which to log. 'trace' is most useful to New Relic when diagnosing
     * issues with the agent, 'info' and higher will impose the least overhead on
     * production applications.
     */
    level: 'info'
  }
}
