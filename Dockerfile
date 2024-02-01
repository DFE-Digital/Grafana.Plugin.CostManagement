ARG grafana_version=latest
ARG grafana_image=grafana

FROM grafana/${grafana_image}:${grafana_version}

# Inject livereload script into grafana index.html
USER root

COPY latest-cost-plugin/grafana.ini /etc/grafana/grafana.ini

# Copy your plugin files into the Grafana plugins directory
COPY latest-cost-plugin/dfe-azurecostbackend-datasource /var/lib/grafana/plugins/dfe-azurecostbackend-datasource
COPY latest-cost-plugin/blackcowmoo-googleanalytics-datasource /var/lib/grafana/plugins/blackcowmoo-googleanalytics-datasource

# Set permissions for the plugin directory
RUN chown -R 472:472 /var/lib/grafana/plugins/dfe-azurecostbackend-datasource
RUN chown -R 472:472 /var/lib/grafana/plugins/blackcowmoo-googleanalytics-datasource