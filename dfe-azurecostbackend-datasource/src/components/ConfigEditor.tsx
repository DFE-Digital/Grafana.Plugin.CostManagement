import React, { ChangeEvent } from 'react';
import { InlineField, SecretInput } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { MyDataSourceOptions, MySecureJsonData } from '../types';

interface Props extends DataSourcePluginOptionsEditorProps<MyDataSourceOptions> {}

export function ConfigEditor(props: Props) {
  const { onOptionsChange, options } = props;

  const onPasswordKeyChange = (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (event.key === 'Enter' || event.type === 'blur') {
      // Save the value on Enter key press
      onOptionsChange({
        ...options,
        secureJsonData: {
          Password: (event.target as HTMLInputElement).value,
        },
      });
    }
  };

  const onResetPasswordKey = () => {
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        Password: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        Password: '',
      },
    });
  };

  const onClientIDKeyChange = (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (event.key === 'Enter' || event.type === 'blur') {
      // Save the value on Enter key press
      onOptionsChange({
        ...options,
        secureJsonData: {
          ClientID: (event.target as HTMLInputElement).value,
        },
      });
    }
  };

  const onResetClientIDKey = () => {
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        ClientID: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        ClientID: '',
      },
    });
  };

  const onTenantIDKeyChange = (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (event.key === 'Enter' || event.type === 'blur') {
      // Save the value on Enter key press
      onOptionsChange({
        ...options,
        secureJsonData: {
          TenantID: (event.target as HTMLInputElement).value,
        },
      });
    }
  };

  const onResetTenantIDKey = () => {
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        TenantID: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        TenantID: '',
      },
    });
  };

  const onSubscriptionIDKeyChange = (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (event.key === 'Enter' || event.type === 'blur') {
      // Save the value on Enter key press
      onOptionsChange({
        ...options,
        secureJsonData: {
          SubscriptionID: (event.target as HTMLInputElement).value,
        },
      });
    }
  };

  const onResetSubscriptionIDKey = () => {
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        SubscriptionID: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        SubscriptionID: '',
      },
    });
  };

  const onRegionKeyChange = (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (event.key === 'Enter' || event.type === 'blur') {
      // Save the value on Enter key press
      onOptionsChange({
        ...options,
        secureJsonData: {
          Region: (event.target as HTMLInputElement).value,
        },
      });
    }
  };

  const onResetRegionKey = () => {
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        Region: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        Region: '',
      },
    });
  };

  //const { jsonData, secureJsonFields } = options;
  const { secureJsonFields } = options;
  const secureJsonData = (options.secureJsonData || {}) as MySecureJsonData;

  
  return (
    <div className="gf-form-group">
      <InlineField label="Password / Client Secret" labelWidth={27}>
        <SecretInput
          isConfigured={(secureJsonFields && secureJsonFields.Password) as boolean}
          value={secureJsonData.Password || ''}
          placeholder="secure Password / Client Secret (backend only)"
          width={100}
          onReset={onResetPasswordKey}
          onChange={onPasswordKeyChange}
        />
      </InlineField>
      <InlineField label="ClientID" labelWidth={12}>
        <SecretInput
          isConfigured={(secureJsonFields && secureJsonFields.ClientID) as boolean}
          value={secureJsonData.ClientID || ''}
          placeholder="secure Client ID (backend only)"
          width={100}
          onReset={onResetClientIDKey}
          onChange={onClientIDKeyChange}
        />
      </InlineField>
      <InlineField label="TenantID" labelWidth={12}>
        <SecretInput
          isConfigured={(secureJsonFields && secureJsonFields.TenantID) as boolean}
          value={secureJsonData.TenantID || ''}
          placeholder="secure Tenant ID (backend only)"
          width={60}
          onReset={onResetTenantIDKey}
          onChange={onTenantIDKeyChange}
        />
      </InlineField>
      <InlineField label="SubscriptionID" labelWidth={17}>
        <SecretInput
          isConfigured={(secureJsonFields && secureJsonFields.SubscriptionID) as boolean}
          value={secureJsonData.SubscriptionID || ''}
          placeholder="secure Subscription ID (backend only)"
          width={100}
          onReset={onResetSubscriptionIDKey}
          onChange={onSubscriptionIDKeyChange}
        />
      </InlineField>
      <InlineField label="Region" labelWidth={12}>
        <SecretInput
          isConfigured={(secureJsonFields && secureJsonFields.Region) as boolean}
          value={secureJsonData.Region || ''}
          placeholder="secure Region (backend only)"
          width={100}
          onReset={onResetRegionKey}
          onChange={onRegionKeyChange}
        />
      </InlineField>
    </div>
  );
}

