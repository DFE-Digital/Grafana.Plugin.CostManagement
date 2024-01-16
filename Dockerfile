ARG grafana_version=latest
ARG grafana_image=grafana

FROM grafana/${grafana_image}:${grafana_version}

# Inject livereload script into grafana index.html
USER root

# Copy your plugin files into the Grafana plugins directory
COPY dfe-azurecostbackend-datasource-v2 /var/lib/grafana/plugins/dfe-azurecostbackend-datasource-v2

# Set permissions for the plugin directory
RUN chown -R 472:472 /var/lib/grafana/plugins/dfe-azurecostbackend-datasource-v2