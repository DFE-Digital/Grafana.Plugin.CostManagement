import React from 'react';
import { InlineField, SecretInput } from '@grafana/ui';
export function ConfigEditor(props) {
    const { onOptionsChange, options } = props;
    const onPasswordKeyChange = (event) => {
        onOptionsChange(Object.assign(Object.assign({}, options), { secureJsonData: {
                Password: event.target.value,
            } }));
    };
    const onResetPasswordKey = () => {
        onOptionsChange(Object.assign(Object.assign({}, options), { secureJsonFields: Object.assign(Object.assign({}, options.secureJsonFields), { Password: false }), secureJsonData: Object.assign(Object.assign({}, options.secureJsonData), { Password: '' }) }));
    };
    const onClientIDKeyChange = (event) => {
        onOptionsChange(Object.assign(Object.assign({}, options), { secureJsonData: {
                ClientID: event.target.value,
            } }));
    };
    const onResetClientIDKey = () => {
        onOptionsChange(Object.assign(Object.assign({}, options), { secureJsonFields: Object.assign(Object.assign({}, options.secureJsonFields), { ClientID: false }), secureJsonData: Object.assign(Object.assign({}, options.secureJsonData), { ClientID: '' }) }));
    };
    const onTenantIDKeyChange = (event) => {
        onOptionsChange(Object.assign(Object.assign({}, options), { secureJsonData: {
                TenantID: event.target.value,
            } }));
    };
    const onResetTenantIDKey = () => {
        onOptionsChange(Object.assign(Object.assign({}, options), { secureJsonFields: Object.assign(Object.assign({}, options.secureJsonFields), { TenantID: false }), secureJsonData: Object.assign(Object.assign({}, options.secureJsonData), { TenantID: '' }) }));
    };
    const onSubscriptionIDKeyChange = (event) => {
        onOptionsChange(Object.assign(Object.assign({}, options), { secureJsonData: {
                SubscriptionID: event.target.value,
            } }));
    };
    const onResetSubscriptionIDKey = () => {
        onOptionsChange(Object.assign(Object.assign({}, options), { secureJsonFields: Object.assign(Object.assign({}, options.secureJsonFields), { SubscriptionID: false }), secureJsonData: Object.assign(Object.assign({}, options.secureJsonData), { SubscriptionID: '' }) }));
    };
    const onRegionKeyChange = (event) => {
        onOptionsChange(Object.assign(Object.assign({}, options), { secureJsonData: {
                Region: event.target.value,
            } }));
    };
    const onResetRegionKey = () => {
        onOptionsChange(Object.assign(Object.assign({}, options), { secureJsonFields: Object.assign(Object.assign({}, options.secureJsonFields), { Region: false }), secureJsonData: Object.assign(Object.assign({}, options.secureJsonData), { Region: '' }) }));
    };
    //const { jsonData, secureJsonFields } = options;
    const { secureJsonFields } = options;
    const secureJsonData = (options.secureJsonData || {});
    return (React.createElement("div", { className: "gf-form-group" },
        React.createElement(InlineField, { label: "Password / Client Secret", labelWidth: 27 },
            React.createElement(SecretInput, { isConfigured: (secureJsonFields && secureJsonFields.Password), value: secureJsonData.Password || '', placeholder: "secure Password / Client Secret (backend only)", width: 100, onReset: onResetPasswordKey, onChange: onPasswordKeyChange })),
        React.createElement(InlineField, { label: "ClientID", labelWidth: 12 },
            React.createElement(SecretInput, { isConfigured: (secureJsonFields && secureJsonFields.ClientID), value: secureJsonData.ClientID || '', placeholder: "secure Client ID (backend only)", width: 100, onReset: onResetClientIDKey, onChange: onClientIDKeyChange })),
        React.createElement(InlineField, { label: "TenantID", labelWidth: 12 },
            React.createElement(SecretInput, { isConfigured: (secureJsonFields && secureJsonFields.TenantID), value: secureJsonData.TenantID || '', placeholder: "secure Tenant ID (backend only)", width: 60, onReset: onResetTenantIDKey, onChange: onTenantIDKeyChange })),
        React.createElement(InlineField, { label: "SubscriptionID", labelWidth: 17 },
            React.createElement(SecretInput, { isConfigured: (secureJsonFields && secureJsonFields.SubscriptionID), value: secureJsonData.SubscriptionID || '', placeholder: "secure Subscription ID (backend only)", width: 100, onReset: onResetSubscriptionIDKey, onChange: onSubscriptionIDKeyChange })),
        React.createElement(InlineField, { label: "Region", labelWidth: 12 },
            React.createElement(SecretInput, { isConfigured: (secureJsonFields && secureJsonFields.Region), value: secureJsonData.Region || '', placeholder: "secure Region (backend only)", width: 100, onReset: onResetRegionKey, onChange: onRegionKeyChange }))));
}
//# sourceMappingURL=ConfigEditor.js.map
