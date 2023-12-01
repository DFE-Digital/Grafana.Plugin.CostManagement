import React from 'react';
import { InlineField, Input } from '@grafana/ui';
export function QueryEditor({ query, onChange, onRunQuery }) {
    const onQueryTextChange = (event) => {
        onChange(Object.assign(Object.assign({}, query), { queryText: event.target.value }));
    };
    const onConstantChange = (event) => {
        onChange(Object.assign(Object.assign({}, query), { constant: parseFloat(event.target.value) }));
        // executes the query
        onRunQuery();
    };
    const { queryText, constant } = query;
    return (React.createElement("div", { className: "gf-form" },
        React.createElement(InlineField, { label: "Constant" },
            React.createElement(Input, { onChange: onConstantChange, value: constant, width: 8, type: "number", step: "0.1" })),
        React.createElement(InlineField, { label: "Query Text", labelWidth: 16, tooltip: "Not used yet" },
            React.createElement(Input, { onChange: onQueryTextChange, value: queryText || '' }))));
}
//# sourceMappingURL=QueryEditor.js.map
