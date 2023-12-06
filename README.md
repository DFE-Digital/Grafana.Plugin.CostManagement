# Grafana.Plugin.CostManagement
This is a Grafana datasource plugin to display Azure cost data

## Requirements

npm, mage, go, and any supported IDE for DEV running (suggested would be VS Code).

## Local Running

This requires a linux environment with npm, mage and go installed

###Build Steps

Clone the repos then cd  dfe-azurecostbackend-datasource
npm install 
npm run dev
mage -v build:linux

###Sign the Plugin

export GRAFANA_ACCESS_POLICY_TOKEN=your token created in the grafana cloud site
npx @grafana/sign-plugin@latest --rootUrls http://localhost:3000/

###Run Steps

Call: docker-compose up
Select the Azure Cost Datasource, confugure it then add it to a pannel