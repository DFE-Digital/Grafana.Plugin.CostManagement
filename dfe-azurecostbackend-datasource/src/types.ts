import { DataQuery, DataSourceJsonData } from '@grafana/schema';

export interface MyQuery extends DataQuery {
  queryText?: string;
  constant: number;
  forecast: boolean;
}

export const DEFAULT_QUERY: Partial<MyQuery> = {
  constant: 6.5,
};

/**
 * These are options configured for each DataSource instance
 */
export interface MyDataSourceOptions extends DataSourceJsonData {
  path?: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface MySecureJsonData {
  //apiKey?: string;
  //AzureCostSubscriptionUrl?: string;
	Password?:                 string;
	ClientID?:                 string;
	TenantID?:                 string;
	SubscriptionID?:           string;
	Region?:                   string;
}

