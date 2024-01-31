import React, { ChangeEvent, useState } from 'react';
import { InlineField, Input, Checkbox } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from '../datasource';
import { MyDataSourceOptions, MyQuery } from '../types';

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;

export function QueryEditor({ query, onChange, onRunQuery }: Props) {
  const [forecast, setForecast] = useState<boolean>(query.forecast || false);

  const onQueryTextChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, queryText: event.target.value });
    // executes the query
    onRunQuery();
  };

  const onConstantChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, constant: parseFloat(event.target.value) });
    // executes the query
    onRunQuery();
  };

  const onForecastChange = () => {
    const newForecast = !forecast;
    setForecast(newForecast);
    onChange({ ...query, forecast: newForecast });
    // executes the query
    onRunQuery();
  };

  const { queryText, constant } = query;

  return (
    <div className="gf-form">
      {false && (
        <InlineField label="Constant">
          <Input onChange={onConstantChange} value={constant} width={8} type="number" step="0.1" />
        </InlineField>
      )}
      
      <InlineField label="Azure Resource Id" labelWidth={26} tooltip="Resource Id">
        <Input onChange={onQueryTextChange} value={queryText || ''} />
      </InlineField>

      <InlineField label="Forecast">
        <Checkbox value={forecast} onChange={onForecastChange} />
      </InlineField>
    </div>
  );
}