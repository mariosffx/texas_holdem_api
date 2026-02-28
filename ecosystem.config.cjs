/* global process */
require("dotenv").config();

module.exports = {
  apps: [
    {
      name: "texas_holdem_api",
      script: "npm run start"
    }
  ],
  deploy: {
    production: {
      user: "root",
      host: "dns-pi",
      ref: "origin/develop",
      repo: "git@github.com:mariosffx/texas_holdem_api.git",
      path: "/root/projects/texas_holdem_api",
      "post-deploy": "./scripts/deploy/post-deploy.sh"
    }
  }
};
