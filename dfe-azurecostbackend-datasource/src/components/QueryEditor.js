import React, { useState } from 'react';
import { InlineField, Input, Checkbox } from '@grafana/ui';
export function QueryEditor({ query, onChange, onRunQuery }) {
    const [forecast, setForecast] = useState(query.forecast || false);
    const onQueryTextChange = (event) => {
        onChange(Object.assign(Object.assign({}, query), { queryText: event.target.value }));
        // executes the query
        onRunQuery();
    };
    const onConstantChange = (event) => {
        onChange(Object.assign(Object.assign({}, query), { constant: parseFloat(event.target.value) }));
        // executes the query
        onRunQuery();
    };
    const onForecastChange = () => {
        const newForecast = !forecast;
        setForecast(newForecast);
        onChange(Object.assign(Object.assign({}, query), { forecast: newForecast }));
        // executes the query
        onRunQuery();
    };
    const { queryText, constant } = query;
    return (React.createElement("div", { className: "gf-form" },
        false && (React.createElement(InlineField, { label: "Constant" },
            React.createElement(Input, { onChange: onConstantChange, value: constant, width: 8, type: "number", step: "0.1" }))),
        React.createElement(InlineField, { label: "Azure Resource Id", labelWidth: 26, tooltip: "Resource Id" },
            React.createElement(Input, { onChange: onQueryTextChange, value: queryText || '' })),
        React.createElement(InlineField, { label: "Forecast" },
            React.createElement(Checkbox, { value: forecast, onChange: onForecastChange }))));
}
//# sourceMappingURL=QueryEditor.js.map
