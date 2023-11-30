import { DataSourceWithBackend } from '@grafana/runtime';
import { DEFAULT_QUERY } from './types';
export class DataSource extends DataSourceWithBackend {
    constructor(instanceSettings) {
        super(instanceSettings);
    }
    getDefaultQuery(_) {
        return DEFAULT_QUERY;
    }
}
//# sourceMappingURL=datasource.js.map