import React from 'react';

import { formatMetric } from '../../utils/string-utils';

function NodeDetailsTableNodeMetric(props) {
  return (
    <td className="node-details-table-node-metric">
      {formatMetric(props.value, props)}
    </td>
  );
}

export default NodeDetailsTableNodeMetric;
