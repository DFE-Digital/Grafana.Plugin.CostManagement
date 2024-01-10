import React from 'react';
import { InlineField, Input } from '@grafana/ui';
export function QueryEditor({ query, onChange, onRunQuery }) {
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
    const { queryText, constant } = query;
    return (React.createElement("div", { className: "gf-form" },
        false && (React.createElement(InlineField, { label: "Constant" },
            React.createElement(Input, { onChange: onConstantChange, value: constant, width: 8, type: "number", step: "0.1" }))),
        React.createElement(InlineField, { label: "Azure Reource Id", labelWidth: 26, tooltip: "Reource Id" },
            React.createElement(Input, { onChange: onQueryTextChange, value: queryText || '' }))));
}
//# sourceMappingURL=QueryEditor.js.map